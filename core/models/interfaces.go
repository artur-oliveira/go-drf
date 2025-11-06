package models

import "gorm.io/gorm"

type IUser interface {
	Active() bool

	Admin() bool

	HasPerm(db *gorm.DB, module string, action string) bool
}
