package main

import (
	"github.com/gocql/gocql"
	"github.com/ngereci/xm_interview/auth"
	"github.com/ngereci/xm_interview/company"
	"github.com/ngereci/xm_interview/env"
	"github.com/ngereci/xm_interview/event"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config/config.env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	viper.AutomaticEnv()
	// Initialize the Cassandra session
	cluster := gocql.NewCluster(viper.GetString(env.COMPANY_CASSANDRA_HOST))
	cluster.Keyspace = viper.GetString(env.COMPANY_CASSANDRA_KEYSPACE)
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()

	if err != nil {
		log.Fatal("Failed to create Cassandra session: ", err)
	}

	kafkaProducer, err := event.NewKafkaAdapter([]string{viper.GetString(env.COMPANY_BROKER_URL)}, viper.GetString(env.COMPANY_BROKER_TOPIC))
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	companyRepo := company.NewRepository(session)
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

	port := viper.GetString(env.COMPANY_SERVER_PORT)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Printf("Server listening on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}

	defer func() {
		session.Close()
		if err := kafkaProducer.Close(); err != nil {
			log.Fatalf("Error closing Kafka producer: %v", err)
		}
	}()
}
