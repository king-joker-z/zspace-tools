# ZSpace Tools 🛠️

极空间 NAS 文件管理工具集。当前版本 (v0.1.0) 提供 **文件后缀批量转换 + 实时监听** 功能。

## ✨ 功能特性

- 📁 **文件夹监听**：配置监听路径，新文件自动按规则转换后缀
- 🔄 **后缀转换**：灵活配置转换规则（如 `.gif` → `.jpg`，仅改后缀不转格式）
- 🔍 **手动扫描**：一键扫描目录，批量转换已有文件
- 📋 **实时日志**：WebSocket 推送转换日志，操作可追溯
- 💾 **配置持久化**：JSON 配置文件，映射到宿主机 volume
- 📱 **响应式 UI**：Apple HIG 风格，手机平板电脑均可使用

## 🚀 快速开始

### Docker Compose（推荐）

```yaml
version: '3'
services:
  zspace-tools:
    image: xieyurong/zspace-tools:latest
    container_name: zspace-tools
    ports:
      - "8080:8080"
    volumes:
      - ./config:/config        # 配置持久化
      - /path/to/watch:/watch   # 替换为你的监听目录
    environment:
      - TZ=Asia/Shanghai
    restart: unless-stopped
```

```bash
docker-compose up -d
```

### Docker Run

```bash
docker run -d \
  --name zspace-tools \
  -p 8080:8080 \
  -v $(pwd)/config:/config \
  -v /path/to/watch:/watch \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  xieyurong/zspace-tools:latest
```

### 访问管理界面

打开浏览器访问 `http://<NAS-IP>:8080`

## 🏗️ 架构

```
zspace-tools/
├── main.go                 # 入口，启动 API + 监听
├── internal/
│   ├── watcher/            # fsnotify 文件监听服务
│   ├── renamer/            # 后缀转换引擎
│   ├── config/             # 配置管理（读写 JSON）
│   └── api/                # REST API + WebSocket
├── web/                    # Vue 3 前端（构建后嵌入 Go 二进制）
├── Dockerfile              # 多阶段构建
└── docker-compose.yml
```

## 📡 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/config` | 获取配置 |
| PUT | `/api/config` | 更新配置 |
| POST | `/api/scan` | 手动触发扫描 |
| GET | `/api/status` | 获取监听状态 |
| GET | `/api/logs` | 获取日志（支持分页） |
| WS | `/api/ws/logs` | WebSocket 实时日志 |
| POST | `/api/watch/start` | 启动监听 |
| POST | `/api/watch/stop` | 停止监听 |

## 🛠️ 本地开发

### 前置要求

- Go 1.22+
- Node.js 20+

### 构建步骤

```bash
# 构建前端
cd web && npm install --include=dev && npx vite build && cd ..

# 构建后端（前端资源通过 go:embed 嵌入）
go build -o zspace-tools .

# 运行
CONFIG_DIR=/tmp/config CONFIG_PATH=/tmp/config/config.json ./zspace-tools
```

### Docker 构建

```bash
docker build -t xieyurong/zspace-tools:latest .
```

## 📋 技术栈

- **后端**: Go + Gin + fsnotify + gorilla/websocket
- **前端**: Vue 3 + Vite（go:embed 嵌入）
- **容器**: Alpine Linux，镜像 < 20MB
- **架构**: 支持 AMD64 + ARM64（极空间 NAS 兼容）

## 📄 License

MIT
