---
id: DES-025
title: M23 受治理多 Agent 组织
date: 2026-07-22
status: completed
author: Codex + User
tags: [m23, agents, organization, governance]
---

# M23 受治理多 Agent 组织

## 1. 目标

把现有确定性岗位 tracer 扩展为可评估的部门/管理 Agent 组织：共享企业目标但拥有不同职责、知识、预算和工具，能交接、质疑、升级和由人类接管。

机器终态：`agent_operating_model_qualified=true`。

## 2. 首版组织

- Planning、Procurement、Production、Quality、Maintenance/EHS、Finance 和 Executive Agent。
- 每个 Agent 有 mandate、objective、allowed tools/capabilities、calendar、knowledge scope、memory refs 和 escalation policy。
- 三条 benchmark：IBP gap、network disruption、field quality/plant incident。
- deterministic/reference policy 与可选 LLM policy 使用同一 contracts；没有模型/凭据时仍可离线验收。

## 3. 安全合同

- Agent 不读取完整 World State，只消费 observation、IAOS Tool 和授权 Knowledge。
- proposal、decision、approval 和 execution 分离；高风险动作永远要求独立人类/角色批准。
- memory 是有来源、版本、时效和 scope 的引用，不是无限自由文本真相库。
- prompt/model/tool/policy/version、input/output、reason、cost、latency 和 takeover 全部进入 evidence。

## 4. 评价

评价任务完成、业务 KPI、约束违反、工具正确性、证据引用、信息泄漏、协作/冲突、人工接管和可重放性。不能只以自然语言“看起来合理”判定成功。

## 5. 完成标准

- 七 Agent 在三个 benchmark 上与 human/reference baseline 可比较。
- tenant/role/knowledge/tool/capability 隔离和 prompt injection/越权失败路径通过。
- 重复调用、并发建议、陈旧知识、模型失败、超预算和 takeover 可恢复。
- Agent 不自动修改策略、计划或业务事实；所有动作沿 IAOS 治理链。

## 6. 非目标

无人值守自治企业、数百 Agent、自由互联网访问、无限长期记忆和未经审批的真实执行不在范围。
