package model

import "gorm.io/gorm"

type BookContent struct {
	gorm.Model

	ChapterId int64

	Content string
}
