---
id: DES-013
title: M12 Project Genesis 产品工业化与量产批准
date: 2026-07-22
status: completed
author: Codex + User
tags: [m12, genesis, industrialization, apqp, ppap, trial-production]
---

# M12 Project Genesis 产品工业化与量产批准

## 1. 目标

M12 从 M11 的 `industrialization_eligible=true` 开始，让华辰苏州制造公司完成首个虚构客户电池冷却板项目的 RFQ、技术/产能/商业可行性、报价与定点，建立受版本治理的产品、BOM、工艺和质量策划，完成供应商/工装/首批物料准备、试生产、问题整改和客户 PPAP 批准，输出 M13 可消费的 `serial_production_eligible=true`。

该终态表示客户项目、产品/工艺版本和量产质量门已经获批，可以接收第一张正式订单；它不表示订单已经下达、正式量产已经发生，也不产生发货、开票、收入或回款。

## 2. 关键现实约束

M11 终态现金为 8,500,000 CNY，同时仍要保留 3,000,000 CNY 工资准备金、5,000,000 CNY 最低现金缓冲以及既有合同义务。M12 可自由使用的现金不足以凭剧情覆盖产品开发、工装和试制，因此必须先冻结项目资金方案：

- 客户定点后允许支付工装预付款或开发费，但到账前不能当现金使用。
- 客户预付款在交付/结算条件满足前是合同负债，不是销售收入或利润。
- 工装、样件、首批物料、外部试验和试生产必须受项目预算、付款节点和现金缓冲约束。
- 报价成本、目标成本、试制实际成本、量产标准成本和现金流是不同数值。

具体 RFQ 量纲、目标价格、开发预算、预付款、样件数量、试制批次、质量门和时间线由 M12 D0 冻结为虚构机器基线。

## 3. 纵向业务链

```text
消费 M11 industrialization_eligible
-> 虚构客户发出 RFQ 与技术要求
-> 技术/产能/供应链/财务可行性和受治理报价
-> 客户定点与项目资金到位
-> 产品 revision、BOM、routing、PFMEA、control plan 发布
-> 供应商、工装、量检具和首批物料准备
-> APQP gate 与试生产计划批准
-> 实际试生产、测量、节拍、良率和追溯结果
-> 焊接/泄漏能力不足 discrepancy
-> 受治理工程变更、整改和复试
-> PPAP package 提交与客户实际批准
-> serial_production_eligible
```

## 4. 与旧 HCTM 场景和 M13 的边界

现有 `scenario-packs/hctm` 已为 O2D 演示预置 `HCTM-BCP-A01`、`BOM-HCTM-BCP-A01-V1` 和 `RT-HCTM-BCP-A01-V1`。这些是兼容 fixture，不是 Genesis 中产品已经开发完成的历史证据。

M12 必须在 Industrialization campaign 中通过 RFQ、工程发布和 PPAP 因果链生成相同 stable code 的版本化定义，并输出 release manifest/hash。运行时不反向改写旧 fixture；M13 adapter 只在 hash、版本和关键字段匹配时把 M12 terminal contract 投影到现有 O2D 对象。任何不匹配必须失败关闭并要求显式 migration。

M12 不接收正式销售订单，不执行正式 MRP/O2D，不形成客户发运、发票、应收或回款；这些属于 M13。

## 5. 三态所有权

| 事实 | World State | IAOS Business State | Actor Knowledge |
| --- | --- | --- | --- |
| 客户真实 RFQ、定点和 PPAP 决定 | Own：外部要求、响应时间、批准/拒绝 | RFQ、报价、客户项目、提交与审批记录 | 销售/项目团队在 observation 后获知 |
| 产品和工艺定义 | 物理可实现性、真实试制表现和版本 evidence ref | Own：产品 revision、BOM、routing、规范、变更 | 按工程角色与发布状态获知 |
| 供应商和材料真实能力 | Own：交期、批次属性、证书真实性和来料结果 | 供应商资格、采购、批次和检验记录 | 采购/质量按 observation 获知 |
| 试生产和质量结果 | Own：实际产出、报废、节拍、测量、泄漏和追溯 | 试制订单、检验、问题、CAPA、APQP/PPAP 状态 | 项目角色按职责和报告获知 |
| 项目资金 | Own：客户预付款实际到账、现金流出 | 预算、合同负债、承诺、应付和付款审批 | CFO/项目负责人按权限获知 |

IAOS 中“工艺已发布”“试制已完成”或“PPAP 已提交”不能伪造物理产品能力或客户批准；外部客户、供应商和实验室是 AESE World 策略，不是 IAOS 用户。

## 6. 最小领域模型

### AESE World

- `CustomerOpportunity`、`CustomerRequirement`、`CustomerDecision`。
- `ProductPrototype`、`ToolingUnit`、`MaterialLot`、`SupplierCapability`。
- `TrialBuild`、`TrialUnit`、`MeasurementResult`、`ProcessCapabilityResult`。
- `QualityIssue`、`ContainmentAction`、`CorrectiveAction`、`VerificationResult`。
- `PPAPSubmission`、`PPAPDecision`、`IndustrializationGate`。

### IAOS 最小管理对象

- `customer_rfq`、`feasibility_review`、`quotation`、`customer_nomination`
- `customer_project`、`apqp_plan`、`apqp_gate`
- `product_revision`、`engineering_specification`
- `engineering_bom`、`manufacturing_bom`、`routing_revision`
- `process_flow`、`pfmea`、`control_plan`、`work_instruction`
- `tooling_project`、`supplier_qualification`、`material_approval`
- `trial_order`、`inspection_plan`、`measurement_study`
- `quality_issue`、`engineering_change`、`corrective_action`
- `ppap_package`、`customer_approval`、`production_release`

文档型对象首版保存结构化状态、版本、owner、审批、内容 hash 和 evidence ref，不在 AESE/IAOS 中重建完整 PLM、QMS 或文档管理产品。

## 7. APQP 与版本治理

首版采用可机器验证的简化 APQP gate：

1. `GATE-RFQ-FEASIBILITY`：客户要求、产能、投资、供应链和风险可行。
2. `GATE-PRODUCT-DESIGN`：产品 revision、设计规范和 EBOM 已审查发布。
3. `GATE-PROCESS-DESIGN`：MBOM、routing、process flow、PFMEA、control plan 和 work instruction 一致。
4. `GATE-SUPPLIER-TOOLING`：关键供应商、材料、工装、量检具和外部试验资源就绪。
5. `GATE-PRODUCT-PROCESS-VALIDATION`：试制、MSA、尺寸/功能、节拍、良率、能力和追溯通过。
6. `GATE-PPAP`：提交内容完整且客户实际批准。

产品/BOM/routing/control plan 必须携带 revision、effective time、状态和 hash。工程变更必须新建 revision 或受治理变更集，禁止原地篡改已用于试制或 PPAP 的版本。

## 8. 供应商、工装和首批物料

首版只覆盖 `HCTM-BCP-A01` 所需关键铝板、密封件、接头、包装和追溯标签，以及焊接工装、检漏/尺寸量检具。供应商和客户全部虚构。

- 已有通用设备不等于产品专用工装就绪。
- 供应商管理记录不等于真实过程能力；必须有样件/批次 observation 和检验结果。
- 未批准材料、无证书/追溯或检验失败的批次不能进入试制。
- 试制物料、样件和报废不能静默成为 M13 可销售库存。

## 9. 试生产、质量与成本不变量

- 没有 `industrialization_eligible=true` 不能启动可执行 M12 项目。
- 未获得客户定点、项目预算和开发资金，不能释放工装订单或试制计划。
- RFQ 数量/价格只是机会和报价，不是正式需求、订单、收入或应收。
- 未发布且相互一致的 product/BOM/routing/control plan 不能发起试制。
- 试制只能消费已批准且实际入库/放行的材料，并使用 M11 accepted 设备和有效资格人员。
- 每个 trial unit/lot 必须保留材料、工装、设备、人员、参数 revision 和检验结果追溯。
- 良率、节拍、Cp/Cpk、MSA 和功能结果由 World 计算；IAOS 状态或 UI 不能直接改写。
- 未关闭重大质量问题、未完成变更复验、PPAP package 不完整或客户未批准时，量产资格为 false。
- 客户预付款计入现金和合同负债，不计销售收入；试制成本不等于量产标准成本。
- PPAP 样件和试制件默认不可销售，不计入 M13 成品库存。

## 10. 最小异常 tracer

M12 固定一条焊接/泄漏过程能力不足链：

```text
IAOS 试制计划与控制计划认为焊接参数 revision A 可达标
-> World 首轮试制出现泄漏超限，Cpk 低于冻结阈值
-> World/IAOS 产生 discrepancy，项目/质量角色起初只知道局部结果
-> measurement observation 送达后形成 containment 和 engineering-change intent
-> IAOS committed outcome 发布参数/工装 revision B 与复试计划
-> AESE 计算材料消耗、报废、设备/人员占用、工期和成本后果
-> 第二轮试制与能力复验通过，关闭问题
-> PPAP 重新提交并等待客户实际决定
```

Agent 只能建议和提交 intent，不能直接改参数、测量结果、Cpk 或客户批准状态。

## 11. 角色与治理

| 角色 | 主要动作 | 关键限制 |
| --- | --- | --- |
| 销售/客户项目负责人 | RFQ、报价、定点和客户沟通 | 不能伪造客户决定或批准自身价格例外 |
| 产品/工艺工程师 | 产品、BOM、routing、PFMEA、控制计划和变更 | 不能覆盖已发布/已试制 revision |
| 质量负责人 | MSA、检验、问题、PPAP 完整性和放行建议 | 不能修改 World 测量值或单独关闭重大问题 |
| 采购负责人 | 供应商、工装和试制物料准备 | 不能绕过材料批准、预算或付款门 |
| 厂长 / CFO | 产能、项目预算和生产放行审批 | 不能把报价、预付款或试制件当收入/库存 |

人类和 Agent 复用同一 Capability、Decision、Policy、Process 和审计；关键工程、质量、商业和生产批准保持职责隔离。

## 12. Bridge payload family

M12 增加严格 allowlist 类型：

```text
genesis.rfq.received.v1
genesis.quotation.approved.v1
genesis.customer.nomination.received.v1
genesis.development.funding.received.v1
genesis.product.revision.released.v1
genesis.process.revision.released.v1
genesis.supplier.material.approved.v1
genesis.trial.build.started.v1
genesis.trial.build.completed.v1
genesis.process.capability.failed.v1
genesis.engineering.change.approved.v1
genesis.corrective.action.verified.v1
genesis.ppap.submitted.v1
genesis.ppap.approved.v1
genesis.production.release.approved.v1
```

所有类型复用 DES-008 envelope、stable ref、tenant journal/cursor、幂等和 committed outcome 语义。

## 13. Pack 与版本策略

- `hctm-genesis` 升级到下一 minor version，新增 `campaigns/industrialization/`。
- 初态验证 M11 terminal hash、`industrialization_eligible`、现金/准备金/承诺、设备能力、人员资格和 shift。
- 终态输出 product/BOM/routing/control-plan/PPAP release manifest 与 canonical hash。
- `scenario-packs/hctm` 保持回归 fixture；M12 通过兼容性测试证明 stable code 与 M13 adapter 可映射，不把 fixture 当 Genesis 输入事实。
- M13 只能消费 M12 机器输出，不从界面、报价、定点或 PPAP 文案推断量产已批准。

## 14. 非目标

- 第二客户、第二产品族、多工厂或复杂平台化产品组合。
- 完整 CRM、CPQ、PLM、APQP/QMS、SRM、项目会计或文档管理产品。
- 正式客户订单、MRP、批量采购、正式量产、交付、开票、应收、回款和项目盈亏。
- 真实客户、供应商、图纸、个人、认证机构或商业价格数据。
- 高精度物理/流体仿真、CAD/CAE、3D 工厂或数字孪生。

## 15. 完成标准

- 单一 run 从 M11 eligibility 确定性推进到 `serial_production_eligible=true`。
- RFQ、可行性、报价、客户定点、开发资金和项目预算因果链完整且金额语义正确。
- 产品/BOM/routing/PFMEA/control plan/工装/供应商/材料版本一致并可机器校验。
- 试制物料、设备、人员、样件、测量、报废、成本和追溯满足守恒。
- 焊接泄漏/Cpk 异常、知情延迟、受治理变更、复试和 PPAP 批准形成完整因果链。
- 人类/Agent 共用治理能力，越权、自批、篡改 revision、伪造结果和未批准放行失败关闭。
- M7-M11 回归、M3/O2D 兼容映射、两仓测试、API/UI、runbook/evidence、revision 和 Atlas 完整。
