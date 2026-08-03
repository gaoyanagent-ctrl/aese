---
id: SOL-014
title: M10 设施项目 Agent JSON 截断与受控修订
date: 2026-08-03
status: completed
author: Codex + User
tags: [aese, m10, agent, minimax, json]
---

# M10 设施项目 Agent JSON 截断与受控修订

## 现象

在“设施项目与 WBS 基线”点击“让 Agent 准备项目方案”，AESE 返回：

```text
facility_project_agent_failed
decode facility project options: unexpected EOF
```

IAOS 需求、场址控制和投资边界均已存在，失败发生在 Agent 候选输出进入业务校验之前。

## 根因

1. 项目 prompt 要求 3 套方案，每套最多 12 个 WBS，输出容易达到 MiniMax completion token 上限。
2. MiniMax adapter 没有把 `finish_reason=length|max_tokens` 视为截断，而是把半截 JSON 当作成功内容。
3. 场址 Proposal 已有一次治理修订，项目/WBS 生成却没有相同恢复路径。

## 修正

- MiniMax adapter 在解析业务 JSON 前识别 token 截断和空 content，保留 request ID、finish reason 和 token usage 证据。
- 首轮改为 2–3 套精简方案，优先 4–6 个 WBS，不放松金额、日期、阶段、序号和分摊合计校验。
- 截断、空输出、JSON 非法或业务校验失败时，只允许一次受控修订：恰好 2 套方案、每套 4 个 WBS、四个阶段各一个。
- 修订轮仍使用 IAOS 权威 Requirement 和 Site Control；不使用静态 fixture，不自动降低投资限额。
- 首轮和修订轮 token usage 合并到 Proposal Evidence；第二轮仍失败则继续失败关闭。

## 验证与恢复

- `TestMiniMaxCompleteJSONRejectsLengthTruncation`
- `TestAIPlanningProviderRepairsTruncatedProjectOptionsOnce`
- `go test ./...`
- `go vet ./...`

部署后刷新 M10 游戏，在项目办公室重新点击“让 Agent 准备项目方案”。该动作只生成候选，选择方案并确认后才进入 IAOS 基线审批。
