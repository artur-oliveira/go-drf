package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PartialUpdate (PATCH /:id)
func (h *GenericController[T, C, U, P, R, F]) PartialUpdate(c *fiber.Ctx) error {
	id := c.Params("id")
	var record T

	// 1. Verificar se existe
	if err := h.DB.First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Registro não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	// 2. Instanciar (fábrica), Bindar e Validar DTO de Patch (P)
	patchInput := h.NewPatchDTO() // P é IPatchDTO

	if err := c.BodyParser(patchInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload inválido: " + err.Error()})
	}
	if err := h.Validator.Struct(patchInput); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validação falhou: " + err.Error()})
	}

	// 3. Converter DTO (ponteiros) para Mapa (via interface)
	patchMap := patchInput.ToPatchMap()

	// 4. Verificar se há atualizações
	if len(patchMap) == 0 {
		// Nada a atualizar, retorna o registro atual
		response := h.MapToResponse(record)
		return c.JSON(response)
	}

	// 5. Aplicar o Patch no DB (Operação segura)
	if err := h.DB.Model(&record).Updates(patchMap).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	// 6. Serializar Resposta (T -> R)
	response := h.MapToResponse(record) // 'record' é atualizado in-place
	return c.JSON(response)
}
