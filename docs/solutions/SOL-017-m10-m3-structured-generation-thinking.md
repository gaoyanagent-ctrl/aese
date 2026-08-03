---
id: SOL-017
title: M10 M3 结构化生成禁用自适应思考
date: 2026-08-03
status: completed
author: Codex + User
tags: [aese, m10, agent, minimax, thinking]
---

# M10 M3 结构化生成禁用自适应思考

## 现象

项目 Agent 已正确返回 retryable 503，但真实 `MiniMax-M3` 请求连续耗尽两个 40 秒窗口。

## 根因

M3 支持按请求切换 thinking；适配器未声明模式，因此严格 JSON 候选生成仍使用默认自适应推理。该任务已有 IAOS 权威输入、固定 Schema、业务校验和人工审批，不需要长链自主推理。实际同规模探针显示自适应模式输出更长且延迟波动更大。

## 修正

- `CompleteJSON` 且模型编码以 `MiniMax-M3` 开头时发送 `thinking: {"type":"disabled"}`。
- 企业命名等非通用 JSON adapter 调用不受本次开关影响。
- 保留 4096 token、40 秒单次窗口、一次完整性恢复、一次治理纠正和最终 retryable 503。
- 非思考模式结果仍必须通过全部治理门，不能直接写入项目事实。

## 证据

- disabled thinking：8.17 秒，`finish_reason=stop`，2682 内容字符，1030 completion tokens。
- 默认模式：15.09 秒，5909 内容字符，1786 completion tokens；真实会话中出现两次超过 40 秒。
- `TestMiniMaxCompleteJSONDisablesM3Thinking`
- `go test ./...`
- `go vet ./...`

MiniMax 官方说明 M3 可按请求关闭 thinking，以适配低延迟场景：<https://www.minimax.io/blog/minimax-m3>。
