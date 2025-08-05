package main

import (
	"flag"
	"github.com/flyhigher139/file-transfer/handler"
	"github.com/flyhigher139/file-transfer/route"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	// 从命令行参数获取文件存储目录
	storagePath := flag.String("storage", "./uploads", "the path to the file storage directory")
	flag.Parse()

	// 确保存储目录存在
	if err := os.MkdirAll(*storagePath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	// 将存储路径传递给 handler
	handler.StoragePath = *storagePath

	router := gin.Default()
	// 静态文件服务
	router.StaticFS("/static", http.Dir("./static"))
	// 注册路由
	route.RegisterRoutes(router)
	// 启动服务
	router.Run(":8080")
}