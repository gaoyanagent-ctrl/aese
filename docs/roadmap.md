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
| M3 可执行场景包 | JSON 场景包、校验器、IAOS apply/replay tracer | Planned / next | DES-001、M3 active plan |
| M4 异常场景运行 | 延期、设备故障、来料不良进入 IAOS 运行链 | Not started | - |
| M5 Agent MVP | 计划、质量、经营分析 Agent | Not started | - |
| M6 2D 企业沙盘 | 事件流、库存、产线、异常、Agent 建议 | Not started | - |

## 2. 当前阶段

当前进入 M3：可执行场景包。

M3 的最小成功标准：

1. HCTM 场景数据不再只存在于 Markdown，而是形成版本化 JSON pack。
2. `aese validate` 能离线发现 schema、引用、数量、时间线和幂等错误。
3. `aese inspect` 能输出对象数、事件数、依赖和故事摘要。
4. 至少一个受控路径可把基础数据导入 IAOS，并触发 `o2d.order.confirmed`。
5. 重复 apply/replay 不产生重复业务对象或重复事件副作用。
6. 提供 API、数据库、事件流和演示故事验收证据。

## 3. M3 范围

包含：

- 场景包 manifest、master data、initial state、events、expected outcomes。
- JSON Schema 和跨文件引用校验。
- Go CLI 的 `validate`、`inspect`、`apply --dry-run`。
- IAOS 兼容性评估和一个端到端 tracer bullet。
- 可重复 reset/verify 设计。

不包含：

- 完整 28 个 entity 全量落地。
- 18 类事件全部接线。
- 三个 Agent 实现。
- 2D/3D 前端。
- 通用仿真时间引擎或复杂离散事件调度。

## 4. M3 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| S0 | 工程治理和 M3 设计 | Completed |
| S1 | HCTM pack + JSON Schema | Pending |
| S2 | Go loader/validator/inspect CLI | Pending |
| S3 | IAOS compatibility report + adapter contract | Pending |
| S4 | 基础数据 dry-run/apply tracer | Pending |
| S5 | OrderConfirmed replay + O2D verification | Pending |
| S6 | Reset、重复执行和 runbook closeout | Pending |

## 5. 风险与依赖

- IAOS 当前平台后端未运行，M3 集成验证前需要恢复 `8082`。
- `/iaos/iaos-go` 当前主分支本地领先远程，任何集成开发必须使用独立 worktree 并确认基线。
- HCTM 业务字段与 IAOS 现有 legacy `sales_order` 物理模型存在差异，需要兼容性报告，不能直接假设可导入。
- IAOS 当前 O2D 自测硬编码 `tenant-001` 和旧订单，需要避免污染 HCTM tracer。
- 正式事件必须走 Outbox 或受治理 ingress，不能把 direct NATS 当最终实现。

## 6. 后续里程碑入口条件

进入 M4 前：

- M3 的 pack、validator、apply/replay 和 reset 均通过。
- `o2d.order.confirmed` tracer 在 `tenant-hctm` 下可复现。
- 事件关联、RLS 和幂等证据完整。

进入 M5 前：

- M4 三类异常均有结构化事件和可查询业务上下文。
- Agent 所需工具通过 IAOS Capability / AI Tool Registry 暴露。

进入 M6 前：

- 业务链和 Agent 输出稳定。
- 2D 沙盘的数据查询合同已经固定。
