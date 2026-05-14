package dbmodel

import (
	utils "BOOK/main/utils_model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 写一个用文章id查文章的逻辑，主要是redis会过期
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
func CtidFindCt(ctid string, db *gorm.DB) (string, error) {
	content, err := utils.Get(ctid)
	if err != nil || content == "" {
		fmt.Println("redis没文章id对应文章内容！")
		var ct Content
		if err = db.Where("ID = ?", ctid).First(&ct).Error; err != nil {
			fmt.Println("MySQL也没有(文章id找文章)")
			return "", err
		}
		fmt.Println("MySQL有(文章id找文章)")
		utils.Set(ctid, ct.Content, time.Hour*24)
		return ct.Content, err
	}
	return content, err
}

// 还有一个如果查用户id查不出文章id的话，就

func UidfindCtid(Uid string, db *gorm.DB) ([]string, error) {
	contentids, err := utils.ListRange(Uid, 0, -1)
	if err != nil || len(contentids) == 0 {
		fmt.Println("redis数据库没有用户id对应文章id！准备去MySQL找")
		var ctid []Content
		var ctid_1 []string
		if err = db.Where("user_id = ?", Uid).Find(&ctid).Error; err != nil {
			fmt.Println("mysql也没有（2）！")
			return []string{}, err
		}
		for _, cid := range ctid {
			utils.ListPush(Uid, cid.ID)
			ctid_1 = append(ctid_1, cid.ID)
			fmt.Println("文章id:(dbmodel)", cid.ID)
		}
		// fmt.Println("dbmodel测试",ctid_1)
		return ctid_1, err
	}
	// fmt.Println("dbmodel测试",contentids) // 这里没问题
	return contentids, err
}

// 上面写一个存入
