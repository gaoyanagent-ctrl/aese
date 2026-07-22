---
id: PLAN-M10-001
title: M10 Project Genesis 工厂选址与设施建设实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m10, genesis, site-selection, facility, construction, project]
---

# M10 Project Genesis 工厂选址与设施建设实施计划

## 1. 交付目标与工期

用 **5 到 7 周**把 M9 的 `plant_project_eligible=true` 推进为：

> 华辰苏州制造公司在 20,000,000 CNY 实际现金和 15,000,000 CNY 首年预算授权内，完成候选场址评估、受治理选址与投资批准、场地控制、设施改造、公用工程接入、异常重排和验收，输出 `capability_build_eligible=true`。

首版固定单法人、苏州区域、一个被选场址、一个设施项目、一个主承包商、单币种 CNY 和一个公用工程延期异常。

## 2. 前置门

- [x] G1 M9 已完成并机器输出 `plant_project_eligible=true`。
- [x] G2 ADR-004、DES-008、DES-009 与 M8 World/Bridge/Knowledge 边界继续有效。
- [x] G3 DES-011 固定 M10 三态所有权、空间层级、选址模型、建设范围和 M11 边界。
- [x] G4 固定三个虚构候选场址、硬约束、权重、报价、项目预算、工期和异常基线。
- [x] G5 完成 IAOS project/investment/contract/payment/acceptance gap audit，冻结两仓 payload 与交付顺序。

G4/G5 完成前只允许 AESE Schema、fixture、离线评分和规则实现，不进入 IAOS 写端点开发。

## 3. 实施切片

### P0 - 场址、空间与项目机器合同（第 1 周）

- [x] T1 定义 SiteOption、Assessment、SpatialNode、UtilityCapacity、FacilityProject、WorkPackage、Milestone 和 InspectionResult stable code/owner。
- [x] T2 固定三个虚构候选场址的成本、周期、物流、人力、公用工程、风险和扩展性基线。
- [x] T3 固定一期 WBS、前置依赖、资源日历、预算/现金/承诺/付款和验收门。
- [x] T4 扩展 JSON Schema、Go 类型、strict parser、canonical hash、合法/破损 fixture 和 payload registry。
- [x] T5 定义 `eligible -> evaluating -> site_selected -> site_controlled -> project_approved -> constructing -> facility_acceptance -> capability_build_eligible` 状态机。
- [x] T6 完成 IAOS gap ledger、权限矩阵、Atlas 依赖和两仓交付矩阵。

验收：候选、空间、项目、资金和三态 owner 无歧义；所有数值、单位、时间和依赖可离线验证。

### P1 - 受约束选址决策（第 1-2 周）

- [x] T7 实现现金、预算、公用工程和最晚可用日期硬约束，硬约束失败先于评分。
- [x] T8 实现版本化可解释评分器，输出每项分值、来源、置信度、权重和敏感性摘要。
- [x] T9 建立项目负责人 observation/Knowledge，确保未调研信息不能被 Agent 读取。
- [x] T10 实现推荐 intent、CEO/CFO 审批、拒绝、重评和幂等 no-op 的离线 bridge tracer。
- [x] T11 验证绿地/超预算候选即使综合分高也不能被批准。

验收：至少三个候选有可解释比较；同一输入产生相同推荐，但最终批准仍经过 IAOS 治理。

### P2 - AESE 设施项目与空间世界（第 2-3 周）

- [x] T12 扩展 `internal/genesis`、`internal/rules` 和 world pack，实现场地交付、空间层级、WBS 和设施资产 reducer。
- [x] T13 实现承包商、园区、公用工程方和验收机构的确定性外部策略与日历。
- [x] T14 实现实际进度、返工、资源占用、合同承诺、里程碑应付、付款和现金规则。
- [x] T15 实现 utility delay -> discrepancy -> delayed Knowledge -> rebaseline -> consequence -> close tracer。
- [x] T16 实现 snapshot/restore/reset、100 次确定性和中途崩溃恢复。
- [x] T17 新增 `world-packs/hctm-genesis/campaigns/plant-build/`，通过 M9 terminal hash/eligibility 初始化。

验收：时间推进不能凭空完成工程；资金、依赖、空间和验收门失败关闭，M9 fixture 不被改写。

### P3 - IAOS 投资与项目治理（第 3-5 周）

- [x] T18 按 IAOS 规则建立新的独立 branch/worktree，读取其 AGENTS、项目上下文和 code map。
- [x] T19 优先使用 metadata/config package 建立 site assessment、investment request、facility project、WBS、agreement、contract、milestone、change、payment 和 acceptance 对象。
- [x] T20 实现选址批准、项目批准、合同承诺、进度重排、里程碑接受、付款批准和设施接受的 allowlist Capability/Process/Decision。
- [x] T21 固定 `genesis.site.evaluate/approve`、`genesis.project.execute/rebaseline/accept`、`genesis.payment.approve` 权限与岗位 mandate。
- [x] T22 保证 business record、audit、journal committed outcome 和 Outbox 同事务；回滚不产生 outcome。
- [x] T23 验证 tenant RLS、越权自批、超预算、未验收付款、重复承诺/付款、并发变更、过期 mandate 和失败无部分写入。
- [x] T24 两仓分别提交、记录 revision、部署并完成 contract tests。

验收：AESE 不直写 IAOS；IAOS 不能伪造现场进度；所有资金和项目动作受权限、Policy、幂等、事务与审计约束。

### P4 - 项目总监岗位与 Plant Build Play（第 4-6 周）

- [x] T25 实现项目负责人选址建议、WBS 管理、延期缓解和验收申请的确定性岗位策略。
- [x] T26 人类接管项目负责人/CEO/CFO 时复用同一 Capability、Decision 和审批链。
- [x] T27 在 World Play 增加 Plant Build campaign、候选比较、项目时间轴、预算/现金/承诺/付款对照。
- [x] T28 增加最小层级厂区图：site/building/office/production/warehouse/quality/utility zones。
- [x] T29 展示 World 实际进度、IAOS 计划、角色 Knowledge、discrepancy 和因果证据，不提供直接改完成度入口。
- [x] T30 完成键盘、ARIA、焦点、移动端、金额单位、数据 owner 和虚拟/真实时间表达。

验收：非研发用户可完成选址、推进工程、处理延期并验收；界面能解释为什么候选不可行、付款被阻止或项目需要重排。

### P5 - 全链验收与交付（第 6-7 周）

- [x] T31 从 M9 clean terminal state 执行 Agent 路径和人类接管路径，终态业务结果/state hash 一致。
- [x] T32 对账 World events、空间/设施、Knowledge、IAOS Decision/Process/Capability、journal、Outbox、合同、里程碑、付款和现金。
- [x] T33 验证断线、SSE 丢失、重启、重复点击、陈旧 cursor、乱序、并发变更和恢复。
- [x] T34 验证 tenant-other、只读用户、超预算、越权自批、未验收付款和重复付款。
- [x] T35 运行 Go、Schema、PostgreSQL、IAOS modules、前端、Playwright 及 M7/M8/M9 回归和三视口验收。
- [x] T36 编写 M10 runbook/evidence，更新设计、计划、Roadmap、Code Map、Progress Log、Atlas 和两仓 revision/部署信息。

验收：只有全部强制设施、utility、消防/EHS、资金和治理门通过时输出 `capability_build_eligible=true`；UI 完成不替代业务证据。

## 4. 首版候选与编码方向

P0 最终冻结前使用以下虚构编码方向：

```text
SITE-SZ-EAST-GREENFIELD
SITE-SZ-NORTH-LEASED-SHELL
SITE-SZ-WEST-BUILD-TO-SUIT
PLANT-HCTM-SZ-01
FACILITY-PROJECT-HCTM-SZ-P1
BUILDING-HCTM-SZ-MAIN
ZONE-HCTM-SZ-OFFICE
ZONE-HCTM-SZ-PRODUCTION
ZONE-HCTM-SZ-WAREHOUSE
ZONE-HCTM-SZ-QUALITY
ZONE-HCTM-SZ-UTILITY
```

候选名称和数据全部虚构，不映射真实园区、政府、承包商或公用工程主体。

## 5. 完成定义

- Plant Build campaign 可验证、初始化、推进、重放、快照恢复和 reset。
- 三个候选通过硬约束和解释性评分；资金不可行方案不能被批准。
- site control、空间、WBS、设施、utility、验收、现金、承诺和付款三态分离。
- utility delay tracer 有 observation/intent/outcome/world consequence/discrepancy close 完整链。
- 人类/Agent 共用治理能力；越权、超预算、未验收付款和重复付款失败关闭。
- IAOS 事务、journal 与 Outbox 原子，tenant/RLS/幂等/并发验收通过。
- M7/M8/M9 全链零回归。
- `capability_build_eligible` 只在设施和治理全部满足时为 true。
- runbook、evidence、Atlas、两仓提交与部署完整。

## 6. 不纳入 M10

- 生产设备、工装和质量仪器采购、安装、调试、验收；属于 M11。
- 人员招聘、培训、认证、班次和岗位补齐；属于 M11。
- 产品工艺、RFQ、APQP、试生产和 PPAP；属于 M12。
- 首批 O2D、开票、回款和实际成本；属于 M13。
- 分支实验和 A/B；属于 M14。
- 完整项目管理、合同管理、工程造价、EAM、总账、BIM 或 3D 产品。

## 7. 并行与所有者规则

- P0/P1 由 AESE owner 先冻结合同、候选和评分规则。
- P2 可在 P0 状态/资金/空间合同稳定后开始。
- P3 只能在 payload、权限和 IAOS gap 冻结后由新的 IAOS worktree owner 开始。
- P4 在 API/view model 稳定后开始；不得为了 UI 直接写 World/IAOS 状态。
- P5 串行收口两仓证据；并行 agent 必须有子计划、owner 和不重叠 worktree。
- 不得覆盖当前共享工作区中的测试修改、截图删除或验收产物。

## 8. 完成记录

- AESE：`hctm-genesis@0.3.0`、Plant Build API/UI、机器合同、确定性 tracer、runbook 与 evidence 已交付。
- IAOS：独立 `feat/m10-plant-governance` worktree，revision `23be02a`，DES-052。
- 终态：`capability_build_eligible=true`；M11 设备与人员能力不在本计划范围。
