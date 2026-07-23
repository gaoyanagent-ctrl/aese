---
id: SOL-001
title: M9 局域网生命周期加载、双向链接与 SSE 截断
date: 2026-07-23
status: completed
author: Codex + User
tags: [aese, m9, lan, sse]
---

# M9 局域网生命周期加载、双向链接与 SSE 截断

## 症状

- 工作台默认请求不存在的 `INC-HCTM-001`，返回 404。
- IAOS/AESE 深链固定指向 localhost/127.0.0.1，远程浏览器打开了错误主机。
- NATS EventSource 每 60 秒收到 `ERR_INCOMPLETE_CHUNKED_ENCODING`。
- Vite 开发服务器在页面进入 BFCache 时报告 HMR WebSocket 关闭。
- 浏览器残留的 `aese_iaos_base_url` 指向 AESE origin 时，生命周期请求错误地
  发往 `:4173/api/v1/incorporations/...` 并返回 404；favicon 也返回 404。
- IAOS 与 AESE 分属 `:3000`、`:4173` 两个 origin，不能共享 localStorage；
  AESE 复用旧租户 JWT 时，RLS 将存在的 case 安全地投影为 404。

## 根因与修复

1. 默认 case 是演示占位符而非持久事实。新增 tenant-scoped
   `GET /api/v1/incorporations/recent`，空输入时加载最近可见 case。
2. 双向链接和 AESE IAOS API fallback 改用 `window.location.hostname`。
3. HTTP server 的 60 秒 `WriteTimeout` 不适合无限 SSE。handler 通过
   `http.ResponseController.SetWriteDeadline` 在握手、消息和每个 15 秒
   heartbeat 时续期。
4. Vite BFCache WebSocket 是开发 HMR 生命周期提示；业务投影和生产构建
   不依赖该连接。
5. IAOS base resolver 拒绝与 AESE 页面相同 origin 的陈旧配置，自动回退到
   当前浏览器 hostname 的 `:8082`；同时提供静态 favicon。
6. IAOS→AESE 深链通过 URL fragment 交接当前 JWT；AESE 接收后立即写入自身
   origin 存储并从地址栏移除 token。服务端仍按 JWT tenant 执行 RLS。

## 回归证据

- `INC-HCTM-001` 可稳定复现 404，clean case
  `INC-E2E-1784787213558168936` 返回 200。
- 以 `http://192.168.50.222:3000` 和 `:4173` 执行双仓三视口
  Playwright，各 3/3 通过；深链 host 保持 `192.168.50.222`。
- AESE Playwright 显式注入 `aese_iaos_base_url=window.location.origin`，
  三个视口仍全部通过，并断言 lifecycle 请求端口全部为 `8082`。
- AESE Playwright 先注入另一租户的陈旧 token，再通过深链接收正确 token；
  首次加载和刷新后三个视口均通过，且 URL 中不再保留 `auth_token`。
- SSE curl 持续 70 秒后由客户端 `--max-time` 主动终止，接收 144 bytes；
  不再由服务端在 60 秒以 incomplete chunked encoding 截断。
