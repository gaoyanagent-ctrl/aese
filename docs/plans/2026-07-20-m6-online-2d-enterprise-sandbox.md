---
id: PLAN-M6-001
title: M6 在线 2D 企业沙盘实施计划
date: 2026-07-20
status: active
author: Codex + User
tags: [m6, frontend, iaos, live, hctm]
---

# M6 在线 2D 企业沙盘实施计划

## 1. 目标与工期

在 **8 到 10 个工作日**内，把 M3V 静态沙盘升级为 Preview/Live 双模式，跑通“IAOS 业务事实 -> 持久场景事件 -> 在线快照 -> 2D 画布 -> Agent 建议”的完整观察链。

本计划是当前唯一 active plan。范围只覆盖华辰苏州基地、电池冷却板 A 线和 `order-expedite-01`，不扩展第二工厂、第二产品族或 3D。

## 2. 完成定义

- IAOS 提供租户隔离的场景 snapshot、cursor event query、SSE 和 recommendation query。
- canonical 事件 17-22 通过 IAOS 受治理业务动作形成完工、入库和两次发运事实。
- 在线库存不超发，重复 replay 不重复入库、扣库、发运或发布事件。
- `IaosScenarioDataSource` 使用现有 `SandboxScenario` 视觉合同，Preview 继续可用。
- Live 模式显示连接状态、最后更新时间、数据完整度、重连和显式刷新。
- 最终在线 KPI 为需求 12,000、累计可供 11,700、累计实发 11,700、期末成品库存 0、缺口 300。
- 三 Agent 建议来自在线工具调用并显示 evidence；发运事实补齐后移除 shipment gap。
- 桌面、移动端、断线恢复、跨租户和全链路回归测试通过。

## 3. 交付切片

### L0 - 合同冻结与基线（第 1 天）

- [ ] T1 核对 IAOS 主分支、运行服务、`tenant-hctm` 数据和 M3-M5 重放基线。
- [ ] T2 明确事件 17-22 的业务对象、自然键、状态机、库存影响和 event payload。
- [ ] T3 补充并批准 HCTM 最小成本影响基线；无法确认金额时将成本完整度明确留为 `partial`。
- [ ] T4 固定 snapshot、event cursor、SSE 和 recommendation JSON 合同及错误码。
- [ ] T5 为 IAOS 跨仓库实现建立独立 worktree、DES 和 contract test 骨架。

验收：所有新增字段有业务来源；Preview expected outcomes 与 Live 目标逐项映射；不存在用 Preview 数字替代在线事实的路径。

### L1 - 完工、入库与发运事实闭环（第 2-4 天）

- [ ] T6 在 IAOS 定义受治理生产完成/完工入库业务动作，稳定解析工单、产品、仓库和批次。
- [ ] T7 同事务写入工单状态、完工入库单、库存事务、库存余额、场景事件日志和 Outbox。
- [ ] T8 定义受治理发运动作，稳定解析订单、产品、成品库存和发运单。
- [ ] T9 实现两次发运、库存扣减、累计发运、短缺和订单 `partially_shipped` 状态更新。
- [ ] T10 实现 event ID/idempotency key 去重、碰撞检测、库存不足失败关闭和事务回滚测试。
- [ ] T11 扩展 AESE replay/client，将 canonical 事件 17-22 路由到正式业务动作，保持默认 dry-run。
- [ ] T12 扩展 reset/verify，清理并验证新增 L2/L3 对象，不删除 L1 主数据。

验收：首次 replay 形成 10,500 件完工入库和 9,000 + 2,700 件发运；第二次 replay 全部 no-op；任何路径均不能把库存扣成负数。

### L2 - 场景观察 API（第 4-6 天）

- [ ] T13 建立 tenant-scoped 场景事件日志和严格递增 cursor，业务写入与日志/Outbox 同事务。
- [ ] T14 实现 snapshot API，聚合场景 ownership 内订单、库存、采购、设备、检验、工单、完工和发运状态。
- [ ] T15 在 snapshot 中返回 `observed_at`、`cursor`、`completeness`、`gaps` 和五项 KPI。
- [ ] T16 实现 `events?after=` 持久补发 API，支持 limit、顺序、event ID 去重和相关 correlation 过滤。
- [ ] T17 实现场景 SSE，支持 `after`、heartbeat、断开取消和慢客户端保护；SSE 仅作增量通道。
- [ ] T18 确保新场景流不继承 `tenant-001` 全局订阅特例，增加权限、RLS、跨租户和多副本测试。

验收：从任意已知 cursor 断开后可完整补发；snapshot 后建立 SSE 不丢失中间事件；其他租户得到空结果或 404，不能观察 HCTM 数据。

### L3 - Agent 建议在线化（第 6-7 天）

- [ ] T19 为 M5 recommendation envelope 增加受治理持久化合同，校验 tenant、correlation 和 Tool Call 归属。
- [ ] T20 扩展 `agent-run --apply`，在 9 次工具调用成功后幂等发布三 Agent 建议。
- [ ] T21 实现 recommendation query，并将建议版本、完整度、data gaps、对象引用和 Tool Call IDs 纳入 snapshot。
- [ ] T22 扩展经营分析：使用真实完工/发运事实计算 11,700 实发和 300 缺口；成本结论严格服从 T3 的基线完整度。

验收：重复 Agent run 更新同一分析版本或 no-op，不生成重复建议；跨租户 Tool Call 引用失败关闭；UI 能追溯每条建议的证据。

### L4 - `IaosScenarioDataSource` 与 Live UX（第 7-9 天）

- [ ] T23 扩展前端类型，区分布局定义、Preview playback 和 Live observation state。
- [ ] T24 实现 IAOS HTTP client、认证配置、错误映射和 `IaosScenarioDataSource`。
- [ ] T25 实现 snapshot-first 加载、cursor 补发、SSE 连接、event ID 去重和指数退避重连。
- [ ] T26 增加 Preview/Live 分段控制，明确显示数据源、连接状态、最后更新时间和完整度。
- [ ] T27 Live 模式将倍速/重置替换为跟随实时、刷新和重连；Preview 控制保持不变。
- [ ] T28 将在线事件映射到 14 节点/13 连线状态，并从在线实体计算对象详情和 KPI。
- [ ] T29 展示在线 Agent 建议、证据、建议状态和数据缺口，不自动执行建议。
- [ ] T30 实现加载、无权限、服务不可用、数据不完整、断线和恢复状态；不得静默降级成 Preview。

验收：同一页面可明确切换两种模式；Live 只显示 IAOS 事实；刷新和重连不清空最后可信状态，也不污染 Preview reducer。

### L5 - 全链路验收与交付（第 9-10 天）

- [ ] T31 增加 TypeScript adapter/reducer/component 测试和 Go API/client/replay 测试。
- [ ] T32 增加 Playwright Preview/Live、断线恢复、空数据、无权限和移动端用例。
- [ ] T33 从 reset 开始执行 apply、O2D、三异常、完工、发运、Agent run 和 Live UI 全链回放。
- [ ] T34 验证重复执行、事件补发、租户隔离、Tool Call 归属、库存守恒和最终 KPI。
- [ ] T35 采集 1440x900、1280x720、390x844 的 Live 截图和非空画布像素证据。
- [ ] T36 编写 M6 runbook/evidence，更新 README、Architecture、Roadmap、Code Map、Agent Context 和 Progress Log。
- [ ] T37 重新部署 IAOS Platform/O2D 与 AESE frontend，记录可访问 URL 和版本 commit。

验收：全部自动化测试通过；在线沙盘完成七幕故事的真实状态闭环；用户可以按 runbook 独立复现。

## 4. 每日可见成果

| 时间 | 可测试成果 |
| --- | --- |
| 第 1 天 | 在线合同、对象状态机和成本完整度决策 |
| 第 3 天 | 完工入库与第一条发运动作可通过 API 验证 |
| 第 4 天 | 两次发运、库存守恒和 300 件缺口在线成立 |
| 第 6 天 | snapshot、cursor 补发和场景 SSE 可测试 |
| 第 7 天 | 在线三 Agent 建议及 Tool Call 证据可查询 |
| 第 9 天 | 2D 页面 Preview/Live 双模式可操作 |
| 第 10 天 | 全链回放、断线恢复、三视口验收与部署完成 |

## 5. 测试重点

业务不变量：

- `opening FG 1,200 + receipt 10,500 - shipment 11,700 = ending FG 0`。
- `demand 12,000 - shipment 11,700 = shortage 300`。
- 发运量不能超过可用库存，失败请求不产生部分业务记录或 Outbox。
- 相同 event/idempotency key 重放不改变库存、工单、发运和事件数量。

在线一致性：

- snapshot cursor 与其状态一致；之后的事件只应用一次。
- SSE 断线期间产生的事件可通过 cursor query 补齐。
- 乱序、重复、未知类型和 cursor 倒退均失败关闭或触发重新取快照。
- Live 错误不会让界面把 Preview 结果标成在线结果。

安全与审计：

- snapshot、events、SSE、recommendations 均验证 JWT、权限和 tenant。
- tenant-other 无法读取 `tenant-hctm` 的对象、事件、建议或 Tool Call。
- 每个完工、入库、发运和 Agent 结果可追溯 actor、correlation、event/tool call 和时间。

## 6. 明确非目标

- 不做第二条演示故事、多工厂或多产品族。
- 不做通用流程编辑器、布局编辑器或自由拖拽保存。
- 不做真实 LLM 自主 Agent、自动审批或自动执行建议。
- 不做完整财务核算；成本只覆盖已批准的 HCTM 演示影响事实。
- 不以轮询多个动态实体接口代替场景快照合同。

## 7. 主要风险与控制

- **当前 SSE 会丢事件**：新增持久 cursor 查询，SSE 只负责低延迟通知。
- **业务动作被误塞进 simulation ingress**：完工和发运必须走 Capability/场景业务动作。
- **聚合快照变成 HCTM 硬编码平台功能**：API 按 scenario ownership、稳定 entity contract 和 projection 组织，视觉布局仍留在 AESE。
- **成本数字被虚构**：T3 未批准前保持 `partial`，不得为了“全绿”填入假金额。
- **跨仓库开发污染主分支**：IAOS 修改按其 AGENTS 要求使用独立 worktree、测试、code map 同步和部署验证。
