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
	authGroup := r.Group("/login")
	{
		//引入jwt中间件
		authGroup.Use(middleware.JWTAuthMiddleware())

		//登录后提问
		authGroup.POST("/new-question", api.CreateQuestion)

		//获取问题列表
		authGroup.GET("/questions", api.GetQuestions)
		//获取问题详情
		authGroup.GET("/questions/:id", api.GetQuestionDetail)
		//查看问题下的答案
		authGroup.GET("/questions/:id/answers", api.GetQuestionWithAnswers)

		//回答问题
		authGroup.POST("/questions/:id/new-answer", api.SubmitAnswer)

		//从用户名搜索问题or答案
		//搜索自己就是查看自己的问题or答案
		authGroup.GET("/questions/users/:username", api.GetQuestionsByUser)
		authGroup.GET("/answers/users/:username", api.GetAnswersByUser)

		//删除自己的问题or回答
		authGroup.DELETE("/questions/:id", api.DeleteMyQuestion)
		authGroup.DELETE("/answers/:id", api.DeleteMyAnswer)

		//修改自己的问题or回答
		authGroup.PUT("/questions/:id", api.UpdateQuestion)
		authGroup.PUT("/answers/:id", api.UpdateAnswer)

		//评论回答

		//给问题点赞

		//注销账户

	}

	return r

}
