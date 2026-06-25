# AESE Agent Instructions

所有进入本项目的 agent 必须先阅读：

1. `README.md`
2. `docs/agent-project-context.md`
3. `docs/AESE_MVP_Blueprint.md`
4. `docs/progress-log.md`

## 工作原则

- AESE 是 IAOS 的智能企业运行仿真环境，不是独立 ERP、MES 或游戏项目。
- 第一阶段优先业务真实性、事件运行和 IAOS 能力映射，视觉沙盘和 3D 表现后置。
- 不要脱离 IAOS 现有架构另造系统。需要开发时，应优先复用 `/iaos/iaos-go` 的 Scenario、Capability、Event、Process、Policy、Decision 和 AI Tool 机制。
- 所有新增设计必须说明它服务于哪类虚拟企业对象、业务事件、仿真行为或 Agent 能力。
- 不要把 AESE 范围扩展到全行业、全模块或全 3D 世界。MVP 以华辰热管理系统集团、苏州制造基地、电池冷却板产品族为锚点。

## 进展更新要求

每次对本项目做出实质性进展后，必须更新 `docs/progress-log.md`。

实质性进展包括：

- 新增或修改项目定位、范围、业务对象、事件流、Agent 设计。
- 新增或修改开发计划、里程碑、任务拆分。
- 新增代码、数据、脚本、场景包或 IAOS 集成设计。
- 做出重要取舍，例如暂缓 3D、调整 MVP 范围、改变虚拟客户设定。
- 发现风险、约束、未决问题或与 IAOS 现有架构的冲突。

更新格式：

```text
## YYYY-MM-DD - 简短标题

- 变更：
- 原因：
- 影响：
- 后续：
```

如果只是阅读、检查或无实质变更，可以不更新进展日志。

