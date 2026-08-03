---
id: SOL-015
title: M10 项目 WBS 治理字段二次纠正
date: 2026-08-03
status: completed
author: Codex + User
tags: [aese, m10, agent, wbs, governance]
---

# M10 项目 WBS 治理字段二次纠正

## 现象

项目 Agent 已能从截断 JSON 恢复，但修订结果仍可能返回：

```text
facility project output invalid after one repair:
project option 1: WBS item "WBS-02" is invalid
```

## 根因

旧实现只有一个共享修订预算。JSON 截断先消耗该预算后，完整结果中的序号、日期、阶段、责任岗位、验收标准或预算分摊错误无法再纠正；校验器还把所有字段错误压缩为同一条信息。

## 修正

- 输出完整性恢复与业务治理纠正分别最多执行一次，总调用次数最多三次。
- WBS 校验逐项返回具体字段和期望值，治理纠正只消费该错误摘要与原 IAOS 权威输入。
- 修订格式固定四阶段数组顺序和连续序号，要求 RFC3339 日期处于项目边界内，并强制非空责任岗位、验收标准及合计 10000 的正数预算分摊。
- 任何预算耗尽后仍非法的输出继续失败关闭，不循环重试、不改写投资权限、不使用固定方案。

## 验证

- `TestAIPlanningProviderSeparatesCompletionAndGovernanceRepairs`
- `TestValidateProjectPlanOptionExplainsInvalidWBSField`
- `go test ./...`
- `go vet ./...`
