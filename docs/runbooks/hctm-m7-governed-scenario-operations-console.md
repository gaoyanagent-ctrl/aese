# HCTM M7 受治理场景运行控制台 Runbook

> 状态：Completed。该文档记录 M7 O4 收口验收步骤、自动采集入口与已验证边界。

## 0. 适用范围

- 租户：`tenant-hctm`
- 场景：`hctm/order-expedite-01`
- 控制边界：AESE 仅做编排与状态查询；浏览器不直接调用 IAOS 写端点。所有实际写入通过 AESE orchestration API 发起，执行身份由 IAOS JWT 透传。
- 所有前后续动作必须共享同一 `Authorization` token，`/dev/token` 仅允许在首次检查时获取，不能在每个 run action 中重复换发。

## 1. 启动与连通性核对（必须先执行）

1. 启动 IAOS（本地/远端环境）：
   - 记录 IAOS 基础地址（`IAOS_BASE`）、UI 地址（`IAOS_UI`）、JWT 机制生效。
   - 记录用于测试的租户与用户角色。
2. 启动 AESE orchestration API：
   - `go run ./cmd/aese-server --listen :8090 --pack-dir scenario-packs/hctm`
   - 健康检查：
     - `curl -i http://127.0.0.1:8090/health`
     - `curl -i http://127.0.0.1:8090/ready`
3. 启动前端：
   - `cd frontend`
   - `npm install`
   - `npm run dev -- --host 0.0.0.0 --port 3000`
4. 联动中心连通页签：
   - base url 目标：IAOS Base（`.../api/v1` 或 IAOS 提供网关路径）
   - orchestrator url：`http://127.0.0.1:8090`
   - tenant：`tenant-hctm`
   - pack key：`hctm`
   - scenario key：`order-expedite-01`
5. 连接检查（Runbook 通过项）：
   - 显示 IAOS 身份/租户信息
   - 显示当前 snapshot（对象数量、KPI 指标）
   - 显示场景事件与建议列表可读

## 2. 标准业务闭环（O3/O4）

### 2.1 联动检查与预检

1. 点击 `连接检查` 完成：
   - 应看到 `scenario`, `snapshot`, `events`, `recommendations`, `permissions`。
2. 点击 `Pre-flight`：
   - 输出 `dry-run` 影响、`plan_hash`、`act_count`、`next_cursor`。
   - 确认未出现无关对象/跨租户对象清单。

### 2.2 初始化与七幕推进

1. 点击 `Initialize`：
   - 记录 `run_id`、`run_version`、`plan_hash`、`cursor`。
2. 点击 `Advance`：按顺序执行第 1~7 幕（或 `Run to End` 一次执行到结束）：
   - 每一步必须记录返回 `status`、`current_act`、`cursor`、`last_action`。
   - `status` 应从 `initialized -> running -> completed/finished` 演化。
3. 中断恢复：
   - 观察刷新后，页面应自动读取 `localStorage` 中的 active run context 并自动恢复至同一 run 状态。

### 2.3 Analyze + Verify + Reset

1. `Analyze`：
   - 触发三类分析，记录 `analysis_id` 与返回报文的版本信息。
2. `Verify`：
   - 对照 `11,700` 实发、`300` 缺口（版本约定值）完成 KPI 验证。
3. `Reset`：
   - 先执行 `Reset` 影响预览；
   - 仅在确认 token 出现后执行复位；
   - 复位成功后清空本地 run 上下文并返回可初始化状态。

## 3. O4 闭环验收（T35/T36/T39）

### T35 全链路 clean reset 验收

- 场景必须从 clean 状态执行以下顺序：
  - connect -> preflight -> initialize -> (advance/run-to-end) -> analyze -> verify -> reset
- 业务侧核对目标：
  - 22 个场景事件数
  - 3 类 Agent 输出
  - 17 条离线断言（来自 `expected`）
  - KPI 指标与 M6 视口一致性
- 采集要求（每次执行都附路径）：
  - `test-results/governed-console-*.png`（控制台视图）
  - `test-results/live-*.png`（Live 视图）
  - `test-results/governed-reset-*.png`（重置预览+执行）

### T36 CLI/数据库/Outbox/Tool Call 一致性复核

- 同一故事使用同一 `run_id`（建议加时间戳后缀）分别跑：
  - AESE Web 受控执行链路
  - AESE CLI `replay/reset/verify`（使用 `--apply` 与 dry-run 对照）
- 复核标准（均需同一 run context）：
  - Run 级别：`run_id`、`plan_hash`、`run_version`、`cursor` 与阶段状态一致。
  - 场景事件级：事件 `event_id`/`business_object` 与顺序一致；重复提交仅有 no-op。
  - 副作用级：DB、Outbox、Tool Call 的增量应与 UI 与 CLI 一致；
    出现重复提交时必须记录无新副作用（`no_op` 或等价幂等痕迹）。
- 数据库与 Outbox 复核已在本机 `iaos-integration-postgres` 完成；最终结果见 M7 evidence。

### T39 部署与健康闭环

- AESE：
  - 记录部署 URL、服务实例、镜像/构建 ID；
  - 验证 `GET /health`、`GET /ready`。
- IAOS：
  - 同步 IAOS 必要 endpoint 与 schema 变更 commit；
- 两仓验收项：
  - 在 report 中记录 AESE 与 IAOS 两个仓 commit；
  - 记录跨仓 `orchestrator` 与 `IAOS` 连通健康路径：
    - AESE -> IAOS `/api/v1/scenarios/...` 可达；
    - 运行命令可持续返回 run 状态。

## 4. 异常路径（故障注入）

- plan hash 不匹配：返回 `plan_hash_mismatch`，前端应提示 `run_context outdated` 并支持刷新并重试。
- cursor 不匹配：返回 `cursor_mismatch`，前端应提示基于最新状态重试。
- 跨租户操作：返回租户或权限拒绝，UI 阻断执行按钮。
- 幂等键缺失/重复：返回 `idempotency_required` 或 no-op；重复提交不得产生二次写入。
- reset token 过期：返回 token 过期错误并阻断 reset。

## 5. 数据留痕清单（本 runbook 通过即保存）

- 终端日志：
  - AESE API logs
  - 前端控制台 logs
  - Playwright `test-results/*` 截图和 trace
- 结果文件（建议目录）：
  - `artifacts/m7-acceptance/{ts}/`
    - `00-plan.json`
    - `01-run-create.json`
    - `02-preflight.json`
    - `03-initialize.json`
    - `04-run-to-end.json`
    - `05-analyze.json`
    - `06-verify.json`
    - `07-reset-plan.json`
    - `08-reset.json`

## 6. T36 CLI 与 UI 副作用一致性脚本化对账

脚本：`scripts/m7-runbook-evidence-collect.sh`

执行方式（按需调整 `DRY_RUN/UI_APPLY/CLI_APPLY`）：

```bash
cd /iaos/aese
export AESE_BASE=http://127.0.0.1:8090
export IAOS_BASE=http://127.0.0.1:8082
export TENANT_ID=tenant-hctm
export OUTPUT_DIR=artifacts/m7-acceptance/$(date -u +%Y%m%dT%H%M%SZ)
export UI_APPLY=1
export CLI_APPLY=0   # 可改 1 做应用链路
export DRY_RUN=0
# 可选：将本次链路 token 固定透传，避免脚本内部重复调用 /dev/token
# export IAOS_TOKEN=<可复用的 IAOS JWT>
./scripts/m7-runbook-evidence-collect.sh
```

脚本输出文件用于 T36 对账：

- `00-plan.json` 与 `01-run-create.json`（plan_hash / run_id / status / cursor）
- `02-preflight.json`、`03-initialize.json`、`04-run-to-end.json`、`05-analyze.json`、`06-verify.json`、`07-reset-plan.json`、`08-reset.json`（UI/编排 API 全链）
- `05-cli-apply-*.json`、`06-cli-replay-*.json`、`07-cli-verify-attempt-*.json`、`07-cli-verify.json`、`08-cli-reset-*.json`（CLI 执行链与异步收敛证据）
- `summary.txt`

对账建议（最少核对）：

1. `plan_hash`、`run_id`、`run_version` 一致。
2. UI/CLI run-to-end 的事件命令字数、`event_id` 顺序与 `correlation_id` 交织一致。
3. no-op 统计、`applied`/`unchanged` 与 DB/Outbox 可观察计数一致（见 T36 要求）。
4. 重置阶段需包含 `reset_confirmation_token` 预览与执行两步产物。

## 7. 完成标志

- T35/T36/T39 已于 2026-07-22 达到交付条件，权威运行产物位于 `artifacts/m7-acceptance/20260722-05/`。
- 场景业务断言为 17 条离线确定性断言；在线 IAOS verify 为 2 条工作单记录断言，两层结果均须保留。
