package auth

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
	Username  string       `json:"username"`
	Password  string       `json:"password"`
	Email     string       `json:"email"`
	IsActive  bool         `json:"is_active"`
	IsAdmin   bool         `json:"is_admin"`
}

func (User) TableName() string {
	return "auth_user"
}
