package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

const CtxUserIDKey = "userID"

func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func getPageInfo(c *gin.Context) (int64, int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var (
		page int64
		size int64
		err  error
	)
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
