# AESE World Store 本地运行与备份合同

World Store 是独立 PostgreSQL database，固定 database `aese_world`、应用账号 `aese_world_app`、本地端口 `55432` 和连接变量 `AESE_WORLD_DATABASE_URL`。禁止使用 IAOS database/账号、跨库查询、跨库外键或保存 IAOS JWT。生产凭据必须由部署环境注入；仓库密码仅供隔离的本地容器。

所有应用事务必须先执行参数化的 `SET LOCAL app.tenant_id = ...`；数据表启用并强制 RLS，未绑定 tenant 的请求看不到也不能写入行。应用账号不拥有建库、迁移或绕过 RLS 权限，迁移由独立管理身份执行。

本地启动与迁移：

```bash
docker compose -f deploy/world-postgres/docker-compose.yml up -d
docker compose -f deploy/world-postgres/docker-compose.yml exec -T world-postgres psql -U postgres -d aese_world < deploy/world-postgres/migrations/000001_world_contract.up.sql
export AESE_WORLD_DATABASE_URL='postgres://aese_world_app:aese_world_local_only@127.0.0.1:55432/aese_world?sslmode=disable'
```

迁移按数字版本严格顺序执行，每版在同一事务中完成并写 `world.schema_migrations`；失败不得留下部分 schema。down migration 只供一次性本地环境，正式环境通过向前修复迁移恢复。

备份必须使用一致性自包含 custom dump，并在独立临时 database 演练恢复后才算有效：

```bash
pg_dump --format=custom --no-owner --no-acl "$AESE_WORLD_DATABASE_URL" > aese-world.dump
createdb aese_world_restore_test
pg_restore --exit-on-error --no-owner --no-acl --dbname=aese_world_restore_test aese-world.dump
```

备份包含 World run、事件、快照、结构化认知、差异和 bridge cursor，不包含 IAOS 业务表或凭据。恢复目标不得是 IAOS database。生产备份要求加密、访问控制、校验和、保留策略与定期恢复演练；这些由部署环境负责，不能把 dump 提交到仓库。
