package orm

import (
	"fmt"
	"reflect"
)

var colMap map[string]map[string]Col
var pkMap map[string]string

func init() {
	colMap = make(map[string]map[string]Col, 100)
	pkMap = make(map[string]string, 100)
}

//Register Register Model
func Register(v interface{}) error {
	var t = reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		t = reflect.ValueOf(v).Elem().Type()
	}
	c := make(map[string]Col, t.NumField())
	var pk = "Id"
	for i := 0; i < t.NumField(); i++ {
		var f = t.Field(i)
		var tag = t.Field(i).Tag.Get("orm")
		if tag == "" {
			tag = toSnake(f.Name)
		}
		c[tag] = Col{
			t: t.Field(i).Type,
			i: i,
		}
		if _, b := f.Tag.Lookup("pk"); b {
			pk = f.Name
		}
	}
	table, ok := v.(TableNameble)
	if ok {
		colMap[table.TableName()] = c
		pkMap[table.TableName()] = pk
		return nil
	}
	return fmt.Errorf("please impl `func TableName() string`, struct name:%s", t.Name())
}
