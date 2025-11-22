# User Demo - go-zero 微服务学习项目

这是一个基于 go-zero 框架的用户管理系统，用于学习微服务架构。项目经历了三个重要阶段的演进，每个阶段都解决了实际开发中会遇到的问题。

## 项目简介

本项目实现了用户管理的完整 CRUD 功能：
- 用户注册
- 查询用户信息
- 用户列表（支持分页）
- 更新用户
- 删除用户

**核心价值**：通过实际操作理解微服务拆分的本质，解答"为什么要拆分"、"怎么拆分"、"拆分后如何通信"等关键问题。

---

## 架构演进历程

### 阶段一：MySQL + sqlx 单体应用（master 分支）

**架构图：**
```
外部请求 → API Handler → Logic → Model (sqlx) → MySQL
```

**特点：**
- 所有功能在一个服务中
- Logic 直接调用 Model 操作数据库
- 简单直接，适合小型项目

### 阶段二：PostgreSQL + GORM 单体应用（postgres 分支）

**架构图：**
```
外部请求 → API Handler → Logic → Model (GORM) → PostgreSQL
```

**改动：**
- 数据库从 MySQL 切换到 PostgreSQL
- ORM 从 sqlx 切换到 GORM

**遇到的问题：**
- PostgreSQL 不支持 `LastInsertId()`，需要用 `RETURNING id`
- SQL 占位符从 `?` 改为 `$1, $2, $3`
- GORM 的用法与 sqlx 完全不同

**核心理解：** 不同数据库和 ORM 有不同的特性，但架构模式不变。

### 阶段三：微服务拆分（当前）

**架构图：**
```
外部请求 → user-api (HTTP Gateway)
              ↓ (gRPC 调用)
          user-rpc (业务服务)
              ↓
          PostgreSQL
```

**改动：**
- 新增 user-rpc：专门处理业务逻辑和数据库
- user 改名为 user-api：只负责接收 HTTP 请求
- API Logic 不再直接调用 Model，改为调用 RPC

**核心理解：** 这就是"多包了一层"，Logic 通过网络调用另一个服务，而不是直接调用本地方法。

---

## 目录结构

```
user-demo/
├── go.mod                 # 统一的 Go 模块管理（重要！）
├── start-services.sh      # 一键启动脚本
├── stop-services.sh       # 一键停止脚本
│
├── user-api/              # API 网关（端口 8888）
│   ├── etc/user-api.yaml  # 配置：包含 RPC 连接地址
│   ├── internal/
│   │   ├── handler/       # HTTP 请求处理
│   │   ├── logic/         # 业务逻辑（调用 RPC，不直接访问数据库）
│   │   └── svc/           # 服务上下文（持有 RPC 客户端）
│   └── user.go            # 主入口
│
└── user-rpc/              # RPC 业务服务（端口 8080）
    ├── etc/user.yaml      # 配置：包含数据库连接
    ├── user.proto         # Protobuf 服务定义（手写）
    ├── internal/
    │   ├── logic/         # 业务逻辑（操作数据库）
    │   ├── server/        # gRPC 服务器（自动生成）
    │   └── svc/           # 服务上下文（持有数据库连接）
    ├── model/             # 数据模型（从 user-api 复制过来的）
    ├── user/              # Protobuf 生成的代码
    ├── userclient/        # RPC 客户端代码（给 API 用）
    └── user.go            # 主入口
```

**重要说明：**
- 只有一个 `go.mod`（在根目录），两个服务共用
- Model 代码完全一样，只是位置从 user-api 移到了 user-rpc
- 90% 的 RPC 代码是自动生成的，你只需要填写业务逻辑

---

## 核心概念详解

### 1. 什么是微服务拆分？

**最简单的理解：把一个大服务拆成多个小服务。**

**单体应用：**
- 所有功能都在一起
- Logic 直接调用 Model
- 一个服务挂了，整个系统挂了

**微服务：**
- 功能分散在不同服务中
- Logic 通过网络调用其他服务
- 一个服务挂了，其他服务还能运行

**本项目的拆分：**
- **user-api**：只管接收 HTTP 请求，把请求转发给 RPC
- **user-rpc**：只管业务逻辑和数据库操作

### 2. Logic 为什么不直接调用 Model 了？

这是很多人的困惑点！

**之前（单体）：**
```
user-api Logic → 直接调用 Model.FindOne() → 数据库
```
这是**函数调用**，在同一个进程内。

**现在（微服务）：**
```
user-api Logic → 通过网络调用 RPC → user-rpc Logic → Model.FindOne() → 数据库
```
这是**网络调用**，跨越两个进程。

**关键理解：**
- Model 代码没有变，还是那些方法
- 只是 Model 不在 user-api 里了，搬到 user-rpc 里了
- user-api 通过网络调用 user-rpc，user-rpc 再调用 Model

**为什么要这样？**
- 职责分离：API 只管接口，RPC 只管业务
- 可扩展：可以启动多个 RPC 实例分担压力
- 可维护：改业务逻辑只需重启 RPC，不影响 API

### 3. "多包了一层"是什么意思？

**之前：**
```
HTTP 请求 → API Logic → Model → 数据库
```
2 层：Logic 直接调 Model

**现在：**
```
HTTP 请求 → API Logic → (网络调用) → RPC Logic → Model → 数据库
```
3 层：API Logic 先调 RPC Logic，RPC Logic 再调 Model

**这一层带来了什么？**
- ✅ 服务分离：可以独立部署
- ✅ 负载均衡：可以多实例
- ✅ 灵活扩展：改 RPC 不影响 API
- ❌ 增加复杂度：多了网络通信
- ❌ 性能损耗：网络调用比函数调用慢

**何时需要拆分？**
- 业务复杂，单体难以维护
- 流量大，需要独立扩展某个服务
- 团队大，多人协作需要拆分模块

**何时不需要拆分？**
- 小项目，几个接口而已
- 流量小，单机足够
- 团队小，维护成本高于收益

### 4. RPC 调用和函数调用有什么区别？

**函数调用（本地）：**
```
result := userModel.FindOne(ctx, id)
```
- 在同一个进程内
- 速度快（纳秒级）
- 不会失败（除非代码 bug）

**RPC 调用（远程）：**
```
result := userRpc.GetUser(ctx, &req)
```
- 跨进程，通过网络
- 速度慢（毫秒级）
- 可能失败（网络问题、服务挂了）

**虽然写法很像，但本质完全不同！**

RPC 调用实际上做了这些事：
1. 把请求参数序列化（变成字节流）
2. 通过 TCP 发送到 RPC 服务器
3. RPC 服务器解析请求
4. 执行业务逻辑
5. 把结果序列化
6. 通过 TCP 返回
7. 客户端解析结果

### 5. Model 代码会变吗？

**不会变！** 这是很重要的理解。

Model 层的代码完全一样，只是：
- **之前**：Model 在 `user/model/`，被 `user/logic/` 调用
- **现在**：Model 在 `user-rpc/model/`，被 `user-rpc/logic/` 调用

**本质：** Model 只是换了个调用者，从 API Logic 变成了 RPC Logic。

### 6. RPC 代码是自动生成的吗？

**90% 是自动生成的！**

**你需要做的：**
1. 编写 `user.proto` 文件（定义服务有哪些方法）
2. 运行 `goctl` 命令生成代码
3. 在 Logic 文件中填写业务逻辑（复制之前的代码）

**自动生成的：**
- gRPC Server 代码
- gRPC Client 代码
- Protobuf 消息定义
- Handler 框架
- 路由注册

**手动填写的：**
- Logic 中的业务逻辑（就是之前的代码）

---

## 服务发现：直连 vs etcd

这是另一个重要概念！

### 直连模式（本项目使用）

**原理：** 在配置文件中直接写死 RPC 服务的地址。

**配置示例：**
```yaml
# user-api/etc/user-api.yaml
UserRpc:
  Endpoints:
  - 127.0.0.1:8080  # 直接写 RPC 的地址
```

**工作流程：**
1. user-rpc 启动在 8080 端口
2. user-api 读配置：RPC 在 127.0.0.1:8080
3. user-api 直接连接这个地址

**优点：**
- 简单，无需额外组件
- 适合学习和测试

**缺点：**
- RPC 地址必须固定
- 不能随意更换机器
- 不支持负载均衡

**适用场景：**
- 同一台机器上的服务
- 固定 IP 的服务器
- 单实例部署

### etcd 服务发现模式

**原理：** RPC 启动时自动注册到 etcd，API 从 etcd 查询 RPC 在哪里。

**配置示例：**
```yaml
# user-rpc/etc/user.yaml
Etcd:
  Hosts:
  - 192.168.1.10:2379  # etcd 地址
  Key: user.rpc        # 服务名称

# user-api/etc/user-api.yaml
UserRpc:
  Etcd:
    Hosts:
    - 192.168.1.10:2379  # etcd 地址
    Key: user.rpc        # 要查询的服务名
```

**工作流程：**
1. user-rpc 启动 → 告诉 etcd："我叫 user.rpc，我在 192.168.1.20:8080"
2. etcd 记录：user.rpc = 192.168.1.20:8080
3. user-rpc 定期发心跳：我还活着
4. user-api 启动 → 问 etcd："user.rpc 在哪？"
5. etcd 返回：192.168.1.20:8080
6. user-api 连接到这个地址

**优点：**
- 动态服务发现（RPC 可以随便换机器）
- 支持多实例负载均衡
- 自动健康检查（RPC 挂了会自动摘除）
- 适合生产环境

**缺点：**
- 需要部署 etcd
- 配置相对复杂

**适用场景：**
- 生产环境
- 多机部署
- 需要高可用

### API 和 RPC 可以在不同机器吗？

**直连模式：** 可以，但需要修改配置中的 IP 地址。

**etcd 模式：** 完全可以，这就是服务发现的威力！

**举例（etcd 模式）：**
- etcd 部署在广州（192.168.1.10）
- user-rpc 部署在北京（192.168.1.20）
- user-api 部署在上海（192.168.1.30）

它们会自动找到彼此，无需手动配置 IP！

**核心理解：**
- 直连模式：你告诉 API："RPC 在 127.0.0.1:8080"
- etcd 模式：API 问 etcd："user.rpc 在哪？"，etcd 自动告诉它

---

## 启动方式

### 重要：必须先启动 RPC，再启动 API！

**为什么？**

因为 user-api 启动时会执行这段代码：
```
创建 RPC 客户端 → 尝试连接 127.0.0.1:8080
```

如果 RPC 还没启动：
- ❌ 连接失败
- ❌ API 启动可能报错或无法正常工作

**正确顺序：**
1. 先启动 user-rpc（监听 8080 端口）
2. 再启动 user-api（连接到 8080 端口）

### 方式一：一键启动（推荐）

```bash
cd /home/su/dome/go/go-zero-docs/test_demo/user-demo
./start-services.sh
```

这个脚本会：
1. 检查 PostgreSQL 是否运行
2. 先启动 user-rpc
3. 再启动 user-api
4. 输出日志文件位置

停止服务：
```bash
./stop-services.sh
```

### 方式二：手动启动（需要两个终端）

**终端 1 - 启动 RPC：**
```bash
cd user-rpc
go run user.go -f etc/user.yaml
```

看到这个表示成功：
```
Starting rpc server at 0.0.0.0:8080...
```

**终端 2 - 启动 API：**
```bash
cd user-api
go run user.go -f etc/user-api.yaml
```

看到这个表示成功：
```
Starting server at 0.0.0.0:8888...
```

---

## 如何测试

### 重要理解

**你只需要测试 user-api（8888 端口）！**

- ✅ 在 Apifox 中访问：`http://localhost:8888`
- ❌ 不需要直接访问 user-rpc（8080 端口）

**为什么？**

因为 user-rpc 是"内部服务"，只供 user-api 调用，不对外暴露。

**请求流程：**
```
Apifox 发送请求
    ↓
http://localhost:8888/user/list (user-api)
    ↓
user-api 内部通过 gRPC 调用 user-rpc
    ↓
127.0.0.1:8080 (user-rpc)
    ↓
user-rpc 查询数据库
    ↓
PostgreSQL
    ↓
原路返回数据
    ↓
Apifox 收到响应
```

**你在 Apifox 中感知不到 RPC 的存在！**

### API 接口

所有接口访问 `http://localhost:8888`：

1. **创建用户** - POST `/user/register`
2. **获取用户** - GET `/user/:id`
3. **更新用户** - PUT `/user`
4. **删除用户** - DELETE `/user/:id`
5. **用户列表** - GET `/user/list`

### 测试方式完全不变

虽然后台变成了微服务架构，但对外接口完全一样！

你在 Apifox 中的所有配置不需要改，继续用之前的就行。

---

## 常见问题（FAQ）

### 1. 为什么有导入错误："package user-demo/user-rpc/xxx is not in std"？

**原因：** user-rpc 最初有独立的 go.mod，模块名设置错了。

**解决：** 删除了 user-rpc 的独立 go.mod，使用项目根目录的统一 go.mod。

**理解：** 整个 user-demo 是一个 Go 模块，user-rpc 和 user-api 都是这个模块下的子目录，不需要各自的 go.mod。

### 2. 为什么必须先启动 RPC 再启动 API？

**原因：** API 启动时会立即尝试连接 RPC。

**类比：**
- RPC 是服务员
- API 是前台
- 前台上班时要确认服务员在岗，否则无法接待客人

**如果先启动 API：**
```
API：我要连接 RPC...
RPC：（还没启动）
API：连接失败！报错或无法正常工作
```

### 3. API 和 RPC 能在不同机器上运行吗？

**直连模式：** 可以，但要改配置中的 IP。

例如 RPC 在 192.168.1.20：
```yaml
UserRpc:
  Endpoints:
  - 192.168.1.20:8080  # 改成实际 IP
```

**etcd 模式：** 完全可以，不需要改配置！

只要所有服务都能连接到同一个 etcd，它们会自动互相发现，无论在哪台机器上。

### 4. 为什么要用 RPC 而不是直接调用 Model？

**核心原因：服务分离！**

**好处：**
- 职责清晰：API 管接口，RPC 管业务
- 独立部署：改业务逻辑只需重启 RPC
- 独立扩展：流量大就多开几个 RPC 实例
- 技术选型：RPC 可以用 Go，API 也可以用其他语言

**代价：**
- 多了网络通信开销
- 配置更复杂
- 运维成本增加

**什么时候拆分？**
- 业务复杂，代码太多
- 流量大，单机扛不住
- 团队大，需要分工

**什么时候不拆？**
- 小项目，简单的 CRUD
- 流量小，单机够用
- 团队小，维护成本高

### 5. 直连和 etcd 哪个更好？

**学习阶段：** 用直连，简单快速上手。

**生产环境：** 必须用 etcd 或类似的服务发现方案。

**原因：**
- 生产环境需要高可用（一个实例挂了，自动切换到其他实例）
- 需要负载均衡（多个实例分担流量）
- 需要动态扩缩容（流量大就加机器，流量小就减机器）

### 6. 90% 的代码是自动生成的，我要写什么？

**自动生成的：**
- gRPC Server 和 Client 代码
- Protobuf 消息定义
- HTTP Handler 框架
- 路由注册

**你要写的：**
- `user.proto`：定义服务有哪些方法
- Logic 业务逻辑：就是把之前 API Logic 的代码复制过来
- 配置文件：数据库连接、RPC 地址等

**理解：** 框架帮你搭好架子，你只需要填充业务逻辑。

### 7. 端口被占用怎么办？

**查看占用情况：**
```bash
lsof -i :8080  # 查看 8080 端口
lsof -i :8888  # 查看 8888 端口
```

**停止旧进程：**
```bash
./stop-services.sh
# 或者
pkill -f 'user.go'
```

### 8. PostgreSQL 连接失败怎么办？

**检查数据库是否运行：**
```bash
pg_isready -h 127.0.0.1 -p 5432
```

**检查配置是否正确：**
```yaml
Postgres:
  DataSource: postgres://用户名:密码@地址:端口/数据库?sslmode=disable
```

**常见错误：**
- 用户名密码错误
- 数据库不存在
- PostgreSQL 没启动

---

## 技术栈

- **框架：** go-zero
- **数据库：** PostgreSQL
- **ORM：** GORM
- **RPC：** gRPC + Protobuf
- **服务发现：** 直连模式（可选 etcd）

---

## 学习路径建议

### 第一步：理解单体应用
切换到 master 分支，理解最基础的 API → Logic → Model 流程。

### 第二步：理解数据库迁移
切换到 postgres 分支，理解如何从 MySQL 迁移到 PostgreSQL。

### 第三步：理解微服务拆分
回到当前分支，理解：
- 为什么要拆分？
- 怎么拆分？
- RPC 调用的原理？

### 第四步：实践服务发现
尝试安装 etcd，把配置改为 etcd 模式，体验动态服务发现。

### 第五步：多实例部署
启动多个 user-rpc 实例（不同端口），看 etcd 如何实现负载均衡。

---

## 下一步可以做什么？

1. **添加更多服务：** 比如订单服务 order-rpc，实现跨服务调用
2. **添加中间件：** 认证、日志、限流、熔断
3. **添加缓存：** Redis 缓存热点数据
4. **添加消息队列：** Kafka 处理异步任务
5. **容器化部署：** Docker + Docker Compose
6. **监控告警：** Prometheus + Grafana
7. **链路追踪：** Jaeger 追踪请求流转

---

## 核心总结

### 微服务拆分的本质

**不是为了拆分而拆分，而是为了解决问题：**
- 单体太大，难以维护 → 拆成多个小服务
- 流量太大，单机扛不住 → 独立扩展某个服务
- 团队太大，协作困难 → 分模块开发

### "多包了一层"的理解

**之前：** Logic 直接调 Model（函数调用）
**现在：** Logic 调 RPC（网络调用） → RPC 调 Model（函数调用）

这一层带来了服务分离的能力，代价是增加了复杂度。

### 服务发现的意义

**直连：** 写死地址，简单但不灵活
**etcd：** 动态发现，复杂但强大

生产环境必须用服务发现，否则无法实现高可用和负载均衡。

### 什么时候该微服务？

**需要拆分：**
- 代码太多（超过 1 万行）
- 流量太大（单机 QPS 超过 1000）
- 团队太大（超过 5 人）

**不需要拆分：**
- 小项目（几百行代码）
- 流量小（几十 QPS）
- 团队小（1-2 人）

**核心原则：** 架构为业务服务，不是为了技术而技术。

---

## 致谢

这个项目记录了从单体到微服务的完整演进过程，以及过程中遇到的真实问题和困惑。希望能帮助正在学习 go-zero 和微服务架构的你！

如果有疑问，建议：
1. 先看完整个 README
2. 对照代码理解每个概念
3. 动手实践，启动服务测试
4. 遇到问题回到 FAQ 找答案

**记住：** 理解比记忆更重要，实践比理论更有效。

---

🤖 本 README 根据实际开发过程中的对话整理而成，记录了所有的困惑点和解决方案。
