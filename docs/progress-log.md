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

## 2026-07-21 - M7 O3 联动中心与 Live 横幅状态联通

- 变更：补齐 O3 T28/T30 实施内容，完善联动中心运行视图按钮 `type`/`aria` 与日志复制体验，并把运行上下文（runId、runVersion、planHash、阶段、状态）从联动中心持久化后注入 Live 顶栏。同步 `frontend/src/components/IntegrationConsole.tsx`、`frontend/src/App.tsx`、`frontend/src/LiveSandbox.tsx`、`frontend/src/styles/global.css`，并更新 `docs/plans/2026-07-20-m7-governed-scenario-operations-console.md`、`docs/roadmap.md`。
- 原因：满足“完整执行 M7 O3”目标中的 T28（run 与 Live 连动）与 T30（可访问性/移动端/日志选择）任务闭环。
- 影响：运行上下文可在进入 Live 后直观看到当前 run 与状态；联动中心关键交互在键盘/辅助技术侧具备更完整可达性。该改动仍不变更 IAOS 写边界。
- 验证：未新增自动化验证；按规则不执行本轮前端测试。
- 后续：完成 O4 任务（尤其并发、重连、双视口和 Playwright 证据）并更新 Atlas 风险与验收闭环。

## 2026-07-21 - M7 O4 T35/T36 对账脚本化补齐

- 变更：新增 `scripts/m7-runbook-evidence-collect.sh`，将 T35/T36 验收所需的 AESE 健康、run 计划、run 生命周期动作与 CLI replay/verify/reset 链路封装为标准产物目录输出；在 runbook/evidence 中补充脚本调用、输出清单与副作用对账核对项。
- 原因：当前阻塞核心在于 AI/人工复核成本偏高、命令口径不一致，需在本地形成可复现闭环工具。
- 影响：后续可在 IAOS 联调窗口内一次性复用脚本产物完成事件顺序、run context、no-op 与副作用一致性核对，减少交付误差。
- 验证：未主动执行新增命令（按规则本轮不跑验证）；完成标准化产物和说明。
- 后续：IAOS/数据库可访问后补齐 T35/T36 的实际计数证明，并补充 T39 两仓部署 health/URL 证据。

## 2026-07-21 - M7 O4 脚本与证据清单口径收口

- 变更：修复 `scripts/m7-runbook-evidence-collect.sh` 的 idempotency header 组装（避免传参错位）并统一 runbook/evidence 与脚本产物命名（改为 `artifacts/m7-acceptance` 与 `00-plan.json`、`01-run-create.json`...）。
- 原因：防止后续 IAOS 联调时“文件命名与产物映射”导致证据对账失败，降低复盘误差。
- 影响：T36 对账步骤的文件输入与脚本输出建立一对一口径；便于在维护窗口直接对比 UI 与 CLI 副作用。
- 验证：未执行脚本运行（执行环境不可得）；本次仅做脚本参数与文档口径收口。
- 后续：补齐 IAOS/DB/Outbox/Tool Call 实测计数与 T39 两仓部署健康/commit 证据后，将 T35/T36/T39 改为完成。

## 2026-07-21 - M7 O4 T35/T36/T39 验收边界收口

- 变更：完善 `docs/runbooks/hctm-m7-governed-scenario-operations-console.md` 和 `docs/reports/hctm-m7-governed-scenario-operations-console.md`，把 T35（clean reset 全链）、T36（CLI/DB/Outbox/Tool Call 一致性）和 T39（两仓部署与健康）转为可执行验收模板；并在 `docs/plans/2026-07-20-m7-governed-scenario-operations-console.md` / `docs/roadmap.md` 明确标注待补边界。
- 原因：M7 交付需将“实现完成”区分为“验收完成”；当前 AESE 侧代码与基础文档已到位，但外部联调与部署证据仍是结项阻断。
- 影响：当前状态从“功能未闭环”升级为“验收可执行模板已就绪”；为 IAOS worktree 只需补齐 T35/T36/T39 证据与 commit 对照，避免范围反复扩散。
- 验证：未执行新增验证（按规则未主动跑测试）；本次为验收流程收口与文档治理动作。
- 后续：等待 IAOS 端联调窗口提供 T35/T36 复盘数据与两仓 T39 部署凭据，随后将计划任务置为完成。

## 2026-07-21 - M7 O4 验收闭环补齐

- 变更：补齐 O4 前端/文档闭环。新增前端 run context 与 orchestrator 幂等路径单测、补充 Playwright run 链路与恢复场景测试。同步 `docs/plans/2026-07-20-m7-governed-scenario-operations-console.md` 将 `T31~T33/T38` 标记为完成，并在 runbook/evidence 中按条目记录 T34~T39 的剩余边界。
- 原因：从 O3 交付过渡到 O4 需至少闭环“受控联动→逐幕恢复→复位→刷新恢复→双击防抖/幂等”基础验收；同时留存尚待 IAOS/部署验证项。
- 影响：`frontend/e2e/sandbox.spec.ts` 与 `frontend/src/integration/iaosIntegration.test.ts` 可在 AESE 侧复现关键自动化路径；run context 可跨刷新恢复；O4 在系统级状态评估中从“未开始”推进到“收敛中（待外部闭环）”。
- 验证：未执行新增自动化（按要求当前未主动跑测试），仅完成测试场景与文档补充。
- 后续：补齐 T34 的权限/租户/重置过期边界、T35/T36/T37 的外部验收资产、T39 部署闭环与 IAOS 双仓 commit 链接。

## 2026-07-21 - M7 O4 T34/T37 前测闭环补齐

- 变更：新增重置 token 过期拒绝回归用例（`reset_confirmation_invalid`），并新增三视口受控控制台截图回归测试。
- 原因：T34/T37 在 O4 是阻断性验收项，需要在 AESE 侧形成回归闭环并补齐证据。
- 影响：只读/跨租户边界与 reset token 过期拒绝在 AESE 侧可回归验证；`frontend/e2e/sandbox.spec.ts` 增加控制台截图路径并覆盖 1440×900 /1280×720 /390×844。
- 验证：未主动执行新增/既有自动化（本轮未跑测试）；测试与文档更新已落盘。
- 后续：补齐 T35/T36 副作用一致性、T39 两仓部署与健康检查闭环，并补充 Atlas 关联记录。

## 2026-07-22 - M7 runbook 采集链路收口

- 变更：修正 `scripts/m7-runbook-evidence-collect.sh` 的 `expected_cursor` 编码为数字类型，并在 M7 runbook/evidence 中补充“单条 run 使用同一 IAOS token，不在每次动作中重复 `/dev/token`”约束。
- 原因：前期对账脚本在 `expected_cursor` 字段上存在 JSON 类型风险，且 token 轮换会引发 run 的 `token_mismatch`，影响 T36/恢复复查。
- 影响：T36 对账脚本可用性与可复现性提升，减少人为干预与误报；AESE 与 IAOS 当前联调链路状态仍以 IAOS/DB/Outbox 结果为阻塞项。
- 验证：本轮未新增自动化回归；本地已确认 AESE `:8090/ready` 与 IAOS `:8082/health、/ready` 可达。
- 后续：在 IAOS 联调窗口补齐 T35/T36（run context 与副作用一致性）与 T39（两仓健康与 commit）并更新证据状态。

## 2026-07-22 - M7 O4 运行恢复容错补齐

- 变更：在 `internal/httpapi/server.go` 修正 run 创建与恢复逻辑对 `ScenarioSnapshot` 缺失（404）的容忍：`run create` 支持场景缺快照时从 `cursor=0` 创建运行，`refreshRunFromFacts` 在快照 404 时保持本地状态并继续（不再失败）；新增 IAOS 404 判定辅助函数。
- 原因：当前 clean reset 后场景首次创建仍被 404 中断，导致 runbook 无法继续推进；问题源于 AESE 将快照缺失视为硬错误。
- 影响：当场景未写入或尚未产生初始 snapshot 时，执行 create/preflight/advance 等链路更容易恢复，不会因仅缺 snapshot 而阻塞。
- 验证：未进行自动化验证；该改动为可恢复路径提供前置修复，后续需结合 runbook 再跑一遍证据链。
- 后续：继续在本地服务稳定运行下重跑 `scripts/m7-runbook-evidence-collect.sh`，并按 T35/T36 对账闭环（UI/CLI、事件、Outbox、Tool Call）以及 T39 部署健康更新记录。

## 2026-07-21 - M7 O2 T21 文档与治理闭环补齐

- 变更：将 M7 plan 的 `T21` 标记为已完成，补齐 AESE-side 文档闭环：更新 `docs/plans/2026-07-20-m7-governed-scenario-operations-console.md`、`docs/roadmap.md`（说明 O2 当前状态与 IAOS 合并依赖）、`docs/code-map.md`（新增 M7 运行 runbook 映射）、`docs/README.md`（新增 M7 运行控制台 runbook 条目），并新建 `docs/runbooks/hctm-m7-governed-scenario-operations-console.md`（草案）。新增 `atlas-updates/2026-07-21-m7-o2-t21-doc-sync.json` 作为记录入口。
- 原因：满足 O2 末尾 T21 的“同步文档、runbook、code map 与状态治理”要求，保证 AESE 侧可追溯性。
- 影响：T21 的 AESE 文档同步点已形成，可让下一位执行者从统一入口理解 M7 O3/O4 目标与当前未完成边界；剩余 IAOS 权限/路由收敛与部署需在独立 IAOS worktree 补齐。
- 验证：未主动执行自动化测试；本次为文档与治理声明更新。
- 后续：待 IAOS worktree 同步 DES/运行端点更新并完成部署后，补一条 merge/deploy 的交付证据并将 O2 风险项降档。


## 2026-07-21 - M7 O2 运行恢复与元数据回归补齐

- 变更：在 `internal/httpapi/server_recovery_test.go` 修复 `TestRunActionAdvanceIsIdempotentUnderConcurrentCalls` 的幂等断言表达，避免将 `committed` 的真值与无关分支绑定；新增 `internal/replay/replay_test.go` 的 `TestReplayImpactIncludesRecoveryMetadataForGovernedActions`，覆盖治理路径在 apply 中返回 `cursor`、`operation_ref`、`correlation_id`、`no_op` 和 `committed` 的回归闭环。并同步 `docs/plans/2026-07-20-m7-governed-scenario-operations-console.md` 标记 T16/T18 完成。
- 原因：M7 O2 阶段要求 IAOS 状态恢复必须可按可重放字段重建，治理阶段返回的影响元信息需要结构化可审计。
- 影响：`T16/T18` 在本次迭代进入可验证完成状态；进度推进不影响 O3/O4 前端与三视口验收主线，但为其恢复和复现实用性提供了稳定接口条件。
- 验证：未执行新增/已存在 go test（本次改动按规则暂不主动运行测试）；按文件级别完成语义变更，并通过 `git diff --check` 检查 whitespace。
- 后续：继续推进 O2 剩余任务（尤其 T21 与运行端到端验收数据）及 O3 可视化控制台联调。

## 2026-07-21 - M7 O3 运行控制台联动与恢复入口收敛

- 变更：修复联动中心运行视图运行时错误并补齐关键交互。具体包括 `frontend/src/components/IntegrationConsole.tsx`：补齐 `PlayArrow/Radio` 图标导入、消除 `runAction/createRun` 闭包顺序隐患、修正 `restoreActiveRun` 依赖、补充创建运行后预检链路、恢复日志显示、控制台状态和动作按钮可用性。`frontend/src/integration/iaosIntegration.ts` 修正 `createScenarioRun` 请求体中未定义 token 引用。新增 `frontend/src/styles/global.css` 的 run 视图样式（mode 切换、stepper、动作面板、运行日志、复位卡片等）。并同步 `docs/plans/2026-07-20-m7-governed-scenario-operations-console.md`、`docs/roadmap.md`。
- 原因：上轮实现后联动控制台编译阻塞点仍存在，且 O3 任务在“可视化 run 控制台”尚未可用。
- 影响：前端从检查视图平滑切换到运行视图；用户可创建/恢复运行、预检、初始化、推进、跑完、分析、验证、复位，并可在浏览器重开后恢复 active run。O3 T22-T29 现有文档定义下完成进展，前端实现进入验收前整形阶段。
- 验证：未执行前端自动化验证（遵循“除非明确要求不额外运行验证”）；本轮为功能结构与状态治理更新。
- 后续：补齐 T30（移动端/可访问性/复制体验）与 O4（Playwright、网络中断/并发、部署与证据清单）并补充 Atlas 更新声明。

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

## 2026-07-19 - M5 受治理 Agent MVP 启动

- 变更：新增 DES-003 和 PLAN-M5-001，将 M5 设为唯一 active plan；并行审计计划、质量和经营分析 Agent 的现有规格、实时数据与 IAOS AI Tool Registry 能力。
- 原因：M4 已提供三类结构化异常和可查询业务状态，但 IAOS query tool 当前仅能列出工具，三个 Agent 尚不能通过受治理工具读取 HCTM 上下文。
- 影响：M5 采用通用 metadata 约束的 `entity.records` dispatcher，HCTM tool manifests 和 Agent 编排保留在 AESE；首版只读建议，不执行业务动作。经营分析必须显式报告尚缺完工入库、发运和实际成本事实，不能复用 Preview 答案冒充在线结果。
- 验证：三个独立审计均确认 AI Tool Registry 的 RLS、权限、schema、call audit 和 milestone event 可复用，并识别当前 `tenant-hctm` tool registry 为空及 shipment/cost 数据缺口。
- 后续：实现通用 query dispatcher、HCTM tool bundle 和三 Agent tracer，完成 live 调用、跨租户及业务零写验收。

## 2026-07-19 - 完成 M5 受治理 Agent MVP

- 变更：IAOS 增加 metadata 约束的通用 `entity.records` query dispatcher；AESE 增加 5 个 metadata schema、9 个 HCTM 只读 AI Tool manifest、`agent-setup`/`agent-run` 命令及计划、质量、经营分析三个 tracer。
- 原因：让 Agent 基于 `tenant-hctm` 的在线业务状态生成可解释建议，并复用 IAOS 的权限、RLS、Tool Registry 和调用审计，而不是读取 Preview 或绕过平台另造运行时。
- 影响：重复 setup 收敛为 9 个工具；重复 live run 不改变 24 条目标业务记录或 39 条 Outbox，只新增 9 条 Tool Call 和 36 条 milestone event；`tenant-other` 看不到 HCTM 工具。计划 Agent 会按当前库存报告 7,600 的物料缺口；质量和经营分析对缺失事实显式返回 `partial`，不虚构 1,700 合格放行、11,700 已发运或成本结果。
- 后续：M6 只消费稳定 recommendation envelope 与 IAOS 在线状态；补齐完工入库、发运和实际成本受治理事实后，再扩展经营分析的最终交付和利润结论。

## 2026-07-20 - M6 在线 2D 企业沙盘规划

- 变更：新增 DES-004 和 PLAN-M6-001，将 M6 设为当前唯一 active plan；计划拆为 L0-L5、T1-T37，并同步 README、MVP Blueprint、Agent Context、Architecture、Roadmap、Code Map 和文档索引。
- 原因：M3V 已有可用 2D Preview，M4/M5 已有在线异常与受治理 Agent tracer，但 IAOS 尚缺完工入库、两次发运、成本影响和可恢复场景事件合同，前端无法把 Preview 安全升级为 Live。
- 影响：M6 采用“快照为真、SSE 作增量提示、按持久 cursor 补齐”的架构。完工和发运走 IAOS 正式业务动作，不复用外生 simulation ingress；AESE 保留布局与视觉映射，通过 `IaosScenarioDataSource` 消费在线事实。成本金额无批准基线时继续标记 `partial`。
- 验证：计划定义 8-10 个工作日的每日成果、业务不变量、断线恢复、跨租户、Agent 证据和三视口验收；当前通用 SSE 已确认无持久 cursor、缓冲满可丢事件，不被误选为 M6 恢复合同。本地 Markdown 相对链接无缺失，active plan 数量为 1，M6 任务为 T1-T37，`git diff --check` 通过。
- 后续：从 L0/T1-T5 开始，冻结事件 17-22、成本完整度和 scenario observation API 合同，并在 IAOS 独立 worktree 建立 contract test。

## 2026-07-20 - 完成 M6 在线 2D 企业沙盘

- 变更：IAOS 新增事件 17-22 的受治理生产/完工入库/发运动作、库存 FIFO 扣减、持久场景游标、snapshot、SSE 和三 Agent 建议持久化；AESE replay/client/agent-run 接入正式合同，前端新增 Preview/Live、认证 HTTP、cursor 补发、SSE 去重重连、完整度和在线 KPI/建议展示。
- 原因：让 2D 沙盘只用 IAOS 业务事实闭合 12,000 件订单故事，并在断线、重复执行和跨租户条件下保持可恢复、可审计。
- 影响：M6 L0-L5/T1-T37 全部完成，当前没有 active plan。在线 KPI 为需求 12,000、累计可供/实发 11,700、期末成品 0、缺口 300；成本继续保留 `cost_actuals` partial gap，建议不自动执行。
- 后续：真实成本、更多场景或通用 projection 需要另立计划，不倒灌进已完成的 M6。

## 2026-07-20 - 修复远程浏览器 Live 回环地址

- 变更：`IaosScenarioDataSource` 默认从浏览器同源 `/api` 读取 IAOS，Vite dev/preview 将 `/api` 代理至本机 Platform 8082；新增回归测试，禁止默认配置重新请求浏览器侧 `127.0.0.1:8082`。
- 原因：远程用户虽然能访问 AESE 前端，但浏览器中的回环地址指向用户自己的机器，导致 Live snapshot `ERR_CONNECTION_REFUSED`；服务器端 Platform 实际健康且监听所有网卡。
- 影响：用户只需访问前端端口，snapshot、cursor 和 SSE 均走同源代理；显式 `VITE_IAOS_BASE_URL` 仍可覆盖默认配置。
- 后续：生产静态部署的反向代理同样需要把 `/api` 转发给 IAOS Platform；开发 token 仅用于本地测试。
## 2026-07-20 - 修复 IAOS 华辰租户可见性与开发工作区切换

- 变更：IAOS SaaS tenant lifecycle 现在把 `tenant_account` 原子投影到主界面使用的 `tenant` 目录，启动 bootstrap 幂等回填历史租户；新增受认证的 `dev-user` tenant token exchange，侧栏切换同时更新 JWT 和本地 tenant id。
- 原因：`tenant-hctm` 只存在于控制面目录，导致 IAOS 主界面租户下拉不可见；原侧栏只改 `localStorage`，不能改变后端强制执行的 JWT tenant claim。
- 影响：华辰租户可从 IAOS 平台工作区直接选择，切换后业务菜单、实体数据与 AESE Live 使用相同 `tenant-hctm` 边界；普通用户不能使用开发切换入口，租户状态仍按 active gate 检查。
- 后续：生产身份体系应使用真实跨租户 membership/SSO，不依赖 dev-user token exchange；AESE 本地演示继续使用该受限开发入口。

## 2026-07-20 - 修复 HCTM 业务菜单与订单明细可用性

- 变更：HCTM metadata bundle 为销售订单补充订单行 `child_list`，为客户、产品和订单引用补充目标实体；IAOS 数据浏览器在缺少独立 Formily UI Schema 时从实体字段生成可用详情表单，忽略请求去重产生的预期 AbortError，并隐藏当前租户不存在 Schema 的核心实体菜单。
- 原因：销售订单虽有头和行数据，但 `/metadata/ui/sales_order` 缺失导致详情抽屉永久停留在加载态，订单头又未声明行关系；全局 `inventory_lot` 菜单被错误投影到只使用 `inventory` 的 HCTM 租户。
- 影响：销售订单列表和明细可稳定查看，客户/产品引用显示业务标签；HCTM 不再展示会报错的“仓储物资”入口，库存继续通过“实物库存与库区”查看；工单、设备和库存页面沿用同一详情 fallback。
- 后续：其他行业包若需要定制表单布局可继续注册 `/metadata/ui/:entity`，不注册时使用字段驱动 fallback；生产环境应逐步将通用核心菜单改为完整的 capability/metadata 可用性投影。

## 2026-07-20 - 补齐 HCTM 客户引用元数据

- 变更：HCTM Agent setup bundle 新增 `customer` Entity Schema，使销售订单的 `customer_id` 引用可以通过 IAOS options API 解析为客户名称。
- 原因：销售订单已引用真实客户记录，但目标租户缺少 `customer` Schema，详情和列表加载客户选项时无法解析业务标签。
- 影响：重新应用 setup 后，销售订单客户字段可显示“星河新能源汽车”等业务名称，不再依赖裸 UUID。
- 后续：持续检查场景包所有 `reference` 字段都同时声明目标 Entity Schema。

## 2026-07-20 - 增加可视化 AESE × IAOS 联动中心

- 变更：AESE 顶栏新增“联动中心”，支持可视化选择 HCTM 租户和订单场景、一键取得本地演示身份、检查 IAOS profile/snapshot/event channel 与销售订单/工单/库存/设备记录，并通过对象映射直接跳转 IAOS 菜单或进入 Live；技术地址收纳在高级设置。
- 原因：原联调手册要求用户操作 Token、curl 和 CLI，不能满足业务用户直接观察配置与跨系统联动的需要。
- 影响：本地演示用户无需浏览器控制台或命令行即可完成只读联动验收；失败状态提供原因和恢复入口，检查过程不重置场景、不写业务数据。
- 后续：场景 reset/apply/replay 仍保持受治理写入边界；若要面向非开发身份开放可视化执行，需要为 AESE 增加有权限、可审计的服务端 orchestration API，而不是把写入逻辑复制到浏览器。

## 2026-07-20 - M7 受治理场景运行控制台规划

- 变更：新增 ADR-003、DES-005 和 PLAN-M7-001，将 M7 设为当前唯一 active plan；计划拆为 O0-O4、T1-T39，并同步 README、MVP Blueprint、Agent Context、Architecture、Roadmap、Code Map 和文档索引。
- 原因：M6 和扩展功能已解决 Live 观察、租户切换、业务对象可见性和一键联动检查，但业务用户仍需 CLI 执行 reset/apply/replay/agent-run。浏览器直接调用 IAOS 写端点会造成权限、恢复和部分执行风险。
- 影响：M7 新增无业务数据库的 AESE 薄编排 API，复用现有 Go 内核并使用调用者 IAOS 身份；页面增加 preflight、初始化、七幕推进、运行到结束、Agent 分析、verify 和一次性确认复位。所有业务事实、权限和审计继续由 IAOS 持有。
- 验证：计划定义 7-9 个工作日的每日成果、状态机、plan hash、cursor、幂等、并发、重启恢复、跨租户、CLI/UI 一致性和三视口验收；本地 Markdown 相对链接无缺失，active plan 数量为 1，M7 任务为 T1-T39，`git diff --check` 通过。
- 后续：从 O0/T1-T6 开始，先把 CLI 编排提取为 application service，并冻结 pack 阶段编译和 run 状态机合同。

## 2026-07-20 - 建立 IAOS + AESE System Atlas 双系统全景

- 变更：IAOS 新增平台级 System Atlas 节点、关系和追加式进展数据库，建立 32 个未来完成体构件基线、权限 API 和进展登记脚本；IAOS 工作台新增双系统动态图谱。AESE 新增聚焦虚拟企业模型、场景、Agent、2D 沙盘、实验和经营评估的动态图谱，并复用同一 IAOS 数据源。
- 原因：现有路线图和进展日志无法同时表达最终完整系统的组成、跨系统依赖、当前完成度和可追溯依据，用户与后续 Agent 难以形成一致全局认知。
- 影响：用户可缩放、筛选和点击构件查看目标、现状、完成度、设计/代码依据与更新历史；后续实质进展除项目日志外还必须登记到 Atlas。全景属于产品治理控制面，不存储或替代 HCTM 仿真业务事实。
- 验证：后端 package/API 单测、IAOS TypeScript 校验、AESE TypeScript 与 24 项前端测试通过；真实 PostgreSQL、生产构建和浏览器三视口验收待部署阶段完成。
- 后续：部署后完成真实 PostgreSQL seed/API 验收和两端浏览器截图；随后按 M7 实施进展持续更新 `aese.operations`，并逐步补充历史 commit 的精细证据。

## 2026-07-20 - 修复 Atlas 空引用下钻

- 变更：IAOS Atlas seed 将缺省设计、代码和证据引用规范化为空 JSON 数组；IAOS 与 AESE 详情面板同时兼容历史 `null` 数据。
- 原因：真实 PostgreSQL 首次 seed 对无引用节点保存 JSON `null`，浏览器点击这类节点时数组展开报错并中断详情渲染。
- 影响：所有节点均可稳定下钻；后端重启后静态引用字段收敛为空数组，前端仍可兼容修复前数据。
- 验证：待重新执行两端生产构建和桌面/移动浏览器节点点击验收。
- 后续：在 API DTO 层继续保持空集合而非 `null` 的响应合同，并为详情下钻增加组件回归测试。

## 2026-07-20 - System Atlas 双端发布与浏览器验收

- 变更：完成 IAOS Platform、IAOS Next.js 工作台和 AESE Vite 沙盘重部署；IAOS 移动端进入 Atlas 时自动收起侧栏与 Copilot，详情面板按实际内容区宽度显示。
- 原因：全景能力必须在真实数据库和实际服务上验证，移动工作台的既有展开面板会遮挡图谱交互。
- 影响：IAOS `3000` 与 AESE `4173` 均可访问动态全景；System Atlas API 在 `8082` 提供 32 个节点、37 条关系和追加式更新历史。
- 验证：真实 PostgreSQL schema/seed/API 读取与 update 写入通过；AESE 桌面/移动各渲染 16 个聚焦节点，IAOS 桌面/移动各渲染 32 个全景节点，四个视口节点下钻和当前状态显示通过且无浏览器异常。
- 后续：M7 每个交付切片完成后同时更新 `aese.operations` 进展；按实际设计逐步拆细业务域和实验评估子节点。

## 2026-07-20 - 调整 AESE System Atlas 入口位置

- 变更：移除页面右上角固定悬浮的“系统全景”入口，将其放入 Preview/Live 共用顶部控制栏，与“联动中心”和播放控制按正常布局排列。
- 原因：固定定位入口覆盖了顶部已有状态和操作信息，在部分桌面宽度和移动视口下影响阅读与点击。
- 影响：系统全景入口不再脱离布局覆盖内容；Preview 和 Live 保持同一入口位置，移动端沿用控制栏自动换行。
- 验证：TypeScript、lint、25 项 Vitest 和 Vite 生产构建通过；桌面与移动浏览器确认图标按钮、联动按钮和状态区边界互不重叠，页面无运行异常。
- 后续：Atlas 数据更新仍按 Agent 强制规则执行：实质进展必须更新进展日志并调用登记脚本；完成度不从 commit 数自动推断。

## 2026-07-20 - System Atlas 声明式自动治理

- 变更：两仓新增版本化 `atlas-updates` 声明、实质变更覆盖检查、GitHub Actions 守门和主分支同步脚本；IAOS API 增加唯一 `update_key`，重复同步返回既有记录而不重复更新节点。
- 原因：仅靠 Agent 规则和直接调用脚本无法阻止漏登记，也无法让 CI 在不能访问内网 IAOS 时验证进展是否已被描述。
- 影响：设计、实现、测试、发布、决策或风险变更必须随代码提交机器可读声明和证据；CI 负责拒绝遗漏，配置 Atlas endpoint secrets 的主分支环境负责自动入库。完成度仍由设计者明确判断，不按 commit 数推断。
- 验证：合法声明检查通过；模拟有实质变更但无声明的提交被拒绝；两仓声明连续同步两次后 5 个 `update_key` 在数据库中均严格只有 1 条，复用 key 修改内容返回 HTTP 409。
- 后续：在 GitHub `system-atlas` environment 配置 `IAOS_ATLAS_BASE_URL` 和 `IAOS_ATLAS_TOKEN`，使云端主分支同步作业连接正式 Atlas API。

## 2026-07-20 - Atlas 同步脚本工作目录加固

- 变更：IAOS 与 AESE 的 Atlas 校验/同步脚本改为从脚本自身路径解析仓库根目录。
- 原因：从 IAOS 目录用绝对路径调用 AESE 同步脚本时，旧实现按调用者当前目录查找声明，可能误同步 IAOS 文件。
- 影响：CI、仓库内调用和跨目录运维调用均读取正确仓库的 `atlas-updates`，避免跨仓声明混用。
- 验证：从 `/tmp` 分别以绝对路径调用两仓校验和同步脚本成功；当前 5 个声明 key 连续同步两次后各只入库一次。
- 后续：所有新增运维脚本默认采用脚本路径确定资源根目录，不依赖调用者 cwd。

## 2026-07-20 - System Atlas 可解释下钻与自由布局

- 变更：AESE 全景图新增 Dagre 自动布局、节点拖动、一跳关系高亮与关系方向列表；详情拆分设计文档、功能入口、代码位置和验证证据，并增加 Markdown 模态阅读器及 hash 深链接导航。
- 原因：原图无法清楚表达构件相关性，静态坐标会重叠且拖动不生效，文档和已实现功能也无法直接进入。
- 影响：用户可从 AESE 构件直接阅读登记文档、进入预览/实时沙盘或联动中心，并可手工调整图谱位置；文档内容仍由 IAOS 安全接口统一提供。
- 验证：`npm run build` 通过；Vitest 7 个测试文件、25 个测试通过；IAOS System Atlas 后端测试、`go vet` 和 Next.js 生产构建通过。
- 后续：部署双端后执行桌面与移动端视觉验收，并根据真实使用补充尚未登记入口的未来构件。

## 2026-07-20 - Atlas 窄内容区可读性加固

- 变更：自动布局命令改为带 tooltip 的固定尺寸图标按钮，图谱初始适配增加 0.4 最小缩放。
- 原因：真实桌面截图显示 IAOS 左侧导航和右侧 Copilot 同时打开时，布局按钮文字会换行，完整适配也会把节点压缩过小。
- 影响：工具栏不再因按钮文案产生重叠；全景保留小地图和平移能力，同时节点文字更容易辨认。
- 验证：IAOS Next.js 与 AESE Vite 生产构建通过，AESE 25 项 Vitest 通过。
- 后续：发布后复核桌面、移动视口和文档阅读器。

## 2026-07-21 - M7 O0 编排内核推进

- 变更：修复 `internal/application/plan_test.go` 的编译问题；`internal/application` 阶段/plan/hash/state 核心类型进入可编译状态；`docs/plans/2026-07-20-m7-governed-scenario-operations-console.md` 标记 O0 的 T2/T3/T5/T6 为完成；`docs/roadmap.md` O0 状态改为 In Progress；`docs/code-map.md` 的 M7 计划路径更新为含 `internal/application/`。
- 原因：保证 M7 第一阶段（状态机与编排合同）形成可复用内核并对齐文档与工程治理要求。
- 影响：下一步可在同一内核上实现 `aese-server` 与运行端点，无需重复定义 plan 与阶段语义；文档与代码地图保持一致。
- 验证：未执行全量单测；完成后未新增 IAOS 业务边界变更。待补充 Atlas 声明一致性检查结果与运行测试。
- 后续：完成 O0 的 application service 拆分（T4）与 HTTP handler 验证（T8/T9），并补齐 run recovery 与幂等恢复用例。

## 2026-07-21 - M7 O0 T4 已完成并完成编译验证

- 变更：将 M7 O0 的 T4 交付完成：CLI 写操作逻辑已稳定进入 `internal/application`，`cmd/aese/main.go` 回归为参数解析/校验/输出层；补齐 O0 执行计划 `correlation_id` 字段引用一致性；更新 O0 计划项 T4 为完成。
- 原因：继续推进下一阶段前先统一 CLI 与服务层边界，避免重复实现及状态机 drift。
- 影响：后续 `aese-server` 可直接复用 `internal/application` 作为 handler 执行核心；plan hash 与状态转换测试在单包范围内稳定通过。
- 验证：`go test ./cmd/aese ./internal/application` 通过；`docs/code-map.md`、`docs/roadmap.md`、`docs/plans/...` 已同步更新；补充 Atlas 声明待校验脚本同步。
- 后续：运行 `scripts/check_system_atlas_tracking.sh` 补齐实质变更声明治理闭环，并继续推进 O0/T6 幂等与恢复测试。

## 2026-07-21 - M7 O1 API 骨架修复与补齐

- 变更：修复了薄编排 API 运行时的两个阻断性缺陷（`schema.Required` 字段引用和 Analyze `pack` 为空导致 `scenario pack is required`），并确认 `aese-server` 启动入口与 M7 API 路由已可复用。
- 原因：防止 M7 O1 阶段在预检/分析阶段出现 500 或 `RunAgents` 运行阻塞，确保前端首轮调用可持续落地。
- 影响：`internal/httpapi/server.go` 预检返回中加入字段级合同，`analyze` 使用同一 pack 上下文；`cmd/aese-server/main.go` 提供独立启动入口，`docs/code-map.md` 增补 API 入口映射。
- 验证：未执行全量回归与集成；修改为编译可达路径修正，未触发 IAOS 侧 schema 字段变更；计划后续对幂等、恢复和权限恢复用例补齐。
- 后续：继续推进 O1 中 run/action 合同、恢复策略与幂等重放；完成 Atlas 声明并走一致性检查脚本。

## 2026-07-21 - M7 O1 运行恢复与幂等补齐

- 变更：在 `internal/httpapi/server.go` 完善 run 恢复逻辑，`refreshRunFromFacts` 改为按 `run.Cursor` 增量重建状态；新增 `refreshRunCursor` 初始化游标基线；补齐 snapshot 事件的 map-to-struct 解析与 cursor 更新规则；`run/version/plan hash/expected cursor/idempotency` 与一次性 reset token 分支保护行为补齐。
- 原因：消除进程重启或重复调用下把历史事件误算入新 run 的风险，降低 reset 重放与恢复误判概率。
- 影响：`docs/plans/2026-07-20-m7-governed-scenario-operations-console.md` 标记 T12/T13/T14 为完成，后续可继续推进 run 并发、权限和前端控制流实现。
- 验证：未执行全量回归与集成，`go test` 未在本轮触发；本轮变更为运行期逻辑修订，需后续在重启/重复点击场景进行契约验证。
- 后续：继续推进 O1 的并发冲突、权限回归和跨服务断线恢复用例；补齐 `atlas-updates` 声明与审计链路。

## 2026-07-21 - M7 O1 同源访问边界补齐

- 变更：在 [internal/httpapi/server.go](/iaos/aese/internal/httpapi/server.go) 增加基础 CORS 与 OPTIONS 预检处理；`refreshRunFromFacts` 的恢复过滤收紧为仅处理与当前场景 run correlation 一致的事件，减少恢复时历史污染风险。
- 原因：前端浏览器访问 thin orchestration API 需要稳定跨域行为，同时恢复逻辑应对齐 run 运行边界。
- 影响：跨域预检返回变为可用；`run` 恢复在 correlation 粗过滤下更严格。
- 验证：未执行回归与 E2E；本项作为 T15 部分收敛，仍需配合同源代理与敏感字段脱敏策略完成。
- 后续：补齐同源代理（前端开发代理或网关层）与脱敏字段，更新 T15 完成状态。

## 2026-07-21 - M7 O1 敏感字段脱敏与错误消息防泄露

- 变更：在前端 IAOS 数据源与连接检查逻辑中加入错误消息脱敏；`Bearer <token>` 将统一输出为 `Bearer [REDACTED]`，并补充单测覆盖。
- 原因：将异常响应中的调用者 token 风险从 UI 回显与调试日志面板剥离。
- 影响：`frontend/src/scenario/iaosDataSource.ts` 与 `frontend/src/integration/iaosIntegration.ts` 增加统一脱敏函数；`frontend/src/scenario/iaosDataSource.test.ts` 与 `frontend/src/integration/iaosIntegration.test.ts` 增加回归用例。
- 验证：未执行全量回归与集成；变更局限于错误消息处理层。
- 后续：对同源代理路径进行收敛（建议统一前端走 `/api/iaos/v1`），然后将 T15 标记为已完成并补齐重启恢复验证。

## 2026-07-21 - M7 O1 T15 脱敏补丁入 System Atlas

- 变更：新增 `atlas-updates/2026-07-21-m7-o1-sensitive-redaction.json`，记录前端异常消息脱敏实现与 T15 进度，便于 System Atlas 追踪。
- 原因：每次实质性变更需留痕声明并可回溯证据。
- 影响：T15 的“敏感字段脱敏”项完成；T15 的“同源代理与恢复测试”仍待补齐。
- 验证：未执行新增声明校验。
- 后续：完成同源代理收敛与服务重启恢复测试后再提交一条 T15 闭环声明。

## 2026-07-21 - M7 O1 服务重启恢复回归补齐

- 变更：新增 `internal/httpapi/server_recovery_test.go`，用可控 IAOS 假服务覆盖 `refreshRunFromFacts` 两次恢复过程，验证基于 cursor 的增量恢复不会重复计算历史事件并可推进完成幕位。
- 原因：T15 的“服务重启恢复测试”要求可复测化。
- 影响：`T15` 逐步收敛，`restart`/重复抓取不再依赖进程内存重建新状态。
- 验证：新增单测未执行（按当前流程暂不跑测试）。
- 后续：新增一条 restart 恢复 System Atlas 更新声明并持续观察真实场景中 AEOS 重启回放行为一致性。

## 2026-07-21 - M7 O1 权限失败关闭补齐

- 变更：在 [internal/httpapi/server.go](/iaos/aese/internal/httpapi/server.go) 完成 T11 其二，`run create` 与 `run status` 现在会把 IAOS 返回的 401/403/409 等错误直接传播为状态码，而不是将快照读取失败静默吞掉；新增两项单测分别验证 `snapshot` 权限拒绝时 create/status 直接返回 403。
- 原因：当前执行控制台仍处于 T11，必须先建立写入前置权限失败关闭，避免用户误以为状态正常。
- 影响：业务运行从“可见但不可执行”切到“可见且可解释失败”；`/api/aese/v1/runs/:run_id` 受到身份/租户/权限失效时会 fail-closed。
- 验证：未执行全量回归；新增测试为逻辑路径提供覆盖。
- 后续：继续完成 O1 的动作并发保护测试与 run 冲突契约补充。

## 2026-07-21 - M7 O2 路径启动与 run 并发约束验证

- 变更：在 [internal/httpapi/server_recovery_test.go](/iaos/aese/internal/httpapi/server_recovery_test.go) 补充并发可写 run 约束回归：验证同租户同场景下仍有可写 run 时拒绝新建（409 conflict），以及已完成 run 后可继续创建新 run（201 created）。同步将 M7 roadmap O2 阶段置为 In Progress。
- 原因：O2 的 `T19` 需要有明确可复现依据，防止前端页面并发点击/重复恢复导致写入冲突。
- 影响：AESE 端薄编排 API 的并发约束可量化回归，避免无感知重建或重复创建 run。
- 验证：未执行全量回归；变更为新增单测。
- 后续：继续补齐 T16/T17 的 IAOS 侧权限行为与运行状态审计对齐。

## 2026-07-21 - M7 O2 补充 stale cursor 防护回归

- 变更：新增 `internal/httpapi/server_recovery_test.go` 稳定性用例：当 `expected_cursor` 过期时，`run action` 返回 `409 cursor_mismatch` 并不继续执行动作前置验证。
- 原因：T20 需要覆盖陈旧游标并发场景，避免用户点击过期操作导致错误推进。
- 影响：`run action` 的前置保护顺序可复测化，便于前端在乐观并发冲突时给出可恢复提示。
- 验证：未执行全量回归；新增测试为逻辑路径提供覆盖。
- 后续：补齐跨租户/权限不足/重放冲突路径后，继续在计划中将 O2 `T20` 标记进行中或完成子项，并开始 T16/T17 的 IAOS 侧审计。

## 2026-07-21 - M7 O2 补充跨租户动作前置回归

- 变更：新增 `internal/httpapi/server_recovery_test.go` 用例：`run status` 在 IAOS profile 回传非同租户时返回 `tenant_mismatch`，验证 `loadProfileForRun` 的租户一致性约束。
- 原因：T20 覆盖跨租户权限不足路径，防止把 tenant-other 的 token 误用于 tenant-hctm run。
- 影响：权限前置失败路径从“运行状态返回 200”变成可预测 403，便于调用方在前端显示清晰失败原因。
- 验证：未执行全量回归；新增测试为逻辑路径提供覆盖。
- 后续：继续补齐 IAOS 侧 `scenario.run.*` 权限枚举与 run 状态契约。

## 2026-07-21 - M7 O1/O2 运行动作幂等与 reset 状态闭环

- 变更：修复 `internal/httpapi/server.go` 中 `handleRunAction` 的幂等缓存路径：移除重复 `cacheKey` 声明；补齐空 idempotency 键直接透传与缓存写入门控；新增 `actionRequiresIdempotency`，并要求 `initialize/advance/run-to-end/analyze/verify/reset` 在 `apply=true` 时必须有 `Idempotency-Key`；`run action` 的错误与成功结果仅在有 idempotency 键时缓存。
- 变更：完善 `reset` 合同：`reset-plan` 走 `application.NextStatus` 进入 `resetting`，返回确认 token；`reset` 执行仅在 `resetting` 状态允许，`apply=true` 才清 token、置 `run.ResetTokenExpiresAt` 为零并进入 `reset`，失败场景保留 token 以允许重试。
- 变更：同步状态机约束 `internal/application/state.go`，允许 `run` 在 `resetting` 时执行 `reset`；并让 `inferRunStatusFromFacts` 对 `reset` 状态进行保留。
- 原因：M7 API 仍缺少空幂等键防呗、重复动作去重边界和 reset 过渡闭环，可能造成重复动作确认、状态错判与误清 token。
- 影响：`run action` 的去重语义更稳定，`reset` 失败不会清 token 丢失一次性确认上下文，状态恢复时不会将 `reset` 回推为其他阶段。
- 验证：未执行全量回归；本次修改为代码级闭环，建议配套继续补齐 `server_recovery_test.go` 的 run/action 幂等和 reset 用例。
- 后续：补齐 T20（权限不足与 reset 冲突/重放）与 O2 前后端联动，更新 frontend run 状态机与控制台恢复流程。

## 2026-07-21 - M7 O2 权限资源提示闭环

- 变更：在 [internal/httpapi/server_recovery_test.go](/iaos/aese/internal/httpapi/server_recovery_test.go) 补充 `required_permission` 回归：`run create` 与 `run status` 权限不足路径返回 `scenario.run.read`；`initialize`/`reset` 返回 `scenario.run.execute` 与 `scenario.run.reset`，并同步补齐 O2 的 T17、T20 里权限与并发/恢复回归项。
- 原因：前端与调用方需要可靠的权限资源提示，才能在 403/tenant 冲突下给出可恢复动作并满足 O2 权限闭环验收标准。
- 影响：错误合约从“状态码+message”上升到“可执行权限资源提示”；同一 API 在读写/复位场景下的失败行为可统一呈现给 UI。
- 验证：未执行自动化；新增回归覆盖包含 `error.required_permission` 字段断言。
- 后续：完成 T16/T21 IAOS 状态恢复字段与 cross-worktree 同步部署，再继续推进 O3 运行控制台联动视图任务。
## 2026-07-21 - M7 证据脚本路径修正与联调中断点固化

- 变更：修复 `scripts/m7-runbook-evidence-collect.sh` 的 IAOS 地址变量语义，明确区分 `IAOS_API_ROOT`、`IAOS_API_BASE` 与 `IAOS_TOKEN_BASE`，并将 token 改为通过 `.../api/v1/dev/token` 获取；IAOS `ready` 证据改为可选采集并输出不可用占位对象，避免 /ready 404 直接中断。
- 原因：当前 M7 收口脚本与 IAOS 路径约定不一致，出现 404 404，导致 T36 自动对账停住，阻塞判断被误认为“服务不可连”。
- 影响：证据链脚本可在 IAOS 健康可达前提下顺畅产出 00-08 与 CLI 工件（UI 部分仍依赖 AESE/IAOS 实际端点行为）。
- 验证：未执行自动化回放；变更为静态路径修订与失败可回退行为。
- 后续：在服务启动后按 DRY-RUN 先跑一遍 `scripts/m7-runbook-evidence-collect.sh`，确认 run/verify/reset 证据文件；将结果直接用于 T35/T36 及 T39 的联调验收。
## 2026-07-21 - M7 证据脚本再修与 409 诊断

- 变更：`scripts/m7-runbook-evidence-collect.sh` 将 `run create` 的 `target` 固定为 IAOS 根地址（`IAOS_API_ROOT`），避免与 `IAOS_API_BASE` 混淆导致的 `/api/v1/api/v1/profile` 404。
- 原因：实际执行检查时仍出现 409，需准确区分“地址拼接问题”与“活跃写入 run 冲突”。
- 影响：`scripts/m7-runbook-evidence-collect.sh` 在地址配置不一致时更稳健；当前 409 为 AESE 约束冲突（同租户同场景已有可写 run）。
- 验证：通过 curl 直接命中 AESE create 接口，返回 `another writable run exists for tenant/story tenant-hctm/hctm/order-expedite-01`。
- 后续：重启 AESE 服务或清理该租户当前场景活动 run 后重新执行；再继续推进 T35/T36。

## 2026-07-22 - M7 O4 快照缺失与 preflight 受阻补齐

- 变更：在 `internal/httpapi/server.go` 对 `ScenarioSnapshot` 的 `404` 做容错后，手工复测 `run create` 已能从 `cursor=0` 成功返回 201。
- 原因：此前清理 run 后 AESE 端仍误返回 `not_found`；修复后确认该层已清。复测继续发现 `preflight` 在当前 IAOS 环境报 `metadata/schema/inventory_transaction` 404。
- 影响：AESE 重连恢复路径已通过；当前受阻点是 IAOS 元数据/seed 完整性，`preflight` 无法继续进入 initialize。
- 验证：在干净进程中手工调用 `run create` 与 `preflight`；create 返回 201，preflight 返回 404（`required_permission: scenario.run.execute`）。
- 后续：补齐 IAOS `inventory_transaction` schema（或在 AESE 明确落化缺 schema 的降级行为）后再继续完整 runbook 跑通并进入 T35/T36/T39。

## 2026-07-22 - M7 受治理场景运行控制台完成

- 变更：修正 O4 证据脚本的数值 cursor、固定 JWT、reset token 路径、CLI apply 前置和异步 verify 有界重试；修复服务重启后按增量 cursor 从已有幕位继续恢复的算法；从 clean reset 完成 `m7-acceptance-20260722-05` 编排 API 与 CLI 对照链，并将 DES-005、PLAN-M7-001、Roadmap O0-O4、runbook/evidence 和项目入口统一转为 Completed。
- 原因：此前阻塞并非 IAOS 不可达，而是历史非 clean 销售订单触发 preflight 409、进程内僵尸 active run 无法寻址，以及证据脚本误读 reset token 并立即执行异步 CLI verify。
- 影响：业务用户可在浏览器受治理完成 22 事件、三 Agent、验证与一次性确认复位；最终单 run 产生 9 次成功 Tool Call、18 条 UI/CLI 对称 Outbox，两个执行面均删除 15 个 L2 对象并保留 12 个 L1 对象。当前没有 active plan。
- 验证：17 条离线业务断言与 2 条 IAOS 在线断言通过；M6 KPI 为需求 12,000、实发 11,700、缺口 300；`go test ./...`、`go vet ./...`、29 项 Vitest、24 项三视口 Playwright、ESLint、TypeScript/Vite production build、Atlas tracking、JSON、Markdown 链接和 diff check 通过。AESE 8090/4173 与 IAOS 8082/3000 健康可达，PostgreSQL/NATS 正常。完整产物位于 `artifacts/m7-acceptance/20260722-05/`。
- 后续：M8 参数化仿真实验必须另立 active plan；生产身份应替换本地 dev token，成本维度在批准基线前继续保持 `partial`。

## 2026-07-22 - AESE 2.0 改造计划与 M8 架构决策门建立

- 变更：评审 `docs/ChatGPT20260722-aese2.0.md`，新增 ADR-004、DES-007 和唯一 active 的 PLAN-M8-001；将企业生命周期方向工程化为 World / IAOS / Actor Knowledge 三态、Go 离散事件内核、受治理双向桥和 Project Genesis 路线，并同步 README、Agent Context、Architecture、文档索引、Roadmap 与 Code Map。
- 原因：原始构思正确指出现有 AESE 仍以预编排业务故事为核心，但其“AESE 负责客观世界”与现有无状态场景层边界、Spring Boot 示例及原 M8 A/B 实验优先级存在冲突，需要先建立可批准、可分期和可回归的系统改造计划。
- 影响：M8 从参数化实验调整为 AESE 2.0 基础里程碑，当前只进入规划与决策门；ADR-004 获批前不实施生产级 World Store 或 IAOS 持久化变更。M7 场景包、CLI、Preview/Live 和控制台继续作为兼容基线。
- 验证：`git diff --check`、active plan 唯一性、System Atlas tracking、全部场景/Atlas JSON 解析和受影响 Markdown 本地链接检查通过；未修改代码或 IAOS 仓库，因此未运行产品测试。
- 后续：与用户确认 ADR-004 的状态所有权、World Store、Actor Knowledge 和最小经济守恒四项决策，再执行 F0 合同与 IAOS gap audit。

## 2026-07-22 - ADR-004 三态所有权与 World Store 决策获批

- 变更：将 ADR-004 转为 accepted，确认 AESE 拥有可持久化仿真事实、World Store 使用独立 PostgreSQL database/账号/迁移边界、Actor Knowledge 首版只保存结构化认知，并将经济守恒限定为现金、承诺支出和应收回款；同步 PLAN-M8-001 的 G1/G2/G3/G5、Architecture、Roadmap、Agent Context 和文档索引。
- 原因：用户接受三态所有权、独立存储和最小经济守恒；对 Actor Knowledge 的首版边界采用可确定性重放、审计和权限控制的结构化方案，避免自由文本长期记忆成为权威状态。
- 影响：M8 可进入 AESE 本地 World Store 和合同原型工作；在 G4 observation/intent/committed_outcome 合同冻结前，仍不得修改 IAOS 主分支或部署生产双向桥。
- 验证：`git diff --check`、决策状态、active plan 唯一性、全部场景/Atlas JSON 解析、受影响 Markdown 本地链接和 System Atlas tracking 均通过；未修改产品代码或 IAOS 仓库。
- 后续：执行 F0 T3-T6，优先完成三态所有权矩阵、PostgreSQL 运行合同、bridge envelope 和 IAOS gap audit。

## 2026-07-22 - G4 World/IAOS 三段式桥接合同获批

- 变更：新增 approved DES-008，固定 observation、intent、committed_outcome 的公共 envelope、strict payload、权限、幂等、错误和跨仓顺序；选择 IAOS tenant journal + cursor query 作为恢复事实，SSE/Outbox 只用于通知。同步 PLAN-M8-001 完成 G4/T5/T6，并更新 README、Architecture、Agent Context、文档索引、Roadmap 与 Code Map。
- 原因：用户授权由架构设计决定 G4。现有 IAOS simulation ingress 会同时改变业务状态、scenario event 仅支持固定 HCTM story，但已有 tenant、RLS、input hash、committed record refs、Outbox 和 cursor 模式可复用；新合同需保持这些优点并分离“看见事实”“提出行动”“事务已提交”。
- 影响：M8 G1-G5 全部通过，可进入 F0 机器可读 Schema/Go 类型和 contract tests。AESE 只根据 committed/no-op outcome 计算世界后果，不根据 intent、HTTP 超时、回滚结果、webhook 或 direct NATS 猜测状态。
- 验证：对照 AESE `internal/iaosclient` 与 IAOS scenario、simulation、Capability execution result、eventdef 当前合同完成只读 gap audit；`git diff --check`、文档 ID、active plan 唯一性、全部场景/Atlas JSON、受影响 Markdown 链接和 System Atlas tracking 均通过。未修改 IAOS 仓库。
- 后续：执行 T3/T4/T7，提交 bridge JSON Schema、fixture、Go 类型、canonical hash 与 mock journal contract tests，再进入 IAOS 独立 worktree。

## 2026-07-22 - M8 F0 机器合同与 World Store 边界完成

- 变更：完成 PLAN-M8-001 T3/T4/T7；新增三态术语、对象所有权和数据分类，八类 World/Bridge JSON Schema、Go strict parser、canonical SHA-256、fixture/测试，独立 PostgreSQL database/账号/迁移/连接/本地启动/备份合同，以及五个 planned System Atlas 节点和依赖声明。
- 原因：在进入确定性内核和 IAOS 独立 worktree 前冻结稳定的状态、事件、认知、差异和桥接语义，避免后续 reducer 与双向桥发生 Schema 漂移。
- 影响：M8 F0 转为 Completed，F1 成为下一切片；未实现完整仿真内核，未修改 IAOS、M7 pack、CLI 或 Preview/Live 行为，也未引入 direct NATS/webhook-only 路径。
- 验证：`go test ./...`、`go vet ./...`、八组 fixture 的 JSON Schema Draft 2020-12 校验、全部新增 JSON 解析、Compose 配置、临时 PostgreSQL 17 的实际 up migration 与 RLS 未绑定拒绝/绑定可写测试、`git diff --check`、Markdown 本地链接和 System Atlas tracking 均通过；临时容器与卷已清理。
- 后续：按 F1 T8-T12 实现纯函数 reducer、虚拟时钟、稳定事件排序、快照/重放和默认 dry-run 命令；F1 前不启动 F2/F3 实现。

## 2026-07-22 - M8 F1 确定性离散事件内核完成

- 变更：完成 PLAN-M8-001 T8-T12；新增虚拟时钟、稳定优先队列、版本化纯函数 reducer、World engine、event log/state hash、snapshot 恢复和严格 replay，并在现有 CLI 增加默认 dry-run 的 `aese world validate|inspect|run|replay`；仅显式 `--apply --output` 写 artifact。
- 原因：先建立可证明重放一致性的最小内核，再扩展设备、人员和 Actor Knowledge 业务模型，避免业务 tracer 掩盖时序与确定性错误。
- 影响：F1 转为 Completed，F2 成为下一切片；当前规则仅包含通用 `state.set.v1` 测试 reducer，不宣称 Genesis 业务规则或 KPI 完成。未连接 PostgreSQL/IAOS，未修改 M7 pack、旧 CLI 命令或 Preview/Live。
- 验证：同一输入 100 次日志 hash 一致；覆盖秒/年时间尺度、相同时间 priority/event ID tie-break、重复 ID/幂等键、因果倒序、未知 rules/payload、时间倒退、损坏快照和恢复 sequence；CLI dry-run/apply/replay、全仓 Go test/vet、JSON/链接/Atlas tracking 与 diff check 均通过。
- 后续：执行 F2 T13-T17，围绕 `LAS-WLD-02` 建立设备实际状态、角色有限认知、只读 IAOS 投影和 discrepancy 生命周期，不提前实现 IAOS 写桥。

## 2026-07-22 - M8 AESE 2.0 foundation 完成

- 变更：完成 PLAN-M8-001 F2-F5；交付 `LAS-WLD-02` World/IAOS/Knowledge 三态 tracer、actor-scoped Knowledge、Genesis world pack、M7 22 事件兼容 adapter、受治理 IAOS observation/intent/committed outcome journal/cursor/SSE、World Play 三态界面和能力缺口台账。IAOS 改动保持在独立 `feat/m8-world-bridge` worktree，revision 为 `e661d9a`。
- 原因：在 F0 合同与 F1 确定性内核冻结后一次完成剩余闭环，同时保持 AESE 不直写 IAOS、通知不替代 cursor 事实和 M7 路径零改写。
- 影响：非研发用户可通过 World 模式观察偏差、推进/复位虚拟时间并查看发现到关闭的因果链；IAOS journal 强制 tenant RLS、权限、幂等和 journal+Outbox 原子提交。人类与 Agent 共用同一 intent/outcome 权限合同，不存在旁路。
- 验证：AESE `go test ./...`、Genesis validate/dry-run、30 项 Vitest、TypeScript/build、World Play 三视口 Playwright 通过；IAOS 全部四个 Go module 测试、bridge SQL mock、Atlas/Code Map 检查通过，后端已从独立 worktree 部署至 8082，健康检查通过，实际 PostgreSQL 表确认 `ENABLE/FORCE RLS` 与 tenant policy。M7 代码路径未修改。
- 后续：M8 已完成；通用 EAM 编排和 Project Genesis 更长生命周期属于后续新里程碑，需另立唯一 active plan。

## 2026-07-22 - M8 Atlas 同步待补登记

- 变更：在 IAOS 8082 健康且 M8 bridge 已部署后执行 `scripts/sync_system_atlas_updates.sh`。
- 原因：按治理要求将六个 M8 声明幂等登记到 System Atlas。
- 影响：同步端点返回 `404 system atlas node not found`；本地声明、证据与 tracking 校验均完整，产品代码和 M8 运行能力不受影响。
- 验证：dev token 获取成功，`POST /api/v1/system-atlas/updates` 可达并明确返回缺少 `aese.world` 节点，而非网络或鉴权失败。
- 后续：IAOS Atlas seed 注册 `atlas/system-atlas-planned.json` 的五个节点后，重跑同步脚本补登记；不得用历史补录脚本绕过缺失节点。

## 2026-07-22 - M8 Atlas 同步补登记完成

- 变更：IAOS seed 注册 `aese.world` 五节点族与 `iaos.world-bridge`，重新部署 8082，并通过受治理 update API 同步 AESE 31 条及 IAOS 11 条声明；同步脚本改为按 `occurred_at` 排序。
- 原因：解除先前 `system atlas node not found`，并防止首次批量同步按文件名应用导致状态回退。
- 影响：M8 World、Time、Knowledge、Genesis、AESE/IAOS Bridge 均在 Atlas 可见；历史声明保持不可变，最终状态使用新 update key 校正。
- 验证：数据库查询确认六个节点均为 `completed/100`，M8 六条原始 update key、AESE 校正记录与 IAOS bridge 完成记录均存在；两仓同步脚本重复执行成功。
- 后续：无。

## 2026-07-22 - M9 Project Genesis 企业成立与治理完成

- 变更：完成 PLAN-M9-001 G4/G5 与 I0-I5/T1-T35；新增 `hctm-genesis@0.2.0` incorporation campaign、8 阶段纯函数 tracer、独立 owner 现金账户、登记/银行策略、三岗位与 Knowledge、预算/mandate/资格不变量、snapshot/restore/reset、统一 human/Agent 授权，以及 World Play 成立视图。IAOS 独立 worktree revision `edcb915` 新增五类 allowlist 治理动作、四个权限、FORCE RLS、幂等与事务 Outbox。
- 原因：把 Genesis 从“企业已经存在”向前推进到法人、资金、管理岗位与预算均具备的机器可验证起点，为 M10 工厂建设立项提供真实前置条件。
- 影响：M9 终态输出 `plant_project_eligible=true`；预算授权与现金、认缴与实缴严格分离。外部登记和银行结果仍由 AESE World 策略产生，IAOS 不能伪造；M7/M8 路径保留。
- 验证：100 次 campaign hash、资金/岗位/mandate/snapshot 失败关闭、全量 AESE Go test/vet、IAOS 四 module test、Schema/JSON、32+ Vitest、生产构建、M7/M8/M9 Playwright 三视口通过。IAOS 真实 API 新建返回 201、重复返回 duplicate，预算自批返回 422；数据库确认 FORCE RLS 与 Outbox。两端 8090/8082/4173 已部署。
- 后续：M10 工厂选址与建设立项应另立唯一 active plan，只消费 M9 机器资格，不从 UI 状态推断。

## 2026-07-22 - Atlas 401 自动恢复

- 变更：System Atlas 的 IAOS 请求在缓存 token 返回 401 时自动清除旧身份、获取新 dev token 并重试一次；同一恢复逻辑覆盖图谱与文档下钻。
- 原因：浏览器 `localStorage.iaos_token` 可能来自旧会话；原实现只判断 token 是否存在，导致 `/#atlas` 永久停在 Unauthorized。
- 影响：用户无需手动清理 localStorage 或先进入联动中心；非 401 错误不重试，避免隐藏权限与服务错误。
- 验证：新增 Vitest 回归覆盖 `stale -> 401 -> refresh -> 200`；生产构建和类型检查通过；在实际 4173/8082 代理链注入失效 token，桌面 1440/1280 与移动 390 三项 Playwright 全部通过。
- 后续：无。

## 2026-07-22 - M9 企业成立与治理计划启动

- 变更：新增 approved DES-010 和唯一 active 的 PLAN-M9-001，将 Project Genesis 下一阶段拆为 I0-I5：业务/机器合同、AESE 成立世界与经济规则、IAOS 法人与治理能力、CEO/CFO 统一岗位、Incorporation Play 和全链验收；同步 README、Agent Context、Architecture、文档索引、Roadmap 与 Code Map，并修正 M8 在文档索引中的陈旧 Active 状态。
- 原因：M8 已完成三态世界和桥接基础，但当前 Genesis 从已有现金、人员和设备开始，仍不能证明企业如何合法成立、获得实际资本、建立管理授权并形成工厂建设资格。
- 影响：M9 成为当前唯一 active plan；目标终态固定为 `plant_project_eligible=true`。M9 不进入工厂选址/建设、设备采购、APQP 或完整财务，且复用 M8 World Runtime/Bridge/Play，不建设第二套引擎。
- 验证：`git diff --check`、DES/PLAN ID 唯一性、active plan 唯一性、全部 Atlas/World/场景 JSON 解析、受影响 Markdown 本地链接和 System Atlas tracking 均通过；本轮未修改产品代码，未运行产品测试。工作区既有测试修改和截图删除均未触碰。
- 后续：执行 I0 T1-T6，先冻结资金/周期/费用、stable code、状态机、payload schema 和 IAOS gap，再允许独立 IAOS worktree 开发。

## 2026-07-22 - World Play null Knowledge 白屏修复

- 变更：Genesis 首帧 Knowledge 从 Go nil slice 改为非 nil 空集合，确保 JSON 输出 `[]`；前端 World API 边界兼容归一化历史 `null`，并新增真实服务 E2E。
- 原因：`/#world` 首帧角色尚未知时 API 返回 `knowledge:null`，组件读取 `.length` 触发运行时异常并白屏。
- 影响：角色未知仍以明确空态呈现，不授予额外 World State；新旧 API 响应均不会导致页面崩溃。
- 验证：修复前 Go 与 Vitest 回归均稳定失败；修复后全量 Go test/vet、32 项 Vitest、TypeScript、production build 通过。8090 API 实际返回 array，4173 的桌面 1440/1280 与移动 390 三项 live Playwright 均无 page error。
- 后续：无。

## 2026-07-22 - M10 工厂选址与设施建设计划启动

- 变更：新增 approved DES-011 和唯一 active 的 PLAN-M10-001，将 M10 拆为 P0-P5：机器合同、受约束选址、AESE 设施世界、IAOS 投资与项目治理、统一角色与 Plant Build Play、全链验收；同步 README、Agent Context、Architecture、文档索引、Roadmap、Code Map 和 System Atlas 声明。
- 原因：M9 已输出 `plant_project_eligible=true`，下一步需要在真实现金、预算、工期、空间、公用工程和外部资源约束下取得场地控制并完成设施验收，为 M11 建设生产能力提供机器资格。
- 影响：M10 成为当前唯一 active plan，目标终态固定为 `capability_build_eligible=true`。首版比较绿地自建、租赁标准厂房改造和定制代建，但不预设评分赢家；生产设备、检测仪器、人员和投产门仍属于 M11。
- 验证：`git diff --check`、DES/PLAN ID 唯一性、active plan 唯一性、全部 Atlas/World/场景 JSON 解析、受影响 Markdown 本地链接和 `scripts/check_system_atlas_tracking.sh` 均通过；本轮只修改设计与计划文档，未运行产品测试。工作区既有测试修改、截图删除和验收产物均未触碰。
- 后续：执行 P0 T1-T6，先关闭 G4 候选/预算/工期基线与 G5 IAOS gap audit，再允许独立 IAOS worktree 开发。

## 2026-07-22 - M10 Genesis 工厂选址与设施建设完成

- 变更：完成 PLAN-M10-001 G4/G5 与 P0-P5/T1-T36；交付三个虚构候选的硬约束和解释评分、10 帧确定性设施项目世界、七节点空间模型、WBS/公用工程/资金守恒、utility delay Knowledge/rebaseline tracer、Plant Build API 与三态 UI。IAOS 独立 worktree revision `23be02a` 交付六类投资/项目/付款治理权限，业务记录、committed outcome journal 与 Outbox 同事务。
- 原因：将 M9 的 `plant_project_eligible` 推进为经过场址、预算、现金、工程、公用工程、消防/EHS、验收和治理门的真实设施载体，为 M11 提供机器可验证入口。
- 影响：`hctm-genesis` 升级为 0.3.0，终态输出 `capability_build_eligible=true`；生产设备、检测仪器、招聘培训与投产仍未实现。M7/M8/M9 的 pack、CLI 和 Preview/Live 行为未改写。
- 验证：AESE 全量 Go test/vet、pack validate、100 次 hash、snapshot/失败门、34 项 Vitest、TypeScript 和生产构建通过；IAOS 四个 Go module test/vet、Code Map/Atlas tracking、真实 201/duplicate/未验收付款拒绝和 8082 部署通过；Plant Build 三视口 Playwright 与最终 Atlas 同步记录在 M10 evidence。
- 后续：M10 已完成，当前无 active 主计划；M11 必须另立计划并只消费机器终态，不从 UI 推断生产能力。

## 2026-07-22 - M11 Genesis 生产能力建设完成

- 变更：完成 PLAN-M11-001 C0-C6/T1-T42；交付资金、三种 acquisition option、七项设备/实验室、十人核心团队、技能/班次联合 gate、检漏校准漂移整改 tracer、Capability Build API/UI。IAOS 独立 revision `789b925` 提供资金、采购、接受、编制、招聘、资格和付款治理。
- 原因：把 M10 设施资格推进为经过资金、实际设备能力、人员到岗、实操资格、安全和治理门的 M12 机器入口。
- 影响：`hctm-genesis@0.4.0` 终态输出 `industrialization_eligible=true`；不包含产品 BOM/routing、APQP、试产、PPAP 或 SOP。
- 验证：两仓 Go test/vet、100 次 hash、Schema/pack、前端 unit/typecheck/build、三视口 Playwright、真实 API/部署与 Atlas 同步通过。
- 后续：当前无 active 主计划；M12 必须另立计划并消费机器终态。

## 2026-07-22 - M12 Genesis 产品工业化与量产批准完成

- 变更：完成 PLAN-M12-001 D0-D7/T1-T53；交付 RFQ/报价/定点、版本化产品/BOM/routing/PFMEA/control plan、供应/工装/物料、两轮试制、泄漏/Cpk 整改、PPAP 与 Industrialization Play。IAOS revision `50a46e2` 交付七类治理权限。
- 原因：把 M11 通用生产能力推进为客户项目、产品工艺和量产质量门均获批的 M13 机器入口。
- 影响：`hctm-genesis@0.5.0` 输出 `serial_production_eligible=true`；RFQ/试制不形成正式订单、库存、收入或回款。
- 验证：两仓 Go test/vet、100 次 hash、Schema/pack、M3/O2D stable-code compatibility、前端 unit/typecheck/build、三视口 Playwright、真实 API/部署与 Atlas 同步通过。
- 后续：当前无 active 主计划；M13 另立计划并只消费机器终态。

## 2026-07-22 - M13 Genesis 第一次完整商业交付完成

- 变更：完成 PLAN-M13-001 E0-E8/T1-T60；交付首单/追加、供应生产、三批接受、300 件恢复、发票/AR/银行核销、实际成本与毛利、First Delivery Play 和 M9-M13 端到端证据。IAOS revision `067bbb4` 交付十类治理权限。
- 原因：把 M12 量产资格推进为第一张订单的真实商业与资金闭环。
- 影响：`hctm-genesis@0.6.0` 输出 `first_commercial_cycle_closed=true`；Genesis 主纵向场景完成，长期多周期经营属于 M14。
- 验证：两仓 Go test/vet、100 次 hash、Schema/pack、旧 O2D compatibility、前端 unit/typecheck/build、三视口 Playwright、真实 API/部署与 Atlas 同步通过。
- 后续：当前无 active 主计划；M14 另立参数化实验计划。

## 2026-07-22 - M11 生产能力建设计划启动

- 变更：新增 approved DES-012 和唯一 active 的 PLAN-M11-001，将 M11 拆为 C0-C6：机器合同、资金与受治理采购、AESE 设备世界、AESE 人员世界、IAOS 采购/资产/组织/资格治理、Capability Build Play 和全链验收；同步 README、Agent Context、Architecture、文档索引、Roadmap、Code Map 和 System Atlas 声明。
- 原因：M10 已输出 `capability_build_eligible=true`，但设施验收不等于拥有设备、实验室、仓储、人员和岗位技能；同时 closing cash 只有 10,000,000 CNY，必须先建立真实资金来源和预算约束。
- 影响：M11 成为当前唯一 active plan，目标终态固定为 `industrialization_eligible=true`。M11 显式处理剩余资本实缴、采购/租赁、设施尾款、工资准备金、设备调试和团队资格，但不进入产品/BOM/工艺、APQP、试生产、PPAP 或 SOP。
- 验证：`git diff --check`、DES/PLAN ID 唯一性、active plan 唯一性、全部 Atlas/World/场景 JSON 解析、受影响 Markdown 本地链接和 `scripts/check_system_atlas_tracking.sh` 均通过；本轮只修改设计与计划文档，未运行产品测试。工作区既有测试修改、截图变化和验收产物均未触碰。
- 后续：执行 C0 T1-T7，关闭 G4-G7 的能力/设备、人员、资金和 IAOS gap 基线，再允许独立 IAOS worktree 开发。

## 2026-07-22 - M12 产品工业化与量产批准计划启动

- 变更：新增 approved DES-013 和唯一 active 的 PLAN-M12-001，将 M12 拆为 D0-D7：机器合同、RFQ/报价/定点、产品/工艺/APQP 版本、供应商/工装/首批物料、试生产/质量/PPAP、IAOS 客户项目与工程质量治理、Industrialization Play 和全链验收；同步 README、Agent Context、Architecture、文档索引、Roadmap、Code Map 和 System Atlas 声明。
- 原因：M11 已输出 `industrialization_eligible=true`，但通用设备和人员能力不等于具体客户产品、工艺和量产质量已获批；现有 O2D fixture 也不能替代 Genesis 的 RFQ、工程和 PPAP 因果链。
- 影响：M12 成为当前唯一 active plan，目标终态固定为 `serial_production_eligible=true`。客户开发预付款按现金与合同负债处理；试制件不可销售。M12 通过 release manifest/hash 与旧 HCTM stable code 兼容，但不进入 M13 的正式订单、量产、交付、开票或回款。
- 验证：`git diff --check`、DES/PLAN ID 唯一性、active plan 唯一性、T1-T53 连续唯一、全部 Atlas/World/场景 JSON 解析、受影响 Markdown 本地链接和 `scripts/check_system_atlas_tracking.sh` 均通过；本轮只修改设计与计划文档，未运行产品测试。工作区既有测试修改、截图变化和验收产物均未触碰。
- 后续：执行 D0 T1-T9，关闭 G4-G8 的客户/资金/工程质量/兼容/IAOS gap 基线，再允许独立 IAOS worktree 开发。

## 2026-07-22 - M13 第一次完整商业交付计划启动

- 变更：新增 approved DES-014 和唯一 active 的 PLAN-M13-001，将 M13 拆为 E0-E8：机器合同、M12→Genesis O2D 适配、正式订单/供应、生产/质量/成本、三批交付与 300 件恢复、开票/应收/回款/利润、IAOS 治理、First Delivery Play 和 Project Genesis 总验收；同步 README、Agent Context、Architecture、文档索引、Roadmap、Code Map 和 System Atlas 声明。
- 原因：M12 已输出 `serial_production_eligible=true`，但量产批准不等于已获得正式订单、完成物理交付、形成现金或证明盈利；M3-M7 的 O2D fixture 也不能替代 Genesis 的真实首单历史。
- 影响：M13 成为当前唯一 active plan，目标终态固定为 `first_commercial_cycle_closed=true`。M13 从零可销售库存完成 12,000 件三批交付，补齐 invoice/AR/cash settlement、actual cost 和项目毛利，并收口 M9-M13；长期多周期经营和参数实验仍属于 M14。
- 验证：`git diff --check`、DES/PLAN ID 唯一性、active plan 唯一性、T1-T60 连续唯一、全部 Atlas/World/场景 JSON 解析、受影响 Markdown 本地链接和 `scripts/check_system_atlas_tracking.sh` 均通过；本轮只修改设计与计划文档，未运行产品测试。工作区既有测试修改、截图变化和验收产物均未触碰。
- 后续：执行 E0 T1-T9，先关闭 G4-G9 的财务结转、订单/履约、成本、兼容和 IAOS gap 基线，再允许独立 IAOS worktree 开发。

## 2026-07-22 - M14 参数化分支经营实验计划启动

- 变更：新增 approved DES-015 和唯一 active 的 PLAN-M14-001，将 M14 拆为 X0-X7：实验机器合同、确定性随机流/矩阵、checkpoint 分支与持久目录、多周期策略、执行器、KPI/EvidenceBundle、IAOS 实验治理和 Scenario Lab；同步 README、Agent Context、Architecture、文档索引、Roadmap、Code Map 和 System Atlas 声明。
- 原因：M13 已证明从成立到首单回款和毛利的单条确定性路径，但单次成功不能回答不同需求、供应、设备、质量和付款假设下哪种经营策略更稳健。
- 影响：M14 成为当前唯一 active 主计划，终态固定为 `strategy_evidence_ready=true`。父 checkpoint、正式 IAOS 数据和兄弟分支必须隔离；策略用共同随机数进行 paired comparison，推荐只能形成证据/治理意图，不能自动投放。
- 验证：`git diff --check`、DES/PLAN ID 唯一性、active plan 唯一性、T1-T64 连续唯一、全部 Atlas JSON 解析、受影响 Markdown 本地链接和 `scripts/check_system_atlas_tracking.sh` 均通过；本轮只修改设计与计划文档，未运行产品测试。工作区既有测试修改、截图变化和验收产物均未触碰。
- 后续：执行 X0 T1-T9，先关闭 G4-G8 的 checkpoint/opening、参数/seed、策略/KPI、运行容量和 IAOS gap 基线，再允许分支持久化与 IAOS 写端点开发。

## 2026-07-22 - M14 参数化分支经营实验完成

- 变更：完成 PLAN-M14-001 X0-X7/T1-T64；交付严格实验合同、命名 PRNG 流、共同随机数、60-run 隔离矩阵、12 周经营规则、run-level KPI、paired delta、Pareto、EvidenceBundle、CLI/API、Scenario Lab 和 IAOS 实验/推荐治理。
- 原因：让 M13 单次成功可在相同外生扰动下公平比较 baseline、lean 和 resilient，同时保持模拟证据与正式经营决策不可绕过的隔离门。
- 影响：`hctm-genesis@0.7.0` 输出 `strategy_evidence_ready=true`；所有 run 的 production writes 为零，推荐状态固定为 proposed_not_applied。当前无 active 主计划。
- 验证：共同随机数配对、命名流独立、100 次 evidence hash、60 个唯一 branch/run、严格 schema/CLI dry-run、Go test/vet、前端 unit/typecheck/build、三视口 Playwright、IAOS tenant/RLS/idempotency/journal/Outbox 与真实 API 重复提交通过。
- 后续：任何策略投放必须另立计划并经过独立 IAOS intent/审批；不得直接消费 M14 推荐修改正式 Policy、订单、预算、采购、排产或现金。

## 2026-07-22 - M15 受治理策略发布与经营试点计划启动

- 变更：新增 approved DES-016 和唯一 active 的 PLAN-M15-001，将 M15 拆为 R0-R7：决策/发布/安全合同、Evidence Review、StrategyRelease 编译、零写入 shadow、canonical pilot、guardrail/回滚/补偿、IAOS 采纳治理和 Strategy Control Room；同步 README、Agent Context、Architecture、文档索引、Roadmap、Code Map 和 System Atlas 声明。
- 原因：M14 已形成完整但不可自动投放的策略证据，下一步必须证明组织能够在不跳过审议、不伪造回滚、不隐藏风险的前提下把 evidence 转化为有限行动并诚实关闭决策。
- 影响：M15 成为当前唯一 active 主计划，机器终态固定为 `strategy_change_cycle_closed=true`，disposition 可以是 adopted、rejected 或 rolled_back。shadow 必须零业务写入；pilot 只在批准 scope/window 内生效；回滚只停止未来动作，既成后果进入 commitment/compensation 链。
- 验证：`git diff --check`、DES/PLAN ID 唯一性、active plan 唯一性、T1-T69 连续唯一、全部 Atlas/World/场景 JSON 解析、受影响 Markdown 本地链接和 `scripts/check_system_atlas_tracking.sh` 均通过；本轮只修改设计与计划文档，未运行产品测试。工作区既有测试修改、截图删除和验收产物均未触碰。
- 后续：执行 R0 T1-T9，先关闭 G4-G8 的 candidate/release、shadow/pilot、guardrail/rollback 和 IAOS gap 基线，再允许 Policy 激活或 canonical pilot 开发。

## 2026-07-22 - M15 受治理策略发布与经营试点完成

- 变更：完成 PLAN-M15-001 R0-R7/T1-T69；交付 immutable StrategyRelease、semantic diff、独立审议、4 周零写入 shadow、4 周 canonical pilot、guardrail、kill switch、rollback、commitment/compensation、Control Room 和 IAOS 治理。
- 原因：把 M14 模拟证据安全推进到可拒绝、可暂停、可回滚且不删除历史的真实决策闭环。
- 影响：`hctm-genesis@0.8.0` 输出 `strategy_change_cycle_closed=true`；adopted 与 rolled_back 都是合法终态，pilot 不被宣称为统计因果证明。当前无 active 主计划。
- 验证：两条路径 100 次 hash、shadow 零写入、exact release、职责分离、真实 IAOS 201/duplicate、Go/Schema/前端/三视口、Atlas 与服务部署通过。
- 后续：M16 必须依据 disposition 另立计划；M15 不授权真实生产租户或无人审批投放。

## 2026-07-22 - M16 持续策略保障与假设校准计划启动

- 变更：新增 approved DES-017 和唯一 active 的 PLAN-M16-001，将 M16 拆为 A0-A7：Assurance/Dataset/Drift 合同、canonical observation lineage、数据质量与 drift、有界校准和防泄漏、holdout/M14 replay、IAOS 到期复审治理、Assurance Observatory 和全链验收；同步 README、Agent Context、Architecture、文档索引、Roadmap、Code Map 和 System Atlas 声明，并修正 M15 Code Map 为实际实现路径。
- 原因：M15 主路径已采纳 resilient release，但 adoption 只对当时 scope/version 有效；必须在复审日期前证明观察数据可信、假设仍受支持，并在不静默改 evidence/Policy 的情况下决定续期、重新实验或退役。
- 影响：M16 成为当前唯一 active 主计划，机器终态固定为 `strategy_assurance_cycle_closed=true`，disposition 可以是 renewed、reexperiment_required 或 retired。数据质量优先于 drift；前 8 周 calibration 与后 4 周 holdout 隔离；校准只形成假设 candidate。
- 验证：`git diff --check`、DES/PLAN ID 唯一性、active plan 唯一性、T1-T68 连续唯一、全部 Atlas/World/场景 JSON 解析、受影响 Markdown 本地链接和 `scripts/check_system_atlas_tracking.sh` 均通过；本轮只修改设计与计划文档，未运行产品测试。工作区既有测试修改、截图删除和验收产物均未触碰。
- 后续：执行 A0 T1-T9，先关闭 G4-G8 的 observation/cutoff、drift、calibration/holdout、replay 和 IAOS assurance 基线，再允许封存 dataset 或改变 release review 状态。

## 2026-07-22 - M16 持续策略保障与假设校准完成

- 变更：完成 A0-A7/T1-T68；交付 12 周 as-of dataset/lineage、质量优先 drift、8/4 防泄漏校准/holdout、新祖先 60-run replay、三种 disposition、Assurance Observatory 和 IAOS 治理。
- 原因：验证 M15 已采纳策略在后续环境中仍受假设支持，而不在线学习或自动修改 release。
- 影响：输出 `strategy_assurance_cycle_closed=true`；renewed、reexperiment_required、retired 均是合法终态。当前无 active 主计划。
- 验证：100 次 cycle hash、zero missing、six-domain drift、holdout lock、旧 evidence hash、Go/Schema/UI/三视口、IAOS 201/duplicate 和 Atlas 通过。
- 后续：按 disposition 另立下一计划；不得直接扩范围或自动发版。

## 2026-07-22 - AESE 3.0 后续完成体设计与 M17 计划启动

- 变更：新增 approved DES-018 总纲及 DES-019 至 DES-026 八个独立里程碑设计，将后续完成体固定为 M17 滚动 IBP、M18 组合扩展、M19 多基地网络、M20 售后质量、M21 工厂韧性、M22 集团价值、M23 多 Agent 组织和 M24 场景平台产品化；新增唯一 active 的 PLAN-M17-001（B0-B7/T1-T66），同步 README、Agent Context、Architecture、文档索引、Roadmap、Code Map 和 System Atlas，并修正 M16 Code Map 为实际实现路径。
- 原因：M16 主路径 renewed 后，系统已完成单产品从 evidence 到持续复审的纵向链；用户要求一次性写出全部后续设计，需要把扩展顺序、依赖、终态和 program 结束边界一次冻结，同时避免多个 active plan 并行造成失控。
- 影响：M17 成为当前唯一 active 主计划，目标为 `integrated_plan_cycle_closed=true`；M18-M24 仅为 approved design，实施状态仍是 Planned。M24 以 `industry_simulation_platform_ready=true` 关闭本轮 AESE 3.0 program；真实生产、法定合规、第二行业和高精度 3D 必须另立新 program。
- 验证：`git diff --check`、DES-018 至 DES-026/PLAN-M17 ID 唯一性、active plan 唯一性、T1-T66 连续唯一、全部 Atlas/World/场景 JSON 解析、受影响 Markdown 本地链接和 `scripts/check_system_atlas_tracking.sh` 均通过；本轮只修改设计与计划文档，未运行产品测试。工作区既有测试修改、截图删除和验收产物均未触碰。
- 后续：执行 M17 B0 T1-T9，先关闭 G4-G8 的 horizon、计划语义、scenario、review gate 和 IAOS gap；M18 只能在 M17 terminal/evidence 完成后另立 active plan。
## 2026-07-22 - M17 滚动 IBP 与 S&OP 完成

- 变更：封存 13 周 weekly、12 月 monthly、三 scenario 和五级 review 的 M17 evidence frame。
- 原因：把 M16 renewed 事实转为不自动执行的唯一 PlanRelease 证据。
- 影响：`integrated_plan_cycle_closed=true`，自动业务写入为零。
- 验证：strict Go validation、100 次 canonical hash、schema/fixture、API/UI 和 IAOS `ibp.release` 门通过。
- 后续：按 M18 扩展产品/客户组合。

## 2026-07-22 - M18 产品与客户组合完成

- 变更：加入第二产品、第二客户和共享 A 线分配证据。
- 原因：验证组合权衡不会重复消费共享能力。
- 影响：2 产品、2 客户、0 capacity violation，`portfolio_operating_model_validated=true`。
- 验证：AESE3 contract/hash 与 IAOS `portfolio.allocate` evidence/approval 门通过。
- 后续：按 M19 扩展履约网络。

## 2026-07-22 - M19 多基地网络完成

- 变更：交付三节点、两 lane、在途 custody 和 disruption/recovery frame。
- 原因：建立不伪造发运/收货事实的网络重排证据。
- 影响：0 unreconciled in-transit，`network_operating_model_validated=true`。
- 验证：数量守恒、strict contract/hash 与 IAOS `network.replan` 门通过。
- 后续：按 M20 关闭客户生命周期。

## 2026-07-22 - M20 售后质保与闭环质量完成

- 变更：交付 complaint、120 件 RMA、containment、8D/CAPA 和 replacement/credit frame。
- 原因：把首交付后的现场质量反馈闭环到可审计客户结果。
- 影响：0 unit variance，`customer_lifecycle_closed=true`。
- 验证：lot/数量 reconciliation 与 IAOS `quality.close` 独立审批门通过。
- 后续：按 M21 验证工厂资源/EHS 韧性。

## 2026-07-22 - M21 工厂资源与 EHS 韧性完成

- 变更：注入 near miss 和 utility outage，完成安全 hard stop 与恢复 frame。
- 原因：证明服务目标不能绕过人员资格和 EHS 约束。
- 影响：0 safety bypass，`plant_resilience_cycle_closed=true`。
- 验证：安全/恢复不变量与 IAOS `resilience.recover` 门通过。
- 后续：按 M22 建立集团价值治理。

## 2026-07-22 - M22 集团财务资金与投资完成

- 变更：交付管理 P&L、cash/working-capital 与 capex decision frame。
- 原因：把经营后果转为不冒充法定会计的管理价值视图。
- 影响：现金/利润 0 conflation，`group_value_cycle_closed=true`。
- 验证：decimal/owner/reconciliation 与 IAOS `finance.invest` 门通过。
- 后续：按 M23 资格化多 Agent 组织。

## 2026-07-22 - M23 受治理多 Agent 组织完成

- 变更：交付七 Agent mandate、三 benchmark、知识/tool 隔离和人工接管 frame。
- 原因：验证多 Agent 协作不会形成 sole approval 或越权写入。
- 影响：0 unauthorized write，`agent_operating_model_qualified=true`。
- 验证：normal/adversarial/recovery evidence 与 IAOS `agent.approve` 门通过。
- 后续：按 M24 封装 reference platform release。

## 2026-07-22 - M24 场景平台产品化完成

- 变更：发布 `hctm-genesis@1.0.0`、统一 AESE3 schema/fixture/API、Completion Room、runbook/evidence 和五 certification gates。
- 原因：把 M17-M23 能力封装为可重复验证且不自动执行正式业务的行业 reference pack。
- 影响：五门 0 failure，`industry_simulation_platform_ready=true`；PLAN-M17 至 PLAN-M24 全部完成，当前无 active 主计划。
- 验证：Go 全量 test/vet、100 次 hash、JSON/Atlas、前端 unit/typecheck/build、IAOS API governance tests 通过。
- 后续：真实生产、第二行业、法定合规或高精度 3D 必须另立 program。

## 2026-07-23 - World 企业生命周期入口与导航修复

- 变更：将 `#world` 从 `LAS-WLD-02` 设备 tracer 重构为 M8-M24 企业生命周期运营中心，明确公司成立、工厂建设、能力建设、产品工业化、商业交付和持续经营主路径；新增各阶段深链接、全程快速导航和独立 `#world-tristate` 架构验证入口，并将主沙盘按钮改为可见的“企业生命周期”。
- 原因：M8 最小三态 tracer 被错误当成整体 World 首页，阶段页面虽存在但缺少统一目录和持久导航，用户无法理解或连续浏览企业运营过程。
- 影响：用户可从首页选择任一过程，进入查看阶段步骤、World/IAOS 交换和治理边界；M8 tracer 回归其架构验证定位。该修复不扩大 IAOS 业务实现范围。
- 验证：前端 38 个 unit tests、typecheck、production build 通过；新生命周期 E2E 在 1440、1280 和 390 三视口共 6 条通过，覆盖完整阶段可见、M9 深链接和 M8 次级入口。
- 后续：M17-M24 仍需从联合 evidence 视图拆为独立运营工作台，并把页面推进接入已部署 IAOS committed outcome。

## 2026-07-23 - M9 IAOS 原生真实闭环设计启动

- 变更：新增 draft DES-027，确认 M9 必须从 IAOS Core/Domain/Tenant 三层语义出发，经 Archetype、Entity、Atomic Ability、Business Capability、Process/Policy/Decision、Runtime Artifact、权限、菜单和工作台进入 World Bridge；确认建立可复用 `enterprise_governance` Domain Semantic Package，而不是 HCTM 专属 CRUD。
- 原因：现有 M9 只有 AESE frame、通用治理 receipt 和失败 Outbox，不能验证 IAOS 是否能从能力语义模型满足现实企业成立需求。
- 影响：DES-010 的所有权原则保留，但 M9 实现边界重新进入设计；DES-027 批准前不开始 IAOS 业务实现。
- 验证：对照 IAOS DES-023 五层语义模型、Semantic Archetype Catalog、Metadata Compiler、Capability Runtime 和动态菜单实现完成现状核验。
- 后续：逐项确认 Core 扩展、Domain graph、能力/流程、World Bridge、UI、seed 和验收边界。

## 2026-07-23 - M9 IAOS 原生真实闭环设计批准

- 变更：完成 DES-027 D1–D18 决策并转为 Approved；正式平台主体确定为 `founder-principal`（显示名称“创始治理者”），`dev-user` 仅保留为本地兼容身份；确认二十个 Business Capability、一主四子 Process、八项 Policy、G1–G7 人工治理门、Trace Spine 和十二项最终验收门。
- 原因：M9 正式验收不能依赖特殊开发账号、AESE frame 或通用治理 receipt，必须证明 IAOS 从三层语义资产到身份、能力、流程、权限、菜单、Agent 与 World Bridge 的真实闭环。
- 影响：DES-027 已具备进入实施计划的批准基线；实现前仍需创建唯一 active 主计划，并在 IAOS 独立 branch/worktree 中完成平台改动。
- 验证：逐项核对断线前会话确认记录与 DES-027 D1–D17，补录 `founder-principal` 身份模型和 D18 十二项验收门；文档状态与索引同步为 Approved。
- 后续：编制跨 AESE/IAOS 的纵向实施计划，明确仓库所有权、依赖顺序、每个切片的真实测试证据和双仓提交边界。

## 2026-07-23 - M9 IAOS 原生真实闭环实施计划启动

- 变更：新增唯一 active 的 PLAN-M9-NATIVE-001，将 DES-027 拆为 P0–P7、T1–T66：基线审计与双仓隔离、平台身份、三层语义与 Runtime Artifact、登记纵向 tracer、完整成立主链、Agent 异常治理、工作台与 Trace Spine、恢复对账与最终验收。
- 原因：按技术层横向铺开会延迟发现身份、事务、Bridge 和运行资产不一致；实施必须从最小登记 tracer 开始，每个切片交付可运行的纵向证据。
- 影响：M9N 成为当前唯一 active 主计划；既有 M9–M24 reference completion 保留，但 M9 的 IAOS 原生真实闭环在 D18 十二项验收门全部通过前不得标记实现完成。
- 验证：按跨仓规则阅读 IAOS AGENTS、Agent Context 和 Code Map；确认 IAOS 已有 Runtime Artifact、Capability、Process/Approval、Policy/Decision、Outbox 等基础，并识别 `dev-user` 特判为 P1 硬门；计划明确双仓所有权、独立 worktree、依赖顺序和每切片交付纪律。
- 后续：执行 P0/T1–T8，先记录双仓基线并完成机器可读资产审计，不开始业务 Runtime Artifact 发布。

## 2026-07-23 - M9N 双仓基线与领域状态机落地

- 变更：从 IAOS `origin/main@8e267f7` 创建 `/iaos/iaos-go-m9-native`、`feat/m9-native-incorporation`，引入既有 M8 Bridge/M9 migration baseline；新增机器资产审计、冻结合同，以及 IAOS `internal/incorporation` 的 20 Capability、5 Process、8 Policy、7 Gate 目录和确定性状态机。
- 原因：通用 `genesis_governance_record` 不能承载企业成立事实；先冻结跨仓合同并以可执行领域状态机锁住正常与异常语义，避免 API、UI 和 Agent 各自实现规则。
- 影响：`tenant-hctm` 保持不变，新实现限定 `tenant-hctm-genesis`；失败 Outbox 不再作为完成证据。计划仍为 Active，尚未满足 D18。
- 验证：IAOS `go test ./internal/incorporation` 通过，覆盖正常终态、重复幂等、登记补正、开户拒绝、出资差异、Agent 自批拒绝和目录基数；新增 JSON 均通过标准解析。
- 后续：把领域状态机接入 PostgreSQL 原子事务、正式身份/Runtime Artifact/Approval/Journal/Outbox 与双向 World Bridge，再完成工作台和联合验收。

## 2026-07-23 - M9N 成立事实原子持久化 tracer

- 变更：IAOS 新增统一成立 command、trace、evidence API，把 `incorporation_case`、领域 Journal 与 `sys_outbox` 放入同一租户事务；新增 FORCE RLS、幂等碰撞、状态 hash 和生产/integration schema，提交 revision `4a76f38`。
- 原因：正式成立事实不能继续停留在 AESE frame 或通用 JSON receipt，且失败转换不得留下部分业务状态或 Outbox。
- 影响：API、后续 UI 与 Agent 已有单一 Capability 入口；旧 `genesis_governance_record` 仅保留为迁移来源，`tenant-hctm` 未迁移。
- 验证：IAOS 领域/API 单测通过；真实 PostgreSQL 验证重复 no-op、同键异载荷 409、两个租户隔离和 bootstrap 重复执行；非法转换的 sqlmock 测试证明事务 rollback 且不写 Outbox。
- 后续：接入 `founder-principal` 正式身份、G1–G7 Approval/Decision、Effective Runtime Artifact 和 registration World 往返。

## 2026-07-23 - M9N 正式主体与治理门接线

- 变更：IAOS 新增默认 dry-run、显式 apply 的 founder bootstrap，持久化 `founder-principal`、平台角色、两个租户可分别授予的访问关系、普通登录绑定、董事长岗位、Mandate、Capability 权限和 Outbox；G1–G7 改由 Approval Runtime 提交、决定并在成功业务事务中 consume，revision `63535a5`。
- 原因：正式验收不能依赖 `dev-user`，请求体中的 `approved_by` 也不能替代真实 Approval 决定。
- 影响：`founder-principal` 可通过普通登录获得 `platform_super_admin`，`/profile` 返回主体、平台角色、租户访问、岗位和 Mandate；Semantic Studio 不再按租户字符串由前端二次隐藏。
- 验证：真实 PostgreSQL 已通过 bootstrap、普通登录、授权摘要、Capability 权限、G1 pending→approved→consumed、重复执行和 RLS 验证；IAOS Go 单测、前端 typecheck/build 通过。
- 后续：发布三层语义和 Effective Runtime Artifact，完成 Process/Decision、World 往返与五 Agent。

## 2026-07-23 - M9N Runtime Artifact、World、Agent 与 Trace Spine

- 变更：IAOS revision `52b11ee` 发布 Core → enterprise_governance → HCTM Extension 三层 Effective Runtime Artifact，正式命令对 missing/stale 失败关闭；revision `13d02e1` 把 registration/bank/appointment Intent 原子写入 Bridge journal、校验可信 Observation，建立五 Agent 岗位/Mandate/Capability allowlist，并扩展 trace/evidence 聚合。AESE 新增离线 reconciliation 分类器。
- 原因：API、UI 和 Agent 必须消费同一已发布运行资产；外部登记、银行和候选人事实只能来自 World；证据必须能按稳定 correlation 恢复和对账。
- 影响：M9 正式路径已具备 Runtime Authority、外部事实信任门、Agent dispatch 前授权门和 Trace Spine 后端；旧 receipt 不再参与完成判断。
- 验证：IAOS Go/API 与真实 PostgreSQL founder/G1/runtime tests 通过；AESE bridge/worldcontract tests 通过；reconciliation 覆盖 converged、missing、duplicate、lagging、terminal_conflict。
- 后续：完成全正常链与四异常 replay、CommittedOutcome 自动回传、生命周期工作台及重启/三视口联合验收。

## 2026-07-23 - M9N 正常闭环、工作台与离线对账验收

- 变更：IAOS `1c2a8d6`、`8fedf25`、`802d4d5`、`2f24866` 完成完整正常链、三层语义投影、五 Agent 服务主体、CommittedOutcome、企业生命周期工作台及出资差异持久证据；AESE reconciliation 增加 `hash_mismatch` 并提供 `aese reconcile` 离线命令。
- 原因：D18 要求正常链进入真实终态、异常拒绝仍保留升级证据，并且双方能在服务不可用时对持久 journal 做确定性复验。
- 影响：`tenant-hctm-genesis` 可由 `founder-principal` 经 G1–G7 到达 `enterprise_operational_ready`；三类 World 往返均形成 Intent/Observation/CommittedOutcome；对账覆盖五类故障。
- 验证：真实 PostgreSQL 完整链验证 7 个 consumed Approval、3/3/3 World exchange 和终态；出资差异验证状态不推进且 Discrepancy/Journal/Outbox 原子写入；IAOS build、Go tests、AESE bridge/CLI/worldcontract tests，以及 1440×900、1280×720、390×844 Playwright 共三条通过。
- 后续：补齐登记补正、开户拒绝和 Agent 自批的联合 replay，完成重启/乱序恢复、全量回归及 D18 最终证据矩阵。

## 2026-07-23 - M9N Process/Decision 与身份安全门关闭

- 变更：IAOS revision `30cc729` 将五个 Process Definition、主 Process Run 和逐 Capability Decision Audit 接入正式运行表；revision `9eebe1e` 移除 M9 对 `dev-user` 的 Runtime Artifact bypass，并批准 founder 身份迁移安全评审。
- 原因：目录中的 Process/Policy 名称和开发账号兼容不能替代真实运行证据与正式授权。
- 影响：正常链在同一业务事务维护 Process trace 和 Decision Audit；Trace API 可返回两者；任何主体（含 dev-user）均必须消费有效 Runtime Artifact。
- 验证：IAOS Go/API 测试通过；真实 PostgreSQL 完整链断言 1 个 completed Process Run、15 条 Decision Audit、7 个 consumed Approval 和最终状态。
- 后续：按 active plan 剩余未勾选项继续完成异常 replay、租户/撤权矩阵、恢复和双仓最终证据。

## 2026-07-23 - M9N 四条受治理异常真实库验收

- 变更：IAOS revision `955f43f` 增加登记补正、开户拒绝、出资差异和 finance Agent 自批四条真实 PostgreSQL 验收；修正补正观察误生成 CommittedOutcome 的问题。
- 原因：异常线必须与正常线共用正式身份、Capability、Process/Decision、World Bridge、Journal 和 Outbox，不能只以领域单测代替。
- 影响：登记 Agent 使用独立 service-only 主体在原 correlation 重提；银行拒绝和出资差异保持原状态并写异常证据；finance Agent 无法取得 G7 自批权限。
- 验证：IAOS 全量 Go tests 通过；异常 integration suite 与正常完整链通过；部署后 founder 三视口 Playwright 3/3 通过。IAOS 前端全量 Vitest 暴露 15 个既有非 M9 测试失败（多为测试超时、旧英文文案断言及未 mock 的 401），M9 Playwright 与 production build 不受影响，但 T63 暂不关闭。
- 后续：清理 IAOS 既有前端测试基线，再完成重启/乱序、租户撤权和 D18 最终矩阵。

## 2026-07-23 - M9N Runtime Artifact 版本升级与安全回退

- 变更：IAOS revision `3da05b5` 将 Artifact/Compiler 升至 1.1.0，并新增默认 dry-run、显式 apply 的已安装版本 rollback API；回退写 Outbox，版本与当前二进制不兼容时正式命令按 stale 失败关闭。
- 原因：重复安装 no-op 之外还必须证明版本升级不覆盖旧内容 hash，并提供有审计、不会静默运行错误资产的回退边界。
- 影响：旧 1.0.0 资产保留为 inactive revision，1.1.0 发布 Process/Decision 投影；回退操作不删除任何法律事实或历史证据。
- 验证：IAOS Go/API 与正常/异常真实 PostgreSQL suite 通过；版本升级 apply 与后续重复 no-op 均执行成功。
- 后续：验证 tenant-other 安装隔离和回退后的兼容二进制演练。

## 2026-07-23 - M9N 合同兼容与 Evidence Bundle

- 变更：双仓增加有效/破损 contract fixtures 和 fail-closed compatibility tests；IAOS revision `938a263` 增加 versioned evidence bundle、稳定引用清单、bundle hash 与离线 verifier；更新 capability gap ledger、Atlas planned dependency 和风险清单。
- 原因：跨仓 schema 漂移和只依赖在线 API 的证据无法满足可重放、可离线审计的完成门。
- 影响：旧 schema、错误租户/目录基数和篡改 evidence 均被机器拒绝；terminal case 可追溯 case、Runtime Artifact、Process、Approval、World、Journal 和 Outbox。
- 验证：AESE worldcontract compatibility test、IAOS incorporation fixture/evidence tests通过；线上 terminal evidence 直接 pipe 至 verifier 验证通过。
- 后续：继续授权矩阵、tenant-other 隔离和重启恢复验证。

## 2026-07-23 - M9N Runtime 菜单、租户隔离与 Agent 授权矩阵

- 变更：IAOS revisions `40f4ddc`、`5ee2c14`、`1c8769c` 完成 Runtime Artifact tenant-other 隔离、正式 acting subject/岗位/Mandate 校验、Agent 有效期/额度/工具开关，以及 Runtime Artifact + tenant access + RBAC 菜单投影。
- 原因：token 中的角色、请求体 actor 或前端隐藏均不能替代服务端有效授权和租户隔离。
- 影响：`tenant-hctm` 未安装或清理 M9N；撤权后菜单消失且写入失败；暂停、过期、撤销、超额、跨租户和禁用工具在 dispatch 前失败关闭。
- 验证：真实 PostgreSQL founder/profile/menu 测试覆盖 active/revoked access；actor authorization suite 覆盖七类失败；tenant-other 与 genesis 各有独立 active artifact，tenant-hctm 计数不变。
- 后续：五 Agent 全权限/并发矩阵和 Runtime 表单/动作投影继续收口。

## 2026-07-23 - M9N World Bridge 失败关闭与恢复矩阵

- 变更：IAOS revision `efdd6e2` 对 M9 Observation 强制校验同租户、同 subject、同 correlation 的既有 governed Intent，并增加旧 schema、错误租户、重复/碰撞、服务重启及 poller 恢复的真实 PostgreSQL测试；AESE reconciliation 增加乱序和延迟到达的最终收敛回归。
- 原因：至少一次投递只有在未知或乱序消息失败关闭、重复效果恰好一次，并且双方重启后从持久事实恢复时才能形成可靠闭环。
- 影响：未关联 Observation 返回 `unknown_or_out_of_order_correlation` 且零写入；相同幂等键与相同 payload 返回原结果，碰撞 payload 被拒绝；重启后的 API 和 poller 从数据库恢复。
- 验证：IAOS `TestIntegrationWorldBridgeRecoveryMatrix` 在真实 PostgreSQL 通过；AESE `go test ./internal/bridge/iaos` 通过并覆盖 shuffled 与 delayed convergence。
- 后续：继续完成正式 override、业务对象守恒/readiness、五 Agent 全矩阵和双仓 UI 联动。

## 2026-07-23 - M9N 业务事实守恒与 Readiness

- 变更：IAOS revision `e037abe` 将法律主体、银行账户、组织、任命、Operating Mandate、出资承诺、核验到账和预算授权固化为独立稳定引用/金额事实；readiness evaluator 对引用完整性、CNY、承诺与到账相等、预算不超过核验现金执行失败关闭。
- 原因：顺序走完 Capability 不能替代企业已具备可运营条件的事实证明。
- 影响：G3–G7 对应的业务事实进入设立案持久 state document；事实缺失或金额不守恒时保持 `initial_budget_approved`，不得进入 `enterprise_operational_ready`。
- 验证：IAOS incorporation 单元测试覆盖完整一致事实和不一致拒绝；真实 PostgreSQL 完整生命周期再次通过并到达终态。
- 后续：实现正式 override/批准失效矩阵、五 Agent 并发审计和完整 UI Runtime 投影。

## 2026-07-23 - M9N 开户拒绝补正与 G3 重新审批

- 变更：IAOS revision `eda6a41` 允许开户拒绝后通过原正式 Capability 重提修改后的受益所有人材料，但强制使用新 correlation 和新的 G3 Approval；旧批准不能复用。
- 原因：外部银行拒绝必须使原申请授权边界失效，补正不能用状态直改或历史批准绕过治理。
- 影响：未取得新 G3 的重提返回 422 且零写入；新批准消费后生成唯一新 Intent，设立案仍保持可审计的开户提交状态。
- 验证：IAOS engine 单元测试与真实 PostgreSQL `TestIntegrationCapitalMismatchCommitsDiscrepancyWithoutStateAdvance` 通过，断言旧 G3 复用失败、新 G3 consumed 且新 correlation 只有一个 Intent。
- 后续：继续正式 override/超时/撤权矩阵、五 Agent 全回归和 UI 联动。

## 2026-07-23 - M9N 五 Agent 权限、并发与审计矩阵

- 变更：IAOS revision `039921e` 对设立、治理、法务合规、财务和审计五个 service-only Agent 执行正式主体、岗位、Mandate、Capability 可见范围、越权拒绝、幂等并发和 Decision/Journal 审计矩阵。
- 原因：单个 finance Agent 的异常测试不能证明五个 Agent 均受各自知识与工具边界约束。
- 影响：每个 Agent 只能执行 Runtime Artifact/RBAC/Agent allowlist 交集内的 Capability；相同命令并发只产生一个业务效果和一组审计证据。
- 验证：真实 PostgreSQL `TestIntegrationFiveAgentPermissionIdempotencyAndAuditMatrix` 通过，断言五条允许、五条拒绝以及并发结果 201/200、Journal=1、Decision=1。
- 后续：继续人工 override、Runtime 全投影、AESE 生命周期页面和最终全量回归。

## 2026-07-23 - M9N 正式 Override、Runtime 全投影与双向 Trace

- 变更：IAOS revisions `7fe197b`、`7d48a43`、`93f2440`、`c66ddd6` 将 founder override 接入 Capability/Approval/Decision/Journal/Outbox，批准绑定 Runtime hash 和 30 分钟有效期；Runtime Artifact 统一 API/人工/Agent/Process 入口及动作阶段；Trace/工作台增加全局搜索、来源影响和 lineage。AESE M9 页面消费 IAOS lifecycle/process projection，并实现 tenant/case/process run/world run/correlation 五参数双向深链。
- 原因：人工特批和 UI 动作不得绕过正式治理；两端不得维护互相脱离的完成状态。
- 影响：超时或 Runtime 版本变化使批准失效；特批缺少原决定引用、理由或正式 G1–G7 Approval 时失败关闭；AESE 页面展示持久 Intent/Observation/CommittedOutcome/Discrepancy。
- 验证：IAOS 单元/API 与真实 PostgreSQL founder override 测试通过；IAOS frontend production build、AESE IncorporationPlay test 和 production build 通过。
- 后续：执行 clean tracer、双方重启恢复、真实库综合矩阵和最终全量回归。

## 2026-07-23 - M9N 双仓最终验收完成

- 变更：完成 clean 正常/补正 tracer、双方服务与浏览器刷新恢复、AESE reset 法律事实保护、真实 PostgreSQL M9 矩阵、双仓三视口 UI、runbook/evidence 和 D18 十二门；计划状态改为 completed。
- 原因：计划只允许在代码、集成、部署、UI、恢复和业务证据全部存在时关闭。
- 影响：M9N 以 `founder-principal`、Effective Runtime Artifact 和正式 World Bridge 作为权威闭环；AESE 与 IAOS 使用同一持久 lifecycle projection。
- 验证：AESE Go 全量、frontend 38/38、build、Playwright 3/3；IAOS Go 全量、M9 PostgreSQL matrix、frontend 332/332（单 worker）、build、Playwright 3/3；JSON、Code Map 和 Atlas checks 通过。
- 后续：生产部署前替换 development fallback secrets；超出单法人/CNY/五 Agent 的范围另立计划。

## 2026-07-23 - 修复 M9 局域网加载与 SSE 60 秒截断

- 变更：移除不存在的 `INC-HCTM-001` 默认值，空输入通过 tenant-scoped recent API 自动加载最近 case；双仓 URL 使用浏览器 hostname；SSE heartbeat 续写 deadline；增加局域网 Playwright 参数化。
- 原因：localhost 只代表浏览器所在机器，且 net/http 固定 WriteTimeout 会截断无限流。
- 影响：从 `192.168.50.222` 访问时 API、IAOS、AESE 与双向深链保持同一主机；SSE 不再每 60 秒异常断开。
- 验证：局域网 IAOS/AESE Playwright 各 3/3；Go/API、AESE 38 tests 和双仓 build 通过；SSE 70 秒探针由客户端超时主动结束。
- 后续：生产环境使用反向代理统一 origin，并关闭 Vite HMR。

## 2026-07-23 - 修复 AESE 陈旧 IAOS 地址导致的生命周期 404

- 变更：AESE IAOS base resolver 识别并拒绝指向当前 AESE origin 的陈旧 localStorage 配置，回退到浏览器 hostname 的 8082；补充 favicon 和陈旧配置网络回归。
- 原因：旧浏览器状态可覆盖新的动态 fallback，使生命周期 API 错误发往 Vite 4173。
- 影响：用户无需清理 localStorage；局域网打开 AESE World 时生命周期请求自动路由到 IAOS。
- 验证：针对性 Vitest 2/2、production build、局域网三视口 Playwright 3/3 通过；Playwright 断言所有 incorporation 请求端口均为 8082。
- 后续：生产环境部署时仍建议通过统一反向代理和显式环境配置消除开发端口依赖。

## 2026-07-23 - 修复 IAOS 到 AESE 跨 Origin 租户会话交接

- 变更：IAOS World 深链在 fragment 中交接当前 JWT；AESE 接收后持久化并立即清除 URL token。
- 原因：`:3000` 与 `:4173` 无法共享 localStorage，AESE 的旧租户 JWT 触发 RLS 404。
- 影响：从企业生命周期点击“打开 AESE World”时使用同一登录租户，且不放宽 RLS。
- 验证：Vitest 3/3、AESE build、陈旧跨租户 token 三视口 Playwright 3/3、刷新恢复通过。
- 后续：生产统一 origin 后可改为后端一次性交接码，避免任何长期 token 进入浏览器 URL。

## 2026-07-23 - 撤回 M9N 误完成并补齐通用平台资产

- 变更：撤回原 completed 结论并恢复 active remediation；DES-027 新增 D19 可发现性机器门。IAOS Runtime 1.2.3 将 11 Entity、20 Capability、5 Process、8 Policy Profile、8 Policy Rule、10 条语义关系和 5 Agent 写入通用注册中心，将现有设立案投影到 11 个租户隔离 Entity 表，并为 founder 开放 Studio 菜单；企业生命周期增加十个业务工作区和新建设立案入口。
- 原因：原验收只证明专用状态机和 trace 闭环，错误地把 artifact 中的声明当成通用 Semantic/Entity/Capability/Process/Policy 注册完成，导致用户无法发现、理解和操作 M9 资产。
- 影响：M9 资产现在可从 Semantic Studio、Entity Explorer、Capability Studio、Process Studio、Governance Studio 和企业生命周期工作台共同查看；占位流程节点和无 rule 的 Policy 已被真实注册数据替代。
- 验证：通用 API 返回 11/20/5/8/8/5/10；11 个 Entity records API 均 HTTP 200 且有投影记录；重复 apply 返回 `no_op=true,writes=0`；IAOS Go 针对测试和 frontend production build 通过。
- 后续：完成双仓文档/Atlas/代码提交和用户 UI 验收后，才允许重新关闭 M9N。

## 2026-07-23 - 补齐 M9N Core Archetype 字段与语义引用

- 变更：IAOS Runtime 1.2.5 为 account、commitment、mandate 安装 7/7/8 个默认字段及对应 Semantic Concept，注册 `stable_business_code`、`lifecycle_status`、`fact_payload`，并修复 Semantic Studio 原型切换时 snapshot identity 串线。
- 原因：原安装器只创建原型主记录，导致原型不可继承复用；Entity 字段引用未注册 concept；前端短暂组合新原型 code 与旧 snapshot UUID 产生 404。
- 影响：六个 Core Archetype 可检查默认/继承字段，M9 Entity 不再产生三项 `semantic_concept_missing`，快速切换原型不再请求错误 History detail。
- 验证：通用 API 字段数 account=7、commitment=7、mandate=8、document=7、document_line=5、document_with_lines 继承=7；快速切换 History 404=0；局域网 Playwright 4/4。
- 后续：继续以 D19 逐页验收通用 Studio，不以注册数量替代字段、关系、History 和分析器质量。

## 2026-07-23 - 修复 Entity Explorer 跨租户默认实体

- 变更：IAOS 数据模型工坊移除 `sales_order` 硬编码初值，等待租户 schema 目录返回后保留有效选择或选择第一项；空目录不加载详情。
- 原因：`tenant-hctm-genesis` 没有销售订单，但组件在目录返回前按示例实体请求 schema/ui，产生无意义 404。
- 影响：工作室默认选择完全由当前租户实际目录驱动，不再把其他业务阶段或租户的实体假设带入 M9。
- 验证：局域网 E2E 断言选中 `/metadata/schemas` 第一项，schema/ui 404=0；完整 M9 Playwright 4/4、production build/deploy 通过。
- 后续：继续检查其他 Studio 是否存在全局示例值覆盖租户目录的同类问题。

## 2026-07-23 - 打通 Archetype 默认字段到 Entity 有效模型

- 变更：DES-027 D19 明确 Core Archetype defaults → Domain Entity fields → Tenant extension 编译规则及非破坏性物理迁移门；IAOS Runtime 1.2.6 已将六类原型默认字段编译到十一项 Entity metadata 和物理投影，并从 payload 增量回填既有数据。
- 原因：原安装器将原型字段与 Entity schema 分开注册，Entity 固定为 business_code/status/payload 三字段，导致语义资产中心与数据模型工坊展示不一致。
- 影响：account、commitment、mandate、organization、role、document 的默认字段成为 Entity effective schema 的真实组成部分；旧设立事实保留，重复安装仍为 no-op。
- 验证：Go incorporation/API 测试通过；真实租户 schema、records API 和 information_schema 一致；bank_account 既有记录 currency 回填为 CNY；IAOS Playwright 4/4。
- 后续：后续 Domain/Tenant 字段扩展必须沿同一编译和迁移路径，不允许旁路直接修改工作室展示数据。

## 2026-07-23 - 增加 M9 Entity 语义发布门

- 变更：DES-027 D19 要求 Runtime 安装在写入前执行语义发布门，并以真实 Analyzer 对十一项 Entity 的零错误零警告作为验收；IAOS Runtime 1.2.7 补齐 enum options 和 system_managed/overridable 继承。
- 原因：1.2.6 编译器只复制字段名、类型和 semantic_id，且安装流程没有 Analyzer 等价门，导致错误直到用户进入数据模型工坊才暴露。
- 影响：enum 无选项或 Archetype 治理属性丢失会在安装前失败关闭，不再发布半有效 Entity schema。
- 验证：修复前 API 稳定复现用户报告；修复后 tenant-hctm-genesis 十一项 Entity 全部 errors=0,warnings=0；Go 测试和 IAOS Playwright 4/4。
- 后续：新增任何 M9 字段规则时必须同时扩展发布门与真实 Analyzer 回归。

## 2026-07-23 - 补齐 Policy Rule UI 与 Founder 编辑闭环

- 变更：DES-027 D19 增加 Profile/Rule 独立明细、业务化解释和持久编辑门；IAOS 规则控制中心新增 Rule 页签，把 Profile JSON 翻译为失败策略、版本和 Runtime 来源，并修复 founder 编辑权限。
- 原因：8 条 Rule 只有 API 无 UI；前端将写权限硬编码为 tenant-001，且 Profile 保存只改内存；原始 definition JSON 对业务用户不可读。
- 影响：founder 可查看 8 个 Profile 和 8 条 Rule 的真实含义，并通过原生 upsert API 持久修改；只读用户仍可打开完整详情。
- 验证：Next production build/deploy、Policy UI Vitest 6/6、Founder governance Playwright 1/1，三视口主体场景 3/3。
- 后续：Policy 变更的 dry-run、审批和审计仍按既有治理链执行，不以 UI 写权限替代业务审批。

## 2026-07-23 - 修复流程工作室菜单业务名称

- 变更：IAOS Runtime 1.2.9 将 M9 十项 menu resource 的权限 code 与展示 name 分离，`menu.process_studio` 显示为“流程编排控制台”。
- 原因：安装器把资源 code 同时写入 name，覆盖业务菜单名称。
- 影响：权限标识保持稳定，用户侧不再暴露机器 code。
- 验证：数据库资源名称正确；Founder 1440×900 Playwright 1/1。
- 后续：新增菜单必须显式提供业务名称，不得用 code 作为默认展示名。

## 2026-07-23 - 修复企业成立与治理页面滚动

- 变更：IAOS 企业成立工作台改为固定高度内的独立纵向滚动容器，并增加真实 scrollTop 三视口回归。
- 原因：原 `min-h-full` 被内容撑高，而 MainLayout 外层锁定 overflow，导致内外层都无法滚动。
- 影响：桌面和移动端可完整访问工作台下部的追踪、资产和审计内容。
- 验证：Next production build/deploy；1440×900、1280×720、390×844 Playwright 3/3。
- 后续：固定 MainLayout 下的新工作区统一使用 `h-full min-h-0 overflow-y-auto`。

## 2026-07-23 - 统一 IAOS 侧栏分组图标

- 变更：业务智造层、智能中枢层、低代码工坊、企业数字治理移除 emoji，改用 IAOS 既有 Lucide 线性图标体系。
- 原因：彩色 emoji 与系统菜单 SVG 的尺寸、描边、颜色和视觉语义不一致，呈现明显模板化观感。
- 影响：四个分组统一使用 Factory、Cpu、LayoutGrid、ShieldCheck，并共享主题颜色和交互风格。
- 验证：Next production build/deploy；三视口 Playwright 断言分组可见且无原 emoji，3/3。
- 后续：结构性导航图标只使用系统图标库，不使用字体 emoji。
## 2026-07-23 - M9 AESE 逐步骤 IAOS 可解释追踪

- 变更：DES-027 增加 D20，固定 8 个 AESE World frame 到 15 次 IAOS Capability transition 的显式映射；World 页面增加 Process、Capability、Entity 影响、治理、Journal/Outbox 与 World Bridge 当前步骤详情，并把 step/capability 带入 IAOS 深链。
- 原因：原页面只显示当前 World frame 和整案 IAOS 原始 trace，无法回答每一步调用了什么能力、处于哪个流程、改变哪些数据。
- 影响：时间线切换时证据同步过滤；未匹配 transition 明示，不再用整案数据伪装当前步骤；技术 JSON 采用可展开的渐进披露。
- 验证：映射单测覆盖 8 frame/15 transition 唯一性与证据过滤；AESE production build 通过；1440×900、1280×720、390×844 Playwright 3/3。
- 后续：保持 M9N active remediation，待双仓治理记录、Atlas 和最终用户验收一致后关闭 T67–T70。
## 2026-07-23 - 修复 M9 Capability DSL Runtime Artifact 访问

- 变更：IAOS Runtime 1.2.11 为 Founder 补齐 `output.template.read`，并为二十项 M9 Capability 安装 active artifact snapshot。
- 原因：Capability Studio 菜单可见不等于 Runtime Artifact 可读；原安装器同时漏掉权限依赖和已发布 DSL snapshot。
- 影响：打开 M9 Capability DSL 不再返回 403 或 `no_active_snapshot`，页面显示实际 Runtime version、compiler 和 artifact hash。
- 验证：Go 针对测试通过；真实租户 apply 后 API 200；Founder Playwright 打开 DSL 并显示 Active v1。
- 后续：后续 Studio 菜单安装必须同时检查全部读侧 API 依赖与发布态 artifact。

## 2026-07-23 - M9 可解释 Capability 与 Process 配置闭环

- 变更：DES-027 增加 D21；IAOS Runtime 1.3.2 为二十项 Capability 发布 Purpose、输入输出、状态、Entity/Journal/Outbox、治理、World 和错误恢复合同；主流程显式引用四个子流程，子流程声明 capability、event/world wait 与 timeout；Capability/Process Studio 增加业务解释层并保留专家 JSON。
- 原因：原通用注册数量可见但配置不可理解，主流程隐藏后四段编排，客户无法知道能力做什么、输出什么、如何治理及如何恢复。
- 影响：业务、治理和集成管理员可按职责理解并配置平台资产；核心语义、事务、幂等、RLS 和审计仍由平台强约束。T71–T75 完成，M9N remediation 重新关闭。
- 验证：Go incorporation/API 测试通过；Process Studio 21/21、TypeScript、production build 通过；真实租户升级到 1.3.2，二十项合同 API 完整，五个 Process Analyzer 均零诊断；Founder Capability 与 Process 浏览器回归 2/2。
- 后续：客户行业扩展应复用同一 Explainable Contract 和发布门，不得退回只有原始 JSON 或隐藏子流程的资产。

## 2026-07-23 - 补齐 M9 流程级业务目的

- 变更：DES-027 D21 增加流程级目的合同；IAOS Runtime 1.3.3 为一主四子流程写入独立业务名称和目的，版本 API 返回描述，Process Studio 在列表和编辑抽屉顶部展示。
- 原因：上一轮只解释了 Capability 和流程节点，没有解释整个流程为什么存在，安装器仍写入统一技术占位描述。
- 影响：用户无需读节点或源码即可先理解流程整体目标，再下钻每一步能力、审批和 World 等待；说明来自正式 `process_definition`，不是前端硬编码。
- 验证：五个流程目的完整性测试、Go incorporation/API、Process Studio 21/21、TypeScript、生产构建、在线 API 与 Founder 浏览器回归通过。
- 后续：新增 Process 的发布门必须同时拒绝空显示名和空业务目的。

## 2026-07-23 - 重新打开 M9 交互式经营闭环

- 变更：DES-027 增加 D22 人类/Agent/审批/World 工作项合同；M9N 计划与路线图恢复为 active；AESE 移除默认自动播放并按 IAOS 已提交状态锁定未来 frame；IAOS 增加 tenant RLS 的 15 节点持久工作项及“我的经营待办”视图。
- 原因：原“十工作区”和自动化集成测试只能证明数据可查及状态机可运行，不能证明真实参与者逐节点输入、审批、等待和恢复。
- 影响：新设立案完成 `incorporation.case.open` 后只解锁 G1 等待项；刷新或服务重启不会丢失任务。M10–M24 的预计算 replay 不再自动等价于交互式完成。
- 验证：IAOS API 单测、两端 TypeScript、IAOS production build 通过；在线创建 `INC-INTERACTIVE-1784818405` 后查询得到 15 项，节点 1 completed、仅节点 2 `waiting_approval`。
- 后续：补 Agent 调度/人工接管动作、审批与 World wait 解锁、独立业务菜单和 Founder+五 Agent 端到端后再关闭 T76–T85。

## 2026-07-23 - M9 工作项执行与 World 人工输入

- 变更：IAOS 增加工作项执行 API、节点 JSON 输入、Gate 提交和正式 Agent 调度；AESE 增加登记、开户、候选人三类 Observation 按钮；IAOS 增加四个独立成立业务菜单。
- 原因：持久工作项如果只能查看，仍不能证明人类和 Agent 实际参与；World 外部结果也不能继续由测试 helper 隐式注入。
- 影响：审批前执行返回 409；Founder 批准后只推进 G1；finance-agent 以自己的 actor、岗位和 Mandate 完成出资承诺；World 按钮只写 Observation，不越权提交正式事实。
- 验证：在线案 `INC-INTERACTIVE-1784818405` 验证 pre-approval 409、G1 approved、节点 2/3 分别由 Founder/finance-agent 提交，流程停在 G2；两仓 TypeScript、AESE build、IAOS Go 测试通过。
- 后续：继续补完整 15 节点浏览器 E2E、人工接管与超时/拒绝恢复。

## 2026-07-23 - M9 交互式主线在线贯通

- 变更：修正 Observation commit 和组织建立节点的主体分类；外部事实由 `world-bridge-runtime` 提交，无判断的组织投影由 `iaos-runtime` 执行，不扩大业务 Agent allowlist；Runtime 升级至 1.3.5。
- 原因：在线 E2E 发现把 Observation commit 分给 incorporation/finance Agent 会被正式 allowlist 正确拒绝；外部事实提交不能为了跑通而给 Agent 越权。
- 影响：15 个节点保留人类、Agent、审批、World wait 和系统事务的真实职责边界；Runtime 升级后旧批准失效并要求重新批准。
- 验证：`INC-INTERACTIVE-1784818405` 达到 `enterprise_operational_ready`；15/15 work item completed、15 journal、G1–G7 均实际批准、三个 Observation 存在；actor 覆盖 Founder、finance/governance/audit Agent、IAOS Runtime 和 World Bridge Runtime。
- 后续：T78 继续补人工接管/升级，T83 补浏览器多主体 E2E，T84–T85 收口未来里程碑合同与最终证据。

## 2026-07-23 - M9 Agent 工作项人工接管

- 变更：IAOS 增加仅针对已解锁 Agent task 的人工接管 API 和 UI；持久保存原 Agent、Founder、理由、接管时间并发布审计 Outbox。
- 原因：Agent 暂停、失效或需治理升级时必须可恢复，但 Founder 不能伪装为 Agent，也不能借接管绕过审批或 World wait。
- 影响：接管后工作项转为 human_task，正式 Capability Journal 的 actor 为 founder-principal；原分派证据不会覆盖。
- 验证：在线案 `INC-TAKEOVER-1784819709` 验证短理由 400、finance-agent → Founder 接管成功、能力提交成功；Playwright 验证独立待办菜单可见原 Agent 与接管理由。
- 后续：T83 仍需覆盖完整浏览器逐节点操作和尚未进入正常主线的 incorporation/legal Agent 参与。

## 2026-07-23 - 纠正 M10–M24 完成口径

- 变更：路线图把 M10–M24 改为 `Reference Replay Complete; D22 Pending`，保留既有确定性证据，但不再声称已完成交互式经营。
- 原因：预计算场景可用于回归和解释业务因果，却不能证明人类、Agent、审批与外部参与者逐节点工作。
- 影响：所有后续里程碑继承 DES-027 D22 完成门；只有持久工作项和真实主体 E2E 通过后才能恢复 Completed。
- 验证：路线图与唯一 active M9N 计划一致，T84 关闭。
- 后续：逐里程碑实施时复用 M9 工作项合同，不复制 M9 专用状态存储。

## 2026-07-24 - M9 案件、流程与证据入口可发现性

- 变更：IAOS 修复“企业设立案件/成立流程运行”菜单状态复用问题；案件编码增加最近记录搜索；案件、Process Run、Journal/Approval/World/Outbox 增加明细下钻。
- 原因：节点完成后用户只能看到计数，无法定位业务数据与事务证据；独立菜单复用组件时没有同步目标 workspace。
- 影响：节点 1 完成后可在“企业设立案件”查看状态和 Entity 数据，在“成立流程运行”查看 Process Run，在证据计数卡直接打开 Journal/Outbox 明细。
- 验证：TypeScript、production build 与 Playwright 菜单/搜索/Journal 下钻回归通过。
- 后续：继续把通用 JSON 输入逐步替换为 Capability Contract 生成的业务表单。

## 2026-07-24 - 隔离 M9 测试案件与业务案件列表

- 变更：最近设立案改为按更新时间倒序，默认排除 E2E、补正、拒绝、金额差异、职责分离等已知测试 fixture；保留 `include_test=true` 排障入口；案件工作区展示可点击业务案件卡片。
- 原因：原查询优先排列 terminal 案件且限制 20 条，大量集成测试数据遮蔽了用户刚创建的 `HCTM-TEST001`。
- 影响：不删除任何审计数据，但日常业务列表不再被测试 fixture 污染。
- 验证：`HCTM-TEST001` Trace API 200；默认最近案件返回 1 条且为该案件，测试 fixture 为 0；显式 include_test 仍可访问；Playwright 通过。
- 后续：长期应为 Case 增加显式 environment/data_classification 字段，替代历史前缀兼容过滤。

## 2026-07-24 - 区分 M9 案件列表搜索与治理案件选择

- 变更：IAOS “设立案件”补充独立列表内搜索，“企业成立与治理”案件编码改为点击即展开的搜索组合框。
- 原因：两个入口承担不同任务，原先把列表可发现性和单案选择混为一个原生 datalist，无法满足直接浏览和搜索。
- 影响：AESE 产生的 `HCTM-TEST001` 可在 IAOS 案件列表搜索，也可在治理工作台直接下拉选择并加载。
- 验证：IAOS TypeScript、生产构建/部署及双入口 Playwright 回归通过。
- 后续：M9 继续按交互式工作项合同验收逐节点操作。

## 2026-07-24 - 清理 M9 通用案件实体测试投影

- 变更：IAOS Runtime 1.3.6 重建通用 Entity 时排除测试 fixture 投影，并重新同步 `HCTM-TEST001`；治理案件选择器增加常显展开说明。
- 原因：此前只过滤最近案件 API，通用 `incorporation_case` Entity 仍保留 60 多条安装时测试投影且缺少后创建的业务案件。
- 影响：通用设立案件列表现在只有真实业务案件，canonical 测试审计事实未删除。
- 验证：在线 Entity API `total=1` 且唯一编码为 `HCTM-TEST001`；Runtime 1.3.6 apply、后端部署、前端生产部署和 Playwright 通过。
- 后续：将通用 Entity 投影从安装时重建演进为 Outbox 驱动的实时投影。

## 2026-07-24 - 修复治理案件下拉候选被二次过滤

- 变更：IAOS 治理案件选择器改为本地受控候选面板，点击固定展示全部业务案件，输入才过滤。
- 原因：第三方 AutoComplete 会按当前案号再次内部过滤，在部分浏览器交互顺序下产生空下拉。
- 影响：候选取数与展示逻辑可预测，`HCTM-TEST001` 可直接展开选择。
- 验证：TypeScript、生产构建/部署和点击展开/选择 Playwright 通过。
- 后续：无。

## 2026-07-24 - 案件下拉改为实时系统目录

- 变更：IAOS 案件下拉每次展开都重新合并 recent API 与通用 `incorporation_case` Entity 数据，并区分加载态和真实空态。
- 原因：页面首次挂载时的单次请求若发生认证时序失败，缓存数组会保持为空，之后点击无法恢复。
- 影响：下拉展示系统当前已有案件，而不是页面初始化快照。
- 验证：生产构建/部署及 `HCTM-TEST001` 展开、选择、加载 Playwright 通过。
- 后续：无。

## 2026-07-24 - M9 工作项改为业务表单弹窗

- 变更：D22 明确禁止业务用户编辑 JSON；IAOS 将 G1–G7、人工节点和 Agent 节点统一改为弹窗表单，金额输入元并自动转换为分。
- 原因：JSON 与 minor unit 属于 Runtime 合同，不是业务用户输入语言。
- 影响：用户从按钮进入表单，查看案件/节点/能力/主体，填写审批意见、金额或任务说明后确认提交。
- 验证：TypeScript、生产构建/部署；Playwright 验证 G1 弹窗、无 JSON 输入和必填审批意见。
- 后续：人工接管确认弹窗继续沿用同一表单框架。

## 2026-07-24 - 修复 M9 准备与校验节点被跳过

- 变更：Runtime 1.3.7 将工作项主线从 15 项修正为 18 项，补回 `founder.resolution.prepare`、`registration.package.validate`、`initial.budget.prepare`；安全重建仅完成开户节点的现有案件。
- 原因：工作项模板手写时只保留状态变更能力，与 Process Definition 不一致。
- 影响：新建设立案后步骤 2 由 incorporation-agent 准备创始人决议，步骤 3 才进入 G1 正式审批。
- 验证：Go 顺序锁测试；在线 Runtime 1.3.7 apply；`HCTM-TEST001` 前六项顺序与状态 API 核对；Playwright 通过。
- 后续：将工作项定义改为由 Process Definition 编译生成，消除双重维护。

## 2026-07-24 - Process Definition 成为工作项唯一事实源

- 变更：新增 ADR-005；IAOS Runtime 1.3.8 从一主四子 Process Definition 递归编译18项工作项，删除 API 手写顺序；登记补正改为 recovery 分支，移除子流程重复 readiness。
- 原因：顺序锁只能发现已知漂移，不能消除 Process 配置与 Runtime 双重维护。
- 影响：客户发布的流程配置决定实际工作项；未知、循环、重复、Gate/主体不一致或不可编译定义在产生业务事实前失败关闭。
- 验证：编译器、发布门、recovery/去重、API顺序测试通过；Runtime 1.3.8 在线安装；`HCTM-TEST001` 18项首尾和状态核对通过。
- 后续：把 process version/artifact hash 持久绑定扩展为所有通用 Process Run 的平台能力。
## 2026-07-24 - M9 节点 2 Agent Runtime 与证据链

- 变更：DES-027 新增 D23，明确 Agent 身份、岗位/Mandate、Capability/工具白名单、Runtime、派发、运行记录、草稿产出、人工审批和外部模型边界；IAOS Runtime 1.3.9 为 `founder.resolution.prepare` 增加独立派发与 `incorporation_agent_run`，并在 Agent 组织/成立审计展示完整配置和运行明细。
- 原因：原工作项虽标为 `agent_task`，实际仍是服务端直接执行 Capability，用户无法知道 Agent 做什么、在哪里定义、调用哪些工具以及外部 API 是否已经连接。
- 影响：节点 2 现在由 `incorporation-agent` 读取案件、形成治理决议草稿并提交受治理能力，然后严格停在节点 3/G1 人工审批；当前明确使用内置确定性 Runtime，不冒充外部 LLM。T83 的五 Agent、G1–G7、三个 World wait 和重启恢复全链仍未完成。
- 验证：IAOS scoped Go tests、TypeScript 编译和生产构建通过；在线案件已验证 Agent Run completed、3 次工具调用、治理决议草稿及步骤 2 completed/步骤 3 waiting_approval；生产后端/前端重部署成功，1440×900、1280×720、390×844 工作台回归通过。Atlas 声明已提交且 tracking check 通过，但线上同步入口对当前迁移/Founder 凭据分别返回 404/403，登记仍待具备 `system.atlas.manage` 的部署主体补同步。
- 后续：完成 T83/T85 全链交替参与和断点恢复验收，再关闭 M9N。
## 2026-07-24 - M9 输入合同与设立案 Entity 必填事实修复

- 变更：节点 2 从无规则“业务说明”升级为决议目标、核心提案、风险限制三字段业务合同，工作项保存并展示已发布 Capability 的完整 Input Contract；新建设立案增加案件名称、拟设企业名称、拟注册地址和拟经营范围表单，Runtime 同事务同步 canonical 案件和通用 Entity 投影。
- 原因：用户发现无法知道“输入合同”由谁定义及在哪里查看，同时数据模型工坊将 document_no/document_date/payload 声明为必填，但现有投影记录中字段为空。
- 影响：合同来源明确为随 Runtime Artifact 发布的 IAOS/领域 Capability Contract，不是用户临时上传；UI 与服务端执行一致的长度校验。设立案件的系统管理字段由 Runtime 生成，不再产生 required 字段为空的新记录，历史记录在 Runtime 升级投影时修复。
- 验证：scoped Go tests、TypeScript 编译与生产构建通过；后端/前端生产重部署和 Runtime 1.3.11 安装成功；真实 API 新建测试案件后，Entity 详情包含非空 document_no、document_date、status、payload 及四项初始资料；无效节点 2 输入返回 400；Playwright 合同与业务表单回归通过。
- 后续：继续完成 T83 五 Agent、G1–G7、三个 World wait 和重启恢复全链验收。

## 2026-07-24 - M9 审批路由与工作分发改造

- 变更：DES-027 新增 D24；IAOS 新增 DES-053、版本化审批方案、冻结 assignment、user/role/position/requester-selected/requester-manager 选择器、串行与 any/all/quorum 聚合、事项快照和分派通知；G1 从 founder 用户名硬编码改为 `position:chair`。
- 原因：原 DES-035 首片明确排除了复杂路由，`approver_role` 未参与授权，任何拥有 `approval.manage` 的人都可审批，审批中心只展示 UUID 和技术资源，不能满足企业权责审批。
- 影响：Process 只引用 flow key；提交时解析实际任职人并冻结。审批人能看到业务事项、发起人、路线和前序证据；非当前阶段受派人不能决策，岗位换人无需修改流程代码。
- 验证：IAOS Approval/Incorporation/API scoped Go tests、Go vet、前端 TypeScript/生产构建、Atlas/Code Map 检查通过；后端和前端已重部署。真实案件 `INC-APPROVAL-1784896345` 验证 G1 `position:chair` 解析、事项/显示名/路线详情、recipient Outbox、决定聚合、Gate consume 及节点 4 解锁。
- 后续：补浏览器三视口回归并继续 T83/T96 Founder 与五 Agent 全链。

## 2026-07-25 - Entity 生命周期与正式审批语义统一

- 变更：DES-027 新增 D25；IAOS Entity transition 区分 direct/approval，正式审批引用
  租户 Flow 与 Capability，首次动作创建并路由请求，批准后核验对象/动作/hash，再提交
  状态并 consume；数据模型工坊改为选择已有 active Flow。
- 原因：原 Entity `approval_flow` 直接改状态，而 DES-053 Approval Runtime 负责人员路由，
  两套机制并存会让客户配置“审批”后仍绕过审批人和审批中心。
- 影响：Entity 配置回答业务对象允许怎样变化，Approval Flow 回答由谁决定；roles 只控制
  发起，assignment 才赋予决定权。AESE 不复制该机制，只消费 IAOS committed evidence。
- 验证：IAOS compiler/approval/API Go tests、TypeScript、生产构建、Code Map 和 Atlas
  tracking 检查通过。
- 后续：部署在线版本并继续 T83/T96 Founder、五 Agent、G1–G7 与 World wait 全链验收。

## 2026-07-25 - 新功能可解释性与配置关系治理

- 变更：DES-027 新增 D26，要求 System Atlas 展示从语义默认生命周期到运行证据的细粒度
  关系；AGENTS.md 增加功能目的、配置、使用、关系和页面“功能说明”的强制交付规则；
  active M9N 计划增加并完成 T98。
- 原因：现有全景图只连接顶层工作台，页面也缺少设计目的和实施步骤，客户仍需询问开发者
  才能理解 Entity 审批、Approval Flow 和 Process 的职责与配置顺序。
- 影响：规则生效后的新功能必须同时交付权威设计、业务表单、页面说明和 Atlas 关系；
  旧页面按后续切片渐进补齐。IAOS 通用实现由 DES-060 承载。
- 验证：Markdown 链接、Atlas 声明、Code Map 与文档治理脚本。
- 后续：完成 M9 T83/T85/T96 全链验收，并为 AESE World 页面补同类功能说明。

## 2026-07-27 - AESE World 失效设立案深链恢复

- 变更：M9 lifecycle trace 返回 404 时不再终止 World；页面保留本地基线、显示结构化警告
  并列出 IAOS recent cases。E2E 改为动态选择真实案件，World Hub 测试更新为当前入口。
- 原因：旧深链固定引用已清理的 `INC-HCTM-001`，IAOS 404 被错误提升为整个 campaign 失败。
- 影响：用户能打开 World 并切换有效案件；401/403/5xx 仍失败关闭，正式进度不会由本地
  frame 冒充。
- 验证：Unit、TypeScript、生产构建与三视口 Playwright。
- 后续：M9 T83/T85/T96 继续使用动态真实案件完成全链验收。

## 2026-07-27 - 企业成立外部确认幂等修复

- 变更：登记、开户和任命外部确认按案件、payload type、result 生成稳定 transport identity；IAOS 同步增加业务事实级去重和案件绑定的可信存在性校验。
- 原因：旧 AESE 每次点击生成新幂等键，重复 Journal 又触发 IAOS `count != 1`，将已有可信登记结果反转为拒绝。
- 影响：重复点击不再增长 Journal；已有重复历史 Observation 的案件也能继续对应 world_wait 节点；其他案件证据不能冒用。
- 验证：IAOS 完整 lifecycle 与 World Bridge recovery 集成测试通过；AESE TypeScript 测试和生产构建。
- 后续：继续按三个 World wait 分别验收登记、开户与任命外部参与者交互。
## 2026-07-27 - Enterprise Genesis 游戏化企业创生体验设计

- 变更：新增 draft DES-028，将 M9 设计为玩家、五个数字员工、IAOS Runtime 与 World 共同推进的企业创生游戏；定义创业构想、AI 命名与 Logo、创始办公室、登记、银行资本、人才治理、开业准备七章，以及世界沙盘、经营桌面、治理证据三层界面、GameProjection、CreativeJob 和 GX0–GX5 交付切片。
- 原因：现有 M9 已能通过持久工作项、审批和 World wait 走通真实过程，但体验仍以流程工作台为主，无法让玩家直观感受“从一个想法创建一家企业”，也没有把生成式 AI 素材、数字员工和 2D/2.5D 世界统一起来。
- 影响：DES-028 保留 IAOS 的事实、权限和审计权威，把游戏画面限定为 committed projection；AI 名称、Logo 和文案先作为候选资产，人工选择并通过 Capability 后才生效。当前 M9N active plan 状态不变，DES-028 批准前不开始实现。
- 验证：对照现有 DES-027 D22–D25、M9 World/Agent/Approval 实现和 Capitalism Lab 官方公司创建、Logo、管理层、政策与 AI 经理玩法；使用 UI/UX 设计规则核对三视口、触摸、键盘、减弱动效、图表替代和性能边界。
- 后续：确认产品命名、美术风格、AI 素材 provider/许可、名称检查范围和首版行业边界，再将 DES-028 转为 Approved 并建立实施计划。

## 2026-07-27 - Enterprise Genesis GX0 开工

- 变更：DES-028 转为 Approved；新增 M9N 下的 active parallel subplan PLAN-GX-001，拆为 GX0–GX5、GXT1–GXT46；完成首个 `GameProjection` Go 合同、JSON Schema、示例 fixture 和严格校验测试。
- 原因：2.5D 画面必须先消费稳定、可追溯的只读投影，不能从现有 frame 或前端状态自行推断业务进度；同时当前 M9N T83/T96 尚在收口，新游戏实现需明确并行边界。
- 影响：游戏层已有 actor、work item、资金、品牌、World exchange、通知和 evidence ref 的机器合同；GX0/GX1 可继续并行，GX2–GX4 的正式推进仍依赖 M9N 交互式合同。
- 验证：GameProjection fixture 同时通过 JSON Schema 与 Go strict parser；非法时间倍率和缺失 evidence 的工作项失败关闭。
- 后续：实现 IncorporationTrace → GameProjection 编译器和只读 projection API，再建立 TypeScript data source 与三层界面线框。

## 2026-07-27 - Enterprise Genesis GX0 投影 API 与游戏壳

- 变更：完成 IncorporationTrace → GameProjection 确定性编译器、按 case/frame 的只读 API、前端 TypeScript data source 和 Enterprise Genesis 三层游戏壳；生命周期 M9 默认入口切换为游戏体验，原 M9 证据入口保留。
- 原因：游戏画面必须从后端 committed projection 恢复，不能由前端播放状态推断公司成立进度。
- 影响：桌面和移动端已能查看等距世界、章节、资金、工作项、Agent 和证据；每个工作项携带 evidence ref。当前场景图形为代码原生原型，PixiJS/AI 素材和真实 IAOS case 聚合仍在后续切片。
- 验证：Go gameprojection/httpapi/worldcontract tests、前端 TypeScript、17 files/42 tests 和 Vite production build 通过；构建提示主 chunk 703.68 kB，已登记为 GX5 动态拆包任务。
- 后续：完成 GX1 FounderIntent、命名、CreativeJob、BrandAsset 与品牌工作室。

## 2026-07-27 - Enterprise Genesis M9 游戏主线完成

- 变更：完成 live GameProjection、AI 企业身份工作室、`incorporation.case.open` 正式选择、PixiJS 2.5D/DOM fallback、原创企业城市素材、IAOS 工作项深链、World Observation 操作、五 Agent/18 工作项/证据视图和开业场景；IAOS Runtime 升级至 1.3.13，将登记材料校验分派法务合规 Agent。
- 原因：静态 frame 和查看型工作台不能证明玩家、五 Agent、G1–G7 与 World 真实交替推进，也无法提供从创业想法进入企业世界的可玩入口。
- 影响：`#enterprise-genesis?tenant=&case=` 读取 IAOS verified evidence bundle，不建立第二业务真相；AI 候选只有经现有 Capability 创建案件后才生效。M9N 与 PLAN-GX-001 均达到完成门，终态可移交 M10。
- 验证：AESE `go test ./...`；前端 17 files/42 tests、TypeScript、production build；mock 与 live Playwright 在 1440×900、1280×720、390×844 各 3/3 通过；live case `INC-WORK-ITEM-E2E-1785163319408212558` 为 18 completed、G1–G7、三个 World wait、六次 Agent Run/五个 distinct Agent、100% `enterprise_operational_ready`；IAOS 集成测试包含中途 Server 重建恢复。
- 后续：外部按需 LLM/Logo connector 需另行配置受治理 provider 与密钥；未配置时继续使用明确标记的 deterministic/文字几何 fallback。

## 2026-07-27 - Enterprise Genesis 玩家原生创建与操作闭环

- 变更：新案件不存在时进入企业创建态；增加“新建企业”、游戏内 Agent 派遣、G1–G7 审批决定、资本/实缴/预算输入、系统 Capability 和三个 World wait 操作卡；GameProjection 增加 Gate；新增从空白 case 完成 18 工作项的三视口 Playwright。
- 原因：旧页面以已完成案件投影和 IAOS 外链为主，玩家不能只通过游戏界面创建并经营自己的企业，无法验证 IAOS 的实际可操作性。
- 影响：玩家选择的名称、地址和经营范围通过 `incorporation.case.open` 成为 IAOS 事实；后续每次操作进入 IAOS Process、Approval、Agent Run、Journal 与 Outbox，前端仅在成功后重新读取 verified evidence。财务 Agent 100 万元授权上限由真实治理拒绝暴露并成为默认金额边界。
- 验证：Go gameprojection/httpapi test 与 vet；TypeScript、17 files/42 tests、production build；Playwright 在 1440×900、1280×720、390×844 从空白案件完成企业创建、18 工作项、G1–G7、三个 World wait 和 `enterprise_operational_ready`。
- 后续：MiniMax/Qwen 已配置为本地密钥，真实 LLM 名称与 Logo provider 仍需在 provider-neutral CreativeJob 后端接入；当前确定性 fallback 不影响企业创建闭环。

## 2026-07-28 - Enterprise Genesis 零起点与真实 AI 方案重构

- 变更：新增 ADR-006、DES-029 和 active PLAN-GXZ-001；将产品入口改为根主页，把 PlayerAccount、GenesisWorkspace、IAOS tenant provisioning、AESE World Run、MiniMax CreativeJob 和 incorporation case 拆成明确阶段；补充双租户隔离、失败恢复、权限和三视口验收门。
- 原因：代码审计确认当前 `DeterministicProvider` 只返回四组固定名称模板；“新建企业”也只在 `tenant-hctm-genesis` 内新增 case，既未调用真实 AI，也未创建独立租户和世界，无法满足真正从零创建企业的产品目标。
- 影响：`tenant-hctm-genesis` 降级为 demo fixture；新企业必须拥有独立 tenant、founder membership、M9 Runtime 和 World Run。浏览器不得调用平台管理员 tenant API；IAOS 需提供受限、幂等、可恢复的 Genesis provisioning saga。MiniMax 调用失败时 fallback 必须显式标记。
- 验证：核对 AESE creative provider 实现、IAOS SaaS Ops tenant lifecycle、Founder bootstrap 和 Runtime 安装边界；对照 MiniMax 官方 M2.7 Text/OpenAI-compatible/rate-limit 文档；DES-029 定义 Z0–Z5、Z1–Z37 和双租户完成门。
- 后续：从 Z0 合同与 provider health 开始；所有 IAOS 代码在独立 worktree 实现和提交。

## 2026-07-28 - Enterprise Genesis 根主页上线

- 变更：修正根路径默认落入旧订单沙盘的问题；新增响应式 Enterprise Genesis 产品主页，提供创建新企业、华辰样板世界、世界地图和 M9 四阶段入口；旧订单演示迁移为显式 `#sandbox` 样板入口。
- 原因：`App.tsx` 将初始模式硬编码为 `preview`，空 hash 未映射到独立路由，而且场景加载状态会抢占根页面，导致对外地址始终显示旧“交付承诺重算”Demo。
- 影响：`http://192.168.50.222:4173/` 现在是游戏主页；主页准确标注当前创建按钮进入 M9 交互开发版，独立租户通道仍按 PLAN-GXZ-001 建设，不把 HCTM fixture 冒充新租户。
- 验证：新增 Vitest 根路由回归和 Playwright 根主页/样板入口验收；前端生产构建通过。
- 后续：继续 Z25–Z30，将“创建新企业”从 M9 开发版切换到真实 GenesisWorkspace provisioning 与独立 tenant。

## 2026-07-28 - MiniMax M3 企业身份真实生成接入

- 变更：新增 MiniMax OpenAI-compatible provider、M3 配置、版本化命名 prompt、响应上限、超时、严格 JSON 与业务字段校验；创意 API 从硬编码 `DeterministicProvider` 切换为服务启动时注入的 provider，上游失败明确返回错误而不静默冒充 AI。
- 原因：企业身份工作室虽然标为 AI，但服务端每次都实例化固定模板 provider，且 `.env` 仍指向 M2.7；选择候选只在既有 tenant 内执行 `incorporation.case.open`，并不创建新 tenant。
- 影响：当前 8090 服务使用账户实际返回的模型 ID `MiniMax-M3`；真实 smoke 生成四个动态公司身份候选。tenant provisioning 仍是独立的前置 saga，不能由公司身份选择隐式完成。
- 验证：MiniMax provider 单元测试覆盖模型、鉴权、reasoning 分离、候选治理字段和 429 显式失败；真实 `/models` 与 completion smoke 通过；creative/httpapi/cmd 测试通过。
- 后续：完成 CreativeJob 持久证据、一次 JSON repair、provider 状态 API，以及 GenesisWorkspace → 独立 IAOS tenant → M9 case 全链。

## 2026-07-28 - Genesis Workspace 独立租户创建纵切

- 变更：新增根主页后的独立空间 onboarding、`/api/aese/v1/genesis/workspaces` BFF、幂等本地 Workspace store 和 IAOS provisioning client；服务端自动生成 workspace、tenant、World Run 和 case 标识，依次创建 tenant、Founder、tenant session、M9 Runtime 并激活后才进入 MiniMax M3 身份工作室。IAOS 新增 DES-062 生产多玩家控制面设计。
- 原因：玩家选择公司身份时才创建 tenant 会混淆隔离边界；由浏览器或 URL 指定 tenant 又会造成串租户风险。tenant 必须在创意身份阶段之前由平台服务端分配。
- 影响：当前本机游戏已经可以从根主页创建真正独立的 IAOS tenant，`tenant-hctm-genesis` 不再是新建主路径。loopback adapter 仍使用本地 dev service identity；正式多人部署必须迁移到 IAOS DES-062 的认证 subject、membership 和 session exchange。
- 验证：真实创建 `tenant-gx-efe42b90684620ce`，状态 active，绑定 `gxw-efe42b90684620ce`、独立 World Run/case、M9 20 项 Capability 资产和 tenant token；同幂等键重放返回同一 tenant。真实浏览器又创建 `tenant-gx-70de7d954e9dbb18` 并自动进入 AI 身份工作室。Go 全量测试、前端 5 项根路由测试、9 项三视口主页/onboarding 测试、1 项 live provisioning 和生产构建通过。
- 后续：实现 IAOS 三张 Genesis 控制面表、普通玩家授权 API、World committed checkpoint、生产 token exchange 和两玩家越权集成测试。

## 2026-07-28 - 修复新企业身份工作室双 502

- 变更：IAOS 新案件 evidence 404 在 AESE 保持 404；前端不再把任意 502 当作“尚未创建”；MiniMax M3 输出预算提升至 8192，非法/截断 JSON 允许一次严格重新生成。
- 原因：新 case 的预期 404 被网关误包装为 502；M3 reasoning 与四组候选超过 2048 token，导致 JSON 中途截断。
- 影响：新 Workspace 能稳定进入身份创建态；MiniMax 上游或 IAOS 真故障仍会显式展示，不会被伪装为空案件或固定候选。
- 验证：两条失败测试修复后通过；原始创业构想真实 M3 请求返回 200 和四个动态公司身份；新 case projection 返回 404。
- 后续：将 CreativeJob request ID、usage、finish reason 和 repair 次数持久化为可见证据。

## 2026-07-28 - 修复创建企业 Founder 会话 422

- 变更：Genesis provisioning 改为 Founder 普通登录后安装 M9 Runtime并返回 Founder token；新增 owner-scoped Workspace session refresh，前端仅对认证主体不匹配的特定 422 自动刷新并重放一次。
- 原因：原流程把 `/dev/switch-tenant` 的管理员 token 交给浏览器，但 `incorporation.case.open` 的 acting subject 固定为 `founder-principal`，IAOS 因认证主体不一致正确拒绝。
- 影响：新旧 Workspace 都能在不新建第二个 tenant、不放宽 IAOS 治理规则的前提下创建企业案件；其他玩家不能刷新该 Workspace 会话。
- 验证：原请求复现 422；同一 Workspace 刷新 Founder session 后案件创建返回 201 committed；AESE 全量 Go 测试、前端 5 项测试和生产构建通过。
- 后续：生产多人环境仍需按 IAOS DES-062 用真实 PlayerAccount/OIDC 替代本机玩家标识和 loopback adapter。

## 2026-07-28 - 增加游戏登录与我的企业大厅

- 变更：根入口增加游戏用户名登录；首次用户名认领浏览器已有 player ID；登录后在 Hero 下方列出 owner 拥有的 Workspace，并可刷新 Founder session 后继续原 case，或创建新企业。
- 原因：原主页只有创建入口，玩家无法发现和恢复之前创建的独立租户，每次只能从零开始。
- 影响：本机玩家拥有明确的“登录 → 我的企业 → 继续游戏/创建新企业”路径；用户名切换使用独立本地 player ID。界面明确说明它不是正式多人认证。
- 验证：前端单元测试覆盖首次登录迁移和登录前置；生产构建与桌面 Playwright 主页、样板、新建路径通过；Workspace list/session 继续复用 owner-scoped 后端测试。
- 后续：按 IAOS DES-062 将本机用户名映射替换为 PlayerAccount/OIDC、服务端 membership 和跨设备会话。

## 2026-07-28 - 创始人办公室 RPG 首章

- 变更：新企业入口从 BrandStudio 表单替换为 PixiJS 创始办公室；增加四个 FounderProfile 头像、数字员工“纪元”、七段主线任务和 RPG 对话选择，在故事中完成产业、客户、产品、品牌性格、MiniMax 身份提案及 IAOS 设立提交；其余 17 个节点增加董事会、政务、银行、人才与经营会议的地点、NPC 和剧情简报。
- 原因：原操作只是 IAOS Work Item 的视觉换皮，玩家缺少角色、场景、任务目标、对话引导和决策反馈，不知道正在扮演谁或为什么操作。
- 影响：玩家先以创始人身份进入世界，通过经营选择形成 command draft，最终签署才写入 IAOS；技术 Capability 留在治理证据层。通用 WorkItem 面板暂保留为后续节点 fallback。
- 验证：Founder Office 桌面 Playwright 完整走过头像、四轮对话、AI 名称、地址/范围和 IAOS commit；前端 7 项单元测试、TypeScript 与生产构建通过。
- 后续：按 D26 将剩余 17 个 M9 节点依次改造成董事会、政务大厅、银行谈判、CEO 会面和经营会议事件。

## 2026-07-28 - 设立案件与正式企业主数据补齐

- 变更：IAOS `incorporation_case` 增加案件名称、拟设企业名称、注册地址和经营范围显式列并按 tenant/RLS 回填历史数据；登记成功时创建正式 `m9_legal_entity` 主数据。
- 原因：第一章输入原先只在 JSONB 快照中，无法直接检索；登记后也没有从拟设资料转成正式企业主体。
- 影响：拟设案件与正式法律主体职责分离；游戏创建信息可直接查询，登记完成后成为可被后续银行、组织和经营流程引用的企业主数据。
- 验证：IAOS 单元测试、完整 18 节点 integration tracer、线上部署通过；三个现存 Genesis 案件已真实回填，正式 tracer 法律主体字段一致。
- 后续：玩家继续当前案件；到登记成功节点后检查其租户内 `m9_legal_entity` active 记录。

## 2026-07-28 - M9 起草结果与审批对象可视化

- 变更：节点 2 明确说明三项对话输入会形成 IAOS 持久化《创始人设立决议草案》；GameProjection 开始读取 Agent Run output，G1 展示原始决议，G2–G7 统一展示待审议文件、关键内容、风险限制、起草人、批准效果和证据引用。
- 原因：原界面丢弃 Agent 输出，只显示通用 Capability 说明，玩家不知道第二步产出了什么，也会在看不到审批对象时直接批准。
- 影响：M9 形成“输入 → 起草文件 → 审阅 → 批准执行”的业务链；审批投影来自 IAOS committed evidence，缺少审阅对象时前端禁止盲目批准。
- 验证：GameProjection 新增 G1 原始 Agent output 与 G1–G7 审阅对象测试；相关 Go tests、18 个前端测试文件/46 项测试、TypeScript 和生产构建通过。
- 后续：把当前治理审阅卡继续融入董事会、政务、银行、人才和经营会议的专用 RPG 场景。

## 2026-07-28 - 登记与银行开户可失败经营事件

- 变更：登记提交增加五项申请资料和审查重点；外部登记支持缺件退回、补正重申及营业执照/三枚印章领取。银行开户增加三家虚拟银行选择、五项尽调资料、拒绝原因、补件重申及基本账户/U 盾领取。
- 原因：原 `registration.submit` 和两个 observation 节点是一键动作，玩家不知道提交内容，也看不到外部反馈、失败恢复或成功后的企业资产。
- 影响：缺件会先写入 `rejected` World Observation，IAOS 案件不推进；补齐后以不同结果幂等重申。证照与账户资产明确标注为虚构沙盘内容，不冒充真实证件。
- 验证：前端 18 个测试文件/46 项测试、TypeScript 与生产构建通过；图像生成服务网络失败，当前交付使用可替换的代码原生虚构凭证卡。
- 后续：图像服务恢复后生成并替换营业执照、公章、账户通知和 U 盾 raster 素材；增加浏览器拒绝—补正—获批专项回归。

## 2026-07-28 - 修复银行开户 Execute 400

- 变更：审批提交输入与 Work Item Execute 输入分离；开户银行和资料说明保留在 G3 intent，严格执行命令只携带 correlation 及受支持的金额字段；重试时若幂等审批已 approved，则跳过重复决定直接恢复 Execute。
- 原因：前端把审批专用 `business_note` 原样传给 IAOS `incorporation.Command`，严格 JSON 解码正确返回 400。
- 影响：当前案件已成功完成的 G3 审批可以直接重试第 8 步，不需重建企业或重新审批；IAOS 严格合同保持不变。
- 验证：真实日志确认 G3 submit/approve 成功、Execute 400，以及重试重复 approve 403；API 回归测试锁定请求体边界和已批准恢复路径，TypeScript 和生产构建通过。
- 后续：继续测试银行外部反馈与开户资产领取。

## 2026-07-28 - 实缴资本差异核对与纠正审批

- 变更：G4 增加认缴、本次到账、差额三栏核对；不一致时禁止提交并可一键修正；金额加入审批 correlation，纠正后的金额必须获得新的 G4 批准。
- 原因：当前案件认缴 100 万、玩家输入实缴 80 万，IAOS 正确提交差异证据并返回 `capital_contribution_mismatch`，但游戏提交前没有清楚暴露该约束。
- 影响：玩家能在批准前理解并纠正资本差异；旧 80 万审批不会被用于 100 万实缴，审批内容与执行事实继续一致。
- 验证：数据库核对认缴 100,000,000 分、拒绝实缴 80,000,000 分及 approved G4 intent；API 测试覆盖金额相关 correlation，TypeScript 和生产构建通过。
- 后续：继续当前案件，按认缴金额 100 万发起新的 G4 审批并完成核验。

## 2026-07-28 - Enterprise Genesis 世界优先场景改造

- 变更：DES-028 新增 D29；城市写实背景中的四个实际区域改为创始办公室、政务中心、合作银行和企业总部热点，点击进入各自室内地点；主线改由建筑事件和 NPC 场景触发，Work Item 降级为治理档案；章节改为只读旅程，移除上一步、后续 Frame 和无 World Tick 的暂停/倍率控件。
- 原因：原通用蓝色房子与城市背景无空间关系，点击只改变文字；技术任务列表和伪时间控件主导体验，使 M9 仍像流程工作台。
- 影响：地点具有独立图形语言、室内家具/NPC/资料台和 committed 资产反馈；登记后办公室显示执照/印章，开户后银行显示账户/U 盾，组织建立后总部显示管理层工位。移动端保留可触摸地点入口。
- 验证：GameProjection/HTTP Go tests、18 个前端测试文件/49 项测试、TypeScript 与生产构建通过；世界优先 Playwright 在 1440×900、1280×720、390×844 验证城市热点、室内进入/返回、三标签和伪控件移除，3/3 通过。
- 后续：为四个室内场景补版本化 raster/atlas 素材、人物行走和地点间转场；M10 开始后再接入真实 World Tick 与时间倍率。

## 2026-07-28 - 城市旅行、室内 NPC 与永久资产反馈

- 变更：DES-028 新增 D30；当前事件在城市中绘制基地到目的地的动态路线，进入地点播放可跳过旅行转场；四类室内场景增加具名 NPC、岗位、工作状态和对白；创始办公室增加执照、印章、账户和团队 trophy 柜。
- 原因：仅能点击地点仍缺少从城市到室内的空间连续性、人物生命感和完成事项后的长期奖励。
- 影响：玩家能够看到“去哪里、见谁、对方在做什么”，刷新后 trophy 仍由 committed Projection 重建；转场和动效不推进 IAOS 事实，减弱动效模式自动停用循环动画。
- 验证：前端 18 个测试文件/49 项测试、TypeScript、生产构建通过；三视口 Playwright 验证旅行状态、室内 NPC、进入/返回和伪控件缺失，3/3 通过。
- 后续：图像生成服务恢复后，用版本化 raster atlas 替换代码原生人物和家具，同时保持 DOM 低性能 fallback。

## 2026-07-28 - 玩家地图化身与场景物件检查

- 变更：DES-028 新增 D31；顶部 FounderProfile 延伸为城市中的玩家化身，旅行时移动到目标地区并支持减弱动效；四类室内地点新增可检查物件和详情卡，展示用途、状态、解锁条件及 IAOS/World committed 来源。
- 原因：旅行转场仍缺少玩家在地图中的持续存在，室内家具也只是背景，不能回答“这个物件现在是什么状态、从哪里来的”。
- 影响：玩家位置与视角位置一致，检查物件不会推进 Capability；正式动作继续从 NPC 当前事件进入，避免再产生一套场景表单。
- 验证：前端 18 个测试文件/49 项测试、TypeScript 和生产构建通过；三视口 Playwright 验证玩家位置、旅行、NPC、物件详情和返回地图，3/3 通过。
- 后续：将角色和物件 fallback 接入版本化 atlas manifest，并为 M10 工厂区域扩展相同地点交互合同。

## 2026-07-28 - OpenAI 生成四类 2.5D 室内场景

- 变更：使用 OpenAI 图片生成功能创建创始办公室、政务服务中心、合作银行和企业总部四张统一风格 2.5D 场景；接入地点背景，并把尺寸、SHA-256、来源和许可写入素材 manifest。
- 原因：代码原生家具只能表达交互结构，缺少经营游戏需要的空间质感、建筑身份和氛围连续性。
- 影响：玩家进入不同建筑后可立即识别地点；NPC、任务和 committed 物件状态继续由 DOM/IAOS 投影驱动，生成图片不成为业务事实。
- 验证：18 个前端测试文件/49 项测试和生产构建通过；三视口 Playwright 3/3 通过；素材 manifest/Atlas JSON、Markdown 链接和 System Atlas tracking 检查通过。
- 后续：生成角色、证照、印章、U 盾等透明 atlas，并为 M10 工厂场景延续同一美术规范。

## 2026-07-28 - 修复室内介绍文字横幅溢出

- 变更：为地点标题、剧情横幅、NPC 对话和场景说明卡补充 border-box、可收缩列、自动换行和最大宽度约束。
- 原因：介绍文字长度变化时，Grid/Flex 子元素的默认最小内容宽度可能超过自身半透明背景。
- 影响：桌面和移动端的介绍文字保持在对应背景条或对话卡内，不改变任务与 IAOS 交互。
- 验证：Playwright 增加房间横幅边界和文字容器无水平溢出断言；前端测试与构建通过。
- 后续：后续新增场景文案继续复用相同文字容器约束。

补充修正：用户指出目标是政务与银行室内的柜台标牌；已将 `.gx-counter` 改为居中、
响应式固定边界，增加最小高度和水平内边距，使“企业登记综合窗口/受理·审查·补正·发照”
及“企业金融服务柜台/开户·尽调·网银·资本金”完整包含在标牌背景内。

## 2026-07-28 - M9 全 RPG、生成精灵与多租户验收

- 变更：17 个工作项增加专属剧情与经营选项；生成并接入五名角色及执照/印章/U 盾透明精灵；增加室内移动、音效开关、里程碑奖励和逐项企业大事记；workspace 失败重试复用原 identity；IAOS 独立 worktree 修复 Founder 多租户发现。
- 原因：通用任务面板、CSS 人物和只有数量的大事记仍不足以形成经营游戏闭环；历史 workspace 密码差异导致 IAOS 登录只能看到一个租户。
- 影响：玩家从剧情进入正式治理动作，完成事项获得可持续反馈；同一 Founder 可选择所有显式授权企业，跨租户 M9 数据仍隔离；失败重试不再创建孤儿 tenant。
- 验证：AESE 49 项前端测试、生产构建、三视口 Playwright 3/3、workspace/HTTP Go tests；IAOS API Go tests、前端 TypeScript/production build；真实登录返回 8 个授权租户，抽查两租户 profile claim 一致且交叉案件读取均 404。
- 后续：实现 IAOS 生产级 `/genesis/workspaces` saga 与持久 CreativeJob；当前 AESE 本机 store 仍是 loopback adapter。

## 2026-07-28 - 校准室内人物与房间尺度

- 变更：创始人与四类地点 NPC 改用响应式室内尺寸，按家具尺度统一身高、脚底缩放原点和接触阴影；移动端单独约束尺寸、位置和热点移动距离。
- 原因：1672×941 房间背景中原玩家仅 64×116、NPC 仅 92×180，透明边距进一步缩小可见身形，人物与椅子、门和室内景深明显失调。
- 影响：常用桌面视口中人物约 200–230px 高，玩家与同景深 NPC 比例小于 1.3；移动端人物保持可辨识且不会因放大越界。
- 验证：TypeScript typecheck 通过；Enterprise Genesis Playwright 在 1440×900、1280×720、390×844 三视口 3/3 通过，并增加最小身高、玩家/NPC 比例和房间截图证据。
- 后续：新增 M10 工厂室内背景时沿用 D32 的家具标尺与三视口比例验收。

## 2026-07-28 - 修复组织节点推进剧情无响应

- 变更：把 `capital_contribution_verified` 的场景章节从银行资本切换为人才治理，使下一项 `organization.establish` 出现时同步开放企业总部；补充当前任务目的地必须开放的设计约束。
- 原因：任务已指向企业总部，但状态章节仍只开放到合作银行；“推进剧情”调用受锁地点导航后被静默返回，表现为按钮无反应。
- 影响：资本核验完成后可正常进入企业总部、打开建立初始组织的受治理操作面板；不改变 IAOS Capability 执行或 committed 状态。
- 验证：先由 Go 回归测试稳定复现“organization mission points to a locked headquarters”，修复后 `internal/gameprojection` 通过；Playwright 验证城市推进剧情、到达总部、室内推进剧情及操作面板打开通过；AESE 后端已重建并在 8090 健康运行。
- 后续：其余章节边界继续遵循“首个可执行 Work Item 的目的地必须开放”不变量。

## 2026-07-28 - 批准 M9–M13 制造企业财务运行体系

- 变更：新增 approved DES-030 和 planned PLAN-M9-FIN-001；将财务组织、账套、科目、会计事件、凭证、总账、AP/AR、资金、固定资产、制造成本、结账、报表和 KPI 纳入 Project Genesis，并更新 DES-027/028、Architecture、Roadmap、Code Map、Gap Ledger 和文档索引。
- 原因：现有 M9–M13 只有现金、预算、承诺、应付、应收、回款、实际成本和项目毛利等经营台账，银行注资没有形成实收资本凭证和开业财务报表，不能满足制造企业业财一体化要求。
- 影响：M9 新增财务组织与七个开业会计工作项，以 `finance_opening_ready` 作为新版本开业条件；M10–M13 按真实工程、采购、资产、生产和销售事实逐步启用财务子账，避免在 M9 虚构未发生交易。
- 验证：设计覆盖财政部企业会计准则、预算/成本管理会计应用指引和 COSO 内控基线；Markdown 链接、计划状态、Atlas JSON 和文档导航检查通过。
- 后续：在独立 IAOS worktree 执行 F0 合同审计和历史案件迁移设计；所有 F1–F35 未实现项保持未勾选。

## 2026-07-28 - M9 开业财务第一纵切

- 变更：IAOS 新增六个财务 Entity、五项 Capability、一个开业财务 Process 和四项 Policy；`capital.contribution.verify` 同事务建立财务组织、CAS-BE/CNY 账套、开放期间、1002/4001 科目及借贷平衡凭证。AESE Projection、企业总部和治理档案接入财务组织、凭证与试算余额。
- 原因：公司现金不是会计账，实缴资本必须形成可审计、可对账的开业凭证，后续费用、资产和制造成本才能有合法起点。
- 影响：原 18 个玩家/Agent/审批/World 工作项顺序保持兼容；财务子步骤作为资本核验的原子编排展开，不重复推进案件状态。数据模型工坊、能力语义工作室和流程编排控制台可发现新增资产。
- 验证：IAOS 离线合同/语义/流程测试通过；PostgreSQL 完整 M9 tracer 验证一组织、一账套、一凭证、两分录及借贷均为 100 万元；AESE Go tests、49 项前端测试和 TypeScript 通过。
- 后续：完成历史案件受治理回填、显式财务工作项/职责 Mandate、开业资产负债表，并进入 M10/M11 AP/资产纵切。

## 2026-07-28 - M9 财务就绪硬门与开业报表

- 变更：IAOS 将财务组织、账套、开放期间、已过账资本凭证和借贷平衡纳入企业运营就绪硬门；开业财务 API 增加银行日记账、总账、试算平衡与资产负债表。AESE 治理档案同步展示这些 committed 视图。
- 原因：仅创建凭证但不参与最终就绪判断，仍可能出现“案件已完成、财务未开业”；只有试算余额也不足以让业务用户理解银行余额、科目余额与开业财务状况。
- 影响：缺失、草稿、零金额或借贷不平的开业凭证会 fail-closed；历史真实案件按租户 RLS 幂等补建。所有报表从同一已过账分录派生，不生成虚构采购、销售或利润。
- 验证：单元测试覆盖缺凭证、零金额与借贷不平；PostgreSQL 全生命周期测试先破坏凭证并确认 422 和状态不推进，再恢复并验证银行日记账 1 条、总账 2 条、资产 100 万元、权益 100 万元且报表平衡。
- 后续：发布两个显式财务 Process，配置 CFO/Controller/GL/出纳/成本/审计 Mandate 与职责冲突，并把七个财务工作项接入持久流程和通知。

补充：IAOS artifact 已升级到 1.5.0，并发布 `finance.foundation.setup.v1` 与
`capital.accounting.v1`；五个现存 `tenant-gx-*` 企业已幂等迁移到 1.5.0，两个流程均为
published。原 `finance.opening.foundation.v1` 只作为旧版本兼容编排保留。

## 2026-07-29 - 总部财务穿透与 Entity 菜单发布恢复

- 变更：修复总部开业财务中心与治理会议桌重叠，把财务中心改为可访问交互对象；治理档案和对象详情增加 IAOS 系统账务/财务报表穿透。IAOS 新增财务账务与报表工作台，并恢复 Entity 发布后的即时菜单刷新和 platform_super_admin 菜单可见性。
- 原因：AESE 财务对象只有装饰没有用途，财务明细继续堆在治理档案无法扩展；动态 Entity 菜单虽然成功生成，却被未识别的平台管理员角色和未刷新 Sidebar 隐藏。
- 影响：AESE 保持经营摘要，长期账务和报表回到 IAOS 权威入口；普通角色仍需显式 menu READ 授权。M9 当前开业切片与未来完整报表系统边界在 DES-030/DES-063 中明确。
- 验证：AESE 18 个测试文件/49 项测试、TypeScript 和总部三视口 Playwright 3/3；IAOS auth/api Go tests 与 Next.js production build。
- 后续：完成财务职责矩阵、七个显式财务工作项、月结/三表/管理 KPI 与全金额证据穿透。

## 2026-07-29 - 会计凭证聚合、多币种与财务数据菜单纠偏

- 变更：DES-030 明确会计凭证必须采用 `document_with_lines`/`document_line` 主子聚合；新增币种定义、日期化汇率、交易币金额、本位币金额与不可变汇率快照合同；要求全部 M9 已发布 Entity 同步发布左侧数据菜单。IAOS 独立 worktree 已实现对应 Entity、RLS 表、历史 CNY 回填和“财务管理”菜单分组。
- 原因：原设计把凭证头和凭证明细都注册成普通 `document`，父子关系仅停留在语义图；同时缺少汇率主数据及运行数据菜单，无法支撑可审计的多币种凭证或让用户发现已创建数据。
- 影响：凭证可在同一聚合中表达任意多行借贷，明细明确保存交易币、本位币和采用汇率；M9 仍限定 CNY，不会未经外币流程治理直接放开资本到账。模型工坊与运行数据入口职责分离。
- 验证：IAOS incorporation/api 定向测试、全平台 Go tests、go vet 和前端 TypeScript 通过；运行安装器后 19/19 个 M9 Entity 均有菜单，币种与汇率各仅保留 1 条有效 CNY 基础记录；DES-030 与 IAOS DES-063 合同一致。
- 后续：完成真实外币银行 Observation、汇兑损益、期末重估和多币种报表折算流程。

## 2026-07-29 - 恢复 M9 Entity 运行数据与财务投影

- 变更：IAOS 修复人工 `INC-INTERACTIVE-*` 案件被夹具清理、启动回填 prepared statement 失败和模型再次发布后 `m9_*` 物理表绑定漂移；DES-030 补充权威账务、通用 Entity 投影与菜单读取的一致性合同。
- 原因：当前租户的财务工作台已有开业语义，但财务 Entity 菜单为空；权威数据被清理、回填整体回滚以及元数据指向空 `bo_*` 表共同造成“模型存在但数据不可见”。
- 影响：实际操作租户 `tenant-gx-f4b3ce3ce8e2712d` 已升级到 M9 Runtime 1.6.0，并从已提交实缴资本事实恢复 1 个财务组织、1 个账套、2 个总账科目、1 个期间、1 张凭证和 2 条分录；19 个 M9 Entity 菜单均能读取该租户自己的投影，模型重新发布不再丢失数据绑定。
- 验证：IAOS 全平台 Go tests、定向 vet、Atlas 检查通过；目标租户安装写入/修复 485 项平台资产；API 验证八个财务 Entity 全部绑定 `m9_*` 且数量为 1/1/2/1/1/1/1/2，开业 API 借贷各 80000000、ready=true；1002 银行存款与 4001 实收资本字段完整，凭证明细分别引用这两个科目。
- 后续：后续 M10–M13 Entity 只在真实业务事实发生后生成数据；实施验收继续区分“业务尚未发生”与“投影/绑定故障”。

## 2026-07-29 - M9 投影迁移为里程碑无关存储

- 变更：IAOS 新增 `entityprojection` 深模块，将 19 张共享 `m9_*` 投影原位迁移为 `entity_projection_*`；Runtime、设立与财务写入统一切换，逐租户修复 118 个 metadata 版本并规范化 RLS policy。
- 原因：M9 是交付阶段，不应成为法律主体、账套、凭证等长期 Entity 的物理命名空间；继续扩散 `m10_*` 会形成平行数据真相。
- 影响：权威 `finance_*`/`incorporation_*` 表不变；全部 UUID、租户、稳定业务编码和业务行保留，M10–M13 复用稳定 Entity code 和同一规范投影。
- 验证：迁移前后 19 张表的行数及 `id:tenant:business_code` 哈希逐项一致；每表 RLS/FORCE RLS 与单一策略有效；目标租户 19 个 records API、80 万元借贷平衡和 Runtime 1.6.0 no-op 通过。
- 后续：M10 新增对象必须遵守 IAOS DES-064，不得创建里程碑前缀投影；继续实施 AP、资产、成本、AR 和结账子账。
## 2026-07-29 - 修复 IAOS 会计凭证主从明细可见性

- 变更：IAOS 修复 `journal_entry → journal_line` 投影父键，启动时幂等回填历史关联；通用业务数据浏览器将当前 Metadata 的 `child_list` 合并进旧 Formily 布局，并明确显示“查看主表及明细”。
- 原因：目标 Genesis 租户已有一张凭证和两条分录，但分录没有引用 Entity 投影凭证头，旧 UI 布局又遗漏 `lines`，导致菜单只能看到主表。
- 影响：AESE M9 银行注资形成的会计凭证可在 IAOS 左侧“会计凭证”菜单穿透查看两条借贷分录；AESE 不复制账务数据或主从渲染逻辑。
- 验证：目标租户 `tenant-gx-f4b3ce3ce8e2712d` 查询得到 `line_count=2`、`all_linked=true`；IAOS 全量 Go、Go vet、前端单测、TypeScript、生产构建和 PostgreSQL 生命周期集成测试通过。
- 后续：M10 及后续业务记账继续复用 `document_with_lines` 主从合同，验收必须检查详情 API 的 `children.lines`，不能只检查头、行分别存在。

## 2026-07-29 - 建立 Capability 唯一业务写入边界

- 变更：IAOS 新增可执行 Implementation Binding、tenant-RLS Capability Execution、事务身份和财务凭证数据库触发器；资本核验改为执行 `capital.contribution.post`，案件 Trace 暴露真实执行，CI 拒绝凭证直写。AESE M9 主计划增加 D24 约束与双仓责任边界。
- 原因：仅在 Catalog、DSL 和 Process 中登记 Capability，不能阻止 Handler 直接调用内部方法，执行、权限、审批、审计和幂等链会因此失真。
- 影响：M9 实缴资本记账现在只能在已发布、已绑定的 Capability Context 中完成；历史修复也必须留下 `system_repair` execution。首个数据库强制范围是凭证头与分录，AESE 继续只消费 IAOS committed 结果。
- 验证：IAOS `main@887c39b` 已部署并推送；完整 PostgreSQL M9 生命周期验证 1 次 execution、1 张凭证、2 条平衡分录；目标租户已有 succeeded execution；无 Context 的裸凭证 UPDATE 被 `governed_write_context_required` 拒绝；Go API、vet、静态治理、Code Map 与 Atlas 检查通过。
- 后续：按 IAOS DES-065 顺序迁移凭证审批/冲销、订单库存、AP/AR、固定资产、制造成本、审批最终决定和 Agent 可写 Tool；共享开发库的既有 Outbox/全量集成测试隔离债务另行治理。

## 2026-07-29 - M9 全量执行绑定、显式财务节点与 Outbox 修复

- 变更：IAOS 将 25 个 M9 Capability 全部绑定到统一执行边界；主流程在资本核验后新增财务组织、账套与期间、科目、资本过账、财务就绪 5 个显式节点，正常路径从 18 项升级为 23 项；Outbox 以行租户和事件类型包装业务 payload。
- 原因：只治理资本过账仍允许其余 M9 命令缺少统一 execution；把完整财务开业隐藏在资本核验内部不符合可解释、可人工/Agent 协作的流程原则；旧 Poller 会把普通业务 JSON 发布为 `iaos..` 空 subject。
- 影响：人工、Agent、Process 和 World Observation 均进入同一 Capability Interface；用户能逐项观察和推进财务开业；Outbox 路由恢复为合法、可重放的完整 Event。AESE 不复制 IAOS 执行器和财务数据库。
- 验证：IAOS API/incorporation/outbox 单元测试通过；真实 PostgreSQL 交互测试以 23 个工作项、7 个审批门、3 个 World wait、6 次 Agent Run/5 个 Agent 到达 `enterprise_operational_ready`；Outbox 测试覆盖 payload 包装、envelope 保留及路由不一致拒绝。
- 后续：部署后只重放活跃业务租户的历史失败 Outbox；继续完成 M9-FIN 职责矩阵，以及 M10–M13 AP/AR、资产、成本和结账子账。

## 2026-07-29 - M9 Runtime 1.7.0 部署与证据校准

- 变更：IAOS `main@a1cc21e` 部署后，为目标 GX 租户安装 Runtime 1.7.0；Artifact 新增 Capability 合同与 Process 定义联合签名。同步更新 M9 最终证据、风险登记，并把财务 F13 拆为显式流程节点与尚未完成的人机协作两项。
- 原因：旧 Artifact hash 只覆盖资产名称，流程图实际变化仍会被安装器误判为 no-op；旧证据也没有区分历史全量验收与本次定向复验。
- 影响：目标租户已获得 25 个 active 执行绑定和新版财务子流程；新建设立案使用 23 个持久工作项，历史已完成案件保持原事实不被改写。文档不再把通知、Mandate 和游戏交互误标为已完成。
- 验证：目标租户 Runtime 1.7.0 为 active、1.6.0 为 superseded；25 个 M9 binding 均 active；目标 Outbox 43/43 `PROCESSED`；PostgreSQL tracer 以 23 工作项、7 门、3 World wait、6 Agent Run/5 Agent 到达运营就绪；IAOS Go、vet、治理检查、Code Map 和 Atlas 检查通过。
- 后续：完成 F10/F13B 财务职责、Mandate、通知和游戏交互；隔离全量 PostgreSQL integration 测试数据库，并按 DES-065 继续扩展非 M9 业务写入治理。

## 2026-07-29 - M9 财务责任岗位、通知与游戏协作闭环

- 变更：IAOS Runtime 1.8.0 安装六个财务责任岗位、两个现任 Agent Mandate 和四条阻断 SoD；五个财务工作项保存执行者、责任岗位、通知对象及空缺升级，并在激活时写幂等 Outbox。财务工作台新增“组织与待办”。AESE 将五个财务节点路由到企业总部，补齐逐项操作说明，并让开业财务中心在建设期即可穿透组织、Mandate、SoD 和待办。
- 原因：只有财务 Capability 和自动执行节点仍无法回答“谁负责、谁能做、岗位空缺怎么办”，游戏中的财务中心也不能只在全部完成后才成为可点击装饰。
- 影响：玩家逐项触发财务组织、账套、科目、资本过账和就绪检查；`iaos-runtime` 仍是确定性执行者但不冒充会计责任人。Controller、总账、资金和成本岗位默认明确为空缺，通知 Founder 补位但不伪造任职。
- 验证：IAOS 主线 `ff259a8` 已推送和部署；目标租户安装 Runtime 1.8.0，API 返回 6 岗位、4 SoD、5 财务节点，重复安装 no-op；资金与总账同主体兼任被数据库阻断且回滚无残留。AESE 前端测试与生产构建通过。
- 后续：M9 财务 F10/F13B 关闭；继续 F15–F35，按 M10–M13 真实工程、资产、采购、生产和销售事实建设 AP/AR、固定资产、制造成本、结账和完整报表。

## 2026-07-29 - M9 零起点生产控制面与财务基线收口

- 变更：IAOS 独立 worktree 实现 owner-scoped Genesis Workspace/member/step 控制面、八步 provisioning、World evidence 激活门和 tenant-only `genesis_owner` 会话；M9 Runtime 安装与人工工作项移除新租户对 `platform_super_admin`、固定 `founder-principal` 和 HCTM 组织编码的依赖。AESE BFF 转发 Player session，五步创建向导展示八 checkpoint 证据；CreativeJob 持久保存 provider/model/host/request/token/latency/validation/fallback，并以真实 IAOS profile 阻断跨租户读取。新增 M9–M13 财务对象清单、五项 M9 财务 Capability 的权限/额度/敏感度/SoD/补偿矩阵及离线校验器。
- 原因：旧 loopback adapter 会以固定 Founder 和平台管理员路径创建租户，无法作为多人交付边界；M9 财务 F1/F5 也缺少机器可读盘点和真实数据证据，不能准确区分 M9 已完成对象与 M10–M13 规划对象。
- 影响：普通认证 Player 能在不持有平台权限、密码或 tenant URL 决策权的情况下创建隔离企业；新案件节点 1 使用真实 Player subject，节点 2 交给设立 Agent。M9 财务范围可以独立判定完成，AP/AR/资产/成本/月结仍明确属于后续里程碑。
- 验证：AESE Go 全量测试、前端 20 文件/56 项测试和生产构建通过；IAOS Go 全量测试通过。真实 PostgreSQL 验证同幂等键 no-op、八 checkpoint、另一 Player 读取 404、JWT 仅 `genesis_owner`、新租户 23 个 M9 工作项及动态 owner participant；目标财务租户验证 1 张凭证/2 条分录、11 类 M9 财务投影计数和 `m9_%` 物理表为 0。
- 后续：完成 PLAN-GXZ-001 Z32——从本次全新 Workspace 逐项跑完 23 个 M9 节点并归档最终 evidence；随后进入 M10，不把 F15–F35 提前算作 M9 欠项。

## 2026-07-29 - M9 全新 Workspace 端到端验收完成

- 变更：通过 IAOS 生产 Genesis Workspace 控制面创建独立租户和 World Run，使用动态 Player subject 创建设立案并完成全部 23 个工作项；同时清除财务岗位空缺通知中残留的固定 `founder-principal` 回退。
- 原因：PLAN-GXZ-001 Z32 要求证明零起点创建的企业能够走完真实 M9，而不只是证明 Runtime 可以安装。
- 影响：GX-ZERO 状态从 Validating 更新为 Completed；M9 完成口径包含 8/8 provisioning checkpoint、23/23 工作项、7/7 正式审批门、3/3 可信 World wait、6 次 Agent run/5 类 Agent 和最终 `enterprise_operational_ready`。
- 验证：全新租户经真实 HTTP API 验收；AESE `go test ./...`、前端 20 files/56 tests、生产构建通过；IAOS `go test ./...` 与 System Atlas tracking 通过。
- 后续：进入 M10 的订单、采购和应付闭环；F15–F35 继续按 M10–M13 财务路线图实施，不回写为 M9 欠项。

## 2026-07-30 - 修复 Genesis 企业列表失效会话误报

- 变更：AESE BFF 保留 IAOS 控制面 HTTP 状态，将上游 401 映射为 `player_session_expired`；前端在持久 Player Token 失效时使用当前 IAOS Token 重试一次，成功后更新 Player Token，两者均失效则清理旧凭据并给出重新登录指引。
- 原因：企业大厅总是优先读取持久 `aese_genesis_player_token`，该 Token 过期或 IAOS 重启后仍被重复转发；BFF 又把上游 401 包装成可重试 500，导致用户误以为后端未启动。
- 影响：服务健康与身份失效被明确区分；有效的当前 IAOS 会话可以无感恢复企业列表，真正过期时不会循环重试、泄露上游响应或绕过 owner membership。
- 验证：curl 稳定复现旧 500；新增 Go 控制面状态保留和 BFF 401 回归测试、前端双 Token 恢复与清理测试；全量验证见本次提交证据。
- 后续：正式多人环境接入 IAOS Player Account/OIDC refresh token，替换当前私人开发环境的本机游戏用户名。

## 2026-07-30 - 恢复控制面上线前的 Genesis 企业会话

- 变更：IAOS 增加旧 Workspace 安全接管 API；AESE 在 session 兑换收到 404 时，仅对当前本地玩家拥有的旧 Workspace 发起接管，随后走正常 session API。
- 原因：`gxw-f4b3ce3ce8e2712d` 已有 active tenant、M9 Runtime、World 和完成的设立案，但创建于新版 Genesis 控制面之前，缺少 `genesis_workspace` 行，继续游戏因此被错误包装成 502。
- 影响：旧企业无需删除、重建或重跑 23 个 M9 节点即可恢复；tenant 和 owner 不由 AESE 声明，IAOS 必须验证 owner access、chair、Founder Mandate、Runtime 与 case 后才登记。
- 验证：IAOS 三个 Go module 全量测试、Code Map、治理写入和 Atlas 检查通过；AESE `go test ./...`、20 个前端测试文件/58 项测试及生产构建通过。目标 `gxw-f4b3ce3ce8e2712d` 经 `:4173` 代理兑换返回 200、8/8 checkpoint completed、tenant token 已签发；重复调用后数据库仍为 1 Workspace、1 member、8 step，原设立案保持 `enterprise_operational_ready`。
- 后续：正式多人环境迁移完成后移除 loopback local adapter 的创建能力，保留只读迁移审计期。

## 2026-07-30 - 阻止旧企业继续游戏重复 Provisioning

- 变更：AESE loopback adapter 对已有 active tenant 先执行 Founder login 并直接返回 session；不再调用 identity bootstrap、Runtime install 或 tenant activate。
- 原因：浏览器没有 IAOS Player Token 时会进入本地 fallback；旧实现把 session refresh 当成首次 provisioning，每次点击都重复 bootstrap，随后 Runtime 1.8.0 重装失败并被包装成 502。
- 影响：旧本地企业可以在开发环境继续游戏，同时既有身份、Runtime、流程和业务数据不会被点击操作改写；active tenant 登录失败将直接失败关闭。
- 验证：回归测试先稳定复现 active tenant 触发 1 次 bootstrap 和 1 次 Runtime install，修复后两者均为 0 且 session 成功签发；AESE 全量 Go 测试通过。现场通过 `:4173` 发起与浏览器一致的无 Authorization session 请求返回 200 和 tenant token；IAOS 日志只出现 tenant GET 与 auth login，未再出现 bootstrap、Runtime install 或 activate。
- 后续：生产部署继续要求正式 Player Account/OIDC；确定性 Founder 凭据仅保留于 loopback 迁移期。

## 2026-07-30 - 修复首页 AI 创意官入口与布局重叠

- 变更：把首页“AI 创意官”从装饰性 `div` 改为可点击、可键盘聚焦的创建企业入口；移除覆盖指挥卡的绝对定位，改为正常网格流布局，并补充悬停、按下和焦点状态。
- 原因：原组件没有交互处理或可访问语义，且 `right: -8px`、`bottom: 50px` 的绝对定位会在桌面和窄屏压住相邻页面区域。
- 影响：点击“AI 创意官”与“创建新企业”进入同一个“先定义创业项目”流程；入口不再覆盖主行动卡或“我的企业”区块。
- 验证：前端 20 个测试文件/58 项测试、TypeScript 类型检查和生产构建通过；Genesis 首页 Playwright 在 1440、1280 和 390 三种视口验证按钮可见、可点击且与指挥卡边界不相交。
- 后续：后续首页角色入口继续复用同一可访问交互和正常流布局规则。

## 2026-07-30 - 重构首页 AI 创意官位置与功能说明

- 变更：将“AI 创意官”从指挥中心外部移入卡片第一步，在企业占位信息与四阶段流程之间展示；入口内增加业务目的、输入范围和权责边界说明。
- 原因：独立悬挂按钮虽然不再重叠，但与主流程缺少空间和语义关联，用户无法判断它为何存在以及点击后会发生什么。
- 影响：用户能直接理解 AI 创意官会把行业、区域和经营目标整理成创业项目草案并引导进入身份工作室；注册、审批和经营决策仍由用户及 IAOS 治理流程负责。
- 验证：定向 Playwright 在 1440、1280、390 三个视口验证入口包含于指挥中心、位于流程列表之前、说明可见且点击后进入“先定义创业项目”；桌面和移动完整页面截图人工检查无重叠。
- 后续：Genesis 后续角色入口均在其所属流程节点附近说明输入、输出和不可越权边界。

## 2026-07-30 - 恢复 Genesis 新企业财务治理安装

- 变更：IAOS 将 M9 财务治理安装拆分为既有共享 schema 的只读合同验证和当前 tenant 的 DML；不再在每次新企业 Runtime 安装时重复执行表所有者 DDL。AESE 将 onboarding 幂等键持久化到浏览器，失败刷新后仍恢复同一 Workspace，成功后才清除。
- 原因：`finance_duty_definition` 等共享对象由数据库初始化所有者持有，普通 `iaos_app` Runtime 执行 `ALTER TABLE` 被 PostgreSQL 拒绝，AESE 创建 Workspace 因而返回 502。
- 影响：创建企业仍保留 FORCE RLS、tenant policy 和 SoD trigger 硬门，不需要向普通 Runtime 授予共享表所有权；部分 schema 合同仍失败关闭。
- 验证：普通 Founder 会话通过 AESE BFF 创建全新 Workspace 返回 201、8/8 checkpoint 完成；原 502 Workspace 使用原幂等键重试后复用同一 Workspace 并恢复 active。前端回归验证页面 remount 复用原键、成功清除后生成新键；IAOS 全量 Go 测试、vet、Code Map、治理写入和 Atlas 检查通过。
- 后续：历史全量 Atlas 同步中个别旧声明仍可能因已废弃节点返回 404；本次声明已单独幂等登记。

## 2026-07-30 - 恢复 MiniMax-M3 真实企业名称生成

- 变更：新增 AESE 标准部署脚本，安全加载权限为 0600 的 `.env`，完整校验 MiniMax key/base/model，重建重启服务并验证 Provider；同步增加无密钥模板、设计约束、运行手册和解决方案。
- 原因：本机已有 MiniMax 配置，但旧手工启动进程没有继承环境变量，名称生成因此明确回退到 deterministic。
- 影响：Enterprise Genesis 身份工作室重新使用真实 MiniMax-M3 生成企业名称；配置缺项、secret 权限过宽或启动状态异常会失败关闭，密钥不进入命令行、日志、证据或提交。
- 验证：直接端点与 4173 BFF 均返回 connected/MiniMax/MiniMax-M3；有效 Workspace tenant session 完成真实调用，CreativeJob 为 completed，耗时 31,805 ms、共 1,903 tokens、校验 valid、无 fallback，并返回 4 个动态候选。
- 后续：正式环境把 MINMAX 配置迁移到 Secret Manager/编排器注入，并为 Provider 增加独立 readiness 和额度告警。

## 2026-07-30 - 恢复资本承诺的创始人治理责任

- 变更：IAOS Runtime 1.9.0 将 `capital.commitment.record` 从 `finance-agent` 自主任务修正为 `founder-principal` 以 `chair` 岗位执行的人工任务；Runtime 公共升级同时按租户解析 `platform_super_admin` 或 `genesis_owner` discovery role。AESE 同步补强 M9 设计、实施计划和路线图。
- 原因：计划投入资本是创始人/投资人的所有者决定，金额因企业而异；原流程编译错误地把该节点交给财务 Agent，导致 50,000,000 CNY 被 Agent 的 1,000,000 CNY 自主授权上限拒绝。
- 影响：Agent 金额上限继续对明确委托的自主金额动作失败关闭，但不再充当企业注册资本或人工资本承诺上限；财务 Agent 负责方案测算、规则校验、到账核验和后续会计处理。
- 验证：IAOS 全量 `go test ./... -count=1`、`go vet ./...`、交互式生命周期和 Agent 授权 PostgreSQL 集成测试、Atlas/Code Map/治理写入检查通过；目标租户已安装 Runtime 1.9.0，节点 4 API 与数据库均显示 `human_task / founder-principal / chair / ready`。未代替用户提交 50,000,000 CNY 的治理决定。
- 后续：法域、行业、股东协议和资本结构的金额约束应进入 Policy Profile；不得通过提高通用 Agent 授权上限绕过责任主体修正。

## 2026-07-30 - 恢复 Founder 对 Genesis 新租户的登录发现

- 变更：IAOS 全局登录租户发现移除目标 tenant-local username 相等条件，按 platform principal 的有效 access assignment、principal binding 和 active user 返回选择项；新增源 `founder-principal`、目标 `genesis-owner` 的真实 PostgreSQL 回归。
- 原因：Genesis Workspace 为租户创建 `genesis-owner` 本地用户，但目标 Workspace 的 owner subject 仍是 `founder-principal`；旧查询错误比较本地用户名，导致正式 owner access 已存在却不显示企业。
- 影响：Founder 重新登录 IAOS 时可以看到并进入自己从 AESE 创建的新企业；每张选择卡仍使用目标 tenant-scoped JWT 与目标角色，同名无授权账户不会被合并。
- 验证：修复前回归稳定遗漏目标租户，修复后通过；IAOS 全量 Go、vet、Atlas、Code Map 与治理写入检查通过。后端部署后真实登录 API 返回 200/multiple，并确认包含 `tenant-gx-bfe7c4374e9340319017`（“还是要赚钱啊”）；健康检查为 UP。
- 后续：正式 Player/OIDC 接入继续使用不可变 platform subject 关联 Workspace，不使用用户名、邮箱或显示名作为跨租户身份键。

## 2026-07-30 - 建立 Enterprise Genesis 真实登录与注册

- 变更：IAOS 增加全局 Genesis PlayerAccount、注册/密码登录/session profile、连续失败锁定和既有 IAOS credential 安全提升；AESE 增加同源认证 BFF 与登录/注册表单，短期 Player Token 使用 sessionStorage，Workspace owner 每次从 IAOS 验证后的 subject 派生。默认 IAOS 模式彻底移除只输用户名登录，本地适配器仅允许显式 `local_dev + loopback`。
- 原因：原 AESE 用户名和 `X-Genesis-Player-Id` 都由浏览器声明，既无注册和密码验证，也能被伪造，不能承担 Workspace ownership、跨租户隔离和正式交付身份。
- 影响：新玩家先注册平台身份，登录后再独立创建企业；现有 Founder 可用原 IAOS 凭据首次提升并保持历史 Workspace；错误密码、锁定、过期/伪造 Token 和跨玩家读取均失败关闭，AESE 不保存密码。
- 验证：IAOS 全量 Go 测试与 vet、Player 真实 PostgreSQL 注册/登录/session/重复账号/五次失败锁定/双 Player Workspace 隔离通过；真实 `founder-principal` 提升后恢复 5 个原 Workspace，伪造 Player Header 且无 Token 返回 401。AESE Go 全量测试、前端 21 文件/61 项测试、TypeScript、定向 ESLint 与生产构建通过；Genesis Home 15 项三视口 Playwright 及真实 Founder 3 项三视口 Playwright 通过且无控制台错误。IAOS 8082 与 AESE 8090 已重建部署。
- 后续：正式互联网部署把短期 Bearer 升级为 HttpOnly/Secure/SameSite Cookie，并补 refresh/revoke、邮箱验证、找回密码、MFA、OIDC 与设备会话管理。
## 2026-07-30 - M9 账套与年度期间可配置

- 变更：AESE“启用会计账套与期间”新增账套名称、会计年度和 12 期可编辑自然月表单；IAOS 能力合同持久化完整期间并在财务开业投影中返回明细。
- 原因：原实现无参数创建编码账套且只生成当前月，用户无法确认账套身份、年度边界或期间范围。
- 影响：新账套必须具名；12 期必须连续、无重叠并覆盖当前会计年度，当前业务日期所在期间开放，其余期间为未来。
- 验证：Go 单元测试覆盖自然月和期间间隙拒绝；前端类型检查、单测、构建及三服务健康检查。
- 后续：13 期制、4-4-5 日历、跨年预建和独立期间开关留给后续 Fiscal Calendar 配置能力。
## 2026-07-30 - 多组织财务与共享主数据基础设计

- 变更：新增 DES-031，定义集团/法人/BU/基地/共享服务中心、多账簿、科目表两段式、Data Set、Business Partner canonical/组织扩展和模块期间状态机；审计目标租户现有账套投影。
- 原因：单法人 M9 tracer 不能支撑集团分子公司、跨组织共享科目/客户/供应商，也无法正确表达共享数据所有权和组织扩展。
- 影响：明确 Tenant 不等于法人、共享配置不共享法定余额、全局身份加组织扩展、消费者只读引用和跨法人双边单据原则；模块期间控制本轮只记录不实现。
- 验证：对照 SAP Company Code/CoA/Controlling Area/MDG 与 Oracle Enterprise Structure/Reference Data Set/BU/Supplier Site 官方实践；检查目标租户权威账套与 Entity 投影字段。
- 后续：先实现组织与 Data Set，再重构账套/科目表/财政日历并迁移 M9 数据，随后实现 Business Partner 和产品组织扩展。

## 2026-07-30 - 财务设计文档模块化

- 变更：建立 `docs/designs/finance/` 独立目录；将 620 行 DES-030 缩为 108 行总览，迁移 DES-031，并新增 DES-032–036 分别承载会计内核、子账资金、制造成本资产、预算关账报表和财务治理 Agent；M9 改为只引用财务模块。
- 原因：财务需求已跨越 M9–M13 和持续经营，单文件继续膨胀会降低阅读、维护、实现追踪和后续 Agent 定位效率。
- 影响：财务规格按稳定功能模块维护，里程碑只描述业务事实和完成门；实现状态仍由路线图和财务计划管理，文档拆分不改变已实现范围。
- 验证：检查 DES frontmatter、目录索引、仓库内旧路径、Markdown 相对链接、Atlas 声明、code map freshness 和 `git diff --check`。
- 后续：实现 DES-031 的组织/Data Set，随后按 DES-032–036 独立推进账套迁移、子账、成本资产、关账报表和治理 UI。

## 2026-07-30 - 财务多组织与共享 Data Set F5B

- 变更：IAOS 新增集团/法人/BU/基地/工厂/财务组织/共享中心、Data Set、决定项分配和主体组织授权四类权威数据，两项可执行 Capability、RLS/复合 tenant 外键/写入触发器、历史账套幂等迁移及财务工作台表单；AESE baseline 新增 HCTM 六组织、两个 Data Set、九个分配和三个岗位访问模板。
- 原因：Tenant 不能继续隐式等同单一法人，科目、伙伴、产品、币种和汇率也不能只有“全租户共享”或“逐公司复制”两种选择。
- 影响：F5B 形成集团财务和后续采购/销售/制造主数据共同依赖的组织与共享边界；AESE 保持场景模板所有权，IAOS 保持运行数据所有权。
- 验证：IAOS 后端单元/Runtime 测试、真实 PostgreSQL 治理/幂等/RLS/跨租户外键/完整回填集成测试、TypeScript 检查和生产构建；目标租户升级 Runtime 2.0.0 后由真实身份 API 读到 5 个组织、2 个 Data Set、9 个分配和 6 个访问授权；AESE `financebaseline` 引用/闭集/稳定编码测试、场景包离线校验、JSON、Atlas 和 Code Map 检查。
- 后续：F5C 重构账套原型、科目表法人扩展、财政日历和多账簿并迁移 M9 数据；F5D 实现 Business Partner 及客户/供应商/产品组织扩展。

## 2026-07-30 - 财务账簿、共享科目表与财政日历 F5C

- 变更：IAOS 新增共享科目表/科目定义、法人科目扩展、财政日历/期间、账簿和账簿集权威模型，发布 `finance.ledger.foundation.configure`、财务工作台“账簿与科目”入口及 M9 数据幂等迁移；AESE baseline 升级到 1.2，固定 HCTM 账簿、1002/4001 科目、自然年 12 期和账簿集模板。
- 原因：M9 临时账套和逐账簿科目不能支撑多法人共享定义、多账簿、法定余额隔离及 M10 后持续经营。
- 影响：稳定 `BOOK-*`、凭证关系和历史事实不变；共享科目定义与法人启用属性分离，账簿明确绑定法人、准则、本位币、科目表和日历；所有变更继续经过 Capability、RLS 和数据库写入硬门。
- 验证：IAOS 全量 Go、真实 PostgreSQL 迁移、治理写入/Atlas、TypeScript、组件测试和生产构建通过；目标租户 API 返回 1 账簿、1 科目表、2 科目定义、2 法人扩展、1 日历、12 期间和 1 账簿集，Runtime 升级到 2.1.0；AESE 离线校验覆盖未知法人、断裂期间和控制矩阵。
- 后续：实施 F5D Business Partner、客户/供应商及产品 canonical identity 和组织扩展；F5E 模块期间控制保持关账阶段范围。

## 2026-07-30 - 财务伙伴与产品共享主数据 F5D

- 变更：IAOS 新增集团唯一 Business Partner、客户/供应商角色、法人/BU 扩展、共享产品和工厂扩展五类权威数据，两项 Capability、RLS/复合 tenant 外键/写入触发器及财务工作台“伙伴与产品”表单；AESE baseline 1.3 新增星河客户、两家铝材供应商、冷却板成品和铝板原材料模板。
- 原因：客户、供应商和产品不能按法人重复建立互不关联的身份，也不能因共享身份而错误共享凭证、余额、库存或成本事实。
- 影响：Data Set 决定集团共享定义，法人/BU/工厂扩展保存本地财务、商业、MRP 和估价属性；相同主体可同时承担客户和供应商角色并保留唯一身份。
- 验证：IAOS 全量 Go/vet、真实 PostgreSQL canonical/role/组织扩展集成测试、治理写入/Atlas、TypeScript、组件测试和生产构建通过；启动历史开业回填无错误；AESE 全量 Go 与 baseline 未知组织、无效成本批量、稳定编码和控制矩阵校验通过。
- 后续：F5E 模块期间控制按关账阶段实施；M10/M11 从该伙伴和产品基础接入 AP、资产、采购和工厂扩展，不再复制单法人主数据。

## 2026-07-30 - 财务菜单拆分与 Semantic/Entity 资产发布 F5D2

- 变更：IAOS 将“财务账务与报表”收敛为账务查询和报表输出，新增财务组织与待办、多组织与共享数据、账簿与会计基础、客户、供应商和产品独立入口；Runtime 2.3.0 把 F5B–F5D 的 16 类权威模型注册到数据模型工坊，并在受治理 Capability 事务内同步 `entity_projection_*`。
- 原因：原财务页混合账务、组织、待办和主数据，且权威 `finance_*` 表虽已实现，却没有完整的 Semantic/Entity 发布，用户无法从数据模型工坊理解字段、关系和菜单。
- 影响：业务工作台与模型工坊各司其职；客户/供应商成为同一 Business Partner 的角色视图，产品和组织扩展不复制身份；通用 Entity 菜单读取权威数据投影，不能绕过 Capability 写入。
- 验证：IAOS 全量 Go/vet、真实 PostgreSQL Runtime 集成、前端 TypeScript/组件测试/生产构建及 Atlas/Code Map/治理写入检查通过；目标租户安装 Runtime 2.3.0 返回 201 和 782 项资产写入，API 可见 35 个 Entity、七个独立工作台，组织 5/5、科目 2/2 权威/投影对账一致。
- 后续：目标租户当前尚无 Business Partner/产品运行记录，后续必须由相应配置 Capability 或 AESE 场景应用创建，不能为了填充菜单直接写投影；F5E 与 F15–F35 按计划继续。

## 2026-07-30 - 建立平台基础语义包与 M9 消费边界

- 变更：IAOS 发布 `iaos_foundation_semantics@1.0.0`、Semantic Governance Compiler 和检索 CLI；M9 Runtime 升级到 2.4.0，只消费基础包并修正 Business Partner/Product 继承；AESE 计划、路线图和 Code Map 增加跨仓合同。
- 原因：旧参考目录、Metadata bootstrap、M9 安装器和租户 Entity 各自定义语义，导致 `business_partner`、`organization`、`party`、`product` 的继承和字段前后不一致，Agent 也缺少发布期硬约束。
- 影响：87 个 Concept、51 个 Archetype 成为版本化机器权威；跨根族、循环、非 active 父类、Core→Domain 和重复 canonical slot 不能安装或发布；新增语义必须先检索并与产品所有者协商，AESE 场景不能覆盖平台资产。
- 验证：IAOS 全量 Go/vet、基础包校验、前端生产构建、真实 PostgreSQL 安装及 Atlas/Code Map/治理写入检查通过；数据库确认 `business_partner → party`、`organization → party`、`product → material`。
- 后续：递归 Effective Archetype Artifact 编译和 Semantic Change Proposal 审批 UI 作为后续平台切片；M10 先复用目录，不得建立同义 seed。

## 2026-07-30 - M9 Effective Runtime Artifact 执行权威闭环

- 变更：IAOS 新增 DES-072 和 Effective Runtime Artifact 深模块，将 Entity Schema/UI/Agent Context、Capability API/Agent Tool、Process 发布与运行统一到不可变编译产物；流程发布冻结能力 artifact version/hash，AESE 的 DES-027、M9 计划、路线图和 Code Map 同步引用该跨仓合同。
- 原因：DES-023 原先只部分实现，正式路径仍可能回读 metadata、authoring DSL 或最新能力版本，无法保证用户看到的配置就是 Runtime 和 Agent 实际执行的合同。
- 影响：缺失、stale、编译器不匹配、哈希不一致和未解析依赖均失败关闭；升级回填按对象隔离，坏对象不再阻断同租户其他有效对象；实施人员可在数据模型工坊和流程编排控制台穿透查看运行权威。
- 验证：IAOS 全量 Go/vet、前端 23 项定向测试、TypeScript、生产构建、语义基础包、Atlas 和 diff 检查通过；真实 PostgreSQL 验证 FORCE RLS、跨租户不可见、350 个 active Entity Artifact 与 128 个 active Process Artifact；IAOS main `945a05b` 已推送，8082/3000 已重建并通过健康检查。
- 后续：M10 及后续场景只声明对 artifact version/hash 的依赖，不得新增 authoring 表或场景 JSON 运行兜底。

## 2026-07-30 - M9 通用资产平台基础包与租户 Edition

- 变更：IAOS 发布 `genesis-m9@1.0.0`，把 M9 Semantic、Entity、Capability、Process、Policy、Menu 和五个 Agent Template 分为三个不可变签名包；tenant-001 安装参考 Edition，Genesis provisioning 和历史租户升级复用同一清单；AESE 的 DES-027、计划、路线图和 Code Map 同步跨仓合同。
- 原因：通用资产过去只是某个 M9 Runtime 安装器的隐式副作用，平台租户、新租户和旧租户无法对账同一产品版本，也容易把 tenant-001 错当成业务数据母租户。
- 影响：平台定义按 SemVer/SHA-256 发布，Agent Template 与租户 Agent/Mandate 实例分离；任何新企业只安装定义并自行形成身份与业务事实，不复制其他租户记录。
- 验证：IAOS 全量 Go、TypeScript、Vitest、Semantic/Foundation、治理写入、Code Map、Atlas 检查通过；真实 PostgreSQL 确认 tenant-001 三个 active installation、Runtime 2.4.0，案件/审批/凭证/World Journal 均为 0；旧 Genesis 租户正确返回 `upgrade_required`，8082/3000 已生产重建。
- 后续：M10 及后续通用资产必须发布新 package/Edition SemVer；历史租户由有权主体在 IAOS 平台基础包页面显式升级。

## 2026-07-31 - M9 Entity 继承与受治理记录浏览闭环

- 变更：IAOS Runtime 2.5.0 修复 M9/历史 Entity 的 Foundation 字段类型漂移，发布门增加继承和类型校验，平台包升级重新编译 Effective Artifact；通用记录 API 兼容原生领域表，菜单显示 direct/capability 写入模式及业务入口。AESE M9 主计划和 Code Map 同步跨仓合同。
- 原因：设立案件等 Entity 出现 `owner_id`、`org_node_id` Archetype 类型冲突，部分原生表菜单无法加载明细，领域投影的只读原因和正确录入入口也不可见。
- 影响：AESE 仍只调用 IAOS Capability/Process，不拥有或直接写运行数据；动态 Entity 可按权限直接维护，领域投影必须经业务工作台和 Capability 写入。
- 验证：IAOS `tenant-001` 46 个 Entity 阻断错误为 0，代表性设立、财务、库存、采购和产品列表 API 返回 200；Go、前端组件测试、生产构建、Atlas、Code Map 和线上健康检查通过，IAOS main `e29e0c2` 已推送。
- 后续：逐步消除 7 条关系展示字段和 10 条未选择 Archetype 的旧制造模型非阻断 warning。

## 2026-07-31 - Genesis 投影会话与 Runtime stale 恢复

- 变更：AESE projection 在旧 tenant/token 造成 401/404 时刷新 Workspace session 并重试；Agent 写请求仅对明确 Runtime stale 422 做一次安全重试。IAOS Workspace session 在签发 token 前幂等对齐 Genesis-managed Edition。
- 原因：已有企业仍安装 Runtime 2.4.0，平台运行代码已是 2.5.0；同时页面首次投影可能读取上一企业的 localStorage 会话，分别造成 Agent 422 和投影 404。
- 影响：受管 Workspace 重新进入时自动收敛，普通 IAOS 租户仍保留显式升级；失败请求不会跳过节点或留下部分业务写入。
- 验证：curl 复现原 422 `effective runtime artifact stale`；正确会话 projection 返回 200；前端新增两条 stale session 回归；IAOS helper 覆盖 stale/current/failure；目标 GX 租户已升级 2.5.0。
- 后续：继续保留平台包控制台作为普通租户的显式版本治理入口。

## 2026-07-31 - M9 设立案件 Entity 字段一致性

- 变更：IAOS Runtime 2.6.0 显式发布设立案件名称、拟设企业名称、注册地址和经营范围，案件 Capability 同步权威表与 Entity 投影；平台安装器新增权威业务列自动发现和字段一致性失败关闭门。AESE M9 计划、路线图、Code Map 与 Atlas 同步该跨仓合同。
- 原因：数据库 `incorporation_case` 已有四个业务列，但旧 Entity 编译器只发布 Archetype 字段和 `payload`，导致数据模型工坊与真实运行数据不一致，原语义分析也无法发现完整性缺口。
- 影响：用户可在数据模型工坊和设立案件菜单按显式字段理解、查询数据；AESE 仍只通过 IAOS Capability/Process 写入，不直接维护投影；未来权威表新增业务列未完成语义分类和投影映射时不能发布。
- 验证：IAOS 单元/API 测试和真实 PostgreSQL 字段闭环测试通过；`tenant-001` 与 `tenant-gx-472324ae8bac4af39519` 已升级 `genesis-m9@1.0.2` / Runtime 2.6.0，两个租户均返回四个字段，现有 GX 案件权威表与投影逐字段一致，语义分析 0 error / 0 warning。
- 后续：新增原生权威 Entity 时复用同一 parity contract；先检索平台语义目录，再提交业务字段的层级、兼容、迁移和受影响资产说明。

## 2026-07-31 - M9 原生 Entity 投影真实性与明细闭环

- 变更：IAOS Runtime 2.8.1 移除把完整设立案件状态复制到所有 Entity 的旧同步，改由 Agent output、Capability journal 和权威领域表按事实所有权逐项物化；治理决议等生命周期 Entity 增加领域显式字段和列表视图，财务基础投影补齐显式列、稳定业务引用及升级期必填完整性门。
- 原因：`governance_resolution` 等逻辑 Entity 虽有正确的 `entity_projection_<code>` 物理读模型，但旧数据只写通用 payload 或提前复制未来状态，导致菜单详情为空、字段语义错误，并伪造尚未发生的后续业务事实。
- 影响：用户继续以逻辑 Entity code 使用模型和菜单；物理表前缀只表示统一读模型。节点 2 完成后只出现治理决议，登记、法人、开户、注资、任命和预算必须等待各自 Capability 提交后才出现；原生投影不能直接 CRUD。
- 验证：IAOS 单元、API、交互式设立集成、财务投影集成、Go vet、语义/治理/Atlas/Code Map 检查通过；`tenant-001` 与目标 GX 租户安装 `genesis-m9@1.0.5` / Runtime 2.8.1，目标案件治理决议 1 条且明细列完整，六类未来 Entity 均为 0，8082 健康检查通过。
- 后续：M10 及后续原生 Entity 必须声明 authority、projection owner 和字段物化映射，并复用“未发生零记录、显式必填不为空”的发布及集成测试门。

## 2026-07-31 - M9 Entity 存储与唯一写入权威

- 变更：IAOS Runtime 2.9.0 / `genesis-m9@1.1.0` 新增三类 Entity 存储合同，把权威来源、唯一写入者和投影维护者编译进 Effective Runtime Artifact；M9 计划、路线图和 Code Map 同步跨仓合同。
- 原因：仅靠 `entity_projection_` 前缀无法区分动态权威表、复杂领域读模型和纯计算汇总，也不能阻止通用 CRUD 与 Capability 同时写同一 Entity。
- 影响：新动态 Entity 统一使用 `entity_record_<code>`；复杂领域对象必须由唯一 Capability/Process 写专用权威表；Journal/Aggregate/Query 汇总只读。AESE 继续只通过 IAOS 受治理入口推进世界，不直接写 Entity 表或投影。
- 验证：IAOS 单元、API、交互式设立集成、Go vet、TypeScript、组件和生产构建通过；真实租户共 35 个 M9 Entity 均为 version 1 合同（27 个领域投影、8 个计算投影），计算投影直接 Create 返回 `405 entity_write_owner_enforced`，8082/3000 健康。
- 后续：历史 contract version 0 动态 Entity 按独立迁移计划冻结写入、对账、切换 metadata 和重编译 Artifact；不得自动改名或用降级规避合同。

## 2026-07-31 - M9 Command Gateway 与 Process Artifact 运行权威

- 变更：AESE 增加同源白名单 Command Gateway，案件、工作项、Agent、审批和 World Observation 写操作不再由浏览器直达 IAOS；IAOS 原生 Capability 校验不可变 Artifact，设立工作项改由 active Process Artifact 递归展开并锁定子流程/能力版本哈希。
- 原因：既有实现违反 ADR-003 的浏览器写边界，且 M9 专用 Runtime 绕过 ADR-005/DES-072，以 Go ProcessDefinition 作为工作项运行源。
- 影响：AESE 仍不保存 IAOS 业务数据或管理员凭据；用户发布内容必须成功编译为 Effective Artifact 才能运行，缺失、漂移和篡改均失败关闭；IAOS 资产发布为 Runtime 2.10.0 / genesis-m9@1.2.0。
- 验证：AESE 全量 Go、网关/客户端测试、前端 10 项 API 测试与生产构建；IAOS 全量 Go/vet、治理写入/Code Map/Atlas；tenant-001 已安装 Runtime 2.10.0 / genesis-m9@1.2.0，主流程 Artifact 固定 5 个子流程，Capability 使用标准 compiler 和 64 位 hash；实时 Gateway 拒绝任意 Entity 路径并正确透传允许路径。
- 后续：用户可用新建 Genesis 企业执行浏览器逐节点验收；浏览器 Network 中所有 M9 POST 应为 `/api/aese/v1/commands/iaos/*`。

## 2026-07-31 - 财务业务来源与通用凭证过账边界

- 变更：IAOS Runtime 2.11.0 / `genesis-m9@1.3.0` 发布
  `finance.journal.entry.post`；`capital.contribution.post` 在同一受治理事务委托过账，
  `journal_entry`、`journal_line` 的 metadata 写入所有者同步修正。AESE 财务设计、
  实施计划、路线图和 Code Map 引用该跨仓合同。
- 原因：旧 metadata 把凭证主子表归给 `finance.foundation.configure`，与真实执行路径
  不一致，导致模型工坊和穿透查询错误解释。
- 影响：M9 资本入账现在可区分业务来源、Posting Profile、通用过账 Execution 和权威
  存储；AESE 仍只通过 IAOS 受治理命令推进，不保存或直接写凭证。
- 验证：IAOS 全量 Go、治理/Code Map/Atlas、三组真实 PostgreSQL 财务测试通过；
  tenant-001 运行库确认两个 Entity owner、两个 published Capability 和 active binding。
- 后续：人工凭证、冲销、汇率审批、期间关账、AP/AR、资产、成本与预算按 F15–F35
  实现；未具备真实运行逻辑前不得注册空 Capability 冒充完成。

## 2026-07-31 - M9 接入 IAOS 可执行原子能力目录

- 变更：IAOS 以 DES-076 重构原子能力口径并交付 19 项真实可执行 V1 目录、统一
  Analyzer/Artifact/Runtime/API、受控 Entity CRUD 原语和 Capability Studio 业务向导；
  AESE M9 计划引用该跨仓执行合同。
- 原因：旧参考文档把审批、过账、完工和 LLM 调用等完整业务动作误列为原子能力，且
  文档、Analyzer 与 Dispatcher 曾发生清单漂移，不适合作为 M9 Agent/Process 的稳定底座。
- 影响：AESE 仍只调用已发布 Business Capability/Process/Bridge，不复制 IAOS Handler；
  原子能力仅在外层租户、权限、Policy、Approval 和事务上下文中执行。
- 验证：IAOS 全仓 Go、真实 PostgreSQL Capability integration、TypeScript、原子目录组件、
  生产构建、8082/3000 部署及 live API（19 total、18 active、1 deprecated）通过；commit
  `67fe3d2` 已推送 main。
- 后续：M10+ 新增业务动作先复用现有原子目录；确需新增原子项时必须先提交 typed contract、
  Handler、RLS/回滚测试和平台版本评审，禁止由 AESE 场景包临时扩展。

## 2026-07-31 - M9 对齐能力生命周期与角色执行授权

- 变更：IAOS 新建能力固定进入规范草稿，角色与权限新增 `capability.<code>/EXECUTE` 授权面，执行器只允许 Active Artifact；AESE M9 计划引用该客户配置闭环。
- 原因：旧 UI 把悬空草稿当作可执行目录项，文档又引用不存在的权限入口，导致用户无法理解 `Test.workorder.create`、`approve_document` 等能力为何不能运行或没有可视化模型。
- 影响：AESE 只能消费 IAOS 已发布、已绑定、已授权的能力；草稿可视化会明确回退 Draft，不能再伪装为运行失败。
- 验证：IAOS 33 项前端定向测试、后端 API/Capability 测试、TypeScript、Go/Next 生产构建及 8082/3000 部署通过；启动迁移在 system tenant 上修复 57 条历史悬空草稿，`tenant-001/Test.workorder.create` 已恢复为规范 draft；Atlas 声明校验通过，运行时同步端点暂返 404，待 System Atlas 写入路由恢复后补同步。
- 后续：M9 测试仅选用平台基础包的 Active Capability；旧 Test 能力若指向 governed projection，应废弃并以声明的领域业务能力替代，不能发布为通用 CRUD。
## 2026-08-01 - AESE 场景知识接入 IAOS 产品知识中枢

- 变更：新增 DES-037、PLAN-KNOWLEDGE-AESE-001 和 M9 企业设立用户手册第一版。
- 原因：现有设计、代码和页面说明缺少面向最终用户的统一检索和 Agent 可引用入口。
- 影响：明确 AESE 只拥有场景知识和 World 映射，IAOS 拥有通用知识 Registry、权限和 Copilot；避免复制平台知识系统。
- 验证：新增文档均带治理头；跨仓引用已核对；Atlas 声明已新增。
- 后续：建立场景知识 Schema/manifest、Knowledge Edition 安装、节点帮助、上下文 Agent 与双侧证据验收。

## 2026-08-01 - M9 场景知识机器合同与节点深链

- 变更：新增场景 Knowledge Edition JSON Schema、`hctm-m9-incorporation-knowledge@1.0.0` manifest、18 节点 purpose/input/output/actor/evidence/IAOS/World 映射、canonical SHA-256 校验、`aese knowledge validate|digest` 和任务弹窗“这一步是什么”入口。
- 原因：人工手册不能自动发现流程节点漏配、错误 Capability/World action 引用或内容漂移，也无法从当前任务直接进入 IAOS 权威知识。
- 影响：AESE 仍不保存知识或 IAOS 业务数据；清单作为待发布场景 Edition 输入，页面通过稳定 article ID 跳转 IAOS。内容哈希不冒充 IAOS 安装完成。
- 验证：`go test ./internal/scenarioknowledge ./cmd/aese`、`aese knowledge validate`、JSON Schema fixture 校验和 AESE `npm run build` 通过；18 节点覆盖率 100%，hash 为 `6a73c7735970e14adbbe198f5ebc07bd7b580cf6fc3f32da3137c3633efd6543`。
- 后续：实现 IAOS Knowledge Edition 幂等安装、Agent workspace/node 上下文和 World/IAOS 双侧运行证据漂移检查。

## 2026-08-01 - HCTM Knowledge Edition 签名发布与安装闭环

- 变更：新增 Markdown 到 IAOS `KnowledgeEditionBundle` 的确定性编译、正文 SHA-256、Package/Edition 双层签名、版本化 dist 产物和默认 dry-run 的 `aese knowledge install`；IAOS 增加租户范围行业文章物化与幂等安装 API。
- 原因：仅校验场景 manifest 不能证明实际正文未被替换，也不能保证行业知识只对已安装租户生效。
- 影响：AESE 继续拥有场景内容、不保存 IAOS 业务数据；IAOS 继续拥有权限、RLS、审计、安装目录和有效解析。M9 深链改用独立行业文章 ID，避免与平台基线无痕覆盖。
- 验证：manifest/schema、编译器、CLI、IAOS 全量 Go、Product Knowledge/API 定向 Go 测试、AESE 全量 Go/生产构建/深链 Vitest 通过；live dry-run 验签通过，缺少 `genesis-m9@1.5.0` 时 409 失败关闭；升级依赖后重复 apply 返回 `no_op=true,writes=0`，`tenant-hctm` 可见行业文章而 `tenant-001` 返回 404，额外行业包不影响基础包 `up_to_date=true`。Edition hash 为 `1afd39e1ce2139c1254aea2e01f250d480796e001121e85bbbe72b9d1e59ec47`。
- 后续：推进 Agent workspace/node 上下文和 World/IAOS 双侧证据漂移。

## 2026-08-01 - M9 知识问答接入封闭场景导航上下文

- 变更：M9 节点知识深链增加 workspace、case、world run、node、actor、task 和 capability；
  IAOS 知识中心可见展示并允许清除，Copilot 前端与 BFF 使用同一封闭字段归一化合同。
- 原因：仅传 tenant/case/capability 无法让用户和 Agent 确认正在解释哪个运行和节点，同时任意
  URL 文本直接进入提示会形成注入与伪造运行事实风险。
- 影响：Agent 可用稳定编码定位下一步应查询的 IAOS/World 证据，但导航上下文本身不证明节点
  状态、审批结果或外部 Observation；未知、超长和提示式字段被丢弃。
- 验证：AESE 深链单测与 TypeScript 通过；IAOS 归一化单测、TypeScript 和生产构建通过；
  详细证据见两仓对应的 2026-08-01 scenario context evidence。
- 后续：S9 接入有权限的 World/IAOS 双侧实际证据读取和配置漂移提示。

## 2026-08-01 - M9 知识中心双侧证据与配置漂移闭环

- 变更：IAOS 知识中心以当前 JWT 读取 M9 Evidence Bundle 和工作项，展示节点实际状态、
  IAOS 证据计数与权威哈希、World Observation 接收回执和逐字段漂移；Copilot BFF 独立重读。
- 原因：导航 ID 和场景手册只能定位，不能证明节点、审批、Agent 或外部世界结果。
- 影响：AESE 仍不保存 IAOS 业务数据；World 证据明确限定为 IAOS 已校验并持久化的 Bridge
  回执。越权、旧合同、缺失节点/Capability/Observation 或验签失败时不再推断状态。
- 验证：IAOS resolver 覆盖一致、漂移、验签失败/缺失和旧合同；TypeScript、lint、生产构建、
  runbook、Code Map 和 Atlas 声明作为 S12 发布门。
- 后续：完成 S12 现场 UI/API 验收并继续通用产品知识覆盖。

## 2026-08-01 - M9 场景知识计划 S1–S12 发布收口

- 变更：完成双侧证据 UI、Copilot 服务端重读、失败关闭、Runbook、Code Map、验证报告、
  System Atlas、生产构建和部署，PLAN-KNOWLEDGE-AESE-001 状态改为 completed。
- 原因：场景导航、知识文章和运行证据必须形成用户可验证的完整闭环，不能停留在设计或单测。
- 影响：M9 场景知识不再阻塞 AESE；通用菜单、错误码和 Agent 引用评测继续由 IAOS 主计划治理。
- 验证：知识定向测试 11 项、场景安全/证据测试 9 项、TypeScript、定向 lint、Next 生产构建、
  3000 部署、无凭证 401、不存在案件 404、四端口 200 和 Atlas 同步通过。
- 后续：用用户新建的首个 M9 案件执行 Runbook 正向现场验收，不为收口预造业务数据。

## 2026-08-01 - M10 改为 Agent 候选、人工选择与金额参数化

- 变更：修订 DES-011 与 PLAN-M10，取消正式运行中的三个固定候选、固定赢家、固定预算/现金、固定 WBS 和固定异常假设；新增 `plant-planning-agent` 结构化 proposal、项目负责人审阅、World 事实补齐、金额来源/可编辑边界及 R1–R10 交互修订门。历史 evidence/runbook 明确降级为 `fixture_only` reference replay。
- 原因：固定剧情会绕过 Agent 与业务人员协作，也会把单一示例企业的资金和选址条件错误推广到所有新企业。
- 影响：正式 M10 必须由当前租户 IAOS 权威现金/预算、用户可编辑投资参数、Agent 建议、人工决定和外部 Observation 共同驱动；Agent 不得虚构报价、容量、权属、许可或自行批准资金动作。
- 验证：Markdown `git diff --check` 通过；Roadmap、Code Map、文档索引、历史 evidence 和 runbook 已同步标明实现边界；System Atlas 跟踪声明已新增。
- 后续：按 PLAN-M10 R1–R10 实现 Effective Process Artifact 工作项、Agent Artifact、表单、版本化金额、人工/外部模型双路径和现场验收；完成前只声明 reference replay 完成。

## 2026-08-01 - M10 交互计划与 Agent proposal 首个代码切片

- 变更：新增 PLAN-M10-INTERACTIVE-001；`internal/plantbuild/interactive.go` 实现设施需求、权威财务快照、金额/工期区间、候选 proposal、人工 review 和生成证据合同；MiniMax 暴露受限 JSON completion adapter，AESE 新增 planning status/proposals API，未配置模型返回显式失败而非固定候选。
- 原因：旧 M10 只有编译在 Go 中的十帧回放，无法让用户输入参数、Agent 生成候选或人员作出选择。
- 影响：候选数量 2–8、区域、日期、选项类型和所有业务输入金额均可参数化；现金/预算必须携带 IAOS 权威来源引用与快照 hash；Agent 输出必须披露假设、待验证事实、风险、来源、模型和输入/输出 hash。
- 验证：`go test ./internal/plantbuild ./internal/creative ./internal/httpapi ./cmd/aese-server` 通过；未配置 provider status 和失败关闭有回归测试。
- 后续：完成 Agent 证据持久化与坏 JSON/provider 测试，在 IAOS 独立 worktree 实现 Entity/Capability/Process 权威链，再接人工审阅 UI。

## 2026-08-01 - M10 设施规划 Agent 与人工审阅首个在线纵切

- 变更：Plant Build Play 增加参数化设施需求、IAOS 只读现金/预算来源、Agent 候选解释卡和人工审阅；AESE 定向 BFF 串联 `facility.requirement.define`、`site.proposal.record` 与 `site.proposal.review`，IAOS 首个权威切片保存 Requirement、ProposalSet、Review、Audit、Outbox 和幂等证据。
- 原因：需求、Agent 技术输出和人工决定必须成为可区分、可追踪的事实，不能继续用三个固定候选、固定金额或浏览器本地状态冒充正式 M10。
- 影响：用户现在可完成“权威财务快照 → 需求 → Agent proposal → 人工采纳/退回/淘汰”；实际现金和已批预算不可在页面编辑，Agent 候选仍是 `candidate_only`，审阅不等于外部调查、投资审批或合同承诺。
- 验证：AESE 全量 Go、合同/HTTP/UI 定向测试、TypeScript、ESLint 与生产构建通过；IAOS 全量 Go 和治理检查通过，后端已部署，`tenant-001`、`tenant-hctm` 均已升级至 `genesis-m9@1.6.0` 且 `up_to_date=true`，三个 Capability 的 Active Artifact 可在线读取。完整业务案件现场 UI/API/DB 证据仍属于 S5。
- 后续：实现人工候选权威提交、`site.investigation.request`、Effective Process Artifact、World 报价/权属/容量/许可 Observation，以及选址审批、项目/WBS、施工、付款、验收和 AP/CIP/总账闭环。部署后的 System Atlas 声明同步端点 `/api/v1/system-atlas/updates` 当前返回 404，三个声明已保留在仓库，待 IAOS 写入路由恢复后补同步。

## 2026-08-01 - M9 预算边界参数化与 Agent 金额效力纠偏

- 变更：M9 设计和计划明确场景 fixture 金额只用于 reference replay；正式预算由租户配置最低金额、可选绝对上限和已核验资金比例。IAOS 将预算准备声明为 proposal，保存 revision/hash，并要求 G7 引用同一草案。
- 原因：finance-agent 的 100 万元自主交易限额被错误套到 100,000,008 元预算草案，且准备与批准可以填写不同金额。
- 影响：企业预算不再受平台固定金额限制；真正形成承诺或资金事实的 binding 能力仍保留 Agent Mandate 限额，最终预算仍不得突破权威已核验资金。
- 验证：IAOS `go test ./internal/incorporation ./internal/api`、governed-write 检查、Atlas tracking 和 `git diff --check` 通过；AESE 文档与计划完成一致性更新。
- 后续：部署 Runtime 2.14.0 / `genesis-m9@1.7.0` 后，以新案件验证节点 21 大额草案、节点 22 只读引用和租户策略调整。

## 2026-08-01 - M10 IAOS 菜单与权威工作台补齐

- 变更：为 `genesis-plant-planning` 包增加 `menu.genesis_plant_planning`，在 IAOS `业务智造层` 发布 `M10 工厂规划` 工作台；工作台提供案件选择、权威资金约束、设施需求、Agent 候选、人工评审、穿透证据和 AESE M10 World 深链。
- 原因：M10 已有 API、Capability 和 Agent，但缺少租户可发现的菜单资产与页面路由，用户无法从 IAOS 导航进入已实现的交互纵切。
- 影响：平台基础包升级至 `1.8.0` 后，既有与新租户通过同一包清单获得 M10 入口；IAOS 继续拥有权威业务事实，AESE 继续拥有 World/Agent 交互。该入口不改变 M10 尚未完成调查、Process、工程执行和财务全链的事实。
- 验证：IAOS 新增菜单包回归测试先失败后通过；`go test ./internal/incorporation ./internal/api`、生产构建、治理写入检查、Atlas tracking 和线上菜单/财务约束 API 验证通过；`tenant-001`、`tenant-hctm` 与当前 GX 租户均已升级到 `1.8.0`。
- 后续：从新工作台完成 Requirement → Agent Proposal → Human Review 现场验收，再实现 Investigation、Effective Process Artifact、World Observation、项目/WBS、施工、付款、验收和工程财务闭环。部署后 Atlas 同步端点仍返回 404，本次声明已保留在仓库，待 IAOS 写入路由恢复后补同步。

## 2026-08-01 - AESE M9 到 M10 阶段交接入口补齐

- 变更：Enterprise Genesis 在 M9 `enterprise_operational_ready`、100% 且无未完成工作项时显示完成卡和 `开始 M10 工厂选址与设施规划` 主按钮；深链携带 tenant、case、workspace，M10 返回同一企业。首页“世界地图”改为 `企业生命周期 · M9–M24` 并增加 M10 可发现说明。
- 原因：M10 路由虽然存在，但 M9 终态只显示“当前章节已完成”，没有下一阶段提示或按钮；用户无法发现入口，且手工路由不能证明沿用当前企业上下文。
- 影响：用户完成 M9 后可以在原页面连续进入 M10，M10 读取同一租户案件的实际现金、预算与法人资料；阶段总览仍只是导航，不能绕过 M9 资格门产生业务事实。
- 验证：新增终态组件回归测试先稳定复现“找不到交接标题/按钮”，修复后前端 23 个测试文件/67 项测试、定向 ESLint、TypeScript、生产构建通过；Playwright 在 1440×900、1280×720 和 390×844 三个视口验证交接标题、按钮和完整上下文 URL。全仓 ESLint 仍有 4 个既有 Fast Refresh 警告，与本次文件无关。
- 后续：执行 live 企业终态浏览器验收，并继续 M10 Investigation、Effective Process Artifact、World Observation 和工程执行闭环。
