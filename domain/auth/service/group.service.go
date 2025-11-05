package service

import (
	"grf/core/exceptions"
	"grf/core/service"
	"grf/domain/auth/dto"
	"grf/domain/auth/filter"
	"grf/domain/auth/mapper"
	"grf/domain/auth/model"

	"gorm.io/gorm"
)

type GroupService struct {
	service.IService[*model.Group, *dto.GroupCreateDTO, *dto.GroupUpdateDTO, *dto.GroupPatchDTO, *dto.GroupResponseDTO, *filter.GroupFilterSet, uint64]

	DB *gorm.DB
}

func NewGroupService(
	config *service.Config[*model.Group, *dto.GroupCreateDTO, *dto.GroupUpdateDTO, *dto.GroupPatchDTO, *dto.GroupResponseDTO, *filter.GroupFilterSet, uint64],
	db *gorm.DB,
) service.IService[*model.Group, *dto.GroupCreateDTO, *dto.GroupUpdateDTO, *dto.GroupPatchDTO, *dto.GroupResponseDTO, *filter.GroupFilterSet, uint64] {

	baseService := service.NewGenericService(config)

	return &GroupService{
		IService: baseService,
		DB:       db,
	}
}

func (s *GroupService) Create(dto *dto.GroupCreateDTO) (*model.Group, error) {
	newRecord := mapper.MapCreateToGroup(dto)

	var permissions []*model.Permission
	if len(dto.PermissionIDs) > 0 {
		if err := s.DB.Where("id IN ?", dto.PermissionIDs).Find(&permissions).Error; err != nil {
			return nil, exceptions.NewBadRequest("error_query_permissions", err)
		}
		if len(permissions) != len(dto.PermissionIDs) {
			return nil, exceptions.NewBadRequest("invalid_permissions", nil)
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

func (s *GroupService) Update(id uint64, dto *dto.GroupUpdateDTO) (*model.Group, error) {
	record, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updatedRecord := mapper.MapUpdateToGroup(dto, record)

	var permissions []*model.Permission
	if len(dto.PermissionIDs) > 0 {
		if err := s.DB.Where("id IN ?", dto.PermissionIDs).Find(&permissions).Error; err != nil {
			return nil, exceptions.NewBadRequest("error_query_permissions", err)
		}
		if len(permissions) != len(dto.PermissionIDs) {
			return nil, exceptions.NewBadRequest("invalid_permissions", nil)
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

func (s *GroupService) PartialUpdate(id uint64, dto *dto.GroupPatchDTO) (*model.Group, error) {
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
			var permissions []*model.Permission
			if len(dto.PermissionIDs) > 0 {
				if err := tx.Where("id IN ?", dto.PermissionIDs).Find(&permissions).Error; err != nil {
					return exceptions.NewBadRequest("error_query_permissions", err)
				}
				if len(permissions) != len(dto.PermissionIDs) {
					return exceptions.NewBadRequest("invalid_permissions", nil)
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
