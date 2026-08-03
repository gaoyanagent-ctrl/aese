# AESE Code Map

本文件把常见任务映射到应优先阅读和修改的文件。M3 离线场景包和工具路径已经创建。

## 1. 快速入口

| 任务 | 先读 |
| --- | --- |
| 理解项目定位 | `README.md`、`docs/agent-project-context.md` |
| 查看当前进度 | `docs/roadmap.md`、`docs/progress-log.md` |
| 理解 AESE/IAOS 边界 | `docs/architecture.md`、ADR-001 |
| 运行或修改 2D 沙盘 | DES-002、M3V completed plan、M3V runbook |
| 修改 M4 异常入口 | M4 completed plan、M4 evidence、`internal/iaosclient/`、`internal/replay/` |
| 修改 M5 Agent tracer | DES-003、M5 completed plan、`internal/agenttrace/`、`scenario-packs/hctm/agent-tools.json` |
| 维护 M6 在线沙盘 | DES-004、M6 completed plan、`frontend/src/scenario/`、IAOS scenario API |
| 维护 M7 场景运行控制台 | ADR-003、DES-005、M7 completed plan、现有 CLI application service |
| 维护 AESE 2.0 World Runtime | ADR-004、DES-007、DES-009、M8 completed plan；先读 `world-contracts/` 和 `internal/worldcontract/` |
| 维护 M9 企业成立与治理 | DES-010、PLAN-M9-001、M9 evidence；先读 `internal/incorporation/`、`world-packs/hctm-genesis/campaigns/incorporation/` 与 World Bridge |
| 维护 M10 工厂选址与设施建设 | DES-011、DES-038、PLAN-M10-INTERACTIVE-001、M10 runbook；先读 `internal/plantbuild/interactive.go`、`frontend/src/components/world/PlantBuildPlay.tsx`、AESE BFF 与 IAOS `plant_interactive.go`。所有 Genesis 主线必须遵守最小玩家输入和无死路门禁；旧 campaign/evidence 只证明 fixture replay |
| 维护 M11 生产能力建设 | DES-012、PLAN-M11-001、M11 evidence；先读 `internal/capabilitybuild/`、Capability Build campaign 与 IAOS DES-053 |
| 维护 M12 产品工业化与量产批准 | DES-013、PLAN-M12-001、M12 evidence；先读 `internal/industrialization/`、Industrialization campaign 与 IAOS DES-054 |
| 维护 M13 第一次完整商业交付 | DES-014、PLAN-M13-001、M13/Genesis evidence；先读 `internal/firstdelivery/` 与 First Delivery campaign |
| 维护 M14 参数化分支经营实验 | DES-015、PLAN-M14-001、M14 evidence；先读 `internal/experiment/`、`internal/genesisexperiment/` 与 Scenario Lab |
| 维护 M15 受治理策略发布与经营试点 | DES-016、PLAN-M15-001、M15 evidence；先读 `internal/strategyrelease/` 与 Strategy Control Room |
| 维护 M16 持续策略保障与假设校准 | DES-017、PLAN-M16-001、M16 evidence；先读 `internal/assurance/` 与 Assurance Observatory |
| 验证 AESE 3.0 M17-M24 完成体 | DES-018 至 DES-026、PLAN-M17 至 PLAN-M24、`internal/aese3/`、`#world-aese3` |
| 发布 HCTM reference pack | `world-packs/hctm-genesis/manifest.json`、`world-packs/hctm-genesis/aese3/`、AESE 3 runbook |
| 查看或维护双系统全景 | DES-006、`frontend/src/components/SystemAtlas.tsx`、IAOS System Atlas API |
| 修改华辰企业设定 | `docs/HCTM_Virtual_Enterprise_Blueprint.md` |
| 修改对象和字段 | `docs/HCTM_Master_Data_Model.md` |
| 修改事件名和 payload | `docs/HCTM_Event_Model.md` |
| 修改 seed 数据 | `docs/HCTM_Seed_Data_Plan.md` |
| 修改演示流程 | `docs/HCTM_Demo_Story_01_Order_Expedite.md` |

## 2. 当前文件地图

| 区域 | 文件 | 职责 |
| --- | --- | --- |
| Agent 规则 | `AGENTS.md` | 工作流、边界、文档和测试规则 |
| 项目入口 | `README.md` | 产品简介和文档入口 |
| System Atlas 进展登记 | `scripts/record_system_atlas_update.sh` | 向 IAOS 追加设计、实现、测试、发布、决策或风险记录 |
| 文档索引 | `docs/README.md` | 文档分类、状态和编号 |
| 项目上下文 | `docs/agent-project-context.md` | Agent 快速入门 |
| 架构 | `docs/architecture.md` | 仓库边界、数据流、安全和可重复性 |
| 路线图 | `docs/roadmap.md` | 里程碑状态；当前 active 主计划为 PLAN-GXZ-001 零起点主页、独立租户与真实 AI |
| 代码导航 | `docs/code-map.md` | 本文件 |
| 历史记录 | `docs/progress-log.md` | 只追加进展日志 |
| 薄编排服务入口 | `cmd/aese-server/main.go` | M7 run orchestration HTTP 服务启动入口 |
| M7 API 实现 | `internal/httpapi/server.go` | 场景运行编排 API 路由与处理器 |
| MVP 范围 | `docs/AESE_MVP_Blueprint.md` | 产品和业务边界 |
| 华辰蓝图 | `docs/HCTM_Virtual_Enterprise_Blueprint.md` | 集团、工厂、产线、产品 |
| 主数据合同 | `docs/HCTM_Master_Data_Model.md` | 28 个对象 |
| 事件合同 | `docs/HCTM_Event_Model.md` | 18 类事件 |
| Seed 规格 | `docs/HCTM_Seed_Data_Plan.md` | 数据清单和 22 步时间线 |
| 演示验收 | `docs/HCTM_Demo_Story_01_Order_Expedite.md` | 七幕演示 runbook |
| M7 运行控制台 runbook | `docs/runbooks/hctm-m7-governed-scenario-operations-console.md` | 已完成的受治理场景控制台验收与对账入口 |
| AESE 2.0 设计输入 | `docs/ChatGPT20260722-aese2.0.md` | 原始构思，仅作为输入；工程边界以 ADR-004/DES-007 为准 |
| AESE 2.0 基础设计 | `docs/designs/DES-007-aese-2-foundation.md` | 三态、离散事件、IAOS 桥与 Genesis 迁移架构 |
| World/IAOS 桥接合同 | `docs/designs/DES-008-world-iaos-bridge-contract.md` | observation/intent/committed outcome、journal/cursor、权限与失败恢复 |
| M8 completed plan | `docs/plans/2026-07-22-m8-aese-2-foundation.md` | 已完成的三态世界、内核、bridge、pack 和 World Play 基线 |
| M9 completed plan | `docs/plans/2026-07-22-m9-genesis-incorporation.md` | 已完成的企业成立、治理、资本和预算基线 |
| M10 Agent 辅助工厂建设设计 | `docs/designs/DES-011-genesis-plant-build.md`、`docs/plans/2026-07-22-m10-genesis-plant-build.md`、`docs/runbooks/genesis-plant-build.md`、`docs/reports/m10-genesis-plant-build-evidence.md` | 场址/设施 reference replay 已完成；正式运行改由用户参数 + Agent Proposal + Human Review + World Observation 驱动，并提供只读事实评分，历史三个候选及金额仅为 fixture |
| M10 交互修订 active plan | `docs/plans/2026-08-01-m10-interactive-agent-plant-build.md`、`internal/plantbuild/interactive.go`、`internal/creative/minimax_provider.go`、`internal/httpapi/server.go`、`internal/iaosclient/client.go`、`frontend/src/world/plantBuild.ts`、`frontend/src/world/plantGame.ts`、`frontend/src/components/world/PlantBuildPlay.tsx`、`frontend/src/components/world/PlantBuildGameScene.tsx`；IAOS `api/plant_interactive.go`、`plant_contract_award.go`、`plant_construction_milestone.go`、`plant_process_run.go` 与 `incorporation/plant_entities.go` | 已落地需求/选址、场址控制、项目/WBS、合同授予，以及施工启动 → World 进度质量事实 → 人工里程碑验收；MiniMax 截断识别见 SOL-014，完整性与治理修复预算分离见 SOL-015。延期/缺陷/变更、工程财务和最终验收仍待实现 |
| M11 completed plan | `docs/plans/2026-07-22-m11-genesis-production-capability-build.md` | C0-C6 跨仓交付与验收记录 |
| M12 completed plan | `docs/plans/2026-07-22-m12-genesis-product-industrialization.md` | D0-D7 跨仓交付与验收记录 |
| M13 completed plan | `docs/plans/2026-07-22-m13-genesis-first-commercial-delivery.md` | E0-E8 跨仓交付与 Genesis 收口记录 |
| M14 completed plan | `docs/plans/2026-07-22-m14-parameterized-branch-experiments.md` | X0-X7 实验合同、分支、执行、证据、治理与 Scenario Lab |
| M15 completed plan | `docs/plans/2026-07-22-m15-governed-strategy-release-and-pilot.md` | R0-R7 证据审议、发布、shadow、pilot、回滚、采纳与 Control Room |
| M16 completed plan | `docs/plans/2026-07-22-m16-continuous-strategy-assurance-and-calibration.md` | A0-A7 dataset、drift、校准、holdout、再实验、治理与 Observatory |
| AESE 3.0 program | `docs/designs/DES-018-aese-3-completion-program.md` | M17-M24 顺序、terminal、依赖和 program 完成边界 |
| M17-M24 completed plans | `docs/plans/2026-07-22-m1*.md`、`docs/plans/2026-07-22-m2*.md` | 八个顺序 terminal 与最终 platform ready |

## 3. M3 实现路径

| 任务 | 路径 |
| --- | --- |
| 场景包 manifest 和数据 | `scenario-packs/hctm/` |
| JSON Schema | `scenario-packs/hctm/schemas/` |
| CLI 入口 | `cmd/aese/` |
| 场景包加载 | `internal/scenariopack/` |
| 离线校验 | `internal/validate/` |
| IAOS API client | `internal/iaosclient/` |
| 事件重放 | `internal/replay/` |
| 单元测试和破损 fixture | 对应 package `_test.go` |

核心实现入口：

| 能力 | 文件/目录 |
| --- | --- |
| CLI 命令分发 | `cmd/aese/main.go` |
| pack 合同、加载与 inspect | `internal/scenariopack/` |
| 结构、引用、时间线与经营不变量 | `internal/validate/` |
| IAOS 认证、schema、upsert、decompose 与 simulation ingress | `internal/iaosclient/` |
| HCTM 到 IAOS DES-047 wire 投影 | `internal/legacyprojection/` |
| dry-run/apply/replay/verify 协调 | `internal/replay/` |
| HCTM machine-readable pack | `scenario-packs/hctm/` |
| IAOS 兼容性证据 | `docs/reports/hctm-iaos-compatibility.md` |
| M3 端到端执行证据 | `docs/reports/hctm-m3-execution-evidence.md` |
| M3 操作手册 | `docs/runbooks/hctm-m3-local-run.md` |

## 4. IAOS 集成地图

AESE 不直接修改下列文件；需要集成时在独立 IAOS worktree 中按 IAOS 规则处理。

| 需求 | IAOS 文件/区域 | 当前事实 |
| --- | --- | --- |
| 统一事件 envelope/constants | `/iaos/iaos-go/shared/eventdef/events.go` | 已有基础 Event 和 O2D 常量 |
| O2D 服务入口 | `/iaos/iaos-go/scenarios/o2d/cmd/o2d/main.go` | 当前订阅 `iaos.*.o2d.order.confirmed` |
| BOM 展开/库存/工单 | `/iaos/iaos-go/scenarios/o2d/internal/mrp/` | decimal BOM 展开和原子幂等 workflow 已实证 |
| 动态实体 schema API | `/iaos/iaos-go/platform/internal/api/router.go` | `GET/POST /api/v1/metadata/schema/:entity` |
| 动态实体 CRUD/import | `/iaos/iaos-go/platform/internal/api/router.go`、`router_entity_*` | `/api/v1/entities/:entity` 和 import 路由 |
| 订单分解入口 | `/iaos/iaos-go/platform/internal/api/router.go` | `POST /api/v1/entities/sales_order/:id/decompose`；commit `0260f28` 增加状态 CAS/no-op 与 trace metadata |
| 场景 apply/reset | `/iaos/iaos-go/platform/internal/api/scenario.go` | `POST /api/v1/scenarios/apply|reset`；M3 allowlist、原子事务、自然键幂等、服务端 UUID resolve |
| 异常事件入口 | `/iaos/iaos-go/platform/internal/api/simulation.go` | `POST /api/v1/simulation/events`；支持设备停机、供应延期和来检失败 |
| M8 World Bridge | `/iaos/iaos-go-m8-world-bridge/platform/internal/api/world_bridge.go` | 已实现 observation ingress、tenant journal、cursor/SSE 和 intent/outcome；M9 IAOS 开发须另建独立 worktree，不直接复用旧工作树提交 |
| O2D workflow 幂等 | `/iaos/iaos-go/platform/pkg/workflow/` | `workflow_run` 去重，DAG/库存/工单/节点 Outbox 单一事务 |
| Outbox 注册 | `/iaos/iaos-go/platform/internal/capability/generic_atomic.go` | `RegisterOutboxMessage` |
| Capability 执行 | `/iaos/iaos-go/platform/internal/capability/` | 受治理业务动作入口 |
| AI Tool 调用 | `/iaos/iaos-go/platform/internal/aitool/` | Agent 安全调用入口 |
| AI Tool entity query | `/iaos/iaos-go/platform/internal/aitool/dispatcher_entity_records.go` | `source_ref=entity.records`；服务端 metadata 固定 entity/fields/filter/order/limit，调用 input 只给值；显式 tenant predicate + RLS |
| 前端业务入口 | `/iaos/iaos-go/frontend/src/app/page.tsx` | IAOS 主工作台 |

## 5. M3V 计划路径

| 任务 | 路径 |
| --- | --- |
| 前端工程与依赖 | `frontend/` |
| 应用壳和根路由 | `frontend/src/App.tsx`；根路径进入 Enterprise Genesis Home，`#sandbox` 保留旧样板沙盘 |
| Enterprise Genesis 产品主页 | `frontend/src/components/game/GenesisHome.tsx`、`GenesisHome.css`；创建企业、可访问且不重叠的 AI 创意官入口、样板世界、M9 旅程和开发状态入口 |
| Genesis Workspace onboarding | `frontend/src/components/game/GenesisOnboarding.tsx`、`frontend/src/game/api.ts`；五步向导、草稿自动保存、生产 Player session 转发，先创建隔离空间再进入身份工作室 |
| Genesis Player 登录、注册与企业大厅 | `frontend/src/components/game/GenesisLogin.tsx`、`GenesisCompanyLobby.tsx`、`frontend/src/game/api.ts`；`POST /api/aese/v1/auth/{register,login}`、`GET /api/aese/v1/auth/session`；IAOS PlayerAccount 密码认证、短期 Player JWT、owner-scoped 企业列表与 Founder tenant session 恢复，浏览器用户名/Player ID 不参与授权 |
| Genesis Workspace BFF | `internal/genesisworkspace/{player_auth_client,control_plane_client}.go`、`internal/httpapi/server.go`、`/api/aese/v1/genesis/workspaces`、`POST /api/aese/v1/genesis/workspaces/:workspace/session`；先由 IAOS session profile 验证 Player JWT 并派生 owner subject，再调用 DES-062 控制面；旧 local adapter 仅在 `AESE_AUTH_MODE=local_dev` 且 loopback 监听时可用，默认 IAOS 模式缺少 Session 返回 401 |
| Enterprise Genesis AI 身份 | `internal/creative/{minimax_provider,provider,job_store}.go`、`cmd/aese-server/main.go`、`GET /api/aese/v1/game/creative/{status,jobs}`；配置 MiniMax-M3 真实生成，持久记录 request/model/host/token/latency/validation/fallback evidence；deterministic provider 仅作明确标记的离线 fallback |
| AESE 后端部署与共享 MiniMax Provider 配置 | `scripts/deploy_aese_server.sh`、`.env.example`；安全加载权限为 0600 的本机 `.env`，校验 MiniMax 配置完整性，重建/重启 8090，并同时验证 M9 Creative 与 M10 Plant Planning Provider 状态；任一未连接都使发布失败 |
| Enterprise Genesis RPG 主线 | `frontend/src/components/game/FounderOfficeRPG.tsx`、`FounderOfficeCanvas.tsx`、`MissionBriefing.tsx`、`RPGEventIntro.tsx`、`LocationScene.tsx`、`WorkItemActionPanel.tsx`、`frontend/src/game/iaosLinks.ts`；23 工作项专属剧情、FounderProfile、城市/室内玩家移动、旅行转场、生成式 NPC 精灵、可检查物件、音效、奖励和 committed trophy；五个财务节点路由企业总部，财务中心穿透 IAOS 组织/Mandate/SoD/待办、账务与报表 |
| 设立与企业主数据 | IAOS `platform/internal/api/incorporation.go`；`incorporation_case` 显式拟设资料列，登记 committed 后物化 `entity_projection_legal_entity`；IAOS SOL-036 / DES-064 / AESE SOL-005 |
| 2D 画布与沙盘组件 | `frontend/src/components/` |
| 播放 reducer 和 Hook | `frontend/src/playback/` |
| 视图模型、校验和静态数据源 | `frontend/src/scenario/` |
| 场景预览数据 | `scenario-packs/hctm/stories/order-expedite-01/preview.json` |
| 单元与组件测试 | `frontend/src/**/*.test.ts(x)` |
| 浏览器验收 | `frontend/e2e/`；根主页回归见 `genesis-home.spec.ts` |
| 固定视口截图 | `frontend/test-results/*-completed.png` |
| 启动和操作手册 | `docs/runbooks/hctm-m3v-2d-sandbox.md` |
| M3V 验收证据 | `docs/reports/hctm-m3v-2d-sandbox-evidence.md` |

## 6. M4 实现路径

| 能力 | 路径 |
| --- | --- |
| AESE simulation request/response 合同 | `internal/iaosclient/client.go` |
| canonical 事件到受治理入口 | `internal/replay/replay.go` |
| M4 采购单与待检验单 DES-047 投影 | `internal/legacyprojection/projection.go`、`scenario-packs/hctm/stories/order-expedite-01/initial-state.json` |
| M4 completed plan | `docs/plans/2026-07-19-m4-governed-simulation-ingress.md` |
| 三类异常验收证据 | `docs/reports/hctm-m4-simulation-ingress-evidence.md` |
| IAOS 入口实现 | `/iaos/iaos-go/platform/internal/api/simulation.go` |

## 7. M5 Agent tracer 实现路径

| 能力 | 路径 |
| --- | --- |
| `agent-setup` / `agent-run` 命令分发 | `cmd/aese/main.go` |
| tool bundle 加载与约束 | `internal/agenttrace/config.go` |
| metadata/tool 创建、更新和启用 | `internal/agenttrace/setup.go` |
| 9 次受审计读取与三 Agent 建议构建 | `internal/agenttrace/run.go` |
| IAOS metadata / AI Tool client 合同 | `internal/iaosclient/client.go` |
| 版本化 HCTM tool / UI metadata manifest | `scenario-packs/hctm/agent-tools.json`；`sales_order.lines` 声明 `sales_order_line` child-list，reference 字段携带 `ref_entity` 供 IAOS 列表与详情解析业务标签 |
| tool bundle JSON Schema | `scenario-packs/hctm/schemas/agent-tools.schema.json` |
| Agent tracer 单元测试 | `internal/agenttrace/run_test.go` |
| IAOS `entity.records` dispatcher | `/iaos/iaos-go/platform/internal/aitool/dispatcher_entity_records.go` |

当前经营分析边界：M6 已补齐完工入库和发运在线事实，11,700 实发与 300 缺口可由 IAOS 证明；成本实际仍无批准基线，因此 `business_analysis` 只在 `cost_actuals` 维度保留 `partial`。

## 8. M6 计划路径

| 能力 | 计划路径 |
| --- | --- |
| 在线沙盘架构 | `docs/designs/DES-004-online-2d-enterprise-sandbox.md` |
| M6 completed plan | `docs/plans/2026-07-20-m6-online-2d-enterprise-sandbox.md` |
| 前端 Live 类型和 adapter | `frontend/src/scenario/` |
| 前端连接/恢复状态 | `frontend/src/LiveSandbox.tsx`、`frontend/src/scenario/iaosDataSource.ts` |
| Preview/Live 应用编排 | `frontend/src/App.tsx` |
| 可视化 IAOS 联动中心 | `frontend/src/components/IntegrationConsole.tsx`、`frontend/src/integration/iaosIntegration.ts`；一键取得 HCTM 本地演示身份，检查 profile/snapshot/events 与销售订单、工单、库存、设备，持久化连接配置并跳入 Live |
| AESE 完工/发运 replay client | `internal/iaosclient/`、`internal/replay/` |
| IAOS 场景业务动作 | `/iaos/iaos-go/platform/internal/api/`、`/iaos/iaos-go/platform/internal/capability/` |
| IAOS snapshot/cursor/SSE | `/iaos/iaos-go/platform/internal/api/` |
| M6 browser E2E | `frontend/e2e/sandbox.spec.ts`、`frontend/test-results/actual-live-*.png` |

## 9. M7 计划路径

| 能力 | 计划路径 |
| --- | --- |
| 无状态编排 API 决策 | `docs/decisions/ADR-003-thin-scenario-orchestration-api.md` |
| M7 设计与 completed plan | `docs/designs/DES-005-governed-scenario-operations-console.md`、`docs/plans/2026-07-20-m7-governed-scenario-operations-console.md` |
| CLI/application service 拆分 | `cmd/aese/`、`internal/application/` |
| HTTP server 与 handlers | `cmd/aese-server/main.go`、`internal/httpapi/server.go` |
| 前端场景运行状态 | `frontend/src/integration/iaosIntegration.ts`、`frontend/src/components/IntegrationConsole.tsx` |
| 联动中心运行视图 | `frontend/src/components/IntegrationConsole.tsx` |
| M7 运行控制台 runbook | `docs/runbooks/hctm-m7-governed-scenario-operations-console.md` |
| IAOS scenario run/permission | `/iaos/iaos-go/platform/internal/api/scenario*.go` |
| M7 browser E2E | `frontend/e2e/` |
| M7 O4 证据采集 | `scripts/m7-runbook-evidence-collect.sh`、`artifacts/m7-acceptance/` |

## 10. System Atlas 视图

| 能力 | 路径 |
| --- | --- |
| AESE 仿真完成体图谱 | `frontend/src/components/SystemAtlas.tsx`、`SystemAtlas.css`；Dagre 自动布局、拖动、关系高亮、Markdown 阅读器与功能入口 |
| AESE Atlas 深链接 | `frontend/src/App.tsx`；`#sandbox`、`#live`、`#integration`、`#atlas` 到真实界面状态 |
| IAOS 数据合同 | `/iaos/iaos-go/platform/internal/systematlas/`、`/api/v1/system-atlas` |
| 双系统视图设计 | `docs/designs/DES-006-system-atlas-aese-projection.md`、IAOS DES-049 |
| 进展登记 | `scripts/record_system_atlas_update.sh` |
| 声明式进展与 CI | `atlas-updates/`、`scripts/check_system_atlas_tracking.sh`、`scripts/sync_system_atlas_updates.sh`、`.github/workflows/system-atlas-governance.yml` |

Atlas 批量同步按声明 `occurred_at` 排序；声明 update key 不可变，状态校正必须新增声明，不能修改已登记 payload。部署单一新声明时通过 `ATLAS_UPDATE_FILE` 精确同步，避免用新的 commit metadata 重放历史声明。

## 11. M8 拟议实现路径

F0 合同入口已创建；F1-F5 目标路径仍须在实现时更新为实际入口：

| 能力 | 目标路径 |
| --- | --- |
| World JSON Schema 与 fixture | `world-contracts/schemas/`、`world-contracts/fixtures/` |
| Go 合同、strict parser、canonical hash | `internal/worldcontract/` |
| World Store 连接边界 | `internal/worldstore/` |
| PostgreSQL compose 与迁移 | `deploy/world-postgres/` |
| World Store runbook | `docs/runbooks/aese-world-store.md` |
| Genesis pack | `world-packs/hctm-genesis/` |
| 三态 tracer / Knowledge | `internal/genesis/`、`internal/knowledge/` |
| IAOS bridge adapter | `internal/bridge/iaos/` |
| M7 World Event adapter | `internal/legacyprojection/` |
| World Play UI | `frontend/src/components/world/`、`frontend/src/world/` |
| 企业生命周期首页与全程导航 | `frontend/src/components/world/WorldLifecycleHub.tsx`、`WorldLifecycleHub.css`、`#world`；M8 tracer 独立为 `#world-tristate` |
| M9 IAOS 原生真实闭环设计 | `docs/designs/DES-027-m9-iaos-native-incorporation-closed-loop.md`；D19 通用资产可发现性、D20 逐步骤追踪、D21 可解释合同、D22 持久工作项与逐节点参与、D23 Agent Runtime/工具/运行证据及外部模型边界、D24 配置驱动审批路由/工作分发/通知、D25 Entity 生命周期与正式审批单一语义、D26 页面功能说明与配置依赖图及财务开业、D27 全量 Capability 执行边界与 23 项显式正常路径；通用实现见 IAOS DES-053/DES-060/DES-065 |
| M9 IAOS 原生真实闭环 active remediation plan | `docs/plans/2026-07-23-m9-iaos-native-incorporation-closed-loop.md`；原完成结论因通用注册中心/工作区缺口撤回，当前以 11/20/5/8/8/5/10 数量、records API、十工作区和重复安装 no-op 为关闭门 |
| M9 游戏化企业创生体验设计 | `docs/designs/DES-028-enterprise-genesis-game-experience.md`；AI 企业身份工作室、2.5D 世界、人工/Agent 协作、GameProjection、CreativeJob 与 GX0–GX5 |
| Genesis 零起点与真实 AI 设计 | `docs/decisions/ADR-006-genesis-workspace-precedes-enterprise.md`、`docs/designs/DES-029-genesis-zero-start-and-ai-identity.md`、`docs/plans/2026-07-28-genesis-zero-start-and-ai-identity.md`；根主页 → PlayerAccount → GenesisWorkspace → IAOS tenant provisioning → AESE World Run → MiniMax naming → incorporation case；生产控制面已实现，剩余全 23 节点新租户验收 |
| 制造企业财务运行体系 | `docs/designs/finance/README.md`（DES-030–036 模块索引）、`docs/plans/2026-07-28-m9-manufacturing-finance-foundation.md`、`scenario-packs/hctm/finance-governance-baseline.json`、`internal/financebaseline/`、`docs/reports/m9-m13-finance-object-inventory.md`；IAOS `platform/internal/api/{finance_opening,finance_foundation,finance_ledger_foundation,finance_master_data_foundation}.go`、`platform/internal/incorporation/{platform_assets,finance_governance}.go`、`frontend/src/components/finance/{FinanceWorkspace,FinanceOperations,FinanceFoundation,FinanceLedgerFoundation,FinanceMasterDataFoundation}.tsx`、`frontend/src/components/finance/financeNavigation.ts`、DES-063/DES-067–070/DES-075；AESE `internal/iaosclient/incorporation.go`、`internal/gameprojection/iaos.go`、`frontend/src/components/game/{LocationScene,EnterpriseGenesisGame,WorkItemActionPanel}.tsx`、`frontend/src/game/iaosLinks.ts` | DES-030 只保留总览；F5B–F5D2 已实现多组织/Data Set、账簿/共享科目/日历、集团唯一伙伴/产品、16 类 Semantic/Entity 发布、同事务投影及七个职责单一财务入口；Runtime 2.11.0 / Edition 1.3.0 以 `finance.journal.entry.post` 唯一拥有凭证主子存储，资本来源能力只负责委托；AESE baseline 1.3 固定 HCTM 模板并校验引用、金额字符串和 12 期连续性；F5E 模块期间留待关账阶段 |
| Enterprise Genesis active subplan | `docs/plans/2026-07-27-enterprise-genesis-game-experience.md`；M9N 下的并行 AESE owner，GX0–GX5 / GXT1–GXT46 |
| GameProjection 合同/API | `internal/gameprojection/`、`internal/iaosclient/incorporation.go`、`world-contracts/schemas/game-projection.schema.json`、`world-contracts/fixtures/game-projection.json`、`GET /api/aese/v1/game/incorporation/:case/projection?frame=`；配置 `aese-server --iaos-base-url` 后读取 IAOS verified evidence bundle、Agent Run output、23 个持久工作项（含 5 个财务开业节点）及 G1–G7 可审阅事项，离线时才使用确定性 trace |
| Enterprise Genesis 游戏入口 | `frontend/src/components/game/EnterpriseGenesisGame.tsx`、`LocationScene.tsx`、`BrandStudio.tsx`、`WorkItemActionPanel.tsx`、`frontend/src/game/`、`#enterprise-genesis?tenant=&case=`；统一城市热点进入四类室内地点，主线通过建筑/NPC/资产推进，Work Item 降级到治理档案，并保留 G1–G7、登记/开户补正和虚构证照资产 |
| Enterprise Genesis 素材与验收 | `frontend/public/assets/enterprise-genesis/manifest.json`、`locations/`、`sprites/`、`frontend/e2e/enterprise-genesis.spec.ts`、`enterprise-genesis-live.spec.ts`、`enterprise-genesis-interactive.spec.ts`、`docs/reports/genesis-rpg-and-multi-tenant-acceptance.md`、`docs/runbooks/enterprise-genesis-game.md`；场景/透明精灵 hash 与许可、三视口、Founder 多租户选择、cross-read 隔离和失败恢复 |
| M9N 冻结合同、审计与风险 | `docs/contracts/m9-native-incorporation-contract.json`、`docs/reports/m9-native-asset-audit.json`、`docs/reports/m9-native-risk-register.json`；独立 IAOS worktree `/iaos/iaos-go-m9-native` |
| M9N Bridge reconciliation | `internal/bridge/iaos/reconcile.go`、`aese reconcile <bridge-journal.json>` | 对持久 journal 按 message/correlation/hash 报告 missing、duplicate、lagging、hash_mismatch、terminal_conflict；稳定排序、只读、可离线复验 |
| M9N IAOS lifecycle 联动页面 | `frontend/src/world/incorporation.ts`、`frontend/src/world/incorporationStepTrace.ts`、`frontend/src/game/api.ts`、`frontend/src/components/world/IncorporationPlay.tsx`；`#world-incorporation?tenant=&case=&process_run=&world_run=&correlation=`；[SOL-009](solutions/SOL-009-genesis-projection-session-and-runtime-recovery.md) | 按 D20 映射 8 个 World frame/15 次 transition；按 D22 默认只解锁 IAOS 已提交阶段；projection 401/404 和明确 Runtime stale 422 可刷新 Genesis managed session 后安全重试一次，其他业务拒绝不重试 |
| M9N 平台基础语义消费 | IAOS `docs/designs/DES-071-semantic-governance-and-foundation-seed.md`、`platform/internal/semantic/foundation/foundation.v1.json`、Runtime 2.9.0、`docs/solutions/SOL-044-entity-semantic-inheritance-and-governed-record-browser.md`、`docs/solutions/SOL-046-incorporation-case-entity-field-parity.md`、`docs/solutions/SOL-047-native-entity-projection-truthfulness-and-details.md`；AESE `docs/plans/2026-07-23-m9-iaos-native-incorporation-closed-loop.md` §20/§22/§24–§26 | AESE/M9 只引用平台 Concept/Archetype；场景 seed 不得覆盖 foundation-owned 资产；新增前先执行目录检索并取得语义变更批准；权威表未分类业务列必须同时进入 Entity 编译合同和显式投影，否则平台包升级失败关闭 |
| M9N Effective Runtime Artifact 执行权威 | IAOS `docs/designs/DES-072-effective-runtime-artifact-authority.md`、`platform/internal/effectiveruntime/`、`docs/runbooks/effective-runtime-artifact.md`；AESE DES-027 D28 / PLAN-M9-NATIVE-001 T24A | Entity/Capability/Process/Agent 正式路径只消费已编译 artifact；Process 冻结 Capability version/hash，任何缺失、stale、compiler/hash mismatch 均失败关闭 |
| M9N 版本化平台基础包与租户 Edition | IAOS `docs/designs/DES-073-platform-baseline-package-and-tenant-editions.md`、`platform/internal/platformpkg/`、`platform/internal/incorporation/platform_package.go`、`docs/runbooks/platform-baseline-packages.md`；AESE DES-027 D29 / PLAN-M9-NATIVE-001 §21–§26 | `genesis-m9@1.3.0` 组合 enterprise-governance、finance-foundation、genesis-incorporation；tenant-001 参考安装、Genesis 新租户和历史租户升级消费同一清单，严禁克隆租户业务数据 |
| M9N 原生 Entity 投影真实性与明细 | IAOS `platform/internal/incorporation/platform_assets.go`、`platform/internal/api/{incorporation,incorporation_work_items}.go`、`docs/solutions/SOL-047-native-entity-projection-truthfulness-and-details.md`；AESE PLAN-M9-NATIVE-001 §25 | 逻辑 code 与 `entity_projection_<code>` 物理读模型一一对应；只有拥有事实的 Agent/Capability/权威表可物化对应 Entity，未发生节点不得生成记录；列表和详情只读取编译后显式列，已有记录必填列为空时 Edition 升级失败关闭 |
| M9N Entity 存储与唯一写入权威 | IAOS `docs/designs/DES-074-entity-storage-and-write-authority-contract.md`、`platform/internal/entitystorage/`、`platform/internal/compiler/entity.go`、`platform/internal/api/router_entity_{schema,crud}.go`；AESE PLAN-M9-NATIVE-001 §26 | 新动态 Entity 使用 `entity_record_<code>`；复杂领域表和纯计算结果分别投影到 `entity_projection_<code>`；每个 Entity 只能声明一个 `storage_write_owner`，领域/计算投影的通用 CRUD 失败关闭 |
| M9N 可执行原子能力与授权合同 | IAOS `docs/designs/DES-076-atomic-capability-library-and-business-authoring.md`、`docs/solutions/SOL-050-capability-lifecycle-permission-and-menu-semantics.md`、`docs/references/Atomic_ability.md`、`platform/internal/capability/atomic_catalog.go`；AESE PLAN-M9-NATIVE-001 §18 | AESE 不维护 Handler；Agent/Process/Bridge 只调用已发布、已绑定且角色获 `capability.<code>/EXECUTE` 的 Business Capability，由 IAOS 19 项 typed 原子能力在外层治理事务内执行；草稿不能进入剧情执行，完整业务动作不得降级成原子能力。 |
| M9N 可调预算与草案审批权威 | IAOS `platform/internal/api/incorporation_budget_policy.go`、`platform/internal/api/incorporation_work_items.go`、`platform/internal/incorporation/capability_contracts.go`、`docs/solutions/SOL-056-m9-budget-proposal-agent-limit.md`；AESE DES-010 / PLAN-M9-NATIVE-001 T36A | `initial.budget.prepare` 是 proposal，不消耗 Agent 交易限额；租户配置最低金额、可选绝对上限和已核验资金比例；节点 21 保存 revision/hash，G7 只批准同一草案，World fixture 金额不得进入正式交互默认值。 |
| M9N World 人工 Observation | `frontend/src/world/incorporation.ts::submitIncorporationObservation`、`IncorporationPlay.tsx`、`docs/solutions/SOL-002-incorporation-observation-click-idempotency.md` | 登记、开户、任命三个外部参与者按钮通过受治理 `/api/v1/world-bridge/observations` 写 Observation；同一案件/类型/结果使用稳定 transport identity；不直接推进 IAOS，用户仍需执行对应 world_wait 工作项 |
| M9N 局域网加载与 SSE 排障 | `docs/solutions/SOL-001-m9-lan-lifecycle-loading-and-sse.md` | 最近 case 自动发现、浏览器 hostname 双向链接、60 秒 SSE write deadline 续期与 BFCache/HMR 边界 |
| M9N 失效案件深链恢复 | `frontend/src/world/incorporation.ts`、`frontend/src/components/world/IncorporationPlay.tsx`、`docs/solutions/SOL-002-aese-world-stale-incorporation-case-recovery.md` | trace 404 保留 AESE World 基线并查询 recent cases；401/403/5xx 仍失败关闭 |
| World Play API | `internal/httpapi/server.go`（`/api/aese/v1/world/genesis`） |
| 验收 runbook / 能力缺口 | `docs/runbooks/aese-world-play.md`、`docs/capability-gap-ledger.md` |
| Atlas planned 投影 | `atlas/system-atlas-planned.json` |
| World CLI（F1） | `cmd/aese/world.go`；`aese world validate|inspect|run|replay` |
| 世界 run、状态投影、快照恢复 | `internal/world/` |
| 虚拟时钟与推进 | `internal/simtime/` |
| 世界事件稳定队列 | `internal/simevent/` |
| 版本化纯函数规则 | `internal/rules/` |
| F1 运行样例 | `world-contracts/runtime-example/` |
| 角色认知与差异 | `internal/knowledge/` |
| IAOS 双向桥 | `internal/bridge/iaos/` |
| Genesis 世界包 | `world-packs/hctm-genesis/` |
| World 前端 | `frontend/src/world/`、`frontend/src/components/world/` |

## 12. M9 已实现路径

M9 I0-I5 的当前入口如下：

| 能力 | 路径 |
| --- | --- |
| M9 设计与 completed plan | `docs/designs/DES-010-genesis-incorporation-and-governance.md`、`docs/plans/2026-07-22-m9-genesis-incorporation.md` |
| 成立领域模型与规则 | `internal/incorporation/` |
| 成立机器合同 | `world-contracts/schemas/incorporation-campaign.schema.json`、对应 fixture |
| Incorporation campaign | `world-packs/hctm-genesis/campaigns/incorporation/` |
| IAOS bridge adapter | `internal/bridge/iaos/` |
| World Play 成立视图 | `frontend/src/components/world/`、`frontend/src/world/` |
| IAOS 法人/治理能力 | `/iaos/iaos-go-m9-genesis/platform/internal/api/genesis_governance.go`；DES-051；revision `edcb915` |
| M9 runbook / evidence | `docs/runbooks/genesis-incorporation.md`、`docs/reports/m9-genesis-incorporation-evidence.md` |

## 13. M10 reference 与交互纵切路径

| 能力 | 路径 |
| --- | --- |
| M10 设计与计划 | `docs/designs/DES-011-genesis-plant-build.md`、`docs/plans/2026-08-01-m10-interactive-agent-plant-build.md`；旧 completed plan 只治理 reference replay |
| 交互 Requirement/Proposal/Review 合同 | `internal/plantbuild/interactive.go`、`world-contracts/schemas/plant-build-interactive.schema.json` |
| Agent provider 与技术证据 | `internal/creative/minimax_provider.go::CompleteJSON`、CreativeJob store；失败时不返回 fixture 候选 |
| M9→M10 阶段交接与全屏 Plant Build 游戏 | `frontend/src/components/game/EnterpriseGenesisGame.tsx` 在 IAOS M9 终态显示携带 tenant/case/workspace 的交接主动作；`frontend/src/App.tsx` 负责双向深链；`frontend/src/components/world/PlantBuildPlay.tsx` 组合 `PlantBuildGameScene.tsx/.css`、NPC 临时任务和按地点场景档案，`frontend/src/world/plantGame.ts` 从 Requirement/Proposal/Review/Investigation/Observation/Recommendation/Decision/SiteControl 编译总部、产业地图、现场、治理会议室、园区权利人和项目准备阶段；`frontend/src/world/plantBuild.ts`、`siteAssessment.ts`、组件/阶段/评分测试、`frontend/e2e/m9-m10-handoff.spec.ts` 与 `/#world-plant-build`。场址控制 POST/Observation 均经过 AESE 命名 Command Gateway，浏览器不能直接写 IAOS 或伪造外部交付 |
| AESE 定向 BFF | `internal/httpapi/server.go`：`planning-status`、`financial-constraints`、Requirement/ProposalSet 读取、`proposals` 写入、`reviews`、`investigations`、`observations`；`internal/iaosclient/client.go` 只暴露命名的 IAOS adapter，禁止通用代理 |
| IAOS M10 菜单与解释型工作台 | `/iaos/iaos-go/frontend/src/components/genesis/PlantPlanningWorkspace.tsx`、`menu.genesis_plant_planning`；由 `genesis-plant-planning` 包随平台基础包 `1.8.0` 安装，入口为 `业务智造层 → M10 工厂规划`，并深链 AESE `/#world-plant-build` |
| 场址外部调研 World 闭环 | `frontend/src/components/world/PlantBuildPlay.tsx`、`frontend/src/world/plantBuild.ts`、`internal/httpapi/server.go`、`internal/plantbuild/interactive.go`；IAOS `/iaos/iaos-go/platform/internal/api/plant_interactive.go`、`world_bridge.go` 与 `facility.site.investigation.v1` | AESE 定向 BFF 发起调研；IAOS 保存请求与持久 `waiting_world` 工作项并创建 Intent；外部角色结构化 Observation 先进入 World Journal，再由受治理能力提交权威事实并完成工作项 |
| 场址比较、修订、推荐、游戏内审批与正式选择 | `frontend/src/world/siteAssessment.ts` 按当前 ProposalSet ID/revision/proposal ID 筛选 Observation；`frontend/src/components/world/PlantBuildPlay.tsx` 展示硬门槛、版本修订、治理会议事项与受派人决策；`frontend/src/world/plantBuild.ts` 读取审批详情并提交 approve/reject；`internal/httpapi/server.go`、`internal/iaosclient/client.go` 提供定向只读审批 BFF 与受控 Command Gateway；IAOS `plant_site_selection.go`、`incorporation/plant_approval.go`、`facility.site.selection.v1` | IAOS 重算 `site-assessment-v1`、冻结 Business Subject/Flow/Assignment 并验证当前决定人；AESE 游戏内审批不指定处理人、不写状态。approved 后由 `site.selection.formalize` 写正式决定；IAOS 审批中心只作审计穿透 |
| IAOS 交互权威切片 | `/iaos/iaos-go/platform/internal/api/plant_interactive.go`、`plant_site_selection.go`、`plant_process_run.go`、`plant_project_baseline.go`、`plant_contract_award.go`：只读财务快照、Requirement/Proposal/Review、Investigation/World Observation、推荐/审批/正式选址、场地控制、项目/WBS、RFQ/可信投标/合同授予、RLS、Capability 写门、幂等、Audit/Outbox；尚未包含施工/变更/工程财务/验收全链 |
| 历史场址/空间/项目 tracer | `internal/plantbuild/`、`world-contracts/schemas/plant-build-campaign.schema.json`、`world-packs/hctm-genesis/campaigns/plant-build/`；只允许 fixture replay |
| 历史 IAOS 投资与项目治理 | `/iaos/iaos-go-m10-plant/platform/internal/api/plant_governance.go`；DES-052；revision `23be02a`，尚未接入新交互 Process |
| M10 runbook / reference evidence | `docs/runbooks/genesis-plant-build.md`、`docs/reports/m10-genesis-plant-build-evidence.md`；evidence 仍只证明 reference replay |

## 14. M11 已实现路径

| 能力 | 路径 |
| --- | --- |
| M11 设计与 completed plan | `docs/designs/DES-012-genesis-production-capability-build.md`、`docs/plans/2026-07-22-m11-genesis-production-capability-build.md` |
| 生产能力领域与 tracer | `internal/capabilitybuild/` |
| 设备/人员/资格机器合同 | `world-contracts/schemas/capability-build-campaign.schema.json`、对应 fixture |
| Capability Build campaign | `world-packs/hctm-genesis/campaigns/capability-build/` |
| IAOS bridge adapter | `internal/bridge/iaos/` |
| World Play 能力建设视图 | `frontend/src/components/world/CapabilityBuildPlay.tsx`、`frontend/src/world/capabilityBuild.ts`、`/#world-capability-build` |
| Capability Build API | `internal/httpapi/server.go`（`GET /api/aese/v1/world/capability-build`） |
| IAOS 资金/采购/资产/组织/资格治理 | `/iaos/iaos-go-m11-capability/platform/internal/api/plant_governance.go`；DES-053；revision `789b925` |
| M11 runbook/evidence | `docs/runbooks/genesis-capability-build.md`、`docs/reports/m11-genesis-production-capability-evidence.md` |

## 15. M12 已实现路径

| 能力 | 路径 |
| --- | --- |
| M12 设计与 completed plan | `docs/designs/DES-013-genesis-product-industrialization.md`、`docs/plans/2026-07-22-m12-genesis-product-industrialization.md` |
| 产品工业化领域与 tracer | `internal/industrialization/` |
| RFQ/APQP/试制/PPAP 机器合同 | `world-contracts/schemas/industrialization-campaign.schema.json`、对应 fixture |
| Industrialization campaign | `world-packs/hctm-genesis/campaigns/industrialization/` |
| 旧 HCTM compatibility | `scenario-packs/hctm/master-data/materials.json`、`scenario-packs/hctm/master-data/manufacturing.json` |
| IAOS bridge adapter | `internal/bridge/iaos/` |
| World Play 工业化视图 | `frontend/src/components/world/IndustrializationPlay.tsx`、`frontend/src/world/industrialization.ts`、`/#world-industrialization` |
| Industrialization API | `internal/httpapi/server.go`（`GET /api/aese/v1/world/industrialization`） |
| IAOS 客户项目/工程/APQP/质量/PPAP 治理 | `/iaos/iaos-go-m12-industrialization/platform/internal/api/plant_governance.go`；DES-054；revision `50a46e2` |
| M12 runbook/evidence | `docs/runbooks/genesis-industrialization.md`、`docs/reports/m12-genesis-industrialization-evidence.md` |

## 16. M13 已实现路径

| 能力 | 路径 |
| --- | --- |
| M13 设计与 completed plan | `docs/designs/DES-014-genesis-first-commercial-delivery.md`、`docs/plans/2026-07-22-m13-genesis-first-commercial-delivery.md` |
| 第一次商业交付领域与 tracer | `internal/firstdelivery/` |
| O2D/发票/现金/成本机器合同 | `world-contracts/schemas/first-delivery-campaign.schema.json`、对应 fixture |
| First Delivery campaign | `world-packs/hctm-genesis/campaigns/first-delivery/` |
| Genesis O2D compatibility | `scenario-packs/hctm/`、`internal/legacyprojection/`、`internal/replay/` |
| IAOS bridge adapter | `internal/bridge/iaos/` |
| World Play 首次交付视图 | `frontend/src/components/world/FirstDeliveryPlay.tsx`、`frontend/src/world/firstDelivery.ts`、`/#world-first-delivery` |
| First Delivery API | `internal/httpapi/server.go`（`GET /api/aese/v1/world/first-delivery`） |
| IAOS O2D/发票/应收/收款/实际成本治理 | `/iaos/iaos-go-m13-delivery/platform/internal/api/plant_governance.go`；DES-055；revision `067bbb4` |
| M13/Genesis evidence | `docs/reports/m13-genesis-first-commercial-delivery-evidence.md`、`docs/reports/project-genesis-m9-m13-e2e.md` |

## 17. M14 参数化实验实现路径

X0-X7 已完成，当前实现入口如下：

| 能力 | 路径 |
| --- | --- |
| M14 设计与 completed plan | `docs/designs/DES-015-parameterized-branch-experiments.md`、`docs/plans/2026-07-22-m14-parameterized-branch-experiments.md` |
| 实验 Schema 与 fixture | `world-contracts/schemas/experiment-*.schema.json`、`world-contracts/fixtures/experiments/` |
| 实验定义、矩阵、随机流与聚合 | `internal/experiment/` |
| checkpoint fork 与 branch/run catalog | `internal/world/`、`internal/worldstore/` |
| 多周期 Genesis 策略/外生规则 | `internal/genesisexperiment/`、`world-packs/hctm-genesis/experiments/` |
| Experiment CLI | `cmd/aese/experiment.go`；`aese experiment validate|inspect|expand|run|compare|evidence|replay` |
| IAOS bridge adapter | `internal/bridge/iaos/`；experiment/recommendation allowlist payload |
| Scenario Lab | `frontend/src/components/world/ScenarioLab.tsx`、`frontend/src/world/experiment.ts`、`/#world-experiments` |
| Experiment API | `internal/httpapi/server.go`（`/api/aese/v1/world/experiments`） |
| M14 runbook / evidence | `docs/runbooks/genesis-scenario-lab.md`、`docs/reports/m14-parameterized-experiment-evidence.md` |

## 18. M15 策略发布实现路径

R0-R7 已完成，当前实现入口如下：

| 能力 | 路径 |
| --- | --- |
| M15 设计与 completed plan | `docs/designs/DES-016-governed-strategy-release-and-pilot.md`、`docs/plans/2026-07-22-m15-governed-strategy-release-and-pilot.md` |
| Strategy Schema 与 fixture | `world-contracts/schemas/strategy-release.schema.json`、`world-contracts/fixtures/strategy-release.json` |
| Evidence review、release、shadow、guardrail | `internal/strategyrelease/` |
| IAOS 治理端点 | `/api/v1/genesis/strategy/actions`；exact evidence/release hash、expected version、幂等和独立审批 |
| Strategy Control Room | `frontend/src/components/world/StrategyControlRoom.tsx`、`frontend/src/world/strategyControl.ts`、`/#world-strategy-control` |
| Strategy API | `internal/httpapi/server.go`（`GET /api/aese/v1/world/strategy-control`） |
| M15 runbook / evidence | `docs/runbooks/genesis-strategy-control.md`、`docs/reports/m15-governed-strategy-release-evidence.md` |

## 19. M16 持续保障实现路径

A0-A7 已完成，当前实现入口如下：

| 能力 | 路径 |
| --- | --- |
| M16 设计与 completed plan | `docs/designs/DES-017-continuous-strategy-assurance-and-calibration.md`、`docs/plans/2026-07-22-m16-continuous-strategy-assurance-and-calibration.md` |
| Assurance Schema 与 fixture | `world-contracts/schemas/assurance-cycle.schema.json`、`world-contracts/fixtures/assurance-cycle.json` |
| Dataset、quality、drift、calibration 与 validation | `internal/assurance/` |
| Strategy Assurance Observatory | `frontend/src/components/world/AssuranceObservatory.tsx`、`frontend/src/world/assurance.ts`、`/#world-assurance` |
| Assurance API | `internal/httpapi/server.go`（`/api/aese/v1/world/strategy-assurance`） |
| M16 runbook / evidence | `docs/runbooks/genesis-strategy-assurance.md`、`docs/reports/m16-strategy-assurance-evidence.md` |

## 20. AESE 3.0 M17-M24 实现路径

B0-B7 实现后必须把“目标路径”更新为实际入口：

| 能力 | 目标路径 |
| --- | --- |
| M17-M24 设计与 completed plans | `docs/designs/DES-018-aese-3-completion-program.md` 至 DES-026、`docs/plans/2026-07-22-m17-*` 至 M24 |
| IBP Schema 与 fixture | `world-contracts/schemas/ibp-*.schema.json`、`world-contracts/fixtures/ibp/` |
| Planning cycle、bucket、reconciliation 与 scenario | `internal/ibp/` |
| IBP CLI | `cmd/aese/ibp.go`；`aese ibp validate|inspect|build|reconcile|evidence` |
| IAOS bridge adapter | `internal/bridge/iaos/`；planning review/decision/release allowlist payload |
| Executive IBP Room | `frontend/src/components/world/ExecutiveIBPRoom.tsx`、`frontend/src/world/ibp.ts`、`/#world-ibp` |
| IBP API | `internal/httpapi/server.go`（`/api/aese/v1/world/ibp`） |
| 严格 Go 类型、校验和 canonical hash | `internal/aese3/program.go`、`internal/aese3/program_test.go` |
| Schema、fixture 与 pack registry | `world-contracts/schemas/aese3-program.schema.json`、`world-contracts/fixtures/aese3-program.json`、`world-packs/hctm-genesis/aese3/` |
| API 与 Completion Room | `internal/httpapi/server.go`、`frontend/src/world/aese3.ts`、`frontend/src/components/world/AESE3CompletionRoom.tsx`、`#world-aese3` |
| runbook / evidence | `docs/runbooks/aese3-reference-release.md`、`docs/reports/aese3-m17-m24-completion-evidence.md` |
| IAOS 治理 | `/iaos/iaos-go-aese3/platform/internal/api/plant_governance.go`、IAOS DES-059 |

## 21. M18-M24 设计导航

| 里程碑 | 设计 |
| --- | --- |
| Program | `docs/designs/DES-018-aese-3-completion-program.md` |
| M18 | `docs/designs/DES-020-product-and-customer-portfolio-expansion.md` |
| M19 | `docs/designs/DES-021-multi-site-supply-and-fulfillment-network.md` |
| M20 | `docs/designs/DES-022-after-sales-warranty-and-closed-loop-quality.md` |
| M21 | `docs/designs/DES-023-plant-resource-and-ehs-resilience.md` |
| M22 | `docs/designs/DES-024-group-finance-treasury-and-investment.md` |
| M23 | `docs/designs/DES-025-governed-multi-agent-organization.md` |
| M24 | `docs/designs/DES-026-scenario-platform-productization.md` |

以上八个设计和计划均已完成；后续扩展必须另立 program，不得静默扩展 M24 reference certification。

## 22. 导航更新触发器

财务文档导航以 `docs/designs/finance/README.md` 为准。多组织财务、Data Set 共享、
科目表/法人扩展和 Business Partner 分层以 DES-031 为准；模块期间控制以 DES-035 为准；
其实现将跨 IAOS Organization、Metadata Entity、RLS、Capability、Approval 与财务运行时。

M9 账套/期间交互入口：`frontend/src/components/game/WorkItemActionPanel.tsx` 采集账套名称、
年度和 12 期日历，经 `frontend/src/game/api.ts` 传入 IAOS `accounting.book.activate`。

M9 受治理写入口：`frontend/src/game/api.ts` 与
`frontend/src/world/incorporation.ts` 只调用
`POST /api/aese/v1/commands/iaos/*`；白名单、身份/租户透传和错误映射位于
`internal/httpapi/server.go::handleIAOSCommandGateway`，窄 IAOS adapter 位于
`internal/iaosclient/client.go::PostGovernedCommand`。IAOS 端的 Capability/Process
Artifact 执行权威见 IAOS SOL-048 和 DES-072。

产品知识与用户手册：AESE 场景知识合同见
`docs/designs/DES-037-aese-scenario-knowledge-and-iaos-product-hub.md`，M9 人工可读入口为
`docs/manuals/m9-incorporation-user-manual.md`，实施计划为
`docs/plans/2026-08-01-product-knowledge-scenario-integration.md`。通用 Knowledge Registry、
知识中心、权限和 Copilot 检索由 IAOS DES-077 与 `/api/v1/knowledge/*` 提供；AESE 不新增
知识数据库。机器 Edition Schema/manifest 位于 `scenario-packs/hctm/knowledge/`；加载、稳定
编码引用、SHA-256 校验和签名 bundle 编译位于 `internal/scenarioknowledge/`；离线入口为
`aese knowledge validate|digest|compile`，受治理安装入口为
`aese knowledge install <bundle> --target <IAOS URL> [--apply]`。编译产物位于
`scenario-packs/hctm/knowledge/dist/`，IAOS 消费 bundle 而不读取 AESE 文件系统。M9 任务的“这一步是什么”深链位于
`frontend/src/components/game/WorkItemActionPanel.tsx`，链接构造位于
`frontend/src/game/iaosLinks.ts::m9KnowledgeUrl`；该函数从 WorldProjection、GameWorkItem 和
Genesis session 生成封闭 workspace/case/world run/node/actor/task/capability 导航上下文。
IAOS 知识中心与 Copilot BFF 负责二次归一化和可见展示。World/IAOS 双侧实际证据和漂移检测
由 IAOS `ScenarioEvidencePanel`、M9 Evidence Bundle/工作项 API 和 Copilot BFF 服务端重读实现；
World 证据特指 IAOS 已持久化的外部 Observation 回执，导航参数不得被解释为运行事实。

以下改动必须更新本文件：

- 新增命令、核心 package、场景包、schema 或脚本。
- 改变场景包目录结构。
- 改变 IAOS API、event subject 或 Capability 集成点。
- 新增前端主要页面或演示入口。
- 删除或替代本文件列出的任何入口。
