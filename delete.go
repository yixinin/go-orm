package orm

import (
	"errors"
	"fmt"
	"reflect"
)

//ForDelete ..
func (o *Orm) ForDelete(v interface{}) *Orm {
	o.op = OpDelete
	o.parseInterface(v)
	return o
}

//DeleteOne ..
func (o *Orm) DeleteOne(s ...interface{}) (int64, error) {
	return o.Limit(1).Delete(s...)
}

//Delete ..
func (o *Orm) Delete(s ...interface{}) (int64, error) {
	o.op = OpDelete
	defer func() {
		exp := fmt.Sprintf("%s, %v", o.delete.queryString, o.delete.args)
		o.delete = &Delete{
			explain: exp,
		}
	}()
	if len(s) > 0 {
		_, _, err := o.parseInterface(s[0])
		if err != nil {
			return 0, err
		}
	}
	if !o.parseDelete() {
		return 0, errors.New("impl TableName")
	}

	res, err := o.mysql.Exec(o.delete.queryString, o.delete.args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

//DeleteMulti ..
func (o *Orm) DeleteMulti(s interface{}, ignoreErrs ...bool) ([]int, error) {
	o.op = OpDelete
	defer func() {
		exp := fmt.Sprintf("%s, %v", o.delete.queryString, o.delete.args)
		o.delete = &Delete{
			explain: exp,
		}
	}()

	var ignoreErr bool
	if len(ignoreErrs) > 0 {
		ignoreErr = ignoreErrs[0]
	}

	var sliceValue, _, err = o.parseInterface(s)
	if err != nil {
		return nil, err
	}
	//读取数据
	var argsSlice [][]interface{}
	for i := 0; i < sliceValue.Len(); i++ {
		var item reflect.Value
		itemV := sliceValue.Index(i)
		switch itemV.Kind() {
		case reflect.Ptr:
			item = reflect.Indirect(itemV)
		case reflect.Struct:
			item = itemV
		default:
			return nil, errors.New("must be *[]Struct or []Struct or []*Struct or *[]*Struct")
		}

		if o.pk == "" {
			var pk, ok = pkMap[o.tableName]
			if ok {
				o.pk = pk
			}
		}

		var pkValue = item.FieldByName(o.pk)
		if !pkValue.IsValid() || pkValue.IsZero() {
			return nil, errors.New("pk value is zero")
		}
		pv := pkValue.Interface()
		if pv == nil {
			return nil, errors.New("pk value is nil")
		}

		argsSlice = append(argsSlice, []interface{}{pv})
	}

	o.delete.queryString = fmt.Sprintf("delete from `%s` where %s = ?", o.tableName, toSnake(o.pk))

	stmt, err := o.mysql.Prepare(o.delete.queryString)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var fails = make([]int, 0)
	for i, args := range argsSlice {
		if _, err1 := stmt.Exec(args...); err1 != nil {
			err = err1
			if ignoreErr {
				fails = append(fails, i)
				continue
			}
			return nil, err1
		}
	}
	return fails, err
}
