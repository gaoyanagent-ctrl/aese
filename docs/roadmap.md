# AESE Roadmap

本文件是 AESE 当前里程碑状态和下一步优先级的权威来源。

最后更新：2026-07-19。

## 1. 里程碑状态

| 里程碑 | 目标 | 状态 | 完成证据 |
| --- | --- | --- | --- |
| M0 项目初始化 | 仓库、背景、规则、GitHub | Completed | README、AGENTS、初始提交 |
| M1 虚拟企业蓝图 | 华辰集团、苏州基地、电池冷却板 A 线 | Completed | HCTM Virtual Enterprise Blueprint |
| M2 业务与技术规格 | 对象、事件、seed、演示故事 | Completed (docs) | 4 份 HCTM 规格文档 |
| M2.5 工程治理 | 架构边界、索引、code map、执行规则 | Completed | 本轮治理文档 |
| M3 可执行场景包 | JSON 场景包、校验器、IAOS apply/replay tracer | Completed | pack、CLI、execution evidence、IAOS commits |
| M3V 快速 2D 沙盘 | 七幕故事、22 事件、A 线画布、KPI 和 Agent 建议预览 | Completed | 前端、preview、18 unit/component tests、9 E2E、3 viewport screenshots |
| M4 异常场景运行 | 延期、设备故障、来料不良进入 IAOS 运行链 | Active | `eam.machine.down` ingress/replay evidence；另两类待实现 |
| M5 Agent MVP | 计划、质量、经营分析 Agent | Not started | - |
| M6 在线 2D 企业沙盘 | IAOS 实时事件、库存、产线、异常和 Agent 运行结果 | Not started | - |

## 2. 当前阶段

M3 和 M3V 已完成，M4 正在执行。`eam.machine.down` 已通过受治理 simulation ingress、事务 Outbox 和 AESE canonical replay 的真实验收；AESE replay 已能把供应商延期和来料检验失败按 canonical metadata 送入同一受治理入口，当前继续完成 IAOS 对这两类事件的对象解析、状态影响和真实验收，再用已验证的 `ScenarioDataSource` 边界推进在线数据源。

M3V 的最小成功标准：

1. 3 到 4 个工作日内提供可访问的 React 2D 沙盘。
2. 页面直接显示苏州制造基地电池冷却板 A 线，不设置落地页。
3. 七幕故事和 22 个事件可播放、暂停、单步、倍速和重置。
4. 画布、事件流、对象详情、五项 KPI 和三类 Agent 建议同步变化。
5. 静态预览数据经过 adapter 提供，后续可替换为 IAOS API/SSE。
6. 桌面与移动端自动化和截图验收通过。

M3 已完成的验收基线：

1. HCTM 场景数据不再只存在于 Markdown，而是形成版本化 JSON pack。
2. `aese validate` 能离线发现 schema、引用、数量、时间线和幂等错误。
3. `aese inspect` 能输出对象数、事件数、依赖和故事摘要。
4. 至少一个受控路径可把基础数据导入 IAOS，并触发 `o2d.order.confirmed`。
5. 重复 apply/replay 不产生重复业务对象或重复事件副作用。
6. 提供 API、数据库、事件流和演示故事验收证据。

## 3. M3V 范围

包含：

- React + TypeScript + Vite 前端。
- 苏州基地 A 线 2D 工艺与物流画布。
- 静态 `preview.json` 和确定性时间线 reducer。
- 七幕/22 事件播放控制、事件流、对象详情、KPI 和 Agent 建议。
- 桌面优先、移动端可用的响应式布局与自动化测试。

不包含：

- IAOS 实时 API/SSE 数据源。
- 真实 Agent 调用和建议执行。
- 业务 CRUD、审批、权限和审计界面。
- 3D 工厂或通用离散事件仿真引擎。

## 4. M3V 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| V0 | 前端骨架、视图模型和 preview 数据 | Completed |
| V1 | 2D 工厂画布和工作台布局 | Completed |
| V2 | 时间线播放内核和视觉状态变化 | Completed |
| V3 | 事件、KPI、对象详情和 Agent 建议 | Completed |
| V4 | 响应式、自动化测试和运行说明 | Completed |

## 5. 风险与依赖

- IAOS Platform、PostgreSQL、NATS 和 O2D 可运行，`tenant-hctm` 的 work_order metadata、workflow config 和 tracer 数据已完成 seed/apply。
- `/iaos/iaos-go` 当前主分支本地领先远程，任何集成开发必须使用独立 worktree 并确认基线。
- HCTM 业务字段与 IAOS 现有 legacy `sales_order` 物理模型存在差异，需要兼容性报告，不能直接假设可导入。
- IAOS 当前 O2D 自测硬编码 `tenant-001` 和旧订单，需要避免污染 HCTM tracer。
- 正式事件必须走 Outbox 或受治理 ingress，不能把 direct NATS 当最终实现。
- IAOS 受治理 scenario apply/reset 已实现 M3 allowlist；订单确认 CAS、workflow/event 去重、跨节点原子事务及 work_order API 已实证。
- legacy 表没有全面 FORCE RLS；M3 scenario adapter 已在所有查询/更新/删除中显式绑定 tenant。平台长期仍应继续推进全表 RLS FORCE hardening。
- 首版 2D 沙盘是确定性预览，不应被描述为 IAOS 实时运行结果；界面必须显示 Preview 数据源状态。
- `preview.json` 只承载视图状态和 delta，不得复制 MRP、排产或 Agent 决策逻辑。

## 6. 后续里程碑入口条件

完成 M3V 前：

- 22 个事件和七幕故事在桌面与移动端均可完整播放。
- 最终 KPI 与 expected outcomes 一致。
- `ScenarioDataSource` 隔离静态预览数据与 UI。
- 自动化测试、截图验收和本地运行说明齐全。

进入 M4 前：

- M3 的 pack、validator、apply/replay 和 reset 均通过。
- `o2d.order.confirmed` tracer 在 `tenant-hctm` 下可复现。
- 事件关联、RLS 和幂等证据完整。

进入 M5 前：

- M4 三类异常均有结构化事件和可查询业务上下文。
- Agent 所需工具通过 IAOS Capability / AI Tool Registry 暴露。

进入 M6 在线沙盘前：

- 业务链和 Agent 输出稳定。
- 2D 沙盘的数据查询合同已经固定。
