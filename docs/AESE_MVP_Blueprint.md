# AESE MVP Blueprint

## 1. 产品定位

AESE 是 IAOS 的智能企业运行仿真环境，用于构建一个可运行、可观察、可干预、可被 AI Agent 操作的虚拟工业企业。

MVP 的核心不是视觉效果，而是验证一件事：

> IAOS 能否让一家虚拟汽车零部件企业从订单到交付真实运转起来，并在异常发生时由系统和 Agent 共同解释、处理和复盘。

## 2. 第一阶段虚拟企业

企业名称：

- 华辰热管理系统集团有限公司。

行业：

- 新能源汽车热管理零部件。

业务模式：

- 主机厂和一级供应商客户。
- 项目制开发加批量供货。
- 订单、预测和 JIT 交付并存。
- 集团统一管控，工厂独立运营。

MVP 工厂：

- 苏州制造基地。

MVP 产品：

- 电池冷却板组件。

## 3. MVP 业务范围

MVP 主线：

```text
客户订单
-> 订单确认
-> MRP 运算
-> 采购需求
-> 采购订单
-> 供应商发货
-> 来料收货
-> 来料检验
-> 原材料入库
-> 生产订单
-> 工序任务
-> 过程检验
-> 完工入库
-> 成品发货
-> 开票
-> 经营分析
```

MVP 不包含：

- 完整 APQP / 新品开发。
- 完整财务总账。
- 完整 HR、EHS、售后服务。
- 多工厂协同优化。
- 真实 3D 工厂。

## 4. 核心对象模型

组织对象：

- 集团。
- 事业部。
- 法人公司。
- 工厂。
- 部门。
- 班组。

资源对象：

- 厂区。
- 车间。
- 产线。
- 工作中心。
- 设备。
- 模具。
- 工装。
- 仓库。
- 库区。
- 库位。

业务对象：

- 客户。
- 供应商。
- 物料。
- BOM。
- 工艺路线。
- 销售订单。
- 采购订单。
- 收货单。
- 检验单。
- 库存事务。
- 生产订单。
- 工序任务。
- 完工入库单。
- 发货单。
- 发票。
- 质量问题。
- 设备故障。

人员和角色：

- 销售经理。
- 计划员。
- 采购员。
- 仓库主管。
- 检验员。
- 班组长。
- 操作工。
- 质量工程师。
- 设备工程师。
- 工厂厂长。
- 财务会计。
- AI Agent 管理员。

## 5. 核心事件流

第一阶段事件建议：

```text
CustomerOrderReceived
SalesOrderConfirmed
MRPGenerated
PurchaseRequirementCreated
PurchaseOrderReleased
SupplierShipmentDispatched
MaterialReceived
IncomingInspectionPassed
IncomingInspectionFailed
InventoryPutawayCompleted
ProductionOrderReleased
OperationStarted
OperationCompleted
ProcessDefectDetected
ProductionOrderCompleted
FinishedGoodsReceived
ShipmentDispatched
InvoiceIssued
BusinessKPIUpdated
```

异常事件：

```text
SupplierDeliveryDelayed
MachineDown
IncomingMaterialRejected
ProductionScrapIncreased
CustomerExpediteRequested
ShipmentShortageDetected
```

事件设计要求：

- 每个事件必须能追溯租户、组织、对象、状态变化和触发原因。
- 事件要能驱动流程、任务、通知、Capability 或 Agent 分析。
- 事件命名应与 IAOS 现有 `shared/eventdef` 风格保持一致。

## 6. IAOS 能力映射

| AESE 需求 | IAOS 能力 |
| --- | --- |
| 虚拟企业主数据 | Metadata / Dynamic Entity |
| 订单到交付链路 | `scenarios/o2d` |
| 状态变化通知 | Outbox + NATS |
| 业务动作封装 | Capability Runtime |
| 规则和约束 | Constraint / Policy |
| 决策解释 | Decision Runtime |
| 审批和流程 | Process Runtime |
| Agent 工具调用 | AI Tool Registry |
| 操作留痕 | Audit |
| 前端操作台 | IAOS Frontend |

## 7. 第一阶段 Agent

### 计划 Agent

输入：

- 销售订单。
- 库存。
- BOM。
- 工艺路线。
- 产能。
- 采购交期。
- 设备状态。

输出：

- 交付风险判断。
- MRP 或排产建议。
- 采购加急建议。
- 产能调整建议。

### 质量 Agent

输入：

- 来料检验记录。
- 过程缺陷记录。
- 批次追溯。
- 供应商批次。
- 工序和设备数据。

输出：

- 缺陷聚类。
- 根因假设。
- 隔离范围建议。
- 供应商风险提示。

### 经营分析 Agent

输入：

- 订单交付。
- 库存。
- 采购。
- 生产。
- 质量。
- 设备。
- 成本。

输出：

- KPI 解释。
- 利润波动原因。
- 异常对交付和成本的影响。
- 管理层行动建议。

## 8. MVP 演示故事线

故事名称：

- 客户追加订单下的交付承诺重算。

背景：

- 主机厂客户临时追加 20% 电池冷却板订单。
- 关键铝材供应商延期 3 天。
- 苏州工厂一台关键焊接设备发生故障。

系统运行：

1. 销售订单确认并触发 MRP。
2. MRP 发现原材料缺口和产能风险。
3. 供应商延期事件进入事件流。
4. 设备故障事件进入事件流。
5. 计划 Agent 生成交期风险和调整建议。
6. 系统生成采购加急、替代产线或加班方案。
7. 经营分析 Agent 解释不同方案对交付率、成本、库存和利润的影响。

演示价值：

- 证明 IAOS 能从静态业务系统变成可运行的经营系统。
- 证明 Agent 可以基于真实上下文参与企业决策。
- 证明 AESE 能作为客户演示、研发测试和产品设计环境。

## 9. 开发里程碑

### M0 - 项目初始化

- 建立 AESE 仓库。
- 建立项目背景、MVP 蓝图和进展跟踪规则。
- 固化虚拟客户和 MVP 范围。

### M1 - 虚拟企业蓝图

- 完成华辰热管理系统集团设定。
- 完成苏州制造基地设定。
- 完成电池冷却板产品、BOM、工艺、设备、仓库、角色设定。

### M2 - 主数据与事件清单

- 定义 MVP 主数据。
- 定义 MVP 事件类型。
- 映射到 IAOS metadata、eventdef 和 scenario package。

### M3 - O2D 仿真主线

- 复用或扩展 IAOS `scenarios/o2d`。
- 跑通订单确认、MRP、采购、生产、入库、发货事件链。

### M4 - 异常场景

- 加入供应商延期、设备故障、来料不良三类异常。
- 让异常驱动流程、任务或 Agent 分析。

### M5 - Agent MVP

- 接入计划 Agent、质量 Agent、经营分析 Agent 的最小能力。
- Agent 先生成建议和解释，再逐步进入受治理动作。

### M6 - 2D 企业沙盘

- 展示订单流、库存流、产线状态、异常事件和 Agent 建议。
- 不做 3D，除非业务链路已经稳定。

## 10. 关键风险

- 范围膨胀：想一次覆盖集团、全流程、全角色、全可视化。
- 视觉先行：过早投入 3D，导致业务运行模型薄弱。
- 与 IAOS 脱节：另建系统，无法验证 IAOS 真实能力。
- Agent 空心化：只做聊天解释，不接业务上下文、Capability 和审计。
- 数据不真实：对象和事件过于玩具化，无法打动制造业客户。

## 11. 下一步建议

优先完成 M1：

- 写出华辰热管理系统集团详细蓝图。
- 选定苏州制造基地的一条电池冷却板生产线。
- 明确 20 到 30 个关键主数据对象。
- 明确 15 到 20 个关键事件。
- 明确第一条演示故事线的数据输入和预期输出。

