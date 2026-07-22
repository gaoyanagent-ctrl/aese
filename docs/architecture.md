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
| 企业业务数据库与 RLS | No | Own |
| 独立仿真事实 PostgreSQL | Own（不含业务台账） | No（只通过合同集成） |
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

## 12. AESE 2.0 演进边界

M8 已把 AESE 从纯场景内容/编排层扩展为最小企业生命周期仿真运行时；状态所有权以 accepted ADR-004 为边界：

```text
AESE World State
  -> 客观空间、资源、时间、物理/经济结果

IAOS Business State
  -> 企业登记、流程、权限、Capability、审计和 Outbox

Actor Knowledge State
  -> 特定人员/Agent 已观察、相信和可访问的信息
```

该边界允许 AESE 使用独立 PostgreSQL 保存仿真事件日志和快照，但不允许复制 IAOS 订单、库存、设备台账、流程、权限或业务数据库。F1 已实现不依赖墙钟的虚拟时钟、稳定事件队列、版本化纯 reducer、事件日志、state hash、快照恢复和离线 replay；PostgreSQL 持久接线仍后置。三态仅通过版本化 observation/intent/outcome 合同、稳定引用、租户、correlation 和幂等键关联；任何正式 IAOS 写入继续走受治理 API。

现有 M7 无状态编排控制面保持独立并作为兼容基线，不直接承担持续 World Runtime。详细方案和实施门见 DES-007、ADR-004 与 PLAN-M8-001。

World/IAOS 桥采用 DES-008 的持久 journal + cursor 模式：AESE 提交 actor-scoped observation；人员或 Agent 通过 IAOS Capability/Process 形成 intent；只有 IAOS 事务已提交或确定 no-op 的 committed outcome 才能驱动 AESE world consequence。SSE/Outbox 用于通知，断线恢复始终读取 journal，不依赖 webhook 或 direct NATS。

## 13. M9 企业成立与治理纵向架构

M9 复用 M8 World Runtime、World Store、三态模型和 IAOS Bridge，不建立第二套项目、组织、预算或 Agent 运行时：

```text
Founder / CEO / CFO human or deterministic role policy
  -> IAOS Capability / Process / Policy
  -> intent + governed business records
  -> committed outcome + journal + Outbox
  -> AESE schedules regulator/bank/world activity
  -> registration/account/capital/appointment world consequence
  -> actor-scoped observation and Knowledge
  -> reconciliation and plant_project_eligible
```

虚构监管机构和银行属于 AESE 外部世界策略，不是 IAOS 用户。法人档案、治理决议、组织岗位、资本记录和预算审批属于 IAOS；登记生效、账户实际开立、资金实际到账和岗位实际接受属于 World State。详细边界见 DES-010 和 PLAN-M9-001。

## 14. M10 工厂选址与设施建设纵向架构

M10 消费 M9 的 `plant_project_eligible=true`，继续复用同一 World Runtime、三态模型和 IAOS Bridge。它不把 IAOS 项目记录误当成现实建设进度：

```text
Project Director / CEO / CFO
  -> IAOS site, investment, contract, project and payment governance
  -> committed outcome + journal + Outbox
  -> AESE contractor, landlord and utility-provider world strategies
  -> actual site control, construction, delay and inspection consequence
  -> actor-scoped observation + discrepancy + rebaseline intent
  -> governed acceptance
  -> capability_build_eligible
```

候选场地、实际占用、公用工程容量、现场施工和检查结果属于 AESE World State；选址评估、投资决议、租赁合同、项目/WBS、变更、付款和验收记录属于 IAOS；角色只通过 observation 获得其可见的 Actor Knowledge。计划到期或 IAOS 中把里程碑标为完成，都不能单独推进实际工程状态。

M10 只建立区域、城市、园区、场地、建筑、楼层和功能区的最小空间层级，不建设 BIM、3D 或布局编辑器。设施验收后仅输出 M11 的能力建设资格；生产设备、检测仪器、招聘培训和投产准备仍由 M11 负责。详细边界见 DES-011 和 PLAN-M10-001。

## 15. M11 生产能力建设纵向架构

M11 消费 M10 的 `capability_build_eligible=true`，把已验收设施转化为“可进入产品工业化”的通用设备与人员能力，但不把资产卡、员工档案或培训记录误当作现实能力：

```text
Project Director / CEO / CFO / Procurement / HR / Equipment / Quality
  -> IAOS capital, budget, procurement, asset, org and qualification governance
  -> committed outcome + journal + Outbox
  -> AESE investor/bank, supplier, logistics, labour-market and training strategies
  -> actual cash, delivery, installation, commissioning, onboarding and learning
  -> actor-scoped observation + discrepancy + remediation intent
  -> governed equipment/person qualification and joint capability acceptance
  -> industrialization_eligible
```

实际资本到账、设备制造/运输/安装/调试、候选接受/到岗和技能掌握属于 AESE World State；预算、订单/租赁、资产、编制、招聘、员工、培训、认证和班次记录属于 IAOS；角色只通过权限和 observation 获得 Actor Knowledge。IAOS 状态或虚拟时间不能直接创造现金、设备能力或人员技能。

M11 复用 M10 空间与 utility 约束，并显式保留设施尾款、工资准备金和现金缓冲。它只建立冷却板产品族的通用能力需求，不发布正式产品 BOM/routing，也不执行 APQP、试生产、PPAP 或 SOP；这些属于 M12。详细边界见 DES-012 和 PLAN-M11-001。

## 16. M12 产品工业化与量产批准纵向架构

M12 消费 M11 的 `industrialization_eligible=true`，建立客户项目、产品/工艺发布、试制验证与客户 PPAP 闭环，但不把报价、工程记录或 PPAP 提交误当成客户批准和现实产品能力：

```text
Sales / Project / Product / Process / Quality / Procurement / Plant / CFO
  -> IAOS RFQ, quotation, project, revision, APQP, quality and release governance
  -> committed outcome + journal + Outbox
  -> AESE customer, supplier, tooling, material and physical trial strategies
  -> actual nomination, material, build, measurement, defect and PPAP decision
  -> actor-scoped observation + discrepancy + engineering-change intent
  -> governed retest, PPAP registration and production release
  -> serial_production_eligible
```

客户实际 RFQ/定点/PPAP 决定、供应商/材料真实能力和试制/测量结果属于 AESE World State；报价、客户项目、产品/BOM/routing revision、APQP、问题/变更、PPAP package 和生产放行记录属于 IAOS；角色只通过权限与 observation 获得 Actor Knowledge。IAOS 状态、UI 或虚拟时间不能直接改写测量值、Cpk、客户决定或物理试制结果。

现有 `scenario-packs/hctm` 中的 `HCTM-BCP-A01`、BOM 和 routing 是 M3/O2D 兼容 fixture，不是 Genesis 已完成工业化的历史事实。M12 在独立 campaign 中生成 release manifest/hash，M13 只有在版本兼容检查通过后才能投影到旧 O2D 对象。正式订单、批量量产、交付、开票和回款属于 M13。详细边界见 DES-013 和 PLAN-M12-001。

## 17. M13 第一次完整商业交付纵向架构

M13 消费 M12 的 `serial_production_eligible=true`，复用 IAOS O2D 运行能力完成首张正式订单到现金和实际毛利的闭环，但不把订单、发票或收款记录误当作客户需求、物理交付或银行到账：

```text
Sales / Planning / Procurement / Production / Quality / Logistics / Finance
  -> IAOS order, MRP, supply, production, inventory, shipment and finance governance
  -> committed outcome + journal + Outbox
  -> AESE customer, supplier, equipment, transport and bank world strategies
  -> actual demand, material, output, delivery, acceptance and cash consequence
  -> actor-scoped observation + 300-unit discrepancy + recovery intent
  -> governed invoice, AR, settlement, actual-cost and margin close
  -> first_commercial_cycle_closed
```

客户正式需求/变更/接受/付款行为、供应/设备/运输现实、实际生产资源消耗和银行到账属于 AESE World State；订单、采购、工单、库存、发运、发票、应收、收款核销、成本和项目损益记录属于 IAOS；角色只通过权限和 observation 获得 Actor Knowledge。发运不等于客户接受，发票不等于现金，毛利不等于现金余额。

M13 从 M12 release manifest 编译 Genesis-specific O2D 输入，期初可销售成品为零，并使用新交易/correlation code；旧 M3/M7 pack 继续作为回归 fixture，不能贡献库存或既成事件。M13 关闭 M9-M13 主纵向场景，长期多周期与参数化实验属于 M14。详细边界见 DES-014 和 PLAN-M13-001。

## 18. M14 参数化分支经营实验架构

M14 不把 M13 的单条闭环复制成第二套业务系统，而是在 AESE World 所有权内建立受治理实验层：

```text
Approved Genesis checkpoint + experiment definition
  -> immutable parameter/policy/seed matrix
  -> isolated World branches with common random numbers
  -> bounded deterministic multi-cycle runs
  -> run-level invariants, snapshots, hashes and KPI
  -> paired comparison + constraint evaluation + Pareto frontier
  -> immutable EvidenceBundle
  -> optional IAOS recommendation intent (never auto-apply)
```

AESE 拥有 checkpoint fork、branch/run、随机流、外生事件、World consequence、artifact 和实验 KPI；IAOS 只拥有实验申请/审批、访问权限、证据接收、推荐与后续业务决策；Actor Knowledge 继续按角色权限投影。分支不创建正式 IAOS 订单、库存、发票或现金，父 checkpoint 和兄弟分支不可变。

随机性必须是版本化输入：固定 PRNG 和命名 stream，同一 scenario profile/seed 的各策略使用共同随机数并做 paired comparison。默认展示硬约束、KPI 分布和 Pareto 前沿，不以单次运行、平均值或未批准权重宣称最优。只有 run 清单完整、hash 可复验、约束计算完成且无未解释失败时输出 `strategy_evidence_ready=true`。详细边界见 DES-015 和 PLAN-M14-001。

## 19. M15 受治理策略发布与经营试点架构

M15 在 M14 “只形成证据、不自动投放”的边界之后增加独立 decision-to-action 控制面：

```text
M14 EvidenceBundle + Pareto candidate
  -> evidence verify + independent ChangeRequest review
  -> immutable StrategyRelease + semantic diff + SafetyEnvelope
  -> zero-write shadow on canonical observations
  -> separate pilot approval
  -> bounded canonical operating cycle through IAOS governed actions
  -> World consequence + guardrail/drift monitoring
  -> adopt | reject | pause/rollback + compensation
  -> strategy_change_cycle_closed
```

IAOS 拥有 ChangeRequest、StrategyRelease、审批、Policy/Capability、incident、rollback、compensation 和 AdoptionDecision；AESE World 拥有 canonical pilot 时间、外生事件、物理/经济后果和 guardrail observation；Actor Knowledge 继续按当时权限和 observation 投影。shadow 只能计算 candidate decision，不能提交 intent 或预占资源。pilot 中每个业务动作必须由 IAOS committed outcome 驱动，IAOS release 状态不能直接创造库存、产能、交付或现金。

回滚只停止 candidate release 对未来决策生效并恢复 prior release，不删除已经提交的采购、生产、库存、发运、发票或现金事实；遗留影响进入 open-commitment ledger 和受治理 compensating action。M15 不用单个 pilot 宣称因果改善，终态允许 `adopted`、`rejected` 或 `rolled_back`，只有过程、breach 和 commitment 全部对账后才输出 `strategy_change_cycle_closed=true`。详细边界见 DES-016 和 PLAN-M15-001。

## 20. M16 持续策略保障与假设校准架构

M16 消费 M15 主路径已采纳 release，但不把 adoption 当作永久有效：

```text
Adopted StrategyRelease + canonical World/IAOS observations
  -> frozen ObservationSpec + as-of cutoff/cursor
  -> immutable AssuranceDataset + quality/correction lineage
  -> data quality gate
  -> input/process/policy-action drift assessment
  -> bounded 8-week CalibrationCandidate
  -> single-use 4-week holdout
  -> new-ancestry M14 common-random-number replay
  -> renew | reexperiment_required | retire
  -> strategy_assurance_cycle_closed
```

AESE World 继续拥有外生事件和物理/经济结果；IAOS 拥有业务记录、active release、AssuranceCycle、finding、审批和 AssuranceDecision；Dataset 只保存允许的观察、stable refs 和 hash，不复制 IAOS 业务数据库。cutoff 后迟到事实通过 correction set 和新 dataset version 处理，不能静默改旧 evidence；Actor Knowledge、UI 或通知不能替代 source facts。

数据质量必须先于 drift 和校准。前 8 周 calibration 与后 4 周 single-use holdout 严格隔离；候选参数只能修改批准的外生假设，不能改 hard constraint、KPI 或 StrategyRelease。新 replay 保留原 M14 EvidenceBundle，终态允许 `renewed`、`reexperiment_required` 或 `retired`。详细边界见 DES-017 和 PLAN-M16-001。

## 21. AESE 3.0 完成体与 M17 滚动计划架构

M16 主路径 `renewed` 后，M17-M24 按以下顺序扩大系统能力，而不是并行建设八套新业务系统：

```text
M17 rolling IBP/S&OP
  -> M18 product/customer portfolio
  -> M19 multi-site network
  -> M20 after-sales and closed-loop quality
  -> M21 plant resource/EHS resilience
  -> M22 group value/treasury/investment
  -> M23 governed multi-agent organization
  -> M24 scenario platform productization
```

M17 首先建立跨需求、供应、产能、库存、交付、成本和现金的版本化计划控制面：

```text
M16 renewed assumptions + actual refs
  -> Demand Review
  -> Supply Review
  -> Financial Reconciliation
  -> Pre-IBP gap/options
  -> Executive IBP decision
  -> immutable PlanRelease
  -> separate governed execution intents
```

IAOS 拥有 PlanningCycle、计划版本、review/decision、权限和 PlanRelease；AESE 拥有未来外生实现、资源约束和 scenario consequence；Knowledge 只包含角色已知假设和版本。Forecast 不等于订单，planned receipt 不等于实际到货，PlanRelease 不自动产生 PO/WO/Shipment/资金动作。

M17-M24 已按顺序实现为 `internal/aese3` 的八个 immutable evidence frame，并由 schema、fixture、100 次 canonical hash、API 和 Completion Room 共同验证。IAOS DES-059 只接收 exact evidence 下的受治理动作；Program approval 不产生自动业务写入。Program 总边界见 DES-018，最终 reference pack 为 `hctm-genesis@1.0.0`。
