package bootstrap

import (
	"grf/bootstrap/database"
	"grf/bootstrap/grf"
	"grf/config"
	"grf/core/exceptions"
	"grf/core/middleware"
	"grf/routes"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewApp() (*grf.App, error) {
	cfg, err := config.LoadConfig("./")
	if err != nil {
		return nil, err
	}
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

	authMw := middleware.NewAuthMiddleware(db, &cfg)
	permMw := middleware.NewPermissionMiddleware(db)

	var bootstrapedApp = &grf.App{
		FiberApp:  app,
		DB:        db,
		Validator: GetValidator(),
		Config:    &cfg,
		AuthMw:    authMw,
		PermMw:    permMw,
	}
	database.RegisterMigrations(
		bootstrapedApp,
	)
	routes.RegisterRoutes(
		bootstrapedApp,
	)
	return bootstrapedApp, nil
}
