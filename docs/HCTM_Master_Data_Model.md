# HCTM Master Data Model

本文把华辰热管理系统集团 MVP 蓝图中的 28 个关键对象转成 IAOS 可落地的主数据和业务对象建模规格。目标是作为后续 metadata entity、seed 数据、页面、权限、事件 payload 和 Agent 工具上下文的共同依据。

## 1. 建模原则

MVP 建模遵循以下原则：

- 先支持订单到交付主线，不追求完整 ERP/MES 覆盖。
- 对 IAOS 暴露业务语义字段，不暴露物理影子列。
- 所有业务对象必须带 `tenant_id`，运行时由 IAOS RLS 注入和隔离。
- 编码字段保持稳定，用于 seed、事件 payload、演示脚本和 Agent 引用。
- 状态字段使用明确枚举，避免用自由文本表达流程状态。
- 关系优先使用业务编码和对象引用，后续落地时可补充内部 ID。
- 事件 payload 中引用对象时优先携带 `*_code`、`*_no` 和 `business_object_id`。

## 2. 命名约定

实体命名：

| 类型 | 约定 | 示例 |
| --- | --- | --- |
| Entity key | snake_case 英文 | `sales_order` |
| Display name | 中文业务名 | 销售订单 |
| Code field | `*_code` | `plant_code` |
| Document number | `*_no` | `order_no` |
| Status field | `status` | `released` |
| Date field | `*_date` | `due_date` |
| Timestamp field | `*_at` | `confirmed_at` |

MVP 编码前缀：

| 对象 | 前缀 | 示例 |
| --- | --- | --- |
| 集团 | `GRP` | `GRP-HCTM` |
| 事业部 | `BU` | `BU-EAST` |
| 法人 | `LE` | `LE-HCTM-SH` |
| 工厂 | `PLT` | `PLT-SZ` |
| 部门 | `DEPT` | `DEPT-SZ-PLAN` |
| 班组 | `TEAM` | `TEAM-SZ-BCP-A-DAY` |
| 客户 | `CUST` | `CUST-SGNEV` |
| 供应商 | `SUP` | `SUP-ALPHA-AL` |
| 物料 | 业务编码 | `HCTM-BCP-A01` |
| 工作中心 | 工艺编码 | `WLD-02` |
| 设备 | 设备编码 | `LAS-WLD-02` |
| 销售订单 | `SO` | `SO-202607-0001` |
| 采购订单 | `PO` | `PO-202607-0001` |
| 生产订单 | `WO` | `WO-202607-0001` |

## 3. 通用字段

所有 entity 建议具备：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `tenant_id` | string | 是 | IAOS 租户。MVP 示例为 `tenant-hctm` |
| `created_at` | datetime | 是 | 创建时间 |
| `created_by` | string | 是 | 创建人或系统用户 |
| `updated_at` | datetime | 是 | 更新时间 |
| `updated_by` | string | 是 | 更新人或系统用户 |
| `status` | enum | 视对象而定 | 主状态 |
| `remark` | text | 否 | 业务备注 |

建议先作为 domain fields 建模，待 IAOS 平台已有系统审计字段可复用时再收敛。

## 4. Entity 总览

| 序号 | Entity key | 中文名 | 类型 | MVP | 主要用途 |
| --- | --- | --- | --- | --- | --- |
| 1 | `enterprise_group` | 集团 | Organization | 必须 | 企业最高组织 |
| 2 | `business_unit` | 事业部 | Organization | 必须 | 区域经营单元 |
| 3 | `legal_entity` | 法人公司 | Organization | 必须 | 签约和财务主体 |
| 4 | `plant` | 工厂 | Facility | 必须 | 制造运营主体 |
| 5 | `department` | 部门 | Organization | 必须 | 组织分工和权限边界 |
| 6 | `production_team` | 班组 | Workforce | 必须 | 生产执行组织 |
| 7 | `customer` | 客户 | Party | 必须 | 销售对象 |
| 8 | `supplier` | 供应商 | Party | 必须 | 采购来源 |
| 9 | `material` | 物料 | Item | 必须 | 原料、半成品、成品、包材 |
| 10 | `bom` | BOM | Engineering | 必须 | 产品结构 |
| 11 | `routing` | 工艺路线 | Engineering | 必须 | 产品制造路径 |
| 12 | `operation` | 工序 | Engineering | 必须 | 路线步骤 |
| 13 | `work_center` | 工作中心 | Resource | 必须 | 能力核算单元 |
| 14 | `equipment` | 设备 | Asset | 必须 | 具体生产资产 |
| 15 | `tooling` | 工装 | Asset | 可选 | 辅助生产资源 |
| 16 | `warehouse` | 仓库 | Inventory | 必须 | 库存管理地点 |
| 17 | `storage_location` | 库位 | Inventory | 必须 | 库存最小位置 |
| 18 | `shift` | 班次 | Workforce | 必须 | 产能时间模型 |
| 19 | `employee` | 员工 | Workforce | 必须 | 岗位执行者 |
| 20 | `sales_order` | 销售订单 | Document | 必须 | 客户需求 |
| 21 | `purchase_order` | 采购订单 | Document | 必须 | 对供应商采购 |
| 22 | `goods_receipt` | 收货单 | Document | 必须 | 供应商到货记录 |
| 23 | `inspection_order` | 检验单 | Document | 必须 | 质量判定 |
| 24 | `inventory_transaction` | 库存事务 | Transaction | 必须 | 库存变化流水 |
| 25 | `production_order` | 生产订单 | Document | 必须 | 工厂生产指令 |
| 26 | `operation_task` | 工序任务 | Task | 必须 | 现场执行任务 |
| 27 | `shipment` | 发货单 | Document | 必须 | 客户交付记录 |
| 28 | `quality_issue` | 质量问题 | Issue | 必须 | 异常质量对象 |

## 5. 组织与资源对象

### 5.1 `enterprise_group` 集团

用途：

- 表达华辰热管理系统集团这个最高组织。

关键字段：

| 字段 | 类型 | 必填 | 示例 | 说明 |
| --- | --- | --- | --- | --- |
| `group_code` | string | 是 | `GRP-HCTM` | 集团编码，唯一 |
| `group_name` | string | 是 | 华辰热管理系统集团有限公司 | 集团名称 |
| `english_name` | string | 否 | Huachen Thermal Management Systems Group Co., Ltd. | 英文名 |
| `industry` | enum | 是 | `nev_thermal_management` | 行业 |
| `headquarter_city` | string | 否 | 上海 | 总部城市 |

关系：

- 一对多 `business_unit`。
- 一对多 `legal_entity`。

状态：

- `active`
- `inactive`

Seed：

- `GRP-HCTM` 华辰热管理系统集团有限公司。

### 5.2 `business_unit` 事业部

用途：

- 表达华东、华南、海外等经营单元。MVP 只启用华东事业部。

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `bu_code` | string | 是 | `BU-EAST` |
| `bu_name` | string | 是 | 华东事业部 |
| `group_code` | ref | 是 | `GRP-HCTM` |
| `region` | string | 是 | 华东 |
| `manager_employee_code` | ref | 否 | `EMP-SZ-GM-001` |

关系：

- 多对一 `enterprise_group`。
- 一对多 `plant`。
- 一对多 `department`。

状态：

- `active`
- `inactive`

Seed：

- `BU-EAST` 华东事业部。

### 5.3 `legal_entity` 法人公司

用途：

- 表达签约、发票和财务责任主体。

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `legal_entity_code` | string | 是 | `LE-HCTM-SH` |
| `legal_entity_name` | string | 是 | 华辰热管理系统上海有限公司 |
| `group_code` | ref | 是 | `GRP-HCTM` |
| `tax_id` | string | 否 | `91310000HCTM000001` |
| `currency` | enum | 是 | `CNY` |
| `invoice_title` | string | 否 | 华辰热管理系统上海有限公司 |

关系：

- 多对一 `enterprise_group`。
- 一对多 `sales_order`。
- 一对多 `purchase_order`。

状态：

- `active`
- `inactive`

Seed：

- `LE-HCTM-SH` 华辰热管理系统上海有限公司。

### 5.4 `plant` 工厂

用途：

- 表达制造基地和工厂运营边界。

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `plant_code` | string | 是 | `PLT-SZ` |
| `plant_name` | string | 是 | 苏州制造基地 |
| `bu_code` | ref | 是 | `BU-EAST` |
| `legal_entity_code` | ref | 是 | `LE-HCTM-SH` |
| `city` | string | 是 | 苏州 |
| `plant_type` | enum | 是 | `manufacturing` |
| `timezone` | string | 是 | `Asia/Shanghai` |

关系：

- 多对一 `business_unit`。
- 一对多 `department`、`warehouse`、`work_center`、`production_order`。

状态：

- `active`
- `maintenance`
- `inactive`

Seed：

- `PLT-SZ` 苏州制造基地。

### 5.5 `department` 部门

用途：

- 表达组织分工、权限边界和岗位归属。

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `dept_code` | string | 是 | `DEPT-SZ-PLAN` |
| `dept_name` | string | 是 | 计划物流部 |
| `plant_code` | ref | 否 | `PLT-SZ` |
| `bu_code` | ref | 否 | `BU-EAST` |
| `dept_type` | enum | 是 | `planning` |
| `parent_dept_code` | ref | 否 | `DEPT-SZ-MGMT` |

关系：

- 可归属 `plant` 或 `business_unit`。
- 一对多 `employee`。

状态：

- `active`
- `inactive`

Seed：

- `DEPT-SZ-PLAN` 计划物流部。
- `DEPT-SZ-QUALITY` 质量管理部。
- `DEPT-SZ-MFG` 生产制造部。
- `DEPT-SZ-WH` 仓储物流部。
- `DEPT-SZ-EAM` 设备工程部。

### 5.6 `production_team` 班组

用途：

- 表达现场生产执行组织和班次绑定。

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `team_code` | string | 是 | `TEAM-SZ-BCP-A-DAY` |
| `team_name` | string | 是 | 电池冷却板 A 线白班 |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `dept_code` | ref | 是 | `DEPT-SZ-MFG` |
| `shift_code` | ref | 是 | `SHIFT-DAY` |
| `line_code` | string | 是 | `SZ-BCP-LINE-A` |
| `leader_employee_code` | ref | 否 | `EMP-SZ-TL-001` |

关系：

- 多对一 `plant`、`department`、`shift`。
- 一对多 `operation_task`。

状态：

- `active`
- `paused`
- `inactive`

Seed：

- `TEAM-SZ-BCP-A-DAY` 电池冷却板 A 线白班。
- `TEAM-SZ-BCP-A-NIGHT` 电池冷却板 A 线夜班。

## 6. 交易方与工程主数据

### 6.1 `customer` 客户

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `customer_code` | string | 是 | `CUST-SGNEV` |
| `customer_name` | string | 是 | 星河新能源汽车有限公司 |
| `customer_type` | enum | 是 | `oem` |
| `priority` | enum | 是 | `high` |
| `payment_terms` | string | 否 | `NET60` |
| `delivery_mode` | enum | 是 | `jit_split` |
| `default_ship_to` | string | 否 | 星河苏州整车工厂 |

关系：

- 一对多 `sales_order`。
- 可关联客户料号到 `material.customer_material_no`。

状态：

- `active`
- `on_hold`
- `inactive`

Seed：

- `CUST-SGNEV` 星河新能源汽车有限公司。
- `CUST-ORBITHM` 欧瑞热管理系统有限公司。

### 6.2 `supplier` 供应商

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `supplier_code` | string | 是 | `SUP-ALPHA-AL` |
| `supplier_name` | string | 是 | 安铝金属材料有限公司 |
| `supplier_type` | enum | 是 | `material` |
| `rating` | enum | 是 | `A` |
| `default_lead_time_days` | integer | 是 | `7` |
| `certified_materials` | string[] | 否 | `AL-PLATE-6061-T6` |
| `inspection_level` | enum | 是 | `normal` |

关系：

- 一对多 `purchase_order`。
- 多对多 `material`，后续可拆 `supplier_material`.

状态：

- `active`
- `qualified`
- `restricted`
- `inactive`

Seed：

- `SUP-ALPHA-AL` 安铝金属材料有限公司。
- `SUP-BETA-AL` 贝塔铝业有限公司。
- `SUP-SEAL-01` 恒密密封科技有限公司。
- `SUP-FIT-01` 鼎联接头制造有限公司。
- `SUP-PKG-01` 诚达包装材料有限公司。

### 6.3 `material` 物料

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `material_code` | string | 是 | `HCTM-BCP-A01` |
| `material_name` | string | 是 | 电池冷却板组件 A01 |
| `material_type` | enum | 是 | `finished_good` |
| `uom` | enum | 是 | `pcs` |
| `lot_control` | boolean | 是 | `true` |
| `serial_control` | boolean | 是 | `false` |
| `customer_material_no` | string | 否 | `SGNEV-CP-7788` |
| `default_plant_code` | ref | 否 | `PLT-SZ` |
| `safety_stock_qty` | decimal | 否 | `1200` |

关系：

- 被 `bom`、`sales_order`、`purchase_order`、`inventory_transaction`、`production_order` 引用。

状态：

- `active`
- `engineering`
- `phase_out`
- `inactive`

Seed：

- `HCTM-BCP-A01` 电池冷却板组件 A01。
- `AL-PLATE-6061-T6` 6061-T6 铝板。
- `SEAL-RING-A01` 密封圈。
- `FITTING-A01` 接头组件。
- `PKG-BOX-A01` 包装箱。
- `LBL-TRACE-A01` 追溯标签。

### 6.4 `bom` BOM

MVP 建议采用头行合一的简化 entity，后续可拆 `bom_header` 和 `bom_line`。

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `bom_code` | string | 是 | `BOM-HCTM-BCP-A01-V1` |
| `parent_material_code` | ref | 是 | `HCTM-BCP-A01` |
| `component_material_code` | ref | 是 | `AL-PLATE-6061-T6` |
| `qty_per` | decimal | 是 | `1.05` |
| `uom` | enum | 是 | `sheet` |
| `scrap_rate` | decimal | 否 | `0.02` |
| `effective_from` | date | 是 | `2026-07-01` |
| `effective_to` | date | 否 |  |

关系：

- 多对一 parent `material`。
- 多对一 component `material`。

状态：

- `draft`
- `active`
- `inactive`

Seed：

- `HCTM-BCP-A01` 对 `AL-PLATE-6061-T6` 用量 `1.05`。
- `HCTM-BCP-A01` 对 `SEAL-RING-A01` 用量 `2`。
- `HCTM-BCP-A01` 对 `FITTING-A01` 用量 `2`。

### 6.5 `routing` 工艺路线

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `routing_code` | string | 是 | `RT-HCTM-BCP-A01-V1` |
| `material_code` | ref | 是 | `HCTM-BCP-A01` |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `version` | string | 是 | `V1` |
| `effective_from` | date | 是 | `2026-07-01` |
| `standard_daily_capacity` | decimal | 是 | `800` |

关系：

- 一对多 `operation`。
- 多对一 `material`、`plant`。

状态：

- `draft`
- `active`
- `inactive`

Seed：

- `RT-HCTM-BCP-A01-V1` 电池冷却板 A01 苏州路线。

### 6.6 `operation` 工序

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `operation_code` | string | 是 | `OP-WLD` |
| `routing_code` | ref | 是 | `RT-HCTM-BCP-A01-V1` |
| `sequence_no` | integer | 是 | `30` |
| `operation_name` | string | 是 | 焊接 |
| `work_center_code` | ref | 是 | `WLD-02` |
| `standard_time_min` | decimal | 是 | `1.2` |
| `inspection_required` | boolean | 是 | `false` |
| `yield_rate` | decimal | 否 | `0.98` |

关系：

- 多对一 `routing`。
- 多对一 `work_center`。
- 一对多 `operation_task`。

状态：

- `active`
- `inactive`

Seed：

- `OP-IQC` 来料检验。
- `OP-CNC` CNC 精加工。
- `OP-WLD` 焊接。
- `OP-CLN` 清洗。
- `OP-LKT` 气密检漏。
- `OP-DIM` 尺寸检测。
- `OP-ASM` 总装包装。

## 7. 产能、库存与人员主数据

### 7.1 `work_center` 工作中心

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `work_center_code` | string | 是 | `WLD-02` |
| `work_center_name` | string | 是 | 焊接工作中心 2 |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `center_type` | enum | 是 | `production` |
| `line_code` | string | 是 | `SZ-BCP-LINE-A` |
| `capacity_per_shift` | decimal | 是 | `400` |
| `bottleneck_flag` | boolean | 是 | `true` |

关系：

- 一对多 `equipment`。
- 一对多 `operation`。

状态：

- `active`
- `limited`
- `down`
- `inactive`

Seed：

- `IQC-01`、`CNC-01`、`WLD-01`、`WLD-02`、`CLN-01`、`LKT-01`、`DIM-01`、`ASM-01`。

### 7.2 `equipment` 设备

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `equipment_code` | string | 是 | `LAS-WLD-02` |
| `equipment_name` | string | 是 | 激光焊接机 2 |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `work_center_code` | ref | 是 | `WLD-02` |
| `equipment_type` | enum | 是 | `laser_welder` |
| `criticality` | enum | 是 | `high` |
| `commissioned_date` | date | 否 | `2024-05-01` |
| `current_status` | enum | 是 | `running` |

关系：

- 多对一 `plant`、`work_center`。
- 被 `operation_task` 和 `quality_issue` 追溯引用。

状态：

- `running`
- `idle`
- `down`
- `maintenance`
- `inactive`

Seed：

- `LAS-WLD-02` 是演示故事中的故障设备。

### 7.3 `tooling` 工装

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `tooling_code` | string | 是 | `FIX-BCP-A01-WLD` |
| `tooling_name` | string | 是 | 电池冷却板焊接夹具 |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `work_center_code` | ref | 否 | `WLD-02` |
| `lifetime_cycles` | integer | 否 | `50000` |
| `used_cycles` | integer | 否 | `12000` |
| `current_status` | enum | 是 | `available` |

关系：

- 多对一 `work_center`。

状态：

- `available`
- `in_use`
- `maintenance`
- `scrapped`

MVP 说明：

- 可选对象。第一阶段可只 seed 一个焊接夹具用于后续扩展。

### 7.4 `warehouse` 仓库

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `warehouse_code` | string | 是 | `WH-SZ-RM` |
| `warehouse_name` | string | 是 | 苏州原材料仓 |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `warehouse_type` | enum | 是 | `raw_material` |
| `managed_by_dept_code` | ref | 否 | `DEPT-SZ-WH` |

关系：

- 一对多 `storage_location`。
- 一对多 `inventory_transaction`。

状态：

- `active`
- `locked`
- `inactive`

Seed：

- `WH-SZ-RM` 原材料仓。
- `WH-SZ-WIP` 半成品暂存区。
- `WH-SZ-FG` 成品仓。
- `WH-SZ-NG` 不合格品区。

### 7.5 `storage_location` 库位

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `location_code` | string | 是 | `RM-A01-01` |
| `warehouse_code` | ref | 是 | `WH-SZ-RM` |
| `location_name` | string | 是 | 原材料 A 区 01-01 |
| `location_type` | enum | 是 | `normal` |
| `material_restriction` | string | 否 | `AL-PLATE-6061-T6` |

关系：

- 多对一 `warehouse`。
- 一对多 `inventory_transaction`。

状态：

- `available`
- `blocked`
- `inactive`

Seed：

- `RM-A01-01` 铝板库位。
- `FG-A01-01` 冷却板成品库位。
- `NG-QA-01` 不合格品隔离库位。

### 7.6 `shift` 班次

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `shift_code` | string | 是 | `SHIFT-DAY` |
| `shift_name` | string | 是 | 白班 |
| `start_time` | time | 是 | `08:00` |
| `end_time` | time | 是 | `16:00` |
| `work_hours` | decimal | 是 | `8` |
| `overtime_allowed` | boolean | 是 | `true` |

关系：

- 一对多 `production_team`。
- 被产能计算引用。

状态：

- `active`
- `inactive`

Seed：

- `SHIFT-DAY` 白班。
- `SHIFT-NIGHT` 夜班。
- `SHIFT-OT` 加班时段。

### 7.7 `employee` 员工

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `employee_code` | string | 是 | `EMP-SZ-PLAN-001` |
| `employee_name` | string | 是 | 李计划 |
| `dept_code` | ref | 是 | `DEPT-SZ-PLAN` |
| `plant_code` | ref | 否 | `PLT-SZ` |
| `role_code` | string | 是 | `planner` |
| `skill_tags` | string[] | 否 | `mrp,scheduling` |

关系：

- 多对一 `department`、`plant`。
- 可作为 `production_team.leader_employee_code`。
- 被任务、审批、审计引用。

状态：

- `active`
- `on_leave`
- `inactive`

Seed：

- `EMP-SZ-PLAN-001` 计划员。
- `EMP-SZ-QE-001` 质量工程师。
- `EMP-SZ-EAM-001` 设备工程师。
- `EMP-SZ-TL-001` 班组长。

## 8. 业务单据与事务对象

### 8.1 `sales_order` 销售订单

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `order_no` | string | 是 | `SO-202607-0001` |
| `customer_code` | ref | 是 | `CUST-SGNEV` |
| `legal_entity_code` | ref | 是 | `LE-HCTM-SH` |
| `material_code` | ref | 是 | `HCTM-BCP-A01` |
| `order_qty` | decimal | 是 | `12000` |
| `due_date` | date | 是 | `2026-07-20` |
| `priority` | enum | 是 | `high` |
| `original_order_no` | ref | 否 | `SO-202607-0001` |
| `confirmed_at` | datetime | 否 |  |

关系：

- 多对一 `customer`、`material`、`legal_entity`。
- 一对多 `production_order`、`shipment`。

状态：

- `draft`
- `confirmed`
- `planned`
- `partially_shipped`
- `shipped`
- `cancelled`

事件：

- `CustomerOrderReceived`
- `SalesOrderConfirmed`

### 8.2 `purchase_order` 采购订单

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `po_no` | string | 是 | `PO-202607-0001` |
| `supplier_code` | ref | 是 | `SUP-ALPHA-AL` |
| `legal_entity_code` | ref | 是 | `LE-HCTM-SH` |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `material_code` | ref | 是 | `AL-PLATE-6061-T6` |
| `order_qty` | decimal | 是 | `5000` |
| `promised_date` | date | 是 | `2026-07-08` |
| `latest_eta` | date | 否 | `2026-07-11` |

关系：

- 多对一 `supplier`、`material`、`plant`。
- 一对多 `goods_receipt`。

状态：

- `draft`
- `released`
- `in_transit`
- `partially_received`
- `received`
- `delayed`
- `closed`
- `cancelled`

事件：

- `PurchaseOrderReleased`
- `SupplierShipmentDispatched`
- `SupplierDeliveryDelayed`

### 8.3 `goods_receipt` 收货单

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `receipt_no` | string | 是 | `GR-202607-0001` |
| `po_no` | ref | 是 | `PO-202607-0001` |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `material_code` | ref | 是 | `AL-PLATE-6061-T6` |
| `supplier_lot_no` | string | 是 | `ALPHA-20260711-01` |
| `received_qty` | decimal | 是 | `5000` |
| `received_at` | datetime | 是 | `2026-07-11T10:00:00+08:00` |

关系：

- 多对一 `purchase_order`、`material`。
- 一对一或一对多 `inspection_order`。

状态：

- `received`
- `inspection_pending`
- `accepted`
- `rejected`
- `closed`

事件：

- `MaterialReceived`

### 8.4 `inspection_order` 检验单

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `inspection_no` | string | 是 | `IQC-202607-0001` |
| `inspection_type` | enum | 是 | `incoming` |
| `object_type` | enum | 是 | `goods_receipt` |
| `object_no` | string | 是 | `GR-202607-0001` |
| `material_code` | ref | 是 | `AL-PLATE-6061-T6` |
| `lot_no` | string | 是 | `ALPHA-20260711-01` |
| `sample_qty` | decimal | 是 | `80` |
| `accepted_qty` | decimal | 否 | `5000` |
| `rejected_qty` | decimal | 否 | `0` |
| `defect_code` | string | 否 |  |
| `inspection_level` | enum | 是 | `normal` |

关系：

- 引用 `goods_receipt` 或 `operation_task`。
- 可生成 `quality_issue`。

状态：

- `pending`
- `in_progress`
- `passed`
- `failed`
- `waived`
- `closed`

事件：

- `IncomingInspectionPassed`
- `IncomingInspectionFailed`
- `ProcessInspectionFailed`

### 8.5 `inventory_transaction` 库存事务

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `transaction_no` | string | 是 | `INV-202607-0001` |
| `transaction_type` | enum | 是 | `putaway` |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `warehouse_code` | ref | 是 | `WH-SZ-RM` |
| `location_code` | ref | 是 | `RM-A01-01` |
| `material_code` | ref | 是 | `AL-PLATE-6061-T6` |
| `lot_no` | string | 是 | `ALPHA-20260711-01` |
| `qty` | decimal | 是 | `5000` |
| `direction` | enum | 是 | `in` |
| `source_object_type` | enum | 是 | `goods_receipt` |
| `source_object_no` | string | 是 | `GR-202607-0001` |

关系：

- 多对一 `material`、`warehouse`、`storage_location`。
- 引用来源单据。

状态：

- `posted`
- `reversed`

事件：

- `InventoryPutawayCompleted`

### 8.6 `production_order` 生产订单

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `work_order_no` | string | 是 | `WO-202607-0001` |
| `sales_order_no` | ref | 否 | `SO-202607-0001` |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `material_code` | ref | 是 | `HCTM-BCP-A01` |
| `routing_code` | ref | 是 | `RT-HCTM-BCP-A01-V1` |
| `planned_qty` | decimal | 是 | `10800` |
| `released_qty` | decimal | 否 | `10800` |
| `due_date` | date | 是 | `2026-07-20` |
| `priority` | enum | 是 | `high` |

关系：

- 多对一 `sales_order`、`material`、`routing`、`plant`。
- 一对多 `operation_task`。

状态：

- `planned`
- `released`
- `in_progress`
- `completed`
- `closed`
- `cancelled`

事件：

- `ProductionOrderReleased`

### 8.7 `operation_task` 工序任务

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `task_no` | string | 是 | `OPT-202607-0001-030` |
| `work_order_no` | ref | 是 | `WO-202607-0001` |
| `operation_code` | ref | 是 | `OP-WLD` |
| `work_center_code` | ref | 是 | `WLD-02` |
| `equipment_code` | ref | 否 | `LAS-WLD-02` |
| `team_code` | ref | 否 | `TEAM-SZ-BCP-A-DAY` |
| `planned_qty` | decimal | 是 | `10800` |
| `completed_qty` | decimal | 否 | `0` |
| `scrap_qty` | decimal | 否 | `0` |
| `planned_start_at` | datetime | 否 |  |
| `planned_end_at` | datetime | 否 |  |

关系：

- 多对一 `production_order`、`operation`、`work_center`、`equipment`、`production_team`。
- 可触发 `inspection_order`。

状态：

- `pending`
- `ready`
- `running`
- `paused`
- `completed`
- `blocked`
- `cancelled`

事件：

- `OperationStarted`
- `OperationCompleted`
- `MachineDown` 影响该任务。

### 8.8 `shipment` 发货单

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `shipment_no` | string | 是 | `SHIP-202607-0001` |
| `sales_order_no` | ref | 是 | `SO-202607-0001` |
| `customer_code` | ref | 是 | `CUST-SGNEV` |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `material_code` | ref | 是 | `HCTM-BCP-A01` |
| `ship_qty` | decimal | 是 | `9000` |
| `ship_date` | date | 是 | `2026-07-19` |
| `carrier` | string | 否 | 星河指定承运商 |

关系：

- 多对一 `sales_order`、`customer`、`material`。
- 引用成品库存批次。

状态：

- `planned`
- `picked`
- `dispatched`
- `delivered`
- `cancelled`

事件：

- `ShipmentDispatched`

### 8.9 `quality_issue` 质量问题

关键字段：

| 字段 | 类型 | 必填 | 示例 |
| --- | --- | --- | --- |
| `issue_no` | string | 是 | `QI-202607-0001` |
| `issue_source` | enum | 是 | `incoming_inspection` |
| `severity` | enum | 是 | `major` |
| `plant_code` | ref | 是 | `PLT-SZ` |
| `material_code` | ref | 是 | `AL-PLATE-6061-T6` |
| `lot_no` | string | 否 | `BETA-20260712-01` |
| `related_object_type` | enum | 是 | `inspection_order` |
| `related_object_no` | string | 是 | `IQC-202607-0002` |
| `defect_code` | string | 是 | `SURFACE_SCRATCH` |
| `containment_action` | text | 否 | 批次隔离，等待质量工程师判定 |
| `owner_employee_code` | ref | 否 | `EMP-SZ-QE-001` |

关系：

- 可引用 `inspection_order`、`operation_task`、`supplier`、`equipment`。

状态：

- `open`
- `contained`
- `analysis`
- `resolved`
- `closed`

事件：

- 由 `IncomingInspectionFailed` 或 `ProcessInspectionFailed` 触发创建。

## 9. MVP 关系图

```text
enterprise_group
-> business_unit
-> plant
-> department
-> employee

plant
-> warehouse
-> storage_location

plant
-> work_center
-> equipment

material
-> bom
-> material

material
-> routing
-> operation
-> work_center

customer
-> sales_order
-> production_order
-> operation_task
-> shipment

supplier
-> purchase_order
-> goods_receipt
-> inspection_order
-> inventory_transaction

inspection_order
-> quality_issue
```

## 10. MVP Seed 最小集合

第一批 seed 不需要覆盖所有字段，但必须能支撑第一条演示故事。

| 类别 | 最小数量 | 必须包含 |
| --- | --- | --- |
| 组织 | 4 | `GRP-HCTM`、`BU-EAST`、`LE-HCTM-SH`、`PLT-SZ` |
| 部门 | 5 | 计划、生产、质量、仓储、设备 |
| 人员 | 5 | 计划员、质量工程师、设备工程师、班组长、仓库主管 |
| 客户 | 1 | `CUST-SGNEV` |
| 供应商 | 2 | `SUP-ALPHA-AL`、`SUP-BETA-AL` |
| 物料 | 6 | 成品、铝板、密封圈、接头、包装箱、标签 |
| BOM | 5 | 成品到关键组件用量 |
| 工艺 | 1 | `RT-HCTM-BCP-A01-V1` |
| 工序 | 7 | IQC、CNC、焊接、清洗、检漏、尺寸检测、总装包装 |
| 工作中心 | 8 | `IQC-01` 到 `ASM-01` |
| 设备 | 8 | 包含 `LAS-WLD-02` |
| 仓库/库位 | 7 | 原材料仓、成品仓、不合格品区及默认库位 |
| 初始库存 | 5 | 成品 1,200，铝板 8,000，密封圈 30,000，接头 18,000，包装箱 1,500 |
| 初始订单 | 1 | 星河追加后总需求 12,000 件 |

## 11. IAOS 落地顺序

建议后续按以下顺序转入 IAOS：

1. 先建组织、交易方、物料、仓库、工作中心、设备。
2. 再建 BOM、工艺路线和工序。
3. 再建销售订单、采购订单、收货单、检验单、库存事务、生产订单和工序任务。
4. 最后把质量问题、Agent 建议和经营分析挂到事件流之后。

第一批 metadata entity 可以先不追求完整字段，只要支持：

- 创建销售订单。
- 计算物料缺口。
- 生成采购订单和生产订单。
- 记录收货、检验和入库。
- 记录工序开始、完成、设备故障。
- 记录成品入库和发货。
- 让 Agent 能读取对象上下文并解释风险。

