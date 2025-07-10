package models

import (
	"database/sql"
	"time"
)

const (
	admin  string = "admin"
	normal string = "normal"
)

type User struct {
	UserID       int64          `db:"user_id"`
	Username     string         `db:"username"`
	Password     string         `db:"password"`
	Nickname     sql.NullString `db:"nickname"`
	UserType     string         `db:"usertype"`
	Profile      sql.NullString `db:"profile"`
	Email        sql.NullString `db:"email"`
	Phone        sql.NullString `db:"phone"`
	Avatar       sql.NullString `db:"avatar"`
	AccessToken  string
	RefreshToken string
	CreateTime   time.Time `db:"create_time"`
	UpdateTime   time.Time `db:"update_time"`
}
type UserProfile struct {
	UserID   int64  `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
	Nickname string `db:"nickname" json:"nickname"`
	UserType string `db:"usertype" json:"usertype"`
	Profile  string `db:"profile" json:"profile"`
	Email    string `db:"email" json:"email"`
	Phone    string `db:"phone" json:"phone"`
	Avatar   string `db:"avatar" json:"avatar"`
}
