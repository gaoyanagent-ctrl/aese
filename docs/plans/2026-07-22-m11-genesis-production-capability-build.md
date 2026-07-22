---
id: PLAN-M11-001
title: M11 Project Genesis 生产能力建设实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m11, genesis, equipment, workforce, capability]
---

# M11 Project Genesis 生产能力建设实施计划

## 1. 交付目标与工期

用 **6 到 8 周**把 M10 的 `capability_build_eligible=true` 推进为：

> 华辰苏州制造公司在可证明的资金来源、空间、公用工程、供应周期和岗位技能约束下，完成电池冷却板 A 线通用设备/实验室/仓储能力与一班制核心团队建设，输出 `industrialization_eligible=true`。

首版固定单法人、单场地、单产线能力集合、单币种 CNY、一班制最小团队，以及一条检漏设备调试失败/整改 tracer。终态只允许进入 M12 产品工业化，不代表试生产、PPAP 或 SOP 完成。

## 2. 前置门

- [x] G1 M10 已完成并机器输出 `capability_build_eligible=true`。
- [x] G2 ADR-004、DES-008、DES-009 与 M8-M10 的 World/Bridge/Knowledge/资金边界继续有效。
- [x] G3 DES-012 固定 M11 三态所有权、设备/人员范围、资金约束和 M12 边界。
- [x] G4 冻结能力需求、设备清单、三个 acquisition option、报价/交期、空间/utility、验收标准和异常基线。
- [x] G5 冻结组织、岗位、最低 headcount、技能矩阵、培训/认证、班次和职责隔离基线。
- [x] G6 冻结资本催缴、采购/租赁组合、新 CAPEX/headcount envelope、设施尾款、工资准备金和最低现金缓冲。
- [x] G7 完成 IAOS procurement/asset/org/recruitment/training/qualification gap audit，冻结两仓 payload、权限和交付顺序。

G4-G7 完成前只允许 Schema、fixture、离线规则和只读 IAOS gap audit；不得开发 IAOS 写端点、签署可执行订单或形成 offer。

## 3. 实施切片

### C0 - 能力、设备、人员与资金机器合同（第 1 周）

- [x] T1 定义 CapabilityRequirement、EquipmentOption/Unit、Commissioning/Calibration、Position、Candidate、Worker、Skill、Qualification、Shift 和 CapabilityGate stable code/owner。
- [x] T2 固定 A 线通用能力边界、设备与实验室清单、空间落位、公用工程需求和验收证据；不得发布产品 BOM/routing。
- [x] T3 固定最小组织、岗位、headcount、技能矩阵、培训、认证、替补和一班制基线。
- [x] T4 固定资本到账、CAPEX/headcount、采购/租赁、设施尾款、付款节点、工资准备金和现金缓冲模型。
- [x] T5 扩展 JSON Schema、Go 类型、strict parser、canonical hash、payload registry 及合法/破损 fixture。
- [x] T6 定义 `eligible -> funded -> sourcing -> ordered -> installing -> commissioning -> staffing -> capability_acceptance -> industrialization_eligible` 状态机和跨 campaign terminal contract。
- [x] T7 完成 IAOS gap ledger、权限/职责冲突矩阵、Atlas 依赖和两仓交付矩阵。

验收：设备、空间、人员、资金、三态 owner 和 M12 边界无歧义，G4-G7 关闭，所有金额/单位/技能/证据可离线验证。

### C1 - 资金与受治理设备采购（第 1-2 周）

- [x] T8 实现剩余认缴资本催缴 -> 银行实际到账 observation -> IAOS 入账 outcome 的资金链。
- [x] T9 实现设施尾款、CAPEX、headcount、工资准备金和最低现金缓冲的联合 hard constraints。
- [x] T10 为关键设备建立至少采购、融资租赁、供应商账期三个组合 option，并输出成本、现金、交期、能力和风险解释。
- [x] T11 实现采购建议、供应商比选、预算审批、订单/租赁批准和付款节点的离线 bridge tracer。
- [x] T12 验证认缴未到账、超预算、侵占工资准备金、自批、重复订单和无验收付款全部失败关闭。

验收：同一输入产生确定性 acquisition recommendation，但签约和付款只能由 IAOS 受治理 outcome 推进。

### C2 - AESE 设备、实验室与仓储能力世界（第 2-4 周）

- [x] T13 新增 `internal/capabilitybuild/`，实现供应商制造、运输、收货、落位、安装、调试、校准、安全和能力验收 reducer。
- [x] T14 绑定 M10 zone、面积、公用工程容量和安全约束，设备不能占用同一不可共享空间或超 utility。
- [x] T15 实现设备合同承诺、应付、付款、租赁义务、设施尾款和 closing cash 守恒。
- [x] T16 建立检漏设备 calibration drift -> discrepancy -> delayed Knowledge -> remediation -> reinspection -> close tracer。
- [x] T17 实现设备供应商、物流方、安装商和校准机构的版本化确定性策略与资源日历。
- [x] T18 实现 snapshot/restore/reset、100 次确定性、中途崩溃恢复和非法状态失败关闭。
- [x] T19 新增 `world-packs/hctm-genesis/campaigns/capability-build/`，只通过 M10 terminal hash/eligibility 初始化。

验收：时间、订单或资产登记不能凭空创造设备能力；所有关键设备和实验室门可重放、解释和对账。

### C3 - AESE 人员、技能与班次世界（第 3-4 周）

- [x] T20 实现虚构劳动力市场、候选出现、筛选、面试、offer 接受/拒绝、到岗和离职前置模型。
- [x] T21 实现岗位、headcount、技能、培训、实操考核、认证有效期、替补和一班制 reducer。
- [x] T22 实现候选/员工 actor-scoped Knowledge 与隐私投影；无权限角色不能读取候选详情。
- [x] T23 实现 equipment availability 对培训、实操认证和班次激活的依赖，禁止纸面培训直接形成技能。
- [x] T24 验证未到岗、培训失败、证书过期、关键岗位空缺、职责冲突和重复 actor 不计入 capability gate。

验收：IAOS 招聘/培训记录与实际接受、到岗和掌握技能明确分离；资格门可确定性重放。

### C4 - IAOS 采购、资产、组织与资格治理（第 3-6 周）

- [x] T25 按 IAOS 规则建立新的独立 branch/worktree，先读其 AGENTS、Agent Context 和 Code Map。
- [x] T26 优先使用 metadata/config package 建立资金、预算、需求、报价、订单/租赁、收货、资产、安装/调试、岗位、编制、候选、offer、员工、培训、资格和班次对象。
- [x] T27 实现资本催缴/入账、CAPEX/headcount、采购选择、订单、收货、资产验收、offer、入职、培训、认证、班次和能力验收的 allowlist Capability/Process/Decision。
- [x] T28 固定 `genesis.capability.fund/procure/accept`、`genesis.workforce.plan/hire/qualify`、`genesis.capability.payment.approve` 权限和岗位 mandate。
- [x] T29 保证 business record、audit、journal committed outcome 与 Outbox 同事务；回滚不产生 outcome。
- [x] T30 验证 tenant RLS、候选隐私、越权自批、超预算、未验收付款、伪造到岗/资格、重复订单/offer、并发更新、过期 mandate 和失败无部分写入。
- [x] T31 两仓分别提交、记录 revision、部署并完成 contract tests。

验收：AESE 不直写 IAOS；IAOS 不能伪造到账、设备能力、人员到岗或技能；所有动作受权限、Policy、幂等、事务和审计约束。

### C5 - 统一岗位与 Capability Build Play（第 5-7 周）

- [x] T32 实现项目负责人、采购、HR、设备和质量岗位的确定性策略；人类接管复用相同治理链。
- [x] T33 在 World Play 增加 Capability Build campaign、资金方案、设备供应/安装时间轴、人员漏斗、技能矩阵和联合 gate。
- [x] T34 在 M10 空间图上展示设备落位、实验室/仓储能力和 utility 占用，不增加自由布局或 3D 编辑。
- [x] T35 展示 World 实际状态、IAOS 管理状态、角色 Knowledge、discrepancy 和因果证据，不提供直接改设备/技能状态入口。
- [x] T36 完成键盘、ARIA、焦点、移动端、金额/数量/能力单位、数据 owner、隐私和虚拟/真实时间表达。

验收：非研发用户可完成资金、采购、安装调试、招聘培训、异常整改和联合能力验收，并能解释阻塞原因。

### C6 - 全链验收与交付（第 7-8 周）

- [x] T37 从 M10 clean terminal state 分别执行 Agent 路径和人类接管路径，终态业务结果/state hash 一致。
- [x] T38 对账 World events、设备/空间/人员/Knowledge、IAOS Decision/Process/Capability、journal、Outbox、预算、合同、付款和现金。
- [x] T39 验证断线、SSE 丢失、重启、重复点击、陈旧 cursor、乱序、并发、快照恢复和安全 reset。
- [x] T40 验证 tenant-other、只读用户、候选隐私、越权自批、资金不足、未验收付款、伪造资格和重复动作。
- [x] T41 运行 Go、Schema、PostgreSQL、IAOS modules、前端、Playwright 及 M7-M10 回归和 390/1280/1440 三视口验收。
- [x] T42 编写 M11 runbook/evidence，更新设计、计划、Roadmap、Code Map、Progress Log、Atlas 和两仓 revision/部署信息。

验收：只有资金、设施尾款、设备、实验室、仓储、人员、技能、班次、安全和治理门全部通过时输出 `industrialization_eligible=true`；UI 完成或 IAOS 记录不能替代 World 证据。

## 4. 首版编码方向

C0 最终冻结前使用以下虚构编码方向：

```text
CAP-COOLING-PLATE-A-LINE
EQ-FORM-01
EQ-CNC-01
EQ-LASER-WELD-01
EQ-WASH-01
EQ-LEAK-TEST-01
EQ-ASSEMBLY-01
LAB-QUALITY-01
POS-PLANT-MANAGER
POS-PLANNING
POS-PROCUREMENT
POS-QUALITY
POS-PROCESS
POS-EQUIPMENT
POS-WAREHOUSE
POS-OPERATOR
POS-INSPECTOR
SHIFT-HCTM-SZ-A
```

供应商、候选人、培训机构和金融主体全部虚构。M8 `LAS-WLD-02` 不默认计入本清单。

## 5. 完成定义

- Capability Build campaign 可验证、初始化、推进、重放、快照恢复和 reset。
- 资金来源、新预算、设施尾款、设备/租赁承诺、工资准备金和现金严格对账。
- 所有关键设备实际到货、安装、调试、校准、安全和能力验收，空间/utility 不超限。
- 最小核心团队实际到岗、资格有效，一班制、关键替补和职责隔离通过。
- 检漏设备失败 tracer 有 observation/intent/outcome/world consequence/discrepancy close 完整链。
- 人类/Agent 共用治理能力；越权、隐私泄露、自批、超预算、未验收付款和伪造资格失败关闭。
- IAOS 事务、journal 与 Outbox 原子，tenant/RLS/幂等/并发验收通过。
- M7-M10 全链零回归；runbook、evidence、Atlas、两仓提交和部署完整。
- `industrialization_eligible` 只表示可进入 M12，不表示产品、试生产、PPAP 或量产就绪。

## 6. 不纳入 M11

- 客户 RFQ、产品定点、产品设计、BOM、routing、成本报价和正式工艺参数；属于 M12。
- APQP、产品专用工装、首件、首批物料、试生产、PPAP 和 SOP；属于 M12。
- 第一张正式订单到交付、开票、回款和实际项目盈亏；属于 M13。
- 参数化分支实验和 A/B；属于 M14。
- 完整采购/SRM、HRIS/LMS、薪酬社保、EAM、固定资产会计、融资、3D 或多工厂产品。

## 7. 并行与所有者规则

- C0/C1 由 AESE owner 先冻结合同、能力、资金和治理门。
- C2 与 C3 只能在 C0 合同稳定后并行，分别拥有设备世界和人员世界文件，公共状态变更由单一 owner 合并。
- C4 只能在 payload、权限和 IAOS gap 冻结后由新的 IAOS worktree owner 开始。
- C5 在 API/view model 稳定后开始，不得为了 UI 直接改 World/IAOS 状态。
- C6 串行收口两仓证据；并行 agent 必须有子计划、owner 和不重叠 worktree。
- 不得覆盖当前共享工作区中的测试修改、截图变化或验收产物。
