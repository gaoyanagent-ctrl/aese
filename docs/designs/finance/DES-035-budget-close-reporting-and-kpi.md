---
id: DES-035
title: 预算、期间关账、报表与 KPI
date: 2026-07-30
status: active
author: Codex + User
tags: [finance, budget, close, reporting, kpi]
---

# 预算、期间关账、报表与 KPI

## 1. 目的

建立“计划—执行—预测—关账—报表—分析”闭环。预算是授权和管理口径，不是会计凭证；
报表只从已过账、已对账和明确完整度的事实生成。

## 2. 预算与预测

Entity：

- `budget_model`
- `budget_version`
- `budget_line`
- `budget_assumption`
- `budget_allocation`
- `budget_consumption`
- `rolling_forecast`
- `forecast_scenario`
- `budget_variance`

预算支持法人、BU、成本中心、项目、产品、科目和期间维度。版本经历编制、汇总、审批、发布、
冻结和修订。承诺、实际和预测分别保存，禁止用预算余额替代现金或总账余额。

## 3. 模块期间控制

`accounting_period_control` 的键为：

```text
(tenant, legal_entity, accounting_book, fiscal_year, period, module)
module = GL | AP | AR | FA | INV | COST | CASH | TAX | PROJECT
status = future | open | soft_closed | closed | reopened
```

子模块先关闭，GL 最后关闭。`soft_closed` 只允许特许调整；重开必须审批、理由、有效窗口和
审计证据。跨期单据由业务日期、会计日期、模块状态和迟到策略共同决定。

## 4. 月结编排

Entity：

- `period_close_run`
- `period_close_task`
- `close_dependency`
- `close_exception`
- `reconciliation_run/item`
- `accrual_run`
- `revaluation_run`
- `consolidation_run`

典型顺序：

1. 锁定主数据和规则版本；
2. 银行、AP、AR、库存、资产、成本子账处理；
3. 暂估、折旧、汇率重估和成本结算；
4. 子账—总账对账；
5. GL 软关闭、调整审批、最终关闭；
6. 生成并审批报表快照；
7. 集团汇总、内部交易匹配和抵销。

失败任务阻断依赖节点，不允许静默跳过。

## 5. 财务报表

- 试算平衡
- 资产负债表
- 利润表
- 现金流量表
- 所有者权益变动表
- 科目明细/总账
- 应收应付账龄
- 资金与银行对账
- 成本、存货和资产报表
- 预算执行和预测差异

每个报表版本记录法人/集团、账簿、期间、准则、币种、汇率、规则、数据完整度、审批和
发布时间。所有金额可下钻余额、凭证、会计事件和业务证据。

## 6. 管理驾驶舱

核心 KPI：

- 收入、毛利、EBITDA、经营现金流；
- DSO、DPO、CCC、逾期应收应付；
- 库存周转、WIP、报废、质量损失；
- 单位标准/实际成本、采购/用量/效率/费用差异；
- 现金可用天数、预算消耗率、预测准确率；
- 资产利用率、CAPEX 执行和维护成本。

KPI 必须显示口径、期间、币种、完整度和来源，不能把 Preview 或 Agent 推测标为实际。

## 7. Capability / Agent

- `budget.prepare/consolidate/approve/revise`
- `forecast.refresh`
- `period.close.plan/execute/approve/reopen`
- `reconciliation.run/resolve`
- `financial.statement.generate/approve/publish`
- `management.report.generate`
- `financial.kpi.calculate`

Agent 可建议暂估、异常优先级、现金压力和预测变化；关账、重开和报表发布需有权主体决定。

## 8. 验收

- 预算、承诺、实际、预测和现金严格分离；
- 每个期间模块状态可解释，关闭后普通业务写入失败；
- 关账任务可恢复、可重跑且无重复凭证；
- 三表勾稽、子账/总账对账和集团抵销有证据；
- 报表/KPI 不因缺失数据静默配平。
