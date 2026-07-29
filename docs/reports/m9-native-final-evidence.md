---
id: REPORT-M9N-001
title: M9 IAOS 原生企业成立闭环最终证据
date: 2026-07-23
status: completed
author: Codex + User
tags: [aese, m9, evidence]
---

# M9 IAOS 原生企业成立闭环最终证据

## Clean tracer

- 正常：`INC-E2E-1784787213558168936`，
  correlation `corr-INC-E2E-1784787213558168936`，
  `enterprise_operational_ready`，G1–G7 consumed=7，
  Intent/Observation/CommittedOutcome=3/3/3，Process=1，Decision=15。
- 登记补正：`INC-CORRECTION-1784787214040461532`，
  correlation `corr-INC-CORRECTION-1784787214040461532`，
  保持 `registration_submitted`，新 Intent=1、错误 Outcome=0、
  incorporation-agent Journal=1。

## D18 十二门

| 门 | 机器证据 | 结果 |
| --- | --- | --- |
| 1 三层语义 | Runtime Artifact 三层 hash、dependency lock、离线 validator | 通过 |
| 2 正式 founder | profile/menu/actor authorization PostgreSQL tests | 通过 |
| 3 20/5/8/G1–G7 | artifact cardinality、Process/Decision/Approval trace | 通过 |
| 4 正常终态 | clean tracer 与 Evidence Bundle | 通过 |
| 5 四异常 | correction、bank rejection/reapproval、capital mismatch、Agent SOD tests | 通过 |
| 6 World 恢复 | recovery matrix、poller cursor、AESE reconciliation | 通过 |
| 7 同投影/深链 | IAOS/AESE 三视口 Playwright 与五参数链接 | 通过 |
| 8 同 Capability | Artifact `entry_points` 四入口同一 command API | 通过 |
| 9 数据库治理 | FORCE RLS、tenant isolation、bigint minor units、并发/幂等/原子性 tests | 通过 |
| 10 操作闭环 | dry-run/apply/no-op、rollback、replay、AESE reset 后法律事实保留 | 通过 |
| 11 Evidence | versioned bundle、canonical/bundle hash、离线 verifier | 通过 |
| 12 founder UI | IAOS 与 AESE 1440×900、1280×720、390×844，刷新恢复 | 通过 |

## 最终回归

以下是 2026-07-23 原始最终验收快照：

- AESE Go：全通过。
- AESE frontend：16 files、38 tests；production build 通过。
- IAOS Go：全通过。
- IAOS M9 PostgreSQL integration matrix：全通过。
- IAOS frontend：52 files、332 tests（单 worker）；production build 通过。
- IAOS UI Playwright：3/3；AESE UI Playwright：3/3。
- Atlas declaration 与 JSON 格式检查通过。

## 2026-07-29 运行边界与财务显式化复验

- IAOS `main@a1cc21e` 已部署到本地联调环境，目标租户
  `tenant-gx-f4b3ce3ce8e2712d` 已从 Runtime 1.6.0 升级到 1.7.0。
- Runtime 1.7.0 的 `definitions_hash` 同时签名 25 项 Capability 合同和
  8 个 Process 定义；仅修改流程图、输入输出或业务目的也会产生新的 Artifact，
  不会再被安装器错误判断为 `no_op`。
- 25 项 M9 Capability 均存在 active `incorporation_command` Implementation
  Binding；正常交互路径为 23 个持久工作项，其中包括财务组织、账套与期间、
  科目、资本过账和财务就绪 5 个显式节点。
- PostgreSQL 定向全生命周期测试到达
  `enterprise_operational_ready`，实测 23 个工作项、7 个审批门、3 个 World
  wait、6 次 Agent Run 和 5 个不同 Agent。
- IAOS `go test ./platform/... ./shared/...`、`go vet`、受治理业务写入检查、
  Code Map 与 System Atlas 检查通过。
- 目标租户 Outbox 43 条均为 `PROCESSED`；业务 payload 会按 Outbox 行中的
  tenant/event routing 包装为完整 Event，World payload 自带 `tenant_id`
  不再被误判为完整 envelope。

当前复验没有重新宣称整个历史 IAOS frontend/Playwright 与共享数据库全量
integration suite 均已复跑。共享库中其他测试租户的历史 FAILED Outbox 和
全量 integration 隔离债务仍按风险登记保留，不能作为目标租户 M9 完成证据。

已知限制记录在
[M9 closed-loop runbook](../runbooks/m9-native-closed-loop.md)。
