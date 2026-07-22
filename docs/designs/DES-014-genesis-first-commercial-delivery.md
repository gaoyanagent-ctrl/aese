---
id: DES-014
title: M13 Project Genesis 第一次完整商业交付
date: 2026-07-22
status: completed
author: Codex + User
tags: [m13, genesis, o2d, delivery, invoice, cash, actual-cost]
---

# M13 Project Genesis 第一次完整商业交付

## 1. 目标

M13 从 M12 的 `serial_production_eligible=true` 开始，让华辰苏州制造公司接收第一张正式客户订单，完成受治理的订单确认、ATP/MRP、采购、来料检验、生产、过程质量、入库、分批发运、客户收货、开票、应收和实际回款，并形成第一张订单和客户项目的实际成本及毛利，输出 `first_commercial_cycle_closed=true`。

这是 Project Genesis M9-M13 主纵向场景的终点。终态必须证明企业从成立、建设、能力形成、工业化走到真实商业闭环，而不是只复放一份预制 O2D 数据。

## 2. M12 财务结转门

M12 终态显示 10,500,000 CNY 现金、2,000,000 CNY 客户开发预付款合同负债和 1,200,000 CNY 试制实际成本，但当前实现没有明确试制成本是已付款、应付还是由预付款覆盖。M13 在任何采购或生产前必须冻结并机器化结转：

- 试制成本的支付状态、对应现金/应付和项目成本归属。
- 2,000,000 CNY 合同负债的履约义务、PPAP/工装接受证据，以及在 M13 是转收入、抵扣订单还是继续递延。
- M11 工资准备金、最低现金缓冲、设备租赁/设施尾款和其他未清承诺。
- 第一张订单的工作资金需求、客户付款条件和最坏现金缺口。

结转不平或缺少依据时，M13 必须失败关闭，不能靠修改 opening cash 或利润配平。

## 3. 首版商业故事

首版复用现有 `order-expedite-01` 的业务张力，但不复用其 opening inventory 或既成事件：

```text
消费 M12 serial_production_eligible
-> 客户发出 10,000 件正式订单
-> 受治理确认、ATP/MRP 和原材料采购
-> 客户追加 2,000 件并保持原交期
-> 供应延期 + 关键设备停机造成缺口风险
-> 计划/质量/经营 Agent 提交受治理建议
-> 备选供应商、加严检验、维修/加班和分批交付
-> 先发 9,000，再发 2,700，最后恢复交付 300
-> 客户实际收货并接受累计 12,000 件
-> 按接受数量和合同条件开票、形成应收
-> 客户实际付款到账并受治理核销
-> 计算实际订单成本、项目毛利和现金变化
-> first_commercial_cycle_closed
```

Genesis 首单使用独立订单/run/correlation/idempotency code，不与 M3/M7 的 `SO-202607-0001` 活跃记录冲突。

## 4. 与现有 O2D 的兼容与事实边界

M13 复用 IAOS 已实现的 scenario apply/reset、订单确认、MRP、采购、检验、工单、受治理完工入库、发运、Outbox/NATS、snapshot/cursor/SSE 和 Agent tracer；不在 AESE 复制 ERP/MES 逻辑。

`scenario-packs/hctm` 仍是 M3-M7 回归 fixture。M13 必须从 M12 release manifest/hash 和 M13 opening contract 编译或投影 Genesis-specific O2D input：

- 零可销售成品期初库存；M12 的试制/PPAP 样件不得迁入。
- 产品/BOM/routing/control-plan revision 必须与 M12 terminal manifest 匹配。
- 正式订单、采购、生产、库存和发运使用新的 stable transaction code。
- 旧 22-event story 可复用事件语义和验证器，但不能把旧事件当作已经发生。

## 5. 三态所有权

| 事实 | World State | IAOS Business State | Actor Knowledge |
| --- | --- | --- | --- |
| 客户正式需求与变更 | Own：订单/追加要求、收货、接受、争议和付款行为 | 销售订单、变更、承诺日期和状态 | 销售/计划按 observation 获知 |
| 供应与生产现实 | Own：供应商交付、材料属性、设备状态、实际产出/报废/消耗 | PO、收货/IQC、工单、报工、库存和质量记录 | 角色按事件、权限和延迟获知 |
| 物流与客户收货 | Own：实际运输、到达、数量/损坏和客户接受 | shipment、delivery note、proof of delivery | 销售/物流在 observation 后获知 |
| 发票、应收和成本账 | 必要外部事实引用；不复制业务账 | Own：invoice、AR、cost transaction、settlement、项目损益 | CFO/经营分析按权限获知 |
| 实际现金 | Own：银行入账/出账与时间 | cash receipt/payment posting 和核销记录 | CFO 在 bank observation 后获知 |

IAOS 订单、发运单、发票或收款记录不能伪造客户需求、物理发运、客户接受或银行到账；AESE 也不能直写 IAOS 账表。

## 6. 最小领域模型

### AESE World

- `CustomerOrderDemand`、`CustomerChange`、`CustomerReceipt`、`CustomerAcceptance`。
- `SupplierDelivery`、`MaterialLotReality`、`MachineAvailability`。
- `ProductionRun`、`OperationActual`、`MaterialConsumption`、`ProductionYield`。
- `TransportLeg`、`ShipmentReality`、`DeliveryDiscrepancy`。
- `BankReceipt`、`CashMovement`。
- `ActualCostFact`：材料、人工、机器/能源、报废、质量、加急和物流实际量。
- `CommercialCycleGate`。

### IAOS 最小管理对象

- `sales_order`、`sales_order_change`、`delivery_commitment`
- `mrp_run`、`purchase_requisition`、`purchase_order`
- `goods_receipt`、`incoming_inspection`、`inventory_transaction`
- `production_order`、`operation_task`、`production_report`、`process_inspection`
- `finished_goods_receipt`、`shipment`、`proof_of_delivery`
- `invoice`、`accounts_receivable`、`cash_receipt`、`settlement`
- `cost_transaction`、`actual_cost`、`order_margin`、`project_profitability`

优先复用现有 IAOS O2D 对象和 Capability；只为发票/应收/收款/实际成本等真实缺口增加最小 metadata/config 与 allowlist，不建设完整总账或税务系统。

## 7. 数量、库存与交付不变量

- 没有 M12 `serial_production_eligible=true` 和匹配 release manifest 不能确认正式订单或释放生产。
- M13 opening saleable finished goods 必须为零；试制件和 PPAP 样件不可销售。
- 正式需求为原单加受治理变更之和；预测、RFQ、定点年量不计正式需求。
- MRP/BOM 需求、采购、收货、IQC、发料、在制、良品、报废、成品和发运数量必须逐单位/批次守恒。
- 未放行材料不能投产，未合格完工不能入库，未入库/未分配库存不能发运。
- 发运不等于客户收货，收货不等于客户接受；只有被接受数量可进入开票门。
- 累计接受数量达到 12,000、争议/短缺关闭后才能关闭交付 gate。
- 重放、重复点击、乱序或网络重试不能重复采购、报工、入库、发运、开票或收款。

## 8. 发票、应收、现金与收入边界

- 发票必须引用客户接受的 shipment/quantity、合同价格、配置化税率和唯一 invoice key。
- 未交付/未接受的 300 件不得提前开票；拆票必须保持累计净额/税额/含税额守恒。
- 发票成立形成应收，不等于现金到账；客户付款由 World/银行策略产生。
- `BankReceipt` observation 与 IAOS cash receipt/settlement committed outcome 对账后才增加已核销回款。
- 客户开发预付款的合同负债只能依据冻结履约规则转收入或抵扣，禁止重复确认。
- 退款、贷项和坏账首版只做失败路径，不进入成功 tracer。
- 配置税率只服务虚构场景计算，不宣称真实税务合规。

## 9. 实际成本与项目盈亏

M13 关闭此前 `cost_actuals=partial` 缺口，首版归集：

- 实际材料消耗与采购/到岸单价。
- 实际人工工时与版本化费率。
- 设备/能源实际用量与费率。
- 报废、返工、加严检验、维修/加班和加急物流。
- 可解释的设备租赁/折旧及制造费用分摊规则。

必须同时展示：报价成本、量产标准成本、订单实际成本和差异。项目毛利只在收入确认与实际成本结转后计算；毛利不等于现金，现金余额不等于利润。首版不生成完整会计分录、总账、资产负债表或法定报表。

## 10. 最小异常 tracer

M13 固定“订单追加 + 供应延期 + 设备停机导致末批 300 件短缺”的闭环：

```text
IAOS 计划承诺 12,000 件按期完成
-> World 供应和设备约束只支持首两批 11,700 件
-> World/IAOS 产生 300 件 discrepancy，客户与部分角色认知有延迟
-> observation 后计划/质量/经营 Agent 提交受治理恢复方案
-> IAOS committed outcomes 批准备选供应、加严检验、维修/加班与第三批计划
-> AESE 计算新的材料、产能、成本和交期后果
-> 第三批 300 件实际生产、发运、客户收货并接受
-> 只对累计已接受 12,000 件完成开票、回款和差异关闭
```

该 tracer 必须保留 M4-M7 异常语义，同时证明 UI、计划状态或发票不能直接创造物理交付和现金。

## 11. 角色与治理

| 角色 | 主要动作 | 关键限制 |
| --- | --- | --- |
| 销售/客户服务 | 订单确认、变更、交付承诺、收货/接受跟踪 | 不能伪造客户要求、POD 或接受 |
| 计划/采购 | ATP/MRP、采购、排产和异常恢复 | 不能绕过材料、预算、产能和审批门 |
| 生产/质量/物流 | 生产、检验、入库、发运和交付证据 | 不能修改 World 数量/质量或重复报工 |
| CFO/财务会计 | 开票、应收、收款核销、成本和项目损益 | 不能把发票当现金或把预付款重复计收入 |
| 经营分析 Agent | 解释交付、成本、毛利和现金差异 | 只能引用完整 evidence，不得填补缺失 actuals |

人类和 Agent 使用同一 Capability、Policy、Decision、Process、AI Tool 和审计链；订单/价格例外、采购、库存、发运、开票、核销和成本结转保持职责隔离。

## 12. Bridge payload family

M13 增加严格 allowlist 类型：

```text
genesis.customer.order.received.v1
genesis.customer.order.changed.v1
genesis.delivery.commitment.approved.v1
genesis.supplier.delivery.delayed.v1
genesis.production.disrupted.v1
genesis.production.completed.v1
genesis.finished.goods.received.v1
genesis.shipment.dispatched.v1
genesis.customer.delivery.received.v1
genesis.customer.delivery.accepted.v1
genesis.delivery.recovery.approved.v1
genesis.invoice.issued.v1
genesis.bank.receipt.observed.v1
genesis.cash.receipt.settled.v1
genesis.actual.cost.closed.v1
genesis.commercial.cycle.closed.v1
```

业务写入继续走 IAOS 受治理 API/Capability/Process 和事务 Outbox；外生客户/供应商/运输/银行事件走 simulation ingress/World Bridge，禁止正式路径 direct NATS。

## 13. Pack、迁移与重置策略

- `hctm-genesis` 升级到下一 minor version，新增 `campaigns/first-delivery/`。
- 初态验证 M12 terminal hash、release manifest、`serial_production_eligible`、现金/合同负债/承诺及零可销售库存。
- Genesis O2D 使用新 transaction code 和 correlation；旧 M3/M7 pack 保持不变并继续回归。
- 兼容 adapter 必须显式声明从 M12 release manifest 到 IAOS O2D master data 的版本映射。
- reset 只删除本次 M13 L2/L3 运行对象，保留 M9-M12 法人、设施、能力和产品发布事实；危险 reset 仍需 dry-run、影响预览和显式确认。

## 14. 非目标

- 第二客户/产品/工厂、多订单长期滚动经营或完整 S&OP。
- 完整财务总账、税务申报、成本会计、资金管理、银行或电子发票真实接口。
- 售后索赔、质保、退货、贷项、坏账和跨期收入确认完整闭环。
- 参数化分支、Monte Carlo、A/B 和长期经营实验；属于 M14。
- 真实客户、价格、发票、银行或生产数据。

## 15. 完成标准

- 单一 run 从 M12 eligibility 确定性推进到 `first_commercial_cycle_closed=true`。
- 正式需求 12,000 件、材料/在制/良品/报废/库存/三批发运/客户接受数量完全守恒。
- 300 件短缺形成 observation、受治理恢复、第三批交付和 discrepancy close 完整链。
- 发票/应收、客户接受、银行到账和收款核销严格分离且金额守恒。
- M12 试制成本/开发预付款结转、订单实际成本、收入、项目毛利和 closing cash 可解释对账。
- 人类/Agent 共用治理能力，越权、自批、重复交易、未交付开票、伪造收款和缺 actuals 盈利结论失败关闭。
- M3-M12 回归、两仓测试、API/UI、runbook/evidence、revision、部署和 Atlas 完整。
