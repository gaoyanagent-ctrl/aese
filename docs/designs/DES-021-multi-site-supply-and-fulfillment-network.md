---
id: DES-021
title: M19 多基地供应与履约网络
date: 2026-07-22
status: completed
author: Codex + User
tags: [m19, network, multi-site, supply-chain]
---

# M19 多基地供应与履约网络

## 1. 目标

把单一苏州基地扩展为最小网络：一个新增制造/外协节点、区域仓或中转节点及多供应商路径，验证跨节点分配、调拨、运输和韧性决策。

机器终态：`network_operating_model_validated=true`。

## 2. 首版范围

- 苏州主厂 + 一个受治理的第二制造/外协节点 + 一个区域物流节点。
- 两条供应路径、跨节点库存/在途、产能分配和客户履约。
- 一个节点中断 tracer：物流/公用工程/供应风险触发 network replan。
- 只做虚构单币种和简化内部结算，不做真实跨境、海关或税务。

## 3. 核心合同

`NetworkNode`、`Lane`、`TransferOrder`、`InTransitLot`、`CapacityAllocation`、`NetworkPromise`、`Disruption`、`RecoveryPlan`。

World 拥有距离、运输、节点实际能力和中断；IAOS 拥有组织/地点台账、调拨/采购/订单/计划和审批。发出不等于到达，调拨单不等于在途物理移动，系统可用量不等于网络实际可用量。

## 4. 完成标准

- 节点、lane、lot、in-transit、ownership 和 tenant/法人边界明确。
- 网络物料、产能、交期、运输成本和现金守恒。
- 中断发现、Knowledge 延迟、重分配、客户承诺和恢复可重放。
- M18 portfolio 和 M17 IBP 可生成/消费 network option，无旁路写入。

## 5. 非目标

全球贸易合规、海关、真实运价、任意规模网络优化和完整 intercompany accounting 不在范围。
