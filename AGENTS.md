# AESE Agent Instructions

本文件约束所有在 `/iaos/aese` 工作的 agent 和开发者。

## 1. 必读顺序

开始任何设计或实现前，必须依次阅读：

1. `README.md`
2. `docs/agent-project-context.md`
3. `docs/README.md`
4. `docs/architecture.md`
5. `docs/roadmap.md`
6. `docs/code-map.md`
7. `docs/progress-log.md`

涉及华辰场景时，再按任务读取对应的 HCTM 文档。不要默认一次加载所有长文档。

## 2. 项目边界

- AESE 是 IAOS 的智能企业运行仿真环境，不是独立 ERP、MES 或游戏项目。
- AESE 仓库负责行业场景包、机器可读场景数据、合同校验、演示编排和仿真工具。
- `/iaos/iaos-go` 负责业务运行时、数据库、RLS、Outbox、NATS、Capability、Process、Agent Tool 和前端工作台。
- 不在 AESE 中复制 IAOS 的元数据引擎、流程引擎、权限体系或业务数据库。
- 第一阶段只覆盖华辰热管理系统集团、苏州制造基地、电池冷却板 A 线和第一条订单加急故事。
- 业务真实性、可重放事件和 IAOS 能力映射优先于 2D/3D 表现。

跨仓库修改规则：

- 修改 `/iaos/iaos-go` 前必须阅读其 `AGENTS.md`、`docs/agent-project-context.md` 和 `docs/code-map.md`。
- IAOS 代码变更必须在独立 branch/worktree 中完成，不得从 AESE 提交混入 IAOS 文件。
- AESE 和 IAOS 的提交、测试和进展记录必须分别维护，并在 AESE 设计/计划文档中互相引用。

## 3. 工作流程

开始工作前：

- 检查 `git status`，保留用户或其他 agent 的现有改动。
- 从 `docs/roadmap.md` 和当前 active plan 确认任务属于哪个里程碑。
- 用 `docs/code-map.md` 定位文件和 IAOS 集成点。
- 实质性实现必须有设计文档或实施计划；小型文档修正可以直接进行。

完成工作前：

- 运行与改动范围相称的验证。
- 更新受影响的设计、计划、`docs/code-map.md` 和 `docs/roadmap.md`。
- 在 `docs/progress-log.md` 追加进展记录。
- 检查所有新增 Markdown 链接和 JSON 文件格式。
- 明确报告已验证内容、未验证内容和剩余风险。

## 4. 文档治理

文档分类：

- `docs/designs/DES-NNN-*.md`：架构或功能设计。
- `docs/decisions/ADR-NNN-*.md`：重要且长期有效的技术决策。
- `docs/plans/YYYY-MM-DD-*.md`：可执行实施计划和任务状态。
- `docs/solutions/SOL-NNN-*.md`：已定位并解决的问题。
- `docs/runbooks/*.md`：运行、演示、排障和发布操作。
- `docs/HCTM_*.md`：华辰领域蓝图和业务规格。

DES、ADR、SOL 文档必须包含：

```yaml
---
id: DES-001
title: 文档标题
date: 2026-07-19
status: draft | approved | active | completed | superseded
author: Codex + User
tags: [aese]
---
```

`docs/README.md` 是文档总索引；新增、重命名、废弃文档时必须同步更新。

## 5. 路线图与进展

- `docs/roadmap.md` 是当前里程碑状态和下一步优先级的权威来源。
- `docs/progress-log.md` 是只追加的历史记录，不作为当前状态唯一来源。
- `docs/plans/` 中只能有一个当前 active 的主实施计划；并行子计划必须说明依赖和所有者。
- 不得把“写完设计”标记为“实现完成”。文档、代码、集成和演示验收是不同状态。

每次实质性进展后，按以下格式更新 `docs/progress-log.md`：

```text
## YYYY-MM-DD - 简短标题

- 变更：
- 原因：
- 影响：
- 验证：
- 后续：
```

## 6. Code Map 规则

以下变化必须在同一提交更新 `docs/code-map.md`：

- 新增、移动、重命名或删除命令入口、核心包、场景包、schema、seed、脚本或前端主要入口。
- 改变 AESE 与 IAOS 的集成端点、事件 subject 或数据流。
- 新增必须由后续 agent 优先阅读的约束或运行入口。

纯文字修正且不改变导航关系时，可以不更新 code map，但最终说明原因。

## 7. 数据与事件规则

- 场景数据必须确定性、可版本化、可重复导入和可按场景重置。
- 主数据按稳定业务编码引用，禁止在场景包中硬编码数据库 UUID。
- 时间使用 RFC 3339；业务时区显式写为 `Asia/Shanghai`。
- 数量和金额不得依赖浮点隐式舍入；schema 必须定义单位和精度。
- 事件必须符合 `iaos.{tenant_id}.{event_type}`，并携带 `correlation_id` 和 `idempotency_key`。
- 持久化业务变化应通过 IAOS API、Capability 或 Process 进入事务和 Outbox，不得在正式路径直接写 IAOS 数据库。
- 外生仿真事件必须通过受治理的 simulation ingress 进入 IAOS；直接发布 NATS 只允许本地开发模式，并必须显式标记。
- 不在提交中放入真实客户、个人或生产凭据；当前 HCTM 数据必须保持虚构。

## 8. 实现与测试规则

- AESE 工具链优先使用 Go，与 IAOS 技术栈保持一致。
- 结构化数据使用 JSON Schema 和标准 JSON 解析，不用字符串拼接实现校验。
- 校验器必须离线运行；导入和重放必须显式指定目标环境。
- 默认命令必须是 dry-run，产生外部写入的命令必须要求显式 `--apply` 或等价参数。
- 单元测试覆盖 schema、引用完整性、事件顺序和幂等键。
- IAOS 集成测试必须验证租户隔离、重复执行和失败后无部分写入。
- 前后端实现完成后必须提供可复现的 API 和 UI 验收步骤。

## 9. Git 与交付

- 不回滚不属于当前任务的改动。
- 不使用破坏性 Git 命令清理工作区。
- AESE 改动提交到 AESE 仓库；IAOS 改动提交到对应 IAOS worktree。
- 提交信息应描述交付结果，例如 `Add M3 scenario package plan`。
- 推送前检查工作区、文档索引、code map、路线图和进展日志是否一致。

## 10. System Atlas 更新要求

- 每次实质性进展必须同时更新 `docs/progress-log.md`、提交一个 `atlas-updates/*.json` 声明，并通过 `scripts/check_system_atlas_tracking.sh`。
- 主分支部署通过 `scripts/sync_system_atlas_updates.sh` 幂等同步声明；`scripts/record_system_atlas_update.sh` 仅用于修复和历史补录。
- 更新必须引用设计文档、测试证据或 commit；完成度是架构判断，不得按提交数量自动计算。
- 如果 IAOS API 暂时不可用，先在进展日志记录待补登记项，服务恢复后立即补录。
