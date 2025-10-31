package controllers

import (
	"grf/core/pagination"

	"github.com/gofiber/fiber/v2"
)

func (h *GenericController[T, C, U, P, R, F]) List(c *fiber.Ctx) error {
	filters := h.NewFilterSet()
	if err := filters.Bind(c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var model T
	query := h.DB.Model(&model)
	query = filters.Apply(query)

	if h.Paginator == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Paginador não configurado"})
	}

	paginatedResponse, err := h.Paginator.Paginate(c, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	responseDTOs := make([]R, len(paginatedResponse.Results))
	for i, item := range paginatedResponse.Results {
		responseDTOs[i] = h.MapToResponse(item)
	}

	finalResponse := pagination.Response[R]{
		Results: responseDTOs,
		HasNext: paginatedResponse.HasNext,
		Count:   paginatedResponse.Count,
	}

	return c.JSON(finalResponse)
}
