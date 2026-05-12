# 特殊教育资源库平台需求与技术方案

更新时间：2026-05-12

## 1. 项目定位

特殊教育资源库平台面向校内、多校联盟、区域以及未来行业级资源共建场景，目标是建设一个可检索、可共享、可审核、可统计的特殊教育数字化资源平台。

参考产品形态类似国家智慧教育平台的资源门户，但本项目不优先考虑 SEO，也不拆独立门户技术栈。前台资源浏览端和后台管理端统一采用 Vue3 技术体系。

核心价值：

- 沉淀校内优质教学资源，提升资源复用率。
- 支持教师备课、课程建设、教研协同和专业成长。
- 支持学生课程学习、资源查看、作业提交和学习进度记录。
- 支持特殊教育个案、IEP、评估、干预方案等专业业务沉淀。
- 支持多学校、区域级数据统计、资源审核和资源下沉。
- 为后续行业资源认证、专家智库、决策支持预留扩展空间。

## 2. 产品发展阶段

| 阶段 | 核心目标 | 描述 |
| --- | --- | --- |
| 校内资源库（0 -> 1） | 校内资源管理和共享 | 统一管理校内资源，服务校内教学 |
| 多校联盟（1 -> 10） | 跨校协作和资源共建 | 联盟学校资源互通，联合教研 |
| 区域资源库（10 -> 100） | 区域统筹与数据下沉 | 区域学校数据统计、资源下沉，支持教师和学生 |
| 行业资源库（100 -> 1000+） | 行业标准化与决策支持 | 行业资源认证、专家智库、辅助决策 |

## 3. 用户角色与权限

| 角色 | 核心权限 | 数据范围 |
| --- | --- | --- |
| 超级管理员 | 系统配置、租户/区域管理、全局资源管理、全局审计 | 全平台 |
| 区域管理员 | 管理区域学校、统计数据、资源审批、区域资源下沉 | 当前区域 |
| 学校管理员 | 校内教师管理、资源管理、资源审核、课程/活动管理 | 当前学校 |
| 教师 | 上传/下载资源、课程管理、教研活动参与、个案/IEP维护 | 个人、班级、学校授权范围 |
| 学生 | 查看资源、课程学习、作业完成、收藏资源 | 个人学习数据与开放资源 |

权限模型建议：

- 使用 RBAC 控制菜单、按钮、接口、资源操作权限。
- 使用数据权限控制区域、学校、部门、班级、个人数据范围。
- 租户模型建议从第一期预留，避免后续从单校扩展到多校/区域时大改。
- 特教案例、IEP、学生评估数据属于敏感数据，必须做严格的数据隔离、脱敏和操作审计。

## 4. 核心业务模块

### 4.1 首页

首页是统一入口，面向区域管理员、学校管理员、教师、学生展示不同内容。

主要内容：

- 顶部导航：Logo、首页、资源中心、专题课程、教研活动、特教案例、名师资源、AI 应用、登录/注册。
- Banner 区：宣传语、搜索框、快捷入口、活动展示。
- 快捷功能区：资源中心、专题课程、教研活动、特教案例、名师资源。
- 专题展示区：热门专题卡片。
- 精品资源推荐：精品课程、热门资源、最新上传。
- 通知公告：教研活动、专家讲座、平台通知。
- 页脚：平台介绍、联系方式、版权信息。

交互要求：

- 搜索框支持关键词搜索和高级筛选。
- 卡片点击进入对应模块详情页。
- Banner 轮播图可跳转活动或课程。

### 4.2 资源中心

面向教师和学生提供课程资源、案例资料、微课视频、教学文档等内容。

核心功能：

- 资源搜索、分类筛选、排序。
- 资源在线预览、播放、下载。
- 资源收藏、评论、点赞、回复。
- 资源上传、标签填写、分类选择、提交审核。
- 相关推荐。

筛选维度：

- 学段。
- 障碍类型。
- 资源类型。
- 能力领域。
- 专题分类。
- 标签。

页面：

- 资源列表页。
- 资源详情页。
- 资源上传页。
- 资源审核页。

### 4.3 专题课程

用于在线课程管理与学习。

核心功能：

- 课程分类、课程列表、课程详情。
- 视频播放。
- 课程介绍、教学目标、配套资源。
- 学习记录、学习进度、收藏。
- 作业提交。
- 教师上传课程并管理资源。

### 4.4 教研活动

用于教师教研活动发布、参与和成果记录。

核心功能：

- 活动发布。
- 活动列表和详情。
- 报名、签到、参与记录。
- 活动资源关联。
- 活动结束后上传成果文档。

### 4.5 特教案例

用于个案管理、IEP 管理、评估工具库和干预方案沉淀。

核心功能：

- 案例列表。
- 案例详情。
- 学生基础信息。
- IEP 记录。
- 评估结果。
- 干预方案。
- 附件上传。
- 学校/区域管理员审核。

注意：

- 该模块涉及敏感学生信息，应默认非公开。
- 需要按学校、班级、教师、授权关系进行数据隔离。
- 操作日志、访问日志、导出日志必须保留。

### 4.6 名师/专家资源

用于汇聚名师课程、专家讲座、指导资料和优质资源。

核心功能：

- 专家列表。
- 专家详情页。
- 专家课程、讲座、资源下载。
- 收藏与分享。

### 4.7 数据中心

用于区域、学校、平台运营统计。

核心指标：

- 资源总数、资源类型分布、资源上传趋势。
- 资源浏览量、下载量、收藏量。
- 学校资源贡献排行。
- 教师活跃度。
- 学生学习人数、学习进度、课程完成率。
- 活动报名、签到、成果提交统计。
- 案例/IEP 数量与状态统计。

导出：

- 支持 Excel 报表导出。
- 后续可支持区域级报告生成。

## 5. 技术选型

### 5.1 总体结论

采用前后端分离架构：

- 前端：Vue3 + Vite + TypeScript。
- UI 组件：Element Plus 优先。
- 后端：Go + Gin。
- 数据库：MySQL 8.0。
- ORM：GORM。
- 权限：Casbin RBAC + 数据权限。
- 缓存：Redis。
- 文件存储：MinIO 起步，生产切换阿里云 OSS。
- 搜索：第一期 MySQL 索引，后续引入 Elasticsearch 或 OpenSearch。
- 部署：Docker / Docker Compose。

### 5.2 前端技术

前台资源浏览端和后台管理端统一使用：

- Vue3。
- Vite。
- TypeScript。
- Pinia。
- Vue Router。
- Element Plus。
- Axios 或基于 Fetch 的请求封装。

建议工程组织：

```text
web/
  apps/
    portal/        # 前台资源浏览端，不做 SSR
    admin/         # 后台管理端
  packages/
    ui/            # 可复用业务组件
    api-client/    # 接口 SDK 或请求封装
    shared/        # 常量、类型、工具函数
```

当前项目已按该方向建立前端 workspace：

```text
web/
  apps/
    admin/         # 后台管理端，vue3_dev 分支
    portal/        # 门户/学习端
  packages/
    shared/        # 跨端共享常量、类型、工具
```

### 5.3 后端开源框架推荐

推荐使用：`go-admin-team/go-admin`

项目地址：

- GitHub: https://github.com/go-admin-team/go-admin
- 文档: https://www.go-admin.pro

推荐原因：

- 基于 Gin，符合后端技术偏好。
- 基于 GORM，适配 MySQL。
- 内置 JWT 鉴权。
- 内置 Casbin RBAC 权限模型。
- 支持菜单、角色、用户、部门、岗位、字典、参数、日志等后台基础能力。
- 项目介绍明确包含多租户支持。
- 支持数据权限，适合区域、学校、部门、个人范围控制。
- 有代码生成器，可加速资源、课程、活动等 CRUD 模块开发。
- 有 Docker 部署示例，便于容器化落地。

备选方案：

- `gin-vue-admin`：生态更活跃，Vue3 + Gin 体验成熟，RBAC、对象存储、Docker 文档完善，但原生多租户不是其最核心优势。如果最终发现 go-admin 的 Vue3/前端栈不符合项目习惯，可考虑选择 gin-vue-admin 并自行补齐租户和数据权限模型。

本项目建议：

- 后端以 `go-admin` 为底座进行二次开发。
- 前端可参考其 Vue3 版本，但最终统一为 Vue3 + TypeScript + Element Plus 的项目规范。
- 不直接把业务强绑定在框架生成代码里，核心业务模块保持清晰分层，便于后续迁移。

### 5.4 后端分层建议

```text
server/
  cmd/                  # 启动入口
  config/               # 配置
  internal/
    auth/               # 登录、JWT、权限
    tenant/             # 租户、区域、学校数据范围
    user/               # 用户、角色、部门、岗位
    resource/           # 资源中心
    course/             # 专题课程
    activity/           # 教研活动
    casefile/           # 特教案例、IEP
    expert/             # 名师/专家
    stats/              # 数据中心
    storage/            # MinIO/OSS 抽象
    search/             # MySQL/ES 搜索抽象
    audit/              # 操作日志、审计
  pkg/                  # 通用包
```

后端接口风格：

- RESTful API。
- Swagger/OpenAPI 文档。
- 统一响应结构。
- 统一错误码。
- 统一分页参数。
- 统一审计日志中间件。

## 6. 存储方案

### 6.1 MySQL

MySQL 作为主业务数据库。

建议关键表：

- sys_user。
- sys_role。
- sys_menu。
- sys_api。
- sys_dept。
- sys_tenant。
- edu_region。
- edu_school。
- edu_resource。
- edu_resource_category。
- edu_resource_tag。
- edu_resource_file。
- edu_resource_review。
- edu_resource_favorite。
- edu_resource_comment。
- edu_course。
- edu_course_chapter。
- edu_course_lesson。
- edu_learning_record。
- edu_assignment。
- edu_activity。
- edu_activity_signup。
- edu_activity_checkin。
- edu_case。
- edu_case_iep。
- edu_case_assessment。
- edu_case_intervention。
- edu_expert。
- edu_expert_resource。
- sys_audit_log。

### 6.2 文件存储

第一期：

- 使用 MinIO。
- 本地 Docker Compose 部署。
- 通过 S3 兼容 API 访问。

生产：

- 切换为阿里云 OSS。
- 后端封装统一 Storage 接口，避免业务代码感知具体实现。

建议抽象：

```go
type Storage interface {
    PutObject(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error
    GetObjectURL(ctx context.Context, objectKey string, expire time.Duration) (string, error)
    DeleteObject(ctx context.Context, objectKey string) error
}
```

### 6.3 文件预览

建议：

- 图片：直接预览。
- PDF：浏览器内嵌预览。
- Word/PPT/Excel：服务端异步转 PDF 后预览。
- 视频：初期可直接 MP4 播放，后续转 HLS。

转码/转换任务建议后续引入异步队列。

## 7. 搜索与性能优化

第一期：

- 使用 MySQL 普通索引。
- 对资源标题、分类、标签、学段、障碍类型、能力领域、上传人、审核状态建立索引。
- 热门资源、首页推荐、统计数字使用 Redis 缓存。

第二期：

- 引入 Elasticsearch 或 OpenSearch。
- 支持中文分词。
- 支持多维筛选聚合。
- 支持相关推荐。
- 支持资源标题、简介、标签、作者、OCR/文档解析内容的全文检索。

Redis 用途：

- 登录态或 Token 黑名单。
- 验证码。
- 接口限流。
- 热点资源缓存。
- 首页推荐缓存。
- 浏览量、下载量、收藏量异步计数。

## 8. Docker 部署

第一期使用 Docker Compose。

建议服务：

```text
mysql
redis
minio
server
portal-web
admin-web
nginx
```

建议目录：

```text
deploy/
  docker-compose.yml
  nginx/
    nginx.conf
  mysql/
    init/
  minio/
  env/
    dev.env
    prod.env.example
```

Nginx 路由建议：

- `/` -> portal-web。
- `/admin/` -> admin-web。
- `/api/` -> server。
- `/storage/` -> MinIO 或 OSS 回源/签名下载。

## 9. 多租户与数据权限设计

推荐逻辑：

- `tenant_id` 表示租户，可对应学校、集团、联盟或区域客户。
- `region_id` 表示区域。
- `school_id` 表示学校。
- `dept_id` 表示组织部门。
- 业务表默认保留 `tenant_id`、`region_id`、`school_id`、`created_by`、`updated_by`。

数据隔离策略：

- 超级管理员：不受租户限制。
- 区域管理员：按 `region_id` 限制。
- 学校管理员：按 `school_id` 限制。
- 教师：按本人创建、授权班级、所属学校开放资源限制。
- 学生：按公开资源、课程授权、个人学习记录限制。

敏感模块：

- 特教案例。
- IEP。
- 学生评估。
- 干预方案。

敏感模块默认不允许跨校访问，跨校共享必须显式授权或脱敏发布。

## 10. 安全与合规要求

必须实现：

- HTTPS。
- 密码加密存储。
- JWT 过期与刷新。
- 接口权限校验。
- 数据权限校验。
- 文件下载鉴权。
- 敏感信息脱敏。
- 操作日志。
- 登录日志。
- 资源审核日志。
- 重要数据导出日志。

建议实现：

- IP 限流。
- 登录失败锁定。
- 管理端操作二次确认。
- OSS/MinIO 私有桶 + 临时签名 URL。

## 11. 无障碍与特殊教育体验要求

由于平台面向特殊教育场景，前端从第一期开始应考虑：

- 颜色对比度。
- 字体缩放。
- 键盘可访问。
- 表单清晰标签。
- 视频字幕。
- 图片替代文本。
- 语义化 HTML。
- 操作反馈明确。
- 避免复杂动效造成干扰。

## 12. 分期实施建议

### 一期 MVP

目标：完成校内资源库 0 -> 1。

范围：

- 登录、用户、角色、菜单、数据权限。
- 学校、教师、学生基础管理。
- 资源中心：上传、分类、标签、审核、列表、详情、预览、下载、收藏。
- 首页基础展示。
- MinIO 文件存储。
- Docker Compose 部署。

### 二期

目标：课程学习和教研协同。

范围：

- 专题课程。
- 章节、视频、学习进度。
- 作业提交。
- 教研活动、报名、签到、成果上传。
- 数据中心基础统计。

### 三期

目标：专业特教业务沉淀。

范围：

- 特教案例。
- IEP 管理。
- 评估工具库。
- 干预方案。
- 敏感数据权限与审计增强。

### 四期

目标：区域化、智能化和搜索增强。

范围：

- 多校联盟/区域资源共享。
- Elasticsearch/OpenSearch。
- 推荐系统。
- AI 工具集。
- 报告导出。
- 行业资源认证。

## 13. 后续 AI 接手须知

本项目当前确定的关键选择：

- 不使用 Nuxt，不优先考虑 SEO。
- 前台和后台都使用 Vue3 技术体系。
- 后端使用 Go + Gin。
- 优先评估并采用 `go-admin-team/go-admin` 作为后端管理框架底座。
- 数据库使用 MySQL。
- 文件存储第一期使用 MinIO，生产切换阿里云 OSS。
- 缓存使用 Redis。
- 搜索第一期使用 MySQL，后续引入 Elasticsearch 或 OpenSearch。
- 部署使用 Docker / Docker Compose。
- 项目必须从第一期预留多租户、数据权限和敏感数据审计。

后续开发时优先关注：

- 不要把资源文件存入 MySQL，只存元数据和对象存储地址。
- 不要把特教案例和普通资源使用同一套公开权限逻辑。
- 所有业务表尽量预留 `tenant_id`、`region_id`、`school_id`、`created_by`。
- 文件存储必须通过抽象接口封装，便于 MinIO 切 OSS。
- 搜索也应通过抽象接口封装，便于 MySQL 切 ES/OpenSearch。
- 首页、资源中心、课程、活动、案例、专家、数据中心是当前主要业务边界。

## 14. 参考资料

- SmartEdu 官方枢纽：https://system.smartedu.cn/
- go-admin：https://github.com/go-admin-team/go-admin
- go-admin 文档：https://www.go-admin.pro
- gin-vue-admin：https://github.com/flipped-aurora/gin-vue-admin
