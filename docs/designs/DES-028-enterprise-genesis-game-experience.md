---
id: DES-028
title: Enterprise Genesis AI 原生企业创生游戏体验
date: 2026-07-27
status: approved
author: Codex + User
tags: [m9, game-experience, 2.5d, ai, agent, human-in-the-loop]
---

# Enterprise Genesis AI 原生企业创生游戏体验

## 1. 产品愿景

把 M9 从“可以逐节点跑通的企业成立流程”升级为“玩家、人类管理者和数字员工共同创建一家企业”的经营模拟开局。

体验借鉴 Capitalism Lab 的公司创建、企业地图、公司详情、管理层、政策和 AI 经理等机制，但不复制其界面。AESE 的差异化是：

- 游戏画面展示企业世界、空间、时间、参与者和结果；
- IAOS 保存正式主体、流程、审批、业务事实、权限和审计；
- AI 负责创意生成、资料准备、分析、建议和受授权执行；
- 人类负责愿景、偏好、关键选择、治理门批准和必要接管；
- World 负责登记机构、银行、候选人和资金结算等外部结果。

> 玩家不是在填一套成立表单，而是在带领一组数字员工，把创业想法逐步变成一家真实可运行、可追溯的虚拟企业。

## 2. 核心原则

1. **游戏化不改变事实所有权**：2D/2.5D 画面、动画和任务卡都是 IAOS/World 持久事实的投影，不能直接推进正式状态。
2. **AI 生成不等于业务生效**：名称、Logo、企业色、章程和材料先进入 `candidate`；人工选择并通过 Capability 后才成为正式资产或事实。
3. **数字员工不是聊天装饰**：Agent 是 IAOS service-only acting principal，具有岗位、Mandate、工具白名单、额度、期限和升级对象。
4. **乐趣来自选择、约束和后果**：有限资本、时间、材料质量、授权范围、补正、拒绝和人工接管共同形成玩法。
5. **快节奏不伪造确定性**：支持暂停、1×、2×、4× 虚拟时间；关键审批自动暂停，外部结果仍由版本化 World policy 产生。

## 3. 核心体验循环

```text
表达创业意图
→ AI 生成候选方案
→ 玩家选择、修改或授权
→ 数字员工准备并执行任务
→ IAOS 校验、审批和提交
→ World 推进外部过程
→ 2.5D 世界产生可见变化
→ 成本、时间、风险和证据反馈
→ 解锁下一项决策
```

每轮必须回答：世界发生了什么、哪个员工在做什么、为什么需要玩家、选择影响什么、正式证据在哪里。

## 4. M9 游戏章节

### Chapter 0 — 创业构想

玩家用自然语言描述企业方向、客户、产品、品牌气质、可投入资金和风险偏好。AI 输出结构化 `FounderIntent`，明确区分用户事实、AI 推断、信息缺口和待确认项。

### Chapter 1 — 企业身份工作室

AI 提供 3–6 个命名方向，每个包含中英文名称、简称、命名逻辑、行业联想、风险提示、口号、品牌关键词、企业色和 Logo brief。

玩家 shortlist 后生成 Logo 候选与应用预览。正式选择通过 `brand.identity.select` Capability 固化；Logo 进入版本化 Asset Registry，不在场景包内联 base64。

### Chapter 2 — 创始办公室

2.5D 场景从空办公室开始。`incorporation-agent` 准备创始人决议，`legal-compliance-agent` 检查，玩家完成 G1。桌面文件、白板、人物与任务气泡映射真实 Work Item、Agent Run、草稿和 Approval。

### Chapter 3 — 政务服务中心

登记提交后，地图出现政务中心和办件队列。玩家查看材料完整度、授权 G2，并应对受理、补正、拒绝或登记成功。普通补正可由 Agent 在既有授权内完成，重大变化必须重新审批。

### Chapter 4 — 银行与资本

地图解锁银行、投资人账户和公司账户。玩家比较虚构银行方案，由 finance/legal Agent 组包，经 G3 开户；随后处理出资、到账差异和 G4。

预算、承诺、实缴和现金必须使用不同图形与文字标签，不能只靠颜色区分。

### Chapter 5 — 人才与治理

注册后总部解锁。玩家查看 CEO/CFO 候选人的能力、薪酬、风险、接受概率和 Mandate 建议。AI 可总结匹配度、生成面试问题、模拟问答和形成任命草稿；候选人接受只能来自 World Observation，最终任命和授权经过 G5、G6。

### Chapter 6 — 开业准备

finance-agent 编制初始预算，audit-agent 独立检查，玩家完成 G7。只有 IAOS readiness evaluator 证明法律主体、账户、实缴、岗位、任命、Mandate 和预算完整一致，才播放企业开业并输出 `enterprise_operational_ready`。

### Chapter 7 — 进入经营世界

M9 是游戏开局：M10 解锁选址建设，M11 解锁产能，M12 解锁工业化，M13 解锁订单交付，M14–M24 形成长期经营、竞争、危机与策略演化。品牌、资本、组织、Agent 和治理策略必须延续。

## 5. 三层界面

### 5.1 世界沙盘

可缩放等距 2.5D 城市/园区地图包含创始办公室、政务中心、银行、人才会面点、总部和事件/资金路径。桌面支持平移缩放；移动端提供场景卡与明确导航，不要求精确拖拽。所有地图对象都有列表替代视图。

### 5.2 经营桌面

上下文面板包含当前目标、我的待办、数字员工、AI 创意工作室、决策对比、资金/预算/时间/风险和 World 新闻。点击建筑、人物或事件切换上下文，不堆叠多层弹窗。

### 5.3 治理证据

任何任务均可打开“为什么/证据”抽屉，查看 Capability、Process、actor、Mandate、Agent Run、Tool Call、Approval、Decision、Intent、Observation、CommittedOutcome、Journal、Outbox、correlation 和 hash。默认使用业务语言，专家可展开合同和 JSON。

## 6. AI 能力分层

- **创意 AI**：名称、口号、Logo brief、Logo、企业色、品牌卡和文案；只产生候选资产。
- **助理 AI**：意图结构化、规则解释、方案比较、会议摘要和缺口提醒；推断必须等待确认。
- **数字员工 Agent**：复用 incorporation、governance、legal-compliance、finance、audit 五 Agent，通过正式上下文、schema、policy 和 Capability 工作。
- **世界 AI**：首版不让自由 LLM 决定登记或银行结果；AI 只生成解释和对话，World policy 决定事实。

LLM 永远不是业务事实源。模型输出先是 proposal，经校验、审批或 World rule 后才能 committed。

## 7. AI 品牌资产管线

```text
FounderIntent
→ NamingProposal[]
→ 玩家 shortlist
→ BrandBrief
→ LogoGenerationJob
→ BrandAssetCandidate[]
→ 玩家选择/编辑
→ 合规检查
→ brand.identity.select
→ Versioned Brand Asset
→ IAOS/AESE 投影
```

资产记录 tenant、case、company stable code、类型、尺寸、prompt、模型/provider/version、seed、参数、source refs、生成者、时间、hash、许可、状态和选择理由。正式 Logo 派生横版、方形、单色、深色背景、浅色背景版本；生成失败时使用文本字标/几何 fallback，不阻塞公司成立。

## 8. 人工与 Agent 协作

工作项显示 owner、动作、虚拟时间、依赖、输入合同、工具、自治模式、产出、风险、证据及暂停/拒绝/接管入口。

Agent 自治等级：

1. 建议：只生成提案；
2. 草稿：准备正式材料；
3. 授权执行：在 Mandate 内调用低风险 Capability；
4. 自动补正：在既有批准范围和次数内重试。

G1–G7 不因自治等级消失。人工接管保留原 Agent 分派和原因，玩家不冒充 Agent。

## 9. 视觉与动效

采用“工业经营沙盘 + 轻拟物 2.5D + 数据驾驶舱”，避免卡通玩具感、赛博霓虹和满屏 AI 紫色渐变。

- 世界层：温暖中性底色、低饱和工业蓝、企业品牌色点缀；
- IAOS 正式状态：深蓝/石墨；
- World：琥珀；Agent：青色；人工等待：紫色；
- 风险/成功必须同时使用图标、文字和颜色。

动效 150–300ms 用于控件，300–500ms 用于场景变化。资金、文件和人物移动只表达真实因果。支持 `prefers-reduced-motion`，动画可中断且不阻塞操作。

过程使用阶段路线与时间线；资金使用流向图加表格；预算/现金/承诺使用数值或 bullet chart；Agent 使用任务队列。所有图表提供文本/表格替代。

## 10. 技术架构

### 10.1 前端渲染

继续复用 AESE React/Vite：

- React DOM：HUD、表单、审批、任务和证据；
- PixiJS/WebGL：首版等距 2.5D 世界、精灵、路径和粒子；
- SVG：流程线、关系和可打印证据；
- 首版不引入完整 Three.js/自由 3D 相机。

渲染层只消费 `GameProjection`，不读取 IAOS 数据库，也不实现业务状态机。

### 10.2 Game Projection

```text
GameProjection {
  world_scene, lifecycle, buildings[], actors[], work_items[],
  resources, exchanges[], brand, notifications[], evidence_refs[], cursor
}
```

AESE 聚合 World snapshot 与 IAOS committed trace。SSE/WebSocket 只提示变化，断线后按 cursor/snapshot 补齐。

### 10.3 命令路径

```text
玩家/Agent
→ AESE Game Command API
→ IAOS Agent Tool / Capability / Approval
→ Journal + Outbox
→ World Bridge（需要时）
→ committed trace
→ GameProjection
→ 2.5D 变化
```

纯创意生成走独立 `CreativeJob`，不能借 Game Command API 直接写 Entity。

### 10.4 素材与性能

使用版本化 atlas manifest；Logo/肖像/建筑图采用 WebP/AVIF并保留原文件；按 Chapter 懒加载；设定 draw call、纹理、actor 和帧预算；低性能设备降级为 2D 场景卡。交互反馈目标小于 100ms，60fps 为目标、30fps 为可接受降级。

## 11. 可玩失败路径

首版支持名称/Logo 重生成、登记补正、开户拒绝、候选人拒绝、出资差异、Agent 越权、Agent 暂停后人工接管、资金不足导致缩减/暂缓，以及玩家拒绝 Agent 建议并保留理由。

失败通常产生新问题和选择，而非立即 Game Over；只有明确终止、资金不可恢复或玩家放弃才结束 run。

## 12. MVP 切片

- **GX0 投影与原型**：冻结 GameProjection、Chapter、地图、状态映射、命令边界、三视口线框和性能预算。
- **GX1 AI 企业身份工作室**：FounderIntent、命名、口号、Logo、版本、fallback 和 `brand.identity.select`。
- **GX2 2.5D 创始办公室**：空办公室、五 Agent 工位、Chapter 0–2、真实 Agent Run、G1/G2 和证据抽屉。
- **GX3 城市与外部过程**：政务、银行、人才、总部、虚拟时间和 Bridge 交换可视化。
- **GX4 完整 M9 主线**：G1–G7、五 Agent、人工接管和 `enterprise_operational_ready`。
- **GX5 质量交付**：恢复、三视口、键盘、触摸、减弱动效、性能、安全、runbook 和 evidence。

## 13. 验收

1. 玩家只通过游戏界面即可从创业想法创建并正式选择名称与 Logo。
2. AI 资产与业务事实严格分离，生成和选择均可追溯。
3. 玩家、五 Agent、IAOS Runtime 和 World 真实交替推进 M9。
4. G1–G7 均等待正确的人类决定。
5. 2.5D 变化只由持久化 committed projection 驱动。
6. 补正、拒绝、差异、候选人拒绝和 Agent 越权可玩且可恢复。
7. 关闭页面、服务重启和重复操作不丢任务、不重复业务变化。
8. 1440×900、1280×720、390×844 可完成主线，移动端无需精确地图操作。
9. 键盘、屏幕阅读器、减弱动效和数据替代视图可用。
10. M9 完成后品牌、资本、组织、Agent 和治理状态可进入 M10。

## 14. 非目标

- 不复制 Capitalism Lab 的 UI、美术或数值系统。
- 不在首版实现开放世界、多人在线、完整竞争经济或全 3D 城市。
- 不让 LLM 自由修改资金、登记、审批、任命或预算事实。
- 不把真实工商/银行结果或法律建议伪装进虚构模拟。
- 不用 AI 素材替代无障碍文本、确定性 fallback 或版本治理。

## 15. 参考

- [Capitalism Lab 官方网站](https://www.capitalismlab.com/)
- [创建新公司](https://www.capitalismlab.com/playing-without-company/)
- [子公司控制与 AI 经理](https://www.capitalismlab.com/subsidiary-dlc/subsidiary-control/)
- [Capitalism Lab 10.0：Logo、经理肖像与 AI](https://www.capitalismlab.com/version100/)

## 16. 待确认

- 产品名采用 `Enterprise Genesis`，还是 AESE 内的“企业创生”模式。
- 美术采用等距插画，还是等距像素/低多边形预渲染。
- AI Logo provider、额度、许可和内容保留策略。
- 名称检查只做虚构世界规则，还是增加真实公开数据检索提示。
- 首版允许自定义行业，还是固定 HCTM 汽车热管理作为引导章节。

DES-028 已批准，由 PLAN-GX-001 实施；当前 M9N active 主计划的完成口径不变。
