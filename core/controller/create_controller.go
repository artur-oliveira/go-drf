package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// Create (POST /)
func (h *GenericController[T, C, U, P, R, F]) Create(c *fiber.Ctx) error {
	// 1. Bind e Validar DTO (C)
	var input C
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload inválido: " + err.Error()})
	}
	if err := h.Validator.Struct(input); err != nil {
		// (Formatar erro de validação)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validação falhou: " + err.Error()})
	}

	// 2. Mapear DTO -> Model (usando o mapeador)
	newRecord := h.MapCreateToModel(input)

	// 3. Salvar no DB
	if err := h.DB.Create(&newRecord).Error; err != nil {
		// (Tratar erros de constraint, ex: unique)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	// 4. Serializar Resposta (T -> R)
	response := h.MapToResponse(newRecord)
	return c.Status(fiber.StatusCreated).JSON(response)
}
