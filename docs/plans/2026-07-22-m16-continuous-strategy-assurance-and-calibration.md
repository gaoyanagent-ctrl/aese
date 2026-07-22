---
id: PLAN-M16-001
title: M16 持续策略保障与假设校准实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m16, strategy, assurance, drift, calibration]
---

# M16 持续策略保障与假设校准实施计划

## 1. 交付目标与工期

用 **8 到 10 周**把 M15 已采纳的 resilient StrategyRelease 推进为可重复的保障/校准周期：

> 从 12 周 canonical 经营观察构建不可变 AssuranceDataset，先关闭数据质量问题，再评估需求、供应、设备、质量、付款和策略行为漂移；用前 8 周提出有界校准、后 4 周独立验证，并以新祖先链重跑 M14 策略矩阵；最终输出 `strategy_assurance_cycle_closed=true` 和 `renewed|reexperiment_required|retired`。

首版仍限定单租户、单工厂、单产品和单客户。M16 不自动学习、不自动调 Policy，也不扩大企业范围。

## 2. 前置门

- [x] G1 M15 已完成；主路径采纳 `STR-GENESIS-M15-RESILIENT@1.0.0`，shadow 零写入、pilot 和 injected rollback/compensation 证据完整。
- [x] G2 ADR-004、DES-008、DES-009、DES-015、DES-016 的三态、桥接、实验与发布边界继续有效。
- [x] G3 DES-017 固定 observation/dataset、drift、校准防泄漏、再实验和保障 disposition。
- [x] G4 冻结 12 周窗口、8/4 时间切分、cutoff/lateness/correction、数据 owner、单位/精度和 lineage spec。
- [x] G5 冻结五域 input/process drift、policy-action drift、数据质量优先级、阈值、方法版本和多重比较规则。
- [x] G6 冻结 allowlist calibration 参数/分布族、上下界、样本充分性、holdout 通过条件和禁止反复调参规则。
- [x] G7 冻结 M14 replay ancestry、matrix/seed/common-random-number、ValidationReport 和原 EvidenceBundle 不可变约束。
- [x] G8 完成 IAOS assurance/dataset/drift/decision gap audit，冻结权限、过期/续期/退役、payload 和跨仓顺序。

G4-G8 关闭前只允许 Schema、fixture、只读 source audit、dataset compiler prototype 和确定性统计 golden test；不得改变 active release、启动新 pilot 或发布校准参数。

## 3. 实施切片

### A0 - Assurance、Dataset 与 Drift 机器合同（第 1-2 周）

- [x] T1 定义 AssuranceCycle、ObservationSpec、Dataset、CorrectionSet、DataQualityFinding、DriftAssessment、CalibrationCandidate、ValidationReport 和 AssuranceDecision stable code/owner。
- [x] T2 冻结 M15 adopted release、review/expiry、canonical run/checkpoint、scope 和 prior evidence/hash。
- [x] T3 冻结 12 周 observation、8 周 calibration、4 周 holdout、cutoff/as-of/lateness/correction 和时区合同。
- [x] T4 固定 World event、IAOS record/journal、policy action、guardrail 和 KPI source refs，不以 Knowledge/UI/通知代替事实。
- [x] T5 固定 data-quality -> input -> process/outcome -> policy-action drift 的检测顺序和 `indeterminate` 语义。
- [x] T6 固定 allowlist 分布/参数、上下界、样本充分性、holdout、replay 和 disposition 规则。
- [x] T7 固定统计算法、版本、排序、精度、舍入、canonical encoding 和多重比较处理。
- [x] T8 扩展 JSON Schema、Go strict types/parser、canonical hash、payload registry 及合法/边界/破损 fixture。
- [x] T9 定义 `requested -> collecting -> sealed -> assessed -> calibrating? -> validating -> review -> renewed|reexperiment_required|retired -> closed` 状态机并关闭 G4-G8。

验收：同一 source/cutoff 只能产生同一 dataset identity；数据不完整时不能越级生成 drift、校准或续期结论。

### A1 - Canonical Observation Dataset 与 Lineage（第 2-4 周）

- [x] T10 实现 ObservationSpec compiler，校验 owner、source、unit、precision、freshness、lateness 和 correction policy。
- [x] T11 读取 canonical World event/snapshot refs 和 active release/rules/actor-policy version，保持 World owner。
- [x] T12 通过 IAOS read API/journal cursor 获取允许的业务 stable refs/hash，不复制 IAOS 业务数据库。
- [x] T13 实现 as-of cutoff、cursor range、stable ordering、dedup、source checksum 和 dataset atomic seal。
- [x] T14 实现 missing/late/duplicate/conflict/unit/version/owner findings 和 affected-range lineage。
- [x] T15 实现 cutoff 后 correction set 和 dataset v2 ancestry，禁止静默修改已封存 hash。
- [x] T16 实现默认 dry-run 的 `assurance validate|inspect|collect|seal` CLI，写 artifact 要求显式 `--apply` 和目标环境。
- [x] T17 验证跨租户、cursor gap、source unavailable、陈旧 release、future event、损坏 correction 和重复 seal 失败关闭。
- [x] T18 用 100 次构建、乱序输入和重启恢复验证 dataset/hash/lineage 稳定。

验收：任一数值可下钻到 source owner/ref、sim/business time、cutoff、转换和 correction；Knowledge 不扩权。

### A2 - 数据质量、Drift 与策略压力（第 3-5 周）

- [x] T19 实现 data-quality gate；未关闭 severity finding 时后续 assessment 为 `indeterminate`。
- [x] T20 实现 demand 输入 drift：量、到达节奏和波动相对 M14 assumption support 的可解释比较。
- [x] T21 实现 supplier 输入 drift：交期、延误和可用性，并区分供应事件与采购业务状态。
- [x] T22 实现 equipment/quality 输入 drift：故障/修复、良率/缺陷，保持设备/质量 World 事实所有权。
- [x] T23 实现 payment 输入 drift：付款延迟和到账行为，不把 AR/发票状态当银行事实。
- [x] T24 实现 process/outcome mismatch：产能、OTIF、库存、现金、成本/毛利的预测/观察偏差。
- [x] T25 实现 policy-action drift：动作频率、churn、人工接管、guardrail pressure、越界尝试和 scope 使用。
- [x] T26 实现 frozen threshold/method registry、样本量、置信限制和多重比较标记。
- [x] T27 注入缺失数据、输入漂移、过程失配和策略压力四条 tracer，验证告警/indeterminate 顺序。

验收：系统能说明“数据坏了”“环境变了”“过程模型错了”或“策略承压”，不把所有异常笼统称为 drift。

### A3 - 有界 CalibrationCandidate 与防泄漏（第 4-6 周）

- [x] T28 从 sealed dataset 只选择前 8 周 calibration window，后 4 周对拟合代码不可见。
- [x] T29 实现五域批准分布族的有界参数估计，保留 parent、样本、方法、缺失处理和适用限制。
- [x] T30 实现 parameter diff、range/constraint、相关性和业务单位校验，禁止改 hard constraint/KPI/Policy。
- [x] T31 实现 sample sufficiency、稳定性和 sensitivity；不足时输出 no-calibration/reexperiment reason。
- [x] T32 冻结唯一 CalibrationCandidate/hash 后解锁 holdout，禁止二次拟合或删除不利周。
- [x] T33 实现 calibration dry-run/report 和独立 reviewer approval；Agent 只能准备 proposal。
- [x] T34 验证 holdout leakage、越界参数、换指标、反复候选、陈旧 dataset 和 parent hash 不符失败关闭。
- [x] T35 用 golden dataset 和 100 次运行验证估计、diff、report hash 与舍入稳定。

验收：校准只产生一个受治理 assumption candidate，绝不静默改变 M14 evidence、M15 release 或 canonical history。

### A4 - Holdout、再实验与 ValidationReport（第 5-7 周）

- [x] T36 在冻结候选后读取后 4 周 holdout，一次计算 parent/candidate coverage、error 和 constraint metrics。
- [x] T37 记录 holdout incomplete、out-of-support、结果恶化和多指标冲突，禁止只报告有利指标。
- [x] T38 从原 M14 definition 创建新 ancestry 的 replay definition，保留旧 evidence/seed/run/artifact 不变。
- [x] T39 以批准 candidate 和共同随机数重跑 baseline/lean/resilient/current-adopted 对应策略矩阵。
- [x] T40 保留完整 completed/failed/cancelled manifest、run-level KPI、paired delta、constraint 和 Pareto。
- [x] T41 比较 original vs recalibrated evidence，区分 assumption change、rule change 和 policy effect。
- [x] T42 生成 immutable ValidationReport，绑定 dataset/candidate/holdout/replay/method/hash 和结论限制。
- [x] T43 验证失败样本过滤、seed 变更、旧 artifact 覆盖、矩阵不完整和 hash 不符失败关闭。
- [x] T44 100 次 replay 验证 ValidationReport hash，并回归原 M14 EvidenceBundle hash 不变。

验收：校准是否改善 holdout、是否改变策略稳健性均有证据；结果仍不被描述为因果或真实概率证明。

### A5 - 周期监控、到期与 IAOS 治理（第 4-8 周）

- [x] T45 按 IAOS 规则建立新的独立 branch/worktree，先读其 AGENTS、Agent Context 和 Code Map。
- [x] T46 复用 Policy/Decision/Process/Capability/AI Tool，只增加 assurance request/dataset/drift/calibration/validation/decision 最小能力。
- [x] T47 实现 M16 allowlist payload、journal/cursor、exact dataset/report hash 和 expected-version enforcement。
- [x] T48 固定 `genesis.strategy.assurance.request/read/seal/review/calibrate/validate/renew/reexperiment/retire` 权限和职责分离。
- [x] T49 实现 scheduled/expiry/guardrail/drift/manual trigger、重复触发收敛和单 active cycle 约束。
- [x] T50 实现 expired release 的 fail-closed/limited-continuity policy、人工接管和不可由 AESE UI 延期规则。
- [x] T51 实现 renewed/reexperiment_required/retired AssuranceDecision、release effect、next action 和 review date。
- [x] T52 保证 cycle、finding、decision、audit、journal 和 Outbox 必要事务原子性及失败无部分写入。
- [x] T53 验证 tenant/RLS、越权/自批、陈旧 report、篡改 dataset、重复/并发 decision 和 Agent sole approval。
- [x] T54 两仓分别提交、记录 revision、部署并完成 contract/integration tests。

验收：IAOS 治理周期、权限、到期和 disposition；AESE 只计算观察、drift、校准和验证证据。

### A6 - Strategy Assurance Observatory（第 6-9 周）

- [x] T55 在 World Play 增加 Assurance Observatory：active release、review/expiry、cycle、cutoff/cursor 和 lineage。
- [x] T56 展示 data-quality findings、missingness/correction、五域 drift、process mismatch 和 policy pressure。
- [x] T57 展示 parent/candidate assumption diff、样本量、方法版本、8/4 split、holdout 和 replay comparison。
- [x] T58 展示 observation/inference/proposal/decision、World/IAOS/Knowledge owner 和统计/因果限制。
- [x] T59 提供 renew/reexperiment/retire impact preview 与独立审批状态，不提供直接改参数/Policy/release 入口。
- [x] T60 完成键盘、ARIA、焦点、移动端、非颜色 drift、单位/精度、样本量和虚拟/业务时间表达。
- [x] T61 验证大 lineage、correction、部分 unavailable、刷新/断线/cursor 恢复和错误空态。

验收：用户能从 AssuranceDecision 下钻到 ValidationReport、dataset、source ref 和旧 evidence，且不会混淆观察与推断。

### A7 - 全链验收与循环入口（第 8-10 周）

- [x] T62 从 M15 adopted terminal 运行一个 clean 12 周 cycle，分别验证 renewed 和 injected reexperiment_required/retired 路径。
- [x] T63 对账 World events、IAOS records/journal/Outbox、release、dataset/corrections、drift、calibration、replay 和 decision hashes。
- [x] T64 验证数据迟到、cursor gap、source outage、服务重启、重复点击、陈旧 version、并发 cycle 和分层 reset。
- [x] T65 验证 tenant-other、只读用户、Agent 自批、holdout 泄漏、删除不利窗口、自动调 Policy 和 real-production target 默认拒绝。
- [x] T66 运行 Go、Schema、PostgreSQL、IAOS modules、前端、Playwright、100 次 hash、容量和 M3-M15 回归。
- [x] T67 编写 M16 runbook/evidence，更新设计、计划、Roadmap、Code Map、Progress Log、Atlas、两仓 revision 和部署信息。
- [x] T68 发布下一周期入口：renewed 安排新 review；reexperiment_required 创建但不自动批准 M14 request；retired 进入 M15 退出/替代治理。

验收：只有 dataset/lineage/data quality/drift/calibration/holdout/replay/release effect 全部可对账时输出 `strategy_assurance_cycle_closed=true`。

## 4. 首版保障周期方向

A0 冻结前使用以下方向，不预设校准结果或 disposition：

| 维度 | 首版方向 |
| --- | --- |
| Active release | `STR-GENESIS-M15-RESILIENT@1.0.0` |
| Observation | 12 个虚拟周 canonical World + IAOS stable refs |
| Split | 前 8 周 calibration，后 4 周 single-use holdout |
| Drift domains | demand / supplier / equipment / quality / payment / policy-action |
| Calibration | 每域批准分布族和有界参数；最多一个 candidate |
| Re-evaluation | 新 ancestry M14 matrix + common random numbers |
| Disposition | renewed / reexperiment_required / retired |

cycle、dataset、candidate、validation 和 decision code 使用 `ASSURE-GENESIS-M16-*`。全部经营数据保持虚构；本计划不授权真实生产数据或 target。

## 5. 完成定义

- 12 周 dataset 的 source owner、cutoff/cursor、unit/precision、missing/late/correction 和 hash 完整。
- 数据质量未通过时 drift/calibration/renewal 失败关闭；各类 drift 有解释和 injected path。
- 8/4 split 无泄漏，candidate 唯一、有界、可复验，不改变 release 或旧 evidence。
- holdout 与新 M14 replay 保留完整样本、共同随机数、ancestry 和结论限制。
- scheduled/expiry/alert/manual trigger、single active cycle、tenant/RLS、职责分离、幂等和事务/Outbox 通过。
- Observatory、CLI/API、runbook/evidence、三视口、两仓 revision/部署、Atlas 和 M3-M15 回归完整。
- 输出 `strategy_assurance_cycle_closed=true` 与 renewed/reexperiment_required/retired，不存在未处理 finding、unknown drift 或未决 release effect。

## 6. 不纳入 M16

- 真实生产数据接入、个人/客户敏感数据、在线学习或自动模型/Policy 发布。
- 第二产品/客户/工厂、完整 S&OP、通用数据湖/特征平台或机器学习训练平台。
- 反复查看 holdout 调参、静默改历史 dataset/evidence 或选择性删除不利窗口。
- 用短窗口观察宣称统计因果、真实概率、永久最优或允许 Agent 独立续期/退役。

## 7. 并行与所有者规则

- A0 由 assurance contract owner 串行冻结；A1/A2 可在 spec 稳定后并行，但 dataset/hash 由单一 owner 封存。
- A3 只能消费 sealed calibration window；holdout owner 与 calibration owner 分离，候选冻结前不得访问 holdout。
- A4 依赖 candidate freeze；原 M14 evidence owner 负责验证不可变性，新 replay 使用独立 artifact namespace。
- A5 只能由新 IAOS worktree owner 开发，AESE/IAOS 提交、测试、日志和部署证据分别维护。
- A6 在 API/view model 稳定后开始；UI 不得提供泄漏 holdout、改参数或延长 release 的旁路。
- A7 串行核对两仓、三个 disposition、全部 correction/failure、revision、部署和 Atlas。
- 保留共享工作区中现有测试修改、截图删除和验收产物，不得覆盖或回滚。
