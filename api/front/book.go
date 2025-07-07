package front

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_novel/pkg/e"
	"go_novel/service"
	"net/http"
	"strconv"
)

func BookCateGoryList(ctx *gin.Context) {
	workDirection := ctx.Query("workDirection")

	atoi, _ := strconv.Atoi(workDirection)
	c := int64(atoi)
	bookService := service.BookService{}

	if err := ctx.ShouldBind(&bookService); err == nil {
		list := bookService.BookCateGoryList(c, ctx.Request.Context())
		ctx.JSON(http.StatusOK, list)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  e.Error,
			"Msg":   e.GetMsg(e.Error),
			"Error": err.Error(),
		})

	}
}

func BookInfo(ctx *gin.Context) {
	bookId := ctx.Query("book_id")
	atoi, _ := strconv.Atoi(bookId)
	c := int64(atoi)
	bookService := service.BookService{}
	if err := ctx.ShouldBind(&bookService); err == nil {
		res := bookService.GetBookInfoById(c, ctx)
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  e.Error,
			"Msg":   e.GetMsg(e.Error),
			"Error": err.Error(),
		})
	}
}

// 增加小说点击量
func AddVisitCount(ctx *gin.Context) {
	book_id := ctx.Query("book_id")
	if book_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  e.Error,
			"Msg":   e.GetMsg(e.Error),
			"Error": "book_id should not be null",
		})
		return
	}
	bookID, err := strconv.ParseUint(book_id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  e.Error,
			"Msg":   e.GetMsg(e.Error),
			"Error": err.Error(),
		})
		return
	}
	// bookService := service.BookVisitService{}
	c := context.Background()
	ba := service.NewBookVisitService(c)

	if err := ctx.ShouldBind(&ba); err == nil {
		res := ba.AddVisitCount(int64(bookID), ctx.Request.Context())
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":  e.Error,
			"Msg":   e.GetMsg(e.Error),
			"Error": err.Error(),
		})
	}
}
