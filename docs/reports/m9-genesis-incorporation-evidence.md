# M9 Genesis Incorporation Evidence

- 业务基线：投资人期初 50,000,000 CNY；认缴 30,000,000；首期实缴 20,000,000；登记费 10,000；公司期末现金 20,000,000；首年预算授权 15,000,000。
- 确定性：8 frame、100 次 hash 一致、snapshot restore/reset 与非法状态失败测试。
- 资格：终态同时满足法人 registered、账户 open、资本 received、三岗位 accepted、mandate active、预算 approved，输出 `plant_project_eligible=true`。
- IAOS：独立 `feat/m9-genesis-governance` worktree revision `edcb915` 提供 allowlist governance endpoint；FORCE RLS、幂等、权限、自批/过期 mandate 拒绝、记录与 Outbox 原子事务。
- UI：桌面 1440/1280、移动 390 均可从初态推进到 M10 eligible；预算明确不等同现金。
- 兼容：M7 全量 Playwright 和 M8 World live 回归保留。
