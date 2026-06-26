# AESE - Agentic Enterprise Simulation Environment

AESE, 中文名 **智能企业运行仿真环境**，是 IAOS 面向工业企业构建的可运行虚拟企业环境。

它不是传统演示账套，也不是单纯的 2D/3D 游戏。AESE 的目标是构建一个能够持续运行、产生业务事件、暴露经营问题、验证 IAOS 能力和 AI Agent 协同能力的工业企业数字沙盘。

第一阶段聚焦汽车零部件行业，虚拟客户暂定为 **华辰热管理系统集团**，主营新能源汽车热管理系统零部件，例如电池冷却板、冷却管路、热管理阀体和铝合金结构件。

## 项目目标

- 为 IAOS 提供一个真实感足够强的行业样板客户。
- 用虚拟企业反向验证 IAOS 的对象模型、流程、事件、规则、Capability 和 Agent 能力。
- 形成可演示、可测试、可扩展的企业运营仿真环境。
- 支撑汽车零部件客户演示、产品设计、研发测试和 Agent 训练。

## 第一阶段范围

AESE MVP 不追求一次覆盖完整集团运营，而是先跑通一个可解释的主线：

```text
客户订单
-> MRP
-> 采购
-> 来料检验
-> 生产排产
-> 工序生产
-> 完工入库
-> 发货
-> 开票
-> 经营分析
```

第一版只建一个集团、一个主工厂、一个产品族、一条主流程、三类异常和三个 Agent。

## 与 IAOS 的关系

AESE 应复用 IAOS 已有能力，而不是另起一套业务系统：

- 元数据驱动实体和表单。
- Outbox + NATS 事件流。
- Scenario Package 场景包。
- Capability Runtime。
- Policy / Decision / Process / Output。
- AI Tool Registry 和治理审计。

IAOS 是企业操作系统，AESE 是运行在 IAOS 上的行业仿真世界。

## 文档入口

- [Agent 项目背景](docs/agent-project-context.md)
- [MVP 蓝图](docs/AESE_MVP_Blueprint.md)
- [华辰热管理系统集团详细蓝图](docs/HCTM_Virtual_Enterprise_Blueprint.md)
- [华辰主数据建模规格](docs/HCTM_Master_Data_Model.md)
- [华辰事件模型规格](docs/HCTM_Event_Model.md)
- [华辰种子数据计划](docs/HCTM_Seed_Data_Plan.md)
- [进展跟踪](docs/progress-log.md)
- [原始构思记录](docs/ChatGPT_20260626.md)
