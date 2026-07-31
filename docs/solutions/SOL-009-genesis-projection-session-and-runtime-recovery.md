---
id: SOL-009
title: Genesis 投影会话竞态与 Runtime stale 自动恢复
date: 2026-07-31
status: completed
author: Codex + User
tags: [aese, genesis, session, projection, runtime-artifact]
---

# Genesis 投影会话竞态与 Runtime stale 自动恢复

## 1. 现象和根因

进入企业后，浏览器可能先请求：

```text
GET /api/aese/v1/game/incorporation/{case}/projection?frame=0
```

若 localStorage 尚保留上一企业的 token/tenant，IAOS RLS 会把当前案件隐藏为 404。页面又把
所有 404 都当成“新案件尚未创建”。随后点击 Agent 操作时，旧 Workspace 的 Runtime 2.4.0
与当前 2.5.0 合同不一致，返回 422 `effective runtime artifact stale`。

## 2. 修复

- 投影遇到 401/404 且存在 Genesis Player + Workspace session 时，先刷新 Workspace
  session，再用服务器返回的 token、tenant、workspace 和 case 重试一次；
- IAOS 写请求只在明确收到 `effective runtime artifact stale` 且第一次请求未产生业务写入
  时刷新 session 并重试一次；
- IAOS Workspace session 在签发 token 前负责对齐受管 Edition，详见
  [IAOS SOL-045](/iaos/iaos-go/docs/solutions/SOL-045-genesis-managed-runtime-session-reconciliation.md)；
- 真正不存在的 case 在刷新后仍返回 404，继续进入新建设立案界面。

## 3. 用户操作和恢复

用户只需重新点击原工作项。系统不会自动跳过节点，也不会把 422 失败请求算作完成。成功后
节点 2 从 `ready` 变为 `completed`，节点 3 解锁为审批状态。

如果 Genesis Player 登录本身已过期，页面会要求重新登录，而不是使用其他租户的旧 token。

## 4. 验证

前端回归测试覆盖：

1. stale tenant 导致 projection 404 → session refresh → 正确 tenant 重试 200；
2. Agent dispatch 返回 Runtime stale 422 → managed session refresh → 原请求重试成功；
3. session 响应的 tenant/token 被原子写回浏览器存储。
