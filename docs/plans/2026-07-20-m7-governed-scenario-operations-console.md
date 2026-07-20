---
id: PLAN-M7-001
title: M7 受治理场景运行控制台实施计划
date: 2026-07-20
status: active
author: Codex + User
tags: [m7, orchestration, frontend, iaos, scenario-run]
---

# M7 受治理场景运行控制台实施计划

## 1. 目标与工期

在 **7 到 9 个工作日**内，让非研发用户仅通过 AESE 页面完成 HCTM 场景的预检、初始化、按幕推进、运行到结束、Agent 分析、结果验证和安全复位。

本计划是当前唯一 active plan。首版只支持内置 `hctm/order-expedite-01`，不提前实现参数化 A/B 实验和第二条故事。

## 2. 完成定义

- 新增无独立数据库的 AESE orchestration API，并复用现有 CLI 编排内核。
- 浏览器不直接调用 IAOS 写接口，不保存管理员凭据。
- 场景运行有明确状态机、阶段、run ID、plan hash、cursor 和幂等键。
- UI 可初始化、逐幕推进、运行到结束、分析、验证和复位。
- 所有业务写入仍经过 IAOS 权限、RLS、事务、Outbox、AI Tool 和审计边界。
- 页面刷新、AESE 服务重启、重复点击和网络断线后可以恢复。
- 运行结束保持 M6 在线不变量，复位后可以重新运行。
- API、Go、前端和 Playwright 验收及运行手册完整。

## 3. 实施切片

### O0 - 合同与编排内核重构（第 1-2 天）

- [ ] T1 核对 M6 CLI、IAOS scenario run/event/recommendation API 和当前部署基线。
- [ ] T2 定义 run 状态机、允许转换、stage result、错误码和可重试语义。
- [ ] T3 定义 pack act/event 到 `preflight/initialize/act-1..7/analyze/verify/reset` 的编译合同。
- [ ] T4 从 `cmd/aese` 提取可复用 application service，CLI 只保留参数解析和输出。
- [ ] T5 实现执行计划、plan hash、对象影响和下一允许动作计算。
- [ ] T6 为阶段编译、非法转换、重复调用和 pack 漂移增加单元测试。

验收：CLI 现有命令行为和输出合同不回归；同一 pack/version 总是生成相同 plan hash。

### O1 - AESE 薄编排 API（第 2-4 天）

- [ ] T7 增加 `aese serve` 启动入口、健康检查、超时、请求体上限和结构化日志。
- [ ] T8 实现 scenario list、run plan、run status 三个只读端点。
- [ ] T9 实现 initialize、advance 和 run-to-end 端点。
- [ ] T10 实现 analyze、verify、reset-plan 和 reset 端点。
- [ ] T11 实现 IAOS JWT 透传、profile/tenant 校验和权限失败关闭。
- [ ] T12 实现 idempotency key、expected cursor、run version 和 plan hash 并发保护。
- [ ] T13 实现一次性 reset confirmation token、过期和重放保护。
- [ ] T14 通过 IAOS run/snapshot/events/recommendations 重建状态，不依赖进程内存。
- [ ] T15 增加 CORS/同源代理、敏感字段脱敏、优雅关闭和服务重启恢复测试。

验收：服务不直连数据库、不持久化 JWT；所有写入均能关联调用者、tenant、run ID 和 correlation。

### O2 - IAOS 运行记录与权限补强（第 4-5 天）

- [ ] T16 审计 IAOS 现有 scenario apply/run/event 状态能否完整恢复 M7 阶段，缺失字段只做最小扩展。
- [ ] T17 固定 `scenario.run.read/execute/reset` 权限资源和 dev-user/普通用户行为。
- [ ] T18 确保每个阶段的 IAOS API 返回稳定 operation ref、cursor、correlation 和 committed/no-op 状态。
- [ ] T19 增加同 tenant 单 active writable run 约束和 409 冲突返回。
- [ ] T20 增加跨租户、权限不足、陈旧 cursor、并发推进和 reset 冲突测试。
- [ ] T21 同步 IAOS DES、code map、runbook，并从独立 worktree 合并和部署。

验收：普通只读用户能观察但不能运行；执行者不能越租户；并发请求不会推进两次或复位错误运行。

### O3 - 可视化运行控制台（第 5-7 天）

- [ ] T22 在联动中心增加“连接检查/运行场景”两个视图，不增加新的落地页。
- [ ] T23 实现 preflight 清单、dry-run 影响摘要和 active run 提示。
- [ ] T24 实现七幕 stepper、状态条、当前动作、下一动作和运行日志。
- [ ] T25 实现初始化、推进下一幕、运行到结束、运行分析和验证命令。
- [ ] T26 实现独立复位确认对话框、影响摘要和一次性 confirmation token。
- [ ] T27 实现执行中、成功、no-op、可重试错误、权限错误、冲突和服务不可用状态。
- [ ] T28 将控制台 run 状态与 Live snapshot/SSE 联动，但不直接修改画布业务状态。
- [ ] T29 支持刷新/重新登录/浏览器重开后的 active run 恢复。
- [ ] T30 完成键盘、焦点、ARIA、移动端布局和日志文本选择/复制。

验收：用户不使用 Token、curl 或 CLI 即可完成完整故事；危险动作和普通命令有明确视觉区分。

### O4 - 全链路验收与交付（第 7-9 天）

- [ ] T31 增加 API application service、handler、IAOS adapter 和恢复逻辑测试。
- [ ] T32 增加前端组件/状态测试及 Playwright 完整运行、逐幕运行、复位和恢复用例。
- [ ] T33 验证双击、并发浏览器、网络中断、服务重启和重复 idempotency key。
- [ ] T34 验证只读用户、无权限用户、tenant-other 和过期 reset token。
- [ ] T35 从 clean reset 执行 UI 全链，确认 22 事件、三 Agent、17 assertions 和 M6 KPI。
- [ ] T36 验证 CLI 与 UI 执行结果一致，且数据库/Outbox/Tool Call 无重复副作用。
- [ ] T37 采集 1440x900、1280x720、390x844 控制台和 Live 截图，检查画布非空与布局重叠。
- [ ] T38 编写 M7 runbook/evidence，更新 README、Architecture、Roadmap、Code Map、Agent Context 和 Progress Log。
- [ ] T39 部署 AESE orchestration API、frontend 和必要 IAOS 更新，记录 URL、健康检查和两仓 commit。

验收：自动化和人工验收全部通过；业务用户可独立运行并复位场景；失败路径不会留下不可解释的部分状态。

## 4. 每日可见成果

| 时间 | 可测试成果 |
| --- | --- |
| 第 1 天 | 状态机、阶段合同和 dry-run plan |
| 第 3 天 | API 可预检、初始化和推进单幕 |
| 第 4 天 | API 可运行到结束、分析、验证和复位 |
| 第 5 天 | 权限、幂等、并发和恢复合同关闭 |
| 第 7 天 | 页面可完成完整场景运行 |
| 第 8 天 | 断线、刷新、重复点击和移动端验收 |
| 第 9 天 | 部署、runbook、evidence 和版本提交 |

## 5. 用户验收故事

1. 计划员进入联动中心，选择华辰和订单加急场景。
2. 系统展示 IAOS 健康、权限、当前数据和初始化影响。
3. 用户确认初始化，Live 画布进入场景初始状态。
4. 用户逐幕推进，在第三至第五幕看到供应、设备和质量风险。
5. 用户选择“运行到结束”，在线 KPI 收敛为实发 11,700、缺口 300。
6. 用户运行三 Agent 分析，并能查看对象和 Tool Call 证据。
7. 系统验证 expected outcomes 并生成运行摘要。
8. 用户查看复位影响，确认后清理 L2/L3，保留华辰主数据。

## 6. 关键测试不变量

- 浏览器不能直接调用 IAOS scenario/simulation/business-action 写端点。
- 同一 stage + idempotency key 最多产生一次业务副作用。
- 只有前一阶段完成才能推进下一阶段。
- UI 状态只在 IAOS committed/no-op 和 snapshot cursor 证实后推进。
- reset 不删除 L1，且不能使用旧 confirmation token 重放。
- AESE 服务重启后，run 状态与重启前一致。
- tenant-other 和只读用户不能运行或复位 `tenant-hctm` 场景。

## 7. 后续路线

M7 完成后再进入 M8“参数化仿真实验”：

- 定义订单增量、供应延期、停机时长、来料不良率等受控参数。
- 串行运行 baseline 与方案 A/B，保存不可变结果摘要。
- 比较交付、库存、成本、质量和 Agent 建议差异。
- 在业务模型稳定后，再评估第二条故事和真实 LLM Agent Runtime。
