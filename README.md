# Go-Zero 用户管理系统 - 学习笔记

一个完整的 go-zero 入门项目，通过用户管理系统学习 go-zero 框架的核心概念、分层架构和开发流程。

---

## 📚 项目简介

本项目实现了用户管理的完整 CRUD 功能：
- 用户注册
- 查询用户信息
- 用户列表（支持分页）
- 更新用户
- 删除用户

**学习目标**：理解 go-zero 的分层架构、数据流动、代码生成工具的使用。

---

## 🏗️ 核心架构理解

### 1. go-zero 的分层设计

go-zero 将应用分为三个核心层次，每层职责明确：

```
┌─────────────────────────────────┐
│  API 层 (user.api)              │  → 定义接口规范（对外的合同）
└─────────────────────────────────┘
            ↓
┌─────────────────────────────────┐
│  Handler 层                      │  → HTTP 请求处理（解析参数、返回响应）
└─────────────────────────────────┘
            ↓
┌─────────────────────────────────┐
│  Logic 层                        │  → 业务逻辑（参数校验、流程控制）
└─────────────────────────────────┘
            ↓
┌─────────────────────────────────┐
│  Model 层                        │  → 数据库操作（SQL 封装）
└─────────────────────────────────┘
            ↓
┌─────────────────────────────────┐
│  数据库 (MySQL)                  │  → 数据存储
└─────────────────────────────────┘
```

**关键原则**：每一层只关心自己的事，通过接口与其他层交互。

---

### 2. 数据流动逻辑

以"查询用户列表"为例，理解数据如何在各层之间流动：

**用户请求**：
```
GET /user/list?page=1&page_size=10
```

**数据流动过程**：

1. **API 定义阶段**（开发前）
   - 在 `.api` 文件中定义接口的"样子"
   - 告诉系统：这个接口接收什么参数、返回什么数据

2. **Handler 层**（请求进来）
   - 接收 HTTP 请求
   - 解析查询参数：`page=1`, `page_size=10`
   - 把参数封装成 Go 结构体
   - 调用 Logic 层

3. **Logic 层**（业务处理）
   - 接收参数结构体
   - 参数校验（page 不能小于 1）
   - 计算偏移量：`offset = (page - 1) * page_size`
   - 调用 Model 层查询数据
   - 把 Model 返回的数据转换成 API 响应格式
   - 返回给 Handler

4. **Model 层**（数据库操作）
   - 接收 offset 和 limit
   - 执行 SQL：`SELECT * FROM user LIMIT 10 OFFSET 0`
   - 查询总数：`SELECT COUNT(*) FROM user`
   - 把数据库记录映射成 Go 结构体
   - 返回给 Logic

5. **Handler 层**（返回响应）
   - 接收 Logic 返回的数据
   - 序列化成 JSON
   - 返回给客户端

**返回数据**：
```json
{
  "list": [...],
  "total": 25,
  "page": 1,
  "page_size": 10
}
```

---

## 🔑 核心概念详解

### API 文件的作用

**API 文件 = 接口的说明书**

定义三件事：
1. **接口地址**：这个功能的 URL 是什么
2. **请求参数**：外部需要传什么数据
3. **返回数据**：外部会得到什么数据

**重要理解**：
- API 字段 ≠ 数据库字段
- API 是给外部看的，数据库是内部存储
- 通过 Logic 层进行转换和映射

示例：
- API 可能只要求用户传 `name` 和 `email`
- 但数据库会自动添加 `id`、`created_at`、`updated_at`
- Logic 层负责补充这些字段

---

### Model 层的两个文件

项目中有两个 Model 文件，分工明确：

**1. `usermodel_gen.go`（自动生成，不要改）**
- goctl 根据数据库表自动生成
- 包含基础的 CRUD 方法
- 封装了所有 SQL 操作
- 每次运行 goctl 会重新生成

**2. `usermodel.go`（可以编辑）**
- 用来添加自定义方法
- 比如复杂的查询、业务相关的数据操作
- 不会被 goctl 覆盖

**为什么分两个文件？**
- 保护你的自定义代码不被覆盖
- 基础功能自动生成，省去重复劳动
- 扩展功能手动添加，保持灵活性

---

### 字段映射的三层关系

**Go 字段名 ↔ Tag ↔ 外部名称**

```
类型定义：
type User struct {
    PageSize int64 `form:"page_size"`
}

三层含义：
1. Go 代码中：req.PageSize（大驼峰，代码中用）
2. Tag 标记：form:"page_size"（蛇形，映射规则）
3. URL 参数：?page_size=10（蛇形，外部传递）
```

**不同场景的 Tag**：
- `form` → 查询参数（`?page=1`）
- `path` → 路径参数（`/user/:id`）
- `json` → JSON 请求体（`{"name": "xxx"}`）
- `db` → 数据库字段（`user.name` 列）

**关键点**：
- Go 字段名必须大写（才能被外部访问）
- Tag 可以定义外部的命名格式
- 数据库字段通过 `db` tag 映射

---

## 🛠️ 开发流程

### 添加新功能的标准流程

以"添加分页功能"为例：

**步骤 1：设计接口**（思考）
- 接口地址：`GET /user/list`
- 需要什么参数：`page`, `page_size`
- 返回什么数据：用户列表、总数、页码

**步骤 2：修改 API 文件**（定义）
- 定义请求结构：`ListUsersReq`
- 定义响应结构：`ListUsersResp`
- 定义路由：`get /list`

**步骤 3：运行 goctl**（生成代码）
```bash
goctl api go -api user.api -dir .
```
- 自动生成 Handler
- 自动生成 Logic 框架
- 自动生成 Types
- 自动注册路由

**步骤 4：实现 Model 方法**（如果需要）
- 在 `usermodel.go` 中添加分页查询方法
- 封装 SQL：`LIMIT` 和 `OFFSET`

**步骤 5：实现 Logic 业务逻辑**（核心）
- 参数校验
- 调用 Model 查询
- 数据转换
- 返回响应

**步骤 6：测试**（验证）
- 启动服务
- 使用 Apifox 测试
- 检查日志
- 验证数据

---

## ⚠️ 常见问题与解决

### 问题 1：400 错误 - 字段小写导致无法解析

**现象**：
```
[HTTP] 400 - POST /user/register
```

**原因**：
API 文件中字段是小写：
```
type RegisterReq {
    name string   // ❌ 小写，JSON 无法解析
}
```

**解决**：
字段必须大写，用 tag 指定 JSON 名称：
```
type RegisterReq {
    Name string `json:"name"`  // ✅ 正确
}
```

**理解**：Go 语言中，小写字段是私有的，外部无法访问。

---

### 问题 2：SQL 查询字段不匹配

**现象**：
```
查询数据库失败: not matching destination to scan
```

**原因**：
SQL 查询的字段少于结构体字段：
```
SQL: select id, name, email from user
结构体: User { Id, Name, Email, Age, CreatedAt, UpdatedAt }
```

**解决**：
SQL 必须查询所有字段：
```
SQL: select id, name, email, age, created_at, updated_at from user
```

**理解**：go-zero 需要把 SQL 结果映射到结构体，字段必须一一对应。

---

### 问题 3：404 错误 - 路由不匹配

**现象**：
```
[HTTP] 404 - PUT /user/update/1
```

**原因**：
URL 路径错误，多加了 `update`：
```
API 定义: put /:id
正确 URL: PUT /user/1
错误 URL: PUT /user/update/1
```

**解决**：
直接使用 `/user/:id`，不要加多余的路径。

**理解**：RESTful 风格，用 HTTP 方法区分操作，而不是路径。

---

### 问题 4：修改代码后还是报错

**现象**：
修改了代码，但测试时还是旧的错误。

**原因**：
服务没有重启，还在运行旧代码。

**解决**：
每次修改代码后，必须重启服务：
1. 按 `Ctrl+C` 停止服务
2. 重新运行 `go run user.go -f etc/user-api.yaml`

**理解**：Go 是编译型语言，代码修改后需要重新运行。

---

### 问题 5：数据库表字段和代码不一致

**现象**：
```
Error: Unknown column 'created_at' in 'field list'
```

**原因**：
- 代码中有 `created_at` 字段
- 但数据库表中没有这个列

**解决方案**：
1. 重新建表（最干净）
2. 或者修改表添加缺失字段

**最佳实践**：
- 先写 SQL 文件定义表结构
- 执行 SQL 创建表
- 用 goctl 从数据库生成 Model
- 这样保证代码和数据库 100% 一致

---

## 💡 核心学习要点

### 1. 分层的意义

**为什么要分层？**
- **职责分离**：每层只做自己的事
- **易于维护**：改数据库不影响业务逻辑
- **便于测试**：可以单独测试每一层
- **代码复用**：Model 方法可以被多个 Logic 调用

### 2. 代码生成的价值

**goctl 生成了什么？**
- 重复的框架代码（HTTP 解析、路由注册）
- 标准的项目结构
- 类型定义

**你需要做什么？**
- 填充业务逻辑（Logic 层）
- 添加自定义查询（Model 层）
- 设计接口规范（API 文件）

### 3. API 设计思维

**API 是对外的承诺**：
- 一旦发布，不能随意改动
- 要考虑扩展性（分页、筛选、排序）
- 要考虑安全性（不要暴露敏感字段）

**内部实现可以随便改**：
- Model 的 SQL 怎么写
- Logic 的业务逻辑怎么处理
- 只要 API 不变，外部无感知

### 4. 数据库管理

**学习阶段**：手动建表
**真实项目**：使用数据库迁移工具（golang-migrate）

**为什么？**
- 版本控制（每次变更都有记录）
- 可回滚（出问题可以退回）
- 团队协作（所有人执行同样的脚本）

---

## 🚀 快速开始

### 1. 环境准备
- Go 1.16+
- MySQL 5.7+
- goctl 工具

### 2. 数据库初始化

创建数据库：
```sql
CREATE DATABASE userdb CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

导入表结构（`user.sql`）：
```sql
USE userdb;

CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户姓名',
  `email` varchar(255) NOT NULL DEFAULT '' COMMENT '用户邮箱',
  `age` int NOT NULL DEFAULT 0 COMMENT '用户年龄',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 3. 配置修改

修改 `user/etc/user-api.yaml` 中的数据库连接信息：
```yaml
Mysql:
  DataSource: root:你的密码@tcp(127.0.0.1:3306)/userdb?charset=utf8mb4&parseTime=true&loc=Local
```

### 4. 启动服务
```bash
cd user
go run user.go -f etc/user-api.yaml
```

### 5. 测试接口

服务启动后访问 `http://localhost:8888`

---

## 📡 API 文档

| 功能 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 注册用户 | POST | `/user/register` | Body: `{"name": "张三", "email": "test@example.com"}` |
| 查询用户 | GET | `/user/info/:id` | 路径参数: `id`（如 `/user/info/1`） |
| 用户列表 | GET | `/user/list` | Query: `?page=1&page_size=10` |
| 更新用户 | PUT | `/user/:id` | 路径参数: `id` + Body: `{"name": "新名字", "email": "new@example.com"}` |
| 删除用户 | DELETE | `/user/:id` | 路径参数: `id`（如 `/user/1`） |

### 请求示例

**注册用户**：
```bash
curl -X POST http://localhost:8888/user/register \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","email":"zhangsan@example.com"}'
```

**查询用户**：
```bash
curl http://localhost:8888/user/info/1
```

**用户列表（分页）**：
```bash
curl http://localhost:8888/user/list?page=1&page_size=10
```

**更新用户**：
```bash
curl -X PUT http://localhost:8888/user/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"李四","email":"lisi@example.com"}'
```

**删除用户**：
```bash
curl -X DELETE http://localhost:8888/user/1
```

---

## 📂 项目结构

```
user-demo/
├── user/                    # 用户服务
│   ├── user.api            # API 定义文件
│   ├── user.go             # 服务入口
│   ├── etc/                # 配置文件
│   │   └── user-api.yaml
│   ├── internal/           # 内部代码
│   │   ├── config/         # 配置结构
│   │   ├── handler/        # HTTP 处理器
│   │   ├── logic/          # 业务逻辑
│   │   ├── svc/            # 服务上下文
│   │   └── types/          # 类型定义
│   └── model/              # 数据模型
│       ├── usermodel.go         # 自定义方法
│       └── usermodel_gen.go     # 自动生成
└── user.sql                # 数据库建表脚本
```

---

## 📖 延伸学习

**下一步可以学习**：
1. RPC 服务开发（微服务间通信）
2. 中间件使用（认证、日志、限流）
3. 缓存集成（Redis）
4. 服务发现（Etcd）
5. 链路追踪（Jaeger）

**推荐资源**：
- go-zero 官方文档：https://go-zero.dev/
- 官方示例项目：https://github.com/zeromicro/go-zero

---

## 📝 总结

通过本项目，你应该理解了：
- ✅ go-zero 的分层架构和每层职责
- ✅ API 文件的作用和字段映射
- ✅ 数据在各层之间的流动逻辑
- ✅ goctl 工具的使用和代码生成
- ✅ 常见错误的排查思路

**核心思想**：
- API 定义接口规范（对外）
- Logic 实现业务逻辑（核心）
- Model 封装数据操作（底层）
- 分层解耦，职责清晰

这就是现代 Web 开发的标准实践！🎉

---

## 📄 License

MIT
