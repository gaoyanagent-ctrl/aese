---
id: DES-027
title: M9 IAOS 原生语义驱动企业成立真实闭环
date: 2026-07-23
status: approved
author: Codex + User
tags: [m9, iaos, semantic-runtime, incorporation, closed-loop]
---

# M9 IAOS 原生语义驱动企业成立真实闭环

## 1. 背景

DES-010 定义了企业成立的 World/IAOS/Knowledge 所有权与业务链，但当前实现主要由 AESE 确定性 frame 和 IAOS `genesis_governance_record` receipt 构成，没有从 IAOS 三层语义库编译正式 Entity、Capability、Process、权限、菜单和工作台，也没有形成可从 Journal/Outbox/World Store 恢复的真实双向闭环。

本设计重新打开 M9 的实现边界。目标不是新增一组 Genesis 专用 CRUD，而是用 M9 验证 IAOS 能否按照自身系统逻辑满足真实企业设立、治理、资本、账户、任命和预算需求。

## 2. 固定生成链

```text
Core Semantic Library
→ Domain Semantic Library
→ Tenant Extension Library
→ Archetype / Entity
→ Atomic Ability / Business Capability
→ Process / Policy / Decision / Approval
→ Permission Resource / Effective Runtime Artifact
→ API / Form / List / Menu / Agent Tool
→ World Bridge observation / intent / committed_outcome
```

任何 M9 正式业务事实不得只保存在 AESE frame、前端状态或通用 JSON receipt 中。IAOS Runtime、UI、Agent 与 API 必须消费同一份 Effective Runtime Artifact。

## 3. 已确认决策

### D1 — 可复用 Enterprise Governance Domain Package

M9 使用可复用的 `enterprise_governance` Domain Semantic Package，不建立 HCTM 专属领域模型。

三层分工：

- Core：跨行业稳定的 party、organization、person、role、document、event、account、commitment、mandate、money、approval 和 workflow 语义。
- Domain：`incorporation_case`、`legal_entity`、`legal_registration`、`governance_resolution`、`executive_appointment`、`capital_commitment`、`capital_contribution`、`bank_account`、`operating_mandate`、`budget_envelope`、`budget_approval`。
- Tenant Extension：华辰稳定编码、虚构登记机构/银行、金额、岗位、审批矩阵、表单/列表/菜单布局和 M10 资格规则。

该包必须能被后续其他企业场景复用；HCTM extension 不得修改或降级 Core/Domain 契约。

### D2 — Core Semantic 扩展边界

复用现有 Core `party`、`organization`、`person`、`role`、`document`、`event`、`approval_policy`、`workflow_instance` 和 `audit_record`。

Core 只新增三个跨行业 Archetype：

- `account`：主体拥有、具有 currency/status 和余额口径的账户原型。
- `commitment`：主体在条件满足后提供资金、资源或履行义务的承诺原型；不得与 actual、payment、receipt 或 budget authorization 混同。
- `mandate`：角色在 scope、limit 和 validity window 内可以执行或批准动作的授权原型。

`money` 是 `value + currency + scale` 的结构化 Semantic Value，不建立独立 Entity 或 Archetype。`bank_account`、`capital_commitment` 和 `operating_mandate` 分别从上述 Core Archetype 派生；其他 M9 对象留在 `enterprise_governance` Domain Library。

### D3 — Domain Entity 图与平台能力复用

`enterprise_governance` 创建 11 个正式 Domain Entity：

1. `incorporation_case extends document`
2. `legal_entity extends organization`
3. `legal_registration extends document`
4. `governance_resolution extends document`
5. `organization_position extends role`
6. `executive_appointment extends document`
7. `capital_commitment extends commitment`
8. `capital_contribution extends document`
9. `bank_account extends account`
10. `operating_mandate extends mandate`
11. `budget_envelope extends document`

不创建 `budget_approval`、`position_assignment`、`workflow_instance`、`approval_history` 或第二套 `audit_record`。审批、流程、组织运行关系、权限和审计分别复用 IAOS Approval、Process、`sys_org_node`/用户组织关系、Role/Permission 和 Audit Runtime。

关键关系：`capital_contribution fulfills capital_commitment` 且 `credited_to bank_account`；`executive_appointment assigns organization_position` 且 `based_on governance_resolution`；`operating_mandate granted_to organization_position`；`budget_envelope governed_by operating_mandate` 并引用 IAOS `approval_request`。

`legal_entity` 注册成功后由业务能力投影到 `sys_org_node`；任命被 World observation 确认接受后，业务能力更新用户组织关系和权限角色。World 实际余额与 IAOS 已核验余额严格分离。

### D4 — 现有资产与数据审计门

实现前必须对 IAOS 当前 Core/Domain/Tenant semantic assets、Archetype、Entity schema/data、Capability、Process、Policy、Decision、Approval、permission resource、menu、Journal 和 Outbox 做机器可读盘点。每项分类为 `reuse|complete|conflict|migrate|create`，并记录 owner、version、status、引用和租户数据质量；未完成审计不得发布 `enterprise_governance` package。

2026-07-23 当前运行库初审：

| 区域 | 当前事实 | 初步分类 |
| --- | --- | --- |
| Semantic Concept | 58 条，只有 42 field、8 capability、8 event；无 entity/archetype concept | `complete` |
| Semantic Relation | 0 条 | `create` |
| Core Archetype | 仅 document/document_line/document_with_lines/party/supplier；party/supplier 为 draft | `complete` |
| Archetype Relation | 1 条 | `complete` |
| tenant-hctm Entity Schema | 10 个既有制造实体，但 archetype/semantic_entity 均为空，所有字段 semantic_id 均为空 | `migrate` |
| tenant-hctm Capability | 0 | `create` |
| tenant-hctm Process | 0 | `create` |
| tenant-hctm Approval | 0 | `create` |
| tenant-hctm Org Node | 0 | `create` |
| tenant-hctm Role | admin/operator 两个系统角色 | `complete` |
| tenant-hctm Menu | 29 个，其中只有 6 个关联现有 Entity，布局元数据多数为空 | `complete|migrate` |
| dev-user / Semantic Studio | API 对 `dev-user@tenant-hctm` 返回 28 个菜单且包含 `semantic_studio`；前端却用 `tenant_id == tenant-001 && role == admin` 再次过滤；dev-user 只在 tenant-001 有账户行，在其他租户由开发 JWT 和后端 bypass 虚拟化；无正式 platform-admin 数据模型 | `conflict|migrate` |
| M9 Governance Receipt | 1 条 budget.approve | `migrate` |
| M9 World Journal | 0 | `create` |
| M9 Outbox | 1 条且 FAILED | `conflict` |

结论：IAOS 已有 Semantic Studio、Compiler、Entity Runtime、Capability/Process/Approval、Org/RBAC、Menu、Journal 和 Outbox 基础设施，但当前 HCTM 数据没有形成五层语义运行模型，不能作为 M9 完成证据。

### D5 — 独立生命周期租户

M9 从干净、可重放的 `tenant-hctm-genesis` 开始，不在承载 M3–M7 兼容数据的 `tenant-hctm` 上原地改造。`tenant-hctm` 保留为既有制造场景和迁移对照；M9–M24 的公司生命周期事实进入新租户。

新租户不是空数据库复制。它必须通过受版本治理的 package/bootstrap 创建，并在启用业务前通过 D4 审计门：

- 安装 Core、`enterprise_governance` Domain 和 HCTM Tenant Extension 的明确版本；
- 编译并校验 Semantic → Archetype → Entity → Capability → Process/Permission/UI 依赖链；
- 创建正式组织、角色、权限、菜单和开发验收身份投影；
- 禁止依赖 `tenant-001`、`dev-user` 字符串硬编码决定前端可见性或管理权；
- seed/import 默认 dry-run，显式 apply，重复执行必须幂等。

`dev-user` 当前只是在本地开发环境中的特殊 subject：数据库仅在 `tenant-001` 存在账户行，跨租户 token 和多数后端权限由特殊 bypass 提供。这不等同于已实现平台管理员。M9 新建平台级主体 `founder-principal`，显示名称为“创始治理者”；M9 实现必须把其平台管理员资格、租户访问关系和有效权限作为可查询的正式投影，并让后端菜单 API 与前端使用同一授权结果。

### D6 — 平台管理员、菜单与工作室授权

平台管理员在其被授权访问的任一业务租户中，默认可见：

- 该租户全部已安装业务菜单；
- Semantic、Entity、Capability、Process、Governance、Permission、Audit、Atlas 等平台工作室；
- 尚未安装但属于可用平台目录的模块，以禁用状态显示并返回不可用原因，不以静默隐藏冒充权限不足。

菜单投影的唯一依据是服务端计算的 `installed runtime artifact + effective permission + tenant access`。前端不得再根据 `tenant-001`、用户名 `dev-user`、`founder-principal` 或 JWT 中单一角色名二次裁剪菜单或管理能力。页面路由和写操作仍必须由服务端逐次鉴权，菜单可见不替代操作权限。

M9 实现必须补齐正式的平台身份目录、平台角色与租户访问投影，并让 `/profile` 或等价身份端点返回可审计的有效授权摘要。`founder-principal` 是 M9 正式验收主体，必须获得 `platform_super_admin`、对 `tenant-hctm` 与 `tenant-hctm-genesis` 的明确租户访问关系，以及租户内董事长/最高管理者岗位和治理 Mandate。`dev-user` 仅保留为本地兼容 bootstrap subject；特殊 bypass 只用于迁移，不能成为 M9 验收依据。原 `tenant-hctm` 的 Semantic Studio 缺失问题归类为前端授权冲突，修复后需用 `founder-principal` 完成 API、菜单、直达路由和写权限四层回归。

### D7 — M9 正式业务状态机

M9 正常主链固定为：

```text
incorporation_case_opened
→ founder_resolution_approved
→ capital_commitments_confirmed
→ registration_submitted
→ legal_entity_registered
→ bank_account_opening_submitted
→ bank_account_opened
→ capital_contribution_verified
→ organization_established
→ executive_appointments_accepted
→ operating_mandates_activated
→ initial_budget_approved
→ enterprise_operational_ready
```

状态不得由通用 Entity Update、前端本地状态或 AESE frame 直接改写。每次转换必须由已发布的 Business Capability 在 IAOS 事务中同时校验前置条件、写业务事实、追加 Journal，并经 Outbox 发布。

以下步骤必须经过 World Bridge：

- 登记机构审查：IAOS 提交 registration Intent；World 返回受理、补正、拒绝或登记 Observation。
- 银行开户审查：IAOS 提交 account-opening Intent；World 返回开户、补件或拒绝 Observation。
- 出资到账：World 返回实际到账 Observation；IAOS 核验后形成 capital-contribution CommittedOutcome，World 余额与 IAOS 核验余额分离。
- 任命接受：IAOS 发出 appointment Intent；World 返回候选人接受或拒绝 Observation。

World Observation 只证明外部事实已被观察，不能绕过 IAOS Policy、Decision、Approval 或 Capability 自动成为正式业务事实。桥接消费者必须做 schema、租户、correlation、causation、幂等键、当前状态和版本校验。

任一步骤均支持明确的 `correction_required`、`rejected`、`withdrawn` 或 `terminated` 分支；补正后以新命令和原 correlation 链重提，不覆盖历史记录。进入 `enterprise_operational_ready` 前必须证明法律主体、账户、已核验出资、组织岗位、任命、有效授权和初始预算全部存在且相互引用一致。

### D8 — 自然人最高管理者与 AI Agent 组织

`founder-principal` 是平台身份目录中的全局唯一 Person/User principal，显示名称为“创始治理者”。该主体代表用户本人，在 `tenant-hctm-genesis` 中担任企业最高管理者和最终治理责任人，同时通过独立的 `platform_super_admin` 角色获得平台管理资格。平台角色与租户内董事长岗位分别授权、分别审计，不因账号代码自动获得权限。该身份必须具有可审计的 Person、User、Platform Role Assignment、Tenant Access Assignment、Organization Position 和业务治理 Mandate 投影。

M9 的日常工作由不同 AI Agent 承担。Agent 是 IAOS 中的正式行动主体，必须具有稳定 `agent_id`、服务身份、组织归属、岗位、角色、Mandate、允许调用的 Atomic Ability/Business Capability、金额或数据范围、有效期和人工升级规则。Agent 不得使用 `dev-user` token，不得通过前端模拟点击绕过 Agent Tool、Capability、Process、Policy、Approval、事务、Journal 或 Outbox。

身份与职责采用双层模型：

- Governance Seat：董事长/最高管理者、治理秘书、财务负责人、法务负责人、审计负责人等企业岗位。
- Acting Principal：自然人用户或 AI Agent 对某一岗位的正式任命；同一岗位可有主责 Agent 和人工监督者，但每次动作只能有一个明确 actor。

所有写操作记录 `actor_type`、`actor_id`、`acting_position`、`mandate_id`、`capability_version`、`correlation_id` 和 `idempotency_key`。AI 建议、草稿、正式提交、审批、外部 Observation 和最终 CommittedOutcome 必须可区分。

`founder-principal` 可基于有效授权查看平台与租户菜单、配置并任免 Agent，也拥有最终治理决定权；但对职责分离事项的介入必须走“显式人工接管/治理特批”，记录理由、影响和审计事件，不能依赖隐藏 bypass 静默改写业务事实。

首批 Agent 固定为：

| Agent | 组织岗位 | 主责 | 禁止事项 |
| --- | --- | --- | --- |
| `incorporation-agent` | 企业设立专员 | 开设设立案、协调材料、提交登记、处理补正 | 不得批准创始人决议或自行确认登记成功 |
| `governance-agent` | 公司治理秘书 | 决议草案、岗位方案、任命和 Mandate 建议 | 不得批准自己的建议或代替候选人接受任命 |
| `legal-compliance-agent` | 法务合规负责人 | 登记/章程/授权合法性与职责冲突检查 | 不得伪造外部机构 Observation |
| `finance-agent` | 财务负责人 | 出资承诺、到账核验、银行账户资料、初始预算 | 不得批准自己编制的预算或自行制造到账事实 |
| `audit-agent` | 内部审计负责人 | 独立核对越权、异常、证据、Journal/Outbox 完整性 | 只读，不得代办业务或修改被审计记录 |

`founder-principal` 担任董事长/最高管理者并保留最终治理权：任免 Agent，最终批准创始人决议、核心高管任命、重大 Mandate 和初始预算，以及执行有理由、有审计的人工接管或特批。

登记机构、银行、候选人和资金结算环境属于 World 外部参与者，不创建为 IAOS Agent。它们只能通过受治理 World Bridge 返回 Observation。

### D9 — Agent 自主边界与七个人工治理门

Agent 可自主执行资料收集、草稿生成、规则检查、状态查询、提醒、证据整理和既有授权范围内的补正，不得自主跨越以下治理门：

| Gate | 受治理决定 | Agent 职责 | 最终决定 |
| --- | --- | --- | --- |
| G1 | 创始人决议 | governance 起草，legal 校验 | `founder-principal` 批准 |
| G2 | 首次登记提交 | incorporation 组包，legal 校验 | `founder-principal` 授权提交 |
| G3 | 银行开户提交 | finance 组包，legal 双检 | `founder-principal` 授权提交 |
| G4 | 出资确认 | World 到账 Observation，finance 核验，audit 复核 | 正常自动确认；异常升级 `founder-principal` |
| G5 | 核心高管任命 | governance 提议，World 候选人接受 | `founder-principal` 批准 |
| G6 | 经营 Mandate 激活 | governance 生成，legal 检查 | `founder-principal` 批准 |
| G7 | 初始预算 | finance 编制，audit 独立检查 | `founder-principal` 批准 |

G2 首次正式提交获批后，`incorporation-agent` 可在原申请范围、原 correlation 链和规定次数/期限内自动处理普通补正并重提；扩大注册范围、改变资本结构或触发 Policy 风险时必须重新申请人工批准。

每个 Gate 由 IAOS Process Runtime 持有等待状态，由 Approval Runtime 持有请求、决定、意见和签署身份。Agent 只能通过 Agent Tool 调用已授权 Capability 创建草稿、提交审批或执行已批准动作。前端按钮同样调用 Capability，不得直接更新状态字段。超时、拒绝、授权撤销、Agent 暂停或版本变化均使尚未执行的批准失效。

### D10 — 企业生命周期导航与双向深链

`tenant-hctm-genesis` 增加“企业生命周期”一级导航，M9 业务区包含：

1. 企业成立总览：状态机、当前阶段、阻塞、下一动作、Agent 与 World 交换。
2. 设立案件：设立案、材料、登记提交和补正。
3. 法律主体：注册信息、证照和主体状态。
4. 公司治理：创始人决议、组织岗位和高管任命。
5. 资本与账户：出资承诺、实际到账、银行账户和差异。
6. 授权与预算：Mandate、额度、有效期、初始预算和审批。
7. Agent 组织：Agent 岗位、状态、授权范围、任务和工具调用。
8. 我的审批：`founder-principal` 的 G1–G7 待办、已办和特批。
9. World 交换：Intent、Observation、CommittedOutcome、correlation 和重试状态。
10. 成立审计：Journal、Outbox、Approval、Agent Tool Call 和证据完整性。

Semantic、Entity、Capability、Process、Permission、Audit 和 Atlas 保持为平台工作室，不混入业务菜单。业务详情必须提供来源深链，可跳转到对应 Semantic Concept/Relation、Archetype、Entity Schema、Capability、Process Version、Approval、Permission Resource 和 Atlas node。

AESE World 首页显示 M9 同一阶段图和交换状态；从阶段、工作步骤或交换记录可携带 tenant、case、process run、correlation 参数深链进入 IAOS。IAOS 业务页也可反向打开关联 World run/observation。两端不得各自维护一套阶段定义，均消费版本化 lifecycle/process projection。

### D11 — 可重放主线与异常场景

M9 验收包至少包含以下确定性场景：

| 场景 | World/Agent 输入 | 必须证明 |
| --- | --- | --- |
| 正常成立 | 登记、开户、到账、任命均成功 | 最终进入 `enterprise_operational_ready`，全部引用和证据完整 |
| 登记补正 | 登记机构返回材料缺失 | incorporation Agent 在 G2 原授权边界内补正重提，保留原 correlation 链 |
| 开户拒绝 | 银行返回受益所有人资料不完整 | 原 G3 授权失效，修订后必须重新审批 |
| 出资差异 | 实际到账与承诺不一致 | 禁止自动确认，生成 Discrepancy 并升级 `dev-user` |
| Agent 越权 | finance Agent 尝试批准自己编制的预算 | 权限/职责分离拒绝，业务状态不变，审计与 Tool Call 记录完整 |

场景使用虚构主体和机构、稳定业务编码、固定随机种子、RFC 3339 时间及 `Asia/Shanghai` 业务时区；金额使用明确 currency/scale，不依赖浮点舍入。每个场景可独立 dry-run、apply、reset、replay 和 verify。

正常与异常场景共享同一 Semantic、Archetype、Entity、Capability、Process、Policy、Approval 和 UI runtime artifact，只允许操作输入和 World Observation 不同。不得为演示分支增加专用 API、直接数据库写入、direct NATS 或前端伪状态。

### D12 — 三层资产安装与现有租户保护

三层安装边界固定如下：

- Core Semantic Library 安装于平台层，只增加已确认的 `account`、`commitment`、`mandate`、结构化 money value 和必要关系；已发布资产不可原地变更。
- `enterprise_governance` 作为带版本、依赖和内容 hash 的 Domain Package 发布到平台目录。
- HCTM Tenant Extension 仅安装到 `tenant-hctm-genesis`，承载企业稳定编码、虚构机构、表单/列表/详情、菜单、审批矩阵、Agent Mandate 和场景数据。

安装前必须输出机器可读审计报告，逐项分类 `reuse|complete|conflict|migrate|create`。`conflict` 阻断发布；`complete` 项先补齐并重新编译；工具不得自动覆盖冲突资产。安装默认 dry-run，apply 必须显式确认，重复 apply 幂等，并生成依赖锁、前后 diff、内容 hash、回滚清单和验证结果。

`tenant-hctm` 暂不迁移或清理，继续承载 M3–M7 兼容数据并作为迁移对照。任何旧租户迁移必须另立计划和验收，不得混入 M9 bootstrap。

`tenant-hctm-genesis` 只有在 Semantic 引用、Archetype 继承、Entity relation、Capability 依赖、Process/Approval、Permission/Menu、Agent Mandate、场景引用和 reset/replay 全部通过验证后方可启用。启用本身是可审计的治理动作。

### D13 — World Bridge 事务、投递与恢复

真实闭环固定为：

```text
IAOS Business Capability
→ 同一事务写业务事实 + Journal + Outbox
→ Bridge Worker 投递 Intent
→ AESE simulation ingress 校验并写 World Store
→ World 产生 Observation
→ IAOS simulation ingress 接收
→ 同一事务校验状态并写业务事实 + Journal + Outbox
→ CommittedOutcome 返回 AESE
→ 双方 reconciliation 校验收敛
```

IAOS 和 AESE 均不得直接写对方数据库；不增加 direct NATS 或 webhook-only 正式路径。NATS 只能分发事务 Outbox 已持久化的事件，不能成为事实来源。

所有交换包含 `tenant_id`、`scenario_id`、`world_run_id`、`schema_version`、`correlation_id`、`causation_id`、`idempotency_key`、`occurred_at` 和 `canonical_hash`。系统采用至少一次投递、效果恰好一次；重复输入返回原 committed result，不产生第二次业务变化。

未知 correlation、非法状态转换、乱序、过期版本、hash 不符或租户不符的消息进入隔离队列并产生 Discrepancy，不推进业务状态。双方 reconciliation API 必须能报告 missing、duplicate、lagging、hash_mismatch 和 terminal_conflict，并可从持久化 Journal/World Store 重放恢复。

UI 只展示持久化 projection。IAOS、AESE、Bridge、浏览器任一方重启后，流程、交换和 Agent 任务均可恢复；不得依赖内存或浏览器 local state 作为完成证据。

### D14 — M9 Business Capability 目录

M9 发布以下 Business Capability：

1. `incorporation.case.open`
2. `founder.resolution.prepare`
3. `founder.resolution.approve`
4. `capital.commitment.record`
5. `registration.package.validate`
6. `registration.submit`
7. `registration.correction.resubmit`
8. `registration.observation.commit`
9. `bank.account.opening.submit`
10. `bank.account.observation.commit`
11. `capital.contribution.verify`
12. `organization.establish`
13. `executive.appointment.propose`
14. `executive.appointment.acceptance.commit`
15. `executive.appointment.approve`
16. `operating.mandate.grant`
17. `initial.budget.prepare`
18. `initial.budget.approve`
19. `enterprise.readiness.evaluate`
20. `incorporation.exception.escalate`

Business Capability 由版本化 Atomic Ability 组合，至少覆盖实体读写、语义/关系校验、草稿、Approval 请求/决定、Intent 发出、Observation 接收、金额核验、Journal、Outbox 和异常升级。事务性业务变化不能拆成由调用方自行拼接的多次原子写入。

API、表单动作、Agent Tool 和流程服务任务必须引用同一个已发布 Capability version。输入/输出 schema、前置/后置条件、Policy、所需 Mandate、职责冲突、幂等作用域和事件契约由 Effective Runtime Artifact 统一生成或解析；禁止在不同入口重复实现业务规则。

### D15 — 一主四子 Process 结构

M9 发布一个主流程和四个子流程：

- `enterprise.incorporation.lifecycle.v1`：持有全局状态机、G1–G7、超时、终止和 readiness。
- `legal.registration.v1`：材料校验、登记提交、World 审查、补正、拒绝和登记成功。
- `banking.and.capitalization.v1`：开户、银行 Observation、出资到账、差异和核验。
- `organization.and.appointments.v1`：组织投影、岗位、高管提议、候选人接受和最终任命。
- `mandate.and.initial.budget.v1`：经营授权、预算编制、独立审计和批准。

主流程只消费子流程持久化 result，不轮询浏览器或 AESE UI。子流程实例引用同一 incorporation case 和 correlation root，可独立暂停、重试、恢复和补正；补正不得重建主流程或覆盖历史。

终局登记拒绝可终止尚未成立的主流程。开户、到账、任命或预算失败进入可恢复阻塞。法律主体已经登记后，系统不得通过删除或普通 rollback 撤销法律事实，只能追加正式撤销、变更、补偿或清算记录，并由相应治理流程处理。

### D16 — Policy、Decision、Approval 最小治理集

M9 发布并版本锁定以下规则：

1. `incorporation.document.completeness`
2. `incorporation.separation.of.duties`
3. `world.observation.trust`
4. `capital.contribution.match`
5. `appointment.eligibility`
6. `mandate.scope.and.limit`
7. `initial.budget.control`
8. `enterprise.operational.readiness`

G1–G7 使用 IAOS Approval Runtime；规则由 Policy/Decision Runtime 执行。Decision result 必须包含 decision、命中条款、输入 artifact/version、事实快照 hash、证据引用和 evaluated_at，不能只返回布尔值。

Policy deny 不得通过通用 Entity Update、管理员直接改数据或 Agent 重试绕过。`founder-principal` 的最终治理权通过正式 override Capability 行使，要求明确目标规则、业务理由、影响范围、有效期和补偿/复查要求，并产生 Approval、Decision Audit、Journal 和 Outbox 记录。不可覆盖原拒绝结果。

### D17 — Trace Spine 与 IAOS 查询入口

所有 M9 正式记录均携带或可确定性解析以下追踪键：

- `tenant_id`
- `incorporation_case_code`
- `legal_entity_code`
- `process_run_id`
- `world_run_id`
- `correlation_id`
- `agent_id`
- `capability_key + capability_version`
- `semantic_artifact_id + semantic_artifact_version`

IAOS 提供四类统一查询：企业成立案时间线、全局追踪搜索、对象详情“来源与影响”、可校验 evidence bundle。业务人员不需要跨工作室人工猜测关联。

最小只读接口：

```text
GET /api/v1/incorporations/:case_code/trace
GET /api/v1/incorporations/:case_code/evidence
GET /api/v1/world-exchanges?correlation_id=...
GET /api/v1/runtime-artifacts/:type/:key/lineage
```

trace 返回业务事实、Process、Approval、Decision、Agent Tool Call、Journal、Outbox、Intent、Observation、CommittedOutcome 和 Discrepancy 的有序投影。evidence bundle 包含 schema/version、canonical hash、引用清单和验证结果。

AESE 与 IAOS 深链使用同一组稳定键；任何一端都不得以数据库 UUID 作为场景包的唯一业务引用。

### D18 — M9 IAOS 原生真实闭环最终验收门

DES-027 的实现只有在以下十二项全部满足后才能判定完成：

1. 三层语义资产可追溯：Core、`enterprise_governance` Domain Package、HCTM Tenant Extension 均有版本、依赖锁和内容 hash。
2. `founder-principal` 是正式平台主体，可查询平台角色、租户访问、董事长岗位和治理 Mandate；验收不依赖 `dev-user` bypass。
3. 二十个 Business Capability、一主四子 Process、八项 Policy 和 G1–G7 Approval 全部由 Effective Runtime Artifact 驱动。
4. 正常成立主线真实进入 `enterprise_operational_ready`，不得只修改 AESE frame 或通用 JSON receipt。
5. 登记补正、开户拒绝、出资差异、Agent 越权四条异常线全部可重放，并产生正确阻塞、升级和审计。
6. IAOS 与 AESE 完成 Intent → Observation → CommittedOutcome 双向闭环；双方重启、重复投递和乱序输入后仍可恢复、幂等并对账收敛。
7. IAOS 业务工作台、平台工作室和 AESE World 使用同一状态投影并支持双向深链。
8. API、前端按钮和 Agent Tool 调用同一 Capability version；不存在专用演示 API、直接数据库写入或 direct NATS 正式路径。
9. 租户隔离、职责分离、金额精度、失败无部分写入和 Journal/Outbox 原子性通过真实 PostgreSQL 集成测试。
10. dry-run、apply、重复 apply、reset、replay、verify 全链可复现；已经成立的法律事实不能依靠删除回滚。
11. evidence bundle 可从设立案追溯到语义资产、Process、Approval、Decision、Agent Tool Call、Journal、Outbox 和 World 交换。
12. 三视口 UI 验收完成，并由 `founder-principal` 实际完成菜单、直达路由、审批和写权限验证。

完成证据必须同时来自 AESE 与 IAOS 两个仓库。设计完成、代码完成、集成完成和业务验收分别记录，不得互相替代。

### D19 — 可发现、可操作与系统数据验收

M9 资产不能只存在于专用 artifact JSON、测试 fixture 或设立案 trace 中。安装完成后，
用户必须能从 IAOS 通用工作室发现并检查同一份有效资产：

- Semantic Studio：Core `account`、`commitment`、`mandate`、`money_value`，
  十一个领域概念及其关系图。Core Archetype 不能只有名称：account、commitment、
  mandate 必须分别提供稳定编码、状态、主体、金额/币种、范围/有效期等可继承默认字段，
  且每个 `semantic_id` 必须在 Semantic Registry 中存在。
- Entity Explorer：十一项 Entity 的 schema、列表定义、物理投影和租户隔离数据；
  已存在的 native 设立案必须投影为可查询记录。
- Capability Studio：二十项已发布 Business Capability，版本和
  Runtime Artifact hash 可追溯。
- Process Studio：一主四子五个已发布流程；每个 capability node 必须引用真实
  Capability key，不接受以 process key 伪装的占位节点。
- Governance Studio：八个 active Policy Profile 与至少一条对应的 enabled
  Policy Rule；只有 profile 没有 rule 不算实现。Profile 与 Rule 必须有独立明细入口：
  Profile 展示业务意图、失败策略、配置版本和 Runtime 来源，Rule 展示判定动作、适用对象、
  阶段、优先级和启用状态；不得要求用户阅读原始 JSON 才能理解含义。查看与编辑权限分离，
  `founder-principal` 的正式平台角色必须能通过原生 API 持久编辑，刷新后不得丢失。
- Agent 组织：五个 service-only Agent 的岗位、Mandate、允许能力、有效期、
  金额上限和升级对象可查询；人类不能冒充 Agent 调用。

`founder-principal` 登录后的“企业成立与治理”必须提供十个业务工作区：总览、设立案件、
法律主体、公司治理、资本与账户、授权与预算、Agent 组织、我的审批、World 交换和成立审计。
工作区展示通用注册中心和 native runtime 的真实数据，并提供到上述 Studio 的直达入口。

完成度采用机器可检查门槛：11 Entity、20 Capability、5 Process、8 Policy Profile、
8 Policy Rule、5 Agent、10 条领域 relation、10 个业务工作区；重复安装必须
`no_op=true,writes=0`。任一数量、通用 API、records API、菜单授权或工作区缺失时，
路线图只能标记为 active remediation，不能标记 completed。

Semantic Studio 的原型列表、详情、History 和 Artifact 必须使用同一原型 code 与
snapshot identity。快速切换原型不得发出“新 code + 旧 snapshot UUID”的请求；
所有正式 Entity schema 的语义分析必须为零 `semantic_concept_missing`。

Core Archetype 是 Entity effective schema 的默认字段来源，不是仅供展示的旁路目录。
编译器必须按 `Core Archetype defaults → Domain Entity fields → Tenant extension`
确定性合并，字段名去重，并把有效 schema 同时写入 Entity metadata 和物理投影。
原型默认字段新增或更名时必须执行幂等、非破坏性投影迁移，并在可安全转换时从
`payload/extension_data` 回填既有记录；不得删除或重建已成立事实。Semantic Studio
显示的默认字段、Entity Explorer 显示的有效字段、records API 查询列和数据库物理列
必须一致。只更新原型目录而不重新编译 Entity，或只更新 Entity metadata 而不迁移
物理列，均视为未完成。

Runtime Artifact 安装必须在任何写入前执行语义发布门，至少覆盖重复字段、Semantic
Concept 存在性与类型兼容、enum options 完整性、Archetype 字段类型和
`system_managed/overridable` 治理继承。任一 error 或治理继承 warning 均不得激活
Runtime；验收还必须对已安装的十一项 Entity 调用平台真实 Semantic Analyzer，并要求
`errors=0,warnings=0`。不能把用户进入数据模型工坊后手工点击“语义分析”作为首个发现点。

Entity Explorer 的默认选择必须来自当前租户实际返回的 schema 目录。目录加载前保持
空选择；目录返回后保留仍有效的用户选择，否则选择第一项。不得硬编码 `sales_order`
或任何示例实体，也不得为不存在的默认实体发出 schema/ui/records 请求。

## 4. 待确认设计树

以下设计项已经确认，实施计划必须逐项映射：

- Core 新增与复用边界。
- Domain Entity 与 relation graph。
- M9 状态机和逐步业务能力。
- World 外部机构模拟与 observation 语义。
- IAOS Process/Policy/Decision/Approval 和职责分离。
- Tenant seed、菜单、表单、列表、详情与工作台。
- Journal/Outbox/World Store 原子性、恢复、幂等与补偿。
- 查询、Agent Tool、验收数据和最终迁移策略。

DES-027 已获批准。开始实现前必须创建唯一 active 主实施计划，并为 IAOS 修改建立独立 branch/worktree。
