---
id: PLAN-M4-001
title: M4 受治理异常事件入口实施计划
date: 2026-07-19
status: completed
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
- [x] 泛化 AESE replay，将供应商延期和来料检验失败按 metadata 业务对象送入同一 simulation ingress，并保持默认 dry-run。
- [x] 将两张 canonical 采购单和 `IQC-202607-0002` pending 检验单投影进 DES-047 scenario apply，为异常对象稳定解析提供前置 fixture。
- [x] 加固 AESE fail-closed 合同：精确 tenant subject、projection 必填字段、required reference 和真实 22 事件路由回归。
- [x] 实现 `o2d.supplier_delivery.delayed` 的稳定采购对象解析和状态影响。
- [x] 实现 `qms.incoming_inspection.failed` 的稳定检验/批次解析和质量状态影响。
- [x] 对三类异常完成统一权限、RLS、审计、Outbox、重复和碰撞验收。
- [x] 固定供 M5 Agent 与 M6 在线数据源消费的查询/事件合同。

## 3. 当前证据

三类异常已通过同一受治理入口完成对象解析、状态影响、幂等、租户隔离、审计和事务 Outbox 验收，详见 `docs/reports/hctm-m4-simulation-ingress-evidence.md`。IAOS 最终增量提交为 `8f683b3`、`0097a68`、`9d050f1`、`a17dc81`、`1cdff23`、`2344af0` 和 `42f51dd`；此前设备 tracer 提交为 `9a8f5ca`、`463abd6`、`153a97a`。

供后续里程碑消费的稳定边界是：`shared/eventdef` 的三项事件常量、`POST /api/v1/simulation/events` 的 committed/duplicate 响应、`iaos.<tenant>.<event-type>` Outbox subject，以及通过 metadata/entity API 查询 `equipment`、`purchase_order` 和 `inspection_order` 的租户业务状态。M4 不包含领域消费者、Agent Runtime 或 `IaosScenarioDataSource` 实现。
