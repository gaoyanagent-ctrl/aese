---
id: DES-001
title: M3 可执行场景包与重放架构
date: 2026-07-19
status: completed
author: Codex + User
tags: [m3, scenario-pack, seed, validation, replay, iaos]
---

# M3 可执行场景包与重放架构

## 1. 背景与目标

M2 已经在 Markdown 中定义 28 个对象、18 类事件和 22 步故事时间线。M3 要把这些规格转换为机器可读、可校验、可导入和可重放的场景包，并完成一个连接 IAOS O2D 的 tracer bullet。

## 2. 设计目标

- 确定性：同一 pack 版本产生一致的解析和 dry-run 结果。
- 可验证：不连接 IAOS 也能发现大部分数据错误。
- 安全：默认 dry-run；正式写入经过 IAOS 权限、RLS 和 Outbox。
- 窄切片：先证明一个订单和一个事件链，不一次落地全部模型。
- 可演进：schema version 和 pack version 独立管理。

## 3. 场景包合同

`manifest.json` 最小字段：

```json
{
  "schema_version": "1.0.0",
  "pack_key": "hctm",
  "pack_version": "0.1.0",
  "display_name": "华辰热管理系统集团 MVP",
  "timezone": "Asia/Shanghai",
  "tenant_template": "tenant-hctm",
  "master_data": [],
  "stories": []
}
```

记录集合：

```json
{
  "schema_version": "1.0.0",
  "entity": "material",
  "natural_key": ["material_code"],
  "records": []
}
```

故事结构：

- `initial-state.json`：L2 业务初始状态。
- `events.json`：有序事件 envelope。
- `expected-outcomes.json`：可机器断言的结果。

## 4. 校验层次

V1 结构校验：

- JSON 可解析。
- 必填字段、类型、枚举和格式正确。
- schema version 支持。

V2 业务引用校验：

- 部门引用存在的工厂。
- BOM parent/component 引用存在物料。
- 工序引用存在路线和工作中心。
- 设备引用存在工作中心。
- 订单引用存在客户、物料和法人。

V3 故事校验：

- 事件 ID 和 idempotency key 唯一。
- correlation ID 一致。
- timestamp 单调或有并发组。
- causation ID 指向先前事件。
- 事件主对象存在于 initial state 或先前事件结果。

V4 经营不变量：

- 数量非负，单位一致。
- BOM 用量和 MRP expected 值可重算。
- 发货总量不能在没有显式补产事件时超过可发库存。
- 期望结果与故事结束状态一致。

## 5. CLI

M3 CLI：

```text
aese validate <pack-dir>
aese inspect <pack-dir>
aese apply <pack-dir> --target <url> [--apply]
aese replay <pack-dir> --story <key> --target <url> [--apply]
aese verify <pack-dir> --story <key> --target <url>
```

约束：

- `validate` 和 `inspect` 完全离线。
- `apply` 和 `replay` 无 `--apply` 时只输出计划。
- 凭据从环境或安全配置读取，不写入 pack。
- 每次写操作产生 run summary，不输出完整 token。

## 6. IAOS 集成策略

第一阶段分两步：

1. Compatibility：读取 IAOS metadata schema，比较 entity、字段、类型、required 和自然键，输出 unsupported/mapping-required/compatible。
2. Tracer：只选择最小对象集合导入，并通过现有订单分解入口触发 `o2d.order.confirmed`。

首个 tracer 最小对象：

- customer
- material/product 映射
- BOM
- inventory
- sales_order

IAOS 缺口必须写入 compatibility report，禁止通过 AESE 直接 SQL 绕过。

## 7. 事件注入策略

- 内生事件：订单确认、库存预留、工单创建，由 IAOS 业务动作和 Outbox 产生。
- 外生事件：供应商延期、设备故障、来料不良，由未来 simulation ingress 或受治理 Capability 产生。
- direct NATS 仅作为本地调试 adapter，必须显式不安全参数，M3 验收不依赖它。

## 8. 测试策略

单元测试：

- manifest 和 record-set 解析。
- schema 校验错误定位。
- 跨文件引用和重复自然键。
- 22 步时间线、correlation 和 causation。
- MRP 数量不变量。

集成测试：

- dry-run 不产生写入。
- apply 重复执行幂等。
- tenant-hctm 只能看到自身数据。
- 订单确认触发 O2D handler。
- 失败时报告具体对象和步骤。

## 9. 非目标

- 通用离散事件仿真调度器。
- 任意脚本执行。
- 生产环境凭据管理。
- 全部 28 个对象和 18 类事件的运行实现。
- 2D/3D 可视化。

## 10. 开放问题

- IAOS metadata schema 是否足以表达所有 HCTM 自然键和引用关系。
- 配置包 `seed_data` 尚未成为受支持对象时，M3 apply 应使用逐实体 API 还是新增批量场景导入端点。
- simulation ingress 应作为通用平台能力还是 O2D 场景能力。
- O2D 当前 legacy sales order 字段与 HCTM 规格的映射方式。

## 11. 实现状态（2026-07-19）

- 已实现 `scenario-packs/hctm`、4 个 JSON Schema、Go loader/validator/inspect、IAOS client、dry-run/apply/replay/verify 基础协调器和安全 reset 计划。
- 实现核对发现原故事“可发 11,700、发运 12,000”的矛盾；合同已修订为第二批请求 3,000、实发 2,700、短缺 300，保持 22 个事件且不允许超发。
- IAOS compatibility 证明现有通用 CRUD 不足以完成正式 scenario apply：需要稳定编码到 UUID 映射、原子 apply/reset、订单确认和 workflow 幂等，以及 decimal BOM 修正。
- IAOS DES-047、DES-048 与 O2D 原子幂等修正已实现所需 M3 切片；`tenant-hctm` 的 dry-run/apply/replay/verify/reset 实证完成。本设计状态更新为 completed。
