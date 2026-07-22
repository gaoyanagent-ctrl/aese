---
id: DES-009
title: AESE World 三态术语、所有权与数据分类
date: 2026-07-22
status: approved
author: Codex + User
tags: [aese, world-state, contract, data-classification]
---

# AESE World 三态术语、所有权与数据分类

## 1. 三态术语

| 术语 | 定义 | 权威来源 |
| --- | --- | --- |
| World State | 仿真世界中已经客观发生的物理、时间、资源和最小经济事实 | AESE World Store |
| IAOS Business State | 企业通过登记、审批、Capability、Process 和事务确认的管理事实 | IAOS |
| Actor Knowledge State | 特定人员或 Agent 在给定时间已获知且有权访问的结构化认知 | AESE World Store（认知事实）；IAOS（访问授权） |
| Observation | 世界事实对指定角色变得可感知的消息，不代表 IAOS 台账已变化 | AESE 产生，IAOS journal 接收 |
| Intent | 人员或 Agent 经 IAOS 治理后提出的行动意图，不代表世界结果 | IAOS |
| Committed Outcome | IAOS 事务已提交或确定 no-op 的结果；唯一可驱动世界后果的 IAOS 输入 | IAOS |
| Discrepancy | World、IAOS、Knowledge 任意两态之间有业务意义的不一致 | AESE 显式记录并以证据关闭 |
| Snapshot | 从不可变事件日志投影出的有界恢复点，不替代事件历史 | AESE World Store |
| Canonical Hash | 对 JSON 语义值规范化后计算的 `sha256:` 摘要，用于幂等和重放对账 | 产生合同的一方 |

## 2. 对象所有权矩阵

| 对象 | AESE | IAOS | Actor Knowledge | 约束 |
| --- | --- | --- | --- | --- |
| WorldRun、branch、seed、虚拟时钟 | Own | stable ref / binding | 可见运行摘要 | 首版只有 `main` 分支 |
| WorldEvent、规则版本、世界快照 | Own | 不镜像 | 只能经 observation 获知 | 不通过 direct NATS 形成正式事实 |
| 设备实际状态、实际产能、实际物料位置 | Own | 登记/台账视图 | 按观察和权限获知 | stable business code 关联，不跨库外键 |
| 订单、库存台账、设备台账、财务凭证 | 不拥有 | Own | 通过 IAOS Tool/API 获知 | AESE 不复制为第二主数据 |
| Observation | 产生与引用 | journal、权限、审计 | 接收后形成 Knowledge | accepted 不等于台账变化 |
| Intent | 只消费引用 | Own、授权、审批 | 发起人可见 | 不允许任意 state patch |
| CommittedOutcome | 消费、去重、计算 consequence | Own、事务提交 | 按权限可见 | rejected/rolled_back 不伪装为 outcome |
| Knowledge | 保存结构化认知 | Owns access decision | 按 actor 隔离 | 不保存自由文本长期记忆 |
| Discrepancy | Own | 提供对账记录引用 | 可成为调查认知 | 关闭必须携带证据链 |
| JWT、RLS、权限、Process、Capability、Outbox | 不拥有 | Own | 不适用 | World Store 禁止保存 JWT |

## 3. 数据分类

| 分类 | 示例 | 存储/日志规则 |
| --- | --- | --- |
| Public | schema、枚举、虚构 pack 版本 | 可版本化提交 |
| Internal | WorldRun ID、规则版本、事件因果、state hash、fixture | 可进入 World Store 和受控日志；不得当作 IAOS 业务真相 |
| Confidential | 角色认知、设备观察、差异、实际资源/资金状态、IAOS record ref | tenant/run/actor 范围控制；日志最小化；备份加密并演练恢复 |
| Restricted | JWT、数据库密码、真实客户/个人数据、生产连接信息 | 不进入合同 payload、World Store、fixture、提交或普通日志 |

所有金额和守恒数量使用十进制字符串并显式 `unit`、`scale`；所有时间使用 RFC 3339，业务时区为 `Asia/Shanghai`。对象关联使用稳定业务编码，环境 UUID 只能作为 IAOS 返回的非权威 `record_refs`。

## 4. 机器合同入口

- JSON Schema：`world-contracts/schemas/`
- 可通过 schema 的最小 fixture：`world-contracts/fixtures/`
- Go 类型、strict parser 与 canonical hash：`internal/worldcontract/`
- PostgreSQL 边界：`deploy/world-postgres/`、`internal/worldstore/`
- 本地启动和备份合同：`docs/runbooks/aese-world-store.md`

这些入口只冻结 F0 合同，不实现 reducer、调度器、bridge adapter 或完整仿真内核。
