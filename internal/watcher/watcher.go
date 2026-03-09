package watcher

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/king-joker-z/zspace-tools/internal/config"
	"github.com/king-joker-z/zspace-tools/internal/renamer"
)

// Watcher 文件监听服务
type Watcher struct {
	mu        sync.Mutex
	fsWatcher *fsnotify.Watcher
	renamer   *renamer.Renamer
	cfgMgr    *config.Manager
	running   bool
	stopChan  chan struct{}
}

// NewWatcher 创建文件监听服务
func NewWatcher(r *renamer.Renamer, cfgMgr *config.Manager) *Watcher {
	return &Watcher{
		renamer: r,
		cfgMgr:  cfgMgr,
	}
}

// Start 启动文件监听
func (w *Watcher) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return nil
	}

	fsW, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	w.fsWatcher = fsW
	w.stopChan = make(chan struct{})
	w.running = true

	// 添加监听目录
	cfg := w.cfgMgr.Get()
	for _, dir := range cfg.WatchDirs {
		if err := fsW.Add(dir); err != nil {
			log.Printf("[监听] 添加目录失败 %s: %v", dir, err)
		} else {
			log.Printf("[监听] 已添加目录: %s", dir)
		}
	}

	// 启动事件处理协程
	go w.eventLoop()

	log.Println("[监听] 服务已启动")
	return nil
}

// Stop 停止文件监听
func (w *Watcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return
	}

	close(w.stopChan)
	w.fsWatcher.Close()
	w.running = false
	log.Println("[监听] 服务已停止")
}

// IsRunning 返回监听状态
func (w *Watcher) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}

// Restart 重启监听（用于配置变更后）
func (w *Watcher) Restart() error {
	w.Stop()
	return w.Start()
}

// eventLoop 事件处理循环
func (w *Watcher) eventLoop() {
	for {
		select {
		case <-w.stopChan:
			return
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			// 只处理创建和重命名事件（新文件出现）
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Rename) {
				cfg := w.cfgMgr.Get()
				w.renamer.RenameFile(event.Name, cfg.Rules)
			}
		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("[监听] 错误: %v", err)
		}
	}
}
