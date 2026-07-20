# AESE Architecture

## 1. 系统定位

AESE 是 IAOS 的行业仿真与场景内容层。它描述企业世界、准备确定性数据、注入受控业务事件，并验证 IAOS 与 Agent 在异常经营场景中的行为。

AESE 不拥有 ERP/MES 运行时，也不复制 IAOS 的安全和治理能力。

## 2. 仓库职责边界

| 能力 | AESE 仓库 | IAOS 仓库 |
| --- | --- | --- |
| 虚拟企业蓝图 | Own | Consume |
| 场景 manifest、seed、事件 fixture | Own | Import/consume |
| JSON Schema 与离线校验 | Own | Contract-test |
| 演示时间线和期望输出 | Own | Execute/render |
| 数据库与 RLS | No | Own |
| Metadata / Dynamic Entity runtime | No | Own |
| Outbox + NATS | No | Own |
| Capability / Process / Decision | No | Own |
| Agent Tool Registry | No | Own |
| 正式业务 UI | Specification | Own |
| 只读 2D 场景预览器 | Own | Integrate later |

长期原则：

> AESE 版本化“企业场景内容”，IAOS 执行“企业业务能力”。

ADR-002 允许 AESE 拥有只读 2D 场景预览器，用于快速验证场景表达和回归。预览器不拥有业务写入、权限、流程和 Agent 运行时；在线模式的数据仍来自 IAOS 受治理 API 和事件流。

## 3. 目标运行链路

```text
HCTM scenario pack
  -> AESE validator
  -> dry-run impact report
  -> IAOS authenticated import/apply
  -> PostgreSQL RLS transaction
  -> sys_outbox
  -> NATS JetStream
  -> O2D / QMS / EAM / WHS handlers
  -> Capability / Process / Agent
  -> IAOS UI / AESE 2D simulation view
```

2D 预览器的数据链路：

```text
preview.json
  -> StaticScenarioDataSource
  -> SandboxScenario view model
  -> deterministic timeline reducer
  -> 2D canvas / event feed / KPI / Agent suggestion

IAOS snapshot API + SSE (later)
  -> IaosScenarioDataSource
  -> same SandboxScenario view model and UI
```

## 4. 场景包结构

M3 目标结构：

```text
scenario-packs/hctm/
├── manifest.json
├── master-data/
│   ├── organization.json
│   ├── parties.json
│   ├── materials.json
│   ├── manufacturing.json
│   └── logistics.json
├── stories/order-expedite-01/
│   ├── initial-state.json
│   ├── events.json
│   ├── expected-outcomes.json
│   └── preview.json
└── schemas/
    ├── manifest.schema.json
    ├── record-set.schema.json
    ├── event-sequence.schema.json
    └── expected-outcomes.schema.json
```

工具结构：

```text
cmd/aese/                    # validate / inspect / apply / replay 命令入口
internal/scenariopack/       # 加载、规范化和引用解析
internal/validate/           # schema、关系、时间线和幂等校验
internal/iaosclient/         # IAOS API 适配器，M3 后半段
internal/replay/             # 受控事件重放协调器
internal/agenttrace/         # M5 tool bundle setup、受审计查询和三 Agent 建议构建
frontend/src/scenario/       # SandboxScenario 与静态数据源边界
frontend/src/playback/       # 确定性重放 reducer 和 React hook
frontend/src/components/     # 2D 画布、控制栏、事件/KPI/详情面板
frontend/e2e/                # 三个目标视口的浏览器验收
```

这些路径已在 M3 中实现；`internal/replay` 负责默认 dry-run 的 apply/replay/verify 协调，正式写入仍受 IAOS API 合同约束。

M5 Agent tracer 的运行链路：

```text
agent-tools.json
  -> aese agent-setup（默认 dry-run）
  -> IAOS metadata schema + AI Tool Registry
  -> source_ref=entity.records
  -> latest tenant metadata + server-owned field/filter/order/limit allowlist
  -> explicit tenant predicate + PostgreSQL RLS
  -> ai_tool_call / call_id audit
  -> aese agent-run deterministic recommendation builder
  -> planning / quality / business_analysis suggestions
```

`entity.records` 不是任意实体 CRUD 或任意 SQL 入口。tool metadata 固定 entity、可返回字段、可过滤字段、排序和最大行数；调用 input 只能给标量 filter value 和 limit。未知字段、不安全物理映射和未注册 source 均失败关闭。AESE 不直接访问 IAOS 数据库，也不在本地复制权限或查询引擎。

## 5. 数据合同

场景包必须满足：

- `manifest` 声明 pack key、版本、租户模板、依赖和故事入口。
- 所有记录使用稳定业务编码，禁止引用环境特定 UUID。
- 每个文件声明 `schema_version`。
- 记录集合声明目标 entity 和自然键。
- 事件使用 IAOS envelope，带 `correlation_id`、`causation_id`、`idempotency_key`。
- 时间线严格单调或显式声明并发组。
- `expected-outcomes` 描述可机器验证的状态，不保存只供展示的散文答案。

## 6. 写入和事件治理

- `validate`、`inspect` 不访问网络。
- `apply` 默认 dry-run，只有显式 `--apply` 才能修改目标环境。
- 主数据和业务单据通过 IAOS 认证 API 或受治理配置包导入。
- 与业务记录同事务产生的领域事件必须由 IAOS Outbox 发布。
- 供应商延期、设备故障等外生仿真事件通过 IAOS simulation ingress 或受治理 Capability 注入。
- 本地开发可提供 direct-NATS replay，但必须使用 `--unsafe-direct-nats`，不得作为演示和生产默认路径。

## 7. 多租户与安全

- AESE pack 不携带真实 JWT、数据库连接串或生产凭据。
- 目标 tenant 在 apply 时绑定，pack 中只提供模板值或演示默认值。
- IAOS 写入必须经过 RLS 事务和权限检查。
- Agent 不能直接执行场景文件中的任意命令；只允许调用已注册的 IAOS Tool/Capability。
- M5 查询工具固定为 `low` risk、`none` confirmation 的只读 `entity.records`；工具注册、启用、调用和 call history 仍经过 IAOS AI Tool Registry 权限与审计边界。
- `agent-run --apply` 的外部副作用仅是 IAOS 受治理的 tool call/audit 记录；当前建议保持 `suggested` 且要求人工确认，不自动执行推荐动作。
- 每次 apply/replay 需要生成 run ID，并记录目标环境、pack 版本、actor 和结果摘要。

## 8. 可重复性

- L1 基础主数据可幂等 upsert。
- L2 故事初始状态按 scenario key 重置。
- L3 事件按 event ID 和 idempotency key 去重。
- 同一 pack 版本在同一 IAOS 版本上应生成一致的 dry-run 报告。
- 每条演示故事必须提供独立 reset 和 verify 步骤。

## 9. 当前架构缺口

- 场景包 JSON、schema、AESE CLI、校验器和 IAOS client 已实现并通过离线测试。
- IAOS 已实现 DES-048 三类 simulation ingress：设备停机、供应延期和来检失败通过权限、tenant、稳定对象解析、状态影响、幂等审计和事务 Outbox 进入运行链。
- IAOS 已实现 M3 allowlist 的原子、自然键幂等 scenario apply/reset，稳定编码在服务端解析 UUID，并显式绑定 tenant。
- order decompose 使用状态 CAS，O2D workflow 以 event/idempotency key 去重并在单一事务内执行；correlation、Outbox 和重复 no-op 已实证。
- HCTM 18 类事件尚未全部进入 IAOS `shared/eventdef`。
- O2D 当前只消费 `o2d.order.confirmed`，其余领域 handler 尚未接线。
- AESE 只读 2D 沙盘已实现：14 节点 A 线、七幕/22 事件、五项 KPI、对象详情和三类 Agent 建议均可确定性播放。
- DES-048 的 M4 三类入口已贯通。M3 的订单 tracer 仍只依赖 IAOS 内生 `o2d.order.confirmed`；异常领域消费者和自动重排产不属于 M4。
- M5 已有 `agent-setup` / `agent-run`、版本化 tool bundle、9 个 metadata 约束查询和 `internal/agenttrace` 三 Agent tracer。计划与质量结论可引用当前受治理业务状态；经营分析因缺少完工入库、发运和成本实际数据而明确返回 `partial`。
- Live 的 11,700 件累计实发和 300 件最终缺口来自 IAOS 受治理完工、库存和发运事实；只有成本实际仍保持数据缺口。
- M3V 静态预览器和 M6 `IaosScenarioDataSource` 均已完成；两种状态隔离，Live 故障不会静默降级为 Preview。

## 10. M6 在线观察架构

M6 采用 snapshot-first 模式：AESE 保留静态布局和视觉映射，IAOS 提供租户级运行事实、持久事件游标和建议证据。

```text
IAOS governed business action
  -> business state + inventory transaction + scenario event log + Outbox (one transaction)
  -> scenario snapshot / events?after=cursor / scenario SSE
  -> IaosScenarioDataSource
  -> existing 2D canvas / KPI / event feed / Agent panel
```

通用 `/api/v1/events/stream` 没有持久游标、断线补发和场景过滤，不作为在线沙盘的事实恢复合同。M6 的场景 SSE 只传递增量，客户端始终以一致性快照和持久 cursor query 恢复状态。

完工、入库和发运属于 IAOS 内生业务动作，必须走 Capability/场景业务动作；simulation ingress 继续只承载受控外生异常。详细合同见 DES-004。

## 11. M7 场景运行控制面

M7 在 AESE 增加无状态 orchestration API，把现有 CLI application service 暴露给浏览器。它负责 pack 加载、阶段编译和调用协调，不拥有业务数据、权限或运行事实。

```text
AESE UI
  -> AESE orchestration API (plan/state machine/idempotency)
  -> IAOS governed APIs (scenario/simulation/business action/AI Tool)
  -> IAOS RLS/audit/Outbox/snapshot/cursor
  -> AESE Live UI
```

运行状态由 IAOS run、snapshot、event 和 recommendation 重建；AESE 服务重启不能丢失或凭内存伪造阶段完成。浏览器只调用 AESE 控制 API 和 IAOS 只读观察 API，不直接调用 IAOS 写端点。详细边界见 ADR-003 和 DES-005。
