package controller

import (
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func PostVoteController(c *gin.Context) {
	//参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParams, errData)
	}
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	logic.VoteForPost(userID, p)
	ResponseSuccess(c, nil)
}
