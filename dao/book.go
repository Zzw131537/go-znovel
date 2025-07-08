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

// 根据类别查询书籍
func (dao *BookDao) FindBooksByCategoryId(categoryId int64) ([]model.BookInfo, error) {
	var bookInfos []model.BookInfo
	err := dao.Model(&model.BookInfo{}).Where("category_id = ?", categoryId).Find(&bookInfos).Error
	if err != nil {
		return nil, err
	}
	return bookInfos, nil
}

// 按小说点击量返回书籍信息
func (dao *BookDao) FindBooksByVisit() ([]model.BookInfo, error) {
	var bookInfos []model.BookInfo
	err := dao.Model(&model.BookInfo{}).Order("visit_count desc").Find(&bookInfos).Error
	if err != nil {
		return nil, err
	}
	return bookInfos, nil
}

// 根据创建时间返回数据
func (dao *BookDao) FindBooksByCreateAt() ([]model.BookInfo, error) {
	var bookInfos []model.BookInfo
	err := dao.Model(&model.BookInfo{}).Order("created_at desc").Find(&bookInfos).Error
	if err != nil {
		return nil, err
	}
	return bookInfos, nil
}

func (dao *BookDao) FindBookByAuthorIdAndBookName(author_id int64, bookName string) (model.BookInfo, error) {
	var bookInfo model.BookInfo
	err := dao.Model(&model.BookInfo{}).Where("author_id = ? and book_name = ?", author_id, bookName).First(&bookInfo).Error
	if err != nil {
		return model.BookInfo{}, err
	}
	return bookInfo, nil
}

func (dao *BookDao) SaveBook(book *model.BookInfo) error {
	return dao.Model(book).Create(book).Error
}

func (dao *BookDao) FindBooks(userId int64) (*[]model.BookInfo, error) {
	var bookInfos []model.BookInfo
	err := dao.Model(&model.BookInfo{}).Where("user_id = ?", userId).Find(&bookInfos).Error
	if err != nil {
		return nil, err
	}
	return &bookInfos, nil
}

func (dao *BookDao) SaveBookChapter(bookChapter *model.BookChapter) error {
	return dao.Model(bookChapter).Create(bookChapter).Error
}

func (dao *BookDao) FindChaptersByBookId(bookId int64) (*[]model.BookChapter, error) {
	var bookChapters []model.BookChapter
	err := dao.Model(&model.BookChapter{}).Where("book_id = ?", bookId).Find(&bookChapters).Error
	if err != nil {
		return nil, err
	}
	return &bookChapters, nil

}
