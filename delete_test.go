package orm

import (
	"fmt"
	"testing"
)

func TestDeleteOne(t *testing.T) {
	InitTest()
	var o = NewOrm()
	n, err := o.Eq("name", "yixin").Delete(&User{Id: 1072})
	fmt.Println(o.Explain())
	fmt.Println(n, err)
}

func TestDeleteMulti(t *testing.T) {
	InitTest()
	var o = NewOrm()
	var users = make([]*User, 1)
	users[0] = &User{
		Id: 1077,
	}
	n, err := o.DeleteMulti(users)
	fmt.Println(o.Explain())
	fmt.Println(n, err)
}
