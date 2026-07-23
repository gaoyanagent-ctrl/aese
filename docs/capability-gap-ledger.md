# M8 Capability Gap Ledger

| 企业活动 | 角色 | 对象 | IAOS 权限/Capability | 当前处理 | 优先级 |
| --- | --- | --- | --- | --- | --- |
| 上报设备观察 | 设备工程师或 Agent | `LAS-WLD-02` | `world.observation.ingest` | World Bridge observation journal + Outbox | P0 已交付 |
| 提出检修 | 设备工程师或 Agent | 维护工单 stable ref | `world.intent.create` | World Bridge intent journal；人和 Agent 共用相同端点与授权资源 | P0 已交付 |
| 批准并提交检修结果 | 有审批权的人员/Agent | 维护工单 committed refs | `world.intent.approve` | 只接受 committed/no-op outcome，journal 与 Outbox 同事务 | P0 已交付 |
| 恢复桥接消费 | AESE runtime | world run | `world.bridge.read` | cursor query 为事实，SSE 仅提示 | P0 已交付 |
| 通用 EAM 工单编排 | EAM 专员 | 任意设备 | 尚无通用 Capability | 本里程碑只覆盖 Genesis 固定 tracer，不扩成 EAM 引擎 | P1 后续 |
| 自动审批高风险维修 | 自主 Agent | 维修预算/停线 | 不允许 | 必须由 IAOS Policy/Approval 决定，AESE 不代批 | 禁止 |

长期 JWT、AESE 直写 IAOS 数据库、正式 direct NATS 和 webhook-only 路径均不是临时方案。

## M9 成立治理（IAOS-native 1.1）

| 企业活动 | 角色 | IAOS 权限 | 实现/缺口 |
| --- | --- | --- | --- |
| 批准设立方案 | `founder-principal` | `founder.resolution.approve` | Effective Runtime Artifact + G1 + Process/Decision/Journal/Outbox，已交付 |
| 工商登记/补正 | incorporation Agent + founder | `registration.*` | G2、World Intent/Observation/CommittedOutcome；原 correlation 补正，已交付 |
| 开户与出资 | finance Agent + founder | `bank.account.*`、`capital.*` | G3/G4；拒绝或差异保持状态并产生 Discrepancy，已交付 |
| 组织/任命 | governance Agent + founder | `organization.*`、`executive.*` | G5 与候选人 World observation，已交付 |
| Mandate/预算/readiness | governance/finance/audit Agent + founder | `operating.mandate.*`、`initial.budget.*`、`enterprise.readiness.*` | G6/G7、职责分离和终态 evaluator；对象级守恒明细仍在最终验收 |
| 双仓恢复对账 | AESE Bridge | `world.bridge.read` | 五类 reconciliation CLI 已交付；重启/延迟矩阵待关闭 |

旧 `genesis.*` allowlist receipt 仅为迁移来源，不再作为 M9N 完成证据。真实监管、
银行、完整法人主数据和总账仍为非目标；外部结果由 AESE 确定性策略产生。

## M10 工厂投资与项目治理

| 企业活动 | 角色 | IAOS 权限 | 实现 |
| --- | --- | --- | --- |
| 场址评估/批准 | 项目负责人、CEO/CFO | `genesis.site.evaluate/approve` | 三候选硬约束与受治理批准 |
| 项目执行/重排 | 项目负责人、审批人 | `genesis.project.execute/rebaseline` | expected version、幂等与 committed outcome |
| 项目/设施接受 | 验收人与项目审批人 | `genesis.project.accept` | 引用 World 验收证据，不直接改现场进度 |
| 付款批准 | CFO/独立审批人 | `genesis.payment.approve` | 预算、现金、里程碑验收与禁止自批 |

通用项目管理、合同管理、总账、EAM、BIM 和 M11 设备/人员能力仍为明确缺口，不由 M10 最小 allowlist 伪装覆盖。

## M11 生产能力治理

`genesis.capability.fund/procure/accept`、`genesis.workforce.plan/hire/qualify` 与 `genesis.capability.payment.approve` 已交付。设备接受和人员资格强制引用 World evidence；候选隐私、职责隔离、预算/现金、验收付款、RLS、幂等与并发失败关闭。通用 SRM/HRIS/LMS/EAM、薪酬、固定资产会计及 M12 产品工业化仍为后续缺口。

## M12 产品工业化治理

`genesis.industrialization.quote/design/process/trial/quality/ppap/release` 已交付。质量、PPAP 与放行强制引用 World/customer evidence；版本篡改、伪造测量/客户决定、未批准试制和重复放行失败关闭。完整 CRM/CPQ/PLM/QMS/SRM 与 M13 正式 O2D/财务闭环仍为后续缺口。

## M13 首次商业交付治理

`genesis.delivery.order/plan/procure/produce/ship/accept/invoice/collect/cost/close` 已交付。客户接受、银行到账、实际成本和周期关闭强制引用 World evidence。完整总账、税务、售后、坏账和 M14 多周期实验仍为后续缺口。
