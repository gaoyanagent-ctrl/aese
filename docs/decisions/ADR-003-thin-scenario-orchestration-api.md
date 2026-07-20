---
id: ADR-003
title: AESE 提供无状态场景编排 API
date: 2026-07-20
status: accepted
author: Codex + User
tags: [architecture, orchestration, api, iaos, m7]
---

# AESE 提供无状态场景编排 API

## 背景

M6 已让用户在浏览器中检查 IAOS 联动并进入 Live，但场景的 `reset`、`apply`、`replay`、`agent-run` 和 `verify` 仍需在服务器命令行执行。把这些写入逻辑复制到浏览器会暴露内部合同、弱化权限和审计，也无法可靠处理断线、重复点击和长任务恢复。

ADR-001 允许 AESE 拥有场景内容、校验和 replay orchestration，但 IAOS 必须保持唯一业务运行时和安全边界。因此需要一个可被 UI 调用、但不成为第二套业务平台的服务边界。

## 决策

AESE 可以提供一个薄、无状态的场景编排 API：

- 从版本化 scenario pack 编译确定性执行计划。
- 复用当前 Go loader、validator、projection、replay、agenttrace 和 IAOS client。
- 使用调用者 IAOS JWT 和 tenant context 调用 IAOS 受治理 API。
- 提供预检、分幕推进、运行到结束、Agent 分析、verify 和 reset 编排。
- 从 IAOS 的 scenario run、event cursor、recommendation 和 audit 事实恢复状态。

该服务不得：

- 建立自己的业务数据库、库存、订单、流程、权限或审计表。
- 直接写 IAOS PostgreSQL 或发布正式 NATS 事件。
- 保存长期用户凭据，或用共享管理员身份替代调用者身份。
- 在浏览器中执行任意命令、上传任意脚本或绕过 pack allowlist。
- 把 HCTM 编排硬编码为不可复用的 HTTP handler。

## 运行边界

```text
AESE UI
  -> AESE orchestration API
  -> load + validate versioned pack
  -> compile allowed stage plan
  -> IAOS scenario / simulation / business-action / AI Tool APIs
  -> IAOS RLS + permission + audit + Outbox
  -> IAOS snapshot/cursor/SSE
  -> AESE Live UI
```

AESE API 的运行状态是 IAOS 事实的投影。进程重启后，服务通过 run ID、scenario snapshot 和 event log 恢复，不依赖进程内存认定任务是否完成。

## 认证和确认

- 浏览器传递 IAOS JWT；AESE API 先调用 IAOS profile/permission 合同确认身份和 tenant。
- 所有 IAOS 写入继续由 IAOS 服务端做权限与 RLS 判定。
- `plan` 和只读状态查询不写业务数据。
- `initialize`、`advance`、`run-to-end` 和 `analyze` 需要显式确认。
- `reset` 属于破坏性动作，必须展示影响摘要并提交一次性 confirmation token。
- 请求必须带 idempotency key；重复请求返回原结果或 no-op。

## 结果

正面影响：

- 业务用户无需 CLI 即可运行完整 HCTM 故事。
- 浏览器不需要了解多套 IAOS 写入 API 和对象 UUID 映射。
- 现有 CLI 与 UI 共用同一编排内核，降低行为漂移。
- 权限、租户、审计和业务事实仍集中在 IAOS。

代价：

- 新增一个需要部署和健康检查的 AESE HTTP 进程。
- 必须设计安全的 token 转发、CORS/反向代理和重启恢复。
- 首版只支持 allowlist pack 和确定性阶段，不是通用工作流引擎。

## 被否决方案

1. 浏览器直接依次调用 IAOS 写入 API：合同泄漏、难以恢复且容易产生部分执行。
2. 在 AESE 新建业务数据库和任务队列：复制 IAOS 平台职责。
3. 把 HCTM pack 硬编码进 IAOS Platform：平台与行业场景耦合。
4. 继续只提供 CLI：无法让非研发用户真正操作仿真环境。
