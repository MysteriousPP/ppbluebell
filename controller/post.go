package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeleteCommentHandler 删除评论
func DeleteCommentHandler(c *gin.Context) {
	param := struct {
		CommentID string `json:"comment_id"`
		UserID    string `json:"user_id"`
	}{}

	if err := c.ShouldBindJSON(&param); err != nil {
		zap.L().Error("delete comment with invalid", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	CommentID, _ := strconv.ParseInt(param.CommentID, 10, 64)
	UserID, _ := strconv.ParseInt(param.UserID, 10, 64)
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	if userID != UserID {
		zap.L().Error("Invalid userID", zap.Error(err))
		zap.L().Error("Invalid userID", zap.Error(err), zap.Int64("userID:", userID))
		zap.L().Error("Invalid userID", zap.Error(err), zap.Int64("UserID:", UserID))
		ResponseError(c, CodeInvalidUserID)
		return
	}
	//
	if err := logic.DeleteCommentById(CommentID, userID); err != nil {
		zap.L().Error("logic.DeleteCommentByID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// GetCommentListHandler 获取帖子评论的处理函数
func GetCommentListHandler(c *gin.Context) {
	// 1.获取参数(从url中获取帖子的id)
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get comment with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 2.根据id去查帖子数据（查数据库）
	data, err := logic.GetCommentListById(pid)
	if err != nil {
		zap.L().Error("logic.GetCommentById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, data)
}
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数的校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	//从 c 取到当前发请求的用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorId = userID
	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// CreateCommentHandler 创建评论的处理函数
func CreateCommentHandler(c *gin.Context) {
	Comment := new(models.Comment)
	if err := c.ShouldBindJSON(Comment); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	Comment.FromID = userID

	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("create comment with invalid postid", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	Comment.PostID = pid

	if err = logic.CreateCommentInPost(Comment); err != nil {
		zap.L().Error("logic.CreateCommentInPost(C) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)

}

// DeletePostHandler 删除帖子的处理函数
func DeletePostHandler(c *gin.Context) {
	// 1.获取参数(从url中获取帖子id)
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("delete post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 2.根据id去更改帖子数据
	err = logic.DeletePostById(pid)
	if err != nil {
		zap.L().Error("logic.DeletePostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回相应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {
	// 1.获取参数(从url中获取帖子的id)
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 2.根据id去查帖子数据（查数据库）
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {

	page, size := getPageInfo(c)
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// 1.获取参数
// 2.去redis查询id列表
// 3.根据id去数据库查询帖子详情信息
func GetPostListHandler2(c *gin.Context) {
	//GET请求参数(query string)： /api/v1/post2?page=1&size=10&order=time
	//获取分页参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, //magic string
	}
	//c.ShouldBind() 根据请求的数据类型选择相应的方法区获取数据
	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	data, err := logic.GetPostListNew(p) //更新：合二为一
	//获取数据
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// func GetCommunityPostListHandler(c *gin.Context) {
// 	//GET请求参数(query string)： /api/v1/post2?page=1&size=10&order=time
// 	//获取分页参数
// 	p := &models.ParamCommunityPostList{
// 		ParamPostList: &models.ParamPostList{
// 			Page:  1,
// 			Size:  10,
// 			Order: models.OrderTime,
// 		},
// 	}
// 	//c.ShouldBind() 根据请求的数据类型选择相应的方法区获取数据
// 	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
// 	if err := c.ShouldBindQuery(p); err != nil {
// 		zap.L().Error("GetCommunityPostListHandler with invalid param", zap.Error(err))
// 		ResponseError(c, CodeInvalidParams)
// 		return
// 	}

// 	//获取数据
// 	data, err := logic.GetCommunityPostList(p)
// 	if err != nil {
// 		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
// 		ResponseError(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuccess(c, data)
// }
