package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const TokenExpireDuration = time.Hour * 24 * 365

var (
	ErrorExpiredAccessToken  error = errors.New("access token expired")
	ErrorExpiredRefreshToken error = errors.New("refresh token expired")
)

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte("ccbccbccb")

// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
	// 可根据需要自行添加字段
	UserID               int64  `json:"userid"`
	Username             string `json:"username"`
	UserType             string `json:"usertype"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

type RefreshClaims struct {
	UserID   int64  `json:"userid"`
	Username string `json:"username"`
	UserType string `json:"usertype"`
	jwt.RegisteredClaims
}

// GenToken 生成JWT
func GenAccessToken(userID int64, username string, usertype string) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		userID,
		username, // 自定义字段
		usertype,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(viper.GetInt("auth.jwt_access_expire")) * time.Second)),
			Issuer:    "bluebell", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

func GenRefreshToken(userID int64, username string, usertype string) (string, error) {
	// 创建一个我们自己的声明
	claims := RefreshClaims{
		userID,
		username, // 自定义字段
		usertype,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(viper.GetInt("auth.jwt_refresh_expire")) * time.Second)),
			Issuer:    "bluebell", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseAccessToken 解析Access JWT
func ParseAccessToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				zap.L().Error("access token expired", zap.Error(err))
				return token.Claims.(*CustomClaims), ErrorExpiredAccessToken
			}
		}
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				zap.L().Error("refresh token expired", zap.Error(err))
				return nil, ErrorExpiredRefreshToken
			}
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
