---
id: DES-007
title: AESE 2.0 企业生命周期仿真基础架构
date: 2026-07-22
status: draft
author: Codex + User
tags: [aese, world-state, discrete-event, actor-knowledge, genesis]
---

# AESE 2.0 企业生命周期仿真基础架构

## 1. 目标

把 AESE 从“版本化业务故事 + IAOS 在线执行 + 2D 观察”升级为可持续推进时间、计算企业客观结果、保留信息差并由人类和 Agent 共同运营的企业生命周期仿真环境。

第一个纵向 Campaign 固定为：

> Project Genesis：华辰热管理集团从成立苏州制造公司，到苏州工厂完成第一批电池冷却板交付。

范围仍固定在华辰、苏州制造基地、电池冷却板 A 线和首批交付，不扩展第二工厂、第二产品族或 3D。

## 2. 现状与目标差距

| 维度 | AESE 1.x 当前事实 | AESE 2.0 目标 |
| --- | --- | --- |
| 运行核心 | 预编排的 22 事件故事 | 可暂停、单步、加速、分支的离散事件内核 |
| 状态 | 主要投影 IAOS 业务事实 | World / IAOS / Actor Knowledge 三态分离 |
| 结果 | expected outcomes 预先声明 | 由资源、时间和行为规则计算，再用不变量验证 |
| 企业范围 | 已存在工厂上的订单到交付 | 公司设立、工厂建设、团队、设备、工业化、首批交付 |
| Agent | 9 次只读查询生成三类建议 | 具备岗位、目标、权限、日历、任务和观察闭环 |
| 2D | A 线固定画布 | 世界层级、时间控制、状态差异和岗位接管视图 |
| 持久化 | AESE 无业务数据库 | 独立仿真事件日志与快照；IAOS 仍拥有管理事实 |

## 3. 目标模块地图

首版保持 Go 模块化单体，不采纳原始构思中的 Spring Boot 目录建议：

```text
cmd/aese-world/              # World Runtime 进程入口
internal/world/              # world/run/branch 聚合与状态投影
internal/simtime/            # 虚拟时钟、调度队列和推进策略
internal/simevent/           # 世界事件 envelope、日志和 reducer
internal/knowledge/          # actor observation/knowledge/belief
internal/rules/              # 资源、产能、资金和因果规则
internal/genesis/            # Project Genesis campaign 编排
internal/bridge/iaos/        # observation/intent/outcome 受治理适配
internal/experiment/         # checkpoint、fork 和方案比较（后置）
world-packs/hctm-genesis/    # 世界模板、campaign、规则与初始条件
frontend/src/world/          # World API adapter 与 view model
frontend/src/components/world/ # 层级地图、时间轴、差异与岗位界面
```

现有 `scenario-packs/hctm`、`internal/replay`、`internal/application` 和 Live UI 继续作为兼容层；M8 不做大爆炸式迁移。

## 4. 三态模型

### 4.1 World State

至少包含：仿真时间、组织存在性、地点与设施状态、资金余额、实际物料位置、设备物理状态、人员实际到岗与技能、项目活动、产能和客户/供应商外部行为。

### 4.2 IAOS Business State

只通过 IAOS 受治理 API、Outbox 和观察 API访问。AESE 保存稳定业务引用与已消费事件游标，不复制 IAOS 记录为新的权威主数据。

### 4.3 Actor Knowledge State

每条认知记录至少携带：`actor_id`、`fact_ref`、`observed_at`、`valid_at`、`confidence`、`source`、`visibility_scope` 和 `supersedes`。角色只能据其知识与 IAOS 权限决策，不能直接读取完整 World State。

### 4.4 差异与对账

三态不一致必须显式建模为 `discrepancy`，并通过盘点、检验、传感、审核、对账或调查关闭。差异不是自动同步错误。

## 5. 离散事件内核

核心实体：

- `WorldRun`：pack/version、timezone、seed、clock、status。
- `ScheduledEvent`：发生时间、优先级、因果链、payload、规则版本。
- `WorldEvent`：不可变事件、幂等键、前后状态摘要。
- `WorldSnapshot`：有界间隔快照，可从日志重建。
- `Checkpoint` / `Branch`：共享父历史，在决策点派生。

确定性要求：同一 world pack、规则版本、seed 和输入事件必须产生同一事件顺序、状态 hash 和 KPI。时间统一 RFC 3339，业务时区固定 `Asia/Shanghai`；金额与数量使用定点/decimal。

推进命令默认 dry-run，外部写入或推进必须显式 `--apply`；首版支持 pause、step、run-until 和固定倍速的服务端语义，UI 动画速度不参与业务计算。

## 6. 资源与经济不变量

M8 首个 tracer 只建立可验证的最小守恒集合：

- 物料：期初 + 收入 = 消耗 + 损耗 + 期末，批次和位置可追溯。
- 资金：期初现金 + 实收 = 实付 + 期末现金；承诺不等于实付。
- 产能：产出不能超过设备、人员、班次、材料、工装和质量门的共同上限。
- 人员：未到岗、未认证、超日历或已被占用的人员不能执行任务。
- 设施/设备：未建设、未安装、未调试、未验收的资源不能形成可用产能。

不变量失败时 run 失败关闭，并输出因果链和最小反例，不通过改写结果“修复剧情”。

## 7. IAOS 桥接合同

具体 envelope、journal/cursor、幂等、权限和失败恢复合同由已批准的 DES-008 固定。本节保留总体方向。

| 合同 | 方向 | 作用 |
| --- | --- | --- |
| `observation` | AESE -> actor/IAOS ingress | 把可感知世界事实暴露给特定角色或治理入口 |
| `intent` | human/Agent via IAOS -> AESE | 表达已授权的管理行动意图 |
| `committed_outcome` | IAOS -> AESE | 证明业务记录、审批或 Capability 已提交/no-op |
| `world_consequence` | AESE internal/outbound | 计算物理、时间与经济后果 |
| `reconciliation` | 双向引用 | 关联三态差异、发现动作与关闭证据 |

所有合同携带 `tenant_id`、`world_run_id`、`correlation_id`、`causation_id`、`idempotency_key`、虚拟发生时间和真实记录时间。正式 IAOS 写入仍必须走 API/Capability/Process 和 Outbox；AESE 不直写 IAOS DB 或直接发布正式 NATS。

## 8. Project Genesis 分阶段范围

1. `Foundation`：三态、时间、事件、规则、存储和桥接 tracer。
2. `Incorporation`：投资主体、注资、法人、治理、初始组织与预算。
3. `Plant Build`：选址、建设项目、设施依赖、预算与工期。
4. `Capability Build`：设备采购安装调试、招聘培训认证、仓储和质量设施。
5. `Industrialization`：产品、工艺、供应商、RFQ/APQP、试生产和 PPAP。
6. `First Delivery`：复用现有 O2D，补齐开票、回款、实际成本和项目盈亏。

每阶段必须交付一个可重放 tracer，而不是先横向建全套模块。

## 9. 迁移策略

- 保留 `order-expedite-01` 作为 Legacy Campaign 中的运营事件片段和回归夹具。
- 先用 adapter 将现有 22 事件映射到 World Event，不修改原 pack 语义。
- 新增 `world-packs/hctm-genesis`，不把生命周期字段塞进现有 record-set schema。
- 2D 前端先增加模式入口和三态差异面板，再演进为层级世界地图。
- M8 只证明单 run；checkpoint/fork 的数据结构进入合同，A/B UI 在确定性和存储验收后实现。

## 10. 非目标

- 多工厂、多产品族、集团全面运营。
- 高精度物理仿真、完整数字孪生或 3D。
- 全量 ERP/PLM/EHS/财务模块。
- 数百 Agent 自治、自由文本长期记忆或未经批准的自动执行。
- 用剧情脚本直接指定财务和生产成功结果。

## 11. 架构验收门

进入大规模业务模块实现前必须满足：

1. ADR-004 获批，World/IAOS/Knowledge 所有权无歧义。
2. 同一输入连续运行两次，event log、state hash 和 KPI 完全一致。
3. 至少演示一条“世界已变化、IAOS 未登记、角色尚未知”的偏差，以及发现和关闭过程。
4. IAOS 写入全部有权限、租户、幂等、Outbox 和审计证据。
5. AESE World Store 删除或重建不会改变 IAOS 业务事实；IAOS reset 不会误删世界历史。
6. 现有 M7 全链和 Preview/Live 回归不受影响。
