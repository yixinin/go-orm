package orm

import (
	"fmt"
	"testing"
	"time"
)

func TestUpdateOne(t *testing.T) {
	InitTest()
	var u = &User{
		Id:         1080,
		Name:       "yn",
		Age:        28,
		CreateTime: time.Now(),
		Price:      1.02,
		Money:      10000000000,
		IsFamale:   true,
		Status:     1,
	}
	var o = NewOrm()

	n, err := o.ForUpdate(u.TableName()).Pk("id", u.Id).UpdateOne(UpdateParam{"name": nil, "age": 19})
	fmt.Println(o.Explain())
	fmt.Println(n, err)
}

func TestUpdateMulti(t *testing.T) {
	InitTest()
	var o = NewOrm()
	var s = make([]User, 0, 2)
	var u1 = User{
		Id:         1077,
		Name:       "xxzz1",
		Age:        101,
		CreateTime: time.Now(),
		Price:      2.02,
		Money:      12000000000,
		IsFamale:   true,
		Status:     1,
	}
	var u2 = User{
		Id:         1078,
		Name:       "zzxx1",
		Age:        108,
		CreateTime: time.Now(),
		Price:      2.02,
		Money:      10000000000,
		IsFamale:   false,
		Status:     2,
	}
	s = append(s, u1, u2)
	fails, err := o.UpdateMulti(&s)
	fmt.Println(o.Explain())
	fmt.Println(fails, err)
}
