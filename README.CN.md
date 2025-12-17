# dnslog_for_go

---
中文 | [English](README.md)

---

## 特性
- 轻量级部署
- 支持 Docker 部署
- 自动校验域名合法性
- 简洁实用的 Web 界面交互

--- 
## 目录结构
```
dnslog_for_go/
├── build/
│   └──docker/                  // dockerfile
├── cmd/
│   └── app/
│       └── main.go             // 主程序
├── internal/
│   └── config/                 // 配置文件
│   └── domain/                 // 域名解析
│   └── log/                    // 日志相关
│   └── router/                 // 路由
│   └── web/
│       ├── templates/          // 模板文件
│       ├── static/             // 静态资源
│       └──resources.go        // embed.FS 资源
├── pkg/                        // 公共包
│   └── utils/                  
│
│
├── test/                       // 测试
│   └── utils/
│
├── go.mod                      // go mod 文件
├── LICENSE                     // 许可证
├── README.CN.md
├── README.md
```

---

## 项目介绍
`dnslog_for_go` 是一个简单的 DNSLog 工具，使用 Go 语言编写，支持 Docker 部署。它可以用于测试和调试 DNS 相关的应用程序，帮助开发者快速进行 DNS 请求的日志记录与分析。

### 使用的技术栈
- [gin 框架](https://github.com/gin-gonic/gin)
- [viper 配置管理](https://github.com/spf13/viper)
- [miekg/dns 库](https://github.com/miekg/dns)
- [zap](https://github.com/uber-go/zap)
- [uuid 自动生成域名](https://github.com/google/uuid)
- [go embed 嵌入静态资源](https://pkg.go.dev/embed) 
- [ini配置管理](https://github.com/go-ini/ini/tree/v1.67.0)

---

## 使用说明
- Go 1.20+ 环境
- Docker 1.12+ 环境
- MySQL 8.0+（日志持久化）

### 配置
- 配置文件：`config/config.yaml`（示例见 `config/config.example.yaml`）
- 重要字段：
  - `rootDomain` / `rootDomains`：单个或多个根域名（例如 `demo.com` 或 `["demo.com","example.com"]`）
  - `captureAll`：设为 `true` 时记录所有域名请求，不再限制根域
  - `dnsListenAddr`：DNS 监听地址（默认 `:15353`）
  - `httpListenAddr`：HTTP 监听地址（默认 `:8080`）
  - `upstreamDNS`：上游 DNS 列表（例如 `["8.8.8.8","223.5.5.5"]`）
  - `protocol`：`udp` / `tcp`（默认 `udp`）
  - `mysqlDSN`：如 `user:pass@tcp(localhost:3306)/dnslog?parseTime=true&loc=Local&charset=utf8mb4`
  - `pageSize` / `maxPageSize`：分页默认与上限
- 环境变量可覆盖（同名大写）：
  - `ROOT_DOMAIN`、`ROOT_DOMAINS`、`CAPTURE_ALL`
  - `DNS_LISTEN_ADDR`、`HTTP_LISTEN_ADDR`
  - `UPSTREAM_DNS`、`DNS_PROTOCOL`
  - `MYSQL_DSN`
  - `PAGE_SIZE`、`MAX_PAGE_SIZE`

### 贡献指南
欢迎大家参与贡献！为了保证项目的质量和协作效率，请遵循以下规范提交 Issue 与 Pull Request。

#### Issue 提交规范
> 仅用于报告 Bug、建议 Feature 或提交设计相关内容。

- 请勿提交无关内容（如“感谢作者”、“求问某环境配置”等），这些 Issue 会被关闭。
- 提交前请先 **搜索** 是否已有相关内容，避免重复提问。
- Bug 提交时，请尽量提供以下信息：
    - 操作系统与版本
    - Go 版本
    - 运行命令与日志输出
    - 若涉及前端，提供截图或视频辅助

##### 示例标题：
```markdown
🐞 dnslog 页面加载报错：WebSocket 无法连接
✨ 支持将 DNS 请求结果导出为 CSV
```

#### Pull Request 提交规范
> 所有 PR 请遵循以下流程提交：

1. **Fork 仓库**：请不要直接在主仓库上创建分支。
2. 在你的仓库中创建新的分支进行开发。
3. 提交前请确保本地代码已通过构建与测试。
4. 每个 PR **仅包含一个功能/问题修复**，避免混杂。
5. **Commit 信息规范**：

格式：
```markdown
[文件/模块名]: 描述信息
```

##### 示例：
```markdown
README.md: 修复示例命令中的端口错误
dnslog.go: 新增导出 CSV 功能
Dockerfile: 更新为多阶段构建，减少镜像体积
```

6. **PR 描述清晰**：
    - 该变更的目的
    - 涉及的文件和模块
    - 是否修复了某个 Issue（建议在描述中引用）

---

### 联系我：
- 可以在 Issue 中附上你的问题和邮箱，或通过邮箱直接联系我。
- 邮箱：
---

## 快速开始

### Docker 部署

```bash
docker run --rm -p 8080:8080 <你的用户名>/dnslog-for-go:latest
```

运行后，打开浏览器访问：
```
http://localhost:8080/dnslog
```

### 本地部署

#### 环境要求：
- Go 1.20+，如果没有，请参考 [Go 安装指南](https://golang.org/doc/install/source)
- Git 环境，若没有，请参考 [Git 安装指南](https://git-scm.com/)

#### 部署步骤：
1. 克隆项目：
   ```bash
   git clone https://github.com/LianPeter/dnslog_for_go.git
   ```

2. 进入项目目录：
   ```bash
   cd dnslog_for_go
   ```

3. 下载依赖：
   ```bash
   go mod download
   ```

4. 运行项目：
   ```bash
   go run main.go
   ```

5. 访问 Web 界面：
   ```
   http://localhost:8080/dnslog
   ```

### Docker 🐳 构建

如果你想构建自己的镜像：

```bash
docker build -t dnslog-for-go .
```

运行测试：

```bash
docker run --rm -p 8080:8080 dnslog-for-go
```

---

## TODO
- [x] 自动生成dnslog 域名
- [x] 支持更换dns服务器
- [ ] 连接数据库，持久化存储（暂定）
- [ ] 美化前端界面
- [ ] 拥有客户端界面，不依赖于浏览器

