package models

import "time"

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"size:60"`
	Password  string `sql:"type:CHAR(60) CHARACTER SET latin1 COLLATE latin1_bin"`
	Roles     []Role `gorm:"many2many:user_roles"`
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
