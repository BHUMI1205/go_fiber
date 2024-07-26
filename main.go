package main

import (
	"database1/database"
	_ "database1/docs"
	"database1/routes"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/limiter"
	"flag"
	"github.com/gofiber/contrib/websocket"
	"github.com/joho/godotenv"
	"log"   
	"os"
	"time"
)

// @title Fiber Swagger Example API
// @version 2.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
// @schemes http

func main() {
	// clear the env file
	// os.Clearenv()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.Connectdb()
	app := fiber.New()
	// app.Static("/", "./public")

	port := os.Getenv("PORT")
	flag.StringVar(&port, "port", ":8080", "port number")
	flag.Parse()
	log.Println("Welcome to logging in Golang!!")

	// 	app.Use(limiter.New(limiter.Config{
	// 		Expiration: 1000 * time.Second,
	// 		Max:        10,
	//   }))

	app.Use(func(c *fiber.Ctx) error {
		log.Println(c.Method(), c.Path(), time.Now().Local())
		return c.Next()
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		log.Println(c.Locals("allowed"))
		log.Println(c.Params("id"))
		log.Println(c.Query("v"))
		log.Println(c.Cookies("session"))
		return fiber.ErrUpgradeRequired
	})
	
	routes.SetupUserRoutes(app)
	log.Fatal(app.Listen(port))
}
