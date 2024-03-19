package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/HudYuSa/mydeen/internal/config"
	"github.com/HudYuSa/mydeen/internal/connection"
	"github.com/HudYuSa/mydeen/pkg/controllers"
	"github.com/HudYuSa/mydeen/pkg/middlewares"
	"github.com/HudYuSa/mydeen/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
	err := config.LoadConfig(".env")

	if err != nil {
		log.Fatal("? Could not load environment variables ", err)
	}

	// connect ke database
	connection.ConnectDB(&config.GlobalConfig)
}

func main() {
	// // middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.GlobalConfig.ClientOrigin, "http://localhost:5173", "http://192.168.1.15:5173"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", " Authorization", " accept", "origin", "Cache-Control", " X-Requested-With", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour}))

	api := router.Group("/api")
	api.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "welcome to this project"})
	})

	api.Handle(http.MethodGet, "/data", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, []byte("this is the response data"))
	})

	// static
	// router.Static("/assets", "./frontend/dist/assets")
	// router.StaticFS("/static", http.Dir("./frontend/dist"))
	// router.Use(static.Serve("/", static.LocalFile("./client/build", true)))
	router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	// middleware
	api.Use(middlewares.DBTimeoutMiddleware(time.Duration(config.GlobalConfig.DatabaseTimeout) * time.Millisecond))
	api.Use(middlewares.AuthenticateUser())

	// melody
	m := melody.New()
	// controller
	controllers.InitializeControllers(m)
	// routes
	routes.InitializeRoutes(api)

	// melody handlers
	m.HandleConnect(controllers.WebSocket.HandleConnect)
	m.HandleDisconnect(controllers.WebSocket.HandleDisconnect)
	m.HandleMessage(controllers.WebSocket.HandleMessage)

	// router.GET("/ws/:sessionId", func(ctx *gin.Context) {
	// 	m.HandleRequest(ctx.Writer, ctx.Request)
	// })

	// m.HandleConnect(func(s *melody.Session) {
	// 	log.Println("connections: ", m.Len()+1)
	// 	log.Println("new connection")
	// })

	// m.HandleDisconnect(func(s *melody.Session) {
	// 	log.Println("removed connection")
	// })

	// m.HandleMessage(func(s *melody.Session, b []byte) {
	// 	m.BroadcastFilter(b, func(q *melody.Session) bool {
	// 		log.Println(q.Request.URL.Path)
	// 		log.Println(s.Request.URL.Path)
	// 		return q.Request.URL.Path == s.Request.URL.Path
	// 	})
	// 	s.Write(b)
	// })

	// Handle all other routes by serving index.html
	router.NoRoute(func(ctx *gin.Context) {
		ctx.File("./dist/index.html")
	})

	// // run app
	fmt.Println(config.GlobalConfig.ServerPort)
	log.Fatal(router.Run(config.GlobalConfig.ServerPort))
}
