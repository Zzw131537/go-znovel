package front

import (
	"github.com/gin-gonic/gin"
	"go_novel/pkg/e"
	"go_novel/pkg/util"
	"go_novel/service"
	"net/http"
	"strconv"
)

func UserRegister(ctx *gin.Context) {
	var userService service.UserRegisterService

	if err := ctx.ShouldBind(&userService); err == nil {
		res := userService.Register(ctx.Request.Context())
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": e.UserRegisterError,
			"msg":  e.GetMsg(e.UserRegisterError),
			"err":  err.Error(),
		})
	}
}

func UserLogin(ctx *gin.Context) {
	var userService service.UserLoginService
	if err := ctx.ShouldBind(&userService); err == nil {
		res := userService.Login(ctx.Request.Context())
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": e.UserLoginError,
			"msg":  e.GetMsg(e.UserLoginError),
			"err":  err.Error(),
		})
	}

}

// 查询当前用户信息
func UserInfo(ctx *gin.Context) {
	var userService service.UserService
	token := ctx.GetHeader("Authorization")
	claims, _ := util.ParseToken(token)
	if err := ctx.ShouldBind(&userService); err == nil {
		res := userService.UserInfo(claims.ID, ctx.Request.Context())
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": e.UserLoginError,
			"msg":  e.GetMsg(e.UserLoginError),
			"err":  err.Error(),
		})
	}
}

// 更新用户信息
func UserUpdate(context *gin.Context) {
	var userService service.UpdateUserService
	token := context.GetHeader("Authorization")
	claims, _ := util.ParseToken(token)
	if err := context.ShouldBind(&userService); err == nil {
		res := userService.UpdateUserInfo(claims.ID, context.Request.Context())
		context.JSON(http.StatusOK, res)

	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": e.Error,
			"msg":  e.GetMsg(e.Error),
			"err":  err.Error(),
		})
	}

}

// 保存用户反馈
func UserFeedback(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	claims, _ := util.ParseToken(token)

	userService := service.UserFeedbackService{}
	if err := ctx.ShouldBind(&userService); err == nil {
		feedback := userService.UserFeedback(claims.ID, ctx)
		ctx.JSON(http.StatusOK, feedback)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": e.Error,
			"msg":  e.GetMsg(e.Error),
			"err":  err.Error(),
		})
	}

}

// 发表评论接口
func UserPublishComment(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	claims, _ := util.ParseToken(token)

	userService := service.UserCommentService{}
	if err := ctx.ShouldBind(&userService); err == nil {
		res := userService.UserComment(claims.ID, ctx.Request.Context())
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": e.Error,
			"msg":  e.GetMsg(e.Error),
			"err":  err.Error(),
		})
	}
}

// 修改评论接口
func UserUpdateComment(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	claims, _ := util.ParseToken(token)
	userService := service.UserCommentService{}
	if err := ctx.ShouldBind(&userService); err == nil {
		res := userService.UpdateUserComment(claims.ID, ctx.Request.Context())
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": e.Error,
			"msg":  e.GetMsg(e.Error),
			"err":  err.Error(),
		})
	}
}

func UserDeleteComment(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	claims, _ := util.ParseToken(token)
	userService := service.UserCommentService{}
	if err := ctx.ShouldBind(&userService); err == nil {
		res := userService.DeleteUserComment(claims.ID, ctx.Request.Context())
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": e.Error,
			"msg":  e.GetMsg(e.Error),
			"err":  err.Error(),
		})
	}
}

func BookShelfStatus(ctx *gin.Context) {
	value := ctx.Query("book_id")
	atoi, _ := strconv.Atoi(value)
	b := atoi
	a := uint(b)
	if value == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": e.Error,
			"msg":  "book_id 是必需的",
			"err":  e.GetMsg(e.Error),
		})
	}
	token := ctx.GetHeader("Authorization")
	claims, _ := util.ParseToken(token)
	userService := service.UserService{}
	if err := ctx.ShouldBind(&userService); err == nil {
		userService.UserBookShelfStatus(claims.ID, a, ctx.Request.Context())
	}
}
