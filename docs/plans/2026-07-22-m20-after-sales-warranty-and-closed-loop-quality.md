---
id: PLAN-M20-001
title: M20 售后质保与闭环质量实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m20, quality, warranty]
---

# M20 售后质保与闭环质量实施计划

依赖 M19 terminal。已完成 field failure、complaint、120 件 RMA、containment、8D/CAPA、replacement/credit 与 lot genealogy 闭环；退回、替换和处置数量差异为 0，输出 `customer_lifecycle_closed=true`。

- [x] 冻结投诉/RMA/质保/8D/CAPA 三态与证据 ancestry。
- [x] 验证 120 件数量、客户接受和财务影响守恒。
- [x] 对齐 IAOS `quality.close` 独立审批与幂等事务。
- [x] 完成 schema、fixture、API/UI 与回归。

边界：World 现场失效不直接关闭 IAOS 质量记录或开具贷项。
