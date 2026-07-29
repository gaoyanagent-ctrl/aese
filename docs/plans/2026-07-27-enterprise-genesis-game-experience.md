---
id: PLAN-GX-001
title: Enterprise Genesis 游戏体验实施计划
date: 2026-07-27
status: completed
plan_type: parallel-subplan
parent_plan: PLAN-M9-NATIVE-001
owner: AESE Game Experience
author: Codex + User
tags: [m9, game, 2.5d, ai, agent]
---

# Enterprise Genesis 游戏体验实施计划

## 1. 定位与依赖

本计划实施 DES-028，是 `PLAN-M9-NATIVE-001` 的并行子计划，不替代当前唯一 active 主计划。

- M9N 继续拥有 IAOS 工作项、审批、Agent Runtime、World wait 和 T83/T96 全链验收。
- GX 只消费 M9N 已提交的 trace/工作项/交换，不修改或复制 IAOS 状态机。
- GX0、GX1 的投影和创意资产可以立即实施。
- GX2–GX4 的正式推进依赖 M9N T83/T96 合同稳定。
- IAOS 新 Capability 或资产注册表修改仍须独立 IAOS worktree。

## 2. GX0 — GameProjection 与交互骨架

- [x] GXT1 定义 DES-028 与首版产品/技术边界。
- [x] GXT2 建立 `GameProjection` Go 合同、JSON Schema 和示例 fixture。
- [x] GXT3 严格拒绝未知字段、非法时间倍率和无 evidence 的工作项。
- [x] GXT4 从现有 IncorporationTrace 编译 GameProjection，禁止前端推断业务状态。
- [x] GXT5 增加 `/api/aese/v1/game/incorporation/:case/projection` 快照 API。
- [x] GXT6 定义 cursor 增量提示与断线补齐合同。
- [x] GXT7 增加 TypeScript 类型、data source 和 reducer。
- [x] GXT8 完成世界沙盘/经营桌面/治理证据三视口线框原型。
- [x] GXT9 建立性能预算、2D fallback、键盘和 reduced-motion 验收。

验收：同一 M9 case 可确定性生成同一 projection；画面状态均有 evidence ref。

## 3. GX1 — AI 企业身份工作室

- [x] GXT10 定义 FounderIntent、NamingProposal、BrandBrief、CreativeJob 和 BrandAsset 合同。
- [x] GXT11 实现 provider-neutral 创意模型接口和本地 deterministic fake。
- [x] GXT12 实现自然语言创业意图结构化与人工确认。
- [x] GXT13 实现 3–6 组公司名称、口号、关键词和命名解释。
- [x] GXT14 接入图像生成 provider，生成首批 2.5D 候选素材；Logo 使用可替换字标/几何 fallback。
- [x] GXT15 建立版本化 Asset Registry、对象存储引用、hash、许可和保留策略。
- [x] GXT16 实现内容安全、名称风险提示、额度、重试和 fallback。
- [x] GXT17 初始品牌选择复用 IAOS `incorporation.case.open`，候选确认后才成为案件事实。
- [x] GXT18 完成品牌工作室 UI、比较、编辑、选择、撤销和证据。

验收：AI 生成与正式品牌事实严格分离；失败不阻塞成立主链。

## 4. GX2 — 2.5D 创始办公室

- [x] GXT19 引入 PixiJS 渲染壳和可替换 scene adapter。
- [x] GXT20 建立等距坐标、点击、键盘导航和列表替代。
- [x] GXT21 实现创始办公室、董事长决策桌和五 Agent 工位。
- [x] GXT22 将 actor/work item/approval 状态映射为人物与场景表现。
- [x] GXT23 实现 Chapter 0–2 与 G1/G2 committed 场景变化。
- [x] GXT24 实现上下文经营桌面和治理证据抽屉。
- [x] GXT25 生成/制作首批建筑、人物、文件和 Logo fallback 展示素材。
- [x] GXT26 验证画面不提前推进、刷新恢复和低性能 2D fallback。

验收：真实 Agent Run 和 G1/G2 能在办公室被观察、操作和追溯。

## 5. GX3 — 城市地图与外部世界

- [x] GXT27 实现政务中心、银行、人才会面点和总部地图。
- [x] GXT28 实现暂停、1×、2×、4× 与关键工作项自动暂停。
- [x] GXT29 在证据视图呈现 Intent、Observation、CommittedOutcome 和 Discrepancy。
- [x] GXT30 复用 IAOS 补正/拒绝/差异恢复路径并投影通知。
- [x] GXT31 实现可访问资金账户与资本/预算表。
- [x] GXT32 实现受证据约束的 World 通知；事实仍由 policy 决定。
- [x] GXT33 验证重复、乱序、重启、断线和 cursor 恢复。

## 6. GX4 — 完整 M9 游戏主线

- [x] GXT34 接通 G1–G7、五 Agent、IAOS Runtime 与三个 World wait。
- [x] GXT35 通过 IAOS 工作项深链提供 Agent 暂停、人工接管和拒绝建议。
- [x] GXT36 实现资本、组织、Mandate、预算和 readiness 场景反馈。
- [x] GXT37 完成 `enterprise_operational_ready` 开业场景。
- [x] GXT38 将品牌、资本、组织和 Agent 状态移交 M10。
- [x] GXT39 完成正常线和可恢复异常线的 IAOS 集成验收。

## 7. GX5 — 质量与交付

- [x] GXT40 运行 Go、Schema、TypeScript、组件和 Playwright 测试。
- [x] GXT41 完成 1440×900、1280×720、390×844 三视口验收。
- [x] GXT42 完成键盘、屏幕阅读器、触摸和 reduced-motion。
- [x] GXT43 游戏入口懒加载，2.5D 独立 chunk，移动/低数据回退 2D。
- [x] GXT44 验证租户、权限、幂等、失败无部分写入和证据完整性。
- [x] GXT45 编写 runbook、evidence、素材许可清单和演示脚本。
- [x] GXT46 更新 Roadmap、Code Map、Progress Log、Atlas 和双仓 revision。

## 8. 完成定义

- GXT1–GXT46 全部关闭。
- DES-028 十项验收标准均有机器或浏览器证据。
- 游戏画面没有独立业务真相，AI 没有未治理写路径。
- M9N T83/T96 及 GX4 全链同时通过。

## 9. GX6 — 玩家原生操作闭环

- [x] GXT47 新案件不存在时进入创建态，不再因 projection 404/502 阻断。
- [x] GXT48 增加“新建企业”，生成独立 case code 并从创业构想开始。
- [x] GXT49 企业身份选择通过 `incorporation.case.open` 创建 IAOS 正式案件。
- [x] GXT50 Agent task 在游戏内填写业务输入并调用 IAOS dispatch-agent。
- [x] GXT51 G1–G7 在游戏内创建 Approval Request、由当前受派审批人决定并执行节点。
- [x] GXT52 资本、实缴与预算使用人民币元输入，提交时转换为精确 minor unit。
- [x] GXT53 三个 World wait 在游戏内产生受治理 Observation 后执行 commit 节点。
- [x] GXT54 系统 Capability 在游戏内执行，完成后重新读取 IAOS verified projection。
- [x] GXT55 超过 Agent 金额授权上限时展示 IAOS 治理拒绝，不生成前端假成功。
- [x] GXT56 Playwright 在 1440、1280、390 三视口从空白 case 完成全部 18 个工作项。
