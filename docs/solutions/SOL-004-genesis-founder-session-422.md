---
id: SOL-004
title: Genesis 创建企业 Founder 会话 422
date: 2026-07-28
status: completed
author: Codex + User
tags: [genesis, identity, iaos, incorporation]
---

# Genesis 创建企业 Founder 会话 422

## 现象

玩家在独立 Workspace 中选择 AI 公司身份并点击“确认身份并创建企业”后，
IAOS `POST /api/v1/incorporations/cases` 返回 422：

```text
governance rejected: authenticated principal does not match acting subject or tenant access
```

## 根因

Workspace provisioning 已创建 `founder-principal`、岗位和 Mandate，但返回给浏览器的是
`/dev/switch-tenant` 签发的租户管理员 token。M9 `incorporation.case.open` 的 acting
subject 是 `founder-principal`，IAOS 因认证主体不一致而正确拒绝请求。

## 修复

- Founder bootstrap 后通过普通 `/api/v1/auth/login` 获取 Founder 会话。
- M9 Runtime 安装及浏览器 tenant session 都使用 Founder token。
- 新增 owner-scoped Workspace session refresh；旧页面遇到该特定 422 时刷新一次 Founder
  会话并重放原请求，不创建第二个 tenant，也不放宽 IAOS 治理校验。
- session refresh 只能由 Workspace owner 的玩家标识调用。

## 验证

- 单元测试验证 provisioning 返回 Founder token、Runtime 由 Founder token 安装，以及
  其他玩家不能刷新 Workspace 会话。
- 原始请求先稳定复现 422；刷新同一 Workspace 会话后，同一
  `incorporation.case.open` 请求返回 201 `committed`。
- AESE 全量 Go 测试、前端根路由测试和生产构建通过。
