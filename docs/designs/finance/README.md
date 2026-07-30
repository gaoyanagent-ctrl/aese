# Finance Design Index

财务是 IAOS 业务运行时的独立子域。本目录按模块保存设计，避免 M9/M10 等里程碑文档
重复或吞入完整财务规格。

## 阅读顺序

1. [DES-030 总览](DES-030-manufacturing-finance-operating-system.md)
2. [DES-031 多组织与共享主数据](DES-031-multi-organization-and-shared-master-data-foundation.md)
3. [DES-032 会计内核与总账](DES-032-accounting-kernel-and-general-ledger.md)
4. [DES-033 子账、应收应付与资金](DES-033-subledgers-receivables-payables-and-treasury.md)
5. [DES-034 制造成本、存货与资产](DES-034-manufacturing-cost-inventory-and-assets.md)
6. [DES-035 预算、关账、报表与 KPI](DES-035-budget-close-reporting-and-kpi.md)
7. [DES-036 财务治理、审批与 Agent](DES-036-finance-governance-approval-and-agents.md)

## 文档边界

- 里程碑文档描述何时发生业务事实，只引用财务模块文件。
- 本目录描述财务语义、Entity、Capability、Process、Policy、审批、Agent 和运行控制。
- IAOS 仓库的 DES-063/064/065 描述规范投影和受治理写入边界；DES-067–069
  描述多组织、账簿和伙伴/产品基础；DES-070 描述财务导航拆分及这些权威模型向
  Semantic/Entity 资产的发布合同。
- 完整状态以 `docs/roadmap.md` 和财务实施计划为准；设计完成不等于运行完成。
