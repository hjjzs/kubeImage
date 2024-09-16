package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"os"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// MinIO 服务的连接信息
	endpoint := "minio-service.default.svc.cluster.local:9000"    // MinIO 服务地址
	accessKeyID := "iMyG9cGsBSOqBN5NEGr4"                         // 替换为你的 Access Key
	secretAccessKey := "CvJ0CwdyOPCBb8g3eSnStxRgfLT0OilLaFFAOqkr" // 替换为你的 Secret Key
	useSSL := false                                               // 是否使用 SSL

	// 初始化 MinIO 客户端
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// 要下载的文件信息
	bucketName := "builder"           // 替换为你的存储桶名称
	objectName := "nginx/nginx.tar"   // 替换为你要下载的对象名称
	filePath := "downloaded-file.tar" // 下载文件保存的路径

	// 创建文件以保存下载内容
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// 下载文件
	err = minioClient.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Successfully downloaded", objectName, "to", filePath)
}
