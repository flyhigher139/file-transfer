package route

import (
	"github.com/flyhigher139/file-transfer/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/files", handler.ListFiles)
		api.POST("/files", handler.UploadFile)
		api.POST("/files/merge", handler.MergeFile)
		api.GET("/files/:filename", handler.DownloadFile)
		api.POST("/simple/upload", handler.SimpleUploadFile)
		api.GET("/simple/download/:filename", handler.SimpleDownloadFile)
	}
}