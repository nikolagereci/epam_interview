package main

import (
	"github.com/gocql/gocql"
	"github.com/ngereci/xm_interview/auth"
	"github.com/ngereci/xm_interview/company"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Initialize the Cassandra session
	cluster := gocql.NewCluster(viper.GetString("cassandra.host"))
	cluster.Keyspace = viper.GetString("cassandra.keyspace")
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()

	if err != nil {
		log.Fatal("Failed to create Cassandra session: ", err)
	}

	//TODO kafkaProducer, err := kafka.NewProducer()
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	companyRepo := company.NewRepository(session)
	companyService := company.NewService(companyRepo /*, kafkaProducer*/)
	companyController := company.NewController(companyService)

	//TODO
	//eventRepo := event.NewRepository(db)
	//eventService := event.NewService(eventRepo, kafkaProducer)
	//eventController := event.NewController(eventService)

	authController := auth.NewAuthController()
	authMiddleware := auth.NewAuthMiddleware(viper.GetString("jwt.secretKey"))

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

	// Event routes
	//eventRouter := apiRouter.Group("/events")
	//eventRouter.Use(authMiddleware.Authenticate())
	//{
	//	eventRouter.POST("", eventController.Create)
	//	eventRouter.GET("/:id", eventController.Get)
	//}

	port := viper.GetString("server.port")
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
		//TODO
		//if err := kafkaProducer.Close(); err != nil {
		//	log.Fatalf("Error closing Kafka producer: %v", err)
		//}
	}()
}
