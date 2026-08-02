---
id: SOL-013
title: M10 Review 恢复与幂等提交
date: 2026-08-02
status: completed
author: Codex + User
tags: [aese, m10, review, idempotency]
---

# M10 Review 恢复与幂等提交

## 现象

候选已经完成“采纳调研”，页面仍可再次提交审阅，并收到 409 Conflict。数据库存在同一
ProposalSet、Proposal 和 revision 的权威 Review。

## 根因与修正

IAOS Review 是不可变人员决定，旧页面却只用本地 `savedReviews` 标记成功，没有在刷新时
读取权威 Review；重新点选还会清除本地 saved 状态。

AESE 现在通过命名 IAOS adapter 暴露只读 Review BFF。页面恢复并锁定已提交 action/reason，
按 action 显示下一步。服务端提交前读取现有 Review：完全相同返回幂等成功；不同决定返回
明确冲突，要求创建新的 ProposalSet revision，绝不覆盖历史决定。
