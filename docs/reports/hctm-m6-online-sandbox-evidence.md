# HCTM M6 在线沙盘验收证据

日期：2026-07-20

本地完整链路从 reset 开始：reset 删除 9 个 L2、保留 12 个 L1；apply 为 9 insert / 12 no-op；首次 replay 为 10 triggered / 12 skipped / 0 failed；稳定状态重复 replay 为 0 triggered / 22 skipped / 0 failed。

在线 snapshot 实测需求 12,000、累计可供 11,700、累计实发 11,700、期末成品 0、交付缺口 300，包含 14 个对象、6 个受治理业务事件和 3 条 Agent 建议。唯一 gap 为 `cost_actuals`。重复 Agent run 执行 9 次受治理 Tool Call，在相同 Agent/correlation scope 更新 3 条建议，不生成重复建议。

游标查询严格升序并返回 `next_cursor`/`has_more`；`tenant-other` 查询返回 404。浏览器 E2E 覆盖 Preview、Live、1440×900、1280×720、390×844；真实 IAOS 链路截图为 `frontend/test-results/actual-live-*.png`。所有建议保持 `suggested` 且要求人工确认；没有写入虚构成本金额。
