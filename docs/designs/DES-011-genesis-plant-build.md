---
id: DES-011
title: M10 Project Genesis 工厂选址与设施建设
date: 2026-07-22
status: completed
author: Codex + User
tags: [m10, genesis, plant, site-selection, construction, project]
---

# M10 Project Genesis 工厂选址与设施建设

## 1. 目标

M10 从 M9 的 `plant_project_eligible=true` 开始，让华辰苏州制造公司在资金、工期、空间、承包商和公用工程约束下完成候选场址评估、项目审批、场地使用权取得、厂房改造和设施验收，输出 M11 可消费的 `capability_build_eligible=true`。

首版只建设苏州制造基地一期。M10 建成的是可以安装设备、布置仓储/质量区域并组建团队的设施载体，不代表电池冷却板 A 线已经形成生产能力。

## 2. 业务现实约束

M9 终态提供：

- 公司实际现金：20,000,000 CNY。
- 首年预算授权：15,000,000 CNY。
- 已生效 CEO、CFO 和工厂项目负责人岗位及 mandate。

因此 M10 不默认选择昂贵绿地自建。候选方案至少包含：

- 绿地购地自建：周期长、资金需求高，首版可被规则判定不可行。
- 标准厂房租赁改造：资金和工期较低，适合一期 tracer。
- 定制代建/长期租赁：成本和控制权居中，但存在交付依赖。

最终选择必须来自版本化多维评估和 IAOS 受治理决策，不能由剧情直接指定赢家。所有候选地点、园区和外部机构保持虚构。

## 3. 纵向业务链

```text
消费 M9 plant_project_eligible
-> 建立产能与设施需求
-> 生成并调研候选场址
-> 评估资金、工期、物流、人力、公用工程和风险
-> 项目总监提出推荐方案
-> CEO / CFO 经 IAOS 审批选址与投资上限
-> 签署租赁/场地使用协议并取得实际场地控制
-> 建立设施建设项目与 WBS
-> 完成设计、许可、改造和公用工程接入
-> 处理公用工程延期并受治理重排
-> 完成消防/EHS/设施验收
-> 输出 capability_build_eligible
```

## 4. 三态所有权

| 事实 | World State | IAOS Business State | Actor Knowledge |
| --- | --- | --- | --- |
| 候选地真实条件 | Own：距离、容量、风险、可用日期 | 调研记录、评分和决策 | 角色只知道已调研/送达的信息 |
| 场地是否实际可占用 | Own：交付、钥匙/控制权、生效时间 | 租赁/使用协议和审批 | 项目总监获交付 observation 后获知 |
| 施工真实进度与质量 | Own：活动、资源消耗、返工、完成度 | 项目/WBS、里程碑、合同、付款记录 | 按现场报告、验收和权限获知 |
| 项目预算和承诺 | 保存实际现金、实际付款和 IAOS envelope 引用 | Own：预算、承诺、付款审批 | CEO/CFO 按 mandate 可见 |
| 设施空间和公用工程 | Own：区域、面积、容量、连接和验收状态 | 设施台账/项目交付记录 | 通过 observation/台账权限获知 |

IAOS 里程碑完成记录不等于现场已经完工。AESE 先计算实际施工结果，验收 observation 进入 IAOS 后，受治理 Capability 才能提交里程碑接受、付款或项目关闭。

## 5. M10 最小领域模型

### AESE World

- `SiteOption`：位置、使用方式、面积、成本、可用日期和风险。
- `SiteAssessment`：指标、来源、观察时间、置信度和评分版本。
- `SpatialNode`：region/city/park/site/building/floor/zone 层级与坐标。
- `UtilityCapacity`：electricity/water/gas/compressed_air/environmental capacity。
- `FacilityProject`：目标、状态、基准工期、预算引用和负责人。
- `WorkPackage`：前置依赖、持续时间、资源需求、成本和交付物。
- `ExternalPartyCapacity`：承包商/园区/公用工程服务能力和日历。
- `FacilityAsset`：建筑、办公区、生产区、仓储区、质量区和公辅区的实际状态。
- `InspectionResult`：消防、EHS、建筑和公用工程验收世界事实。

### IAOS 最小管理对象

- `site_option`、`site_assessment`
- `investment_request`、`site_decision`
- `land_or_lease_agreement`
- `facility_project`、`project_wbs`
- `contractor_contract`
- `project_milestone`、`change_request`
- `payment_request`
- `facility_acceptance`

优先用 metadata/config package、Process、Policy、Decision 和 Capability。M10 不把 IAOS 扩展成通用建筑项目管理产品。

## 6. 空间与设施范围

M10 最小空间层级：

```text
China / Jiangsu / Suzhou
-> fictional industrial park
-> HCTM Suzhou site
-> main building
-> office / production / warehouse / quality / utility zones
```

每个节点至少携带 stable code、parent、二维坐标/边界、面积、用途、占用状态、容量和验收状态。M10 不做自由布局编辑器、BIM、3D、物流路径优化或设备级摆放。

设施交付边界：

- 厂房主体/租赁空间已交付并完成必要改造。
- 办公、生产、仓储、质量和公辅区域已划分。
- 基础电、水、气、消防、EHS 和网络条件已验收。
- 生产设备、工装、实验室仪器、具体货架与人员尚未安装/到岗，属于 M11。

## 7. 选址决策模型

指标至少包含：

- 总现金需求、预算承诺和租赁/建设成本。
- 可用日期和建设/改造周期。
- 客户、供应商和物流距离。
- 人力成本与人才可获得性。
- 电、水、气、环保和消防容量。
- 自然灾害、供应商和审批风险。
- 扩展性、控制权和退出成本。

评分器只生成可解释 candidate comparison，不自动替代批准。硬约束（现金、预算、最低公用工程、最晚可用日期）先于加权评分；不满足硬约束的候选不能因高综合分被选中。

## 8. 工期、资金和项目不变量

- 没有 `plant_project_eligible=true` 不能创建可执行项目。
- 未批准场址和投资上限不能签署协议或开工。
- 累计 committed amount 不能超过批准 budget envelope。
- 实际付款不能超过公司可用现金，也不能先于对应审批和验收条件。
- 预算授权、合同承诺、应付、已付款和现金是不同数值。
- WorkPackage 只有在全部 predecessor 完成且所需场地/资源可用时才能开始。
- 时间推进本身不能自动完成工作；活动完成由资源、日历、外部能力和规则计算。
- 里程碑付款必须引用已接受的世界交付证据；返工不重复创造资产。
- 设施 acceptance 需要所有强制区域、公用工程、消防和 EHS 门通过。
- `capability_build_eligible=true` 不代表可生产，只代表可以进入 M11 设备/人员能力建设。

## 9. 最小异常 tracer

M10 固定一条公用工程延期 tracer：

```text
IAOS 项目计划认为增容按期
-> 公用工程服务方实际延期
-> World schedule 与 IAOS plan 产生 discrepancy
-> 项目总监尚未知
-> observation 送达后项目总监获知
-> 提交 rebaseline / mitigation intent
-> IAOS committed outcome 更新计划与审批
-> AESE 计算新的施工顺序、工期和成本结果
-> 验收后关闭 discrepancy
```

该 tracer 必须证明“计划变更”不会直接让现场恢复，也不能通过 UI 手工改完成度。

## 10. 角色与治理

| 角色 | 职责 | 关键限制 |
| --- | --- | --- |
| 工厂项目负责人 | 场址调研、推荐、WBS、异常重排和验收申请 | 不能批准自身越权投资/付款 |
| CEO | 批准选址与项目目标 | 不能绕过预算硬约束 |
| CFO | 审查现金、预算、承诺、变更和付款 | 不能把未验收里程碑当作付款证据 |
| 外部承包商/园区/公用工程方 | 由 World policy 提供报价、施工和外部结果 | 不是 IAOS 内部用户 |

人类与 Agent 共用同一 IAOS Capability、Decision、Policy、Process 和审计；Agent 输出只形成 recommendation/intent。

## 11. Bridge payload family

M10 增加严格 allowlist 类型：

```text
genesis.site.assessment.completed.v1
genesis.site.selection.requested.v1
genesis.site.selection.approved.v1
genesis.site.control.delivered.v1
genesis.facility.project.approved.v1
genesis.work_package.started.v1
genesis.work_package.completed.v1
genesis.utility.connection.delayed.v1
genesis.project.rebaseline.requested.v1
genesis.project.rebaseline.approved.v1
genesis.milestone.accepted.v1
genesis.payment.approved.v1
genesis.facility.accepted.v1
```

所有类型复用 DES-008 envelope、stable ref、tenant journal、cursor、幂等和 committed outcome 语义。

## 12. Pack 与版本策略

- `hctm-genesis` 升级到下一 minor version，新增 `campaigns/plant-build/`。
- Plant Build 初态显式引用并验证 M9 incorporation 终态 hash/eligibility，不复制或手填成立结果。
- M8 设备 tracer 与 M9 incorporation 继续可独立运行和回归。
- M11 只能消费 M10 机器输出，不从画布或项目文案推断设施已就绪。

## 13. 非目标

- 生产设备、工装、实验室仪器采购安装和调试。
- 人员招聘、培训、认证与班次。
- 完整 EAM、采购、合同、工程造价或财务总账产品。
- BIM、CAD、3D、自由布局编辑和高精度施工物理仿真。
- 真实园区、政府、公用工程或承包商接口。

## 14. 完成标准

- 单一 run 从 M9 eligibility 确定性推进到 `capability_build_eligible=true`。
- 至少三个场址候选通过硬约束和可解释评分，资金不可行方案被正确拒绝。
- 空间、设施、项目、现金、预算、承诺、付款和知识三态可对账。
- 公用工程延期、知情延迟、受治理 rebaseline 和最终验收形成完整因果链。
- 人类/Agent 共用治理能力，越权、自批、超预算、未验收付款全部失败关闭。
- M7/M8/M9 回归、两仓测试、部署、runbook/evidence 和 Atlas 完整。
