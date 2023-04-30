package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ngereci/xm_interview/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
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

	if !a.authenticate(request) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	tokenString, err := a.createToken(request.Username)
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
func (a *Controller) createToken(username string) (string, error) {
	expiration := time.Duration(rand.Int31n(viper.GetInt32(env.COMPANY_JWT_EXPIRE_TIME))) * time.Second
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString(env.COMPANY_JWT_SECRET_KEY)))
}

func (a *Controller) authenticate(request LoginRequest) bool {
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
	validPasswordHash, err := bcrypt.GenerateFromPassword([]byte(validPassword), 0)
	if err != nil {
		log.Errorf("unable to generate password hash, error: %v", err)
		return false
	}
	err = bcrypt.CompareHashAndPassword(validPasswordHash, []byte(request.Password))
	fmt.Print(err)
	return err == nil
}
