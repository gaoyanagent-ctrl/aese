# M13 Genesis First Commercial Delivery Evidence

- 正式需求 10,000 + 追加 2,000；零期初可销售库存，Genesis transaction code 不与旧 O2D 冲突。
- 12,360 供应/投入形成 12,000 良品与 360 报废；三批 9,000/2,700/300 全部客户接受，300 件 discrepancy 经治理恢复关闭。
- 净收入 14,400,000、含税发票 16,272,000、抵扣合同负债后 AR/银行核销 14,272,000、AR 归零；实际成本 10,200,000、项目毛利 4,200,000 CNY。
- 13 帧、100 次 hash、snapshot/restore、数量/财务失败关闭和旧 O2D compatibility 通过，终态 `first_commercial_cycle_closed=true`。
- IAOS revision `067bbb4` 提供十类 delivery 权限，强制 World evidence、RLS、幂等、并发和原子 journal/Outbox。
