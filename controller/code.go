package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 * iota
	CodeInvalidParams
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidToken
	CodeExpiredAccessToken
	CodeNeedLogin
	CodeInvalidUserID
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户名重复",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken:       "无效token",
	CodeExpiredAccessToken: "access token过期",
	CodeNeedLogin:          "需要登录",
	CodeInvalidUserID:      "非法用户ID",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
