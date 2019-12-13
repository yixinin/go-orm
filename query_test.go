package orm

import (
	"fmt"
	"testing"
)

func TestCount(t *testing.T) {
	InitTest()
	var o = NewOrm()
	count, err := o.ForQuery(&User{}).In("id", []int{1072, 1073}).Count()
	fmt.Println(o.Explain())
	fmt.Println(count, err)
}

func TestFindOne(t *testing.T) {
	InitTest()
	var o = NewOrm()
	var user User
	user.Id = 1078
	// var cond = NewCond().Like("name", "yi", "<>")

	err := o.FindOne(&user)
	fmt.Println(o.Explain())
	fmt.Printf("%+v,%v\n", user, err)
}

func TestFindMany(t *testing.T) {
	InitTest()
	var o = NewOrm()
	var users []*User

	// var query = o.Table(User{}.TableName()).Gt("id", 1078).IsNull("name", true)
	// fmt.Println(query.Count())

	err := o.Find(&users, "name")
	fmt.Println(o.Explain())
	for _, name := range users {
		fmt.Printf("%+v\n", name)
	}
	fmt.Println(err)
}

func TestFindMap(t *testing.T) {
	InitTest()
	var o = NewOrm()
	var m = make(map[int]string, 3)
	m[1079] = ""
	m[1089] = ""
	m[1081] = ""
	err := o.Table(User{}.TableName()).Gt("id", 1078).FindMap(&m, "id", "create_time")
	fmt.Println(o.Explain())
	fmt.Println(m, err)
}
