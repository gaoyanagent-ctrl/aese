---
id: PLAN-M12-001
title: M12 Project Genesis 产品工业化与量产批准实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m12, genesis, industrialization, apqp, ppap, trial-production]
---

# M12 Project Genesis 产品工业化与量产批准实施计划

## 1. 交付目标与工期

用 **8 到 10 周**把 M11 的 `industrialization_eligible=true` 推进为：

> 华辰苏州制造公司完成首个虚构客户 `CUST-SGNEV` 的 `HCTM-BCP-A01` 项目定点、产品/工艺发布、供应商与工装准备、APQP、两轮试生产、问题整改和 PPAP 批准，输出 `serial_production_eligible=true`。

首版固定单客户、单产品、单客户项目、单工厂/产线、一套发布版本、关键供应商集合、两轮试制和一条焊接泄漏/Cpk 不足异常。终态只允许 M13 接收第一张正式订单，不产生正式需求、销售库存、发票、收入或回款。

## 2. 前置门

- [x] G1 M11 已完成并机器输出 `industrialization_eligible=true`。
- [x] G2 ADR-004、DES-008、DES-009 与 M8-M11 的 World/Bridge/Knowledge/资金边界继续有效。
- [x] G3 DES-013 固定 M12 三态所有权、APQP/PPAP、旧 HCTM 兼容和 M13 边界。
- [x] G4 冻结虚构 RFQ、客户要求、年量/爬坡、目标价格、项目时间、定点规则和报价假设。
- [x] G5 冻结开发预算、客户工装预付款/合同负债、付款节点、样件/试制成本和现金缓冲。
- [x] G6 冻结 product/BOM/routing/PFMEA/control-plan revision、供应商/工装/物料、试制批次、MSA/Cp/Cpk/良率/节拍/PPAP 门和异常基线。
- [x] G7 冻结 `scenario-packs/hctm` stable-code/revision/hash 兼容矩阵和 M13 terminal contract。
- [x] G8 完成 IAOS RFQ/quotation/project/product/process/APQP/quality/PPAP gap audit，冻结两仓 payload、权限和交付顺序。

G4-G8 完成前只允许 Schema、fixture、离线规则、兼容性分析和只读 IAOS gap audit；不得开发 IAOS 写端点或把旧 HCTM fixture 标记为 Genesis 已完成事实。

## 3. 实施切片

### D0 - 客户项目、产品、质量与资金机器合同（第 1 周）

- [x] T1 定义 RFQ、Requirement、Quotation、CustomerProject、ProductRevision、BOM、Routing、APQP、Tooling、SupplierMaterial、TrialBuild、Measurement、Issue、PPAP 和 ReleaseManifest stable code/owner。
- [x] T2 固定客户要求、年量/爬坡、目标节拍、关键特性、报价、定点、开发周期和项目角色基线。
- [x] T3 固定 EBOM/MBOM、routing、process flow、PFMEA、control plan、work instruction 的最小结构、revision、hash 和一致性规则。
- [x] T4 固定供应商、工装、量检具、首批物料、批次追溯、试制数量、报废和不可销售规则。
- [x] T5 固定 MSA、尺寸/功能、泄漏、节拍、良率、Cp/Cpk、PPAP 内容和客户批准门。
- [x] T6 固定开发预算、客户预付款/合同负债、承诺、付款、试制成本、工资准备金和现金缓冲。
- [x] T7 扩展 JSON Schema、Go 类型、strict parser、canonical hash、payload registry 及合法/破损 fixture。
- [x] T8 定义 `eligible -> rfq -> nominated -> product_design -> process_design -> supplier_tooling -> trial_1 -> remediation -> trial_2 -> ppap -> serial_production_eligible` 状态机。
- [x] T9 完成旧 HCTM compatibility ledger、IAOS gap ledger、职责冲突/权限矩阵、Atlas 依赖和两仓交付矩阵。

验收：客户、产品、工艺、供应、质量、资金、三态 owner 和 M13 边界无歧义，G4-G8 关闭，所有单位/阈值/版本/证据可离线验证。

### D1 - RFQ、报价、定点与项目资金（第 1-2 周）

- [x] T10 实现虚构客户 RFQ/要求产生、送达延迟和 actor-scoped Knowledge。
- [x] T11 实现技术、产能、供应链、投资、质量和工期 hard constraints，可行性失败先于报价。
- [x] T12 实现材料、加工、工装摊销、质量、物流和风险的版本化成本/报价解释，金额使用定点。
- [x] T13 实现报价审批、价格例外、提交、客户定点/拒绝的 observation/intent/outcome/world consequence tracer。
- [x] T14 实现客户工装预付款实际到账、IAOS 入账和合同负债对账；禁止计收入/利润。
- [x] T15 验证无可行性、无预算、自批价格例外、虚构客户决定、预付款未到账和重复定点全部失败关闭。

验收：客户决定由 World 产生；报价、定点、项目预算和现金/合同负债因果完整，RFQ 不会变成正式订单。

### D2 - 产品、工艺与 APQP 版本治理（第 2-4 周）

- [x] T16 实现 `HCTM-BCP-A01` 产品 revision 与客户规范/关键特性模型。
- [x] T17 实现 EBOM -> MBOM 转换和 component/uom/decimal/scrap/reference integrity。
- [x] T18 实现 routing/process flow 与 M11 accepted equipment/work center/qualified skill 的能力匹配。
- [x] T19 实现 PFMEA 风险、control plan 检查频次/反应计划、work instruction 和 inspection plan 一致性校验。
- [x] T20 实现 APQP 六个 gate、owner、evidence、审批、退回和重新提交状态机。
- [x] T21 实现 engineering change set：已发布/已用于试制 revision 不可原地修改，变更必须保留 supersedes 和影响分析。
- [x] T22 生成 product/BOM/routing/control-plan release manifest/hash，并与旧 HCTM stable code 做离线 compatibility check。

验收：所有试制输入锁定到同一发布基线；缺引用、版本漂移、未审批变更和旧 fixture 不匹配全部失败关闭。

### D3 - 供应商、工装与首批物料世界（第 3-5 周）

- [x] T23 实现关键材料供应商的报价、样件、过程审核、材料证书和实际能力策略。
- [x] T24 实现焊接工装、量检具和试验资源的设计、制造、交付、校准与验收依赖。
- [x] T25 实现首批材料采购、运输、收货、批次、证书、IQC、隔离和放行 world reducer。
- [x] T26 实现项目预算、工装/材料承诺、应付、付款、合同负债、试制成本和 closing cash 守恒。
- [x] T27 验证未批准供应商、证书/追溯缺失、不合格批次、未验收工装和超预算材料不能进入试制。
- [x] T28 实现 snapshot/restore/reset、100 次确定性、中途崩溃恢复和非法状态失败关闭。

验收：IAOS 供应商/库存/工装记录不能伪造实际能力、材料属性或到货；每个试制投入可追溯。

### D4 - 试生产、质量异常与 PPAP 世界（第 4-7 周）

- [x] T29 实现 trial order、物料发放、设备/工装/人员占用、工序执行、样件谱系、报废和退料 reducer。
- [x] T30 实现尺寸、功能、泄漏、MSA、节拍、良率和过程能力的版本化确定性计算。
- [x] T31 实现首轮焊接泄漏超限/Cpk 不足 -> discrepancy -> delayed Knowledge -> containment tracer。
- [x] T32 实现 engineering-change intent/outcome 后的参数/工装 revision B、再培训、复试计划和资源/成本后果。
- [x] T33 实现第二轮试制、复验、问题关闭、PPAP package completeness 和 PSW 状态。
- [x] T34 实现客户 PPAP 批准/拒绝/有条件批准的外部世界策略；首版成功路径最终完全批准。
- [x] T35 验证未关闭重大问题、Cpk/良率/节拍未达标、MSA 失败、包不完整或客户未批准时 production release 为 false。

验收：试制和客户批准是 World 事实；工程/质量记录、时间推进或 UI 操作不能凭空达标。

### D5 - IAOS 客户项目、工程、质量与放行治理（第 4-8 周）

- [x] T36 按 IAOS 规则建立新的独立 branch/worktree，先读其 AGENTS、Agent Context 和 Code Map。
- [x] T37 优先使用 metadata/config package 建立 RFQ、报价、客户项目、APQP、产品/revision、BOM/routing、风险/控制计划、供应商/工装、试制、测量、问题/变更、PPAP 和生产放行对象。
- [x] T38 实现可行性、报价、定点登记、项目预算、工程发布、供应商/物料批准、试制放行、问题/变更、PPAP 提交/登记和生产放行的 allowlist Capability/Process/Decision。
- [x] T39 固定 `genesis.industrialization.quote/design/process/trial/quality/ppap/release` 权限和销售/工程/质量/采购/厂长/CFO mandate。
- [x] T40 保证 business record、audit、journal committed outcome 与 Outbox 同事务；回滚不产生 outcome。
- [x] T41 验证 tenant RLS、越权自批、revision 篡改、未审批试制、伪造测量/客户批准、重复定点/PPAP、并发变更、过期 mandate 和失败无部分写入。
- [x] T42 两仓分别提交、记录 revision、部署并完成 contract tests。

验收：AESE 不直写 IAOS；IAOS 不能伪造客户、供应商、试制、测量或 PPAP 事实；所有动作受权限、职责隔离、幂等、事务和审计约束。

### D6 - 统一岗位与 Industrialization Play（第 6-9 周）

- [x] T43 实现销售/项目、产品、工艺、质量、采购、厂长和 CFO 的确定性岗位策略；人类接管复用相同治理链。
- [x] T44 在 World Play 增加 Industrialization campaign、RFQ/报价/定点、APQP gate、版本关系、供应/工装、试制和 PPAP 时间线。
- [x] T45 展示 World 实际结果、IAOS 管理状态、角色 Knowledge、discrepancy、成本/现金和因果证据，不提供直接改测量/Cpk/批准入口。
- [x] T46 展示 product/BOM/routing/control-plan release manifest 与旧 HCTM compatibility 状态。
- [x] T47 完成键盘、ARIA、焦点、移动端、金额/数量/质量单位、版本、数据 owner 和虚拟/真实时间表达。

验收：非研发用户可完成 RFQ 到 PPAP/生产放行，并能解释报价、版本、质量或资金阻塞原因。

### D7 - 全链验收与交付（第 9-10 周）

- [x] T48 从 M11 clean terminal state 分别执行 Agent 路径和人类接管路径，终态业务结果/state hash 一致。
- [x] T49 对账 World events、物料/样件/测量/Knowledge、IAOS Decision/Process/Capability、journal、Outbox、预算、合同负债、成本和现金。
- [x] T50 验证断线、SSE 丢失、重启、重复点击、陈旧 cursor、乱序、并发、快照恢复和安全 reset。
- [x] T51 验证 tenant-other、只读用户、越权自批、revision 篡改、未批准材料/试制/PPAP、伪造结果和重复动作。
- [x] T52 运行 Go、Schema、PostgreSQL、IAOS modules、前端、Playwright 及 M7-M11、M3/O2D 兼容回归和 390/1280/1440 三视口验收。
- [x] T53 编写 M12 runbook/evidence，更新设计、计划、Roadmap、Code Map、Progress Log、Atlas 和两仓 revision/部署信息。

验收：只有 RFQ/定点、资金、产品/工艺版本、供应商/工装/物料、APQP、两轮试制、质量问题关闭、PPAP 和治理门全部通过时输出 `serial_production_eligible=true`；M13 尚未产生正式订单或可销售库存。

## 4. 首版编码方向

D0 最终冻结前使用以下虚构编码方向：

```text
RFQ-SGNEV-BCP-A01-01
NOM-SGNEV-BCP-A01-01
PROJECT-SGNEV-BCP-A01
HCTM-BCP-A01
PRODUCT-HCTM-BCP-A01-REV-A
BOM-HCTM-BCP-A01-V1
RT-HCTM-BCP-A01-V1
PFMEA-HCTM-BCP-A01-V1
CP-HCTM-BCP-A01-V1
TOOL-HCTM-BCP-A01-WELD-01
TRIAL-HCTM-BCP-A01-T1
TRIAL-HCTM-BCP-A01-T2
ISSUE-HCTM-BCP-A01-LEAK-01
PPAP-HCTM-BCP-A01-V1
RELEASE-HCTM-BCP-A01-SOP
```

客户、供应商、价格、规格和试验结果全部虚构。stable code 与旧 HCTM fixture 的一致只用于受测兼容，不代表历史事实。

## 5. 完成定义

- Industrialization campaign 可验证、初始化、推进、重放、快照恢复和 reset。
- RFQ、报价、定点、开发资金、合同负债、项目预算、试制成本和现金严格对账。
- 产品/BOM/routing/PFMEA/control-plan/工装/供应商/材料 revision 一致且发布链可审计。
- 试制物料、设备、人员、样件、报废、测量、成本和追溯满足守恒。
- 焊接泄漏/Cpk tracer 有 observation/intent/outcome/world consequence/change/retest/close/PPAP 完整链。
- 人类/Agent 共用治理能力；越权、自批、版本篡改、伪造测量/客户决定和未批准放行失败关闭。
- IAOS 事务、journal 与 Outbox 原子，tenant/RLS/幂等/并发验收通过。
- M7-M11 零回归，M3/O2D stable code/hash compatibility 通过。
- `serial_production_eligible` 只允许进入 M13，不表示正式订单、量产、交付、收入或回款已发生。
- runbook、evidence、Atlas、两仓提交和部署完整。

## 6. 不纳入 M12

- 正式销售订单、需求变更、MRP、批量采购、正式量产和客户交付；属于 M13。
- 发票、应收、回款、实际量产成本、项目利润和经营分析闭环；属于 M13。
- 第二客户、第二产品族、多工厂、多版本并行量产和复杂工程配置管理。
- 完整 CRM/CPQ/PLM/APQP/QMS/SRM、项目会计、CAD/CAE 或真实外部接口。
- 参数化分支实验和 A/B；属于 M14。

## 7. 并行与所有者规则

- D0/D1 由 AESE owner 先冻结合同、客户项目、资金和治理门。
- D2、D3 可在 D0 合同稳定后并行，但 release manifest 和共享资金状态由单一 owner 合并。
- D4 依赖 D2/D3 发布基线；不能用假测量或硬编码 PPAP 跳过实际试制。
- D5 只能在 payload、权限和 IAOS gap 冻结后由新的 IAOS worktree owner 开始。
- D6 在 API/view model 稳定后开始，不得为了 UI 直接改 World/IAOS 状态。
- D7 串行收口两仓证据；并行 agent 必须有子计划、owner 和不重叠 worktree。
- 不得覆盖当前共享工作区中的测试修改、截图变化或验收产物。
