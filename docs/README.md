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
| 当前实施任务 | 无；PLAN-M17 至 PLAN-M24 已完成 |

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
| [HCTM M4 Simulation Ingress Evidence](reports/hctm-m4-simulation-ingress-evidence.md) | 三类异常入口、状态、幂等、跨租户和 Outbox 证据 | Completed |
| [HCTM M5 Agent Tracer Runbook](runbooks/hctm-m5-governed-agent-tracers.md) | Tool setup、Agent run、重复和跨租户验收步骤 | Active |
| [HCTM M5 Agent Tracer Evidence](reports/hctm-m5-agent-evidence.md) | 三 Agent 在线结果、零业务写入和调用审计证据 | Completed |
| [HCTM M6 Online Sandbox Runbook](runbooks/hctm-m6-online-sandbox.md) | 在线事实、游标、SSE 和 Live UI 操作 | Active |
| [HCTM M6 Online Sandbox Evidence](reports/hctm-m6-online-sandbox-evidence.md) | KPI、幂等、租户隔离和三视口证据 | Completed |
| [HCTM M7 Governed Scenario Operations Console Runbook](runbooks/hctm-m7-governed-scenario-operations-console.md) | 预检、初始化、七幕推进、分析、验证、复位与重启恢复验收清单 | Completed |
| [HCTM M7 Governed Scenario Operations Console Evidence](reports/hctm-m7-governed-scenario-operations-console-evidence.md) | O3/O4 运行控制台联调、可恢复运行、权限边界与三视口证据 | Completed |
| [AESE World F1 Runbook](runbooks/aese-world-f1.md) | 确定性内核 validate/inspect/run/replay、dry-run 和 artifact 合同 | Active |
| [AESE World Play Runbook](runbooks/aese-world-play.md) | Genesis 三态界面、bridge、恢复与全链验收 | Completed |
| [M8 Capability Gap Ledger](capability-gap-ledger.md) | World Bridge 企业活动、角色、权限和后续缺口 | Completed |
| [Original Chat](ChatGPT_20260626.md) | 原始构思记录，仅供追溯 | Archive |
| [AESE 2.0 原始设计构思](ChatGPT20260722-aese2.0.md) | 企业生命周期仿真方向的输入材料；以 DES-007/ADR-004 为工程化解释 | Draft |

### 设计、决策和计划

| ID | 文档 | 状态 |
| --- | --- | --- |
| ADR-001 | [AESE 与 IAOS 仓库边界](decisions/ADR-001-aese-iaos-repository-boundary.md) | Accepted |
| ADR-002 | [AESE 只读 2D 场景预览界面](decisions/ADR-002-aese-2d-preview-ui.md) | Accepted |
| ADR-003 | [AESE 无状态场景编排 API](decisions/ADR-003-thin-scenario-orchestration-api.md) | Accepted |
| ADR-004 | [AESE 仿真世界事实所有权](decisions/ADR-004-aese-world-state-ownership.md) | Accepted |
| ADR-005 | [Process Definition 是工作项运行时的唯一事实源](decisions/ADR-005-process-definition-as-runtime-source.md) | Accepted |
| DES-001 | [M3 可执行场景包与重放架构](designs/DES-001-m3-executable-scenario-package.md) | Completed |
| DES-002 | [快速 2D 企业沙盘设计](designs/DES-002-fast-track-2d-enterprise-sandbox.md) | Completed |
| DES-003 | [M5 受治理 Agent Tracer 设计](designs/DES-003-governed-agent-tracers.md) | Completed |
| DES-004 | [M6 在线 2D 企业沙盘架构](designs/DES-004-online-2d-enterprise-sandbox.md) | Completed |
| DES-005 | [M7 受治理场景运行控制台](designs/DES-005-governed-scenario-operations-console.md) | Completed |
| DES-006 | [AESE System Atlas 双系统全景投影](designs/DES-006-system-atlas-aese-projection.md) | Completed |
| DES-007 | [AESE 2.0 企业生命周期仿真基础架构](designs/DES-007-aese-2-foundation.md) | Draft |
| DES-008 | [AESE World 与 IAOS 三段式桥接合同](designs/DES-008-world-iaos-bridge-contract.md) | Approved |
| DES-009 | [AESE World 三态术语、所有权与数据分类](designs/DES-009-world-contract-model.md) | Approved |
| DES-010 | [M9 Project Genesis 企业成立与治理](designs/DES-010-genesis-incorporation-and-governance.md) | Approved |
| DES-011 | [M10 Project Genesis 工厂选址与设施建设](designs/DES-011-genesis-plant-build.md) | Completed |
| DES-012 | [M11 Project Genesis 生产能力建设](designs/DES-012-genesis-production-capability-build.md) | Completed |
| DES-013 | [M12 Project Genesis 产品工业化与量产批准](designs/DES-013-genesis-product-industrialization.md) | Completed |
| DES-014 | [M13 Project Genesis 第一次完整商业交付](designs/DES-014-genesis-first-commercial-delivery.md) | Completed |
| DES-015 | [M14 参数化分支经营实验](designs/DES-015-parameterized-branch-experiments.md) | Completed |
| DES-016 | [M15 受治理策略发布与经营试点](designs/DES-016-governed-strategy-release-and-pilot.md) | Completed |
| DES-017 | [M16 持续策略保障与假设校准](designs/DES-017-continuous-strategy-assurance-and-calibration.md) | Completed |
| DES-018 | [AESE 3.0 M17-M24 完成体规划](designs/DES-018-aese-3-completion-program.md) | Completed |
| DES-019 | [M17 滚动 IBP 与 S&OP](designs/DES-019-rolling-ibp-and-sop.md) | Completed |
| DES-020 | [M18 多产品与多客户组合运营](designs/DES-020-product-and-customer-portfolio-expansion.md) | Completed |
| DES-021 | [M19 多基地供应与履约网络](designs/DES-021-multi-site-supply-and-fulfillment-network.md) | Completed |
| DES-022 | [M20 售后质保与闭环质量](designs/DES-022-after-sales-warranty-and-closed-loop-quality.md) | Completed |
| DES-023 | [M21 资产人员 EHS 与工厂韧性](designs/DES-023-plant-resource-and-ehs-resilience.md) | Completed |
| DES-024 | [M22 集团财务资金与投资治理](designs/DES-024-group-finance-treasury-and-investment.md) | Completed |
| DES-025 | [M23 受治理多 Agent 组织](designs/DES-025-governed-multi-agent-organization.md) | Completed |
| DES-026 | [M24 场景平台产品化与行业交付](designs/DES-026-scenario-platform-productization.md) | Completed |
| DES-027 | [M9 IAOS 原生语义驱动企业成立真实闭环](designs/DES-027-m9-iaos-native-incorporation-closed-loop.md) | Approved |
| PLAN-M3 | [M3 实施计划](plans/2026-07-19-m3-executable-scenario-package.md) | Completed |
| PLAN-M3V | [快速 2D 企业沙盘实施计划](plans/2026-07-19-fast-track-2d-enterprise-sandbox.md) | Completed |
| PLAN-M4 | [M4 受治理异常事件入口实施计划](plans/2026-07-19-m4-governed-simulation-ingress.md) | Completed |
| PLAN-M5 | [M5 受治理 Agent MVP 实施计划](plans/2026-07-19-m5-governed-agent-mvp.md) | Completed |
| PLAN-M6 | [M6 在线 2D 企业沙盘实施计划](plans/2026-07-20-m6-online-2d-enterprise-sandbox.md) | Completed |
| PLAN-M7 | [M7 受治理场景运行控制台实施计划](plans/2026-07-20-m7-governed-scenario-operations-console.md) | Completed |
| PLAN-M8 | [M8 AESE 2.0 三态世界与仿真内核实施计划](plans/2026-07-22-m8-aese-2-foundation.md) | Completed |
| PLAN-M9 | [M9 Project Genesis 企业成立与治理实施计划](plans/2026-07-22-m9-genesis-incorporation.md) | Completed |
| PLAN-M9-NATIVE | [M9 IAOS 原生语义驱动企业成立真实闭环实施计划](plans/2026-07-23-m9-iaos-native-incorporation-closed-loop.md) | Active — D22 人机协同补强 |
| PLAN-M10 | [M10 Project Genesis 工厂选址与设施建设实施计划](plans/2026-07-22-m10-genesis-plant-build.md) | Completed |
| PLAN-M11 | [M11 Project Genesis 生产能力建设实施计划](plans/2026-07-22-m11-genesis-production-capability-build.md) | Completed |
| PLAN-M12 | [M12 Project Genesis 产品工业化与量产批准实施计划](plans/2026-07-22-m12-genesis-product-industrialization.md) | Completed |
| PLAN-M13 | [M13 Project Genesis 第一次完整商业交付实施计划](plans/2026-07-22-m13-genesis-first-commercial-delivery.md) | Completed |
| PLAN-M14 | [M14 参数化分支经营实验实施计划](plans/2026-07-22-m14-parameterized-branch-experiments.md) | Completed |
| PLAN-M15 | [M15 受治理策略发布与经营试点实施计划](plans/2026-07-22-m15-governed-strategy-release-and-pilot.md) | Completed |
| PLAN-M16 | [M16 持续策略保障与假设校准实施计划](plans/2026-07-22-m16-continuous-strategy-assurance-and-calibration.md) | Completed |
| PLAN-M17 | [M17 滚动 IBP 与 S&OP 实施计划](plans/2026-07-22-m17-rolling-ibp-and-sop.md) | Completed |
| PLAN-M18 | [M18 多产品与多客户组合运营实施计划](plans/2026-07-22-m18-product-and-customer-portfolio-expansion.md) | Completed |
| PLAN-M19 | [M19 多基地供应与履约网络实施计划](plans/2026-07-22-m19-multi-site-supply-and-fulfillment-network.md) | Completed |
| PLAN-M20 | [M20 售后质保与闭环质量实施计划](plans/2026-07-22-m20-after-sales-warranty-and-closed-loop-quality.md) | Completed |
| PLAN-M21 | [M21 资产人员 EHS 与工厂韧性实施计划](plans/2026-07-22-m21-plant-resource-and-ehs-resilience.md) | Completed |
| PLAN-M22 | [M22 集团财务资金与投资治理实施计划](plans/2026-07-22-m22-group-finance-treasury-and-investment.md) | Completed |
| PLAN-M23 | [M23 受治理多 Agent 组织实施计划](plans/2026-07-22-m23-governed-multi-agent-organization.md) | Completed |
| PLAN-M24 | [M24 场景平台产品化与行业交付实施计划](plans/2026-07-22-m24-scenario-platform-productization.md) | Completed |
| AESE 3 Runbook | [M17-M24 Reference Release](runbooks/aese3-reference-release.md) | Completed |
| AESE 3 Evidence | [M17-M24 Completion Evidence](reports/aese3-m17-m24-completion-evidence.md) | Completed |
| M9 Runbook | [Genesis Incorporation Runbook](runbooks/genesis-incorporation.md) | Completed |
| M9 Evidence | [Genesis Incorporation Evidence](reports/m9-genesis-incorporation-evidence.md) | Completed |
| M9N Asset Audit | [IAOS-native M9 machine-readable asset audit](reports/m9-native-asset-audit.json) | Active |
| M9N Frozen Contract | [IAOS-native M9 lifecycle contract](contracts/m9-native-incorporation-contract.json) | Active |
| M9N Final Evidence | [IAOS-native M9 final evidence](reports/m9-native-final-evidence.md) | Completed |
| M9N Closed-loop Runbook | [IAOS-native M9 operations](runbooks/m9-native-closed-loop.md) | Active |
| SOL-001 | [M9 LAN lifecycle loading and SSE continuity](solutions/SOL-001-m9-lan-lifecycle-loading-and-sse.md) | Completed |
| M9N Risk Register | [IAOS-native M9 risk register](reports/m9-native-risk-register.json) | Active |
| M10 Runbook | [Genesis Plant Build Runbook](runbooks/genesis-plant-build.md) | Completed |
| M10 Evidence | [Genesis Plant Build Evidence](reports/m10-genesis-plant-build-evidence.md) | Completed |
| M11 Runbook | [Genesis Capability Build Runbook](runbooks/genesis-capability-build.md) | Completed |
| M11 Evidence | [Genesis Capability Build Evidence](reports/m11-genesis-production-capability-evidence.md) | Completed |

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
