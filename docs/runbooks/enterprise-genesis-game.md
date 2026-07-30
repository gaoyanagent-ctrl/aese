# Enterprise Genesis 企业创生游戏

## 从零创建独立企业空间

打开 `http://192.168.50.222:4173/`：

1. 输入游戏用户名登录。首次登录会关联当前浏览器此前创建企业所用的本机 player ID。
2. 在“我的企业”点击“继续游戏”恢复已有企业；或点击“创建新企业”。
3. 按五步向导填写创业项目、制造模板、区域/时区、真实性等级和数据保留确认。
4. 点击“创建空间并进入身份工作室”。
5. IAOS 服务端分配 Workspace、tenant、World Run 和案件编号，绑定真实 Player
   subject，安装 M9 Runtime；AESE 写入 World binding evidence 后才完成 smoke 和激活。
6. 成功后浏览器分别保存 Player 控制面会话与 `genesis_owner` tenant-scoped 会话；
   两者都没有 `platform_super_admin` 或平台管理权限。

进入新企业后不再直接填写设立表单：

1. 在 PixiJS 创始办公室选择虚拟头像；
2. 与数字员工“纪元”对话，依次选择产业、目标客户、核心产品和品牌性格；
3. 补充创业宣言，让数字员工调用 MiniMax M3 形成四组身份提案；
4. 选择名称并确认注册地址、经营范围；
5. 签署创始人指令后才调用 IAOS `incorporation.case.open`。

页面右侧主线任务显示当前目标；IAOS capability 与 evidence 属于治理层，不要求玩家
先理解技术节点。

AESE 的 Workspace/World 稳定引用保存在 `.aese-data/genesis-workspaces.json`，权限为
`0600`；tenant、owner membership、Runtime、checkpoint 和 session 的权威数据在 IAOS
生产控制面。同一认证 Player 和 idempotency key 重试返回同一个 tenant；不同创建动作
得到不同 tenant。

没有 IAOS Player session 时，loopback 本地开发仍可显式使用 dev adapter。正常产品
路径必须把当前 Player session 转发给受限的 Genesis Workspace API；不得把平台 token
写进 AESE 状态、Vite 环境变量或浏览器。

当前限制：游戏用户名和 player ID 映射保存在浏览器 localStorage，只有本机体验意义，
没有密码校验且不能跨设备找回。正式多人
部署必须使用 IAOS DES-062 的 Player Account/OIDC 与服务端 membership 授权，不能把
用户名或 localStorage player ID 当作安全身份。

## 1. 启动

先启动 IAOS `:8082` 与前端 `:3000`，再让 AESE 以 live projection 模式连接 IAOS：

```bash
go run ./cmd/aese-server \
  --listen :8090 \
  --pack-dir scenario-packs/hctm \
  --iaos-base-url http://127.0.0.1:8082

cd frontend
npm run dev
```

不传 `--iaos-base-url` 时，AESE 只提供确定性离线演示投影；它不会伪装成 IAOS
持久业务状态。

## 2. 进入游戏

从 IAOS 企业生命周期工作台打开 AESE，或直接访问：

```text
http://127.0.0.1:4173/#enterprise-genesis?tenant=tenant-hctm-genesis&case=<CASE_CODE>
```

浏览器需要当前 Workspace 的 tenant session。游戏把 token 传给 AESE，AESE 再读取
IAOS 的已校验 evidence bundle 和 23 个持久工作项；前端不访问数据库。

## 3. 操作主线

1. 在“企业身份工作室”输入创业构想，生成四组名称、英文名、口号、色彩与风险提示。
2. 选择候选，补充注册地址和经营范围，点击“确认身份并创建企业”。正式写入复用
   IAOS `incorporation.case.open`；AI 候选本身不改变业务事实。
3. 在任务页点击“进入操作”，直接在游戏内完成数字员工派遣、业务金额输入、G1–G7
   审批和系统任务；所有写入仍由 IAOS Runtime 治理，不需要跳转其他页面。
4. 到登记、开户或任命 World wait 时，在游戏进入政务、银行或任命办理。AESE 发送受治理
   Observation，再执行对应 commit Capability；刷新后只展示新的 committed projection。
5. 23 个工作项完成后进入“经营世界”，必须显示
   `enterprise_operational_ready`、100%、五 Agent、资金与证据。

## 4. AI 与素材边界

- 名称生成当前使用 provider-neutral 合同和确定性 fallback，保证离线演示可重复；
  未配置外部模型时不得显示为“已连接 LLM”。
- 首批企业城市背景由 OpenAI 图像生成工具制作，登记在
  `frontend/public/assets/enterprise-genesis/manifest.json`，包含尺寸、hash、来源、
  项目许可和保留规则。
- 品牌选择复用 IAOS 案件创建 Capability。首版未把按需 Logo 外部模型凭据写入仓库；
  Logo 生成失败必须使用文字/几何 fallback，不能阻塞企业成立。

## 5. 恢复与异常

- 刷新或重新打开相同 URL：按 evidence bundle + work items 重建，不依赖前端内存。
- AESE 重启：重新带 token 请求 projection；cursor 取 IAOS Journal/World 最大值。
- “我的企业”返回 `player_session_expired`：页面会先使用当前 IAOS 会话重试一次；成功
  后自动更新 Player 会话。若仍失败，返回 IAOS 重新登录并再次打开 Enterprise Genesis。
  不需要删除企业数据，也不要反复点击“刷新”。
- 旧企业点击“继续游戏”曾返回 Workspace session 404/502：新版本会读取当前玩家拥有的
  `0600` 本地 Workspace，并请求 IAOS 安全接管。IAOS 会验证当前 tenant 的 owner、
  chair、Founder Mandate、M9 Runtime 和原设立案；通过后自动兑换 session。失败时按
  页面给出的前置条件修复身份/Runtime/案件，禁止手工补数据库或重新执行 Founder
  bootstrap。
- 重复 World Observation：IAOS 业务事实级幂等返回同一事实，不新增业务变化。
- 登记/开户/任命拒绝、资本差异、Agent 越权或停用：业务状态保持在原节点；从 IAOS
  工作项进入补正、人工接管或重新审批，不在游戏画面私自跳步。
- 低性能、移动或 reduced-data：使用 DOM 列表与 2D 网格 fallback；所有任务和证据
  不依赖 PixiJS 点击。

## 6. 可复现验收

```bash
go test ./...

cd frontend
npm run typecheck
npm test -- --run
npm run build
npx eslint src/components/game src/game e2e/enterprise-genesis*.spec.ts --max-warnings 0
npx playwright test e2e/enterprise-genesis.spec.ts
npx playwright test e2e/enterprise-genesis-interactive.spec.ts

GX_LIVE=1 \
M9_EXPECT_READY=1 \
M9_CASE_CODE=<通过23工作项完成的CASE_CODE> \
npx playwright test e2e/enterprise-genesis-live.spec.ts
```

Playwright 固定验证 `1440×900`、`1280×720`、`390×844`，并在每个项目保存
`enterprise-genesis-live.png`。IAOS 全链测试见其
`TestIntegrationInteractiveWorkItemsFiveAgentsSevenGatesThreeWorldWaitsAndRestart`：
18 completed work items、G1–G7、三个 World wait、六次 Run/五个 distinct Agent，
并在第二个 World wait 后重建 Server 继续执行。
# RPG、精灵与多租户补充验收（2026-07-28）

1. 从企业大厅进入任一进行中的企业，进入带事件标记的地点。
2. 确认任务弹层显示地点、具名 NPC、三项经营选择和正式 IAOS 操作。
3. 点击室内物件，确认创始人移动后才打开状态卡；开启/关闭顶部音效按钮。
4. 完成任务后确认出现里程碑奖励；“治理档案”逐项显示 committed 企业大事记。
5. 登记/开户获批时确认显示虚构营业执照、印章套装和 U 盾生成素材。
6. 在 IAOS `http://127.0.0.1:3000/login` 使用 Founder 凭据登录，确认出现所有明确授权
   企业；进入任一企业后 `/profile` tenant 与选择项一致。
7. provisioning 失败后使用同一幂等键重试，workspace ID、tenant ID、world run 和 case
   code 必须保持不变。

完整证据见 `docs/reports/genesis-rpg-and-multi-tenant-acceptance.md`。
# 开业财务验证

完成“注入并核验实缴资本”后：

1. 进入企业总部，检查“开业财务中心”物件，状态应为“账套已启用”。
2. 打开右侧“治理档案”，应看到“实缴资本开业入账”。
3. 凭证必须显示借 `1002 银行存款`、贷 `4001 实收资本`，两个合计相等。
4. 证据引用必须以 `iaos:finance:` 开头；页面没有财务事实时会显示等待资本到账，不会根据公司现金伪造凭证。
5. 同一档案继续检查银行日记账、总账和开业资产负债表：
   - 银行日记账显示“实收资本到账”、本次收入和滚动余额；
   - 总账显示 `1002` 与 `4001` 的期末余额；
   - 开业资产负债表显示资产等于所有者权益、负债为零，并标记“资产 = 负债 + 所有者权益”。
6. 这些视图必须来自同一 IAOS 已过账凭证；若报表不平，企业最终就绪节点应返回
   `finance_opening_readiness_failed`，不得继续完成 M9。
7. 点击总部“开业财务中心”，确认它不会遮挡“治理与经营会议桌”，并打开用途说明。
8. 点击“查看系统账务”，应进入 IAOS `#finance_workspace` 的银行日记账/总账视图；
   点击“查看财务报表”，应进入同一页面的开业资产负债表视图。两个入口自动携带当前
   tenant 和 case，IAOS 仍会用当前登录身份做权限校验。
9. 在 IAOS 数据模型工坊启用某 Entity 的“在侧边栏显示菜单”并发布：
   `platform_super_admin` 应无需刷新看到菜单；普通角色须先在“角色与数据范围”授予
   `menu.<entity_code>` READ。
