# HCTM M6 在线沙盘运行手册

## 前置与在线事实

IAOS Platform 使用 `http://127.0.0.1:8082`，AESE frontend 使用 `http://127.0.0.1:4173`。本地 token：

```bash
export IAOS_TOKEN="$(curl -fsS 'http://127.0.0.1:8082/api/v1/dev/token?tenant_id=tenant-hctm&roles=admin' | jq -r .token)"
```

依次执行 `reset --apply`、`apply --apply`、`replay --apply` 和 `agent-run --apply`。均使用 pack `./scenario-packs/hctm`、story `order-expedite-01`、target `http://127.0.0.1:8082` 和 tenant `tenant-hctm`；replay 的 `--order-id` 取 apply 回显中的 sales order UUID。重复 replay 必须为 `0 triggered / 22 skipped / 0 failed`。

```bash
curl -fsS -H "Authorization: Bearer $IAOS_TOKEN" -H 'X-Tenant-ID: tenant-hctm' \
  http://127.0.0.1:8082/api/v1/scenarios/hctm/order-expedite-01/snapshot | jq
```

补发使用 `/events?after=<cursor>&limit=200`；SSE 使用 `/events/stream?after=<cursor>`。SSE 仅作增量提示，重连必须先补发。

## UI 与验收

```bash
cd frontend
npm install
npm run dev
```

浏览器控制台执行 `localStorage.setItem('iaos_token', '<JWT>')` 后刷新并切换 Live。浏览器请求默认走前端同源 `/api`，Vite 将其代理到 Platform 8082，因此远程浏览器不应直接访问 `127.0.0.1:8082`。生产环境必须由真实登录提供 JWT，并在入口代理配置同样的 `/api` 转发。Live 应显示需求 12,000、累计可供/实发 11,700、期末成品 0、缺口 300 和 `cost_actuals` gap；错误不会静默降级为 Preview。

自动化命令：AESE 根目录执行 `go test ./...`、`go vet ./...`；frontend 执行 `npm run typecheck`、`npm test`、`npm run lint`、`npm run build`、`npm run test:e2e`。
