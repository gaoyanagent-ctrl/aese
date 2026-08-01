# M9 知识场景上下文验证证据

日期：2026-08-01

- AESE 只从当前 Genesis workspace、WorldProjection 和 GameWorkItem 生成封闭导航字段。
- 深链同时保留既有 tenant/case/article 路由参数，不改变知识文章租户权限边界。
- IAOS 前端与 Copilot BFF 复用同一归一化器；未知、超长和提示式值不会进入 Agent 提示。
- 知识中心明确标记“导航上下文不是运行证据”并允许用户清除。
- 当前切片不宣称已读取 World Observation、IAOS Runtime、Journal 或 Outbox；该能力属于 S9。
- AESE Vitest 覆盖全部深链字段，TypeScript、Go 测试和生产构建作为发布门。

