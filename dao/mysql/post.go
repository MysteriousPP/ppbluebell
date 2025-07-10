package mysql

import (
	"bluebell/models"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func DeleteCommentById(commentID, userID int64) (err error) {
	sqlStr := `update comment
	set status = 0
	where comment_id = ? and from_id = ? 
	`
	// 执行更新
	result, err := db.Exec(sqlStr, commentID, userID)
	if err != nil {
		zap.L().Error("delete comment failed", zap.Error(err))
		return
	}

	// 检查是否有行受影响
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Error("failed to get rows affected", zap.Error(err))
		return
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func GetCommentListById(pid int64) (data []*models.Comment, err error) {
	sqlStr := `select comment_id, post_id, from_id, to_id, content, create_time
	from comment
	where status = 1 and post_id = ?
	order by create_time
	desc
	`
	data = make([]*models.Comment, 1)
	err = db.Select(&data, sqlStr, pid)
	return
}
func CreateComment(C *models.Comment) (err error) {
	sqlStr := `insert into comment(
	comment_id, post_id, from_id, to_id, content)
	values(?,?,?,?,?)`

	_, err = db.Exec(sqlStr, C.CommentID, C.PostID, C.FromID, C.ToID, C.Content)
	if err != nil {
		zap.L().Error("insert comment failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
	post_id, title, content, author_id, community_id)
	values(?,?,?,?,?)`

	_, err = db.Exec(sqlStr, p.PostID, p.Title, p.Content, p.AuthorId, p.CommunityID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

// DeletePostById 根据id删除某个帖子
func DeletePostById(pid int64) (err error) {
	sqlStr := `update post
	set status = 0
	where post_id = ?
	`
	// 执行更新
	result, err := db.Exec(sqlStr, pid)
	if err != nil {
		zap.L().Error("delete post failed", zap.Error(err))
		return
	}

	// 检查是否有行受影响
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Error("failed to get rows affected", zap.Error(err))
		return
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetPostById 根据id查询单个帖子数据
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id = ?
	`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where status = 1
	order by create_time
	desc
	limit ?,?
	`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)`

	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)

	err = db.Select(&postList, query, args...) //！！！！！ ...是什么意思？

	return
}
