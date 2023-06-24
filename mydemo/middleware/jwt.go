package middleware

import (
	"demo/model"
	"demo/utils"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

var Secret = []byte("YYY")

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//假设Token放在Header的Authorization中，使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		//如果没有auth的错误处理
		if authHeader == "" {
			utils.FailRes(c, "empty authorization in Header")
			c.Abort()
			return
		}
		//如果有，那么继续判断是否使用Bearer开头
		//按空格分开成两份，存入parts数组中
		parts := strings.SplitN(authHeader, " ", 2)

		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.FailRes(c, "wrong authorization formation")
			c.Abort()
			return
		}

		//parts[1]是我们获取到的token，通过我们定义的Parse-token函数来解析
		//mc变量存储了一个Claim声明
		mc, err := Parsetoken(parts[1])
		if err != nil {
			utils.FailRes(c, "Invalid Token")
			c.Abort()
			return
		}
		//将返回的声明中的username存入上下文中
		c.Set("username", mc.Username)
		c.Next()
	}
}

func Parsetoken(tokenString string) (*model.MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
