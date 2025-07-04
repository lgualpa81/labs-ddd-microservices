package routes

import (
	"poc-product-svc/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	productHandler *handlers.ProductHandler
}

func NewRouter(productHandler *handlers.ProductHandler) *Router {
	return &Router{
		productHandler: productHandler,
	}
}

func (r *Router) Setup(app *fiber.App) {
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":  "ok",
			"service": "product_svc",
		})
	})

	api := app.Group("/api/v1")
	products := api.Group("/products")

	products.Post("/", r.productHandler.CreateProduct)
	products.Get("/", r.productHandler.GetAllProducts)
	products.Get("/:id", r.productHandler.GetProduct)
	products.Put("/:id", r.productHandler.UpdateProduct)
	products.Delete("/:id", r.productHandler.DeleteProduct)
	products.Patch("/:id/stock", r.productHandler.UpdateStock)

	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Route not found",
			"path":    ctx.Path(),
		})
	})
}
