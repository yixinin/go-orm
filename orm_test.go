package orm

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id         int
	Name       string
	Age        uint8
	Status     int8
	CreateTime time.Time
	Price      float64
	Money      int64
	IsFamale   bool
}

func (u User) TableName() string {
	return "user"
}

func InitTest() {
	Init(&MysqlConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "123456",
		DB:       "orm",
		MaxIdle:  10,
		MaxConn:  10,
	})

	var u = &User{}
	Register(u)
}

func TestParseInterface(t *testing.T) {
	InitTest()
	var user = User{
		Id:         1078,
		Name:       "xx1",
		Age:        108,
		CreateTime: time.Now(),
		Price:      2.02,
		Money:      10000000000,
		IsFamale:   false,
		Status:     2,
	}
	var o = NewOrm()
	var users = []User{
		user,
	}
	_, _, err := o.parseInterface(&users)
	fmt.Println(o.tableName, err)
}
