---
id: ADR-001
title: AESE 与 IAOS 仓库职责边界
date: 2026-07-19
status: accepted
author: Codex + User
tags: [architecture, repository, iaos, scenario-pack]
---

# AESE 与 IAOS 仓库职责边界

## 背景

AESE 需要版本化虚拟企业、场景数据和演示故事，同时必须复用 IAOS 的元数据、事件、权限、流程和 Agent 运行时。如果所有能力都在 AESE 重建，会形成第二套 ERP/MES；如果所有内容都直接硬编码进 IAOS，行业场景又无法独立演进。

## 决策

AESE 仓库拥有：

- 虚拟企业和业务故事规格。
- 机器可读 scenario pack、seed 和 expected outcomes。
- JSON Schema、离线 validator、inspect 和 replay orchestration 工具。
- IAOS 兼容性报告和演示 runbook。

IAOS 仓库拥有：

- 数据库、RLS、metadata runtime 和动态实体。
- Outbox、NATS、Scenario Package handler。
- Capability、Process、Policy、Decision、AI Tool 和审计。
- 业务 UI 和运行期查询 API。

AESE 通过版本化合同和受认证接口驱动 IAOS，不直接写 IAOS 数据库。正式事件通过 IAOS Outbox 或受治理 simulation ingress 产生。

## 结果

正面影响：

- 行业场景可独立版本化和测试。
- IAOS 保持唯一业务运行时和安全边界。
- 同一 AESE pack 可用于开发、演示和回归测试。

代价：

- 需要维护 AESE schema 与 IAOS DTO/metadata 的合同测试。
- 跨仓库功能需要两个独立提交和清晰依赖。
- IAOS 缺少的 simulation ingress 需要单独设计和实现。

## 被否决方案

1. 在 AESE 中建立独立数据库和业务服务：会复制 IAOS 能力，否决。
2. 把所有 HCTM seed 硬编码进 IAOS bootstrap：耦合平台和演示客户，否决。
3. 直接向 NATS 发布所有事件：绕过事务、权限和审计，只允许本地显式调试模式。
