package plugins

import (
	"builder/pkg/downloader/downloaderPlugin"
	"io"
	"net/http"
	"os"
)

type HTTPDownloader struct {
}

const httpType = "http"

func (d *HTTPDownloader) Download(url string, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	outFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}

func (d *HTTPDownloader) GetType() downloaderPlugin.PluginType {
	return httpType
}

// 在 init 函数中注册 HTTP 下载器
func init() {
	downloaderPlugin.RegisterDownloader(httpType, &HTTPDownloader{})
}
