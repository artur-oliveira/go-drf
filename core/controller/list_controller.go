package controllers

import (
	"grf/core/pagination"

	"github.com/gofiber/fiber/v2"
)

// List (GET /)
func (h *GenericController[T, C, U, P, R, F]) List(c *fiber.Ctx) error {
	// 1. Instanciar Filtros (usando a fábrica)
	filters := h.NewFilterSet()
	if err := filters.Bind(c); err != nil {
		// Erro de Bind (ex: bad request)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 2. Preparar Query
	var model T
	query := h.DB.Model(&model)
	query = filters.Apply(query)

	// 3. Paginação (Injetada)
	if h.Paginator == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Paginador não configurado"})
	}

	// O Paginator retorna T (Model)
	paginatedResponse, err := h.Paginator.Paginate(c, query)
	if err != nil {
		// Erro de DB (na Paginação, ex: Count() ou Find())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro de banco de dados: " + err.Error()})
	}

	// 4. Serializar Resposta (Mapear T[] -> R[])
	responseDTOs := make([]R, len(paginatedResponse.Results))
	for i, item := range paginatedResponse.Results {
		responseDTOs[i] = h.MapToResponse(item)
	}

	// 5. Recriar a resposta paginada com os DTOs
	finalResponse := pagination.Response[R]{
		Results: responseDTOs,
		HasNext: paginatedResponse.HasNext,
		Count:   paginatedResponse.Count,
	}

	return c.JSON(finalResponse)
}
