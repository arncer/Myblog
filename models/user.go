package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model        //使用gorm内嵌的一个结构体,它包含了ID、CreatedAt、UpdatedAt、DeletedAt等字段
	UserName   string `gorm:"unique"`
	Password   string
}
