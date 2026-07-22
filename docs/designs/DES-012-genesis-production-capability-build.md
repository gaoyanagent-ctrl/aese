---
id: DES-012
title: M11 Project Genesis 生产能力建设
date: 2026-07-22
status: completed
author: Codex + User
tags: [m11, genesis, equipment, workforce, capability]
---

# M11 Project Genesis 生产能力建设

## 1. 目标

M11 从 M10 的 `capability_build_eligible=true` 开始，让华辰苏州制造公司在资金、空间、供应周期、设备能力和人员技能约束下，完成电池冷却板 A 线所需设备与检测仪器的采购、交付、安装、调试、校准，以及核心团队的编制、招聘、培训、认证和班次建立，输出 M12 可消费的 `industrialization_eligible=true`。

该终态表示“设施、通用制造资源和人员能力已具备，可以启动产品工业化”，不表示某个产品已经完成 BOM/工艺发布、APQP、试生产、PPAP 或 SOP。

## 2. 现实资金与范围决策

M10 终态只有 10,000,000 CNY 现金，同时仍有设施合同未付款承诺；现有首年预算也不能自动覆盖整条生产线。因此 M11 必须先冻结资金来源和新 CAPEX/人员预算，禁止凭剧情生成设备：

- M9 尚有 10,000,000 CNY 认缴资本未实缴，可经受治理资本催缴和真实到账链使用。
- 设备方案允许采购、融资租赁和供应商账期的组合，但首版不实现通用贷款、利息或复杂融资产品。
- 必须保留设施尾款、至少六个月核心团队人工成本和最低营运现金缓冲。
- 设备合同承诺、应付、已付款、租赁义务、工资预算和现金分别记录。

具体金额、交付周期、人员数量和缓冲值由 M11 C0 冻结为虚构机器基线；未关闭资金缺口前不得签署可执行订单或发出录用通知。

## 3. 纵向业务链

```text
消费 M10 capability_build_eligible
-> 冻结通用能力/产能需求与资金方案
-> 批准 CAPEX、编制和采购策略
-> 设备询价、比选、订单/租赁与付款门
-> 供应商制造、运输、到货和安装
-> 单机调试、校准、安全与能力验收
-> 组织定岗、招聘、入职、培训、认证和班次
-> 设备故障/招聘缺口 observation
-> 受治理纠正、重排与复验
-> 设备、实验室、仓储和人员资格联合门
-> industrialization_eligible
```

## 4. M11 与 M12 的边界

M11 使用版本化 `CapabilityRequirement` 表达冷却板产品族所需的成形、机加工、焊接、清洗、检漏、装配、包装、仓储和质量测量能力，只用于设备选型、空间落位和培训准备。

M11 不发布产品 BOM、正式 routing、控制计划或工艺参数，不以空跑/标准件调试替代产品试制。具体客户产品、RFQ、成本报价、APQP、工装、首件、试生产、PPAP 和 SOP 属于 M12。

## 5. 三态所有权

| 事实 | World State | IAOS Business State | Actor Knowledge |
| --- | --- | --- | --- |
| 供应商真实交期与设备状态 | Own：制造、运输、到货、安装、故障和能力 | 询价、订单、资产与项目记录 | 角色只知道已送达的进展/异常 |
| 现金和实际付款 | Own：银行到账、现金流出和租赁现实义务 | 资本催缴、预算、承诺、应付和付款审批 | CFO 按 observation 与权限获知 |
| 人员真实可用性与技能 | Own：候选接受、到岗、出勤、训练结果和技能 | 编制、招聘、员工、培训与资格记录 | HR/主管按职责获知，不暴露无关候选信息 |
| 设备/人员管理状态 | 实际状态的 stable ref 与必要投影 | Own：采购、固定资产、组织、岗位、班次和认证 | 通过受治理查询和 observation 获知 |

采购订单、资产卡、员工档案或培训记录不等于设备已经可用、人员已经到岗或技能已经掌握。World consequence 产生后，IAOS 才能受治理接受、登记或放行。

## 6. 最小领域模型

### AESE World

- `CapabilityRequirement`：能力编码、目标范围、精度/节拍边界和验证方法。
- `EquipmentOption`：供应商、购置方式、价格、交期、能力、能耗和风险。
- `EquipmentUnit`：stable code、所在 zone、生命周期、状态和能力证据。
- `DeliveryActivity`、`InstallationTask`、`CommissioningResult`、`CalibrationResult`。
- `CandidateProfile`：完全虚构的技能、可用日期、期望条件和接受概率输入。
- `Worker`、`Skill`, `TrainingSession`、`QualificationResult`、`ShiftAssignment`。
- `CapabilityGate`：设备、实验室、仓储、人员、资金和安全门的可解释结果。

### IAOS 最小管理对象

- `capital_call`、`capex_budget`、`headcount_budget`
- `equipment_requirement`、`supplier_bid`、`procurement_request`
- `purchase_or_lease_order`、`goods_receipt`、`equipment_asset`
- `installation_work_order`、`commissioning_acceptance`、`calibration_record`
- `org_position`、`headcount_requisition`、`candidate_application`
- `employment_offer`、`employee_assignment`、`training_record`
- `qualification_record`、`shift_roster`

优先使用 IAOS metadata/config package、Capability、Process、Policy 和 Decision。M11 不把 IAOS 扩展为完整采购套件、HRIS、LMS、EAM 或财务系统。

## 7. 设备与空间范围

首版至少覆盖：成形、CNC、激光焊接、清洗、检漏、装配/包装，以及最小计量/质量实验室。设备必须落位到 M10 已验收的 zone，并满足面积、公用工程、安全间距和容量约束。

M8 的 `LAS-WLD-02` 是独立设备偏差 tracer，不能静默当作 M11 已购资产；若选择复用其编码，必须建立明确的采购、交付、验收和 pack migration 证据。

设备状态至少经过：

```text
planned -> sourced -> ordered -> in_transit -> received
-> installed -> commissioning -> accepted | remediation
```

## 8. 组织、技能与隐私范围

首版岗位族覆盖厂长、计划、采购、质量、工艺、设备、仓储、操作工和检验员。只实现一班制最小团队及关键替补约束；多班复杂排班、薪酬社保和绩效后置。

- 候选人和员工全部虚构，使用 stable actor code，不保存真实个人信息。
- 招聘记录属于 IAOS；实际接受、到岗、出勤和学习结果属于 World State。
- 岗位资格由必需培训、考试/实操结果、有效期和授权范围共同决定。
- 设备供应商、培训机构和候选劳动力市场是 AESE 外部世界策略，不是 IAOS 用户。

## 9. 资金、能力和治理不变量

- 没有 `capability_build_eligible=true` 不能创建可执行 M11 能力项目。
- 未批准资金来源、CAPEX 和 headcount envelope，不能提交设备订单或 offer。
- 订单/租赁承诺与设施遗留承诺之和不能超过批准 envelope；实际付款不能超过可用现金。
- 资本承诺不是现金；只有 World 银行到账 observation 和 IAOS 入账 outcome 对账后才可使用。
- 未到货不能安装，未安装不能调试，未通过安全/校准/能力验收不能 accepted。
- IAOS 资产状态不能伪造 World 设备可用性；虚拟时间推进本身不能完成供应、安装或培训。
- 未实际到岗、未完成必需培训或资格过期的人员不能计入 capability gate。
- 设备和人员数量必须覆盖关键能力、最低安全角色和首版班次；任一强制门失败时资格为 false。
- `industrialization_eligible=true` 不代表产品或量产已批准。

## 10. 最小异常 tracer

M11 固定一条关键检漏设备调试失败链：

```text
IAOS 计划认为检漏设备将按期验收
-> World 调试发现校准漂移，实际能力不合格
-> World/IAOS 产生 discrepancy，质量负责人起初未知
-> observation 送达后项目负责人、设备工程师和质量负责人获知
-> 提交 remediation / reinspection intent
-> IAOS committed outcome 更新工单、预算和计划
-> AESE 计算供应商整改、人员占用、复验时间与现金后果
-> 校准和能力复验通过后关闭 discrepancy
```

该链必须证明“关闭工单”“登记资产”或 UI 点击不能直接创造设备能力。

## 11. 角色与治理

| 角色 | 主要动作 | 关键限制 |
| --- | --- | --- |
| 工厂项目负责人 | 能力项目、综合进度和联合资格申请 | 不能批准自身越权预算/付款 |
| CEO / CFO | 资本、CAPEX、编制和重大合同审批 | 不能把认缴或授信当现金 |
| 采购负责人 | 询价、比选、订单和供应商跟催 | 不能自批供应商选择或付款 |
| HR 负责人 | 编制、招聘、offer、入职与培训组织 | 不能伪造候选接受/到岗/考试结果 |
| 设备工程师 / 质量负责人 | 安装调试、校准、能力和安全验收 | 职责冲突时需要独立接受者 |

岗位尚未招聘到位时，只允许已有项目负责人在 mandate 内发起准备动作；受治理批准仍由已生效 CEO/CFO 承担。人类与 Agent 使用同一 Capability、Policy、Decision 和审计链。

## 12. Bridge payload family

M11 增加严格 allowlist 类型：

```text
genesis.capital.call.requested.v1
genesis.capital.received.v1
genesis.capability.budget.approved.v1
genesis.equipment.order.approved.v1
genesis.equipment.delivered.v1
genesis.equipment.commissioning.failed.v1
genesis.equipment.remediation.approved.v1
genesis.equipment.accepted.v1
genesis.headcount.approved.v1
genesis.candidate.offer.approved.v1
genesis.worker.onboarded.v1
genesis.worker.qualified.v1
genesis.shift.activated.v1
genesis.capability.accepted.v1
```

所有类型复用 DES-008 envelope、stable ref、tenant journal/cursor、幂等和 committed outcome 语义。

## 13. Pack 与版本策略

- `hctm-genesis` 升级到下一 minor version，新增 `campaigns/capability-build/`。
- 初态显式验证 M10 terminal hash、`capability_build_eligible`、closing cash、设施遗留承诺、空间和 utility capacity。
- M8/M9/M10 campaign 必须继续可独立运行；旧设备 tracer 不得被无证据迁移成生产资产。
- M12 只消费 M11 的机器资格，不从画布、人员数量或资产文案推断工业化已就绪。

## 14. 非目标

- 客户 RFQ、正式产品设计、BOM、routing、控制计划、APQP、PPAP 和 SOP。
- 产品专用工装、首批原材料采购和带产品试生产。
- 完整采购、SRM、HRIS、LMS、排班、工资、固定资产会计、EAM 或融资产品。
- 真实候选人、员工、供应商或金融机构数据和接口。
- 多工厂、多条产线、多班复杂排程和 3D 设备布局。

## 15. 完成标准

- 单一 run 从 M10 eligibility 确定性推进到 `industrialization_eligible=true`。
- 资金来源、新预算、设施尾款、设备/租赁承诺、工资准备金和现金可机器对账。
- 关键设备全部实际到货、安装、调试、校准、安全及能力验收，且空间/utility 守恒。
- 最小核心团队实际到岗并具备有效岗位资格，一班制和职责隔离门通过。
- 检漏设备失败、知情延迟、受治理整改、复验和差异关闭形成完整因果链。
- 人类/Agent 共用治理能力，越权、自批、超预算、未验收付款和伪造资格全部失败关闭。
- M7-M10 回归、两仓测试、API/UI、runbook/evidence、revision 和 Atlas 完整。
