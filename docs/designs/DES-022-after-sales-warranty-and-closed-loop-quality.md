---
id: DES-022
title: M20 售后质保与闭环质量
date: 2026-07-22
status: completed
author: Codex + User
tags: [m20, warranty, quality, after-sales]
---

# M20 售后质保与闭环质量

## 1. 目标

把“客户接受”后的生命周期补齐为现场问题、投诉、遏制、退货/换货、质保、8D/CAPA、追溯和设计/过程反馈闭环。

机器终态：`customer_lifecycle_closed=true`。

## 2. 首版范围

- 一个虚构 field failure、一个受影响 lot 范围和一个未受影响对照 lot。
- Customer Complaint、RMA、Containment、Trace、8D、CAPA、Warranty Decision 和 Credit/Replacement。
- 质量 Agent、人类接管、供应商责任和 product/process revision feedback。
- 简化 warranty reserve/actual cost，不实现完整法定会计。

## 3. 三态边界

客户现场实际故障、退回物流、物理检验和 replacement 到达属于 World；投诉、RMA、质量问题、8D、CAPA、贷项/换货业务记录属于 IAOS；角色只看到其权限和 observation 范围。

发出遏制通知不等于现场已隔离，RMA 不等于货已退回，CAPA closed 不等于故障原因现实消失。

## 4. 完成标准

- serial/lot/product/process/supplier/customer trace 完整且保密隔离。
- complaint quantity、return、scrap、replacement、credit 和 warranty cost 守恒。
- 假阳性扩大召回、漏召回、重复赔付、越权关闭和无 evidence root cause 失败关闭。
- 经验进入 M18 revision、M17 plan 和 M16 assurance，而不静默改历史 release。

## 5. 非目标

真实法规召回、法律诉讼、保险、全生命周期 PLM 和任意复杂收入冲销不在范围。
