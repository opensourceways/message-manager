package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 测试 SendUnauthorized
func TestSendUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试常规错误
	err := errors.New("unauthorized access")
	SendUnauthorized(c, err)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"code":"unauthorized","msg":"unauthorized access","data":null}`, w.Body.String())
}

// 测试 SendBadRequestParam
func TestSendBadRequestParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试常规错误
	err := errors.New("parameter invalid")
	SendBadRequestParam(c, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"code":"bad_request_param","msg":"parameter invalid","data":null}`, w.Body.String())
}

// 测试 SendError
func TestSendError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试常规错误
	err := errors.New("parameter invalid")
	SendError(c, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"code":"system_error","msg":"parameter invalid","data":null}`, w.Body.String())
}
