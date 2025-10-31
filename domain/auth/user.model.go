package auth

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Password    string `gorm:"size:128;not null"`
	LastLogin   *time.Time
	IsSuperuser bool   `gorm:"default:false"`
	Username    string `gorm:"size:150;uniqueIndex;not null"`
	FirstName   string `gorm:"size:150"`
	LastName    string `gorm:"size:150"`
	Email       string `gorm:"size:254;uniqueIndex;not null"`
	IsStaff     bool   `gorm:"default:false"`
	IsActive    bool   `gorm:"default:true"`

	Groups          []*Group      `gorm:"many2many:auth_user_groups;"`
	UserPermissions []*Permission `gorm:"many2many:auth_user_permissions;"`
}

func (u *User) TableName() string { return "auth_user" }

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
func (u *User) HasPerm(
	db *gorm.DB,
	module,
	action string,
) bool {
	// 1. Superusuários ativos têm todas as permissões
	if u.IsActive && u.IsSuperuser {
		return true
	}
	// 2. Usuários inativos não têm permissões
	if !u.IsActive {
		return false
	}

	// 3. Encontrar o ID da permissão
	var permID uint64
	err := db.Model(&Permission{}).
		Select("id").
		Where("module = ? AND action = ?", module, action).
		First(&permID).Error

	if err != nil {
		// A permissão nem existe
		return false
	}

	// 4. Verificar Permissões Diretas do Usuário
	// (Usa nomes de tabela/coluna M2M padrão do GORM)
	var directPermissionCount int64
	err = db.Table("auth_user_permissions").
		Where("user_id = ? AND permission_id = ?", u.ID, permID).
		Count(&directPermissionCount).Error

	if err == nil && directPermissionCount > 0 {
		return true // Encontrada permissão direta
	}

	// 5. Verificar Permissões de Grupo
	var groupPermissionCount int64
	err = db.Table("auth_user_groups").
		Select("auth_user_groups.group_id").
		// Junta as tabelas de associação
		Joins("INNER JOIN auth_group_permissions ON auth_user_groups.group_id = auth_group_permissions.group_id").
		// Filtra pelo usuário
		Where("auth_user_groups.user_id = ?", u.ID).
		// Filtra pela permissão
		Where("auth_group_permissions.permission_id = ?", permID).
		Count(&groupPermissionCount).Error

	// Se a contagem for > 0, o usuário tem a permissão através de um grupo
	return err == nil && groupPermissionCount > 0
}
