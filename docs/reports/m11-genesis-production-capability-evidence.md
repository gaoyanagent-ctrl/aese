# M11 Genesis Production Capability Evidence

- 资金：M10 期初现金 10,000,000，剩余认缴资本 10,000,000 经实际到账链转为现金；保留设施尾款、3,000,000 工资准备金和最低现金缓冲。
- 采购：采购、融资租赁、供应商账期三个组合先过现金/工期硬约束，选择唯一可行租赁组合。
- 设备：成形、CNC、激光焊接、清洗、检漏、装配和质量实验室七项能力完成到货、落位、安装、调试、校准、安全与验收。
- 人员：10 名虚构核心人员实际到岗并完成安全与岗位实操资格，一班制和关键替补通过；候选详情保持 actor/private scope。
- 异常：检漏校准漂移形成 World/IAOS discrepancy，设备与质量负责人仅在 observation 后获知，经治理整改与复验关闭。
- 确定性：10 帧、100 次 hash、snapshot/restore 和联合 gate 失败关闭通过；终态 `industrialization_eligible=true`。
- IAOS：独立 revision `789b925` 提供七类权限 allowlist，World evidence、FORCE RLS、幂等、并发和 business record/journal/Outbox 原子事务。
