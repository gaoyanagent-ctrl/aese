---
id: PLAN-M18-001
title: M18 多产品与多客户组合运营实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m18, portfolio]
---

# M18 多产品与多客户组合运营实施计划

依赖 M17 `integrated_plan_cycle_closed=true`。已完成第二产品、第二客户、组合需求/利润/服务视图、共享 A 线能力分配、资格与现金硬约束、独立审批和无自动业务写入验证。机器合同位于 AESE3 Program M18 frame；2 产品、2 客户、0 shared-capacity violation，输出 `portfolio_operating_model_validated=true`。

- [x] 冻结产品/客户/组合/分配 stable code、单位、owner 与 source refs。
- [x] 验证共享 BOM/routing/capacity/inventory 不被重复使用。
- [x] 对齐 World 后果与 IAOS `portfolio.allocate` 治理动作。
- [x] 完成 schema、fixture、100 次 hash、API/UI 和回归。

边界：组合推荐不创建订单、排产、采购或客户承诺。
