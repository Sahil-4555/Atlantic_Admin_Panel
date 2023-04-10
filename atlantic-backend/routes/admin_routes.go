package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github/sahil/atlantic-backend/controllers/admin_controller"
)

func AdminRoutes(app *fiber.App) {
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Get("/user", controllers.User)
	app.Post("/logout", controllers.Logout)
}
