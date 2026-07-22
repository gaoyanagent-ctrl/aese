---
id: DES-017
title: M16 持续策略保障与假设校准
date: 2026-07-22
status: completed
author: Codex + User
tags: [m16, strategy, assurance, drift, calibration]
---

# M16 持续策略保障与假设校准

## 1. 决策摘要

M16 消费 M15 主路径已采纳的 `STR-GENESIS-M15-RESILIENT@1.0.0`，建立一个可重复的策略保障周期：按冻结 cutoff 收集 canonical 经营观察，验证数据质量和 lineage，区分输入/过程/结果/策略行为漂移，提出有界假设校准，在独立 holdout 和重跑实验中验证，最后决定维持、重新实验或退役策略。

机器终态固定为：

```text
strategy_assurance_cycle_closed=true
disposition=renewed | reexperiment_required | retired
```

M16 不在线自学习、不自动修改 StrategyRelease、不覆写 M14 EvidenceBundle，也不把 12 周虚构样本包装成真实统计规律。`reexperiment_required` 与 `retired` 都是合法完成结果。

## 2. 为什么是 M16

M15 已证明 evidence-to-action 可以受治理关闭，但 adopted 只对当时 scope/version 有效。持续运行后仍需回答：

- M14 的需求、供应、设备、质量和付款假设是否仍覆盖当前 canonical World。
- KPI 变化来自数据迟到/缺失、外生输入变化、过程能力变化，还是策略动作本身。
- 策略是否频繁贴近 guardrail、依赖人工接管或在支持域外运行。
- 新观察能否形成一个更好的假设版本，而不发生数据泄漏和事后挑选。
- 到复审日期时应续期、返回 M14 重新实验，还是受治理退役。

因此 M16 优先补齐 evidence-to-action 之后的 learning/assurance loop，而不是立即扩展第二产品、客户或工厂。

## 3. 首版范围

首版固定：

- `tenant-hctm`、苏州基地、电池冷却板 A 线、单客户/产品和已采纳 M15 release。
- 12 个虚拟周 canonical observation window：前 8 周 calibration、后 4 周 holdout；具体 cutoff 在 A0 冻结。
- 需求、供应、设备、质量、付款五个 M14 命名域和 M15 policy-action/guardrail 域。
- 一个 immutable AssuranceDataset、一个原假设版本、最多一个 CalibrationCandidate 和一个 AssuranceDecision。
- 保留 M14 common-random-number matrix/replay 作为再评估工具，但不得改写原 60-run evidence。

不包含：

- 真实生产数据接入、个人/客户敏感数据、在线学习或模型自动部署。
- 自动调整 Policy/StrategyRelease、自动启动新 pilot 或 Agent sole approval。
- 第二产品/客户/工厂、完整 S&OP、通用特征平台或机器学习训练平台。
- 以 12 周单一虚构样本宣称长期因果、真实概率或永久最优。

## 4. 核心合同

| 对象 | 所有者 | 关键字段 |
| --- | --- | --- |
| AssuranceCycle | IAOS 治理；AESE 执行观察/分析 | cycle code、release、scope、cutoff、windows、status |
| ObservationSpec | 共同版本合同 | metric/event refs、owner、unit、precision、lateness/freshness rule |
| AssuranceDataset | AESE artifact；IAOS 保存 hash/ref | cutoff、World event refs、IAOS record/journal refs、missingness、hash |
| DataQualityFinding | 源 owner + IAOS incident | kind、severity、affected range、resolution、decision impact |
| DriftAssessment | AESE | domain、method/version、baseline/current、threshold、support status |
| CalibrationCandidate | AESE proposal | parent assumption、bounded parameter diff、fit window、hash、limitations |
| ValidationReport | AESE evidence | holdout、replay matrix、metrics、failure list、comparison、hash |
| AssuranceDecision | IAOS | disposition、reason、approvers、release effect、next action、review date |

Dataset 只保存仿真观察、允许的度量和 IAOS stable refs/hash，不复制 IAOS 订单、库存、财务或权限数据库。数量和金额继续使用显式精度；统计输出固定算法、版本、排序、舍入和 canonical encoding。

## 5. Observation 与 as-of 语义

每个 AssuranceDataset 必须绑定：

- canonical World run/branch、active release、pack/rules/actor-policy version。
- IAOS tenant、journal cursor range、business stable refs 和 committed-outcome range。
- `window_start`、`window_end`、`cutoff_at`、`built_at` 和 `Asia/Shanghai`。
- 每个指标的数据 owner、单位、精度、freshness、允许 lateness 和 correction policy。
- 缺失、迟到、重复、冲突、修订和排除记录。

Dataset 是 as-of cutoff 的不可变快照。cutoff 后到达的事实不能静默回填旧 hash；必须产生 correction set 和新 dataset version。Actor Knowledge、UI 状态和通知不能替代 World/IAOS source facts。

## 6. Drift 分类与决策顺序

检测顺序固定，避免把坏数据误判成业务变化：

1. Data quality/freshness drift：缺失、迟到、重复、单位/版本/owner 不符。
2. Input drift：需求、供应、设备、质量、付款输入相对 M14 assumption support 的变化。
3. Process/outcome mismatch：相同输入区间下，产能、质量、交付、现金或成本结果与规则预测的系统差异。
4. Policy-action drift：动作频率、人工接管、guardrail pressure、churn 和越界尝试变化。

数据质量未通过时，后续 drift 默认 `indeterminate`，不能生成校准结论。阈值、检验方法和多重比较规则必须在看当前窗口结果前冻结；首版使用可解释的 bounded rules 和版本化统计，不引入黑盒异常检测。

## 7. 校准与防泄漏

CalibrationCandidate 只能修改 DES-015 allowlist 中的外生假设参数，不能修改业务事实、hard constraint、KPI 定义或已采纳 release：

- 使用前 8 周 calibration window 拟合批准的分布族和有界参数。
- 后 4 周 holdout 在候选冻结后只评估一次，不回流调参。
- 保留 parent assumption、参数 diff、样本量、缺失处理、估计方法和适用限制。
- 参数越界、样本不足、holdout 不完整、结果不稳定或改进只来自一个指标时失败关闭。
- 不允许通过反复查看 holdout、删除不利周或更换指标获得通过结果。

校准成功只表示新 assumption candidate 更适合该观察窗口；它不自动取代 M14 assumption，也不自动证明当前策略应变更。

## 8. 再实验与策略保障

M16 使用两个独立检查：

1. Holdout validation：比较 parent 与 calibration candidate 对未见 4 周数据的覆盖和误差。
2. M14 replay：在新 assumption candidate 上重建版本化矩阵，以共同随机数比较 baseline/lean/resilient 和当前 adopted release 对应策略。

原 M14 matrix、seed、run 和 EvidenceBundle 保持不可变；新结果使用新 experiment/evidence code 和祖先引用。任何失败/cancelled run 进入 completeness gate，不得过滤。replay 仍是 simulation evidence，不是 causal proof。

## 9. AssuranceDecision

决策门定义：

- `renewed`：数据完整，当前假设仍受支持或校准不改变策略结论，release 按明确期限续期。
- `reexperiment_required`：输入/过程/策略支持发生实质变化，或校准后策略比较不再稳定；保持/暂停当前 release 由 SafetyEnvelope 决定，并回到 M14/M15 新周期。
- `retired`：release 已越出允许 scope、持续触发 hard risk 或已被受治理替代；退役沿用 M15 rollback/commitment/compensation 语义。

AssuranceDecision 必须记录 evidence、data-quality findings、drift、calibration/validation、未确定性、approvers、release effect、下一动作和复审日期。Agent 可准备报告，不能独立续期或退役。

## 10. 持续运行边界

M16 实现的是“可重复保障周期”，不是无人值守 daemon：

- 周期由计划、到期、guardrail/drift alert 或人工 request 触发。
- 每个周期冻结 spec/cutoff，并生成独立 dataset/evidence/decision。
- 监控可持续追加 observation，但只有关闭后的 AssuranceDecision 能改变 release review 状态。
- 监控故障、cursor gap、source unavailable 或 dataset mismatch 自动 pause review，不自动续期。
- expired release 的业务行为由 IAOS Policy 决定 fail-closed/limited continuity，不能由 AESE UI 擅自延长。

## 11. 三态与 Bridge

AESE World 拥有外生事件、物理/经济结果和 World observations；IAOS 拥有业务记录、active release、AssuranceCycle、data-quality incident、审批和 AssuranceDecision；Actor Knowledge 只包含角色在当时被允许看到的数据质量、drift 和保障摘要。

Bridge 增加 allowlist payload family：

```text
genesis.strategy.assurance.requested.v1
genesis.strategy.dataset.sealed.v1
genesis.strategy.drift.detected.v1
genesis.strategy.calibration.proposed.v1
genesis.strategy.validation.completed.v1
genesis.strategy.assurance.decided.v1
genesis.strategy.assurance.closed.v1
```

通知仍不是事实。Dataset、ValidationReport 和 Decision 必须经 exact hash、expected version、tenant、journal/cursor 和幂等验证。

## 12. Strategy Assurance Observatory

World Play 增加 Strategy Assurance Observatory，展示 active release/expiry、observation lineage、cutoff/cursor、数据质量、五域 drift、calibration diff、holdout、replay comparison、限制和 AssuranceDecision。

页面必须区分 source observation、推断、校准 proposal 和已批准 decision；不提供改写 dataset、反复查看 holdout 调参、删除不利窗口或直接修改 StrategyRelease 的入口。所有图表提供非颜色表达、单位、样本量、missingness 和方法版本。

## 13. 完成标准

- 从 adopted M15 release 建立一个 12 周 assurance cycle，dataset as-of/hash/lineage 可重建。
- World/IAOS owner、cursor、cutoff、迟到/修订/缺失和 correction set 合同通过。
- 数据质量先于 drift；五个外生域和 policy-action 域均有可解释 assessment 与失败路径。
- calibration/holdout 严格时间隔离，参数有界，方法/舍入/version 固定且 100 次 hash 一致。
- 新 M14 replay 保留原 evidence，以共同随机数、完整 run manifest 和新 ancestry 形成 ValidationReport。
- IAOS tenant/RLS、职责分离、幂等、事务/Outbox、过期/续期/退役和失败无部分写入通过。
- Observatory、CLI/API、runbook/evidence、三视口、两仓 revision/部署和 M3-M15 回归完整。
- 输出 `strategy_assurance_cycle_closed=true` 与明确 disposition，不存在未处理数据质量问题、未知 drift 或未决 release effect。

## 14. 后续边界

M16 完成后，若 `renewed`，M17 可考虑滚动 S&OP 或第二产品/客户扩展；若 `reexperiment_required`，返回新的 M14/M15 evidence-to-action 周期；若 `retired`，先完成替代/退出治理。不得绕过 disposition 直接扩展范围。
