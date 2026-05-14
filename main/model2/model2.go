package model2

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type Content struct {
	ID        string    `gorm:"column:ID"`
	User_id   string    `gorm:"column:user_id"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Content) TableName() string {
	return "Content"
}

func DeteleContent(cid, uid string, db *gorm.DB) {
	result := db.Where("`user_id` = ? AND `ID` = ?", uid, cid).Delete(&Content{})
	if result.Error != nil {
		fmt.Println("在删除数据这里的mysql删除错误!", result.Error)
		return
	}
}

func CreateContent(id, user_id, content string, db *gorm.DB) bool { // 创建用户逻辑
	user_content := Content{
		ID:      id,
		User_id: user_id,
		Content: content,
	}
	err := db.Create(&user_content).Error
	if err != nil {
		log.Println("保存文章失败:", err)
		return false
	}
	return true
}
