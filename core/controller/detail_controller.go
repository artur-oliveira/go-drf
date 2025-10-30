package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Retrieve (GET /:id)
func (h *GenericController[T, C, U, P, R, F]) Retrieve(c *fiber.Ctx) error {
	id := c.Params("id")
	var record T

	if err := h.DB.First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Registro não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	// Mapear T -> R
	response := h.MapToResponse(record)
	return c.JSON(response)
}
