package dao

import (
	"context"
	"go_novel/model"
	"gorm.io/gorm"
)

type BookDao struct {
	*gorm.DB
}

func NewBookDao(ctx context.Context) *BookDao {
	return &BookDao{NewDBClient(ctx)}
}

func NewBookDaoByDB(db *gorm.DB) *BookDao {
	return &BookDao{db}
}

// 保存用户对书的评论
func (dao *BookDao) SaveUserComment(bookComment *model.BookCommment) error {
	return dao.Model(&bookComment).Create(bookComment).Error
}

// 查看数据库是否存在该评论
func (dao *BookDao) FindUserCommentByUserIdAndBookId(userId, bookId int64) bool {
	var count int64
	count = 0
	err := dao.Model(&model.BookCommment{}).Where("user_id = ? and book_id = ?", userId, bookId).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// 更新数据
func (dao *BookDao) UpdateUserCommentByUserIdAndBookId(val *model.BookCommment) error {
	return dao.Model(&model.BookCommment{}).Where("user_id = ? and book_id = ?", val.UserId, val.BookId).Updates(&val).Error
}

// 删除评论
func (dao *BookDao) DeleteUserCommentByUserIdAndBookId(userId, bookId int64) error {
	return dao.Model(&model.BookCommment{}).Where("user_id = ? and book_id = ?", userId, bookId).Delete(&model.BookCommment{}).Error
}

// 查询小说类别
func (dao *BookDao) FindCategoryByWorkDirection(workDirection int64) ([]*model.BookCategory, error) {
	var categories []*model.BookCategory
	err := dao.Model(&model.BookCategory{}).Where("work_direction = ?", workDirection).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// 根据id 返回小说详细信息
func (dao *BookDao) GetBookInfoById(id int64) (*model.BookInfo, error) {
	var bookInfo model.BookInfo
	err := dao.Model(&model.BookInfo{}).Where("id = ?", id).First(&bookInfo).Error
	if err != nil {
		return nil, err
	}
	return &bookInfo, nil
}
