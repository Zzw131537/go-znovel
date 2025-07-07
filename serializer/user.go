package serializer

import (
	"go_novel/model"
	"go_novel/pkg/util"
)

type User struct {
	UserName string `json:"user_name"`

	NickName string `json:"nick_name"`

	UserPhoto string `json:"user_photo"`

	UserSex int64 `json:"user_sex"`

	AccountBalance int64 `json:"account_balance"`
}

type User2 struct {
	UserName string `json:"user_name"`

	Password string `json:"password"`

	NickName string `json:"nick_name"`

	UserPhoto string `json:"user_photo"`

	UserSex int64 `json:"user_sex"`

	AccountBalance int64 `json:"account_balance"`

	CreateTime int64 `json:"create_time"`
}

func BuildUser(user *model.UserInfo) *User {
	return &User{
		UserName:       user.UserName,
		NickName:       user.NickName,
		UserPhoto:      user.UserPhoto,
		AccountBalance: user.AccountBalance,
		UserSex:        user.UserSex,
	}
}

func BuildUser2(user *model.UserInfo) *User2 {
	return &User2{
		UserName:       user.UserName,
		Password:       util.Encrypt.AesDecoding(user.Password),
		NickName:       user.NickName,
		UserPhoto:      user.UserPhoto,
		UserSex:        user.UserSex,
		AccountBalance: user.AccountBalance,
	}
}
