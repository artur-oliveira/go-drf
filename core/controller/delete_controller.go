package controllers

import "github.com/gofiber/fiber/v2"

// Delete (DELETE /:id)
func (h *GenericController[T, C, U, P, R, F]) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var record T // Necessário para GORM inferir a tabela

	// GORM executa a deleção (Soft ou Hard)
	tx := h.DB.Delete(&record, id)
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + tx.Error.Error()})
	}
	if tx.RowsAffected == 0 {
		// Se RowsAffected for 0, o registro não existia
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Registro não encontrado"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
