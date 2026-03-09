package api

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/king-joker-z/zspace-tools/internal/config"
	"github.com/king-joker-z/zspace-tools/internal/renamer"
	"github.com/king-joker-z/zspace-tools/internal/watcher"
)

// Handler API 路由处理器
type Handler struct {
	cfgMgr   *config.Manager
	watcher  *watcher.Watcher
	rnm      *renamer.Renamer
	logChan  chan renamer.LogEntry
	logs     []renamer.LogEntry
	logMu    sync.RWMutex
	wsClients map[*websocket.Conn]bool
	wsMu     sync.Mutex
}

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（内网使用）
	},
}

// NewHandler 创建 API 处理器
func NewHandler(cfgMgr *config.Manager, w *watcher.Watcher, r *renamer.Renamer, logChan chan renamer.LogEntry) *Handler {
	h := &Handler{
		cfgMgr:    cfgMgr,
		watcher:   w,
		rnm:       r,
		logChan:   logChan,
		logs:      make([]renamer.LogEntry, 0, 1000),
		wsClients: make(map[*websocket.Conn]bool),
	}

	// 启动日志收集协程
	go h.collectLogs()

	return h
}

// RegisterRoutes 注册 API 路由
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/config", h.getConfig)
		api.PUT("/config", h.updateConfig)
		api.POST("/scan", h.manualScan)
		api.GET("/status", h.getStatus)
		api.GET("/logs", h.getLogs)
		api.GET("/ws/logs", h.wsLogs)
		api.POST("/watch/start", h.startWatch)
		api.POST("/watch/stop", h.stopWatch)
	}
}

// getConfig 获取当前配置
func (h *Handler) getConfig(c *gin.Context) {
	cfg := h.cfgMgr.Get()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": cfg,
	})
}

// updateConfig 更新配置
func (h *Handler) updateConfig(c *gin.Context) {
	var cfg config.Config
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := h.cfgMgr.Update(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2,
			"message": "保存配置失败: " + err.Error(),
		})
		return
	}

	// 如果监听正在运行，重启以应用新配置
	if h.watcher.IsRunning() {
		if err := h.watcher.Restart(); err != nil {
			log.Printf("[API] 重启监听失败: %v", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "配置已保存",
	})
}

// manualScan 手动触发扫描
func (h *Handler) manualScan(c *gin.Context) {
	cfg := h.cfgMgr.Get()
	totalCount := 0

	for _, dir := range cfg.WatchDirs {
		count, err := h.rnm.ScanDir(dir, cfg.Rules)
		if err != nil {
			log.Printf("[API] 扫描目录 %s 失败: %v", dir, err)
		}
		totalCount += count
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "扫描完成",
		"data": gin.H{
			"processed": totalCount,
		},
	})
}

// getStatus 获取监听状态
func (h *Handler) getStatus(c *gin.Context) {
	cfg := h.cfgMgr.Get()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"watching":   h.watcher.IsRunning(),
			"watch_dirs": cfg.WatchDirs,
			"rules":      len(cfg.Rules),
			"log_count":  len(h.logs),
		},
	})
}

// getLogs 获取日志（支持分页）
func (h *Handler) getLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 200 {
		size = 50
	}

	h.logMu.RLock()
	total := len(h.logs)

	// 倒序返回（最新的在前）
	start := total - page*size
	end := total - (page-1)*size

	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}

	// 取出切片并翻转
	slice := make([]renamer.LogEntry, 0)
	if start < end {
		for i := end - 1; i >= start; i-- {
			slice = append(slice, h.logs[i])
		}
	}
	h.logMu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total": total,
			"page":  page,
			"size":  size,
			"items": slice,
		},
	})
}

// startWatch 启动监听
func (h *Handler) startWatch(c *gin.Context) {
	if err := h.watcher.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2,
			"message": "启动监听失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "监听已启动",
	})
}

// stopWatch 停止监听
func (h *Handler) stopWatch(c *gin.Context) {
	h.watcher.Stop()
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "监听已停止",
	})
}

// wsLogs WebSocket 实时日志推送
func (h *Handler) wsLogs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] 升级失败: %v", err)
		return
	}

	h.wsMu.Lock()
	h.wsClients[conn] = true
	h.wsMu.Unlock()

	// 发送最近 20 条日志作为初始数据
	h.logMu.RLock()
	start := len(h.logs) - 20
	if start < 0 {
		start = 0
	}
	recent := make([]renamer.LogEntry, len(h.logs)-start)
	copy(recent, h.logs[start:])
	h.logMu.RUnlock()

	for _, entry := range recent {
		conn.WriteJSON(entry)
	}

	// 保持连接，监听关闭
	defer func() {
		h.wsMu.Lock()
		delete(h.wsClients, conn)
		h.wsMu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// collectLogs 收集日志并广播到 WebSocket 客户端
func (h *Handler) collectLogs() {
	for entry := range h.logChan {
		// 存储日志
		h.logMu.Lock()
		h.logs = append(h.logs, entry)
		// 保留最近 10000 条
		if len(h.logs) > 10000 {
			h.logs = h.logs[len(h.logs)-10000:]
		}
		h.logMu.Unlock()

		// 广播到 WebSocket 客户端
		h.wsMu.Lock()
		for conn := range h.wsClients {
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := conn.WriteJSON(entry); err != nil {
				conn.Close()
				delete(h.wsClients, conn)
			}
		}
		h.wsMu.Unlock()
	}
}
