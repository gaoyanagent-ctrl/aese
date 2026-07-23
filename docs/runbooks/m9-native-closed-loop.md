# M9 IAOS 原生企业成立闭环运行手册

## 固定环境

- IAOS revision：`1832810`（branch `feat/m9-native-incorporation`）
- AESE revision：以包含本手册的提交为准
- tenant：`tenant-hctm-genesis`
- business timezone：`Asia/Shanghai`
- IAOS API/UI：`http://127.0.0.1:8082`、`http://127.0.0.1:3000`
- AESE UI：`http://127.0.0.1:4173`

正式验收使用普通登录主体 `founder-principal`。密码由环境 seed/bootstrap
提供，不写入 evidence bundle 或提交日志。

## 发布、回退与恢复

1. 调用 `POST /api/v1/incorporations/runtime/install`，空 body 验证 dry-run。
2. 以 `{"apply":true}` 发布；再次 apply 必须返回 `no_op=true,writes=0`。
3. 回退先 dry-run，再显式调用 runtime rollback；不兼容版本使业务命令按
   stale artifact 失败关闭。
4. 已成立法律事实、Approval、Journal、Outbox 和 World exchange 不参与
   scenario reset；reset 只复位 AESE 的可重放 World frame。
5. 使用原 correlation/idempotency replay；相同键不同 payload 必须冲突。
6. 重启 IAOS、AESE 或浏览器后，从 trace API/World Store 恢复，不从
   local state 推断完成。

## 验收命令

```bash
cd /iaos/iaos-go-m9-native/platform
go test ./...
IAOS_INTEGRATION_DATABASE_URL='postgres://iaos_app:iaos_app@127.0.0.1:5433/iaos?sslmode=disable' \
  go test -tags=integration ./internal/api \
  -run 'TestIntegration(CompleteIncorporationLifecycle|RegistrationCorrection|CapitalMismatch|FormalActor|WorldBridgeRecovery|FiveAgent|FounderOverride)'

cd ../frontend
npm test -- --run --maxWorkers=1 --no-file-parallelism
npm run build

cd /iaos/aese
go test ./...
cd frontend
npm test -- --run
npm run build
npx playwright test e2e/m9-native-lifecycle.spec.ts
```

Evidence bundle 从
`GET /api/v1/incorporations/:case_code/evidence` 获取，并用
`platform/cmd/incorporation-evidence-verify` 离线验证。Bridge 离线对账使用
`aese reconcile <bridge-journal.json>`。

## 风险边界

- 当前只支持单法人、CNY、单设立案主线和五个 Agent。
- 本地部署使用 development fallback secret，不得直接作为生产配置。
- IAOS 前端测试需单 worker 执行；并行 happy-dom/Ant Design 测试会争用
  全局 DOM 调度器，并非支持的确定性验收模式。
