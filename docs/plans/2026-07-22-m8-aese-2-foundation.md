---
id: PLAN-M8-001
title: M8 AESE 2.0 三态世界与仿真内核实施计划
date: 2026-07-22
status: completed
author: Codex + User
tags: [m8, aese-2, world-state, simulation, genesis]
---

# M8 AESE 2.0 三态世界与仿真内核实施计划

## 1. 交付目标

用 **4 到 6 周**完成 AESE 2.0 的最小架构闭环，而不是一次实现完整企业生命周期：

> 在华辰苏州工厂单一 world run 中，确定性推进虚拟时间，制造一条设备实际退化但 IAOS 和计划员尚未知的偏差；经过观察、受治理 IAOS 动作和世界结果计算后关闭偏差，并可完整重放。

该 tracer 是 Project Genesis 的地基。M8 不以“页面更多”或“文档完成”作为实现完成。

## 2. 前置决策门

- [x] G1 用户批准 ADR-004 的状态所有权边界。
- [x] G2 World Store 采用独立 PostgreSQL database、账号、迁移和连接边界；备份与本地启动细节在 F0 固化。
- [x] G3 Actor Knowledge 首版采用结构化事实引用、观察/有效时间、来源、置信度、可见范围和替代关系，不保存自由文本长期记忆。
- [x] G4 以 DES-008 固定 `observation/intent/committed_outcome`、World Bridge journal/cursor 和跨仓交付顺序。
- [x] G5 批准 M8 tracer 的设备、人员和产能不变量；经济面限定为现金、承诺支出和应收回款。

G1-G5 与 F0-F5 已全部通过。M8 foundation 已完成；后续企业生命周期扩展应另立里程碑和唯一 active plan。

## 3. 实施切片

### F0 - 基线与合同冻结（第 1 周）

- [x] T1 评审 AESE 2.0 原始构思，识别与 Go 技术栈、ADR-001/003 和当前 M8 定义的冲突。
- [x] T2 建立 ADR-004、DES-007 和本 active plan。
- [x] T3 建立三态术语表、对象所有权矩阵和数据分类清单，并固化独立 PostgreSQL 的迁移、备份与本地启动合同。
- [x] T4 定义 `WorldRun/Event/Snapshot/Knowledge/Discrepancy` JSON Schema 与 Go 类型；同时按 DES-008 落地 Observation/Intent/CommittedOutcome strict Schema、fixture 与 canonical hash。
- [x] T5 通过 DES-008 定义 observation/intent/outcome envelope、错误、重试、幂等、权限和 cursor 恢复语义。
- [x] T6 审计 IAOS scenario/simulation/Capability/Event 能力并在 DES-008 记录复用点与缺口；未修改 IAOS 实现。
- [x] T7 以 `atlas/system-atlas-planned.json` 声明 AESE World、Time、Knowledge、Genesis 和 IAOS Bridge 的 planned 节点/依赖，并提交 Atlas update 声明。

验收：Completed。五个决策门、八类机器合同、独立 World Store 运维边界和 planned Atlas 投影均已落地；两仓职责和交付顺序可由新 agent 无歧义复述。

### F1 - 确定性仿真内核（第 1-2 周）

- [x] T8 新增 `internal/world`、`simtime`、`simevent` 与 `rules`，实现纯函数 reducer。
- [x] T9 实现 pause、step、run-until、队列优先级、因果链和稳定 tie-break。
- [x] T10 实现 seed、规则版本、事件日志、快照、state hash、快照恢复和离线 replay。
- [x] T11 实现默认 dry-run 的 `aese world validate|inspect|run|replay` 命令；写 artifact 必须显式 `--apply --output`。
- [x] T12 为跨秒/年时间尺度、重复事件、100 次确定性、损坏快照恢复和非法时间倒退建立单元测试。

验收：Completed。相同输入重放 100 次的事件顺序、日志 hash 和 state hash 一致；当前 F1 尚无业务 KPI，故不伪造 KPI 验收。未知事件和规则版本失败关闭。

### F2 - 三态与偏差 tracer（第 2-3 周）

- [x] T13 建立世界设备、班次、人员认证和实际产能的最小模型。
- [x] T14 建立 Actor、Observation、Knowledge、置信度、可见范围和信息延迟模型。
- [x] T15 建立 IAOS 设备台账/维护记录只读投影和 stable ref 映射。
- [x] T16 模拟 `LAS-WLD-02` 实际退化 -> 传感观察 -> 设备工程师获知 -> IAOS 登记 -> 偏差关闭。
- [x] T17 验证计划员在获知前后得到不同可解释结论，且不能读取未授权 World State。

验收：UI/API 可同时展示三态和差异时间线；差异关闭保留发现来源、IAOS correlation 和因果链。

### F3 - 受治理 IAOS 双向桥（第 3-4 周）

- [x] T18 实现 AESE-side bridge adapter、cursor、outcome 去重和故障恢复。
- [x] T19 在独立 IAOS worktree 实现最小 observation ingress / governed intent 或复用 Capability，保留 RLS、权限、审计与 Outbox。
- [x] T20 验证 tenant-other、重复 outcome、乱序、超时、部分失败和重连。
- [x] T21 建立能力缺口台账，记录企业活动、角色、对象、Capability、临时处理和优先级。
- [x] T22 完成两仓合同测试并分别提交、部署和记录 revision。

验收：AESE 不直写 IAOS DB、不保存长期 JWT、不以 direct NATS 作为正式路径；重复执行无重复业务副作用。

### F4 - Genesis world pack 与兼容迁移（第 4-5 周）

- [x] T23 建立 `world-packs/hctm-genesis` manifest、世界初态、角色、规则和 tracer 数据。
- [x] T24 将现有 `order-expedite-01` 的 22 事件通过 adapter 映射为兼容 World Event。
- [x] T25 加入设备、人员、材料和最小现金守恒验证器。
- [x] T26 证明旧 M7 场景可独立运行，也可作为 Genesis 后期运营片段挂接。
- [x] T27 为 pack 校验、稳定编码、单位/精度、Asia/Shanghai 和引用完整性增加测试。

验收：旧场景零破坏；Genesis pack 可离线验证、确定性初始化、推进、复位和重放。

### F5 - World Play 最小界面与验收（第 5-6 周）

- [x] T28 增加 World 模式入口、虚拟时钟、step/run-until 控制和运行状态。
- [x] T29 增加 World/IAOS/Knowledge 三栏对照、差异列表和因果时间线。
- [x] T30 在现有 A 线画布显示设备实际状态、系统记录和角色是否已知，不新增 3D。
- [x] T31 验证人类接管设备工程师岗位后，与 Agent 使用同一 IAOS Capability 和权限合同。
- [x] T32 完成 Go、Schema、合同、前端、Playwright、重启恢复、租户和幂等验收。
- [x] T33 编写 runbook/evidence，更新路线图、架构、code map、进展日志和 Atlas。

验收：非研发用户可观察偏差、推进时间、执行受治理发现/处置并重放相同结果；M7 全链回归通过。

## 4. 跨仓交付顺序

| 顺序 | AESE | IAOS | 阻断条件 |
| --- | --- | --- | --- |
| 1 | Schema、fixture、离线 contract tests | gap audit | ADR-004 未批准 |
| 2 | World 内核和 mock bridge | 最小 API/Capability 设计 | envelope 未冻结 |
| 3 | adapter 与失败路径 | 独立 worktree 实现 | 权限/Outbox/幂等未通过 |
| 4 | Genesis tracer 与 UI | 部署合同端点 | 两仓 revision 未记录 |
| 5 | 全链验收和 evidence | 集成验收 | tenant、重放或回归失败 |

IAOS 修改前必须按其 `AGENTS.md` 读取上下文并建立独立 branch/worktree；AESE 提交不得混入 IAOS 文件。

## 5. 完成定义

- ADR-004 accepted，三态所有权和存储边界已落定。
- World Runtime 可确定性运行、崩溃恢复、重放并产生稳定 state hash。
- 设备退化 tracer 展示三态不一致、有限认知、受治理处置和偏差关闭。
- 所有数量、金额、时间、事件和 stable ref 满足 AESE 合同。
- IAOS 双向桥通过租户、权限、幂等、乱序和无部分写入测试。
- 现有 M7 22 事件、三 Agent、Preview/Live 和 reset 验收不回归。
- runbook、evidence、两仓 revision、Roadmap、Code Map、Progress Log 和 Atlas 一致。

## 6. 后续路线（不纳入 M8 完成条件）

| 里程碑 | 纵向成果 |
| --- | --- |
| M9 Genesis Incorporation | 注资、法人、治理、初始组织、预算和 CEO/CFO 岗位闭环 |
| M10 Genesis Plant Build | 选址、建设项目、设施依赖、工期、预算和项目总监闭环 |
| M11 Genesis Capability Build | 设备采购安装调试、招聘培训认证、仓储/质量设施和投产门 |
| M12 Genesis Industrialization | RFQ、报价、APQP、供应商、试生产、PPAP 和 SOP |
| M13 Genesis First Delivery | 复用 O2D，补齐开票、回款、实际成本与首批项目盈亏 |
| M14 Experiments | checkpoint/fork、参数化 baseline/A/B 和多方案比较 |

每个里程碑只扩展一个可重放纵向闭环；任何阶段均不以横向模块数量作为完成度。

## 7. 主要风险与控制

- 双系统事实混淆：所有对象标注 owner，UI 明示来源与时间；禁止共享表。
- 过度建模：M8 只做一台设备、两个角色、一个偏差闭环和最小守恒。
- 非确定性 Agent：世界结果只由版本化规则计算；Agent 输出作为 intent，不直接改状态。
- 状态爆炸：首版单 tenant、单 writable world run、单分支；分支只冻结数据合同。
- IAOS 范围膨胀：用能力缺口台账排序，只为 tracer 实现最小平台能力。
- 旧功能回归：M7 evidence 作为强制兼容门，不迁移或删除现有 pack。
