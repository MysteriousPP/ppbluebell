package models

import "time"

type Post struct {
	PostID      int64     `json:"post_id,string" db:"post_id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	AuthorId    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"-" db:"create_time"`
}

type ApiPostDetail struct {
	AuthorName string `json:"author_name"`
	VoteNum    int64  `json:"vote_num"`
	*Post
	*CommunityDetail `json:"community"`
}

type Comment struct {
	CommentID  int64 `json:"comment_id,string" db:"comment_id"`
	PostID     int64 `json:"post_id,string" db:"post_id" binding:"required"`
	FromID     int64 `json:"from_id,string" db:"from_id" binding:"required"`
	ToID       int64 `json:"to_id,string" db:"to_id"`
	Status     int32
	Content    string    `json:"content" db:"content" binding:"required"`
	FromName   string    `json:"from_name"`
	ToName     string    `json:"to_name"`
	FromAvatar string    `json:"from_avatar"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}
