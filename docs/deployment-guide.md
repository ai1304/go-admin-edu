# 特殊教育资源库部署指南

更新时间：2026-05-16

本文档记录当前项目的本地启动、Docker 部署、远程服务器运行态部署，以及后续代码更新流程。当前项目包含：

- 后端：`go-admin/`
- 前端工作区：`web/`
- 管理端前端：`web/apps/admin`
- 门户端前端：`web/apps/portal`
- 运行态部署配置：`deploy/docker-compose.runtime.yml`

## 1. 当前推荐部署方式

当前远程服务器推荐使用“本地构建产物，服务器只运行”的方式。

原因：

- 远程服务器配置较小，直接在服务器上构建 Go / Node 镜像容易卡在依赖下载或编译阶段。
- 本地已经能稳定完成 Go 编译和前端构建。
- 服务器只运行 MySQL、Redis、MinIO、API 二进制和 Nginx 静态前端，部署更快、更可控。

当前已验证的远程服务器：

- 服务器 IP：`117.72.200.80`
- 管理端端口：`18080`
- 门户端端口：`18081`
- API 端口：`8000`
- MinIO API：`9000`
- MinIO Console：`9001`

安全组建议放行：

- `18080`
- `18081`
- `8000`
- `9000`
- `9001`

不要对公网放行 `3306` 和 `6379`，除非明确需要外部直连 MySQL / Redis。

## 2. 端口约定

本地开发端口：

- 后端 API：`8000`
- 管理端前端：`1798`
- 门户端前端：`1799`

Docker / 远程部署端口：

- 后端 API：`8000`
- 管理端前端：`18080`
- 门户端前端：`18081`
- MySQL：`3306`
- Redis：`6379`
- MinIO API：`9000`
- MinIO Console：`9001`

## 3. 本地开发启动

### 3.1 后端

```bash
cd go-admin
go mod tidy
go build ./...
go run main.go migrate -c config/settings.yml
go run main.go server -c config/settings.yml
```

默认访问：

```text
http://127.0.0.1:8000
```

### 3.2 管理端

```bash
cd web
pnpm install
pnpm dev:admin
```

访问：

```text
http://127.0.0.1:1798
```

### 3.3 门户端

```bash
cd web
pnpm dev:portal
```

访问：

```text
http://127.0.0.1:1799
```

## 4. 本地 Docker Compose 启动

根目录提供了 `docker-compose.yml`，可用于本地 Docker 联调。

```bash
docker compose up -d mysql redis minio
docker compose run --rm api /main migrate -c /config/settings.yml
docker compose up -d api admin-web portal-web
```

访问：

- API：`http://127.0.0.1:8000`
- 管理端：`http://127.0.0.1:18080`
- 门户端：`http://127.0.0.1:18081`
- MinIO Console：`http://127.0.0.1:9001`

MinIO 默认账号：

- 用户名：`minioadmin`
- 密码：`minioadmin`

## 5. 远程服务器首次部署

### 5.1 服务器准备

服务器需要安装：

- Docker
- Docker Compose v2
- Git

如果服务器内存较小，建议增加 2G swap，避免 MySQL / 构建 / 解压时内存紧张。

```bash
fallocate -l 2G /swapfile || dd if=/dev/zero of=/swapfile bs=1M count=2048
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
echo '/swapfile none swap sw 0 0' >> /etc/fstab
```

### 5.2 拉取代码

```bash
mkdir -p /opt
git clone https://github.com/ai1304/go-admin-edu.git /opt/go-admin-edu
cd /opt/go-admin-edu
```

如果仓库已经存在：

```bash
cd /opt/go-admin-edu
git fetch origin main
git reset --hard origin/main
```

### 5.3 配置 MinIO 公网地址

生产或远程部署时，必须把 `publicEndpoint` 改成浏览器可访问的地址。

```bash
sed -i 's#publicEndpoint: http://localhost:9000#publicEndpoint: http://117.72.200.80:9000#' go-admin/config/settings.docker.yml
```

说明：

- `endpoint: minio:9000` 是容器内部访问 MinIO 的地址。
- `publicEndpoint` 是后端生成文件访问 URL 时使用的公网地址。
- MinIO 预签名 URL 会把 Host 纳入签名计算，不能先用容器内部地址签名再替换成公网地址，否则门户端图片可能报 `SignatureDoesNotMatch`。

### 5.4 本地构建运行态产物

在本地项目根目录执行。

构建 Linux API 二进制：

```powershell
$env:GOOS='linux'
$env:GOARCH='amd64'
$env:CGO_ENABLED='0'
go build -trimpath -ldflags='-s -w' -o D:\project\go-admin-edu\deploy\artifacts\api\go-admin .\go-admin
```

构建管理端和门户端：

```powershell
npm.cmd --prefix web/apps/admin run build
npm.cmd --prefix web/apps/portal run build
```

复制前端 dist 到运行态目录：

```powershell
New-Item -ItemType Directory -Force deploy\artifacts\admin | Out-Null
New-Item -ItemType Directory -Force deploy\artifacts\portal | Out-Null
Copy-Item -Path web\apps\admin\dist\* -Destination deploy\artifacts\admin -Recurse -Force
Copy-Item -Path web\apps\portal\dist\* -Destination deploy\artifacts\portal -Recurse -Force
```

打包产物：

```powershell
tar -czf deploy-artifacts.tgz -C deploy artifacts
```

### 5.5 上传产物

示例：

```bash
scp deploy-artifacts.tgz root@117.72.200.80:/opt/go-admin-edu/deploy-artifacts.tgz
```

服务器解压：

```bash
cd /opt/go-admin-edu
rm -rf deploy/artifacts
tar -xzf deploy-artifacts.tgz -C deploy
chmod +x deploy/artifacts/api/go-admin
```

### 5.6 启动服务

```bash
cd /opt/go-admin-edu
docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml up -d mysql redis minio
```

首次部署或有新增迁移时执行：

```bash
docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml run --rm api sh -c 'mkdir -p temp/logs && /main migrate -c /config/settings.yml'
```

启动 API、管理端、门户端：

```bash
docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml up -d api admin-web portal-web
```

查看状态：

```bash
docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml ps
```

### 5.7 健康检查

服务器本机检查：

```bash
curl -fsS http://127.0.0.1:9000/minio/health/live
curl -fsS http://127.0.0.1:8000/api/v1/app-config
curl -fsSI http://127.0.0.1:18080
curl -fsSI http://127.0.0.1:18081
```

公网检查：

```text
http://117.72.200.80:18080
http://117.72.200.80:18081
http://117.72.200.80:8000/api/v1/app-config
http://117.72.200.80:9001
```

如果服务器本机检查通过，但公网访问超时，优先检查云服务器安全组。

## 6. 后续更新代码怎么做

后续更新时，不需要每次重装 Docker / MySQL / MinIO。按变更范围选择操作即可。

### 6.1 后端代码有变化

本地重新构建 Linux 二进制：

```powershell
$env:GOOS='linux'
$env:GOARCH='amd64'
$env:CGO_ENABLED='0'
go build -trimpath -ldflags='-s -w' -o D:\project\go-admin-edu\deploy\artifacts\api\go-admin .\go-admin
```

上传并替换服务器文件：

```bash
scp deploy/artifacts/api/go-admin root@117.72.200.80:/opt/go-admin-edu/deploy/artifacts/api/go-admin
ssh root@117.72.200.80 'chmod +x /opt/go-admin-edu/deploy/artifacts/api/go-admin'
```

如果后端新增了迁移文件，先执行迁移：

```bash
ssh root@117.72.200.80 "cd /opt/go-admin-edu && docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml run --rm api sh -c 'mkdir -p temp/logs && /main migrate -c /config/settings.yml'"
```

重启 API：

```bash
ssh root@117.72.200.80 "cd /opt/go-admin-edu && docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml up -d api"
```

### 6.2 管理端前端有变化

本地构建：

```powershell
npm.cmd --prefix web/apps/admin run build
```

上传新的 dist：

```bash
scp -r web/apps/admin/dist/* root@117.72.200.80:/opt/go-admin-edu/deploy/artifacts/admin/
```

重启管理端 Nginx：

```bash
ssh root@117.72.200.80 "cd /opt/go-admin-edu && docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml restart admin-web"
```

### 6.3 门户端前端有变化

本地构建：

```powershell
npm.cmd --prefix web/apps/portal run build
```

上传新的 dist：

```bash
scp -r web/apps/portal/dist/* root@117.72.200.80:/opt/go-admin-edu/deploy/artifacts/portal/
```

重启门户端 Nginx：

```bash
ssh root@117.72.200.80 "cd /opt/go-admin-edu && docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml restart portal-web"
```

### 6.4 配置文件或部署文件有变化

先推送代码，再在服务器拉取：

```bash
git push origin main
ssh root@117.72.200.80 "cd /opt/go-admin-edu && git fetch origin main && git reset --hard origin/main"
```

重新设置 MinIO 公网地址：

```bash
ssh root@117.72.200.80 "cd /opt/go-admin-edu && sed -i 's#publicEndpoint: http://localhost:9000#publicEndpoint: http://117.72.200.80:9000#' go-admin/config/settings.docker.yml"
```

重启相关服务：

```bash
ssh root@117.72.200.80 "cd /opt/go-admin-edu && docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml up -d"
```

### 6.5 全量更新

如果后端、管理端、门户端都有变化，推荐重新打包并上传：

```powershell
tar -czf deploy-artifacts.tgz -C deploy artifacts
scp deploy-artifacts.tgz root@117.72.200.80:/opt/go-admin-edu/deploy-artifacts.tgz
```

服务器执行：

```bash
cd /opt/go-admin-edu
git fetch origin main
git reset --hard origin/main
sed -i 's#publicEndpoint: http://localhost:9000#publicEndpoint: http://117.72.200.80:9000#' go-admin/config/settings.docker.yml
rm -rf deploy/artifacts
tar -xzf deploy-artifacts.tgz -C deploy
chmod +x deploy/artifacts/api/go-admin
docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml run --rm api sh -c 'mkdir -p temp/logs && /main migrate -c /config/settings.yml'
docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml up -d api admin-web portal-web
```

## 7. 默认账号

管理端默认账号：

- 用户名：`admin`
- 密码：`123456`

如果登录失败，以数据库 `sys_user` 表中的实际数据为准。

## 8. 常见问题

### 8.1 服务器内部正常，公网访问超时

优先检查云服务器安全组是否放行：

- `18080`
- `18081`
- `8000`
- `9000`
- `9001`

服务器内部可用以下命令确认服务是否监听：

```bash
ss -lntp | grep -E ':8000|:18080|:18081|:9000|:9001'
```

### 8.2 门户端图片显示 `SignatureDoesNotMatch`

检查 `go-admin/config/settings.docker.yml`：

```yml
extend:
  storage:
    endpoint: minio:9000
    publicEndpoint: http://117.72.200.80:9000
```

`publicEndpoint` 必须是浏览器访问图片时使用的地址。

### 8.3 迁移找不到 `config/db-begin-mysql.sql`

运行态 compose 必须满足：

- API 容器 `working_dir: /app`
- 挂载 `../go-admin/config:/app/config:ro`

当前 `deploy/docker-compose.runtime.yml` 已包含这些配置。

### 8.4 单个远程命令超过 1 分钟没有进度

建议中断后分步执行：

1. `docker compose ps`
2. `docker ps -a`
3. `docker logs go-admin-edu-api --tail=100`
4. 只重启单个服务，例如 `docker compose ... restart api`

不要在小内存服务器上反复执行完整镜像构建，优先使用本文的运行态产物部署流程。
