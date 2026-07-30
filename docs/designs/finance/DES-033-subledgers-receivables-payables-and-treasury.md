---
id: DES-033
title: 子账、应收应付与资金
date: 2026-07-30
status: active
author: Codex + User
tags: [finance, subledger, accounts-receivable, accounts-payable, treasury]
---

# 子账、应收应付与资金

## 1. 目的

建立采购到付款、订单到收款和银行资金的业务子账。子账保留业务明细和未清项，总账保留
控制科目余额；两者必须可对账，但不能互相替代。

## 2. Entity

应付：

- `supplier_invoice`
- `payable_open_item`
- `invoice_match_result`
- `payment_proposal`
- `payment_instruction`
- `payment_execution`
- `supplier_advance`

应收：

- `customer_invoice`
- `receivable_open_item`
- `collection`
- `cash_application`
- `credit_limit`
- `dunning_case`
- `bad_debt_assessment`

资金：

- `bank_account`
- `bank_statement`
- `bank_statement_line`
- `cash_transaction`
- `bank_reconciliation`
- `cash_position`
- `cash_forecast`
- `liquidity_plan`

客户、供应商和银行敏感数据的共享/组织扩展遵循 DES-031。

## 3. AP 与付款

供应商发票必须关联法人、供应商组织扩展、采购/收货事实、币种、税额、付款条件和凭证。
三单匹配比较订单、收货和发票的数量、价格、税和容差。差异进入异常/审批，不自动配平。

付款流程：

```text
到期未清项 → payment.proposal.generate → 审批/双人释放
→ 银行指令 → World 银行结果 → 清账 → 银行对账
```

维护供应商银行账户、创建付款义务、批准付款和释放银行指令必须职责分离。

## 4. AR 与收款

客户接受/开票规则形成应收未清项。支持部分收款、预收、贷项、争议、坏账、核销和多币种。
银行到账只有经可信 Observation 或受治理导入后才能形成现金事实；收款核销不能伪造到账。

信用额度由法人/信用控制范围维护，销售订单占用、发货、开票和回款按规则释放额度。

## 5. Open Item 与对账

每个未清项保存原始金额、剩余金额、到期日、账龄、结算引用和状态。部分结算生成分配记录，
不覆盖原始金额。AP/AR 控制科目余额必须与未清项合计一致。

银行对账以银行流水为外部事实，匹配账面现金交易。未达项、手续费、利息和重复流水进入
明确异常流程。

## 6. Capability / Process

- `supplier.invoice.capture/validate/match/post`
- `payable.open_item.settle`
- `payment.proposal.generate/approve/release`
- `customer.invoice.issue/post`
- `receivable.collect/apply`
- `bank.statement.import/commit`
- `bank.reconciliation.run/resolve`
- `cash.position.calculate`
- `cash.forecast.refresh`

Agent 可以建议匹配、付款批次和催收优先级；不得自行修改银行账户、释放重大付款或伪造回单。

## 7. 验收

- AP/AR 未清项与总账控制科目一致；
- 银行账面余额与可信银行余额可解释对账；
- 重复发票、重复回单和重复付款失败关闭；
- 跨法人付款、收款与内部交易不混账；
- 所有付款可穿透审批、银行指令、Observation、清账和凭证。
