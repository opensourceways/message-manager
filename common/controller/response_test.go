package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 测试 SendBadRequestBody
func TestSendBadRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试常规错误
	err := errors.New("something went wrong")
	SendBadRequestBody(c, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"code":"bad_request_body","msg":"something went wrong","data":null}`, w.Body.String())
}

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

// 测试 SendRespOfPut
func TestSendRespOfPut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试 nil 数据响应
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendRespOfPut(c, nil)
	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.JSONEq(t, `{"code":"","msg":"success","data":null}`, w.Body.String())

	// 测试有效数据响应
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	SendRespOfPut(c, "data")
	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.JSONEq(t, `{"code":"","msg":"","data":"data"}`, w.Body.String())
}

// 测试 SendRespOfGet
func TestSendRespOfGet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendRespOfGet(c, "data")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"code":"","msg":"","data":"data"}`, w.Body.String())
}

// 测试 SendRespOfPost
func TestSendRespOfPost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendRespOfPost(c, nil)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"code":"","msg":"success","data":null}`, w.Body.String())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	SendRespOfPost(c, "data")
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"code":"","msg":"","data":"data"}`, w.Body.String())
}

// 测试 SendRespOfDelete
func TestSendRespOfDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendRespOfDelete(c)
	assert.Equal(t, http.StatusNoContent, w.Code)
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
