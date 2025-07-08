package e

const (
	//redis - key
	USER_PREFIX_KEY = "znovel:user:"

	USERINFO_CACHE_EXPIRATION = 0

	BOOK_PREFIX_KEY = "znovel:book:"

	AUTHOR_PREFIX_KEY = "znovel:author:"

	DefaultAesKey = "0123456789012345"

	Success                    = 200
	Error                      = 500
	InvalidParams              = 400
	ErrorAuthToken             = 401
	ErrorAuthCheckTokenTimeout = 402

	UserRegisterError      = 50001
	UserExistError         = 50002
	ErrorExistUserNotFound = 50003
	UserPasswordError      = 50004
	UserTokenError         = 50005
	UserLoginError         = 50006
)
