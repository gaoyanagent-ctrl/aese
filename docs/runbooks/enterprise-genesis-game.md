# Enterprise Genesis 企业创生游戏

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

浏览器需要已有 `founder-principal` 登录 token。游戏把 token 传给 AESE，AESE 再读取
IAOS 的已校验 evidence bundle 和 18 个持久工作项；前端不访问数据库。

## 3. 操作主线

1. 在“企业身份工作室”输入创业构想，生成四组名称、英文名、口号、色彩与风险提示。
2. 选择候选，补充注册地址和经营范围，点击“确认身份并创建企业”。正式写入复用
   IAOS `incorporation.case.open`；AI 候选本身不改变业务事实。
3. 在任务页点击“在 IAOS 处理”，进入对应工作项业务表单。Agent task、G1–G7、
   人工接管和审批均由 IAOS Runtime 治理。
4. 到登记、开户或任命 World wait 时，可在游戏点击“模拟世界同意”。AESE 发送受治理
   Observation，再执行对应 commit Capability；刷新后只展示新的 committed projection。
5. 18 个工作项完成后进入“经营世界”，必须显示
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

GX_LIVE=1 \
M9_EXPECT_READY=1 \
M9_CASE_CODE=<通过18工作项完成的CASE_CODE> \
npx playwright test e2e/enterprise-genesis-live.spec.ts
```

Playwright 固定验证 `1440×900`、`1280×720`、`390×844`，并在每个项目保存
`enterprise-genesis-live.png`。IAOS 全链测试见其
`TestIntegrationInteractiveWorkItemsFiveAgentsSevenGatesThreeWorldWaitsAndRestart`：
18 completed work items、G1–G7、三个 World wait、六次 Run/五个 distinct Agent，
并在第二个 World wait 后重建 Server 继续执行。
