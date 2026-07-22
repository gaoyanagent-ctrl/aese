---
id: PLAN-M17-001
title: M17 滚动 IBP 与 S&OP 实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m17, ibp, sop, planning]
---

# M17 滚动 IBP 与 S&OP 实施计划

## 1. 交付目标与工期

用 **8 到 10 周**把 M16 renewed 策略与校准假设转化为第一个受治理滚动经营计划周期：

> 在单工厂、单产品、单客户范围内，以 13 周 weekly execution horizon 和 12 个月 monthly financial horizon 对齐需求、供应、产能、库存、交付、成本和现金，经 Demand Review、Supply Review、Financial Reconciliation、Pre-IBP 和 Executive IBP 形成唯一 PlanRelease，输出 `integrated_plan_cycle_closed=true`。

终态 disposition 为 `approved|replan_required|deferred`。批准只发布计划版本，不自动产生订单、采购、工单、发运或资金动作。

## 2. 前置门

- [x] G1 M16 主路径 `renewed`，12 周 dataset、校准/holdout、60-run replay 和 resilient support evidence 完整。
- [x] G2 DES-018 冻结 M17-M24 program 顺序，DES-019 冻结 M17 三态、计划/现实分离和完成边界。
- [x] G3 ADR-004、DES-008、DES-009 及 M14-M16 evidence/decision ancestry 继续有效。
- [x] G4 冻结 13 周/12 月 horizon、calendar/bucket、cutoff、frozen/slushy/liquid fence、单位/精度和 opening balances。
- [x] G5 冻结 demand/forecast/order、supply/capacity/inventory、cost/cash/budget 的 owner、公式和 source refs。
- [x] G6 冻结 baseline/upside/downside assumptions、scenario inheritance、gap/option、service/cash/quality hard constraints。
- [x] G7 冻结五个 review gate、角色/职责、version/CAS、change reason、PlanRelease 和 replan trigger。
- [x] G8 完成 IAOS planning/forecast/policy/decision gap audit，冻结 payload、权限、事务、Outbox 和跨仓顺序。

G4-G8 关闭前只允许 Schema、fixture、calendar/decimal、纯计划算法、只读 gap audit 和离线 prototype；不得创建 IAOS PlanRelease 或执行任何业务动作。

## 3. 实施切片

### B0 - Planning Cycle 与集成计划机器合同（第 1-2 周）

- [x] T1 定义 PlanningCycle、AssumptionSet、DemandPlan、SupplyPlan、CapacityPlan、InventoryPlan、FinancialPlan、Gap、Option、Decision 和 PlanRelease stable code/owner。
- [x] T2 冻结 horizon、bucket/calendar、cutoff/as-of、version、fence、source actual 和 opening reconciliation。
- [x] T3 固定 forecast/order/commitment/actual、available/capacity、budget/cash/profit 的语义分离。
- [x] T4 固定 quantity/unit/precision、currency/rounding、capacity/time 和 aggregation/disaggregation 规则。
- [x] T5 固定 baseline/upside/downside assumption ancestry，绑定 M16 dataset/calibration/decision hash。
- [x] T6 固定 service、quality、cash、capacity、inventory 和 mandate hard constraint 与 gap severity。
- [x] T7 定义 `draft -> demand_reviewed -> supply_reviewed -> financially_reconciled -> pre_ibp -> approved|replan_required|deferred -> closed` 状态机。
- [x] T8 扩展 JSON Schema、Go strict types/parser、canonical hash、payload registry 及合法/边界/破损 fixture。
- [x] T9 完成权限/职责、IAOS gap、Atlas 依赖、跨仓矩阵和 G4-G8 评审。

验收：每个计划数值有 bucket、owner、assumption/source 和 version；任何未对齐单位、时间或 owner 的计划失败关闭。

### B1 - Demand Plan、订单与假设（第 2-4 周）

- [x] T10 从 M16 renewed assumption、M13-M16 actual 和客户 observation 编译 frozen DemandInputSnapshot。
- [x] T11 实现 forecast、confirmed order、promotion/one-off、backlog 和 customer commitment 的独立 bucket。
- [x] T12 实现 baseline/upside/downside 需求曲线、assumption diff、confidence/limitation 和 source freshness。
- [x] T13 实现 weekly 13 周到 monthly 12 月 aggregation，禁止重复累计或跨 cutoff 混用。
- [x] T14 实现销售/计划协作、change reason、comment/evidence 和 optimistic concurrency。
- [x] T15 实现 Demand Review gate、独立 approver、unresolved gap 和 stale input 失败关闭。
- [x] T16 验证 forecast 不变订单、客户意向不变承诺、重复订单/跨租户/陈旧版本和越权修改失败关闭。
- [x] T17 用 100 次 build 和乱序输入验证 DemandPlan/hash 稳定。

验收：需求计划说明已确认需求与假设需求，并能追溯到 evidence；不会凭计划创造客户订单。

### B2 - Supply、Capacity、Inventory 与交付计划（第 3-5 周）

- [x] T18 读取已发布 BOM/routing、库存/在途、供应周期、设备/人员/班次能力和当前 StrategyRelease refs。
- [x] T19 实现 rough-cut material/supply plan，区分 planned receipt、PO commitment 和 World actual arrival。
- [x] T20 实现设备/人员/班次/维护/质量门共同约束的 rough-cut capacity plan。
- [x] T21 实现 RM/WIP/FG inventory projection、安全库存、backlog、allocation 和 ATP/CTP 计划语义。
- [x] T22 实现 frozen/slushy/liquid fence 对 change、expedite、overtime 和 reschedule option 的限制。
- [x] T23 实现 delivery/service plan 与 demand/supply/capacity/inventory 数量守恒。
- [x] T24 实现供应/产能/库存 Gap、root constraint、option、trade-off 和 owner。
- [x] T25 实现 Supply Review gate；缺材料、资格、维护窗口、质量 release 或 cash support 时失败关闭。
- [x] T26 验证超产能、负库存、计划收货即实际到货、重复 allocation、跨产品引用和越权 approval 失败关闭。

验收：供给计划只表达未来可行性和选择，不创建现实材料、产能、生产或交付。

### B3 - Financial Reconciliation 与经营价值（第 4-6 周）

- [x] T27 读取 M13-M16 actual cost/rate、price/terms、AR/AP/cash/commitment refs 和管理口径 opening。
- [x] T28 实现 volume/mix/price、材料/人工/设备/质量/物流成本和 gross margin plan。
- [x] T29 实现 inventory/AR/AP/contract liability、working capital、cash-in/out 和 cash trough plan。
- [x] T30 区分 budget、forecast、commitment、actual、cash 和 profit；禁止现金/利润或发票/收款混同。
- [x] T31 实现三个 scenario 的 P&L/cash/working-capital bridge 和与 operational plan 的数量金额 reconciliation。
- [x] T32 实现最低现金、信用、采购/加班预算和 margin guardrail 及 funding option。
- [x] T33 实现 Financial Reconciliation gate、CFO review、unexplained variance 和 missing actual failure。
- [x] T34 验证重复收入/成本、浮点舍入、伪造银行事实、超预算、自批和错币种失败关闭。

验收：每个 operational option 的服务、成本、利润、现金和资金风险可解释；计划财务不冒充会计事实。

### B4 - Gap、Scenario、Pre-IBP 与 Executive Decision（第 5-7 周）

- [x] T35 编译 cross-functional Gap Register，合并但不掩盖 demand/supply/quality/cash owner 和 severity。
- [x] T36 实现 baseline/upside/downside PlanBundle，各 scenario 共享固定事实并只覆盖显式 assumption。
- [x] T37 实现 option impact：服务、库存、产能、质量、成本、现金、风险和不可逆 commitment。
- [x] T38 实现 Pareto/constraint view，默认不以隐藏权重输出单一最优方案。
- [x] T39 实现 Pre-IBP gate，要求所有 hard gap 有 option/owner/decision request 或明确 deferred reason。
- [x] T40 实现 Executive IBP 决策：approve option、request replan 或 defer，并保留 dissent/limitation。
- [x] T41 实现 immutable PlanRelease、exact component hashes、effective cycle、review date 和 next-cycle carryover。
- [x] T42 实现 replan trigger：重大需求/供应/质量/现金变化、M16 drift 或陈旧 actual；不得静默改 release。
- [x] T43 验证跳 gate、自批、陈旧 component、并发 approve、hidden gap、PlanRelease 自动执行业务动作失败关闭。

验收：管理层看到的是一致又可下钻的计划与权衡，批准后仍需独立执行治理。

### B5 - IAOS Planning 治理与 Bridge（第 4-8 周）

- [x] T44 按 IAOS 规则建立新的独立 branch/worktree，先读其 AGENTS、Agent Context 和 Code Map。
- [x] T45 复用 metadata/Policy/Decision/Process/Capability/AI Tool，只增加 planning cycle/review/gap/option/release 最小能力。
- [x] T46 实现 M17 allowlist observation/intent/outcome payload、journal/cursor 和 exact plan/component hash。
- [x] T47 固定 `genesis.ibp.create/read/edit/review/reconcile/preibp/approve/replan/release` 权限和职责分离。
- [x] T48 保证 plan version、review/decision、audit、journal 和 Outbox 必要事务原子性及 CAS/幂等。
- [x] T49 实现 PlanRelease 与 execution intent 强隔离；任何 PO/WO/Shipment/资金动作必须另行授权。
- [x] T50 验证 tenant/RLS、跨职能越权、自批、陈旧/重复 release、并发 review 和失败无部分写入。
- [x] T51 两仓分别提交、记录 revision、部署并完成 contract/integration tests。

验收：IAOS 拥有计划/审批，AESE 提供 World/assumption/scenario 后果；两边都不能用 plan 状态伪造现实。

### B6 - Executive IBP Room（第 6-9 周）

- [x] T52 在 World Play 增加 Executive IBP Room：cycle、horizon、fence、assumption、version 和五个 gate。
- [x] T53 展示 demand/supply/capacity/inventory/financial plan、actual、gap 和 owner 的 bucket 对齐。
- [x] T54 展示三个 scenario、option、constraint/Pareto、服务/现金/利润权衡和限制。
- [x] T55 展示 World/IAOS/Knowledge owner、forecast/order/actual、plan/release/execution 的状态差异。
- [x] T56 提供 review/replan/approve impact preview、exact hash 和独立审批状态，不提供直接业务执行入口。
- [x] T57 完成键盘、ARIA、焦点、移动端、非颜色 gap、单位/精度、weekly/monthly zoom 和大表性能。
- [x] T58 验证刷新、断线、陈旧 cursor/version、部分 unavailable、并发编辑和错误空态。

验收：非研发用户可完成一个 IBP cycle，并解释所有数字、冲突、选择和未自动执行边界。

### B7 - 全链验收与 M18 入口（第 8-10 周）

- [x] T59 从 M16 renewed terminal 执行 approved、injected replan_required 和 deferred 三条完整路径。
- [x] T60 对账 assumption/dataset、plans、gaps/options、decisions、journal/Outbox、World refs、Knowledge 和 hashes。
- [x] T61 验证 cutoff 后迟到 actual、重大需求/供应/质量/现金变化触发 replan 而不改旧 release。
- [x] T62 验证服务重启、SSE 丢失、cursor 恢复、重复点击、陈旧 version、并发 cycle 和分层 reset。
- [x] T63 验证 tenant-other、只读用户、Agent 自批、隐藏 gap、自动执行、直接 DB/NATS 和 real-production target 拒绝。
- [x] T64 运行 Go、Schema、PostgreSQL、IAOS modules、前端、Playwright、100 次 hash、容量和 M3-M16 回归。
- [x] T65 编写 M17 runbook/evidence，更新设计、计划、Roadmap、Code Map、Progress Log、Atlas、两仓 revision 和部署信息。
- [x] T66 发布 M18 machine handoff：approved PlanRelease、shared-capacity/portfolio extension points 和未改变的单产品 regression baseline。

验收：只有五个 review gate、跨职能 reconciliation、hard gap、PlanRelease 和未执行边界全部可对账时输出 `integrated_plan_cycle_closed=true`。

## 4. 首版业务方向

| 维度 | 首版范围 |
| --- | --- |
| Scope | 苏州基地 / `HCTM-BCP-A01` / 当前客户 / resilient release |
| Execution horizon | 13 周 weekly |
| Financial horizon | 12 个月 monthly |
| Scenarios | baseline / upside / downside |
| Gates | Demand / Supply / Finance / Pre-IBP / Executive IBP |
| Hard constraints | quality、service、capacity、inventory、cash、mandate |
| Terminal | approved / replan_required / deferred |

所有预测、价格、成本和经营数据保持虚构。cycle/plan/gap/option/release 使用 `IBP-GENESIS-M17-*` 稳定编码。

## 5. 完成定义

- horizon/calendar/fence/cutoff/unit/version 和 opening reconciliation 完整。
- demand/supply/capacity/inventory/financial plan 跨 bucket 数量金额守恒。
- 五 gate、职责分离、gap/option/decision 和 exact PlanRelease hash 可审计。
- 计划、actual、commitment 和 execution 分离，approved plan 不自动产生业务副作用。
- tenant/RLS、CAS/幂等、事务/Outbox、恢复、replan 和失败无部分写入通过。
- IBP Room、CLI/API、runbook/evidence、三视口、两仓 revision/部署、Atlas 和 M3-M16 回归完整。
- 输出 `integrated_plan_cycle_closed=true` 与 approved/replan_required/deferred。

## 6. 不纳入 M17

- 第二产品/客户/工厂、多基地网络和完整 portfolio optimization。
- 真实预测模型、自动求解器、计划自动执行和 real-production target。
- 法定财务、税务、真实银行、完整 S&OP SaaS 或无限 horizon。

## 7. 并行与所有者规则

- B0 串行冻结；B1/B2 可在合同稳定后并行，但 bucket/calendar/version 由单一 owner 合并。
- B3 依赖 B1/B2 数量口径；现金/利润/成本由单一 financial owner 对账。
- B4 只消费冻结 component versions，不允许为达成 Executive 决策回写历史 component。
- B5 由独立 IAOS worktree owner 开发，两仓提交/测试/日志分别维护。
- B6 在 API/view model 稳定后开始，UI 不得执行 PO/WO/Shipment/资金动作。
- B7 串行核对三条 disposition、两仓、全部版本、revision、部署和 Atlas。
- 保留共享工作区中现有测试修改、截图删除和验收产物，不得覆盖或回滚。
