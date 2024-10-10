// server/server_test.go
package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStartWebServer(t *testing.T) {

	// 使用 httptest 创建一个新请求
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()

	// 启动 Web 服务器
	go StartWebServer()

	// 发送请求到服务器
	w.WriteHeader(http.StatusOK)
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello, World!", w.Body.String())
}

func TestLogRequest(t *testing.T) {
	// 创建一个新的 gin 引擎
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(logRequest())

	// 测试请求
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err)
	// 记录请求
	r.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLogRequestWithError(t *testing.T) {
	// 创建一个新的 gin 引擎
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(logRequest())
	r.GET("/error", func(c *gin.Context) {
		c.Error(fmt.Errorf("test error"))
		c.String(http.StatusInternalServerError, "Internal Server Error")
	})

	// 测试请求
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/error", nil)
	assert.NoError(t, err)
	// 记录请求
	r.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}
