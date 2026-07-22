---
id: DES-020
title: M18 多产品与多客户组合运营
date: 2026-07-22
status: completed
author: Codex + User
tags: [m18, portfolio, product, customer]
---

# M18 多产品与多客户组合运营

## 1. 目标

在 M17 计划合同上增加第二虚构客户和第二产品，使共享产线、工装、材料、质量能力、信用和现金约束产生真实组合权衡。

机器终态：`portfolio_operating_model_validated=true`。

## 2. 首版范围

- 保留 `HCTM-BCP-A01`，新增一个热管理产品和一个虚构客户。
- 新产品独立 RFQ、revision、BOM/routing、APQP/PPAP、定价和成本证据。
- 两客户订单、优先级、服务等级、信用/付款条件和共享 A 线能力。
- 一个组合冲突 tracer：旺季需求导致容量/材料/现金无法同时满足。

## 3. 所有权与模型

AESE 拥有客户外生需求、物理资源竞争和交付后果；IAOS 拥有 customer/product/project/order/plan/policy；Knowledge 按客户、岗位和保密范围隔离。

新增 `PortfolioItem`、`CustomerProgram`、`SharedConstraint`、`AllocationOption`、`PortfolioDecision` 和稳定 revision/release ancestry。禁止以产品平均值掩盖单客户质量、成本或服务风险。

## 4. 完成标准

- 两产品 master/release/lot/cost 全链稳定编码且互不串线。
- 共享材料、工装、设备、人员、库存、信用和现金守恒。
- allocation 决策有职责、客户承诺、机会成本和 World consequence。
- 单产品 M9-M17 回归不变；组合计划、执行和损益可独立下钻。

## 5. 非目标

第二工厂、跨法人调拨、完整产品组合优化、真实客户定价和大规模 SKU 管理属于后续。
