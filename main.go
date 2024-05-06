package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 创建一个默认的Gin引擎
	router := gin.Default()

	// 定义路由和处理函数
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Message Manager!",
		})
	})

	// 启动HTTP服务器，监听在本地的 8080 端口
	router.Run(":8080")
}
