---
id: DES-026
title: M24 场景平台产品化与行业交付
date: 2026-07-22
status: completed
author: Codex + User
tags: [m24, platform, authoring, release]
---

# M24 场景平台产品化与行业交付

## 1. 目标

把 M3-M23 的 HCTM 工程资产整理成可复制交付的 AESE 3.0 场景平台：有 authoring、合同验证、兼容、benchmark、发布、部署和运维门，而不是只能由原作者维护的单一 demo。

机器终态：`industry_simulation_platform_ready=true`。

## 2. 首版产品面

- Scenario/World Pack SDK：manifest、schema、stable code、rules、actors、campaign、experiment、expected evidence。
- Authoring Studio：模板创建、引用检查、时间线/地图预览、diff、lint 和 dry-run；不允许任意代码上传。
- Certification Pipeline：determinism、invariant、tenant、security、performance、accessibility 和 compatibility benchmark。
- Release Registry：pack/rules/schema/API/IAOS compatibility、签名/hash、provenance、promotion/rollback 和 deprecation。
- Operator Console：部署目标、preflight、apply、health、quota、artifact retention、backup/restore 和 evidence export。

## 3. HCTM Reference Release

HCTM pack 必须包含 M9-M23 可独立运行的 campaign、完整 terminal ancestry、三态说明、runbook、sample data、benchmark 和已知限制。所有主体保持虚构，任何生产凭据、客户数据和不可再分发资产禁止进入发布物。

## 4. 扩展合同

插件点只允许声明式 entity/event/rule/actor/view schema 和受控 Go extension interface；扩展必须明确 owner、权限、资源配额、确定性级别和 sandbox。未知扩展、任意脚本、网络访问和直接 IAOS DB/NATS 写入失败关闭。

## 5. 发布与运维

- dev/staging/demo target 明确；真实 production 不在默认 allowlist。
- 版本升级有 impact report、migration、rollback、兼容窗口和旧 pack regression。
- World Store、IAOS dependency、NATS、artifact、Atlas 和前端有健康/容量/恢复 runbook。
- release bundle 和 evidence 可离线验证；默认命令 dry-run，危险动作显式确认。

## 6. 完成标准

- 新维护者可从模板创建一个最小虚构场景，validate、run、replay、UI 预览并生成 evidence。
- HCTM reference release 通过全量 certification 和 clean-environment install/upgrade/rollback。
- API/schema/pack compatibility matrix、SBOM/依赖、权限、tenant、性能和无障碍门完整。
- 文档、Code Map、Atlas、runbook、release notes 和支持边界可交付。
- M3-M23 全链回归通过并输出 `industry_simulation_platform_ready=true`。

## 7. Program 结束边界

M24 关闭 AESE 3.0 program。真实生产接入、第二行业、高保真 3D/数字孪生、法定合规或托管 SaaS 必须另立新 program 和安全/商业决策，不在 M24 中隐式承诺。
