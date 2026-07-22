# AESE Roadmap

本文件是 AESE 当前里程碑状态和下一步优先级的权威来源。

最后更新：2026-07-22。

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
| M7 受治理场景运行控制台 | 浏览器预检、初始化、逐幕运行、分析、验证和复位 | Completed | ADR-003、DES-005、PLAN-M7-001、M7 evidence |
| M8 AESE 2.0 基础 | 三态世界、确定性离散事件内核、IAOS 双向桥和最小 Genesis tracer | Completed | PLAN-M8-001、World Play runbook、两仓测试与部署证据 |
| X1 System Atlas 全景治理 | 最终完成体、当前状态、依赖与进展历史 | Completed | DES-006、IAOS DES-049、双端动态图谱 |

## 2. 当前阶段

M3、M3V、M4、M5、M6、M7 和跨里程碑的 X1 System Atlas 已完成。联动中心已支持联动检查与受治理场景运行，不依赖 CLI 完成 preflight、initialize、七幕推进、Agent 分析、verify 与 reset。

PLAN-M8-001 已完成，当前没有 active 主实施计划。M8 F0-F5 已交付机器合同、确定性内核、三态设备偏差 tracer、独立 IAOS World Bridge、Genesis pack、M7 兼容 adapter 与 World Play；持续扩展 World 服务不属于本里程碑。

M7 O0-O4 已完成。最终 `m7-acceptance-20260722-05` 从 clean reset 跑通编排 API 与 CLI 对照链：22 个事件、三 Agent、17 条离线业务断言、2 条在线 IAOS 断言和 M6 KPI 均通过；单 run 产生 9 次成功 Tool Call 与两套一致的 O2D Outbox 副作用，UI/CLI 均安全复位。AESE 8090/4173 与 IAOS 8082/3000 的本机部署和健康检查已记录在 M7 evidence。该基线由 M8 强制保留。

M7 的最小成功标准：

1. 浏览器不直接调用 IAOS 写端点，由 AESE 薄编排 API 复用现有 Go 内核。
2. 运行具有 run ID、阶段状态机、plan hash、cursor 和 idempotency key。
3. 用户可从页面初始化、逐幕推进、运行到结束、分析、验证和安全复位。
4. 刷新、断线、重复点击和 AESE 服务重启不会产生重复业务副作用。
5. 权限不足、跨租户、陈旧 cursor 和非法状态转换全部失败关闭。
6. UI 与 CLI 对同一 pack 产生一致的 22 事件、Agent 建议、断言和 KPI。

## 3. M8 当前范围

包含：

- World run、虚拟时钟、离散事件、规则版本、日志、快照和确定性 replay。
- World / IAOS / Actor Knowledge 三态与显式 discrepancy。
- 单台 `LAS-WLD-02` 设备退化的最小发现、登记、处置和关闭 tracer。
- observation / intent / committed outcome 的受治理 IAOS 合同。
- 设备、人员、物料和最小现金守恒。
- `hctm-genesis` world pack 和现有 `order-expedite-01` 兼容适配。
- World Play 的时间控制、三态对照和差异时间线。

不包含：

- 一次完成公司设立、工厂建设、APQP、财务和全生命周期的所有模块。
- 多工厂、多产品族、数百 Agent、自由文本长期记忆和自主批准。
- 3D、高精度物理仿真或完整设备数字孪生。
- 绕过 IAOS 权限/Capability/Outbox 的业务写入。
- 在 AESE 中镜像 IAOS 企业管理数据库。

M8 决策门与 F0-F5 的任务、验收和跨仓顺序以 PLAN-M8-001 为准。后续 Project Genesis 分解为 M9-M13，参数化分支实验后移至 M14。

## 4. M8 当前交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| F0 | 基线、状态所有权、存储和桥接合同冻结 | Completed |
| F1 | 确定性仿真内核 | Completed |
| F2 | 三态与设备偏差 tracer | Completed |
| F3 | 受治理 IAOS 双向桥 | Completed |
| F4 | Genesis world pack 与旧场景兼容 | Completed |
| F5 | World Play 最小界面与全链验收 | Completed |

## 5. M8 架构风险与依赖

- ADR-004 已 accepted；实现必须遵守独立 PostgreSQL database/账号/迁移边界，禁止跨库查询和外键。
- AESE 仿真事实和 IAOS 管理事实必须物理/逻辑隔离，禁止共享表和跨库写入。
- 原始设计稿中的 Spring Boot 仅是模块示意；AESE 实现继续使用 Go，与现有工具链保持一致。
- 世界结果必须由版本化规则和资源守恒计算；Agent 只提交 intent，不能直接改写 World State。
- Actor Knowledge 必须遵守可见范围，不能为了方便让 Agent 读取全量客观世界。
- IAOS 修改必须在独立 worktree，并先完成权限、RLS、Outbox、幂等和无部分写入设计。
- M7 22 事件、三 Agent、Preview/Live 与 reset 是强制回归门。

## 6. M8 完成条件

- ADR-004 accepted，World/IAOS/Knowledge 所有权和 World Store 选型明确。
- 相同 pack、规则版本、seed 和输入可重复产生相同 event log、state hash 与 KPI。
- 设备退化 tracer 可展示世界变化、IAOS 未登记、角色未知及其发现/关闭过程。
- IAOS 双向桥通过租户、权限、幂等、乱序、失败恢复和 Outbox 审计验收。
- Genesis pack 可离线验证、初始化、推进、复位和 replay，旧 M7 场景不回归。
- API/UI/runbook/evidence 与两仓 revision 完整。

## 7. M7 已完成范围（保留基线）

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

## 8. M7 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| O0 | 状态机、阶段合同和编排内核重构 | Completed |
| O1 | AESE 薄编排 API | Completed |
| O2 | IAOS 运行记录、权限和并发补强 | Completed |
| O3 | 可视化场景运行控制台 | Completed |
| O4 | 全链路、恢复、安全和三视口验收 | Completed |

## 9. 历史风险与依赖

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

## 10. M7 完成条件

- 非研发用户可从浏览器完整运行并复位第一条故事。
- UI 状态只在 IAOS committed/no-op 和 snapshot cursor 证实后推进。
- 双击、并发、断线、刷新和服务重启均不重复推进阶段。
- reset 影响可预览，一次性 confirmation token 不能重放，L1 始终保留。
- IAOS 与 AESE 两仓权限、测试、部署、runbook 和 evidence 完整。

## 11. M6 完成证据

M6 已满足：

- 在线库存、完工、发运和订单状态满足 expected outcomes 与库存守恒。
- snapshot 与 cursor 来自一致观察边界，断线补发和事件去重可复现。
- Preview/Live 数据源明示，错误时不静默混用。
- Agent 建议、对象引用和 Tool Call 证据属于同一 tenant/correlation。
- IAOS 与 AESE 两仓测试、部署、runbook 和 evidence 完整。
