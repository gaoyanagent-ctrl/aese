---
id: DES-011
title: M10 Project Genesis 工厂选址与设施建设
date: 2026-07-22
status: active
author: Codex + User
tags: [m10, genesis, plant, site-selection, construction, project]
---

# M10 Project Genesis 工厂选址与设施建设

## 1. 目标

M10 从 M9 的 `plant_project_eligible=true` 开始，让华辰苏州制造公司在资金、工期、空间、承包商和公用工程约束下完成候选场址评估、项目审批、场地使用权取得、厂房改造和设施验收，输出 M11 可消费的 `capability_build_eligible=true`。

首版只建设苏州制造基地一期。M10 建成的是可以安装设备、布置仓储/质量区域并组建团队的设施载体，不代表电池冷却板 A 线已经形成生产能力。

## 2. 业务现实约束

M9 终态在每个企业自己的账簿、预算和治理记录中提供运行时约束，而不是由 M10 场景写死：

- 公司实际可用现金及币种，由 IAOS 现金/银行账簿实时读取。
- 已批准预算 envelope、剩余额度和期间，由 IAOS 预算事实实时读取。
- 已生效 CEO、CFO 和工厂项目负责人岗位及 mandate。

用户必须能在受治理表单中填写或修订设施需求、计划投用日期、候选数量、目标区域、投资申请额、预算上限、现金保留下限和评估权重。所有金额采用 decimal string + ISO 4217 币种，保存修订版本、修改人、修改理由和审批状态；UI 不得把演示金额设为不可修改常量。

候选方案不再由剧情固定为三种模式。`plant-planning-agent` 根据上述约束生成可解释的候选方案草稿；项目负责人可以要求补充、删除、重生成或人工新增候选，再选择进入正式调研的方案。Agent 可以建议绿地自建、租赁改造、定制代建或其他合理模式，但这些只是建议，不是预设必选项。

Agent 不能虚构已经发生的报价、权属、公用工程容量或政府许可。候选草稿中的估算、假设和缺失事实必须明确标注；真实报价和场址条件由 AESE World 外部参与者 Observation 或经验证资料补齐。最终选择必须来自版本化多维评估和 IAOS 受治理决策，不能由 Agent 或剧情直接指定赢家。所有演示地点、园区和外部机构保持虚构。

## 3. 纵向业务链

```text
消费 M9 plant_project_eligible
-> 建立产能与设施需求
-> plant-planning-agent 生成候选方案草稿、假设和待调研项
-> 项目负责人审阅、补充并选择需要调研的候选
-> 外部参与者返回报价与场址事实
-> 评估资金、工期、物流、人力、公用工程和风险
-> Agent 形成解释性比较与推荐，项目负责人决定是否提交
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

另需保存以下受治理对象：

- `facility_requirement`：用户可编辑的需求、目标日期、金额边界、权重和版本。
- `site_option_proposal`：Agent/人工候选草稿、假设、置信度、待验证事实和生成运行引用。
- `site_option_review`：项目负责人的采纳、退回、补充、淘汰及理由。
- `agent_run` / `agent_output_evidence`：模型、提示模板版本、输入快照、结构化输出、校验结果和人工决定；密钥不进入业务数据。

优先用 metadata/config package、Process、Policy、Decision 和 Capability。M10 不把 IAOS 扩展成通用建筑项目管理产品。

## 6. 空间与设施范围

M10 最小空间层级是结构约束，节点名称和地点不是固定值：

```text
country / region / city（用户需求）
-> candidate park / site（Agent 建议，外部事实验证）
-> selected site（人员批准）
-> building / floor（Agent 或人工规划）
-> office / production / warehouse / quality / utility zone（Agent 建议，人员确认）
```

`China / Jiangsu / Suzhou`、`HCTM Suzhou site` 和既有 zone 编码仅是 reference fixture。每个正式节点至少携带 stable code、parent、二维坐标/边界、面积、用途、占用状态、容量和验收状态。制造工厂必须满足办公、生产、仓储、质量、公辅等功能门，但具体建筑数量、区域组合、面积和编码由 Agent/人工方案决定。M10 不做自由布局编辑器、BIM、3D、物流路径优化或设备级摆放。

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

当前 S4.2 评分实现是 AESE 中的只读派生视图，不新增第二个业务写入权威。它只消费已由 IAOS
匹配 Intent、写入 World Journal 并完成工作项的 Observation，按以下顺序计算：

1. 先检查正式报价币种与投资申请额、实际面积、实际电力、正式可用日期、权属结论和许可结论。
2. 任一硬约束失败时标记为不合格，`total_score=null`；调整权重不能把不合格方案变成合格方案。
3. 对合格方案计算成本、工期、容量、控制四个 0–100 分的可解释分项，并按用户当前输入的非负权重自动归一化。
4. Agent 估算额和预计日期只与 Observation 并列展示；正式报价、可用日期、面积、电力、权属和许可始终来自 Observation，二者不得覆盖。
5. 每个比较卡必须展示 Observation ID、硬约束结果、分项分数和证据引用。该比较不是推荐、选址批准或投资批准。

页面权重仍只控制预览，但正式推荐纵切已经交付：`site.selection.recommend` 不信任页面分数，
而是在 IAOS 事务中从权威 Requirement、ProposalSet revision 和已提交 Observation 重算
`site-assessment-v1`，冻结权重、结果、Observation IDs、证据、人员理由和输入 hash，再进入
`genesis.site.selection.approval`。推荐、审批和正式落地分离；只有审批状态 approved 时，
`site.selection.formalize` 才能消费该请求并写正式决定。少于两个合格候选必须填写单一来源例外说明。

IAOS 同时把 Requirement、ProposalSet/Proposal Line、Review、Investigation Request/Observation、持久 Work Item、Recommendation 和 Decision 发布为九个 `domain_projection` Entity。它们只用于数据模型工坊、左侧业务菜单和穿透查询；权威写入仍分别属于七个 Capability 或调查 Process。投影与权威变化在同一事务同步，菜单不提供通用新增、修改和删除，避免用户或 Agent 绕过 M10 业务能力直接改事实。

### 7.1 Agent 生成与人工选择合同

`plant-planning-agent` 的输入只能来自当前租户有权读取的事实和本次需求版本：

- 企业、目标区域、设施用途、面积/产能与公用工程需求；
- 用户填写的目标日期、候选数量范围、投资申请额、现金保留下限和评分偏好；
- IAOS 返回的实际现金、已批准预算、岗位 mandate 和既有承诺；
- 已送达 Actor Knowledge 的市场、场址和外部 Observation。

结构化输出至少包含：

```text
proposal_id / option_type / business_rationale
estimated_amount { amount, currency, basis, range }
estimated_schedule { earliest, likely, latest }
assumptions[] / facts_required[] / risks[]
score_preview[] / confidence / source_refs[]
```

Agent 输出先进入 `proposed`，不得直接创建协议、承诺、付款或批准结果。项目负责人必须在工作项中逐项执行 `adopt_for_investigation`、`request_revision`、`add_manual_option` 或 `discard`；只有采纳后的候选才能请求外部调研。若外部模型未启用、输出不符合 Schema 或缺少来源，工作项保持可人工处理，系统只提供同构空白表单。reference fixture 只能由测试人员在显式 replay 模式加载，并显示醒目标识，不能出现在正式候选下拉框或自动成为业务事实。

### 7.2 金额来源与可编辑边界

| 金额 | 来源 | 用户是否可改 | 治理规则 |
| --- | --- | --- | --- |
| 实际可用现金 | IAOS 银行/现金账簿 | 不可直接改 | 只能由受治理收付款/记账能力改变 |
| 预算 envelope | 已批准预算事实 | 不能在 M10 页面直接覆盖 | 可发起预算新增/调整审批 |
| 投资申请额、现金保留下限 | 用户表单或 Agent 建议 | 可修改 | 每次修改形成新版本并重新校验/审批 |
| 候选估算额 | Agent/人工估算 | 可修改 | 必须保留依据、区间、币种和假设 |
| 外部报价/合同金额 | World Observation/供应商报价 | 可在谈判草稿中修订 | 正式承诺必须引用有效报价和审批 |
| 付款申请额 | 已验收里程碑与合同 | 可在余额内调整 | 不得超过验收、合同、预算及现金门 |

候选数量、投资规模、预算、合同、变更和付款均不得依赖代码常量。演示包可以提供预填值，但 UI 必须标注“演示建议”，允许用户在首次提交前修改；提交后按版本化变更和审批处理。

### 7.3 所有方案选择点

| 选择点 | Agent 负责 | 人员负责 | 权威事实来源 |
| --- | --- | --- | --- |
| 场址候选 | 生成多种可解释草案 | 选择调研项、允许人工新增/淘汰 | World 调研 Observation |
| 空间/功能区方案 | 建议建筑和区域组合、面积及容量 | 确认业务适用性和提交审批 | 设施需求与现场测量 |
| 承包/租赁策略 | 比较总包、分包、租赁、代建等选项 | 选择询价与谈判策略 | 外部报价、资质和合同 |
| WBS 与基准计划 | 生成任务、依赖、工期和资源建议 | 调整、确认基线并提交批准 | 已批准范围、合同和资源日历 |
| 延期缓解方案 | 生成重排、临时供给或范围调整建议 | 评估风险、金额并决定提交 | World discrepancy 与最新 Observation |
| 投资/变更方案 | 比较金额、现金和工期影响 | 选择申请额并承担审批责任 | IAOS 现金、预算、承诺和报价 |

任何 Agent 建议都必须留下“为什么、基于什么、还不知道什么、谁选择了什么”的穿透证据。系统规则可以固定合规门和数据结构，不能固定业务人员应选择的方案。

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

M10 必须覆盖至少一条公用工程或外部交付异常 tracer。异常类型、发生时间和影响值由版本化 World policy + seed 决定；下列公用工程延期仅是确定性回归 fixture，不是每个企业都会发生的固定剧情：

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
| 设施规划 Agent | 生成候选、WBS/缓解建议和解释性比较 | 不得伪造外部事实、批准自身建议或提交资金动作 |

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
- pack 中固定候选、固定金额和固定赢家只能位于 `fixtures/reference-replay/`，并标记 `fixture_only=true`；正式运行从用户输入、IAOS 权威事实、Agent proposal 和 World Observation 建立数据。
- 相同 seed、输入快照、Agent artifact 版本和外部 observation 应可重放；真实外部 LLM 的原始自由文本不能直接成为确定性状态，必须经过 Schema 校验、规范化并以已接受 proposal snapshot 进入运行。
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
- 用户可让 Agent 生成可配置数量的候选、人工增删/重生成并选择正式调研项；至少两个有效候选进入比较，或由有权人员记录单一来源例外理由。
- 空间、设施、项目、现金、预算、承诺、付款和知识三态可对账。
- 公用工程延期、知情延迟、受治理 rebaseline 和最终验收形成完整因果链。
- 人类/Agent 共用治理能力，越权、自批、超预算、未验收付款全部失败关闭。
- M7/M8/M9 回归、两仓测试、部署、runbook/evidence 和 Atlas 完整。
- 无任何候选、金额、赢家或异常被运行时代码写死；关闭外部模型时可人工完成同一流程，且 Agent 不可用状态清晰可见。

## 15. 2026-08-01 设计修订与实现状态

原 `hctm-genesis@0.3.0` 的三个候选、1,500 万预算和固定延期仍作为历史确定性 reference replay 证据保留。它们不再代表正式 M10 产品合同。

交互修订当前已形成第一个可操作纵切：

1. Plant Build Play 读取 IAOS 已过账银行科目和已批准预算形成的只读财务快照，显示来源引用和 snapshot hash；用户只能编辑投资申请额、现金保留下限等业务输入，不能覆盖现金或预算事实。
2. 用户填写 `FacilityRequirement` 后，浏览器只调用 AESE 同源 BFF；AESE 严格校验合同、调用已配置的 `plant-planning-agent` provider，并保存模型、prompt、request、token、输入/输出 hash 与校验结果等 CreativeJob 技术证据。
3. AESE 分别以 `facility.requirement.define` 和 `site.proposal.record` 向 IAOS 提交 Requirement 与 candidate-only ProposalSet。IAOS 使用当前 tenant/actor、FORCE RLS、Capability Execution、幂等键、版本检查、Audit 和 Outbox 保存业务事实。
4. 项目负责人可对候选执行采纳调研、退回重生成或淘汰，并填写业务理由；AESE 从 IAOS profile 解析实际人员身份，再以 `site.proposal.review` 保存审阅，浏览器不能指定 `reviewed_by`。
5. 对已保存且审阅结论为“采纳调研”的候选，人员可发起外部调研。IAOS 创建 `facility.site.investigation.v1` 持久 `waiting_world` 工作项，并通过 World Bridge 写入带 correlation/subject 的 Intent。
6. AESE 外部参与者使用结构化表单返回权属、可用面积、电力容量、正式报价、可用日期、许可、证据引用和备注。AESE 必须先写受信 World Journal Observation；IAOS 的 `site.investigation.observation.commit` 只消费与 Intent/correlation 匹配的 Journal 事实，保存权威 Observation 并完成工作项。
7. 人员在可解释比较下选择合格候选并填写推荐理由、替代方案比较；AESE Command Gateway 调用 `site.selection.recommend`，IAOS 重算并创建统一审批。用户从深链进入 IAOS 审批中心决定，批准后返回点击“同步批准并正式选址”，由 `site.selection.formalize` 原子消费审批。
8. 外部模型未配置、财务快照缺失或已变化、Agent 输出不满足 Schema、候选超过投资申请上限、身份/权限不足、revision 冲突、World Observation 不受信或审批不匹配时失败关闭，不用固定候选或页面 JSON 兜底。

Requirement 到正式选址现已统一为 `facility.plant.planning.v1@1.0.0` 单一持久 Effective Process Run。七个业务 Capability 成功时在同一事务追加人、Agent 或 World 节点证据；调查请求进入 `waiting_event`，推荐进入 `waiting_approval`，审批决定恢复到正式选址或失败结束，正式选择后为 `succeeded`。流程工作室只允许查看该业务命令驱动流程，不能一键运行而重复创建审批或伪造外部事实。

以下仍未实现：人工新增候选的权威提交、场地控制、项目/WBS/合同/施工/变更/付款/验收，以及 AP/CIP/总账财务闭环。因此 M10 状态仍只能表述为“Reference Replay Complete; Interactive Revision Pending”，不能声明完整 M10 已完成。

## 16. 用户配置、操作、验证与恢复

### 16.0 IAOS 菜单入口与双系统职责

M10 的 IAOS 用户入口是左侧导航的 `业务智造层 → M10 工厂规划`，不是隐藏 API，
也不要求用户先记住 AESE URL。该工作台提供八个可解释标签页：

- `权威资金约束`：读取当前案件已过账现金和已批准预算，只读展示来源与快照哈希。
- `设施需求`：查看人员提交的版本化 Requirement。
- `Agent 候选方案`：查看 candidate-only ProposalSet、假设、风险和待验证事实。
- `人工评审`：查看采纳调研、退回或淘汰及其理由；评审不等于投资审批。
- `外部调研工作项`：查看调查请求、`facility.site.investigation.v1`、等待能力、World 状态与已返回 Observation。
- `推荐与审批`：查看权威评分输入 hash、推荐候选、审批流/请求状态和正式选址决定。
- `流程运行`：查看同一 Process Run 的状态、当前节点、context 和完整 trace，并深链流程工作室。
- `穿透证据`：解释 IAOS 业务事实、AESE Agent 技术证据和后续 World Observation 的关系。

用户先在工作台选择一个满足 M9 前置条件的设立案件，再点击 `打开 AESE M10 World` 进入
`/#world-plant-build` 完成需求录入、Agent 生成和人工评审。AESE 负责交互式世界与 Agent
技术运行，IAOS 负责资金约束、Requirement、Proposal、Review、权限、审计和 Outbox；
工作台不得在浏览器本地伪造这些权威事实。

AESE 内部还有同一条受治理入口：当 Enterprise Genesis 游戏投影满足
`enterprise_operational_ready`、进度 100% 且没有未完成 M9 工作项时，终态卡显示
`开始 M10 工厂选址与设施规划`。该按钮把当前 tenant、case 和 workspace 放入 M10 深链并
同步会话上下文；M10 返回动作回到同一 M9 企业，而不是 reference replay 或空白默认案件。
AESE 首页同时把原“世界地图”命名为 `企业生命周期 · M9–M24`，作为阶段总览入口，但总览
入口不能绕过上述 M9 机器资格校验创建业务事实。

该菜单随 `genesis-plant-planning` 平台包安装，平台基础包版本从 `1.8.0` 起包含
`menu.genesis_plant_planning`。升级后的既有租户若仍看不到菜单，应重新登录或强制刷新，
使会话重新加载租户菜单投影；不能以手工 URL 代替缺失的菜单授权。

### 16.1 前置条件

- 当前用户已通过 IAOS 登录并进入正确 Genesis Workspace；AESE 只能透传该用户的 token 和 tenant context。
- 对应 M9 案件已经形成法人编码、币种、已批准预算和已过账银行现金；任一事实缺失时不得手工填写替代值。
- `plant-planning-agent` provider 已配置并返回 connected；未启用外部模型时页面明确显示不可生成，用户可建立本地人工草稿，但当前版本不能把该草稿提交为权威候选。

### 16.2 当前可操作步骤

1. 在 IAOS 进入 `业务智造层 → M10 工厂规划`，选择案件并确认权威资金约束；再点击 `打开 AESE M10 World`，确认“可用现金快照”和“已批准预算快照”均为只读，并能看到 `gl:`、`budget:` 来源引用。
2. 填写目标区域、设施用途、最小面积、电力容量、目标日期、2–8 个候选、允许的方案类型、投资申请额、最低现金保留额、偏好和修订原因。
3. 点击“保存需求并让 Agent 生成候选”。成功表示 Requirement 与 ProposalSet 均已进入 IAOS；这不表示投资已获批。
4. 阅读每个候选的业务理由、金额/工期区间、估算依据、假设、待验证事实、风险、来源和置信度。
5. 选择“采纳调研”“退回重生成”或“淘汰”，填写至少 6 个字符的理由并提交。成功表示 Review 已保存；“采纳调研”仍不等于外部事实已确认。
6. 对已采纳候选点击“发起外部调研工作项”。页面显示 `waiting_world` 后，模拟登记/园区/产权等外部角色填写权属、面积、电力、正式报价、可用日期、许可和证据引用并提交；只有 IAOS 显示工作项 `completed` 才代表事实已进入权威闭环。
7. 至少一个 Observation 完成后，页面显示“外部事实比较”。用户可调整成本、工期、容量和控制权重；系统先显示六类硬约束，再显示分项与综合分。Agent 估算与 World 正式事实分栏显示，任何权重都不能抵消硬约束失败。
8. 选择合格候选，填写至少 12 个字符的推荐理由和替代方案比较；少于两个合格候选时再填写至少 20 个字符的单一来源例外说明，提交 IAOS 选址审批。
9. 点击“打开 IAOS 审批中心”查看事项、证据与冻结处理人并作出决定。approved 只表示有权主体同意；返回 AESE 点击“同步批准并正式选址”后，看到正式选择 ID 才表示落地。

### 16.3 关系与证据

```text
M9 已过账现金 + 已批预算
-> FinancialConstraint（IAOS 只读 snapshot）
-> FacilityRequirement（人员输入，IAOS 权威版本）
-> plant-planning-agent / CreativeJob（AESE 技术证据）
-> ProposalSet（IAOS candidate-only 业务事实）
-> ProposalReview（IAOS 人工决定、Audit、Outbox）
-> Investigation Request / persistent World wait / World Observation（IAOS 权威事实）
-> SiteAssessment derived comparison（AESE 只读、Observation-only）
-> site.selection.recommend（IAOS 权威重算 + 版本化推荐）
-> genesis.site.selection.approval（冻结事项、版本和处理人）
-> site.selection.formalize（消费批准并形成正式决定）
-> facility.plant.planning.v1 Process Run（全程参与、等待、审批与 Artifact 证据）
-> [待实现] Project / WBS / Contract / Payment / Acceptance / Finance
```

页面上的 Agent 估算不是报价，Review 不是审批，派生评分不是推荐或批准，reference replay 不是当前企业运行事实。验收时必须分别核对 AESE CreativeJob 技术证据、IAOS Requirement/Proposal/Review/Audit/Outbox 业务证据和 Observation 的 Journal/工作项证据。

### 16.4 失败恢复

- 财务快照变化：重新打开页面读取新 snapshot，以新 revision 和修订理由保存 Requirement，不能复用旧快照继续提交。
- 模型失败或坏 JSON：保留已填写需求，修复 provider 后重试；系统不得生成模拟候选冒充成功。
- 重复点击：相同 object/revision 使用同一幂等键返回既有结果；同键不同输入必须 409 冲突。
- 审阅版本冲突：刷新 ProposalSet/Review 后按最新 revision 重新决定，不能覆盖他人已提交审阅。
- IAOS 不可用或权限不足：页面保持当前表单和候选展示，但不得显示“已保存”；恢复会话或权限后再次提交。
