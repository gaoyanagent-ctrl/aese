---
id: PLAN-M5-001
title: M5 受治理 Agent MVP 实施计划
date: 2026-07-19
status: completed
author: Codex + User
tags: [m5, agent, ai-tool, hctm]
---

# M5 受治理 Agent MVP 实施计划

## 1. IAOS 通用工具边界

- [x] 增加 metadata 约束的 `entity.records` query dispatcher。
- [x] 覆盖字段/filter allowlist、SQL identifier、limit、tenant/RLS 和未知 metadata 失败关闭。
- [x] 更新 DES-045 runbook、code map，并执行平台 test/vet、real-PostgreSQL integration 和部署。

## 2. HCTM Tool Bundle

- [x] 增加版本化 metadata schema 和只读 AI Tool manifests。
- [x] 增加默认 dry-run、显式 `--apply` 的 `agent-setup`。
- [x] 验证重复 setup 可收敛，其他 tenant 看不到 HCTM tools。

## 3. 三个 Agent Tracer

- [x] 实现计划 Agent 的订单、库存、BOM、采购、设备和工单上下文读取与建议。
- [x] 实现质量 Agent 的 IQC、采购和批次证据读取与隔离/追溯建议。
- [x] 实现经营分析 Agent 的交付/成本风险解释和数据完整性失败关闭。
- [x] 输出统一 recommendation envelope、object refs、tool call IDs 和 human-confirmation 标记。

## 4. 验收与收口

- [x] 执行三 Agent live tracer、重复调用和跨租户验证。
- [x] 确认业务表/Outbox 零增量，AI Tool audit/event 可追溯。
- [x] 更新 runbook、evidence、roadmap、architecture、code map 和 progress log。
- [x] 明确 M6 仅消费稳定 Agent 输出，不复制 Agent 规则。

经营分析 tracer 的 `partial` 是本计划定义的失败关闭完成态：当前 IAOS 缺少完工入库、发运和实际成本事实，因此不得用 Preview 的 11,700/300 代替在线结果。补齐这些业务事实属于后续业务链扩展，不阻塞只读 Agent MVP。
