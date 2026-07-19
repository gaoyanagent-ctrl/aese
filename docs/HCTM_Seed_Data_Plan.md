# HCTM Seed Data Plan

本文定义华辰热管理系统集团 AESE MVP 的 seed 数据计划。它不是最终导入脚本，而是后续 JSON、SQL、Go seed、IAOS metadata seed 和演示事件序列的来源清单。

## 1. Seed 目标

第一批 seed 必须支撑第一条演示故事：

> 客户追加订单下的交付承诺重算。

Seed 后系统应具备：

- 一个可识别的租户、集团、事业部、法人和工厂。
- 苏州制造基地的部门、班次、班组、人员、仓库、库位、工作中心和设备。
- 电池冷却板组件 `HCTM-BCP-A01` 的物料、BOM、工艺路线和工序。
- 星河新能源汽车客户、主供应商和备选供应商。
- 初始库存、初始订单、在途采购订单和演示异常数据。
- 一条可按时间顺序重放的事件序列。

## 2. Seed 分层

Seed 建议分三层执行：

| 层级 | 名称 | 内容 | 可重复执行要求 |
| --- | --- | --- | --- |
| L1 | 基础主数据 | 组织、交易方、物料、BOM、工艺、资源、人员 | 必须幂等 |
| L2 | 演示初始业务数据 | 初始库存、客户订单、采购订单、生产计划草稿 | 必须幂等，可按场景重置 |
| L3 | 演示事件序列 | 订单确认、MRP、供应商延期、设备故障、检验、入库、发货 | 可重放，需使用固定 `correlation_id` |

建议所有 seed 记录带：

- `tenant_id = tenant-hctm`
- `scenario = hctm_order_expedite_01`
- `seed_version = 2026-06-26-mvp`

## 3. 租户和组织

### 3.1 Tenant

| 字段 | 值 |
| --- | --- |
| `tenant_id` | `tenant-hctm` |
| `tenant_name` | 华辰热管理系统集团 |
| `timezone` | `Asia/Shanghai` |
| `currency` | `CNY` |
| `scenario` | `hctm_order_expedite_01` |

### 3.2 `enterprise_group`

| group_code | group_name | english_name | industry | headquarter_city | status |
| --- | --- | --- | --- | --- | --- |
| `GRP-HCTM` | 华辰热管理系统集团有限公司 | Huachen Thermal Management Systems Group Co., Ltd. | `nev_thermal_management` | 上海 | `active` |

### 3.3 `business_unit`

| bu_code | bu_name | group_code | region | status |
| --- | --- | --- | --- | --- |
| `BU-EAST` | 华东事业部 | `GRP-HCTM` | 华东 | `active` |

### 3.4 `legal_entity`

| legal_entity_code | legal_entity_name | group_code | tax_id | currency | status |
| --- | --- | --- | --- | --- | --- |
| `LE-HCTM-SH` | 华辰热管理系统上海有限公司 | `GRP-HCTM` | `91310000HCTM000001` | `CNY` | `active` |

### 3.5 `plant`

| plant_code | plant_name | bu_code | legal_entity_code | city | plant_type | timezone | status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `PLT-SZ` | 苏州制造基地 | `BU-EAST` | `LE-HCTM-SH` | 苏州 | `manufacturing` | `Asia/Shanghai` | `active` |

## 4. 部门、人员、班次和班组

### 4.1 `department`

| dept_code | dept_name | plant_code | dept_type | parent_dept_code | status |
| --- | --- | --- | --- | --- | --- |
| `DEPT-SZ-MGMT` | 工厂管理部 | `PLT-SZ` | `management` |  | `active` |
| `DEPT-SZ-PLAN` | 计划物流部 | `PLT-SZ` | `planning` | `DEPT-SZ-MGMT` | `active` |
| `DEPT-SZ-MFG` | 生产制造部 | `PLT-SZ` | `manufacturing` | `DEPT-SZ-MGMT` | `active` |
| `DEPT-SZ-QUALITY` | 质量管理部 | `PLT-SZ` | `quality` | `DEPT-SZ-MGMT` | `active` |
| `DEPT-SZ-WH` | 仓储物流部 | `PLT-SZ` | `warehouse` | `DEPT-SZ-MGMT` | `active` |
| `DEPT-SZ-EAM` | 设备工程部 | `PLT-SZ` | `equipment` | `DEPT-SZ-MGMT` | `active` |
| `DEPT-SZ-FIN` | 财务驻厂组 | `PLT-SZ` | `finance` | `DEPT-SZ-MGMT` | `active` |

### 4.2 `shift`

| shift_code | shift_name | start_time | end_time | work_hours | overtime_allowed | status |
| --- | --- | --- | --- | --- | --- | --- |
| `SHIFT-DAY` | 白班 | `08:00` | `16:00` | 8 | true | `active` |
| `SHIFT-NIGHT` | 夜班 | `16:00` | `00:00` | 8 | true | `active` |
| `SHIFT-OT` | 加班时段 | `00:00` | `02:00` | 2 | true | `active` |

### 4.3 `employee`

| employee_code | employee_name | dept_code | plant_code | role_code | skill_tags | status |
| --- | --- | --- | --- | --- | --- | --- |
| `EMP-SZ-GM-001` | 周厂长 | `DEPT-SZ-MGMT` | `PLT-SZ` | `plant_manager` | `operation,kpi` | `active` |
| `EMP-SH-SALES-001` | 陈销售 | `DEPT-SZ-PLAN` | `PLT-SZ` | `sales_manager` | `customer_order` | `active` |
| `EMP-SZ-PLAN-001` | 李计划 | `DEPT-SZ-PLAN` | `PLT-SZ` | `planner` | `mrp,scheduling` | `active` |
| `EMP-SZ-BUY-001` | 王采购 | `DEPT-SZ-PLAN` | `PLT-SZ` | `buyer` | `purchase,expedite` | `active` |
| `EMP-SZ-WH-001` | 赵仓管 | `DEPT-SZ-WH` | `PLT-SZ` | `warehouse_supervisor` | `receipt,putaway,shipment` | `active` |
| `EMP-SZ-IQC-001` | 钱检验 | `DEPT-SZ-QUALITY` | `PLT-SZ` | `incoming_inspector` | `iqc,aluminum` | `active` |
| `EMP-SZ-QE-001` | 孙质量 | `DEPT-SZ-QUALITY` | `PLT-SZ` | `quality_engineer` | `root_cause,containment` | `active` |
| `EMP-SZ-TL-001` | 吴班长 | `DEPT-SZ-MFG` | `PLT-SZ` | `team_leader` | `line_a,welding` | `active` |
| `EMP-SZ-EAM-001` | 郑设备 | `DEPT-SZ-EAM` | `PLT-SZ` | `equipment_engineer` | `laser_welder,maintenance` | `active` |
| `EMP-SZ-FIN-001` | 冯财务 | `DEPT-SZ-FIN` | `PLT-SZ` | `accountant` | `invoice,cost` | `active` |

### 4.4 `production_team`

| team_code | team_name | plant_code | dept_code | shift_code | line_code | leader_employee_code | status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `TEAM-SZ-BCP-A-DAY` | 电池冷却板 A 线白班 | `PLT-SZ` | `DEPT-SZ-MFG` | `SHIFT-DAY` | `SZ-BCP-LINE-A` | `EMP-SZ-TL-001` | `active` |
| `TEAM-SZ-BCP-A-NIGHT` | 电池冷却板 A 线夜班 | `PLT-SZ` | `DEPT-SZ-MFG` | `SHIFT-NIGHT` | `SZ-BCP-LINE-A` | `EMP-SZ-TL-001` | `active` |

## 5. 客户和供应商

### 5.1 `customer`

| customer_code | customer_name | customer_type | priority | payment_terms | delivery_mode | default_ship_to | status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `CUST-SGNEV` | 星河新能源汽车有限公司 | `oem` | `high` | `NET60` | `jit_split` | 星河苏州整车工厂 | `active` |
| `CUST-ORBITHM` | 欧瑞热管理系统有限公司 | `tier1` | `normal` | `NET45` | `weekly_delivery` | 欧瑞合肥工厂 | `active` |

### 5.2 `supplier`

| supplier_code | supplier_name | supplier_type | rating | default_lead_time_days | certified_materials | inspection_level | status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `SUP-ALPHA-AL` | 安铝金属材料有限公司 | `material` | `A` | 7 | `AL-PLATE-6061-T6` | `normal` | `qualified` |
| `SUP-BETA-AL` | 贝塔铝业有限公司 | `material` | `B` | 5 | `AL-PLATE-6061-T6` | `tightened` | `qualified` |
| `SUP-SEAL-01` | 恒密密封科技有限公司 | `material` | `A` | 5 | `SEAL-RING-A01` | `normal` | `qualified` |
| `SUP-FIT-01` | 鼎联接头制造有限公司 | `material` | `B` | 6 | `FITTING-A01` | `normal` | `qualified` |
| `SUP-PKG-01` | 诚达包装材料有限公司 | `packaging` | `A` | 3 | `PKG-BOX-A01,LBL-TRACE-A01` | `skip_lot` | `qualified` |

## 6. 物料、BOM、工艺和工序

### 6.1 `material`

| material_code | material_name | material_type | uom | lot_control | serial_control | customer_material_no | default_plant_code | safety_stock_qty | status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `HCTM-BCP-A01` | 电池冷却板组件 A01 | `finished_good` | `pcs` | true | false | `SGNEV-CP-7788` | `PLT-SZ` | 1200 | `active` |
| `AL-PLATE-6061-T6` | 6061-T6 铝板 | `raw_material` | `sheet` | true | false |  | `PLT-SZ` | 6000 | `active` |
| `SEAL-RING-A01` | 密封圈 A01 | `purchased_part` | `pcs` | true | false |  | `PLT-SZ` | 20000 | `active` |
| `FITTING-A01` | 接头组件 A01 | `purchased_part` | `set` | true | false |  | `PLT-SZ` | 12000 | `active` |
| `PKG-BOX-A01` | 包装箱 A01 | `packaging` | `box` | false | false |  | `PLT-SZ` | 800 | `active` |
| `LBL-TRACE-A01` | 追溯标签 A01 | `packaging` | `pcs` | false | false |  | `PLT-SZ` | 10000 | `active` |

### 6.2 `bom`

| bom_code | parent_material_code | component_material_code | qty_per | uom | scrap_rate | effective_from | status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `BOM-HCTM-BCP-A01-V1` | `HCTM-BCP-A01` | `AL-PLATE-6061-T6` | 1.05 | `sheet` | 0.02 | `2026-07-01` | `active` |
| `BOM-HCTM-BCP-A01-V1` | `HCTM-BCP-A01` | `SEAL-RING-A01` | 2 | `pcs` | 0.005 | `2026-07-01` | `active` |
| `BOM-HCTM-BCP-A01-V1` | `HCTM-BCP-A01` | `FITTING-A01` | 2 | `set` | 0.005 | `2026-07-01` | `active` |
| `BOM-HCTM-BCP-A01-V1` | `HCTM-BCP-A01` | `PKG-BOX-A01` | 0.1 | `box` | 0 | `2026-07-01` | `active` |
| `BOM-HCTM-BCP-A01-V1` | `HCTM-BCP-A01` | `LBL-TRACE-A01` | 1 | `pcs` | 0 | `2026-07-01` | `active` |

### 6.3 `routing`

| routing_code | material_code | plant_code | version | effective_from | standard_daily_capacity | status |
| --- | --- | --- | --- | --- | --- | --- |
| `RT-HCTM-BCP-A01-V1` | `HCTM-BCP-A01` | `PLT-SZ` | `V1` | `2026-07-01` | 800 | `active` |

### 6.4 `operation`

| operation_code | routing_code | sequence_no | operation_name | work_center_code | standard_time_min | inspection_required | yield_rate | status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `OP-IQC` | `RT-HCTM-BCP-A01-V1` | 10 | 来料检验 | `IQC-01` | 0.3 | true | 1.0 | `active` |
| `OP-CNC` | `RT-HCTM-BCP-A01-V1` | 20 | CNC 精加工 | `CNC-01` | 1.0 | false | 0.99 | `active` |
| `OP-WLD` | `RT-HCTM-BCP-A01-V1` | 30 | 焊接 | `WLD-02` | 1.2 | false | 0.98 | `active` |
| `OP-CLN` | `RT-HCTM-BCP-A01-V1` | 40 | 清洗 | `CLN-01` | 0.5 | false | 0.995 | `active` |
| `OP-LKT` | `RT-HCTM-BCP-A01-V1` | 50 | 气密检漏 | `LKT-01` | 0.7 | true | 0.99 | `active` |
| `OP-DIM` | `RT-HCTM-BCP-A01-V1` | 60 | 尺寸检测 | `DIM-01` | 0.4 | true | 0.995 | `active` |
| `OP-ASM` | `RT-HCTM-BCP-A01-V1` | 70 | 总装包装 | `ASM-01` | 0.8 | false | 0.995 | `active` |

## 7. 工作中心、设备、仓库和库位

### 7.1 `work_center`

| work_center_code | work_center_name | plant_code | center_type | line_code | capacity_per_shift | bottleneck_flag | status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `IQC-01` | 来料检验工作中心 | `PLT-SZ` | `quality` | `SZ-BCP-LINE-A` | 6000 | false | `active` |
| `CNC-01` | CNC 加工工作中心 | `PLT-SZ` | `production` | `SZ-BCP-LINE-A` | 420 | false | `active` |
| `WLD-01` | 焊接工作中心 1 | `PLT-SZ` | `production` | `SZ-BCP-LINE-A` | 400 | true | `active` |
| `WLD-02` | 焊接工作中心 2 | `PLT-SZ` | `production` | `SZ-BCP-LINE-A` | 400 | true | `active` |
| `CLN-01` | 清洗工作中心 | `PLT-SZ` | `production` | `SZ-BCP-LINE-A` | 900 | false | `active` |
| `LKT-01` | 气密检漏工作中心 | `PLT-SZ` | `quality` | `SZ-BCP-LINE-A` | 850 | false | `active` |
| `DIM-01` | 尺寸检测工作中心 | `PLT-SZ` | `quality` | `SZ-BCP-LINE-A` | 1000 | false | `active` |
| `ASM-01` | 总装包装工作中心 | `PLT-SZ` | `production` | `SZ-BCP-LINE-A` | 950 | false | `active` |

### 7.2 `equipment`

| equipment_code | equipment_name | plant_code | work_center_code | equipment_type | criticality | commissioned_date | current_status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `CNC-A01` | CNC 加工中心 A01 | `PLT-SZ` | `CNC-01` | `cnc` | `medium` | `2023-04-01` | `running` |
| `CNC-A02` | CNC 加工中心 A02 | `PLT-SZ` | `CNC-01` | `cnc` | `medium` | `2023-05-01` | `running` |
| `LAS-WLD-01` | 激光焊接机 1 | `PLT-SZ` | `WLD-01` | `laser_welder` | `high` | `2024-03-01` | `running` |
| `LAS-WLD-02` | 激光焊接机 2 | `PLT-SZ` | `WLD-02` | `laser_welder` | `high` | `2024-05-01` | `running` |
| `USC-CLN-01` | 超声波清洗机 | `PLT-SZ` | `CLN-01` | `cleaning` | `medium` | `2024-01-10` | `running` |
| `AIR-LKT-01` | 气密检漏台 1 | `PLT-SZ` | `LKT-01` | `leak_tester` | `high` | `2024-02-01` | `running` |
| `AIR-LKT-02` | 气密检漏台 2 | `PLT-SZ` | `LKT-01` | `leak_tester` | `high` | `2024-02-15` | `running` |
| `LAS-MRK-01` | 激光打标机 | `PLT-SZ` | `ASM-01` | `laser_marker` | `medium` | `2024-04-01` | `running` |

### 7.3 `tooling`

| tooling_code | tooling_name | plant_code | work_center_code | lifetime_cycles | used_cycles | current_status |
| --- | --- | --- | --- | --- | --- | --- |
| `FIX-BCP-A01-WLD` | 电池冷却板焊接夹具 | `PLT-SZ` | `WLD-02` | 50000 | 12000 | `available` |

### 7.4 `warehouse`

| warehouse_code | warehouse_name | plant_code | warehouse_type | managed_by_dept_code | status |
| --- | --- | --- | --- | --- | --- |
| `WH-SZ-RM` | 苏州原材料仓 | `PLT-SZ` | `raw_material` | `DEPT-SZ-WH` | `active` |
| `WH-SZ-WIP` | 苏州半成品暂存区 | `PLT-SZ` | `wip` | `DEPT-SZ-WH` | `active` |
| `WH-SZ-FG` | 苏州成品仓 | `PLT-SZ` | `finished_goods` | `DEPT-SZ-WH` | `active` |
| `WH-SZ-NG` | 苏州不合格品区 | `PLT-SZ` | `nonconforming` | `DEPT-SZ-QUALITY` | `active` |

### 7.5 `storage_location`

| location_code | warehouse_code | location_name | location_type | material_restriction | status |
| --- | --- | --- | --- | --- | --- |
| `RM-A01-01` | `WH-SZ-RM` | 原材料 A 区 01-01 | `normal` | `AL-PLATE-6061-T6` | `available` |
| `RM-B01-01` | `WH-SZ-RM` | 原材料 B 区 01-01 | `normal` | `SEAL-RING-A01` | `available` |
| `RM-C01-01` | `WH-SZ-RM` | 原材料 C 区 01-01 | `normal` | `FITTING-A01` | `available` |
| `PKG-A01-01` | `WH-SZ-RM` | 包材 A 区 01-01 | `normal` | `PKG-BOX-A01` | `available` |
| `WIP-A01-01` | `WH-SZ-WIP` | 半成品 A 区 01-01 | `normal` |  | `available` |
| `FG-A01-01` | `WH-SZ-FG` | 成品 A 区 01-01 | `normal` | `HCTM-BCP-A01` | `available` |
| `NG-QA-01` | `WH-SZ-NG` | 不合格品隔离区 01 | `blocked` |  | `available` |

## 8. 初始库存

建议用 `inventory_transaction` 的 opening balance 类型导入初始库存。

| transaction_no | transaction_type | warehouse_code | location_code | material_code | lot_no | qty | direction | source_object_type | source_object_no | status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `INV-OPEN-0001` | `opening_balance` | `WH-SZ-FG` | `FG-A01-01` | `HCTM-BCP-A01` | `FG-OPEN-20260701` | 1200 | `in` | `seed` | `SEED-HCTM-MVP` | `posted` |
| `INV-OPEN-0002` | `opening_balance` | `WH-SZ-RM` | `RM-A01-01` | `AL-PLATE-6061-T6` | `AL-OPEN-20260701` | 8000 | `in` | `seed` | `SEED-HCTM-MVP` | `posted` |
| `INV-OPEN-0003` | `opening_balance` | `WH-SZ-RM` | `RM-B01-01` | `SEAL-RING-A01` | `SEAL-OPEN-20260701` | 30000 | `in` | `seed` | `SEED-HCTM-MVP` | `posted` |
| `INV-OPEN-0004` | `opening_balance` | `WH-SZ-RM` | `RM-C01-01` | `FITTING-A01` | `FIT-OPEN-20260701` | 18000 | `in` | `seed` | `SEED-HCTM-MVP` | `posted` |
| `INV-OPEN-0005` | `opening_balance` | `WH-SZ-RM` | `PKG-A01-01` | `PKG-BOX-A01` | `PKG-OPEN-20260701` | 1500 | `in` | `seed` | `SEED-HCTM-MVP` | `posted` |

## 9. 演示初始业务数据

### 9.1 `sales_order`

| order_no | customer_code | legal_entity_code | material_code | order_qty | due_date | priority | original_order_no | status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `SO-202607-0001` | `CUST-SGNEV` | `LE-HCTM-SH` | `HCTM-BCP-A01` | 10000 | `2026-07-20` | `high` |  | `confirmed` |
| `SO-202607-0001-ADD1` | `CUST-SGNEV` | `LE-HCTM-SH` | `HCTM-BCP-A01` | 2000 | `2026-07-20` | `high` | `SO-202607-0001` | `confirmed` |

MVP 运行时可以把两张订单合并为总需求 12,000 件，也可以保留追加订单用于展示需求变化。

### 9.2 `purchase_order`

| po_no | supplier_code | legal_entity_code | plant_code | material_code | order_qty | promised_date | latest_eta | status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `PO-202607-0001` | `SUP-ALPHA-AL` | `LE-HCTM-SH` | `PLT-SZ` | `AL-PLATE-6061-T6` | 5000 | `2026-07-08` | `2026-07-08` | `in_transit` |
| `PO-202607-0002` | `SUP-BETA-AL` | `LE-HCTM-SH` | `PLT-SZ` | `AL-PLATE-6061-T6` | 2000 | `2026-07-12` | `2026-07-12` | `draft` |

说明：

- `PO-202607-0001` 是主供应商在途订单，事件序列中会延期到 `2026-07-11`。
- `PO-202607-0002` 是备选供应商草稿，计划 Agent 可建议释放。

### 9.3 `production_order`

| work_order_no | sales_order_no | plant_code | material_code | routing_code | planned_qty | released_qty | due_date | priority | status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `WO-202607-0001` | `SO-202607-0001` | `PLT-SZ` | `HCTM-BCP-A01` | `RT-HCTM-BCP-A01-V1` | 10800 | 0 | `2026-07-20` | `high` | `planned` |

说明：

- 计划生产数量为 10,800 件，因为已有成品库存 1,200 件。
- 释放动作由事件 `proc.production_order.released` 表达。

### 9.4 `operation_task`

| task_no | work_order_no | operation_code | work_center_code | equipment_code | team_code | planned_qty | status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `OPT-202607-0001-020` | `WO-202607-0001` | `OP-CNC` | `CNC-01` | `CNC-A01` | `TEAM-SZ-BCP-A-DAY` | 10800 | `pending` |
| `OPT-202607-0001-030` | `WO-202607-0001` | `OP-WLD` | `WLD-02` | `LAS-WLD-02` | `TEAM-SZ-BCP-A-DAY` | 10800 | `pending` |
| `OPT-202607-0001-050` | `WO-202607-0001` | `OP-LKT` | `LKT-01` | `AIR-LKT-01` | `TEAM-SZ-BCP-A-DAY` | 10800 | `pending` |
| `OPT-202607-0001-070` | `WO-202607-0001` | `OP-ASM` | `ASM-01` | `LAS-MRK-01` | `TEAM-SZ-BCP-A-DAY` | 10800 | `pending` |

## 10. 演示事件序列 Seed

统一字段：

| 字段 | 值 |
| --- | --- |
| `correlation_id` | `corr-so-202607-0001` |
| `scenario` | `hctm_order_expedite_01` |
| `source` | `aese.hctm.simulator` |
| `tenant_id` | `tenant-hctm` |
| `plant_code` | `PLT-SZ` |

事件序列：

| seq | event_id | timestamp | type | business_object_type | business_object_id | 摘要 |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | `evt-20260701-0001` | `2026-07-01T09:00:00+08:00` | `o2d.order.received` | `sales_order` | `SO-202607-0001` | 收到原始 10,000 件订单 |
| 2 | `evt-20260701-0002` | `2026-07-01T09:30:00+08:00` | `o2d.order.received` | `sales_order` | `SO-202607-0001-ADD1` | 收到追加 2,000 件订单 |
| 3 | `evt-20260701-0003` | `2026-07-01T10:10:00+08:00` | `o2d.order.confirmed` | `sales_order` | `SO-202607-0001` | 总需求 12,000 件确认 |
| 4 | `evt-20260701-0004` | `2026-07-01T10:20:00+08:00` | `o2d.mrp.generated` | `sales_order` | `SO-202607-0001` | 初始 MRP 生成 |
| 5 | `evt-20260701-0005` | `2026-07-01T10:25:00+08:00` | `o2d.purchase_requirement.created` | `sales_order` | `SO-202607-0001` | 创建铝板时点风险采购需求 |
| 6 | `evt-20260701-0006` | `2026-07-01T11:00:00+08:00` | `o2d.purchase_order.released` | `purchase_order` | `PO-202607-0001` | 主供应商铝板 PO 发布 |
| 7 | `evt-20260702-0001` | `2026-07-02T15:00:00+08:00` | `o2d.supplier_shipment.dispatched` | `purchase_order` | `PO-202607-0001` | 主供应商发货，原 ETA 7 月 8 日 |
| 8 | `evt-20260708-0001` | `2026-07-08T09:00:00+08:00` | `o2d.supplier_delivery.delayed` | `purchase_order` | `PO-202607-0001` | 主供应商延期 3 天 |
| 9 | `evt-20260708-0002` | `2026-07-08T09:30:00+08:00` | `eam.machine.down` | `equipment` | `LAS-WLD-02` | 焊接设备停机 8 小时 |
| 10 | `evt-20260708-0003` | `2026-07-08T09:40:00+08:00` | `o2d.mrp.generated` | `sales_order` | `SO-202607-0001` | 异常触发重算，建议启用备选供应商和加班 |
| 11 | `evt-20260708-0004` | `2026-07-08T10:20:00+08:00` | `o2d.purchase_order.released` | `purchase_order` | `PO-202607-0002` | 释放备选供应商加急 PO |
| 12 | `evt-20260711-0001` | `2026-07-11T10:00:00+08:00` | `whs.material.received` | `goods_receipt` | `GR-202607-0001` | 主供应商铝板到货 5,000 张 |
| 13 | `evt-20260711-0002` | `2026-07-11T14:00:00+08:00` | `qms.incoming_inspection.passed` | `inspection_order` | `IQC-202607-0001` | 主供应商铝板检验合格 |
| 14 | `evt-20260711-0003` | `2026-07-11T15:00:00+08:00` | `whs.inventory.putaway_completed` | `inventory_transaction` | `INV-202607-0001` | 铝板 5,000 张入库 |
| 15 | `evt-20260712-0001` | `2026-07-12T11:00:00+08:00` | `whs.material.received` | `goods_receipt` | `GR-202607-0002` | 备选供应商铝板到货 2,000 张 |
| 16 | `evt-20260712-0002` | `2026-07-12T15:00:00+08:00` | `qms.incoming_inspection.failed` | `inspection_order` | `IQC-202607-0002` | 备选供应商批次发现 300 张表面划伤 |
| 17 | `evt-20260713-0001` | `2026-07-13T08:00:00+08:00` | `proc.production_order.released` | `production_order` | `WO-202607-0001` | 生产订单释放 |
| 18 | `evt-20260714-0001` | `2026-07-14T08:00:00+08:00` | `proc.operation.started` | `operation_task` | `OPT-202607-0001-030` | 焊接工序开始 |
| 19 | `evt-20260717-0001` | `2026-07-17T20:00:00+08:00` | `proc.operation.completed` | `operation_task` | `OPT-202607-0001-030` | 焊接工序完成，报废 80 件 |
| 20 | `evt-20260718-0001` | `2026-07-18T10:00:00+08:00` | `whs.finished_goods.received` | `inventory_transaction` | `INV-FG-202607-0001` | 成品入库 10,500 件 |
| 21 | `evt-20260719-0001` | `2026-07-19T16:00:00+08:00` | `o2d.shipment.dispatched` | `shipment` | `SHIP-202607-0001` | 第一批发货 9,000 件 |
| 22 | `evt-20260720-0001` | `2026-07-20T11:00:00+08:00` | `o2d.shipment.dispatched` | `shipment` | `SHIP-202607-0002` | 第二批请求 3,000 件，实际发货 2,700 件，短缺 300 件 |

## 11. 关键事件 Payload Seed

### 11.1 初始 MRP `evt-20260701-0004`

```json
{
  "mrp_run_id": "MRP-202607-0001",
  "run_reason": "initial_plan",
  "order_no": "SO-202607-0001",
  "material_code": "HCTM-BCP-A01",
  "demand_qty": 12000,
  "available_finished_goods_qty": 1200,
  "net_production_qty": 10800,
  "material_shortages": [
    {
      "material_code": "AL-PLATE-6061-T6",
      "required_qty": 12600,
      "available_qty": 8000,
      "in_transit_qty": 5000,
      "shortage_qty": 0,
      "timing_risk": "none"
    }
  ],
  "capacity_risks": [],
  "recommended_actions": [
    "release_purchase_order_PO-202607-0001",
    "prepare_production_order_WO-202607-0001"
  ]
}
```

### 11.2 异常重算 MRP `evt-20260708-0003`

```json
{
  "mrp_run_id": "MRP-202607-0002",
  "run_reason": "exception_replan",
  "order_no": "SO-202607-0001",
  "material_code": "HCTM-BCP-A01",
  "demand_qty": 12000,
  "available_finished_goods_qty": 1200,
  "net_production_qty": 10800,
  "material_shortages": [
    {
      "material_code": "AL-PLATE-6061-T6",
      "required_qty": 12600,
      "available_qty": 8000,
      "in_transit_qty": 5000,
      "shortage_qty": 0,
      "timing_risk": "supplier_eta_delayed"
    }
  ],
  "capacity_risks": [
    {
      "work_center_code": "WLD-02",
      "equipment_code": "LAS-WLD-02",
      "risk_type": "machine_down",
      "impact_qty": 280,
      "risk_level": "high"
    }
  ],
  "recommended_actions": [
    "release_backup_supplier_po_PO-202607-0002",
    "enable_overtime_SHIFT-OT",
    "split_shipment_9000_plus_3000",
    "tighten_incoming_inspection_for_SUP-BETA-AL"
  ]
}
```

### 11.3 供应商延期 `evt-20260708-0001`

```json
{
  "po_no": "PO-202607-0001",
  "supplier_code": "SUP-ALPHA-AL",
  "material_code": "AL-PLATE-6061-T6",
  "original_eta": "2026-07-08",
  "new_eta": "2026-07-11",
  "delay_days": 3,
  "reason": "供应商熔炼排产延迟",
  "affected_order_no": "SO-202607-0001"
}
```

### 11.4 设备故障 `evt-20260708-0002`

```json
{
  "equipment_code": "LAS-WLD-02",
  "work_center_code": "WLD-02",
  "plant_code": "PLT-SZ",
  "down_at": "2026-07-08T09:30:00+08:00",
  "expected_repair_at": "2026-07-08T17:30:00+08:00",
  "expected_down_hours": 8,
  "capacity_loss_ratio": 0.35,
  "reason_code": "LASER_POWER_ALARM",
  "affected_task_no": "OPT-202607-0001-030"
}
```

### 11.5 来料检验失败 `evt-20260712-0002`

```json
{
  "inspection_no": "IQC-202607-0002",
  "receipt_no": "GR-202607-0002",
  "supplier_code": "SUP-BETA-AL",
  "material_code": "AL-PLATE-6061-T6",
  "lot_no": "BETA-20260712-01",
  "rejected_qty": 300,
  "defect_code": "SURFACE_SCRATCH",
  "severity": "major"
}
```

## 12. Agent 期望 Seed 输出

这些不是数据库主数据，而是演示验收时可比对的期望输出。

### 12.1 计划 Agent

输入事件：

- `o2d.order.confirmed`
- `o2d.mrp.generated`
- `o2d.supplier_delivery.delayed`
- `eam.machine.down`
- 第二次 `o2d.mrp.generated`

期望输出摘要：

```text
订单总需求 12,000 件，扣除成品库存 1,200 件后需生产 10,800 件。
铝板总需求约 12,600 张，当前 8,000 张加主供应商在途 5,000 张数量上可覆盖，但主供应商延期导致时点风险。
LAS-WLD-02 停机 8 小时导致焊接产能短期下降，建议释放备选供应商 PO、启用加班，并采用分批发运。当前可发库存只能支持 9,000 + 2,700 件，最终 300 件需补产或重新承诺。
```

### 12.2 质量 Agent

输入事件：

- `whs.material.received`
- `qms.incoming_inspection.failed`

期望输出摘要：

```text
备选供应商 SUP-BETA-AL 首批铝板触发加严检验。IQC-202607-0002 发现 300 张表面划伤，建议隔离 BETA-20260712-01 批次，禁止不合格数量投产，并评估剩余合格数量是否满足紧急生产。
```

### 12.3 经营分析 Agent

输入事件：

- `o2d.supplier_delivery.delayed`
- `eam.machine.down`
- `whs.finished_goods.received`
- `o2d.shipment.dispatched`

期望输出摘要：

```text
通过启用备选供应商、加班和分批发运，交付风险下降但尚未关闭：累计发运 11,700 件，仍有 300 件缺口。成本上升主要来自备选供应商采购、加严检验和加班；质量风险集中在备选供应商批次，需要质量隔离、补产或交期重承诺决策支撑。
```

## 13. 后续脚本化建议

后续可以把本文拆成：

```text
seed/hctm/master_data.json
seed/hctm/opening_balances.json
seed/hctm/demo_story_01_objects.json
seed/hctm/demo_story_01_events.json
```

导入顺序：

1. 组织和人员。
2. 客户和供应商。
3. 物料、BOM、工艺、工序。
4. 工作中心、设备、仓库、库位。
5. 初始库存。
6. 初始订单和采购订单。
7. 生产订单和工序任务。
8. 演示事件序列。

幂等策略：

- 主数据按 `tenant_id + *_code` upsert。
- 单据按 `tenant_id + *_no` upsert。
- 事件按 `event_id` 或 `metadata.idempotency_key` 去重。
- 演示重置时只清理 `scenario = hctm_order_expedite_01` 的 L2 和 L3 数据，不删除 L1 基础主数据。
