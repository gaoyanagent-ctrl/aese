---
id: SOL-011
title: M10 Agent 候选越过投资上限并被永久重放
date: 2026-08-02
status: completed
author: Codex + User
tags: [aese, m10, minimax, governance, idempotency]
---

# M10 Agent 候选越过投资上限并被永久重放

## 1. 现象

M10 保存 Requirement 后 MiniMax 成功返回候选，但 IAOS `site.proposal.record` 返回 422：

```text
proposal ... exceeds user-approved investment request
```

同一按钮再次执行仍返回相同错误。

## 2. 根因

IAOS 正确要求每个候选的最高估算额不得超过用户提交的投资申请额。AESE 原合同只校验
金额区间顺序，没有在提交前校验该上限；同时 CreativeJob 在 IAOS 写入前就被标记为
`completed`，幂等重试因此永久重放非法输出。

## 3. 修正

- `plant-planning-v2` 把币种、精度和最高金额写成模型硬约束。
- AESE `ValidateProposalSet` 执行与 IAOS 相同的金额上限校验。
- 首次模型输出不合法时携带具体原因进行一次完整 JSON 修订，不静默改写金额。
- 两次调用的 Token 计入同一 Agent 运行证据。
- IAOS 原子提交成功后才保存 completed；提交失败保存 failed，可安全重试。
- 旧版本 completed 证据若不符合当前合同，自动作废并重新生成。

## 4. 验收

刷新 M10 页面并对 active Requirement 重试。成功时 Proposal 每一项的 maximum 均不超过
投资申请额，CreativeJob 和 IAOS Agent Run 为 completed/valid，且 Proposal revision 可在
IAOS M10 工作台穿透查询。若修订后仍无上限内可行方案，页面应显示模型校验失败，用户可
提高并重新审批投资申请额或填写人工候选。
