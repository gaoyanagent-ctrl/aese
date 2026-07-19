# HCTM M3V 2D 企业沙盘验收证据

日期：2026-07-19  
范围：PLAN-M3V-001 / DES-002

## 1. 实现结果

- React + TypeScript + Vite 只读工作台完成。
- `preview.json` 映射 canonical 七幕、22 个事件、14 个节点、13 条连线、14 个可检查实体和三类 Agent 输出。
- 纯 reducer 支持确定性重放、前后单步、跳转、播放/暂停、1×/2×/4×和重置。
- 画布、事件流、对象详情、五项 KPI 和 Agent 建议随事件同步变化。
- 桌面固定为单视口工作台；事件面板内部滚动。移动端使用三个标签切换主区域，页面无横向溢出。

## 2. 自动化结果

| 检查 | 结果 |
| --- | --- |
| TypeScript | 通过 |
| ESLint | 通过，0 warning |
| Vitest | 5 files / 18 tests 通过 |
| Vite production build | 通过 |
| Playwright | 3 projects / 9 tests 通过 |
| npm audit | 0 vulnerabilities |
| Go `test ./...` / `vet ./...` | 通过 |

播放内核的真实 fixture 测试确认最终累计实发 11,700、第二批实发 2,700、短缺 300，且交付风险为严重。

验收时开发服务已绑定 `0.0.0.0:4173`，`http://127.0.0.1:4173/` 返回 HTTP 200；同网段访问地址为 `http://192.168.50.222:4173/`。

## 3. 视觉证据

- [1440×900 完成态](../../frontend/test-results/desktop-1440-completed.png)
- [1280×720 完成态](../../frontend/test-results/desktop-1280-completed.png)
- [390×844 完成态](../../frontend/test-results/mobile-390-completed.png)

三个截图尺寸经过文件元数据核对。画布非空，控制栏、七幕进度和 KPI 没有相互覆盖；右侧长事件流在桌面端内部滚动；移动端 KPI 在自身区域滚动，不造成整页横向溢出。

## 4. 数据源和架构边界

React 组件不直接读取 IAOS DTO。`StaticScenarioDataSource` 负责加载并校验 `SandboxScenario`，组件只消费内部视图模型。浏览器仅应用预计算 delta，不实现第二套 MRP、排产、质量或 Agent Runtime。

## 5. 后续

M3V 已达到完成定义。下一步是 M4：通过受治理 simulation ingress 将供应商延期、设备停机和来料不良接入 IAOS；之后用 `IaosScenarioDataSource` 替换静态数据源。
