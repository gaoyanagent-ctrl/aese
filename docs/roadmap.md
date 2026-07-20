# AESE Roadmap

本文件是 AESE 当前里程碑状态和下一步优先级的权威来源。

最后更新：2026-07-20。

## 1. 里程碑状态

| 里程碑 | 目标 | 状态 | 完成证据 |
| --- | --- | --- | --- |
| M0 项目初始化 | 仓库、背景、规则、GitHub | Completed | README、AGENTS、初始提交 |
| M1 虚拟企业蓝图 | 华辰集团、苏州基地、电池冷却板 A 线 | Completed | HCTM Virtual Enterprise Blueprint |
| M2 业务与技术规格 | 对象、事件、seed、演示故事 | Completed (docs) | 4 份 HCTM 规格文档 |
| M2.5 工程治理 | 架构边界、索引、code map、执行规则 | Completed | 本轮治理文档 |
| M3 可执行场景包 | JSON 场景包、校验器、IAOS apply/replay tracer | Completed | pack、CLI、execution evidence、IAOS commits |
| M3V 快速 2D 沙盘 | 七幕故事、22 事件、A 线画布、KPI 和 Agent 建议预览 | Completed | 前端、preview、18 unit/component tests、9 E2E、3 viewport screenshots |
| M4 异常场景运行 | 延期、设备故障、来料不良进入 IAOS 运行链 | Completed | 三类 ingress、状态影响、事务 Outbox、租户/幂等及 canonical replay evidence |
| M5 Agent MVP | 计划、质量、经营分析 Agent | Completed | 9 个受治理只读工具、三 Agent live tracer、跨租户与零业务写入证据 |
| M6 在线 2D 企业沙盘 | IAOS 实时事件、库存、产线、异常和 Agent 运行结果 | Active | DES-004、PLAN-M6-001 |

## 2. 当前阶段

M3、M3V、M4 和 M5 已完成。当前进入 M6：补齐完工入库、两次发运和最小成本影响事实，建立可恢复的场景 snapshot/cursor/SSE 合同，实现 `IaosScenarioDataSource` 和 Preview/Live 双模式。

M6 的最小成功标准：

1. 在线状态只来自 IAOS 租户事实，不复用 Preview expected outcomes 冒充 Live。
2. canonical 事件 17-22 形成幂等、事务化的完工、入库和发运记录。
3. 场景事件可按持久 cursor 查询和补发，SSE 断线不丢失状态。
4. Live 最终显示需求 12,000、实发 11,700、缺口 300 和 `partially_shipped`。
5. 三 Agent 建议能追溯 correlation、业务对象和 IAOS Tool Call。
6. Preview/Live、跨租户、断线恢复、桌面和移动端验收通过。

## 3. M6 范围

包含：

- 事件 17-22 的 IAOS 完工、入库、发运和库存事实闭环。
- tenant-scoped snapshot、持久 cursor event query、场景 SSE 和 recommendation query。
- M5 recommendation envelope 的在线持久化与 Tool Call 证据关联。
- `IaosScenarioDataSource`、Preview/Live 双模式、断线恢复和完整度展示。
- 全链路回放、跨租户、桌面与移动端自动化验收。

不包含：

- 第二场景、多工厂和多产品族。
- 真实 LLM 自主 Agent 和建议自动执行。
- 完整成本核算、总账和利润中心。
- 3D 工厂、布局编辑器或通用离散事件仿真引擎。

## 4. M6 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| L0 | 合同冻结、成本完整度决策和基线 | Pending |
| L1 | 完工、入库与两次发运事实闭环 | Pending |
| L2 | snapshot、持久 cursor 和场景 SSE | Pending |
| L3 | Agent 建议持久化与在线查询 | Pending |
| L4 | `IaosScenarioDataSource` 和 Live UX | Pending |
| L5 | 全链路、断线、跨租户和三视口验收 | Pending |

## 5. 风险与依赖

- IAOS Platform、PostgreSQL、NATS 和 O2D 可运行，`tenant-hctm` 的 work_order metadata、workflow config 和 tracer 数据已完成 seed/apply。
- `/iaos/iaos-go` 当前主分支本地领先远程，任何集成开发必须使用独立 worktree 并确认基线。
- HCTM 业务字段与 IAOS 现有 legacy `sales_order` 物理模型存在差异，需要兼容性报告，不能直接假设可导入。
- IAOS 当前 O2D 自测硬编码 `tenant-001` 和旧订单，需要避免污染 HCTM tracer。
- 正式事件必须走 Outbox 或受治理 ingress，不能把 direct NATS 当最终实现。
- IAOS 受治理 scenario apply/reset 已实现 M3 allowlist；订单确认 CAS、workflow/event 去重、跨节点原子事务及 work_order API 已实证。
- legacy 表没有全面 FORCE RLS；M3 scenario adapter 已在所有查询/更新/删除中显式绑定 tenant。平台长期仍应继续推进全表 RLS FORCE hardening。
- M4 的显式 tenant predicates 已关闭当前入口越界路径，但 tenant-safe composite foreign key 和 metadata `version` 的平台级排序仍需后续 hardening。
- 首版 2D 沙盘是确定性预览，不应被描述为 IAOS 实时运行结果；界面必须显示 Preview 数据源状态。
- `preview.json` 只承载视图状态和 delta，不得复制 MRP、排产或 Agent 决策逻辑。
- 当前通用 `/api/v1/events/stream` 无持久 cursor 且缓冲满会丢事件，只能作为监控通道，不能直接作为 M6 恢复合同。
- 完工和发运是内生业务动作，不能为了复用现有接口而错误接入 simulation ingress。
- HCTM 尚无已批准的成本金额基线；在基线确认前，在线经营分析的成本部分必须保持 `partial`。

## 6. 后续里程碑入口条件

完成 M6 前：

- 在线库存、完工、发运和订单状态满足 expected outcomes 与库存守恒。
- snapshot 与 cursor 来自一致观察边界，断线补发和事件去重可复现。
- Preview/Live 数据源明示，错误时不静默混用。
- Agent 建议、对象引用和 Tool Call 证据属于同一 tenant/correlation。
- IAOS 与 AESE 两仓测试、部署、runbook 和 evidence 完整。
