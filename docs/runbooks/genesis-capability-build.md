# M11 Genesis Capability Build Runbook

访问 `/#world-plant-build`，点击“生产能力 Campaign”，单步 9 次。终态必须显示 `M12 eligible`、七项设备/实验室能力、10 名核心人员、九项联合 gate 和关闭的检漏校准差异。

离线运行 `go test ./... && go vet ./...`、`go run ./cmd/aese world validate world-packs/hctm-genesis`；前端运行 test/typecheck/build/Playwright。IAOS `POST /api/v1/genesis/capability/actions` 的 accept/qualify 必须携带 World evidence，付款必须携带验收，且受 RLS、权限、预算、mandate、自批、幂等与 expected version 约束。

页面刷新从确定性 API 恢复；线上恢复以 journal cursor 为事实，SSE 只通知。资产卡、员工档案、培训记录或 UI 点击均不能创造设备能力、实际到岗或技能。
