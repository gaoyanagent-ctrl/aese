# Enterprise Genesis RPG 与多租户验收

日期：2026-07-28

## 验收范围

- 全 23 M9 工作项专属 RPG 剧情入口；
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
| IAOS production Workspace checkpoints | 8/8 completed |
| Fresh Workspace M9 work items | 23/23 completed |
| Formal approval gates | 7/7 consumed |
| Trusted AESE World waits | 3/3 completed |
| Agent runs | 6 runs / 5 service Agents |
| Final lifecycle state | `enterprise_operational_ready` |

测试没有输出或保存 JWT、密码及其他凭据。8 个租户数量是当前本机 Founder assignment
状态，不是产品固定值。

## 生产控制面复验（2026-07-29）

从 IAOS `/api/v1/genesis/workspaces` 创建全新 Workspace，经 AESE World evidence
确认后签发仅含 `genesis_owner` 的租户会话。该会话创建设立案后产生 23 个持久工作项，
逐项完成 7 个正式审批门、3 个可信 World Observation 和 6 次 Agent run（覆盖 5 个
service Agent），最终进入 `enterprise_operational_ready`。

Workspace catalog、8 个 provisioning checkpoint、成员关系和失败重试均已迁入 IAOS
全局控制面；AESE 本地 adapter 仅作为显式开发回退。CreativeJob 使用持久 store，并记录
provider、model、request id、token usage、prompt version、workspace 与 correlation 证据。

本报告不保存 Workspace ID、JWT、密码或其他凭据。F15–F35 的完整 AP/AR、资产、成本、
月结与管理报表属于 M10–M13，不属于 M9 完成口径。
