# AESE Code Map

本文件把常见任务映射到应优先阅读和修改的文件。当前 AESE 尚处于文档到实现的过渡期，标记为“planned”的路径尚未创建。

## 1. 快速入口

| 任务 | 先读 |
| --- | --- |
| 理解项目定位 | `README.md`、`docs/agent-project-context.md` |
| 查看当前进度 | `docs/roadmap.md`、`docs/progress-log.md` |
| 理解 AESE/IAOS 边界 | `docs/architecture.md`、ADR-001 |
| 开始 M3 实现 | DES-001、M3 active plan |
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

## 3. M3 Planned 路径

| 任务 | Planned 路径 |
| --- | --- |
| 场景包 manifest 和数据 | `scenario-packs/hctm/` |
| JSON Schema | `scenario-packs/hctm/schemas/` |
| CLI 入口 | `cmd/aese/` |
| 场景包加载 | `internal/scenariopack/` |
| 离线校验 | `internal/validate/` |
| IAOS API client | `internal/iaosclient/` |
| 事件重放 | `internal/replay/` |
| 单元测试 fixture | `testdata/` 或对应 package `_test.go` |
| 文档一致性检查 | `scripts/check_project_docs.sh` |

## 4. IAOS 集成地图

AESE 不直接修改下列文件；需要集成时在独立 IAOS worktree 中按 IAOS 规则处理。

| 需求 | IAOS 文件/区域 | 当前事实 |
| --- | --- | --- |
| 统一事件 envelope/constants | `/iaos/iaos-go/shared/eventdef/events.go` | 已有基础 Event 和 O2D 常量 |
| O2D 服务入口 | `/iaos/iaos-go/scenarios/o2d/cmd/o2d/main.go` | 当前订阅 `iaos.*.o2d.order.confirmed` |
| BOM 展开/库存/工单 | `/iaos/iaos-go/scenarios/o2d/internal/mrp/` | 可复用现有 handler |
| 动态实体 schema API | `/iaos/iaos-go/platform/internal/api/router.go` | `GET/POST /api/v1/metadata/schema/:entity` |
| 动态实体 CRUD/import | `/iaos/iaos-go/platform/internal/api/router.go`、`router_entity_*` | `/api/v1/entities/:entity` 和 import 路由 |
| 订单分解入口 | `/iaos/iaos-go/platform/internal/api/router.go` | `POST /api/v1/entities/sales_order/:id/decompose` |
| Outbox 注册 | `/iaos/iaos-go/platform/internal/capability/generic_atomic.go` | `RegisterOutboxMessage` |
| Capability 执行 | `/iaos/iaos-go/platform/internal/capability/` | 受治理业务动作入口 |
| AI Tool 调用 | `/iaos/iaos-go/platform/internal/aitool/` | Agent 安全调用入口 |
| 前端业务入口 | `/iaos/iaos-go/frontend/src/app/page.tsx` | IAOS 主工作台 |

## 5. 导航更新触发器

以下改动必须更新本文件：

- 新增命令、核心 package、场景包、schema 或脚本。
- 改变场景包目录结构。
- 改变 IAOS API、event subject 或 Capability 集成点。
- 新增前端主要页面或演示入口。
- 删除或替代本文件列出的任何入口。
