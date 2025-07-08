package service

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"go_novel/cache"
	"go_novel/dao"
	"go_novel/model"
	"go_novel/mq"
	"go_novel/pkg/e"
	"go_novel/serializer"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type BookService struct{}

type BookVisitService struct {
	db         *gorm.DB
	buffer     map[int64]int64 // 内存缓冲区
	bufferLock sync.Mutex
}

// 初始化服务
func NewBookVisitService(ctx context.Context) *BookVisitService {
	vs := &BookVisitService{
		db:     dao.NewDBClient(ctx),
		buffer: make(map[int64]int64),
	}

	// 启动后台同步服务
	go vs.syncBufferToRedis()
	go vs.syncRedisToDB()

	return vs
}

type BookChapterService struct {
	BookId int64 `json:"book_id" form:"book_id"`

	ChapterNum int64 `json:"chapter_num" form:"chapter_num"`

	ChapterName string `json:"chapter_name" form:"chapter_name"`

	WordCount int64 `json:"word_count" form:"word_count"`

	IsVip int64 `json:"is_vip" form:"is_vip"`
}

// 发布小说章节
func (s *BookChapterService) PublishBookChapter(userId int64, ctx context.Context) serializer.Response {
	bookDao := dao.NewBookDao(ctx)
	chapter := model.BookChapter{
		BookId:      s.BookId,
		ChapterNum:  s.ChapterNum,
		ChapterName: s.ChapterName,
		WordCount:   s.WordCount,
		IsVip:       s.IsVip,
	}
	err := bookDao.SaveBookChapter(&chapter)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  e.GetMsg(e.Error),
			Data: "保存到数据库失败",
		}
	}

	err = mq.PublishChapterUpdate(s.BookId, s.ChapterNum, s.ChapterName)
	if err != nil {
		log.Printf("faild to publish chapter update: %v", err)
	}
	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: "发布小说章节成功",
	}
}

// 增加小说点击量
func (s *BookVisitService) AddVisitCount(id int64, ctx context.Context) serializer.Response {
	// 现在内存缓冲区累加
	s.bufferLock.Lock()
	s.buffer[id] += 1
	s.bufferLock.Unlock()

	// 异步更新到redis,后台任务处理
	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: "点击量增加成功",
	}
}

// 定时将内存缓冲区内的数据同步到缓存
func (s *BookVisitService) syncBufferToRedis() {
	ticker := time.NewTicker(5 * time.Second) // 每5秒同步一次
	defer ticker.Stop()
	for range ticker.C {
		s.bufferLock.Lock()
		if len(s.buffer) == 0 {
			s.bufferLock.Unlock()
			continue
		}

		// 复制当前缓冲区
		localBuffer := make(map[int64]int64)
		for k, v := range s.buffer {
			localBuffer[k] = v
		}
		// 清空缓冲区
		s.buffer = make(map[int64]int64)
		s.bufferLock.Unlock()

		// 批量更新到redis
		pipe := cache.RedisClient.Pipeline()
		key := e.BOOK_PREFIX_KEY + "book_visit"
		for k, v := range localBuffer {
			//fmt.Println(k, v, "-->v1")
			pipe.HIncrBy(key, strconv.Itoa(int(k)), v)
		}
		if _, err := pipe.Exec(); err != nil {
			// 更新失败，将数据加回缓冲区
			s.bufferLock.Lock()
			for k, v := range localBuffer {
				//fmt.Println(k, v, "---->v2")
				s.buffer[k] = v
			}
			s.bufferLock.Unlock()
		}
	}
}

// 将数据从redis 加载进数据库
func (s *BookVisitService) syncRedisToDB() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	key := e.BOOK_PREFIX_KEY + "book_visit"

	for range ticker.C {
		result, err := cache.RedisClient.HGetAll(key).Result()
		//fmt.Println(result, "-->v4")
		if err != nil || len(result) == 0 {
			continue
		}

		// 准备批量更新
		type updateData struct {
			BookId int64
			Count  int64
		}
		var updates []updateData
		for k, v := range result {
			//fmt.Println(k, v, "--->v7")
			bookId, _ := strconv.ParseInt(k, 10, 64)
			count, _ := strconv.ParseInt(v, 10, 64)

			if bookId > 0 && count > 0 {
				updates = append(updates, updateData{
					BookId: bookId,
					Count:  count,
				})
			}
		}

		//fmt.Println(len(updates), "-->v10")
		// 数据库批量更新
		if len(updates) > 0 {
			//fmt.Println("THis 1")
			s.db.Transaction(func(tx *gorm.DB) error {
				//fmt.Println("this 2")
				for _, data := range updates {
					//fmt.Println(data.BookId, data.Count, "---->v3")
					if err := tx.Model(&model.BookInfo{}).Where("id = ?", data.BookId).Update("visit_count", gorm.Expr("visit_count + ?", data.Count)).Error; err != nil {
						continue
					}
				}
				return nil
			})
		}

		//清空redis
		cache.RedisClient.Del(key)
	}
}

func (s *BookService) BookCateGoryList(workDirection int64, ctx context.Context) serializer.Response {
	// 查看缓存中是否存在数据
	key := e.BOOK_PREFIX_KEY + "workDirection:" + strconv.Itoa(int(workDirection))
	redisClient := cache.RedisClient
	result, err := redisClient.Get(key).Result()
	if err == nil && result != "" {
		var cateList []model.BookCategory
		json.Unmarshal([]byte(result), &cateList)
		return serializer.Response{
			Code: e.Success,
			Msg:  e.GetMsg(e.Success),
			Data: cateList,
		}
	}

	// 查询数据库
	bookDao := dao.NewBookDao(ctx)

	cates, err := bookDao.FindCategoryByWorkDirection(workDirection)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  e.GetMsg(e.Error),
			Data: "查询数据库出错",
		}
	}

	marshal, _ := json.Marshal(cates)
	redisClient.Set(key, string(marshal), time.Hour*24)
	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: cates,
	}
}

func (s *BookService) GetBookInfoById(id int64, ctx context.Context) serializer.Response {
	// 查询缓存是否存在
	key := e.BOOK_PREFIX_KEY + "id:" + strconv.Itoa(int(id))
	redisClient := cache.RedisClient
	result, err := redisClient.Get(key).Result()
	if err == nil && result != "" {
		var bookInfo model.BookInfo
		json.Unmarshal([]byte(result), &bookInfo)
		return serializer.Response{
			Code: e.Success,
			Msg:  e.GetMsg(e.Success),
			Data: bookInfo,
		}
	}
	bookDao := dao.NewBookDao(ctx)
	bookInfo, err := bookDao.GetBookInfoById(id)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  e.GetMsg(e.Error),
			Data: "查询数据出错",
		}
	}
	marshal, _ := json.Marshal(bookInfo)
	redisClient.Set(key, string(marshal), time.Hour*24)
	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: bookInfo,
	}
}

// 返回推荐小说列表
func (s *BookService) BookRecList(book_id int64, ctx context.Context) serializer.Response {
	bookDao := dao.NewBookDao(ctx)
	res, err := bookDao.GetBookInfoById(book_id)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  e.GetMsg(e.Error),
			Data: "查询数据库出错",
		}
	}
	category := res.CategoryId
	// 先看缓存中有没有
	key := e.BOOK_PREFIX_KEY + "book_rec:" + strconv.Itoa(int(category))
	redisClient := cache.RedisClient
	result, err := redisClient.Get(key).Result()
	var books []model.BookInfo
	var ok bool
	ok = true
	if err == nil && result != "" {
		err := json.Unmarshal([]byte(result), &books)
		if err != nil {
			ok = false
		}
	} else {
		ok = false
	}
	if !ok {
		books, err := bookDao.FindBooksByCategoryId(category)
		if err != nil {
			return serializer.Response{
				Code: e.Error,
				Msg:  e.GetMsg(e.Error),
				Data: "数据库查询失败",
			}
		}
		marshal, _ := json.Marshal(books)
		redisClient.Set(key, string(marshal), time.Hour*24)
	}

	if len(books) > 0 {
		res := make([]model.BookInfo, len(books))
		vis := make([]bool, len(books)+10)
		var count int
		count = 0
		rand.Seed(time.Now().UnixNano())
		for count < 3 {
			idx := rand.Intn(len(books))
			if vis[idx] == false {
				res = append(res, books[idx])
				vis[idx] = true
				count++
			}
		}
		return serializer.Response{
			Code: e.Success,
			Msg:  e.GetMsg(e.Success),
			Data: res,
		}
	} else {
		return serializer.Response{
			Code: e.Success,
			Msg:  e.GetMsg(e.Success),
			Data: "",
		}
	}
}

// 根据点击量查询小说排行榜
func (s *BookService) BookVisitRank(ctx context.Context) serializer.Response {
	key := e.BOOK_PREFIX_KEY + "visit_rank"
	result, err := cache.RedisClient.ZRevRangeWithScores(key, 0, -1).Result()
	var books []model.BookInfo
	if err == nil && result != nil {

		for _, v := range result {
			book := model.BookInfo{}
			err := json.Unmarshal([]byte(v.Member.(string)), &book)
			if err == nil {
				books = append(books, book)
			}
		}
	}
	if len(books) == 0 {
		// 从数据库中
		bookDao := dao.NewBookDao(ctx)
		books, err = bookDao.FindBooksByVisit()
		if err != nil {
			return serializer.Response{
				Code: e.Error,
				Msg:  e.GetMsg(e.Error),
				Data: "查询数据库错误",
			}
		}
		// 使用管道批量写入redis
		pip := cache.RedisClient.Pipeline()
		for _, book := range books {
			marshal, _ := json.Marshal(book)
			pip.ZAdd(key, redis.Z{
				Score:  float64(book.VisitCount),
				Member: marshal,
			})
		}
		pip.Expire(key, time.Hour*24)
		pip.Exec()
	}
	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: books,
	}
}

// 根据创建时间查询小说排行榜
func (s *BookService) BookNewRank(ctx context.Context) serializer.Response {
	key := e.BOOK_PREFIX_KEY + "newbook_rank"
	result, err := cache.RedisClient.ZRevRangeWithScores(key, 0, -1).Result()
	var books []model.BookInfo
	if err == nil && result != nil {

		for _, v := range result {
			book := model.BookInfo{}
			err := json.Unmarshal([]byte(v.Member.(string)), &book)
			if err == nil {
				books = append(books, book)
			}
		}
	}
	if len(books) == 0 {
		// 从数据库中
		bookDao := dao.NewBookDao(ctx)
		books, err = bookDao.FindBooksByCreateAt()
		if err != nil {
			return serializer.Response{
				Code: e.Error,
				Msg:  e.GetMsg(e.Error),
				Data: "查询数据库错误",
			}
		}
		// 使用管道批量写入redis
		pip := cache.RedisClient.Pipeline()
		for _, book := range books {
			marshal, _ := json.Marshal(book)
			pip.ZAdd(key, redis.Z{
				Score:  float64(book.VisitCount),
				Member: marshal,
			})
		}
		pip.Expire(key, time.Hour*24)
		pip.Exec()
	}
	return serializer.Response{
		Code: e.Success,
		Msg:  e.GetMsg(e.Success),
		Data: books,
	}
}

func (s *BookService) GetChapterList(bookId int64, ctx context.Context) serializer.Response {
	key := e.BOOK_PREFIX_KEY + "chapter_list:" + strconv.Itoa(int(bookId))
	result, err := cache.RedisClient.Get(key).Result()
	if err == nil && result != "" {
		var chapters []model.BookChapter
		json.Unmarshal([]byte(result), &chapters)
		return serializer.Response{
			Code: e.Success,
			Data: chapters,
			Msg:  e.GetMsg(e.Success),
		}
	}
	bookDao := dao.NewBookDao(ctx)
	res, err := bookDao.FindChaptersByBookId(bookId)
	if err != nil {
		return serializer.Response{
			Code: e.Error,
			Msg:  e.GetMsg(e.Error),
			Data: "查询数据库失败",
		}
	}
	marshal, _ := json.Marshal(res)
	cache.RedisClient.Set(key, string(marshal), time.Hour*24)
	return serializer.Response{
		Code: e.Success,
		Data: res,
		Msg:  e.GetMsg(e.Success),
	}
}
