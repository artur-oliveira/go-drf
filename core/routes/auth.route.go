package routes

import (
	"grf/core/middleware"
	"grf/core/permission"
	"grf/core/server"
	"grf/domain/auth/controller"
	"grf/domain/auth/model"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(
	router fiber.Router,
	app *server.App,
) {
	Check := middleware.Check
	IsAuthenticated := app.IsAuthenticated
	IsAdmin := app.IsAdmin

	userController := controller.NewDefaultUserController(app.DB, app.Validator)
	groupController := controller.NewDefaultGroupController(app.DB, app.Validator)
	permissionController := controller.NewDefaultPermissionController(app.DB, app.Validator)
	authController := controller.NewAuthController(app.DB, app.Config, app.Validator)

	authRoutes := router.Group("/auth")
	authRoutes.Post("/token", authController.ObtainToken)
	authRoutes.Post("/refresh", authController.ObtainTokenRefresh)
	authRoutes.Get("/me", Check(IsAuthenticated), authController.GetMe)
	authRoutes.Post("/change-password", Check(IsAuthenticated), authController.ChangePassword)

	RegisterModelController(&RegisterModelOptions{
		App:        app,
		Router:     router,
		Path:       "/users",
		Model:      new(model.User),
		Controller: userController,
	})

	RegisterModelController(&RegisterModelOptions{
		App:        app,
		Router:     router,
		Path:       "/groups",
		Model:      new(model.Group),
		Controller: groupController,
	})

	adminOnlyPerm := permission.NewAnd(IsAuthenticated, IsAdmin)
	RegisterModelController(&RegisterModelOptions{
		App:        app,
		Router:     router,
		Path:       "/permissions",
		Model:      new(model.Permission),
		Controller: permissionController,
		Permission: adminOnlyPerm,
	})
}
