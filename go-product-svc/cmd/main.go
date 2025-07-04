package main

import (
	"context"
	"fmt"
	"log"
	"poc-product-svc/internal/infrastructure/config"
	"poc-product-svc/internal/interfaces/di"
	"poc-product-svc/internal/interfaces/http/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or error loading it:", err)

	}
	app := di.NewContainer()
	if err := app.Start(context.Background()); err != nil {
		log.Fatal("Failed to start application:", err)
	}

	defer func() {
		if err := app.Stop(context.Background()); err != nil {
			log.Printf("Failed to stop application gracefully: %v", err)
		}
	}()

	// Mantener la aplicación corriendo
	<-app.Done()
}

// startHTTPServer inicia el servidor HTTP Fiber
func startHTTPServer(lc fx.Lifecycle, cfg *config.Config, router *routes.Router) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Crear instancia de Fiber
			app := fiber.New(fiber.Config{
				AppName:      "Product Management API",
				ServerHeader: "Fiber",
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					code := fiber.StatusInternalServerError
					if e, ok := err.(*fiber.Error); ok {
						code = e.Code
					}

					return c.Status(code).JSON(fiber.Map{
						"error":   true,
						"message": err.Error(),
					})
				},
			})

			// Middleware globales
			app.Use(recover.New())
			app.Use(logger.New(logger.Config{
				Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
			}))
			app.Use(cors.New(cors.Config{
				AllowOrigins: "*",
				AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
				AllowHeaders: "Origin,Content-Type,Accept,Authorization",
			}))

			router.Setup(app)

			// Iniciar servidor en goroutine
			go func() {
				if err := app.Listen(":" + cfg.Server.Port); err != nil {
					log.Printf("Failed to start server: %v", err)
				}
			}()

			fmt.Printf("🚀 Server running on port %s\n", cfg.Server.Port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("🛑 Shutting down server...")
			return nil
		},
	})
}

func init() {
	// Agregar el hook del servidor HTTP al contenedor
	di.HTTPModule = fx.Module("http",
		fx.Provide(routes.NewRouter),
		fx.Invoke(startHTTPServer),
	)
}
