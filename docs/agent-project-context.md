# AESE Agent Project Context

本文是给后续 agent 的项目背景入口。进入 `/iaos/aese` 后，先按 `AGENTS.md` 的顺序阅读项目级入口，再按任务从 `docs/code-map.md` 选择 HCTM 领域文档。

## 1. 项目定位

AESE 是 **Agentic Enterprise Simulation Environment**，中文名 **智能企业运行仿真环境**。

它是 IAOS 的行业场景运行环境，用一个可运行的虚拟汽车零部件企业来验证 IAOS 的企业对象模型、业务流程、事件驱动、Capability、AI Agent、权限、审计和经营分析能力。

核心判断：

- AESE 不是传统样板账套。样板账套偏静态数据，AESE 要持续产生业务事件。
- AESE 不是独立游戏。视觉沙盘可以借鉴游戏表达，但业务运行和 IAOS 集成是核心。
- AESE 不是另一个 ERP/MES。它应该复用 IAOS 平台能力，成为 IAOS 的场景包、数据包、仿真包和演示环境。

一句话：

> IAOS 是企业操作系统，AESE 是一个在 IAOS 上运行、可被人和 AI Agent 共同操作的虚拟工业企业。

## 2. 虚拟客户设定

第一阶段虚拟客户暂定为：

**华辰热管理系统集团有限公司**

行业：

- 新能源汽车热管理系统零部件。

核心产品：

- 电池冷却板。
- 冷却管路。
- 热管理阀体。
- 铝合金结构件。

第一阶段主工厂：

- 苏州制造基地。

第一阶段主产品：

- 新能源汽车电池冷却板。

详细虚拟企业蓝图：

- `docs/HCTM_Virtual_Enterprise_Blueprint.md`

主数据建模规格：

- `docs/HCTM_Master_Data_Model.md`

事件模型规格：

- `docs/HCTM_Event_Model.md`

种子数据计划：

- `docs/HCTM_Seed_Data_Plan.md`

第一条演示故事脚本：

- `docs/HCTM_Demo_Story_01_Order_Expedite.md`

典型工艺：

```text
铝材来料
-> 压铸 / 成形
-> 机加工
-> 焊接
-> 清洗
-> 检漏测试
-> 总装
-> 包装
-> 入库
```

## 3. MVP 边界

MVP 只追求一个窄而完整的经营闭环：

```text
客户订单
-> MRP
-> 采购计划
-> 供应商交付
-> 来料检验
-> 生产排产
-> 工序执行
-> 过程检验
-> 完工入库
-> 客户发运
-> 开票
-> 经营分析
```

第一阶段必须避免范围失控：

- 不做完整集团所有工厂。
- 不做所有产品族。
- 不做完整财务、人力、EHS 闭环。
- 不先做 3D 游戏。
- 不为 AESE 另建一套与 IAOS 割裂的平台。

## 4. IAOS 依赖关系

AESE 的开发应优先参考 `/iaos/iaos-go`。

当前可复用能力：

- `scenarios/o2d`：订单到交付场景包雏形。
- `shared/eventdef`：统一事件定义。
- `platform/internal/metadata`：元数据与动态实体。
- `platform/internal/capability`：Capability Runtime。
- `platform/internal/process`：流程定义和运行。
- `platform/internal/policy`：策略规则。
- `platform/internal/decision`：决策、审计和解释。
- `platform/internal/aitool`：AI Tool Registry 和治理调用。
- `frontend`：IAOS 操作台。

AESE 的正确开发方向：

```text
虚拟企业蓝图
-> 主数据与场景数据
-> 事件流
-> 场景包
-> Capability
-> Agent 工具调用
-> 2D 企业沙盘
-> 经营仿真和决策反馈
```

## 5. Agent 角色

第一阶段只定义和验证三个 Agent：

- 计划 Agent：处理订单变化、库存、产能和交期约束。
- 质量 Agent：分析来料不良、过程缺陷、客户投诉和追溯链路。
- 经营分析 Agent：解释交付率、成本、库存、利润和异常对经营结果的影响。

Agent 不能只是聊天窗口。它们必须逐步具备：

- 读取业务上下文。
- 调用 IAOS Capability。
- 生成建议或草稿。
- 触发受治理的业务动作。
- 留下审计记录。
- 解释决策依据。

## 6. 当前状态

截至 2026-07-19：

- M0、M1 已完成。
- M2 的业务与技术规格文档已完成：28 个对象、18 类事件、Seed 数据计划和第一条演示故事均有明确合同。
- 工程治理已补齐文档索引、架构边界、code map、路线图、ADR、DES 和 active plan。
- 已生成 `hctm@0.1.0` 机器可读场景包、4 个 JSON Schema、Go loader/validator/inspect CLI，以及默认 dry-run 的 IAOS client/replay 基础实现。
- M3“可执行 HCTM 场景包”已完成：受治理 scenario apply/reset、tenant isolation、O2D Outbox/NATS/workflow、重复 no-op 和在线 verify 均有实际证据，详见 `docs/reports/hctm-m3-execution-evidence.md`。
- M3V 快速 2D 企业沙盘已完成：静态 `preview.json`、14 节点 A 线画布、七幕/22 事件播放、五项 KPI、对象详情和三类 Agent 建议已通过桌面与移动验收。
- M3V 实现位于 `frontend/`，运行见 `docs/runbooks/hctm-m3v-2d-sandbox.md`，验收证据见 `docs/reports/hctm-m3v-2d-sandbox-evidence.md`。
- 下一最高优先级是 M4：实现受治理 simulation ingress，让供应商延期、设备停机和来料不良进入 IAOS 运行链；之后增加 `IaosScenarioDataSource`。
- `/iaos/iaos-go` 已提供 DES-047 scenario apply/reset、原子幂等 O2D workflow 和 HCTM work_order/workflow fixture；Platform、PostgreSQL、NATS 和 O2D 当前运行正常。

## 7. 后续 agent 必须维护的信息

任何实质性进展都要更新 active plan、`docs/roadmap.md` 和 `docs/progress-log.md`，具体规则见 `AGENTS.md`。

更新内容包括：

- 做了什么。
- 为什么这样做。
- 对 MVP 范围或架构有什么影响。
- 留下了什么后续问题。

新增文档、代码入口或 IAOS 集成点时，还必须同步 `docs/README.md` 和 `docs/code-map.md`。
