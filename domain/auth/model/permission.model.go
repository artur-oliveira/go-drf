package model

import (
	"time"
)

type Permission struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Module      string `gorm:"size:100;not null;uniqueIndex:unique_permission"`
	Action      string `gorm:"size:100;not null;uniqueIndex:unique_permission"`
	Description string `gorm:"size:255;not null"`
}

func (Permission) TableName() string { return "auth_permission" }
