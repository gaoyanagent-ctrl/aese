---
id: PLAN-M10-INTERACTIVE-001
title: M10 Agent 候选与参数化工厂建设交互闭环实施计划
date: 2026-08-01
status: active
author: Codex + User
tags: [m10, agent, plant-build, interactive, governance]
---

# M10 Agent 候选与参数化工厂建设交互闭环实施计划

## 1. 目标与边界

把既有 `hctm-genesis@0.3.0` 固定十帧 reference replay 升级为真实交互闭环：用户填写设施需求和金额边界，`plant-planning-agent` 生成候选，项目负责人选择调研项，World 返回外部事实，人员与审批流决定选址、投资、项目、WBS、变更、付款和验收。旧回放迁入 `fixture_only`，不再自动成为业务事实。

AESE 拥有 Agent 候选生成适配、World 调研/施工事实和游戏交互；IAOS 拥有设施需求、proposal/review、Capability、Process、Approval、Agent Run、资金事实、审计和 Outbox。浏览器写操作统一经过 AESE Command Gateway。

## 2. 纵向切片

### S0 计划与差距审计

- [x] S0.1 修订 DES-011，冻结 Agent/人工/World/IAOS 所有权与金额边界。
- [x] S0.2 识别固定候选、金额、赢家、WBS、空间和异常所在代码、fixture、API 与 UI。
- [x] S0.3 建立本计划、路线图、文档索引、Code Map 和 Atlas 入口。

### S1 参数化需求与 Agent proposal 合同

- [x] S1.1 定义 `FacilityRequirement`、`SiteOptionProposal`、`ProposalSet`、`ProposalReview` 严格 JSON 合同。
- [x] S1.2 金额使用 decimal string + currency；日期 RFC3339；候选数量、权重和枚举有边界校验。
- [x] S1.3 接入 MiniMax JSON provider 和 status/proposal API；未配置模型时失败关闭，不返回伪 Agent 候选。人工表单在 S3 接入。
- [x] S1.4 记录模型、prompt 版本、输入/输出 hash、request ID、token、校验和错误证据；AESE CreativeJob 保存生成技术证据，业务 proposal 后续由 IAOS 保存。
- [x] S1.5 单元测试覆盖合法、金额非法、重复候选、无来源、模型坏 JSON、重试和未配置模型。

### S2 IAOS 权威对象与受治理能力

- [ ] S2.1 在 IAOS 独立 worktree 发布 Facility Requirement、Proposal、Review 与调查请求 Entity/storage contract。首个切片已实现前三类 FORCE RLS 权威表、Capability 写入触发器和版本/幂等表；调查请求及平台 Entity/Edition 发布仍未完成。
- [ ] S2.2 发布 `facility.requirement.define`、`site.proposal.record`、`site.proposal.review`、`site.investigation.request` Capability Artifact 与实现绑定。前三项已随 `genesis-m9@1.6.0` 安装为 Published Active Artifact/native binding；调查请求仍未完成。
- [ ] S2.3 Process Artifact 展开为 human task → agent task → human review → world wait，不允许 Go 固定帧跳步。
- [ ] S2.4 proposal/review、Agent Run、audit、journal 和 Outbox 同事务；RLS、幂等、版本冲突和失败无部分写入。Requirement/Proposal/Review 已具备事务 Audit、Outbox、RLS、幂等与版本冲突门；Agent 技术证据当前由 AESE CreativeJob 保存，尚未纳入 IAOS Agent Run/Process 事务。

### S3 用户表单与人工选择

- [x] S3.0 M9 `enterprise_operational_ready` 终态显示 M9→M10 交接卡与主按钮，携带当前 tenant/case/workspace；首页提供可发现的 `企业生命周期 · M9–M24` 总览入口。
- [x] S3.1 Plant Build Play 增加需求表单、字段解释、来源标识和金额/日期/候选边界校验；实际现金与已批预算只读显示 IAOS 来源引用。
- [x] S3.2 候选卡展示依据、金额/工期区间、假设、未知事实、风险、来源、置信度和 Agent 生成证据。
- [ ] S3.3 支持采纳调研、退回重生成、人工新增和淘汰；所有决定要求理由。采纳/退回/淘汰已通过 BFF 持久化到 IAOS；人工新增仍是明确标识的本地草稿，尚不能提交调查或成为权威候选。
- [x] S3.4 页面提供可见“功能说明”，解释 M9 资格/资金 → Requirement → Agent Proposal → Human Review → IAOS/World 后续链；真正的 Entity/Capability/Process/Evidence 深链随 S2.3/S4 补齐。

### S4 World 调研、参数化项目和资金治理

- [ ] S4.1 World 外部参与者返回报价、权属、容量、许可和可用日期 Observation。
- [ ] S4.2 评分只消费已送达事实；Agent 估算与外部事实并列显示，不得互相覆盖。
- [ ] S4.3 空间、承包策略、WBS、延期缓解和投资变更均采用 Agent 建议 + 人员决定。
- [x] S4.4a IAOS 从已过账银行科目与设立案件已批预算读取只读财务快照，返回来源引用和 snapshot hash；AESE BFF 不接受页面伪造该快照。
- [ ] S4.4b 投资、合同、变更和付款金额可修订但必须重新校验/审批；资金变化使旧 Requirement snapshot 失效并要求修订。
- [ ] S4.5 AP/CIP/付款/验收按财务 DES-033/034 接通，不以治理 JSON 冒充会计事实。

### S5 Fixture 隔离、迁移与全链验收

- [ ] S5.1 将三个旧候选和固定数值迁入显式 `fixtures/reference-replay/`，生产 API 禁止默认加载。
- [ ] S5.2 外部模型与纯人工两条路径达到相同业务终态；模型输出不同不要求相同 hash，已接受 snapshot 重放必须确定。
- [ ] S5.3 覆盖断线、重启、重复点击、并发审阅、模型失败、无有效候选、单一来源例外和人工接管。
- [ ] S5.4 完成两仓测试、构建、部署、Runbook、现场 UI/API/DB 证据、Atlas、提交和推送。

## 3. 发布门

- 未完成 S2，不得把 AESE 内存/文件草稿称为 IAOS 业务记录。
- 未完成 S3，不得把 JSON API 称为用户可操作功能。
- 未完成 S4，不得把 Agent 估算称为报价、现金、预算、合同或付款事实。
- 未完成 S5，不得把 M10 状态从 `Interactive Revision Pending` 改为 Completed。

## 4. 当前执行

S0、S1、S3.0、S3.1、S3.2、S3.4 和 S4.4a 已完成。当前首个在线纵切为：M9 机器终态交接 → IAOS 只读财务快照 → 人员填写 Facility Requirement → AESE 调用真实 Agent → IAOS 保存 Requirement/Proposal → 人员提交采纳/退回/淘汰 Review。该纵切不等于完整 M10：S2 的调查请求与 Effective Process Artifact、S3 的人工候选权威提交、S4 的 World 调研/选址审批/项目/WBS/合同/施工/会计、S5 的两路径和现场验收仍未完成。既有未跟踪 M7 验收产物不属于本计划，不修改、不提交。

## 5. 当前接口与事实所有权

| 用户动作 | 浏览器调用 | AESE 职责 | IAOS 权威结果 |
| --- | --- | --- | --- |
| 打开交互规划 | `GET /api/aese/v1/world/plant-build/financial-constraints?case_code=...` | 使用当前 IAOS token/tenant 定向转发，不缓存业务事实 | 从已过账银行科目和已批预算形成只读 snapshot、来源引用与 hash |
| 保存需求并生成候选 | `POST /api/aese/v1/world/plant-build/proposals` | 严格校验 Requirement、调用已配置模型、保存 CreativeJob 技术证据 | `facility.requirement.define` 与 `site.proposal.record` 分别提交 Requirement 和 candidate-only ProposalSet |
| 审阅候选 | `POST /api/aese/v1/world/plant-build/reviews` | 从 IAOS profile 解析实际 actor，禁止浏览器伪造 reviewer | `site.proposal.review` 保存 action、理由、revision、reviewer、Audit 和 Outbox |

上述 BFF 不是通用 IAOS 代理，也不拥有权威业务表。外部模型未配置、财务快照不完整或变更、候选 Schema 不合法、权限不足、重复键输入不同、审阅版本冲突时均失败关闭。
