---
id: DES-031
title: IAOS 多组织财务与共享主数据基础
date: 2026-07-30
status: active
author: Codex + User
tags: [finance, multi-organization, master-data, data-sharing, ledger, governance]
---

# IAOS 多组织财务与共享主数据基础

## 1. 业务问题与设计目的

当前 M9 财务纵切以单租户、单法人、单账套为参考实现，不能直接扩展成集团财务：

- 集团公司、分子公司、事业部、业务单元、基地和财务共享中心没有统一组织模型；
- `accounting_book` 同时承担账套配置、账户原型和法人归属，边界混乱；
- 科目、客户、供应商等数据只能“每公司复制”或“整个租户共用”，缺少受治理共享范围；
- 主数据的全局身份、法人属性、采购/销售组织属性没有分层；
- 会计期间只有账套级状态，尚未记录 GL、AP、AR、资产、库存、成本等模块的独立开关。

本设计建立 IAOS 后续财务、采购、销售、制造和主数据治理共同使用的最小多组织底座。
它不在 AESE 建立第二套组织或财务引擎。

### 1.1 第一实施切片（F5B，2026-07-30）

F5B 已由 IAOS
`docs/designs/DES-067-finance-multi-organization-and-reference-data-sets.md` 实现：

- `finance_organization_unit`、`finance_reference_data_set`、
  `finance_reference_data_assignment`、`finance_organization_access` 四张权威表；
- `finance.enterprise.structure.configure` 与 `finance.reference.data.configure`
  两项受治理业务 Capability；
- FORCE RLS、带 tenant 的复合外键和 `require_capability_execution()` 写入硬门；
- 财务工作台“多组织与共享数据”业务表单，以及按已有 M9 账套生成默认组织/Data Set
  的幂等迁移；
- AESE `scenario-packs/hctm/finance-governance-baseline.json` 只保存确定性场景模板，
  不复制 IAOS 权威数据，position access template 在应用时解析为平台 subject。

### 1.2 第二实施切片（F5C，2026-07-30）

F5C 已由 IAOS
`docs/designs/DES-068-finance-ledger-chart-calendar-foundation.md` 实现：

- 财政日历/期间、共享科目表/科目定义、法人科目扩展、账簿与账簿集合权威表；
- `finance.ledger.foundation.configure` 受治理 Capability、FORCE RLS、复合 tenant
  外键和数据库写入硬门；
- 财务工作台“账簿与科目”查询和业务配置表单；
- M9 历史 `BOOK-*`、凭证关系、1002/4001 科目和 12 期日历的原位幂等迁移；
- 目标租户 `tenant-gx-f4b3ce3ce8e2712d` 已迁移为 1 账簿、1 科目表、2 科目定义、
  2 法人扩展、1 日历/12 期间和 1 账簿集。

AESE baseline 1.2 固定同构的 HCTM 确定性模板并离线校验所有引用和期间连续性。
F5D Business Partner/产品组织扩展仍未实现，不能因 F5C 完成而标记为完整集团财务。

## 2. 设计原则

1. **Tenant 是客户安全边界，不等于法人公司。** 一个集团及其受控分子公司通常位于同一
   tenant；跨 tenant 分享属于后续企业网络/数据交换，不走本设计的内部共享捷径。
2. **法人是法定记账和交易责任边界。** 业务单元、基地和共享中心不能替代法人归属。
3. **共享配置，不共享法定余额。** 科目表、支付条件等可共享；凭证、余额、应收应付、
   库存和税务事实必须有明确法人及账簿所有者。
4. **一个全局身份，多个组织扩展。** 客户、供应商、物料等建立集团级 canonical identity，
   法人、采购组织、销售组织和工厂只维护自己的扩展属性。
5. **引用优先于复制。** 共享数据通过 Data Set assignment 被消费；复制只用于有版本、
   来源、同步状态和冲突处理的发布/订阅场景。
6. **共享不等于可修改。** 消费组织默认只读引用；只有数据所有者或授权 Steward 可改。
7. **组织、数据范围和权限分别建模。** 组织树不能替代 RLS，角色也不能隐式决定数据归属。

## 3. 多组织模型

```text
Tenant / Enterprise Boundary
└─ Enterprise Group
   ├─ Legal Entity (母公司)
   │  ├─ Business Unit
   │  │  ├─ Site / Plant / Warehouse
   │  │  ├─ Sales Organization
   │  │  └─ Purchasing Organization
   │  └─ Finance Organization
   ├─ Legal Entity (子公司)
   │  └─ ...
   ├─ Shared Service Center
   ├─ Management Accounting Area
   └─ Consolidation Group / Consolidation Unit
```

### 3.1 组织对象及责任

| 对象 | 目的 | 是否拥有法定交易/余额 |
| --- | --- | --- |
| `enterprise_group` | 集团治理、共享策略和合并范围 | 否 |
| `legal_entity` | 合同、税务、法定报表和资金责任主体 | 是 |
| `business_unit` | 交易处理、管理责任、权限和共享数据选择上下文 | 否，交易仍带法人 |
| `site` / `plant` | 实物、产能、库存和成本发生地点 | 否 |
| `finance_organization` | 财务岗位、职责分离和服务关系 | 否 |
| `shared_service_center` | 为多个法人提供 AP/AR/GL 等服务 | 否，不取得客户法人的数据所有权 |
| `management_accounting_area` | 跨法人管理会计范围 | 只拥有管理口径 |
| `consolidation_unit` | 集团合并、抵销和报表映射 | 只拥有合并调整 |

所有业务单据至少携带 `tenant_id`、`legal_entity_id`、`business_unit_id`；涉及实物时再携带
`site_id/plant_id`，涉及核算时携带 `accounting_book_id`。这些字段是稳定引用，不允许自由文本。

## 4. 财务账簿、科目表和组织关系

### 4.1 账簿模型

一个法定账簿 `accounting_book` 必须归属一个法人，但一个法人可以有多个账簿：

- `primary`：主要法定/经营账簿；
- `local_statutory`：当地准则账簿；
- `reporting`：集团或其他准则报告账簿；
- `tax`：税务调整账簿；
- `management`：管理口径账簿。

多个账簿可组成 `ledger_set` 用于统一查询或期间操作，但不能因此混合凭证序列和法定余额。
平行账簿通过版本化 `accounting_rule_set` 从同一会计事件派生，不复制原始业务单据。

账簿必填字段：

| 字段 | 含义 |
| --- | --- |
| `book_code`, `book_name` | 稳定编码和用户名称 |
| `legal_entity_id` | 唯一法人所有者 |
| `book_role` | primary/local/reporting/tax/management |
| `accounting_standard_code` | CAS-BE、IFRS 等 |
| `functional_currency_code` | 本位币 |
| `chart_of_accounts_id` | 运营科目表 |
| `fiscal_calendar_id` | 财政日历 |
| `balancing_segment_rule` | 法人/利润中心等平衡段 |
| `retained_earnings_account_id` | 本年利润结转科目 |
| `intercompany_rule_set_id` | 内部交易平衡规则 |
| `effective_from`, `effective_until`, `status`, `version` | 生效和版本治理 |

### 4.2 科目表分层

```text
Group Chart of Accounts
        ↓ mapping
Operating Chart of Accounts ── shared by one or more legal entities
        ↓ company extension
GL Account Legal Entity Settings
        ↓ optional mapping
Local / Statutory Chart of Accounts
```

- `chart_of_accounts` 和 `gl_account_definition` 是可共享定义；
- `gl_account_legal_entity` 保存法人特有的启用状态、统驭科目、税务、币种、字段状态、
  允许过账类型等；
- `account_mapping` 保存运营科目、集团科目和当地科目的有效期映射；
- 同一管理会计范围内的法人原则上使用同一运营科目表和兼容财政日历；
- 已有余额的科目不能通过修改共享定义改变历史含义，变更必须新版本生效。

这采用“科目表段 + 公司段”的思路，避免把 1002 银行存款为每家公司复制成互不关联的
主数据，同时保留各法人不同的控制属性。

## 5. 共享主数据架构

### 5.1 Data Set / Set ID

新增：

- `reference_data_set`：`COMMON`、集团、区域、事业部或组织专用的数据集合；
- `reference_data_set_assignment`：把数据类型的数据集分配给业务单元、法人、采购组织、
  销售组织、工厂、账簿等 determinant；
- `master_data_ownership`：数据所有范围、Steward、治理状态和有效期；
- `master_data_share_policy`：消费者范围、`read/reference/extend` 权限和冲突策略；
- `master_data_change_request`：新增、变更、合并、停用和解除共享的审批对象。

交易选择值为：

```text
Common Set
+ 交易上下文所分配的数据集
+ 显式授权的集团/区域数据集
- 已停用、未生效、无权限和不适用于当前法人的记录
```

### 5.2 主数据分层

| 数据 | 集团/共享层 | 组织扩展层 |
| --- | --- | --- |
| 总账科目 | 科目编号、名称、类型、集团映射 | 法人启用、统驭科目、币种、字段控制 |
| Business Partner | 法定名称、统一身份、地址、关系、去重键 | 法人客户/供应商角色、税务、付款条件、信用、对账科目 |
| 供应商 | canonical supplier identity | 采购组织站点、收货/容差/币种；法人付款与预提税属性 |
| 客户 | canonical customer identity | 销售区域、定价/交付；法人信用、开票、收款与对账属性 |
| 产品/物料 | 编码、描述、基础单位、分类 | 工厂 MRP/成本/库存属性；销售/采购组织属性 |
| 支付条件/交易类型 | 可放入 Common/集团/区域 Data Set | BU 专用覆盖或扩展 |
| 汇率/币种/单位 | 平台或集团共享参考数据 | 法人采用策略和允许类型 |

客户和供应商不再各自保存一份互不关联的公司身份，而以 `business_partner` 为 canonical
identity，再挂 `customer_role`、`supplier_role` 和组织扩展。银行账户、税号、信用额度等
敏感字段独立授权，不能因为主体身份共享就自动向所有组织开放。

### 5.3 共享与交易数据边界

- 可共享：主数据定义、参考数据、规则模板、科目表、日历模板、报表模板。
- 可共享但需组织扩展：客户、供应商、产品、科目、成本中心模板。
- 不共享可写行：凭证、余额、发票、应收应付、付款、银行流水、库存批次、成本结算。
- 跨法人交易必须产生双方业务单据、内部交易伙伴和可对账引用；集团查询使用授权视图，
  合并使用 consolidation journal，不直接修改成员公司的凭证。

## 6. 数据安全、RLS 与治理

每条 set-enabled 主数据至少携带：

`tenant_id`、`data_set_id`、`owner_scope_type`、`owner_scope_id`、`steward_position_id`、
`governance_status`、`effective_from/until`、`version`、`source_ref`。

读取条件为 tenant 相同且：

1. 调用者拥有交易上下文组织的数据权限；
2. 记录位于 Common Set、该数据类型分配的数据集或显式共享策略中；
3. 字段级敏感权限允许读取。

写入条件为 owner/steward 权限与已批准 Change Request 同时成立。消费者对共享数据的
组织扩展写入自己的 extension 表，不能修改 canonical row。Agent 使用相同 Capability 和
数据范围，不获得额外旁路。

## 7. 会计期间模块控制（已记录，后续实现）

未来新增 `accounting_period_control`：

```text
(tenant, legal_entity, accounting_book, fiscal_year, period, module)
module = GL | AP | AR | FA | INV | COST | CASH | TAX | PROJECT
status = future | open | soft_closed | closed | reopened
```

要求：

- 子模块可以先关闭，GL 最后关闭；
- `soft_closed` 只允许特许角色/调整来源；
- 重开必须有审批、理由、有效时间窗和审计证据；
- 跨期单据按业务日期、会计日期、模块状态和迟到策略共同决定；
- 关账编排以 `period_close_run/task` 驱动，不能靠直接更新 status。

本轮只冻结该模型和状态机，不实现 UI、Capability、Process 或数据库迁移。

## 8. 原账套模型问题与已完成迁移

`tenant-gx-f4b3ce3ce8e2712d` 的权威账套行已有编码、名称、法人、本位币、准则、启用时间、
年度起始月和状态；但通用投影存在以下结构性问题：

1. `accounting_book` 错用 `account` 原型，出现 `account_name/account_type/currency/opened_at`
   等账户语义，账套自身字段只作为附加字段；
2. 权威账套名称与 Entity 投影显示名称可能不同；
3. payload 没有完整携带账套名称、科目表、财政日历、账簿角色和版本；
4. 缺少集团、法人 ID、共享科目表、财政日历及多账簿关系；
5. “字段非空”不等于模型正确，不能继续用默认值掩盖错误原型。

F5C 已按以下合同完成迁移：

- 新建正确的 `accounting_book` 财务配置原型，不复用银行/总账账户原型；
- 先补组织、科目表、财政日历和账簿 assignment，再迁移现有行；
- 保留稳定 `BOOK-*` 编码、UUID、凭证外键和审计证据；
- 权威表与 Entity 投影在同一受治理 Capability 中原子更新；
- 对现有租户生成差异报告，禁止用空字符串或 `other` 作为必填字段的静默默认；
- 迁移前后核对账簿、12 期、科目、凭证、余额和报表。

## 9. 实施顺序

1. 组织基础：集团、法人、BU、基地/工厂、共享服务中心、管理会计范围；
2. Data Set：共享集、组织分配、所有权、共享策略、Change Request；
3. 财务基础：科目表/科目两段式、财政日历、账簿/账簿集、多准则映射；
4. 主数据：Business Partner canonical + 客户/供应商组织扩展，产品组织扩展；
5. 迁移：修正当前账套原型和必填字段，迁移现有 M9 数据；
6. 运行控制：组织上下文、RLS、Capability、Approval、Agent Tool 和穿透 UI；
7. 后续：模块期间开关、关账编排、跨法人交易和合并。

M10 业务模型不得在第 1–5 步未形成稳定合同前继续复制单法人字段。

## 10. 验收标准

- 一个 tenant 可建立集团、母公司、子公司、BU、工厂和共享服务中心；
- 每个法定交易、凭证和余额可唯一定位法人和账簿；
- 两个法人可共享同一运营科目表，但有独立法人科目属性和余额；
- 同一 Business Partner 可被多个法人/采购/销售组织扩展，不复制 canonical identity；
- Common/集团/BU 专用数据选择符合分配，未授权组织不可见且不可修改；
- 共享主数据变更有 Steward、审批、版本、影响分析和订阅者通知；
- 跨法人业务不混账，集团查询和合并不改写成员公司法定凭证；
- 现有 M9 账套不再出现错误账户原型字段或静默默认值；
- 不需要读取数据库或 JSON 才能理解组织归属和共享来源。

## 11. 参考实践

- [SAP：Company Code 是可形成完整独立法定报表的最小组织单元](https://help.sap.com/docs/SAP_S4HANA_ON-PREMISE/5e23dc8fe9be4fd496f8ab556667ea05/b541de531ed3424de10000000a174cb4.html)
- [SAP：运营、集团和国家/地区科目表及多 Company Code 共享](https://help.sap.com/docs/SAP_S4HANA_ON-PREMISE/8fbeed5f2046489696a50ac7fd76f9c6/9f53c2531bb9b44ce10000000a174cb4.html)
- [SAP：G/L Account 的科目表层与 Company Code 层](https://help.sap.com/docs/s4hana-best-practices/bc5ee9baf05e4fe6820985e2dbfbae63/content?locale=en-US)
- [SAP：Controlling Area 跨 Company Code 的共同科目表要求](https://help.sap.com/docs/SAP_S4HANA_ON-PREMISE/5e23dc8fe9be4fd496f8ab556667ea05/dd3dd253913e4608e10000000a174cb4.html)
- [SAP：Business Partner/Supplier 中央治理与复制](https://help.sap.com/docs/SAP_MASTER_DATA_GOVERNANCE/db97296fe85d45f9b846e8cd2a580fbd/73940ddfd14f41c7bc808bba05f674f4.html)
- [Oracle：Legal Entity、Ledger、Business Unit、科目表和会计日历的企业结构](https://docs.oracle.com/en/cloud/saas/financials/24b/faigl/implementing-enterprise-structures-and-general-ledger.pdf)
- [Oracle：Reference Data Set 在 BU 之间共享或隔离参考数据](https://docs.oracle.com/en/cloud/saas/financials/26b/faigl/reference-data-sets.html)
- [Oracle：Business Unit 用于交易处理、数据安全和参考数据共享](https://docs.oracle.com/en/cloud/saas/financials/25d/fafcf/business-units.html)
- [Oracle：Supplier Site 保存采购 BU 与供应商的组织特有关系](https://docs.oracle.com/en/cloud/saas/procurement/25c/oaprc/supplier-sites-and-supplier-site-assignments.html)
