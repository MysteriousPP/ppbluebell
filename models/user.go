package models

const (
	admin  string = "admin"
	normal string = "normal"
)

type User struct {
	UserID       int64  `db:"user_id"`
	Username     string `db:"username"`
	Password     string `db:"password"`
	UserType     string `db:"usertype"`
	AccessToken  string
	RefreshToken string
}
