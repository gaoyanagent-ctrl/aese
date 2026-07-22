---
id: DES-015
title: M14 参数化分支经营实验
date: 2026-07-22
status: completed
author: Codex + User
tags: [m14, simulation, experiment, branch, decision-support]
---

# M14 参数化分支经营实验

## 1. 决策摘要

M14 把 M13 已验证的单次商业闭环升级为可重复的经营实验系统：从受批准的 Genesis checkpoint 创建相互隔离的 World 分支，在相同外生扰动下比较不同经营策略，形成可追溯的多周期 KPI、差异解释和决策证据包。

M14 的产品终态为：

> 指定基线、参数、策略、随机流、目标和约束后，AESE 能稳定执行一组隔离实验，重复运行得到相同结果；IAOS 能受治理地登记实验意图、审批需要治理的实验和接收证据/建议，但任何建议都不会自动修改基线业务状态。

机器终态固定为 `strategy_evidence_ready=true`。它表示证据足以供人或 Agent 发起后续治理决策，不表示某策略已获批准、已部署或必然最优。

## 2. 为什么是 M14

M9-M13 已证明企业从成立到首单回款和实际毛利关闭的单条确定性路径，但仍不能回答：

- 需求、供应、设备、质量和付款条件变化时，当前策略是否仍稳健。
- 安全库存、双源、维护、加班和现金保护策略各自付出什么代价。
- 某个结果来自策略本身，还是来自一组更有利的随机事件。
- 哪些建议跨多次运行仍满足交付、质量、现金和治理硬约束。

因此 M14 优先建设实验方法和证据治理，而不是扩展第二工厂、完整 ERP 模块或更华丽的 2D/3D 表现。

## 3. 范围

首版固定：

- `tenant-hctm`、苏州基地、电池冷却板 A 线、单客户和单产品。
- 从批准的 M13 checkpoint 延伸 12 个虚拟周、12 个订单周期；具体业务数值在 X0 冻结。
- 三类策略族：库存/供应、产能/维护、资金/回款保护。
- 五类外生变量：需求、供应交期与可靠性、设备故障与修复、质量良率、客户付款延迟。
- 基准、精益、韧性三种首批 policy variant，并允许受 schema 约束的新增 variant。
- 固定版本 PRNG、命名随机流、共同随机数和批准的 seed set。
- 交付、积压、库存/营运资金、现金低点、毛利、加班、报废、加急和恢复时间 KPI。

不包含：

- 自动把胜出策略写回正式订单、采购、排产、价格、预算或 Policy。
- 用单次运行宣称因果、概率、最优策略或真实世界预测精度。
- 第二客户/产品/工厂、完整 S&OP、数字孪生校准、机器学习训练平台或通用分布式计算平台。
- 真实客户数据、真实概率分布或生产环境直接试验。

## 4. 核心领域合同

| 对象 | 所有者 | 关键字段 |
| --- | --- | --- |
| ExperimentDefinition | AESE；IAOS 保存治理登记 | experiment code、目标、约束、checkpoint、horizon、metric set |
| ParameterSet | AESE | 外生变量、单位/精度、范围、distribution assumption、版本/hash |
| PolicyVariant | AESE 定义仿真规则；IAOS 保存业务建议/批准 | policy code、阈值、动作边界、版本/hash |
| SeedSet | AESE | PRNG/version、stream names、seed values、pairing rule、hash |
| ExperimentBranch | AESE World | parent snapshot/hash、branch ID、parameter/policy/seed refs、状态 |
| ExperimentRun | AESE World | run ID、event log、snapshot、state hash、metrics、failure reason |
| Comparison | AESE | paired delta、分位数、约束违反、Pareto 状态、解释 refs |
| EvidenceBundle | AESE 生成；IAOS 接收不可变引用 | 输入/输出 hash、完整性、审计、结论限制、推荐草稿 |

所有引用使用稳定业务编码；时间使用 RFC 3339 和 `Asia/Shanghai`；数量、金额和比率显式定义单位、精度和舍入。实验合同必须严格拒绝未知字段、无界参数、无版本规则、缺失 seed 或不兼容 checkpoint。

## 5. Checkpoint 与分支语义

首版只允许从 allowlist checkpoint 创建分支：

1. M12 `serial_production_eligible`：用于比较首个商业周期策略。
2. M13 `first_commercial_cycle_closed`：用于连续经营和工作资金实验。
3. 经批准的 M14 周期末 checkpoint：用于滚动延伸，但必须保留完整祖先链。

父 checkpoint 和父分支不可变。创建分支时记录 parent state hash、pack/rules version、parameter hash、policy hash、seed-set hash 和 actor-policy version。分支拥有独立事件序列、cursor、快照、幂等命名空间和 artifact 路径，不得复用会导致 IAOS 业务写入冲突的 transaction code。

实验默认在隔离 World 中运行。除实验登记、治理审批和 evidence receipt 外，不向 IAOS 创建正式订单、库存、发票、现金或成本记录。需要验证治理动作时使用明确标记的 simulation scope/namespace；从实验结论到正式业务动作必须由后续独立 intent 和审批完成。

## 6. 参数与策略边界

参数分三类：

- 固定合同：产品/BOM/routing、初始资产、权限、会计口径和不可突破的质量门。
- 外生假设：订单量/到达、供应交期/失败、设备故障/修复、良率、付款延迟。
- 可控策略：安全库存、采购分配、维护触发、加班门、现金缓冲和信用保护。

外生假设和策略动作不能混为一个字段。参数必须具有合理上下界和约束，例如概率位于 `[0,1]`、交期非负、策略不能允许负库存、越权自批、未接受开票或伪造银行到账。

首版不默认计算单一综合分数。比较先展示硬约束，再展示各 KPI 的 paired delta、分布和 Pareto 前沿。只有实验定义显式给出版本化权重并获批准时才展示综合评分。

## 7. 随机性与公平比较

随机性是可重放输入，不是隐藏的运行副作用：

- 固定 PRNG 算法和版本，禁止依赖 Go map 顺序、墙钟或全局随机源。
- 按 `demand`、`supplier`、`equipment`、`quality`、`payment` 命名独立随机流。
- 同一 scenario profile 的各 policy variant 使用共同随机数，进行 paired comparison。
- seed set 在执行前冻结；失败后重试仍使用原 seed，不允许只重跑不利结果。
- 同一输入重复 100 次必须产生一致 event-log hash、state hash 和 KPI hash。

概率分布只是版本化场景假设，不声称来自真实数据。报告必须展示假设、样本数、缺失数据和外推限制；单次运行只用于 tracer，不进入稳健性结论。

## 8. 执行与调度

```text
Experiment Definition + approved checkpoint
  -> offline validate / preflight / run-count and quota estimate
  -> immutable matrix expansion
  -> isolated branch creation
  -> bounded worker execution with deterministic seed streams
  -> per-run invariants, snapshots, hashes and metrics
  -> paired aggregation and constraint evaluation
  -> evidence bundle
  -> optional IAOS recommendation intent (no auto-apply)
```

默认 CLI 和 UI 操作为 dry-run/preflight。真正创建分支或运行矩阵要求显式 `--apply` 或等价确认，并明确目标环境、tenant、run 数、时间/存储预算。执行器支持限流、取消、失败隔离和从已完成 run 继续；取消不能删除已有证据，重试不能改变 run identity。

首版是单机有界并发执行器和持久 run catalog，不建设通用分布式队列。矩阵超过批准配额、估算资源不足或参数笛卡尔积意外膨胀时失败关闭。

## 9. KPI 与证据

每个 run 至少输出：

- OTIF/客户接受服务水平、backlog 峰值和恢复时间。
- 原料/在制/成品库存、库存周转和营运资金占用。
- closing cash、cash trough、应收天数和现金缓冲违反次数。
- 收入、actual cost、gross margin、报废、返工、加班和加急成本。
- 质量门、权限门、数量/金额守恒及其他 invariant 违反。

Comparison 必须保留 run-level 数据，不能只保存平均值。EvidenceBundle 包含定义、checkpoint、矩阵、所有版本/hash、完整 run 清单、失败/取消清单、聚合方法、KPI 分布、paired deltas、Pareto 集、结论限制和推荐草稿。部分运行失败时默认 `incomplete`，不得通过忽略失败样本生成 `strategy_evidence_ready=true`。

## 10. 三态与 IAOS 合同

AESE World 拥有分支、随机流、外生事件、物理/经济后果、run artifact 和实验 KPI。IAOS 拥有实验申请/审批、组织 mandate、预算/风险约束、推荐和后续业务决策。Actor Knowledge 只包含角色已获准看到的参数、运行状态和结果摘要。

Bridge 只增加实验治理所需的 allowlist payload family：

```text
genesis.experiment.requested.v1
genesis.experiment.approved.v1
genesis.experiment.evidence.received.v1
genesis.strategy.recommendation.proposed.v1
genesis.strategy.recommendation.decided.v1
```

这些 outcome 只控制实验和建议的治理状态，不创造 World 结果，也不修改正式经营记录。观察、意图和 committed outcome 继续遵守 DES-008；通知可用 SSE/Outbox，恢复必须读取 journal/cursor。

## 11. Scenario Lab

World Play 增加 Scenario Lab：选择 checkpoint、配置参数和策略、预览矩阵/资源、启动/暂停/取消、观察进度、比较 KPI 分布和下钻单一 run 因果链。界面必须清楚标记：

- Simulation / Not production。
- 参数假设和 seed set。
- World、IAOS 和 Knowledge 的数据 owner。
- 基线、策略差异、硬约束违反和不确定性。
- “推荐草稿”与“已批准策略”的状态差异。

UI 不提供绕过 schema 的任意 JSON、直接改 run 结果、删除不利样本或一键应用正式策略的入口。

## 12. 完成标准

- 从批准 checkpoint 建立的分支相互隔离，父状态与其他分支不被修改。
- 参数、策略、seed、规则和 actor policy 全部版本化并进入 canonical hash。
- 同输入 100 次重放 hash 一致；跨策略共同随机数和 paired comparison 可验证。
- 至少三种策略在批准的多周期矩阵上完成，并输出 run-level、聚合和 Pareto 证据。
- 数量、质量、资金、权限和事件顺序 invariant 在每个 run 都执行；失败/取消不被静默过滤。
- CLI 默认 dry-run；配额、取消、继续、重试、崩溃恢复、tenant isolation 和 artifact retention 通过验收。
- Scenario Lab、API、runbook、evidence、M3-M13 回归、Atlas 和两仓 revision/部署记录完整。
- 只有证据完整、约束计算完成且无未解释 run 缺失时输出 `strategy_evidence_ready=true`。

## 13. 后续边界

M14 完成后再决定 M15。候选方向包括：用经过治理的真实/合成历史校准分布、扩展第二产品/客户、建立滚动 S&OP、或把已批准策略以独立变更流程投放到新的经营周期。任何方向都不能把 M14 的模拟推荐自动当作业务批准。
