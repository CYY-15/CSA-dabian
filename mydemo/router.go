package main

import (
	"demo/api"
	"demo/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoutes(r *gin.Engine) *gin.Engine {

	//注册
	r.POST("/register", api.Register)

	//登录
	r.POST("/login", api.Login)

	// 在 authGroup 中定义需要进行身份验证的路由
	// 创建一个路由组，并将 JWT 中间件应用于该组
	rouGroup := r.Group("/login")
	{
		//引入jwt中间件
		rouGroup.Use(middleware.JWTAuthMiddleware())

		//登录后提问
		rouGroup.POST("/new-question", api.CreateQuestion)

		//获取问题列表
		rouGroup.GET("/questions", api.GetQuestions)

		//获取问题详情
		rouGroup.GET("/questions/:id", api.GetQuestionDetail)

		//回答问题
		rouGroup.POST("/questions/:id/new-answer", api.SubmitAnswer)

		//查看自己的问题or答案
		rouGroup.GET("/questions/users/:username", api.GetQuestionsByUser)
		rouGroup.GET("/answers/users/:username", api.GetAnswersByUser)

		//删除自己的问题or回答
		rouGroup.DELETE("/questions/:id", api.DeleteMyQuestion)
		rouGroup.DELETE("/answers/:id", api.DeleteMyAnswer)

		//修改自己的问题or回答
		rouGroup.PUT("/questions/:id", api.UpdateQuestion)
		rouGroup.PUT("/answers/:id", api.UpdateAnswer)

	}

	return r

}
