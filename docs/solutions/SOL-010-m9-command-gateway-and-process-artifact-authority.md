---
id: SOL-010
title: M9 Command Gateway 与 Process Artifact 运行权威修复
date: 2026-07-31
status: completed
author: Codex + User
tags: [aese, m9, bff, process-artifact]
---

# M9 Command Gateway 与 Process Artifact 运行权威修复

## 问题

M9 页面读取游戏投影时经过 AESE，但案件、工作项、审批和 World Observation 写操作
由浏览器直接调用 IAOS；同时 IAOS 工作项由 Go 中的 `ProcessDefinition` 展开。前者
泄漏跨系统写接口，后者使用户发布的流程定义不能成为唯一运行事实。

## 修复

- AESE 增加同源 Command Gateway，只接受企业设立所需的精确 POST 路径，透传调用者
  JWT 和 `X-IAOS-Tenant-Id`，不接受任意代理路径。
- 前端所有 M9 写操作切换到 `/api/aese/v1/commands/iaos/*`。
- IAOS 安装器用标准 Runtime Artifact envelope 发布 M9 Capability Artifact。
- 原生 Capability 执行前校验 active publication、source version、compiler 和哈希。
- Process Artifact 固定 Capability 和 subprocess 的 artifact version/hash；设立工作项
  从主 Process Artifact 递归展开，缺少锁定依赖即失败关闭。

## 使用与恢复

用户仍在 AESE 界面创建企业、派发 Agent、审批和提交外部观察，无需输入新配置。
发生 `effective runtime artifact` 错误时，应在 IAOS 发布或升级对应平台包/流程，
不得由 AESE 降级为直写或使用 Go 默认流程继续执行。

