---
id: DES-030
title: M9–M13 制造企业财务运行体系
date: 2026-07-28
status: approved
author: Codex + User
tags: [m9, finance, accounting, costing, treasury, manufacturing]
---

# M9–M13 制造企业财务运行体系

## 1. 业务问题与设计目的

当前 Project Genesis 已分别记录资本到账、现金、预算、承诺、应付、开票、应收、回款、
实际成本和项目毛利，但这些对象仍是经营台账，不构成完整会计系统。尤其是 M9 银行注资
只改变资本与现金状态，没有形成实收资本会计凭证、总账余额、银行日记账和开业资产负债表。

本设计把财务能力建设提前到企业成立阶段：

1. M9 建立财务组织、账套、科目、会计政策、期间、审批权限和期初账；
2. 实收资本到账形成不可变、可追溯、借贷平衡的正式会计记录；
3. M10–M13 的工程、采购、资产、生产、销售事件逐步启用对应子账与自动记账规则；
4. 业务单据、会计事件、凭证、总账、管理成本和报表彼此关联，但不互相冒充；
5. 玩家、财务人员和 Agent 在游戏与 IAOS 中共同完成建账、审核、记账、对账和结账。

目标是形成适合离散制造企业的业财一体化财务运行底座，而不是在 M9 虚构尚未发生的采购、
生产或销售交易。

## 2. 规范基线与适用边界

首个参考实现以中国境内、单法人、人民币本位币、权责发生制和借贷记账法为基线：

- [财政部《企业会计准则——基本准则》](https://xj.mof.gov.cn/caizhengjiancha/200805/t20080524_40447.htm)；
- [财政部企业会计准则专题](https://kjs.mof.gov.cn/zt/kjzzss/kuaijizhunzeshishi/)；
- [财政部企业会计准则应用指南](https://www.mof.gov.cn/zhengwuxinxi/caizhengxinwen/200805/t20080519_24574.htm)；
- [管理会计应用指引第 200 号——预算管理](https://kjs.mof.gov.cn/gongzuotongzhi/201612/P020161222582925321535.pdf)；
- [管理会计应用指引第 300 号——成本管理](https://m.mof.gov.cn/tzgg/201612/P020161226556825846042.pdf)；
- [财政部管理会计应用指引体系](https://kjs.mof.gov.cn/zhengcefabu/201710/t20171018_2727363.htm)；
- [COSO Internal Control—Integrated Framework](https://www.coso.org/internal-control)。

首版不宣称替代持证会计师的专业判断，不连接真实税控、银行、征信或监管申报系统，不直接
生成可对外报送的法定报表。会计准则、科目模板、税率和报表模板必须版本化并由有权财务
主体批准后生效。

## 3. 系统边界与权威数据

```text
AESE World 现实事实
  银行到账/付款、客户接受、实物收发、设备投用、人员服务
             ↓ Observation
IAOS 业务子账
  资本、采购、库存、生产、销售、资产、薪酬、资金
             ↓ Accounting Event
IAOS Accounting Kernel
  会计政策 → 记账规则 → 凭证草稿 → 审核/过账 → 总账/明细账
             ↓
对账、关账、财务报表、管理报表、KPI、审计证据
```

- AESE 拥有外部及物理世界事实，不建立第二套总账。
- IAOS 拥有会计政策、凭证、总账、子账、成本计算、关账和报表。
- Agent Knowledge 只通过受治理查询获得必要财务信息，不保存影子账本。
- 预算不是凭证，承诺不是应付，发票不是现金，认缴不是实缴，毛利不是现金。
- 业务记录不得由会计凭证反向伪造；会计差错通过冲销、反向凭证或调整凭证修正。

## 4. M9 财务组织

游戏在“建立企业初始组织”后增加财务组织建设章节。最小可运行组织如下：

| 岗位 | 主要职责 | 禁止事项 |
| --- | --- | --- |
| CFO / 财务负责人 | 会计政策、资金、预算、报表和重大审批 | 不得独自创建并批准自己的付款或凭证 |
| 财务经理 / Controller | 总账、关账、对账、报表质量 | 不得修改已过账凭证 |
| 总账会计 | 凭证、期间、科目余额、结账 | 不得审批自己录入的重大手工凭证 |
| 应付会计 | 供应商发票、三单匹配、应付和付款建议 | 不得维护供应商银行账户并释放付款 |
| 应收会计 | 客户开票、应收、账龄和收款核销 | 不得修改客户信用并核销坏账 |
| 资金出纳 | 银行账户、收付款、银行对账 | 不得建立付款义务或过账总账调整 |
| 成本会计 | 标准成本、实际成本、在制品和差异 | 不得修改已批准 BOM/工艺事实 |
| 资产会计 | 在建工程、转固、折旧、盘点和处置 | 不得同时验收资产和批准付款 |
| FP&A / 经营分析 | 预算、滚动预测、经营分析和 KPI | 管理口径不得冒充法定会计事实 |
| 内部审计 / 独立复核 | 控制测试、异常审阅和审计轨迹 | 只读，不执行原始业务交易 |

小企业允许一人兼岗，但职责冲突必须通过外部复核、Founder 审批或系统限制补偿。组织、
岗位、用户、Agent、Mandate、额度、有效期和替代人均为版本化配置。

## 5. 三层财务语义

### 5.1 Core 语义

`economic_event`、`monetary_amount`、`accounting_period`、`account`、`debit`、`credit`、
`journal_entry`、`ledger_balance`、`document`、`party`、`organization`、`asset`、
`liability`、`equity`、`income`、`expense`、`cash_flow`、`approval`、`evidence`。

### 5.2 Domain 语义

`chart_of_accounts`、`accounting_book`、`fiscal_calendar`、`posting_rule`、
`accounting_event`、`journal_batch`、`subledger`、`reconciliation`、`close_task`、
`cost_element`、`cost_center`、`profit_center`、`cost_object`、`work_in_process`、
`standard_cost`、`actual_cost`、`cost_variance`、`receivable`、`payable`、
`bank_statement`、`fixed_asset`、`construction_in_progress`、`depreciation`、
`financial_statement`、`management_report`、`financial_kpi`。

### 5.3 HCTM 扩展语义

`plant_cost_center`、`production_line_cost_center`、`product_cost_object`、
`customer_project_profitability`、`tooling_asset`、`manufacturing_equipment_asset`、
`quality_loss_cost`、`downtime_cost`、`scrap_cost`、`warranty_cost`、
`energy_cost_pool`、`logistics_cost_pool`。

关键关系包括：

```text
business_document generates accounting_event
accounting_event compiled_by posting_rule
journal_entry contains journal_line
journal_line posts_to account
journal_line attributed_to cost_center/profit_center/project/product/order
subledger reconciles_to control_account
bank_statement_line matches cash_transaction
production_order consumes material/labor/overhead
cost_object accumulates cost_element
fixed_asset originates_from purchase_or_CIP
report_line aggregates ledger_balance
kpi derives_from governed_measure
```

## 6. 正式 Entity 目录

所有已发布财务 Entity 必须同时具备两类入口：数据模型工坊用于查看和配置模型；左侧
“财务管理”菜单用于查询和维护运行数据。场景包安装必须幂等创建对应
`menu.{entity_code}` 和 CRUD 权限资源，不允许出现“模型存在但业务数据无入口”的状态。
租户显式关闭侧边栏显示时除外。

### 6.1 基础设置

1. `finance_organization`
2. `accounting_book`
3. `fiscal_calendar`
4. `accounting_period`
5. `chart_of_accounts`
6. `gl_account`
7. `accounting_policy`
8. `posting_rule`
9. `document_number_range`
10. `currency_and_exchange_rate`
11. `financial_dimension`
12. `approval_matrix`

首版维度至少包含法人、基地、部门、成本中心、利润中心、项目、客户、供应商、产品、订单、
生产订单、资产和现金流项目。维度组合必须经过有效性校验，禁止依赖自由文本。

### 6.2 会计内核

1. `accounting_event`
2. `journal_entry`
3. `journal_line`
4. `journal_batch`
5. `ledger_balance`
6. `subledger_entry`
7. `reconciliation_run`
8. `reconciliation_item`
9. `manual_adjustment_request`
10. `period_close_run`
11. `period_close_task`
12. `accounting_exception`

`journal_entry` 必须包含来源、规则版本、制单人、审核人、过账人、业务日期、过账日期、
期间、币种、借方合计、贷方合计、状态、correlation、idempotency key 和 evidence ref。
已过账凭证不可更新或删除。

`journal_entry` 不是普通单表 Entity，必须使用 `document_with_lines`；`journal_line`
必须使用 `document_line`，通过 `parent_document_id` 形成受治理的主子表聚合。凭证头的
`lines` 必须编译为 `child_list`，支持一张凭证任意多行借贷分录。凭证提交、审批和过账时
至少两行、每行只能借或贷、借贷合计必须相等；头部合计由明细计算，不得作为独立可编辑真相。

多币种凭证同时保存：

- 交易币种及交易币借贷金额；
- 账套本位币及折算后的本位币借贷金额；
- 汇率类型、汇率日期、采用的汇率记录和不可变汇率快照；
- 币种最小货币单位，金额继续使用整数，汇率使用受控高精度 decimal，禁止浮点隐式舍入。

基础设置拆分为 `currency_definition` 和 `exchange_rate`。汇率按类型、源币种、目标币种和
生效日期版本化；修改未来汇率不能重算或改变已过账凭证。M9 的开业资本仍限定人民币，
但数据模型必须支持 M10–M13 的外币采购、销售、银行和重估。

模型工坊负责定义和发布模型，左侧菜单负责查询运行数据。M9 发布的全部 Entity（不只财务
Entity）必须同步生成受权限控制的数据菜单；财务 Entity 统一归入“财务管理”，不得要求
用户记住 Entity 编码后再从模型工坊绕行查询。

运行数据不得由 AESE 或 IAOS 前端静态构造。`finance_*` 权威账务与
`entity_projection_*` 通用 Entity
投影必须在同一业务事务或受治理的幂等历史回填中保持一致；例如 `gl_account` 菜单展示
1002 银行存款和 4001 实收资本，凭证明细通过 `account_code` 引用同一科目主数据。
`INC-INTERACTIVE-*` 人工演练案件是可审计业务案件，不能按 E2E/UI 自动化夹具清理。
模型工坊再次发布平台 Entity 必须保留既有 `physical_table_name`，不能把读取绑定切换到
新的空 `bo_*` 表。没有发生对应业务事实的后续 Entity 可以为空，但已提交事实存在而菜单
为空时必须诊断权威表、投影和元数据绑定三层，不得用演示数据掩盖。

`m9_*` 仅是首版投影的历史名称，不是财务或 Project Genesis 的长期存储边界。
IAOS DES-064 已将 19 张投影原位迁移为 `entity_projection_*`，保留 UUID、租户、
稳定业务编码、RLS 和全部业务行。M10–M13 必须复用稳定 Entity code 与规范投影；
禁止创建 `m10_*`、`m11_*` 等按里程碑复制的平行表。

### 6.3 资金、应收与应付

1. `bank_account`
2. `bank_statement`
3. `bank_statement_line`
4. `cash_transaction`
5. `cash_position`
6. `cash_forecast`
7. `customer_invoice`
8. `receivable_open_item`
9. `collection_and_application`
10. `credit_limit`
11. `supplier_invoice`
12. `payable_open_item`
13. `payment_proposal`
14. `payment_instruction`
15. `payment_execution`
16. `advance_and_deposit`
17. `bad_debt_assessment`

应收应付采用 open-item 管理，支持部分结算、预收预付、核销、账龄、争议、贷项和冲销。
银行余额只能由可信银行 Observation 或受治理导入形成，账面余额与银行余额必须对账。

### 6.4 成本与存货

1. `cost_element`
2. `cost_center`
3. `profit_center`
4. `cost_object`
5. `cost_pool`
6. `allocation_cycle`
7. `standard_cost_version`
8. `standard_cost_rollup`
9. `material_ledger_entry`
10. `production_cost_transaction`
11. `work_in_process_balance`
12. `cost_variance`
13. `inventory_valuation`
14. `landed_cost`
15. `scrap_and_rework_cost`
16. `product_cost_statement`

离散制造首版采用“标准成本用于计划与过程控制，期间实际成本用于结算”的双口径。成本
对象支持产品、批次、生产订单、客户订单、项目和资产；成本要素至少包括直接材料、直接
人工、机器、能源、制造费用、质量损失、外协、物流和折旧。

### 6.5 固定资产

1. `asset_class`
2. `fixed_asset`
3. `asset_component`
4. `asset_book`
5. `construction_in_progress`
6. `asset_acquisition`
7. `asset_transfer`
8. `depreciation_run`
9. `asset_impairment`
10. `asset_count`
11. `asset_disposal`
12. `asset_maintenance_cost_link`

设备采购不等于固定资产启用。采购、到货、安装、验收、在建工程归集、达到预定可使用
状态、转固、折旧、减值和处置是不同事实。

### 6.6 报表与指标

1. `financial_statement_definition`
2. `financial_statement_run`
3. `financial_statement_line`
4. `management_report_definition`
5. `management_report_run`
6. `financial_kpi_definition`
7. `financial_kpi_snapshot`
8. `variance_analysis`
9. `reporting_package`
10. `disclosure_and_note`

报表定义必须引用科目、维度、公式和口径版本；报表运行保存期间、数据快照、生成者、
批准者和 drill-through 路径。

## 7. 记账引擎与不可替代控制

```text
业务事实 committed
→ 生成 Accounting Event
→ 按有效日期选择 Posting Rule
→ 校验业务来源、维度、期间、币种和金额
→ 生成凭证草稿
→ 自动/人工复核与审批
→ 原子过账 Journal + Ledger + Outbox
→ 子账/总账对账
→ 报表与 KPI
```

硬约束：

- 每张凭证借贷相等，金额使用整数最小货币单位或显式 Decimal；
- 每个源业务事实和规则版本只能产生一次会计效果；
- 关闭期间禁止普通过账，重开必须独立审批；
- 自动凭证引用业务单据和规则 hash，手工凭证必须说明原因并强化审批；
- 子账控制科目禁止普通用户直接手工记账；
- 禁止删除已过账凭证，以冲销和更正凭证保留完整链；
- 业务日期、凭证日期、过账日期和 World 发生时间分别保存；
- Journal、Ledger、审计记录和 Outbox 在同一事务中提交。

## 8. M9 开业会计

### 8.1 M9 新增工作项

实缴资本核验后、初始组织建立前增加五个已发布工作项：

1. `finance.organization.configure`：安装财务负责人、Controller、总账、资金、成本和内审岗位模板；
2. `accounting.book.activate`：选择账套、CAS-BE、CNY 本位币、年度和首个开放期间；
3. `chart.of.accounts.activate`：启用开业科目及控制属性；
4. `capital.contribution.post`：从已验证银行到账和资本事实生成并过账双分录凭证；
5. `finance.opening.readiness.evaluate`：验证岗位模板、账套、期间、凭证和借贷平衡。

这五项都是 23 节点正式主流程的一部分，不再是
`capital.contribution.verify` 的隐藏内部副作用。玩家在 AESE 企业总部逐项发起，
IAOS `iaos-runtime` 执行确定性 Capability，但责任分别归属 `finance-lead`、
`finance-controller`、`general-ledger-accountant` 和 `internal-audit`。执行者不等于责任人。
需要人工审批的资本事实仍由既有 G4 冻结并决定；M9 不额外虚构一个“自动批准期初凭证”的门。

M9 终态增加 `finance_opening_ready=true`，并作为
`enterprise_operational_ready` 的必要条件。既有案件需要版本化迁移或补建流程，禁止
静默改变历史终态含义。

### 8.2 资本与成立费用示例

实收资本实际到账 1,000,000 CNY：

```text
借：银行存款                         1,000,000
  贷：实收资本                                   1,000,000
```

银行手续费 100 CNY：

```text
借：财务费用                               100
  贷：银行存款                                         100
```

尚未支付的设立服务费 5,000 CNY：

```text
借：管理费用                             5,000
  贷：其他应付款                                     5,000
```

认缴资本只形成资本承诺和披露，不自动形成上述实收资本凭证。只有可信银行到账 Observation
和已批准资本事项同时满足时才允许记账。

### 8.3 当前已实现的开业读取与硬门

当前 Runtime 1.8.0 已把资本核验后的五个财务子步骤编译为第 11–15 项持久工作项。
历史已完成案件不改写事实；新案件和仅完成首节点的可安全迁移案件使用 23 节点流程。

已实现读取模型：

- 银行日记账：按 `1002 银行存款` 展示逐笔收支、滚动余额和业务证据；
- 总账：按科目展示期初、借方、贷方和期末余额；
- 试算平衡：展示科目类别、余额方向及借贷发生额；
- 开业资产负债表：展示资产、负债、所有者权益和会计等式。

四个视图均由同一组 IAOS posted Journal 派生，不允许用户分别修改。最终
`enterprise.readiness.evaluate` 必须同时确认财务组织、账套、开放期间、资本凭证、
借贷平衡及资产负债表平衡，否则返回 `finance_opening_readiness_failed` 且不推进案件。

## 9. M10–M13 自动记账覆盖

| 阶段 | 业务事实 | 主要会计效果 |
| --- | --- | --- |
| M10 工厂建设 | 工程合同、进度确认、付款、验收 | 在建工程、应付、预付、银行付款 |
| M11 能力建设 | 设备采购、安装、投用、人员服务 | 在建工程/固定资产、折旧、薪酬应付、费用 |
| M12 产品工业化 | 工装、样件、材料、研发和客户预付款 | 项目成本、存货/费用、合同负债、应付 |
| M13 商业交付 | 采购、领料、生产、完工、接受、开票、回款 | AP、WIP、FG、COGS、收入、AR、现金 |

典型制造分录：

```text
材料入库：       借 原材料               贷 暂估应付/应付账款
生产领料：       借 生产成本-WIP         贷 原材料
人工与制造费用： 借 生产成本-WIP/制造费用 贷 应付职工薪酬/累计折旧/应付
产品完工：       借 库存商品             贷 生产成本-WIP
客户接受并确认收入：
                 借 应收账款             贷 主营业务收入/相关税项
                 借 主营业务成本         贷 库存商品
客户实际回款：   借 银行存款             贷 应收账款
```

具体确认时点必须由业务合同、客户接受事实和有效会计政策决定，游戏按钮不得直接指定会计
结论。

## 10. 成本核算

### 10.1 计划与标准成本

- BOM 数量、采购价格、工艺工时、人工费率、机器费率、能耗和制造费用预算；
- 按产品/版本/工厂/生效日期滚算标准成本；
- 标准成本发布需要工程、生产和财务联合审批；
- 已发布版本只读，新版本通过差异和生效日期替换。

### 10.2 实际成本

- 直接材料来自批次领退料；
- 直接人工来自获批工时或产量驱动；
- 机器、能源、折旧和间接费用先进入成本中心/成本池；
- 按机器小时、人工小时、产量、面积、能耗或作业动因分摊；
- 月末计算 WIP、完工成本、销售成本及采购价差、用量差异、效率差异、产能差异；
- 支持产品、订单、批次、项目、客户和基地多维穿透。

### 10.3 核算方法

首版支持标准成本 + 分批/订单实际成本；后续可配置移动加权平均、个别计价、作业成本、
目标成本、变动成本和生命周期成本。方法必须按成本对象和期间版本化，禁止中途无痕切换。

## 11. 流程、审批与职责分离

核心 Process：

1. `finance.foundation.setup.v1`
2. `capital.accounting.v1`
3. `record.to.report.v1`
4. `procure.to.pay.accounting.v1`
5. `order.to.cash.accounting.v1`
6. `bank.reconciliation.v1`
7. `fixed.asset.lifecycle.v1`
8. `manufacturing.cost.close.v1`
9. `period.close.and.report.v1`
10. `budget.forecast.control.v1`

审批 Flow：

- 科目/政策/记账规则发布；
- 手工凭证、跨期凭证、冲销和期间重开；
- 供应商主数据与银行账户变更；
- 采购付款、员工报销和资金调拨；
- 信用额度、坏账、贷项和核销；
- 固定资产转固、减值和处置；
- 标准成本发布、费用分摊和重大成本差异；
- 月结、报表发布和期初余额确认。

路由支持金额阈值、风险、法人、部门、项目、付款方式和异常类型；支持串行、并行、会签、
或签、上级、岗位、指定主体及替代人。制单、审核、付款、银行账户维护、对账和审计之间
实施职责分离，Agent 不得批准自己的高风险动作。

## 12. Capability 与 Agent Tool

最小 Capability 族：

- `finance.organization.configure`
- `accounting.book.create`
- `chart.of.accounts.activate`
- `posting.rule.publish`
- `accounting.event.recognize`
- `journal.draft.generate`
- `journal.entry.post`
- `journal.entry.reverse`
- `opening.balance.post`
- `supplier.invoice.match`
- `payment.proposal.prepare`
- `payment.execute`
- `customer.invoice.issue`
- `receivable.collect.apply`
- `bank.statement.import`
- `bank.reconciliation.execute`
- `asset.capitalize`
- `depreciation.run`
- `standard.cost.rollup`
- `manufacturing.cost.settle`
- `period.close.execute`
- `financial.statement.generate`
- `finance.kpi.calculate`

Finance Agent 只能通过这些已发布能力和只读查询 Tool 工作。每个 Tool 暴露 purpose、
input/output contract、权限、额度、数据范围、失败升级、证据和幂等语义。自动建议、异常
解释和预测不得直接过账。

## 13. 结账与对账

月结任务至少包括：

1. 银行对账；
2. AP/AR 子账与总账控制科目对账；
3. 存货数量、价值和总账对账；
4. WIP、完工和销售成本计算；
5. 固定资产折旧及资产台账对账；
6. 薪酬、费用、预提和摊销；
7. 外币重估（启用多币种后）；
8. 关联交易和跨组织对账（启用多法人后）；
9. 异常凭证、悬账和未匹配事项清理；
10. 试算平衡、报表生成、复核和期间关闭。

每项任务保存 owner、依赖、截止时间、状态、差异、处理证据和复核结果。存在不平衡、
未解释重大差异或关键子账未对账时必须失败关闭。

## 14. 报表、图表与 KPI

### 14.1 财务报表

- 试算平衡表、科目余额表、总账和明细账；
- 资产负债表、利润表、现金流量表、所有者权益变动表；
- 银行日记账、应收/应付账龄、固定资产及折旧表；
- 成本中心、利润中心、产品、订单和项目损益；
- 预算执行、滚动预测和差异分析；
- 报表附注、口径、快照和审计穿透。

### 14.2 管理驾驶舱

- 收入、毛利、EBITDA、经营利润、净利润；
- 经营现金流、自由现金流、现金可支撑月数、现金预测偏差；
- DSO、DPO、DIO、现金转换周期；
- 应收逾期率、坏账风险、应付到期结构、付款及时率；
- 存货周转、呆滞库存、WIP 周期、库存账实差；
- 标准/实际成本、采购价差、材料用量差异、人工效率差异、制造费用差异；
- 产品/客户/订单/项目贡献毛利；
- CAPEX 预算执行、资产利用率、在建工程超期、折旧负担；
- 预算偏差、费用率、单位制造成本、质量损失成本、停机成本；
- 关账周期、自动凭证率、手工凭证率、对账未清项和控制异常数。

所有图表必须显示期间、币种、口径版本、数据更新时间和 drill-through。KPI 阈值和公式是
版本化语义资产，不在前端硬编码。

## 15. 游戏体验

M9 企业总部新增“财务中心”，玩家依次完成：

1. 查看六个财务责任岗位，识别 Controller、总账、资金和成本岗位空缺；
2. 检查四条职责分离规则，冲突兼任由 IAOS 数据库拒绝；
3. 发起 `finance.organization.configure`；
4. 发起账套/期间与科目启用；
5. 在银行 Observation 和 G4 资本事实成立后发起资本过账；
6. 发起内部审计财务就绪检查；
7. 穿透查看“银行到账 → 资本事项 → Capability Execution → 凭证 → 总账 → 报表”；
8. 在 IAOS“组织与待办”查看执行者、责任岗位、Mandate、通知对象和岗位空缺升级。

后续 M10–M13 每个业务章节同时出现财务后果，但不要求玩家逐张手工制证。正常自动凭证由
规则生成；异常、重大金额、手工调整、关账和政策选择才进入玩家/财务人员决策。

## 16. 菜单与可解释入口

IAOS 目标菜单：

- 财务工作台
- 总账与凭证
- 应收管理
- 应付管理
- 资金与银行
- 固定资产
- 存货与成本
- 预算与预测
- 结账中心
- 财务报表
- 管理驾驶舱
- 财务基础设置

每个页面必须提供“功能说明”，解释目的、角色、前置配置、操作、上下游关系、校验、
恢复和权威设计。业务单据、凭证、科目余额、报表和 KPI 之间提供双向穿透。

### 16.1 当前 M9 已实现入口与边界

当前已交付 IAOS `财务账务与报表`（`#finance_workspace`）第一纵切，可按设立案查看
财务组织、账套、期间、开业凭证、银行日记账、总账/试算平衡和开业资产负债表。AESE
企业总部与治理档案通过“查看组织与待办 / 查看系统账务 / 查看财务报表”进入该页面。
组织视图来自 `GET /api/v1/finance/opening/:case_code/operations`，账务和报表来自
`GET /api/v1/finance/opening/:case_code` 的同一已过账凭证，不是游戏投影的复制。

Runtime 1.8.0 安装六条 `finance_duty_definition`、四条阻断型
`finance_sod_rule`，并只为已有任职主体的 `finance-agent@finance-lead` 与
`audit-agent@internal-audit` 保留 active Mandate。其余岗位明确为 vacant；工作项激活时
生成幂等 `incorporation.work_item.assigned` Outbox。空缺责任岗位会通知 Founder 补位，
但不会把 Founder 伪装成财务任职人。

完整财务报表按“业务事实/子账 → 凭证与总账 → 期间报表读取模型 → 已审批发布快照 →
管理 KPI/驾驶舱”五层建设。M9 当前只完成开业凭证和开业资产负债表，不包含利润表、
现金流量表、所有者权益变动表、月结、合并、预算和制造成本驾驶舱；这些仍属于
F25–F29，不能把当前工作台描述成完整财务系统。

## 17. 数据、安全与验收

- 全部财务 Entity 启用 tenant RLS；
- 法人、账套、期间、币种、金额和维度不可缺失；
- 关键主数据与规则版本化，发布后不可原地修改；
- 银行账户、身份证明、工资和客户信用属于敏感数据；
- 导入默认 dry-run，过账、付款、关账要求显式 apply；
- 验证重复事件不重复记账、失败无部分写入、跨租户不可见；
- 验证业务子账、总账、World 现金和报表可对账；
- 验证重启、乱序、迟到事件、冲销、期间关闭及重开；
- 100% 已过账凭证可追溯到业务事实或批准的手工调整；
- 报表每一金额可穿透至余额、凭证、会计事件和原始业务证据。

## 18. 分阶段完成定义

| 阶段 | 完成条件 |
| --- | --- |
| M9-F0 | 财务语义、Entity、Capability、Process、Policy 和职责矩阵发布 |
| M9-F1 | 财务组织、账套、科目、期间和期初资本凭证完成 |
| M10-F2 | 工程承诺、应付、付款和在建工程会计完成 |
| M11-F3 | 设备转固、折旧、薪酬费用和能力建设成本完成 |
| M12-F4 | 项目成本、工装、样件和合同负债会计完成 |
| M13-F5 | P2P、库存/WIP/FG、O2C、收入、应收、回款和实际成本完成 |
| FIN-F6 | 月结、三表、管理报表、KPI、恢复、安全和三视口验收完成 |

设计批准不代表上述实现完成。任何阶段只有在 Semantic Studio、数据模型工坊、能力工作室、
流程工作室、审批工作台、Agent、API、数据库、业务 UI 和游戏中均有反向证据后才能关闭。
