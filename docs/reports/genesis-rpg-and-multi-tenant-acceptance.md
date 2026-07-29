# Enterprise Genesis RPG 与多租户验收

日期：2026-07-28

## 验收范围

- 全 17 M9 工作项专属 RPG 剧情入口；
- 五名透明角色精灵及执照、印章、U 盾 collectible；
- 室内移动、物件接近、可关闭音效、里程碑奖励和大事记；
- Founder IAOS 新租户选择；
- 两租户 JWT/profile 隔离、交叉案件读取和 provisioning 失败恢复。

## 结果

| 检查 | 结果 |
| --- | --- |
| AESE frontend unit | 18 files / 49 tests passed |
| AESE production build | passed |
| 1440×900 / 1280×720 / 390×844 | Playwright 3/3 passed |
| Workspace independent allocation / owner session | passed |
| Injected provisioning failure retry | same workspace and tenant recovered |
| IAOS API package tests | passed |
| IAOS frontend TypeScript / production build | passed |
| Founder global login | HTTP 200, `multiple`, 8 explicitly assigned tenants |
| Two selected tenant tokens | JWT and `/profile` tenant claims matched |
| Tenant A reads tenant B case | 404 |
| Tenant B reads tenant A case | 404 |

测试没有输出或保存 JWT、密码及其他凭据。8 个租户数量是当前本机 Founder assignment
状态，不是产品固定值。

## 已知边界

当前 workspace catalog/checkpoint 仍保存在 AESE 本机版本化 JSON store；IAOS
`/genesis/workspaces` 生产控制面、持久 CreativeJob 和完整分布式补偿 saga 仍属于
PLAN-GXZ-001 的后续平台化工作，不影响当前单机沙盘验收。

