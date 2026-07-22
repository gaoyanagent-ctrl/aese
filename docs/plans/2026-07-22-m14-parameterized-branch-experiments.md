---
id: PLAN-M14-001
title: M14 参数化分支经营实验实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m14, simulation, experiment, branch, scenario-lab]
---

# M14 参数化分支经营实验实施计划

## 1. 交付目标与工期

用 **8 到 10 周**把 M13 的单次确定性商业闭环升级为受治理、可复现的多周期经营实验能力：

> 从批准的 Genesis checkpoint 创建隔离分支，在相同外生随机流下比较基准、精益和韧性策略，形成包含 run-level 证据、KPI 分布、paired delta、约束违反和 Pareto 结果的 EvidenceBundle，只有证据完整时输出 `strategy_evidence_ready=true`。

首版仍限定单租户、单工厂、单产品和单客户；默认 12 个虚拟周/12 个订单周期。M14 解决“策略在一组明确假设下是否稳健”，不宣称真实预测精度，不自动应用建议。

## 2. 前置门

- [x] G1 M13 已完成，`hctm-genesis@0.6.0` 输出 `first_commercial_cycle_closed=true`，M9-M13 端到端证据完整。
- [x] G2 ADR-004、DES-008、DES-009 的 World / IAOS / Actor Knowledge 所有权和桥接合同继续有效。
- [x] G3 DES-015 固定 M14 的实验对象、分支隔离、公平比较、IAOS 治理和禁止自动投放边界。
- [x] G4 冻结 allowlist checkpoint、12 周经营基线、订单周期、期初余额/承诺和 M13→M14 opening reconciliation。
- [x] G5 冻结外生参数、分布假设、单位/范围/相关性、PRNG/version、命名随机流、seed-set 数量和共同随机数规则。
- [x] G6 冻结基准/精益/韧性 policy variant、动作边界、硬约束、KPI、paired comparison、Pareto 和证据完整性规则。
- [x] G7 完成 branch/run catalog、持久 artifact、并发/配额/保留/清理、崩溃恢复和性能容量基线。
- [x] G8 完成 IAOS experiment/recommendation gap audit，冻结 payload、权限、simulation namespace 和跨仓交付顺序。

G4-G8 关闭前只允许 Schema、fixture、纯函数随机流/矩阵/统计实现、只读 gap audit 和离线 prototype；不得新增 IAOS 写端点或把实验运行投影成正式经营记录。

## 3. 实施切片

### X0 - 实验方法、基线和机器合同（第 1-2 周）

- [x] T1 定义 ExperimentDefinition、ParameterSet、PolicyVariant、SeedSet、Branch、Run、Metric、Comparison 和 EvidenceBundle stable code/owner。
- [x] T2 冻结 M12/M13/M14 checkpoint allowlist、parent/ancestor hash、opening reconciliation 和失效规则。
- [x] T3 冻结 12 周/12 订单周期的需求、供应、生产、质量、交付、开票、回款和成本时间粒度。
- [x] T4 定义外生参数与可控策略的分离合同、单位/精度/范围/相关性和跨字段约束。
- [x] T5 固定版本 PRNG、`demand/supplier/equipment/quality/payment` 随机流、seed set 和共同随机数 pairing。
- [x] T6 固定硬约束、KPI、paired delta、分位数、失败样本、Pareto 和可选权重规则。
- [x] T7 扩展 JSON Schema、Go strict types/parser、canonical hash、payload registry 及合法/边界/破损 fixture。
- [x] T8 定义 `draft -> validated -> approved -> expanded -> running -> aggregating -> evidence_ready | incomplete | cancelled` 状态机。
- [x] T9 完成 IAOS gap ledger、权限/职责矩阵、Atlas 依赖、两仓交付矩阵和 X0 方法评审。

验收：G4-G8 关闭；相同输入矩阵可获得相同 branch/run identity，任何无界参数、隐藏随机性或结果选择规则都失败关闭。

### X1 - 确定性随机流与参数矩阵（第 2-3 周）

- [x] T10 实现独立命名随机流，证明增加 payment draw 不改变 demand/supplier/equipment/quality 序列。
- [x] T11 实现固定值、离散枚举和首版批准分布的 schema/采样器，禁止墙钟/global rand/map 顺序影响。
- [x] T12 实现 parameter constraint compiler，验证概率、数量、金额、交期、相关性和策略阈值边界。
- [x] T13 实现不可变 matrix expansion、稳定排序、run-count/storage/time estimate 和 matrix hash。
- [x] T14 实现共同随机数 pair assignment，确保同 profile/seed 下各 policy 共享外生 draw 而动作保持独立。
- [x] T15 实现默认 dry-run 的 `experiment validate|inspect|expand` CLI 和可审计 preflight 输出。
- [x] T16 用 golden vector、跨进程和 100 次重复测试验证 stream、matrix 和 input hash 稳定。

验收：任一 run 的全部随机输入可由版本、stream 和 seed 重建，策略比较不因随机样本不同而失真。

### X2 - Checkpoint fork、隔离分支与持久运行目录（第 2-4 周）

- [x] T17 实现 checkpoint compatibility、祖先链、terminal hash 和 M13 opening reconciliation 校验。
- [x] T18 实现 copy-on-write/等价不可变 fork，生成独立 branch ID、cursor、event sequence、snapshot 和 idempotency namespace。
- [x] T19 扩展 World Store 的 experiment/branch/run catalog、状态转换和 artifact manifest，不复制 IAOS 业务数据库。
- [x] T20 实现 branch-local transaction/correlation code 和 simulation namespace，阻止与正式 Genesis/M3/M7 code 冲突。
- [x] T21 实现 artifact 原子发布、checksum、失败隔离、保留策略、配额和安全清理预览。
- [x] T22 验证父 checkpoint 不变、兄弟分支无泄漏、tenant isolation、陈旧 checkpoint、重复 fork 和损坏 snapshot 失败关闭。
- [x] T23 验证中途崩溃后的 catalog/artifact 对账和从最后一致 checkpoint 恢复。

验收：任何分支运行、取消、重试或 reset 都不能修改父状态、兄弟分支或正式 IAOS 经营事实。

### X3 - 多周期 World、策略执行和经济守恒（第 3-6 周）

- [x] T24 从 M13 closed checkpoint 延伸 12 个订单周期，保持产品、产线、财务 opening 和承诺连续。
- [x] T25 实现需求量/到达和客户付款延迟的外生 World policy。
- [x] T26 实现供应交期/可靠性、设备故障/修复和质量良率的外生 World policy。
- [x] T27 实现基准、精益、韧性三种库存/双源/维护/加班/现金保护 policy variant。
- [x] T28 将策略动作建模为 observation -> intent -> governed/simulated outcome -> World consequence，禁止策略函数直接改结果。
- [x] T29 扩展跨周期材料/WIP/成品、订单/接受、应收/现金、actual cost/margin 和资源占用守恒。
- [x] T30 实现硬约束违反、不可行策略、现金破线、质量门失败和恢复/终止语义。
- [x] T31 对单一 seed 完成三策略 tracer，对账 event log、Knowledge、动作、KPI 和 causal refs。
- [x] T32 验证 100 次 replay、snapshot/restore、旧 M13 terminal hash 和 M3-M13 零回归。

验收：多周期后果由相同 World 规则和受治理动作产生；策略不能凭空创造库存、产能、质量、现金或客户接受。

### X4 - 有界实验执行器与生命周期治理（第 4-6 周）

- [x] T33 实现持久 experiment executor、稳定 run 调度、有界 worker pool 和资源配额。
- [x] T34 实现显式 `--apply` 的 `experiment run`、目标环境/tenant 确认、进度和机器可读输出。
- [x] T35 实现 pause/cancel/resume/retry；已完成 run 不重复，失败 run 保留原 identity/seed/evidence。
- [x] T36 实现 per-run timeout、panic/error 捕获、失败隔离和 experiment incomplete 判定。
- [x] T37 实现并发抢占/lease 或等价单执行保证，验证重复点击和多进程竞争不重复产生 run。
- [x] T38 实现运行预算、预估/实际 CPU 时间、artifact 大小和保留到期可观察性。
- [x] T39 完成小/标准/上限矩阵容量测试，冻结首版最大 run 数、并发数和 artifact 配额。

验收：执行器可安全中断和继续，矩阵超限先于写入失败，失败/取消样本不会被静默删除或替换。

### X5 - KPI 聚合、比较与 EvidenceBundle（第 5-7 周）

- [x] T40 实现 run-level 服务、积压、库存/资金、现金、毛利、质量、加班、报废、加急和恢复 KPI。
- [x] T41 实现 paired delta、分位数、约束违反率、missing/failed run 和 denominator 合同。
- [x] T42 实现 Pareto frontier；默认不输出单一“最优”，显式权重必须版本化且进入 evidence hash。
- [x] T43 实现 sensitivity/assumption slice，区分参数假设、运行观察和策略结论。
- [x] T44 生成不可变 EvidenceBundle：定义、输入/hash、完整 run manifest、统计、因果引用、限制和推荐草稿。
- [x] T45 实现 evidence completeness gate；缺 run、hash 不符、聚合版本不符或约束未计算时保持 `incomplete`。
- [x] T46 实现 `experiment compare|evidence|replay` CLI 和 JSON/Markdown 导出。
- [x] T47 用手算小样本、golden dataset、顺序置换和失败样本注入验证聚合不偏移。

验收：用户能从任何聚合结论下钻到 paired run、输入 draw、策略动作和 World 因果链。

### X6 - IAOS 实验治理与推荐边界（第 5-8 周）

- [x] T48 按 IAOS 规则建立新的独立 branch/worktree，先读其 AGENTS、Agent Context 和 Code Map。
- [x] T49 复用 Capability/Policy/Decision/Process/AI Tool，只增加 experiment request/approval/evidence receipt/recommendation 的最小治理能力。
- [x] T50 实现 M14 allowlist payload、journal/cursor 和 simulation namespace，不创建正式订单/库存/发票/现金记录。
- [x] T51 固定 `genesis.experiment.define/run/cancel/read/evidence/propose/decide` 权限、mandate、职责分离和只读角色投影。
- [x] T52 保证 request/approval/evidence/recommendation、audit、journal committed outcome 和 Outbox 必要原子性。
- [x] T53 验证 tenant/RLS、越权审批、自批推荐、篡改 evidence ref、重复 receipt、陈旧 bundle 和失败无部分写入。
- [x] T54 验证推荐只形成独立 intent/decision，不能自动更改正式 Policy、预算、订单、采购、排产或现金。
- [x] T55 两仓分别提交、记录 revision、部署并完成 contract/integration tests。

验收：IAOS 治理“是否运行、谁能看、是否采纳”，AESE 计算“发生什么”；实验结论与正式业务动作之间存在不可绕过的决策门。

### X7 - Scenario Lab 与全链验收（第 7-10 周）

- [x] T56 在 World Play 增加 Scenario Lab：checkpoint、参数、策略、seed、矩阵/资源 preflight 和启动控制。
- [x] T57 展示进度、失败/取消、KPI 分布、paired delta、约束、Pareto 和 run-level 因果下钻。
- [x] T58 清晰标记 Simulation、假设、owner、证据完整性以及 recommendation proposed/decided，禁止一键应用正式策略。
- [x] T59 完成键盘、ARIA、焦点、移动端、图表非颜色表达、数量/金额/概率单位和大矩阵性能。
- [x] T60 执行人类路径和 Agent 推荐草稿路径，对账 World、Knowledge、IAOS journal/Decision/Outbox 和 evidence hash。
- [x] T61 验证断线、SSE 丢失、重启、重复点击、陈旧 cursor、乱序、并发、取消/继续、配额和分层 reset。
- [x] T62 验证 tenant-other、只读用户、跨分支读取、结果篡改、选择性删除、越权采纳和正式业务零污染。
- [x] T63 运行 Go、Schema、PostgreSQL、IAOS modules、前端、Playwright、100 次 hash、容量和 M3-M13 回归。
- [x] T64 编写 M14 runbook/evidence，更新设计、计划、Roadmap、Code Map、Progress Log、Atlas、两仓 revision 和部署信息。

验收：只有矩阵完整、所有 run 可追溯、统计/约束已计算、无未解释失败且 EvidenceBundle hash 可复验时输出 `strategy_evidence_ready=true`。

## 4. 首版实验矩阵

X0 冻结前使用以下设计方向，不把暂定数值当完成合同：

| 维度 | 首版候选 |
| --- | --- |
| Checkpoint | M13 `first_commercial_cycle_closed` |
| Horizon | 12 个虚拟周 / 12 个订单周期 |
| Strategy | baseline / lean / resilient |
| Scenario profile | nominal / demand-spike / supplier-risk / equipment-risk / cash-delay |
| Seed set | 每 profile 固定且跨 strategy 配对；数量由容量基线批准 |
| Hard constraints | 质量门、权限门、数量/金额守恒、最低现金、加班上限 |
| Primary KPI | OTIF、cash trough、gross margin、working capital、recovery time |

实验 code 使用 `EXP-GENESIS-M14-*`，分支、run、correlation 和 artifact 名称由输入 hash 派生。全部主体、订单、概率、价格和财务数据保持虚构。

## 5. 完成定义

- 分支父状态不可变、兄弟隔离、tenant 隔离和正式 IAOS 业务零污染。
- checkpoint、参数、策略、PRNG、seed、规则、聚合和 actor policy 全部版本化、可 hash、可重放。
- 共同随机数和 paired comparison 经过自动化验证；不以未配对样本宣称策略差异。
- 多周期数量、质量、资源、收入/应收/现金和 actual cost/margin 守恒。
- 执行器默认 dry-run，显式 apply、有界并发、配额、取消、继续、重试和崩溃恢复通过。
- EvidenceBundle 保留全部 run/失败/取消，不以平均值掩盖硬约束或缺失样本。
- 推荐不会自动应用；IAOS 权限、职责分离、事务/Outbox、幂等和 RLS 通过。
- Scenario Lab、CLI/API、runbook/evidence、三视口、容量、M3-M13 回归和 Atlas 完整。
- `strategy_evidence_ready=true` 只代表在明确假设与实验范围内证据就绪，不代表已批准或真实世界保证。

## 6. 不纳入 M14

- 真实历史数据校准、在线学习、机器学习训练或真实概率预测。
- 第二客户/产品/工厂、完整滚动 S&OP、供应网络优化或完整财务计划。
- 通用分布式计算平台、任意用户代码执行、无界参数搜索或自动超参数优化。
- 自动修改正式业务 Policy/订单/采购/排产/预算，或让 Agent 自批并实施建议。

## 7. 并行与所有者规则

- X0 由实验合同 owner 串行冻结；X1/X2 可在合同稳定后并行，但 input/branch identity 由单一 owner 合并。
- X3 依赖 X1 随机流和 X2 分支；经济守恒与 M13 opening reconciliation 由单一领域 owner 负责。
- X4 可与 X3 后半段并行，不能为了吞吐放松确定性、隔离或证据完整性。
- X5 只消费冻结的 run artifact；统计/聚合版本改变必须使旧 EvidenceBundle 明确失效或迁移。
- X6 只能在 G8 关闭后由独立 IAOS worktree owner 开始，AESE 和 IAOS 提交/测试/日志分别维护。
- X7 在 API/view model 稳定后开始；图表不得代替 run-level evidence，UI 不得直接修改结果。
- 最终收口串行核对两仓、全部 seed、失败样本、revision、部署和 Atlas；并行 agent 必须有子计划、owner 和不重叠 worktree。
- 保留共享工作区中现有测试修改、截图变化和验收产物，不得覆盖或回滚。
