package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 成功后的响应

func SucRes(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": message,
	})
}

// 失败后的响应

func FailRes(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  500,
		"message": message,
	})
}
