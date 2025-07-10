package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetUserProfileHandler(c *gin.Context) {
	user_id_from_token, _ := getCurrentUserID(c)
	user_id_from_param_str := c.Param("user_id")
	user_id_from_param, _ := strconv.ParseInt(user_id_from_param_str, 10, 64)

	if user_id_from_token != user_id_from_param {
		zap.L().Error("invalid user_id")
		ResponseError(c, CodeInvalidUserID)
		return
	}

	//2.业务处理
	user, err := logic.GetUserProfile(user_id_from_param)
	if err != nil {
		zap.L().Error("logic.GetUserProfile failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, user)
}
func UpdateUserProfileHandler(c *gin.Context) {
	user_id_from_token, _ := getCurrentUserID(c)
	user_id_from_param_str := c.Param("user_id")
	user_id_from_param, _ := strconv.ParseInt(user_id_from_param_str, 10, 64)
	user_profile := new(models.UserProfile)

	c.ShouldBindJSON(user_profile)
	if user_id_from_token != user_id_from_param {
		zap.L().Error("invalid user_id")
		ResponseError(c, CodeInvalidUserID)
		return
	}

	//2.业务处理
	err := logic.UpdateUserProfile(user_profile)
	if err != nil {
		zap.L().Error("logic.UpdateUserProfile failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

func SignUpHandler(c *gin.Context) {
	//1.获取参数参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validator.ValidtaionErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": removeTopStruct(errs.Translate(trans)),
		// })
		return
	}
	fmt.Println(p)
	//2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	//1.获取参数
	p := new(models.ParamLogin)
	err := c.ShouldBindJSON(p)
	if err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.验证账号密码
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("login.Login failed", zap.String("username", p.Username))
		if errors.Is(err, mysql.ErrorUserNotExit) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	c.SetCookie("refresh_token",
		user.RefreshToken,
		viper.GetInt("auth.refresh_token_expire"),
		"/",
		"",
		false,
		true)
	//3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":       fmt.Sprintf("%d", user.UserID),
		"user_name":     user.Username,
		"access_token":  user.AccessToken,
		"refresh_token": user.RefreshToken,
	})
}
