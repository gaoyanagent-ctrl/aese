---
id: PLAN-GXZ-001
title: Enterprise Genesis 零起点与真实 AI 身份实施计划
date: 2026-07-28
status: active
author: Codex + User
tags: [genesis, tenant, homepage, minimax, onboarding]
---

# Enterprise Genesis 零起点与真实 AI 身份实施计划

## 1. 交付目标

从根主页创建独立 IAOS tenant、AESE World Run 和 M9 案件，并用真实 MiniMax 调用生成
企业身份候选。`tenant-hctm-genesis` 只保留为 demo/测试 fixture。

## 2. Z0 — 合同与风险

- [ ] Z1 冻结 PlayerAccount、GenesisWorkspace、ProvisioningCheckpoint 和证据合同。
- [ ] Z2 冻结 opaque tenant ID、workspace ID、case code 与 correlation 生成规则。
- [ ] Z3 建立 IAOS/AESE 跨仓矩阵、权限矩阵、补偿矩阵和 failure injection 表。
- [x] Z4 验证 MiniMax endpoint/model/key health，禁止在日志输出密钥。
- [ ] Z5 为 DES-029 建立 API、Schema 和 contract tests。

## 3. Z1 — IAOS Genesis Provisioning Saga

- [ ] Z6 在独立 IAOS worktree 实现 `genesis.workspace.create/read` 资源。
- [ ] Z7 新增控制平面 GenesisWorkspace/step 表和 RLS/owner 读取边界。
- [ ] Z8 服务端编排 tenant create、plan、founder membership、M9 Runtime install。
- [ ] Z9 完成 RLS/login/registry/process smoke check 后才激活 tenant。
- [ ] Z10 实现 retry、幂等冲突、checkpoint 恢复和失败 tenant 隔离。
- [ ] Z11 移除对全局固定 `founder-principal` 和 HCTM organization code 的新租户依赖。
- [ ] Z12 两玩家/两租户并发、越权、重试和部分失败集成测试。

## 4. Z2 — AESE World 与 BFF

- [ ] Z13 接收 tenant activated outcome 并创建独立 World Run。
- [ ] Z14 增加生产 workspace 列表/详情 BFF，不保存平台管理员凭据。（loopback local adapter 已完成）
- [ ] Z15 绑定 workspace、tenant、world run、creative job、case correlation。
- [ ] Z16 验证 IAOS/AESE 重启、消息重复与 World 创建失败恢复。

## 5. Z3 — MiniMax Provider

- [x] Z17 实现 `MiniMaxProvider`、配置校验、HTTP 超时和响应上限。
- [x] Z18 版本化 prompt，要求严格 JSON naming proposal。
- [ ] Z19 分离 reasoning、Schema 校验、一次 repair 和业务校验。
- [ ] Z20 持久 CreativeJob 调用证据、token/latency/fallback reason。
- [ ] Z21 实现限流、幂等、取消、重试和显式 deterministic fallback。
- [ ] Z22 Provider 状态 API 区分 connected/degraded/fallback/not_configured。
- [ ] Z23 真实 key smoke 与敏感日志测试。

## 6. Z4 — 产品主页与 onboarding

- [x] Z24 根路径新增 Enterprise Genesis Home。
- [x] Z25 实现创建新企业、继续经营、样板世界和功能说明。
- [x] Z26 实现玩家登录/本机 owner 身份入口；明确标记为非生产认证。
- [ ] Z27 实现五步创建向导和草稿自动保存。（首步独立空间 onboarding 已完成）
- [ ] Z28 实现八 checkpoint provisioning 进度、失败恢复和证据 ID。
- [x] Z29 workspace active 后获取 Founder tenant session 并进入 AI 身份工作室；旧管理员会话可 owner-scoped 刷新。
- [ ] Z30 移除主路径对 tenant/case/auth_token URL 参数的依赖。
- [ ] Z31 AI/fallback/Logo fallback 使用准确标签。

## 7. Z5 — 全链验收

- [ ] Z32 从根地址创建 tenant A 并完成 M9。
- [x] Z33 创建/使用 tenant A、B 完成 Founder 选择、JWT claim、profile 和 M9 cross-read 隔离验证。
- [ ] Z34 验证浏览器 token 无平台管理员权限。
- [x] Z35 完成 workspace provisioning 失败注入与同 identity 重试恢复；runtime/world/AI 的既有失败测试保留。
- [x] Z36 完成三视口、键盘/触摸、刷新恢复和重启验收；慢网由 loading/degraded 路径覆盖。
- [ ] Z37 更新 runbook/evidence/roadmap/code map/progress/Atlas 和两仓 revision。

## 8. 完成门

- 根地址是唯一对外推荐入口。
- 新公司不复用 fixture tenant 或已有业务数据。
- MiniMax 调用有 request/model/token/validation 证据；fallback 不冒充 AI。
- 两个独立 tenant 的 RLS、身份、Runtime、World 和 M9 全链都通过。

## 9. Z6 — 游戏化主交互

- [x] Z38 实现 FounderProfile、四个本机头像候选与角色持久化。
- [x] Z39 用 PixiJS 建立创始办公室、创始人和数字员工场景。
- [x] Z40 用 RPG 对话选择收集产业、客户、产品、品牌性格与创业宣言。
- [x] Z41 将 MiniMax 命名和 `incorporation.case.open` 嵌入对话任务并保持人工确认。
- [x] Z42 把剩余 17 个 M9 Work Item 逐章改造成董事会、政务、银行、人才和经营会议专属 RPG 事件。
- [x] Z43 将节点 2 输出解释为创始决议草案，并为 G1–G7 增加先审阅审批对象再批准的统一交互。
- [x] Z44 将登记和银行开户拆为资料包、机构审查、拒绝反馈、补正重申与获批资产领取事件。
- [x] Z45 世界优先场景导航：统一城市热点、四类室内地点、移除伪时间/帧控件并将 Work Item 降级为治理档案。
- [x] Z46 场景生命感：任务路线、旅行转场、室内 NPC、工作状态动效和 committed trophy。
- [x] Z47 玩家化身沿任务路线移动，室内物件可检查并显示 committed 状态和解锁条件。
- [x] Z48 使用 OpenAI 图片生成能力制作并接入四类统一 2.5D 室内场景，登记素材来源、许可和 hash。
- [x] Z49 生成并接入创始人、四名 NPC、营业执照、印章套装和 U 盾透明精灵。
- [x] Z50 增加室内热点移动、检查动作、旅行/检查/成功音效和减弱动效 fallback。
- [x] Z51 增加里程碑奖励和逐项 committed 企业大事记。
- [x] Z52 修复 Founder IAOS 多租户发现，并完成双租户与 provisioning 失败恢复验收。
- [x] Z53 校准室内人物与房间家具比例，并增加桌面、移动端最小身高及角色比例回归断言。
- [x] Z54 修复资本核验完成后组织节点与总部解锁不同步，并增加投影及浏览器回归测试。
