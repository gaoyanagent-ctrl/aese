# HCTM M4 Simulation Ingress Evidence

日期：2026-07-19。范围：DES-048 首个 `eam.machine.down` tracer。

## 实现边界

- IAOS 仅允许 `eam.machine.down`，服务端生成 subject 和事件 ID。
- 设备通过 tenant entity metadata 的物理表和稳定 code 解析；状态只允许 `running|idle -> maintenance`。
- advisory lock、状态 CAS、ingress audit/idempotency 和事务 Outbox 位于同一 tenant transaction。
- AESE replay 从 canonical event metadata 构造 business object、correlation、causation 和 idempotency，不直发 NATS；IAOS 的 2xx 回显若缺少稳定事件/对象标识或未实际提交，AESE 会失败关闭。

## 真实验收

| 验收 | 结果 |
| --- | --- |
| 首次提交 | HTTP 200，`committed=true`，event `evt-sim-43b222fe83807fc9e68c23b9` |
| 状态迁移 | `LAS-WLD-02` 从 `running` 到 `maintenance` |
| 完全重复 | HTTP 200，同一 event ID，`duplicate=true` |
| 同键不同意图 | HTTP 409，`idempotency_key_collision` |
| 跨租户对象 | `tenant-other` HTTP 404，`business_object_not_found` |
| 持久化 | `tenant-hctm` 仅 1 条 ingress、1 条 `eam.machine.down` Outbox（`PROCESSED`） |
| 跨租户持久化 | `tenant-other` ingress 0 条 |
| AESE canonical replay | 22 事件重放成功；设备事件显示 `action=duplicate` 并返回相同设备 UUID |

subject 为 `iaos.tenant-hctm.eam.machine.down`，设备 UUID 为 `876321d8-7ead-4d21-a8f1-dd5de328016d`。首次真实执行还暴露并修复了动态实体物理表解析和 PostgreSQL text 不接受 NUL advisory-lock key 两个问题。

## 验证命令

- IAOS Platform：`go test ./...`、`go vet ./...`、部署脚本健康检查。
- AESE：`go test ./...`、`go vet ./...`，simulation HTTP wire/malformed success 回归测试，以及带 `--order-id` 的真实 `aese replay --apply`。
- 数据库只读取证：查询 `simulation_event_ingress`、`sys_outbox` 和跨租户计数。

## 剩余范围

供应商延期和来料检验失败尚未进入 allowlist；M4 因此仍为 Active。O2D 也尚未消费设备停机事件，本切片证明的是受治理事实入口及可消费 Outbox，而不是完整重排产流程。
