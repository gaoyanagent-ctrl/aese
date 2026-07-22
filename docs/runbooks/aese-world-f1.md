# AESE World F1 确定性内核运行手册

F1 输入目录包含 `world-run.json`、`initial-state.json` 和 `scheduled-events.json`。仓库最小样例位于 `world-contracts/runtime-example/`。F1 只提供通用确定性内核合同，不是 Genesis pack，也不连接 IAOS。

```bash
go run ./cmd/aese world validate world-contracts/runtime-example
go run ./cmd/aese world inspect world-contracts/runtime-example
go run ./cmd/aese world run world-contracts/runtime-example
```

以上命令均不写文件。`run` 可用 `--until 2026-07-08T10:15:00+08:00` 推进至指定虚拟时间。只有显式 apply 才产生事件日志和快照：

```bash
go run ./cmd/aese world run world-contracts/runtime-example --apply --output /tmp/aese-world-run
go run ./cmd/aese world replay world-contracts/runtime-example --log /tmp/aese-world-run/event-log.json
```

排序固定为 `sim_occurred_at`、数值较小的 `priority`、`event_id`。相同 pack、rules version、seed 和输入必须产生完全一致的事件日志、状态 hash 和 snapshot。未知 rules/payload、重复 event/idempotency key、因果引用倒序、虚拟时间倒退、日志 hash 不一致或损坏快照均失败关闭。

`--apply` 当前只允许写调用者显式指定的输出目录，使用临时文件后原子 rename；它不写 PostgreSQL、不调用 IAOS、不发布 NATS。World Store 持久化接线和 IAOS bridge 属于后续切片。
