package auth

import (
	"errors"
	"fmt"
	controllers "grf/core/controller"
	"grf/core/pagination"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GroupController struct {
	*controllers.GenericController[
		*Group,
		*GroupCreateDTO,
		*GroupUpdateDTO,
		*GroupPatchDTO,
		*GroupResponseDTO,
		*GroupFilterSet,
	]
}

func NewDefaultGroupController(
	db *gorm.DB,
	validator *validator.Validate,
) *GroupController {
	GroupPaginator := pagination.NewLimitOffsetPagination[*Group](
		10,
		100,
	)

	GroupConfig := &controllers.ControllerConfig[
		*Group,
		*GroupCreateDTO,
		*GroupUpdateDTO,
		*GroupPatchDTO,
		*GroupResponseDTO,
		*GroupFilterSet,
	]{
		DB:        db,
		Validator: validator,
		Paginator: GroupPaginator,

		NewFilterSet: func() *GroupFilterSet {
			return new(GroupFilterSet)
		},
		NewPatchDTO: func() *GroupPatchDTO {
			return new(GroupPatchDTO)
		},

		MapToResponse:    MapGroupToResponse,
		MapCreateToModel: MapCreateToGroup,
		MapUpdateToModel: MapUpdateToGroup,
	}
	return NewGroupController(GroupConfig)
}

func NewGroupController(
	config *controllers.ControllerConfig[
		*Group,
		*GroupCreateDTO,
		*GroupUpdateDTO,
		*GroupPatchDTO,
		*GroupResponseDTO,
		*GroupFilterSet,
	],
) *GroupController {
	generic := controllers.NewGenericController(config)

	return &GroupController{
		GenericController: generic,
	}
}

func (h *GroupController) Create(c *fiber.Ctx) error {
	var input GroupCreateDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload inválido: " + err.Error()})
	}
	if err := h.Validator.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validação falhou: " + err.Error()})
	}

	newRecord := h.MapCreateToModel(&input)

	var permissions []Permission
	if len(input.PermissionIDs) > 0 {
		if err := h.DB.Where("id IN ?", input.PermissionIDs).Find(&permissions).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro ao buscar permissões: " + err.Error()})
		}
		if len(permissions) != len(input.PermissionIDs) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Uma ou mais PermissionIDs são inválidas"})
		}
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newRecord).Error; err != nil {
			return err
		}

		if len(permissions) > 0 {
			if err := tx.Model(&newRecord).Association("Permissions").Replace(permissions); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	h.DB.Preload("Permissions").First(&newRecord, newRecord.ID)
	response := h.MapToResponse(newRecord)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *GroupController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var record Group
	if err := h.DB.First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Registro não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	var input GroupUpdateDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload inválido: " + err.Error()})
	}
	if err := h.Validator.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validação falhou: " + err.Error()})
	}

	updatedRecord := h.MapUpdateToModel(&input, &record)

	var permissions []Permission
	if len(input.PermissionIDs) > 0 {
		if err := h.DB.Where("id IN ?", input.PermissionIDs).Find(&permissions).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro ao buscar permissões: " + err.Error()})
		}
		if len(permissions) != len(input.PermissionIDs) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Uma ou mais PermissionIDs são inválidas"})
		}
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&updatedRecord).Error; err != nil {
			return err
		}

		if err := tx.Model(&updatedRecord).Association("Permissions").Replace(permissions); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	h.DB.Preload("Permissions").First(&updatedRecord, updatedRecord.ID)
	response := h.MapToResponse(updatedRecord)
	return c.JSON(response)
}

func (h *GroupController) PartialUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	var record Group
	if err := h.DB.Preload("Permissions").First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Registro não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	var rawInput map[string]interface{}
	if err := c.BodyParser(&rawInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload inválido: " + err.Error()})
	}

	patchInput := h.NewPatchDTO()
	if err := c.BodyParser(patchInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload inválido: " + err.Error()})
	}
	if err := h.Validator.Struct(patchInput); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validação falhou: " + err.Error()})
	}

	patchMap := patchInput.ToPatchMap()
	_, permissionsSent := rawInput["permission_ids"]

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if len(patchMap) > 0 {
			if err := tx.Model(&record).Updates(patchMap).Error; err != nil {
				return err
			}
		}

		if permissionsSent {
			var permissions []Permission
			if len(patchInput.PermissionIDs) > 0 {
				if err := tx.Where("id IN ?", patchInput.PermissionIDs).Find(&permissions).Error; err != nil {
					return fmt.Errorf("erro ao buscar permissões: %v", err)
				}
				if len(permissions) != len(patchInput.PermissionIDs) {
					return fmt.Errorf("uma ou mais PermissionIDs são inválidas")
				}
			}
			if err := tx.Model(&record).Association("Permissions").Replace(permissions); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	h.DB.Preload("Permissions").First(&record, record.ID)
	response := h.MapToResponse(&record)
	return c.JSON(response)
}
