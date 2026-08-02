---
id: SOL-008
title: MiniMax 企业名称生成未启用
date: 2026-07-30
status: completed
author: Codex + User
tags: [aese, genesis, minimax, deployment, configuration]
---

# MiniMax 企业名称生成未启用

## 1. 现象

Enterprise Genesis 身份工作室可以生成候选企业名称，但
`GET /api/aese/v1/game/creative/status` 返回：

```json
{"state":"fallback","provider":"deterministic","model":"built-in"}
```

这表示请求没有调用 MiniMax，界面得到的是离线确定性候选。

同一根因也会使 M10 `GET /api/aese/v1/world/plant-build/planning-status`
返回 `not_configured / none`，随后“保存需求并让 Agent 生成候选”在 Requirement
已保存或校验完成后返回 503。M9 的 deterministic 命名 fallback 不能冒充 M10
规划 Agent；M10 没有外部模型时只能显式走人工候选表单。

## 2. 根因

本机权限为 `0600` 的 `.env` 已配置 `MINMAX_API_KEY`、
`MINMAX_API_BASE=https://api.minimax.chat/v1` 和 `MINMAX_MODEL=MiniMax-M3`，
但旧的手工启动命令没有加载 `.env`。AESE 只读取服务进程环境，因此安全地回退到了
`DeterministicProvider`。配置存在不等于运行进程已经使用配置。

## 3. 修正

新增 `scripts/deploy_aese_server.sh` 作为标准后端部署入口：

- 非执行式解析 `.env`，不会 `source` 其中的 shell 内容；
- 要求 secret 文件权限为 `0600`；
- MiniMax key、base 和 model 必须同时配置，否则失败关闭；
- 重建并重启 `:8090` 上的 AESE 服务，只停止命令行为 `aese-server` 的进程；
- 启动后同时校验 M9 Creative Provider 与 M10 Plant Planning Provider 状态；任一未连接都使发布失败；
- 密钥只进入子进程环境，不进入命令行、日志或 Git。

`.env.example` 提供无密钥模板。运行手册和 Code Map 同步记录该入口。

## 4. 验证证据

2026-07-30 使用有效 Genesis Workspace tenant session 经 `:4173` BFF 完成真实名称
生成：

- Provider：`MiniMax`
- Model：`MiniMax-M3`
- Host：`api.minimax.chat`
- 状态：`completed`
- 延迟：31,805 ms
- Token：prompt 470、completion 1,433、total 1,903
- 业务校验：`valid`
- Fallback：无
- 候选数量：4

CreativeJob 记录了 request ID、内容哈希和以上证据，没有保存 API key。直接访问
`:8090` 与通过 `:4173` 代理读取 Provider 状态均返回 `connected`。

## 5. 操作

```bash
cp .env.example .env
chmod 600 .env
# 填写三项 MINMAX_* 配置
scripts/deploy_aese_server.sh
curl -s http://127.0.0.1:8090/api/aese/v1/game/creative/status | jq
curl -s http://127.0.0.1:8090/api/aese/v1/world/plant-build/planning-status | jq
```

两个端点均应为 `connected / MiniMax / MiniMax-M3`。正式验收还必须分别在身份工作室
和 M10 工厂规划执行一次生成，并检查 CreativeJob/Agent Run，而不是只检查状态端点。
