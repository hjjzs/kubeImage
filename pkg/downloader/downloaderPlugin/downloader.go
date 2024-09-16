package downloaderPlugin

import (
	"fmt"
	"sync"
)

type PluginType string

// Downloader 接口
type Downloader interface {
	Download(url string, destination string) error
	GetType() PluginType
}

// 下载器注册表
var (
	downloaders = make(map[string]Downloader)
	mu          sync.Mutex
)

// RegisterDownloader 注册下载器
func RegisterDownloader(protocol string, downloader Downloader) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := downloaders[protocol]; exists {
		panic(fmt.Sprintf("Downloader for protocol %s already exists", protocol))
	}
	downloaders[protocol] = downloader
}

// GetDownloader 根据 URL 获取下载器
func GetDownloader(url string) (Downloader, error) {
	mu.Lock()
	defer mu.Unlock()

	for protocol, downloader := range downloaders {
		if startsWithProtocol(url, protocol) {
			return downloader, nil
		}
	}

	return nil, fmt.Errorf("unsupported protocol")
}

// GetDownloaderByType 根据type 获取下载器
func GetDownloaderByType(t string) (Downloader, error) {
	mu.Lock()
	defer mu.Unlock()

	for protocol, downloader := range downloaders {
		if protocol == t {
			return downloader, nil
		}
	}
	return nil, fmt.Errorf("unsupported protocol")
}

// 辅助函数：判断 URL 前缀是否是协议
func startsWithProtocol(url string, protocol string) bool {
	return len(url) > len(protocol) && url[:len(protocol)+3] == protocol+"://"
}
