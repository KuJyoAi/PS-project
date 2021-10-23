package main

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/develop1024/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	//"lib_jwt/lib_jwt"
)

const (
	databaseName = "root:KuJyo2796747392@tcp(127.0.0.1:3306)/test?charset=utf8"
)

var db *gorm.DB

type user struct {
	username   string
	passwordHS string
}
type userinfo struct {
	Signature string
	Name      string
}

func (userinfo) TableName() string {
	return "jwt"
}
func main() {
	var err error
	db, err = gorm.Open(mysql.Open(databaseName), &gorm.Config{})
	if err != nil {
		fmt.Println("err:", err)
	}

	r := gin.Default()
	r.GET("/rigster", func(c *gin.Context) {
		user_claims := user{
			username:   c.Query("name"),
			passwordHS: c.Query("psw"),
		}
		//查询重名
		if !FindDatabaseExist(user_claims.username) {
			//key设置为名字 base64url encoding
			sig := EncodeToken(base64.URLEncoding.EncodeToString([]byte(user_claims.username)), &user_claims)
			fmt.Println("SigPrint", userinfo{
				Signature: sig,
				Name:      user_claims.username,
			})
			SaveToDatabase(userinfo{
				Signature: sig,
				Name:      user_claims.username,
			})
			c.JSON(200, "签名:"+sig)
		} else {
			c.JSON(200, "名称已经被注册")
		}

	})

	r.GET("/login", func(c *gin.Context) {
		user_claims := user{
			username:   c.Query("name"),
			passwordHS: c.Query("sig"),
		}
		psw := c.Query("psw")
		if psw == "" {
			if FindDatabaseExist(user_claims.username) {
				//fmt.Printf("username:%s,psw:%s", user_claims.username, user_claims.passwordHS)
				token := DecodeToken(user_claims.passwordHS, base64.URLEncoding.EncodeToString([]byte(user_claims.username)))
				//fmt.Println(token)
				if token.Valid {
					c.JSON(200, fmt.Sprintf("使用签名登录成功,%s", user_claims.username))
				} else {
					c.JSON(200, fmt.Sprintf("签名匹配失败,%s", user_claims.username))
				}
			} else {
				c.JSON(200, fmt.Sprintf("没有此用户:%s", user_claims.username))
			}
		} else {
			if FindDatabaseExist(user_claims.username) {
				//fmt.Printf("username:%s,psw:%s", user_claims.username, user_claims.passwordHS)
				user_claims.passwordHS = EncodeToken(base64.URLEncoding.EncodeToString([]byte(user_claims.username)), &user_claims)
				token := DecodeToken(user_claims.passwordHS, base64.URLEncoding.EncodeToString([]byte(user_claims.username)))
				//fmt.Println(token)
				if token.Valid {
					c.JSON(200, fmt.Sprintf("使用密码登录成功,%s,签名是:%s", user_claims.username, user_claims.passwordHS))
				} else {
					c.JSON(200, fmt.Sprintf("签名匹配失败,%s", user_claims.username))
				}
			} else {
				c.JSON(200, fmt.Sprintf("没有此用户:%s", user_claims.username))
			}
		}
	})
	r.Run()
}

func FindDatabaseExist(name string) bool {
	//fmt.Println(usf)
	err := db.Where("name = ?", name).Take(&userinfo{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("Query None")
		return false
	} else if err == nil {
		fmt.Println("Query Success")
		return true
	} else {
		panic(err)
	}

}
func SaveToDatabase(user userinfo) {
	db.Create(&user)
}
func EncodeToken(key string, users *user) string {
	claims := jwt.MapClaims{
		"username":   users.username,
		"passwordHS": users.passwordHS,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(key))
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println(claims)
	//fmt.Println(token)
	//fmt.Println(ss)
	return ss
}
func DecodeToken(tokenString string, key string) *jwt.Token {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	//fmt.Printf("header:%s,claims:%s\n", token.Header, token.Claims)
	//fmt.Println(token.Valid)
	//fmt.Println(token)
	//claims := token.Claims.(jwt.MapClaims)
	//fmt.Println(claims)
	return token
}
