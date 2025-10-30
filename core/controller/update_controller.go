package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Update (PUT /:id)
func (h *GenericController[T, C, U, P, R, F]) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var record T

	// 1. Verificar se existe
	if err := h.DB.First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Registro não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	// 2. Bind e Validar DTO de Update (U)
	var input U
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload inválido: " + err.Error()})
	}
	if err := h.Validator.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validação falhou: " + err.Error()})
	}

	// 3. Mapear DTO -> Model (usando o mapeador de atualização)
	// O mapeador mescla 'input' em 'record', preservando ID/CreatedAt
	updatedRecord := h.MapUpdateToModel(input, record)

	// 4. Salvar (Substituição total)
	if err := h.DB.Save(&updatedRecord).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	// 5. Serializar Resposta (T -> R)
	response := h.MapToResponse(updatedRecord)
	return c.JSON(response)
}
