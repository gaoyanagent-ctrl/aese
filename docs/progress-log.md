# AESE Progress Log

本文件用于记录 AESE 项目的实质性进展。后续 agent 每次推进项目后，都必须在这里追加记录，方便其他 agent 快速掌握上下文。

更新格式：

```text
## YYYY-MM-DD - 简短标题

- 变更：
- 原因：
- 影响：
- 验证：
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

## 2026-06-26 - 华辰种子数据计划

- 变更：新增 `docs/HCTM_Seed_Data_Plan.md`，将 HCTM MVP 的基础主数据、演示初始业务数据、初始库存、订单、采购、生产任务、22 步演示事件序列、关键事件 payload 和 Agent 期望输出整理为可脚本化 seed 清单；同时更新 `README.md` 和 `docs/agent-project-context.md` 的文档入口。
- 原因：AESE 需要从模型规格推进到可初始化的演示数据，后续才能生成 JSON/SQL/Go seed 并在 IAOS 中重放第一条演示故事。
- 影响：M2 的 seed 数据计划已成形，主数据模型、事件模型和演示初始数据之间有了统一编码和导入顺序。
- 后续：编写 `docs/HCTM_Demo_Story_01_Order_Expedite.md`，把第一条演示故事转成面向用户操作、系统事件、Agent 输出和验收标准的可执行演示脚本。

## 2026-06-26 - 第一条演示故事脚本

- 变更：新增 `docs/HCTM_Demo_Story_01_Order_Expedite.md`，把“客户追加订单下的交付承诺重算”整理为可执行演示 runbook，包含演示前置条件、角色、视图建议、七幕流程、系统事件、页面展示、Agent 输出、事件流验收、Agent 验收、页面验收和失败条件；同时更新 `README.md` 和 `docs/agent-project-context.md` 的文档入口。
- 原因：AESE 需要一条能被产品、研发、销售和后续 agent 共同理解和执行的端到端故事，作为从文档模型进入可运行仿真的验收锚点。
- 影响：M2 文档闭环基本完成，已经具备虚拟企业蓝图、主数据模型、事件模型、种子数据计划和第一条演示脚本。
- 后续：进入实现准备阶段，建议先生成 `seed/hctm/*.json` 数据文件或设计 IAOS metadata/entity seed 转换方案。

## 2026-07-19 - 工程治理完善与 M3 开发规划

- 变更：重写 `AGENTS.md`，新增 `docs/README.md`、`docs/architecture.md`、`docs/code-map.md`、`docs/roadmap.md`、ADR-001、DES-001 和 M3 active plan；同步更新 README、Agent Context 和 MVP Blueprint 的当前状态与下一步。
- 原因：原仓库只有轻量规则和高层里程碑，缺少与 `iaos-go` 类似的文档索引、代码导航、架构边界、状态权威来源和可执行任务计划，且多处状态已过期。
- 影响：AESE 与 IAOS 的仓库职责已经固定；M3 被拆成场景包、校验器、兼容性报告、IAOS apply、O2D replay 和 closeout 六个切片，共 30 项可追踪任务。
- 验证：核对 AESE 全部现有文档和 `/iaos/iaos-go` 的 AGENTS、agent context、code map、eventdef、O2D 入口及动态实体 API；本地 Markdown 相对链接检查无缺失，M3 plan 确认包含 T1-T30 共 30 项任务，`git diff --check` 通过。
- 后续：从 M3 S1/T1-T4 开始，创建 `scenario-packs/hctm` 的 manifest、record sets、故事数据和 JSON Schema。

## 2026-07-19 - M3 场景包与离线执行链完成，在线 tracer 缺口固化

- 变更：新增 `hctm@0.1.0` 场景包（80 条 L1 主数据、14 条故事初始记录、22 个事件、17 条离线断言、2 条 IAOS 断言和 4 个 JSON Schema）；新增 Go CLI、loader、分层 validator、inspect、IAOS client、dry-run/apply/replay/verify 协调与安全 reset 计划；新增 compatibility report 和本地 runbook；更新项目入口、架构、code map、roadmap、DES-001 和 M3 checklist。IAOS `main` fast-forward 合并 commit `0260f28`，增加 decimal BOM、订单确认 CAS/no-op、trace metadata 及 DES-047/DES-048。
- 原因：把 M2 Markdown 合同转成可解析、可计算、默认零写入的执行资产，并通过真实平台取证区分 AESE 已实现能力和 IAOS 尚缺平台合同。
- 影响：实现核对发现原故事可发库存 11,700 却安排发运 12,000 的矛盾；为保持业务真实性，第 22 个事件修订为请求 3,000、实发 2,700、短缺 300，最终状态为 `partially_shipped`。M3 的 S1-S3 已完成，S4-S6 的代码和文档主体已完成，但 T20-T24 在线实证仍未完成，因此 M3 保持 active。
- 验证：`go test ./...`、`go vet ./...`、真实 pack `validate`/`inspect`、Draft 2020-12 Schema 校验和 `git diff --check` 通过；inspect 输出 80 master、14 initial、22 events、17 assertions。Platform `/health` 为 UP、`/ready` 的 DB/EventBus 为 OK。`tenant-hctm` apply dry-run 以 customer schema 404 fail closed，verify 的 work_order/inventory 两条断言同样明确失败；只读数据库复核 schema/customer/outbox 均为 0，证明 dry-run 零写入。IAOS platform 与 O2D 全量 Go 测试由独立 worktree 验证通过；commit 合并后又从 IAOS `main` 重新构建部署，Platform 与 O2D 进程均运行主 checkout 二进制并成功连接 DB/NATS。
- 后续：实现 DES-047 的受治理 scenario apply/reset、HCTM 稳定编码到 legacy UUID 映射和 tenant-hctm schema/workflow seed；补齐 workflow/event 去重、跨节点失败原子性和 work_order API 对齐，然后执行 T20-T24，保存 correlation、Outbox subject/event ID、O2D 日志、库存/工单结果和第二次运行 no-op 证据。

## 2026-07-19 - M3 受治理 O2D tracer 完成

- 变更：AESE 新增 HCTM→IAOS legacy projection，`apply` 改用 IAOS DES-047 原子 scenario endpoint，`replay` 支持 scenario apply 返回的 order UUID 并传递固定 correlation/idempotency，`reset` 接入服务端 L2/派生状态清理；expected outcomes 增加 2 条 IAOS work order 断言。IAOS `main` 合并 scenario apply/reset、O2D workflow 原子幂等、work_order metadata/workflow seed，以及真实执行发现的 dry-run、reset、correlation 和 tenant 显式绑定修正。新增 `docs/reports/hctm-m3-execution-evidence.md`，同步 README、Agent Context、architecture、code map、roadmap、DES-001、M3 plan 和 runbook。
- 原因：用户要求不分阶段完整执行计划；离线实现后继续关闭真实平台阻塞，并用实际 dry-run/apply/replay/reset 暴露和修正单测无法发现的跨服务问题。
- 影响：M3 T1-T30、S1-S6 全部完成。`tenant-hctm` 当前保留 1 customer、6 product、5 BOM、5 inventory、1 order/line、3 work orders 和 completed workflow，可直接演示。scenario reset 能删除 6 个 L2 对象与本轮派生工单/workflow，同时保留 12 个 L1。legacy 表未全面 FORCE RLS 的长期风险仍存在，但 scenario adapter 已通过每条 SQL 显式 tenant 条件关闭本合同的越界路径。
- 验证：dry-run 18 insert 且数据库零写；首次 apply 18 insert；第二 run 18 no-op 且对象数不变；`tenant-other` dry-run 18 insert 且目标租户零数据；O2D 完成 `corr-so-202607-0001` / event `evt-conf-d2f7c859b9e7d9fd10a7bd1a` / workflow `af706c43-b080-42de-8c98-b421d1b9e815`，decimal BOM 得到铝板 12,600，生成 3 个工单；第二次 replay 返回 `already_confirmed`，确认 Outbox 数不变；在线 verify 2/2；reset dry-run/apply 均显示删除 6、保留 12，reset 后恢复再次通过。AESE `go test ./...`、`go vet ./...`、Schema/Markdown 链接/diff checks 通过；IAOS platform/O2D tests、real-PG atomic/idempotency integration 和主 checkout 重部署通过。
- 后续：进入 M4 时实现 DES-048 的外生 simulation ingress，把供应商延期、设备故障和来料不良接入同样的权限、RLS、幂等、审计与 Outbox 边界；继续推动 legacy 表 FORCE RLS 平台 hardening。

## 2026-07-19 - 快速 2D 企业沙盘提升为当前里程碑

- 变更：新增 ADR-002、DES-002 和 PLAN-M3V-001，将只读 2D 场景预览器明确为 AESE 可拥有的产品验证界面；把 M3V 插入为当前 active 里程碑，并同步 README、Agent Context、Architecture、Roadmap、Code Map 和文档索引。
- 原因：原路线要到 M6 才出现 2D 沙盘，用户看到可用产品形态的时间过晚。现有 M3 已具备 80 条主数据、22 个事件和确定性结果，足以先形成可见、可操作的预览版。
- 影响：下一步不等待 M4/M5 完成，先在 3 到 4 个工作日内实现 React 2D 工作台、A 线画布、时间线、事件流、KPI、对象详情和 Agent 建议。首版不新增业务后端，不复制 IAOS 运行时，并通过 `ScenarioDataSource` 为后续 IAOS API/SSE 接入保留边界。
- 验证：M3V 计划拆为 V0-V4、T1-T26，定义每日可见成果、功能/视觉/边界测试和完成标准；本地 Markdown 相对链接检查无缺失，active plan 数量为 1，`git diff --check` 通过。
- 后续：从 V0/T1 开始创建 `frontend/` 和 `preview.json`，第 1 天结束前交付可缩放、可点击的苏州基地 A 线画布。

## 2026-07-19 - M3V 快速 2D 企业沙盘完成

- 变更：一次性完成 PLAN-M3V-001 的 V0–V4、T1–T26；新增 React + TypeScript + Vite 前端、14 节点/13 连线 A 线画布、七幕/22 canonical 事件 `preview.json`、`ScenarioDataSource` 合同、确定性播放 reducer、事件/KPI/对象详情和计划/质量/经营分析 Agent 面板；新增 M3V runbook、验收报告和三个固定视口截图，并同步 README、Agent Context、Architecture、Code Map、Roadmap、Blueprint、DES-002 和文档索引。
- 原因：用户要求不分阶段直接执行完整快速 2D 沙盘计划，并明确允许把可独立的场景数据、播放内核和画布实现交给 sub agent 并行推进。
- 影响：M3V 从计划态转为 Completed。AESE 现在已有可访问的产品预览界面，但仍严格保持只读 Preview 边界；浏览器只应用预计算 delta，不复制 IAOS 的 MRP、流程、权限或 Agent Runtime。下一优先级转为 M4 受治理异常事件入口和后续 `IaosScenarioDataSource`。
- 验证：`npm run typecheck`、ESLint 0 warning、Vitest 5 files/18 tests、Vite build、Playwright 3 projects/9 tests、npm audit 0 vulnerabilities、Go test/vet、preview 七幕/22 事件/3 Agent 合同检查、Markdown links 和 `git diff --check` 通过；1440×900、1280×720、390×844 截图人工检查无阻塞性重叠或整页横向溢出；开发服务绑定 `0.0.0.0:4173`，本机 HTTP 探测返回 200。
- 后续：M4 为供应商延期、设备停机和来料不良实现 IAOS simulation ingress；基于同一 `SandboxScenario` 视图模型增加 `IaosScenarioDataSource`，保留 Preview/Live 明示与受治理写入边界。

## 2026-07-19 - M4 设备停机受治理入口贯通

- 变更：IAOS `main` 合并 `9a8f5ca`、`463abd6`、`153a97a`，新增 DES-048 `POST /api/v1/simulation/events` 首个 `eam.machine.down` allowlist、动态设备解析、状态 CAS、幂等审计和事务 Outbox；AESE client/replay 接入 canonical 设备停机事件，对不完整或未提交的 2xx 响应失败关闭，并新增 M4 active plan 与执行证据。
- 原因：M4 需要证明外生仿真事实可进入 IAOS 现有权限、租户和事件治理边界，而不是从 AESE 直接发布 NATS；真实执行还发现设备属于 tenant 动态物理表，以及 PostgreSQL text advisory lock 不能包含 NUL。
- 影响：`LAS-WLD-02` 已可由 HCTM 事件稳定解析并从 `running` 转为 `maintenance`；重复重放返回相同事件，碰撞和跨租户失败关闭。当前只完成设备停机，供应商延期和来料检验失败仍是 M4 active 范围，O2D 尚未消费该事件。
- 验证：首次 HTTP 200/committed、重复 HTTP 200/duplicate、碰撞 409、跨租户 404；数据库仅 1 条 ingress 和 1 条 `PROCESSED` Outbox、目标外租户 0 条；AESE 22 事件 canonical replay 成功并将设备事件识别为 duplicate；client wire contract 与 malformed success 回归测试、IAOS Platform 测试/vet、部署健康检查与 AESE Go 测试/vet 通过。
- 后续：按同一合同实现 `o2d.supplier_delivery.delayed` 和 `qms.incoming_inspection.failed`，完成三类异常统一验收后再固定 M5/M6 消费合同。

## 2026-07-19 - M4 canonical 异常 replay 泛化

- 变更：AESE replay 将 `o2d.supplier_delivery.delayed`、`eam.machine.down` 和 `qms.incoming_inspection.failed` 统一投影为 `IngestSimulationEvent` 请求；业务对象类型和稳定编码只从 canonical metadata 获取，source 固定为 `aese:<pack>/<story>`，payload 原样透传；订单确认继续独立调用 decompose，其余事件保持 unsupported。
- 原因：供应商延期和来料检验失败需要复用设备停机已经验证的受治理入口与 fail-closed 响应边界，避免为每类异常复制请求构造和成功判定逻辑。
- 影响：AESE 已具备三类 M4 异常的统一 replay 适配；IAOS 对新增两类事件的采购/检验对象解析、状态变化和真实运行验收仍是 M4 未完成项。
- 验证：新增两类 canonical 请求、三类 dry-run 零写入、metadata 缺失/错配、完整 duplicate、malformed duplicate、其他事件 unsupported 以及原有 machine/order 回归测试；`go test ./...`、`go vet ./...` 和 `git diff --check` 通过。
- 后续：在独立 IAOS worktree 完成两类事件的入口实现后，执行统一权限、RLS、审计、Outbox、重复和碰撞验收。

## 2026-07-19 - M4 采购与来检对象 projection 前置

- 变更：canonical initial-state 新增 `IQC-202607-0002` pending 来料检验单，故事初始记录从 14 条增至 15 条；DES-047 legacy projection 新增两张 `purchase_order` 和该 `inspection_order`，稳定自然键分别为 `po_no` 和 `inspection_no`。采购日期保持 DateOnly，待检验单不虚构尚未产生的 receipt/lot。
- 原因：供应商延期和来料检验失败的 simulation ingress 必须在当前租户内解析到稳定采购单和检验单，真实 replay 前需先由 scenario apply 原子准备这些业务对象。
- 影响：AESE scenario request 现在包含 21 个对象；inspection 的 `po_no` 与 `material_code` 引用加入离线完整性校验。2D preview 已预置同编码 pending 检验对象，视图数据无需改动。IAOS fixture 已确认 7 字段采购 wire、DateOnly 日期以及 receipt/lot 可空的预分配检验单合同。
- 验证：projection 测试覆盖对象数量、自然键、全部 wire 字段、DateOnly、可选 receipt/lot 和 dropped-field warning；真实 pack `validate`、`inspect`（80 master、15 initial、22 events、17 assertions）、`go test ./...`、`go vet ./...` 和 `git diff --check` 通过。
- 后续：使用更新后的 scenario apply fixture 执行两类异常的首次提交、重复、碰撞、跨租户、状态变化和 Outbox 统一验收。

## 2026-07-19 - M4 replay 与 projection fail-closed 加固

- 变更：simulation success response 改为精确验证目标 tenant subject；无显式 tenant 的内部调用也只接受合法 `iaos.<tenant>.<event-type>`。DES-047 projection 对采购 7 字段和检验 8 字段显式必填，validator 只对 purchase/inspection 的合同必需引用报告缺失，不改变其他 optional reference 语义。
- 原因：只校验 subject 后缀会误接受其他 tenant 或任意前缀的成功回显；通用 mapping/引用逻辑会静默省略 M4 wire 必填字段，两者都会削弱真实 replay 的失败关闭边界。
- 影响：canonical pack 字段漂移会在离线 validator 或 projection 阶段失败；replay 不再把错误 tenant subject 计为成功。runbook 已按当前 21-object apply 和 9-L2 reset 计划更新，兼容性报告明确采购/检验对象已进入 M4 窄合同。
- 验证：新增 wrong-tenant、wrong-prefix、PO/inspection 缺字段、缺 required reference 和真实 pack 22 事件路由测试；真实 pack dry-run 识别 3 个 simulation candidate、1 个 decompose 和 18 个 unsupported，apply fake 路径触发 3+1。`go test ./...`、`go vet ./...`、pack validate/inspect、JSON Schema 和 `git diff --check` 通过。
- 后续：在 IAOS fixture 与两类 ingress 合并后执行真实 apply/reset/replay，并用实际回显确认 21-object apply 和 9-L2 reset 计数。

## 2026-07-19 - M4 三类受治理异常入口完成

- 变更：IAOS 完成供应延期和来检失败的稳定对象解析、严格 payload、状态影响、幂等审计和事务 Outbox，并补齐采购/检验 scenario fixture；AESE 完成 21-object projection、三类 canonical replay 和精确 tenant subject 失败关闭。M4 plan、roadmap、evidence、architecture、code map 和项目入口同步转为 Completed。
- 原因：M4 的完成标准不是直接发布消息，而是让三类外生事实在同一权限、租户、事务和可重复性边界内形成可查询业务上下文，供后续 Agent 和在线沙盘消费。
- 影响：`LAS-WLD-02`、`PO-202607-0001` 和 `IQC-202607-0002` 均有稳定受治理状态；事件常量、simulation response、租户 subject、metadata/entity query 和 Outbox 构成 M5/M6 的消费合同。领域消费者、自动重排产、Agent Runtime 和 `IaosScenarioDataSource` 未提前计入 M4。
- 验证：21-object dry-run 为 9 insert/12 no-op，apply 后第二次为 21 no-op；首次 canonical replay 3 triggered/19 skipped/0 failed，第二次 0/22/0；三类 ingress/Outbox 各 1 条，采购 ETA、IQC 数量/缺陷/批次/严重度及设备状态均落库；O2D workflow completed 并生成 3 张工单。AESE test/vet/validate/inspect、IAOS 各模块 test/vet、real-PostgreSQL `-race` integration、部署健康检查和 diff checks 通过。
- 后续：进入 M5 时只通过 IAOS Capability / AI Tool Registry 为计划、质量和经营分析 Agent 暴露受治理读写工具；M6 再实现 `IaosScenarioDataSource`。继续推进 legacy FORCE RLS、tenant-safe composite foreign key 和 metadata 版本排序 hardening。
