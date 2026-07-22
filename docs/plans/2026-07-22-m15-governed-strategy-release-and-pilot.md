---
id: PLAN-M15-001
title: M15 受治理策略发布与经营试点实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m15, strategy, governance, pilot, rollback]
---

# M15 受治理策略发布与经营试点实施计划

## 1. 交付目标与工期

用 **8 到 10 周**把 M14 的不可变策略证据推进为可审议、可 shadow、可有限试点、可暂停/回滚的治理闭环：

> 选择一个 M14 Pareto candidate，经过证据复验和独立审批，编译为版本化 StrategyRelease；先运行零写入 shadow，再在新的 Genesis canonical operating cycle 中按安全包络有限生效；最终以 adopted、rejected 或 rolled_back 关闭决策周期，输出 `strategy_change_cycle_closed=true`。

首版仍限定单租户、单工厂、单产品和单客户。M15 验证的是“组织能否安全地把证据转成受控行动”，不是强行证明候选策略更优。

## 2. 前置门

- [x] G1 M14 已完成，60-run EvidenceBundle 输出 `strategy_evidence_ready=true`，推荐状态为 `proposed_not_applied` 且 production writes 为零。
- [x] G2 ADR-004、DES-008、DES-009 和 DES-015 的三态、bridge、实验证据及禁止自动投放边界继续有效。
- [x] G3 DES-016 固定 M15 的 candidate/release、双阶段审批、shadow/pilot、guardrail、回滚/补偿和诚实终态。
- [x] G4 冻结 candidate 选择程序、M14 evidence/hash、assumption support、审议角色、利益冲突和 rejection 规则；不得预设 winner。
- [x] G5 冻结 StrategyRelease schema、semantic diff、Policy/Capability 映射、目标 canonical checkpoint、适用对象和 release compatibility。
- [x] G6 冻结 shadow/pilot 窗口、active/candidate 对照、hard/pause/informational guardrail、数据 freshness 和 assumption drift 规则。
- [x] G7 冻结 prior release、kill switch、open-commitment ledger、rollback/compensation、恢复点和不可逆动作分类。
- [x] G8 完成 IAOS change/release/incident/adoption gap audit，冻结权限、职责分离、payload、事务边界和跨仓顺序。

G4-G8 关闭前只允许 Schema、fixture、evidence verifier、semantic diff、shadow 纯函数、只读 gap audit 和离线 prototype；不得激活 Policy、提交 canonical 业务动作或创建 pilot 事实。

## 3. 实施切片

### R0 - 决策、发布和安全机器合同（第 1-2 周）

- [x] T1 定义 StrategyCandidate、ChangeRequest、StrategyRelease、SafetyEnvelope、ShadowRun、PilotCycle、GuardrailBreach、RollbackPlan、Commitment 和 AdoptionDecision stable code/owner。
- [x] T2 冻结 M14 EvidenceBundle verify、candidate eligibility、Pareto/constraint refs、assumption support 和 evidence expiry。
- [x] T3 冻结 proposer/reviewer/approver、业务 owner、CFO/风险角色、Agent 边界和利益冲突规则。
- [x] T4 固定 active release、candidate release、target checkpoint、scope、effective window、prior release 和 expected version。
- [x] T5 固定 hard stop、pause/review、informational guardrail、measurement owner、freshness 和 missing-data 行为。
- [x] T6 固定 shadow/pilot 窗口、promotion gates、adopt/reject/rollback disposition 和 `strategy_change_cycle_closed` 规则。
- [x] T7 固定 kill switch、rollback、open commitments、compensating action 和不可逆 consequence 语义。
- [x] T8 扩展 JSON Schema、Go strict types/parser、canonical hash、payload registry 及合法/边界/破损 fixture。
- [x] T9 定义 `candidate -> evidence_verified -> review -> shadow -> pilot -> adopted|rejected|rolled_back -> closed` 状态机并完成 G4-G8 评审。

验收：任何 candidate 都不能跳过 evidence、owner、监控、回滚或独立审批；拒绝/回滚与采纳均能合法关闭周期。

### R1 - Evidence Review 与 ChangeRequest（第 1-3 周）

- [x] T10 实现 M14 bundle checksum、run completeness、policy/input/version、constraint/Pareto 和 freshness 复验。
- [x] T11 实现 candidate dossier，保留收益、代价、假设、限制、失败样本、敏感性和不支持域。
- [x] T12 实现 candidate 选择说明和 rejection reason，禁止按平均值、单一 run 或隐藏权重挑选。
- [x] T13 实现 StrategyChangeRequest dry-run/impact report，列出目标对象、动作、审批、风险和不可逆后果。
- [x] T14 实现 proposer/approver conflict、mandate、证据访问和 actor-scoped Knowledge 投影。
- [x] T15 实现 evidence expired/changed/unknown、candidate 非 Pareto、缺约束和跨租户引用失败关闭。
- [x] T16 实现人类提议与 Agent-prepared proposal 两条路径，Agent 路径强制 human accountable owner。
- [x] T17 对账 ChangeRequest、Decision、journal、Outbox、evidence refs 和 idempotency。

验收：进入发布编译前，候选为何被选、谁负责、证据支持什么/不支持什么均机器可验证。

### R2 - StrategyRelease 编译、差异与预检（第 2-4 周）

- [x] T18 定义安全库存/双源、维护/加班和现金保护 allowlist 参数到 IAOS Policy/Capability 的映射。
- [x] T19 实现 immutable StrategyRelease manifest、version/hash、effective window、scope 和 prior-release link。
- [x] T20 实现 current vs candidate semantic diff，显示默认值、单位、上下界、动作/权限和 owner 变化。
- [x] T21 实现 AESE World rule refs 与 IAOS release manifest compatibility，禁止 release 直接改 World state。
- [x] T22 实现 preflight：受影响对象、预计动作量、库存/现金/产能边界、未结承诺和 rollback 可行性。
- [x] T23 实现默认 dry-run 的 `strategy validate|inspect|diff|preflight` CLI，所有写入要求显式目标和 `--apply`。
- [x] T24 验证未知参数、任意代码/API、越界阈值、陈旧 checkpoint、hash/version 不符和不可回滚 release 失败关闭。
- [x] T25 100 次编译验证 release/diff/preflight hash 稳定，M14 evidence 和 M13 parent 不被修改。

验收：同一 candidate 只能产生一个确定性 release；所有业务影响和退出路径在激活前可见。

### R3 - 零写入 Shadow 与运行兼容性（第 3-5 周）

- [x] T26 从批准 canonical checkpoint 创建 shadow observer，复用同一 observation/cursor 而不 fork 正式业务事实。
- [x] T27 运行 active 与 candidate decision evaluator，保留 input freshness、Knowledge scope、reason 和 decision diff。
- [x] T28 强制 shadow 禁止 intent、Capability 写调用、committed outcome 和 World resource reservation。
- [x] T29 实现 candidate action frequency、churn、越界、不可见事实依赖、assumption drift 和数据缺口指标。
- [x] T30 实现 shadow pause/resume/replay，断线和陈旧 cursor 从 journal 恢复且不重复 decision。
- [x] T31 验证 4 周或批准窗口的完整性、零业务写入、跨租户、权限降级和监控缺失失败关闭。
- [x] T32 生成 ShadowEvidence，并由独立角色作 `reject` 或 `approve_for_pilot` 决策。
- [x] T33 100 次 shadow replay 验证 decision/evidence hash，active canonical result 不受 candidate 影响。

验收：shadow 只证明候选能否安全读取和给出建议；没有任何采购、生产、库存、发运、发票或现金副作用。

### R4 - Canonical Pilot 与受治理动作（第 4-7 周）

- [x] T34 从新的 Genesis canonical operating checkpoint 创建唯一 PilotCycle，绑定 exact release、scope、window 和 correlation namespace。
- [x] T35 实现独立 pilot approval、激活 CAS、effective time 和 prior release snapshot，重复/并发激活收敛。
- [x] T36 将 candidate decision 转成 IAOS intent；所有动作走 Policy/Capability/Process/Decision、RLS、事务和 Outbox。
- [x] T37 由 AESE 计算需求、供应、设备、质量、运输、银行和资源/资金 consequence，IAOS 状态不能直接创造结果。
- [x] T38 实现 active release 可观察性、每次 decision/action/world consequence 因果引用和 actor Knowledge 时效。
- [x] T39 执行最多 4 周的库存/供应、维护/加班和现金保护 allowlist pilot，不扩展其他业务范围。
- [x] T40 对账订单、采购、生产、库存、交付、AR/现金、actual cost/margin、guardrail 和 open commitments。
- [x] T41 验证越权、自批、陈旧 release、重复动作、负库存、超预算、无 observation 和失败无部分写入。
- [x] T42 验证暂停/重启、SSE 丢失、cursor 恢复和 PilotCycle reset 不删除已提交业务/World 事实。

验收：策略只影响允许的未来决策；每个现实后果仍来自 World，所有正式业务变化均来自 IAOS committed outcome。

### R5 - Guardrail、暂停、回滚与补偿（第 5-8 周）

- [x] T43 实现 hard/pause/informational guardrail evaluator，固定 measurement source、freshness、severity 和 owner。
- [x] T44 实现 assumption-support/drift 检查，超出 M14 支持域时自动 pause 并要求新决策。
- [x] T45 实现 emergency kill switch，只停止未来 intent，不删除 committed outcome 或 World consequence。
- [x] T46 实现 rollback CAS，原子恢复 prior release 并生成 rollback committed outcome/journal/Outbox。
- [x] T47 实现 open-commitment ledger，列出 PO、工单、库存、发运、发票、现金和其他不可逆/待处理事实。
- [x] T48 实现取消、变更、消耗、隔离等 allowlist compensating action，并保持原始事实与审计。
- [x] T49 注入质量 hard stop、现金 pause、监控缺失和 drift 四类失败，验证无后续越界动作。
- [x] T50 验证并发 pause/rollback、重复 kill、补偿失败、服务崩溃和恢复后状态/commitment 一致。

验收：回滚不等于历史删除；release、未结承诺和补偿结果可逐项对账，任何未知风险阻止周期关闭。

### R6 - IAOS 治理、决策与采纳（第 4-8 周）

- [x] T51 按 IAOS 规则建立新的独立 branch/worktree，先读其 AGENTS、Agent Context 和 Code Map。
- [x] T52 复用 Policy/Capability/Process/Decision/AI Tool，只增加 change/release/shadow/pilot/incident/rollback/adoption 最小能力。
- [x] T53 实现 M15 allowlist payload、journal/cursor 和 StrategyRelease exact-hash/expected-version enforcement。
- [x] T54 固定 `genesis.strategy.propose/review/shadow/pilot/activate/pause/rollback/compensate/adopt/read` 权限和职责分离。
- [x] T55 保证 release activation、business change、incident、rollback、compensation、audit、journal 和 Outbox 的必要事务原子性。
- [x] T56 实现 Adoption Review：M14 evidence、shadow、pilot、guardrail、drift、commitment、人工接管和限制完整投影。
- [x] T57 实现 adopted/rejected/rolled_back 决策、复审日期、失效条件和关闭门；禁止 Agent sole approval。
- [x] T58 验证 tenant/RLS、越权/自批、篡改 evidence/release、重复/并发决策、陈旧 review 和失败无部分写入。
- [x] T59 两仓分别提交、记录 revision、部署并完成 contract/integration tests。

验收：IAOS 对谁能提议、审批、激活、暂停、回滚和采纳负责；AESE 不复制 Policy/权限/业务数据库。

### R7 - Strategy Control Room 与全链验收（第 7-10 周）

- [x] T60 在 World Play 增加 Strategy Control Room：evidence、semantic diff、审批、shadow、pilot、guardrail、drift、commitment 和 disposition。
- [x] T61 展示 World/IAOS/Knowledge owner、Simulation 与 canonical pilot、推荐/批准/激活/采纳的状态差异。
- [x] T62 危险动作显示 impact preview、exact release、目标环境、不可逆后果和二次确认；不提供跳门或直接改 World 入口。
- [x] T63 完成键盘、ARIA、焦点、移动端、非颜色告警、金额/数量/概率单位和实时/虚拟时间表达。
- [x] T64 执行 adopted 路径和 injected rollback 路径，终态均能诚实关闭且保留全部 consequence/compensation。
- [x] T65 对账 EvidenceBundle、ChangeRequest、Release、Decision、journal、Outbox、Policy、业务记录、World events、Knowledge 和 hashes。
- [x] T66 验证断线、重启、重复点击、陈旧 cursor/version、乱序、并发、监控故障、pause/rollback 和分层 reset。
- [x] T67 验证 tenant-other、只读用户、Agent 自批、跨 scope 动作、结果删除、直接投放和 real-production target 默认拒绝。
- [x] T68 运行 Go、Schema、PostgreSQL、IAOS modules、前端、Playwright、100 次 hash 和 M3-M14 回归。
- [x] T69 编写 M15 runbook/evidence，更新设计、计划、Roadmap、Code Map、Progress Log、Atlas、两仓 revision 和部署信息。

验收：`strategy_change_cycle_closed=true` 必须携带 adopted/rejected/rolled_back disposition、完整证据、已处理 breach 和已对账 commitment；不得用“里程碑完成”强迫策略采纳。

## 4. 首版试点方向

R0 冻结前只使用下列设计方向，不预先选择 M14 winner 或业务阈值：

| 维度 | 首版方向 |
| --- | --- |
| Candidate | M14 Pareto set 中通过 evidence review 的一个版本 |
| Shadow | 4 个虚拟周，复用 canonical observation，业务写入为零 |
| Pilot | 最多 4 个虚拟周，单厂/单产品/单客户，allowlist 动作 |
| Release scope | safety stock / dual source / maintenance / overtime / cash guard |
| Hard stop | 质量、权限、数量金额守恒、负库存、最低现金 |
| Pause | OTIF/backlog、加班、供应/设备风险、assumption drift、监控缺失 |
| Disposition | adopted / rejected / rolled_back |

策略、release、pilot、incident 和 rollback code 使用 `STR-GENESIS-M15-*`。所有客户、订单、概率、价格和财务数据保持虚构；真实生产 target 不属于本计划授权范围。

## 5. 完成定义

- candidate 与 exact M14 EvidenceBundle/hash 绑定，选择过程、限制和责任人可审计。
- StrategyRelease、semantic diff、SafetyEnvelope、RollbackPlan 和 AdoptionDecision 可版本化、可 hash、可重放。
- shadow 完整窗口零业务写入；candidate/active decision diff 和数据 freshness 可下钻。
- pilot 只在批准 scope/window 内通过 IAOS 治理动作，World consequence 与业务记录严格对账。
- hard stop、pause、kill switch、rollback、open commitment 和 compensation 失败路径通过。
- 回滚恢复未来策略但保留历史事实；任何未结承诺都有 owner、状态和处理证据。
- tenant/RLS、职责分离、幂等、并发、事务/Outbox、恢复和 real-production 默认拒绝通过。
- Control Room、CLI/API、runbook/evidence、三视口、两仓 revision/部署、Atlas 和 M3-M14 回归完整。
- 最终输出 `strategy_change_cycle_closed=true` 与明确 disposition，不宣称 pilot 提供统计因果证明。

## 6. 不纳入 M15

- 真实客户或生产租户投放、无人审批自动发布、Agent sole approval。
- 第二客户/产品/工厂、完整 S&OP、动态定价、组织/资本重大变更。
- 删除业务事实式回滚、任意脚本/Policy 代码执行或绕过 IAOS Capability 的动作。
- 真实数据校准、在线学习、统计因果平台或永久无人值守策略优化。

## 7. 并行与所有者规则

- R0 由合同 owner 串行冻结；R1/R2 可在 evidence/release identity 稳定后并行，hash/compatibility 由单一 owner 合并。
- R3 依赖 R1/R2，但可与 R6 的只读治理合同并行；shadow owner 不得自行批准 pilot。
- R4 只能在 shadow evidence 和 pilot approval 完成后开始；canonical business/World consequence 由各自单一 owner 维护。
- R5 在 R0 即设计、R4 联调；不得等发生 breach 后才补 rollback/compensation。
- R6 只能由新的 IAOS worktree owner 开发，AESE/IAOS 提交、测试、日志和部署证据分别维护。
- R7 在 API/view model 稳定后开始；UI 不得弱化职责分离或为演示跳过审批/guardrail。
- 最终收口串行核对两仓、adopted 与 rollback 路径、所有 commitment、revision、部署和 Atlas。
- 保留共享工作区中现有测试修改、截图删除和验收产物，不得覆盖或回滚。
