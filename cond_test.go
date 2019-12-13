package orm

import (
	"fmt"
	"testing"
)

func TestCond(t *testing.T) {
	var cond = NewCond()
	var c1 = cond.Eq("id", 10).Gt("name", "xx")
	var c2 = NewCond().Eq("age", 9)
	var c3 = c1.Or(c2)
	fmt.Println(c3.wheres, c3.whereArgs)
}
