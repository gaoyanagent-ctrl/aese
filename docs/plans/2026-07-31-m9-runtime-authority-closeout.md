---
id: PLAN-M9-AUTHORITY-001
title: M9 Command Gateway 与 Effective Process Artifact 运行权威收口
date: 2026-07-31
status: completed
author: Codex + User
tags: [m9, command-gateway, effective-runtime-artifact, process]
---

# M9 Command Gateway 与 Effective Process Artifact 运行权威收口

## 目标

兑现 ADR-003、ADR-005、DES-027 和 IAOS DES-072 已批准但尚未完全落地的两条边界：

1. AESE 浏览器的企业设立写操作只调用同源、白名单 Command Gateway；Gateway
   携带当前用户 JWT 与租户上下文调用 IAOS，AESE 不保存业务数据和特权凭据。
2. IAOS 企业设立运行时只从已发布 Effective Process Artifact 生成工作项；父流程
   固定子流程及 Capability Artifact 的版本和哈希，缺失、漂移或完整性失败时关闭执行。

## 实施任务

- [x] C1 增加 AESE `/api/aese/v1/commands/iaos/*` 白名单命令网关。
- [x] C2 将案件创建、工作项、Agent、审批和 World Observation 写操作切换到网关。
- [x] C3 M9 Capability 安装生成可校验的不可变 Artifact，并让原生适配器执行前校验。
- [x] C4 Process Artifact 锁定子流程和 Capability Artifact 依赖。
- [x] C5 工作项运行时从已发布主流程 Artifact 递归展开，不再读取 Go 目录序列。
- [x] C6 增加网关白名单、身份透传、Artifact 完整性和流程展开回归测试。
- [x] C7 更新设计、Code Map、Roadmap、进展日志与 System Atlas。

## 完成门

- 浏览器 Network 中不存在直达 IAOS 的 M9 `POST /api/v1/...`。
- 任意非白名单 IAOS 路径经 AESE Gateway 返回 403。
- 缺少或篡改 Capability/Process Artifact 时不创建工作项、不执行领域写入。
- 用户重新发布子流程不会静默改变已发布父流程；父流程重新发布后才采用新依赖。

