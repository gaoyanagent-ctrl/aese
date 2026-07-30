---
id: DES-030
title: 制造企业财务运行体系总览
date: 2026-07-28
status: approved
author: Codex + User
tags: [finance, accounting, manufacturing, index]
---

# 制造企业财务运行体系总览

## 1. 目的

本文件只定义财务子域的总体边界、模块导航和分阶段完成门。详细字段、流程、控制和验收
分别放在本目录的模块设计中，M9–M13 计划只引用相关文件，不再复制财务设计正文。

IAOS 财务目标是把资本、采购、库存、生产、销售、资产、薪酬、资金等经营事实编译为
可审计的会计事件、凭证、总账、子账、成本、报表和 KPI。AESE 只提供外部与物理世界事实，
不建立第二套账。

## 2. 模块导航

| 模块 | 权威设计 | 主要范围 |
| --- | --- | --- |
| 多组织与共享主数据 | [DES-031](DES-031-multi-organization-and-shared-master-data-foundation.md) | 集团、法人、BU、共享中心、多账簿、Data Set、BP/产品组织扩展 |
| 会计内核与总账 | [DES-032](DES-032-accounting-kernel-and-general-ledger.md) | 会计事件、凭证主子表、科目、汇率、余额、过账、冲销 |
| 子账、应收应付与资金 | [DES-033](DES-033-subledgers-receivables-payables-and-treasury.md) | AP、AR、open item、付款、收款、银行、资金预测与对账 |
| 制造成本、存货与资产 | [DES-034](DES-034-manufacturing-cost-inventory-and-assets.md) | 标准/实际成本、WIP、差异、存货、CIP、固定资产和折旧 |
| 预算、关账、报表与 KPI | [DES-035](DES-035-budget-close-reporting-and-kpi.md) | 预算、预测、模块期间、月结、三表、管理报表和驾驶舱 |
| 财务治理、审批与 Agent | [DES-036](DES-036-finance-governance-approval-and-agents.md) | 财务组织、SoD、权限、审批、Capability、Agent Tool、菜单和解释入口 |

目录索引和阅读顺序见 [Finance Design Index](README.md)。

## 3. 规范与系统边界

首个 tracer 以中国境内、人民币本位币、权责发生制和借贷记账法为基线，但单法人不是
长期架构。会计准则、科目模板、税率、报表模板和记账规则必须版本化并由有权主体批准。

```text
AESE World 现实事实
  银行到账/支付、客户接受、实物收发、设备投用、人员服务
                    ↓ Observation
IAOS 业务子账
  资本、采购、库存、生产、销售、资产、薪酬、资金
                    ↓ Accounting Event
IAOS Accounting Kernel
  会计政策 → 记账规则 → 凭证 → 总账/子账 → 对账/关账 → 报表/KPI
```

关键边界：

- 预算不是凭证，承诺不是应付，发票不是现金，认缴不是实缴；
- 业务记录不能由会计凭证反向伪造；
- 已过账凭证不可编辑或删除，只能冲销或调整；
- 财务 Entity 必须同时有模型入口和运行数据入口；
- `finance_*` 权威表与 `entity_projection_*` 投影必须同事务或受治理幂等迁移；
- 不得按 M9/M10 等里程碑复制平行财务表。

## 4. Project Genesis 分段接入

| 里程碑 | 财务后果 | 对应设计 |
| --- | --- | --- |
| M9 | 财务组织、账簿、科目、期间、实缴资本凭证、开业试算平衡 | DES-031、032、036 |
| M10 | 工程合同、暂估、供应商发票、应付、付款、在建工程 | DES-033、034 |
| M11 | 设备转固、折旧、薪酬费用、能力建设成本 | DES-034 |
| M12 | 项目成本、工装、样件、客户付款和合同负债 | DES-033、034 |
| M13 | P2P、库存/WIP/FG、生产成本、收入、应收、回款、销售成本 | DES-033、034 |
| 持续运营 | 预算、月结、三表、管理报表、KPI、合并与恢复 | DES-035、031 |

业务章节产生财务后果，但正常自动记账不要求玩家逐张制证。异常、重大金额、手工调整、
政策选择、关账与重开才进入人工或受治理 Agent 决策。

## 5. 当前 M9 实现边界

当前实现包含：

- 六个财务责任岗位、Mandate 与四条阻断型职责分离规则；
- 财务组织、CAS-BE/CNY 账套、12 个期间、1002/4001 科目；
- 实缴资本借银行存款、贷实收资本的已过账双分录；
- 银行日记账、总账、试算平衡和开业资产负债表读取；
- 五个财务工作项和 `finance_opening_ready` 硬门；
- AESE 财务中心到 IAOS 账务/报表的穿透入口。

当前未完成：

- DES-031 的多组织、Data Set、BP/产品组织扩展和账套原型迁移；
- AP、AR、资金、资产、库存和成本完整子账；
- 模块期间、月结、完整三表、合并、预算和经营驾驶舱。

## 6. 总体完成门

- 每个财务事实可追溯到业务事件、规则版本、Capability Execution 和证据；
- 每张凭证借贷平衡，已过账记录不可变；
- AP、AR、资金、资产、库存、成本与总账可对账；
- 每个余额明确 tenant、法人、账簿、期间、币种和财务维度；
- 共享主数据不共享法定余额，跨法人业务不混账；
- 月结和报表失败关闭，不用演示数据静默配平；
- 每项金额可穿透至凭证、会计事件、业务单据和 World evidence；
- 用户不需要源码、数据库或原始 JSON 才能配置和理解财务功能。

## 7. 规范参考

- [财政部企业会计准则专题](https://kjs.mof.gov.cn/zt/kjzzss/kuaijizhunzeshishi/)
- [财政部企业会计准则应用指南](https://www.mof.gov.cn/zhengwuxinxi/caizhengxinwen/200805/t20080519_24574.htm)
- [COSO Internal Control—Integrated Framework](https://www.coso.org/internal-control)
- [SAP Chart of Accounts](https://help.sap.com/docs/SAP_S4HANA_ON-PREMISE/8fbeed5f2046489696a50ac7fd76f9c6/9f53c2531bb9b44ce10000000a174cb4.html)
- [Oracle Enterprise Structures and General Ledger](https://docs.oracle.com/en/cloud/saas/financials/24b/faigl/implementing-enterprise-structures-and-general-ledger.pdf)
