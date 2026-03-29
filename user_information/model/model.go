package model

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"user_information/common"
)

var B bool // 是用在Selection_1函数里

type user struct { // 用户信息
	name      string
	gender    string
	age       int16
	call_id   string
	e_mail_id string
}

func Jgtpath() *user { // 返回用户结构体地址
	var a user
	return &a
}

type User_list struct { // 链节
	User *user
	Next *User_list
}

func Lj_list() *User_list { // 返回一个链节结构体的地址
	var a User_list
	return &a
}

// 搞一个头链条和全局变量的指针，首字母大写
var U_list User_list = User_list{User: Jgtpath(), Next: nil}
var P *User_list = &U_list // 这个p就是头链条
var Q *User_list = &U_list // 这个是会移动的链条指针(尾指针)

func File_create(path string) { // 创造用户信息文件(一次性)
	_, err := os.Stat(common.File_path)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("File_create err=%v\n", err)
			return
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		writer.WriteString("not\n")
		writer.Flush()
		common.File_exist = true
		return
	}
	fmt.Println("File_create err = ", err)
}

func Create_user_jgt() *user { // 建一个用户信息结构体，返回这个结构体的地址
	var a *user = Jgtpath()
	var age_2 string
	fmt.Println("客户信息：")
	fmt.Println("请输入客户名字：")
	fmt.Scanln(&(*a).name)
	fmt.Println("请输入客户性别：")
	fmt.Scanln(&(*a).gender)
	fmt.Println("请输入客户年龄：")
	for {
		fmt.Scanln(&age_2)
		age_3, err := strconv.ParseInt(age_2, 10, 16)
		if err == nil {
			(*a).age = int16(age_3)
			break
		} else {
			fmt.Println("年龄格式不规范！")
			fmt.Println("请重新输入客户年龄：")
		}
	}

	fmt.Println("请输入客户电话：")
	fmt.Scanln(&(*a).call_id)
	fmt.Println("请输入客户邮箱：")
	fmt.Scanln(&(*a).e_mail_id)
	return a
}

func Make_list_file() { // 这里需要打开文件，然后用一个string来接受里面的数据，再把他里面的数据分开装入单链表中
	file, err := os.OpenFile(common.File_path, os.O_RDONLY, 0666) //只读
	var str_2 [5]string
	if err != nil {
		fmt.Printf("Make_list_file err=%v\n", err)
		return
	}
	defer file.Close()
	const (
		defaultBufSize = 4096
	)
	reader := bufio.NewReader(file)
	str_2[0], err = reader.ReadString('\n')

	if str_2[0] == "not\n" { // 是无信息的话就直接退出，
		return // 不用参与下面的通过文件内容创建单链表了
	}

	//  需要一个去掉\n的函数
	str_2[0] = strings.TrimSuffix(str_2[0], "\r\n")

	var err_2 error
	for {
		if err_2 == io.EOF {
			break
		}

		var user_in *user
		var i int8 = 1
		for i < 5 {
			str_2[i], err_2 = reader.ReadString('\n')
			str_2[i] = strings.TrimSuffix(str_2[i], "\r\n")
			i++
		}

		user_age, _ := strconv.ParseInt(str_2[2], 10, 16)
		user_age_1 := int16(user_age)

		user_in = Jgtpath()

		*user_in = user{name: str_2[0], gender: str_2[1], age: user_age_1, call_id: str_2[3], e_mail_id: str_2[4]}
		// 先做一个链节，把数据录进去
		var user_list *User_list = Lj_list()
		*user_list = User_list{User: user_in, Next: nil}
		//  接下来就是把结构体录入单链条中了
		(*Q).Next = user_list
		Q = user_list

		str_2[0], err_2 = reader.ReadString('\n')
		//  需要一个去掉\n的函数
		str_2[0] = strings.TrimSuffix(str_2[0], "\r\n")
	}
}

func travel_list() {
	content, err := ioutil.ReadFile(common.File_path)
	if err != nil {
		fmt.Println("travel_list err = ", err)
	}
	if string(content) == "not\n" { // 是无信息的话就直接退出，
		fmt.Println("目前客户信息系统没有信息储存")
		return
	}
	var q *User_list = U_list.Next
	var i int = 1
	fmt.Println("-------------------------------")
	for {
		fmt.Printf("下面是第%d位的客户数据:\n", i)
		fmt.Println("名字：", (*(*q).User).name)
		fmt.Println("性别：", (*(*q).User).gender)
		fmt.Println("年龄：", (*(*q).User).age)
		fmt.Println("电话：", (*(*q).User).call_id)
		fmt.Println("邮箱：", (*(*q).User).e_mail_id)
		fmt.Println("-------------------------------")
		if (*q).Next == nil {
			break
		}
		q = (*q).Next
		i++
	}
	i = 1
}

func Make_String_From_List() string { // 把单链条转换成string
	var q *User_list = U_list.Next // 这个是会移动的链条指针,指向第一个数据链节结构体，而不是链条头
	if U_list.Next == nil {
		return "not\n"
	}
	var U_I_list_string string
	for {
		Age_str := fmt.Sprintf("%d", (*(*q).User).age)
		U_I_list_string += (*(*q).User).name + "\r\n" + (*(*q).User).gender + "\r\n" + Age_str + "\r\n" +
			(*(*q).User).call_id + "\r\n" + (*(*q).User).e_mail_id + "\r\n"

		if (*q).Next == nil {
			break
		}
		q = (*q).Next
	}
	return U_I_list_string
}

func Add_user_to_file(s string) { //这里的要求是打开文件，把数据录入进去
	file, err := os.OpenFile(common.File_path, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Add_user:err = ", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(s)
	writer.Flush()
}

func Add_U_I_T_list(U *user) { // 这里的目的是：1、把新建的结构体录入到一个新的链节里去，2、让Q的Next等于它
	// 先建一个链节结构体
	U_L := Lj_list()
	*U_L = User_list{User: U, Next: nil}
	// 剩下的就是把它们接起来了
	(*Q).Next = U_L
	Q = (*Q).Next
}

func Main_page() int8 { //主页面
	var selection int8
	fmt.Println("----------客户信息管理软件----------")
	fmt.Println("")
	fmt.Println("1 添加客户")
	fmt.Println("2 修改客户")
	fmt.Println("3 删除客户")
	fmt.Println("4 客户列表")
	fmt.Println("5 退   出")
	fmt.Println("")
	fmt.Printf("请选择(1-5):")
	fmt.Scanln(&selection)
	fmt.Println("")
	return selection
}

func del_or_chance_user_in(j int8) {
	var i int = 1
	var q *User_list = &U_list
	var select_2 int

	if U_list.Next == nil {
		fmt.Println("没有客户信息可以操作")
		return
	}

	travel_list()

	switch j {
	case 3:
		fmt.Printf("请问要删除第几个客户的信息:(结束请输入0)")
	case 2:
		fmt.Println("请问要修改第几个客户的信息")
	}

	fmt.Scanln(&select_2)

	for select_2 > i {
		q = (*q).Next
		if (*q).Next == nil {
			fmt.Println("超出范围！")
			return
		}
		i++
	}

	switch j {
	case 3:
		for {
			(*q).Next = (*(*q).Next).Next
			if U_list.Next == nil {
				fmt.Println("客户信息已删除完")
				return
			}
			travel_list()
			//  逻辑代码头
			i = 1
			q = &U_list
			fmt.Printf("请问要删除第几个客户的信息:(结束请输入0)")
			fmt.Scanln(&select_2)

			if select_2 == 0 {
				return
			}

			for select_2 > i {
				if (*q).Next == nil {
					fmt.Println("超出范围！")
					return
				}
				q = (*q).Next
				i++
			}
		}

	case 2:
		U_L := Lj_list()
		*U_L = User_list{User: Create_user_jgt(), Next: (*(*q).Next).Next}
		(*q).Next = U_L
	}

}

func Selection_1(selection int8) { //返回选择值
	if !B {
		Make_list_file()
		B = true
	}
	switch selection {
	case 1: // 添加用户
		Add_U_I_T_list(Create_user_jgt())
		Add_user_to_file(Make_String_From_List())
		fmt.Println("添加完成!")
	case 2: // 修改用户
		del_or_chance_user_in(selection)
		Add_user_to_file(Make_String_From_List())
	case 3: // 删除用户
		del_or_chance_user_in(selection)
		Q = &U_list
		for {
			if (*Q).Next == nil {
				break
			}
			Q = (*Q).Next
		}
		Add_user_to_file(Make_String_From_List())
	case 4:
		travel_list() // 客户列表
	case 5: // 退出程序
		common.Xun = true
	default:
		fmt.Println("输入错误")
	}
}
