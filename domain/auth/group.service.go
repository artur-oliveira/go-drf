package auth

import (
	"fmt"
	"grf/core/exceptions"
	"grf/core/service"

	"gorm.io/gorm"
)

type IDType = uint64

type GroupService struct {
	service.IService[*Group, *GroupCreateDTO, *GroupUpdateDTO, *GroupPatchDTO, *GroupResponseDTO, *GroupFilterSet, IDType]

	DB *gorm.DB
}

func NewGroupService(
	config *service.ServiceConfig[*Group, *GroupCreateDTO, *GroupUpdateDTO, *GroupPatchDTO, *GroupResponseDTO, *GroupFilterSet, IDType],
	db *gorm.DB,
) service.IService[*Group, *GroupCreateDTO, *GroupUpdateDTO, *GroupPatchDTO, *GroupResponseDTO, *GroupFilterSet, IDType] {

	baseService := service.NewGenericService(config)

	return &GroupService{
		IService: baseService,
		DB:       db,
	}
}

func (s *GroupService) Create(dto *GroupCreateDTO) (*Group, error) {
	newRecord := MapCreateToGroup(dto)

	var permissions []*Permission
	if len(dto.PermissionIDs) > 0 {
		if err := s.DB.Where("id IN ?", dto.PermissionIDs).Find(&permissions).Error; err != nil {
			return nil, exceptions.NewInternal(fmt.Errorf("erro ao buscar permissões: %w", err))
		}
		if len(permissions) != len(dto.PermissionIDs) {
			return nil, exceptions.NewBadRequest("Uma ou mais PermissionIDs são inválidas", nil)
		}
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newRecord).Error; err != nil {
			return err
		}
		if len(permissions) > 0 {
			if err := tx.Model(newRecord).Association("Permissions").Replace(permissions); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, exceptions.NewInternal(err)
	}

	s.DB.Preload("Permissions").First(newRecord, newRecord.ID)
	return newRecord, nil
}

func (s *GroupService) Update(id IDType, dto *GroupUpdateDTO) (*Group, error) {
	record, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updatedRecord := MapUpdateToGroup(dto, record)

	var permissions []*Permission
	if len(dto.PermissionIDs) > 0 {
		if err := s.DB.Where("id IN ?", dto.PermissionIDs).Find(&permissions).Error; err != nil {
			return nil, exceptions.NewInternal(fmt.Errorf("erro ao buscar permissões: %w", err))
		}
		if len(permissions) != len(dto.PermissionIDs) {
			return nil, exceptions.NewBadRequest("Uma ou mais PermissionIDs são inválidas", nil)
		}
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(updatedRecord).Error; err != nil {
			return err
		}
		if err := tx.Model(updatedRecord).Association("Permissions").Replace(permissions); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, exceptions.NewInternal(err)
	}

	s.DB.Preload("Permissions").First(updatedRecord, updatedRecord.ID)
	return updatedRecord, nil
}

func (s *GroupService) PartialUpdate(id IDType, dto *GroupPatchDTO) (*Group, error) {
	record, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	patchMap := dto.ToPatchMap()

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if len(patchMap) > 0 {
			if err := tx.Model(record).Updates(patchMap).Error; err != nil {
				return err
			}
		}

		if dto.PermissionIDs != nil {
			var permissions []*Permission
			if len(dto.PermissionIDs) > 0 {
				if err := tx.Where("id IN ?", dto.PermissionIDs).Find(&permissions).Error; err != nil {
					return exceptions.NewBadRequest("Erro ao buscar permissões", err)
				}
				if len(permissions) != len(dto.PermissionIDs) {
					return exceptions.NewBadRequest("Uma ou mais PermissionIDs são inválidas", nil)
				}
			}
			if err := tx.Model(record).Association("Permissions").Replace(permissions); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, exceptions.NewInternal(err)
	}

	s.DB.Preload("Permissions").First(record, record.ID)
	return record, nil
}
