---
id: SOL-012
title: M10 LAN HTTP 缺少 randomUUID
date: 2026-08-02
status: completed
author: Codex + User
tags: [aese, m10, frontend, compatibility]
---

# M10 LAN HTTP 缺少 randomUUID

## 现象与根因

在 `http://192.168.50.222:4173` 点击“发起外部调研工作项”没有网络请求，控制台显示
`crypto.randomUUID is not a function`。部分浏览器只在安全上下文暴露 randomUUID；M10
页面直接调用该 API，导致事件处理器提前中断。

## 修正

前端统一使用 `createClientRequestId`：优先 randomUUID；不可用时由 getRandomValues 生成
符合 v4/variant 位的 UUID；极老环境使用时间戳、进程内序列和随机片段兜底。调研请求、
World Observation、场址推荐和 Workspace 幂等键均不再直接访问 randomUUID。

客户端编码不是身份或授权凭证。AESE 仍从 IAOS Session 派生 actor，并由 Capability、
World Journal、RLS、Audit 和 Outbox 约束正式写入。
