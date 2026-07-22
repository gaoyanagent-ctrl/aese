---
id: DES-008
title: AESE World 与 IAOS 三段式桥接合同
date: 2026-07-22
status: approved
author: Codex + User
tags: [aese, iaos, world-bridge, observation, intent, outcome]
---

# AESE World 与 IAOS 三段式桥接合同

## 1. 决策摘要

AESE 与 IAOS 使用 `observation -> intent -> committed_outcome` 三段式合同。IAOS 持久化租户级 World Bridge journal 并提供单调 cursor；AESE 通过 cursor query 恢复，SSE/Outbox 只作为低延迟通知，不作为唯一事实来源。

三段语义固定为：

- `observation`：世界中可被特定角色感知的事实，不等同于 IAOS 业务台账已经变化。
- `intent`：人员或 Agent 通过 IAOS 权限、流程和 Capability 提出的行动意图，不等同于世界结果已经发生。
- `committed_outcome`：IAOS 事务已提交或确定 no-op 的管理结果，是 AESE 计算世界后果的唯一 IAOS 输入。

AESE 不接收 `committed=false` 作为世界输入，不根据 HTTP 超时猜测提交结果，也不让 Agent 直接提交任意世界状态 patch。

## 2. 选择理由

采用 journal + cursor，而不是 webhook-only 或 direct NATS：

- 对齐现有 M6 snapshot/cursor/SSE 恢复模型。
- AESE 或网络中断后可以从持久 cursor 补发。
- IAOS 事务提交与 journal/outbox 可以原子完成。
- 重复投递、乱序通知和进程重启不会重复产生世界后果。
- 不要求 IAOS 保存 AESE 回调凭据，也不把 NATS 暴露为正式写入口。

## 3. 公共 envelope

三个消息共享以下字段：

| 字段 | 类型 | 规则 |
| --- | --- | --- |
| `schema_version` | string | 首版固定 `1.0`；未知 major 失败关闭 |
| `message_id` | string | 生产方生成的全局稳定 ID，建议 UUID/ULID |
| `kind` | enum | `observation`、`intent`、`committed_outcome` |
| `tenant_id` | string | 必填，必须与 IAOS 认证 tenant 完全一致 |
| `world_pack_key` | string | 首版 `hctm-genesis` |
| `world_pack_version` | string | 产生消息的不可变 pack 版本 |
| `world_run_id` | string | AESE WorldRun 稳定 ID |
| `branch_id` | string | 首版 `main`；为后续 fork 保留 |
| `sim_occurred_at` | RFC 3339 | 世界事实或意图在虚拟时间中的发生时间，时区语义为 `Asia/Shanghai` |
| `recorded_at` | RFC 3339 | 持有方服务端写入的真实 UTC 时间；请求方不得指定权威值 |
| `correlation_id` | string | 一次业务目标或问题处理链保持不变 |
| `causation_id` | string | 直接上游 message/world event ID；根消息可为空 |
| `idempotency_key` | string | 同 tenant + kind 范围唯一 |
| `producer` | object | `system`、`component`、`version`，不得伪造 actor 身份 |
| `subject_ref` | object | 稳定对象引用，禁止环境 UUID 作为唯一引用 |
| `payload_type` | string | allowlist 类型，如 `equipment.condition.observed.v1` |
| `payload` | object | 由 payload_type 对应 JSON Schema 严格校验 |

除 `recorded_at` 外，这些字段必须出现在持久 journal entry 中。写入 observation 的请求不携带 `recorded_at`，由 IAOS 在事务中生成；intent 和 committed outcome 同样由 IAOS 生成权威 `recorded_at`。

`subject_ref` 首版结构：

```json
{
  "namespace": "hctm",
  "type": "equipment",
  "code": "LAS-WLD-02"
}
```

IAOS 在需要时解析自己的 UUID，并只在响应的 `record_refs` 中返回；AESE World State 始终以稳定业务编码关联。

## 4. Observation 合同

### 4.1 语义

Observation 代表“某个世界事实现在对哪些角色可见”。它可以生成 IAOS inbox、任务、告警或调查入口，但默认不得直接把设备台账改为故障、修改库存或生成财务凭证。

首版入口：

```text
POST /api/v1/world-bridge/observations
```

必需权限：`world.observation.ingest`。调用者使用短期服务身份或用户委托身份；tenant 从认证上下文确认，body 中 tenant 不匹配返回 403。

### 4.2 Payload

```json
{
  "schema_version": "1.0",
  "message_id": "obs-01J...",
  "kind": "observation",
  "tenant_id": "tenant-hctm",
  "world_pack_key": "hctm-genesis",
  "world_pack_version": "0.1.0",
  "world_run_id": "world-run-01J...",
  "branch_id": "main",
  "sim_occurred_at": "2026-07-08T10:15:00+08:00",
  "correlation_id": "corr-equipment-degradation-01",
  "causation_id": "world-event-vibration-rise-01",
  "idempotency_key": "world-run-01J:obs:vibration-rise:1",
  "producer": {
    "system": "aese",
    "component": "world-runtime",
    "version": "0.1.0"
  },
  "subject_ref": {
    "namespace": "hctm",
    "type": "equipment",
    "code": "LAS-WLD-02"
  },
  "payload_type": "equipment.condition.observed.v1",
  "payload": {
    "fact_type": "vibration_above_baseline",
    "measurement": {"value": "7.20", "unit": "mm/s", "scale": 2},
    "threshold": {"value": "6.00", "unit": "mm/s", "scale": 2},
    "source_ref": "sensor:VIB-LAS-WLD-02-01",
    "available_at": "2026-07-08T10:16:00+08:00",
    "recipient_refs": ["position:EQUIPMENT-ENGINEER-SZ"],
    "confidence": "0.95",
    "visibility_scope": "assigned_recipients"
  }
}
```

数值以字符串 + unit + scale 表达，不使用 JSON 浮点作为金额或守恒数量。

### 4.3 响应

成功返回 `201`；完全重复返回 `200`：

```json
{
  "message_id": "obs-01J...",
  "journal_cursor": 41,
  "accepted": true,
  "duplicate": false,
  "recorded_at": "2026-07-22T04:00:00Z",
  "operation_ref": "world-bridge-observation/obs-01J..."
}
```

`accepted` 只证明 IAOS 已持久化观察，不证明设备台账或世界状态已经改变。

## 5. Intent 合同

### 5.1 语义

Intent 由人员或 Agent 在 IAOS 内通过 Capability/Process 产生。调用者不能直接指定 actor、批准状态或权限结果；IAOS 从认证、岗位、审批和 Capability execution 生成这些字段。

首版不要求 AESE 调用公开的 intent 写端点。IAOS 内部 Capability 使用受治理原语原子写入 intent；管理 UI 可以通过业务 Capability 创建：

```text
POST /api/v1/capabilities/:capability_key/execute
```

示例 intent journal payload：

```json
{
  "kind": "intent",
  "message_id": "intent-01J...",
  "causation_id": "obs-01J...",
  "payload_type": "equipment.inspection.requested.v1",
  "payload": {
    "actor_ref": "employee:EMP-SZ-EAM-001",
    "role_ref": "position:EQUIPMENT-ENGINEER-SZ",
    "action": "inspect_equipment",
    "target_ref": {"namespace": "hctm", "type": "equipment", "code": "LAS-WLD-02"},
    "requested_start_at": "2026-07-08T10:30:00+08:00",
    "preconditions": {
      "world_state_hash": "sha256:...",
      "valid_until": "2026-07-08T11:00:00+08:00"
    },
    "authority": {
      "capability_key": "eam.request-equipment-inspection",
      "process_run_id": "",
      "approval_status": "not_required"
    }
  }
}
```

Intent 只允许 allowlist `action` 和严格 schema，不允许 SQL、脚本、HTTP URL、任意 state patch 或任意 Capability key。

## 6. Committed Outcome 合同

### 6.1 语义

Committed Outcome 只能由 IAOS 在业务事务成功提交后写入 journal，并与业务记录、审计和 Outbox 同事务。它是 AESE 可以消费并计算 `world_consequence` 的唯一 IAOS 消息。

允许状态只有：

- `committed`：业务变化已经提交。
- `no_op`：同一意图已满足或幂等重复，没有新增副作用。

失败、拒绝、回滚和待审批保留在 IAOS intent/process/capability 状态中，不伪装为 committed outcome。AESE 可以展示这些状态，但不得据此改变 World State。

### 6.2 Payload

```json
{
  "kind": "committed_outcome",
  "message_id": "outcome-01J...",
  "causation_id": "intent-01J...",
  "payload_type": "equipment.inspection.scheduled.v1",
  "payload": {
    "intent_id": "intent-01J...",
    "status": "committed",
    "committed_at": "2026-07-22T04:01:12Z",
    "operation_ref": "capability-run/caprun-01J...",
    "iaos_cursor": 43,
    "record_refs": [
      {
        "entity": "maintenance_work_order",
        "id": "2d2e...",
        "kind": "work_order",
        "committed": true
      }
    ],
    "result_facts": {
      "inspection_order_no": "EAM-WO-20260708-001",
      "scheduled_start_at": "2026-07-08T10:30:00+08:00"
    }
  }
}
```

`record_refs[*].committed` 必须全部为 true。`result_facts` 只包含该 outcome schema 允许的结果，不是完整 IAOS 记录镜像。

## 7. Journal 读取与恢复

```text
GET /api/v1/world-bridge/entries?world_run_id={id}&branch_id=main&after={cursor}&limit=200
GET /api/v1/world-bridge/entries/stream?world_run_id={id}&branch_id=main&after={cursor}
```

必需权限：`world.bridge.read`。响应包含 `items`、`next_cursor` 和 `has_more`。cursor 在 tenant journal 内单调递增；消息的 `sim_occurred_at` 不能替代 cursor 排序。

AESE 消费规则：

1. 从本地已提交 `last_iaos_cursor` 开始分页读取。
2. 按 cursor 顺序验证 schema、tenant、run、branch 和因果引用。
3. Observation/intent 只更新桥接投影或角色认知，不直接改变客观世界。
4. 只有 committed outcome 进入世界规则 reducer。
5. World event、world state、outcome message ID 和新 cursor 在 AESE PostgreSQL 单事务提交。
6. 崩溃重试时由 message ID 和 cursor 去重。

SSE 断开、缓冲溢出或 NATS 重投后始终回到步骤 1，不从内存推断缺失消息。

## 8. 幂等、并发与错误

- 唯一键：`(tenant_id, kind, idempotency_key)` 和 `(tenant_id, message_id)`。
- 相同 idempotency key + 相同 canonical input hash：返回原结果并标记 `duplicate=true`。
- 相同 idempotency key + 不同 input hash：`409 idempotency_conflict`。
- world run/branch 不存在或未绑定 tenant：`404 world_run_binding_not_found`。
- tenant 与认证上下文不一致：`403 tenant_mismatch`。
- 不支持的 payload type/version：`422 payload_type_unsupported`。
- 过期 precondition 或 world state hash：intent 返回 `409 world_precondition_stale`，要求重新观察和决策。
- journal cursor 超出当前边界：`409 cursor_ahead`；过旧 cursor 仍必须可分页恢复，除非有显式保留策略和 snapshot handoff。
- 请求体上限首版 256 KiB；未知字段失败关闭。

HTTP 5xx/超时表示结果未知。调用方必须用同一 idempotency key 重试或按 message ID 查询，不得生成新 key 猜测执行。

## 9. 权限与审计

| 权限 | 主体 | 能力 |
| --- | --- | --- |
| `world.observation.ingest` | AESE World 服务身份 | 提交 allowlist observation |
| `world.bridge.read` | AESE World 服务身份、授权观察者 | 按 tenant/run 读取 journal |
| `world.intent.create` | 岗位用户或 Agent | 通过 Capability 创建意图 |
| `world.intent.approve` | 指定审批岗位 | 批准高风险意图 |

IAOS 服务身份不等于管理员身份；只能访问绑定 tenant 和 world run。日志不得包含 JWT、自由文本隐私数据或未授权的完整 World State。

每个 journal entry 可追溯 actor/service identity、permission decision、input hash、correlation、IAOS operation ref、Outbox ref 和 commit time。

## 10. 当前 IAOS 能力与缺口

| 能力 | 当前可复用 | M8 缺口 |
| --- | --- | --- |
| 外生事实入口 | `/api/v1/simulation/events` 有 allowlist、tenant、幂等、事务 Outbox | 当前入口同时改变业务状态，不适合作为纯 observation journal |
| 场景业务事件 | scenario events 有 strict payload、cursor、SSE | 仅支持 `hctm/order-expedite-01` 和固定事件类型 |
| Capability | 有权限、流程、事务、record refs、committed 标记 | 尚无 World intent 原语和统一 world context |
| 结果恢复 | scenario snapshot/events 支持 cursor | 尚无跨 observation/intent/outcome 的 tenant journal |
| 失败安全 | 已有 RLS、input hash、duplicate/no-op 和审计模式 | 需补 world run binding、branch、payload registry 和权限资源 |

结论：不扩展现有 simulation ingress 为万能入口。IAOS 新增窄 World Bridge journal 和 Capability 内部 intent/outcome 写入原语；现有 scenario/simulation API 保持兼容。

## 11. 跨仓交付顺序

1. AESE 先提交 envelope JSON Schema、fixture、canonical hash 和 mock journal contract tests。
2. IAOS 独立 worktree 实现 journal schema、RLS、payload registry、observation endpoint、cursor query/SSE 和 Capability intent/outcome 原语。
3. IAOS 验证事务 Outbox、权限、tenant、幂等、回滚、乱序和 cursor 恢复后独立提交。
4. AESE 实现 bridge adapter、PostgreSQL cursor checkpoint 和 world reducer 接线。
5. 两仓部署后用 `LAS-WLD-02` tracer 对账 message、journal、Outbox、world event 和 state hash。

## 12. 合同验收标准

- 相同 observation 重试不会产生第二条 journal 或任务副作用。
- Capability 回滚不会出现 committed outcome。
- IAOS 提交成功但 AESE 响应丢失时，AESE 能从旧 cursor 恢复并只计算一次后果。
- SSE/NATS 丢失、重复或乱序不改变最终 World State。
- tenant-other 无法提交或读取 HCTM journal。
- Actor 在 observation 可用时间前不能获得该知识，也不能绕过 IAOS 权限提交 intent。
- 每个 world consequence 能追溯唯一 committed outcome、intent、observation 和原始 world event。
