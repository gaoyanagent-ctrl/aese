---
id: SOL-005
title: 设立案件与正式企业主数据物化
date: 2026-07-28
status: completed
author: Codex + User
tags: [genesis, iaos, incorporation, master-data]
---

# 设立案件与正式企业主数据物化

## 问题

游戏第一章提交了公司名称、注册地址和经营范围，但 IAOS 原先只保存于
`incorporation_case.state_document` 和 `m9_incorporation_case.payload`，案件表没有
显式列；登记成功后也没有形成正式法律主体主数据。

## 修复

IAOS commits `89f6c6f`、`0e02d64`：

- `incorporation_case` 增加案件名称、拟设企业名称、注册地址和经营范围列；
- 历史数据按 tenant 设置 RLS context 后从 JSONB 幂等回填；
- 新案件在同一事务同步显式列；
- `registration.observation.commit=registered` 后创建正式
  `m9_legal_entity`，物化企业名称、类型、地址、经营范围和 active 状态；
- JSONB 继续保留为状态机快照和审计证据。

## 验证

- 三个现存 `INC-GX-*` 案件已真实回填显式列。
- IAOS integration tracer 完成 18 个节点后验证案件字段和正式法律主体完全一致。
- 线上 8082 已部署。
