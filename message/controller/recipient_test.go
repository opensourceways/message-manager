package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/xerrors"

	"github.com/opensourceways/message-manager/message/app"
)

// Mock for the MessageRecipientAppService
type MockMessageRecipientAppService struct {
	mock.Mock
}

func (m *MockMessageRecipientAppService) GetRecipientConfig(countPerPage,
	pageNum int, userName string) ([]app.MessageRecipientDTO, int64, error) {
	args := m.Called(countPerPage, pageNum, userName)
	return args.Get(0).([]app.MessageRecipientDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageRecipientAppService) AddRecipientConfig(userName string,
	cmd *app.CmdToAddRecipient) error {
	return m.Called(userName, cmd).Error(0)
}

func (m *MockMessageRecipientAppService) UpdateRecipientConfig(userName string,
	cmd *app.CmdToUpdateRecipient) error {
	return m.Called(userName, cmd).Error(0)
}

func (m *MockMessageRecipientAppService) RemoveRecipientConfig(userName string,
	cmd *app.CmdToDeleteRecipient) error {
	return m.Called(userName, cmd).Error(0)
}

func (m *MockMessageRecipientAppService) SyncUserInfo(cmd *app.CmdToSyncUserInfo) (uint, error) {
	args := m.Called(cmd)
	return args.Get(0).(uint), args.Error(1)
}

func TestGetRecipientConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockAppService := new(MockMessageRecipientAppService)

	AddRouterForMessageRecipientController(router, mockAppService)

	// Successful case
	mockAppService.On("GetRecipientConfig", 10, 1, "testUser").
		Return([]app.MessageRecipientDTO{{}}, int64(1), nil)

	req, err := http.NewRequest(http.MethodGet,
		"/message_center/config/recipient?count_per_page=10&page=1", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer testToken")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Error case
	mockAppService.On("GetRecipientConfig", 10, 1, "testUser").
		Return(nil, int64(0), xerrors.New("db error"))

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAddRecipientConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockAppService := new(MockMessageRecipientAppService)
	AddRouterForMessageRecipientController(router, mockAppService)
	// Successful case
	mockAppService.On("AddRecipientConfig", "testUser", mock.Anything).
		Return(nil)

	reqBody := `{"name":"recipient1"}`
	req, err := http.NewRequest(http.MethodPost, "/message_center/config/recipient",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer testToken")
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Error case
	mockAppService.On("AddRecipientConfig", "testUser", mock.Anything).
		Return(xerrors.New("add error"))

	reqBody = `{"name":"recipient1"}`
	req, err = http.NewRequest(http.MethodPost, "/message_center/config/recipient",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestUpdateRecipientConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockAppService := new(MockMessageRecipientAppService)
	AddRouterForMessageRecipientController(router, mockAppService)

	// Successful case
	mockAppService.On("UpdateRecipientConfig", "testUser", mock.Anything).
		Return(nil)

	reqBody := `{"id":"1","name":"updatedRecipient"}`
	req, err := http.NewRequest(http.MethodPut, "/message_center/config/recipient",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer testToken")
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Error case
	mockAppService.On("UpdateRecipientConfig", "testUser", mock.Anything).
		Return(xerrors.New("update error"))

	reqBody = `{"id":"1","name":"updatedRecipient"}`
	req, err = http.NewRequest(http.MethodPut, "/message_center/config/recipient",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestRemoveRecipientConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockAppService := new(MockMessageRecipientAppService)
	AddRouterForMessageRecipientController(router, mockAppService)

	// Successful case
	mockAppService.On("RemoveRecipientConfig", "testUser", mock.Anything).
		Return(nil)

	reqBody := `{"id":"1"}`
	req, err := http.NewRequest(http.MethodDelete, "/message_center/config/recipient",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer testToken")
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Error case
	mockAppService.On("RemoveRecipientConfig", "testUser", mock.Anything).
		Return(xerrors.New("delete error"))

	reqBody = `{"id":"1"}`
	req, err = http.NewRequest(http.MethodDelete, "/message_center/config/recipient",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestSyncUserInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockAppService := new(MockMessageRecipientAppService)
	AddRouterForMessageRecipientController(router, mockAppService)

	// Successful case
	mockAppService.On("SyncUserInfo", mock.Anything).Return(uint(1), nil)

	reqBody := `{"user_info":"example"}`
	req, err := http.NewRequest(http.MethodPost, "/message_center/config/recipient/sync",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer testToken")
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusAccepted, recorder.Code)

	// Error case
	mockAppService.On("SyncUserInfo", mock.Anything).Return(uint(0), xerrors.New("sync error"))

	reqBody = `{"user_info":"example"}`
	req, err = http.NewRequest(http.MethodPost, "/message_center/config/recipient/sync",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusAccepted, recorder.Code)
}
