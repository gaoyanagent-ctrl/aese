# HCTM M3 本地运行手册

日期：2026-07-19  
适用 pack：`hctm@0.1.0`

## 1. 目的和安全边界

本手册用于离线校验、检查、IAOS dry-run、显式 apply、O2D replay、verify 和 reset 计划。

- `validate`、`inspect` 完全离线。
- `apply`、`replay` 默认只读 dry-run；只有显式 `--apply` 才发起写请求。
- token 只从 `IAOS_TOKEN` 或临时 `--token` 参数读取，不进入 pack 或 run summary。
- 正式事件只走 IAOS 订单业务入口与 Outbox；不提供 direct NATS 捷径。
- reset 通过 IAOS 受治理 endpoint 只清理 L2 故事状态和派生 O2D 状态，保留 L1 公共主数据。

## 2. 离线验收

在 `/iaos/aese` 执行：

```bash
go test ./...
go vet ./...
go run ./cmd/aese validate ./scenario-packs/hctm
go run ./cmd/aese inspect ./scenario-packs/hctm --json
```

预期关键结果：

```text
valid: hctm@0.1.0 (19 record sets, 1 stories)
master_records: 80
initial_records: 15
events: 22
assertions: 17
correlation_id: corr-so-202607-0001
```

Schema 独立校验：

```bash
python3 - <<'PY'
import json
from pathlib import Path
from jsonschema import Draft202012Validator

root = Path("scenario-packs/hctm")
pairs = [
    ("schemas/manifest.schema.json", "manifest.json"),
    ("schemas/record-set.schema.json", "master-data/organization.json"),
    ("schemas/record-set.schema.json", "stories/order-expedite-01/initial-state.json"),
    ("schemas/event-sequence.schema.json", "stories/order-expedite-01/events.json"),
    ("schemas/expected-outcomes.schema.json", "stories/order-expedite-01/expected-outcomes.json"),
]
for schema_name, instance_name in pairs:
    schema = json.loads((root / schema_name).read_text())
    instance = json.loads((root / instance_name).read_text())
    Draft202012Validator(schema).validate(instance)
    print("valid", instance_name)
PY
```

## 3. IAOS 前置检查

```bash
curl -fsS http://127.0.0.1:8082/health
curl -fsS http://127.0.0.1:8082/ready
```

本地开发环境通过显式 dev endpoint 获取 `tenant-hctm` token，并只在当前 shell 保存：

```bash
export IAOS_TOKEN="$(curl -fsS \
  'http://127.0.0.1:8082/api/v1/dev/token?tenant_id=tenant-hctm&roles=admin' \
  | jq -r .token)"
test -n "$IAOS_TOKEN" && test "$IAOS_TOKEN" != null
```

注意：dev token endpoint 只适用于本地开发，演示或生产环境必须使用真实登录/服务身份。JWT 的 tenant 必须是 `tenant-hctm`，不能用请求头把其他租户 token 伪装成 HCTM 租户。

## 4. Dry-run 和 Apply

先执行 dry-run；命令不得产生 POST/PUT：

```bash
go run ./cmd/aese apply ./scenario-packs/hctm \
  --story order-expedite-01 \
  --run-id hctm-m3-preview-001 \
  --target http://127.0.0.1:8082 \
  --tenant tenant-hctm
```

确认 impact report、IAOS schema、权限、自然键和 tenant 后，才可显式写入：

```bash
go run ./cmd/aese apply ./scenario-packs/hctm \
  --story order-expedite-01 \
  --run-id hctm-m3-apply-001 \
  --target http://127.0.0.1:8082 \
  --tenant tenant-hctm \
  --apply
```

AESE 会把 HCTM 稳定业务编码投影为 DES-047 wire；IAOS 在单一 tenant 事务中解析 UUID，并对 customer/product/BOM/inventory/order 执行 insert/update/no-op。输出中的 mapping warnings 明示未进入 legacy tracer 的字段。

## 5. Replay 和 Verify

只读解析订单并生成 replay impact：

```bash
go run ./cmd/aese replay ./scenario-packs/hctm \
  --story order-expedite-01 \
  --order-id <scenario-apply 返回的 sales_order object_id> \
  --target http://127.0.0.1:8082
```

显式触发受治理订单分解：

```bash
go run ./cmd/aese replay ./scenario-packs/hctm \
  --story order-expedite-01 \
  --order-id <scenario-apply 返回的 sales_order object_id> \
  --target http://127.0.0.1:8082 \
  --apply
```

最小在线断言：

```bash
go run ./cmd/aese verify ./scenario-packs/hctm \
  --story order-expedite-01 \
  --target http://127.0.0.1:8082
```

必须保存 run ID、HTTP 结果、`correlation_id`、Outbox event ID/subject、O2D 日志、库存和 work order 查询结果。第二次 replay 必须返回 no-op/明确拒绝，不能重复扣库或创建工单。

## 6. Reset

安全查看 reset 计划：

```bash
go run ./cmd/aese reset ./scenario-packs/hctm \
  --story order-expedite-01 \
  --run-id hctm-reset-preview-001 \
  --target http://127.0.0.1:8082
```

确认 dry-run 只包含 1 个 sales order、5 个 inventory 和对应派生状态后，显式执行：

```bash
go run ./cmd/aese reset ./scenario-packs/hctm \
  --story order-expedite-01 \
  --run-id hctm-reset-apply-001 \
  --target http://127.0.0.1:8082 \
  --apply
```

预期删除 6 个 L2 对象并保留 12 个 L1 customer/product/BOM。不得用直接 SQL 删除或重置租户数据。

## 7. 完成证据

截至 2026-07-19，M3 已完成。实际 run IDs、input hash、correlation/event/workflow IDs、对象计数、第二次 no-op、tenant isolation 和 reset 恢复结果见 `docs/reports/hctm-m3-execution-evidence.md`。
