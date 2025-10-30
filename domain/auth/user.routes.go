package auth

import (
	controllers "grf/core/controller"
	"grf/core/pagination"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterUserRoutes encapsula a configuração e registro do UserController.
func RegisterUserRoutes(router fiber.Router, db *gorm.DB, validate *validator.Validate) {

	// --- Toda a lógica de configuração do `main.go` anterior está AQUI ---

	// 1. Definir Tipos Genéricos Concretos
	// T = models.User
	// C = dtos.UserCreateDTO
	// U = dtos.UserUpdateDTO
	// P = *UserPatchDTO
	// R = dtos.UserResponseDTO
	// F = *UserFilterSet

	// 2. Instanciar Paginador Específico
	userPaginator := pagination.NewLimitOffsetPagination[User](10, 50)

	// 3. Criar a Configuração do Controlador
	userConfig := &controllers.ControllerConfig[
		User,
		UserCreateDTO,
		UserUpdateDTO,
		*UserPatchDTO,
		UserResponseDTO,
		*UserFilterSet,
	]{
		DB:        db,
		Validator: validate,
		Paginator: userPaginator,

		// Fábricas
		NewFilterSet: func() *UserFilterSet {
			return new(UserFilterSet)
		},
		NewPatchDTO: func() *UserPatchDTO {
			return new(UserPatchDTO)
		},

		// Mapeadores (importados do pacote de mappers)
		MapToResponse:    MapUserToResponse,
		MapCreateToModel: MapCreateToUser,
		MapUpdateToModel: MapUpdateToUser,
	}

	// 4. Instanciar o Controlador Genérico
	userController := controllers.NewGenericController(userConfig)

	// 5. Registrar Rotas (no grupo /users)
	userRoutes := router.Group("/users")

	// (Aqui você também adicionaria Middlewares de autenticação/permissão)
	// ex: userRoutes.Use(middlewares.MwAuthenticate)

	userRoutes.Get("/", userController.List)
	userRoutes.Post("/", userController.Create)
	userRoutes.Get("/:id", userController.Retrieve)
	userRoutes.Put("/:id", userController.Update)
	userRoutes.Patch("/:id", userController.PartialUpdate)
	userRoutes.Delete("/:id", userController.Delete)
}
