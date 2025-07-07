package model

import "gorm.io/gorm"

type BookChapter struct {
	gorm.Model

	BookId int64

	ChapterNum int64

	ChapterName string

	WordCount int64

	IsVip int64
}
