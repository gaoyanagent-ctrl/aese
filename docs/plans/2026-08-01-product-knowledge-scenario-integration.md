---
id: PLAN-KNOWLEDGE-AESE-001
title: AESE 场景知识与 IAOS 产品知识中枢集成计划
date: 2026-08-01
status: active
author: Codex + User
tags: [knowledge, m9, iaos, aese]
parent_plan: PLAN-GXZ-001
---

# 定位

本计划是 PLAN-GXZ-001 下并行知识子计划，不替代 M9 财务计划。IAOS 主实现见 DES-077 和
PLAN-KNOWLEDGE-001；本仓只维护场景知识、映射、校验和 AESE 页面上下文。

## S0 合同

- [x] S1 固定 AESE/IAOS 内容所有权和三类事实边界。
- [x] S2 固定场景 Article 扩展、来源、版本和 Agent 回答合同。
- [x] S3 建立 M9 用户手册第一版。

## S1 机器发布

- [ ] S4 建立场景知识 JSON Schema 和 M9 article manifest。
- [ ] S5 校验 World action、IAOS Capability/Process/Entity 稳定编码引用。
- [ ] S6 生成签名 Knowledge Edition 并经 IAOS 平台包幂等安装。

## S2 页面与 Agent

- [ ] S7 World 节点详情增加“这一步是什么”上下文入口。
- [ ] S8 将 workspace/case/world run/node/actor 传入受治理知识问答。
- [ ] S9 展示 World 与 IAOS 双侧实际证据和配置漂移。

## S3 验收

- [ ] S10 M9 全节点人工、Agent、Approval、World 分类覆盖率 100%。
- [ ] S11 越权、旧版本、缺失资产和无运行证据失败关闭。
- [ ] S12 UI/API/runbook、Atlas、code map 和 evidence 全部收口。
