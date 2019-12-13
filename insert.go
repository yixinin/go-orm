package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	//DefaultInsertMulti ..
	DefaultInsertMulti = 1000
)

//InsertOne ..
func (o *Orm) InsertOne(v interface{}) (int64, error) {
	o.op = OpInsert
	defer func() {
		exp := fmt.Sprintf("%s, %v", o.insert.queryString, o.insert.args)
		o.insert = &Insert{
			explain: exp,
		}
	}()
	var pk = "Id"
	var cols, cs, remarks []string
	var colM = make(map[string]bool, 0)
	var args []interface{}

	if err := o.parseTableName(v); err != nil {
		return 0, err
	}
	cs = o.parseCols(nil)

	if len(cs) == 0 {
		return 0, errors.New("no insert field")
	}

	if k, ok := pkMap[o.tableName]; ok {
		pk = k
	}

	var t = reflect.TypeOf(v)
	var value reflect.Value

	switch t.Kind() {
	case reflect.Ptr:
		value = reflect.ValueOf(v).Elem()
		if value.Kind() != reflect.Struct {
			return 0, errors.New("insert value must be struct or *struct")
		}
	case reflect.Struct:
		value = reflect.ValueOf(v)
	}

	for _, col := range cs {
		f := value.FieldByName(fromSnake(col))
		v := f.Interface()
		if !f.IsZero() {
			args = append(args, v)
			colM[col] = true
		}
	}

	for _, v := range cs {
		if b, ok := colM[v]; ok && b {
			cols = append(cols, fmt.Sprintf("`%s`", v))
			remarks = append(remarks, "?")
		}
	}

	var query = fmt.Sprintf("insert into `%s` (%s) values (%s)", o.tableName, strings.Join(cols, ","), strings.Join(remarks, ","))
	o.insert.queryString = query
	o.insert.args = args
	res, err := o.mysql.Exec(query, args...)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if t.Kind() == reflect.Ptr {
		value.FieldByName(pk).SetInt(id)
	}
	return id, err
}

//InsertMany ..
//atom if true, one fail all fail
func (o *Orm) InsertMany(s interface{}, ignoreErrs ...bool) ([]int, error) {
	defer func() {
		o.op = OpInsert
		exp := fmt.Sprintf("%s, %v", o.insert.queryString, o.insert.args)
		o.insert = &Insert{
			explain: exp,
		}
	}()

	var cols, cs, remarks []string
	var colM = make(map[string]bool, 0)

	var sliceValue reflect.Value
	var t = reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		sliceValue = reflect.Indirect(reflect.ValueOf(s))
	} else {
		sliceValue = reflect.ValueOf(s)
	}

	var structType reflect.Type
	sliceElementType := sliceValue.Type().Elem()
	if sliceElementType.Kind() == reflect.Ptr {
		structType = sliceElementType.Elem()
		if structType.Kind() == reflect.Struct {
			if o.tableName == "" {
				pv := reflect.New(structType)
				if err := o.parseTableNameBySlice(sliceElementType, pv); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, errors.New("slice must be struct")
		}
	} else {
		structType = sliceElementType
		pv := reflect.New(structType)
		if err := o.parseTableNameBySlice(sliceElementType, pv); err != nil {
			return nil, err
		}
	}

	cs = o.parseCols(nil)

	//读取数据
	var argsSlice [][]interface{}
	for i := 0; i < sliceValue.Len(); i++ {

		var item reflect.Value
		itemV := sliceValue.Index(i)
		if itemV.Kind() == reflect.Ptr {
			item = reflect.Indirect(itemV)
		} else {
			item = itemV
		}

		var args = make([]interface{}, 0, item.NumField())
		for _, col := range cs {
			f := item.FieldByName(fromSnake(col))
			v := f.Interface()

			if f.IsZero() && (col == pkMap[o.tableName] || f.Type().Name() == "Time") {
				colM[col] = false
				args = append(args, nil)
			} else {
				colM[col] = true
				args = append(args, v)
			}
		}
		argsSlice = append(argsSlice, args)
	}

	var removeCount = 0
	for i, v := range cs {
		if b, ok := colM[v]; ok && b {
			cols = append(cols, fmt.Sprintf("`%s`", v))
			remarks = append(remarks, "?")
		} else {
			for j := range argsSlice {
				var index = i - removeCount
				if index == len(cs)-1 {
					argsSlice[j] = (argsSlice[j])[:index]
				} else {
					argsSlice[j] = append((argsSlice[j])[:index], (argsSlice[j])[index+1:]...)
				}
			}
			removeCount++
		}
	}

	var query = fmt.Sprintf("insert into `%s` (%s) values (%s)", o.tableName, strings.Join(cols, ","), strings.Join(remarks, ","))
	o.insert.queryString = query
	stmt, err := o.mysql.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var failIndexs = make([]int, 0, len(argsSlice))
	var ignoreErr bool
	if len(ignoreErrs) > 0 {
		ignoreErr = ignoreErrs[0]
	}

	for i, args := range argsSlice {
		if _, err1 := stmt.Exec(args...); err1 != nil {
			err = err1
			if ignoreErr {
				failIndexs = append(failIndexs, i)
				continue
			}
			return nil, err1
		}
	}
	return failIndexs, nil
}
