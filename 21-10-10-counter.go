package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	//SQL

	_ "github.com/go-sql-driver/mysql"
)

const (
	apiName   = "/api"
	DbName    = "root:@tcp(127.0.0.1:3306)/test?charset=utf8"
	TableName = "User"
)

var Db *sqlx.DB

func init() {
	DB, err := sqlx.Open("mysql", DbName)
	if err != nil {
		fmt.Println("Open database failed! err:", err)
		return
	}
	Db = DB
}

func main() {
	r := gin.Default()
	r.POST(apiName, func(c *gin.Context) {
		SaveToDatabase(c.Query("id"), c.Query("time"))
		c.JSON(200, gin.H{
			//返回信息
			"message": "返回成功",
		})
	})
	r.Run()

	defer Db.Close()
}

type User struct {
	ID   int
	Time int
}

/*
params收不到,暂废
func SaveToDatabase(par gin.Params) {
	fmt.Println(par)
}
*/
func SaveToDatabase(id string, time string) {
	//留作params的解析与处理
	ID, _ := strconv.Atoi(id)
	Time, err := strconv.Atoi(time)

	if err != nil {
		tmp, _ := strconv.ParseFloat(time, 32)
		InsertDB(User{ID, int(tmp + 1)})
		return
	}
	InsertDB(User{ID, int(Time)})
}
func InsertDB(emp User) {
	_, err := Db.Exec("insert into users(ID,Time)values(?,?)", emp.ID, emp.Time)
	if err != nil {
		fmt.Println("Insert failed! err:", err)
		return
	}
}
