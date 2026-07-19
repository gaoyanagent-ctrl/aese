# HCTM → IAOS M3 兼容性报告

日期：2026-07-19  
范围：M3 T12–T16，只读调查 `/iaos/iaos-go` 当前 checkout 与本地运行环境。

## 1. 结论

IAOS 现有能力可以支撑一条很窄的 O2D tracer，但当前环境还不能直接执行 HCTM pack。第一批应复用 legacy `customer` / `product` / `bom` / `sales_order` + `sales_order_line` / `inventory`，通过 adapter 把 HCTM 稳定业务编码解析成 IAOS UUID。不应在 M3 为 28 个对象全部创建动态 entity。

当前硬阻塞是：

- `tenant-hctm` 数据库中没有 HCTM metadata schema、O2D workflow config 或 tracer 数据。`GET /metadata/schemas` 返回的唯一 `sales_order` 只是 handler 的空列表 fallback，不是已注册 schema。
- 现有 entity create/import 是 insert-only，没有基于自然键的 dry-run/upsert，不能满足 pack apply 幂等、整包原子性和可恢复性。
- `POST /entities/sales_order/:id/decompose` 可在一个事务中确认订单并写 Outbox，但没有动作权限检查、业务幂等键或重复确认防护。
- O2D coordinator 按 node 分开提交事务；库存扣减、工单创建之间没有 workflow-run 幂等记录。重放可能重复扣库，后续工单唯一键冲突时还会留下部分副作用。
- HCTM `qty_per=1.05` 在 O2D BOM 展开中被 `int(qtyPerUnit) * orderQty` 截断为 `1 * orderQty`，对 12,000 件需求会丢失 600 张铝板的基础用量，还未计入 `scrap_rate=0.02`。该算法必须在 tracer 前修复并用 decimal 合同验证。

因此，T16 结论是：**需要 IAOS 侧新设计，不能仅依赖现有通用 CRUD**。建议先建两份独立 DES：“Scenario Import/Apply Contract”和“Governed Simulation Ingress”。如果确认两者是平台级长期能力，再用一份 IAOS ADR 固定“场景内容由 AESE 版本化，IAOS 拥有导入、权限、事务、Outbox 和幂等运行时”的边界。

## 2. 运行环境证据（T12）

| 组件 | 2026-07-19 实测 | 启动/检查入口 | 判断 |
| --- | --- | --- | --- |
| Platform | `platform_bin` PID 4016907 监听 `8082`；`GET /health` 返回 `200 {"status":"UP"}`；`GET /ready` 返回 DB/EventBus 均 `OK` | `DATABASE_URL=postgres://iaos_app:***@localhost:5433/iaos?sslmode=disable ./scripts/deploy_platform.sh` | 运行中，二进制时间 2026-07-19 03:24 |
| PostgreSQL | `iaos-integration-postgres` healthy，host `5433`；以 `iaos_app` 成功执行只读查询 | `docker ps`；container 内 `psql -U iaos_app -d iaos` | 运行中 |
| NATS JetStream | `iaos-nats` 运行，`4222` 监听；Platform readiness 的 EventBus 为 `OK`；O2D 日志记录 NATS 连接成功 | compose 中 `nats:2.10-alpine` + `-js` | 客户端端口可用；容器虽映射 `8222`，启动命令未启用 `-m 8222`，因此 monitoring API 不可用 |
| O2D | `o2d-bin` PID 4114193 运行，日志显示已订阅 `iaos.*.o2d.order.confirmed` | `DATABASE_URL=postgres://postgres:***@127.0.0.1:5433/iaos?sslmode=disable ./scripts/deploy_o2d.sh` | 进程活着，但二进制是 2026-06-18 构建；启动自测因种子订单不存在而失败，无端到端成功证据 |

运行环境里 `tenant-hctm` 的事实：

- `entity_metadata_schema` 0 行。
- `workflow_config` 0 行。
- `sys_outbox` 0 行。
- `customer` / `product` / `bom` / `sales_order` / `sales_order_line` / `inventory` / `work_order` 各 0 行。
- 全库 `workflow_config` 也是 0 行；`deploy/postgres/seed.sql` 虽定义 `wf-o2d-standard`，当前 integration DB 没有导入它。

## 3. IAOS metadata/API 合同（T13）

### 3.1 已验证 API

| API | 认证/权限 | 行为 | HCTM 用法/限制 |
| --- | --- | --- | --- |
| `GET /api/v1/metadata/schemas` | JWT + tenant/RLS；无对象级 read 检查 | 列 schema | 可做 compatibility discovery；空库会伪造 `sales_order` fallback，adapter 必须再 GET 具体 schema |
| `GET /api/v1/metadata/schema/:entity` | JWT + tenant/RLS | 返回 fields/indexes/events/physical table，另计算 create/update/delete 权限 | dry-run 的主要读取入口 |
| `POST /api/v1/metadata/register` | 当前仅 JWT + tenant/RLS，handler 内无显式 RBAC | 编译 schema、DDL、RLS 和权限资源 | 不应由 AESE apply 默认调用；需受治理配置包/管理员步骤 |
| `POST /api/v1/entities/:entity` | JWT + `btn.<entity>.create:EXECUTE` + RLS | 单行 insert；无 upsert/dry-run | 可作小型 tracer adapter，但 adapter 需先查自然键并处理冲突 |
| `GET /api/v1/entities/:entity/records` / `GET .../:id` | JWT + RLS；字段/ABAC 过滤，无对象级 read 检查 | 分页/单记录读取 | 可用于 dry-run/verify，但需修补 read 权限合同 |
| `PUT` / `DELETE /api/v1/entities/:entity/:id` | JWT + `btn.<entity>.update/delete:EXECUTE` + RLS | 按 UUID 更新/删除 | 不能直接用业务编码；adapter 需先 resolve |
| `POST /api/v1/entities/:entity/import` | JWT + RLS，handler 内无 create/import RBAC | CSV/XLSX，每个 parent 单独事务 | 不接受 JSON pack，无全包原子性、upsert 或 dry-run，不适合正式 scenario apply |
| `POST /api/v1/entities/sales_order/:id/decompose` | JWT + RLS，handler 内无动作 RBAC | 更新 `confirmed` + 同事务写 `sys_outbox` | 可作 OrderConfirmed 正式入口，但必须先增加权限、幂等和 correlation 合同 |

metadata 可表达 `required`、`unique`、composite `indexes`、`reference/relation`、enum、decimal、date/datetime、字段 layer 和影子列。但 `unique` 只是 schema/DDL 属性，现有 CRUD 没有“按自然键 upsert”语义；AESE adapter 不能把 metadata 能表达唯一键误解为 apply 已幂等。

### 3.2 数据与 RLS

legacy tracer 表均开启 RLS：`customer`、`product`、`bom`、`sales_order`、`sales_order_line`、`inventory`、`work_order`、`workflow_config`、`sys_outbox`。当前 policy 使用 `tenant_id = current_setting('app.current_tenant_id', true)`，但上述表的 `FORCE ROW LEVEL SECURITY` 均为 false。应用路径依赖 `store.Transact` 注入 tenant；AESE 正式路径不得直连数据库。

## 4. 28 个 HCTM 对象分类（T14）

`C` = 可直接复用 legacy entity 的核心字段；`M` = 需 adapter/字段裁剪；`U` = 当前无运行实体或不进入 M3 tracer。

| HCTM 对象 | 等级 | IAOS 对象/API | 字段和缺口 | M3 策略 |
| --- | --- | --- | --- | --- |
| `enterprise_group` | U | 无 HCTM entity | 无 `group_code/name/industry` | pack 保留，不 apply |
| `business_unit` | U | 可候选 IAOS org node | HCTM `bu_code/region` 未映射 | 不进 tracer |
| `legal_entity` | U | 无 tracer entity | 订单 legacy 不存法人外键 | pack 保留，订单丢失此维度 |
| `plant` | U | 可候选 org/facility | O2D 表无 plant | 固定 `PLT-SZ` 作 run summary 上下文，不入 legacy 表 |
| `department` | U | IAOS org node 候选 | 缺业务编码映射 | 不进 tracer |
| `production_team` | U | 无 | 无 line/shift/team 合同 | 不进 tracer |
| `customer` | C/M | `customer` + entity CRUD | `customer_code→code`，`customer_name→name`，`customer_type→oem_category`；priority/payment/delivery/ship-to 缺失 | 导入 1 家，缺失字段仅作 extension 或不导入 |
| `supplier` | U | 无已验证 entity | 全部 HCTM 字段无对应 | 不进首个 OrderConfirmed tracer |
| `material` | C/M | `product` + entity CRUD | `material_code→code`，`material_name→name`；可用 `specification/bom_version`；type/uom/lot/plant/safety stock 缺失 | 成品与 5 个组件统一映射为 product，adapter 维护 code→UUID |
| `bom` | C/M | `bom` + entity CRUD | parent/component code 需 resolve 到 UUID；`qty_per→quantity_required`；uom/scrap/effective/version 缺失；无自然唯一键 | 导入 5 行；修复 decimal 展开后才能验收 |
| `routing` | U | 无 O2D runtime 映射 | 无产能/版本合同 | pack 保留，不 apply |
| `operation` | U | 无 O2D runtime 映射 | 无工序顺序/节拍/良率 | 不进 tracer |
| `work_center` | U | 无已注册 schema | 无班产能/瓶颈标识 | 不进 tracer |
| `equipment` | U | 物理表候选，当前 metadata API 404 | 无可用 HCTM schema | 不进 tracer |
| `tooling` | U | 无 | 全部缺失 | 不进 tracer |
| `warehouse` | M | `inventory.warehouse_name` | 只有自由文本，无 warehouse 主数据/类型/工厂 | 仅把 `warehouse_code` 写入 `warehouse_name` |
| `storage_location` | U | legacy inventory 无 location | 库位和限制缺失 | 不进 tracer |
| `shift` | U | 无 O2D 映射 | 无时间/产能模型 | 不进 tracer |
| `employee` | U | IAOS user/org 候选 | 业务员工编码与账号未映射 | 不进 tracer |
| `sales_order` | M | `sales_order` + `sales_order_line` | header: `order_no`、`customer_code→customer_id`、`due_date→required_date`、status；line: `material_code→product_id`、`order_qty→quantity`；无 legal entity/priority/original order，`unit_price` 必填但 HCTM 未给 | 只导入一张 12,000 件合并 tracer 订单；价格使用明示的演示值，不调用 legacy 特殊 create 的默认交期/价格逻辑 |
| `purchase_order` | U→M4 | M4 scenario fixture + simulation ingress | M3 无完整 header；M4 已固定 supplier/material/quantity/promised date/latest ETA/status 的窄合同 | M3 首个 tracer 未导入；M4 已投影两张采购单，供供应商延期稳定解析 |
| `goods_receipt` | U/M | 只有 `purchase_receipt_line` | 无 receipt header，line 字段仅 product/qty/batch/warehouse | 不进首个 tracer |
| `inspection_order` | M4 | M4 scenario fixture + simulation ingress | pending 预分配允许 receipt/lot 为空；固定 inspection/PO/material/sample/accepted/rejected/status 窄合同 | M3 未导入；M4 已投影 `IQC-202607-0002`，供来料检验失败稳定解析 |
| `inventory_transaction` | M | `inventory` 当前余额 + `inventory_transaction` 流水 | 余额表可对应 product/warehouse/qty/batch；流水表只有 product/warehouse/qty/direction/reference/batch；无 plant/location/lot/source type/status | tracer 初始库存导入 `inventory`；不把 opening balance 假装成完整库存事务 |
| `production_order` | M/U | O2D `work_order` | 物理字段是 `wo_no/sales_order_id/product_id/quantity/planned_*`；当前 metadata schema 却主要暴露 `document_no/document_date/...`，API 合同与 O2D SQL 不一致 | 只由 O2D 产生，AESE 不预先 apply；verify 需新的受治理查询或修正 schema |
| `operation_task` | U | 无 | 无现场任务合同 | 不进 tracer |
| `shipment` | U | 无 | 无发货对象/事件 handler | 不进 tracer |
| `quality_issue` | U | 无 HCTM entity | 无缺陷、处置、owner 合同 | 不进 tracer |

## 5. 第一批 legacy 映射决策（T15）

### 5.1 `sales_order`

HCTM 的一行订单必须拆成 legacy header + line：

| HCTM | IAOS | 转换 |
| --- | --- | --- |
| `order_no` | `sales_order.order_no` | 原值，tenant 内唯一 |
| `customer_code` | `sales_order.customer_id` | 先按 `(tenant_id, customer.code)` resolve UUID |
| `due_date` | `sales_order.required_date` | `Asia/Shanghai` 当地日期转 RFC 3339 timestamp，不使用 `NOW()+30 days` |
| `status` | `sales_order.status` | tracer apply 初始为 `draft`，由 decompose 转 `confirmed` |
| `material_code` | `sales_order_line.product_id` | 先按 `(tenant_id, product.code)` resolve UUID |
| `order_qty` | `sales_order_line.quantity` | 必须是正整数；12,000 可兼容 |
| 无对应 | `sales_order_line.unit_price` | pack/adapter 必须显示给出演示价格，不接受 handler 硬编码 `128.5000` 作业务真值 |

`legal_entity_code`、`priority`、`original_order_no`、`confirmed_at` 暂不进 legacy O2D 运算。它们可在后续受控 metadata 扩展中保留，但不能影响 tracer 的核心 SQL 列。

### 5.2 `material` / `product`

M3 将 finished good、raw material、purchased part 和 packaging 统一映射为 `product`。自然键是 `(tenant_id, code)`。只把 `material_code/name` 作核心必需字段；`uom`、lot/serial control、safety stock 暂不参与 O2D，但必须在 compatibility 输出中标记丢失，不可默默丢弃。

### 5.3 BOM

`parent_material_code` / `component_material_code` 分别 resolve 为 product UUID，`qty_per` 写入 `quantity_required NUMERIC(12,4)`。M3 不把 `bom_code`、`uom`、`scrap_rate`、`effective_from/to` 假装成 legacy 字段。

必须补的最小 IAOS 修正：

1. O2D BOM 展开全程使用确定 decimal，不得先 `int(qtyPerUnit)`。
2. 对 `(tenant_id, parent_product_id, child_product_id)` 给出明确的唯一/版本策略，否则重复 apply 会重复展开用量。
3. 在报告中明示 `scrap_rate` 未被 legacy MRP 计入；若 tracer 验收要求 12,600 张，则必须扩展 BOM 合同或把演示期望拆分为基础用量与损耗。

### 5.4 Inventory

HCTM opening balance 在 tracer 中投影为 legacy `inventory(product_id, warehouse_name, quantity, batch_no)`：`material_code` resolve UUID，`warehouse_code→warehouse_name`，`qty→quantity`，`lot_no→batch_no`。

当前 `inventory` 无自然唯一键，且 O2D reserve 直接减少 quantity。所以 adapter 不能用“找不到就 insert”实现可靠重置。至少需要 `(tenant_id, product_id, warehouse_name, batch_no)` 的受治理 upsert/reset 语义，并保证 reset 不删 L1 主数据。

## 6. O2D 事件、权限和幂等缺口

### 6.1 现有正向链路

`POST .../sales_order/:id/decompose` 在订单事务中生成 `eventdef.Event`，并把完整 event JSON 写入 `sys_outbox`。Outbox poller 发布到 `iaos.{tenant_id}.o2d.order.confirmed`，O2D durable group `o2d-mrp-group` 订阅 `iaos.*.o2d.order.confirmed`。这一边界符合“持久化业务变化通过 IAOS 事务 + Outbox”原则。

### 6.2 不符合 HCTM 合同的地方

- decompose event 未写 `Metadata`，因而没有 HCTM 必需 `correlation_id`、`causation_id`、`idempotency_key`、scenario 和 business object metadata。HTTP middleware 的 correlation ID 也没有被传入 event。
- Event ID 使用 `time.Now().UnixNano()`，重复调用必然生成新事件。`sys_outbox` 没有对业务幂等键的唯一约束。
- decompose 对已 `confirmed` 订单仍执行 UPDATE 并写新 Outbox，无 compare-and-set 状态机。
- decompose 没有 `capability.*`、`btn.*` 或独立 order-confirm 权限检查，任何带 tenant JWT 的用户都可进入 handler。
- O2D `LoadAndRun` 没有 workflow execution/run 表或输入 event 去重。每个 node 独立事务，不是整个 DAG 原子。
- O2D 缺少 HCTM 18 类事件中的大部分 constant/handler；M3 只能验证 `o2d.order.confirmed` 及现有 inventory/work-order 副作，不能宣称完成 22 步故事重放。

## 7. 需要的 IAOS 设计决定（T16）

### 7.1 DES-A：Scenario Import / Apply Contract（M3 前置）

建议在 IAOS 独立 worktree 中设计，最小合同包括：

- 显式 target tenant/environment，服务端默认 dry-run，apply 需独立权限。
- `pack_key + pack_version + scenario_key + run_id`、actor、correlation ID 和结果摘要的审计记录。
- 稳定业务编码到 UUID 的服务端 resolve，不让 AESE pack 携带环境 UUID。
- 按对象定义自然键的 upsert/no-op/conflict 结果，不使用通用 insert-only CRUD 暗示幂等。
- 每个可观测 apply unit 的事务边界、失败时无部分写入证据，以及 reset 不删 L1 的安全边界。
- 一个面向 JSON pack 的 API/Capability，而不是将 CSV/XLSX import 当作场景合同。

现有 API 可以作 DES 的 adapter 内部实现依据，但不足以直接对 AESE 公开正式 apply。

### 7.2 DES-B：Governed Simulation Ingress（M4 前置，M3 先固定合同）

用于供应商延期、设备故障、来料不良等外生事件，最小要求：

- 只接受允许的 event type/schema，校验 tenant、actor、source、business object reference。
- 必填 `correlation_id` 和 `idempotency_key`，数据库唯一约束保证重试 no-op/返回原结果。
- 通过受治理 Capability/API 进入事务、审计和 Outbox，不将 direct NATS 作演示验收路径。
- 对“仅通知”和“同时改变设备/PO/检验状态”的事件分别定义事务语义。

### 7.3 O2D tracer 还需独立实现修正

这些是代码修正，不必各自建 ADR，但应在 IAOS DES/plan 中成为验收项：

1. decimal BOM 展开与 HCTM 数量回归测试。
2. `tenant-hctm` 的 `wf-o2d-standard` 受治理 seed/config apply。
3. order confirm 动作权限、状态 compare-and-set、correlation/idempotency 写入 Outbox。
4. workflow execution/event 去重，以及失败后不重复扣库。
5. `work_order` metadata schema 与 O2D 物理 SQL 字段收敛，使 API verify 能读到 `wo_no/product_id/quantity/sales_order_id`。

## 8. 建议的 S4/S5 执行顺序

1. 先审批 IAOS Scenario Import/Apply DES，固定最小 compatible object 和自然键。
2. 在 IAOS 独立 worktree 修复 BOM decimal、order-confirm 权限/幂等、workflow 去重与 work-order schema。
3. 为 `tenant-hctm` 导入真实 metadata/resource/workflow config，不依赖 list API fallback。
4. AESE adapter dry-run 先输出 code→UUID resolve 计划、insert/update/no-op/conflict 计数，确认零写入。
5. 显式 `--apply` 导入 customer/product/BOM/inventory/order，立即重复执行验证 no-op。
6. 调用受治理 order-confirm，保存 HTTP response、Outbox payload/subject、O2D log、inventory/work-order 结果和第二次调用结果。

## 9. 取证入口

本报告主要依据以下源码和实测：

- `/iaos/iaos-go/platform/internal/api/router.go`
- `/iaos/iaos-go/platform/internal/api/router_entity_schema.go`
- `/iaos/iaos-go/platform/internal/api/router_entity_register.go`
- `/iaos/iaos-go/platform/internal/api/router_entity_crud.go`
- `/iaos/iaos-go/platform/internal/api/import_entity.go`
- `/iaos/iaos-go/platform/internal/metadata/engine_schema.go`
- `/iaos/iaos-go/shared/metadata/types.go`
- `/iaos/iaos-go/deploy/postgres/init.sql`
- `/iaos/iaos-go/deploy/postgres/seed.sql`
- `/iaos/iaos-go/scenarios/o2d/cmd/o2d/main.go`
- `/iaos/iaos-go/scenarios/o2d/internal/mrp/`
- `/iaos/iaos-go/platform/pkg/workflow/coordinator.go`
- `/iaos/iaos-go/shared/eventdef/events.go`
- `curl /health`、`curl /ready`、authenticated metadata GET、`docker ps`、`ss`、容器内只读 `psql`、Platform/O2D log tail。

本次没有调用任何写 API，没有直接修改 IAOS 数据库，也没有修改 `/iaos/iaos-go` 工作树。

## 10. 报告后的最小修正（2026-07-19）

兼容性取证完成后，IAOS 独立 worktree 实现并测试了首批 tracer 前置，随后 fast-forward 合并到 IAOS `main`，commit 为 `0260f28`：

- BOM 展开改为精确 decimal，并对离散需求向上取整；回归测试证明 `1.05 × 12,000 = 12,600`。
- sales order decompose 使用状态 CAS；重复确认返回 `already_confirmed` 且不新增 Outbox。
- Outbox event metadata 写入 correlation、request、idempotency 和业务对象信息。
- 新增 IAOS DES-047 Scenario Import/Apply Contract 和 DES-048 Governed Simulation Ingress。

这些修正当时只关闭了 decimal 和订单重复确认的局部缺口。随后 M3 继续实现并实证了正式 scenario apply/reset API、稳定编码到 UUID 映射、workflow/event 去重与跨节点原子性、work_order metadata/API 对齐；最终状态以 `docs/reports/hctm-m3-execution-evidence.md` 为准。本报告前九节保留为实现前取证快照。
