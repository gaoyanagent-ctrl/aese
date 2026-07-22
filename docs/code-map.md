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
| 设计或实现 AESE 2.0 World Runtime | ADR-004、DES-007、DES-009、M8 active plan；先读 `world-contracts/` 和 `internal/worldcontract/` |
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
| 路线图 | `docs/roadmap.md` | M0-M6 当前状态 |
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
| M8 active plan | `docs/plans/2026-07-22-m8-aese-2-foundation.md` | 当前唯一 active 计划与跨仓交付门 |

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
| M8 World Bridge（拟议） | `/iaos/iaos-go/platform/internal/api/` 新窄模块，实际文件待 IAOS 独立 worktree 确定 | DES-008 规定 observation ingress、tenant journal、cursor/SSE 和 Capability intent/outcome 原语；尚未实现 |
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
| 应用壳和响应式布局 | `frontend/src/App.tsx`、`frontend/src/styles/global.css` |
| 2D 画布与沙盘组件 | `frontend/src/components/` |
| 播放 reducer 和 Hook | `frontend/src/playback/` |
| 视图模型、校验和静态数据源 | `frontend/src/scenario/` |
| 场景预览数据 | `scenario-packs/hctm/stories/order-expedite-01/preview.json` |
| 单元与组件测试 | `frontend/src/**/*.test.ts(x)` |
| 浏览器验收 | `frontend/e2e/` |
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

Atlas 批量同步按声明 `occurred_at` 排序；声明 update key 不可变，状态校正必须新增声明，不能修改已登记 payload。

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

## 12. 导航更新触发器

以下改动必须更新本文件：

- 新增命令、核心 package、场景包、schema 或脚本。
- 改变场景包目录结构。
- 改变 IAOS API、event subject 或 Capability 集成点。
- 新增前端主要页面或演示入口。
- 删除或替代本文件列出的任何入口。
