---
id: DES-004
title: M6 在线 2D 企业沙盘架构
date: 2026-07-20
status: completed
author: Codex + User
tags: [m6, frontend, live, iaos, observation]
---

# M6 在线 2D 企业沙盘架构

## 1. 背景

M3V 已验证 2D 沙盘的信息布局和交互，但其状态来自 `preview.json`。M4 已把设备停机、供应延期和来检失败送入 IAOS，M5 已让三个 Agent 通过 AI Tool Registry 读取在线事实。当前仍有三个断点：

- 前端没有 `IaosScenarioDataSource`，看不到 IAOS 当前状态。
- IAOS 缺少完工入库、两次发运和成本影响事实，经营分析只能返回 `partial`。
- 通用 `/api/v1/events/stream` 是无游标的 NATS PubSub 转发，缓冲满会丢事件，断线后不能补发，不足以作为沙盘事实源。

## 2. 目标

把同一套 2D 沙盘组件接到 IAOS 的租户级在线事实，让用户可以：

1. 在 Preview 和 Live 间明确切换。
2. 从 IAOS 快照恢复当前企业状态。
3. 按持久游标接收和补齐场景事件。
4. 查看带业务对象和 Tool Call 证据的 Agent 建议。
5. 在线验证 12,000 件需求、11,700 件实发和 300 件缺口。

## 3. 核心决策

### 3.1 静态布局与在线事实分离

AESE 继续拥有苏州基地 A 线布局、节点业务编码和视觉映射。IAOS 只返回运行事实，不保存 React Flow 坐标或颜色。

```text
AESE layout definition
        +
IAOS scenario snapshot / ordered events / recommendations
        -> IaosScenarioDataSource
        -> SandboxScenario / LiveScenarioState
        -> existing 2D components
```

### 3.2 快照为真，事件为增量

Live 首次进入必须先读取一致性快照，再从快照返回的 `cursor` 继续消费事件。SSE 断线后，客户端先通过 `after=<cursor>` 补齐持久事件，再恢复流式连接。不能只依赖浏览器连接期间收到的 NATS 消息计算库存或交付结果。

### 3.3 场景级观察 API

IAOS 新增只读、租户隔离的窄合同：

```text
GET /api/v1/scenarios/hctm/order-expedite-01/snapshot
GET /api/v1/scenarios/hctm/order-expedite-01/events?after=<cursor>&limit=<n>
GET /api/v1/scenarios/hctm/order-expedite-01/events/stream?after=<cursor>
GET /api/v1/scenarios/hctm/order-expedite-01/recommendations
```

所有端点要求 JWT、`scenario.read` 权限和 tenant context。响应只包含该场景 ownership/correlation 范围内的数据。通用 `/api/v1/events/stream` 保留给操作台监控，但不作为 M6 的恢复合同。

### 3.4 在线业务事实必须来自受治理动作

事件 17-22 中的工单释放、工序完成、完工入库和发运不是外生异常，不进入 simulation ingress。它们通过 IAOS Capability/场景业务动作处理，并满足：

- 业务状态、库存事务、场景事件日志和 Outbox 同事务提交。
- 以 tenant、event ID 和 idempotency key 去重。
- 库存不足时失败关闭，不能超发。
- 重复调用返回 no-op，不重复入库、扣库或发运。

成本只做 HCTM MVP 的成本影响事实，不建设总账或完整标准成本引擎。成本基线和金额必须先在场景规格中显式批准，未批准前经营分析继续标记成本部分 `partial`。

### 3.5 Agent 建议是可审计在线结果

M5 recommendation envelope 保持不变。M6 将建议作为 IAOS 场景观察结果持久化，并校验其引用的 Tool Call 均属于同一 tenant/correlation。建议保持 `suggested`、`requires_human_confirmation=true`，在线沙盘不自动执行建议。

## 4. 在线视图模型

在现有 `SandboxScenario` 基础上新增：

```ts
interface LiveScenarioState {
  snapshotVersion: string;
  observedAt: string;
  cursor: string;
  completeness: "complete" | "partial";
  entities: SandboxEntity[];
  kpis: KpiSnapshot;
  events: ScenarioEvent[];
  agentOutputs: AgentOutput[];
  gaps: string[];
}

interface LiveScenarioDataSource extends ScenarioDataSource {
  refresh(): Promise<LiveScenarioState>;
  connect(after: string, signal: AbortSignal): AsyncIterable<ScenarioEvent>;
}
```

Preview 仍使用确定性播放 reducer。Live 使用 snapshot + ordered event reducer，并提供“跟随实时”“刷新”“重连”操作；Live 不显示倍速播放，也不能回写业务状态。

## 5. 一致性与恢复

- 持久 cursor 在 tenant/场景内严格递增，event ID 唯一。
- 快照返回的数据与 cursor 必须来自同一数据库观察边界。
- 客户端按 event ID 去重，拒绝 cursor 倒退和未知事件类型。
- SSE 心跳不推进 cursor。
- 网络中断时保留最后可信快照并显示“已断开”，不得静默切回 Preview。
- 恢复失败时允许用户显式选择 Preview，但两种数据源必须持续明示。

## 6. 安全边界

- 浏览器 token 只来自登录态或显式开发配置，不写入场景包和仓库。
- 服务端从认证上下文确定 tenant，不接受请求体覆盖 tenant。
- 新场景流不得继承通用 SSE 对 `tenant-001` 的全局订阅特例；平台管理员跨租户观察需要独立权限与端点。
- 快照、事件和建议查询均使用显式 tenant predicate 与 RLS。
- UI 保持只读；apply、replay 和 Agent 运行继续由受治理命令/API 触发。

## 7. 非目标

- 不建设通用数字孪生平台或多工厂自由建模器。
- 不把全部 18 类事件一次性在线化，只覆盖第一条故事所需合同。
- 不实现真实 LLM 自主执行、自动排产或自动采购审批。
- 不建设完整财务总账、成本核算和利润中心。
- 不移除 Preview；Preview 是演示降级与回归基线。

## 8. 完成标准

- Live 模式从 IAOS 快照显示订单、库存、采购、设备、质量、工单、完工和发运事实。
- 断开并恢复 SSE 后不丢事件、不重复应用状态。
- 在线最终状态为需求 12,000、入库 10,500、累计实发 11,700、缺口 300、`partially_shipped`。
- 三 Agent 建议带 correlation 和 Tool Call 证据；交付事实补齐后经营分析不再报告 shipment data gap。
- 成本金额只有在成本基线批准并落库后才标记完整，否则界面明确显示成本数据缺口。
- Preview/Live 在桌面和移动端均可区分、可操作、无状态串扰。
