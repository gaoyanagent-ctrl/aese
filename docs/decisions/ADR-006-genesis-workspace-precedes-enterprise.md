---
id: ADR-006
title: Genesis 创业空间先于企业法律主体
date: 2026-07-28
status: accepted
author: Codex + User
tags: [genesis, tenant, onboarding, identity, isolation]
---

# Genesis 创业空间先于企业法律主体

## 背景

现有 Enterprise Genesis 以 `tenant-hctm-genesis` 和已知 `founder-principal`
登录态为前置条件。玩家可以创建新的 incorporation case，但所有公司仍共享同一租户、
同一创始人账号和同一套基线数据。这只能验证案件隔离，不能验证 IAOS 的租户创建、
身份、RLS、Runtime 安装和企业全生命周期隔离。

租户也不能直接等同于已经成立的公司：在玩家选择名称、登记法人之前，IAOS 已经需要
一个隔离空间保存创业项目、草稿、模型调用和设立流程。

## 决策

1. 玩家先拥有平台级 `PlayerAccount`，其身份不依赖任何企业租户。
2. 点击“创建新企业”先创建 `GenesisWorkspace`，再由 IAOS 控制平面生成独立 tenant。
3. tenant 是稳定的安全与数据隔离边界，不是法律主体；名称选择不会改变 tenant ID。
4. `legal_entity` 只在 `incorporation.case.open` 之后进入 IAOS Business State。
5. 浏览器不得持有 `platform.manage` 或 `platform.identity.bootstrap` 权限。IAOS 提供
   受限的 Genesis provisioning API，在服务端执行租户创建、身份成员关系、Runtime
   安装、RLS smoke check 和激活。
6. AESE 创建独立 World Run，并只保存 tenant/workspace/case 的稳定引用，不写 IAOS DB。
7. 同一玩家可以拥有多个 workspace；每个 workspace 对应一个 tenant。一个 tenant
   首版只允许一个 Genesis 根企业，后续集团法人仍属于同一 tenant。

## 创建顺序

```text
PlayerAccount
-> GenesisWorkspace(requested)
-> tenant_account(provisioning)
-> founder membership + tenant-scoped credential
-> M9 semantic/entity/capability/process/policy/agent runtime install
-> RLS + login + runtime smoke checks
-> tenant_account(active)
-> AESE World Run
-> FounderIntent
-> AI NamingProposal
-> incorporation.case.open
-> legal_entity
```

## 失败与恢复

- provisioning 每一步保存 checkpoint 和幂等键，刷新后继续同一 workspace。
- 失败 tenant 保持 `provisioning`；不发放可写业务 token，不静默改成 active。
- Genesis provisioning record 保存失败阶段、可重试性和补偿结果；不为
  `tenant_account` 发明现有状态机之外的新枚举。
- 已激活 tenant 不因后续 AI 生成失败而删除；玩家可使用 fallback 或稍后重试。
- 租户删除/归档不属于创建向导，必须进入 IAOS SaaS Ops 的高风险治理路径。

## 影响

- `tenant-hctm-genesis` 降级为验收 fixture，不再作为新游戏默认租户。
- URL 的 `tenant`、`case` 参数只用于深链和恢复，不再是进入游戏的前置知识。
- 需要 IAOS 新增 self-service Genesis provisioning contract；不能把现有平台管理员
  API 原样暴露给前端。

