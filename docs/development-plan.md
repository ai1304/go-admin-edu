# 特殊教育资源库平台开发计划

更新时间：2026-05-15

## 1. 当前技术基线

后端：

- 项目：`go-admin`
- 框架：Go + Gin
- ORM：GORM
- 权限：JWT + Casbin RBAC
- 数据权限：基于部门、角色、创建人，可扩展为区域/学校/租户数据权限
- 数据库：MySQL
- 部署：Docker

后台前端：

- 项目：`web/apps/admin`
- 当前分支：`vue3_dev`
- 技术：Vue3 + Vite + Pinia + Vue Router 4 + Arco Design

门户前端：

- 项目：`web/apps/portal`
- 技术：Vue3 + Vite + Pinia + Vue Router 4 + Arco Design
- 与后台前端共享接口类型、请求封装和部分业务组件
- 不考虑 SEO，不使用 Nuxt

前端共享包：

- 项目：`web/packages/shared`
- 用途：跨前台/后台共享常量、类型、请求基础配置和通用工具。

文件存储：

- 开发/测试：MinIO
- 生产：阿里云 OSS
- 后端必须封装统一 Storage 接口，业务代码不直接依赖 MinIO 或 OSS SDK

搜索：

- 一期：MySQL 索引与条件筛选
- 后续：Elasticsearch 或 OpenSearch

缓存：

- Redis

## 2. 总体开发原则

- 先完成管理后台和核心资源流转，再做门户体验 polish。
- 先做字段级多租户和数据权限预留，不完全依赖 go-admin 的按域名/多库租户模式。
- 所有业务表默认预留 `tenant_id`、`region_id`、`school_id`、`created_by`、`updated_by`、`created_at`、`updated_at`、`deleted_at`。
- 特教案例、IEP、学生评估、干预方案按敏感数据处理，不能套用普通资源公开逻辑。
- 文件只存对象存储，MySQL 只存文件元数据和对象 Key。
- 搜索、存储、预览、统计都要留接口抽象，方便二期替换或增强。

## 3. 里程碑规划

## M0：项目基线整理

目标：让两个已克隆项目可以稳定启动，并形成可扩展工程规范。

任务：

- 确认 `go-admin` 后端可连接 MySQL 并完成初始化迁移。
- 确认 `go-admin-ui` 当前 `vue3_dev` 分支可安装依赖、启动、登录。
- 补充 Docker Compose：MySQL、Redis、MinIO、后端、后台前端。
- 整理本地开发环境变量。
- 梳理 go-admin 现有用户、角色、菜单、部门、数据权限表结构。
- 确认代码生成器是否适合生成业务 CRUD。

交付物：

- 可启动的本地开发环境。
- 基础 Docker Compose。
- 初始化账号。
- 开发说明文档。

## M1：组织、租户与权限改造

目标：支撑区域、学校、教师、学生的数据隔离。

后端任务：

- 新增或改造租户模型：`sys_tenant`。
- 新增区域表：`edu_region`。
- 新增学校表：`edu_school`。
- 明确 `sys_dept` 与区域/学校的关系，是继续作为组织树使用，还是学校内部门使用。
- 扩展用户表或用户档案表，支持用户绑定 `tenant_id`、`region_id`、`school_id`、用户类型。
- 扩展数据权限中间件，支持按区域、学校、本人、授权范围过滤。
- 定义角色：超级管理员、区域管理员、学校管理员、教师、学生。

前端任务：

- 后台增加区域管理。
- 后台增加学校管理。
- 用户管理增加用户类型、区域、学校字段。
- 角色管理保留菜单权限和接口权限，增强数据范围配置。

交付物：

- 区域/学校/用户/角色可管理。
- 基础数据权限可验证。

## M2：文件存储与资源基础能力

目标：完成资源中心的底层文件能力。

后端任务：

- 新增 Storage 接口。
- 实现 MinIO Storage。
- 预留 OSS Storage。
- 新增文件元数据表：`edu_file` 或 `edu_resource_file`。
- 文件上传接口支持私有桶上传。
- 文件下载/预览接口支持鉴权和临时 URL。
- 设计对象 Key 规则，例如：`tenant/{tenant_id}/resource/{yyyy}/{mm}/{uuid}.{ext}`。
- 增加文件类型、大小、原始文件名、hash、上传人等字段。

前端任务：

- 后台封装 Arco Upload 组件。
- 支持资源封面上传。
- 支持附件上传。
- 支持上传进度和失败重试提示。

交付物：

- MinIO 可上传、下载、删除。
- 文件元数据入库。
- 可切换 OSS 的配置结构。

## M3：资源中心 MVP

目标：完成平台最核心的资源上传、审核、检索、查看闭环。

后端任务：

- 资源分类表：学段、障碍类型、资源类型、能力领域、专题分类。
- 标签表。
- 资源主表：标题、简介、封面、分类、标签、作者、学校、状态、统计字段。
- 资源附件表。
- 资源审核表。
- 收藏表。
- 评论表。
- 浏览、下载、收藏计数。
- 资源列表接口。
- 资源详情接口。
- 资源上传接口。
- 资源审核接口。
- 资源发布/下架接口。
- 资源收藏、取消收藏。
- 评论、回复、点赞。

后台前端任务：

- 资源分类管理。
- 标签管理。
- 资源管理列表。
- 资源上传/编辑页。
- 资源审核页。
- 资源详情管理页。

门户前端任务：

- 首页基础版。
- 资源中心列表。
- 筛选栏：学段、障碍类型、资源类型、能力领域、专题分类。
- 搜索框。
- 排序：最新、浏览量、下载量。
- 资源详情页。
- 收藏、下载、评论。

交付物：

- 教师上传资源。
- 学校管理员审核资源。
- 学生/教师检索和查看资源。

## M4：专题课程

目标：支持课程建设和学生在线学习。

后端任务：

- 课程表。
- 课程章节表。
- 课时/视频表。
- 课程资源关联表。
- 学习记录表。
- 学习进度表。
- 作业表。
- 作业提交表。
- 视频播放完成更新学习进度。

后台前端任务：

- 课程管理。
- 章节管理。
- 课时管理。
- 配套资源管理。
- 作业管理。

门户前端任务：

- 课程分类页。
- 课程列表页。
- 课程详情页。
- 学习页。
- 作业提交页。
- 学习进度展示。

交付物：

- 教师可建课。
- 学生可学习课程。
- 系统记录学习进度。

## M5：教研活动

目标：支持教师教研活动发布、报名、签到和成果沉淀。

后端任务：

- 活动表。
- 报名表。
- 签到表。
- 活动成果表。
- 活动资源关联表。
- 报名、取消报名、签到、成果上传接口。

后台前端任务：

- 活动管理。
- 报名管理。
- 签到管理。
- 成果管理。

门户前端任务：

- 活动列表。
- 活动详情。
- 报名入口。
- 成果展示。

交付物：

- 管理员可发布活动。
- 教师可报名和上传成果。

## M6：特教案例与 IEP

目标：沉淀特殊教育专业业务，严格控制敏感数据访问。

后端任务：

- 案例表。
- 学生基础信息表或案例学生档案表。
- IEP 表。
- 评估记录表。
- 干预方案表。
- 案例附件表。
- 案例审核表。
- 敏感数据访问日志。
- 数据脱敏输出能力。
- 授权访问机制。

后台前端任务：

- 案例管理。
- IEP 编辑。
- 评估记录管理。
- 干预方案管理。
- 案例审核。
- 敏感访问日志查看。

门户前端任务：

- 教师工作台中的个案入口。
- 案例列表。
- 案例详情。
- IEP 查看/编辑。

交付物：

- 教师可维护案例。
- 学校/区域管理员可审核。
- 非授权人员无法访问敏感案例。

## M7：名师/专家资源

目标：建设专家资源和优质课程聚合入口。

后端任务：

- 专家表。
- 专家领域标签。
- 专家课程/讲座关联。
- 专家资源附件。
- 收藏与分享。

后台前端任务：

- 专家管理。
- 专家资源管理。
- 讲座管理。

门户前端任务：

- 专家列表。
- 专家详情页。
- 专家课程与讲座列表。

交付物：

- 专家主页可展示。
- 专家资源可下载和收藏。

## M8：数据中心

目标：支持学校、区域、平台运营统计。

后端任务：

- 资源统计接口。
- 课程学习统计接口。
- 活动统计接口。
- 学校贡献排行。
- 教师活跃统计。
- 学生学习统计。
- Excel 导出。
- Redis 缓存热点统计。

后台前端任务：

- 数据概览页。
- 区域统计页。
- 学校统计页。
- 资源统计页。
- 学习统计页。
- 导出按钮。

交付物：

- 管理员可查看统计。
- 可导出基础报表。

## M9：搜索增强与性能优化

目标：从 MySQL 检索平滑升级到 ES/OpenSearch。

任务：

- 定义 Search 接口。
- MySQL Search 实现。
- ES/OpenSearch Search 实现。
- 资源索引同步。
- 中文分词。
- 分类聚合筛选。
- 热门资源缓存。
- 浏览量/下载量异步计数。

交付物：

- 搜索体验提升。
- 多维筛选性能稳定。

## 4. 推荐开发顺序

1. 先跑通后端和后台前端。
2. 做组织、租户、区域、学校、用户和数据权限。
3. 做 MinIO 文件存储。
4. 做资源中心 MVP。
5. 新建门户前端，并接资源中心接口。
6. 做课程。
7. 做教研活动。
8. 做特教案例/IEP。
9. 做专家资源。
10. 做数据中心。
11. 做 ES/OpenSearch 和性能优化。

## 5. 第一阶段详细任务拆分

第一阶段建议只做 M0 到 M3，形成可演示闭环。

### 后端

- 配置 MySQL、Redis、MinIO。
- 初始化 go-admin 基础表。
- 新增区域、学校、租户相关表。
- 扩展用户数据字段。
- 实现数据权限扩展。
- 实现 MinIO Storage。
- 实现资源分类、标签、资源、资源文件、资源审核。
- 实现资源上传、列表、详情、审核、发布、下载。

### 后台前端

- 确认 `vue3_dev` 分支能启动。
- 后台项目入口统一使用 `web/apps/admin`。
- 配置后端 API 地址。
- 登录、菜单、权限跑通。
- 区域管理。
- 学校管理。
- 用户管理字段扩展。
- 资源分类管理。
- 标签管理。
- 资源管理。
- 资源审核。

### 门户前端

- 使用 `web/apps/portal`。
- 首页基础布局。
- 资源列表。
- 资源详情。
- 登录态复用。
- 收藏、下载、评论入口。

## 6. 风险与注意事项

- `go-admin-ui` 的 `vue3_dev` 使用 Arco Design，不是 Element Plus，后续后台页面应统一 Arco 风格。
- `go-admin` 的多租户是按 Host/数据库选择的思路，不能完全等同本项目需要的区域/学校字段级数据隔离。
- 资源文件、案例附件、课程视频需要统一文件模型，避免每个模块各写一套上传逻辑。
- 特教案例和 IEP 是高敏感模块，必须等权限模型稳定后再开发。
- 课程视频如果直接 MP4 播放，一期可以接受；后续需要 HLS 转码。
- 代码生成器可用于基础 CRUD，但核心业务流程、权限过滤、审核流不要完全依赖生成代码。

## 7. 当前进度快照

当前分支：

- `main`
- 远程仓库：`https://github.com/ai1304/go-admin-edu.git`
- 最新业务提交：待提交本轮课程、活动、案例、专家、统计与后台菜单初始化代码。

当前阶段：

- 已完成项目基线整理的一部分：仓库已迁移为 `go-admin-edu`，后端、后台前端、门户前端、共享包已纳入同一个根仓库。
- 已开始落地 M1-M3：组织/学校、资源中心、门户资源展示的第一批业务骨架已提交。
- 资源中心尚未形成完整可演示闭环，原因是资源上传/审核 UI、完整 CRUD 表单、数据权限还未完成。

已完成：

- 后端新增 `app/edu` 业务模块。
- 后端新增模型：区域、学校、资源分类、资源标签、资源、资源文件、资源审核。
- 后端新增迁移：`go-admin/cmd/migrate/migration/version/2026051200010_edu_tables.go`。
- 新增启动与部署说明文档：`docs/deployment-guide.md`。
- 新增根级 `docker-compose.yml`，包含 MySQL、Redis、MinIO、后端、后台前端、门户前端。
- 新增后端 Docker 配置：`go-admin/config/settings.docker.yml`。
- 后端新增 MinIO 对象存储封装：`go-admin/common/objectstorage`。
- 后端新增资源文件真实上传接口：`POST /api/v1/edu/resource-files/upload`。
- 后端新增课程、教研活动、特教案例/IEP、专家资源、数据中心模型与接口骨架。
- 后端新增业务迁移：`go-admin/cmd/migrate/migration/version/2026051300010_edu_business_tables.go`。
- 后端新增后台菜单迁移：`go-admin/cmd/migrate/migration/version/2026051500010_edu_menus.go`，用于初始化「特殊教育」目录及教育业务页面菜单。
- 修复后台前端 Dockerfile 中遗留的合并冲突标记。
- 新增门户前端 Dockerfile。
- 后端新增后台管理接口：
  - `/api/v1/edu/regions`
  - `/api/v1/edu/schools`
  - `/api/v1/edu/resource-categories`
  - `/api/v1/edu/resource-tags`
  - `/api/v1/edu/resources`
  - `/api/v1/edu/resource-files`
  - `/api/v1/edu/resource-files/upload`
  - `/api/v1/edu/courses`
  - `/api/v1/edu/activities`
  - `/api/v1/edu/cases`
  - `/api/v1/edu/experts`
  - `/api/v1/edu/stats/overview`
- 后端新增门户公开资源接口：
  - `/api/v1/portal/resources`
  - `/api/v1/portal/resources/:id`
  - `/api/v1/portal/courses`
  - `/api/v1/portal/courses/:id`
  - `/api/v1/portal/activities`
  - `/api/v1/portal/activities/:id`
  - `/api/v1/portal/experts`
  - `/api/v1/portal/experts/:id`
- 后台前端新增教育业务 API 封装：
  - `web/apps/admin/src/api/edu/region.js`
  - `web/apps/admin/src/api/edu/school.js`
  - `web/apps/admin/src/api/edu/resource.js`
  - `web/apps/admin/src/api/edu/course.js`
  - `web/apps/admin/src/api/edu/activity.js`
  - `web/apps/admin/src/api/edu/case.js`
  - `web/apps/admin/src/api/edu/expert.js`
  - `web/apps/admin/src/api/edu/stats.js`
- 后台前端新增页面骨架：
  - `web/apps/admin/src/views/edu/region/index.vue`
  - `web/apps/admin/src/views/edu/school/index.vue`
  - `web/apps/admin/src/views/edu/resource/index.vue`
  - `web/apps/admin/src/views/edu/course/index.vue`
  - `web/apps/admin/src/views/edu/activity/index.vue`
  - `web/apps/admin/src/views/edu/case/index.vue`
  - `web/apps/admin/src/views/edu/expert/index.vue`
  - `web/apps/admin/src/views/edu/stats/index.vue`
- 后台菜单已规划为数据库驱动，管理员角色会自动看到 `sys_menu` 中 `M/C` 类型菜单；普通角色需要在角色菜单权限中勾选对应教育业务菜单。
- 门户前端新增资源 API：`web/apps/portal/src/api/resources.js`。
- 门户前端新增课程、活动、专家 API：`web/apps/portal/src/api/courses.js`、`activities.js`、`experts.js`。
- 门户资源、课程、活动、专家列表已接入公开接口，资源详情已接入公开资源详情接口。

已验证：

- 后端曾执行 `go build ./...` 通过；本轮菜单迁移补充后需再次验证。
- 当前本地仓库状态曾在提交后保持 `main...origin/main`。

未验证：

- 前端依赖尚未安装，未执行 `pnpm dev:admin`、`pnpm dev:portal`、`pnpm build`。
- 尚未实际连接 MySQL 执行迁移。
- 当前机器 Docker CLI 不可用，尚未执行 `docker compose config` 或启动完整 Docker 环境。

## 8. 近期最小可执行清单

下一步建议按这个顺序继续，避免业务代码越写越散：

1. 在有 Docker 的机器上执行 `docker compose config` 并启动 MySQL、Redis、MinIO。
2. 安装前端依赖并验证 `web/apps/admin`、`web/apps/portal` 能启动。
3. 执行后端迁移，确认教育业务表能成功创建。
4. 后台资源页面接入 `POST /api/v1/edu/resource-files/upload`。
5. 将当前后台列表骨架升级为完整 CRUD 表单：区域、学校、资源、课程、活动、案例、专家。
6. 补充教育业务按钮级权限和普通角色授权策略。
7. 完成资源上传/编辑/提交审核/审核页面。
8. 完成门户详情页：课程详情、活动详情、专家详情、资源预览、下载入口。
9. 扩展用户与数据权限：`tenant_id`、`region_id`、`school_id`、用户类型。

## 9. 当前前端目录结构

```text
web/
  package.json
  pnpm-workspace.yaml
  apps/
    admin/          # 后台管理端
    portal/         # 门户/学习端
  packages/
    shared/         # 前后台共享类型、常量、工具
```

启动命令：

```bash
cd web
pnpm dev:admin
pnpm dev:portal
```
