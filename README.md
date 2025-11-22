# Docker 容器化部署问题记录

本文档记录了将 go-zero 微服务项目容器化部署过程中遇到的所有问题和解决方案，帮助理解 Docker 网络、容器通信和配置管理。

---

## 一、项目背景

### 初始状态
- 项目已完成开发：user-api（HTTP 网关）+ user-rpc（gRPC 服务）
- 数据库：PostgreSQL
- 本地已有一个运行中的 PostgreSQL 容器：`my_postgres`（端口 5432）
- 目标：将服务容器化，方便部署和测试

### 核心需求
1. 将 user-api 和 user-rpc 打包成 Docker 镜像
2. 通过 docker-compose 一键启动所有服务
3. 服务间能正常通信
4. 可以在 Apifox 中测试 API

---

## 二、遇到的问题及解决方案

### 问题 1：Go 版本不匹配

**现象：**
```
go: go.mod requires go >= 1.25.2 (running go 1.21.13; GOTOOLCHAIN=local)
```

**原因分析：**
- Dockerfile 中使用 `golang:1.21-alpine` 镜像
- go.mod 文件要求 Go 1.25.2
- Go 官方最新版本是 1.24，1.25.2 不存在

**解决方案：**
1. 修改 Dockerfile：将 `golang:1.21-alpine` 改为 `golang:1.24-alpine`
2. 修改 go.mod：将 `go 1.25.2` 改为 `go 1.24`

**关键理解：**
- Docker 镜像的 Go 版本必须满足 go.mod 的要求
- go.mod 中的版本号必须是真实存在的 Go 版本

---

### 问题 2：go mod 依赖需要更新

**现象：**
```
go: updates to go.mod needed; to update it:
  go mod tidy
```

**原因分析：**
构建时直接编译，但 go.mod 和 go.sum 与实际代码依赖不同步

**解决方案：**
在 Dockerfile 的编译步骤前添加 `go mod tidy`

**关键理解：**
- `go mod tidy` 会清理未使用的依赖，添加缺失的依赖
- 在容器构建时需要确保依赖完整

---

### 问题 3：PostgreSQL 端口冲突

**现象：**
```
Bind for 0.0.0.0:5432 failed: port is already allocated
```

**原因分析：**
- docker-compose.yaml 中定义了新的 PostgreSQL 容器，端口 5432
- 系统中已有运行中的 `my_postgres` 容器占用了 5432 端口

**决策点：是否需要独立数据库？**

**方案 A：每个项目一个独立 PostgreSQL**
- 优点：隔离性好，不同项目互不影响
- 缺点：占用资源多，管理复杂

**方案 B：共用现有 PostgreSQL**
- 优点：节省资源，数据集中管理
- 缺点：需要处理容器间通信

**最终选择：方案 B**

理由：
1. 测试阶段不需要完全隔离
2. 可以在 Navicat 中统一查看所有数据
3. 节省系统资源

**解决方案：**
1. 删除 docker-compose.yaml 中的 postgres 服务定义
2. 修改 user-rpc 配置，连接到现有的 `my_postgres`

---

### 问题 4：容器间网络通信失败（第一次尝试）

**现象：**
```
invalid config for network bridge: invalid endpoint settings:
network-scoped alias is supported only for containers in user defined networks
```

**尝试的方案：**
使用 `networks: default` 并设置 `external: true, name: bridge`

**原因分析：**
Docker 默认的 bridge 网络不支持容器别名（alias），无法通过容器名互相访问

**关键理解：**
- 默认 bridge 网络：容器只能通过 IP 通信，不能用容器名
- 自定义 bridge 网络：支持容器名解析（DNS）
- 容器别名是自定义网络的特性

---

### 问题 5：容器无法通过容器名通信

**现象：**
```
rpc error: desc = "transport: Error while dialing: dial tcp: lookup user-rpc on 192.168.65.7:53: no such host"
```

**原因分析：**
- user-api 尝试通过 `user-rpc:8080` 连接 RPC 服务
- 在默认 bridge 网络中无法解析容器名

**尝试的方案：**
使用 `network_mode: host`，让容器直接使用宿主机网络

**问题：**
在 WSL2 + Docker Desktop 环境下，`network_mode: host` 不生效，端口无法映射到宿主机

**关键理解：**
- `network_mode: host` 在 Linux 原生环境有效
- WSL2 使用虚拟化，Docker Desktop 运行在虚拟机中
- host 模式在这种环境下无法正常工作

---

### 问题 6：如何让容器访问宿主机服务

**需求：**
- user-rpc 需要访问宿主机上的 `my_postgres`（5432）
- user-api 需要访问 user-rpc（8080）

**最终方案：使用 `host.docker.internal`**

**配置要点：**

1. **docker-compose.yaml 配置：**
   - 使用端口映射（`ports: 8080:8080`, `8888:8888`）
   - 添加 `extra_hosts` 配置 `host.docker.internal`

2. **服务配置文件：**
   - user-rpc 连接数据库：`host.docker.internal:5432`
   - user-api 连接 RPC：`host.docker.internal:8080`

**工作原理：**
```
宿主机 (localhost)
  ├─ my_postgres (5432)
  ├─ user-rpc 容器 (8080) → 连接 host.docker.internal:5432
  └─ user-api 容器 (8888) → 连接 host.docker.internal:8080
```

**关键理解：**
- `host.docker.internal` 是 Docker Desktop 提供的特殊域名
- 它会自动解析为宿主机的 IP 地址
- 通过端口映射，容器服务暴露到宿主机端口
- 其他容器通过 `host.docker.internal:端口` 访问

---

## 三、最终架构方案

### 容器架构

```
┌─────────────────────────────────────────────┐
│            宿主机 (Windows/WSL2)              │
│                                              │
│  ┌────────────────────────────────────────┐ │
│  │   my_postgres 容器                      │ │
│  │   端口: 5432                            │ │
│  │   数据库: dex_db                        │ │
│  └────────────────────────────────────────┘ │
│                    ↑                         │
│                    │ (host.docker.internal)  │
│  ┌────────────────────────────────────────┐ │
│  │   user-rpc 容器                         │ │
│  │   端口: 8080 → 8080                     │ │
│  │   连接: host.docker.internal:5432       │ │
│  └────────────────────────────────────────┘ │
│                    ↑                         │
│                    │ (host.docker.internal)  │
│  ┌────────────────────────────────────────┐ │
│  │   user-api 容器                         │ │
│  │   端口: 8888 → 8888                     │ │
│  │   连接: host.docker.internal:8080       │ │
│  └────────────────────────────────────────┘ │
│                    ↑                         │
└────────────────────┼─────────────────────────┘
                     │
              Apifox 测试工具
           (http://localhost:8888)
```

### 配置文件关键点

**docker-compose.yaml：**
- 移除独立的 PostgreSQL 定义
- 使用端口映射而非 host 网络
- 添加 `extra_hosts` 支持 `host.docker.internal`

**user-rpc/etc/user.yaml：**
- 数据库连接：`host.docker.internal:5432`

**user-api/etc/user-api.yaml：**
- RPC 连接：`host.docker.internal:8080`

---

## 四、核心知识点总结

### 1. Docker 网络模式对比

| 网络模式 | 容器间通信 | 容器访问宿主机 | 宿主机访问容器 | 适用场景 |
|---------|----------|--------------|--------------|---------|
| bridge (默认) | 仅 IP，不支持容器名 | 不便 | 需端口映射 | 简单场景 |
| 自定义 bridge | 支持容器名解析 | 不便 | 需端口映射 | 多容器项目 |
| host | N/A（共享网络） | 直接访问 | 直接访问 | Linux 原生 |
| host.docker.internal | 通过宿主机端口 | 直接访问 | 需端口映射 | WSL2/Mac |

### 2. 容器间通信方案选择

**场景 A：所有服务都在容器中**
- 使用自定义 bridge 网络
- 通过容器名通信（如 `user-rpc:8080`）

**场景 B：部分服务在宿主机**
- 使用 `host.docker.internal`
- 通过宿主机端口通信

**场景 C：生产环境**
- 使用 Kubernetes 或 Docker Swarm
- 服务发现和负载均衡

### 3. 为什么选择共用数据库

**技术角度：**
1. 容器销毁不影响数据（数据持久化在宿主机）
2. 减少网络层级（少一次容器通信）
3. 便于调试（Navicat 直连查看）

**管理角度：**
1. 统一备份和管理
2. 节省资源
3. 简化配置

**何时需要独立数据库：**
1. 生产环境（隔离性）
2. 多租户场景
3. 数据安全要求高

---

## 五、常见问题排查

### 问题：容器启动但无法访问

**排查步骤：**
1. 检查容器状态：`docker ps`
2. 查看容器日志：`docker logs <容器名>`
3. 检查端口映射：确认端口是否暴露
4. 检查防火墙：WSL2 可能需要配置防火墙

### 问题：数据库连接失败

**排查步骤：**
1. 确认数据库容器运行中
2. 检查配置文件中的连接字符串
3. 验证数据库用户名、密码
4. 确认数据库已创建

### 问题：RPC 调用失败

**排查步骤：**
1. 确认 user-rpc 已启动
2. 检查 user-rpc 端口是否映射
3. 查看 user-api 配置中的 RPC 地址
4. 检查日志中的具体错误

---

## 六、经验教训

### 1. 环境差异要注意

- Linux 原生 Docker 和 Docker Desktop 行为不同
- WSL2 的网络是虚拟化的，不是真正的 Linux
- 跨平台部署时要考虑网络模式兼容性

### 2. 配置文件很关键

- Dockerfile 的 Go 版本要与 go.mod 匹配
- 配置文件（yaml）会打包进镜像，修改后需重新构建
- 使用环境变量可以避免频繁重新构建

### 3. 循序渐进调试

- 不要一次性改太多配置
- 每次修改后验证日志
- 理解问题根因，不要盲目尝试

### 4. 文档化很重要

- 记录每个决策的原因
- 记录遇到的问题和解决方案
- 帮助团队成员理解架构

---

## 七、下一步优化方向

### 短期优化
1. **使用环境变量**：数据库连接信息改为环境变量，避免硬编码
2. **健康检查**：添加容器健康检查，确保服务正常启动
3. **日志管理**：配置日志输出，便于排查问题

### 长期优化
1. **完全容器化**：将 PostgreSQL 也纳入 docker-compose，使用 volume 持久化
2. **服务发现**：引入 etcd 或 Consul 实现动态服务发现
3. **CI/CD**：自动化构建和部署流程
4. **监控告警**：添加 Prometheus + Grafana 监控

---

## 八、总结

本次容器化过程的核心收获：

1. **理解 Docker 网络**：不同网络模式适用不同场景
2. **灵活应对环境差异**：WSL2 需要特殊处理
3. **资源合理利用**：不必为每个项目创建独立数据库
4. **配置管理的重要性**：配置文件变更需重新构建镜像

关键思路：
- 先理解问题本质，再选择解决方案
- 不要被工具限制思路，灵活组合使用
- 文档化决策过程，帮助他人理解

最终实现：
- 通过 `host.docker.internal` 优雅解决容器访问宿主机服务
- 共用 PostgreSQL 节省资源
- 配置清晰，易于维护
