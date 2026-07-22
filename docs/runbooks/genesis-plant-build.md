# M10 Genesis Plant Build Runbook

## 页面验收

确认 IAOS `:8082`、AESE `:8090` 和前端 `:4173` 健康，访问 `/#world-incorporation`，点击“工厂建设 Campaign”。初态必须显示消费 M9 机器资格；单步 9 次后显示 `M11 eligible`、七个空间节点、World/IAOS 进度、现金/承诺/应付/实付及已关闭 discrepancy。

## 离线验证

```bash
go test ./...
go vet ./...
go run ./cmd/aese world validate world-packs/hctm-genesis
cd frontend
npm test -- --run
npm run typecheck
npm run build
npx playwright test e2e/plant-build.spec.ts
```

`internal/plantbuild` 验证三个虚构候选、先硬约束后评分、10 帧状态机、100 次 hash、snapshot/restore/reset、资金和验收门。绿地及代建方案因预算/日期失败；租赁标准厂房方案在 1,500 万预算内获选。

## IAOS 治理与恢复

`POST /api/v1/genesis/plant/actions` 使用严格合同与 `expected_version`，支持 site evaluate/approve、project execute/rebaseline/accept、payment approve。相同幂等输入为 no-op；并发版本冲突为 409；越权、自批、超预算、过期 mandate 或未验收付款失败关闭。治理记录、committed outcome journal 和 Outbox 同事务。

World 页面刷新后从确定性 campaign API 恢复；正式在线恢复以 IAOS journal cursor 为事实，SSE 只作通知。AESE 不直写 IAOS 数据库，也不把 IAOS 计划进度当作施工事实。
