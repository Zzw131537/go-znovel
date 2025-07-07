package dao

import (
	"context"
	"fmt"
	"go_novel/model"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

func NewUserDaoByDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// 判断用户是否存在
func (dao *UserDao) ExistOrNotByUserName(name string) (user *model.UserInfo, exist bool, err error) {
	var count int64
	err = dao.Model(&model.UserInfo{}).Where("user_name = ?", name).First(&user).Count(&count).Error
	fmt.Println(count, err)
	if count == 0 && err != nil {
		return nil, false, nil
	}
	return user, true, nil
}

// 保存用户信息
func (dao *UserDao) SaveUserInfo(user *model.UserInfo) error {
	return dao.Model(&model.UserInfo{}).Save(user).Error
}

// 更新用户信息
func (dao *UserDao) UpdateUserInfo(user *model.UserInfo) error {
	return dao.Model(&model.UserInfo{}).Where("id = ?", user.ID).Updates(user).Error
}

// 根据id 返回用户信息
func (dao *UserDao) GetUserInfoById(id uint) (user *model.UserInfo, err error) {
	err = dao.Model(&model.UserInfo{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 保存用户反馈
func (dao *UserDao) SaveUserFeedback(userFeedback *model.UserFeedback) error {
	return dao.Model(&model.UserFeedback{}).Save(userFeedback).Error
}

// 查询书籍是否在该用户的书架中 0-不在 1-在
func (dao *UserDao) GetUserBookShelfStatus(userId, bookId uint) bool {
	var count int64
	count = 0
	err := dao.Model(&model.UserBookshelf{}).Where("user_id = ? AND book_id = ?", userId, bookId).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
