package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barkhayot/auth-fiber/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Get("/user", controllers.User)
	app.Post("logout", controllers.Logout)
}
