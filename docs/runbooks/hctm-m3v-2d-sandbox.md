# HCTM M3V 2D 企业沙盘运行手册

## 1. 用途

本手册用于本地启动和验收华辰热管理系统集团苏州制造基地的只读 2D 企业沙盘。该界面播放 `order-expedite-01` 的七幕、22 个确定性事件，不写 IAOS，也不代表实时业务结果。

## 2. 环境要求

- Node.js 22 或兼容的当前 LTS。
- npm 9 或更高版本。
- 首次运行需要安装前端依赖。

## 3. 启动

```bash
cd /iaos/aese/frontend
npm install
npm run dev
```

默认地址：`http://localhost:4173/`。开发服务绑定 `0.0.0.0`，同一网络中的访问者可使用宿主机地址和端口 `4173`。

页面右上方必须显示 `PREVIEW`。这表示当前数据来自场景包中的静态 `preview.json`，不是 IAOS 在线快照。

## 4. 操作

- 顶部：播放/暂停、前后单步、1×/2×/4×、重置。
- 中央：缩放或适配 A 线画布，点击节点查看对象详情。
- 七幕条：直接跳转到订单、MRP、供应延期、设备停机、质量、生产或发运场景。
- 右侧：查看事件严重级别和计划、质量、经营分析 Agent 的确定性建议。
- 底部：核对订单需求、可用成品、缺料、产能和交付风险。
- 移动端：通过“企业 / A 线画布 / 事件与 Agent”三个标签切换区域。

最终状态应显示总需求 12,000、累计可发/实发 11,700、第二批实发 2,700、短缺 300，并将交付风险标记为严重。

## 5. 验证

```bash
cd /iaos/aese/frontend
npm run typecheck
npm run lint
npm test
npm run build
npm run test:e2e
npm audit --audit-level=high
```

Playwright 自动覆盖 1440×900、1280×720 和 390×844，并在 `frontend/test-results/` 生成三个完成态截图。

## 6. 故障排查

- 显示“场景无法加载”：检查 `scenario-packs/hctm/stories/order-expedite-01/preview.json`，运行 `npm test` 获取具体合同错误。
- 画布为空：确认 preview 数据仍有 14 个节点和 13 条连线，并确认容器不是零高度。
- Playwright 浏览器缺失：执行 `npx playwright install chromium`。
- 端口冲突：用 `npm run dev -- --port <port>` 临时换端口；自动化测试仍按配置使用 4173。

## 7. 边界

- 不在浏览器中重新计算 MRP、排产、质检或 Agent 决策。
- 不从该页面写数据库、发布 NATS 事件或执行建议。
- 在线化时新增 `IaosScenarioDataSource`，组件继续只依赖 `SandboxScenario`。
