---
id: DES-010
title: M9 Project Genesis 企业成立与治理
date: 2026-07-22
status: completed
author: Codex + User
tags: [m9, genesis, incorporation, governance, budget, agent]
---

# M9 Project Genesis 企业成立与治理

## 1. 目标

M9 把 `hctm-genesis` 从“已经有现金、人员和设备的 tracer 世界”向前扩展到企业诞生阶段：华辰创始投资人从尚未成立企业的状态开始，完成出资决策、苏州制造公司注册、运营账户开户、首期资本到账、管理层任命、初始组织和首年预算批准，为 M10 工厂选址与建设立项建立合法、资金和授权前提。

首版只覆盖一个虚构投资主体、一个集团、一个苏州制造法人、CEO/CFO/工厂项目负责人三个岗位和一份启动预算。

## 2. 纵向业务链

```text
投资人形成创始意图
-> IAOS 创建并批准设立方案
-> AESE 安排外部登记活动并推进虚拟时间
-> 虚构监管机构批准法人登记
-> IAOS 登记法人和治理决议
-> 董事会任命 CEO / CFO / 工厂项目负责人
-> 银行开户并完成首期资本实缴
-> IAOS 建立初始组织、岗位、授权和预算
-> CEO / CFO 获知职责与资金边界
-> 验证企业具备 M10 工厂建设立项资格
```

## 3. 三态所有权

| 事实 | World State | IAOS Business State | Actor Knowledge |
| --- | --- | --- | --- |
| 企业是否获得登记批准 | Own：外部登记结果和生效时间 | 法人档案与登记记录 | 获授权角色收到批准 observation 后获知 |
| 投资人资金与公司实际到账 | Own：账户实际余额和资金移动 | 出资方案、收款/资金记录引用 | CFO 只知道已送达且有权查看的资金信息 |
| 股权和治理安排 | 世界规则保存生效后的最小权属事实 | Own：治理决议、角色、授权与审计 | 董事长、CEO、CFO 按职责可见 |
| 岗位是否有人实际接受任命 | Own：接受时间、可用性 | 任命决议、人员/Agent 岗位绑定 | 被任命者收到任命后获知 |
| 预算是否批准 | 不复制预算明细，只保存可支用上限引用 | Own：预算版本、审批和冻结状态 | CEO/CFO 按权限看到批准额度 |

法人登记、开户和资金到账是外部世界后果，不能由 IAOS 创建一条记录就自动视为完成。IAOS 的 committed outcome 只触发 AESE 安排后续世界活动；最终世界结果再以 observation 回到 IAOS。

## 4. M9 最小领域模型

### AESE 新增或扩展

- `FoundingCase`：发起人、目标法人、阶段、相关 intent/outcome。
- `LegalRegistration`：申请、受理、批准/拒绝、生效时间和虚构监管机构。
- `CashAccount`：owner stable ref、币种、实际余额和账户状态。
- `CapitalCommitment`：认缴、实缴、到期时间和资金来源。
- `Appointment`：岗位、受任主体、接受状态、生效时间和任期。
- `OperatingMandate`：角色可承诺金额、可批准类别和有效期的世界引用。
- `BudgetEnvelopeRef`：IAOS 预算 stable ref、批准额度和可支用状态，不复制预算明细。

### IAOS 最小管理对象

- `incorporation_case`
- `legal_entity`
- `governance_resolution`
- `organization_position`
- `position_assignment`
- `capital_contribution`
- `bank_account`
- `budget_envelope`
- `budget_approval`

优先使用 IAOS metadata entity、Process、Policy 和 Capability；只在现有平台合同无法满足原子性、权限或 Outbox 时做最小平台扩展。

## 5. 角色与 Agent

| 角色 | M9 目标 | 权限边界 |
| --- | --- | --- |
| 创始投资人/董事长 | 提出设立方案、提名管理层、批准资本与预算 | 不能伪造监管批准或银行到账 |
| CEO | 接受任命、确认三年目标、提交启动预算 | 不能批准自己的越权预算 |
| CFO | 核验到账、形成现金视图、审查预算守恒 | 不能把承诺金额当作实付或到账 |
| 工厂项目负责人 | 接受任命并准备 M10 立项输入 | M9 不创建厂址、合同或设备采购 |
| 外部监管机构/银行 | 由确定性世界策略产生批准、拒绝、开户和到账结果 | 不是 IAOS 用户，不调用内部 Capability |

Agent 首版采用版本化确定性岗位策略，不引入自由自治 LLM。人类接管任一岗位时使用相同 IAOS Capability、Policy、审批和审计合同。

## 6. 规则与不变量

- 法人未获 World registration approval 前，IAOS 不得把其用于签约或预算执行。
- 银行账户未实际 opened 前，资本不得到账。
- `investor opening cash = capital paid + fees paid + investor closing cash`。
- `company opening cash + capital received = company paid + company closing cash`。
- 认缴、实缴、可用现金和预算额度是四个不同数值，禁止互相替代。
- 预算批准只形成支出授权，不消耗现金；M10 实际付款才消耗现金。
- 岗位任命需要决议 committed outcome 和受任者接受，两者缺一不可。
- CEO/CFO 未生效或无相应 mandate 时，预算 intent 必须失败关闭。
- 所有金额使用十进制字符串、`CNY` 和显式 scale；时间使用 RFC 3339 与 `Asia/Shanghai`。

## 7. 世界事件与桥接类型

M9 新增 allowlist payload family：

```text
genesis.incorporation.requested.v1
genesis.registration.submitted.v1
genesis.registration.approved.v1
genesis.registration.rejected.v1
genesis.account.opening.requested.v1
genesis.account.opened.v1
genesis.capital.transfer.requested.v1
genesis.capital.received.v1
genesis.executive.appointment.requested.v1
genesis.executive.appointment.accepted.v1
genesis.budget.submitted.v1
genesis.budget.approved.v1
```

这些类型复用 DES-008 envelope、journal/cursor 和 committed outcome 语义。拒绝是合法世界结果或 IAOS process 状态，但不得伪装为 committed success；重试必须复用 idempotency key。

## 8. Pack 与兼容策略

- 将 `world-packs/hctm-genesis` 升级为新 minor version，增加 `campaigns/incorporation/`，不覆盖 M8 设备 tracer fixture。
- 增加显式 pre-incorporation 初态；M8 已有 10,000,000 CNY 不得被静默解释为公司现金，必须迁移为有 owner 的账户和可追溯资本事件。
- M7 `scenario-packs/hctm/order-expedite-01` 继续独立运行；M9 不要求旧故事从企业成立阶段开始。
- M10 只能消费 M9 机器可验证的 `plant_project_eligible=true` 输出，不能凭 UI 状态开始。

## 9. 非目标

- 真实工商、税务或银行接口。
- 完整公司法、复杂股权结构、融资轮次、贷款、利息和外汇。
- 完整总账、会计科目、税务、工资或集团合并报表。
- 工厂选址、土地协议、建设合同、设备采购和招聘流水线。
- 自动批准高风险预算或让 Agent 伪造外部机构结果。

## 10. 完成标准

- 单一 world run 可从 pre-incorporation 确定性推进到 `plant_project_eligible`。
- 法人、账户、资金、任命和预算的 World/IAOS/Knowledge 三态可分别观察并对账。
- 至少包含一次信息延迟或登记差异，以及通过 observation/intent/outcome 关闭的证据链。
- 人类与 Agent 对 CEO/CFO 岗位使用同一治理能力且越权失败关闭。
- 资金、任命和预算不变量通过，重复执行无重复资金或业务副作用。
- M7 和 M8 全链回归通过，runbook/evidence 与两仓 revision 完整。
