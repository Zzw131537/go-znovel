package router

import (
	"github.com/gin-gonic/gin"
	api "go_novel/api/front"
	"go_novel/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())

	v1 := r.Group("api/v1")
	{
		// 前台模块
		front := v1.Group("/front")

		user := front.Group("/user")

		// 注册接口
		user.POST("/register", api.UserRegister)

		user.POST("/login", api.UserLogin)

		userJwt := user
		userJwt.Use(middleware.JWT())
		{
			userJwt.GET("", api.UserInfo)
			userJwt.POST("/update", api.UserUpdate)
			userJwt.POST("/feedback", api.UserFeedback)
			userJwt.POST("/publishcomment", api.UserPublishComment)
			userJwt.POST("/updatecomment", api.UserUpdateComment)
			userJwt.POST("/deletecomment", api.UserDeleteComment)
			userJwt.GET("/bookshelf_status", api.BookShelfStatus)
		}

		book := front.Group("/book")
		book.GET("/category/list", api.BookCateGoryList)
		book.GET("", api.BookInfo)
		book.GET("/visit", api.AddVisitCount)
	}
	return r
}
