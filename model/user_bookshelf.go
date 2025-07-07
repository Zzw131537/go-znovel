package model

import "gorm.io/gorm"

type UserBookshelf struct {
	gorm.Model

	UserId int64

	BookId int64

	PreContentId int64
}
