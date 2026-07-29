---
id: SOL-003
title: Genesis 新案件投影与 MiniMax JSON 截断 502
date: 2026-07-28
status: completed
author: Codex + User
tags: [genesis, minimax, projection, error-mapping]
---

# Genesis 新案件投影与 MiniMax JSON 截断 502

## 现象

新 Workspace 进入身份工作室时，浏览器同时看到：

- `game/incorporation/:case/projection` 返回 502；
- `game/creative/names` 返回 502。

## 根因

两条请求无共同网络故障：

1. 新 Workspace 尚无 incorporation case，IAOS evidence 正常返回 404，但 AESE 将所有
   IAOS trace 错误统一包装为 502。
2. MiniMax M3 的 reasoning 与四组候选超过 `max_tokens=2048`，content 在 JSON 中途
   截断，严格解析报 `unexpected end of JSON input`。

## 修复

- IAOS trace 404 映射为 AESE 404 `incorporation_case_not_found`；前端只把 404 解释为
  “进入身份创建态”，不再把真实 502 当作空案件。
- M3 输出预算提升到 8192。
- 第一次 JSON 非法时允许一次低温度严格 JSON 重新生成；第二次失败显式报错。
- reasoning 仍与业务 JSON 分离，不从 `<think>` 中提取候选。

## 验证

- 回归测试覆盖 IAOS 404 状态保持和截断 JSON 单次 repair。
- 原创业构想真实 M3 调用返回 200，生成四个动态候选。
- 新案件 projection 返回 404，不再返回 502。
