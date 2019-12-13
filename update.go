package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//UpdateParam ..
type UpdateParam map[string]interface{}

//ForUpdate ..
func (o *Orm) ForUpdate(v interface{}) *Orm {
	o.op = OpUpdate
	o.parseInterface(v)
	return o
}

//Update map[string]interface{} struct
func (o *Orm) Update(v interface{}) (int64, error) {
	o.op = OpUpdate
	defer func() {
		exp := fmt.Sprintf("%s, %v", o.update.queryString, o.update.args)
		o.update = &Update{
			explain: exp,
		}
	}()

	var value = reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Ptr:
		value = reflect.Indirect(value)
		switch value.Kind() {
		case reflect.Struct:
			o.structValue = value
			o.parseTableName(v)
			o.parsePk()
		}
	case reflect.Struct:
		o.structValue = value
		o.parseTableName(v)
		o.parsePk()
	}
	if v != nil {
		if p, ok := v.(UpdateParam); ok {
			o.parseUpdateParams(p)
		} else {
			o.parseUpdateSets(v)
		}
	}
	if len(o.update.sets) == 0 {
		return 0, errors.New("nothing for update")
	}

	o.parseUpdate()

	res, err := o.mysql.Exec(o.update.queryString, o.update.args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

//UpdateOne ..
func (o *Orm) UpdateOne(v interface{}) (int64, error) {
	return o.Limit(1).Update(v)
}

//UpdateMulti ..
func (o *Orm) UpdateMulti(s interface{}, ignoreErrs ...bool) ([]int, error) {
	o.op = OpUpdate
	defer func() {
		exp := fmt.Sprintf("%s, %v", o.update.queryString, o.update.args)
		o.update = &Update{
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
	var cs = make([]string, 0)
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

		if i == 0 {
			cols := colMap[o.tableName]
			for k := range cols {
				cs = append(cs, k)
			}
		}

		var args = make([]interface{}, 0, item.NumField())
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

		var snakePk = toSnake(o.pk)
		for _, col := range cs {
			if col == snakePk {
				continue
			}
			if i == 0 {
				o.update.sets = append(o.update.sets, fmt.Sprintf("`%s` = ?", col))
			}
			args = append(args, item.FieldByName(fromSnake(col)).Interface())
		}
		args = append(args, pv)
		argsSlice = append(argsSlice, args)
	}

	o.update.queryString = fmt.Sprintf("update `%s` set %s where %s = ?", o.tableName, strings.Join(o.update.sets, ","), toSnake(o.pk))

	stmt, err := o.mysql.Prepare(o.update.queryString)
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

//Set ...
//Set(k,v) or Set("k = v")
func (o *Orm) Set(k string, v ...interface{}) *Orm {
	if l := len(v); l == 0 {
		o.update.sets = append(o.update.sets, k)
	} else {
		o.update.sets = append(o.update.sets, fmt.Sprintf("%s = ?", k))
		o.update.setArgs = append(o.update.setArgs, v[0])
	}
	return o
}
