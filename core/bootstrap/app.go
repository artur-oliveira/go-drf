package bootstrap

import (
	"grf/core/auth"
	"grf/core/config"
	"grf/core/database"
	"grf/core/exceptions"
	"grf/core/i18n"
	"grf/core/middleware"
	"grf/core/routes"
	"grf/core/server"
	"grf/core/validator"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewApp(cfg config.Config) (*server.App, error) {

	db, err := database.ConnectDB(&cfg)
	if err != nil {
		return nil, err
	}

	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: exceptions.GlobalErrorHandler,
		IdleTimeout:  time.Duration(cfg.ServerIdleTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ServerIdleTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerIdleTimeout) * time.Second,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(logger.New())

	authMw := middleware.NewAuthenticationMiddleware(
		auth.NewJWTAuthBackend(db, &cfg),
		auth.NewBasicAuthBackend(db),
	)
	permMw := middleware.NewPermissionMiddleware(db)
	i18nMw := middleware.NewI18NMiddleware(i18n.NewI18nService())

	var bootstrapedApp = &server.App{
		FiberApp:  app,
		DB:        db,
		Validator: validator.GetValidator(),
		I18nMw:    i18nMw,
		Config:    &cfg,
		AuthMw:    authMw,
		PermMw:    permMw,
	}

	i18nMw.UseMiddleWare(
		app,
	)
	database.RegisterMigrations(
		bootstrapedApp,
	)
	routes.RegisterRoutes(
		bootstrapedApp,
	)
	return bootstrapedApp, nil
}
