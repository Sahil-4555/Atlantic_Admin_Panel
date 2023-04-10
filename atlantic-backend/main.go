package main

import (
	"github/sahil/atlantic-backend/configs"
	"github/sahil/atlantic-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	configs.ConnectDB()

	print("Sahil_4555! Connecting You To The Database")

	routes.AdminRoutes(app)
	routes.AppUserRoute(app)
	routes.ProductRoute(app)
	app.Listen(":5000")
}
