---
id: DES-019
title: M17 滚动 IBP 与 S&OP
date: 2026-07-22
status: completed
author: Codex + User
tags: [m17, ibp, sop, planning]
---

# M17 滚动 IBP 与 S&OP

## 1. 目标

M17 消费 M16 `renewed` 的策略与校准假设，建立第一个受治理滚动经营计划周期：把需求、供应、产能、库存、交付、成本和现金放到同一版本化计划中，经职能审议和管理层决策形成唯一 approved operating plan。

机器终态：

```text
integrated_plan_cycle_closed=true
disposition=approved | replan_required | deferred
```

## 2. 首版范围

- 单工厂、单产品、单客户、当前 resilient release。
- 13 周 weekly execution horizon + 12 个月 monthly financial horizon。
- baseline、upside、downside 三个 planning scenario；不等同于 M14 stochastic run。
- Demand Review、Supply Review、Financial Reconciliation、Pre-IBP 和 Executive IBP 五个 gate。
- 一个 approved plan version、一个 frozen zone、一个 replan tracer。

## 3. 核心模型

`PlanningCycle`、`AssumptionSet`、`DemandPlan`、`SupplyPlan`、`CapacityPlan`、`InventoryPlan`、`FinancialPlan`、`Gap`、`ScenarioOption`、`PlanDecision`、`PlanRelease`。

IAOS 拥有计划版本、审批、任务、权限和 release；AESE World 拥有未来外生实现与资源后果；Knowledge 只保存角色已知假设和版本。计划不得创建客户需求、材料、产能或现金。

## 4. 关键合同

- Bucket、单位、精度、calendar、cutoff、frozen/slushy/liquid time fence 明确。
- Forecast、order、commitment、capacity、budget、cash 和 actual 分离。
- 每个职能计划绑定 assumption/evidence/hash，冲突形成显式 Gap。
- Executive approval 只发布 plan，不自动转 PO/WO/Shipment；执行必须另走 Capability。
- actual feedback 通过 M16 dataset/assurance 进入下一 cycle，不静默改已批准版本。

## 5. 完成标准

- 数量、产能、库存、交付、成本和现金计划跨 horizon 可对账。
- 需求/供应/财务三方冲突、陈旧版本、并发编辑和越权审批失败关闭。
- downside 场景产生可解释 gap、选项、权衡和 replan decision。
- Executive IBP Room、API、runbook/evidence、IAOS 治理和 M3-M16 回归完整。

## 6. 非目标

第二产品/客户、多基地网络、自动优化器、真实预测模型和计划自动执行属于后续里程碑。
