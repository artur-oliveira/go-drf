package routes

import (
	"grf/core/server"
	"grf/domain/auth"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(
	router fiber.Router,
	app *server.App,
) {
	userController := auth.NewDefaultUserController(app.DB, app.Validator)
	groupController := auth.NewDefaultGroupController(app.DB, app.Validator)
	permissionController := auth.NewDefaultPermissionController(app.DB, app.Validator)
	authController := auth.NewAuthController(app.DB, app.Config, app.Validator)

	permMw := app.PermMw
	authMw := app.AuthMw

	authRoutes := router.Group("/auth")
	authRoutes.Post("/token", authController.ObtainToken)
	authRoutes.Post("/refresh", authController.ObtainTokenRefresh)

	authRoutes.Get("/me", authMw.RequireAuth, authController.GetMe)
	authRoutes.Post("/change-password", authMw.RequireAuth, authController.ChangePassword)

	userRoutes := router.Group("/users")
	userRoutes.Use(authMw.RequireAuth)

	userRoutes.Get("/", permMw.RequirePerm("auth", "view_user"), userController.List)
	userRoutes.Post("/", permMw.RequirePerm("auth", "add_user"), userController.Create)
	userRoutes.Get("/:id", permMw.RequirePerm("auth", "view_user"), userController.Retrieve)
	userRoutes.Put("/:id", permMw.RequirePerm("auth", "change_user"), userController.Update)
	userRoutes.Patch("/:id", permMw.RequirePerm("auth", "change_user"), userController.PartialUpdate)
	userRoutes.Delete("/:id", permMw.RequirePerm("auth", "delete_user"), userController.Delete)

	groupRoutes := router.Group("/groups")
	groupRoutes.Use(app.AuthMw.RequireAuth)

	groupRoutes.Get("/", permMw.RequirePerm("auth", "view_group"), groupController.List)
	groupRoutes.Post("/", permMw.RequirePerm("auth", "add_group"), groupController.Create)
	groupRoutes.Get("/:id", permMw.RequirePerm("auth", "view_group"), groupController.Retrieve)
	groupRoutes.Put("/:id", permMw.RequirePerm("auth", "change_group"), groupController.Update)
	groupRoutes.Patch("/:id", permMw.RequirePerm("auth", "change_group"), groupController.PartialUpdate)
	groupRoutes.Delete("/:id", permMw.RequirePerm("auth", "delete_group"), groupController.Delete)

	permissionRoutes := router.Group("/permissions")
	permissionRoutes.Use(app.AuthMw.RequireAuth)

	permissionRoutes.Get("/", permMw.RequirePerm("auth", "view_permission"), permissionController.List)
	permissionRoutes.Post("/", permMw.RequirePerm("auth", "add_permission"), permissionController.Create)
	permissionRoutes.Get("/:id", permMw.RequirePerm("auth", "view_permission"), permissionController.Retrieve)
	permissionRoutes.Put("/:id", permMw.RequirePerm("auth", "change_permission"), permissionController.Update)
	permissionRoutes.Patch("/:id", permMw.RequirePerm("auth", "change_permission"), permissionController.PartialUpdate)
	permissionRoutes.Delete("/:id", permMw.RequirePerm("auth", "delete_permission"), permissionController.Delete)
}
