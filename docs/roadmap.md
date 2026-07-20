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
| M6 在线 2D 企业沙盘 | IAOS 实时事件、库存、产线、异常和 Agent 运行结果 | Completed | DES-004、PLAN-M6-001、M6 evidence |
| M7 受治理场景运行控制台 | 浏览器预检、初始化、逐幕运行、分析、验证和复位 | Active | ADR-003、DES-005、PLAN-M7-001 |
| X1 System Atlas 全景治理 | 最终完成体、当前状态、依赖与进展历史 | Completed | DES-006、IAOS DES-049、双端动态图谱 |

## 2. 当前阶段

M3、M3V、M4、M5、M6 和跨里程碑的 X1 System Atlas 已完成。当前进入 M7：为现有联动中心增加受治理场景运行控制，让业务用户不依赖 CLI 完成 preflight、initialize、七幕推进、Agent 分析、verify 和 reset。未来完成体与每个构件的动态状态由 IAOS System Atlas 数据库统一提供。

M7 的最小成功标准：

1. 浏览器不直接调用 IAOS 写端点，由 AESE 薄编排 API 复用现有 Go 内核。
2. 运行具有 run ID、阶段状态机、plan hash、cursor 和 idempotency key。
3. 用户可从页面初始化、逐幕推进、运行到结束、分析、验证和安全复位。
4. 刷新、断线、重复点击和 AESE 服务重启不会产生重复业务副作用。
5. 权限不足、跨租户、陈旧 cursor 和非法状态转换全部失败关闭。
6. UI 与 CLI 对同一 pack 产生一致的 22 事件、Agent 建议、断言和 KPI。

## 3. M7 范围

包含：

- 无独立数据库的 AESE scenario orchestration API。
- pack 阶段编译、dry-run 影响、状态机、幂等和恢复。
- initialize、advance、run-to-end、analyze、verify 和 reset 编排。
- 联动中心场景运行视图、七幕 stepper、日志和危险确认。
- 权限、跨租户、并发、断线、重启及 CLI/UI 一致性验收。

不包含：

- 参数化 A/B 实验、并行分支和第二条故事。
- 真实 LLM 自主 Agent 和建议自动执行。
- AESE 业务数据库、通用任务队列或工作流引擎。
- 完整成本核算、3D 工厂和布局编辑器。

## 4. M7 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| O0 | 状态机、阶段合同和编排内核重构 | Pending |
| O1 | AESE 薄编排 API | Pending |
| O2 | IAOS 运行记录、权限和并发补强 | Pending |
| O3 | 可视化场景运行控制台 | Pending |
| O4 | 全链路、恢复、安全和三视口验收 | Pending |

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
- M7 新增 HTTP 服务但不得拥有业务数据库；运行恢复必须以 IAOS run/snapshot/event/recommendation 为事实。
- 当前场景使用固定自然键，同一 tenant/scenario 首版只能有一个可写 active run。
- 浏览器不得直接编排多个 IAOS 写 API；所有危险动作需要服务端权限、幂等和确认合同。

## 6. M7 完成条件

- 非研发用户可从浏览器完整运行并复位第一条故事。
- UI 状态只在 IAOS committed/no-op 和 snapshot cursor 证实后推进。
- 双击、并发、断线、刷新和服务重启均不重复推进阶段。
- reset 影响可预览，一次性 confirmation token 不能重放，L1 始终保留。
- IAOS 与 AESE 两仓权限、测试、部署、runbook 和 evidence 完整。

## 7. M6 完成证据

M6 已满足：

- 在线库存、完工、发运和订单状态满足 expected outcomes 与库存守恒。
- snapshot 与 cursor 来自一致观察边界，断线补发和事件去重可复现。
- Preview/Live 数据源明示，错误时不静默混用。
- Agent 建议、对象引用和 Tool Call 证据属于同一 tenant/correlation。
- IAOS 与 AESE 两仓测试、部署、runbook 和 evidence 完整。
