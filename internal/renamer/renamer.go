package renamer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/king-joker-z/zspace-tools/internal/config"
)

// LogEntry 转换日志条目
type LogEntry struct {
	Time     string `json:"time"`      // 时间戳
	File     string `json:"file"`      // 文件路径
	From     string `json:"from"`      // 原后缀
	To       string `json:"to"`        // 新后缀
	Status   string `json:"status"`    // 状态：success / error
	Message  string `json:"message"`   // 详细信息
}

// Renamer 后缀转换引擎
type Renamer struct {
	logChan chan LogEntry // 日志通道
}

// NewRenamer 创建转换引擎
func NewRenamer(logChan chan LogEntry) *Renamer {
	return &Renamer{
		logChan: logChan,
	}
}

// RenameFile 对单个文件执行后缀转换
func (r *Renamer) RenameFile(filePath string, rules []config.Rule) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		return
	}

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		fromExt := strings.ToLower(rule.From)
		if !strings.HasPrefix(fromExt, ".") {
			fromExt = "." + fromExt
		}

		toExt := rule.To
		if !strings.HasPrefix(toExt, ".") {
			toExt = "." + toExt
		}

		if ext == fromExt {
			// 构建新文件名
			dir := filepath.Dir(filePath)
			baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
			newPath := filepath.Join(dir, baseName+toExt)

			// 检查目标文件是否已存在
			if _, err := os.Stat(newPath); err == nil {
				r.emitLog(filePath, fromExt, toExt, "skipped", "目标文件已存在: "+newPath)
				return
			}

			// 执行重命名
			if err := os.Rename(filePath, newPath); err != nil {
				r.emitLog(filePath, fromExt, toExt, "error", err.Error())
				return
			}

			r.emitLog(filePath, fromExt, toExt, "success",
				fmt.Sprintf("%s → %s", filepath.Base(filePath), filepath.Base(newPath)))
			return
		}
	}
}

// ScanDir 扫描目录并执行转换
func (r *Renamer) ScanDir(dir string, rules []config.Rule) (int, error) {
	count := 0

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过无法访问的文件
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		for _, rule := range rules {
			if !rule.Enabled {
				continue
			}
			fromExt := strings.ToLower(rule.From)
			if !strings.HasPrefix(fromExt, ".") {
				fromExt = "." + fromExt
			}
			if ext == fromExt {
				r.RenameFile(path, rules)
				count++
				break
			}
		}
		return nil
	})

	return count, err
}

// emitLog 发送日志
func (r *Renamer) emitLog(file, from, to, status, message string) {
	entry := LogEntry{
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		File:    file,
		From:    from,
		To:      to,
		Status:  status,
		Message: message,
	}

	// 非阻塞发送
	select {
	case r.logChan <- entry:
	default:
	}
}
