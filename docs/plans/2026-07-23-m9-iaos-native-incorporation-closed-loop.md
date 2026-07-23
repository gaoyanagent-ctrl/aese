---
id: PLAN-M9-NATIVE-001
title: M9 IAOS 原生语义驱动企业成立真实闭环实施计划
date: 2026-07-23
status: completed
author: Codex + User
tags: [m9, iaos, semantic-runtime, incorporation, world-bridge]
---

# M9 IAOS 原生语义驱动企业成立真实闭环实施计划

## 1. 目标

本计划以 DES-027 为批准基线，重新打开 M9 的实现边界，把现有 AESE frame、通用 `genesis_governance_record` receipt 和失败 Outbox 升级为 IAOS 原生、可审计、可恢复的企业成立闭环：

```text
三层语义资产
→ Effective Runtime Artifact
→ Entity / Capability / Process / Policy / Approval
→ founder-principal 与受治理 Agent
→ IAOS Intent
→ AESE World Observation
→ IAOS CommittedOutcome
→ 双方 reconciliation 与 evidence
```

目标终态为 `enterprise_operational_ready`。本计划不重写已完成的 M9 World 规则，而是迁移其正式 IAOS 事实承载和联动合同；既有 M3–M24 回归必须保持通过。

## 2. 状态、工期与所有权

- 当前唯一 active 主实施计划：`PLAN-M9-NATIVE-001`。
- 预计工期：8–12 周，以真实纵向 tracer 和双仓证据为完成依据，不以日历自动判定完成。
- AESE 仓库拥有：World pack、外部机构策略、simulation ingress、World Store、CommittedOutcome 接收、双向深链、replay/reset/verify 和联合 evidence。
- IAOS 仓库拥有：平台身份、三层语义与 Runtime Artifact、业务 Entity、Capability、Process/Policy/Decision/Approval、权限、Agent Tool、Journal、Outbox、业务工作台和 Trace Spine。
- IAOS 实现必须在基于 `origin/main` 的独立 branch/worktree 中完成；不得从 AESE 提交混入 IAOS 文件。
- 两仓分别提交、测试、记录 revision 和 Atlas 声明；跨仓依赖只用版本化合同和 commit/revision 引用。

## 3. 固定边界

- 正式验收主体为 `founder-principal`，显示名称“创始治理者”；`dev-user` 只作迁移兼容。
- 新生命周期事实安装到 `tenant-hctm-genesis`；`tenant-hctm` 保留为 M3–M7 兼容和迁移对照。
- 首版为单法人、单币种 CNY、单设立案主线、五个 Agent、七个人工治理门。
- 外部登记机构、银行、候选人和资金结算环境属于 AESE World，不创建为 IAOS Agent。
- API、UI 和 Agent Tool 必须执行同一已发布 Capability version。
- 正式路径禁止直接写对方数据库、direct NATS、前端直改状态或演示专用旁路 API。
- 法律事实一旦成立只能追加变更、撤销、补偿或清算记录，不可通过 reset 删除；reset 只允许清理明确标识的测试 run 和尚未生效草稿。

## 4. 前置硬门

- [x] G1 DES-027 D1–D18 已批准。
- [x] G2 `founder-principal`、五 Agent、G1–G7、二十个 Business Capability、一主四子 Process、八项 Policy 和十二项验收门已确认。
- [x] G3 AESE 当前无其他 active 主计划。
- [x] G4 已阅读 AESE 与 IAOS 的 AGENTS、Agent Context 和 Code Map。
- [x] G5 IAOS 独立 worktree、branch、基线 commit 和干净状态已记录。
- [x] G6 现有资产审计报告完成，所有项归类为 `reuse|complete|conflict|migrate|create`，无未决 conflict。
- [x] G7 双仓合同版本、稳定编码、状态机、消息 envelope 和 schema hash 冻结。
- [x] G8 `founder-principal` 身份迁移和 `dev-user` 兼容退出方案通过安全评审。

G5–G8 未关闭前，不得发布业务 Runtime Artifact 或执行正式 M9 apply。

## 5. P0 — 基线审计与双仓隔离

- [x] T1 记录 AESE 与 IAOS `git status`、HEAD、远程差异、运行服务和现有 M9 数据快照。
- [x] T2 从 IAOS `origin/main` 创建专用 branch/worktree，并在 AESE plan 中记录路径、owner 和基线 revision。
- [x] T3 审计 Core Semantic、Archetype、Entity、Capability、Process、Policy、Decision、Approval、Permission、Menu、Journal、Outbox、World Bridge 和 Runtime Artifact。
- [x] T4 审计现有 `genesis_governance_record`、M9 receipt、World Journal、失败 Outbox 和 `tenant-hctm`/`tenant-hctm-genesis` 数据。
- [x] T5 输出机器可读资产审计报告，逐项标记 `reuse|complete|conflict|migrate|create`、owner、依据和阻塞原因。
- [x] T6 冻结稳定编码、money 结构、RFC 3339/`Asia/Shanghai` 时间、状态机、correlation/causation/idempotency 和 canonical hash。
- [x] T7 为 IAOS 与 AESE 分别建立 contract fixtures、破损 fixtures 和 schema compatibility tests。
- [x] T8 更新 capability gap ledger、双仓 code map、Atlas planned dependency 和风险清单。

验收：G5–G7 关闭；没有未分类资产或隐式覆盖；dry-run 对两仓业务数据零写入。

## 6. P1 — 平台身份与正式授权

- [x] T9 设计并迁移平台级 Person/User principal、Platform Role Assignment 和 Tenant Access Assignment 的持久模型。
- [x] T10 幂等创建 `founder-principal` 与显示名称“创始治理者”，授予 `platform_super_admin` 和两个目标租户的显式访问关系。
- [x] T11 在 `tenant-hctm-genesis` 建立董事长/最高管理者岗位、任命和治理 Mandate，平台角色与业务岗位分别审计。
- [x] T12 扩展 `/profile` 或等价接口，返回平台主体、当前租户、租户访问、平台角色、岗位、Mandate 和有效权限摘要。
- [x] T13 将菜单投影统一为 `installed runtime artifact + effective permission + tenant access`，移除前端按租户名、用户名或 token 单角色裁剪。
- [x] T14 逐项替换 M9 路径中的 `userID == "dev-user"` 特判；兼容 token 仅能映射到已登记主体且不得用于正式验收。
- [x] T15 验证 API、菜单、直达路由、写操作四层一致授权，以及无访问关系、撤权、过期 Mandate 和跨租户失败关闭。
- [x] T16 验证 `founder-principal` 的人工接管/特批必须走正式 Capability、Decision、Approval、Journal 和 Outbox。

验收：D18.2 通过；M9 正式链零 `dev-user` bypass；身份和有效授权均可查询、可撤销、可审计。

## 7. P2 — 三层语义包与 Effective Runtime Artifact

- [x] T17 在 Core Semantic Library 复用或新增 `account`、`commitment`、`mandate`、结构化 money value 和必要关系，禁止原地改已发布版本。
- [x] T18 建立版本化 `enterprise_governance` Domain Package，定义 M9 concept、relation、Archetype、Entity 和 invariant graph。
- [x] T19 建立 HCTM Tenant Extension，承载稳定编码、虚构机构、表单、列表、详情、菜单、审批矩阵、Agent Mandate 和 seed。
- [x] T20 定义 package manifest、依赖锁、内容 hash、compiler version、source version、diff、rollback 清单和签名/来源信息。
- [x] T21 扩展编译器，生成 Entity Schema、Capability/Process 引用、Permission Resource、Form/List/Menu、Agent Tool 和 lineage。
- [x] T22 发布前执行引用、继承、relation、状态、权限、职责分离、菜单和 seed 完整性检查；`conflict` 必须阻断。
- [x] T23 实现默认 dry-run、显式 apply、重复 apply no-op、版本升级、stale artifact 阻断和安全回退策略。
- [x] T24 验证 `tenant-hctm-genesis` 与 `tenant-other` 的安装隔离，确认 `tenant-hctm` 未被迁移或清理。

验收：D18.1、D18.3 的资产部分通过；Runtime 只能消费已发布且未 stale 的 Effective Runtime Artifact。

## 8. P3 — 最小成立纵向 tracer

- [x] T25 先实现 `incorporation.case.open`、`founder.resolution.prepare/approve`、`registration.package.validate/submit/observation.commit`。
- [x] T26 发布 `enterprise.incorporation.lifecycle.v1` 和 `legal.registration.v1` 的最小可运行版本，持有 G1、G2 与补正分支。
- [x] T27 实现 `incorporation.document.completeness`、`incorporation.separation.of.duties` 和 `world.observation.trust` Decision 证据。
- [x] T28 让业务事实、Process transition、Approval、Decision Audit、Journal 和 Outbox 在同一 IAOS 事务提交。
- [x] T29 实现 IAOS registration Intent → AESE ingress/World Store → registration Observation → IAOS CommittedOutcome。
- [x] T30 对未知 correlation、错租户、乱序、旧 schema、hash 不符和重复 Observation 失败关闭或返回原结果。
- [x] T31 提供最小 IAOS 页面和 AESE 深链，可观察设立案、G1/G2、World 交换和 trace。
- [x] T32 从 clean tenant 执行正常登记与登记补正两条 tracer，保存数据库、Journal、Outbox、World Store 和 API 证据。

验收：第一个真实纵向切片贯通；UI、API 和 Agent Tool 使用同一 Capability；失败无部分写入。

## 9. P4 — 资本、账户、组织、Mandate 与预算

- [x] T33 实现剩余 Business Capability，补齐开户、出资核验、组织、任命、Mandate、预算、readiness 和异常升级。
- [x] T34 完成 `banking.and.capitalization.v1`、`organization.and.appointments.v1`、`mandate.and.initial.budget.v1`。
- [x] T35 完成 `capital.contribution.match`、`appointment.eligibility`、`mandate.scope.and.limit`、`initial.budget.control` 和 `enterprise.operational.readiness`。
- [x] T36 固定承诺、到账、核验余额、预算授权和实际现金的独立对象、单位、精度、引用和守恒规则。
- [x] T37 实现开户与候选人 Intent/Observation/CommittedOutcome；外部结果只能由 World Bridge 输入。
- [x] T38 实现 G3–G7、批准失效、拒绝、补正、超时、撤权、版本变化和正式 override。
- [x] T39 实现 readiness evaluator，只有法律主体、账户、核验出资、组织岗位、任命、Mandate 和预算引用一致时进入终态。
- [x] T40 从 clean run 执行完整正常主线，保存 `enterprise_operational_ready` 和全链状态 hash。

验收：D18.4 通过；正式业务事实不再只存在于 AESE frame 或通用 receipt。

## 10. P5 — Agent 组织与异常治理

- [x] T41 建立五个 Agent 的稳定主体、服务身份、组织岗位、Mandate、Capability/Tool allowlist、额度、期限和升级规则。
- [x] T42 让 Agent Tool、人工按钮和 API 进入同一 Capability/Process/Approval 路径，并区分建议、草稿、提交、批准和 committed outcome。
- [x] T43 验证登记补正：Agent 可在原 G2、原 correlation、限次限期内补正，扩大范围必须重新审批。
- [x] T44 验证开户拒绝：原 G3 失效，修改受益所有人材料后重新审批。
- [x] T45 验证出资差异：禁止自动确认，生成 Discrepancy 并升级 `founder-principal`。
- [x] T46 验证 Agent 越权：finance Agent 自批预算被职责分离拒绝，业务状态不变且 Tool Call/Decision/Journal 完整。
- [x] T47 验证 Agent 暂停、Mandate 撤销/过期、金额超限、跨租户和工具禁用均在 dispatch 前失败关闭。
- [x] T48 对五 Agent 执行权限矩阵、知识可见范围、幂等、并发和审计回归。

验收：D18.5 通过；Agent 是正式 acting principal，不使用 `founder-principal` 或 `dev-user` token 冒充人工主体。

## 11. P6 — 工作台、Trace Spine 与双向联动

- [x] T49 在 IAOS 增加“企业生命周期”业务导航及 DES-027 D10 的十个工作区。
- [x] T50 让全部菜单、表单、列表、详情、动作和 Agent Tool 从 Runtime Artifact 与有效权限投影。
- [x] T51 实现企业成立时间线、全局追踪搜索、对象“来源与影响”和 Runtime Artifact lineage。
- [x] T52 实现 trace/evidence/world-exchange/lineage 只读 API，统一稳定键并禁止以数据库 UUID 作为场景唯一引用。
- [x] T53 实现 evidence bundle 的 schema/version、canonical hash、引用清单和离线验证器。
- [x] T54 更新 AESE World M9 页面，消费同一 lifecycle/process projection，展示 Intent/Observation/CommittedOutcome/Discrepancy。
- [x] T55 实现 AESE → IAOS 与 IAOS → AESE 双向深链，携带 tenant、case、process run、world run 和 correlation。
- [x] T56 验证 IAOS/AESE/Bridge/浏览器重启与刷新后页面只从持久事实恢复，不依赖 local state 作为完成证据。

验收：D18.7、D18.11 通过；业务人员可从一个设立案追溯全部语义、治理和 World 证据。

## 12. P7 — 恢复、对账与最终验收

- [x] T57 实现 reconciliation API/CLI，报告 missing、duplicate、lagging、hash_mismatch 和 terminal_conflict。
- [x] T58 验证至少一次投递、效果恰好一次：重复、乱序、延迟、断线、双方重启和 poller recovery 后收敛。
- [x] T59 在真实 PostgreSQL 验证租户隔离、FORCE RLS、职责分离、decimal 精度、并发、幂等碰撞和失败无部分写入。
- [x] T60 验证 Journal/Outbox 与业务事实原子提交，rollback 不产生 committed outcome，失败 Outbox 可恢复且不重复业务变化。
- [x] T61 执行 dry-run、apply、重复 apply、reset、replay、verify；证明生效法律事实不被 reset 删除。
- [x] T62 执行 API、IAOS UI 和 AESE UI 验收，由 `founder-principal` 完成菜单、直达路由、G1–G7 审批和授权写操作。
- [x] T63 运行 AESE/IAOS Go、Schema、数据库、前端、Playwright、Markdown links、JSON、Code Map、Atlas 和 M3–M24 回归。
- [x] T64 在 1440×900、1280×720、390×844 三视口完成无阻塞重叠、可复制文本、错误恢复和双向深链验收。
- [x] T65 编写双仓 runbook、evidence、迁移/回滚说明，分别提交并记录 revision、部署版本和未验证风险。
- [x] T66 逐项关闭 DES-027 D18 十二项验收门；只有全部有机器证据时才将计划和实现标记 Completed。

验收：D18.1–D18.12 全部通过；设计、代码、集成、部署和业务验收状态分别有证据。

## 13. 依赖顺序

```text
P0
→ P1
→ P2
→ P3
→ P4
→ P5
→ P6
→ P7
```

- P1 可在 P0 审计报告初稿后开始，但正式身份迁移必须等待 G8。
- P2 必须消费 P0 冻结合同；P3 必须消费已发布的最小 Runtime Artifact。
- P4 只在 P3 的登记 tracer、事务和 Bridge 恢复通过后扩展。
- P5 可在 P4 Capability schema 冻结后开发测试 fixture，但不得提前绕过业务链。
- P6 可先做只读 projection 骨架，正式动作必须等待对应 Capability 稳定。
- P7 串行收口，期间冻结新范围。

## 14. 每个切片的交付纪律

每个 P0–P7 切片都必须：

1. 先补合同测试或失败用例，再实现。
2. 提供默认 dry-run 和明确 apply 边界。
3. 更新受影响设计、计划、roadmap、code map、progress log 和 Atlas 声明。
4. 分别运行 AESE 与 IAOS 范围相称的测试。
5. 分别提交两仓改动，记录对方依赖 revision。
6. 报告已验证、未验证和剩余风险。

## 15. 非目标

- 不实现真实工商、银行、支付、税务或法定报送集成。
- 不实现完整总账、复杂股权、多币种、融资、清算或集团合并。
- 不迁移或清理 `tenant-hctm` 的 M3–M7 兼容数据。
- 不重写 M10–M24 业务世界；仅保证回归和未来消费新 M9 terminal 的迁移入口。
- 不以真实客户、个人、证件、银行账号或生产凭据作为 seed。
- 不以 UI 展示、文档完成或提交数量替代真实闭环证据。

## 16. 完成定义

- T1–T66 全部完成，G5–G8 全部关闭。
- DES-027 D18 十二项验收门逐项有机器证据。
- `enterprise_operational_ready` 可从 clean tenant 确定性运行、恢复、重放和验证。
- 正常主线与四个异常场景共享同一 Runtime Artifact 和正式业务路径。
- `founder-principal`、五 Agent、G1–G7、Capability/Process/Policy/Approval、World Bridge 和 Trace Spine 全部可审计。
- 两仓提交、测试、部署、runbook、evidence、roadmap、code map、progress log 和 Atlas 状态一致。
