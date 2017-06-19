package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:60"`
	Password string `sql:"type:CHAR(60) CHARACTER SET latin1 COLLATE latin1_bin"`
	Roles    []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	ID          uint
	Name        string       `gorm:"size:60"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	ID   uint
	Name string `gorm:"size:60"`
}
