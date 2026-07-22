---
id: PLAN-M9-001
title: M9 Project Genesis 企业成立与治理实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m9, genesis, incorporation, governance, budget, agent]
---

# M9 Project Genesis 企业成立与治理实施计划

## 1. 交付目标与工期

用 **4 到 6 周**交付第一个真正从“企业尚不存在”开始的生命周期纵向闭环：

> 创始投资人通过 IAOS 提出并批准华辰苏州制造公司的设立方案；AESE 推进登记、开户和资本到账的客观世界过程；董事会任命 CEO、CFO 和工厂项目负责人；IAOS 批准初始组织与启动预算，最终形成机器可验证的 M10 工厂建设立项资格。

首版固定为单 tenant、单 world run、单法人、单币种 CNY、单预算版本和三个管理岗位。

## 2. 前置基线

- [x] G1 M8 F0-F5、World Store、确定性内核、World Bridge 和 World Play 已完成。
- [x] G2 ADR-004 三态所有权与 DES-008 bridge 合同继续有效。
- [x] G3 DES-010 固定 M9 业务链、对象所有权、角色、事件族和非目标。
- [x] G4 固定虚构资金基线、登记周期/费用、初始岗位和预算额度，并由 invariants fixture 验证。
- [x] G5 完成 IAOS metadata/Process/Policy/Capability gap audit，冻结两仓实施顺序。

G1-G5 与 I0-I5 已完成；当前没有 active 主实施计划。

## 3. 实施切片

### I0 - 业务基线与机器合同（第 1 周）

- [x] T1 定义投资主体、法人、登记、账户、资本、任命、mandate 和 budget envelope 的稳定编码与所有权矩阵。
- [x] T2 固定虚构资金、登记费用/周期、资本到账、首年预算和三年目标的基线，不使用未注明来源的展示数字。
- [x] T3 扩展 World JSON Schema、Go 类型、strict parser、fixture 和 canonical hash。
- [x] T4 为 M9 observation/intent/outcome payload family 建立 strict schema registry 和破损 fixture。
- [x] T5 定义 `pre_incorporation -> registering -> registered -> capitalizing -> organizing -> budgeted -> plant_project_eligible` 状态机。
- [x] T6 更新能力缺口台账、System Atlas 节点/依赖和 M9 两仓交付矩阵。

验收：同一业务概念只有一个 owner；金额、时间、stable ref、权限和状态转换合同可离线验证。

### I1 - AESE 成立世界与经济规则（第 1-2 周）

- [x] T7 扩展 `internal/genesis` 与 `internal/rules`，实现登记、开户、资本转移、任命接受和预算资格的纯函数 reducer。
- [x] T8 实现虚构监管机构和银行的版本化确定性策略，包括批准、拒绝、处理时间和费用。
- [x] T9 实现 investor/company 独立 CashAccount、CapitalCommitment 和资金流水，禁止无 owner 的现金。
- [x] T10 实现资金、登记前置、任命双确认、mandate 和预算不变量。
- [x] T11 实现 snapshot/replay/reset 和 100 次确定性回归，覆盖拒绝、超时、重复与状态倒退。
- [x] T12 将 `world-packs/hctm-genesis` 升级并增加 `campaigns/incorporation/`，保留 M8 tracer。

验收：完全离线即可从 pre-incorporation 推进或回放到可验证终态；非法资金、任命或预算路径失败关闭。

### I2 - IAOS 成立与治理能力（第 2-4 周）

- [x] T13 按 IAOS `AGENTS.md` 建立独立 branch/worktree，审计现有 legal entity、organization、Process、Policy、Capability 和预算对象。
- [x] T14 以 metadata/config package 提供最小管理对象；仅在无法满足原子性或治理时增加平台代码。
- [x] T15 实现设立方案、法人登记、管理层任命、资本核验和预算批准的 allowlist Capability/Process。
- [x] T16 固定 `genesis.incorporation.execute`、`genesis.governance.appoint`、`genesis.budget.submit/approve` 权限资源及岗位映射。
- [x] T17 让 intent、业务记录、audit、committed outcome、journal 和 Outbox 在同一事务边界提交。
- [x] T18 验证 tenant RLS、权限不足、越权自批、幂等冲突、回滚无 outcome、重复资金 no-op 和并发审批。
- [x] T19 分别提交 IAOS 和 AESE 合同更新，记录依赖 revision 后再部署联调。

验收：AESE 不直写 IAOS DB；外部登记/银行结果不能由 IAOS 伪造；所有正式管理动作可审计且失败无部分写入。

### I3 - CEO/CFO 统一岗位运行（第 3-4 周）

- [x] T20 建立创始投资人、CEO、CFO、工厂项目负责人和外部机构 actor/policy fixture。
- [x] T21 实现 appointment observation、角色接受、Knowledge 可见范围、信息延迟和 supersedes。
- [x] T22 实现 CEO 三年目标与启动预算提交、CFO 资金核验与预算审查的确定性岗位策略。
- [x] T23 实现人类岗位接管；人类和 Agent 使用相同 Capability、Policy、权限与审计，不提供旁路按钮。
- [x] T24 验证未获知、未任命、未接受、越权、过期 mandate 和自我审批均失败关闭。

验收：同一 world state 下，不同角色因知识和权限不同得到不同可解释动作；Agent 输出只形成 intent。

### I4 - Genesis Incorporation Play（第 4-5 周）

- [x] T25 在 World Play 增加“企业成立”campaign 入口，不新增独立前端应用。
- [x] T26 展示投资人、法人状态、资金账户、治理岗位、预算 envelope 和 M10 资格。
- [x] T27 增加虚拟时间、阶段 stepper、World/IAOS/Knowledge 对照和因果时间线。
- [x] T28 实现设立、任命、预算等动作的岗位接管、确认、权限错误、no-op 和恢复状态。
- [x] T29 保持金额来源、单位、World/IAOS owner 和虚拟/真实时间可见，不把预算画成现金。

验收：非研发用户可只通过页面完成或观察成立闭环，并能解释每一步为什么允许、等待或失败。

### I5 - 全链验收与交付（第 5-6 周）

- [x] T30 从 clean pre-incorporation 执行 Agent 路径和人类接管路径，结果 state hash 与业务结果一致。
- [x] T31 对账 World event、Knowledge、IAOS journal、Capability/Process、Audit、Outbox、资金和预算引用。
- [x] T32 验证断线、刷新、AESE/IAOS 重启、SSE 丢失、重复点击、陈旧 cursor 和恢复。
- [x] T33 验证只读用户、tenant-other、越权审批、重复资本到账和回滚无部分写入。
- [x] T34 运行 Go、Schema、数据库、前端、Playwright、M7/M8 回归及三目标视口验收。
- [x] T35 编写 M9 runbook/evidence，更新设计、计划、Roadmap、Code Map、Progress Log、Atlas 和两仓 revision。

验收：终态满足 `plant_project_eligible=true`；所有完成条件有机器证据，UI 完成不替代业务或集成验收。

## 4. 建议交付编码

首版稳定编码在 I0 最终冻结，建议起点：

```text
INV-HCTM-FOUNDERS
HCTM-GROUP
HCTM-SZ-MFG
POS-HCTM-CEO
POS-HCTM-CFO
POS-HCTM-SZ-PROJECT-DIRECTOR
BUDGET-HCTM-SZ-Y1
BANK-HCTM-SZ-OPERATING-01
```

这些是虚构业务编码，不携带真实工商、银行或个人信息。

## 5. 完成定义

- M9 world pack 可验证、初始化、推进、重放、快照恢复和 reset。
- 法人登记、开户、资本到账、任命和预算三态分离且有完整因果链。
- 资金守恒区分投资人现金、公司现金、资本承诺、实缴、预算与实际支出。
- CEO/CFO 人类与 Agent 使用相同 IAOS 治理能力，权限和知识边界可回归。
- IAOS journal/Outbox 与业务提交原子，回滚不产生 committed outcome。
- tenant、幂等、并发、断线、重启和重复执行测试通过。
- M7 订单故事与 M8 设备 tracer 零回归。
- `plant_project_eligible` 只在法人、资金、岗位、mandate 和预算全部满足时为 true。
- runbook、evidence、两仓提交/部署和 Atlas 登记完整。

## 6. 不纳入 M9

- 工厂选址、土地/租赁协议和建设项目执行；属于 M10。
- 设备采购、安装调试、招聘培训和投产门；属于 M11。
- RFQ、APQP、试生产和 PPAP；属于 M12。
- 首批 O2D、开票、回款和实际成本；属于 M13。
- checkpoint/fork 和 A/B；属于 M14。
- 完整总账、税务、真实监管/银行集成和复杂融资。

## 7. 并行与所有者规则

- I0/I1 由 AESE owner 先完成合同和离线 tracer。
- I2 只能在 I0 payload/state contracts 冻结后，由独立 IAOS worktree owner 开始。
- I3 可在 I2 权限和 Capability 名称冻结后与 IAOS 测试有限并行。
- I4 在 API/view model 稳定后开始；I5 必须串行收口两仓证据。
- 并行 agent 必须使用明确子计划、owner 和不重叠 worktree；不得在共享 dirty worktree 覆盖其他人的测试或截图改动。
