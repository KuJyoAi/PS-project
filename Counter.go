package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	//SQL
	/*
			go get github.com/go-sql-driver/mysql
		    go get github.com/jmoiron/sqlx
	*/
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	//需要更改的量:
	PostName  = "/api"
	DbName    = "root:@tcp(127.0.0.1:3306)/counter?charset=utf8"
	TableName = "counter"
)
/*
create table counter(id int,time int,type int,saved int);
*/

var Db *sqlx.DB

func init() {
	DB, err := sqlx.Open("mysql", DbName)
	if err != nil {
		fmt.Println("Open database failed! err:", err)
		panic(err)
	}
	Db = DB

}

func main() {
	r := gin.Default()
	r.POST(PostName, func(c *gin.Context) {
		SaveToDatabase(c.Query("id"), c.Query("time"), c.Query("rqtype"), c.Query("saved"))
		c.JSON(200, gin.H{
			//返回信息
			"message": "Success",
		})
	})
	r.Run()

	defer Db.Close()
}

//数据库存储
type User struct {
	ID     int
	Time   int
	RqType int
	Saved  int
}

func SaveToDatabase(id, time, RqType, Saved string) {
	//留作params的解析与处理
	ID, _ := strconv.Atoi(id)
	Time, err := strconv.Atoi(time)
	//转换整数出错,是小数
	if err != nil {
		tmp, _ := strconv.ParseFloat(time, 32)
		Time = int(tmp + 1)
	}
	//rquest
	var RqtypeTmp int = 0

	switch RqType {
	case "":
		RqtypeTmp = 0
	default:
		var err error
		RqtypeTmp, err = strconv.Atoi(RqType)
		if err != nil {
			fmt.Println("RqType converse failed")
		}
	}
	//Saved 默认false
	IsSaved, _ := strconv.Atoi(Saved)
	//fmt.Println("121", TableName, User{ID, int(Time), RqtypeTmp, IsSaved})
	InsertDB(User{ID, int(Time), RqtypeTmp, IsSaved})

	/*
		def:
		Rqtype:0=未知
		Save:0=false 1=true

	*/

}
func InsertDB(emp User) {
	var tmp string = fmt.Sprintf("insert into %v(id,time,type,saved)values(%v,%v,%v,%v)\n",
		TableName, emp.ID, emp.Time, emp.RqType, emp.Saved)
	//fmt.Print(tmp)
	_, err := Db.Exec(tmp)
	if err != nil {
		fmt.Println("Insert failed! err:", err)
		return
	}
}
