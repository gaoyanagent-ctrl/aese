# AESE Documentation

本目录是 AESE 的项目知识库。当前状态、架构边界、代码导航和实施计划分别由明确文件维护，避免依赖聊天记录恢复上下文。

## 1. 权威信息来源

| 信息 | 权威文件 |
| --- | --- |
| 项目定位和快速入口 | `README.md` |
| Agent 入门上下文 | `docs/agent-project-context.md` |
| 系统边界与数据流 | `docs/architecture.md` |
| 当前里程碑和优先级 | `docs/roadmap.md` |
| 任务到文件导航 | `docs/code-map.md` |
| 历史进展 | `docs/progress-log.md` |
| 当前实施计划 | `docs/plans/2026-07-19-m4-governed-simulation-ingress.md` |

若内容冲突，当前状态以 `docs/roadmap.md` 为准，长期架构决策以 ADR 为准。

## 2. 文档分类

### 项目级文档

| 文档 | 用途 | 状态 |
| --- | --- | --- |
| [Agent Project Context](agent-project-context.md) | 新 agent 速读入口 | Active |
| [Architecture](architecture.md) | AESE/IAOS 边界和运行架构 | Active |
| [Roadmap](roadmap.md) | 里程碑、当前状态和下一步 | Active |
| [Code Map](code-map.md) | 任务到文件、跨仓库集成点 | Active |
| [Progress Log](progress-log.md) | 只追加的历史记录 | Active |
| [MVP Blueprint](AESE_MVP_Blueprint.md) | MVP 产品与业务范围 | Approved |

### 华辰场景文档

| 文档 | 用途 | 状态 |
| --- | --- | --- |
| [Virtual Enterprise Blueprint](HCTM_Virtual_Enterprise_Blueprint.md) | 集团、工厂、产线和业务蓝图 | Approved |
| [Master Data Model](HCTM_Master_Data_Model.md) | 28 个对象的建模规格 | Approved |
| [Event Model](HCTM_Event_Model.md) | 18 类事件和 IAOS subject/payload | Approved |
| [Seed Data Plan](HCTM_Seed_Data_Plan.md) | Seed 数据和 22 步事件序列 | Approved |
| [Demo Story 01](HCTM_Demo_Story_01_Order_Expedite.md) | 第一条可执行演示 runbook | Approved |
| [M3 Local Runbook](runbooks/hctm-m3-local-run.md) | validate/inspect/apply/replay/verify/reset 操作与安全边界 | Active |
| [HCTM → IAOS Compatibility](reports/hctm-iaos-compatibility.md) | 28 对象映射、API/权限/幂等缺口和运行证据 | Active |
| [HCTM M3 Execution Evidence](reports/hctm-m3-execution-evidence.md) | apply/replay/verify/reset、租户和幂等验收证据 | Completed |
| [HCTM M3V 2D Sandbox Runbook](runbooks/hctm-m3v-2d-sandbox.md) | 前端启动、操作、验证与故障排查 | Active |
| [HCTM M3V 2D Sandbox Evidence](reports/hctm-m3v-2d-sandbox-evidence.md) | 功能、自动化和三目标视口截图证据 | Completed |
| [HCTM M4 Simulation Ingress Evidence](reports/hctm-m4-simulation-ingress-evidence.md) | 设备停机入口、幂等、跨租户和 Outbox 真实证据 | Active |
| [Original Chat](ChatGPT_20260626.md) | 原始构思记录，仅供追溯 | Archive |

### 设计、决策和计划

| ID | 文档 | 状态 |
| --- | --- | --- |
| ADR-001 | [AESE 与 IAOS 仓库边界](decisions/ADR-001-aese-iaos-repository-boundary.md) | Accepted |
| ADR-002 | [AESE 只读 2D 场景预览界面](decisions/ADR-002-aese-2d-preview-ui.md) | Accepted |
| DES-001 | [M3 可执行场景包与重放架构](designs/DES-001-m3-executable-scenario-package.md) | Completed |
| DES-002 | [快速 2D 企业沙盘设计](designs/DES-002-fast-track-2d-enterprise-sandbox.md) | Completed |
| PLAN-M3 | [M3 实施计划](plans/2026-07-19-m3-executable-scenario-package.md) | Completed |
| PLAN-M3V | [快速 2D 企业沙盘实施计划](plans/2026-07-19-fast-track-2d-enterprise-sandbox.md) | Completed |
| PLAN-M4 | [M4 受治理异常事件入口实施计划](plans/2026-07-19-m4-governed-simulation-ingress.md) | Active |

## 3. 命名与状态

- 设计：`docs/designs/DES-NNN-kebab-case.md`
- 决策：`docs/decisions/ADR-NNN-kebab-case.md`
- 解决方案：`docs/solutions/SOL-NNN-kebab-case.md`
- 计划：`docs/plans/YYYY-MM-DD-kebab-case.md`
- Runbook：`docs/runbooks/kebab-case.md`

状态闭集：

- `draft`：仍在讨论，不作为稳定合同。
- `approved` / `accepted`：已确认，可以指导实现。
- `active`：正在执行。
- `completed`：实现和验收均完成。
- `superseded`：已被其他文档替代，必须给出替代链接。
- `archive`：仅保留历史参考。

## 4. 更新规则

- 新增文档时更新本索引。
- 改变里程碑时更新 `roadmap.md` 和 `progress-log.md`。
- 改变文件入口或集成点时更新 `code-map.md`。
- 形成长期架构取舍时新增或更新 ADR。
- 完成实施任务时更新对应 plan 的 checklist 和验收证据。
