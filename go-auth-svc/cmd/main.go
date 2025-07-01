package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"poc-auth-svc/internal/application/usecases"
	"poc-auth-svc/internal/domain/services"
	"poc-auth-svc/internal/infrastructure/database"
	"poc-auth-svc/internal/infrastructure/http/handlers"
	"poc-auth-svc/internal/infrastructure/http/routes"
	"poc-auth-svc/internal/infrastructure/persistence"
	"poc-auth-svc/internal/infrastructure/security"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or error loading it:", err)

	}
	// Configuración desde variables de entorno
	mongoUri := getEnv("MONGO_URI", "")
	dbName := getEnv("DB_NAME", "auth_svc")
	jwtSecret := getEnv("JWT_SECRET", "12454sd32")
	port := getEnv("PORT", "3000")
	fmt.Println(jwtSecret)
	// Conectar a mongoDB
	mongoClient, err := database.NewMongoClient(mongoUri)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}
	defer mongoClient.Disconnect(context.Background())
	db := mongoClient.Database(dbName)

	// Inicializar dependencias (Dependency Injection)
	hasher := security.NewBcryptHasher()
	jwtExpirationHours, _ := strconv.ParseInt(getEnv("JWT_EXPIRATION_HOURS", "2"), 10, 64)
	jwtWrapper := usecases.JwtWrapper{
		SecretKey:       jwtSecret,
		Issuer:          getEnv("JWT_ISSUER", "go"),
		ExpirationHours: jwtExpirationHours,
	}
	userRepo := persistence.NewMongoUserRepository(db)
	authService := services.NewAuthService(userRepo, hasher)
	authUseCase := usecases.NewAuthUseCase(authService, jwtWrapper)
	authHandler := handlers.NewAuthHandler(authUseCase)

	// Configurar fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middlewares
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth_service",
		})
	})

	routes.SetupRoutes(app, authHandler)
	log.Printf("Auth service running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
