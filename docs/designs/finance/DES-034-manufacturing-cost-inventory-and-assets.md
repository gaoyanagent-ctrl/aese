---
id: DES-034
title: 制造成本、存货与固定资产
date: 2026-07-30
status: active
author: Codex + User
tags: [finance, manufacturing-cost, inventory, fixed-assets]
---

# 制造成本、存货与固定资产

## 1. 目的

把物料、人工、机器、能耗、制造费用、在制品、完工、质量损失和设备建设转换为可追踪的
制造成本与资产事实。生产实物事实仍由业务系统/World 拥有，会计只记录价值后果。

## 2. 成本与存货 Entity

- `cost_element`
- `cost_center`
- `profit_center`
- `activity_type`
- `cost_object`
- `cost_pool`
- `allocation_cycle`
- `standard_cost_version`
- `standard_cost_component`
- `material_ledger_entry`
- `production_cost_entry`
- `wip_balance`
- `manufacturing_variance`
- `inventory_valuation`
- `cost_settlement_run`

维度至少包含法人、工厂、部门、成本中心、利润中心、项目、产品、订单、生产订单和批次。

## 3. 标准成本

标准成本由有效 BOM、工艺、采购价、人工/机器费率、能耗和制造费用率滚算：

```text
材料 + 直接人工 + 机器 + 能源 + 变动制造费用 + 固定制造费用
```

版本经历 draft → simulated → approved → released → frozen。发布后不能原地修改，后续版本按
生效日切换。Agent 可解释差异和缺失输入，不能无审批发布标准成本。

## 4. 实际成本与 WIP

实际成本按受治理业务事件归集：

- 领料/退料、报废、返工、替代料；
- 人员工时、机器工时、停机、能耗；
- 外协、运费、检验和质量损失；
- 生产入库、WIP 结转和订单结算。

支持分批法、订单法和标准成本加差异。禁止用销售价格反推生产成本，禁止用 Preview KPI
冒充实际成本。

## 5. 存货核算

存货数量来自受治理实物交易，价值来自估价规则。至少支持原材料、WIP、产成品、在途、
委外和质量冻结库存。成本层、批次和估价范围必须明确法人/工厂/账簿。

库存子账与总账控制科目按期间对账；负库存、无成本层、跨法人无内部交易引用等情况失败关闭。

## 6. 固定资产与 CIP

Entity：

- `asset_class`
- `fixed_asset`
- `asset_book`
- `construction_in_progress`
- `asset_component`
- `depreciation_rule`
- `depreciation_run`
- `asset_transfer`
- `asset_impairment`
- `asset_disposal`

设备采购、安装、调试先进入 CIP；达到可使用状态且验收后转固。资产可按账簿采用不同折旧、
税务和估值规则。验收、转固、折旧、盘点、减值和处置均保留证据与审批。

## 7. Capability / Process

- `standard.cost.rollup/simulate/approve/release`
- `production.cost.capture`
- `wip.calculate`
- `manufacturing.variance.calculate`
- `cost.settlement.run`
- `inventory.valuation.run`
- `asset.cip.capture`
- `asset.capitalize`
- `asset.depreciation.run`
- `asset.transfer/impair/dispose`

## 8. 验收

- 生产订单成本可分解至物料、人工、机器、能耗和费用；
- WIP、完工、差异与库存/总账可对账；
- 设备从采购/CIP/验收到转固/折旧全链可追溯；
- 重复事件不重复计价，迟到事件按期间政策处理；
- World 实物、业务子账、成本子账和总账数量/金额边界清晰。
