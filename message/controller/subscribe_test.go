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

// Mock for the MessageSubscribeAppService
type MockMessageSubscribeAppService struct {
	mock.Mock
}

func (m *MockMessageSubscribeAppService) UpdateSubsConfig(userName string, cmd *app.CmdToUpdateSubscribe) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageSubscribeAppService) GetAllSubsConfig(userName string) (
	[]app.MessageSubscribeDTO, error) {
	args := m.Called(userName)
	return args.Get(0).([]app.MessageSubscribeDTO), args.Error(1)
}

func (m *MockMessageSubscribeAppService) GetSubsConfig(userName string) (
	[]app.MessageSubscribeDTOWithPushConfig, int64, error) {
	args := m.Called(userName)
	return args.Get(0).([]app.MessageSubscribeDTOWithPushConfig), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageSubscribeAppService) AddSubsConfig(userName string,
	cmd *app.CmdToAddSubscribe) ([]uint, error) {
	args := m.Called(userName, cmd)
	return args.Get(0).([]uint), args.Error(1)
}

func (m *MockMessageSubscribeAppService) RemoveSubsConfig(userName string,
	cmd *app.CmdToDeleteSubscribe) error {
	return m.Called(userName, cmd).Error(0)
}

func TestGetAllSubsConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockAppService := new(MockMessageSubscribeAppService)
	AddRouterForMessageSubscribeController(router, mockAppService)

	// Successful case
	mockAppService.On("GetAllSubsConfig", "testUser").
		Return([]app.MessageSubscribeDTO{{}}, nil)

	req, err := http.NewRequest(http.MethodGet, "/message_center/config/subs/all", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer testToken")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Error case
	mockAppService.On("GetAllSubsConfig", "testUser").
		Return(nil, xerrors.New("db error"))

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAddSubsConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockAppService := new(MockMessageSubscribeAppService)
	AddRouterForMessageSubscribeController(router, mockAppService)

	// Successful case
	mockAppService.On("AddSubsConfig", "testUser", mock.Anything).
		Return([]uint{1}, nil)

	reqBody := `{"name":"new subscription"}`
	req, err := http.NewRequest(http.MethodPost, "/message_center/config/subs",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer testToken")
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Error case (binding error)
	reqBody = `{"invalid_field":"new subscription"}`
	req, err = http.NewRequest(http.MethodPost, "/message_center/config/subs",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Error case (add error)
	mockAppService.On("AddSubsConfig", "testUser", mock.Anything).
		Return(nil, xerrors.New("add error"))

	reqBody = `{"name":"new subscription"}`
	req, err = http.NewRequest(http.MethodPost, "/message_center/config/subs",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestRemoveSubsConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockAppService := new(MockMessageSubscribeAppService)
	AddRouterForMessageSubscribeController(router, mockAppService)

	// Successful case
	mockAppService.On("RemoveSubsConfig", "testUser", mock.Anything).
		Return(nil)

	reqBody := `{"id":1}`
	req, err := http.NewRequest(http.MethodDelete, "/message_center/config/subs",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer testToken")
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Error case (binding error)
	reqBody = `{"invalid_field":1}`
	req, err = http.NewRequest(http.MethodDelete, "/message_center/config/subs",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// Error case (remove error)
	mockAppService.On("RemoveSubsConfig", "testUser", mock.Anything).
		Return(xerrors.New("remove error"))

	reqBody = `{"id":1}`
	req, err = http.NewRequest(http.MethodDelete, "/message_center/config/subs",
		strings.NewReader(reqBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}
