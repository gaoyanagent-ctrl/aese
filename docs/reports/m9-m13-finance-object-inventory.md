---
id: REPORT-M9-FIN-001
title: M9–M13 财务对象与历史数据盘点
date: 2026-07-29
status: completed
author: Codex + User
tags: [m9, finance, inventory, migration]
---

# M9–M13 财务对象与历史数据盘点

## 结论

M9 开业会计使用 IAOS 通用 `entity_projection_*` 存储，不再建立或继续扩散
`m9_*` 里程碑前缀表。M10–M13 的 AP、资产、成本、AR 和利润对象尚未进入本次
M9 实现，状态保持 `planned`，不能因已有演示投影而宣称子账完成。

机器可读对象清单与治理矩阵位于
`scenario-packs/hctm/finance-governance-baseline.json`，由
`internal/financebaseline` 离线校验。

## 对象边界

| 财务对象 | 里程碑 | IAOS 权威对象 | 当前状态 |
|---|---|---|---|
| 银行/现金账户 | M9 | `bank_account` | 已迁移 |
| 资本承诺 | M9 | `capital_commitment` | 已迁移 |
| 初始预算 | M9 | `budget_envelope` | 已迁移 |
| 会计凭证 | M9 | `journal_entry` + `journal_line` | 已迁移 |
| 应付与工程承诺 | M10–M11 | `accounts_payable_open_item` | 规划 |
| 固定资产 | M10–M11 | `fixed_asset` | 规划 |
| 制造成本 | M12–M13 | `cost_object_actual` | 规划 |
| 应收 | M13 | `accounts_receivable_open_item` | 规划 |
| 管理利润 | M13 | `management_profit_measure` | 规划 |

## 2026-07-29 数据证据

对真实租户 `tenant-gx-f4b3ce3ce8e2712d` 以 RLS tenant context 查询：

| Entity | 记录数 |
|---|---:|
| `bank_account` | 1 |
| `capital_commitment` | 1 |
| `budget_envelope` | 1 |
| `finance_organization` | 1 |
| `accounting_book` | 1 |
| `gl_account` | 2 |
| `accounting_period` | 1 |
| `currency_definition` | 1 |
| `exchange_rate` | 1 |
| `journal_entry` | 1 |
| `journal_line` | 2 |

`journal_entry:journal_line = 1:2` 证明开业凭证按主从单据保存两条借贷分录。
数据库 `public` schema 中 `m9_%` 表数量为 0；上述 Entity 的有效 metadata 均指向
`entity_projection_*`。同一 Entity 出现多个 metadata version 属于版本历史，不代表
重复物理表或重复业务记录。

新建但只完成设立案第 1 节点的
`tenant-gx-54048d38a61540739b27` 只有币种和汇率基线各 1 条，其余开业财务记录为 0；
这验证了数据不会在对应 M9 财务节点执行前被静默伪造。

## 重跑查询

使用 `iaos_app` 并先在事务内设置：

```sql
SELECT set_config('app.current_tenant_id', '<tenant-id>', true);
```

随后分别查询 `entity_projection_bank_account`、
`entity_projection_capital_commitment`、`entity_projection_budget_envelope`、
`entity_projection_finance_organization`、`entity_projection_accounting_book`、
`entity_projection_gl_account`、`entity_projection_accounting_period`、
`entity_projection_currency_definition`、`entity_projection_exchange_rate`、
`entity_projection_journal_entry` 和 `entity_projection_journal_line`。

禁止把本报告中的计数当作固定 seed 断言；验收应校验业务引用、借贷平衡、租户隔离和
节点状态，而不是依赖总行数。
