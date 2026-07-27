---
id: SOL-002
title: AESE World 关联设立案失效时的可恢复加载
date: 2026-07-27
status: completed
author: Codex + User
tags: [aese, world, m9, recovery]
---

# 问题

从旧深链打开 `#world-incorporation` 时，链接中的 `case` 可能已经被清理或不属于当前租户。
IAOS trace 返回 404，AESE 将其当成整个 campaign 加载失败，导致本地 World 基线也无法查看。
默认 E2E 又固定使用已失效的 `INC-HCTM-001`，未能及时暴露数据生命周期变化。

# 根因

`loadIncorporation` 把所有非 2xx lifecycle 响应统一抛错，没有区分：

- 404：关联业务对象已不存在，可恢复；
- 401/403：身份或权限错误，必须失败关闭；
- 5xx：IAOS 服务错误，必须显式失败。

# 修复

- IAOS trace 为 404 时继续返回 AESE 本地 frame，并写入结构化
  `iaos_lifecycle_warning`。
- 查询 `/api/v1/incorporations/recent`，在页面提示中给出当前可用案件的切换入口。
- 页面明确说明本地 World 可查看，但正式进度仍必须来自有效 IAOS 案件。
- 401/403/5xx 继续阻止加载，避免掩盖授权或平台故障。
- M9 E2E 从 recent API 动态选择真实案件；World Hub 测试按当前链接文案进入，不依赖旧按钮。

# 验证

- Unit：404 → World frame 保留、warning 结构正确、recent case 可发现。
- E2E：无案件深链可打开本地 World；真实 recent case 可加载 IAOS projection。
- 三视口：1440、1280、390。

# 操作

若页面提示关联案件不存在：

1. 点击提示中的可用案件；
2. 或从 IAOS“企业成立与治理/企业设立案件”重新点击“打开 AESE World”；
3. 若提示 401/403，重新从 IAOS 打开以刷新 token，不应继续使用旧链接。
