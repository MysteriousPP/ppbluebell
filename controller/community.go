package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区相关
func CommunityHandler(c *gin.Context) {
	//查询到所有的社区（community_id, community_name）以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	ResponseSuccess(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	// 1.获取社区id
	idStr := c.Param("id")
	//参数校验
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
