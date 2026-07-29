---
id: SOL-007
title: 实缴资本核验差异 422 与纠正审批
date: 2026-07-28
status: completed
author: Codex + User
tags: [aese, genesis, capital, approval]
---

# 实缴资本核验差异 422 与纠正审批

## 现象与证据

案件 `INC-GX-695553068F0AE561` 的认缴资本为 100,000,000 分，本次 G4 intent 和
Execute 输入为 80,000,000 分。IAOS 返回
`capital_contribution_mismatch`，保留差异、Decision、Journal 和 Outbox 证据，案件仍在
`bank_account_opened`。

## 修复

- G4 页面同时展示认缴、本次到账和差额；
- 金额不一致时禁止核验，并提供一键按认缴金额修正；
- 金额进入审批 correlation，使 80 万改为 100 万时创建新的 G4 审批，不能复用旧金额
  的批准；
- 同一金额重试仍保持幂等。

IAOS 规则保持“首版实缴必须等于认缴”，未通过前端绕过。
