---
id: DES-024
title: M22 集团财务资金与投资治理
date: 2026-07-22
status: completed
author: Codex + User
tags: [m22, finance, treasury, investment]
---

# M22 集团财务资金与投资治理

## 1. 目标

在多产品、多节点和持续运营之上形成管理口径的集团价值闭环：计划/实际利润、现金、营运资金、融资、投资项目和情景压力决策。

机器终态：`group_value_cycle_closed=true`。

## 2. 首版范围

- 法人/基地/客户/产品/项目维度的 management P&L、cash flow 和 working capital。
- AR/AP/库存/合同负债、内部调拨简化抵销、资金池和融资额度。
- 一个 cash stress + capex choice tracer，比较扩产、外协和延期投资。
- 管理合并和决策支持，不声称法定报表、税务或审计合规。

## 3. 边界

银行实际到账/支付和外部融资结果属于 World；IAOS 拥有预算、凭证/台账、结算、资金计划、投资审批和管理报表。AESE 不建立第二总账，不用 World cash 替代 IAOS 会计记录。

利润、现金、预算、承诺和可用额度必须分离；内部交易双方、在途和抵销有稳定引用和金额/币种/汇率版本。

## 4. 完成标准

- 经营、投资、融资现金和管理利润可跨产品/基地对账。
- 计划/实际、standard/actual、cash/profit 和 legal/management view 明确。
- 重复抵销、伪造银行事实、超授权融资、自批投资和缺 evidence 盈利结论失败关闭。
- M17 IBP、M19 network、M20 warranty、M21 resilience 的价值影响可解释。

## 5. 非目标

法定总账、税务申报、真实汇率/银行、审计意见和完整合并会计不在范围。
