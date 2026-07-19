---
id: PLAN-M4-001
title: M4 受治理异常事件入口实施计划
date: 2026-07-19
status: active
author: Codex + User
tags: [m4, simulation, iaos, hctm]
---

# M4 受治理异常事件入口实施计划

## 1. 目标

让 HCTM 的设备停机、供应商延期和来料检验失败通过 IAOS 的权限、租户、幂等、审计和事务 Outbox 边界进入运行链；AESE 只负责提供场景事实，不接受或拼接消息 subject。

## 2. 工作项

- [x] 在 IAOS 增加 `POST /api/v1/simulation/events` 和事件 allowlist。
- [x] 深度实现 `eam.machine.down`：稳定设备解析、状态 CAS、审计/幂等与 Outbox 单事务提交。
- [x] 验证首次提交、完全重复、幂等碰撞、未知对象和跨租户失败关闭。
- [x] 在 AESE IAOS client 和 replay 中接入设备停机事件。
- [x] 用 canonical HCTM 事件执行真实 replay 幂等验证。
- [ ] 实现 `o2d.supplier_delivery.delayed` 的稳定采购对象解析和状态影响。
- [ ] 实现 `qms.incoming_inspection.failed` 的稳定检验/批次解析和质量状态影响。
- [ ] 对三类异常完成统一权限、RLS、审计、Outbox、重复和碰撞验收。
- [ ] 固定供 M5 Agent 与 M6 在线数据源消费的查询/事件合同。

## 3. 当前证据

设备停机 tracer 已完成，详见 `docs/reports/hctm-m4-simulation-ingress-evidence.md`。IAOS 提交为 `9a8f5ca`、`463abd6`、`153a97a`；M4 保持 active，直到另外两类异常和统一验收完成。
