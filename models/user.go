package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Username    string `gorm:"size:60"`
	Password    string `sql:"type:CHAR(60) CHARACTER SET latin1 COLLATE latin1_bin"`
	Roles       []Role `gorm:"many2many:user_roles"`
	Permissions []Permission
}

type Role struct {
	ID          uint         `gorm:"primary_key"`
	Name        string       `gorm:"size:60"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"size:60"`
}

func (user *User) GetPermissions(db *gorm.DB) error {
	return db.Joins(`
	JOIN role_permissions
		ON permissions.id = role_permissions.permission_id
	JOIN user_roles
		USING (role_id)
	`).Find(&user.Permissions, fmt.Sprintf("user_id = %d", user.ID)).Error
}
