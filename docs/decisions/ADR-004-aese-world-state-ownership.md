---
id: ADR-004
title: AESE 拥有仿真世界事实，IAOS 拥有企业管理事实
date: 2026-07-22
status: accepted
author: Codex + User
tags: [aese, architecture, world-state, iaos, simulation]
---

# AESE 拥有仿真世界事实，IAOS 拥有企业管理事实

## 背景

AESE 1.x 把版本化企业场景交给 IAOS 执行，能够证明订单到交付、异常入口、Agent 查询和在线沙盘，但无法表达“设备已经异常而系统尚未登记”“实物已到货而系统库存未更新”等真实企业偏差。

AESE 2.0 需要区分三类状态：

- `World State`：企业客观世界实际发生的状态。
- `IAOS Business State`：企业在管理系统中登记、审批和审计的状态。
- `Actor Knowledge State`：人员或 Agent 在特定时刻可见、相信和记忆的状态。

这要求 AESE 从纯内容与无状态编排层升级为具有自身运行事实的仿真系统，同时不能在 AESE 中复制 ERP/MES、权限或流程引擎。

## 决策

AESE 拥有且只拥有：

- 仿真时钟、run、branch、随机种子和离散事件队列。
- 空间、设施、设备物理能力、实际资源、实际资金流等客观世界事实。
- 角色的观察、消息、认知、记忆引用和信息时效。
- 外生剧情、世界规则、结果计算、世界快照与可重放日志。

IAOS 继续唯一拥有：

- 订单、库存台账、设备台账、人员档案、财务凭证等企业管理记录。
- Metadata、Capability、Process、Policy、Decision、权限、RLS、Outbox、NATS、AI Tool 和审计。
- 人类与 Agent 对正式企业管理动作的授权和执行入口。

AESE 与 IAOS 不共享表，不跨库写入。世界事实通过版本化 `observation`、`intent`、`outcome` 合同与 IAOS 交互：

```text
World event
-> observation exposed to an actor
-> actor/human invokes governed IAOS capability
-> IAOS commits business record + audit + outbox
-> AESE consumes committed outcome
-> world rules calculate physical/economic consequence
```

首版 AESE World Store 使用独立 PostgreSQL 数据库和独立迁移边界，只保存仿真事实，不保存 IAOS JWT，不把 IAOS 业务表镜像为第二份主数据。它可以在本地与 IAOS PostgreSQL 由同一容器编排启动，但必须使用独立 database、独立账号和独立连接配置，禁止跨库查询和外键。

首版 Actor Knowledge 只保存结构化认知：事实引用、观察时间、事实有效时间、来源、置信度、可见范围和替代关系。不保存自由文本长期记忆；自然语言解释可以作为可再生输出或审计附件，但不是确定性仿真的权威状态。

首版经济守恒限定为：现金、承诺支出和应收回款。利润、完整总账、税务和复杂融资不作为 M8 完成条件。

## 与既有决策的关系

- ADR-001 继续约束业务运行时、权限和 IAOS 数据库所有权；其“AESE 不建立业务数据库”解释为“不建立第二套企业管理数据库”，不禁止独立的仿真事实存储。
- ADR-003 继续约束现有 HCTM 场景运行控制台。无状态 orchestration API 不直接演变为 World Runtime；两者共享合同和 IAOS client，但具有不同运行责任。
- 现有 `order-expedite-01` 场景包、CLI、Live snapshot 和验收证据继续保留为 AESE 1.x 兼容基线。

## 已确认事项

1. AESE 拥有独立、可持久化的仿真事实。
2. World Store 采用独立 PostgreSQL 数据库。
3. Actor Knowledge 首版只保存结构化认知，不保存自由文本长期记忆。
4. M8 经济守恒覆盖现金、承诺支出和应收回款。

`observation/intent/committed_outcome` 的字段合同和 IAOS gap audit 仍属于 PLAN-M8-001 F0，不改变本 ADR 的所有权结论。
