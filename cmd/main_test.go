//go:build integration
// +build integration

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/ngereci/xm_interview/auth"
	"github.com/ngereci/xm_interview/company"
	"github.com/ngereci/xm_interview/env"
	"github.com/ngereci/xm_interview/event"
	"github.com/ngereci/xm_interview/model"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Setup(t *testing.T) *httptest.Server {
	viper.SetConfigFile("../config/config_test.env")
	if err := viper.ReadInConfig(); err != nil {
		t.Error(err)
	}
	viper.AutomaticEnv()
	// Initialize the Cassandra session
	cluster := gocql.NewCluster(viper.GetString(env.COMPANY_CASSANDRA_HOST))
	cluster.Keyspace = viper.GetString(env.COMPANY_CASSANDRA_KEYSPACE)
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()

	if err != nil {
		t.Error(err)
	}

	kafkaProducer, err := event.NewKafkaAdapter([]string{viper.GetString(env.COMPANY_BROKER_URL)}, viper.GetString(env.COMPANY_BROKER_TOPIC))
	if err != nil {
		t.Error(err)
	}

	companyRepo := company.NewRepository(session)
	// empty test keyspace
	query := session.Query(`TRUNCATE companies_test.company;`)
	err = query.Exec()
	if err != nil {
		t.Error(err)
	}
	companyService := company.NewService(companyRepo, kafkaProducer)
	companyController := company.NewController(companyService)

	authController := auth.NewAuthController()
	authMiddleware := auth.NewAuthMiddleware(viper.GetString(env.COMPANY_JWT_SECRET_KEY))

	router := gin.Default()

	loginRouter := router.Group("/api/v1")
	loginRouter.POST("/login", authController.Login)

	apiRouter := router.Group("/api/v1")
	apiRouter.Use(authMiddleware.Authenticate())
	// Company routes
	apiRouter.POST("/companies", companyController.CreateCompany)
	apiRouter.PATCH("/companies/:id", companyController.UpdateCompany)
	apiRouter.DELETE("/companies/:id", companyController.DeleteCompany)
	apiRouter.GET("/companies/:id", companyController.GetCompany)

	return httptest.NewServer(router)
}

func Test_Login(t *testing.T) {
	server := Setup(t)
	testUser := auth.LoginRequest{
		Username: "admin",
		Password: "admin",
	}
	requestBody, _ := json.Marshal(testUser)
	// Execute the request to the service
	resp, err := http.Post(fmt.Sprintf("%s/api/v1/login", server.URL), "application/json", strings.NewReader(string(requestBody)))
	assert.NoError(t, err)
	defer resp.Body.Close()
	// Check the response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(body)
	token := `{"token":"`
	assert.Truef(t, strings.Contains(bodyString, token), "response body:%v does not contain %v", bodyString, token)
	defer server.Close()
}

func Test_Companies_CRUD(t *testing.T) {
	server := Setup(t)
	token, err := login(server)
	assert.NoError(t, err)
	var companyUUID uuid.UUID
	newCompany := &model.Company{
		Name:        "New Test Company",
		Description: "Description of a new test company",
		Employees:   100,
		Registered:  false,
		Type:        model.Corporation,
	}
	updatedCompany := &model.Company{
		Name:        "New Test Company UPDATED",
		Description: "Description of a new test company UPDATED",
		Employees:   200,
		Registered:  true,
		Type:        model.NonProfit,
	}
	t.Run("it should successfully insert", func(t *testing.T) {
		requestBody, _ := json.Marshal(newCompany)
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/companies", server.URL), strings.NewReader(string(requestBody)))
		assert.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var responseCompany model.Company
		err = json.NewDecoder(resp.Body).Decode(&responseCompany)
		assert.NoError(t, err)
		assert.Equal(t, newCompany.Name, responseCompany.Name)
		assert.Equal(t, newCompany.Type, responseCompany.Type)
		assert.Equal(t, newCompany.Employees, responseCompany.Employees)
		assert.Equal(t, newCompany.Description, responseCompany.Description)
		assert.Equal(t, newCompany.Registered, responseCompany.Registered)
		assert.NotEmpty(t, responseCompany.ID)
		companyUUID = responseCompany.ID
	})
	t.Run("inserted item should be available", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/companies/%v", server.URL, companyUUID), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		// Execute the request to the service
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		// Check the response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var responseCompany model.Company
		err = json.NewDecoder(resp.Body).Decode(&responseCompany)
		assert.NoError(t, err)
		assert.Equal(t, newCompany.Name, responseCompany.Name)
		assert.Equal(t, newCompany.Type, responseCompany.Type)
		assert.Equal(t, newCompany.Employees, responseCompany.Employees)
		assert.Equal(t, newCompany.Description, responseCompany.Description)
		assert.Equal(t, newCompany.Registered, responseCompany.Registered)
		assert.Equal(t, companyUUID, responseCompany.ID)
	})
	t.Run("it should successfully update", func(t *testing.T) {
		requestBody, _ := json.Marshal(updatedCompany)
		req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/companies/%v", server.URL, companyUUID), strings.NewReader(string(requestBody)))
		assert.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var responseCompany model.Company
		err = json.NewDecoder(resp.Body).Decode(&responseCompany)
		assert.NoError(t, err)
		assert.Equal(t, companyUUID, responseCompany.ID)
		assert.Equal(t, updatedCompany.Name, responseCompany.Name)
		assert.Equal(t, updatedCompany.Type, responseCompany.Type)
		assert.Equal(t, updatedCompany.Employees, responseCompany.Employees)
		assert.Equal(t, updatedCompany.Description, responseCompany.Description)
		assert.Equal(t, updatedCompany.Registered, responseCompany.Registered)
	})
	t.Run("updated item should be available", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/companies/%v", server.URL, companyUUID), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		// Execute the request to the service
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		// Check the response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var responseCompany model.Company
		err = json.NewDecoder(resp.Body).Decode(&responseCompany)
		assert.NoError(t, err)
		assert.Equal(t, updatedCompany.Name, responseCompany.Name)
		assert.Equal(t, updatedCompany.Type, responseCompany.Type)
		assert.Equal(t, updatedCompany.Employees, responseCompany.Employees)
		assert.Equal(t, updatedCompany.Description, responseCompany.Description)
		assert.Equal(t, updatedCompany.Registered, responseCompany.Registered)
		assert.Equal(t, companyUUID, responseCompany.ID)
	})
	t.Run("it should successfully delete", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/companies/%v", server.URL, companyUUID), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		// Execute the request to the service
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		// Check the response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("deleted item should NOT be available", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/companies/%v", server.URL, companyUUID), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		// Execute the request to the service
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		// Check the response
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		bodyString := string(body)
		responseExpected := `{"error":"company not found"}`
		assert.Equal(t, responseExpected, bodyString)
	})

	fmt.Print(companyUUID)
	defer server.Close()
}

func login(server *httptest.Server) (token string, err error) {
	testUser := auth.LoginRequest{
		Username: "admin",
		Password: "admin",
	}
	type tokenResponse struct {
		Token string `json:"token"`
	}
	requestBody, _ := json.Marshal(testUser)
	// Execute the request to the service
	resp, err := http.Post(fmt.Sprintf("%s/api/v1/login", server.URL), "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Check the response
	if err != nil {
		return "", err
	}
	var response tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}
	token = response.Token
	return
}
