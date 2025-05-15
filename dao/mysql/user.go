package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

const secret = "MysteriousPP"

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func GetPassword(username string) (password string, err error) {
	sqlStr := `select password from user where username = ?`
	if err := db.Get(&password, sqlStr, username); err != nil {
		return password, err
	}
	return password, nil
}

func CheckPassword(username string, password string) (bool, error) {
	encrypted_password := encryptPassword(password)
	var user_password string
	var err error
	if user_password, err = GetPassword(username); err != nil {
		return false, err
	}
	if encrypted_password == user_password {
		return true, nil
	} else {
		return false, errors.New("wrong username or password")
	}
}

func Login(user *models.User) (err error) {
	originPassword := user.Password // 记录一下原始密码
	sqlStr := "select user_id, username, password from user where username = ?"
	err = db.Get(user, sqlStr, user.Username)
	// if err != nil && err != sql.ErrNoRows {
	// 	// 查询数据库出错
	// 	return
	// }
	if err == sql.ErrNoRows {
		// 用户不存在
		return ErrorUserNotExit
	}
	// 生成加密密码与查询到的密码比较
	password := encryptPassword(originPassword)
	if user.Password != password {
		return ErrorPasswordWrong
	}
	return
}

func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}
