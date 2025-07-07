package service

import (
	"context"
	"encoding/json"
	"go_novel/cache"
	"go_novel/dao"
	"go_novel/model"
	"go_novel/pkg/e"
	"go_novel/pkg/util"
	"go_novel/serializer"
	"strconv"

	"regexp"
)

type UserService struct{}
type UserRegisterService struct {
	UserName        string `json:"user_name" form:"user_name"`
	Password        string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}

type UserLoginService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

type UpdateUserService struct {
	NickName string `json:"nick_name" form:"nick_name"`

	UserPhoto string `json:"user_photo" form:"user_photo"`

	UserSex int64 `json:"user_sex" form:"user_sex"`

	Password string `json:"password" form:"password"`
}

type UserFeedbackService struct {
	Content string `json:"content" form:"content"`
}

type UserCommentService struct {
	BookId         int64  `json:"book_id" form:"book_id"`
	CommentContent string `json:"comment_content" form:"comment_content"`
}

// 用户注册
func (s *UserRegisterService) Register(ctx context.Context) serializer.Response {
	// 用户名就是手机号，唯一

	// 1.判断用户名是否符合要求
	matched, _ := regexp.Match("^1[3|4|5|6|7|8|9][0-9]{9}$", []byte(s.UserName))

	if !matched {
		return serializer.Response{
			Code: e.Error,
			Msg:  "用户名不符合电话号码要求",
		}
	}

	// 2.  判断两次输入的密码是否正确
	if s.Password != s.ConfirmPassword {
		return serializer.Response{
			Code: e.Error,
			Msg:  "两次输入的密码不同",
		}
	}

	//3. 判断该用户是否已经注册过
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(s.UserName)

	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  "用户注册出现错误",
		}
	}
	if exist {
		return serializer.Response{
			Code: e.UserExistError,
			Msg:  e.GetMsg(e.UserRegisterError),
		}
	}
	encoding := util.Encrypt.AesEncoding(s.Password)
	user := &model.UserInfo{
		UserName:       s.UserName,
		Password:       encoding,
		Salt:           "0",
		NickName:       s.UserName,
		UserPhoto:      "",
		UserSex:        0,
		AccountBalance: 0,
		Status:         0,
	}

	// 保存注册信息
	err = userDao.SaveUserInfo(user)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  "保存注册信息失败",
		}
	}

	return serializer.Response{
		Code: e.Success,
		Msg:  "注册成功",
		Data: user,
	}
}

// 用户登录
func (s *UserLoginService) Login(ctx context.Context) serializer.Response {
	// 1. 判断用户是否已经注册
	dao := dao.NewUserDao(ctx)
	user, exist, err := dao.ExistOrNotByUserName(s.UserName)
	if !exist || err != nil {
		return serializer.Response{
			Code: e.ErrorExistUserNotFound,
			Msg:  e.GetMsg(e.ErrorExistUserNotFound),
			Data: "用户不存在，请先注册",
		}
	}

	// 校验密码
	if s.Password != util.Encrypt.AesDecoding(user.Password) {
		return serializer.Response{
			Code: e.UserPasswordError,
			Msg:  e.GetMsg(e.UserPasswordError),
			Data: "密码错误，请重新登录",
		}
	}

	// token 签发
	token, err := util.GenerateToken(user.ID, user.UserName, 0)
	if err != nil {
		return serializer.Response{
			Code: e.UserTokenError,
		}
	}

	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: serializer.TokenData{
			User:  serializer.BuildUser(user), // 序列化同时进行数据脱敏
			Token: token,
		},
	}

}

// 返回当前用户信息
func (s *UserService) UserInfo(userId uint, ctx context.Context) serializer.Response {
	redisClient := cache.RedisClient
	//1. 判断缓存中是否存在
	key := e.USER_PREFIX_KEY + strconv.Itoa(int(userId))
	result, err2 := redisClient.Get(key).Result()

	if err2 == nil && result != "" {
		// 反序列化成对象并返回
		var user model.UserInfo
		if err3 := json.Unmarshal([]byte(result), &user); err3 == nil {
			return serializer.Response{
				Code: e.Success,
				Msg:  e.GetMsg(e.Success),
				Data: serializer.BuildUser2(&user),
			}
		}
	}

	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserInfoById(userId)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  "查询用户信息错误",
		}
	}

	// 将数据写入缓存
	marshal, _ := json.Marshal(user)
	redisClient.Set(key, marshal, e.USERINFO_CACHE_EXPIRATION)
	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: serializer.BuildUser2(user),
	}
}

// 修改用户信息
func (s *UpdateUserService) UpdateUserInfo(userId uint, ctx context.Context) serializer.Response {
	// 先更新数据库，再删除缓存
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserInfoById(userId)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  "查询数据失败",
		}
	}
	if s.NickName != "" {
		user.NickName = s.NickName
	}
	if s.UserPhoto != "" {
		user.UserPhoto = s.UserPhoto
	}
	if s.Password != "" {
		user.Password = util.Encrypt.AesDecoding(s.Password)
	}
	if s.UserSex == 0 || s.UserSex == 1 {
		user.UserSex = s.UserSex
	}

	// 保存数据到数据库
	err = userDao.UpdateUserInfo(user)
	if err != nil {
		return serializer.Response{
			Code:  e.Error,
			Msg:   "保存到数据库失败",
			Error: err.Error(),
		}
	}

	// 删除缓存
	redisClient := cache.RedisClient
	key := e.USER_PREFIX_KEY + strconv.Itoa(int(userId))
	_, err = redisClient.Del(key).Result()
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  "删除用户信息缓存失败",
		}
	}
	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: "更新用户信息成功",
	}
}

// 用户提交反馈
func (s *UserFeedbackService) UserFeedback(userId uint, ctx context.Context) serializer.Response {
	userDao := dao.NewUserDao(ctx)
	userFeedback := model.UserFeedback{
		UserId:  int64(userId),
		Content: s.Content,
	}
	err := userDao.SaveUserFeedback(&userFeedback)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  "保存用户反馈失败",
		}
	}
	return serializer.Response{
		Code: e.Success,
		Msg:  "保存用户反馈成功",
		Data: userFeedback,
	}
}

// 一个用户对一本书发表一个评论
func (s *UserCommentService) UserComment(userId uint, ctx context.Context) serializer.Response {
	redisClient := cache.RedisClient
	key := e.USER_PREFIX_KEY + "comment:" + strconv.Itoa(int(userId)) + "bookId-" + strconv.Itoa(int(s.BookId))
	result, err2 := redisClient.Get(key).Result()
	if err2 == nil && result != "" {
		// 这本书该用户之前已经评价过了
		return serializer.Response{
			Code: e.Error,
			Msg:  e.GetMsg(e.Error),
			Data: "一个用户对一本书只能评价一次",
		}
	}
	// 保存该用户评论
	bookComment := model.BookCommment{
		BookId:         s.BookId,
		UserId:         int64(userId),
		CommentContent: s.CommentContent,
		ReplyCount:     0,
		AuditStatus:    0,
	}
	dao := dao.NewBookDao(ctx)
	err2 = dao.SaveUserComment(&bookComment)
	if err2 != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  "保存用户评论失败",
		}
	}
	marshal, _ := json.Marshal(bookComment)
	redisClient.Set(key, marshal, e.USERINFO_CACHE_EXPIRATION)
	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: "保存用户评论成功",
	}
}

func (s *UserCommentService) UpdateUserComment(userId uint, ctx context.Context) serializer.Response {
	key := e.USER_PREFIX_KEY + "comment:" + strconv.Itoa(int(userId)) + "bookId-" + strconv.Itoa(int(s.BookId))
	// 查询数据库是否存在数据
	dao := dao.NewBookDao(ctx)

	ok := dao.FindUserCommentByUserIdAndBookId(int64(userId), s.BookId)
	if !ok {
		return serializer.Response{
			Code: e.Error,
			Msg:  "数据库中没有要修改的书籍评论",
		}
	}
	v := model.BookCommment{
		BookId:         s.BookId,
		UserId:         int64(userId),
		CommentContent: s.CommentContent,
		ReplyCount:     0,
		AuditStatus:    0,
	}
	err := dao.UpdateUserCommentByUserIdAndBookId(&v)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  "数据更新到数据库失败",
		}
	}

	// 更新缓存
	marshal, _ := json.Marshal(v)
	redisCient := cache.RedisClient
	redisCient.Set(key, marshal, e.USERINFO_CACHE_EXPIRATION)

	return serializer.Response{
		Code: e.Success,
		Msg:  "更新数据成功",
		Data: v,
	}
}

// 删除用户评论
func (s *UserCommentService) DeleteUserComment(userId uint, ctx context.Context) serializer.Response {

	key := e.USER_PREFIX_KEY + "comment:" + strconv.Itoa(int(userId)) + "bookId-" + strconv.Itoa(int(s.BookId))
	// 查询数据库是否存在数据
	dao := dao.NewBookDao(ctx)

	ok := dao.FindUserCommentByUserIdAndBookId(int64(userId), s.BookId)
	if !ok {
		return serializer.Response{
			Code: e.Error,
			Msg:  "数据库中没有要删除的书籍评论",
		}
	}

	// 删除数据库数据
	err := dao.DeleteUserCommentByUserIdAndBookId(int64(userId), s.BookId)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  "删除数据库·数据·失败",
		}
	}

	// 删除缓存
	redisCient := cache.RedisClient
	redisCient.Del(key)
	return serializer.Response{
		Code: e.Success,
		Msg:  "删除数据库数据成功",
	}
}

func (s *UserService) UserBookShelfStatus(userId, bookId uint, ctx context.Context) serializer.Response {
	dao := dao.NewUserDao(ctx)
	status := dao.GetUserBookShelfStatus(userId, bookId)
	var c int
	if status == false {
		c = 0
	} else {
		c = 1
	}
	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: c,
	}
}
