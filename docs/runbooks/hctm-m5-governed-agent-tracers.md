# HCTM M5 受治理 Agent Tracer 运行手册

日期：2026-07-19  
适用 pack：`hctm@0.1.0`

## 1. 目的和边界

本手册用于配置并运行计划、质量和经营分析三个只读 Agent tracer。

- 九个上下文工具均通过 IAOS AI Tool Registry 调用，`source_ref` 固定为 `entity.records`。
- 工具的实体、返回字段、过滤字段、排序和最大行数由 `scenario-packs/hctm/agent-tools.json` 固定；调用者不能传入 tenant、表名、字段名或 SQL。
- `agent-setup` 和 `agent-run` 默认 dry-run；只有显式 `--apply` 才访问 IAOS 写端点或执行真实 Tool call。
- Agent 只输出建议，不释放采购单、不调整排产、不隔离库存；所有建议均标记 `requires_human_confirmation=true`。
- 真实 `agent-run --apply` 不修改业务表或 Outbox，但会为九次 Tool call 追加九条 `ai_tool_call` 和三十六条 milestone event。

## 2. 前置条件

### 2.1 服务健康

```bash
curl -fsS http://127.0.0.1:8082/health
curl -fsS http://127.0.0.1:8082/ready
```

### 2.2 正式创建并激活租户

`tenant-hctm` 必须先进入 SaaS Operations 的正式租户生命周期。不能只生成一个带任意 tenant claim 的开发 token，再把它当作租户已完成 provision 的证据。

本地环境可用 platform admin token 创建租户；若租户已存在，先查询并确认其状态，不要重复创建：

```bash
export PLATFORM_TOKEN="$(curl -fsS \
  'http://127.0.0.1:8082/api/v1/dev/token?tenant_id=tenant-001&roles=admin' \
  | jq -r .token)"

curl -fsS -X POST http://127.0.0.1:8082/api/v1/platform/tenants \
  -H "Authorization: Bearer $PLATFORM_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"tenant_id":"tenant-hctm","display_name":"华辰热管理系统集团有限公司"}' \
  | jq .

curl -fsS -X POST \
  http://127.0.0.1:8082/api/v1/platform/tenants/tenant-hctm/activate \
  -H "Authorization: Bearer $PLATFORM_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"reason":"HCTM M5 smoke checks passed"}' \
  | jq .
```

已存在环境用以下命令确认 `status=active`：

```bash
curl -fsS \
  -H "Authorization: Bearer $PLATFORM_TOKEN" \
  http://127.0.0.1:8082/api/v1/platform/tenants/tenant-hctm \
  | jq '{tenant_id,status,display_name}'
```

开发 token 仅用于本地验收。演示或生产环境必须换成真实登录或服务身份：

```bash
export IAOS_TOKEN="$(curl -fsS \
  'http://127.0.0.1:8082/api/v1/dev/token?tenant_id=tenant-hctm&roles=admin' \
  | jq -r .token)"
test -n "$IAOS_TOKEN" && test "$IAOS_TOKEN" != null
```

### 2.3 场景状态

先按 M3/M4 手册完成 scenario apply 和 canonical replay。Agent tracer 的最小在线前置状态包括：

- `SO-202607-0001` 及订单行可查询。
- `PO-202607-0001` 为延期采购单。
- `LAS-WLD-02` 为 `maintenance`。
- `IQC-202607-0002` 为 `failed`，含批次、缺陷、不合格数量和严重度。
- O2D 已生成关联 work order。

## 3. 离线验证

```bash
cd /iaos/aese

go test ./...
go vet ./...
go run ./cmd/aese validate ./scenario-packs/hctm
go run ./cmd/aese inspect ./scenario-packs/hctm --json
jq empty \
  scenario-packs/hctm/agent-tools.json \
  scenario-packs/hctm/schemas/agent-tools.schema.json
```

`agent-tools.json` 应包含五份最小 metadata schema 和九个只读工具：

```text
hctm.product.read
hctm.sales_order.read
hctm.sales_order_line.read
hctm.inventory.read
hctm.bom.read
hctm.purchase_order.read
hctm.equipment.read
hctm.inspection_order.read
hctm.work_order.read
```

## 4. Tool Bundle dry-run 和 apply

先查看 setup 计划；此命令不创建 metadata 或 Tool：

```bash
go run ./cmd/aese agent-setup ./scenario-packs/hctm \
  --target http://127.0.0.1:8082 \
  --tenant tenant-hctm
```

预期关键字段：

```text
mode: dry-run
metadata_schemas: 5
tools_created: 9
tool_keys: 9 entries
```

确认目标、租户和工具清单后显式应用：

```bash
go run ./cmd/aese agent-setup ./scenario-packs/hctm \
  --target http://127.0.0.1:8082 \
  --tenant tenant-hctm \
  --apply
```

首次空环境预期创建九个 Tool。重复执行同一命令必须收敛为：

```text
tools_created: 0
tools_updated: 9
```

重复 setup 会更新 operator-owned manifest 并重新启用工具，不应创建重名 Tool。

可通过 IAOS 查询清单：

```bash
curl -fsS \
  -H "Authorization: Bearer $IAOS_TOKEN" \
  'http://127.0.0.1:8082/api/v1/ai/tools?include_disabled=true&limit=200' \
  | jq '{total,keys:[.items[].tool_key]}'
```

## 5. Agent tracer dry-run 和 apply

dry-run 只列出九个预期工具，不调用 IAOS：

```bash
go run ./cmd/aese agent-run ./scenario-packs/hctm \
  --story order-expedite-01 \
  --run-id hctm-m5-agent-dry-001 \
  --target http://127.0.0.1:8082 \
  --tenant tenant-hctm
```

真实运行：

```bash
go run ./cmd/aese agent-run ./scenario-packs/hctm \
  --story order-expedite-01 \
  --run-id hctm-m5-agent-apply-001 \
  --target http://127.0.0.1:8082 \
  --tenant tenant-hctm \
  --apply \
  | tee /tmp/hctm-m5-agent-apply-001.json
```

成功输出必须满足：

- `mode=apply`、`correlation_id=corr-so-202607-0001`。
- `tool_evidence` 有九项，每项包含非空 `call_id`。
- 三个 Agent 均为 `status=suggested`、`requires_human_confirmation=true`。
- 计划 Agent 根据当前在线库存和 `PO-202607-0001` 报告铝板仍缺 `7600`，不能复用 Preview 的静态答案。
- 质量 Agent 为 `partial`：可确认 `SUP-BETA-AL`、`BETA-20260712-01`、`SURFACE_SCRATCH` 和不合格 300；由于未记录加严检验且合格数量未放行，必须列出数据缺口。
- 经营分析 Agent 为 `partial`：缺少完工入库、发运和实际成本事实时，必须明确“不能判断”最终交付数量、缺口或利润影响。

快速检查：

```bash
jq '{
  calls:(.tool_evidence|length),
  agents:[.agents[]|{
    agent_key,completeness,requires_human_confirmation,summary,data_gaps
  }]
}' /tmp/hctm-m5-agent-apply-001.json
```

## 6. 重复执行和零业务写入

用新的 run ID 再运行一次：

```bash
go run ./cmd/aese agent-run ./scenario-packs/hctm \
  --story order-expedite-01 \
  --run-id hctm-m5-agent-apply-002 \
  --target http://127.0.0.1:8082 \
  --tenant tenant-hctm \
  --apply \
  | tee /tmp/hctm-m5-agent-apply-002.json
```

第二次调用应产生新的调用审计，但不得改变业务状态。现场验收的前后计数为：

| 计数 | 运行前 | 运行后 | 预期增量 |
| --- | ---: | ---: | ---: |
| 受检业务对象合计 | 24 | 24 | 0 |
| `sys_outbox` | 39 | 39 | 0 |
| `ai_tool_call` | 11 | 20 | 9 |
| `ai_tool_call_event` | 44 | 80 | 36 |

调用记录可通过受治理 API 查询：

```bash
curl -fsS \
  -H "Authorization: Bearer $IAOS_TOKEN" \
  'http://127.0.0.1:8082/api/v1/ai/tool-calls?limit=50' \
  | jq '[.items[]|{id,tool_key,status,session_id,created_at}]'
```

从输出选择 `call_id` 后核对四个 milestone：

```bash
export CALL_ID='<call_id>'
curl -fsS \
  -H "Authorization: Bearer $IAOS_TOKEN" \
  "http://127.0.0.1:8082/api/v1/ai/tool-calls/$CALL_ID/events" \
  | jq .
```

预期顺序为 `queued`、`validated`、`dispatch_started`、`succeeded`。

## 7. 跨租户失败关闭

使用另一个正式租户的 token 查询 Tool 清单：

```bash
export OTHER_TOKEN="$(curl -fsS \
  'http://127.0.0.1:8082/api/v1/dev/token?tenant_id=tenant-other&roles=admin' \
  | jq -r .token)"

curl -fsS \
  -H "Authorization: Bearer $OTHER_TOKEN" \
  'http://127.0.0.1:8082/api/v1/ai/tools?include_disabled=true&limit=200' \
  | jq '{total,items}'
```

现场预期为 `total=0`。直接调用 `hctm.*` Tool 应返回 404；不得通过请求头或输入字段覆盖 JWT tenant。

## 8. 故障排查

| 现象 | 原因与处理 |
| --- | --- |
| setup 返回 metadata 404/500 | IAOS Platform 未部署包含 `entity.records` dispatcher 的最终版本，或 metadata schema upsert 失败；先检查 `/ready` 和 Platform 日志。 |
| Tool call 404 | 未执行 `agent-setup --apply`、目标 tenant 错误，或跨租户 RLS 正常隐藏该 Tool。 |
| Tool call `invalid_input` | 输入含 manifest 未允许的 filter/字段；不要把 SQL、tenant 或字段清单传给 Tool。 |
| 计划 Agent 未报告 7600 缺口 | 在线库存或 PO 已变化；以当前受治理事实重算，不要强行匹配历史证据。 |
| 质量 Agent 不是 `partial` | 检查 `inspection_level` 和 `accepted_qty`；当前现场分别为 `normal` 和 `0`，不能声称已加严或剩余数量可用。 |
| 经营分析给出 11,700/300 或利润数值 | 这是错误地复用了 Preview 结果；在线缺少 shipment/cost facts 时必须失败关闭。 |

实际执行证据见 [HCTM M5 Agent Evidence](../reports/hctm-m5-agent-evidence.md)。
