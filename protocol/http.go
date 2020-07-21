package protocol

// UserReq 登录/注册 请求
type UserReq struct {
	UserName string `form:"UserName" json:"UserName"`
	PassWord string `form:"PassWord" json:"PassWord"`
}

// UserResp ...
type UserResp struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Err  string `json:"error"`
}
