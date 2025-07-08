package dao

import (
	"fmt"
	"go_novel/model"
)

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		&model.UserInfo{},
		&model.AuthorInfo{},
		&model.BookCategory{},
		&model.BookChapter{},
		&model.BookCommment{},
		&model.UserFeedback{},
		&model.BookContent{},
		&model.UserBookshelf{},
		&model.BookInfo{},
		&model.UserNotification{},
	)
	if err != nil {
		fmt.Println("err ", err)
	}
	return
}
