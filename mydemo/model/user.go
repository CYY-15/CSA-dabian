package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// 定义一个用户结构体

type User struct {
	gorm.Model
	Name     string `gorm:"varchar(20);not null;unique" form:"name" json:"name" binding:"required"`
	Password string `gorm:"size:255;not null" form:"password" json:"password" binding:"required"`
}

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 定义一个问题结构体

type Question struct {
	gorm.Model
	Title    string   `gorm:"type:varchar(100);not null"`
	Content  string   `gorm:"type:text;not null" form:"content"`
	UserName string   `gorm:"type:varchar(20);not null"`
	Answers  []Answer `gorm:"ForeignKey:QuestionID"`
}

// 定义一个回答结构体

type Answer struct {
	gorm.Model
	Content    string `gorm:"type:text;not null" form:"content"`
	UserName   string `gorm:"varchar(20);not null"`
	QuestionID uint   `gorm:"type:int unsigned;not null"`
}
