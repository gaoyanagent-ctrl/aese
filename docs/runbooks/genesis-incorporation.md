# M9 Genesis Incorporation Runbook

## 启动与页面验收

确保 IAOS 8082、AESE 8090 和前端 4173 健康，然后访问 `/#world-incorporation`。从 `pre_incorporation` 单步推进 7 次，终态必须显示 `M10 eligible`；投资人现金、公司现金、认缴、实缴和预算授权必须分别展示。

## 离线验证

```bash
go test ./...
go vet ./...
go run ./cmd/aese world validate world-packs/hctm-genesis
cd frontend && npm test -- --run && npm run typecheck && npm run build
npx playwright test e2e/incorporation.spec.ts
```

重复构建 campaign 100 次 hash 必须一致；损坏 snapshot、资金不守恒、未接受任命、失效 mandate、自批预算必须失败关闭。

## IAOS 治理动作

`POST /api/v1/genesis/governance/actions` 只接受成立批准、任命、资本核验、预算提交/批准。浏览器或 AESE 不直写 IAOS DB；外部登记/银行结果只能通过 World observation 返回。相同 idempotency key 与输入返回 duplicate no-op，异输入冲突。

## 恢复与对账

World 页面刷新后从确定性 campaign API 恢复；正式 bridge 消费仍以 journal cursor 为事实，SSE 仅作通知。对账 `genesis_governance_record`、`sys_outbox`、World causation、Knowledge、现金 owner 和 budget stable ref。
