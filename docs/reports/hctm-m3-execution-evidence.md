# HCTM M3 执行证据

日期：2026-07-19  
Pack：`hctm@0.1.0`  
目标租户：`tenant-hctm`

## 1. 离线合同

- `aese validate`：`valid: hctm@0.1.0 (19 record sets, 1 stories)`。
- `aese inspect`：当前 pack 为 80 条 L1 主数据、15 条故事初始记录、22 个事件、17 条离线断言；其中 1 条 pending 检验单是 M4 为受治理异常 replay 补充的前置 fixture。
- 4 个 Draft 2020-12 Schema 对 9 个实例文件校验通过。
- Go 单测覆盖路径逃逸、schema 版本、重复自然键、缺失引用、时间乱序、重复幂等键、负数量、BOM/MRP、库存/发运和 expected outcomes。

## 2. IAOS 实现基线

M3 使用并部署了以下 IAOS `main` 提交：

- `0260f28`：decimal BOM、订单确认 CAS/no-op、trace metadata、DES-047/DES-048。
- `127a3c0`：受治理 scenario apply/reset endpoint。
- `d96636f`：O2D workflow 原子事务、事件幂等、work_order metadata 和 HCTM workflow fixture。
- `afecb4f`、`0d2d0e3`：真实 dry-run 发现并修复 missing-row 计划错误。
- `19483ff`、`4d1aa60`：reset 清理派生工单和当前 correlation 的 workflow run。
- `aeaeccb`：所有 scenario legacy 查询显式绑定 tenant，关闭非 FORCE RLS 表的跨租户命中。

Platform 与 O2D 均从 IAOS 主 checkout 重新构建部署；`/health` 为 UP，`/ready` 的 database/event_bus 均为 OK。

## 3. Dry-run、Apply 和幂等

Run `hctm-m3-20260719-001`：

- dry-run：18 insert、0 update/no-op/conflict/unsupported，`committed=false`。
- input hash：`464da3331578afc1ee2c66525cccc09c910c48657979f88379620b4f2712bf47`。
- dry-run 后只读数据库计数：customer/product/BOM/inventory/sales order/scenario run 全部为 0。
- 显式 apply：18 insert，`committed=true`。

Run `hctm-m3-20260719-002`：

- 第二次显式 apply：0 insert/update、18 no-op、0 conflict。
- 对象数保持 customer=1、product=6、BOM=5、inventory=5、sales order=1、line=1。

Reset 后 run `hctm-m3-20260719-004` 再次恢复故事状态：6 insert、12 L1 no-op，证明 reset 后可重复 apply。

## 4. 租户隔离

使用 `tenant-other` JWT 对相同 pack 执行 dry-run：

- 结果为 18 insert、0 no-op/conflict，说明未解析到 `tenant-hctm` 自然键。
- `tenant-other` 的 customer/product/order/scenario run 数据库计数均为 0。

该检查曾发现 legacy 表未 FORCE RLS 时的跨租户命中，已由 IAOS commit `aeaeccb` 通过所有 adapter SQL 显式 `tenant_id` 条件修复并回归验证。

## 5. O2D tracer

最终可演示运行：

- AESE correlation：`corr-so-202607-0001`。
- IAOS order-confirmed event：`evt-conf-d2f7c859b9e7d9fd10a7bd1a`。
- Workflow run：`af706c43-b080-42de-8c98-b421d1b9e815`。
- Workflow：`wf-o2d-hctm`，状态 `completed`。
- BOM 证据：铝板 `1.05 × 12,000 = 12,600`，没有整数截断。
- 库存节点识别 3 个短缺，预留包装箱 1,200、密封圈 24,000。
- 生成 3 个 pending work orders，数量分别为 4,600、6,000、12,000。
- `aese verify`：work order exists 与 count=3 两条在线断言均通过。

重复 replay 返回：

```text
action=already_confirmed
triggered=0
applied=false
```

重复前后 `o2d.order.confirmed` Outbox 数量保持 3，不新增事件；work order 数保持 3。

## 6. Reset 与恢复

Run `hctm-reset-20260719-002`：

- dry-run：计划删除 6 个 L2 对象，保留 12 个 L1 对象，零写入。
- 显式 reset：`committed=true`，删除 1 个 sales order 和 5 个 inventory opening balance。
- 同事务清理该订单派生的 work orders 和当前 correlation workflow run。
- reset 后 customer=1、product=6、BOM=5；inventory/order/work order=0。
- 随后 apply/replay/verify 再次成功，环境最终停在可演示状态。

首次端到端试跑在 reset 派生清理修复前留下 1 条历史 workflow run 和对应 Outbox 审计记录；它们不参与最终自然键、订单确认或工单计数，也不影响重复运行。该历史证据保留用于说明修复过程，没有用直接 SQL 篡改。

## 7. 验收结论

T1-T30 已完成。dry-run、apply、tenant isolation、第二次 apply no-op、订单确认、Outbox/NATS/O2D、重复 replay no-op、verify、reset 和 reset 后恢复均有实际证据。M3 可以标记 completed。
