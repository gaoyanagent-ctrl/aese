---
id: DES-029
title: Enterprise Genesis 零起点主页、独立租户与真实 AI 企业身份
date: 2026-07-28
status: approved
author: Codex + User
tags: [genesis, homepage, tenant-provisioning, minimax, ai, onboarding]
---

# Enterprise Genesis 零起点主页、独立租户与真实 AI 企业身份

## 1. 修正的问题

当前实现有三个根本缺口：

1. `DeterministicProvider.GenerateNames` 返回四组代码模板；“AI 候选”标签不准确。
2. 用户必须持有 tenant、case 和 token 拼成的复杂 hash URL，缺少产品主页。
3. “新建企业”只在 `tenant-hctm-genesis` 中创建新案件，不会创建新的租户、身份、
   Runtime 和 World Run，因此不是从零开始。

DES-029 替代 DES-028 中关于默认入口、租户前置和外部模型已连接状态的假设；已经完成
的 2.5D 投影、18 工作项和 IAOS 治理链继续复用。

## 2. 产品入口

局域网和正式部署只公布根地址：

```text
http://192.168.50.222:4173/
```

根页面是 Enterprise Genesis Home，不要求用户理解 hash 参数。页面包含：

- 主行动：“创建新企业”；
- “AI 创意官”：嵌入“创始人指挥中心”的第一步，作为与主行动等价的可访问按钮
  进入创业项目定义；页面必须说明它负责把行业、区域和经营目标整理为创业项目草案，
  并明确它不代替注册、审批或经营决策；
- “继续经营”：列出当前玩家拥有的 workspace、企业名称、阶段和最后活动时间；
- “体验样板企业”：明确标记 demo，进入 HCTM fixture；
- 服务状态：IAOS、World、AI Provider，不能把 fallback 显示为在线 AI；
- 功能说明：解释 tenant、企业、World 和 AI 的边界。

首页的角色入口必须参与正常文档流，保持至少 44px 的点击区域、键盘焦点和明确的
accessible name；不得通过负边距或绝对定位覆盖主行动卡、企业列表或下一页面区块。
入口应位于“你的下一家公司”和四阶段流程之间，表达“先由用户发起构想，再进入 AI
身份工作室和 IAOS 设立”的顺序。桌面 1440/1280 与移动 390 视口均需验证入口可点击、
包含于指挥中心且不覆盖后续流程。

深链仍保留为恢复和支持入口，但由应用生成：

```text
/#genesis/workspaces/{workspace_id}
/#genesis/workspaces/{workspace_id}/incorporation
```

tenant ID、case code 和 token 不出现在主导航 URL；token 不放 query/hash。

## 3. 零起点 onboarding

### Step 0 — 玩家身份

玩家必须先在 IAOS 注册或登录平台级 PlayerAccount。AESE 不保存密码，不创建第二套
账号目录，也不接受“只输入游戏用户名”作为认证。IAOS 成功认证后签发短期
`genesis_player` Session；AESE BFF 只转发注册/登录请求和 Bearer Session，Workspace
控制面从 Session 中解析稳定 Player subject。

浏览器中的昵称、头像和界面偏好不参与授权。新租户 owner、chair 任职、Founder
Mandate 和 M9 人工工作项都引用 IAOS Player subject。`X-Genesis-Player-Id` 只能由
AESE BFF 从已验证 Session/IAOS profile 派生；正式模式不得信任浏览器声明的 player ID。

现有 IAOS 用户可以使用原凭据首次登录 Genesis：IAOS 验证任一 active tenant
credential 后，将其已绑定的 platform principal 提升为 PlayerAccount，确保既有
Workspace ownership 不丢失。全新用户通过 Genesis 注册入口建立独立平台 PlayerAccount，
此时不创建业务 tenant；只有用户点击“创建新企业”才进入 Workspace/tenant provisioning。

本机用户名映射只保留在显式 `AESE_AUTH_MODE=local_dev`。该模式必须同时绑定
`127.0.0.1`，不由产品登录页暴露，只供本机测试客户端显式携带开发 Header，且不得部署
到局域网或正式环境。默认 `AESE_AUTH_MODE=iaos` 在缺少有效 Session 时失败关闭。

### Step 1 — 创建创业空间

玩家只输入：

- 创业项目临时名称，例如“我的第一家制造企业”；
- 行业模板或空白模板；
- 业务区域与时区；
- 难度/真实性等级；
- 数据保留确认。

此时不要求公司正式名称、注册地址或经营范围。

### Step 2 — IAOS provisioning

IAOS `GenesisProvisioningSaga` 服务端执行：

1. 生成 `workspace_id`、opaque `tenant_id` 和 provisioning key；
2. 创建 `tenant_account(status=provisioning)` 和 tenant directory；
3. 建立当前玩家的 owner membership 与 tenant-scoped founder subject；
4. 安装 M9 Semantic、Entity、Capability、Policy、Approval、Process 和 Agent Runtime；
5. 创建五个 service-only Genesis Assistant，但不把它们冒充公司正式员工；
6. 等待 AESE 创建并持久化独立 World Run binding evidence；
7. 执行 RLS、owner binding、Runtime 和 World smoke，之后才激活 tenant；
8. 签发仅含 `genesis_owner` 的 tenant-scoped session。

进度页以八个可恢复 checkpoint 展示，不使用无限 spinner。失败显示阶段、原因、重试和
支持证据 ID。

### Step 3 — 创业构想

只有 workspace active 后才进入 FounderIntent。AI 输入不再要求前端提供 tenant；
服务端从 workspace membership 解析 tenant 和 case identity。

### Step 4 — 企业身份与正式案件

真实 AI 生成候选，玩家选择或编辑后，才通过 IAOS `incorporation.case.open` 创建
正式设立案。此时产生 proposed company name，后续登记成功才产生 legal entity。

## 4. 控制平面合同

### 4.0 当前生产纵切

IAOS DES-062 控制面与 AESE BFF 已实现：

```text
根主页
-> 创建隔离创业空间
-> 服务端生成 workspace/tenant/world/case 标识
-> IAOS tenant_account(provisioning)
-> Player owner/chair/Mandate bind
-> M9 Runtime install
-> AESE World evidence
-> tenant active
-> tenant-scoped genesis_owner session
-> MiniMax M3 企业身份工作室
```

Workspace/World 引用以 `0600` 本地状态文件保存；tenant、membership、Runtime 和
checkpoint 以 IAOS 为权威。浏览器只提交项目配置和幂等键，不能提交 tenant ID；
AESE 转发当前 Player session，不保存平台凭据或 Founder 密码。loopback dev adapter
仅在明确缺少生产 session 的本地开发路径保留。

loopback adapter 必须区分“首次 provisioning”和“既有 active tenant 的 session
refresh”。后者只能使用旧 Workspace 的确定性本地 Founder 凭据调用 login，成功后直接
返回 session；不得再次执行 Founder bootstrap、Runtime install、tenant activate 或任何
业务流程。active tenant 登录失败时 fail closed，并提示使用正式 IAOS 会话，禁止通过
bootstrap 重置身份来掩盖失败。

Player 控制面会话和 Workspace tenant 会话必须使用两个独立浏览器存储键。短期 Player
Token 只保存在当前标签页的 `sessionStorage`；tenant session 只在进入具体企业后用于
该企业业务调用。列表、创建和恢复 Workspace 只能使用 Player session，绝不拿 tenant
Token 或管理员 Token 轮询尝试。Player Token 被 IAOS 拒绝时，AESE BFF 透传 401，前端
清除失效 Player session 并返回登录入口；不得把身份失效包装成可重试的 500，也不得降级
为仅凭 localStorage player ID 访问生产控制面。

控制面上线前由 loopback adapter 创建的企业可能只有 AESE `0600` Workspace 记录和
完整的 IAOS tenant/Runtime/case，没有新版 IAOS Workspace 行。恢复 session 收到 404
时，AESE 只允许对当前本地 player 拥有的同 ID Workspace 发起一次
`legacy-adoptions`；不发送 tenant/owner 声明。IAOS 必须从当前 JWT 推导身份并验证
active tenant、owner、chair、Founder Mandate、M9 Runtime 和设立案，然后幂等登记
控制面。401/403/409/5xx 不得触发该恢复路径，也不得重新 bootstrap 或重放业务流程。

### 4.1 对外 API

```text
POST /api/v1/genesis/workspaces
POST /api/v1/genesis/workspaces/legacy-adoptions
GET  /api/v1/genesis/workspaces
GET  /api/v1/genesis/workspaces/{workspace_id}
POST /api/v1/genesis/workspaces/{workspace_id}/retry
POST /api/v1/genesis/workspaces/{workspace_id}/world-ready
POST /api/v1/genesis/workspaces/{workspace_id}/session
POST /api/v1/genesis/workspaces/{workspace_id}/members
```

调用者只需 `genesis.workspace.create/read`；平台内部服务身份才拥有 tenant create、
identity bootstrap、runtime install 和 activate 权限。

### 4.2 Provisioning record

记录至少包含：

- workspace ID、owner subject、tenant ID、world run ID；
- requested template/version、region、timezone；
- current checkpoint、status、attempt、error code；
- idempotency key、input hash、created/updated/completed time；
- tenant/runtime/world evidence refs。

创建重试必须返回同一 workspace；相同幂等键不同输入返回冲突。激活前任何失败不得留下
可登录但未完整安装的 tenant。

### 4.3 AESE 边界

AESE 只接收 `tenant_activated` committed outcome 后创建 World Run。主页聚合展示可以
由 AESE BFF 返回，但租户、身份和 Runtime 事实都来自 IAOS。AESE 不保存密码、平台
管理员 token 或租户业务表。

### 4.3.1 IAOS Player 认证合同

AESE 对浏览器公开同源 BFF：

```text
POST /api/aese/v1/auth/register
POST /api/aese/v1/auth/login
GET  /api/aese/v1/auth/session
```

BFF 分别代理 IAOS：

```text
POST /api/v1/genesis/auth/register
POST /api/v1/genesis/auth/login
GET  /api/v1/genesis/auth/session
```

注册输入为 username、password、display_name；用户名规范化后全局唯一，密码长度和复杂度
由 IAOS 校验。注册与登录错误不得泄露账号是否存在以外的敏感数据，不记录请求体或密码。
登录成功返回 Player subject、显示名、短期 JWT 和过期时间。后续 Workspace API 必须携带
该 JWT；IAOS 使用 JWT user ID 查 `genesis_player_account`，不能使用客户端 header
作为 owner 权威。

首个可交付版本允许浏览器会话存储短期 Token；正式互联网部署应由 BFF 兑换为
`HttpOnly + Secure + SameSite` Cookie，并接入 refresh/revoke、邮箱验证、找回密码、
MFA 和 OIDC。无论采用哪种载体，Token、Cookie、密码和 refresh secret 都不得进入
URL、World evidence、CreativeJob 或 AESE 持久文件。

登录页必须提供“登录 / 注册”两条清晰路径、密码显隐、字段约束、提交中状态和可理解的
错误。注册只产生平台 Player 身份；“创建企业”是登录后的独立动作。连续五次密码错误后
IAOS 临时锁定账号 15 分钟。Workspace BFF 在每次请求中用 Player Token 调 IAOS session
profile，并以 IAOS 返回的 subject 为 owner，浏览器声明的 subject 仅可作为调试信息。

### 4.4 跨仓、权限与失败恢复矩阵

| 事实/动作 | 权威仓库 | 浏览器权限 | 失败状态与恢复 |
|---|---|---|---|
| Workspace、member、checkpoint | IAOS | 当前 Player membership | 保持 provisioning/failed；同幂等键 retry；401 清理旧会话并要求重新登录；旧记录 404 经 IAOS 前置条件核验后一次接管 |
| Tenant create/activate | IAOS SaaS Ops | 无 platform 权限 | World/smoke 前不可 active |
| Owner user/role/position/Mandate | IAOS Identity/Governance | `genesis_owner` | 幂等 upsert；绑定失败不签 session |
| M9 Runtime | IAOS Runtime | 只读/执行已授权资产 | content hash no-op；失败停在 runtime checkpoint |
| World Run binding | AESE | 当前 Workspace | 稳定 workspace/world key 重试；无 evidence 不激活 |
| CreativeJob | AESE | tenant session 必须匹配请求 tenant | 同输入互斥/幂等；失败保留 reason 后重试 |
| Incorporation case | IAOS M9 Runtime | chair + Founder Mandate | Capability 事务回滚；不产生部分业务事实 |

失败注入至少覆盖 tenant 创建失败、owner 绑定失败、Runtime 安装失败、World evidence
缺失、smoke 失败、session 兑换失败、模型超时/429/5xx、非法模型 JSON 和重复请求。
任何失败都不得把平台 token、密码或模型密钥写入响应、日志或持久证据。

## 5. MiniMax 真实 AI Provider

### 5.1 Provider 选择

`creative.Provider` 增加：

- `MiniMaxProvider`：生产首选；
- `DeterministicProvider`：离线测试和显式 fallback；
- `ProviderRouter`：按配置、健康、额度和请求模式选择。

本项目 Coding Plan 账户的 `GET /v1/models` 于 2026-07-28 返回
`MiniMax-M3`，产品版本称 M3.0；代码必须使用账户返回的精确模型 ID，而不能把
“M3.0”直接作为请求值。该账户的 OpenAI-compatible endpoint
`https://api.minimax.chat/v1/chat/completions` 已完成真实 smoke。MiniMax 公开文档
当前仍主要列出 M2.7，因此启动和发布验收必须以账户 models API 与真实 completion
为准，并在主页显示实际 endpoint、model 和状态，不显示密钥。

### 5.2 调用链

```text
FounderIntent
-> PromptTemplate(versioned)
-> MiniMax chat completion
-> strip/separate reasoning
-> strict JSON decode
-> JSON Schema validation
-> business validation + duplicate/risk checks
-> NamingProposal[]
-> CreativeJob audit
```

模型不得直接写 IAOS。返回必须包含 4–6 个差异化候选：

- 中文全称、简称、英文名；
- 命名理由、品牌承诺、口号；
- 关键词和主色建议；
- 与行业/客户/产品输入的依据；
- 现实工商核名和商标风险提示。

禁止从 `<think>` 内容提取业务结果。无法得到合法 JSON 时最多一次修复调用；仍失败则
显示“AI 暂不可用”，由玩家明确选择重试或 deterministic fallback，不能把 fallback
继续标成“AI 生成”。

### 5.3 CreativeJob 与可观察性

每次调用记录 provider、model、base URL host、prompt version、input hash、request ID、
latency、token usage、finish reason、validation result、fallback reason 和 content hash。
不记录 API key；创业原文按 tenant 的数据保留策略处理。

限流和超时：

- 单 workspace 同时一个 naming job；
- 15 秒软超时、30 秒硬超时；
- 429/5xx 指数退避，最多一次自动重试；
- 用户点击使用稳定 idempotency key；
- 模型额度不可用时主页与工作室显示 degraded。

MiniMax 官方公布文本模型的 RPM/TPM 限制，但实际账户额度和 Token Plan 还可能受滚动窗口
影响，因此不能硬编码“可用”结论。

### 5.4 部署配置与启用判定

Provider 配置只从 AESE 服务进程环境读取。开发机使用仓库根目录 `.env` 保存
`MINMAX_API_KEY`、`MINMAX_API_BASE` 和 `MINMAX_MODEL`，文件权限必须为 `0600`；
`.env` 不提交。标准启动入口是 `scripts/deploy_aese_server.sh`，它以非执行式解析器
加载配置，三项缺一项、权限过宽或启动后 Provider 状态不符合预期时失败关闭，且不得把
密钥放入命令行、日志或状态响应。

启用验收分为两层，不能只凭配置文件或 UI 标签判定：

1. `GET /api/aese/v1/game/creative/status` 返回 `state=connected`、
   `provider=MiniMax` 和账户实际模型 ID；
2. 用有效 Workspace tenant session 完成一次名称生成，并在 CreativeJob 中看到
   MiniMax provider/model/host、非零 latency 与 token usage、`validation_result=valid`
   且没有 `fallback_reason`。

`AnalyzeIntent` 当前仍由确定性逻辑完成输入规范化和合同校验；外部 MiniMax 只负责
`GenerateNames`。页面应准确说明这一边界，不能把整个 onboarding 都描述成模型生成。
详细启动和恢复步骤见
[Enterprise Genesis Runbook](../runbooks/enterprise-genesis-game.md)，故障根因见
[SOL-008](../solutions/SOL-008-minimax-provider-startup-config.md)。

## 6. Logo 与 Embedding 的边界

- MiniMax M3 是文本模型，本阶段用于意图结构化、名称、品牌 brief 和解释。
- Logo 需要独立 image provider/模型；没有图像 provider 时只显示“几何字标 fallback”，
  不能称为 AI Logo。
- Qwen Embedding 用于未来知识检索，不参与公司命名，也不是 provisioning 前置条件。
- 名称和 Logo 资产仍为 candidate；只有玩家选择及 IAOS Capability 才能成为正式引用。

## 7. 主页信息架构

桌面：

```text
顶部：Enterprise Genesis / 我的企业 / 样板世界 / 功能说明 / 服务状态
Hero：从一个想法，创建一家真正运行在 IAOS 上的企业
主按钮：创建新企业
继续经营：workspace cards
创建过程：创业空间 -> AI 身份 -> 登记 -> 资本 -> 团队 -> 开业
验证 IAOS：租户隔离、流程审批、数字员工、事件审计
```

移动端保留一个主 CTA；workspace card 提供“继续”而不是整卡隐藏点击。所有触点至少
44×44 px，表单可返回且自动保存草稿。Provisioning 和 AI 调用分别显示进度，不混为
同一 loading 状态。

登录成功后，“我的企业”紧跟 Hero 展示；列表由 owner-scoped Workspace API 返回。
点击“继续游戏”必须先取得该 Workspace 的 Founder tenant session，再用服务端返回的
tenant/case/workspace 绑定进入游戏，不允许由玩家手填 tenant ID。

## 8. 安全和隔离验收

必须证明：

1. 两次“创建新企业”产生不同 tenant、workspace、World Run 和 case。
2. 两个 tenant 可使用相同 company short name，但互不可见。
3. 玩家 A 不能读取玩家 B 的 workspace/provisioning 状态。
4. 浏览器 token 没有 `platform.manage`、`platform.identity.bootstrap` 权限。
5. provisioning 任一步失败后 tenant 不可写、不可普通登录；重试无重复资产。
6. 只有 active tenant 能创建 FounderIntent 和 incorporation case。
7. CreativeJob、IAOS case 和 World Run 均携带同一 workspace correlation。

## 9. 完成标准

- 根地址进入产品主页，不需要复杂 URL。
- 从平台玩家身份创建全新的 IAOS tenant 和 AESE World Run。
- MiniMax 实际调用证据可见，fallback 明确标记。
- 玩家从空白 workspace 完成公司身份选择和 23 个 M9 工作项。
- 第二个独立 tenant 重复全链并通过 RLS 隔离验证。
- 主页、provisioning、身份工作室和游戏在 1440、1280、390 三视口通过。

### 9.1 财务治理 Runtime 安装恢复证据

新企业 provisioning 曾在 IAOS `runtime_installed` 检查点因共享财务治理表由数据库
初始化所有者持有、普通 Runtime 重复执行 owner-only DDL 而返回 502。IAOS
`main@732a9b1` 将既有 schema 路径改为只读验证 FORCE RLS、tenant policy 与 SoD
trigger，tenant Runtime 只执行 DML；全新和原失败幂等请求均已通过 AESE BFF 完成
8/8 checkpoint。Onboarding 将 provisioning 幂等键与草稿一同保存在浏览器，失败或
刷新后继续使用原键恢复同一 Workspace，仅在创建成功后清除。详细根因和恢复步骤见
IAOS `SOL-042`。

## 10. 非目标

- 把 tenant 当成现实工商注册主体。
- 允许匿名用户无限创建租户。
- 从浏览器调用 SaaS Ops 管理员 API。
- 在创建失败时自动删除已激活 tenant。
- 首版复杂集团多租户、租户合并或跨租户交易。
- 用 Qwen Embedding 替代 LLM，或用文本模型冒充 Logo 生成器。

## 11. 外部接口依据

- [MiniMax Text Generation](https://platform.minimax.io/docs/guides/text-generation)
- [MiniMax OpenAI-compatible API](https://platform.minimax.io/docs/api-reference/text-openai-api)
- [MiniMax Rate Limits](https://platform.minimax.io/docs/guides/rate-limits)
