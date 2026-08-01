---
id: MANUAL-M9-001
title: M9 企业设立用户手册
date: 2026-08-01
status: active
author: Codex + User
tags: [manual, m9, incorporation, aese, iaos]
---

# 这份手册解决什么问题

本手册帮助创始治理者理解 M9 不是“一键播放动画”，而是一组由人、Agent、审批流和 AESE
外部参与者共同完成的企业活动。完整节点合同以
[M9 原生闭环设计](../designs/DES-027-m9-iaos-native-incorporation-closed-loop.md) 为准；IAOS
通用概念和配置方式从“帮助与知识中心”查询。

# 四类节点

| 类型 | 谁推进 | 典型输入 | 结果在哪里看 |
| --- | --- | --- | --- |
| human_task | 当前岗位用户 | 表单、选择、确认 | 我的经营待办、Entity、Journal |
| agent_task | 有 Mandate 的 IAOS Agent | 业务说明和已有事实 | Agent 任务、草稿输出、Tool Call |
| approval | Approval Flow 解析的实际审批人 | 事项快照、意见、决定 | 审批中心、Decision、通知 |
| world_wait | AESE 外部机构或客观世界 | Observation | World Journal、IAOS World Bridge |

# 推荐测试路径

1. 在 AESE 创建或选择 Workspace，并进入企业世界。
2. 在 IAOS“企业生命周期”新建设立案。
3. 在“我的经营待办”只处理当前 waiting 节点；locked 节点不能提前推进。
4. Agent 节点先审阅输入来源和授权，再点击执行并查看草稿。
5. Approval 节点到“审批中心”查看完整事项和实际处理人。
6. World wait 节点回 AESE 由登记机构、银行等外部参与者提交一次幂等 Observation。
7. 回 IAOS 打开节点全链，核对 Capability、Process、Entity、Journal、Outbox 和 World 证据。
8. 完成账套、期间、资本凭证和经营就绪检查。

# 常见误解

- AESE 画面中的公司代码不应是脱离 IAOS 案件的演示常量。
- Agent 输出是草稿或建议，不等于审批决定。
- 点击外部确认多次不应增加不同事实；相同幂等键应返回同一结果。
- `completed` 必须有实际输出和证据，不能仅因为剧情帧前进。
- 设计文档说明预期行为；当前是否可执行必须查看 IAOS Active Runtime Artifact。

# 知识查询示例

在 IAOS Copilot 中可以询问：

- “当前节点为什么由 Agent 执行，它读取哪些数据？”
- “这个审批人的来源是什么？”
- “登记机构确认后为什么 IAOS 仍在等待？”
- “capital.commitment.record 的输入、输出和写入 Entity 是什么？”
- “当前 Process Artifact 与用户手册是否一致？”
