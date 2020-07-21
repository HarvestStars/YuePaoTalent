package protocol

import "github.com/jinzhu/gorm"

// UserInfo 账户信息
type UserInfo struct {
	gorm.Model
	UserName string
	PassWord string
	Salt     string
}
