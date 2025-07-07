package model

import (
	"gorm.io/gorm"
	"time"
)

// 书籍信息
type BookInfo struct {
	gorm.Model

	// 作品方向 0-男频 1 女频
	WorkDirection int64

	CategoryId int64

	CategoryName string

	PicUrl string

	BookName string

	AuthorId int64

	AuthorName string

	BookDesc string

	Score int64

	// 书籍状态 0-连载中 1-已完结
	BookStatus int64

	VisitCount int64

	WordCount int64

	CommentCount int64

	LastChapterId int64

	LastChapterName string

	LastChapterUpdateTime time.Time

	IsVip int64
}
