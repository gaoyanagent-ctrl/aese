---
id: SOL-006
title: 银行开户审批后 Work Item Execute 400
date: 2026-07-28
status: completed
author: Codex + User
tags: [aese, genesis, banking, iaos]
---

# 银行开户审批后 Work Item Execute 400

## 现象

G3 银行开户审批提交和批准成功，随后
`POST /incorporations/:case/work-items/8/execute` 返回 400。

## 根因

开户银行和资料清单使用 `business_note` 写入审批 intent。前端在批准后又把同一对象原样
传给 Work Item Execute，而 IAOS 使用 `DisallowUnknownFields` 严格解码
`incorporation.Command`；该执行合同没有 `business_note`，因此返回
`invalid work item input`。

## 修复

`approveAndExecuteWorkItem` 分离审批快照输入和执行输入：

- Gate submit 保留完整业务说明；
- Execute 只发送 `correlation_id` 以及该 Capability 真正支持的金额和币种；
- 不放宽 IAOS 严格合同，也不丢失审批事项中的开户银行和资料证据。

回归测试锁定 G3 Gate intent 包含银行说明，但第 8 步 Execute body 不含
`business_note`。

## 重试时的重复审批

首次请求已经把 G3 审批置为 `approved`，随后才在 Execute 解码处失败。修复请求体后
重试会幂等取得同一审批；旧前端仍再次调用 approve，IAOS 因审批已终结而正确返回 403。
前端现在读取 Gate submit 返回的审批状态：已 approved 时跳过重复决定，直接重试
Execute。回归测试锁定该恢复路径只发 Gate submit 和 Execute 两个请求。
