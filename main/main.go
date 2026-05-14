package main

import (
	"BOOK/main/dbmodel"
	"BOOK/main/model"
	"BOOK/main/model2"
	"BOOK/main/utils_model"
	"BOOK/main/utils_model/redis"
	"bytes"
	"fmt"
	"html/template"

	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	// "strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func markdownToHTML(mdContent string) template.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(mdContent), &buf); err != nil {
		return template.HTML("解析错误")
	}
	return template.HTML(buf.String())
}

func main() {

	if err := redis.InitRedis(); err != nil {
		log.Fatalf("Redis 初始化失败: %v", err)
	}
	fmt.Print("redis连接成功!")

	r := gin.Default()
	// r.Static("/", "../html")
	r.Static("/css", "../html/css")
	r.Static("/images", "../html/images")

	dbUser := "root"
	dbPass := "ROOT"
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "users"
	encodedPass := url.QueryEscape(dbPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, encodedPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("连接初始化失败: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("获取底层数据库对象失败: %v", err))
	}

	fmt.Println("安全连接成功！")
	defer sqlDB.Close()

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if path == "/" {
			c.File("../html/index.html")
			return
		}

		filePath := "../html" + path + ".html"

		if _, err := os.Stat(filePath); err == nil {
			c.File(filePath)
			return
		}

		c.File("../html/index.html")
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		if bl, user := model.User_find(username, password, db); bl {
			claims := jwt.MapClaims{
				"id":   float64(user.ID),
				"name": user.Username,
				"exp":  time.Now().Add(time.Hour * 72).Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, _ := token.SignedString([]byte("123miyao")) // 密钥

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "登录成功!",
				"token":   tokenString,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "用户名或密码错误!",
			})}})

	r.POST("/register", func(c *gin.Context) {
		new_username := c.PostForm("username")
		new_password := c.PostForm("password")
		// 这里可以添加注册逻辑，保存用户名和密码到数据库
		if model.CreateUser(new_username, new_password, db) {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "注册成功!"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "注册失败!"})
		}
	})
	type RequestBody struct { // write
		Content string `json:"content"`
		ID      int    `json:"id"`
		RGO     string `json:"rgo"`
	}

	type RequestBody2 struct { // home
		ID uint `json:"id"`
	}

	type RequestBody_view struct {
		ContentId string `json:"content_id"`
	}

	type detele_1 struct {
		Uid int    `json:"uid"`
		Cid string `json:"cid"`
	}

	aapp := r.Group("/api")
	aapp.Use(model.JWTAuthMiddleware())
	{
		r.POST("/api/write", func(c *gin.Context) {
			var reqBody RequestBody
			if err := c.ShouldBindJSON(&reqBody); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			fmt.Println("请求体内容:", c.Request.Body)
			content := reqBody.Content
			id := reqBody.ID
			Rgo := reqBody.RGO
			htmlContent := markdownToHTML(content)

			fmt.Println("id:", id)

			log.Println("生成的 HTML 是：", htmlContent)
			c.JSON(http.StatusOK, gin.H{
				"Content": htmlContent,
			})

			// 这里的redis的key是content_id，值是文章详细内容
			if Rgo == "y" {
				id_time := time.Now().UnixMilli()
				content_id := fmt.Sprintf("%d", int(id_time)) + "00" + fmt.Sprintf("%d", id)
				fmt.Println("时间戳:", id_time)
				fmt.Println("文章id:", content_id)
				//  第一个是文章id，第二个是：文章内容
				if err := utils.Set(content_id, string(htmlContent), 24*time.Hour); err != nil {
					fmt.Println("redis保存失败")
				}
				fmt.Println("redis保存成功!")
				fmt.Println("准备在redis存文章id到用户id下")
				// 第一个是用户id，第二个是文章id
				fmt.Println("存入的用户id是", id)
				fmt.Println("存入的文章id是", content_id)
				utils.ListPush(fmt.Sprintf("%d", id), content_id)
				model2.CreateContent(content_id, fmt.Sprintf("%d", id), string(htmlContent), db) // 存入mysql
			}
		})

		aapp.POST("/check", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		})

		aapp.POST("/home", func(c *gin.Context) {
			var reqBody RequestBody2
			var contents_2 []string
			if err := c.ShouldBindJSON(&reqBody); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			id := reqBody.ID

			// 拿一下文章id
			contents_id, err := dbmodel.UidfindCtid(fmt.Sprintf("%d", id), db)
			if err != nil {
				fmt.Println("主页面从redis取文章id失败!")
				return
			}

			for _, cot := range contents_id {
				content_1, err_2 := dbmodel.CtidFindCt(cot, db)
				if err_2 != nil {
					fmt.Println("从文章id—文章数据库取数据失败！", err)
					return 
				}
				contents_2 = append(contents_2, content_1)
			}

			c.JSON(200, gin.H{
				"contents_2": contents_2,
				"content_id": contents_id,
			})
		})
		aapp.POST("/view", func(c *gin.Context) {
			var reqBody RequestBody_view
			if err := c.ShouldBindJSON(&reqBody); err != nil {
				fmt.Println("111")
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			contentid := reqBody.ContentId
			fmt.Println(contentid)
			// 接下来取用文章id取文章内容
			content, err2 := dbmodel.CtidFindCt(contentid, db)
			if err2 != nil {
				fmt.Println("view取文件有问题！")
				return
			}
			c.JSON(200, gin.H{
				// "ok": true,
				"Content": content,
			})
		})
		aapp.DELETE("/view", func(c *gin.Context) {
			var msg detele_1

			fmt.Println("111")
			if err := c.ShouldBindJSON(&msg); err != nil {
				fmt.Println("删除处接收json错误！", err)
				return
			}

			uid := fmt.Sprintf("%d", msg.Uid)
			cid := msg.Cid
			fmt.Println("uid和cid", uid, cid)

			err_2 := redis.RedisClient.Del(c.Request.Context(), cid).Err()
			if err_2 != nil {
				fmt.Println("在detele中redis的单键删除的err:", err_2)
				return
			}

			_, err_3 := redis.RedisClient.LRem(c.Request.Context(), uid, 0, cid).Result()
			if err_3 != nil {
				fmt.Println("在detele中的redis的list删除失败!", err_3)
				return
			}

			model2.DeteleContent(cid, uid, db)

			c.JSON(200, gin.H{
				"msg": "key已删除!",
			})
		})
	}
	r.Run(":8080")
}
