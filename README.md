# dnslog_for_go

---
[中文](README.CN.md) | English

--- 
## Features
- Lightweight deployment
- Docker support
- Automatic domain validity checking
- Simple and user-friendly web interface

--- 
## Project Structure
```
├── build/
│   └── docker/                 // Dockerfile and related files
├── cmd/
│   └── app/
│       └── main.go             // Main application entry point
├── internal/
│   └── config/                 // Configuration handling
│   └── domain/                 // Domain resolution logic
│   └── log/                    // Logging utilities
│   └── router/                 // Route definitions
│   └── web/
│       ├── templates/          // HTML templates
│       ├── static/             // Static assets (CSS, JS, images)
│       └── resources.go        // Embedded FS resources
├── pkg/                        // Shared utility packages
│   └── utils/                  
├── test/                       // Test files
│   └── utils/
├── go.mod                      // Go module definition
├── LICENSE                     // License file
├── README.CN.md                // Chinese documentation
├── README.md                   // English documentation
```

---

## Project Overview
`dnslog_for_go` is a simple DNSLog tool written in Go, with Docker support. It can be used for testing and debugging DNS-related applications, helping developers efficiently log and analyze DNS requests.

### Tech Stack
- [gin framework](https://github.com/gin-gonic/gin)
- [viper for configuration](https://github.com/spf13/viper)
- [miekg/dns library](https://github.com/miekg/dns)
- [zap](https://github.com/uber-go/zap)
- ['uuid' Automatically generate domain names](https://github.com/google/uuid)
- ['go embed' Embedding static resources](https://pkg.go.dev/embed)
- ['ini' Configuration Management](https://github.com/go-ini/ini/tree/v1.67.0)


---

## Requirements
- Go 1.20+
- Docker 1.12+
- MySQL 8.0+ (for DNS log persistence)

### Configuration
- File: `config/config.yaml` (see `config/config.example.yaml`)
- Key options:
  - `rootDomain` / `rootDomains`: one or multiple root domains to capture (e.g. `demo.com`, `["demo.com","example.com"]`)
  - `captureAll`: set `true` to log all DNS queries, regardless of domain
  - `dnsListenAddr`: DNS listen address (default `:15353`)
  - `httpListenAddr`: HTTP listen address (default `:8080`)
  - `upstreamDNS`: list of upstream DNS (e.g. `["8.8.8.8","223.5.5.5"]`)
  - `protocol`: `udp` / `tcp` (default `udp`)
  - `mysqlDSN`: e.g. `user:pass@tcp(localhost:3306)/dnslog?parseTime=true&loc=Local&charset=utf8mb4`
  - `pageSize`, `maxPageSize`: pagination defaults
- Env overrides (same keys in upper snake case):
  - `ROOT_DOMAIN`, `ROOT_DOMAINS`, `CAPTURE_ALL`
  - `DNS_LISTEN_ADDR`, `HTTP_LISTEN_ADDR`
  - `UPSTREAM_DNS`, `DNS_PROTOCOL`
  - `MYSQL_DSN`
  - `PAGE_SIZE`, `MAX_PAGE_SIZE`

---

## Contribution Guide
Contributions are welcome! To ensure quality and efficient collaboration, please follow the guidelines below when submitting issues and pull requests.

### Issue Guidelines
> Use Issues only to report bugs, suggest features, or provide design-related feedback.

- Please avoid submitting irrelevant content (e.g., “Thanks!”, or “How do I configure this on X?”). Such issues will be closed.
- Search existing issues before opening a new one to avoid duplication.
- When reporting bugs, try to provide the following information:
    - Operating system and version
    - Go version
    - Command used and logs
    - Screenshots or videos (if front-end related)

#### Example Issue Titles:
```markdown
🐞 WebSocket connection error when loading dnslog page
✨ Add support to export DNS logs to CSV
```

### Pull Request Guidelines
> Please follow the steps below when submitting a PR:

1. **Fork the repository** – Do not create branches directly on the main repo.
2. Create a new branch in your fork for development.
3. Ensure your code passes build and tests before submitting.
4. Each PR should address only one feature or bugfix.
5. **Commit message format**:

```markdown
[File/Module]: Description
```

#### Examples:
```markdown
README.md: Fix incorrect port in sample command
dnslog.go: Add support to export DNS logs to CSV
Dockerfile: Use multi-stage build to reduce image size
```

6. Include a clear description:
    - The purpose of the change
    - Files and modules affected
    - If it fixes an issue, reference it in the PR

---

### Contact
- You can open an Issue and leave your email there, or contact me directly via email.
- Email: *(fill in your contact)*

---

## Quick Start

### Docker Deployment

```bash
docker run --rm -p 8080:8080 <your-username>/dnslog-for-go:latest
```

After running, open your browser and visit:
```
http://localhost:8080/dnslog
```

---

### Local Deployment

#### Requirements:
- Go 1.20+ ([Go installation guide](https://golang.org/doc/install/source))
- Git ([Git installation guide](https://git-scm.com/))

#### Steps:
1. Clone the repository:
   ```bash
   git clone https://github.com/LianPeter/dnslog_for_go.git
   ```

2. Enter the project directory:
   ```bash
   cd dnslog_for_go
   ```

3. Download dependencies:
   ```bash
   go mod download
   ```

4. Run the project:
   ```bash
   go run main.go
   ```

5. Access the web UI:
   ```
   http://localhost:8080/dnslog
   ```

---

### 🐳 Docker Build

If you want to build your own image:

```bash
docker build -t dnslog-for-go .
```

To run and test:

```bash
docker run --rm -p 8080:8080 dnslog-for-go
```

---

## TODO
- [x] Automatically generate dnslog domains
- [x] Support replacement of dns server
- [ ] Connect to a database for persistent storage (TBD)
- [ ] Improve front-end UI
- [ ] Provide a native client UI (without relying on browser)

---