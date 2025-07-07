package e

var MsgFlags = map[int]string{
	Success:                    "ok",
	Error:                      "fail",
	InvalidParams:              "参数错误",
	ErrorAuthToken:             "token错误",
	ErrorAuthCheckTokenTimeout: "token已过期",
	UserRegisterError:          "用户注册失败",
	UserExistError:             "注册的用户已经存在",
	ErrorExistUserNotFound:     "用户不存在",
	UserPasswordError:          "用户密码错误",
	UserTokenError:             "签发token失败",
	UserLoginError:             "用户登录失败",
}

// GetMsg 获取状态码对应的信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[code]
	} else {
		return msg
	}
}
