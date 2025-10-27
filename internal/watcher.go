package internal

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// fsnotifyWatcher fsnotify 文件监听器实现
type fsnotifyWatcher struct {
	watcher *fsnotify.Watcher
	path    string
	closed  bool
}

// NewFileWatcher 创建文件监听器
func NewFileWatcher(path string, reloadCallback func() error) FileWatcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("[i18n] Watcher init failed: %v", err)
		return &NoOpWatcher{}
	}

	err = watcher.Add(path)
	if err != nil {
		log.Printf("[i18n] Add watcher failed: %v", err)
		watcher.Close()
		return &NoOpWatcher{}
	}

	fw := &fsnotifyWatcher{
		watcher: watcher,
		path:    path,
	}

	// 启动监听协程
	go fw.watch(reloadCallback)

	log.Printf("[i18n] File watcher started for: %s", path)
	return fw
}

// watch 监听文件变化
func (w *fsnotifyWatcher) watch(reloadCallback func() error) {
	defer w.watcher.Close()

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// 检查是否为写入或创建事件
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				// 检查是否为 JSON 文件
				if filepath.Ext(event.Name) == ".json" {
					log.Printf("[i18n] Reloading locales due to file change: %s", filepath.Base(event.Name))

					// 执行重载回调
					if err := reloadCallback(); err != nil {
						log.Printf("[i18n] Failed to reload locales: %v", err)
					}
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("[i18n] Watcher error: %v", err)
		}
	}
}

// Close 关闭监听器
func (w *fsnotifyWatcher) Close() error {
	if w.closed {
		return nil
	}

	w.closed = true
	if w.watcher != nil {
		return w.watcher.Close()
	}
	return nil
}

// NoOpWatcher 空操作监听器
type NoOpWatcher struct{}

func (w *NoOpWatcher) Close() error {
	return nil
}

// WatcherConfig 监听器配置
type WatcherConfig struct {
	Enable bool   `yaml:"enable" json:"enable"`
	Path   string `yaml:"path" json:"path"`
}

// NewFileWatcherWithConfig 使用配置创建文件监听器
func NewFileWatcherWithConfig(config WatcherConfig, reloadCallback func() error) FileWatcher {
	if !config.Enable {
		return &NoOpWatcher{}
	}

	return NewFileWatcher(config.Path, reloadCallback)
}