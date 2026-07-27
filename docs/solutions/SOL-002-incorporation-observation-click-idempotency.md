---
id: SOL-002
title: 企业成立外部确认重复点击幂等
date: 2026-07-27
status: completed
author: Codex + User
tags: [aese, incorporation, world-bridge, idempotency]
---

# 问题

用户在企业成立 World 中重复点击“登记机构确认登记”时，每次都会生成新的
Observation，IAOS Journal 的 JSON 持续增长。重复事实还触发旧 IAOS
`count != 1` 缺陷，使节点 7 误报缺少可信 Observation。

# 修复

AESE 不再把当前时间写入同一业务结果的 transport identity。同一案件、payload type
和 result 形成稳定 `message_id`/`idempotency_key`。IAOS 同时按业务事实执行服务端
去重，以兼容已部署的旧 AESE 页面和网络重试。

World 按钮只负责登记外部事实，仍不会直接越过 IAOS 的 `world_wait` 工作项。用户回到
IAOS 执行节点时，Runtime 消费已登记的可信事实并提交状态转换。

# 验证

- TypeScript 测试与生产构建。
- IAOS World Bridge recovery matrix：不同 transport id 重发相同事实返回 duplicate。
- IAOS 完整 lifecycle：存在重复历史 Observation 时节点 7 仍可提交。
