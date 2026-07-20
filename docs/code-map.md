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
| 开始 M6 在线沙盘 | DES-004、M6 active plan、`frontend/src/scenario/`、IAOS scenario API |
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
| 文档索引 | `docs/README.md` | 文档分类、状态和编号 |
| 项目上下文 | `docs/agent-project-context.md` | Agent 快速入门 |
| 架构 | `docs/architecture.md` | 仓库边界、数据流、安全和可重复性 |
| 路线图 | `docs/roadmap.md` | M0-M6 当前状态 |
| 代码导航 | `docs/code-map.md` | 本文件 |
| 历史记录 | `docs/progress-log.md` | 只追加进展日志 |
| MVP 范围 | `docs/AESE_MVP_Blueprint.md` | 产品和业务边界 |
| 华辰蓝图 | `docs/HCTM_Virtual_Enterprise_Blueprint.md` | 集团、工厂、产线、产品 |
| 主数据合同 | `docs/HCTM_Master_Data_Model.md` | 28 个对象 |
| 事件合同 | `docs/HCTM_Event_Model.md` | 18 类事件 |
| Seed 规格 | `docs/HCTM_Seed_Data_Plan.md` | 数据清单和 22 步时间线 |
| 演示验收 | `docs/HCTM_Demo_Story_01_Order_Expedite.md` | 七幕演示 runbook |

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
| 版本化 HCTM tool manifest | `scenario-packs/hctm/agent-tools.json` |
| tool bundle JSON Schema | `scenario-packs/hctm/schemas/agent-tools.schema.json` |
| Agent tracer 单元测试 | `internal/agenttrace/run_test.go` |
| IAOS `entity.records` dispatcher | `/iaos/iaos-go/platform/internal/aitool/dispatcher_entity_records.go` |

当前经营分析边界：在线工具可读订单、订单行、库存、BOM、采购、设备、检验和工单；完工入库、发运及成本实际仍没有受治理在线事实，所以 `business_analysis` 必须保留 `partial` / data gaps，不能复述 Preview 的 11,700/300 为在线结果。

## 8. M6 计划路径

| 能力 | 计划路径 |
| --- | --- |
| 在线沙盘架构 | `docs/designs/DES-004-online-2d-enterprise-sandbox.md` |
| M6 active plan | `docs/plans/2026-07-20-m6-online-2d-enterprise-sandbox.md` |
| 前端 Live 类型和 adapter | `frontend/src/scenario/` |
| 前端连接/恢复状态 | `frontend/src/live/` |
| Preview/Live 应用编排 | `frontend/src/App.tsx` |
| AESE 完工/发运 replay client | `internal/iaosclient/`、`internal/replay/` |
| IAOS 场景业务动作 | `/iaos/iaos-go/platform/internal/api/`、`/iaos/iaos-go/platform/internal/capability/` |
| IAOS snapshot/cursor/SSE | `/iaos/iaos-go/platform/internal/api/` |
| M6 browser E2E | `frontend/e2e/` |

## 9. 导航更新触发器

以下改动必须更新本文件：

- 新增命令、核心 package、场景包、schema 或脚本。
- 改变场景包目录结构。
- 改变 IAOS API、event subject 或 Capability 集成点。
- 新增前端主要页面或演示入口。
- 删除或替代本文件列出的任何入口。
