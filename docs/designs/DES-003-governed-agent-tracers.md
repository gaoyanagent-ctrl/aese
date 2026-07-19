---
id: DES-003
title: M5 受治理 Agent Tracer 设计
date: 2026-07-19
status: completed
author: Codex + User
tags: [m5, agent, ai-tool, hctm]
---

# M5 受治理 Agent Tracer 设计

## 1. 目标

让计划、质量和经营分析三个 Agent 从 `tenant-hctm` 的真实 IAOS 业务状态读取最小上下文，通过 AI Tool Registry 留下调用审计，并输出带 `correlation_id`、业务对象引用、证据、建议和数据完整性声明的只读结果。

## 2. 边界

- IAOS 提供通用 `entity.records` query dispatcher。entity、字段、过滤字段和最大行数由 operator-owned tool metadata 固定；调用者只能提交过滤值，不能提交表名、字段名、SQL 或 tenant。
- AESE 持有 HCTM metadata/tool manifests、Agent 编排和场景验收合同，不把 HCTM 规则写入 IAOS 平台代码。
- 首版 Agent 只生成建议和草稿动作，不释放采购单、不调整排产、不隔离库存。
- 每次真实 Agent tracer 都通过 `/api/v1/ai/tools/:key/call`，业务表零写入，`ai_tool_call` 和四个 milestone events 追加审计。
- 经营分析在缺少完工入库、发运和实际成本事实时必须输出 `partial`，不得把 Preview 的 11,700/300 或成本估计冒充在线事实。

## 3. 数据流

```text
HCTM agent tool bundle
  -> AESE agent-setup --apply
  -> IAOS metadata schema + AI Tool Registry
  -> AESE agent-run --apply
  -> registered entity.records tools
  -> tenant RLS transaction
  -> deterministic evidence normalizer
  -> planning / quality / business recommendation envelopes
  -> ai_tool_call + milestone audit
```

## 4. 完成标准

- 三个 Agent 均能引用 `corr-so-202607-0001` 和具体订单、采购单、设备、检验单、批次或工单。
- 计划 Agent 回答需求、净生产、物料时点、设备风险和建议方案。
- 质量 Agent 回答供应商批次、缺陷、影响数量、是否可直接投产及隔离/追溯建议。
- 经营分析 Agent 只输出当前事实可支撑的结论，并显式列出 shipment/cost 数据缺口。
- 跨租户、非法字段、非法 filter、disabled/unauthorized tool 失败关闭。
- 重复 tracer 不改变业务状态；每次调用产生独立审计记录并可追溯。
