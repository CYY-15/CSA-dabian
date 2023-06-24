package database

import (
	"demo/model"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() {

	driverName := "mysql"
	host := "127.0.0.1"
	port := "3306"
	database := "demo"
	username := "root"
	password := "123456"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB")

	//迁移三大模型
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Question{})
	db.AutoMigrate(&model.Answer{})

	// 定义外键关系
	db.Model(&model.Answer{}).AddForeignKey("question_id", "questions(id)", "CASCADE", "CASCADE")

	DB = db
	return
}

func GetCurrentUserName(c *gin.Context) any {
	username, exists := c.Get("username")
	print(username, 1)
	if exists == false {
		return nil
	}

	return username
}
