---
id: DES-037
title: AESE 场景知识与 IAOS 产品知识中枢集成
date: 2026-08-01
status: active
author: Codex + User
tags: [aese, knowledge, manual, m9, iaos]
---

# 1. 目标

AESE 用户不仅要知道游戏中“点击什么”，还要理解每个剧情节点对应的真实企业活动、参与者、
IAOS Capability/Process/Entity、输入输出和证据查询路径。IAOS DES-077 负责通用产品知识
Registry、权限、页面帮助和 Copilot；AESE 负责场景手册源及 World/IAOS 映射，不复制知识
服务、权限体系或运行数据库。

# 2. 内容所有权

| 内容 | 所有者 |
| --- | --- |
| IAOS 概念、菜单、配置和错误恢复 | IAOS Product Knowledge Hub |
| AESE 剧情、角色、World Observation 和操作步骤 | AESE 场景手册 |
| 当前 IAOS 租户 Active Artifact 和业务记录 | IAOS Runtime |
| 当前 AESE World 客观事实和 actor knowledge | AESE World |

场景手册不能宣称某 IAOS 节点已完成；必须通过 Bridge/Runtime evidence 证明。

# 3. 场景知识 Article 扩展

AESE 场景文章沿用 IAOS `Knowledge Article` 合同，并在 `related_assets` 中增加：

```json
{
  "aese_routes": ["world-incorporation"],
  "world_actions": ["registrar.confirm_registration"],
  "iaos_processes": ["enterprise.incorporation.lifecycle.v1"],
  "iaos_capabilities": ["registration.observation.commit"],
  "iaos_entities": ["incorporation_case"],
  "evidence_types": ["journal", "outbox", "world_observation"]
}
```

每个节点说明：业务目的、前置事实、责任主体、人工/Agent/外部世界分类、输入表单、校验、预期
输出、实际输出查询、失败恢复、IAOS 深链和 World evidence。

# 4. 发布方式

AESE 平台基线手册保留在 Git 并受场景版本审核；构建时生成 Knowledge Edition manifest，
由 IAOS 平台包安装器发布为 `scenario` 文章。浏览器不直接读取仓库文件。文章必须带
`applies_to_version`、AESE pack hash 和 IAOS Edition 依赖。

人工可读源为 `docs/manuals/m9-incorporation-user-manual.md`；机器 Edition 位于
`scenario-packs/hctm/knowledge/m9-incorporation.json`，Schema 位于同目录 `schemas/`。
`internal/scenarioknowledge` 使用严格 JSON 解码、稳定编码 allowlist、18 节点顺序/完整性和
canonical SHA-256 校验；`aese knowledge validate` 是离线质量门。`aese knowledge compile`
读取受版本控制的 Markdown，将正文摘要写入 `knowledge_article` 资产版本，再签名 Package 和
Edition，产物位于 `scenario-packs/hctm/knowledge/dist/`。`aese knowledge install` 默认仅向
IAOS 做 dry-run 校验，显式 `--apply` 后才按 JWT 租户安装。IAOS 保存 FORCE RLS 行业 Active
版本、审计和安装清单；相同 bundle 重试 no-op，同版本换正文失败关闭。

# 5. Agent 合同

AESE 对话框将当前 workspace、case、world run、剧情节点和用户角色传给 IAOS 知识检索；
回答必须分区显示：

1. 场景中下一步怎么做；
2. 这一步在真实企业中的含义；
3. 对应 IAOS 配置和当前运行状态；
4. World 与 IAOS 两侧证据；
5. 来源文章 ID 和版本。

对话框只提供知识查询时不得调用写 API；推进剧情继续走 AESE Command Gateway、IAOS
Capability、Approval 和 World Observation。

## 5.1 场景导航上下文合同

AESE 从当前 `WorldProjection`、`GameWorkItem` 和 Genesis session 生成以下封闭字段，并通过
知识深链传给 IAOS：

| 字段 | 来源 | 含义 |
| --- | --- | --- |
| `workspace_id` | 当前 Genesis workspace | 企业空间稳定编码 |
| `case_code` | `projection.case_id` | 当前设立案编码 |
| `world_run_id` | `projection.world_run_id` | 当前 World 运行编码 |
| `node_id` | `work_item_id` | 当前剧情/流程节点 |
| `actor_id` | `owner_id` | 当前责任主体 |
| `actor_type` | `owner_type` | 人、Agent、审批或外部主体类型 |
| `task_type` | `kind` | 当前任务类型 |
| `capability` | `capability_code` | 对应 IAOS Capability |

该对象是**导航上下文**，不是 IAOS 运行事实或 World Observation。IAOS 必须：

- 只接受上述字段和稳定编码格式，忽略未知、超长或提示注入文本；
- 在 BFF 再次归一化，不能信任浏览器传入对象；
- 在知识中心可见展示并允许用户清除，不能将其作为隐藏 Agent 指令；
- Copilot 回答当前状态前，重新读取有权访问的 IAOS Runtime/业务 API；无证据时明确失败关闭；
- S9 前不宣称已完成双侧证据核验或配置漂移检测。

# 6. 验收

- 用户从任一 M9 节点可打开对应场景说明；
- 文章能跳转 IAOS Capability、Process、Entity 和证据入口；
- Agent 能回答“谁做、为什么、输入输出、在哪里查看”，并显示来源；
- 文章与运行事实冲突时明确提示漂移；
- AESE 不新增知识数据库，不绕过 IAOS 权限读取跨租户内容。

# 7. 当前实现入口

- `WorkItemActionPanel` 在每个 M9 任务显示“这一步是什么”，携带 tenant、case、capability 和
  `KB-M9-INCORPORATION` 打开 IAOS 知识中心；深链同时携带 5.1 节定义的封闭场景导航上下文；
- IAOS 知识中心明确展示该上下文，Copilot 请求经前端和 BFF 双重归一化后使用；它不能替代
  Runtime、Journal、Outbox 或 World Observation 证据；
- 清单逐节点固定 purpose、inputs、outputs、actor、task type、gate、evidence、IAOS menu 和
  World action；
- 内容哈希不等同于生产签名或安装完成；只有 IAOS 安装器登记成功后才能宣称 Edition 已发布。
