---
id: DES-005
title: M7 受治理场景运行控制台
date: 2026-07-20
status: completed
author: Codex + User
tags: [m7, orchestration, frontend, iaos, scenario-run]
---

# M7 受治理场景运行控制台

## 1. 产品目标

把 AESE 从“可观察的单故事沙盘”推进为“业务用户可安全运行的仿真环境”。用户在联动中心选择华辰场景后，可以看到运行影响、初始化环境、按七幕推进、运行 Agent、验证结果和复位环境，全程不使用终端。

## 2. 首版用户流程

```text
选择租户与场景
-> 联动健康检查
-> 生成 dry-run 计划和影响摘要
-> 确认初始化
-> 按幕推进 / 运行到结束
-> 查看 Live 画布与实时事件
-> 运行三 Agent 分析
-> 验证 expected outcomes
-> 查看运行报告
-> 确认并复位
```

场景控制和业务观察使用同一 run ID、pack version、scenario key 和 correlation ID。

## 3. 场景阶段模型

首版把 22 个事件编译为稳定阶段：

| 阶段 | 内容 | IAOS 动作 |
| --- | --- | --- |
| `preflight` | pack、IAOS、权限、租户和对象合同检查 | 只读 |
| `initialize` | L1/L2 数据、订单和 O2D 基线 | scenario apply + decompose |
| `act-1` | 客户订单追加 | 已初始化事实/事件确认 |
| `act-2` | 初始 MRP 与采购 | 受治理 replay |
| `act-3` | 供应延期 | simulation ingress |
| `act-4` | 设备停机 | simulation ingress |
| `act-5` | 备选供应与来检失败 | scenario/business action + simulation ingress |
| `act-6` | 生产完成与完工入库 | scenario business action |
| `act-7` | 两次发运与交付复盘 | scenario business action |
| `analyze` | 三 Agent 读取在线事实并发布建议 | AI Tool Registry + recommendation API |
| `verify` | expected outcomes 与在线断言 | 只读 |
| `reset` | 清理场景 L2/L3 并保留 L1 | scenario reset |

阶段编译来自 pack 中的 act/event metadata，handler 只执行 allowlist stage 类型。首版不允许用户上传脚本或编辑 HTTP 请求。

## 4. API 合同

建议入口：

```text
GET  /api/aese/v1/scenarios
POST /api/aese/v1/runs/plan
POST /api/aese/v1/runs
GET  /api/aese/v1/runs/:run_id
POST /api/aese/v1/runs/:run_id/advance
POST /api/aese/v1/runs/:run_id/run-to-end
POST /api/aese/v1/runs/:run_id/analyze
POST /api/aese/v1/runs/:run_id/verify
POST /api/aese/v1/runs/:run_id/reset-plan
POST /api/aese/v1/runs/:run_id/reset
```

所有写请求要求：

- `Authorization: Bearer <IAOS JWT>`。
- 明确 `Idempotency-Key`。
- run version 或 expected cursor，用于防止陈旧页面写入。
- plan hash，保证用户确认的影响与实际执行一致。
- reset 额外要求一次性 confirmation token。

API 使用问题详情结构返回 stage、错误码、可否重试和 IAOS correlation。成功响应包含当前阶段、已完成阶段、snapshot cursor、对象影响和下一允许动作。

## 5. 状态与恢复

运行状态闭集：

```text
planned -> initializing -> ready -> running -> awaiting_analysis
-> analyzing -> awaiting_verification -> completed -> reset
```

失败不单独成为终态，而是记录 `last_error` 和 `retryable`，当前阶段保持未完成。浏览器刷新或 AESE 服务重启后，通过 IAOS run/event/recommendation 状态重建视图。

同一 tenant + scenario 首版只允许一个可写运行。第二个运行请求返回 409 和当前 active run，避免固定自然键场景并发互相覆盖。

## 6. 控制台交互

在现有“联动中心”内增加“运行场景”视图：

- 场景状态条：未初始化、运行中、待分析、已完成、可复位。
- 预检清单：身份、权限、Platform/O2D、pack、snapshot 和 active run。
- 影响摘要：将新增、复用、更新和删除的对象数量。
- 七幕 stepper：当前幕、已完成幕、下一步动作和实时事件数量。
- 主命令：初始化、推进下一幕、运行到结束、运行分析、验证、复位。
- 运行日志：时间、阶段、actor、correlation、结果和错误；业务文本可选择复制。
- 危险确认：复位使用独立确认对话框，不与普通运行按钮混放。

Live 画布继续只消费 IAOS snapshot/cursor/SSE。控制台动作成功后不手工修改画布状态，而是等待在线事实到达。

## 7. 失败和并发语义

- 重复点击同一动作使用相同 idempotency key，必须返回同一结果。
- 浏览器断线不取消已提交的 IAOS 事务；恢复后查询 run 状态。
- 阶段中部分 IAOS 调用失败时不得把阶段标记完成；可重试动作只执行缺失步骤。
- plan hash 或 expected cursor 不匹配返回 409，要求重新预检。
- 非法阶段跳转返回 409，不允许从未初始化直接发运或分析。
- 其他租户不能查看、推进或复位 HCTM run。

## 8. 非目标

- 不支持任意参数化实验、A/B 方案比较或并行分支；建议作为 M8。
- 不支持第二条故事、多工厂和用户自定义场景编辑器。
- 不引入 AESE 业务数据库或通用任务队列。
- 不实现真实 LLM 自主执行和自动批准建议。
- 不补造未经批准的成本实际金额。

## 9. 完成标准

- 新用户从浏览器完成 preflight、initialize、七幕推进、analyze、verify 和 reset。
- CLI 和 UI 对同一 pack 产生一致对象、事件、KPI 和断言结果。
- 重复点击、刷新、AESE 服务重启和 SSE 断线均不产生重复业务副作用。
- 权限不足、跨租户、陈旧 cursor、plan hash 变化和非法阶段转换全部失败关闭。
- 所有写入可在 IAOS 审计、scenario run、Outbox、Tool Call 和 correlation 中追溯。
- 2D Live 最终仍显示需求 12,000、实发 11,700、缺口 300，reset 后回到可重新初始化状态。
