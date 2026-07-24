---
id: ADR-005
title: Process Definition 是工作项运行时的唯一事实源
date: 2026-07-24
status: approved
author: Codex + User
tags: [aese, iaos, process, runtime, governance]
---

# Process Definition 是工作项运行时的唯一事实源

## 背景

M9 曾同时维护 Process Studio 的 `process_definition` 和 API 包内手写的工作项顺序。
两者漂移后，流程图包含 `founder.resolution.prepare`，实际运行却直接进入
`founder.resolution.approve`。这种“双重事实源”会使客户配置、审计视图和业务执行不一致。

## 决策

1. 已发布 Process Definition 是可执行节点、顺序、子流程和分支的唯一事实源。
2. Runtime 必须从指定 Process version 的不可变 Artifact 编译工作项，禁止再维护手写
   Capability 顺序。
3. 主流程中的 subprocess 必须递归展开；normal、recovery、timeout 等边必须保留明确
   语义，恢复节点不得混入正常主线。
4. 发布前必须失败关闭校验循环、未知 Capability、重复正常节点、无参与主体、
   Capability Gate 不一致和不可编译节点。
5. Process Run 必须绑定 `process_key`、`process_version`、`artifact_hash` 和编译器版本。
   已开始的 Run 不因新版本发布而隐式改变。
6. 迁移只允许显式执行：未推进案件可重建；已推进案件必须保留原版本，另行生成迁移
   预览和审批，不得静默重排审计历史。

## 结果

- Process Studio 中发布的配置决定 Runtime 实际创建的工作项。
- Analyzer/发布门可以在业务事实产生前发现配置错误。
- UI、Agent、审批、World wait 和审计共享同一版本化流程语义。
- 流程变更需要版本发布和迁移治理，不能靠修改后端数组热替换。

## M9 落地

Runtime 1.3.8 增加 `CompileInteractiveWorkItems`，递归展开一主四子流程生成18个正常
工作项；`registration.correction.resubmit` 是恢复分支，不进入正常主线；
`enterprise.readiness.evaluate` 只在主流程执行一次。Runtime 安装调用同一编译器和
发布门，API 不再持有手写顺序。
