package main

import (
	_ "fmt"
	"user_information/common"
	"user_information/model"
)

func main() {
	if !common.File_exist { // 创造一个文件
		model.File_create(common.File_path)
	}
	for {
		// 主界面，以及返回一个选择数
		model.Selection_1(model.Main_page()) // 传入一个选择数，看看用什么函数
		if common.Xun {
			break
		}
	}
}
