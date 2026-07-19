# HCTM M5 Governed Agent Tracer Evidence

日期：2026-07-19。范围：计划、质量和经营分析三个只读 Agent tracer，以及 IAOS metadata 约束的 `entity.records` AI Tool 调用链。

## 实现边界

- `scenario-packs/hctm/agent-tools.json` 持有五份最小 metadata schema 和九个 HCTM Tool manifest；IAOS 平台代码只实现通用、metadata 约束的 `entity.records` query dispatcher，不包含 HCTM 专用业务规则。
- Tool metadata 固定 entity、返回字段、允许 filter、排序和最大行数。调用输入不能选择 tenant、表、列或 SQL。
- `agent-setup` 默认 dry-run；显式 `--apply` 后通过 IAOS metadata API 和 AI Tool Registry 创建或更新 Tool。
- `agent-run` 默认 dry-run；显式 `--apply` 后逐一调用九个 Tool，使用真实 IAOS 结果构建三个 recommendation envelope。
- Agent 输出只包含建议和待人工确认动作，不修改订单、采购、库存、设备、检验或工单，不产生业务 Outbox。
- 每个 Tool call 通过 Registry 的 enabled、权限、输入 Schema 和 dispatcher 边界，并持久化一条 `ai_tool_call` 及四条 milestone event。

## 现场前置

本次验收先通过 SaaS Operations 正式创建并激活 `tenant-hctm`，再使用绑定该租户的本地开发 token。仅生成包含任意 tenant claim 的 token 不作为租户 provision 完成证据。

M3/M4 场景已完成 apply/replay，在线状态至少包括销售订单和订单行、产品/BOM/库存、两张采购单、`LAS-WLD-02`、`IQC-202607-0002` 和三张 O2D 工单。三类异常继续保持 M4 已验证状态：供应延期、设备停机和来检失败均可查询。

## Tool Bundle 验收

| 验收 | 结果 |
| --- | --- |
| Bundle 结构 | 5 份 metadata schema，9 个只读 `entity.records` Tool |
| 默认安全模式 | `agent-setup` dry-run 只报告计划，不写 IAOS |
| 首次 apply | 九个 Tool 在 `tenant-hctm` 注册并启用 |
| 重复 setup | `tools_created=0`、`tools_updated=9`，未产生重名 Tool |
| Tool 类型 | 全部为 `query`、`risk_level=low`、`confirmation_mode=none` |
| 跨租户 | `tenant-other` 查询 `include_disabled=true` 仍为 `total=0` |

九个 Tool 覆盖 `product`、`sales_order`、`sales_order_line`、`inventory`、`bom`、`purchase_order`、`equipment`、`inspection_order` 和 `work_order`。这些是三类 Agent 回答第一条 HCTM 故事所需的最小在线上下文，不扩展为通用 ERP/MES Agent 数据面。

## 三 Agent 实际输出

### 计划 Agent

状态：`complete_for_current_risk`，置信度 `high`，需要人工确认。

计划 Agent 从当前在线订单、成品库存、BOM、铝板库存、主采购单和设备状态重新计算风险。现场当前事实显示：

- 订单需求、成品库存和净生产需求均来自 IAOS Tool 返回值。
- 铝板需求由 BOM 用量计算。
- 当前铝板库存与 `PO-202607-0001` 合计后仍缺 `7600`。
- `PO-202607-0001` 已延期，`LAS-WLD-02` 为 `maintenance`。
- 建议只形成“评估备选采购、准备加班、保留分批发运/重承诺草稿”，不自动执行业务动作。

这与 M3V Preview 的预计算答案不同，证明 Agent tracer 使用的是当前 IAOS 事实，而不是复制静态沙盘结果。

### 质量 Agent

状态：`partial`，置信度 `high`，需要人工确认。

已验证事实：

- 供应商：`SUP-BETA-AL`。
- 检验单：`IQC-202607-0002`。
- 收货号：`GR-202607-0002`。
- 批次：`BETA-20260712-01`。
- 缺陷：`SURFACE_SCRATCH`。
- 不合格数量：300。
- 严重度：`major`。
- 状态：`failed`。

质量 Agent 建议隔离不合格数量、禁止直接投产、保留供应商/收货/批次追溯，并只在质量放行后使用合格数量。输出保持 `partial`，因为现场 `inspection_level=normal`，没有证明已执行 HCTM 规格中的加严检验；同时 `accepted_qty=0`，不能声称剩余 1,700 张已放行。根因证据也不足，因此只允许发起供应商纠正措施草稿，不能断言根因。

### 经营分析 Agent

状态：`partial`，置信度 `medium`，需要人工确认。

当前 IAOS 可以证明订单需求以及供应延期、设备停机和来料不良三类风险，但没有可供本 tracer 使用的完工入库、发运和实际成本事实。因此输出显式声明：

- 不能判断最终交付数量和缺口。
- 不能判断利润影响。
- 成本影响只能作采购加急、加班和质量成本压力的定性提示。
- 数据缺口固定报告为 `finished_goods_receipt`、`shipment_dispatch` 和 `cost_actuals`。

经营分析没有把 Preview 的累计发运 11,700、短缺 300 或成本估计冒充在线事实，满足失败关闭要求。

## 重复执行与审计证据

第二次完整 `agent-run --apply` 调用九个 Tool。执行前后计数：

| 计数 | 执行前 | 执行后 | 增量 | 结论 |
| --- | ---: | ---: | ---: | --- |
| 受检业务对象合计 | 24 | 24 | 0 | 只读 tracer 未修改业务对象 |
| `sys_outbox` | 39 | 39 | 0 | 未伪造领域事件或业务动作 |
| `ai_tool_call` | 11 | 20 | 9 | 每个 Tool 调用均有独立审计 |
| `ai_tool_call_event` | 44 | 80 | 36 | 每个调用追加四个 milestone |

每个成功调用的 milestone 顺序为：

```text
queued
validated
dispatch_started
succeeded
```

重复运行使用新的 run/session ID，因此预期追加新的调用审计；“幂等”在此指业务状态和 Outbox 零变化，不是抹掉合法的只读调用历史。

## 租户隔离证据

使用 `tenant-other` token 查询 AI Tool Registry，包含 disabled Tool 的总数仍为 0。`tenant-hctm` 的九个 Tool 和调用记录由 RLS 隔离，不能通过请求输入传入 tenant 或动态 SQL 绕过。

这一结果证明 Tool Registry visibility 的租户边界。更完整的跨租户 dispatcher SQL、非法 field/filter、unsafe identifier 和 metadata 缺失路径由 IAOS real-PostgreSQL integration 与单元测试负责锁定。

## 验证命令

AESE：

```bash
go test ./...
go vet ./...
go run ./cmd/aese validate ./scenario-packs/hctm
go run ./cmd/aese inspect ./scenario-packs/hctm --json
go run ./cmd/aese agent-setup ./scenario-packs/hctm \
  --target http://127.0.0.1:8082 --tenant tenant-hctm
go run ./cmd/aese agent-setup ./scenario-packs/hctm \
  --target http://127.0.0.1:8082 --tenant tenant-hctm --apply
go run ./cmd/aese agent-run ./scenario-packs/hctm \
  --story order-expedite-01 --run-id hctm-m5-agent-evidence \
  --target http://127.0.0.1:8082 --tenant tenant-hctm --apply
```

IAOS 最终验证范围包括 Platform/AI Tool 单元测试和 vet、真实 PostgreSQL integration、code-map freshness、Platform redeploy 及 `/health`、`/ready` 检查。具体操作步骤见 [M5 运行手册](../runbooks/hctm-m5-governed-agent-tracers.md)。

## 剩余边界

- 质量建议仍为 `partial`：加严检验执行事实、合格数量放行和根因证据未齐。
- 经营分析仍为 `partial`：完工入库、发运和实际成本 Tool 尚未进入 M5 最小合同。
- 本切片没有自动 NATS Agent consumer、模型生成或自主行动。AESE CLI 负责确定性演示编排，业务变更仍必须进入后续 IAOS Capability、Process、Policy/Decision 和审批边界。
- M6 可以消费稳定 recommendation envelope，但不得在 2D 浏览器端复制 Agent 计算逻辑或把 Preview 数据冒充在线结果。
