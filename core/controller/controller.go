package controller

import (
	"errors"
	"grf/core/dto"
	"grf/core/exceptions"
	"grf/core/filterset"
	"grf/core/models"
	"grf/core/pagination"
	"grf/core/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type GenericController[T models.IModel, C any, U any, P dto.IPatchDTO, R any, F filterset.IFilterSet, ID comparable] struct {
	Service   service.IService[T, C, U, P, R, F, ID]
	Validator *validator.Validate
	Paginator pagination.IPagination[T]

	MapToResponse func(model T) R

	NewFilterSet func() F
	NewPatchDTO  func() P

	ParseID func(s string) (ID, error)
}

type ControllerConfig[T models.IModel, C any, U any, P dto.IPatchDTO, R any, F filterset.IFilterSet, ID comparable] struct {
	Service   service.IService[T, C, U, P, R, F, ID]
	Validator *validator.Validate
	Paginator pagination.IPagination[T]

	MapToResponse func(model T) R

	NewFilterSet func() F
	NewPatchDTO  func() P
	ParseID      func(s string) (ID, error)
}

func NewGenericController[T models.IModel, C any, U any, P dto.IPatchDTO, R any, F filterset.IFilterSet, ID comparable](
	config *ControllerConfig[T, C, U, P, R, F, ID],
) *GenericController[T, C, U, P, R, F, ID] {

	if config.Service == nil || config.Validator == nil || config.ParseID == nil || config.MapToResponse == nil {
		panic("Controlador Genérico: Service, Validator, ParseID e MapToResponse não podem ser nulos")
	}

	return &GenericController[T, C, U, P, R, F, ID]{
		Service:       config.Service,
		Validator:     config.Validator,
		Paginator:     config.Paginator,
		MapToResponse: config.MapToResponse,
		NewFilterSet:  config.NewFilterSet,
		NewPatchDTO:   config.NewPatchDTO,
		ParseID:       config.ParseID,
	}
}

func (h *GenericController[T, C, U, P, R, F, ID]) List(c *fiber.Ctx) error {
	filters := h.NewFilterSet()
	if err := filters.Bind(c); err != nil {
		return exceptions.NewBadRequest("Parâmetros de filtro inválidos", err)
	}
	if h.Paginator == nil {
		return exceptions.NewInternal(errors.New("paginador não configurado"))
	}
	if err := h.Paginator.Bind(c); err != nil {
		return exceptions.NewBadRequest("Parâmetros de paginação inválidos", err)
	}

	paginatedResponse, err := h.Service.List(filters, h.Paginator)
	if err != nil {
		return err
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

func (h *GenericController[T, C, U, P, R, F, ID]) Create(c *fiber.Ctx) error {
	var input C
	if err := c.BodyParser(&input); err != nil {
		return exceptions.NewBadRequest("Payload inválido", err)
	}
	if err := h.Validator.Struct(input); err != nil {
		return err
	}

	newRecord, err := h.Service.Create(input)
	if err != nil {
		return err
	}

	response := h.MapToResponse(newRecord)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *GenericController[T, C, U, P, R, F, ID]) Retrieve(c *fiber.Ctx) error {
	id, err := h.ParseID(c.Params("id"))
	if err != nil {
		return exceptions.NewBadRequest("ID inválido", err)
	}

	record, err := h.Service.GetByID(id)
	if err != nil {
		return err
	}

	response := h.MapToResponse(record)
	return c.JSON(response)
}

func (h *GenericController[T, C, U, P, R, F, ID]) Update(c *fiber.Ctx) error {
	id, err := h.ParseID(c.Params("id"))
	if err != nil {
		return exceptions.NewBadRequest("ID inválido", err)
	}

	var input U
	if err := c.BodyParser(&input); err != nil {
		return exceptions.NewBadRequest("Payload inválido", err)
	}
	if err := h.Validator.Struct(input); err != nil {
		return err
	}

	updatedRecord, err := h.Service.Update(id, input)
	if err != nil {
		return err
	}

	response := h.MapToResponse(updatedRecord)
	return c.JSON(response)
}

func (h *GenericController[T, C, U, P, R, F, ID]) PartialUpdate(c *fiber.Ctx) error {
	id, err := h.ParseID(c.Params("id"))
	if err != nil {
		return exceptions.NewBadRequest("ID inválido", err)
	}

	patchInput := h.NewPatchDTO()
	if err := c.BodyParser(patchInput); err != nil {
		return exceptions.NewBadRequest("Payload inválido", err)
	}
	if err := h.Validator.Struct(patchInput); err != nil {
		return err
	}

	updatedRecord, err := h.Service.PartialUpdate(id, patchInput)
	if err != nil {
		return err
	}

	response := h.MapToResponse(updatedRecord)
	return c.JSON(response)
}

func (h *GenericController[T, C, U, P, R, F, ID]) Delete(c *fiber.Ctx) error {
	id, err := h.ParseID(c.Params("id"))
	if err != nil {
		return exceptions.NewBadRequest("ID inválido", err)
	}

	if err := h.Service.Delete(id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
