package main

import (
	"log"
	"net/http"
	"time"

	"github.com/HudYuSa/mydeen/internal/config"
	"github.com/HudYuSa/mydeen/internal/connection"
	"github.com/HudYuSa/mydeen/pkg/controllers"
	"github.com/HudYuSa/mydeen/pkg/middlewares"
	"github.com/HudYuSa/mydeen/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine
)

func init() {
	server = gin.Default()
	err := config.LoadConfig(".env")

	if err != nil {
		log.Fatal("? Could not load environment variables ", err)
	}

	// connect ke database
	connection.ConnectDB(&config.GlobalConfig)
}

func main() {
	// // middleware
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.GlobalConfig.ClientOrigin, "http://localhost:5173"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", " accept", "origin", "Cache-Control", " X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           600}))

	router := server.Group("/api")
	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "welcome to this project"})
	})

	// middleware
	router.Use(middlewares.DBTimeoutMiddleware(time.Duration(config.GlobalConfig.DatabaseTimeout) * time.Millisecond))

	// controller
	controllers.InitializeControllers()
	// routes
	routes.InitializeRoutes(router)

	// run app
	log.Fatal(server.Run(":" + config.GlobalConfig.ServerPort))
}
