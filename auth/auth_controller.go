package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Controller struct{}

func NewAuthController() *Controller {
	return &Controller{}
}

// Login authenticates a user and returns a JWT token
func (a *Controller) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if authenticate(request) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	tokenString, err := createToken(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	response := LoginResponse{
		Token: tokenString,
	}

	c.JSON(http.StatusOK, response)
}

// CreateToken generates a JWT token for the given username
func createToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      viper.GetString("jwt.expireTime"),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("jwt.secretKey")))
}

func authenticate(request LoginRequest) bool {
	// TODO
	// In a real application, the username and password would be validated
	// against a database of users. For simplicity, we will just use a hardcoded
	// username and password.
	const (
		validUsername = "admin"
		validPassword = "admin"
	)
	if request.Username != validUsername {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(validPassword), []byte(request.Password))
	return err == nil
}
