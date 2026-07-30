---
id: DES-032
title: 会计内核与总账
date: 2026-07-30
status: active
author: Codex + User
tags: [finance, accounting-kernel, general-ledger, journal]
---

# 会计内核与总账

## 1. 目的与边界

会计内核把已提交业务事实转换为可审计会计事件和不可变凭证。它不创建采购、收货、生产、
发运或银行事实，也不允许 HTTP handler、Process、Agent 或事件消费者绕过 Capability
直接写账务表。

## 2. 核心语义与 Entity

核心语义：

`economic_event`、`accounting_event`、`accounting_policy`、`posting_rule`、
`journal_entry`、`journal_line`、`ledger_balance`、`account`、`debit`、`credit`、
`monetary_amount`、`exchange_rate`、`accounting_period`、`reversal`。

正式 Entity：

- `accounting_event`
- `accounting_book`
- `chart_of_accounts`
- `gl_account_definition`
- `gl_account_legal_entity`
- `posting_rule`
- `journal_entry`
- `journal_line`
- `journal_batch`
- `ledger_balance`
- `currency_definition`
- `exchange_rate`
- `account_mapping`
- `accounting_exception`

组织、账簿、科目表和财政日历归属由
[DES-031](DES-031-multi-organization-and-shared-master-data-foundation.md) 定义。

## 3. 凭证主子表

`journal_entry` 使用 `document_with_lines`，`journal_line` 使用 `document_line`。
明细通过系统管理的 `parent_document_id` 指向 Entity 投影凭证头。凭证至少两行；
每行只能有借方或贷方；头部借贷合计由明细计算，不允许独立编辑。

凭证头必须保存：

- 法人、账簿、会计期间；
- 来源类型/引用、会计事件、规则版本；
- 业务日期、会计日期、过账日期；
- 制单、审核、过账主体和职责岗位；
- 交易币/本位币、借贷合计、状态；
- correlation、idempotency、Capability Execution 和 evidence。

已过账凭证不可更新或删除。错误通过反向凭证、冲销或批准的调整凭证修正。

## 4. 多币种

凭证明细保存交易币金额、本位币金额、汇率类型、汇率日期、汇率记录和不可变汇率快照。
金额使用整数最小货币单位，汇率使用受控高精度 decimal。修改后续汇率不能改变历史凭证。

M9 开业只接受 CNY，并登记 CNY→CNY 恒等汇率；外币资本到账需要独立流程与政策。

## 5. 记账管线

```text
业务事实 committed
→ accounting.event.capture
→ accounting.policy.resolve
→ posting.rule.compile
→ journal.draft.generate
→ policy / approval / period checks
→ journal.post
→ ledger.balance.update
→ reconciliation / report events
```

不变量：

- 同一业务事实与规则版本只生成一次会计事件；
- 同一会计事件/账簿只产生一次有效过账结果；
- 规则、科目、期间、法人、币种或维度缺失时失败关闭；
- 凭证、余额、Entity 投影和 Outbox 在同一事务提交；
- 规则升级只影响生效日后的事件。

## 6. M9 开业会计

当前 M9 能力：

1. `finance.organization.configure`
2. `accounting.book.activate`
3. `chart.of.accounts.activate`
4. `capital.contribution.post`
5. `finance.opening.readiness.evaluate`

实缴资本示例：

```text
借：1002 银行存款
贷：4001 实收资本
```

`finance_opening_ready` 要求财务组织、有效账簿、开放期间、必需科目、已过账资本凭证和
借贷平衡同时成立。`accounting.book.activate` 必须采集用户账套名称、当前会计年度和
12 个连续期间；不能用无参数默认值代替用户确认。

## 7. 迁移与验收

- 修复当前 `accounting_book` 错用 `account` 原型的问题；
- 保留稳定 BOOK 编码、UUID、凭证外键和审计证据；
- 迁移前后核对账簿、科目、期间、凭证、余额和报表；
- 重复执行 no-op，失败无部分凭证；
- 每个余额可下钻凭证、会计事件、业务事实和 World evidence；
- RLS、Capability write boundary 和已过账不可变触发器必须同时通过。
