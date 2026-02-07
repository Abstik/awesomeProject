# 项目说明 - 西游移动开发实验室后端 API

## 项目概述

基于 Go + Gin + GORM + MySQL 的 RESTful API 后端，用于管理实验室成员、活动、团队方向、捐款、视频等数据。
运行端口：8080，API 前缀：`/xupt-web/api`

## 项目结构

```
handler/     — HTTP 请求处理层（接收请求、参数校验、调用 service）
service/     — 业务逻辑层（核心逻辑、数据组装）
dao/         — 数据访问层（GORM 数据库操作）
model/       — 数据模型（请求体、数据库实体）
middleware/  — 中间件（JWT 鉴权、管理员鉴权）
router/      — 路由配置
utils/       — 工具函数（JWT、日志、响应构建、URL 处理、密码哈希、验证码）
dao/main/    — 一次性数据库迁移脚本（独立 main 包）
```

## 代码规范

- 分层架构：handler -> service -> dao，禁止跨层调用（handler 不直接调用 dao）
- 模型命名：请求体用 `XxxRequest` / `XxxReq`，数据库实体用 `XxxPO`，业务视图用 `XxxVO`
- 请求体标签规范：必填字段用 `binding:"required"`，可选字段不加标签；`omitempty` 仅用于响应/PO 的 JSON 序列化，Request 结构体不应使用
- 客户端错误响应使用 `utils.BuildErrorResponse(c, code, msg)`（不记录日志）
- 服务端错误响应使用 `utils.BuildServerError(c, msg, err)`（自动记录 ERROR 日志并返回 500）
- 成功响应统一使用 `utils.BuildSuccessResponse(c, data)`
- 面向客户端的错误信息不应暴露内部实现细节（如 err.Error()）
- 所有写操作路由需加 JWT 鉴权中间件，管理员操作额外加 IsAdminAuthMiddleware
- 数据库字段使用下划线命名，JSON 字段使用驼峰命名
- 密码使用 bcrypt 哈希（`utils.HashPassword` / `utils.CheckPassword`）
- JWT 密钥硬编码在 `utils/jwt.go` 中

## 密码体系

- 新用户注册：默认密码为 `用户名 + "123"`，bcrypt 哈希后存库
- 密码重置：管理员通过 `/reset_password` 接口重置为 `用户名 + "123"`，bcrypt 哈希
- 登录验证：bcrypt 校验
- 修改用户信息接口（`ChangeMemberInfo`）不支持修改密码
- 批量重置脚本：`dao/main/main.go` → `InitAllMemberPassword()`，随机密码 + Excel 导出

## 环境变量

- `GO_ENV` — 运行模式（debug/release/test）
- `DB_USER` / `DB_PASSWORD` / `DB_HOST` / `DB_PORT` / `DB_NAME` — 数据库配置（未设置时使用本地默认值）
- `DOMAIN_NAME` — 域名，用于拼接资源 URL
- `DEFAULT_PORTRAIT` / `DEFAULT_MIEN_IMG` / `DEFAULT_GRADUATE_IMG` — 默认图片路径

## 角色体系

- status = 0：管理员
- status = 1：普通用户
- JWT 有效期通过 `TokenExpireDuration` 配置（当前 1 天）

## 依赖

- gin v1.10.1 — Web 框架
- gorm v1.25.12 — ORM
- golang-jwt/jwt/v5 — JWT 令牌
- golang.org/x/crypto/bcrypt — 密码哈希
- zap + lumberjack — 结构化日志 + 按大小轮转（50MB/文件，保留 30 个备份）
- base64Captcha — 验证码（内存存储，1 分钟过期）
- excelize/v2 — Excel 导出
- godotenv — 环境变量加载

---

## 已修复问题

| 编号 | 问题 | 修复内容 |
|------|------|----------|
| S1 | 图片删除路径穿越 | `handler/img.go` 增加 `filepath.Abs` 二次校验 |
| S4 | 无 ReadHeaderTimeout | `main.go` 添加 `ReadHeaderTimeout: 10s` |
| S5 | 内部错误信息泄露 | 所有 handler 统一返回通用错误信息 |
| B1 | 密码逻辑 | 注册和重置统一使用 `用户名 + "123"` + bcrypt |
| B2 | 多处缺少 return | contact/introduction/team/trainplan handler 已修复 |
| B3 | DeleteActivity 忽略 ParseInt | 添加错误检查和 400 响应 |
| B4 | GORM Count 查询污染 | `dao/member.go` 使用 `Session(&gorm.Session{})` 克隆 |
| B5 | ActivityPO.Title 空 column 标签 | 修正为 `column:title` |
| L1 | DeleteMember 跨层调用 | handler 改为调用 `service.DeleteMember()` |
| L2 | GetMemberByName 应用层过滤 | SQL 增加 `AND status != 0` |
| L3 | GetTeams 排序不确定 | `dao/team.go` 添加 `ORDER BY tid ASC` |
| L5 | Register 重复代码 | 移除 RegisterWithResult，统一为 Register |
| P1 | 无连接池配置 | `dao/connector.go` 配置 MaxIdleConns/MaxOpenConns/ConnMaxLifetime |
| P3 | 批量重置无事务 | `InitAllMemberPassword` 包裹在 `db.Transaction` 中 |
| M1 | `ActivityReq` 带多余 gorm 标签 | 移除 `model/activity.go` 中请求体的 gorm 标签 |
| M2 | 捐款时间校验不一致 | 统一 `AddDonationReq` 和 `AddDonationsReq` 的 binding 为 `datetime=2006` |
| D1 | `GetActivityByAid` 查不到不报错 | 改用 `First`，handler 增加 404 处理 |
| M4 | `TrainPlan` 未遵循 PO 命名 | 重命名为 `TrainPlanPO`，更新所有引用 |
| M5 | `TeamPO.Delay` 混入非 DB 字段 | 新增 `TeamVO` 承载业务计算字段，`TeamPO` 仅保留数据库字段 |
| M6 | `ContactPO.Tid` 字段名误导 | 重命名为 `ID` |
| D2 | `UpdateTeam` DAO 层依赖 Request 模型 | 改为接收 `*TeamPO`，service 层负责 Req→PO 转换 |
| L6 | `zap.L()` 全局无效 | `utils/logger.go` 补充 `zap.ReplaceGlobals(Logger)`，使 ginzap 等全局调用生效 |
| L7 | 错误信息中英混杂 | 统一所有 handler 面向客户端的错误信息为中文 |
| L8 | 500 响应未记录日志 | `handler/member.go`、`handler/captcha.go` 改用 `BuildServerError` 记录 ERROR |
| L9 | 登录失败无审计日志 | `handler/member.go` 添加 `zap.L().Warn("登录失败", ...)` 含用户名、IP、原因 |
| L10 | `init()` 使用 `log.Fatalf` | `handler/img.go` 改为 `panic()`，与 `handler/video.go` 风格统一 |
| L11 | `main.go` 英文日志 | 翻译为中文：数据库初始化失败、服务启动失败、服务强制关闭 |
| M8 | `model/user.go` 文件名与内容不符 | 重命名为 `model/member.go` |
| M9 | `AddDonationsReq` 内联重复定义 | 改为引用 `AddDonationReq` |
| M10 | `ActivityListReq` 死代码 | 删除未使用的结构体 |
| M11 | `TrainPlanPO` GORM 标签错误 | `primarykey` → `primaryKey`，`gorm:"content"` → `gorm:"column:content"` |
| N1 | 空指针风险：`donation.Money` | `handler/donation.go` 添加 nil 检查 |
| N2 | 空指针风险：`teamPOs.IsExist` | `service/team.go` 添加 nil 检查 |
| N3 | 空指针风险：`res.Content` | `service/activity.go` 添加 nil 检查 |
| D3 | Docker 环境变量不匹配 | `docker-compose.prod.yml` 将 `MYSQL_DATABASE` 改为 `DB_NAME` 与代码一致 |
| S3 | 数据库密码硬编码 | `dao/connector.go` 通过 `envOrDefault()` 读取环境变量，本地保留默认值 |
| M3 | `MemberRequest` 职责过多 | 拆分为 `RegisterReq` + `LoginReq` + `UpdateMemberRequest`，删除旧 `MemberRequest` |

## 待处理/已知问题

| 编号 | 问题 | 说明 |
|------|------|------|
| S2 | `Rand5Digits` 使用 `math/rand` | 仅用于文件名生成，非安全敏感场景，可接受 |
| S6 | 活动内容 XSS 风险 | 仅管理员可写入，建议前端渲染时过滤 |
| L4 | FullURL/OldFullURL 代码重复 | 可优化但不影响功能 |
| P2 | GetMemberList 两次查询 | 当前规模影响较小 |

## 待处理：Model 层和 DAO 层优化项

| 编号 | 问题 | 说明 |
|------|------|------|
| M7 | Req 结构体 `omitempty` 滥用 | `AddTeamReq` 等请求体加了 `omitempty`，应清除；部分 Req 缺少 `binding:"required"` |
| D4 | `UpdateTrainPlan` 硬编码 `id = 1` | 单例表设计，功能无误，仅 magic number 不够清晰 |
