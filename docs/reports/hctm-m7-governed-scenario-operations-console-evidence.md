# HCTM M7 受治理场景运行控制台 Evidence

## 1. 文档与验收状态

- `docs/designs/DES-005-governed-scenario-operations-console.md`：Approved
- `docs/plans/2026-07-20-m7-governed-scenario-operations-console.md`：Completed
- `docs/runbooks/hctm-m7-governed-scenario-operations-console.md`：Completed

## 2. O3 已完成项（已闭环）

- T28：联动中心 run 与 Live 横幅联动（AESE 前端/状态链路）
- T30：键盘焦点、按钮 `type/ARIA`、日志文本选择和移动端布局补齐
- T31/T32：API/application 与前端状态测试（unit/e2e）补齐
- T33：双击、并发提交、刷新恢复与服务重启回归覆盖
- T34：只读限制、跨租户拒绝、reset token 过期处理覆盖
- T37：三视口截图（1440×900、1280×720、390×844）生成并与 Live 视口复核
- T38：runbook 与 evidence 的主干已建立；后续以章节追加证据

## 3. O4 闭环结果（2026-07-22）

- T35：`m7-acceptance-20260722-05` 从 clean reset 完成 preflight、initialize、run-to-end、analyze、verify、reset。22 个 pack 事件映射为 10 个受治理动作与 12 个 unsupported/no-op 事实，失败 0；三 Agent 产生 9 条 Tool 证据并发布 3 条 recommendation；17 条离线业务断言通过，2 条 IAOS 在线断言通过。M6 KPI 为需求 12,000、可供/实发 11,700、缺口 300、交付状态 `partially_shipped`。
- T36：UI 与 CLI 均为 initialize/apply 9 insert + 12 no-op、replay 10 triggered + 12 skipped + 0 failed、verify 2 passed + 0 failed、reset 15 deleted + 12 L1 preserved。单 run 数据库窗口新增 9 次成功 Tool Call、18 条 Outbox（两套各 9 条的 UI/CLI O2D 副作用）和 2 条 apply run；reset 后 scenario event/recommendation 归零。CLI 首次 verify 在 O2D 异步窗口内失败，第二次有界重试通过，两个 attempt 均保留。
- T39：AESE API `http://127.0.0.1:8090` 的 `/health` 为 `UP`、`/ready` 为 `OK`；AESE frontend `http://127.0.0.1:4173` 返回 200。IAOS API `http://127.0.0.1:8082` 的 `/health` 为 `UP`、`/ready` 的 database/event_bus 均为 `OK`；IAOS frontend `http://127.0.0.1:3000` 返回 200。M7 完成提交为 AESE `a4c9af356028822131fe2060ae68158512818f36`，联调部署使用 IAOS `8e267f74a8d5e9e92d779f0247c8065faf2ccb13`。

## 4. 证据清单（按章节累积）

### 4.1 API 与状态链路

- 示例文件：
  - `artifacts/m7-acceptance/<ts>/00-plan.json`
  - `artifacts/m7-acceptance/<ts>/01-run-create.json`
  - `artifacts/m7-acceptance/<ts>/02-preflight.json`
  - `artifacts/m7-acceptance/<ts>/03-initialize.json`
  - `artifacts/m7-acceptance/<ts>/04-run-to-end.json`
  - `artifacts/m7-acceptance/<ts>/05-analyze.json`
  - `artifacts/m7-acceptance/<ts>/06-verify.json`
  - `artifacts/m7-acceptance/<ts>/07-reset-plan.json`
- `artifacts/m7-acceptance/<ts>/08-reset.json`
- 待补：
  - `plan_hash_mismatch`、`cursor_mismatch`、`run_not_found`、`not_found`、`run_version_mismatch` 的完整报文样本。
- `scripts/m7-runbook-evidence-collect.sh` 运行产物（建议目录）
- `expected_cursor` 在采集脚本中按数字类型发送，避免字符串类型导致 `expected_cursor` 反序列化失败；`IAOS_TOKEN` 在单条 run 链路中保持不变。
- 运行摘要：
  - `run_id`、`tenant`、`plan_hash`、`run_version`、`run_cursor`

### 4.2 运行视图

- 控制台截图：
  - `test-results/governed-console-*.png`
  - `test-results/governed-reset-*.png`
- Live 视图截图：
  - `test-results/live-*.png`
- Playwright trace：
  - 如 `test-results/run-to-end/`

### 4.3 部署与健康

- `AESE`
  - `GET /health`
  - `GET /ready`
  - `GET /api/v1/scenarios/{pack}/{story}`
- `IAOS`
  - 与 AESE orchestrator 协同的联通结果、健康返回和 auth token 测试

## 5. 已知边界与后续

- IAOS O2D 由 Outbox/NATS 异步驱动，CLI verify 必须使用有界重试，不能把首次暂不可见误判为最终失败。
- reset 会清除 scenario event 与 recommendation 业务投影，但 Tool Call 与 Outbox 审计保留；副作用对账应在 reset 前采业务结果、reset 后采审计增量。
- `scripts/m7-runbook-evidence-collect.sh` 是后续回归主入口；默认仍建议先 `DRY_RUN=1`，维护窗口再开启 `UI_APPLY=1/CLI_APPLY=1`。
