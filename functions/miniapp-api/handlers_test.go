package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	telegramparser "github.com/kd3n1z/go-telegram-parser"
	"github.com/stretchr/testify/assert"
)

type mockEnvProvider struct {
	token string
}

func (m *mockEnvProvider) GetBotToken() string {
	return m.token
}

var testEnvProvider = &mockEnvProvider{token: "test"}

type mockTelegramParser struct {
	shouldSucceed bool
	userID        int64
}

func (m *mockTelegramParser) Parse(_ string) (telegramparser.WebAppInitData, error) {
	if !m.shouldSucceed {
		return telegramparser.WebAppInitData{}, fmt.Errorf("mock validation failed")
	}
	return telegramparser.WebAppInitData{
		User: telegramparser.WebAppUser{Id: m.userID},
	}, nil
}

var failMockParser ParserFactory = func(botToken string) ParserInterface {
	var t = mockTelegramParser{shouldSucceed: false, userID: 123456789}
	return &t
}

var successMockParser ParserFactory = func(botToken string) ParserInterface {
	var t = mockTelegramParser{shouldSucceed: true, userID: 123456789}
	return &t
}

func TestGetUserID_NoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create request without Authorization header
	req, _ := http.NewRequest("GET", "/test", nil)
	c.Request = req

	userID := getUserID(c, testEnvProvider, failMockParser)

	assert.Nil(t, userID)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserID_EmptyAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create request with empty Authorization header
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "")
	c.Request = req

	userID := getUserID(c, testEnvProvider, failMockParser)

	assert.Nil(t, userID)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserID_NoBotToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer some_init_data")
	c.Request = req

	userID := getUserID(c, testEnvProvider, failMockParser)

	assert.Nil(t, userID)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserID_InvalidTelegramData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid_telegram_data")
	c.Request = req

	userID := getUserID(c, testEnvProvider, failMockParser)

	assert.Nil(t, userID)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserID_ValidTelegramData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// This is a sample valid Telegram WebApp init data format
	// In real tests, you'd need actual valid signed data or mock the parser
	validInitData := "user=%7B%22id%22%3A123456789%2C%22first_name%22%3A%22Test%22%7D&hash=valid_hash"

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validInitData)
	c.Request = req
	userID := getUserID(c, testEnvProvider, successMockParser)
	assert.NotNil(t, userID)
}

func TestGetTag_ID_NoParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// No tagId in request
	req, _ := http.NewRequest("GET", "/test", nil)
	c.Request = req

	tagID := getTagID(c)

	assert.Nil(t, tagID)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTag_ID_Param(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	randomTag := rand.Int63()

	c.AddParam("tagId", fmt.Sprintf("%d", randomTag))
	req, _ := http.NewRequest("GET", "/test", nil)
	c.Request = req

	tagID := getTagID(c)

	assert.NotNil(t, tagID)
	assert.Equal(t, randomTag, *tagID)
	assert.Equal(t, http.StatusOK, w.Code)
}
