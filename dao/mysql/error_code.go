package mysql

import "errors"

var (
	ErrorUserExist     = errors.New("用户已存在")
	ErrorUserNotExit   = errors.New("用户不存在")
	ErrorPasswordWrong = errors.New("用户名或密码错误")
	ErrorInvalidID     = errors.New("无效的ID")
	ErrorInsertFailed  = errors.New("插入数据失败")
)
