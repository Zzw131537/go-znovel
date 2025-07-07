package model

import "gorm.io/gorm"

// 用户反馈
type UserFeedback struct {
	gorm.Model

	UserId int64

	Content string
}
