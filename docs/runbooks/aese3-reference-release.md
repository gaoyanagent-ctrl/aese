# AESE 3.0 M17-M24 Reference Release Runbook

## 验证

1. `go test ./...` 与 `go vet ./...`。
2. `jq empty world-contracts/schemas/aese3-program.schema.json world-contracts/fixtures/aese3-program.json world-packs/hctm-genesis/aese3/payload-registry.json`。
3. 启动 `go run ./cmd/aese-server`，访问 `GET /api/aese/v1/world/aese3`；确认八个 milestone 顺序、terminal 均为 true、`automatic_business_writes=0`。
4. `cd frontend && npm test && npm run typecheck && npm run build`，访问 `#world-aese3`；移动端和桌面端均应显示 M17-M24 与最终 ready=true。
5. 在 IAOS 独立 worktree 运行 `go test ./platform/internal/api`；向 `/api/v1/genesis/aese3/actions` 提交 `platform.promote`，必须包含 tenant、exact evidence hash、不同 actor/approver、mandate、CAS 和幂等键。
6. 重复同一请求应返回 duplicate；缺 evidence、自批、陈旧 version 和跨租户必须失败关闭。

## 发布与恢复

`hctm-genesis@1.0.0` 是虚构 HCTM reference pack。promotion 只登记受治理发布，不执行订单、采购、工单、发运、资金或 Policy 变更。失败时保留旧 pack/version 和 journal，修复合同后以新 idempotency key、正确 expected version 重试；禁止直接写 IAOS 数据库或 NATS。
