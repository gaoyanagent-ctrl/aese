# M10 Genesis Plant Build Evidence

- 候选与决策：三个虚构候选先执行 1,500 万预算、2,200 kVA 和 2026-09-01 可用日期硬约束，再输出可解释评分；`SITE-SZ-NORTH-LEASED-SHELL` 是唯一可行方案。
- 确定性世界：10 帧从 M9 terminal eligibility 推进到 `capability_build_eligible=true`；100 次 canonical hash 一致，snapshot restore 和破损状态失败关闭。
- 设施与经济：七个 site/building/zone 节点、六项 WBS、公用工程/消防/EHS 验收；承诺 13,500,000、应付/已付 10,000,000、期末现金 10,000,000 CNY。
- 异常链：utility delay 形成 World 48% / IAOS 65% discrepancy，项目负责人仅在 observation 送达后获知并提交 rebaseline，committed outcome 后由 World 计算新进度，验收关闭差异。
- IAOS：独立 worktree `feat/m10-plant-governance` revision `23be02a`；六类权限 allowlist、FORCE RLS、预算/现金/验收/自批/mandate/版本/幂等门，治理记录、journal 与 Outbox 原子事务。真实 API 首次 201、重复 200 duplicate、未验收付款 422。
- UI：Plant Build Play 显示候选对比、虚拟时间、World/IAOS 进度、Knowledge、资金和空间图；390/1280/1440 三视口验收覆盖。
- 回归：AESE 全量 Go test/vet、Schema/pack validate、34 项 Vitest、TypeScript 与生产构建通过；IAOS 四个 Go module test/vet 通过。
