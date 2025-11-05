package model

import (
	"time"
)

type Group struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string `gorm:"size:50;uniqueIndex;not null"`

	Permissions []*Permission `gorm:"many2many:auth_group_permissions;"`
}

func (Group) TableName() string { return "auth_group" }
