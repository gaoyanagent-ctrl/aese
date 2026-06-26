# AESE Progress Log

本文件用于记录 AESE 项目的实质性进展。后续 agent 每次推进项目后，都必须在这里追加记录，方便其他 agent 快速掌握上下文。

更新格式：

```text
## YYYY-MM-DD - 简短标题

- 变更：
- 原因：
- 影响：
- 后续：
```

## 2026-06-25 - 项目文档初始化

- 变更：新增 `README.md`、`AGENTS.md`、`docs/agent-project-context.md`、`docs/AESE_MVP_Blueprint.md` 和本进展日志。
- 原因：需要让后续 agent 能快速理解 AESE 的定位、MVP 范围、与 IAOS 的关系，以及每次进展如何记录。
- 影响：AESE 从原始构思文档进入项目化阶段，MVP 暂定聚焦华辰热管理系统集团、苏州制造基地、电池冷却板产品族、订单到交付主线、三类异常和三个 Agent。
- 后续：初始化 Git 仓库并连接 GitHub remote；继续补充华辰热管理系统集团详细虚拟企业蓝图。

## 2026-06-25 - 华辰虚拟企业蓝图

- 变更：新增 `docs/HCTM_Virtual_Enterprise_Blueprint.md`，详细定义华辰热管理系统集团、苏州制造基地、电池冷却板 A 线、28 个关键主数据对象、18 个关键事件，以及第一条演示故事线的输入和预期输出；同时在 `README.md` 中加入该文档入口。
- 原因：AESE 需要从概念蓝图推进到可建模、可 seed、可事件化、可演示的业务蓝图，方便后续 IAOS metadata、Scenario Package、Capability 和 Agent 设计接力。
- 影响：M1 虚拟企业蓝图的主体已经成形，MVP 锚点进一步收敛为 `HCTM-BCP-A01` 电池冷却板组件、苏州制造基地电池冷却板 A 线、客户追加订单下的交付承诺重算。
- 后续：将 28 个对象转为 IAOS metadata/entity 草案，将 18 个事件映射到 IAOS event subject 和 payload 规范，并补充种子数据清单。

## 2026-06-26 - 华辰主数据建模规格

- 变更：新增 `docs/HCTM_Master_Data_Model.md`，把 28 个关键对象整理为 IAOS entity 建模规格，包含命名约定、通用字段、字段定义、关系、状态、seed 示例、MVP 关系图和最小 seed 集合；同时更新 `README.md` 和 `docs/agent-project-context.md` 的文档入口。
- 原因：后续要把虚拟企业蓝图转成 IAOS metadata、seed 数据和事件 payload，需要先统一对象编码、字段、关系和状态模型。
- 影响：M2 的第一部分已经完成，后续 agent 可以基于该文档继续编写事件模型、种子数据计划，或开始设计 IAOS metadata entity seed。
- 后续：编写 `docs/HCTM_Event_Model.md`，把 18 个事件转成 IAOS subject、payload schema、上下游对象和 Agent/Capability 触发规格。

## 2026-06-26 - 华辰事件模型规格

- 变更：新增 `docs/HCTM_Event_Model.md`，把 18 个关键事件映射为 IAOS dotted event type、NATS subject、payload 字段、幂等键、上下游对象、Agent 触发矩阵、Capability / Process 接线建议和订阅建议；同时更新 `README.md` 和 `docs/agent-project-context.md` 的文档入口。
- 原因：AESE 后续要复用 IAOS Outbox + NATS + Scenario Package 机制，必须先让事件命名、payload 和触发关系与 IAOS 当前事件模型对齐。
- 影响：M2 的事件规格已成形，后续可以继续编写 seed 数据计划，或基于事件规格准备 `shared/eventdef` 常量和 payload struct 草案。
- 后续：编写 `docs/HCTM_Seed_Data_Plan.md`，把组织、客户、供应商、物料、BOM、工艺、设备、仓库、库存、订单和第一条演示事件序列整理成可 seed 的数据清单。
