# M10 Genesis Plant Build Runbook

> 本 Runbook 分成两个边界明确的部分：第一部分验收当前可操作的“设施需求 → Agent 候选 → 人工审阅”在线纵切；第二部分只验收历史 reference replay。World 调研、选址/投资审批、项目/WBS、施工、付款、验收和工程财务尚未接通，不能据此宣称完整 M10 完成。

## 交互规划纵切验收

### 前置条件

- IAOS `:8082`、AESE `:8090` 和前端 `:4173` 健康。
- 用户从 Genesis 企业入口进入当前企业，浏览器存在有效 IAOS token、tenant、workspace 和 `aese_genesis_case_code`。
- M9 案件已产生法人编码、已批准预算、币种和已过账银行现金。页面不提供替代这些事实的手工输入。
- AESE 启动时已配置 MiniMax provider；若页面显示“未启用外部模型”，Agent 路径应失败关闭。

### UI 操作

1. 进入“工厂建设 Campaign”。页面顶部应出现“设施需求与候选方案”和“功能说明”。
2. 检查“可用现金快照”“已批准预算快照”为只读，并分别显示 `gl:` 与 `budget:` 来源；若无数据，先回 M9 完成预算和资本入账，不能在 M10 页面伪造。
3. 填写目标区域、设施用途、最小面积、最小电力、目标日期、候选数量、允许类型、投资申请额、最低现金保留额、业务偏好和修订原因。
4. 点击“保存需求并让 Agent 生成候选”。候选数量必须等于填写值（2–8），每张卡必须展示金额/工期区间、依据、假设、待验证事实、风险、来源、置信度和模型/prompt 证据。
5. 对一个候选选择“采纳调研”，填写至少 6 个字符的业务理由并点击“提交审阅到 IAOS”。按钮变为“已保存审阅”。再用另一候选验证“退回重生成”或“淘汰”。
6. 展开页面底部的“已封存的确定性参考回放”，确认它有 `fixture-only` 提示且不会自动写入上方候选列表。

### API 与权威证据

浏览器 Network 中应出现：

| 方法与路径 | 预期 | 权威含义 |
| --- | --- | --- |
| `GET /api/aese/v1/world/plant-build/planning-status` | `connected` 或明确 `not_configured` | 只说明 Agent provider 状态 |
| `GET /api/aese/v1/world/plant-build/financial-constraints?case_code=...` | 200，带 source refs 与 snapshot hash | IAOS 权威现金/预算只读快照 |
| `POST /api/aese/v1/world/plant-build/proposals` | 200，`authority_status=committed` | IAOS 已分别保存 Requirement 与 candidate-only ProposalSet |
| `POST /api/aese/v1/world/plant-build/reviews` | 201，`status=committed` | IAOS 已保存当前用户的 ProposalReview |

写操作不得从浏览器直接请求 IAOS `:8082`。AESE BFF 使用当前用户身份调用 IAOS
`POST /api/v1/genesis/plant/interactive/actions`，只允许：

- `facility.requirement.define`
- `site.proposal.record`
- `site.proposal.review`

IAOS 侧应能读取最新 Requirement、ProposalSet 和 Review；每次提交同时产生 tenant-scoped Audit 与 Outbox。Agent 模型、prompt、request、token 和输入/输出 hash 保存在 AESE CreativeJob 技术证据中。两类证据必须分别存在，不能用页面展示或 Agent 文本替代。

### 失败与恢复验收

- 清除/失效 token：财务读取和审阅返回 401，页面不得显示已保存。
- 切换 tenant 或伪造 tenant：IAOS RLS/身份校验失败，不得跨企业读取或写入。
- 关闭模型配置：状态显示未启用，生成按钮不可用；不得出现固定候选。
- 提交同一 revision 相同输入：幂等返回既有结果；同一幂等键不同输入返回冲突。
- 在生成候选后改变 IAOS 现金或预算：保存 Proposal 时拒绝旧 snapshot，用户必须重读并创建 Requirement 新 revision。
- 同一 ProposalSet revision 并发审阅：冲突方刷新最新 Review 后重新决定，不得覆盖。

### 当前限制

- “人工新增候选”当前只加入本地审阅列表，不是 IAOS 权威记录，不能发起外部调查。
- “采纳调研”只保存人工意图；`site.investigation.request`、World Observation 和等待节点尚未实现。
- 尚无交互式选址/投资审批、场地控制、项目/WBS、合同、施工、变更、付款、验收、AP/CIP/总账闭环。
- 现场部署、纯人工路径和完整断线/重启/并发证据属于 S5，完成前 M10 保持 `Interactive Revision Pending`。

## 历史 reference replay 页面验收

## 页面验收

确认 IAOS `:8082`、AESE `:8090` 和前端 `:4173` 健康，访问 `/#world-incorporation`，点击“工厂建设 Campaign”。初态必须显示消费 M9 机器资格；单步 9 次后显示 `M11 eligible`、七个空间节点、World/IAOS 进度、现金/承诺/应付/实付及已关闭 discrepancy。

## 离线验证

```bash
go test ./...
go vet ./...
go run ./cmd/aese world validate world-packs/hctm-genesis
cd frontend
npm test -- --run
npm run typecheck
npm run build
npx playwright test e2e/plant-build.spec.ts
```

`internal/plantbuild` 同时验证交互合同和 reference fixture。交互合同覆盖金额精度、日期、候选数量、来源、重复候选、坏 JSON、未配置模型和审阅理由；HTTP 测试覆盖财务快照 BFF、Requirement/Proposal Capability 提交、CreativeJob 幂等证据和 Review 身份/提交。reference fixture 继续验证三个虚构候选、先硬约束后评分、10 帧状态机、100 次 hash、snapshot/restore/reset、资金和验收门。绿地及代建方案因该 fixture 的预算/日期失败；租赁标准厂房方案在该 fixture 的 1,500 万预算内获选。此结果不能推断其他企业应选择同一模式或金额。

## 历史 IAOS 项目治理与恢复

`POST /api/v1/genesis/plant/actions` 使用严格合同与 `expected_version`，支持 site evaluate/approve、project execute/rebaseline/accept、payment approve。相同幂等输入为 no-op；并发版本冲突为 409；越权、自批、超预算、过期 mandate 或未验收付款失败关闭。治理记录、committed outcome journal 和 Outbox 同事务。

World 页面刷新后从确定性 campaign API 恢复；正式在线恢复以 IAOS journal cursor 为事实，SSE 只作通知。AESE 不直写 IAOS 数据库，也不把 IAOS 计划进度当作施工事实。
