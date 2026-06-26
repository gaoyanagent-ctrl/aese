# HCTM Event Model

本文把华辰热管理系统集团 MVP 蓝图中的 18 个关键事件转成 IAOS 可落地的事件规格。目标是统一 AESE 事件命名、NATS subject、payload、上下游对象、幂等键、Agent 触发和 Capability / Process 接线方式。

## 1. IAOS 对齐约定

IAOS 当前事件模型使用：

```text
subject = iaos.{tenant_id}.{event_type}
event_type = {scenario}.{entity}.{action}
```

示例：

```text
iaos.tenant-hctm.o2d.order.confirmed
```

AESE MVP 事件沿用 IAOS 现有规则：

- 租户：`tenant-hctm`
- 主场景包：`o2d`
- 辅助领域前缀：`proc`、`qms`、`eam`、`whs`
- subject 不直接使用 CamelCase 名称。
- CamelCase 名称保留为业务蓝图事件名。

## 2. 通用事件 Envelope

所有 AESE 事件建议使用统一 envelope：

```json
{
  "id": "evt-20260701-000001",
  "tenant_id": "tenant-hctm",
  "type": "o2d.order.confirmed",
  "source": "aese.hctm.simulator",
  "timestamp": "2026-07-01T10:00:00+08:00",
  "data": {},
  "metadata": {
    "correlation_id": "corr-so-202607-0001",
    "causation_id": "evt-20260701-000000",
    "scenario": "hctm_order_expedite_01",
    "plant_code": "PLT-SZ",
    "business_object_type": "sales_order",
    "business_object_id": "SO-202607-0001",
    "idempotency_key": "tenant-hctm:o2d.order.confirmed:SO-202607-0001:v1"
  }
}
```

Envelope 字段：

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `id` | 是 | 事件实例唯一 ID |
| `tenant_id` | 是 | IAOS 租户 |
| `type` | 是 | dotted event type |
| `source` | 是 | 事件来源服务 |
| `timestamp` | 是 | 事件产生时间 |
| `data` | 是 | 事件业务 payload |
| `metadata.correlation_id` | 是 | 串联同一业务链路 |
| `metadata.causation_id` | 否 | 上游事件 ID |
| `metadata.scenario` | 是 | AESE 演示故事或仿真场景 |
| `metadata.plant_code` | 否 | 工厂范围 |
| `metadata.business_object_type` | 是 | 主业务对象类型 |
| `metadata.business_object_id` | 是 | 主业务对象编号或 ID |
| `metadata.idempotency_key` | 是 | 消费端幂等键 |

## 3. 事件命名映射

| 蓝图事件 | IAOS event type | NATS subject 示例 | 主对象 |
| --- | --- | --- | --- |
| `CustomerOrderReceived` | `o2d.order.received` | `iaos.tenant-hctm.o2d.order.received` | `sales_order` |
| `SalesOrderConfirmed` | `o2d.order.confirmed` | `iaos.tenant-hctm.o2d.order.confirmed` | `sales_order` |
| `MRPGenerated` | `o2d.mrp.generated` | `iaos.tenant-hctm.o2d.mrp.generated` | `sales_order` |
| `PurchaseRequirementCreated` | `o2d.purchase_requirement.created` | `iaos.tenant-hctm.o2d.purchase_requirement.created` | `sales_order` |
| `PurchaseOrderReleased` | `o2d.purchase_order.released` | `iaos.tenant-hctm.o2d.purchase_order.released` | `purchase_order` |
| `SupplierShipmentDispatched` | `o2d.supplier_shipment.dispatched` | `iaos.tenant-hctm.o2d.supplier_shipment.dispatched` | `purchase_order` |
| `SupplierDeliveryDelayed` | `o2d.supplier_delivery.delayed` | `iaos.tenant-hctm.o2d.supplier_delivery.delayed` | `purchase_order` |
| `MaterialReceived` | `whs.material.received` | `iaos.tenant-hctm.whs.material.received` | `goods_receipt` |
| `IncomingInspectionPassed` | `qms.incoming_inspection.passed` | `iaos.tenant-hctm.qms.incoming_inspection.passed` | `inspection_order` |
| `IncomingInspectionFailed` | `qms.incoming_inspection.failed` | `iaos.tenant-hctm.qms.incoming_inspection.failed` | `inspection_order` |
| `InventoryPutawayCompleted` | `whs.inventory.putaway_completed` | `iaos.tenant-hctm.whs.inventory.putaway_completed` | `inventory_transaction` |
| `ProductionOrderReleased` | `proc.production_order.released` | `iaos.tenant-hctm.proc.production_order.released` | `production_order` |
| `OperationStarted` | `proc.operation.started` | `iaos.tenant-hctm.proc.operation.started` | `operation_task` |
| `MachineDown` | `eam.machine.down` | `iaos.tenant-hctm.eam.machine.down` | `equipment` |
| `OperationCompleted` | `proc.operation.completed` | `iaos.tenant-hctm.proc.operation.completed` | `operation_task` |
| `ProcessInspectionFailed` | `qms.process_inspection.failed` | `iaos.tenant-hctm.qms.process_inspection.failed` | `inspection_order` |
| `FinishedGoodsReceived` | `whs.finished_goods.received` | `iaos.tenant-hctm.whs.finished_goods.received` | `inventory_transaction` |
| `ShipmentDispatched` | `o2d.shipment.dispatched` | `iaos.tenant-hctm.o2d.shipment.dispatched` | `shipment` |

说明：

- `SalesOrderConfirmed` 保持现有 IAOS `o2d.order.confirmed`，可直接触发当前 O2D workflow。
- `proc`、`qms`、`eam`、`whs` 是 AESE MVP 建议领域前缀，后续若 IAOS 已有正式 scenario package，可迁移到正式前缀。
- `PurchaseRequirementCreated` 在主数据模型中暂未作为独立 entity，可第一阶段作为 MRP 结果 payload 或后续新增 `purchase_requirement`。

## 4. 事件链路

第一条演示故事线的标准链路：

```text
o2d.order.received
-> o2d.order.confirmed
-> o2d.mrp.generated
-> o2d.purchase_requirement.created
-> o2d.purchase_order.released
-> o2d.supplier_shipment.dispatched
-> o2d.supplier_delivery.delayed
-> eam.machine.down
-> o2d.mrp.generated
-> proc.production_order.released
-> whs.material.received
-> qms.incoming_inspection.passed / qms.incoming_inspection.failed
-> whs.inventory.putaway_completed
-> proc.operation.started
-> proc.operation.completed
-> whs.finished_goods.received
-> o2d.shipment.dispatched
```

关键规则：

- 第二次 `o2d.mrp.generated` 的 `data.run_reason` 应为 `exception_replan`。
- 供应商延期和设备故障都必须复用同一个 `correlation_id`，便于计划 Agent 汇总影响。
- 所有库存和生产事件必须携带 `plant_code`、`material_code` 和相关对象编号。

## 5. 事件详细规格

### 5.1 `CustomerOrderReceived`

Event type：

- `o2d.order.received`

触发条件：

- 销售经理收到客户正式订单或追加订单。

主对象：

- `sales_order`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `order_no` | string | 是 | `SO-202607-0001` |
| `customer_code` | string | 是 | `CUST-SGNEV` |
| `material_code` | string | 是 | `HCTM-BCP-A01` |
| `order_qty` | number | 是 | `12000` |
| `additional_qty` | number | 否 | `2000` |
| `due_date` | date | 是 | `2026-07-20` |
| `priority` | string | 是 | `high` |
| `delivery_policy` | string | 否 | `split_allowed_final_due_fixed` |

下游：

- 创建或更新 `sales_order` 草稿。
- 触发销售确认任务。

Agent：

- 经营分析 Agent 可记录需求变化。

幂等键：

- `{tenant_id}:o2d.order.received:{order_no}:{order_qty}:{due_date}`

### 5.2 `SalesOrderConfirmed`

Event type：

- `o2d.order.confirmed`

触发条件：

- 销售订单确认，状态从 `draft` 变为 `confirmed`。

主对象：

- `sales_order`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `order_no` | string | 是 | `SO-202607-0001` |
| `customer_code` | string | 是 | `CUST-SGNEV` |
| `material_code` | string | 是 | `HCTM-BCP-A01` |
| `confirmed_qty` | number | 是 | `12000` |
| `due_date` | date | 是 | `2026-07-20` |
| `plant_code` | string | 是 | `PLT-SZ` |
| `confirmed_by` | string | 是 | `EMP-SH-SALES-001` |
| `confirmed_at` | datetime | 是 | `2026-07-01T10:10:00+08:00` |

下游：

- 触发 O2D MRP workflow。
- 调用或编排 `o2d.mrp.bom_expand`、`o2d.mrp.inventory_check`、`o2d.mrp.create_workorder`。

Agent：

- 计划 Agent 开始评估交期风险。

幂等键：

- `{tenant_id}:o2d.order.confirmed:{order_no}:{confirmed_qty}:{due_date}`

### 5.3 `MRPGenerated`

Event type：

- `o2d.mrp.generated`

触发条件：

- MRP 运算完成。
- 初始订单确认或异常触发重算。

主对象：

- `sales_order`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `mrp_run_id` | string | 是 | `MRP-202607-0001` |
| `run_reason` | enum | 是 | `initial_plan` |
| `order_no` | string | 是 | `SO-202607-0001` |
| `material_code` | string | 是 | `HCTM-BCP-A01` |
| `demand_qty` | number | 是 | `12000` |
| `net_production_qty` | number | 是 | `10800` |
| `material_shortages` | array | 是 | 见下方 |
| `capacity_risks` | array | 否 | 见下方 |
| `recommended_actions` | array | 否 | 见下方 |

`material_shortages` item：

```json
{
  "material_code": "AL-PLATE-6061-T6",
  "required_qty": 12600,
  "available_qty": 8000,
  "in_transit_qty": 5000,
  "shortage_qty": 0,
  "timing_risk": "supplier_eta_delayed"
}
```

`capacity_risks` item：

```json
{
  "work_center_code": "WLD-02",
  "equipment_code": "LAS-WLD-02",
  "risk_type": "machine_down",
  "impact_qty": 280,
  "risk_level": "high"
}
```

下游：

- 生成采购需求。
- 生成生产订单建议。
- 触发计划 Agent 输出方案。

Agent：

- 计划 Agent 主触发事件。
- 经营分析 Agent 可监听重算结果。

幂等键：

- `{tenant_id}:o2d.mrp.generated:{mrp_run_id}`

### 5.4 `PurchaseRequirementCreated`

Event type：

- `o2d.purchase_requirement.created`

触发条件：

- MRP 发现物料数量、时点或风险缺口。

主对象：

- `sales_order` 或未来 `purchase_requirement`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `requirement_no` | string | 是 | `PR-202607-0001` |
| `mrp_run_id` | string | 是 | `MRP-202607-0001` |
| `material_code` | string | 是 | `AL-PLATE-6061-T6` |
| `required_qty` | number | 是 | `2000` |
| `need_by_date` | date | 是 | `2026-07-12` |
| `preferred_supplier_code` | string | 否 | `SUP-BETA-AL` |
| `reason` | enum | 是 | `timing_risk` |

下游：

- 采购员创建或释放采购订单。
- 可触发采购审批或加急流程。

Agent：

- 计划 Agent。

幂等键：

- `{tenant_id}:o2d.purchase_requirement.created:{requirement_no}`

### 5.5 `PurchaseOrderReleased`

Event type：

- `o2d.purchase_order.released`

触发条件：

- 采购订单发布给供应商。

主对象：

- `purchase_order`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `po_no` | string | 是 | `PO-202607-0001` |
| `supplier_code` | string | 是 | `SUP-ALPHA-AL` |
| `material_code` | string | 是 | `AL-PLATE-6061-T6` |
| `order_qty` | number | 是 | `5000` |
| `promised_date` | date | 是 | `2026-07-08` |
| `plant_code` | string | 是 | `PLT-SZ` |
| `expedite_flag` | boolean | 否 | `false` |

下游：

- 等待供应商发货。
- 更新供应计划。

Agent：

- 计划 Agent 监听交期承诺。

幂等键：

- `{tenant_id}:o2d.purchase_order.released:{po_no}`

### 5.6 `SupplierShipmentDispatched`

Event type：

- `o2d.supplier_shipment.dispatched`

触发条件：

- 供应商发货或仿真器生成发货事件。

主对象：

- `purchase_order`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `po_no` | string | 是 | `PO-202607-0001` |
| `supplier_code` | string | 是 | `SUP-ALPHA-AL` |
| `supplier_lot_no` | string | 是 | `ALPHA-20260711-01` |
| `shipped_qty` | number | 是 | `5000` |
| `eta` | date | 是 | `2026-07-08` |
| `carrier` | string | 否 | 安铝物流 |

下游：

- 仓库准备收货。
- 计划 Agent 更新预计可用库存。

幂等键：

- `{tenant_id}:o2d.supplier_shipment.dispatched:{po_no}:{supplier_lot_no}`

### 5.7 `SupplierDeliveryDelayed`

Event type：

- `o2d.supplier_delivery.delayed`

触发条件：

- 供应商更新 ETA。
- 到期未到货。

主对象：

- `purchase_order`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `po_no` | string | 是 | `PO-202607-0001` |
| `supplier_code` | string | 是 | `SUP-ALPHA-AL` |
| `material_code` | string | 是 | `AL-PLATE-6061-T6` |
| `original_eta` | date | 是 | `2026-07-08` |
| `new_eta` | date | 是 | `2026-07-11` |
| `delay_days` | integer | 是 | `3` |
| `reason` | string | 否 | 供应商熔炼排产延迟 |
| `affected_order_no` | string | 否 | `SO-202607-0001` |

下游：

- 触发 MRP 重算。
- 触发采购加急或备选供应商建议。

Agent：

- 计划 Agent 强触发。
- 经营分析 Agent 计算交付和成本影响。

幂等键：

- `{tenant_id}:o2d.supplier_delivery.delayed:{po_no}:{new_eta}`

### 5.8 `MaterialReceived`

Event type：

- `whs.material.received`

触发条件：

- 仓库完成供应商到货收货。

主对象：

- `goods_receipt`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `receipt_no` | string | 是 | `GR-202607-0001` |
| `po_no` | string | 是 | `PO-202607-0001` |
| `supplier_code` | string | 是 | `SUP-ALPHA-AL` |
| `material_code` | string | 是 | `AL-PLATE-6061-T6` |
| `supplier_lot_no` | string | 是 | `ALPHA-20260711-01` |
| `received_qty` | number | 是 | `5000` |
| `received_at` | datetime | 是 | `2026-07-11T10:00:00+08:00` |

下游：

- 创建来料检验单。
- 如果物料免检，后续可直接入库；MVP 铝板必须检验。

Agent：

- 质量 Agent 监听。

幂等键：

- `{tenant_id}:whs.material.received:{receipt_no}`

### 5.9 `IncomingInspectionPassed`

Event type：

- `qms.incoming_inspection.passed`

触发条件：

- 来料检验单判定合格。

主对象：

- `inspection_order`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `inspection_no` | string | 是 | `IQC-202607-0001` |
| `receipt_no` | string | 是 | `GR-202607-0001` |
| `material_code` | string | 是 | `AL-PLATE-6061-T6` |
| `lot_no` | string | 是 | `ALPHA-20260711-01` |
| `accepted_qty` | number | 是 | `5000` |
| `inspection_level` | enum | 是 | `normal` |
| `inspector_employee_code` | string | 是 | `EMP-SZ-IQC-001` |

下游：

- 允许生成入库事务。

幂等键：

- `{tenant_id}:qms.incoming_inspection.passed:{inspection_no}`

### 5.10 `IncomingInspectionFailed`

Event type：

- `qms.incoming_inspection.failed`

触发条件：

- 来料检验单判定失败。

主对象：

- `inspection_order`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `inspection_no` | string | 是 | `IQC-202607-0002` |
| `receipt_no` | string | 是 | `GR-202607-0002` |
| `supplier_code` | string | 是 | `SUP-BETA-AL` |
| `material_code` | string | 是 | `AL-PLATE-6061-T6` |
| `lot_no` | string | 是 | `BETA-20260712-01` |
| `rejected_qty` | number | 是 | `300` |
| `defect_code` | string | 是 | `SURFACE_SCRATCH` |
| `severity` | enum | 是 | `major` |

下游：

- 创建 `quality_issue`。
- 隔离批次。
- 触发计划重算或让步接收流程。

Agent：

- 质量 Agent 强触发。
- 计划 Agent 评估物料可用性影响。

幂等键：

- `{tenant_id}:qms.incoming_inspection.failed:{inspection_no}:{defect_code}`

### 5.11 `InventoryPutawayCompleted`

Event type：

- `whs.inventory.putaway_completed`

触发条件：

- 合格物料或成品完成上架。

主对象：

- `inventory_transaction`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `transaction_no` | string | 是 | `INV-202607-0001` |
| `source_object_type` | enum | 是 | `goods_receipt` |
| `source_object_no` | string | 是 | `GR-202607-0001` |
| `warehouse_code` | string | 是 | `WH-SZ-RM` |
| `location_code` | string | 是 | `RM-A01-01` |
| `material_code` | string | 是 | `AL-PLATE-6061-T6` |
| `lot_no` | string | 是 | `ALPHA-20260711-01` |
| `qty` | number | 是 | `5000` |

下游：

- 更新可用库存。
- 触发生产领料或计划重算。

Agent：

- 计划 Agent 更新物料可用性。

幂等键：

- `{tenant_id}:whs.inventory.putaway_completed:{transaction_no}`

### 5.12 `ProductionOrderReleased`

Event type：

- `proc.production_order.released`

触发条件：

- 计划员释放生产订单。

主对象：

- `production_order`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `work_order_no` | string | 是 | `WO-202607-0001` |
| `sales_order_no` | string | 否 | `SO-202607-0001` |
| `material_code` | string | 是 | `HCTM-BCP-A01` |
| `routing_code` | string | 是 | `RT-HCTM-BCP-A01-V1` |
| `planned_qty` | number | 是 | `10800` |
| `due_date` | date | 是 | `2026-07-20` |
| `plant_code` | string | 是 | `PLT-SZ` |

下游：

- 生成工序任务。
- 产线状态进入待生产。

Agent：

- 计划 Agent 监听。

幂等键：

- `{tenant_id}:proc.production_order.released:{work_order_no}`

### 5.13 `OperationStarted`

Event type：

- `proc.operation.started`

触发条件：

- 班组开始执行工序任务。

主对象：

- `operation_task`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `task_no` | string | 是 | `OPT-202607-0001-030` |
| `work_order_no` | string | 是 | `WO-202607-0001` |
| `operation_code` | string | 是 | `OP-WLD` |
| `work_center_code` | string | 是 | `WLD-02` |
| `equipment_code` | string | 否 | `LAS-WLD-02` |
| `team_code` | string | 否 | `TEAM-SZ-BCP-A-DAY` |
| `started_at` | datetime | 是 | `2026-07-14T08:00:00+08:00` |

下游：

- 更新工序任务状态。
- 更新沙盘产线状态。

幂等键：

- `{tenant_id}:proc.operation.started:{task_no}:{started_at}`

### 5.14 `MachineDown`

Event type：

- `eam.machine.down`

触发条件：

- 设备故障停机，或仿真器注入设备异常。

主对象：

- `equipment`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `equipment_code` | string | 是 | `LAS-WLD-02` |
| `work_center_code` | string | 是 | `WLD-02` |
| `plant_code` | string | 是 | `PLT-SZ` |
| `down_at` | datetime | 是 | `2026-07-12T09:30:00+08:00` |
| `expected_repair_at` | datetime | 否 | `2026-07-12T17:30:00+08:00` |
| `expected_down_hours` | number | 是 | `8` |
| `capacity_loss_ratio` | number | 是 | `0.35` |
| `reason_code` | string | 否 | `LASER_POWER_ALARM` |
| `affected_task_no` | string | 否 | `OPT-202607-0001-030` |

下游：

- 更新设备状态为 `down`。
- 触发计划重算。
- 可生成维修任务，MVP 可先不建独立维修单。

Agent：

- 计划 Agent 强触发。
- 经营分析 Agent 计算加班和交付影响。

幂等键：

- `{tenant_id}:eam.machine.down:{equipment_code}:{down_at}`

### 5.15 `OperationCompleted`

Event type：

- `proc.operation.completed`

触发条件：

- 工序任务完成报工。

主对象：

- `operation_task`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `task_no` | string | 是 | `OPT-202607-0001-030` |
| `work_order_no` | string | 是 | `WO-202607-0001` |
| `operation_code` | string | 是 | `OP-WLD` |
| `work_center_code` | string | 是 | `WLD-02` |
| `completed_qty` | number | 是 | `10520` |
| `scrap_qty` | number | 否 | `80` |
| `completed_at` | datetime | 是 | `2026-07-17T20:00:00+08:00` |

下游：

- 推进下一工序。
- 若 `scrap_qty` 超阈值，可触发质量分析。

Agent：

- 质量 Agent 监听报废和不良波动。

幂等键：

- `{tenant_id}:proc.operation.completed:{task_no}:{completed_at}`

### 5.16 `ProcessInspectionFailed`

Event type：

- `qms.process_inspection.failed`

触发条件：

- 过程检验失败，例如气密检漏不合格。

主对象：

- `inspection_order`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `inspection_no` | string | 是 | `PQC-202607-0001` |
| `task_no` | string | 是 | `OPT-202607-0001-050` |
| `work_order_no` | string | 是 | `WO-202607-0001` |
| `material_code` | string | 是 | `HCTM-BCP-A01` |
| `lot_no` | string | 否 | `FG-20260717-01` |
| `affected_qty` | number | 是 | `120` |
| `defect_code` | string | 是 | `LEAK_RATE_HIGH` |
| `severity` | enum | 是 | `major` |

下游：

- 创建质量问题。
- 隔离批次。
- 触发返工或报废决策。

Agent：

- 质量 Agent 强触发。

幂等键：

- `{tenant_id}:qms.process_inspection.failed:{inspection_no}:{defect_code}`

### 5.17 `FinishedGoodsReceived`

Event type：

- `whs.finished_goods.received`

触发条件：

- 成品完成入库。

主对象：

- `inventory_transaction`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `transaction_no` | string | 是 | `INV-FG-202607-0001` |
| `work_order_no` | string | 是 | `WO-202607-0001` |
| `material_code` | string | 是 | `HCTM-BCP-A01` |
| `lot_no` | string | 是 | `FG-20260717-01` |
| `warehouse_code` | string | 是 | `WH-SZ-FG` |
| `location_code` | string | 是 | `FG-A01-01` |
| `received_qty` | number | 是 | `10500` |
| `received_at` | datetime | 是 | `2026-07-18T10:00:00+08:00` |

下游：

- 增加可发货库存。
- 触发发货计划。

Agent：

- 经营分析 Agent 计算交付可行性。

幂等键：

- `{tenant_id}:whs.finished_goods.received:{transaction_no}`

### 5.18 `ShipmentDispatched`

Event type：

- `o2d.shipment.dispatched`

触发条件：

- 成品发货出厂。

主对象：

- `shipment`

Payload：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `shipment_no` | string | 是 | `SHIP-202607-0001` |
| `sales_order_no` | string | 是 | `SO-202607-0001` |
| `customer_code` | string | 是 | `CUST-SGNEV` |
| `material_code` | string | 是 | `HCTM-BCP-A01` |
| `ship_qty` | number | 是 | `9000` |
| `ship_date` | date | 是 | `2026-07-19` |
| `batch_lots` | string[] | 否 | `FG-20260717-01` |
| `split_shipment_flag` | boolean | 是 | `true` |

下游：

- 更新销售订单发货状态。
- 触发开票或收入预测。
- 更新经营 KPI。

Agent：

- 经营分析 Agent 强触发。

幂等键：

- `{tenant_id}:o2d.shipment.dispatched:{shipment_no}`

## 6. Agent 触发矩阵

| 事件 | 计划 Agent | 质量 Agent | 经营分析 Agent |
| --- | --- | --- | --- |
| `o2d.order.confirmed` | 强触发 | - | 监听 |
| `o2d.mrp.generated` | 强触发 | - | 监听 |
| `o2d.supplier_delivery.delayed` | 强触发 | - | 强触发 |
| `qms.incoming_inspection.failed` | 监听 | 强触发 | 监听 |
| `eam.machine.down` | 强触发 | - | 强触发 |
| `qms.process_inspection.failed` | 监听 | 强触发 | 监听 |
| `whs.finished_goods.received` | 监听 | - | 监听 |
| `o2d.shipment.dispatched` | - | - | 强触发 |

Agent 输出约束：

- 第一阶段 Agent 只生成建议、风险解释和草稿动作。
- 任何会改变业务状态的动作必须通过 IAOS Capability / Process 受治理执行。
- Agent 输出必须引用 `correlation_id` 和相关业务对象编号。

## 7. Capability 和 Process 接线建议

| 事件 | 建议 Capability / Process | 说明 |
| --- | --- | --- |
| `o2d.order.confirmed` | `o2d.mrp.run` | 包装 BOM 展开、库存检查和产能检查 |
| `o2d.mrp.generated` | `planning.generate_supply_plan` | 生成采购需求和生产订单建议 |
| `o2d.purchase_requirement.created` | `procurement.create_purchase_order_draft` | 创建采购订单草稿 |
| `o2d.supplier_delivery.delayed` | `planning.recalculate_delivery_risk` | 重算交付风险 |
| `qms.incoming_inspection.failed` | `quality.create_containment_issue` | 创建质量问题和隔离建议 |
| `whs.inventory.putaway_completed` | `inventory.update_available_stock` | 更新可用库存 |
| `proc.production_order.released` | `manufacturing.create_operation_tasks` | 生成工序任务 |
| `eam.machine.down` | `planning.recalculate_capacity_risk` | 重算产能风险 |
| `qms.process_inspection.failed` | `quality.create_process_issue` | 创建过程质量问题 |
| `o2d.shipment.dispatched` | `finance.prepare_invoice_draft` | 开票草稿，MVP 可只做事件记录 |

说明：

- 上表是 AESE 建议能力名，不表示 IAOS 当前已经全部实现。
- 已存在的 O2D 能力可优先复用当前 `o2d.mrp.*` handler。
- 新能力应先以只读或草稿形式落地，避免 Agent 直接改变关键业务状态。

## 8. 订阅建议

MVP 服务可以按领域订阅：

```text
O2D 场景服务:
iaos.*.o2d.order.confirmed
iaos.*.o2d.mrp.generated
iaos.*.o2d.supplier_delivery.delayed

仓储服务:
iaos.*.whs.material.received
iaos.*.qms.incoming_inspection.passed
iaos.*.whs.finished_goods.received

质量服务:
iaos.*.whs.material.received
iaos.*.qms.incoming_inspection.failed
iaos.*.qms.process_inspection.failed

设备 / 计划服务:
iaos.*.eam.machine.down
iaos.*.proc.operation.*

经营分析服务:
iaos.*.o2d.*
iaos.*.whs.finished_goods.received
iaos.*.eam.machine.down
iaos.*.qms.*.failed
```

## 9. MVP 示例事件

示例：供应商延期事件。

```json
{
  "id": "evt-20260708-0007",
  "tenant_id": "tenant-hctm",
  "type": "o2d.supplier_delivery.delayed",
  "source": "aese.hctm.simulator",
  "timestamp": "2026-07-08T09:00:00+08:00",
  "data": {
    "po_no": "PO-202607-0001",
    "supplier_code": "SUP-ALPHA-AL",
    "material_code": "AL-PLATE-6061-T6",
    "original_eta": "2026-07-08",
    "new_eta": "2026-07-11",
    "delay_days": 3,
    "reason": "供应商熔炼排产延迟",
    "affected_order_no": "SO-202607-0001"
  },
  "metadata": {
    "correlation_id": "corr-so-202607-0001",
    "causation_id": "evt-20260701-0005",
    "scenario": "hctm_order_expedite_01",
    "plant_code": "PLT-SZ",
    "business_object_type": "purchase_order",
    "business_object_id": "PO-202607-0001",
    "idempotency_key": "tenant-hctm:o2d.supplier_delivery.delayed:PO-202607-0001:2026-07-11"
  }
}
```

## 10. 后续落地顺序

建议下一步：

1. 基于本文创建 `shared/eventdef` 常量草案。
2. 为 18 个事件补 JSON Schema 或 Go payload struct。
3. 扩展 O2D scenario 订阅 `o2d.mrp.generated`、`o2d.supplier_delivery.delayed` 和 `eam.machine.down`。
4. 在 AESE seed 数据中准备第一条演示故事的事件序列。
5. 将计划 Agent 的触发上下文限制在 `correlation_id = corr-so-202607-0001` 的事件集合内。

