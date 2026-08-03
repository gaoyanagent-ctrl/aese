---
id: PLAN-M10-INTERACTIVE-001
title: M10 Agent 候选与参数化工厂建设交互闭环实施计划
date: 2026-08-01
status: active
author: Codex + User
tags: [m10, agent, plant-build, interactive, governance]
---

# M10 Agent 候选与参数化工厂建设交互闭环实施计划

## 1. 目标与边界

把既有 `hctm-genesis@0.3.0` 固定十帧 reference replay 升级为真实交互闭环：`plant-planning-agent` 先基于 IAOS 权威资金边界生成设施需求草案，人员选择并少量调整后再生成场址候选；项目负责人选择调研项，World 生成外部事实，人员与审批流决定选址、投资、项目、WBS、变更、付款和验收。旧回放迁入 `fixture_only`，不再自动成为业务事实。

AESE 拥有 Agent 候选生成适配、World 调研/施工事实和游戏交互；IAOS 拥有设施需求、proposal/review、Capability、Process、Approval、Agent Run、资金事实、审计和 Outbox。浏览器写操作统一经过 AESE Command Gateway。

## 2. 纵向切片

### S0 计划与差距审计

- [x] S0.1 修订 DES-011，冻结 Agent/人工/World/IAOS 所有权与金额边界。
- [x] S0.2 识别固定候选、金额、赢家、WBS、空间和异常所在代码、fixture、API 与 UI。
- [x] S0.3 建立本计划、路线图、文档索引、Code Map 和 Atlas 入口。

### S1 参数化需求与 Agent proposal 合同

- [x] S1.1 定义 `FacilityRequirement`、`SiteOptionProposal`、`ProposalSet`、`ProposalReview` 严格 JSON 合同。
- [x] S1.2 金额使用 decimal string + currency；日期 RFC3339；候选数量、权重和枚举有边界校验。
- [x] S1.3 接入 MiniMax JSON provider 和 status/proposal API；未配置模型时失败关闭，不返回伪 Agent 候选。人工表单在 S3 接入。
- [x] S1.4 记录模型、prompt 版本、输入/输出 hash、request ID、token、校验和错误证据；AESE CreativeJob 保存生成技术证据，业务 proposal 后续由 IAOS 保存。
- [x] S1.5 单元测试覆盖合法、金额非法、重复候选、无来源、模型坏 JSON、重试和未配置模型。

### S2 IAOS 权威对象与受治理能力

- [x] S2.1 在 IAOS 独立 worktree 发布 Facility Requirement、Proposal/Proposal Line、Review、Investigation Request/Observation、Work Item、Recommendation 与 Decision 九个 Entity/storage contract；八张 FORCE RLS 权威/工作项表仍由逐对象 Capability/Process 唯一写入，同事务生成只读 `entity_projection_*`，数据模型工坊和左侧菜单不暴露通用 CRUD。
- [x] S2.2 发布 `facility.requirement.define`、`site.proposal.record`、`site.proposal.review`、`site.investigation.request`、`site.investigation.observation.commit` Capability Artifact 与实现绑定。五项均进入 `genesis-m9@1.9.0` 平台包并使用 Published Active Artifact/native binding。
- [x] S2.3 Process Artifact 展开为 human task → agent task → human review → world wait → recommendation → approval → formalization；`facility.plant.planning.v1` 由七个受治理 Capability 和审批决定在原事务推进同一 `process_run`。该流程禁止通用一键运行，不能重复创建审批或伪造 World Observation；IAOS M10 工作台提供按案件穿透页和流程工作室深链。
- [x] S2.4 proposal/review、Agent Run、audit、生命周期 trace 和 Outbox 同事务；RLS、幂等、版本冲突和失败无部分写入。AESE CreativeJob 继续保存模型调用技术日志；Agent Proposal 成功时 IAOS 要求匹配的 completed/valid Agent Run，并与 ProposalSet、Capability Execution、Process trace、Audit、Outbox、投影原子提交。模型调用不是 World Observation，因此不写 World Journal。

### S3 用户表单与人工选择

- [x] S3.0 M9 `enterprise_operational_ready` 终态显示 M9→M10 交接卡与主按钮，携带当前 tenant/case/workspace；首页提供可发现的 `企业生命周期 · M9–M24` 总览入口。
- [x] S3.1 Plant Build Play 增加需求表单、字段解释、来源标识和金额/日期/候选边界校验；实际现金与已批预算只读显示 IAOS 来源引用。
- [x] S3.2 候选卡展示依据、金额/工期区间、假设、未知事实、风险、来源、置信度和 Agent 生成证据。
- [x] S3.3 支持采纳调研、退回重生成、人工新增和淘汰；所有决定要求理由。人工新增使用完整业务表单，经 BFF 读取最新 Requirement/ProposalSet 后调用唯一写 owner `site.proposal.record`；没有 Agent 候选时可建立第 1 版人工候选，已有版本时只允许追加一项，并固定旧候选、认证人员和输入/输出 hash。
- [x] S3.4 页面提供可见“功能说明”，解释 M9 资格/资金 → Requirement → Agent Proposal → Human Review → IAOS/World 后续链；真正的 Entity/Capability/Process/Evidence 深链随 S2.3/S4 补齐。
- [x] S3.5 M10 复用 Enterprise Genesis 全屏游戏壳，新增总部规划中心、产业园地图、候选现场和治理会议室四个内部/外部机构地点；玩家、NPC、当前事件和旅程状态首屏可见，正式玩家路由不显示 reference replay。
- [x] S3.6 建立只读 M10 Game Projection 阶段编译器，以 Requirement、Proposal、Review、Investigation、Observation、Recommendation 和 Decision 决定场景，不建立前端业务状态机。
- [x] S3.7 受治理输入只在点击当前 NPC 后以临时经营任务对话框出现，游戏主画面不常驻表单；地点切换和动画不提交业务事实。
- [x] S3.8 场景地点、主任务和候选标记支持键盘与文本替代；触摸目标、焦点和 reduced-motion 满足可访问性合同。
- [x] S3.8a Requirement、Proposal/Review、Investigation/Observation、Recommendation/Approval、Decision 按内部/外部地点进入只读场景档案，输入表单不兼任历史查询；治理会议室读取 IAOS 冻结 Subject/Flow/Assignment，在游戏内由当前受派人批准或驳回，IAOS 工作台只作审计穿透。
- [ ] S3.9 在 S4.3–S4.5 实现权威项目、施工、付款和验收对象后，增加项目办公室、施工现场和验收移交场景。

### S4 World 调研、参数化项目和资金治理

- [x] S4.1 World 外部参与者形成正式报价、权属、面积、电力容量、许可、可用日期和证据引用 Observation；AESE 先写 World Journal，再由 IAOS 受治理能力匹配 Intent/correlation 并完成持久工作项。2026-08-02 已移除玩家外部事实表单，改为服务端重读权威请求并确定性生成。
- [x] S4.2 评分只消费已送达且同时匹配当前 ProposalSet ID、revision 和 proposal ID 的事实；Agent 估算与外部事实并列显示，不得互相覆盖。已交付 Observation-only 派生比较、六类 Requirement/Observation/差额逐项对照、四维可调权重、默认来源与公式解释和证据引用；无合格候选时权重禁用，并可从失败摘要直接带入当前 Requirement、创建下一 revision、重新生成候选。旧候选集 Observation 只进入历史档案，不推进阶段或进入推荐。评分策略/结果由 S2.3/选址审批纵切在 IAOS 权威固化。
- [x] S4.2b 发布 `site.selection.recommend`、`site.selection.formalize`、`facility.site.selection.v1` 和 `genesis.site.selection.approval`；IAOS 重算评分并冻结推荐/审批证据，AESE 治理会议室通过只读审批详情与受控 Command Gateway 完成受派人批准/驳回，approved 后由独立 Capability 正式落地。
- [x] S4.2c 正式选址后发布 `site.control.request`、`site.control.observation.commit` 与 `facility.plant.delivery.v1`；项目负责人发起协议/交接请求，玩家只确认接收，园区权利人的协议、交接、占有授权和生效时间由 World 引擎从权威请求确定性生成，延迟/拒绝失败关闭且允许保留历史后重新发起。
- [ ] S4.3 空间、承包策略、WBS、延期缓解和投资变更均采用 Agent 建议 + 人员决定。
- [x] S4.3a 按 DES-038 把首次需求改为 Agent 草案、人员选择和少量调整；完整专业参数只在人工接管中展开。
- [x] S4.3b 把场址调研 Observation 改为最小玩家确认命令，由 World 引擎从 IAOS 权威请求确定性生成并归档全部外部事实。
- [x] S4.3c 发布设施项目/WBS 的 Effective Capability/Process Artifact 后再解锁下一 NPC；`facility.project.baseline.v1` 依次执行 Agent 草案、人员确认、审批和人员激活，形成项目与 WBS 权威投影；AESE 项目办公室只让玩家生成、选择和确认方案，治理会议室完成审批；MiniMax 严格 JSON 候选禁用 M3 自适应 thinking，分别提供一次输出完整性恢复和一次具体治理字段纠正，最终超时按可重试 503 失败关闭。
- [x] S4.3d 发布首个工程合同授予纵切：从 active WBS 生成 RFQ，World 提供可信投标，Agent 给出可解释推荐，玩家确认并在游戏会议室审批，随后显式归档权威合同；全链由 `facility.contract.award.v1`、四项 Capability、Approval、Audit、Outbox 和 Entity 投影承载。
- [x] S4.3e 发布首个施工里程碑纵切：正式合同启动施工包，施工现场 World 确定性生成进度/质量/安全/证据 Observation，项目负责人另行验收；`facility.construction.milestone.v1` 明确验收不等于付款。
- [x] S4.4a IAOS 从已过账银行科目与设立案件已批预算读取只读财务快照，返回来源引用和 snapshot hash；AESE BFF 不接受页面伪造该快照。
- [ ] S4.4b 投资、合同、变更和付款金额可修订但必须重新校验/审批；资金变化使旧 Requirement snapshot 失效并要求修订。
- [ ] S4.5 AP/CIP/付款/验收按财务 DES-033/034 接通，不以治理 JSON 冒充会计事实。

### S5 Fixture 隔离、迁移与全链验收

- [ ] S5.1 将三个旧候选和固定数值迁入显式 `fixtures/reference-replay/`，生产 API 禁止默认加载。
- [ ] S5.2 外部模型与纯人工两条路径达到相同业务终态；模型输出不同不要求相同 hash，已接受 snapshot 重放必须确定。
- [ ] S5.3 覆盖断线、重启、重复点击、并发审阅、模型失败、无有效候选、单一来源例外和人工接管。
- [ ] S5.4 完成两仓测试、构建、部署、Runbook、现场 UI/API/DB 证据、Atlas、提交和推送。

## 3. 发布门

- 未完成 S2，不得把 AESE 内存/文件草稿称为 IAOS 业务记录。
- 未完成 S3，不得把 JSON API 称为用户可操作功能。
- 未完成 S4，不得把 Agent 估算称为报价、现金、预算、合同或付款事实。
- 未完成 S5，不得把 M10 状态从 `Interactive Revision Pending` 改为 Completed。

## 4. 当前执行

S0、S1、S2、S3.0–S3.8a、S4.1、S4.2、S4.2b、S4.2c、S4.3a–S4.3e 和 S4.4a 已完成。当前在线纵切已从正式选址延伸到场址控制、设施项目/WBS、工程合同，以及“施工启动 → World 进度质量事实 → 独立里程碑验收”。该纵切不等于完整 M10：延期/缺陷/变更、工程发票/AP/CIP/付款、最终设施验收、对应异常场景和 S5 全链验收仍未完成。既有未跟踪 M7 验收产物不属于本计划，不修改、不提交。

## 5. 当前接口与事实所有权

| 用户动作 | 浏览器调用 | AESE 职责 | IAOS 权威结果 |
| --- | --- | --- | --- |
| 打开交互规划 | `GET /api/aese/v1/world/plant-build/financial-constraints?case_code=...` | 使用当前 IAOS token/tenant 定向转发，不缓存业务事实 | 从已过账银行科目和已批预算形成只读 snapshot、来源引用与 hash |
| 恢复候选估算 | `GET /api/aese/v1/world/plant-build/proposals?requirement_id=...` | 使用命名 IAOS adapter 读取最新 ProposalSet，不代理任意路径 | 返回 IAOS 已保存的 candidate-only ProposalSet，使刷新后仍可与 Observation 并列比较 |
| 保存需求并生成候选 | `POST /api/aese/v1/world/plant-build/proposals` | 严格校验 Requirement、调用已配置模型、保存 CreativeJob 技术日志并提交匹配的 Agent Run Evidence | `facility.requirement.define` 提交 Requirement；`site.proposal.record` 原子提交 candidate-only ProposalSet、Agent Run、Capability/Process/Audit/Outbox |
| 人工新增候选 | `POST /api/aese/v1/world/plant-build/proposals/manual` | 读取最新 Requirement/ProposalSet、覆盖浏览器身份字段并生成版本证据 | `site.proposal.record` 可在无候选集时建立第 1 版；已有版本时只允许追加一个人员来源候选，既有候选不可修改 |
| 审阅候选 | `POST /api/aese/v1/world/plant-build/reviews` | 从 IAOS profile 解析实际 actor，禁止浏览器伪造 reviewer | `site.proposal.review` 保存 action、理由、revision、reviewer、Audit 和 Outbox |
| 发起外部调研 | `POST /api/aese/v1/world/plant-build/investigations` | 只接受已保存且结论为 `adopt_for_investigation` 的候选，服务端解析 actor | `site.investigation.request` 创建调查请求、`facility.site.investigation.v1` 持久 `waiting_world` 工作项和 World Intent |
| 外部参与者反馈 | `POST /api/aese/v1/world/plant-build/observations` | 玩家只确认接收报告；BFF 重读 Requirement、ProposalSet、Investigation Request，生成可重放 World Observation 并先提交 Journal | `site.investigation.observation.commit` 只消费受信 Journal，保存外部事实并完成对应工作项；浏览器不能提供或覆盖外部事实 |
| 比较已送达事实 | 页面内只读派生计算 | 只读取当前 Requirement 与已完成 Observation；先硬约束、后按可调权重评分并展示证据 | 不产生 IAOS 写入、不创建推荐或审批结果；后续由版本化 Capability/Approval 固化正式决定 |
| 提交场址推荐 | `POST /api/aese/v1/world/plant-build/site-selections` | 采集候选、权重、推荐理由和例外说明；不接收审批人 | `site.selection.recommend` 重算并固化评分，创建路由后的 Approval Request |
| 正式选址 | `POST /api/aese/v1/world/plant-build/site-selections/finalize` | 只提交推荐和审批引用 | `site.selection.formalize` 校验/消费 approved 请求并写 `site_selection_decision` |
| 发起场址控制 | `POST /api/aese/v1/world/plant-build/site-controls` | 服务端解析项目负责人并提交期望交付日、方式和必要证据范围 | `site.control.request` 创建权威请求、交付 Process Run、World wait 与 Intent |
| 玩家确认接收场址 | `POST /api/aese/v1/world/plant-build/site-controls/observations` | 浏览器只提交案件、控制请求和 `accept_delivery`；BFF 重读 IAOS 权威请求，由 World 引擎生成稳定协议/交接/占有授权 Observation 后先写 Journal | `site.control.observation.commit` 只消费受信 Observation，完成或失败关闭交付 Run 并投影 Entity |
| 让 Agent 生成项目方案 | `POST /api/aese/v1/world/plant-build/project-options` | 浏览器只提交案件；BFF 重读 Requirement 与可信 Site Control，生成 2–3 套项目/WBS 方案 | 暂不写业务事实；返回模型、prompt、request 和输入/输出 hash 证据 |
| 选择并提交项目基线 | `POST /api/aese/v1/world/plant-build/facility-projects` | 玩家只选择 Agent 方案；BFF 固化 plan 后依次调用 record 与 submit | `facility.project.plan.record` 写 Agent 草案，`facility.project.baseline.submit` 创建并路由项目审批 |
| 激活项目基线 | `POST /api/aese/v1/world/plant-build/facility-projects/activate` | 只提交案件、plan 和已批准请求 | `facility.project.baseline.activate` 消费 approved 请求，原子写设施项目和 WBS，并形成 Audit、Outbox、Process trace 与只读 Entity 投影 |
| 发布工程采购邀请 | `POST /api/aese/v1/world/plant-build/contract-rfqs` | 玩家只选择 WBS 包和策略；BFF 重读 active project 并推导合同上限/日期 | `contractor.rfq.issue` 创建 RFQ、World Intent 和 `facility.contract.award.v1` wait |
| 收取正式投标 | `POST /api/aese/v1/world/plant-build/contract-bids/confirm` | 玩家只确认；World 从权威 RFQ 生成可重放虚构投标 | Journal Observation 经 `contractor.bid.observation.commit` 核验并保存 |
| Agent 比选并提交 | `POST /api/aese/v1/world/plant-build/contract-recommendations/agent`、`.../contract-recommendations` | Agent 仅比较可信投标；玩家确认建议 | `contractor.award.recommend` 校验 Evidence hash 并路由合同审批 |
| 归档正式合同 | `POST /api/aese/v1/world/plant-build/contracts/award` | 只提交推荐和批准请求 | `contractor.contract.award` 消费批准并写权威合同；不自动开票、应付或付款 |

上述 BFF 不是通用 IAOS 代理，也不拥有权威业务表。Requirement 和 ProposalSet 的 GET 只用于恢复 IAOS 已保存事实，不从 AESE 本地重建。外部模型未配置、财务快照不完整或变更、候选 Schema 不合法、权限不足、重复键输入不同、审阅版本冲突时均失败关闭。
