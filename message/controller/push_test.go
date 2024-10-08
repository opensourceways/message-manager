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

// MockMessagePushAppService 是 MessagePushAppService 的模拟实现
type MockMessagePushAppService struct {
	mock.Mock
}

func (m *MockMessagePushAppService) GetPushConfig(countPerPage, pageNum int,
	userName string, subsIds []string) ([]app.MessagePushDTO, error) {
	args := m.Called(countPerPage, pageNum, userName, subsIds)
	return args.Get(0).([]app.MessagePushDTO), args.Error(1)
}

func (m *MockMessagePushAppService) AddPushConfig(cmd *app.CmdToAddPushConfig) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func (m *MockMessagePushAppService) UpdatePushConfig(cmd *app.CmdToUpdatePushConfig) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func (m *MockMessagePushAppService) RemovePushConfig(cmd *app.CmdToDeletePushConfig) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func TestGetPushConfig_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	req, err := http.NewRequest("GET",
		"/message_center/config/push?count_per_page=10&page=1&subscribe_id=sub1,sub2", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetPushConfig_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	req, err := http.NewRequest("GET", "/message_center/config/push", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetPushConfig_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	userName := "testUser"
	subsIds := []string{"sub1", "sub2"}
	countPerPage, pageNum := 10, 1
	mockService.On("GetPushConfig", countPerPage, pageNum, userName, subsIds).
		Return(nil, xerrors.New("service error"))

	req, err := http.NewRequest("GET",
		"/message_center/config/push?count_per_page=10&page=1&subscribe_id=sub1,sub2", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddPushConfig_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	reqBody := app.CmdToAddPushConfig{
		SubscribeId:      1,
		RecipientId:      123456,
		NeedMessage:      true,
		NeedPhone:        false,
		NeedMail:         true,
		NeedInnerMessage: true,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("AddPushConfig", mock.Anything).Return(nil)

	req, err := http.NewRequest("POST", "/message_center/config/push", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockService.AssertExpectations(t)
}

func TestAddPushConfig_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	req, err := http.NewRequest("POST", "/message_center/config/push", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddPushConfig_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	reqBody := app.CmdToAddPushConfig{
		SubscribeId:      1,
		RecipientId:      123456,
		NeedMessage:      true,
		NeedPhone:        false,
		NeedMail:         true,
		NeedInnerMessage: true,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("AddPushConfig", mock.Anything).Return(xerrors.New("service error"))

	req, err := http.NewRequest("POST", "/message_center/config/push", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdatePushConfig_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	reqBody := app.CmdToUpdatePushConfig{
		SubscribeId:      []string{"1", "2"},
		RecipientId:      "recipient@example.com",
		NeedMessage:      true,
		NeedPhone:        false,
		NeedMail:         true,
		NeedInnerMessage: true,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("UpdatePushConfig", mock.Anything).Return(nil)

	req, err := http.NewRequest("PUT", "/message_center/config/push", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdatePushConfig_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	req, err := http.NewRequest("PUT", "/message_center/config/push", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePushConfig_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	reqBody := app.CmdToUpdatePushConfig{
		SubscribeId:      []string{"1", "2"},
		RecipientId:      "recipient@example.com",
		NeedMessage:      true,
		NeedPhone:        false,
		NeedMail:         true,
		NeedInnerMessage: true,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("UpdatePushConfig", mock.Anything).Return(xerrors.New("service error"))

	req, err := http.NewRequest("PUT", "/message_center/config/push", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRemovePushConfig_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	reqBody := app.CmdToDeletePushConfig{
		SubscribeId: 1,
		RecipientId: 123456,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("RemovePushConfig", mock.Anything).Return(nil)

	req, err := http.NewRequest("DELETE", "/message_center/config/push", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockService.AssertExpectations(t)
}

func TestRemovePushConfig_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	req, err := http.NewRequest("DELETE", "/message_center/config/push", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRemovePushConfig_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessagePushAppService)
	AddRouterForMessagePushController(r, mockService)

	reqBody := app.CmdToDeletePushConfig{
		SubscribeId: 1,
		RecipientId: 123456,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("RemovePushConfig", mock.Anything).Return(xerrors.New("service error"))

	req, err := http.NewRequest("DELETE", "/message_center/config/push", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
