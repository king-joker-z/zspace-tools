package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/king-joker-z/zspace-tools/internal/api"
	"github.com/king-joker-z/zspace-tools/internal/config"
	"github.com/king-joker-z/zspace-tools/internal/renamer"
	"github.com/king-joker-z/zspace-tools/internal/watcher"
)

//go:embed web/dist/*
var webFS embed.FS

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("🚀 ZSpace Tools v0.1.0 启动中...")

	// 配置文件路径（支持环境变量覆盖）
	configPath := "/config/config.json"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		configPath = p
	}

	// 确保配置目录存在
	configDir := "/config"
	if d := os.Getenv("CONFIG_DIR"); d != "" {
		configDir = d
	}
	os.MkdirAll(configDir, 0755)

	// 初始化配置管理器
	cfgMgr, err := config.NewManager(configPath)
	if err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}
	log.Println("✅ 配置已加载")

	// 创建日志通道
	logChan := make(chan renamer.LogEntry, 1000)

	// 创建转换引擎
	rnm := renamer.NewRenamer(logChan)

	// 创建文件监听服务
	w := watcher.NewWatcher(rnm, cfgMgr)

	// 如果配置了自动监听，启动监听
	cfg := cfgMgr.Get()
	if cfg.AutoWatch {
		if err := w.Start(); err != nil {
			log.Printf("⚠️ 自动启动监听失败: %v", err)
		} else {
			log.Println("✅ 自动监听已启动")
		}
	}

	// 设置 Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// 注册 API 路由
	handler := api.NewHandler(cfgMgr, w, rnm, logChan)
	handler.RegisterRoutes(r)

	// 嵌入前端静态文件
	distFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Fatalf("加载前端资源失败: %v", err)
	}

	// 前端路由：先尝试静态文件，找不到则返回 index.html（SPA 支持）
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 尝试打开静态文件
		f, err := http.FS(distFS).Open(path)
		if err == nil {
			f.Close()
			c.FileFromFS(path, http.FS(distFS))
			return
		}

		// SPA fallback: 返回 index.html
		c.FileFromFS("/", http.FS(distFS))
	})

	// 启动服务器
	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	log.Printf("🌐 服务已启动，监听端口 %s", port)
	log.Printf("📖 访问 http://localhost%s 打开管理界面", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
