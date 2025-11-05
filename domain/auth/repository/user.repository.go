package repository

import (
	"grf/core/repository"
	"grf/domain/auth/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	repository.IRepository[*model.User, uint64]

	DB *gorm.DB
}

func NewUserRepository(
	db *gorm.DB,
) *UserRepository {
	return &UserRepository{
		IRepository: repository.NewGenericRepository(&repository.Config[*model.User, uint64]{
			DB: db,
			NewModel: func() *model.User {
				return new(model.User)
			},
		}),
		DB: db,
	}
}

func (r *UserRepository) FindUserByEmailOrUsername(login string) (model.User, error) {
	var user model.User
	if err := r.DB.Where("username = ? OR email = ?", login, login).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
