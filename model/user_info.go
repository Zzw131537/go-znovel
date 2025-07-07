package model

import (
	"gorm.io/gorm"
)

// 用户信息
type UserInfo struct {
	gorm.Model

	UserName string

	Password string

	Salt string

	NickName string

	UserPhoto string

	UserSex int64

	AccountBalance int64

	Status int64
}
