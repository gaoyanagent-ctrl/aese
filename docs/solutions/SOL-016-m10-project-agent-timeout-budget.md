---
id: SOL-016
title: M10 项目 Agent 模型超时预算
date: 2026-08-03
status: completed
author: Codex + User
tags: [aese, m10, agent, minimax, timeout]
---

# M10 项目 Agent 模型超时预算

## 现象

项目/WBS 候选生成返回不可重试的 422：

```text
MiniMax request failed: context deadline exceeded
```

## 根因

MiniMax 客户端单次最多等待 75 秒，而 AESE 整体请求只有 90 秒。第一次慢响应后，适配器重复同一份 8192-token 请求，但第二次只剩约 15 秒父上下文时间；最终网络超时又被项目 API 当成业务治理失败。

## 修正

- MiniMax 默认单次等待限制为 40 秒，为有界恢复保留窗口。
- 通用 JSON 规划调用的 timeout 不再由 adapter 重复同一提示，而由项目 Agent 使用精简修订提示恢复一次。
- 两套四阶段 WBS 使用 4096 completion token 上限。
- 最终仍超时时返回 `503 facility_project_agent_timeout` 和 `retryable=true`，并明确尚未写入项目事实。

## 验证

- `TestMiniMaxDefaultTimeoutPreservesRetryWindow`
- `TestAIPlanningProviderUsesCompletionRepairForTimeout`
- `TestAIPlanningProviderSeparatesCompletionAndGovernanceRepairs`
- `TestPlantProjectOptionsReportsProviderTimeoutAsRetryable`
- `go test ./...`
- `go vet ./...`
