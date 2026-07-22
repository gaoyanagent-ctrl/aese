# M12 Genesis Industrialization Runbook

访问 `/#world-capability-build`，进入“产品工业化 Campaign”，单步 10 次。终态必须显示 `M13 eligible`、六个 APQP gate、PPAP approved、revision B release hashes、两轮试制和旧 HCTM stable-code compatible。

运行全仓 Go test/vet、`aese world validate world-packs/hctm-genesis`、前端 test/typecheck/build/Playwright。IAOS `/api/v1/genesis/industrialization/actions` 的 quality/ppap/release 必须引用 World/customer evidence；受权限、职责隔离、RLS、幂等、expected version 和事务 journal/Outbox 约束。

试制件不可销售；RFQ 不是订单，客户预付款是合同负债而非收入。M12 不产生正式库存、发运、发票、应收或回款。
