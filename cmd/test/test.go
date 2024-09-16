package main

import (
	"fmt"
	_ "test/pkg/downloader"
	"test/pkg/downloader/downloaderPlugin"
)

func main() {
	downloader, err := downloaderPlugin.GetDownloader("http://www.baidu.com")
	if err != nil {
		panic(err)
	}

	byType, err := downloaderPlugin.GetDownloaderByType("http")
	if err != nil {
		panic(err)
	}
	fmt.Printf("byType: %v\n", byType)
	//downloader.Download()
	fmt.Println(downloader)
}
