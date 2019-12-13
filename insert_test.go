package orm

import (
	"fmt"
	"testing"
	"time"
)

func TestInsertOne(t *testing.T) {
	InitTest()
	var u = &User{
		Name:       "yixin",
		Age:        18,
		CreateTime: time.Now(),
		Price:      1.02,
		// Money:      10000000000,
		IsFamale: true,
		Status:   1,
	}
	var o = NewOrm()
	id, err := o.InsertOne(u)
	fmt.Println(o.Explain())
	fmt.Printf("id=%v err=%v", id, err)
}

func TestInsertMany(t *testing.T) {
	InitTest()
	var u1 = &User{
		Name: "zz",
		Age:  18,
		// CreateTime: time.Now(),
		Price: 1.02,
		// Money:      10000000000,
		IsFamale: true,
		Status:   1,
	}
	var u2 = &User{
		Name:       "xx",
		Age:        18,
		CreateTime: time.Now(),
		Price:      1.02,
		// Money:      10000000000,
		IsFamale: false,
		Status:   1,
	}
	var users []*User
	users = append(users, u1, u2)
	var o = NewOrm()
	fails, err := o.InsertMany(users, true)
	fmt.Println(o.Explain())
	fmt.Println(fails, err)
}
