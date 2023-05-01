package auth

import (
	"encoding/json"
	"github.com/ngereci/xm_interview/env"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestController_Login(t *testing.T) {
	// Initialize the Gin engine
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Create a test user
	testUser := LoginRequest{
		Username: "admin",
		Password: "admin",
	}

	// Create a new auth controller
	authController := NewAuthController()

	//set up test env
	viper.Set(env.COMPANY_JWT_SECRET_KEY, "test-key")
	viper.Set(env.COMPANY_JWT_EXPIRE_TIME, 3600)
	// Mount the auth controller's routes on the router
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
	}

	// Create a test request
	requestBody, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that the response has the expected status code and body
	assert.Equal(t, http.StatusOK, w.Code)
	//expectedResponseBody := `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI4NzkwMzQsInVzZXJuYW1lIjoiYWRtaW4ifQ.rhcSdAGchxPoXiOUbyb5ZfZ7l-Wp7PckO4m1zUfDhWI"}`
	//assert.Equal(t, expectedResponseBody, w.Body.String())
	expectedResponseBodyPart := `"token":`
	bodyString := w.Body.String()
	assert.True(t, strings.Contains(bodyString, expectedResponseBodyPart))

}

func TestController_Login_InvalidRequest(t *testing.T) {
	// Initialize the Gin engine
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Create a new auth controller
	authController := NewAuthController()

	// Mount the auth controller's routes on the router
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
	}

	// Create a test request with an invalid request body
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader("invalidrequest"))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that the response has the expected status code and body
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expectedResponseBody := `{"error":"invalid character 'i' looking for beginning of value"}`
	assert.Equal(t, expectedResponseBody, w.Body.String())
}

func TestController_Login_InvalidCredentials(t *testing.T) {
	// Initialize the Gin engine
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Create a test user with invalid credentials
	testUser := LoginRequest{
		Username: "wronguser",
		Password: "wrongpassword",
	}
	// Create a new auth controller
	authController := NewAuthController()

	// Mount the auth controller's routes on the router
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
	}

	// Create a test request with an invalid request body
	requestBody, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that the response has the expected status code and body
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	expectedResponseBody := `{"error":"Invalid username or password"}`
	assert.Equal(t, expectedResponseBody, w.Body.String())
}
