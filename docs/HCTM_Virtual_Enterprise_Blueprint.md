# HCTM Virtual Enterprise Blueprint

本文定义 AESE 第一阶段虚拟企业：**华辰热管理系统集团有限公司**。目标是给后续 IAOS metadata、seed 数据、Scenario Package、事件流、Agent 能力和 2D 企业沙盘提供可落地的业务蓝图。

## 1. 蓝图边界

本文不是完整企业咨询报告，而是 MVP 建模蓝图。它只覆盖第一条可运行主线：

```text
客户订单
-> 订单确认
-> MRP
-> 采购
-> 来料检验
-> 生产执行
-> 质量检验
-> 完工入库
-> 发货
-> 开票
-> 经营分析
```

MVP 锚点：

- 虚拟集团：华辰热管理系统集团有限公司。
- 主工厂：苏州制造基地。
- 主产线：电池冷却板 A 线。
- 主产品：电池冷却板组件 `HCTM-BCP-A01`。
- 主演示故事：客户追加订单下的交付承诺重算。

## 2. 集团基本设定

中文名称：

- 华辰热管理系统集团有限公司。

英文名称：

- Huachen Thermal Management Systems Group Co., Ltd.

英文简称：

- HCTM.

行业定位：

- 新能源汽车热管理系统核心零部件供应商。

核心产品族：

- 电池冷却板。
- 冷却管路。
- 热管理阀体。
- 铝合金结构件。

客户类型：

- 新能源主机厂。
- 一级热管理系统集成商。
- 海外平台项目客户。

业务模式：

- 项目定点后批量供货。
- 预测订单、正式订单和临时追加订单并存。
- 部分客户要求 JIT 或分批发运。
- 关键原材料依赖认证供应商。

经营目标：

- 准时交付率不低于 96%。
- 来料批次合格率不低于 98.5%。
- 过程一次合格率不低于 97%。
- 关键设备稼动率不低于 85%。
- 成品库存周转天数控制在 12 天以内。
- 异常发生后 2 小时内形成可执行处置方案。

## 3. 组织架构

集团结构：

```text
华辰热管理系统集团
├── 集团总部
│   ├── 战略与经营管理中心
│   ├── 财务共享中心
│   ├── IT 与数字化中心
│   ├── 集团采购中心
│   ├── 质量管理中心
│   └── 运营管理中心
├── 华东事业部
│   ├── 上海销售公司
│   ├── 苏州制造基地
│   └── 宁波制造基地
├── 华南事业部
│   ├── 广州销售公司
│   └── 佛山制造基地
└── 海外事业部
    ├── 泰国制造基地
    └── 欧洲销售办事处
```

MVP 只建模：

- 集团总部。
- 华东事业部。
- 上海销售公司。
- 苏州制造基地。

苏州制造基地部门：

- 工厂管理部。
- 计划物流部。
- 采购执行组。
- 生产制造部。
- 质量管理部。
- 工艺工程部。
- 设备工程部。
- 仓储物流部。
- 财务驻厂组。

## 4. 苏州制造基地

定位：

- 华东区域主力生产基地。
- MVP 中唯一实际运行的制造工厂。
- 负责电池冷却板组件 `HCTM-BCP-A01` 的批量生产和发运。

厂区结构：

```text
苏州制造基地
├── 办公区
├── 生产区
│   ├── 机加工车间
│   ├── 焊接车间
│   ├── 清洗测试车间
│   └── 总装包装车间
├── 仓储区
│   ├── 原材料仓
│   ├── 半成品暂存区
│   ├── 成品仓
│   └── 不合格品区
├── 质量区
│   ├── 来料检验区
│   ├── 过程检验区
│   ├── 气密测试区
│   └── 计量检测室
└── 公辅区
    ├── 空压站
    ├── 冷却水系统
    └── 配电房
```

MVP 生产节拍假设：

- 单班 8 小时。
- 常规双班生产。
- A 线常规日产能 800 件。
- 紧急双班加加班日产能 1,100 件。
- 焊接工作中心是关键瓶颈。
- 检漏测试能力在正常情况下略高于焊接能力。

## 5. 电池冷却板 A 线

产线名称：

- 苏州制造基地 - 电池冷却板 A 线。

产线编码：

- `SZ-BCP-LINE-A`.

主产品：

- 电池冷却板组件。

产品编码：

- `HCTM-BCP-A01`.

工艺流程：

```text
铝板来料
-> 来料检验
-> CNC 精加工
-> 焊接
-> 清洗
-> 气密检漏
-> 尺寸检测
-> 附件总装
-> 激光打标
-> 包装
-> 成品入库
```

工作中心：

| 编码 | 名称 | 类型 | MVP 作用 |
| --- | --- | --- | --- |
| `IQC-01` | 来料检验工作中心 | 质量 | 判断关键铝板批次是否可入库 |
| `CNC-01` | CNC 加工工作中心 | 生产 | 加工冷却板流道和安装面 |
| `WLD-01` | 焊接工作中心 1 | 生产 | 常规焊接能力 |
| `WLD-02` | 焊接工作中心 2 | 生产 | MVP 异常中的故障设备 |
| `CLN-01` | 清洗工作中心 | 生产 | 清洗油污和残留 |
| `LKT-01` | 气密检漏工作中心 | 质量 | 核心质量检测 |
| `DIM-01` | 尺寸检测工作中心 | 质量 | 尺寸抽检 |
| `ASM-01` | 总装包装工作中心 | 生产 | 装附件、打标和包装 |

关键设备：

| 编码 | 名称 | 所属工作中心 | 风险点 |
| --- | --- | --- | --- |
| `CNC-A01` | CNC 加工中心 A01 | `CNC-01` | 刀具磨损导致尺寸偏差 |
| `CNC-A02` | CNC 加工中心 A02 | `CNC-01` | 排队造成 WIP 积压 |
| `LAS-WLD-01` | 激光焊接机 1 | `WLD-01` | 焊缝强度和热变形 |
| `LAS-WLD-02` | 激光焊接机 2 | `WLD-02` | MVP 故障设备 |
| `USC-CLN-01` | 超声波清洗机 | `CLN-01` | 清洁度不达标 |
| `AIR-LKT-01` | 气密检漏台 1 | `LKT-01` | 漏率超标判定 |
| `AIR-LKT-02` | 气密检漏台 2 | `LKT-01` | 检测节拍波动 |
| `LAS-MRK-01` | 激光打标机 | `ASM-01` | 条码追溯失败 |

## 6. 产品和 BOM

主产品：

| 字段 | 值 |
| --- | --- |
| 产品编码 | `HCTM-BCP-A01` |
| 产品名称 | 电池冷却板组件 A01 |
| 产品类型 | 成品 |
| 计量单位 | 件 |
| 批次管理 | 是 |
| 序列号管理 | MVP 可不启用，后续扩展 |
| 客户料号 | `SGNEV-CP-7788` |
| 默认工厂 | 苏州制造基地 |

MVP BOM：

| 物料编码 | 名称 | 类型 | 单台用量 | 单位 | 关键约束 |
| --- | --- | --- | --- | --- | --- |
| `AL-PLATE-6061-T6` | 6061-T6 铝板 | 原材料 | 1.05 | 张 | 关键物料，认证供应商 |
| `BCP-UPPER-PLATE-A01` | 冷却板上板 | 半成品 | 1 | 件 | CNC 后形成 |
| `BCP-LOWER-PLATE-A01` | 冷却板下板 | 半成品 | 1 | 件 | CNC 后形成 |
| `SEAL-RING-A01` | 密封圈 | 外购件 | 2 | 个 | 批次追溯 |
| `FITTING-A01` | 接头组件 | 外购件 | 2 | 套 | 供应商质量风险 |
| `PKG-BOX-A01` | 包装箱 | 包材 | 0.1 | 个 | 10 件一箱 |
| `LBL-TRACE-A01` | 追溯标签 | 包材 | 1 | 张 | 打标和发运追溯 |

MVP 简化说明：

- 可先把上板和下板作为工序产出，不必建立完整多层半成品库存。
- 关键计算聚焦 `AL-PLATE-6061-T6`、焊接产能和成品库存。

## 7. 客户和供应商

MVP 客户：

| 编码 | 名称 | 类型 | 业务规则 |
| --- | --- | --- | --- |
| `CUST-SGNEV` | 星河新能源汽车有限公司 | 主机厂 | 交付优先级高，要求准时分批发运 |
| `CUST-ORBITHM` | 欧瑞热管理系统有限公司 | 一级供应商 | 预测稳定，价格压力高 |

MVP 供应商：

| 编码 | 名称 | 供应物料 | 规则和风险 |
| --- | --- | --- | --- |
| `SUP-ALPHA-AL` | 安铝金属材料有限公司 | 6061-T6 铝板 | 主供应商，MVP 中发生延期 |
| `SUP-BETA-AL` | 贝塔铝业有限公司 | 6061-T6 铝板 | 备选供应商，需加严检验 |
| `SUP-SEAL-01` | 恒密密封科技有限公司 | 密封圈 | 批次质量稳定 |
| `SUP-FIT-01` | 鼎联接头制造有限公司 | 接头组件 | 偶发尺寸偏差 |
| `SUP-PKG-01` | 诚达包装材料有限公司 | 包装箱和标签 | 低风险 |

## 8. 关键岗位

| 角色 | 所属部门 | MVP 职责 |
| --- | --- | --- |
| 工厂厂长 | 工厂管理部 | 查看交付、成本和异常处置方案 |
| 销售经理 | 上海销售公司 | 接收客户订单和追加需求 |
| 计划员 | 计划物流部 | 执行 MRP、评估交期、释放生产订单 |
| 采购员 | 采购执行组 | 释放采购订单、跟催供应商 |
| 仓库主管 | 仓储物流部 | 收货、入库、发料、成品入库 |
| 来料检验员 | 质量管理部 | 执行 IQC，判定批次合格或不合格 |
| 质量工程师 | 质量管理部 | 分析不良、追溯批次、提出处置 |
| 班组长 | 生产制造部 | 接收工序任务、组织生产和报工 |
| 操作工 | 生产制造部 | 执行工序、报工、报废记录 |
| 设备工程师 | 设备工程部 | 处理设备故障、评估停机影响 |
| 财务会计 | 财务驻厂组 | 开票、成本口径核对 |
| AI Agent 管理员 | IT 与数字化中心 | 配置 Agent 权限和工具边界 |

## 9. 关键主数据对象

下列 28 个对象是 MVP 建模优先级最高的主数据和业务对象。后续可映射为 IAOS metadata entity、seed 数据或场景运行对象。

| 序号 | 对象 | 业务含义 | 关键字段 | 关系 | MVP |
| --- | --- | --- | --- | --- | --- |
| 1 | 集团 | 企业最高组织 | group_code, group_name | 下辖事业部和法人 | 必须 |
| 2 | 事业部 | 区域经营单元 | bu_code, bu_name | 下辖销售公司和工厂 | 必须 |
| 3 | 法人公司 | 签约和财务主体 | legal_entity_code, tax_id | 关联客户、供应商、发票 | 必须 |
| 4 | 工厂 | 制造运营主体 | plant_code, plant_name, location | 下辖车间、仓库、产线 | 必须 |
| 5 | 部门 | 组织分工 | dept_code, dept_name | 关联人员和权限 | 必须 |
| 6 | 班组 | 生产执行组织 | team_code, shift_pattern | 关联员工和工序任务 | 必须 |
| 7 | 客户 | 销售对象 | customer_code, priority, payment_terms | 关联销售订单 | 必须 |
| 8 | 供应商 | 采购来源 | supplier_code, rating, lead_time | 关联采购订单和来料批次 | 必须 |
| 9 | 物料 | 原料、半成品、成品 | material_code, type, uom, lot_control | 关联 BOM、库存、订单 | 必须 |
| 10 | BOM | 产品结构 | bom_code, parent_material, component, qty | 关联物料和 MRP | 必须 |
| 11 | 工艺路线 | 产品制造路径 | routing_code, operations | 关联工作中心和工序任务 | 必须 |
| 12 | 工序 | 路线中的步骤 | operation_code, sequence, standard_time | 关联工作中心 | 必须 |
| 13 | 工作中心 | 能力核算单元 | work_center_code, capacity_per_shift | 关联设备、工序 | 必须 |
| 14 | 设备 | 具体生产资产 | equipment_code, status, work_center | 关联故障和产能 | 必须 |
| 15 | 工装 | 辅助生产资源 | tooling_code, lifetime, status | 关联工序和设备 | 可选 |
| 16 | 仓库 | 库存管理地点 | warehouse_code, type | 下辖库区和库位 | 必须 |
| 17 | 库位 | 库存最小位置 | location_code, warehouse_code | 关联库存批次 | 必须 |
| 18 | 班次 | 产能时间模型 | shift_code, start_time, end_time | 关联班组和产能 | 必须 |
| 19 | 员工 | 真实岗位执行者 | employee_code, role, dept | 关联任务和权限 | 必须 |
| 20 | 销售订单 | 客户需求 | order_no, customer, material, qty, due_date | 触发 O2D | 必须 |
| 21 | 采购订单 | 对供应商采购 | po_no, supplier, material, qty, promised_date | 关联收货和延期 | 必须 |
| 22 | 收货单 | 供应商到货记录 | receipt_no, po_no, lot_no, qty | 触发 IQC | 必须 |
| 23 | 检验单 | 质量判定 | inspection_no, object_type, result, defect_code | 关联收货或工序 | 必须 |
| 24 | 库存事务 | 库存变化流水 | transaction_no, material, lot, qty, direction | 关联仓库和业务单据 | 必须 |
| 25 | 生产订单 | 工厂生产指令 | work_order_no, material, qty, due_date, status | 生成工序任务 | 必须 |
| 26 | 工序任务 | 现场执行任务 | task_no, operation, work_center, planned_qty | 关联报工和检验 | 必须 |
| 27 | 发货单 | 客户交付记录 | shipment_no, customer, qty, ship_date | 关联销售订单 | 必须 |
| 28 | 质量问题 | 异常质量对象 | issue_no, source, severity, containment | 触发质量 Agent | 必须 |

后续扩展对象：

- 设备故障。
- 发票。
- 成本对象。
- 客户投诉。
- 供应商评分。
- CAPA。

说明：设备故障和发票在 MVP 事件中出现，但可以第一阶段先作为业务事件 payload 或轻量对象处理，待异常闭环和财务闭环扩大时再提升为完整主对象。

## 10. 关键事件

下列 18 个事件构成 MVP 第一阶段事件流。事件命名用于业务蓝图，落地到 IAOS 时应与 `shared/eventdef` 的 subject 和 payload 规范对齐。

| 序号 | 事件 | 触发条件 | 核心 payload | 下游影响 | Agent |
| --- | --- | --- | --- | --- | --- |
| 1 | `CustomerOrderReceived` | 销售接收客户订单或追加订单 | customer_code, material_code, qty, due_date | 创建销售订单草稿 | 经营分析 Agent |
| 2 | `SalesOrderConfirmed` | 销售订单审核确认 | sales_order_id, confirmed_qty, due_date | 触发 MRP | 计划 Agent |
| 3 | `MRPGenerated` | MRP 运算完成 | demand_id, material_shortages, capacity_risks | 生成采购需求和生产建议 | 计划 Agent |
| 4 | `PurchaseRequirementCreated` | MRP 发现物料缺口 | material_code, required_qty, need_by_date | 采购员释放 PO | 计划 Agent |
| 5 | `PurchaseOrderReleased` | 采购订单发布给供应商 | po_id, supplier_code, material_code, qty, promised_date | 等待供应商发货 | 计划 Agent |
| 6 | `SupplierShipmentDispatched` | 供应商发货 | po_id, lot_no, shipped_qty, eta | 仓库准备收货 | 计划 Agent |
| 7 | `SupplierDeliveryDelayed` | 供应商更新 ETA 或超期未到 | po_id, delay_days, reason | 重算交期风险 | 计划 Agent |
| 8 | `MaterialReceived` | 仓库完成收货 | receipt_id, material_code, lot_no, received_qty | 触发 IQC | 质量 Agent |
| 9 | `IncomingInspectionPassed` | 来料检验合格 | inspection_id, lot_no, accepted_qty | 入库可用库存 | 质量 Agent |
| 10 | `IncomingInspectionFailed` | 来料检验不合格 | inspection_id, lot_no, rejected_qty, defect_code | 隔离、退货或让步 | 质量 Agent |
| 11 | `InventoryPutawayCompleted` | 物料完成上架 | material_code, lot_no, location_code, qty | 更新可用库存 | 计划 Agent |
| 12 | `ProductionOrderReleased` | 计划员释放生产订单 | work_order_id, material_code, planned_qty, due_date | 生成工序任务 | 计划 Agent |
| 13 | `OperationStarted` | 班组开始工序任务 | task_id, operation_code, work_center_code, start_time | 更新产线状态 | 计划 Agent |
| 14 | `MachineDown` | 设备停机 | equipment_code, work_center_code, down_time, expected_repair_time | 重算产能和交期 | 计划 Agent |
| 15 | `OperationCompleted` | 工序报工完成 | task_id, completed_qty, scrap_qty, end_time | 推进下一工序或检验 | 质量 Agent |
| 16 | `ProcessInspectionFailed` | 过程检验失败 | task_id, defect_code, affected_qty | 隔离、返工或报废 | 质量 Agent |
| 17 | `FinishedGoodsReceived` | 完工入库 | work_order_id, material_code, lot_no, qty | 形成可发货库存 | 经营分析 Agent |
| 18 | `ShipmentDispatched` | 成品发货 | shipment_id, sales_order_id, qty, ship_date | 更新交付和收入预测 | 经营分析 Agent |

事件最低字段要求：

- `event_id`
- `event_type`
- `tenant_id`
- `occurred_at`
- `source_system`
- `correlation_id`
- `business_object_type`
- `business_object_id`
- `plant_code`
- `payload`

## 11. 第一条演示故事线

故事名称：

- 客户追加订单下的交付承诺重算。

演示目标：

- 展示 AESE 如何让虚拟企业在异常下继续运行。
- 展示 IAOS 如何通过事件流、MRP、Capability 和 Agent 形成建议。
- 展示管理层如何看到交付、成本、库存和质量风险的影响。

### 11.1 初始数据输入

组织：

| 字段 | 值 |
| --- | --- |
| 集团 | 华辰热管理系统集团 |
| 工厂 | 苏州制造基地 |
| 产线 | 电池冷却板 A 线 |
| 租户 | `tenant-hctm` |

客户订单：

| 字段 | 值 |
| --- | --- |
| 客户 | 星河新能源汽车有限公司 |
| 客户编码 | `CUST-SGNEV` |
| 产品 | 电池冷却板组件 A01 |
| 产品编码 | `HCTM-BCP-A01` |
| 原订单数量 | 10,000 件 |
| 追加订单数量 | 2,000 件 |
| 总需求数量 | 12,000 件 |
| 要求交期 | 2026-07-20 |
| 交付策略 | 可分两批发运，但最终交期不变 |

库存：

| 物料 | 编码 | 当前可用库存 |
| --- | --- | --- |
| 电池冷却板组件 A01 | `HCTM-BCP-A01` | 1,200 件 |
| 6061-T6 铝板 | `AL-PLATE-6061-T6` | 8,000 张 |
| 密封圈 | `SEAL-RING-A01` | 30,000 个 |
| 接头组件 | `FITTING-A01` | 18,000 套 |
| 包装箱 | `PKG-BOX-A01` | 1,500 个 |

采购和供应：

| 字段 | 值 |
| --- | --- |
| 主铝板供应商 | 安铝金属材料有限公司 |
| 主供应商编码 | `SUP-ALPHA-AL` |
| 在途铝板数量 | 5,000 张 |
| 原 ETA | 2026-07-08 |
| 新 ETA | 2026-07-11 |
| 延期天数 | 3 天 |
| 备选供应商 | 贝塔铝业有限公司 |
| 备选供应商规则 | 首批来料需加严检验 |

产能：

| 字段 | 值 |
| --- | --- |
| A 线常规日产能 | 800 件 |
| A 线加班日产能 | 1,100 件 |
| 焊接工作中心 | 关键瓶颈 |
| 故障设备 | `LAS-WLD-02` |
| 预计停机 | 8 小时 |
| 停机影响 | 当日焊接能力下降约 35% |

质量规则：

| 场景 | 规则 |
| --- | --- |
| 主供应商正常批次 | 常规抽检 |
| 备选供应商首批供货 | 加严抽检 |
| 气密检漏失败 | 批次隔离并触发质量分析 |
| 铝板批次异常 | 禁止直接投产，需质量放行 |

### 11.2 预期事件过程

```text
CustomerOrderReceived
-> SalesOrderConfirmed
-> MRPGenerated
-> PurchaseRequirementCreated
-> PurchaseOrderReleased
-> SupplierDeliveryDelayed
-> MachineDown
-> MRPGenerated
-> ProductionOrderReleased
-> MaterialReceived
-> IncomingInspectionPassed / IncomingInspectionFailed
-> InventoryPutawayCompleted
-> OperationStarted
-> OperationCompleted
-> FinishedGoodsReceived
-> ShipmentDispatched
```

说明：

- 第二次 `MRPGenerated` 表示异常发生后的重算。
- 若启用备选供应商，来料必须经过加严检验。
- 若焊接设备恢复晚于预期，计划 Agent 应重新提示交付风险。

### 11.3 预期系统输出

系统应输出：

| 输出 | 内容 |
| --- | --- |
| 交付风险 | 原交期存在风险，风险来源为铝板延期和焊接产能下降 |
| 物料缺口 | 按 12,000 件需求和 1.05 张用量计算，铝板总需求约 12,600 张；当前可用 8,000 张，考虑在途后仍受延期影响 |
| 产能缺口 | `LAS-WLD-02` 停机导致当日焊接瓶颈扩大 |
| 建议方案 A | 主供应商加急，A 线加班，维持单供应商 |
| 建议方案 B | 启用备选供应商补 2,000 张铝板，加严检验，A 线双班加班 |
| 建议方案 C | 分批发运，先发 8,000 到 9,000 件，剩余延后 1 到 2 天 |
| 推荐方案 | 方案 B + 分批发运兜底 |
| 管理层解释 | 成本上升来自加急费、加班费和加严检验；交付风险下降但质量检验负荷上升 |

### 11.4 Agent 预期输出

计划 Agent：

- 识别销售订单追加导致需求增加 20%。
- 识别铝板供给时间和焊接产能同时受限。
- 给出三种交付方案。
- 推荐启用备选供应商并安排 A 线加班。
- 说明每个方案对交期、成本和库存的影响。

质量 Agent：

- 提醒备选供应商首批铝板需要加严检验。
- 若检验失败，建议隔离批次并禁止投产。
- 若气密检漏失败率升高，建议追溯铝板批次、焊接设备和操作班组。

经营分析 Agent：

- 解释准时交付率、加班成本、采购成本和库存周转变化。
- 对比三个方案的经营影响。
- 输出给工厂厂长的简短决策摘要。

### 11.5 演示成功标准

MVP 演示成功必须满足：

- 能从一张客户订单触发完整事件链。
- 能在供应商延期和设备故障后重算风险。
- 能展示至少一个 Agent 基于业务上下文生成建议。
- 能把建议关联到具体业务对象，例如订单、采购订单、生产订单、设备和库存。
- 能留下事件和操作审计痕迹。
- 能给管理层展示可解释的经营影响。

## 12. IAOS 落地提示

后续开发可以按以下顺序推进：

1. 把第 9 节对象转成 IAOS metadata entity 设计。
2. 把第 6、7、8 节内容转成 seed 数据。
3. 把第 10 节事件映射到 `shared/eventdef` 和 scenario subject。
4. 扩展或复用 `scenarios/o2d`，先跑订单到工单和库存链路。
5. 把计划 Agent 的建议先做成只读解释，再接入受治理 Capability。
6. 事件流、对象关系和演示故事已在 M3 跑通；当前按 M3V 计划先实现只读 2D 预览，再接 IAOS 在线数据源。
