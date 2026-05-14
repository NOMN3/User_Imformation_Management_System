package model

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User_information struct {
	ID        int       `gorm:"column:id"`
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User_information) TableName() string {
	return "user_information"
}

func User_find(username, password string, db *gorm.DB) (bool, User_information) {
	// 实现查找用户逻辑
	var user User_information
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Println("用户名不存在:", err)
		return false, User_information{}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) //////////////////////////
	if err != nil {
		log.Println("用户密码不正确:", username)
		return false, User_information{}
	}
	// if user.Password != string(hashedPassword) {
	// 	log.Println("用户密码不正确:", username,password,user.Password,string(hashedPassword))
	// 	return false,User_information{}
	// }
	return true, user
}

func CreateUser(username, password string, db *gorm.DB) bool { // 创建用户逻辑
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) /////////////////////////
	if err != nil {                                                                          // 错误处理
		log.Println("密码哈希失败:", err)
		return false
	} ///////////////
	user := User_information{
		Username: username,
		Password: string(hashedPassword), //////////////////////////
	}
	if found, _ := User_find(username, password, db); found {
		log.Println("用户已存在:", username)
		return false
	}
	err = db.Create(&user).Error
	if err != nil {
		log.Println("创建用户失败:", err)
		return false
	}
	return true
}

// 检查用户有没有jwt的中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 JWT token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("签名算法是非法: %v", token.Header["alg"])
			}
			return []byte("123miyao"), nil // 密钥暂时是这个，之后会改成复杂又容易记的
		})
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "signature is invalid"):
				c.Redirect(http.StatusFound, "/login")
			case strings.Contains(err.Error(), "token is expired"):
				c.Redirect(http.StatusFound, "/login")
			default:
				c.Redirect(http.StatusFound, "/login")
			}
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["id"].(float64))
			c.Set("userID", userID)
			// 后续控制器中
			// userID, _ := c.Get("userID") // 获取到 1001
			// user, _ := db.GetUser(userID) // 查询用户数据
			c.Next()
		} else {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
	}
}
