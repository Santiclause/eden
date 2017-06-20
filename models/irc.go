package models

type IrcUser struct {
	ID       uint   `gorm:"primary_key"`
	Nickname string `gorm:"size:60"`
	User     User
	UserID   uint
}
