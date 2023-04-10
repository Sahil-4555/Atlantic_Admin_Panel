package routes

import (
	controllers "github/sahil/atlantic-backend/controllers/app_controller"

	"github.com/gofiber/fiber/v2"
)

func AppUserRoute(app *fiber.App) {
	app.Post("/users", controllers.CreateUser)
	app.Get("/users/:uid", controllers.GetAUser)
	app.Put("/users/:uid", controllers.EditAUser)
	app.Delete("/users/:uid", controllers.DeleteAUser)
	app.Get("/allusers", controllers.GetAllUsers)
}
