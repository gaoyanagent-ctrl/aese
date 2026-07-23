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
