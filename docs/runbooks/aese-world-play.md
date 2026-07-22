# M8 World Play 验收运行手册

## 离线合同与确定性

```bash
go test ./...
go vet ./...
go run ./cmd/aese world validate world-packs/hctm-genesis
go run ./cmd/aese world run world-packs/hctm-genesis --until 2026-07-22T10:00:00+08:00
```

预期：pack 有效；dry-run 不写 artifact；重复执行得到相同 state hash。

## 前端

```bash
cd frontend
npm test -- --run
npm run typecheck
npm run build
npx playwright test e2e/world.spec.ts
```

打开 `/#world`。验证虚拟时钟、运行/暂停/单步/复位、World/IAOS/Knowledge 三栏与偏差因果时间线；桌面 1440/1280 和移动 390 均可操作。

## IAOS bridge

IAOS revision `e661d9a` 在独立 `feat/m8-world-bridge` worktree。部署后检查：

```bash
curl -fsS http://127.0.0.1:8082/health
docker exec iaos-integration-postgres psql -U iaos_app -d iaos -Atc \
  "SELECT relrowsecurity,relforcerowsecurity FROM pg_class WHERE relname='world_bridge_journal'"
```

写端点为 observations/intents/outcomes；读端点为 cursor entries 与 SSE stream。所有写入要求 IAOS tenant/authz，journal 与 Outbox 同事务；重复 idempotency key 同输入返回原 cursor，异输入返回冲突。人类接管和 Agent 不使用旁路，均调用相同 intent/outcome 合同。

## 恢复与边界

SSE 中断后使用最后 durable cursor 调用 entries；只消费 committed/no-op outcome。不得依据 intent、HTTP 超时或通知推断世界后果。M7 pack、CLI、Preview/Live 均保持原路径。
