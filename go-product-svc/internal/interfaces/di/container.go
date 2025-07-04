package di

import (
	"context"
	"fmt"
	"poc-product-svc/internal/application/usecases"
	"poc-product-svc/internal/domain/repositories"
	"poc-product-svc/internal/domain/services"
	"poc-product-svc/internal/infrastructure/config"
	"poc-product-svc/internal/infrastructure/database"
	"poc-product-svc/internal/infrastructure/persistence/adapters"
	"poc-product-svc/internal/infrastructure/persistence/postgres"
	"poc-product-svc/internal/interfaces/http/handlers"
	"poc-product-svc/internal/interfaces/http/routes"
	"poc-product-svc/internal/shared/utils"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var ConfigModule = fx.Module("config", fx.Provide(config.Load))
var DatabaseModule = fx.Module("database",
	fx.Provide(
		database.NewPostgresDB,
		func(db *database.PostgresDB) adapters.ORMAdapter {
			return adapters.NewGormAdapter(db.DB)
		},
		func(db *database.PostgresDB) adapters.DatabaseAdapter {
			return adapters.NewGormDatabaseAdapter(db.DB)
		},
	),
)

var RepositoryModule = fx.Module("repository",
	fx.Provide(
		fx.Annotate(
			postgres.NewProductRepository,
			fx.As(new(repositories.ProductRepository)),
		),
	),
)

var ServiceModule = fx.Module("service",
	fx.Provide(
		services.NewProductService,
	),
)

var UseCaseModule = fx.Module("usecase",
	fx.Provide(usecases.NewProductUseCase),
)

var HandlerModule = fx.Module("handler",
	fx.Provide(utils.NewValidationService),
	fx.Provide(handlers.NewProductHandler),
)

var HTTPModule = fx.Module("http", fx.Provide(routes.NewRouter))

func NewContainer() *fx.App {
	return fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),

		fx.Provide(zap.NewProduction),

		ConfigModule,
		DatabaseModule,
		RepositoryModule,
		ServiceModule,
		UseCaseModule,
		HandlerModule,
		HTTPModule,

		fx.Invoke(registerHooks),
	)
}

func registerHooks(
	lc fx.Lifecycle,
	dbAdapter adapters.DatabaseAdapter,
	logger *zap.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting application...")

			// Verificar conexión a la base de datos
			if err := dbAdapter.Ping(ctx); err != nil {
				return fmt.Errorf("failed to ping database: %w", err)
			}

			logger.Info("Database connection verified")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping application...")

			// Cerrar conexiones
			if err := dbAdapter.Close(); err != nil {
				logger.Error("Failed to close database connection", zap.Error(err))
				return err
			}

			logger.Info("Application stopped gracefully")
			return nil
		},
	})
}
