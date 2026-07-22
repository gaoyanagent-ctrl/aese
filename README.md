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

## 当前状态

- M0 项目初始化：完成。
- M1 华辰虚拟企业蓝图：完成。
- M2 主数据、事件、Seed 和演示规格：文档完成。
- M2.5 工程治理：完成。
- M3 可执行 HCTM 场景包：完成；pack、Schema、CLI、受治理 apply/reset、O2D replay/verify 和幂等证据均已落地。
- M3V 快速 2D 企业沙盘：完成；七幕/22 事件、A 线画布、KPI、对象详情和三类 Agent 建议已通过桌面与移动端验收。
- M5 Agent MVP：完成；默认 dry-run 的 `agent-setup` / `agent-run`、9 个低风险只读 AI Tool bundle 和三 Agent 确定性 tracer 已通过 live、重复调用、零业务写入与跨租户验收。
- M6 在线 2D 企业沙盘：完成；受治理完工/入库/发运、snapshot/cursor/SSE、在线 Agent 建议和 Preview/Live 双模式已通过验收。
- M7 受治理场景运行控制台：完成；业务用户可从页面完成预检、初始化、逐幕运行、Agent 分析、验证和安全复位，CLI/UI 与 IAOS 审计副作用已完成对账。
- M8 AESE 2.0 基础：Completed；F0-F5 已交付机器合同、确定性内核、三态 tracer、受治理 IAOS bridge、Genesis pack 与 World Play。
- M9 Project Genesis 企业成立与治理：Completed；成立、登记、开户、资本、三岗位、mandate、预算与 `plant_project_eligible` 已通过两仓和三视口验收。
- M10 Project Genesis 工厂选址与设施建设：Completed；场址决策、设施世界、IAOS 治理、Plant Build Play 与 `capability_build_eligible` 已验收。

当前同一页面可运行 Preview 或 IAOS Live 沙盘，并可在联动中心通过无业务数据库的 AESE 薄编排 API 受治理运行和复位场景；浏览器不直接调用 IAOS 写接口。

## M5 Agent tracer

`agent-setup` 从 `scenario-packs/hctm/agent-tools.json` 安装 HCTM 最小 metadata 和 9 个 `source_ref=entity.records` 查询工具；`agent-run` 通过 IAOS Tool API 读取上下文并返回三类结构化建议。两个命令默认都不写入，只有显式 `--apply` 才会注册工具或产生受审计的 tool calls：

```bash
go run ./cmd/aese agent-setup ./scenario-packs/hctm --target http://127.0.0.1:8082
go run ./cmd/aese agent-run ./scenario-packs/hctm --story order-expedite-01 --target http://127.0.0.1:8082
```

这条 tracer 是可验证的只读建议链，不是独立 Agent Runtime，也不调用真实 LLM 或自动执行业务动作。

## 快速启动 2D 沙盘

```bash
cd frontend
npm install
npm run dev
```

访问 `http://localhost:4173/`。点击顶栏“联动中心”，选择华辰租户与订单场景后执行“一键连接并检查”，页面会自动验证 IAOS 身份、快照、事件通道及销售订单/工单/库存/设备数据；通过后直接进入 Live，无需手工复制 JWT 或调用 API。环境变量方式仍保留给自动化部署。详细操作见 [M6 运行手册](docs/runbooks/hctm-m6-online-sandbox.md)。

## 文档入口

- [Agent 项目背景](docs/agent-project-context.md)
- [文档总索引](docs/README.md)
- [系统架构与仓库边界](docs/architecture.md)
- [项目路线图](docs/roadmap.md)
- [Code Map](docs/code-map.md)
- [M3 实施计划](docs/plans/2026-07-19-m3-executable-scenario-package.md)
- [快速 2D 企业沙盘实施计划](docs/plans/2026-07-19-fast-track-2d-enterprise-sandbox.md)
- [M6 在线 2D 企业沙盘实施计划](docs/plans/2026-07-20-m6-online-2d-enterprise-sandbox.md)
- [M7 受治理场景运行控制台实施计划](docs/plans/2026-07-20-m7-governed-scenario-operations-console.md)
- [AESE 2.0 企业生命周期仿真基础设计](docs/designs/DES-007-aese-2-foundation.md)
- [AESE World 与 IAOS 三段式桥接合同](docs/designs/DES-008-world-iaos-bridge-contract.md)
- [M8 AESE 2.0 三态世界与仿真内核实施计划](docs/plans/2026-07-22-m8-aese-2-foundation.md)
- [M9 Project Genesis 企业成立与治理设计](docs/designs/DES-010-genesis-incorporation-and-governance.md)
- [M9 Project Genesis 企业成立与治理实施计划](docs/plans/2026-07-22-m9-genesis-incorporation.md)
- [M10 Project Genesis 工厂选址与设施建设设计](docs/designs/DES-011-genesis-plant-build.md)
- [M10 Project Genesis 工厂选址与设施建设实施计划](docs/plans/2026-07-22-m10-genesis-plant-build.md)
- [M10 Genesis Plant Build 运行手册](docs/runbooks/genesis-plant-build.md)
- [M10 Genesis Plant Build 验收证据](docs/reports/m10-genesis-plant-build-evidence.md)
- [M3 本地运行手册](docs/runbooks/hctm-m3-local-run.md)
- [M3V 2D 沙盘运行手册](docs/runbooks/hctm-m3v-2d-sandbox.md)
- [M5 受治理 Agent Tracer 运行手册](docs/runbooks/hctm-m5-governed-agent-tracers.md)
- [HCTM → IAOS 兼容性报告](docs/reports/hctm-iaos-compatibility.md)
- [M3 执行证据](docs/reports/hctm-m3-execution-evidence.md)
- [M3V 2D 沙盘验收证据](docs/reports/hctm-m3v-2d-sandbox-evidence.md)
- [M5 受治理 Agent Tracer 验收证据](docs/reports/hctm-m5-agent-evidence.md)
- [MVP 蓝图](docs/AESE_MVP_Blueprint.md)
- [华辰热管理系统集团详细蓝图](docs/HCTM_Virtual_Enterprise_Blueprint.md)
- [华辰主数据建模规格](docs/HCTM_Master_Data_Model.md)
- [华辰事件模型规格](docs/HCTM_Event_Model.md)
- [华辰种子数据计划](docs/HCTM_Seed_Data_Plan.md)
- [演示故事 01：客户追加订单下的交付承诺重算](docs/HCTM_Demo_Story_01_Order_Expedite.md)
- [进展跟踪](docs/progress-log.md)
- [原始构思记录](docs/ChatGPT_20260626.md)
