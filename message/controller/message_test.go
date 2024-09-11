package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/xerrors"

	"github.com/opensourceways/message-manager/message/app"
)

// MockMessageListAppService 是 MessageListAppService 的模拟实现
type MockMessageListAppService struct {
	mock.Mock
}

func (m *MockMessageListAppService) GetInnerMessageQuick(userName string,
	cmd *app.CmdToGetInnerMessageQuick) ([]app.MessageListDTO, int64, error) {
	args := m.Called(userName, cmd)
	return args.Get(0).([]app.MessageListDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageListAppService) GetInnerMessage(userName string,
	cmd *app.CmdToGetInnerMessage) ([]app.MessageListDTO, int64, error) {
	args := m.Called(userName, cmd)
	return args.Get(0).([]app.MessageListDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageListAppService) CountAllUnReadMessage(userName string) ([]app.CountDTO, error) {
	args := m.Called(userName)
	return args.Get(0).([]app.CountDTO), args.Error(1)
}

func (m *MockMessageListAppService) SetMessageIsRead(cmd *app.CmdToSetIsRead) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func (m *MockMessageListAppService) RemoveMessage(cmd *app.CmdToSetIsRead) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func TestGetInnerMessageQuick_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	req, err := http.NewRequest("GET", "/message_center/inner_quick?source=test", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetInnerMessageQuick_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	req, err := http.NewRequest("GET", "/message_center/inner_quick", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetInnerMessageQuick_ConvertError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	mockService.On("GetInnerMessageQuick", "testUser", mock.Anything).
		Return(nil, int64(0), xerrors.New("error"))

	req, err := http.NewRequest("GET", "/message_center/inner_quick?source=test", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetInnerMessage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	req, err := http.NewRequest("GET", "/message_center/inner?source=test", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetInnerMessage_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	req, err := http.NewRequest("GET", "/message_center/inner", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetInnerMessage_ConvertError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	mockService.On("GetInnerMessage", "testUser", mock.Anything).
		Return(nil, int64(0), xerrors.New("error"))

	req, err := http.NewRequest("GET", "/message_center/inner?source=test", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCountAllUnReadMessage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	req, err := http.NewRequest("GET", "/message_center/inner/count", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCountAllUnReadMessage_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	mockService.On("CountAllUnReadMessage", "testUser").
		Return(0, xerrors.New("error"))

	req, err := http.NewRequest("GET", "/message_center/inner/count", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSetMessageIsRead_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	messages := []messageStatus{
		{Source: "source1", EventId: "event1"},
	}
	body, err := json.Marshal(messages)
	if err != nil {
		t.Fatal("Failed to marshal messages", err)
	}
	mockService.On("SetMessageIsRead", mock.Anything).Return(nil)

	req, err := http.NewRequest("PUT", "/message_center/inner", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockService.AssertExpectations(t)
}

func TestSetMessageIsRead_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	req, err := http.NewRequest("PUT", "/message_center/inner", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetMessageIsRead_ConvertError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	messages := []messageStatus{
		{Source: "source1", EventId: "event1"},
	}
	body, err := json.Marshal(messages)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("SetMessageIsRead", mock.Anything).Return(xerrors.New("error"))

	req, err := http.NewRequest("PUT", "/message_center/inner", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRemoveMessage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	messages := []messageStatus{
		{Source: "source1", EventId: "event1"},
	}
	body, err := json.Marshal(messages)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("RemoveMessage", mock.Anything).Return(nil)

	req, err := http.NewRequest("DELETE", "/message_center/inner", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockService.AssertExpectations(t)
}

func TestRemoveMessage_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	req, err := http.NewRequest("DELETE", "/message_center/inner", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRemoveMessage_ConvertError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	messages := []messageStatus{
		{Source: "source1", EventId: "event1"},
	}
	body, err := json.Marshal(messages)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("RemoveMessage", mock.Anything).Return(xerrors.New("error"))

	req, err := http.NewRequest("DELETE", "/message_center/inner", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
