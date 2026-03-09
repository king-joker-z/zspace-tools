package config

import (
	"encoding/json"
	"os"
	"sync"
)

// Rule 定义一条后缀转换规则
type Rule struct {
	From    string `json:"from"`    // 原始后缀，如 ".gif"
	To      string `json:"to"`      // 目标后缀，如 ".jpg"
	Enabled bool   `json:"enabled"` // 是否启用
}

// Config 应用配置
type Config struct {
	WatchDirs []string `json:"watch_dirs"` // 监听的文件夹列表
	Rules     []Rule   `json:"rules"`      // 转换规则列表
	AutoWatch bool     `json:"auto_watch"` // 是否自动开启监听
}

// Manager 配置管理器
type Manager struct {
	mu       sync.RWMutex
	config   Config
	filePath string
}

// NewManager 创建配置管理器
func NewManager(filePath string) (*Manager, error) {
	m := &Manager{
		filePath: filePath,
		config: Config{
			WatchDirs: []string{"/watch"},
			Rules:     []Rule{},
			AutoWatch: false,
		},
	}

	// 尝试加载已有配置
	if err := m.load(); err != nil {
		// 文件不存在则创建默认配置
		if os.IsNotExist(err) {
			return m, m.Save()
		}
		return nil, err
	}

	return m, nil
}

// load 从文件加载配置
func (m *Manager) load() error {
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &m.config)
}

// Save 保存配置到文件
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0644)
}

// Get 获取当前配置（返回副本）
func (m *Manager) Get() Config {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cfg := Config{
		AutoWatch: m.config.AutoWatch,
		WatchDirs: make([]string, len(m.config.WatchDirs)),
		Rules:     make([]Rule, len(m.config.Rules)),
	}
	copy(cfg.WatchDirs, m.config.WatchDirs)
	copy(cfg.Rules, m.config.Rules)
	return cfg
}

// Update 更新配置
func (m *Manager) Update(cfg Config) error {
	m.mu.Lock()
	m.config = cfg
	m.mu.Unlock()
	return m.Save()
}
