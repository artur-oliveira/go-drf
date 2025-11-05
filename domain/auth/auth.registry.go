package auth

import "grf/domain/auth/model"

func GetModels() []interface{} {
	return []interface{}{
		&model.Permission{},
		&model.Group{},
		&model.User{},
	}
}
