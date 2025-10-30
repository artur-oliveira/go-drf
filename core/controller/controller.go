package controllers

import (
	"grf/core/dto"
	"grf/core/filterset"
	"grf/core/models"
	"grf/core/pagination"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type GenericController[T models.IModel, C any, U any, P dto.IPatchDTO, R any, F filterset.IFilterSet] struct {
	DB        *gorm.DB
	Validator *validator.Validate
	Paginator pagination.IPagination[T] // Paginador injetado

	NewFilterSet func() F
	NewPatchDTO  func() P

	MapToResponse    func(model T) R
	MapCreateToModel func(dto C) T
	MapUpdateToModel func(dto U, model T) T
}

// ControllerConfig é usado para construir o GenericController (Injeção de Dependência)
type ControllerConfig[T models.IModel, C any, U any, P dto.IPatchDTO, R any, F filterset.IFilterSet] struct {
	DB        *gorm.DB
	Validator *validator.Validate
	Paginator pagination.IPagination[T]

	NewFilterSet func() F
	NewPatchDTO  func() P

	MapToResponse    func(model T) R
	MapCreateToModel func(dto C) T
	MapUpdateToModel func(dto U, model T) T
}

// NewGenericController cria o controlador com todas as dependências.
func NewGenericController[T models.IModel, C any, U any, P dto.IPatchDTO, R any, F filterset.IFilterSet](
	config *ControllerConfig[T, C, U, P, R, F],
) *GenericController[T, C, U, P, R, F] {

	if config.DB == nil || config.Validator == nil ||
		config.NewFilterSet == nil || config.NewPatchDTO == nil ||
		config.MapToResponse == nil || config.MapCreateToModel == nil ||
		config.MapUpdateToModel == nil {
		panic("Controlador Genérico: DB, Validator, Fábricas e Mapeadores não podem ser nulos")
	}

	return &GenericController[T, C, U, P, R, F]{
		DB:               config.DB,
		Validator:        config.Validator,
		Paginator:        config.Paginator,
		NewFilterSet:     config.NewFilterSet,
		NewPatchDTO:      config.NewPatchDTO,
		MapToResponse:    config.MapToResponse,
		MapCreateToModel: config.MapCreateToModel,
		MapUpdateToModel: config.MapUpdateToModel,
	}
}
