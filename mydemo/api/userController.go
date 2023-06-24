package api

import (
	"demo/database"
	"demo/middleware"
	"demo/model"
	"demo/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Register 注册函数
func Register(ctx *gin.Context) {

	//获取参数
	//此处使用Bind()函数，可以处理不同格式的前端数据
	var requestUser model.User
	//错误判断
	if err := ctx.ShouldBindWith(&requestUser, binding.Form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数解析错误",
		})
		return
	}
	name := requestUser.Name
	password := requestUser.Password

	//无误时
	ctx.JSON(http.StatusOK, gin.H{
		"name":     name,
		"password": password,
	})

	//数据验证
	if len(name) == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户名不能为空",
		})
		return
	}

	//判断用户名是否存在
	var u1 model.User
	database.DB.Where("username = ?", name).First(&u1) //查询
	if u1.ID != 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户已存在",
		})
		return
	}

	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码不能少于6位",
		})
		return
	}

	//创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    500,
			"message": "密码加密错误",
		})
		return
	}
	newUser := model.User{
		Name:     name,
		Password: string(hashedPassword),
	}

	database.DB.Create(&newUser)

	//返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
	})
}

// Login 登录函数
func Login(ctx *gin.Context) {

	//获取参数
	//此处使用Bind()函数，可以处理不同格式的前端数据
	var requestUser model.User
	if err := ctx.ShouldBindWith(&requestUser, binding.Form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数解析错误",
		})
		return
	}
	username := requestUser.Name
	password := requestUser.Password

	//数据验证
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码不能少于6位",
		})
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(requestUser.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码错误",
		})
	}

	//返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
	})

	// 正确则登录成功
	// 创建一个我们自己的声明
	claim := model.MyClaims{
		Username: username, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // 过期时间
			Issuer:    "YYY",                                // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, _ := token.SignedString(middleware.Secret)
	utils.SucRes(ctx, tokenString)
}

// CreateQuestion 创建问题
func CreateQuestion(c *gin.Context) {

	// 从请求中解析问题的标题和内容
	var request struct {
		Title   string `form:"title" json:"title" binding:"required"`
		Content string `form:"content" json:"content" binding:"required"`
	}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户的名字，已经通过jwt实现了用户认证和授权部分
	username := database.GetCurrentUserName(c) // 自定义函数，获取当前用户的ID

	if username == nil {
		utils.FailRes(c, "获取用户名失败!")
	}

	// 创建问题
	question := model.Question{
		Title:    request.Title,
		Content:  request.Content,
		UserName: username.(string),
	}

	// 保存问题到数据库
	if err := database.DB.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, question)
}

// GetQuestions 获取所有问题
func GetQuestions(c *gin.Context) {

	var questions []model.Question
	database.DB.Find(&questions)

	c.JSON(200, questions)

}

// GetQuestionDetail 获取问题详情
func GetQuestionDetail(c *gin.Context) {

	questionID := c.Param("id")

	var question model.Question
	if err := database.DB.First(&question, questionID).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Question not found",
		})
		return
	}
	c.JSON(200, question)
}

// SubmitAnswer 提交回答
func SubmitAnswer(c *gin.Context) {

	questionID := c.Param("id")

	var question model.Question
	if err := database.DB.First(&question, questionID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Question not found"})
		return
	}

	username := database.GetCurrentUserName(c)
	if username == "" {
		utils.FailRes(c, "获取用户名失败!")
	}

	var answer model.Answer
	answer.UserName = username.(string)
	if err := c.ShouldBind(&answer); err != nil {
		c.JSON(400, gin.H{"error": "Invalid answer data"})
		return
	}

	answer.QuestionID = question.ID
	database.DB.Create(&answer)

	c.JSON(200, answer)
}

// GetQuestionsByUser 根据提问者的用户名获取问题列表
func GetQuestionsByUser(c *gin.Context) {

	username := c.Param("username")

	var questions []model.Question
	if err := database.DB.Where("user_name = ?", username).Find(&questions).Error; err != nil {
		c.JSON(404, gin.H{"error": "No questions found for the given username"})
		return
	}

	c.JSON(200, questions)
}

// GetAnswersByUser 根据提问者的用户名获取回答列表
func GetAnswersByUser(c *gin.Context) {

	username := c.Param("username")

	var answers []model.Answer
	if err := database.DB.Where("user_name = ?", username).Find(&answers).Error; err != nil {
		c.JSON(404, gin.H{"error": "No answers found for the given username"})
		return
	}

	c.JSON(200, answers)
}

// UpdateQuestion 修改自己的问题
func UpdateQuestion(c *gin.Context) {

	questionID := c.Param("id")

	var question model.Question
	if err := database.DB.First(&question, questionID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Question not found"})
		return
	}

	// 检查问题是否属于当前用户
	if question.UserName != database.GetCurrentUserName(c) {
		c.JSON(403, gin.H{"error": "Can not update this question"})
		return
	}

	// 从请求正文中获取更新后的问题内容
	var updatedQuestion model.Question
	if err := c.ShouldBind(&updatedQuestion); err != nil {
		c.JSON(400, gin.H{"error": "Invalid question data"})
		return
	}

	// 更新问题内容
	question.Title = updatedQuestion.Title
	question.Content = updatedQuestion.Content
	database.DB.Save(&question)

	c.JSON(200, question)
}

// UpdateAnswer 修改自己的回答
func UpdateAnswer(c *gin.Context) {

	answerID := c.Param("id")

	var answer model.Answer
	if err := database.DB.First(&answer, answerID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Question not found"})
		return
	}

	// 检查该回答是否属于当前用户
	if answer.UserName != database.GetCurrentUserName(c) {
		c.JSON(403, gin.H{"error": "Not authorized to update this answer"})
		return
	}

	// 从请求正文中获取更新后的回答内容
	var updatedAnswer model.Answer
	if err := c.ShouldBind(&updatedAnswer); err != nil {
		c.JSON(400, gin.H{"error": "Invalid answer data"})
		return
	}

	// 更新回答
	answer.Content = updatedAnswer.Content
	database.DB.Save(&answer)

	c.JSON(200, answer)
}

// DeleteMyQuestion 删除自己的问题
func DeleteMyQuestion(c *gin.Context) {

	questionID := c.Param("id")

	var question model.Question
	if err := database.DB.First(&question, questionID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Question not found"})
		return
	}

	// 检查当前用户是否是问题的作者
	if question.UserName != database.GetCurrentUserName(c) {
		c.JSON(403, gin.H{"error": "You don't have permission to delete this question"})
		return
	}

	// 执行删除操作
	if err := database.DB.Unscoped().Delete(&question).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete question"})
		return
	}

	c.JSON(200, gin.H{"message": "Question deleted successfully"})
}

// DeleteMyAnswer 删除自己的回答
func DeleteMyAnswer(c *gin.Context) {

	AnswerID := c.Param("id")

	var answer model.Answer
	if err := database.DB.First(&answer, AnswerID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Question not found"})
		return
	}

	// 检查当前用户是否是问题的作者
	if answer.UserName != database.GetCurrentUserName(c) {
		c.JSON(403, gin.H{"error": "Can not  delete this answer"})
		return
	}

	// 执行删除操作
	if err := database.DB.Delete(&answer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete answer"})
		return
	}

	c.JSON(200, gin.H{"message": "Answer deleted successfully"})
}
