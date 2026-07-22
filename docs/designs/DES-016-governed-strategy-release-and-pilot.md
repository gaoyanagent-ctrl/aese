---
id: DES-016
title: M15 受治理策略发布与经营试点
date: 2026-07-22
status: completed
author: Codex + User
tags: [m15, strategy, governance, pilot, rollback]
---

# M15 受治理策略发布与经营试点

## 1. 决策摘要

M15 把 M14 的 `proposed_not_applied` 推荐推进为一条可审议、可拒绝、可分阶段试运行、可停止和可追责的策略变更链。它不把实验排名直接变成正式 Policy，也不把一次试点表现包装成因果证明。

M15 的产品终态为：

> 人或 Agent 可以引用不可变 EvidenceBundle 提交策略变更请求；独立角色审议目标、适用范围、风险包络和退出方案；获批版本先 shadow、再在新的 Genesis canonical operating cycle 中有限试点；监控触发采纳、拒绝或回滚，并保留完整业务后果和补偿记录。

机器终态固定为 `strategy_change_cycle_closed=true`，同时必须给出 `disposition=adopted|rejected|rolled_back`。里程碑完成不要求策略一定获胜；只有决策过程和结果均可解释、可重放、无未处理风险时才能关闭。

## 2. 为什么是 M15

M14 已证明在虚构假设和 60-run 矩阵中可公平比较策略，但仍存在最后一公里：

- EvidenceBundle 完整不等于组织已经审议或批准变更。
- 模拟 PolicyVariant 不等于 IAOS 中可授权、可版本化和可撤销的业务 Policy。
- 推荐在隔离分支表现良好，不等于进入 canonical operation 后仍在适用范围内。
- 回滚只能停止未来决策，不能删除已经提交的采购、生产、库存、交付和现金后果。
- 单个 pilot 没有随机对照，不足以宣称真实因果提升。

因此 M15 优先建设 decision-to-action 的治理、发布、shadow、pilot、监控和退出闭环，而不是继续增加实验参数或扩大企业范围。

## 3. 首版范围

首版固定：

- 单租户、苏州基地、电池冷却板 A 线、单客户和单产品。
- 从 M14 EvidenceBundle 的 Pareto 候选中审议一个 StrategyCandidate；不在设计阶段预设 winner。
- 策略范围只覆盖安全库存/双源、维护/加班和现金保护的 allowlist 阈值与决策动作。
- 4 个虚拟周 shadow + 最多 4 个虚拟周受控 pilot；周期数在 R0 冻结。
- 一个 immutable StrategyRelease、一个 canonical pilot cycle、一个明确 rollback target。
- OTIF、质量、现金、毛利、库存、加班和策略动作合规 guardrail。

不包含：

- 真实生产租户投放、无人审批自动发布或 Agent 自批。
- 第二客户/产品/工厂、完整 S&OP、价格优化、组织裁员或资本投资自动化。
- 用 pilot 单样本宣称统计因果、真实预测精度或永久最优。
- 删除或改写已提交业务事实来伪造回滚。

## 4. 核心合同

| 对象 | 所有者 | 关键字段 |
| --- | --- | --- |
| StrategyCandidate | AESE evidence projection；IAOS 治理登记 | evidence hash、policy hash、assumptions、Pareto refs、limitations |
| StrategyChangeRequest | IAOS | requester、target scope、objective、risk class、approvers、status |
| StrategyRelease | IAOS；AESE 消费只读 manifest | release code/version/hash、effective window、thresholds、action allowlist |
| SafetyEnvelope | IAOS Policy + AESE invariant | hard limits、warning/stop thresholds、owners、measurement source |
| ShadowRun | AESE | canonical observation refs、candidate decisions、zero committed action、metrics |
| PilotCycle | AESE World + IAOS business runtime | checkpoint、release、window、business correlations、status |
| GuardrailBreach | World 事实 + IAOS incident/decision | observed value、threshold、severity、action、evidence refs |
| RollbackPlan | IAOS | prior release、stop boundary、open commitments、compensating actions、approvers |
| AdoptionDecision | IAOS | disposition、review evidence、effective version、conditions、audit refs |

所有对象使用稳定编码和 canonical hash。StrategyRelease 必须绑定一个 M14 EvidenceBundle、准确的目标 checkpoint、规则/Policy 版本、适用窗口、权限、监控指标、停止条件和回滚目标；任一引用陈旧或不兼容时失败关闭。

## 5. 状态机与决策门

```text
candidate
  -> evidence_verified
  -> under_review
  -> rejected
  | approved_for_shadow
      -> shadow_running
      -> shadow_failed | approved_for_pilot
          -> pilot_running
          -> paused
          -> adopted | rejected | rolled_back
              -> strategy_change_cycle_closed
```

提议者不能审批自己的变更；Agent 可以准备材料和提交 intent，但不能成为唯一审批人。涉及库存/供应、产能/维护和现金保护时，分别要求业务 owner，并由风险/财务角色完成独立批准。shadow 到 pilot、pilot 到 adopted 是两个不同决策门，前一阶段批准不能自动穿透下一阶段。

## 6. StrategyRelease 与编译边界

M14 PolicyVariant 是模拟输入；M15 StrategyRelease 是受治理业务配置。发布过程必须显式编译和校验：

- 将 allowlist 阈值和动作映射到 IAOS Policy/Capability/Process，不允许任意代码或任意 API。
- 将 World 所需规则引用编译为只读 release manifest；IAOS committed outcome 仍是动作唯一入口。
- 生成 semantic diff：当前 release、候选 release、默认值、单位、上下界、权限和影响对象。
- preflight 计算受影响的 policy、未结承诺、预计动作量、现金/库存上限和 rollback 可行性。
- 签名/hash 不符、未知字段、越界阈值、缺 owner、缺监控或不可回滚动作失败关闭。

CLI 和 UI 默认 dry-run；创建、批准、激活、暂停、采纳和回滚均要求明确目标环境、tenant、release hash、expected version 和幂等键。

## 7. Shadow 语义

shadow 使用 canonical cycle 的同一 observation 和 Actor Knowledge，计算“候选策略会建议什么”，但：

- 不提交业务 intent，不产生 committed outcome，不调用写 Capability。
- 不预占库存、产能、预算或现金，不改变客户/供应商/设备 World 行为。
- 将 candidate decision 与 active-policy decision 并排记录，保留输入 freshness、权限可见性和 reason code。
- 检查候选是否频繁越界、依赖不可见事实、产生不稳定动作或超出 M14 assumption support。

shadow 只能证明运行兼容性和建议质量，不能证明实际业务效果。只有完整窗口、零写入、全部 guardrail 可计算且独立审议通过时，才能申请 pilot。

## 8. 受控 Pilot 语义

pilot 在新的 canonical Genesis operating cycle 中运行，而不是 M14 隔离实验分支。候选策略只在明确时间窗、对象范围和动作 allowlist 内生效；所有业务动作仍走 IAOS Capability/Process/Policy/Decision、RLS、幂等、事务和 Outbox。

AESE 继续拥有需求、供应、设备、质量、运输、银行和物理/经济 consequence。IAOS Policy 决定是否允许 intent，不能直接创造库存、产能、交付或现金。pilot 必须保留 active release、candidate release、actor、observation、intent、committed outcome 和 World consequence 的因果链。

试点指标与 M14 假设范围并行监控。输入漂移超出支持域、关键 evidence 过期、硬 guardrail 违反、数据不完整或监控失效时自动 pause，不自动回退或继续。

## 9. Guardrail、停止和回滚

Guardrail 分三层：

- Hard stop：质量/权限/数量金额守恒、最低现金、负库存、未批准动作等不可突破约束。
- Pause and review：OTIF、积压、库存、加班、设备/供应风险或 assumption drift 达到审议阈值。
- Informational：成本、毛利、周转和建议采纳率等趋势指标。

回滚的语义是停止候选 release 对未来决策生效，并原子恢复已批准的 prior release。已经提交的采购单、生产任务、库存移动、发运、发票和现金不能删除；必须生成 open-commitment ledger，由受治理的取消、变更、消耗或其他 compensating action 处理。紧急 kill switch 只暂停未来动作，不能绕过事后审批、审计和补偿。

## 10. 评价与采纳

Adoption Review 必须同时展示：

- 原 M14 evidence、假设和限制是否仍适用。
- shadow 的零写入、decision diff 和兼容性结果。
- pilot 的实际 World/IAOS/Knowledge 因果链、guardrail 和 open commitments。
- baseline forecast、实际结果和偏差，但明确它不是随机对照因果估计。
- 数据缺口、异常、人工接管和任何补偿动作。

`adopted` 表示在当前 scope/version 下批准继续使用，不表示永久最优；必须带复审日期、owner 和失效条件。`rejected` 或 `rolled_back` 同样是有效完成结果，并进入后续实验/规则改进输入。

## 11. 三态与 Bridge 合同

AESE World 拥有 canonical pilot 时间、外生事件、物理/经济后果、guardrail observation 和试点事实。IAOS 拥有 ChangeRequest、StrategyRelease、审批、Policy/Capability、incident、rollback、compensation 和 AdoptionDecision。Actor Knowledge 只包含角色在当时已获准看到的 evidence、观察和告警。

Bridge 增加严格 allowlist payload family：

```text
genesis.strategy.change.requested.v1
genesis.strategy.shadow.approved.v1
genesis.strategy.pilot.approved.v1
genesis.strategy.release.activated.v1
genesis.strategy.guardrail.breached.v1
genesis.strategy.release.paused.v1
genesis.strategy.rollback.committed.v1
genesis.strategy.adoption.decided.v1
genesis.strategy.change.closed.v1
```

所有 outcome 继续遵守 DES-008 的 journal/cursor 恢复合同。通知不是事实；任何 release activation 必须验证 IAOS committed outcome、expected version 和 exact release hash。

## 12. Strategy Control Room

World Play 增加 Strategy Control Room，展示 EvidenceBundle、semantic diff、审批职责、shadow decision diff、pilot 因果时间线、guardrail、assumption drift、open commitments、暂停/回滚和最终 disposition。

页面不提供跳过 shadow、合并审批、改写 evidence、删除不利周期或直接修改 World 状态的入口。危险动作必须显示影响预览、exact release、目标环境、不可逆后果和确认步骤；移动端同样不能弱化这些信息。

## 13. 完成标准

- 一个 M14 candidate 完成 evidence verify、独立审议、shadow 和 bounded pilot，或在任一门被可解释地拒绝。
- StrategyRelease、SafetyEnvelope、RollbackPlan 和 AdoptionDecision 全部版本化、hash 绑定和可审计。
- shadow 完整窗口零业务写入，candidate/active decision diff 可下钻。
- pilot 的每个动作通过 IAOS 正式治理链，World consequence、Knowledge 和业务记录可对账。
- hard stop、pause、kill switch、rollback 和 compensation 至少各有自动化失败路径；回滚不删除既成事实。
- tenant/RLS、职责分离、陈旧版本、重复/并发动作、失败无部分写入和恢复通过。
- Strategy Control Room、CLI/API、runbook/evidence、三视口、两仓 revision/部署和 M3-M14 回归完整。
- 最终输出 `strategy_change_cycle_closed=true` 和明确 disposition；不存在未处理 breach、未知 release 或未对账 commitment。

## 14. 后续边界

M15 关闭决策到行动链后，M16 再根据 disposition 决定：持续策略复审和漂移校准、扩展第二产品/客户、建立滚动 S&OP，或返回 M14 重新实验。M15 不预先承诺任何一个方向。
