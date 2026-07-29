---
id: DES-028
title: Enterprise Genesis AI 原生企业创生游戏体验
date: 2026-07-27
status: completed
author: Codex + User
tags: [m9, game-experience, 2.5d, ai, agent, human-in-the-loop]
---

# Enterprise Genesis AI 原生企业创生游戏体验

> 2026-07-28：主页、租户前置、平台玩家身份与真实外部 AI provider 由
> [DES-029](DES-029-genesis-zero-start-and-ai-identity.md) 修订。本文已经交付的
> GameProjection、2.5D 世界和 M9 操作闭环继续有效。

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

游戏端不得把“在 IAOS 处理”的外链当作可玩闭环。新案件不存在时直接进入创业构想和
企业身份工作室；案件创建后，Agent 派遣、G1–G7 提交与决定、金额输入、系统 Capability
和三个 World wait 均在当前游戏界面内完成。每次成功操作后重新读取 verified evidence
和持久工作项，前端不乐观推进 lifecycle。

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
- **GX6 Founder-operated loop**：从空白 case 创建企业，在游戏内逐项完成 18 个 IAOS
  Work Item，并以三视口浏览器测试证明终态来自 IAOS。

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

## 17. D26 — 对话驱动的创始人玩法

仅把 IAOS Work Item 显示成任务卡不构成游戏。玩家主路径采用三层交互：

1. **场景层**：PixiJS 渲染创始办公室、玩家角色、数字员工、家具和因果变化；
2. **叙事层**：数字员工以 RPG 对话逐问产业、客户、产品、品牌性格和风险偏好；
3. **治理层**：对话选择先形成 command draft，明确确认后才转换为 IAOS Capability
   输入；技术 capability、hash 和 evidence 默认收纳在证据抽屉。

首章必须先创建 FounderProfile 与虚拟头像，再由创业顾问“纪元”引导完成企业身份。
玩家看到的是主线任务、NPC 对话、选项后果和世界反馈，不是字段名与流程节点。AI 生成
名称属于故事中的数字员工协作动作；真实 IAOS commit 仍只发生在玩家签署创始人指令后。

后续 17 个工作项按同一模式改造为可玩事件：决议通过董事会对话、登记通过政务大厅、
开户通过银行谈判、CEO 任命通过候选人会面、预算通过经营会议。通用
`WorkItemActionPanel` 降级为治理/排障 fallback，不再是默认体验。

## 18. D27 — 审批对象先于审批动作

M9 的所有人类审批统一采用“起草结果 → 审阅文件 → 批准并执行”语义，禁止只显示
Capability 名称就允许批准。Agent 工作项完成后，游戏投影必须携带其 committed output
或由 IAOS 已提交事实编译出的审批事项：

- 节点 2 的三项输入形成持久化《创始人设立决议草案》，G1 展示原始决议目标、核心
  提案、风险限制、起草人及 Agent Run 证据；
- G2–G7 分别展示登记申请、开户申请、实缴核验、任命议案、经营授权书和启动预算；
- 每张审批事项必须包含摘要、关键字段、风险/限制、批准后的效果和 evidence ref；
- 缺少可审阅事项时，游戏界面失败关闭，不允许盲目批准。

审批 UI 是 IAOS committed evidence 的业务化表达，不创建第二份审批真相。

## 19. D28 — 登记与银行是可失败的外部经营事件

`registration.submit` 必须展示拟提交的登记申请及资料清单，不再只显示技术动作。
`registration.observation.commit` 必须允许外部机构给出合格或退回补正反馈；缺件时记录
`rejected` World Observation，案件保持等待状态，玩家补齐资料后以新的幂等结果重新申请。
通过后发放明确标注为沙盘虚构资产的营业执照、公章、财务章和法定代表人章。

银行开户采用同一结构：玩家先在多家虚拟银行中选择服务、速度和尽调偏好，审阅开户
资料包；缺少执照、身份、章程/决议、印鉴或地址证明时，银行可拒绝并说明原因。补件后
重新申请，通过后发放基本账户信息与虚构企业网银 U 盾资产。准确业务字段由 UI 文本层
和 IAOS committed state 提供，生成式图片只作表现素材，不作为业务证据。

## 20. D29 — 世界优先场景导航

当前原型把同一种几何建筑按钮叠在写实城市背景上，并由右侧 18 个 Work Item 与无实际
World Tick 的帧/倍率按钮主导操作。这违反本设计第 3、5、9 节，必须按以下合同重构：

1. 城市背景中的真实建筑体量就是可进入地点；交互热点必须与政务、金融、创业园和总部
   的视觉区域对齐，不再额外悬浮一组无关的通用蓝色房子。
2. 地图只显示地点名称、状态和当前事件标记。点击地点进入对应室内场景，NPC、资料台、
   资产柜和会议桌承载主交互。
3. 当前主线从建筑事件气泡进入；完整 Work Item 只保留在“治理档案”辅助视图，不作为
   游戏主导航。
4. M9 暂无真实连续 World Tick，删除暂停/1×/2×/4×；删除把 committed history 当成
   可前后切换 frame 的“上一步/查看后续”。历史改为只读企业大事记。
5. 章节条改为只读旅程状态，不允许通过点击章节切换业务状态。
6. 城市、室内场景与地点资产变化必须由当前 committed lifecycle、work item 和 evidence
   编译；前端选择地点只改变视角，不推进事实。
7. Projection 必须保证当前首个可执行 Work Item 所属地点已经开放；生命周期刚完成上一
   章节、下一节点已 ready 时，应按下一节点解锁目的地，禁止显示可点击任务却把地点保持
   locked。前端也不得以无反馈的方式吞掉受锁地点导航。

首个纵向切片覆盖创始办公室、政务服务中心、合作银行和企业总部：四个热点拥有不同
图形语言与室内布置；登记成功后办公室出现执照/印章资产，开户成功后出现账户/U 盾，
组织建立后总部出现管理层工位。移动端以地点卡替代精确地图点击。

## 21. D30 — 场景生命感与因果转场

地图热点可进入只是空间导航的第一步。地点切换和 committed 状态变化必须拥有游戏化
反馈，但不得用动画伪造业务推进：

- 从城市进入地点前显示 300–500ms 可跳过旅行转场，包含目的地、当前事件和路线；
- 城市绘制从当前基地到任务地点的因果路线，只有当前事件地点显示脉冲标记；
- 室内至少有一名与当前事件相符的 NPC，并用角色名、岗位、状态和对白表达 Work Item；
- Agent `ready/running/waiting`、World wait 和人类审批使用不同的姿态、气泡和动效；
- 执照、印章、账户/U 盾和管理层席位是 committed trophy，刷新和重新进入地点仍存在；
- 所有动效只改变表现状态，Capability 成功后重新读取 Projection 才能改变世界资产；
- `prefers-reduced-motion` 下取消位移与循环动画，旅行转场缩短为即时淡入；
- 低性能和移动端使用 CSS/DOM 角色与路线，不要求 WebGL 才能完成操作。

## 22. D31 — 玩家化身与可检查场景物件

世界中必须持续存在玩家化身，而不是只有顶部头像。城市层的化身位于当前地点，前往任务
地点时沿因果路线移动；室内层以场景热点表示玩家可以检查的桌面、窗口、资产柜和会议桌。

- 城市角色只表达视角位置，不自行触发 Capability；
- 旅行开始后角色沿当前路线移动，到达后才切换室内场景；
- 无路线或同地点进入时使用短距离淡入，不绘制虚假路径；
- 室内物件点击打开轻量详情卡，显示用途、当前状态、解锁条件和 committed 来源；
- 物件详情不得成为第二套表单，正式动作仍由 NPC 当前事件入口发起；
- 键盘焦点、触摸目标和文本列表必须提供地图热点的等价操作；
- 当前 SVG/CSS 角色和物件属于版本化 fallback，后续 raster atlas 替换不得改变交互合同。

## 23. D32 — 统一 2.5D 室内场景美术

创始办公室、政务服务中心、合作银行和企业总部使用 OpenAI 图片生成能力制作同一视觉
体系的 1672×941 2.5D 场景。场景图片只负责空间、氛围和家具表现，任务、NPC、物件热点、
解锁状态和 IAOS evidence 仍由可访问的 DOM 交互层承载。

- 图片必须版本化存放于 `frontend/public/assets/enterprise-genesis/locations/`；
- `manifest.json` 记录尺寸、SHA-256、来源、许可和保留策略；
- 禁止把生成图片中的文字、证照或账户信息当作业务事实；
- 背景加载失败时保留颜色底图、NPC、任务和物件操作，核心流程仍可完成；
- 宽屏采用 cover 居中裁切，移动端不得依赖图片中的精确像素位置完成正式操作。
- 透明人物精灵按房间家具校准：完整身高应约为单人椅可见高度的 1.7–2 倍，玩家与
  同一景深 NPC 的身高比不得超过 1.3；桌面与移动端分别设置响应式尺寸和移动距离。
- 人物脚底以地面为缩放原点并保留接触阴影；人物放大后仍不得遮挡正式操作、越出场景
  或使 NPC 对话卡溢出。三类视口的端到端测试必须守住最小可见身高与角色比例。

## 24. D33 — 全 M9 专属 RPG 事件合同

17 个设立工作项统一采用“章节地点 → 具名人物开场 → 三项经营思考 → 文件/资料/金额
操作 → IAOS 正式行动 → 奖励反馈”。经营选项帮助玩家理解责任与取舍，但不能覆盖金额、
审批对象、外部观察或 Capability 合同。`WorkItemActionPanel` 继续承担受治理提交，
`RPGEventIntro` 承担每个 capability 的专属剧情，不再把技术字段作为第一信息层。

## 25. D34 — 精灵、移动、声音与永久奖励

- 创始人、纪元、林岚、周衡和顾远使用独立透明角色精灵；
- 营业执照、印章套装和 U 盾使用独立 collectible 精灵并明确标注虚构；
- 点击室内物件后创始人先移动到热点，再打开状态卡；
- 旅行、检查和成功采用短促 Web Audio 反馈，默认开启且可随时关闭；
- 每个 committed 工作项进入只读企业大事记，完成时显示短暂里程碑奖励；
- 所有移动遵守 `prefers-reduced-motion`，图片失败时 CSS fallback 仍可操作。

## 26. D35 — 多租户与失败恢复验收

Genesis workspace 必须以玩家 owner、幂等键和 tenant assignment 隔离。相同幂等键在
provisioning 失败后重试必须复用原 workspace、tenant、world run 和 case code，禁止生成
孤儿租户。IAOS Founder 登录以 platform principal 的显式 tenant access assignment 展开
企业选择；每个选择项签发目标租户 JWT，跨租户读取 M9 案件必须返回 404/403。

## 27. D36 — 财务组织与开业会计章节

DES-030 扩展 M9 后，企业总部增加财务中心。玩家在组织骨架建立后配置 CFO/Controller/
会计/出纳/成本/审计岗位，处理职责冲突，选择账套和科目模板，并审阅银行注资产生的期初
凭证。游戏必须展示：

- 银行到账、资本事项、记账规则、凭证、总账和开业报表的穿透链；
- 认缴、实缴、现金、费用、应付、预算和利润之间的明确区别；
- 自动凭证由已发布规则生成，正常经营不要求玩家逐张手工制证；
- 重大手工调整、付款、政策、关账和财务异常由有权玩家/财务角色处理；
- `finance_opening_ready` 达成后才允许显示新版本企业开业完成。

## 28. D37 — 世界摘要与 IAOS 财务穿透

AESE 不承担长期账务浏览和报表工作台。企业总部的“开业财务中心”使用账套/凭证图形
表达财务系统入口，必须是可聚焦、可点击的交互对象，并与治理会议桌保持非重叠布局。

- 点击财务中心打开说明卡，解释目的、账套状态和 committed evidence 来源；
- “查看系统账务”进入 IAOS `#finance_workspace` 的 ledger 视图；
- “查看财务报表”进入同一页面的 reports 视图；
- 治理档案只保留当前章节的小型开业摘要和两个穿透按钮；
- 未来业务增长不得把完整银行流水、总账和多期报表继续堆入游戏侧栏；
- 链接携带 tenant/case 定位，IAOS 必须用当前 JWT 重新执行 tenant/RBAC 校验。

三视口测试必须验证财务中心与会议桌不相交、按钮可通过键盘访问、两个链接携带当前
tenant/case 且落到 IAOS 正式财务页。
