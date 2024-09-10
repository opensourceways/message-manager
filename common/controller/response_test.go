package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockError struct {
	code string
	msg  string
}

func (e mockError) Error() string {
	return e.msg
}

func (e mockError) ErrorCode() string {
	return e.code
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSendBadRequestBody(t *testing.T) {
	router := setupRouter()
	router.POST("/test", func(ctx *gin.Context) {
		SendBadRequestBody(ctx, mockError{"bad_request", "Invalid body"})
	})

	w := performRequest(router, "POST", "/test", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"code":"bad_request","msg":"Invalid body","data":null}`, w.Body.String())
}

func TestSendUnauthorized(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(ctx *gin.Context) {
		SendUnauthorized(ctx, mockError{"unauthorized", "access_denied"})
	})

	w := performRequest(router, "GET", "/test", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"code":"unauthorized","msg":"access_denied","data":null}`, w.Body.String())
}

func TestSendBadRequestParam(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(ctx *gin.Context) {
		SendBadRequestParam(ctx, mockError{"bad_param", "Invalid parameter"})
	})

	w := performRequest(router, "GET", "/test", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"code":"bad_param","msg":"Invalid parameter","data":null}`, w.Body.String())
}

func TestSendRespOfPut(t *testing.T) {
	router := setupRouter()
	router.PUT("/test", func(ctx *gin.Context) {
		SendRespOfPut(ctx, nil)
	})

	w := performRequest(router, "PUT", "/test", nil)
	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.JSONEq(t, `{"code":"","msg":"success","data":null}`, w.Body.String())
}

func TestSendRespOfGet(t *testing.T) {
	router := setupRouter()
	router.GET("/test", func(ctx *gin.Context) {
		SendRespOfGet(ctx, "data")
	})

	w := performRequest(router, "GET", "/test", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"code":"","msg":"","data":"data"}`, w.Body.String())
}

func TestSendRespOfPost(t *testing.T) {
	router := setupRouter()
	router.POST("/test", func(ctx *gin.Context) {
		SendRespOfPost(ctx, "new data")
	})

	w := performRequest(router, "POST", "/test", nil)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"code":"","msg":"","data":"new data"}`, w.Body.String())
}

func TestSendRespOfDelete(t *testing.T) {
	router := setupRouter()
	router.DELETE("/test", func(ctx *gin.Context) {
		SendRespOfDelete(ctx)
	})

	w := performRequest(router, "DELETE", "/test", nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}
