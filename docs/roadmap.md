# AESE Roadmap

本文件是 AESE 当前里程碑状态和下一步优先级的权威来源。

最后更新：2026-07-22。

## 1. 里程碑状态

| 里程碑 | 目标 | 状态 | 完成证据 |
| --- | --- | --- | --- |
| M0 项目初始化 | 仓库、背景、规则、GitHub | Completed | README、AGENTS、初始提交 |
| M1 虚拟企业蓝图 | 华辰集团、苏州基地、电池冷却板 A 线 | Completed | HCTM Virtual Enterprise Blueprint |
| M2 业务与技术规格 | 对象、事件、seed、演示故事 | Completed (docs) | 4 份 HCTM 规格文档 |
| M2.5 工程治理 | 架构边界、索引、code map、执行规则 | Completed | 本轮治理文档 |
| M3 可执行场景包 | JSON 场景包、校验器、IAOS apply/replay tracer | Completed | pack、CLI、execution evidence、IAOS commits |
| M3V 快速 2D 沙盘 | 七幕故事、22 事件、A 线画布、KPI 和 Agent 建议预览 | Completed | 前端、preview、18 unit/component tests、9 E2E、3 viewport screenshots |
| M4 异常场景运行 | 延期、设备故障、来料不良进入 IAOS 运行链 | Completed | 三类 ingress、状态影响、事务 Outbox、租户/幂等及 canonical replay evidence |
| M5 Agent MVP | 计划、质量、经营分析 Agent | Completed | 9 个受治理只读工具、三 Agent live tracer、跨租户与零业务写入证据 |
| M6 在线 2D 企业沙盘 | IAOS 实时事件、库存、产线、异常和 Agent 运行结果 | Completed | DES-004、PLAN-M6-001、M6 evidence |
| M7 受治理场景运行控制台 | 浏览器预检、初始化、逐幕运行、分析、验证和复位 | Completed | ADR-003、DES-005、PLAN-M7-001、M7 evidence |
| M8 AESE 2.0 基础 | 三态世界、确定性离散事件内核、IAOS 双向桥和最小 Genesis tracer | Completed | PLAN-M8-001、World Play runbook、两仓测试与部署证据 |
| M9 Genesis Incorporation | 注资、法人登记、治理、管理岗位、初始组织与预算 | Completed | hctm-genesis@0.2.0、M9 evidence、IAOS DES-051 |
| M10 Genesis Plant Build | 选址、场地控制、设施项目、公用工程、异常重排与验收 | Completed | hctm-genesis@0.3.0、M10 evidence、IAOS DES-052 |
| M11 Genesis Capability Build | 资金补足、设备/实验室/仓储能力、核心团队与岗位资格 | Completed | hctm-genesis@0.4.0、M11 evidence、IAOS DES-053 |
| M12 Genesis Industrialization | RFQ/定点、产品/工艺、供应商/工装、APQP、试制、PPAP 与量产批准 | Completed | hctm-genesis@0.5.0、M12 evidence、IAOS DES-054 |
| M13 Genesis First Delivery | 正式 O2D、三批交付、客户接受、开票/回款、实际成本与项目毛利 | Completed | hctm-genesis@0.6.0、M13/Genesis evidence、IAOS DES-055 |
| M14 Parameterized Branch Experiments | checkpoint 分支、多周期参数/策略、共同随机数、实验执行与决策证据 | Completed | hctm-genesis@0.7.0、M14 evidence、IAOS DES-056 |
| M15 Governed Strategy Release & Pilot | evidence 审议、版本化发布、shadow、canonical pilot、guardrail 与回滚/采纳 | Completed | hctm-genesis@0.8.0、M15 evidence、IAOS DES-057 |
| X1 System Atlas 全景治理 | 最终完成体、当前状态、依赖与进展历史 | Completed | DES-006、IAOS DES-049、双端动态图谱 |

## 2. 当前阶段

M3、M3V、M4、M5、M6、M7 和跨里程碑的 X1 System Atlas 已完成。联动中心已支持联动检查与受治理场景运行，不依赖 CLI 完成 preflight、initialize、七幕推进、Agent 分析、verify 与 reset。

PLAN-M8-001 至 PLAN-M15-001 均已完成，当前无 active 主实施计划。M15 adopted 与 injected rolled_back 路径均诚实关闭，M14 EvidenceBundle 和 M13 parent checkpoint 保持不可变。

M7 O0-O4 已完成。最终 `m7-acceptance-20260722-05` 从 clean reset 跑通编排 API 与 CLI 对照链：22 个事件、三 Agent、17 条离线业务断言、2 条在线 IAOS 断言和 M6 KPI 均通过；单 run 产生 9 次成功 Tool Call 与两套一致的 O2D Outbox 副作用，UI/CLI 均安全复位。AESE 8090/4173 与 IAOS 8082/3000 的本机部署和健康检查已记录在 M7 evidence。该基线由 M8 强制保留。

M7 的最小成功标准：

1. 浏览器不直接调用 IAOS 写端点，由 AESE 薄编排 API 复用现有 Go 内核。
2. 运行具有 run ID、阶段状态机、plan hash、cursor 和 idempotency key。
3. 用户可从页面初始化、逐幕推进、运行到结束、分析、验证和安全复位。
4. 刷新、断线、重复点击和 AESE 服务重启不会产生重复业务副作用。
5. 权限不足、跨租户、陈旧 cursor 和非法状态转换全部失败关闭。
6. UI 与 CLI 对同一 pack 产生一致的 22 事件、Agent 建议、断言和 KPI。

## 3. M15 当前范围

包含：

- M14 EvidenceBundle/hash 复验、Pareto candidate 审议、利益冲突与独立审批。
- immutable StrategyRelease、semantic diff、Policy/Capability 映射、SafetyEnvelope 和 preflight。
- 4 周或批准窗口的零写入 shadow、active/candidate decision diff 和 assumption drift。
- 新 Genesis canonical operating cycle 中最多 4 周的 allowlist pilot 和正式 IAOS 治理动作。
- hard stop、pause/review、kill switch、prior-release rollback、open commitments 和 compensating action。
- adopted/rejected/rolled_back AdoptionDecision、Strategy Control Room 和完整两仓证据。

不包含：

- 真实客户/生产租户投放、无人审批自动发布或 Agent sole approval。
- 第二客户/产品/工厂、完整 S&OP、动态定价或组织/资本重大变更。
- 删除既成事实式回滚，或用单一 pilot 宣称统计因果和永久最优。

## 4. M15 当前交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| R0 | 决策、发布和安全机器合同 | Completed |
| R1 | Evidence Review 与 ChangeRequest | Completed |
| R2 | StrategyRelease 编译、差异与预检 | Completed |
| R3 | 零写入 Shadow 与运行兼容性 | Completed |
| R4 | Canonical Pilot 与受治理动作 | Completed |
| R5 | Guardrail、暂停、回滚与补偿 | Completed |
| R6 | IAOS 治理、决策与采纳 | Completed |
| R7 | Strategy Control Room 与全链验收 | Completed |

## 5. M15 完成条件

- candidate 与 exact M14 EvidenceBundle/hash 绑定，选择理由、限制、责任人和独立审批可审计。
- StrategyRelease、diff、SafetyEnvelope、RollbackPlan 和 AdoptionDecision 版本化、可 hash、可重放。
- shadow 完整窗口业务写入为零；candidate/active decision diff 和数据 freshness 可下钻。
- pilot 只在批准 scope/window 内通过 IAOS 治理动作，World consequence 与正式业务记录严格对账。
- hard stop、pause、kill switch、rollback、open commitment 和 compensation 失败路径通过。
- tenant/RLS、职责分离、幂等、事务/Outbox、恢复、Control Room、两仓测试及 M3-M14 回归完整。
- 最终输出 `strategy_change_cycle_closed=true` 与 `adopted|rejected|rolled_back`，不存在未处理 breach 或未对账 commitment。

## 6. M15 风险与依赖

- G4-G8 未冻结 candidate、release、shadow/pilot、guardrail/rollback 和 IAOS gap 前，不得激活 Policy 或创建 pilot 业务事实。
- M14 推荐仍只是模拟证据；不得为了推进计划预设 winner、隐藏限制或把 Pareto 误写成唯一最优。
- shadow 必须零业务写入，shadow approval 不能自动穿透 pilot approval。
- canonical pilot 会产生真实的虚拟企业业务后果；reset/rollback 不能删除已提交订单、库存、发运、发票或现金。
- kill switch 只停止未来动作；未结承诺必须进入 ledger 并用受治理 compensation 处理。
- 单一 pilot 无随机对照，不得宣称因果提升；超出 M14 assumption support 时必须 pause/review。
- 当前工作区已有其他人的测试修改、截图删除和验收产物，实施 agent 必须保留并避免重叠修改。

## 7. M14 已完成范围

包含：

- M12/M13/M14 批准 checkpoint allowlist、祖先/hash 和 opening reconciliation。
- 需求、供应、设备、质量和付款外生参数，以及库存/供应、产能/维护、资金保护策略。
- 固定版本 PRNG、命名随机流、seed set、共同随机数和 paired comparison。
- 12 个虚拟周/订单周期、分支隔离、持久 run catalog、有界执行、取消/继续/重试和配额。
- OTIF、积压、库存/营运资金、现金低点、毛利、质量、加班、报废、加急和恢复 KPI。
- Constraint/Pareto/EvidenceBundle、IAOS 实验治理与推荐边界，以及 Scenario Lab。

不包含：

- 自动把推荐应用到正式 Policy、预算、订单、采购、排产或现金。
- 用单次运行或未经校准的参数宣称真实因果、概率、预测精度或最优策略。
- 第二客户/产品/工厂、完整 S&OP、真实数据校准、机器学习训练或通用分布式计算平台。

## 8. M14 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| X0 | 实验方法、基线和机器合同 | Completed |
| X1 | 确定性随机流与参数矩阵 | Completed |
| X2 | Checkpoint fork、隔离分支与持久运行目录 | Completed |
| X3 | 多周期 World、策略执行和经济守恒 | Completed |
| X4 | 有界实验执行器与生命周期治理 | Completed |
| X5 | KPI 聚合、比较与 EvidenceBundle | Completed |
| X6 | IAOS 实验治理与推荐边界 | Completed |
| X7 | Scenario Lab 与全链验收 | Completed |

## 9. M14 完成条件

- 父 checkpoint 不变、兄弟分支/租户隔离，正式 IAOS 经营事实零污染。
- checkpoint、参数、策略、PRNG、seed、规则和聚合全部版本化、可 hash、可重放。
- 同输入 100 次 hash 一致，共同随机数和 paired comparison 可自动验证。
- 多周期数量、质量、资源、应收/现金和 actual cost/margin 守恒，失败样本不被过滤。
- 执行器默认 dry-run；显式 apply、有界并发、配额、取消、继续、重试和崩溃恢复通过。
- EvidenceBundle、IAOS 权限/职责/Outbox/RLS、Scenario Lab、runbook/evidence 和 M3-M13 回归完整。
- 只有证据完整、约束完成且无未解释运行缺失时输出 `strategy_evidence_ready=true`。

## 10. M14 历史风险与控制

- M14 实施时在 G4-G8 冻结 checkpoint、分布/seed、策略/KPI、运行容量和 IAOS gap 后才进入 IAOS 写端点开发。
- 场景分布是虚构假设，不能包装成真实概率；单次 run 只用于 tracer，不用于稳健性结论。
- 参数笛卡尔积可能失控；所有 apply 前必须预估 run 数、时间和存储并受配额约束。
- 不同策略必须使用共同随机数并保留失败/取消样本，否则比较存在选择偏差。
- 分支、transaction/correlation/idempotency namespace 必须隔离，不能污染父 checkpoint、兄弟分支或正式 IAOS 数据。
- 推荐与批准/投放必须分离；Agent、人或 UI 都不能从实验结果直接修改正式业务策略。
- 当前工作区已有其他人的测试、截图和生成物改动，实施 agent 必须保留并避免重叠修改。

## 11. M13 已完成范围

包含：

- 从零可销售成品库存接收 10,000 件正式订单和 2,000 件追加需求。
- M12 release manifest 到 Genesis-specific O2D 的兼容适配，以及 ATP/MRP、采购、IQC、生产、质量和库存。
- 供应延期、设备停机和末批 300 件短缺的受治理恢复。
- 9,000、2,700、300 三批实际发运、客户收货/接受和 delivery close。
- 开票、应收、客户实际付款、银行到账、收款核销和合同负债处理。
- 材料、人工、设备/能源、报废、质量、加急、物流和分摊实际成本，以及订单/项目毛利。
- First Delivery Play、Project Genesis M9-M13 terminal-hash 链和 M3-M12 强制回归。

不包含：

- 第二订单/客户/产品、多工厂或长期滚动 S&OP/MPS。
- 完整总账、税务申报、资金管理、真实银行/电子发票接口和法定报表。
- 售后、退货、质保、贷项、坏账和复杂跨期收入/成本会计。
- 参数化分支、Monte Carlo、A/B 和长期经营实验；属于 M14。

## 12. M13 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| E0 | 订单、履约、财务与成本机器合同 | Completed |
| E1 | M12 terminal 到 Genesis O2D 兼容适配 | Completed |
| E2 | 正式订单、ATP/MRP 与供应执行 | Completed |
| E3 | 正式生产、质量、库存与实际成本 | Completed |
| E4 | 分批发运、客户接受与 300 件恢复 | Completed |
| E5 | 开票、应收、银行回款与项目盈亏 | Completed |
| E6 | IAOS O2D、财务与成本治理 | Completed |
| E7 | 统一 Agent 与 First Delivery Play | Completed |
| E8 | 全链验收与 Project Genesis 收口 | Completed |

## 13. M13 完成条件

- 从 M12 terminal contract 到 `first_commercial_cycle_closed` 可确定性运行、恢复、重放和分层复位。
- 正式需求 12,000、采购/生产/报废/库存、三批发运和客户接受数量严格守恒。
- 300 件短缺形成 observation、受治理恢复、第三批交付和 discrepancy close 完整链。
- 客户接受、发票/应收、银行到账、收款核销、收入和现金严格分离且金额守恒。
- M12 财务结转、订单实际成本、标准/实际差异、项目毛利和 closing cash 可解释对账。
- 两仓权限、RLS、Outbox、幂等、API/UI、runbook/evidence、Project Genesis 总报告以及 M3-M12 回归全部通过。

## 14. M13 历史风险与控制

- M12 1,200,000 CNY 试制成本的支付/应付状态和 2,000,000 CNY 合同负债履约处理已在 G4 冻结；全程禁止通过改 opening cash/利润静默配平。
- 旧 HCTM 场景包含 1,200 件 opening inventory 和已发生事件，M13 只能复用语义/能力，不能继承库存或交易历史。
- 发运不等于客户接受，发票不等于现金，毛利不等于现金余额；UI 和 IAOS 状态不能伪造 World/银行事实。
- 实施时在 G4-G9 关闭后才进入 IAOS 写端点开发，E6 使用了独立 IAOS branch/worktree。
- actual cost 缺任一强制要素时经营分析必须失败或保持不完整，不能回退到估算后宣称盈利。
- M13 reset 必须保留 M9-M12 L1 事实和旧场景数据，只清理本次运行的 L2/L3 对象。
- 当前工作区已有其他人的测试、截图和生成物改动，实施 agent 必须保留并避免重叠修改。

## 15. M12 已完成范围

包含：

- 单一虚构客户的 RFQ、可行性、成本/报价、客户定点、开发资金和客户项目。
- `HCTM-BCP-A01` 产品 revision、EBOM/MBOM、routing、process flow、PFMEA、control plan 和 work instruction。
- 关键供应商、产品专用工装/量检具、首批材料、批次/证书、IQC 和追溯。
- 简化 APQP 六个 gate、两轮试生产、尺寸/功能/泄漏、MSA、节拍、良率和 Cp/Cpk。
- 一条焊接泄漏/Cpk 不足的 discrepancy、Knowledge、工程变更、复试和 PPAP tracer。
- Industrialization Play、旧 HCTM stable-code/hash 兼容，以及 M7-M11 强制回归。

不包含：

- 第一张正式销售订单、正式 MRP、批量采购和正式量产；属于 M13。
- 客户发运、发票、应收、回款、实际量产成本和项目盈亏；属于 M13。
- 第二客户/产品、多工厂、多版本并行量产或完整 CRM/CPQ/PLM/QMS/SRM。
- CAD/CAE、高精度物理仿真、真实外部接口或参数化分支实验。

## 16. M12 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| D0 | 客户项目、产品、质量与资金机器合同 | Completed |
| D1 | RFQ、报价、定点与项目资金 | Completed |
| D2 | 产品、工艺与 APQP 版本治理 | Completed |
| D3 | 供应商、工装与首批物料世界 | Completed |
| D4 | 试生产、质量异常与 PPAP 世界 | Completed |
| D5 | IAOS 客户项目、工程、质量与放行治理 | Completed |
| D6 | 统一岗位与 Industrialization Play | Completed |
| D7 | 全链、安全、恢复和回归验收 | Completed |

## 17. M12 完成条件

- 从 M11 terminal contract 到 `serial_production_eligible` 可确定性运行、恢复、重放和安全复位。
- RFQ/报价/定点、开发资金/合同负债、项目预算、试制成本和现金语义正确且守恒。
- 产品/BOM/routing/PFMEA/control-plan/工装/供应商/材料 revision 一致并可机器校验。
- 试制物料、样件、设备、人员、测量、报废、成本和追溯满足守恒。
- 焊接泄漏/Cpk 异常形成 observation、受治理变更、第二轮试制、问题关闭和客户 PPAP 批准完整链。
- 两仓权限、RLS、Outbox、幂等、API/UI、runbook/evidence、M3/O2D 兼容以及 M7-M11 回归全部通过。

## 18. M12 风险与依赖

- M11 closing cash 为 8,500,000 CNY，并需保留 3,000,000 工资准备金和 5,000,000 最低缓冲；G5 必须冻结客户工装预付款/开发预算，且预付款只能记合同负债，不能计收入。
- 现有 `scenario-packs/hctm` 已预置同名产品/BOM/routing，但只是兼容 fixture；M12 必须产生独立 release manifest/hash，禁止把旧 seed 当完成证据。
- IAOS 工程/APQP/PPAP 状态不等于实际试制能力或客户批准；只有 World consequence 和客户 observation 可推进现实状态。
- G4-G8 未关闭前不得进入 IAOS 写端点开发；D5 必须建立新的独立 IAOS branch/worktree。
- 试制件、PPAP 样件和报废默认不可销售，不能静默成为 M13 库存或收入。
- M12 不得提前接收正式订单、运行正式 O2D 或宣称第一批交付完成。
- 当前工作区已有其他人的测试、截图和生成物改动，实施 agent 必须保留并避免重叠修改。

## 19. M11 已完成范围

包含：

- M9 剩余认缴资本的受治理催缴/实际到账、新 CAPEX/headcount envelope、设施尾款、工资准备金和现金缓冲。
- 冷却板 A 线通用能力需求，以及成形、CNC、激光焊接、清洗、检漏、装配/包装、质量实验室的采购/租赁、交付、安装、调试、校准和验收。
- M10 空间/utility 约束下的设备落位、实验室、仓储和最小一班制能力。
- 厂长、计划、采购、质量、工艺、设备、仓储、操作工和检验员的编制、招聘、到岗、培训、资格和班次。
- 一条检漏设备校准漂移的 discrepancy、Knowledge、整改、复验和关闭 tracer。
- Capability Build Play，以及 M7-M10 强制回归。

不包含：

- 客户 RFQ、正式产品/BOM/routing、成本报价和 APQP；属于 M12。
- 产品专用工装、首批物料、试生产、PPAP、SOP 和量产；属于 M12。
- 第一张正式订单、开票、回款和实际项目盈亏；属于 M13。
- 完整采购/SRM、HRIS/LMS、薪酬、EAM、固定资产会计、融资或 3D 产品。

## 20. M11 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| C0 | 能力、设备、人员与资金机器合同 | Completed |
| C1 | 资金与受治理设备采购 | Completed |
| C2 | AESE 设备、实验室与仓储能力世界 | Completed |
| C3 | AESE 人员、技能与班次世界 | Completed |
| C4 | IAOS 采购、资产、组织与资格治理 | Completed |
| C5 | 统一岗位与 Capability Build Play | Completed |
| C6 | 全链、安全、恢复和回归验收 | Completed |

## 21. M11 完成条件

- 从 M10 terminal contract 到 `industrialization_eligible` 可确定性运行、恢复、重放和安全复位。
- 资本实际到账、新预算、设施尾款、设备/租赁承诺、工资准备金和现金严格守恒。
- 关键设备全部实际交付、安装、调试、校准、安全和能力验收，空间/utility 不超限。
- 最小核心团队实际到岗且岗位资格有效，一班制、替补和职责隔离门通过。
- 检漏设备失败形成 observation、Knowledge 差异、受治理整改、复验和关闭的完整因果链。
- 两仓权限、RLS、隐私、Outbox、幂等、API/UI、runbook/evidence 以及 M7-M10 回归全部通过。

## 22. M11 风险与依赖

- M10 closing cash 只有 10,000,000 CNY 且存在设施遗留承诺，不能无资金来源生成整线；G6 必须先冻结剩余资本实缴、采购/租赁组合和现金缓冲。
- 资产登记、员工档案和培训记录不等于实际设备能力、人员到岗或技能掌握；只有 World consequence 和验收事实可推进能力状态。
- M8 `LAS-WLD-02` 是独立设备 tracer，未经采购/交付/验收迁移证据不能计入 M11 生产资产。
- G4-G7 未关闭前不得进入 IAOS 写端点开发；C4 必须建立新的独立 IAOS branch/worktree。
- 候选人与员工完全虚构且按 actor scope 投影；不得引入真实个人信息或让无关角色读取候选详情。
- M11 只交付通用设备和人员能力，不得提前宣称产品、APQP、试生产、PPAP 或 SOP 完成。
- 当前工作区已有其他人的测试、截图和生成物改动，实施 agent 必须保留并避免重叠修改。

## 23. M10 已完成范围

包含：

- 至少三个虚构候选场址及资金、工期、物流、人力、公用工程、风险和扩展性评估。
- 项目负责人推荐、CEO/CFO 受治理的选址与投资批准，以及租赁/场地使用控制。
- 设施项目、WBS、承包商资源日历、合同承诺、变更、里程碑、付款和验收。
- 区域/城市/园区/场地/建筑/楼层/功能区的最小空间层级和公用工程容量。
- 一个固定公用工程接入延期的 discrepancy、Knowledge、重排和关闭 tracer。
- World Play 工厂建设 campaign，以及 M7/M8/M9 强制回归。

不包含：

- 生产设备和检测仪器采购、安装及调试；属于 M11。
- 招聘、培训、工艺能力、试生产和投产门；属于 M11。
- BIM、3D、自由布局编辑器、真实地图/园区/承包商接口。
- 完整地产、总账、税务、融资或通用项目管理产品。

## 24. M10 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| P0 | 场址、空间与项目机器合同 | Completed |
| P1 | 受约束选址决策 | Completed |
| P2 | AESE 设施项目与空间世界 | Completed |
| P3 | IAOS 投资与项目治理 | Completed |
| P4 | 统一角色与 Plant Build Play | Completed |
| P5 | 全链、安全、恢复和回归验收 | Completed |

## 25. M10 完成条件

- 至少三个候选先经过硬约束，再产生版本化、可解释的多维评分和受治理决策。
- 从 M9 terminal contract 到 `capability_build_eligible` 可确定性运行、恢复、重放和安全复位。
- IAOS 项目记录与 AESE 现场事实严格分离；时间到期或管理记录不能凭空完成工程。
- 预算、合同承诺、实际付款和现金分别守恒，未验收、越权或超预算动作失败关闭。
- 公用工程延期能形成 observation、Knowledge 差异、受治理重排和新实际结果的完整因果链。
- 两仓权限、RLS、Outbox、幂等、API/UI、runbook/evidence 以及 M7/M8/M9 回归全部通过。

## 26. M10 风险与依赖

- M9 只有 20,000,000 CNY 实际现金和 15,000,000 CNY 首年预算，首版不得假设可支付绿地自建；候选基线预计收敛到租赁标准厂房改造，但必须由规则和审批得出。
- IAOS project/milestone 状态不等于现场实际进度；只有 AESE World consequence 和验收事实可以推进物理状态。
- G4/G5 已关闭；P3 在独立 `feat/m10-plant-governance` worktree 完成，revision `23be02a`。
- 承包商、公用工程方、园区和验收机构是确定性外部世界策略，不是 IAOS 用户，也不能绕过 observation/intent/outcome 合同。
- M10 只交付设施载体，不得提前引入生产设备、人员或投产能力，避免侵入 M11。
- 当前工作区已有其他人的测试、截图和生成物改动，实施 agent 必须保留并避免重叠修改。

## 27. M9 已完成范围

包含：

- 单一虚构投资主体、集团和苏州制造法人。
- 设立方案、外部登记、开户、首期资本到账和成立费用。
- CEO、CFO、工厂项目负责人任命、接受、mandate 和知识边界。
- 初始组织、首年 budget envelope、审批和可支用资格。
- World/IAOS/Knowledge 三态、资金守恒和受治理 bridge 全链。
- World Play 企业成立 campaign 与 M7/M8 强制回归。

不包含：

- 工厂选址、土地/租赁和建设执行；属于 M10。
- 设备采购、招聘培训和投产门；属于 M11。
- RFQ/APQP/PPAP、首批交付和参数化实验。
- 完整总账、税务、复杂融资或真实监管/银行接口。

## 28. M9 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| I0 | 业务基线、所有权和机器合同 | Completed |
| I1 | AESE 成立世界、外部策略与经济规则 | Completed |
| I2 | IAOS 法人、治理、岗位和预算能力 | Completed |
| I3 | CEO/CFO 统一岗位运行 | Completed |
| I4 | Genesis Incorporation Play | Completed |
| I5 | 全链、恢复、安全和回归验收 | Completed |

## 29. M9 完成条件

- pre-incorporation 到 `plant_project_eligible` 可确定性运行、恢复和重放。
- 法人、账户、资本、任命和预算的三态及因果链可机器验证。
- 投资人/公司现金、认缴、实缴、预算和实付严格区分并满足守恒。
- 人类与 Agent 使用相同 IAOS Capability、Policy、权限和审计。
- IAOS journal/Outbox 与业务提交原子，回滚不产生 committed outcome。
- tenant、幂等、并发、断线和重复执行测试通过，M7/M8 零回归。

## 30. M9 风险与依赖

- M8 初态中的 10,000,000 CNY 当前没有显式 owner；M9 必须通过 pack version 和资本事件迁移，禁止静默当作公司现金。
- 法人档案 committed outcome 不等于外部登记已经生效；监管和银行结果由 AESE 确定性世界策略产生。
- 预算是支出授权，不是现金；认缴不是实缴；管理界面必须明确区分。
- M9 IAOS 修改必须新建独立 worktree，并在 I0 合同冻结后启动。
- 自动 Agent 不得批准自身预算、伪造外部结果或绕过岗位 mandate。
- 当前工作区已有其他人的测试与截图改动，实施 agent 必须保留并使用不重叠 worktree。

## 31. M8 已完成范围

包含：

- World run、虚拟时钟、离散事件、规则版本、日志、快照和确定性 replay。
- World / IAOS / Actor Knowledge 三态与显式 discrepancy。
- 单台 `LAS-WLD-02` 设备退化的最小发现、登记、处置和关闭 tracer。
- observation / intent / committed outcome 的受治理 IAOS 合同。
- 设备、人员、物料和最小现金守恒。
- `hctm-genesis` world pack 和现有 `order-expedite-01` 兼容适配。
- World Play 的时间控制、三态对照和差异时间线。

不包含：

- 一次完成公司设立、工厂建设、APQP、财务和全生命周期的所有模块。
- 多工厂、多产品族、数百 Agent、自由文本长期记忆和自主批准。
- 3D、高精度物理仿真或完整设备数字孪生。
- 绕过 IAOS 权限/Capability/Outbox 的业务写入。
- 在 AESE 中镜像 IAOS 企业管理数据库。

M8 决策门与 F0-F5 的任务、验收和跨仓顺序以 PLAN-M8-001 为准。后续 Project Genesis 分解为 M9-M13，参数化分支实验后移至 M14。

## 32. M8 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| F0 | 基线、状态所有权、存储和桥接合同冻结 | Completed |
| F1 | 确定性仿真内核 | Completed |
| F2 | 三态与设备偏差 tracer | Completed |
| F3 | 受治理 IAOS 双向桥 | Completed |
| F4 | Genesis world pack 与旧场景兼容 | Completed |
| F5 | World Play 最小界面与全链验收 | Completed |

## 33. M8 架构风险与依赖

- ADR-004 已 accepted；实现必须遵守独立 PostgreSQL database/账号/迁移边界，禁止跨库查询和外键。
- AESE 仿真事实和 IAOS 管理事实必须物理/逻辑隔离，禁止共享表和跨库写入。
- 原始设计稿中的 Spring Boot 仅是模块示意；AESE 实现继续使用 Go，与现有工具链保持一致。
- 世界结果必须由版本化规则和资源守恒计算；Agent 只提交 intent，不能直接改写 World State。
- Actor Knowledge 必须遵守可见范围，不能为了方便让 Agent 读取全量客观世界。
- IAOS 修改必须在独立 worktree，并先完成权限、RLS、Outbox、幂等和无部分写入设计。
- M7 22 事件、三 Agent、Preview/Live 与 reset 是强制回归门。

## 34. M8 完成条件

- ADR-004 accepted，World/IAOS/Knowledge 所有权和 World Store 选型明确。
- 相同 pack、规则版本、seed 和输入可重复产生相同 event log、state hash 与 KPI。
- 设备退化 tracer 可展示世界变化、IAOS 未登记、角色未知及其发现/关闭过程。
- IAOS 双向桥通过租户、权限、幂等、乱序、失败恢复和 Outbox 审计验收。
- Genesis pack 可离线验证、初始化、推进、复位和 replay，旧 M7 场景不回归。
- API/UI/runbook/evidence 与两仓 revision 完整。

## 35. M7 已完成范围（保留基线）

包含：

- 无独立数据库的 AESE scenario orchestration API。
- pack 阶段编译、dry-run 影响、状态机、幂等和恢复。
- initialize、advance、run-to-end、analyze、verify 和 reset 编排。
- 联动中心场景运行视图、七幕 stepper、日志和危险确认。
- 权限、跨租户、并发、断线、重启及 CLI/UI 一致性验收。

不包含：

- 参数化 A/B 实验、并行分支和第二条故事。
- 真实 LLM 自主 Agent 和建议自动执行。
- AESE 业务数据库、通用任务队列或工作流引擎。
- 完整成本核算、3D 工厂和布局编辑器。

## 36. M7 交付切片

| Slice | 内容 | 状态 |
| --- | --- | --- |
| O0 | 状态机、阶段合同和编排内核重构 | Completed |
| O1 | AESE 薄编排 API | Completed |
| O2 | IAOS 运行记录、权限和并发补强 | Completed |
| O3 | 可视化场景运行控制台 | Completed |
| O4 | 全链路、恢复、安全和三视口验收 | Completed |

## 37. 历史风险与依赖

- IAOS Platform、PostgreSQL、NATS 和 O2D 可运行，`tenant-hctm` 的 work_order metadata、workflow config 和 tracer 数据已完成 seed/apply。
- `/iaos/iaos-go` 当前主分支本地领先远程，任何集成开发必须使用独立 worktree 并确认基线。
- HCTM 业务字段与 IAOS 现有 legacy `sales_order` 物理模型存在差异，需要兼容性报告，不能直接假设可导入。
- IAOS 当前 O2D 自测硬编码 `tenant-001` 和旧订单，需要避免污染 HCTM tracer。
- 正式事件必须走 Outbox 或受治理 ingress，不能把 direct NATS 当最终实现。
- IAOS 受治理 scenario apply/reset 已实现 M3 allowlist；订单确认 CAS、workflow/event 去重、跨节点原子事务及 work_order API 已实证。
- legacy 表没有全面 FORCE RLS；M3 scenario adapter 已在所有查询/更新/删除中显式绑定 tenant。平台长期仍应继续推进全表 RLS FORCE hardening。
- M4 的显式 tenant predicates 已关闭当前入口越界路径，但 tenant-safe composite foreign key 和 metadata `version` 的平台级排序仍需后续 hardening。
- 首版 2D 沙盘是确定性预览，不应被描述为 IAOS 实时运行结果；界面必须显示 Preview 数据源状态。
- `preview.json` 只承载视图状态和 delta，不得复制 MRP、排产或 Agent 决策逻辑。
- 当前通用 `/api/v1/events/stream` 无持久 cursor 且缓冲满会丢事件，只能作为监控通道，不能直接作为 M6 恢复合同。
- 完工和发运是内生业务动作，不能为了复用现有接口而错误接入 simulation ingress。
- HCTM 尚无已批准的成本金额基线；在基线确认前，在线经营分析的成本部分必须保持 `partial`。
- M7 新增 HTTP 服务但不得拥有业务数据库；运行恢复必须以 IAOS run/snapshot/event/recommendation 为事实。
- 当前场景使用固定自然键，同一 tenant/scenario 首版只能有一个可写 active run。
- 浏览器不得直接编排多个 IAOS 写 API；所有危险动作需要服务端权限、幂等和确认合同。

## 38. M7 完成条件

- 非研发用户可从浏览器完整运行并复位第一条故事。
- UI 状态只在 IAOS committed/no-op 和 snapshot cursor 证实后推进。
- 双击、并发、断线、刷新和服务重启均不重复推进阶段。
- reset 影响可预览，一次性 confirmation token 不能重放，L1 始终保留。
- IAOS 与 AESE 两仓权限、测试、部署、runbook 和 evidence 完整。

## 39. M6 完成证据

M6 已满足：

- 在线库存、完工、发运和订单状态满足 expected outcomes 与库存守恒。
- snapshot 与 cursor 来自一致观察边界，断线补发和事件去重可复现。
- Preview/Live 数据源明示，错误时不静默混用。
- Agent 建议、对象引用和 Tool Call 证据属于同一 tenant/correlation。
- IAOS 与 AESE 两仓测试、部署、runbook 和 evidence 完整。
