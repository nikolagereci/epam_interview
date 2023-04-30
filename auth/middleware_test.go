package auth

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_Authenticate_Success(t *testing.T) {
	secretKey := "secret"
	middleware := NewAuthMiddleware(secretKey)

	// Create a test JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = "123"
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatal(err)
	}

	// Create a test request with the JWT token in the Authorization header
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a test response recorder and a Gin context
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Request = req

	// Invoke the Authenticate middleware
	middleware.Authenticate()(c)

	// Verify that the middleware set the userId in the Gin context
	userId, ok := c.Get("userId")
	if !ok || userId.(string) != "123" {
		t.Errorf("Authenticate middleware did not set userId in context")
	}

	// Verify that the middleware called the next handler
	if res.Code != http.StatusOK {
		t.Errorf("Authenticate middleware did not call next handler")
	}
}

func TestAuthMiddleware_Authenticate_MissingAuthorizationHeader(t *testing.T) {
	secretKey := "secret"
	middleware := NewAuthMiddleware(secretKey)

	// Create a test request without an Authorization header
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test response recorder and a Gin context
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Request = req

	// Invoke the Authenticate middleware
	middleware.Authenticate()(c)

	// Verify that the middleware returned a 401 Unauthorized error
	if res.Code != http.StatusUnauthorized {
		t.Errorf("Authenticate middleware did not return 401 Unauthorized")
	}

	// Verify that the middleware did not set the userId in the Gin context
	_, ok := c.Get("userId")
	if ok {
		t.Errorf("Authenticate middleware set userId in context despite missing Authorization header")
	}
}

func TestAuthMiddleware_Authenticate_InvalidTokenFormat(t *testing.T) {
	secretKey := "secret"
	middleware := NewAuthMiddleware(secretKey)

	// Create a test request with an invalid Authorization header format
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "invalid_token_format")

	// Create a test response recorder and a Gin context
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Request = req

	// Invoke the Authenticate middleware
	middleware.Authenticate()(c)

	// Verify that the middleware returned a 401 Unauthorized error
	if res.Code != http.StatusUnauthorized {
		t.Errorf("Authenticate middleware did not return 401 Unauthorized")
	}

	// Verify that the middleware did not set the userId in the Gin context
	_, ok := c.Get("userId")
	if ok {
		t.Errorf("Authenticate middleware set userId in context despite invalid token format")
	}
}

func TestAuthMiddleware_Authenticate_InvalidToken(t *testing.T) {
	secretKey := "secret"
	middleware := NewAuthMiddleware(secretKey)

	// Create a mock Gin context
	router := gin.New()
	router.Use(middleware.Authenticate())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	req, _ := http.NewRequest("GET", "/", nil)

	// Test with invalid token
	req.Header.Set("Authorization", "Bearer invalidtoken")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check response status code
	assert.Equal(t, http.StatusUnauthorized, resp.Code)

	// Check response body
	//expectedBody := gin.H{"error": "Failed to parse token: token contains an invalid number of segments"}
	expectedBody := `{"error":"token contains an invalid number of segments"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}
