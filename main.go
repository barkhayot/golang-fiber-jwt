package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/barkhayot/auth-fiber/database"
	"github.com/barkhayot/auth-fiber/routes"
)

func main() {

	database.Connect()

	app := fiber.New()
	
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":8000")
}