package routes

import (
	controllers "github/sahil/atlantic-backend/controllers/product_controller"

	"github.com/gofiber/fiber/v2"
)

func ProductRoute(app *fiber.App) {
	app.Post("/product", controllers.CreateProduct)
	app.Get("/product/:productID", controllers.GetAProduct)
	app.Get("/productphoto/:productID", controllers.GetAProductPhoto)
	app.Put("/product/:productID", controllers.UpdateProduct)
	app.Delete("/product/:productID", controllers.DeleteAProduct)
	app.Get("/allproducts", controllers.GetAllProducts)
}
