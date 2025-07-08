package router

import (
	"github.com/gin-gonic/gin"
	api "go_novel/api"
	api2 "go_novel/api"
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
		user.POST("/register", api2.UserRegister)

		user.POST("/login", api2.UserLogin)

		userJwt := user
		userJwt.Use(middleware.JWT())
		{
			userJwt.GET("", api2.UserInfo)
			userJwt.POST("/update", api2.UserUpdate)
			userJwt.POST("/feedback", api2.UserFeedback)
			userJwt.POST("/publishcomment", api2.UserPublishComment)
			userJwt.POST("/updatecomment", api2.UserUpdateComment)
			userJwt.POST("/deletecomment", api2.UserDeleteComment)
			userJwt.GET("/bookshelf_status", api2.BookShelfStatus)
		}

		book := front.Group("/book")
		book.GET("/category/list", api2.BookCateGoryList)
		book.GET("", api2.BookInfo)
		book.GET("/visit", api2.AddVisitCount)
		book.GET("/rec_list", api2.BookRecList)
		book.GET("/visit_rank", api2.VisitRank)
		book.GET("/newbook_rank", api2.NewBookRank)

		author := front.Group("/author")
		author.GET("/book/list_chapter", api.ListChapters)
		author.Use(middleware.JWT())
		{
			// 作家注册
			author.POST("/register", api.AuthorRegister)
			author.GET("/status", api.AuthorStatus)
			author.POST("/book", api.PublishBook)
			author.GET("/books", api.ListBooks)
			author.POST("/book/chapter", api.PublishBookChapter)
		}
	}
	return r
}
