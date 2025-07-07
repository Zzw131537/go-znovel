package serializer

type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
