---
id: PLAN-M13-001
title: M13 Project Genesis 第一次完整商业交付实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m13, genesis, o2d, delivery, invoice, cash, actual-cost]
---

# M13 Project Genesis 第一次完整商业交付实施计划

## 1. 交付目标与工期

用 **8 到 10 周**把 M12 的 `serial_production_eligible=true` 推进为：

> 华辰苏州制造公司从零可销售库存接收首张 10,000 件订单和 2,000 件追加需求，经 MRP、采购、生产、质量和三批交付完成累计 12,000 件客户接受，随后完成开票、应收、实际回款、订单实际成本和项目毛利对账，输出 `first_commercial_cycle_closed=true`。

首版固定单客户、单产品、单订单族、单工厂/产线、单币种 CNY、三批交付和一条 300 件末批短缺恢复 tracer。M13 完成 Project Genesis M9-M13 主纵向场景；多周期/参数化经营属于 M14。

## 2. 前置门

- [x] G1 M12 已完成并机器输出 `serial_production_eligible=true`。
- [x] G2 ADR-004、DES-008、DES-009 与 M8-M12 的 World/Bridge/Knowledge/资金边界继续有效。
- [x] G3 DES-014 固定 M13 三态所有权、O2D/财务边界、旧场景兼容和 Genesis 完成条件。
- [x] G4 冻结 M12 closing cash、1,200,000 试制成本支付/应付状态、2,000,000 合同负债履约处理、遗留承诺和 opening balance migration。
- [x] G5 冻结正式订单/追加、单价、税率、付款条件、交付批次、客户接受和银行到账基线。
- [x] G6 冻结 BOM/MRP、采购/来检、产能/班次、生产/报废、库存/发运和 300 件恢复时间线。
- [x] G7 冻结材料、人工、设备/能源、报废、质量、加急、物流、租赁/折旧和制造费用的 actual-cost 口径与项目毛利公式。
- [x] G8 冻结 M12 release manifest -> Genesis O2D master-data 映射、新 transaction/correlation code 和 M3/M7 regression boundary。
- [x] G9 完成 IAOS order/MRP/procurement/production/inventory/shipment/invoice/AR/cash/cost gap audit，冻结 payload、权限和跨仓交付顺序。

G4-G9 完成前只允许 Schema、fixture、离线财务/数量规则、兼容性分析和只读 IAOS gap audit；不得开发 IAOS 写端点或创建正式业务数据。

## 3. 实施切片

### E0 - 订单、履约、财务与成本机器合同（第 1 周）

- [x] T1 定义 OrderDemand、OrderChange、Plan/MRP、SupplyLot、ProductionRun、InventoryLot、Shipment/Acceptance、Invoice/AR/Receipt、CostFact、Margin 和 CommercialCycleGate stable code/owner。
- [x] T2 固定订单 10,000 + 追加 2,000、交期、三批 9,000/2,700/300、价格/税率/付款条件和客户接受规则。
- [x] T3 固定 M12 财务结转、opening cash/contract liability/payable/commitment、工作资金和最低现金缓冲。
- [x] T4 固定 BOM/MRP、供应/来检、生产/报废、库存/发运/接受的单位、精度和守恒公式。
- [x] T5 固定 actual-cost 元素、费率版本、归集对象、分摊、收入/合同负债和毛利公式。
- [x] T6 固定 300 件短缺、Knowledge 延迟、恢复动作、客户接受、开票和回款时间线。
- [x] T7 扩展 JSON Schema、Go 类型、strict parser、canonical hash、payload registry 及合法/破损 fixture。
- [x] T8 定义 `eligible -> ordered -> planned -> supplied -> producing -> shipment_1 -> shipment_2 -> recovery -> delivered -> invoiced -> collected -> cost_closed -> first_commercial_cycle_closed` 状态机。
- [x] T9 完成 IAOS gap ledger、职责冲突/权限矩阵、Atlas 依赖和两仓交付矩阵。

验收：数量、金额、税、成本、现金、三态 owner 和 Genesis/M14 边界无歧义，G4-G9 关闭，M12 opening reconciliation 可机器验证。

### E1 - M12 terminal 到 Genesis O2D 兼容适配（第 1-2 周）

- [x] T10 验证 M12 product/BOM/routing/control-plan/PPAP release manifest/hash，不匹配时失败关闭。
- [x] T11 编译 Genesis-specific O2D master data/input，明确零可销售成品库存和 M11 accepted resource refs。
- [x] T12 建立新订单、采购、工单、批次、发运、发票、收款和 correlation/idempotency stable code，避免与旧 `SO-202607-0001` 冲突。
- [x] T13 复用现有 22-event 语义、validator 和 IAOS adapter，但不复制旧 opening inventory、运行记录或已发生事件。
- [x] T14 实现默认 dry-run 的 migration/impact report，外部写入要求显式 `--apply` 和目标环境。
- [x] T15 验证 M3/M7 原 pack hash/行为不变，Genesis reset 不删除 M9-M12 L1 事实或旧场景数据。

验收：M13 以 M12 机器终态启动，旧 O2D 只提供受测兼容能力，不提供历史事实或库存。

### E2 - 正式订单、ATP/MRP 与供应执行（第 2-4 周）

- [x] T16 实现客户正式订单/追加要求的 World 产生、送达延迟、Knowledge 和 IAOS 受治理确认/变更。
- [x] T17 实现 release manifest、库存、BOM、产能、供应周期和现金约束下的 ATP/MRP/交付承诺。
- [x] T18 实现采购申请/订单、供应商确认、运输、收货、批次、证书、IQC、隔离和放行链。
- [x] T19 实现供应延期与备选供应商加严检验 observation/intent/outcome/world consequence。
- [x] T20 验证 RFQ/预测不能转正式需求，未放行材料不能 ATP/发料，重复订单/PO 和跨租户引用失败关闭。
- [x] T21 对账需求、采购、收货、拒收、可用量、承诺和供应现金流。

验收：正式需求和物料来源完整；计划只读取已知事实并通过 IAOS 治理提交动作。

### E3 - 正式生产、质量、库存与实际成本（第 3-6 周）

- [x] T22 实现生产订单/工序任务、材料发放、设备/人员/班次占用、报工和过程检验 world reducer。
- [x] T23 实现设备停机、维修、加班、替代资源和能力恢复，禁止 IAOS 状态直接修复 World 设备。
- [x] T24 实现良品、报废、返工、在制、完工入库、批次/序列和不可销售状态守恒。
- [x] T25 实现材料、人工、设备/能源、报废、质量、维修/加班和制造费用 actual facts 与 rate-version 归集。
- [x] T26 实现 standard vs actual cost、quantity/rate/efficiency/scrap/expedite variance。
- [x] T27 验证缺 World evidence、未资格人员、未校准设备、未发布 revision、重复报工/入库和负库存失败关闭。
- [x] T28 实现 snapshot/restore/reset、100 次确定性、中途崩溃恢复和事件/成本 hash 一致。

验收：生产、质量、库存和成本由同一批次/工单/资源证据驱动，数量与金额均可重放。

### E4 - 分批发运、客户接受与 300 件恢复（第 5-7 周）

- [x] T29 实现 allocation、pick、pack、dispatch、transport、customer receipt、acceptance 和 POD 状态机。
- [x] T30 执行 9,000 和 2,700 两批实际发运，形成累计 11,700 与 300 discrepancy；未发货数量不可接受/开票。
- [x] T31 observation 送达后运行计划/质量/经营 Agent，提交备选供应、维修/加班和第三批恢复 intent。
- [x] T32 committed outcomes 后由 World 计算第三批 300 的材料、产能、成本、发运和到达后果。
- [x] T33 实现客户数量/质量接受、争议、拒收失败路径和累计 12,000 delivery close。
- [x] T34 验证超库存发运、重复 dispatch/POD、发运即接受、未解决争议关闭和旧 shipment code 冲突失败关闭。

验收：三批客户实际接受合计 12,000，短缺差异关闭且每批来源/成本/时间可追溯。

### E5 - 开票、应收、银行回款与项目盈亏（第 6-8 周）

- [x] T35 实现基于 accepted quantity、合同价格和配置税率的 invoice plan、审批、开具和拆票/合计守恒。
- [x] T36 实现 invoice -> AR，不把发票当现金；未接受数量和重复 invoice key 失败关闭。
- [x] T37 实现客户付款/延期/重复通知的 World bank policy 和 actor-scoped observation。
- [x] T38 实现 bank receipt -> IAOS cash receipt -> AR settlement committed outcome，金额、币种和 reference 必须匹配。
- [x] T39 冻结并实现开发预付款合同负债的履约转收入/抵扣，防止重复确认。
- [x] T40 汇总 product/order/project actual cost、recognized revenue、gross margin、cash movement 和 variance explanation。
- [x] T41 验证部分收款、超额核销、错客户/币种、伪造银行到账、负毛利和缺失 actuals 的失败/解释路径。

验收：收入、应收、收款、现金和毛利彼此独立又能对账；经营分析不再以 `partial` 掩盖成本缺口。

### E6 - IAOS O2D、财务与成本治理（第 4-8 周）

- [x] T42 按 IAOS 规则建立新的独立 branch/worktree，先读其 AGENTS、Agent Context 和 Code Map。
- [x] T43 复用现有 scenario/O2D/metadata/Capability/Process；只为订单变更、POD/接受、invoice/AR/cash settlement、cost actuals 和 profitability 缺口增加最小实现。
- [x] T44 实现订单/变更、承诺、采购、生产、入库、发运、接受、开票、核销、成本关闭和周期关闭的 allowlist Capability/Process/Decision。
- [x] T45 固定 `genesis.delivery.order/plan/procure/produce/ship/accept/invoice/collect/cost/close` 权限与销售/计划/采购/生产/质量/物流/财务 mandate。
- [x] T46 保证 business record、inventory/cost posting、audit、journal committed outcome 与 Outbox 的必要事务原子性；回滚不产生 outcome。
- [x] T47 验证 tenant RLS、越权自批、重复/并发订单、负库存、重复发运/开票/收款、未接受开票、伪造成本和失败无部分写入。
- [x] T48 两仓分别提交、记录 revision、部署并完成 contract/integration tests。

验收：AESE 不直写 IAOS；IAOS 不伪造客户/供应商/运输/银行事实；数量、金额和成本动作受权限、幂等、事务和审计约束。

### E7 - 统一 Agent 与 First Delivery Play（第 7-9 周）

- [x] T49 将计划、质量和经营分析 Agent 接入同一订单/correlation，继续使用 IAOS Tool/Capability 和 human takeover 治理链。
- [x] T50 在 World Play 增加 First Delivery campaign、需求/MRP、采购/生产、批次库存、三批交付、发票/应收/现金和利润视图。
- [x] T51 展示 World 实际状态、IAOS 业务状态、角色 Knowledge、300 件 discrepancy、成本差异和因果证据，不提供直接改库存/接受/现金入口。
- [x] T52 展示从 M9 成立到 M13 商业闭环的 Genesis milestone/terminal-hash 链，并保留各 campaign 独立重放。
- [x] T53 完成键盘、ARIA、焦点、移动端、金额/数量/成本单位、数据 owner 和虚拟/真实时间表达。

验收：非研发用户可从正式订单推进到收款/利润关闭，并能解释短缺恢复、成本差异和现金/利润不同步。

### E8 - 全链验收与 Project Genesis 收口（第 9-10 周）

- [x] T54 从 M12 clean terminal state 分别执行 Agent 路径和人类接管路径，终态业务结果/state hash 一致。
- [x] T55 对账 World events、Knowledge、IAOS Decision/Process/Capability、journal、Outbox、订单、采购、工单、库存、发运、invoice/AR、cash 和 cost ledger。
- [x] T56 验证断线、SSE 丢失、重启、重复点击、陈旧 cursor、乱序、并发、快照恢复和分层 reset。
- [x] T57 验证 tenant-other、只读用户、越权自批、负库存、未接受开票、重复收款、伪造银行/成本和跨场景污染。
- [x] T58 运行 Go、Schema、PostgreSQL、IAOS modules、前端、Playwright 及 M3-M12 回归和 390/1280/1440 三视口验收。
- [x] T59 编写 M13 runbook/evidence，更新设计、计划、Roadmap、Code Map、Progress Log、Atlas、能力缺口和两仓 revision/部署信息。
- [x] T60 发布 Project Genesis M9-M13 端到端验收报告，证明所有 terminal contract、资金/数量/成本和旧 O2D 兼容链可重放。

验收：只有累计 12,000 客户接受、发票/应收正确、银行实际到账并核销、actual cost 完整和项目毛利可解释时输出 `first_commercial_cycle_closed=true`。

## 4. 首版编码方向

E0 最终冻结前使用以下虚构编码方向：

```text
SO-GENESIS-0001
SO-GENESIS-0001-ADD1
MRP-GENESIS-0001
PO-GENESIS-AL-0001
WO-GENESIS-BCP-A01-0001
LOT-GENESIS-BCP-A01-001
SHIP-GENESIS-0001-A
SHIP-GENESIS-0001-B
SHIP-GENESIS-0001-C
POD-GENESIS-0001
INV-GENESIS-0001
AR-GENESIS-0001
RCPT-GENESIS-0001
COST-GENESIS-0001
MARGIN-PROJECT-SGNEV-BCP-A01
corr-genesis-first-delivery-001
```

客户、供应商、订单、价格、税率、发票和银行数据全部虚构；不得与真实主体或旧场景交易 code 混用。

## 5. 完成定义

- First Delivery campaign 可验证、初始化、推进、重放、快照恢复和分层 reset。
- M12 financial opening、合同负债、遗留承诺和试制成本结转无未解释差额。
- 正式需求、采购、生产、报废、库存、三批发运和客户接受数量严格守恒。
- 300 件恢复 tracer 有 observation/intent/outcome/world consequence/final delivery/close 完整链。
- invoice/AR、客户接受、bank receipt、cash settlement 和收入确认严格分离。
- actual cost 覆盖材料、人工、设备/能源、报废、质量、加急、物流和分摊，并解释 standard/actual variance。
- 人类/Agent 共用治理链；越权、重复交易、负库存、未接受开票、伪造银行/成本失败关闭。
- IAOS 事务/Outbox、tenant/RLS、幂等、并发和失败无部分写入通过。
- M3-M12 零回归；Project Genesis M9-M13 terminal hash 链和端到端证据完整。
- `first_commercial_cycle_closed=true` 只代表首单商业闭环，不代表长期持续经营或多周期稳定盈利。

## 6. 不纳入 M13

- 第二订单/客户/产品、多工厂或完整滚动 S&OP/MPS。
- 完整总账、税务申报、资金管理、银行/电子发票真实接口和法定财务报表。
- 售后、退货、质保、贷项、坏账、跨期收入和复杂成本会计。
- 长期经营、参数化分支、Monte Carlo 和 A/B；属于 M14。

## 7. 并行与所有者规则

- E0/E1 由 AESE owner 先冻结财务、数量、兼容和迁移合同。
- E2/E3 可在 E0/E1 稳定后并行，但库存、现金和 cost ledger 由单一 owner 合并。
- E4 依赖可验证成品批次；E5 依赖客户 acceptance，不能以 UI 或 shipment status 跳过。
- E6 只能在 payload、权限和 IAOS gap 冻结后由新的 IAOS worktree owner 开始。
- E7 在 API/view model 稳定后开始，不得为了 UI 直接改 World/IAOS 状态。
- E8 串行收口两仓和 M9-M13 证据；并行 agent 必须有子计划、owner 和不重叠 worktree。
- 不得覆盖当前共享工作区中的测试修改、截图变化或验收产物。
