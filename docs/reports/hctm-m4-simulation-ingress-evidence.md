# HCTM M4 Simulation Ingress Evidence

日期：2026-07-19。范围：DES-048 的 `eam.machine.down`、`o2d.supplier_delivery.delayed` 和 `qms.incoming_inspection.failed` 三类受治理事实入口。

## 实现边界

- IAOS allowlist 仅包含上述三类 M4 事件，服务端生成 subject 和事件 ID。
- 设备、采购单和来料检验单分别通过 tenant metadata/fixture 和稳定 code 解析；状态影响为设备进入 `maintenance`、采购单进入 `delayed` 并更新 ETA、检验单进入 `failed` 并记录不良数量、缺陷、批次、收货号和严重度。
- advisory lock、状态 CAS、ingress audit/idempotency 和事务 Outbox 位于同一 tenant transaction。
- AESE replay 从 canonical event metadata 构造 business object、correlation、causation 和 idempotency，不直发 NATS；IAOS 的 2xx 回显若缺少稳定事件/对象标识、未实际提交，或 subject 不精确匹配 `iaos.<target-tenant>.<event-type>`，AESE 会失败关闭。

## 真实验收

| 验收 | 结果 |
| --- | --- |
| 场景准备 | reset 删除 6 个旧 L2 对象；21-object dry-run 为 9 insert/12 no-op；apply 为 9 insert/12 no-op；第二次 apply 为 21 no-op |
| 设备停机 | `LAS-WLD-02` 从 `running` 到 `maintenance`；event `evt-sim-43b222fe83807fc9e68c23b9` |
| 供应延期 | `PO-202607-0001` 为 `delayed`，`latest_eta=2026-07-11`；event `evt-sim-ef2642b944ad516c0319e0bf` |
| 来检失败 | `IQC-202607-0002` 为 `failed`，`rejected_qty=300`、`defect_code=SURFACE_SCRATCH`、`severity=major`，并保留 lot/receipt |
| 完全重复 | 第二次 canonical 22-event replay 为 0 triggered/22 skipped/0 failed；三类异常均为 `duplicate`，订单为 `already_confirmed` |
| 同键不同意图 | real-PostgreSQL integration 返回 HTTP 409 `idempotency_key_collision`，无业务状态或 Outbox 增量 |
| 未知/跨租户对象 | real-PostgreSQL integration 返回 404 `business_object_not_found`，目标外租户 ingress 为 0 |
| 持久化 | `tenant-hctm` 三类 ingress 各 1 条、三类事务 Outbox 各 1 条 |
| O2D tracer | 订单确认 workflow completed，并为新订单生成 3 张 pending work order；重复 replay 不再确认订单 |

三类 subject 分别为 `iaos.tenant-hctm.eam.machine.down`、`iaos.tenant-hctm.o2d.supplier_delivery.delayed` 和 `iaos.tenant-hctm.qms.incoming_inspection.failed`。真实执行还暴露并修复了动态实体物理表解析、PostgreSQL text advisory-lock key、历史设备输入 hash 兼容、精确 decimal 和检验严重度持久化问题。

上述首次/重复 canonical replay、状态和 Outbox 计数来自 2026-07-19 本地 live Platform/PostgreSQL/NATS/O2D。碰撞、未知对象、跨租户、事务回滚和数量边界来自同一最终代码的 real-PostgreSQL integration suite；单元测试另覆盖严格 payload、subject 和 malformed-success 失败关闭。

## 验证命令

- IAOS Platform/shared/O2D/AI：各模块 `go test ./...`、`go vet ./...`；Platform real-PostgreSQL integration 使用 `-tags=integration -race`；部署脚本健康检查通过。
- AESE：`go test ./...`、`go vet ./...`、pack validate/inspect、simulation HTTP wire/malformed success 回归测试，以及带 `--order-id` 的两次真实 `aese replay --apply`。
- 数据库只读取证：查询 `simulation_event_ingress`、`sys_outbox` 和跨租户计数。

## 后续边界

M4 证明的是受治理事实入口及可查询业务上下文，不代表 EAM/O2D/QMS 下游异常消费者、自动重排产、Agent Runtime 或在线 `IaosScenarioDataSource` 已实现。legacy 全表 FORCE RLS、tenant-safe composite foreign key 和 metadata `version` 的平台级版本排序仍是后续 hardening，不改变本入口显式 tenant predicate 与事务边界的验收结论。
