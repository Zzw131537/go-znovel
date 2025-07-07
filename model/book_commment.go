package model

import "gorm.io/gorm"

type BookCommment struct {
	gorm.Model

	BookId int64

	UserId int64

	CommentContent string

	ReplyCount int64

	// 审核状态 0-待审核 1-审核通过 2-审核不通过
	AuditStatus int64
}
