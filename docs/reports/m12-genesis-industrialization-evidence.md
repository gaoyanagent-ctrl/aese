# M12 Genesis Industrialization Evidence

- 单一虚构客户 `CUST-SGNEV`、产品 `HCTM-BCP-A01`、RFQ/定点和 2,000,000 CNY 开发预付款合同负债。
- product/BOM/routing/PFMEA/control-plan revision/hash 一致；revision A 试制后只能由受治理 revision B 取代。
- 两轮试制：T1 100 件、良率 82%、泄漏 7、Cpk 0.91；T2 120 件、良率 96.67%、泄漏 0、Cpk 1.67，全部可追溯且不可销售。
- 焊接泄漏/Cpk discrepancy 经 delayed Knowledge、containment、engineering change、复试关闭；客户 PPAP 实际批准。
- 11 帧、100 次 hash、snapshot/restore、PPAP 失败关闭和旧 HCTM stable code compatibility 通过，终态 `serial_production_eligible=true`。
- IAOS 独立 revision `50a46e2` 提供 quote/design/process/trial/quality/ppap/release allowlist，强制 evidence、RLS、幂等、并发与原子 journal/Outbox。
