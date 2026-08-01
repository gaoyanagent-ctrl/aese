# AESE Roadmap

## 2026-08-01 产品知识与用户手册

- 已完成 AESE/IAOS 知识所有权、场景 Article 扩展和 Agent 回答合同设计（DES-037）。
- 已建立 M9 企业设立用户手册第一版，明确 human/agent/approval/world 四类节点与证据入口。
- 已交付场景知识 Schema、带 canonical hash 的 M9 manifest、18 节点稳定引用校验和任务弹窗
  “这一步是什么”入口；IAOS 已实现租户知识发布和正式只读 Agent Tool。
- 已完成 workspace/case/world run/node/actor/capability 封闭导航上下文深链；IAOS 知识中心
  可见展示、允许清除并在 Copilot BFF 二次校验，明确不能把浏览器参数当作运行事实。
- 下一步按 PLAN-KNOWLEDGE-AESE-001 完成 World/IAOS 双侧实际证据、配置漂移和失败关闭验收。
- 该计划是 PLAN-GXZ-001 的并行知识治理计划，不改变 M9 财务未完成项状态。

本文件是 AESE 当前里程碑状态和下一步优先级的权威来源。

最后更新：2026-08-01。

2026-07-31 M9 原子能力运行合同：IAOS DES-076 已交付 19 项真实可执行 V1 原子能力
（18 active、1 deprecated），同一目录驱动 Analyzer、Artifact、Runtime、API 和 Studio。
AESE 只消费已发布 Business Capability、Process 和 Bridge，不复制原子 Handler，也不把
审批、过账、完工或 LLM 调用误当成原子能力；M10+ 继续复用该目录并按平台评审新增。

2026-07-31 M9 运行权威收口：浏览器的案件、工作项、Agent、审批和 World
Observation 写操作已统一进入 AESE 同源白名单 Command Gateway；IAOS 原生 Capability
执行校验标准 Artifact，设立工作项只从 active 主 Process Artifact 递归展开，并固定
Capability 与子流程的版本/哈希。Go 流程目录只作为平台包发布源，不再作为运行兜底。
该合同对应 IAOS Runtime 2.10.0 / `genesis-m9@1.2.0`。

2026-07-31 Entity 存储与唯一写入权威：IAOS Runtime 2.9.0 /
`genesis-m9@1.1.0` 把三类存储合同编译进 Effective Runtime Artifact。普通动态
Entity 以 `entity_record_<code>` 作为唯一权威存储；复杂领域对象从专用权威表投影到
`entity_projection_<code>`；Journal/Aggregate/Query 汇总只生成只读投影。通用 CRUD
对领域和计算投影失败关闭，M9 的 35 个 Entity 已明确分类为 27 个领域投影和 8 个计算
投影。

2026-07-31 原生 Entity 投影真实性：IAOS Runtime 2.8.1 / `genesis-m9@1.0.5`
取消把完整设立案件状态复制到所有投影的旧路径，改为由拥有事实的 Agent output、
Capability journal 或权威领域表逐项物化。逻辑 Entity code 保持
`governance_resolution`，物理读模型按 `entity_projection_<code>` 命名；节点未执行时
对应 Entity 必须为零记录，已执行节点的菜单与详情必须返回编译后的显式字段。

2026-07-30 补强：M9 账套启用已支持用户输入账套名称、确认当前会计年度并逐期检查
12 个默认自然月期间；IAOS 对全年连续覆盖失败关闭。

2026-07-30 财务文档模块化：`docs/designs/finance/` 成为独立财务设计目录；DES-030 为
短总览，DES-031–036 分别维护多组织共享数据、会计内核、子账资金、成本资产、预算关账
报表和治理 Agent。多组织基线已冻结集团/法人/BU/基地/共享服务中心、多账簿、
共享科目表、Data Set、Business Partner canonical/组织扩展；模块期间转由 DES-035。
F5B–F5D 已实现，F5E 保持关账阶段范围，不能把设计完成标记为持续经营财务完成。

2026-07-30 F5B 实现：IAOS 已交付集团/法人/BU/基地/工厂/财务组织/共享中心、
Data Set、决定项分配和主体组织访问基础，使用 FORCE RLS、复合 tenant 外键和
Capability 写入硬门；AESE 已把 HCTM 确定性组织与共享数据模板纳入离线校验。
2026-07-30 F5C 实现：IAOS 已交付共享科目表/科目定义、法人科目扩展、财政日历/期间、
账簿/账簿集和受治理配置 UI；目标租户保留原账簿与凭证关联并得到完整 12 期。AESE
baseline 1.2 提供可版本化 HCTM 模板和离线引用/期间连续性校验。

2026-07-30 F5D 实现：IAOS 已交付集团唯一 Business Partner、客户/供应商角色、法人/BU
扩展、共享产品和工厂扩展，以及两项受治理 Capability 和“伙伴与产品”业务表单；AESE
baseline 1.3 固定星河客户、两家铝材供应商、冷却板成品和铝板原材料模板。

2026-07-30 F5D2 资产发布：IAOS Runtime 2.3.0 已把 F5B–F5D 的 16 类权威模型发布到
数据模型工坊，Runtime catalog 达到 35 个 Entity/33 条语义关系；“财务账务与报表”
收敛为查询与输出，组织待办、多组织共享、账簿基础、客户、供应商、产品成为独立入口，
权威表与 Entity 投影由受治理 Capability 同事务同步。

2026-07-30 平台基础语义治理：IAOS 已发布
`iaos_foundation_semantics@1.0.0`（87 Concept/51 Archetype）和 Semantic Governance
Compiler；M9 Runtime 2.4.0 只消费平台基础包。`business_partner → party`、
`organization → party`、`product → material` 成为机器校验合同。后续 AESE/M10
新增语义必须先检索现有目录、记录不可复用原因并取得产品所有者批准，场景包不得覆盖
foundation-owned 资产。

2026-07-30 Effective Runtime Artifact 权威闭环：IAOS DES-072 已将 Entity
Schema/UI/Agent Context、Capability API/Agent Tool 和 Process Run 收敛到同一不可变
编译产物；Process 发布冻结 Capability version/hash，缺失、过期、编译器或哈希不一致
均失败关闭。M9 不再允许从 authoring DSL、metadata 表或旧运行回执做正式路径兜底。

2026-07-30 平台基础包 Edition：IAOS 已发布 `genesis-m9@1.0.0`，由
`enterprise-governance`、`finance-foundation`、`genesis-incorporation` 三个不可变
签名包组成，覆盖七类通用资产。tenant-001 安装同一参考 Edition，Genesis 新租户的
`runtime_installed` checkpoint 和历史租户显式升级消费同一清单；参考安装不复制任何
租户身份或业务事实。

## 1. 里程碑状态

| 里程碑 | 目标 | 状态 | 完成证据 |
| --- | --- | --- | --- |
| M0 项目初始化 | 仓库、背景、规则、GitHub | Completed | README、AGENTS、初始提交 |
| M1 虚拟企业蓝图 | 华辰集团、苏州基地、电池冷却板 A 线 | Completed | HCTM Virtual Enterprise Blueprint |
| M2 业务与技术规格 | 对象、事件、seed、演示故事 | Completed (docs) | 4 份 HCTM 规格文档 |
| M2.5 工程治理 | 架构边界、索引、code map、执行规则 | Completed | 本轮治理文档 |
| M3 可执行场景包 | JSON 场景包、校验器、IAOS apply/replay tracer | Completed | pack、CLI、execution evidence、IAOS commits |
| M3V 快速 2D 沙盘 | 七幕故事、22 事件、A 线画布、KPI 和 Agent 建议预览 | Completed | 前端、preview、18 unit/component tests、9 E2E、3 viewport screenshots |
| M4 异常场景运行 | 延期、设备故障、来料不良进入 IAOS 运行链 | Completed | 三类 ingress、状态影响、事务 Outbox、租户/幂等及 canonical replay evidence |
| M5 Agent MVP | 计划、质量、经营分析 Agent | Completed | 9 个受治理只读工具、三 Agent live tracer、跨租户与零业务写入证据 |
| M6 在线 2D 企业沙盘 | IAOS 实时事件、库存、产线、异常和 Agent 运行结果 | Completed | DES-004、PLAN-M6-001、M6 evidence |
| M7 受治理场景运行控制台 | 浏览器预检、初始化、逐幕运行、分析、验证和复位 | Completed | ADR-003、DES-005、PLAN-M7-001、M7 evidence |
| M8 AESE 2.0 基础 | 三态世界、确定性离散事件内核、IAOS 双向桥和最小 Genesis tracer | Completed | PLAN-M8-001、World Play runbook、两仓测试与部署证据 |
| M9 Genesis Incorporation | 注资、法人登记、治理、管理岗位、初始组织与预算 | Completed | hctm-genesis@0.2.0、M9 evidence、IAOS DES-051 |
| M9N IAOS-native Incorporation Closed Loop | 用三层语义、正式身份、Runtime Artifact、Capability/Process/Approval、持久工作项和 World Bridge 重建 M9 真实闭环 | Completed | Runtime 2.11.0；`genesis-m9@1.3.0` 三包 Edition；31 个 Capability；23 个显式工作项（含 5 个财务开业节点）、G1–G7、三个 World wait、五次 Agent Run/五 Agent；35 个 Entity 具备唯一存储/写入权威合同，资本来源能力委托通用凭证过账能力 |
| GX Enterprise Genesis Game Experience | AI 企业身份、人工/Agent 协作和 2.5D 世界中的 M9 游戏化开局 | Completed | DES-028/DES-029、从空白 case 创建、23 工作项游戏内操作、G1–G7、三个 World wait、首页 AI 创意官可访问入口与三视口 live 浏览器验收 |
| GX-ZERO Zero-start Enterprise Genesis | 从产品主页创建独立 tenant、World Run、真实 AI 企业身份和 M9 企业 | Completed | IAOS Player 注册/密码登录与既有账号安全提升、生产 Workspace 控制面、八 checkpoint、tenant-only owner session、旧 local Workspace 安全接管、五步向导、持久 CreativeJob、动态 Player subject，以及全新 Workspace 的 23/23 节点、7 审批门、3 World wait、6 Agent run 验收 |
| M9-FIN Manufacturing Finance Foundation | 在 M9 建立财务组织、账套、科目、期初资本会计和开业报表，并在 M10–M13 接通 AP/AR/资金/资产/成本/总账 | F11A Complete; F5E Deferred; M10–M13 Planned | M9 开业纵切、多组织/共享数据、账簿/伙伴/产品已交付；Runtime 2.11.0 新增通用凭证过账能力并修复凭证主子写入所有者；模块期间和 F15–F35 仍不得计入完成 |
| M10 Genesis Plant Build | 选址、场地控制、设施项目、公用工程、异常重排与验收 | Reference Replay Complete; D22 Pending | hctm-genesis@0.3.0 与既有 evidence 仅证明确定性 replay；交互式工作项未验收 |
| M11 Genesis Capability Build | 资金补足、设备/实验室/仓储能力、核心团队与岗位资格 | Reference Replay Complete; D22 Pending | hctm-genesis@0.4.0 与既有 evidence 仅证明确定性 replay；交互式工作项未验收 |
| M12 Genesis Industrialization | RFQ/定点、产品/工艺、供应商/工装、APQP、试制、PPAP 与量产批准 | Reference Replay Complete; D22 Pending | hctm-genesis@0.5.0 与既有 evidence 仅证明确定性 replay；交互式工作项未验收 |
| M13 Genesis First Delivery | 正式 O2D、三批交付、客户接受、开票/回款、实际成本与项目毛利 | Reference Replay Complete; D22 Pending | hctm-genesis@0.6.0 与既有 evidence 仅证明确定性 replay；交互式工作项未验收 |
| M14 Parameterized Branch Experiments | checkpoint 分支、多周期参数/策略、共同随机数、实验执行与决策证据 | Reference Replay Complete; D22 Pending | 既有实验 replay 保留；交互式决策工作项未验收 |
| M15 Governed Strategy Release & Pilot | evidence 审议、版本化发布、shadow、canonical pilot、guardrail 与回滚/采纳 | Reference Replay Complete; D22 Pending | 既有策略 replay 保留；交互式治理工作项未验收 |
| M16 Continuous Strategy Assurance & Calibration | canonical observation、数据质量、drift、8/4 校准/留出和策略复审 | Reference Replay Complete; D22 Pending | 既有校准 replay 保留；交互式复审工作项未验收 |
| M17 Rolling IBP & S&OP | 13 周执行/12 月财务计划、五级 review、PlanRelease | Reference Replay Complete; D22 Pending | 既有 evidence 保留；交互式五级 review 未验收 |
| M18 Product & Customer Portfolio | 第二产品/客户、共享能力和组合权衡 | Reference Replay Complete; D22 Pending | 既有 evidence 保留；交互式组合决策未验收 |
| M19 Multi-site Network | 第二制造/外协节点、物流节点和网络韧性 | Reference Replay Complete; D22 Pending | 既有 evidence 保留；交互式网络决策未验收 |
| M20 Closed-loop Customer Quality | 售后、质保、RMA、8D/CAPA 和追溯 | Reference Replay Complete; D22 Pending | 既有 evidence 保留；交互式质量工作项未验收 |
| M21 Plant Resilience | 资产、人员、EHS、能源和业务连续性 | Reference Replay Complete; D22 Pending | 既有 evidence 保留；交互式应急工作项未验收 |
| M22 Group Value Governance | 集团管理财务、资金、营运资本和投资 | Reference Replay Complete; D22 Pending | 既有 evidence 保留；交互式价值治理未验收 |
| M23 Governed Multi-Agent Organization | 七 Agent、协作、评估、接管与安全 | Reference Replay Complete; D22 Pending | 既有 evidence 保留；真实多 Agent 工作项尚未验收 |
| M24 Scenario Platform Productization | SDK、authoring、certification、registry 和发布 | Reference Replay Complete; D22 Pending | 既有 SDK evidence 保留；D22 交互式认证门尚未验收 |
| X1 System Atlas 全景治理 | 最终完成体、当前状态、依赖与进展历史 | Completed | DES-006、IAOS DES-049、双端动态图谱 |

## 2. 当前阶段

M3、M3V、M4、M5、M6、M7 和跨里程碑的 X1 System Atlas 已完成。联动中心已支持联动检查与受治理场景运行，不依赖 CLI 完成 preflight、initialize、七幕推进、Agent 分析、verify 与 reset。

PLAN-M8-001 至 PLAN-M24-001 均已完成。M17-M24 严格消费前一 terminal，统一证据输出 `industry_simulation_platform_ready=true`。

当前唯一 active 主计划为 `PLAN-GXZ-001`。`PLAN-M9-NATIVE-001` 曾因通用平台资产缺口恢复为 active remediation。专用状态机、
事务、身份、G1–G7、World Bridge 与 trace 证据保留；完成 D19–D21 的通用资产注册、
十工作区、逐步骤追踪及可解释配置合同已经交付，但用户验收证明这些查看入口不能替代
持久化工作项和真实参与者推进。M9N 因 DES-027 D22 再次打开；只有人工、Agent、审批、
World wait 和断点恢复的交互式证据齐备后才能重新关闭。

D23/T83 已完成完整交互主线：Runtime 1.3.13 将登记校验正式分派
`legal-compliance-agent`，23 个持久工作项依次经过 Founder、五 Agent、G1–G7、三个
World wait；集成测试在第二个 World wait 后重建 Server 并恢复至
`enterprise_operational_ready`。当前 Agent 执行器仍明确为 IAOS 内置确定性 Runtime，
未配置外部模型时不冒充已连接 LLM。

D25 已消除数据模型工坊 Entity 生命周期与正式审批运行时的语义分叉：Entity 只声明
状态和受治理动作，需审批的动作引用租户版本化 Flow，审批决定来自冻结 assignment，
批准后才提交状态并 consume。该平台收敛不替代 M9 T83/T96 的完整交互式验收。
D26 进一步要求 System Atlas 展示 Semantic → Entity → Capability → Approval Flow →
Process → Runtime Evidence 的配置与运行依赖，并为新菜单/主要页面提供“功能说明”。
IAOS DES-060 已交付流程/审批/Entity 首批入口；AESE World 后续按同一合同渐进补齐。
M9 World 深链已补充失效案件恢复：trace 404 时保留 World 基线并列出当前 recent cases，
不再让已清理的旧 case 阻断页面；身份、权限和平台错误仍失败关闭。

PLAN-GX-001 已完成 MVP：游戏从 IAOS verified evidence bundle 与 23 个工作项编译
GameProjection；当前 deterministic 企业身份候选通过 `incorporation.case.open` 才成为正式案件事实；
PixiJS 2.5D、DOM/2D fallback、工作项深链、World Observation、治理证据和
`enterprise_operational_ready` 已在 1440×900、1280×720、390×844 三视口通过 live
浏览器恢复验收。它没有完成独立 tenant provisioning 或真实外部 LLM。PLAN-GXZ-001
以根主页、独立 workspace/tenant/World Run、MiniMax CreativeJob 和双租户 RLS 证据
关闭该缺口；凭据只从本地 secret/env 注入，不放入 AESE 仓库。

DES-030 已批准把完整制造企业财务底座纳入 Project Genesis。该扩展不撤销现有 M9/M9N
闭环证据，但新增 `finance_opening_ready` 完成门：M9 建立财务组织、账套、科目、期间、
实收资本凭证和开业资产负债；M10–M13 再按真实工程、采购、资产、生产和销售事件启用
AP、AR、固定资产、制造成本和总账。PLAN-M9-FIN-001 已进入 active：M9 已完成首批
财务 Entity/Capability/Process、资本双分录、开业报表、就绪硬门和 IAOS 财务穿透页；
25 个 M9 Capability 已按 M9 D27 进入统一执行边界；财务开业 5 个动作成为显式工作项，无 Capability Execution
Context 的凭证 DML 由数据库拒绝。F10、F13 已完成；F15–F35 和非财务受治理写入迁移仍未
完成，不得把当前开业切片描述为完整会计系统或全平台 Capability 强制迁移完成。

M7 O0-O4 已完成。最终 `m7-acceptance-20260722-05` 从 clean reset 跑通编排 API 与 CLI 对照链：22 个事件、三 Agent、17 条离线业务断言、2 条在线 IAOS 断言和 M6 KPI 均通过；单 run 产生 9 次成功 Tool Call 与两套一致的 O2D Outbox 副作用，UI/CLI 均安全复位。AESE 8090/4173 与 IAOS 8082/3000 的本机部署和健康检查已记录在 M7 evidence。该基线由 M8 强制保留。

M7 的最小成功标准：

1. 浏览器不直接调用 IAOS 写端点，由 AESE 薄编排 API 复用现有 Go 内核。
2. 运行具有 run ID、阶段状态机、plan hash、cursor 和 idempotency key。
3. 用户可从页面初始化、逐幕推进、运行到结束、分析、验证和安全复位。
4. 刷新、断线、重复点击和 AESE 服务重启不会产生重复业务副作用。
5. 权限不足、跨租户、陈旧 cursor 和非法状态转换全部失败关闭。
6. UI 与 CLI 对同一 pack 产生一致的 22 事件、Agent 建议、断言和 KPI。

## 3. AESE 3.0 M17-M24 Program

| 顺序 | 主题 | Terminal | 当前状态 |
| --- | --- | --- | --- |
| M17 | 滚动 IBP 与 S&OP | `integrated_plan_cycle_closed` | Completed |
| M18 | 多产品与多客户组合 | `portfolio_operating_model_validated` | Completed |
| M19 | 多基地供应履约网络 | `network_operating_model_validated` | Completed |
| M20 | 售后质保与闭环质量 | `customer_lifecycle_closed` | Completed |
| M21 | 资产人员 EHS 韧性 | `plant_resilience_cycle_closed` | Completed |
| M22 | 集团财务资金与投资 | `group_value_cycle_closed` | Completed |
| M23 | 受治理多 Agent 组织 | `agent_operating_model_qualified` | Completed |
| M24 | 场景平台产品化 | `industry_simulation_platform_ready` | Completed |

Program 总边界以 DES-018 为准。M24 关闭本轮 AESE 3.0 规划；真实生产、法定合规、第二行业和高精度 3D 必须另立 program，不继续隐式追加里程碑。

## 4. M17 已完成范围

包含：

- M16 renewed assumption/actual refs，13 周 weekly execution 和 12 月 monthly financial horizon。
- Demand、Supply、Capacity、Inventory、Delivery、Cost、Margin、Working Capital 和 Cash 计划。
- baseline/upside/downside 三个 scenario、frozen/slushy/liquid fence 和跨职能 Gap/Option。
- Demand Review、Supply Review、Financial Reconciliation、Pre-IBP、Executive IBP 五个 gate。
- immutable PlanRelease、replan trigger、IAOS 权限/职责/事务/Outbox 和 Executive IBP Room。

不包含：

- 第二产品/客户/工厂、多基地网络和 portfolio optimization。
- 真实预测模型、自动求解器、计划自动执行或 real-production target。
- 法定财务、税务、真实银行和无限 horizon。

## 5. M17 已完成交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| B0 | Planning Cycle 与集成计划机器合同 | Completed |
| B1 | Demand Plan、订单与假设 | Completed |
| B2 | Supply、Capacity、Inventory 与交付计划 | Completed |
| B3 | Financial Reconciliation 与经营价值 | Completed |
| B4 | Gap、Scenario、Pre-IBP 与 Executive Decision | Completed |
| B5 | IAOS Planning 治理与 Bridge | Completed |
| B6 | Executive IBP Room | Completed |
| B7 | 全链验收与 M18 入口 | Completed |

## 6. M17 完成条件

- horizon/calendar/fence/cutoff/unit/version 和 opening reconciliation 完整。
- demand/supply/capacity/inventory/financial plan 跨 bucket 数量金额守恒。
- 五 gate、职责分离、gap/option/decision 和 exact PlanRelease hash 可审计。
- approved plan 与 order/PO/WO/shipment/cash execution 强隔离。
- tenant/RLS、CAS/幂等、事务/Outbox、恢复、replan、IBP Room、两仓测试和 M3-M16 回归完整。
- 输出 `integrated_plan_cycle_closed=true` 与 approved/replan_required/deferred。

## 7. M17 风险与依赖

- G4-G8 未冻结 horizon、计划语义、scenario、review gate 和 IAOS gap 前，不得创建 PlanRelease。
- forecast 不等于 order，planned receipt 不等于 actual arrival，capacity plan 不等于物理能力，financial plan 不等于会计事实。
- weekly/monthly aggregation、unit/precision 和 time fence 错误会制造隐性不守恒，必须先于 UI 冻结。
- Executive approval 不得自动创建业务动作；执行必须另走 Capability/Process/Policy。
- downside 不能通过隐藏 hard gap 或改 assumption 变成“批准”；replan/deferred 是合法终态。
- 当前工作区已有其他人的测试修改、截图删除和验收产物，实施 agent 必须保留并避免重叠修改。

## 8. M16 已完成范围

包含：

- adopted release、12 周 canonical observation、World/IAOS stable refs、as-of cutoff/cursor 和 correction lineage。
- missing/late/duplicate/conflict/unit/version/owner 数据质量门。
- demand/supplier/equipment/quality/payment 输入 drift、process mismatch 和 policy-action pressure。
- 前 8 周有界 CalibrationCandidate、后 4 周 single-use holdout 和防泄漏约束。
- 新 ancestry 的 M14 common-random-number replay、ValidationReport 和原 evidence 不可变验证。
- scheduled/expiry/alert/manual assurance trigger、renew/reexperiment/retire 决策和 Assurance Observatory。

不包含：

- 真实生产数据、在线学习、自动调参/Policy 发布或 Agent sole approval。
- 第二客户/产品/工厂、完整 S&OP 或通用数据湖/机器学习平台。
- 反复查看 holdout、静默改历史 evidence，或用短窗口宣称真实概率/因果/永久最优。

## 9. M16 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| A0 | Assurance、Dataset 与 Drift 机器合同 | Completed |
| A1 | Canonical Observation Dataset 与 Lineage | Completed |
| A2 | 数据质量、Drift 与策略压力 | Completed |
| A3 | 有界 CalibrationCandidate 与防泄漏 | Completed |
| A4 | Holdout、再实验与 ValidationReport | Completed |
| A5 | 周期监控、到期与 IAOS 治理 | Completed |
| A6 | Strategy Assurance Observatory | Completed |
| A7 | 全链验收与循环入口 | Completed |

## 10. M16 完成条件

- 12 周 dataset 的 source owner、cutoff/cursor、unit/precision、missing/late/correction 和 hash 可重建。
- 数据质量先于 drift；五个外生域、process mismatch 和 policy-action 域有解释及失败路径。
- 8/4 split 无泄漏，candidate 唯一、有界、可复验且不修改 StrategyRelease/旧 evidence。
- holdout 和新 M14 replay 保留完整样本、common random numbers、ancestry 与结论限制。
- scheduled/expiry/alert/manual trigger、tenant/RLS、职责分离、幂等、事务/Outbox 和失败恢复通过。
- Observatory、两仓测试/部署、runbook/evidence、Atlas 和 M3-M15 回归完整。
- 输出 `strategy_assurance_cycle_closed=true` 与 `renewed|reexperiment_required|retired`，不存在未处理 finding、unknown drift 或未决 release effect。

## 11. M16 历史风险与控制

- M16 实施时在 G4-G8 冻结 observation、drift、校准/holdout、replay 和 IAOS gap 后才封存 dataset 并改变 release review 状态。
- canonical observation 仍来自虚构世界，不能包装成真实客户概率或外部统计规律。
- 数据质量问题必须先关闭；不得用校准掩盖 source 缺失、迟到、单位或 owner 错误。
- calibration owner 在候选冻结前不得访问 holdout；不得通过反复拟合、换指标或删除不利周过门。
- 新 replay 必须使用独立 ancestry/artifact，原 M14 seed/run/EvidenceBundle 不得覆盖。
- renewed 不表示永久有效；reexperiment_required/retired 同样是合法完成结果且不能自动修改 Policy。
- 当前工作区已有其他人的测试修改、截图删除和验收产物，实施 agent 必须保留并避免重叠修改。

## 12. M15 已完成范围

包含：

- M14 EvidenceBundle/hash 复验、Pareto candidate 审议、利益冲突与独立审批。
- immutable StrategyRelease、semantic diff、Policy/Capability 映射、SafetyEnvelope 和 preflight。
- 4 周或批准窗口的零写入 shadow、active/candidate decision diff 和 assumption drift。
- 新 Genesis canonical operating cycle 中最多 4 周的 allowlist pilot 和正式 IAOS 治理动作。
- hard stop、pause/review、kill switch、prior-release rollback、open commitments 和 compensating action。
- adopted/rejected/rolled_back AdoptionDecision、Strategy Control Room 和完整两仓证据。

不包含：

- 真实客户/生产租户投放、无人审批自动发布或 Agent sole approval。
- 第二客户/产品/工厂、完整 S&OP、动态定价或组织/资本重大变更。
- 删除既成事实式回滚，或用单一 pilot 宣称统计因果和永久最优。

## 13. M15 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| R0 | 决策、发布和安全机器合同 | Completed |
| R1 | Evidence Review 与 ChangeRequest | Completed |
| R2 | StrategyRelease 编译、差异与预检 | Completed |
| R3 | 零写入 Shadow 与运行兼容性 | Completed |
| R4 | Canonical Pilot 与受治理动作 | Completed |
| R5 | Guardrail、暂停、回滚与补偿 | Completed |
| R6 | IAOS 治理、决策与采纳 | Completed |
| R7 | Strategy Control Room 与全链验收 | Completed |

## 14. M15 完成条件

- candidate 与 exact M14 EvidenceBundle/hash 绑定，选择理由、限制、责任人和独立审批可审计。
- StrategyRelease、diff、SafetyEnvelope、RollbackPlan 和 AdoptionDecision 版本化、可 hash、可重放。
- shadow 完整窗口业务写入为零；candidate/active decision diff 和数据 freshness 可下钻。
- pilot 只在批准 scope/window 内通过 IAOS 治理动作，World consequence 与正式业务记录严格对账。
- hard stop、pause、kill switch、rollback、open commitment 和 compensation 失败路径通过。
- tenant/RLS、职责分离、幂等、事务/Outbox、恢复、Control Room、两仓测试及 M3-M14 回归完整。
- 最终输出 `strategy_change_cycle_closed=true` 与 `adopted|rejected|rolled_back`，不存在未处理 breach 或未对账 commitment。

## 15. M15 历史风险与控制

- M15 实施时在 G4-G8 冻结 candidate、release、shadow/pilot、guardrail/rollback 和 IAOS gap 后才激活 Policy 并创建 pilot 业务事实。
- M14 推荐仍只是模拟证据；不得为了推进计划预设 winner、隐藏限制或把 Pareto 误写成唯一最优。
- shadow 必须零业务写入，shadow approval 不能自动穿透 pilot approval。
- canonical pilot 会产生真实的虚拟企业业务后果；reset/rollback 不能删除已提交订单、库存、发运、发票或现金。
- kill switch 只停止未来动作；未结承诺必须进入 ledger 并用受治理 compensation 处理。
- 单一 pilot 无随机对照，不得宣称因果提升；超出 M14 assumption support 时必须 pause/review。
- 当前工作区已有其他人的测试修改、截图删除和验收产物，实施 agent 必须保留并避免重叠修改。

## 16. M14 已完成范围

包含：

- M12/M13/M14 批准 checkpoint allowlist、祖先/hash 和 opening reconciliation。
- 需求、供应、设备、质量和付款外生参数，以及库存/供应、产能/维护、资金保护策略。
- 固定版本 PRNG、命名随机流、seed set、共同随机数和 paired comparison。
- 12 个虚拟周/订单周期、分支隔离、持久 run catalog、有界执行、取消/继续/重试和配额。
- OTIF、积压、库存/营运资金、现金低点、毛利、质量、加班、报废、加急和恢复 KPI。
- Constraint/Pareto/EvidenceBundle、IAOS 实验治理与推荐边界，以及 Scenario Lab。

不包含：

- 自动把推荐应用到正式 Policy、预算、订单、采购、排产或现金。
- 用单次运行或未经校准的参数宣称真实因果、概率、预测精度或最优策略。
- 第二客户/产品/工厂、完整 S&OP、真实数据校准、机器学习训练或通用分布式计算平台。

## 17. M14 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| X0 | 实验方法、基线和机器合同 | Completed |
| X1 | 确定性随机流与参数矩阵 | Completed |
| X2 | Checkpoint fork、隔离分支与持久运行目录 | Completed |
| X3 | 多周期 World、策略执行和经济守恒 | Completed |
| X4 | 有界实验执行器与生命周期治理 | Completed |
| X5 | KPI 聚合、比较与 EvidenceBundle | Completed |
| X6 | IAOS 实验治理与推荐边界 | Completed |
| X7 | Scenario Lab 与全链验收 | Completed |

## 18. M14 完成条件

- 父 checkpoint 不变、兄弟分支/租户隔离，正式 IAOS 经营事实零污染。
- checkpoint、参数、策略、PRNG、seed、规则和聚合全部版本化、可 hash、可重放。
- 同输入 100 次 hash 一致，共同随机数和 paired comparison 可自动验证。
- 多周期数量、质量、资源、应收/现金和 actual cost/margin 守恒，失败样本不被过滤。
- 执行器默认 dry-run；显式 apply、有界并发、配额、取消、继续、重试和崩溃恢复通过。
- EvidenceBundle、IAOS 权限/职责/Outbox/RLS、Scenario Lab、runbook/evidence 和 M3-M13 回归完整。
- 只有证据完整、约束完成且无未解释运行缺失时输出 `strategy_evidence_ready=true`。

## 19. M14 历史风险与控制

- M14 实施时在 G4-G8 冻结 checkpoint、分布/seed、策略/KPI、运行容量和 IAOS gap 后才进入 IAOS 写端点开发。
- 场景分布是虚构假设，不能包装成真实概率；单次 run 只用于 tracer，不用于稳健性结论。
- 参数笛卡尔积可能失控；所有 apply 前必须预估 run 数、时间和存储并受配额约束。
- 不同策略必须使用共同随机数并保留失败/取消样本，否则比较存在选择偏差。
- 分支、transaction/correlation/idempotency namespace 必须隔离，不能污染父 checkpoint、兄弟分支或正式 IAOS 数据。
- 推荐与批准/投放必须分离；Agent、人或 UI 都不能从实验结果直接修改正式业务策略。
- 当前工作区已有其他人的测试、截图和生成物改动，实施 agent 必须保留并避免重叠修改。

## 20. M13 已完成范围

包含：

- 从零可销售成品库存接收 10,000 件正式订单和 2,000 件追加需求。
- M12 release manifest 到 Genesis-specific O2D 的兼容适配，以及 ATP/MRP、采购、IQC、生产、质量和库存。
- 供应延期、设备停机和末批 300 件短缺的受治理恢复。
- 9,000、2,700、300 三批实际发运、客户收货/接受和 delivery close。
- 开票、应收、客户实际付款、银行到账、收款核销和合同负债处理。
- 材料、人工、设备/能源、报废、质量、加急、物流和分摊实际成本，以及订单/项目毛利。
- First Delivery Play、Project Genesis M9-M13 terminal-hash 链和 M3-M12 强制回归。

不包含：

- 第二订单/客户/产品、多工厂或长期滚动 S&OP/MPS。
- 完整总账、税务申报、资金管理、真实银行/电子发票接口和法定报表。
- 售后、退货、质保、贷项、坏账和复杂跨期收入/成本会计。
- 参数化分支、Monte Carlo、A/B 和长期经营实验；属于 M14。

## 21. M13 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| E0 | 订单、履约、财务与成本机器合同 | Completed |
| E1 | M12 terminal 到 Genesis O2D 兼容适配 | Completed |
| E2 | 正式订单、ATP/MRP 与供应执行 | Completed |
| E3 | 正式生产、质量、库存与实际成本 | Completed |
| E4 | 分批发运、客户接受与 300 件恢复 | Completed |
| E5 | 开票、应收、银行回款与项目盈亏 | Completed |
| E6 | IAOS O2D、财务与成本治理 | Completed |
| E7 | 统一 Agent 与 First Delivery Play | Completed |
| E8 | 全链验收与 Project Genesis 收口 | Completed |

## 22. M13 完成条件

- 从 M12 terminal contract 到 `first_commercial_cycle_closed` 可确定性运行、恢复、重放和分层复位。
- 正式需求 12,000、采购/生产/报废/库存、三批发运和客户接受数量严格守恒。
- 300 件短缺形成 observation、受治理恢复、第三批交付和 discrepancy close 完整链。
- 客户接受、发票/应收、银行到账、收款核销、收入和现金严格分离且金额守恒。
- M12 财务结转、订单实际成本、标准/实际差异、项目毛利和 closing cash 可解释对账。
- 两仓权限、RLS、Outbox、幂等、API/UI、runbook/evidence、Project Genesis 总报告以及 M3-M12 回归全部通过。

## 23. M13 历史风险与控制

- M12 1,200,000 CNY 试制成本的支付/应付状态和 2,000,000 CNY 合同负债履约处理已在 G4 冻结；全程禁止通过改 opening cash/利润静默配平。
- 旧 HCTM 场景包含 1,200 件 opening inventory 和已发生事件，M13 只能复用语义/能力，不能继承库存或交易历史。
- 发运不等于客户接受，发票不等于现金，毛利不等于现金余额；UI 和 IAOS 状态不能伪造 World/银行事实。
- 实施时在 G4-G9 关闭后才进入 IAOS 写端点开发，E6 使用了独立 IAOS branch/worktree。
- actual cost 缺任一强制要素时经营分析必须失败或保持不完整，不能回退到估算后宣称盈利。
- M13 reset 必须保留 M9-M12 L1 事实和旧场景数据，只清理本次运行的 L2/L3 对象。
- 当前工作区已有其他人的测试、截图和生成物改动，实施 agent 必须保留并避免重叠修改。

## 24. M12 已完成范围

包含：

- 单一虚构客户的 RFQ、可行性、成本/报价、客户定点、开发资金和客户项目。
- `HCTM-BCP-A01` 产品 revision、EBOM/MBOM、routing、process flow、PFMEA、control plan 和 work instruction。
- 关键供应商、产品专用工装/量检具、首批材料、批次/证书、IQC 和追溯。
- 简化 APQP 六个 gate、两轮试生产、尺寸/功能/泄漏、MSA、节拍、良率和 Cp/Cpk。
- 一条焊接泄漏/Cpk 不足的 discrepancy、Knowledge、工程变更、复试和 PPAP tracer。
- Industrialization Play、旧 HCTM stable-code/hash 兼容，以及 M7-M11 强制回归。

不包含：

- 第一张正式销售订单、正式 MRP、批量采购和正式量产；属于 M13。
- 客户发运、发票、应收、回款、实际量产成本和项目盈亏；属于 M13。
- 第二客户/产品、多工厂、多版本并行量产或完整 CRM/CPQ/PLM/QMS/SRM。
- CAD/CAE、高精度物理仿真、真实外部接口或参数化分支实验。

## 25. M12 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| D0 | 客户项目、产品、质量与资金机器合同 | Completed |
| D1 | RFQ、报价、定点与项目资金 | Completed |
| D2 | 产品、工艺与 APQP 版本治理 | Completed |
| D3 | 供应商、工装与首批物料世界 | Completed |
| D4 | 试生产、质量异常与 PPAP 世界 | Completed |
| D5 | IAOS 客户项目、工程、质量与放行治理 | Completed |
| D6 | 统一岗位与 Industrialization Play | Completed |
| D7 | 全链、安全、恢复和回归验收 | Completed |

## 26. M12 完成条件

- 从 M11 terminal contract 到 `serial_production_eligible` 可确定性运行、恢复、重放和安全复位。
- RFQ/报价/定点、开发资金/合同负债、项目预算、试制成本和现金语义正确且守恒。
- 产品/BOM/routing/PFMEA/control-plan/工装/供应商/材料 revision 一致并可机器校验。
- 试制物料、样件、设备、人员、测量、报废、成本和追溯满足守恒。
- 焊接泄漏/Cpk 异常形成 observation、受治理变更、第二轮试制、问题关闭和客户 PPAP 批准完整链。
- 两仓权限、RLS、Outbox、幂等、API/UI、runbook/evidence、M3/O2D 兼容以及 M7-M11 回归全部通过。

## 27. M12 风险与依赖

- M11 closing cash 为 8,500,000 CNY，并需保留 3,000,000 工资准备金和 5,000,000 最低缓冲；G5 必须冻结客户工装预付款/开发预算，且预付款只能记合同负债，不能计收入。
- 现有 `scenario-packs/hctm` 已预置同名产品/BOM/routing，但只是兼容 fixture；M12 必须产生独立 release manifest/hash，禁止把旧 seed 当完成证据。
- IAOS 工程/APQP/PPAP 状态不等于实际试制能力或客户批准；只有 World consequence 和客户 observation 可推进现实状态。
- G4-G8 未关闭前不得进入 IAOS 写端点开发；D5 必须建立新的独立 IAOS branch/worktree。
- 试制件、PPAP 样件和报废默认不可销售，不能静默成为 M13 库存或收入。
- M12 不得提前接收正式订单、运行正式 O2D 或宣称第一批交付完成。
- 当前工作区已有其他人的测试、截图和生成物改动，实施 agent 必须保留并避免重叠修改。

## 28. M11 已完成范围

包含：

- M9 剩余认缴资本的受治理催缴/实际到账、新 CAPEX/headcount envelope、设施尾款、工资准备金和现金缓冲。
- 冷却板 A 线通用能力需求，以及成形、CNC、激光焊接、清洗、检漏、装配/包装、质量实验室的采购/租赁、交付、安装、调试、校准和验收。
- M10 空间/utility 约束下的设备落位、实验室、仓储和最小一班制能力。
- 厂长、计划、采购、质量、工艺、设备、仓储、操作工和检验员的编制、招聘、到岗、培训、资格和班次。
- 一条检漏设备校准漂移的 discrepancy、Knowledge、整改、复验和关闭 tracer。
- Capability Build Play，以及 M7-M10 强制回归。

不包含：

- 客户 RFQ、正式产品/BOM/routing、成本报价和 APQP；属于 M12。
- 产品专用工装、首批物料、试生产、PPAP、SOP 和量产；属于 M12。
- 第一张正式订单、开票、回款和实际项目盈亏；属于 M13。
- 完整采购/SRM、HRIS/LMS、薪酬、EAM、固定资产会计、融资或 3D 产品。

## 29. M11 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| C0 | 能力、设备、人员与资金机器合同 | Completed |
| C1 | 资金与受治理设备采购 | Completed |
| C2 | AESE 设备、实验室与仓储能力世界 | Completed |
| C3 | AESE 人员、技能与班次世界 | Completed |
| C4 | IAOS 采购、资产、组织与资格治理 | Completed |
| C5 | 统一岗位与 Capability Build Play | Completed |
| C6 | 全链、安全、恢复和回归验收 | Completed |

## 30. M11 完成条件

- 从 M10 terminal contract 到 `industrialization_eligible` 可确定性运行、恢复、重放和安全复位。
- 资本实际到账、新预算、设施尾款、设备/租赁承诺、工资准备金和现金严格守恒。
- 关键设备全部实际交付、安装、调试、校准、安全和能力验收，空间/utility 不超限。
- 最小核心团队实际到岗且岗位资格有效，一班制、替补和职责隔离门通过。
- 检漏设备失败形成 observation、Knowledge 差异、受治理整改、复验和关闭的完整因果链。
- 两仓权限、RLS、隐私、Outbox、幂等、API/UI、runbook/evidence 以及 M7-M10 回归全部通过。

## 31. M11 风险与依赖

- M10 closing cash 只有 10,000,000 CNY 且存在设施遗留承诺，不能无资金来源生成整线；G6 必须先冻结剩余资本实缴、采购/租赁组合和现金缓冲。
- 资产登记、员工档案和培训记录不等于实际设备能力、人员到岗或技能掌握；只有 World consequence 和验收事实可推进能力状态。
- M8 `LAS-WLD-02` 是独立设备 tracer，未经采购/交付/验收迁移证据不能计入 M11 生产资产。
- G4-G7 未关闭前不得进入 IAOS 写端点开发；C4 必须建立新的独立 IAOS branch/worktree。
- 候选人与员工完全虚构且按 actor scope 投影；不得引入真实个人信息或让无关角色读取候选详情。
- M11 只交付通用设备和人员能力，不得提前宣称产品、APQP、试生产、PPAP 或 SOP 完成。
- 当前工作区已有其他人的测试、截图和生成物改动，实施 agent 必须保留并避免重叠修改。

## 32. M10 已完成范围

包含：

- 至少三个虚构候选场址及资金、工期、物流、人力、公用工程、风险和扩展性评估。
- 项目负责人推荐、CEO/CFO 受治理的选址与投资批准，以及租赁/场地使用控制。
- 设施项目、WBS、承包商资源日历、合同承诺、变更、里程碑、付款和验收。
- 区域/城市/园区/场地/建筑/楼层/功能区的最小空间层级和公用工程容量。
- 一个固定公用工程接入延期的 discrepancy、Knowledge、重排和关闭 tracer。
- World Play 工厂建设 campaign，以及 M7/M8/M9 强制回归。

不包含：

- 生产设备和检测仪器采购、安装及调试；属于 M11。
- 招聘、培训、工艺能力、试生产和投产门；属于 M11。
- BIM、3D、自由布局编辑器、真实地图/园区/承包商接口。
- 完整地产、总账、税务、融资或通用项目管理产品。

## 33. M10 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| P0 | 场址、空间与项目机器合同 | Completed |
| P1 | 受约束选址决策 | Completed |
| P2 | AESE 设施项目与空间世界 | Completed |
| P3 | IAOS 投资与项目治理 | Completed |
| P4 | 统一角色与 Plant Build Play | Completed |
| P5 | 全链、安全、恢复和回归验收 | Completed |

## 34. M10 完成条件

- 至少三个候选先经过硬约束，再产生版本化、可解释的多维评分和受治理决策。
- 从 M9 terminal contract 到 `capability_build_eligible` 可确定性运行、恢复、重放和安全复位。
- IAOS 项目记录与 AESE 现场事实严格分离；时间到期或管理记录不能凭空完成工程。
- 预算、合同承诺、实际付款和现金分别守恒，未验收、越权或超预算动作失败关闭。
- 公用工程延期能形成 observation、Knowledge 差异、受治理重排和新实际结果的完整因果链。
- 两仓权限、RLS、Outbox、幂等、API/UI、runbook/evidence 以及 M7/M8/M9 回归全部通过。

## 35. M10 风险与依赖

- M9 只有 20,000,000 CNY 实际现金和 15,000,000 CNY 首年预算，首版不得假设可支付绿地自建；候选基线预计收敛到租赁标准厂房改造，但必须由规则和审批得出。
- IAOS project/milestone 状态不等于现场实际进度；只有 AESE World consequence 和验收事实可以推进物理状态。
- G4/G5 已关闭；P3 在独立 `feat/m10-plant-governance` worktree 完成，revision `23be02a`。
- 承包商、公用工程方、园区和验收机构是确定性外部世界策略，不是 IAOS 用户，也不能绕过 observation/intent/outcome 合同。
- M10 只交付设施载体，不得提前引入生产设备、人员或投产能力，避免侵入 M11。
- 当前工作区已有其他人的测试、截图和生成物改动，实施 agent 必须保留并避免重叠修改。

## 36. M9 已交付范围（M9N 平台化验收完成）

以下业务范围与 D19–D21 平台化证据共同构成 M9N 完成结论。M9N 当前状态以顶部里程碑表
和 `PLAN-M9-NATIVE-001` 为准；专用闭环证据不能脱离通用资产与可解释配置验收单独使用。

包含：

- 单一虚构投资主体、集团和苏州制造法人。
- 设立方案、外部登记、开户、首期资本到账和成立费用。
- CEO、CFO、工厂项目负责人任命、接受、mandate 和知识边界。
- 初始组织、首年 budget envelope、审批和可支用资格。
- World/IAOS/Knowledge 三态、资金守恒和受治理 bridge 全链。
- World Play 企业成立 campaign 与 M7/M8 强制回归。

不包含：

- 工厂选址、土地/租赁和建设执行；属于 M10。
- 设备采购、招聘培训和投产门；属于 M11。
- RFQ/APQP/PPAP、首批交付和参数化实验。
- 完整总账、税务、复杂融资或真实监管/银行接口。

## 37. M9 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| I0 | 业务基线、所有权和机器合同 | Completed |
| I1 | AESE 成立世界、外部策略与经济规则 | Completed |
| I2 | IAOS 法人、治理、岗位和预算能力 | Active remediation |
| I3 | CEO/CFO 统一岗位运行 | Active remediation |
| I4 | Genesis Incorporation Play | Completed |
| I5 | 全链、恢复、安全和回归验收 | Active remediation |

## 38. M9 完成条件

- pre-incorporation 到 `plant_project_eligible` 可确定性运行、恢复和重放。
- 法人、账户、资本、任命和预算的三态及因果链可机器验证。
- 投资人/公司现金、认缴、实缴、预算和实付严格区分并满足守恒。
- 人类与 Agent 使用相同 IAOS Capability、Policy、权限和审计。
- IAOS journal/Outbox 与业务提交原子，回滚不产生 committed outcome。
- tenant、幂等、并发、断线和重复执行测试通过，M7/M8 零回归。

## 39. M9 风险与依赖

- M8 初态中的 10,000,000 CNY 当前没有显式 owner；M9 必须通过 pack version 和资本事件迁移，禁止静默当作公司现金。
- 法人档案 committed outcome 不等于外部登记已经生效；监管和银行结果由 AESE 确定性世界策略产生。
- 预算是支出授权，不是现金；认缴不是实缴；管理界面必须明确区分。
- M9 IAOS 修改必须新建独立 worktree，并在 I0 合同冻结后启动。
- 自动 Agent 不得批准自身预算、伪造外部结果或绕过岗位 mandate。
- 当前工作区已有其他人的测试与截图改动，实施 agent 必须保留并使用不重叠 worktree。

## 40. M8 已完成范围

包含：

- World run、虚拟时钟、离散事件、规则版本、日志、快照和确定性 replay。
- World / IAOS / Actor Knowledge 三态与显式 discrepancy。
- 单台 `LAS-WLD-02` 设备退化的最小发现、登记、处置和关闭 tracer。
- observation / intent / committed outcome 的受治理 IAOS 合同。
- 设备、人员、物料和最小现金守恒。
- `hctm-genesis` world pack 和现有 `order-expedite-01` 兼容适配。
- World Play 的时间控制、三态对照和差异时间线。

不包含：

- 一次完成公司设立、工厂建设、APQP、财务和全生命周期的所有模块。
- 多工厂、多产品族、数百 Agent、自由文本长期记忆和自主批准。
- 3D、高精度物理仿真或完整设备数字孪生。
- 绕过 IAOS 权限/Capability/Outbox 的业务写入。
- 在 AESE 中镜像 IAOS 企业管理数据库。

M8 决策门与 F0-F5 的任务、验收和跨仓顺序以 PLAN-M8-001 为准。后续 Project Genesis 分解为 M9-M13，参数化分支实验后移至 M14。

## 41. M8 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| F0 | 基线、状态所有权、存储和桥接合同冻结 | Completed |
| F1 | 确定性仿真内核 | Completed |
| F2 | 三态与设备偏差 tracer | Completed |
| F3 | 受治理 IAOS 双向桥 | Completed |
| F4 | Genesis world pack 与旧场景兼容 | Completed |
| F5 | World Play 最小界面与全链验收 | Completed |

## 42. M8 架构风险与依赖

- ADR-004 已 accepted；实现必须遵守独立 PostgreSQL database/账号/迁移边界，禁止跨库查询和外键。
- AESE 仿真事实和 IAOS 管理事实必须物理/逻辑隔离，禁止共享表和跨库写入。
- 原始设计稿中的 Spring Boot 仅是模块示意；AESE 实现继续使用 Go，与现有工具链保持一致。
- 世界结果必须由版本化规则和资源守恒计算；Agent 只提交 intent，不能直接改写 World State。
- Actor Knowledge 必须遵守可见范围，不能为了方便让 Agent 读取全量客观世界。
- IAOS 修改必须在独立 worktree，并先完成权限、RLS、Outbox、幂等和无部分写入设计。
- M7 22 事件、三 Agent、Preview/Live 与 reset 是强制回归门。

## 43. M8 完成条件

- ADR-004 accepted，World/IAOS/Knowledge 所有权和 World Store 选型明确。
- 相同 pack、规则版本、seed 和输入可重复产生相同 event log、state hash 与 KPI。
- 设备退化 tracer 可展示世界变化、IAOS 未登记、角色未知及其发现/关闭过程。
- IAOS 双向桥通过租户、权限、幂等、乱序、失败恢复和 Outbox 审计验收。
- Genesis pack 可离线验证、初始化、推进、复位和 replay，旧 M7 场景不回归。
- API/UI/runbook/evidence 与两仓 revision 完整。

## 44. M7 已完成范围（保留基线）

包含：

- 无独立数据库的 AESE scenario orchestration API。
- pack 阶段编译、dry-run 影响、状态机、幂等和恢复。
- initialize、advance、run-to-end、analyze、verify 和 reset 编排。
- 联动中心场景运行视图、七幕 stepper、日志和危险确认。
- 权限、跨租户、并发、断线、重启及 CLI/UI 一致性验收。

不包含：

- 参数化 A/B 实验、并行分支和第二条故事。
- 真实 LLM 自主 Agent 和建议自动执行。
- AESE 业务数据库、通用任务队列或工作流引擎。
- 完整成本核算、3D 工厂和布局编辑器。

## 45. M7 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| O0 | 状态机、阶段合同和编排内核重构 | Completed |
| O1 | AESE 薄编排 API | Completed |
| O2 | IAOS 运行记录、权限和并发补强 | Completed |
| O3 | 可视化场景运行控制台 | Completed |
| O4 | 全链路、恢复、安全和三视口验收 | Completed |

## 46. 历史风险与依赖

- IAOS Platform、PostgreSQL、NATS 和 O2D 可运行，`tenant-hctm` 的 work_order metadata、workflow config 和 tracer 数据已完成 seed/apply。
- `/iaos/iaos-go` 当前主分支本地领先远程，任何集成开发必须使用独立 worktree 并确认基线。
- HCTM 业务字段与 IAOS 现有 legacy `sales_order` 物理模型存在差异，需要兼容性报告，不能直接假设可导入。
- IAOS 当前 O2D 自测硬编码 `tenant-001` 和旧订单，需要避免污染 HCTM tracer。
- 正式事件必须走 Outbox 或受治理 ingress，不能把 direct NATS 当最终实现。
- IAOS 受治理 scenario apply/reset 已实现 M3 allowlist；订单确认 CAS、workflow/event 去重、跨节点原子事务及 work_order API 已实证。
- legacy 表没有全面 FORCE RLS；M3 scenario adapter 已在所有查询/更新/删除中显式绑定 tenant。平台长期仍应继续推进全表 RLS FORCE hardening。
- M4 的显式 tenant predicates 已关闭当前入口越界路径，但 tenant-safe composite foreign key 和 metadata `version` 的平台级排序仍需后续 hardening。
- 首版 2D 沙盘是确定性预览，不应被描述为 IAOS 实时运行结果；界面必须显示 Preview 数据源状态。
- `preview.json` 只承载视图状态和 delta，不得复制 MRP、排产或 Agent 决策逻辑。
- 当前通用 `/api/v1/events/stream` 无持久 cursor 且缓冲满会丢事件，只能作为监控通道，不能直接作为 M6 恢复合同。
- 完工和发运是内生业务动作，不能为了复用现有接口而错误接入 simulation ingress。
- HCTM 尚无已批准的成本金额基线；在基线确认前，在线经营分析的成本部分必须保持 `partial`。
- M7 新增 HTTP 服务但不得拥有业务数据库；运行恢复必须以 IAOS run/snapshot/event/recommendation 为事实。
- 当前场景使用固定自然键，同一 tenant/scenario 首版只能有一个可写 active run。
- 浏览器不得直接编排多个 IAOS 写 API；所有危险动作需要服务端权限、幂等和确认合同。

## 47. M7 完成条件

- 非研发用户可从浏览器完整运行并复位第一条故事。
- UI 状态只在 IAOS committed/no-op 和 snapshot cursor 证实后推进。
- 双击、并发、断线、刷新和服务重启均不重复推进阶段。
- reset 影响可预览，一次性 confirmation token 不能重放，L1 始终保留。
- IAOS 与 AESE 两仓权限、测试、部署、runbook 和 evidence 完整。

## 48. M6 完成证据

M6 已满足：

- 在线库存、完工、发运和订单状态满足 expected outcomes 与库存守恒。
- snapshot 与 cursor 来自一致观察边界，断线补发和事件去重可复现。
- Preview/Live 数据源明示，错误时不静默混用。
- Agent 建议、对象引用和 Tool Call 证据属于同一 tenant/correlation。
- IAOS 与 AESE 两仓测试、部署、runbook 和 evidence 完整。
