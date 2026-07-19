---
id: PLAN-M3-001
title: M3 可执行 HCTM 场景包实施计划
date: 2026-07-19
status: active
author: Codex + User
tags: [m3, implementation, hctm, iaos]
---

# M3 可执行 HCTM 场景包实施计划

## 1. 目标

把 HCTM Markdown 规格转换为机器可读场景包，提供离线校验和 inspect CLI，并完成一个通过 IAOS 受治理接口导入数据、触发 `o2d.order.confirmed`、验证 O2D 结果的 tracer bullet。

## 2. 完成定义

M3 只有同时满足以下条件才可标记 completed：

- HCTM pack 数据和 JSON Schema 已提交。
- Go validator/inspect 命令及单元测试通过。
- IAOS compatibility report 已生成并归档。
- dry-run 与 apply 行为清晰分离。
- `tenant-hctm` tracer 成功触发一次 O2D 订单分解。
- 第二次执行没有重复对象或重复副作用。
- reset、replay、verify 有 runbook 和实际验证证据。
- AESE 与 IAOS 两边的 code map、设计、测试和进展日志已同步。

## 3. 任务与依赖

### S1 - 场景包与 Schema

- [ ] T1 创建 `scenario-packs/hctm/manifest.json`。
- [ ] T2 把组织、交易方、物料、制造和物流主数据转换为 JSON record sets。
- [ ] T3 创建 `order-expedite-01` initial state、events 和 expected outcomes。
- [ ] T4 创建 manifest、record set、event sequence、expected outcomes JSON Schema。
- [ ] T5 人工核对 Markdown 规格与 JSON 数量、编码和时间线。

验收：所有 JSON 能被标准解析器读取；22 个事件、固定 correlation ID 和关键数量与 Seed Data Plan 一致。

### S2 - Go Loader、Validator 和 Inspect

- [ ] T6 初始化最小 Go module 和 `cmd/aese`。
- [ ] T7 实现场景包加载、路径安全和 schema version 检查。
- [ ] T8 实现结构、自然键、跨文件引用和事件时间线校验。
- [ ] T9 实现 BOM/MRP/库存/发货经营不变量。
- [ ] T10 实现 `aese validate` 和 `aese inspect`。
- [ ] T11 增加有效 pack 和破损 fixture 的表驱动测试。

验收：离线命令退出码稳定；错误包含文件、记录和字段位置；测试覆盖重复键、缺失引用、乱序事件、重复幂等键和数量不变量。

### S3 - IAOS Compatibility

- [ ] T12 启动并验证 IAOS platform `8082`、PostgreSQL、NATS 和 O2D。
- [ ] T13 读取 IAOS metadata schema/API，生成 HCTM → IAOS compatibility report。
- [ ] T14 固定第一批 compatible/mapped/unsupported 对象清单。
- [ ] T15 对 legacy `sales_order`、product/material、BOM、inventory 字段作出映射决策。
- [ ] T16 判断是否需要 IAOS 新的 scenario import 或 simulation ingress；若需要，单独建立 IAOS DES/ADR。

验收：报告逐对象列出字段映射、API、权限、缺口和处理策略；不得使用直接 SQL 作为正式解决方案。

### S4 - Dry-run 与 Apply Tracer

- [ ] T17 实现 IAOS client 的认证、schema 查询和记录 upsert adapter。
- [ ] T18 实现 `aese apply` 默认 dry-run impact report。
- [ ] T19 实现显式 `--apply` 和 run summary。
- [ ] T20 导入最小 customer/material/BOM/inventory/sales order 数据。
- [ ] T21 验证租户隔离和第二次 apply 幂等。

验收：dry-run 零写入；apply 只写 `tenant-hctm`；重复 apply 无重复自然键；失败报告可定位到具体记录。

### S5 - O2D Replay Tracer

- [ ] T22 通过 IAOS 订单业务入口触发 `o2d.order.confirmed`。
- [ ] T23 验证 NATS subject、workflow 执行和 work order/inventory 结果。
- [ ] T24 保存 correlation、事件 ID、运行日志和查询证据。
- [ ] T25 实现 `aese verify` 的最小断言。

验收：订单确认到 O2D 结果可重复；事件属于 `tenant-hctm`；第二次触发受幂等保护或明确拒绝。

### S6 - Reset 与 Closeout

- [ ] T26 实现场景级 reset 计划和安全确认，不删除 L1 公共主数据。
- [ ] T27 编写 `docs/runbooks/hctm-m3-local-run.md`。
- [ ] T28 更新 README、architecture、code map、roadmap 和 progress log。
- [ ] T29 记录测试命令、结果和剩余风险。
- [ ] T30 将 DES-001 状态从 draft 更新为 approved/completed（按实际情况）。

## 4. 跨仓库工作安排

AESE 提交：

- scenario pack、schema、CLI、validator、IAOS client、runbook。

IAOS 独立 worktree 提交：

- eventdef 扩展、simulation ingress、O2D handler 或 metadata/API 兼容改动。

若 IAOS 无需修改，compatibility report 必须解释为什么现有 API 足够。

## 5. 验证命令目标

命令将在实现后固化，目标形式：

```bash
go test ./...
go run ./cmd/aese validate ./scenario-packs/hctm
go run ./cmd/aese inspect ./scenario-packs/hctm
go run ./cmd/aese apply ./scenario-packs/hctm --target http://127.0.0.1:8082
go run ./cmd/aese replay ./scenario-packs/hctm --story order-expedite-01 --target http://127.0.0.1:8082
go run ./cmd/aese verify ./scenario-packs/hctm --story order-expedite-01 --target http://127.0.0.1:8082
```

前三个阶段不得依赖 IAOS 在线服务；S3 以后需要平台和基础设施。

## 6. 当前状态

- S0 工程治理：completed。
- S1-S6：pending。
- 当前第一项开发任务：T1-T4，建立 HCTM machine-readable pack 和 JSON Schema。
