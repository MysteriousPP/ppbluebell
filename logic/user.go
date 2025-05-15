package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// // 存放业务逻辑的代码
// func SignUpHandler() {
// 	//1.判断用户存不存在
// 	mysql.CheckUserExist()
// 	//2.生成UID
// 	snowflake.GenID()
// 	//3.保存进数据库
// 	mysql.InsertUser()
// }

func SignUp(p *models.ParamSignUp) (err error) {
	//1.判断用户存不存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	//2.生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3.保存进数据库
	return mysql.InsertUser(user)
}

//	func Login(p *models.ParamLogin) (err error) {
//		if err = mysql.CheckUserExist(p.Username); err != nil {
//			_, err := mysql.CheckPassword(p.Username, p.Password)
//			if err != nil {
//				return err
//			} else {
//				return nil
//			}
//		}
//		return err
//	}
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	if err = mysql.Login(user); err != nil {
		return nil, err
	}

	token, err := jwt.GenToken(user.UserID, p.Username)
	if err != nil {
		return
	}

	user.Token = token
	return
	//生成jwt
}
