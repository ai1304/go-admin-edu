# 特殊教育资源库平台启动与部署说明

更新时间：2026-05-13

## 1. 适用范围

本文档用于指导当前仓库的本地开发启动、测试环境初始化以及后续部署准备。

当前仓库结构：

- 后端：`go-admin/`
- 前端工作区：`web/`
- 后台前端：`web/apps/admin`
- 门户前端：`web/apps/portal`

## 2. 当前部署结论

当前项目已经具备本地启动和数据库初始化能力，但还不是完整的一键部署形态。

已具备：

- 后端 `go-admin` 可通过配置文件连接 MySQL。
- 后端具备迁移命令，可初始化基础表和教育业务表。
- 后台前端和门户前端已统一到 `web/` workspace。
- 两个前端均已配置反向代理到本地后端 `http://127.0.0.1:8000`。

当前限制：

- `go-admin/docker-compose.yml` 目前只包含后端容器，不包含 MySQL、Redis、MinIO、前端服务。
- 生产部署前仍需补齐完整 Compose 或其他部署脚本。

## 3. 环境要求

建议环境：

- Go
- Node.js
- pnpm 9
- MySQL 8.x

说明：

- `web/package.json` 指定 `packageManager` 为 `pnpm@9.0.0`。
- 后端默认数据库类型为 MySQL。
- 当前仓库默认本地开发端口为 `8000`、`1798`、`1799`。

## 4. 端口约定

- 后端 API：`8000`
- 后台前端：admin：`1798`
- 门户前端：portal：`1799`

当前前端代理关系：

- `web/apps/admin/vite.config.js` 将 `/api/v1` 代理到 `http://127.0.0.1:8000`
- `web/apps/portal/vite.config.js` 将 `/api/v1` 代理到 `http://127.0.0.1:8000`

## 5. 数据库准备

### 5.1 创建数据库

先在 MySQL 中创建一个空库，例如：

```sql
CREATE DATABASE go_admin_edu DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5.2 修改后端配置

编辑文件：

- `go-admin/config/settings.yml`

将数据库连接修改为你自己的环境，例如：

```yml
settings:
  database:
    driver: mysql
    source: root:123456@tcp(127.0.0.1:3306)/go_admin_edu?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms
```

注意：

- 仓库中当前没有 `config/settings.dev.yml`。
- 当前实际可直接使用的是 `config/settings.yml`。
- 如需为本地、测试、生产拆分配置，建议后续新增 `settings.local.yml`、`settings.test.yml`、`settings.prod.yml`。

## 6. 后端启动流程

在仓库根目录执行：

```bash
cd go-admin
go mod tidy
go build
```

### 6.1 初始化数据库

首次启动前执行迁移：

```bash
./go-admin migrate -c config/settings.yml
```

如果你不想先编译，也可以直接执行：

```bash
go run main.go migrate -c config/settings.yml
```

迁移会完成以下工作：

- 初始化 `sys_migration` 迁移记录表
- 初始化 go-admin 基础表结构
- 导入基础 SQL 数据
- 执行教育业务模块迁移

教育业务迁移文件位于：

- `go-admin/cmd/migrate/migration/version/2026051200010_edu_tables.go`

### 6.2 启动后端服务

```bash
./go-admin server -c config/settings.yml
```

或：

```bash
go run main.go server -c config/settings.yml
```

启动成功后，后端默认访问地址：

```text
http://127.0.0.1:8000
```

## 7. 前端启动流程

在仓库根目录执行：

```bash
cd web
pnpm install
```

### 7.1 启动后台前端

```bash
pnpm dev:admin
```

等价命令：

```bash
pnpm --dir apps/admin dev
```

访问地址：

```text
http://127.0.0.1:1798
```

### 7.2 启动门户前端

```bash
pnpm dev:portal
```

等价命令：

```bash
pnpm --dir apps/portal dev
```

访问地址：

```text
http://127.0.0.1:1799
```

## 8. 数据库初始化脚本说明

当前项目的数据库初始化不是单纯依赖一个 `.sql` 文件，而是由迁移命令统一调度。

相关文件如下：

基础 SQL：

- `go-admin/config/db.sql`
- `go-admin/config/db-begin-mysql.sql`
- `go-admin/config/db-end-mysql.sql`
- `go-admin/config/pg.sql`
- `go-admin/config/db-sqlserver.sql`

迁移入口：

- `go-admin/cmd/migrate/server.go`
- `go-admin/cmd/migrate/migration/init.go`

教育业务迁移：

- `go-admin/cmd/migrate/migration/version/2026051200010_edu_tables.go`

因此推荐的初始化方式始终是：

```bash
./go-admin migrate -c config/settings.yml
```

而不是手工逐个执行 SQL 文件。

## 9. 默认账号

基础 SQL 中已包含管理员账号初始化数据。

默认可尝试使用：

- 用户名：`admin`
- 密码：`123456`

说明：

- 该账号来自当前仓库内置初始化数据和 go-admin 默认演示账号约定。
- 若后续初始化脚本被修改，请以实际数据库内容为准。

## 10. 本地启动顺序建议

推荐按以下顺序启动：

1. 启动 MySQL，并创建数据库。
2. 修改 `go-admin/config/settings.yml` 中的数据库连接。
3. 执行 `go-admin migrate -c config/settings.yml` 完成初始化。
4. 启动后端 `go-admin server -c config/settings.yml`。
5. 在 `web/` 下执行 `pnpm install`。
6. 启动后台前端 `pnpm dev:admin`。
7. 启动门户前端 `pnpm dev:portal`。

## 11. 启动验证

建议最少做以下验证：

### 11.1 后端验证

- 访问 `http://127.0.0.1:8000`
- 观察后端控制台是否报数据库连接错误
- 确认迁移后数据库已生成基础表和教育业务表

### 11.2 后台前端验证

- 打开 `http://127.0.0.1:1798`
- 尝试使用 `admin / 123456` 登录
- 确认接口请求已成功代理到后端

### 11.3 门户前端验证

- 打开 `http://127.0.0.1:1799`
- 确认资源列表和详情页接口能够正常返回

## 12. 常见问题

### 12.1 `settings.dev.yml` 不存在

原因：

- 当前仓库并没有这个文件，历史文档里提到的是上游项目常见命名。

处理方式：

- 直接使用 `config/settings.yml`

### 12.2 Docker Compose 不能直接起完整环境

原因：

- 当前 `go-admin/docker-compose.yml` 只定义了后端容器。

处理方式：

- 本地开发请先手动启动 MySQL。
- 后续需要补齐 MySQL、Redis、MinIO、前端服务的 Compose 配置。

### 12.3 迁移成功但前端仍无法访问数据

优先检查：

- 后端是否已启动在 `8000`
- 前端代理目标是否仍为 `http://127.0.0.1:8000`
- 浏览器接口请求是否报 401、404、500
- 数据库连接串中的库名、账号、密码是否正确

## 13. 后续部署建议

为了让该项目具备可重复部署能力，建议下一步补齐以下内容：

1. 根目录统一 `docker-compose.yml`，包含 MySQL、Redis、MinIO、后端、admin、portal。
2. 增加环境拆分配置：本地、测试、生产。
3. 为前端增加 `.env` 配置，而不是将后端地址完全写死在 Vite 配置中。
4. 为数据库初始化增加部署脚本或 Makefile 命令，减少人工步骤。
5. 增加反向代理配置示例，例如 Nginx。

