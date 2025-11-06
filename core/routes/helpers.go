package routes

import (
	"grf/core/controller"
	"grf/core/middleware"
	"grf/core/models"
	"grf/core/permission"
	"grf/core/server"

	"github.com/gofiber/fiber/v2"
)

type RegisterModelOptions struct {
	App    *server.App
	Router fiber.Router
	Path   string

	Controller controller.ICRUDController

	Model models.IModel

	Permission permission.IPermission
}

func RegisterModelController(opts *RegisterModelOptions) {

	if opts.App == nil || opts.Router == nil || opts.Controller == nil || opts.Path == "" {
		panic("RegisterModelController: App, Router, Path e Controller são obrigatórios")
	}

	var perm permission.IPermission

	if opts.Permission != nil {
		perm = opts.Permission
	} else {
		if opts.Model == nil {
			panic("RegisterModelController: Model é obrigatório se a permissão customizada não for fornecida")
		}
		perm = permission.NewAnd(
			opts.App.IsAuthenticated,
			permission.NewModelPermissions(opts.App.DB, opts.Model),
		)
	}

	routes := opts.Router.Group(opts.Path)
	routes.Use(middleware.Check(perm))
	RegisterCRUDController(routes, opts.Controller)
}

func RegisterCRUDController(
	router fiber.Router,
	controller controller.ICRUDController,
) {
	router.Get("/", controller.List)
	router.Post("/", controller.Create)
	router.Get("/:id", controller.Retrieve)
	router.Put("/:id", controller.Update)
	router.Patch("/:id", controller.PartialUpdate)
	router.Delete("/:id", controller.Delete)
}
