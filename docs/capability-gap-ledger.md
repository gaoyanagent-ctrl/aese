# M8 Capability Gap Ledger

| 企业活动 | 角色 | 对象 | IAOS 权限/Capability | 当前处理 | 优先级 |
| --- | --- | --- | --- | --- | --- |
| 上报设备观察 | 设备工程师或 Agent | `LAS-WLD-02` | `world.observation.ingest` | World Bridge observation journal + Outbox | P0 已交付 |
| 提出检修 | 设备工程师或 Agent | 维护工单 stable ref | `world.intent.create` | World Bridge intent journal；人和 Agent 共用相同端点与授权资源 | P0 已交付 |
| 批准并提交检修结果 | 有审批权的人员/Agent | 维护工单 committed refs | `world.intent.approve` | 只接受 committed/no-op outcome，journal 与 Outbox 同事务 | P0 已交付 |
| 恢复桥接消费 | AESE runtime | world run | `world.bridge.read` | cursor query 为事实，SSE 仅提示 | P0 已交付 |
| 通用 EAM 工单编排 | EAM 专员 | 任意设备 | 尚无通用 Capability | 本里程碑只覆盖 Genesis 固定 tracer，不扩成 EAM 引擎 | P1 后续 |
| 自动审批高风险维修 | 自主 Agent | 维修预算/停线 | 不允许 | 必须由 IAOS Policy/Approval 决定，AESE 不代批 | 禁止 |

长期 JWT、AESE 直写 IAOS 数据库、正式 direct NATS 和 webhook-only 路径均不是临时方案。
