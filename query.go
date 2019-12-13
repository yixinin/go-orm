package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//ForQuery ..
func (o *Orm) ForQuery(v interface{}) *Orm {
	o.op = OpQuery
	o.parseInterface(v)
	return o
}

//Sort ..  Sort("key") or Sort("-key")
func (o *Orm) Sort(k string) *Orm {
	if len(k) > 0 {
		if k[0] == '-' {
			o.query.sorts = append(o.query.sorts, fmt.Sprintf("%s DESC", k))
		} else {
			o.query.sorts = append(o.query.sorts, k)
		}
	}

	return o
}

//FindOne ..
func (o *Orm) FindOne(v interface{}, cols ...string) error {
	o.op = OpQuery
	defer func() {
		exp := fmt.Sprintf("%s, %v", o.query.queryString, o.query.args)
		o.query = &Query{
			explain: exp,
		}
	}()
	interfaceValue, structType, err := o.parseInterface(v)
	if err != nil {
		return err
	}

	fmt.Println(structType.Name())

	cs := o.parseCols(cols)
	if len(cs) == 0 {
		return errors.New("no select fields")
	}
	o.query.cols = cs

	if o.pk == "" {
		o.parsePk()
	}

	if !o.parseQuery() {
		return errors.New("parse query error, no table name")
	}

	var scaners = make([]interface{}, len(o.query.cols))
	for i := range o.query.cols {
		var v interface{}
		scaners[i] = &v
	}

	err = o.mysql.QueryRow(o.query.queryString, o.query.args...).Scan(scaners...)
	if err != nil {
		return err
	}

	if structType.Kind() == reflect.Struct && structType.Name() != typeTime {
		o.loadData(interfaceValue, scaners)
	} else {
		o.loadSingleData(interfaceValue, scaners[0].(*interface{}), structType.Name())
	}
	return err
}

//Find ..
func (o *Orm) Find(s interface{}, cols ...string) error {
	o.op = OpQuery
	defer func() {
		exp := fmt.Sprintf("%s, %v", o.query.queryString, o.query.args)
		o.query = &Query{
			explain: exp,
		}
	}()

	sliceValue, structType, err := o.parseInterface(s)
	if err != nil {
		return err
	}

	cs := o.parseCols(cols)
	if len(cs) == 0 {
		return errors.New("no select field")
	}
	o.query.cols = cs

	if !o.parseQuery() {
		return errors.New("parse query error, no table name")
	}

	var scaners []interface{}

	if structType.Kind() != reflect.Struct {
		var v interface{} = reflect.New(structType)
		scaners = []interface{}{
			&v,
		}
	} else {
		scaners = make([]interface{}, len(o.query.cols))
		for i := range o.query.cols {
			var v interface{}
			v = reflect.New(o.query.colm[o.query.cols[i]].t)
			scaners[i] = &v
		}
	}

	rows, err := o.mysql.Query(o.query.queryString, o.query.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(scaners...); err != nil {
			return err
		}
		if structType.Kind() == reflect.Struct && structType.Name() != typeTime {
			var itemptr = reflect.New(structType)
			item := reflect.Indirect(itemptr)
			o.loadData(item, scaners)
			sliceValue.Set(reflect.Append(sliceValue, itemptr))
		} else {
			var value = reflect.New((structType))
			var indValue = reflect.Indirect(value)
			var inf, ok = scaners[0].(*interface{})
			if ok && inf != nil {
				o.loadSingleData(indValue, inf, structType.Name())
			}
			if structType.Name() == typeTime {
				sliceValue.Set(reflect.Append(sliceValue, value))
			} else {
				sliceValue.Set(reflect.Append(sliceValue, indValue))
			}

		}
	}
	return nil
}

//FindMap ..
func (o *Orm) FindMap(m interface{}, kCol, vCol string) error {
	o.op = OpQuery
	defer func() {
		exp := fmt.Sprintf("%s, %v", o.query.queryString, o.query.args)
		o.query = &Query{
			explain: exp,
		}
	}()

	var mKeyType, mValueType reflect.Type
	// var mKvalue, mVvalue reflect.Value
	var mapValue reflect.Value
	var kvs []interface{}
	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		mapValue = reflect.Indirect(reflect.ValueOf(m))
		if mapValue.Kind() == reflect.Map {
			mKeyType = mapValue.Type().Key()

			mValueType = mapValue.Type().Elem()
			if mValueType.Kind() == reflect.Ptr {
				mValueType = mValueType.Elem()
			}
			keys := mapValue.MapKeys()
			if len(keys) > 0 {
				kvs = make([]interface{}, 0, len(keys))
				for _, key := range keys {
					kvs = append(kvs, key.Interface())
				}
			} else {
				return errors.New("no map content")
			}
		} else {
			return errors.New("not map")
		}
	} else {
		return errors.New("not ptr")
	}

	var in = joinSlice(kvs)
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` in(%s))", kCol, in))
	if len(o.whereArgs) > 0 {
		o.query.args = o.whereArgs
	}
	var query = fmt.Sprintf("select `%s`,`%s` from `%s` where %s", kCol, vCol, o.tableName, strings.Join(o.wheres, " AND "))
	o.query.queryString = query
	rows, err := o.mysql.Query(o.query.queryString, o.query.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	var valueType = mValueType.Name()
	var keyType = mKeyType.Name()
	for rows.Next() {
		var key = reflect.Indirect(reflect.New(mKeyType))
		var value = reflect.New(mValueType)
		var indValue = reflect.Indirect(value)
		var k, v interface{}
		if err := rows.Scan(&k, &v); err != nil {
			continue
		}
		if infK, ok := k.(*interface{}); ok {
			o.loadSingleData(key, infK, keyType)
		} else {
			o.loadSingleData(key, &k, keyType)
		}
		if infV, ok := v.(*interface{}); ok {
			o.loadSingleData(indValue, infV, valueType)
		} else {
			o.loadSingleData(indValue, &v, valueType)
		}
		if valueType == typeTime {
			mapValue.SetMapIndex(key, value)
		} else {
			mapValue.SetMapIndex(key, indValue)
		}

	}
	return nil
}

func (o *Orm) loadData(item reflect.Value, scaners []interface{}) {
	for i, col := range o.query.cols {
		c := o.query.colm[col]
		if inf, ok := scaners[i].(*interface{}); ok {
			if inf != nil {
				if buf, ok := (*inf).([]byte); ok && len(buf) > 0 {
					var str = string(buf)
					switch c.t.Name() {
					case typeString:
						item.Field(c.i).SetString(str)
					case typeTime:
						v, _ := time.ParseInLocation(timeLayout, str, time.Local)
						item.Field(c.i).Set(reflect.ValueOf(v))
					case typeBool:
						v, _ := strconv.ParseBool(str)
						item.Field(c.i).SetBool(v)
					case typeUint8:
						v, _ := strconv.ParseUint(str, 10, 8)
						item.Field(c.i).SetUint(v)
					case typeUint16:
						v, _ := strconv.ParseUint(str, 10, 16)
						item.Field(c.i).SetUint(v)
					case typeUint, typeUint32:
						v, _ := strconv.ParseUint(str, 10, 32)
						item.Field(c.i).SetUint(v)
					case typeUint64:
						v, _ := strconv.ParseUint(str, 10, 64)
						item.Field(c.i).SetUint(v)
					case typeInt8:
						v, _ := strconv.ParseInt(str, 10, 8)
						item.Field(c.i).SetInt(v)
					case typeInt16:
						v, _ := strconv.ParseInt(str, 10, 16)
						item.Field(c.i).SetInt(v)
					case typeInt, typeInt32:
						v, _ := strconv.ParseInt(str, 10, 32)
						item.Field(c.i).SetInt(v)
					case typeInt64:
						v, _ := strconv.ParseInt(str, 10, 64)
						item.Field(c.i).SetInt(v)
					case typeFloat32:
						v, _ := strconv.ParseFloat(str, 32)
						item.Field(c.i).SetFloat(v)
					case typeFloat64:
						v, _ := strconv.ParseFloat(str, 64)
						item.Field(c.i).SetFloat(v)
					}
				} else if v, ok := (*inf).(int64); ok {
					switch c.t.Name() {
					case typeUint, typeUint8, typeUint16, typeUint32, typeUint64:
						item.Field(c.i).SetUint(uint64(v))
					case typeFloat32, typeFloat64:
						item.Field(c.i).SetFloat(float64(v))
					case typeBool:
						item.Field(c.i).SetBool(v == 1)
					default:
						item.Field(c.i).SetInt(v)
					}
				} else if v, ok := (*inf).(float64); ok {
					switch c.t.Name() {
					case typeBool:
						item.Field(c.i).SetBool(v == 1)
					default:
						item.Field(c.i).SetFloat(v)
					}
				} else if v, ok := (*inf).(uint64); ok {
					switch c.t.Name() {
					case typeInt, typeInt8, typeInt16, typeInt32, typeInt64:
						item.Field(c.i).SetInt(int64(v))
					case typeFloat32, typeFloat64:
						item.Field(c.i).SetFloat(float64(v))
					case typeBool:
						item.Field(c.i).SetBool(v == 1)
					default:
						item.Field(c.i).SetUint(v)
					}
				} else if v, ok := (*inf).(bool); ok {
					item.Field(c.i).SetBool(v)
				}
			}
		}
	}
}

func (o *Orm) loadSingleData(interfaceValue reflect.Value, inf *interface{}, structType string) {
	if buf, ok := (*inf).([]byte); ok {
		var str string
		if len(buf) == 0 {
			str = ""
		} else {
			str = string(buf)
		}

		switch structType {
		case typeString:
			interfaceValue.SetString(str)
		case typeTime:
			v, _ := time.ParseInLocation(timeLayout, str, time.Local)
			interfaceValue.Set(reflect.ValueOf(v))
		case typeBool:
			v, _ := strconv.ParseBool(str)
			interfaceValue.SetBool(v)
		case typeUint8:
			v, _ := strconv.ParseUint(str, 10, 8)
			interfaceValue.SetUint(v)
		case typeUint16:
			v, _ := strconv.ParseUint(str, 10, 16)
			interfaceValue.SetUint(v)
		case typeUint, typeUint32:
			v, _ := strconv.ParseUint(str, 10, 32)
			interfaceValue.SetUint(v)
		case typeUint64:
			v, _ := strconv.ParseUint(str, 10, 64)
			interfaceValue.SetUint(v)
		case typeInt8:
			v, _ := strconv.ParseInt(str, 10, 8)
			interfaceValue.SetInt(v)
		case typeInt16:
			v, _ := strconv.ParseInt(str, 10, 16)
			interfaceValue.SetInt(v)
		case typeInt, typeInt32:
			v, _ := strconv.ParseInt(str, 10, 32)
			interfaceValue.SetInt(v)
		case typeInt64:
			v, _ := strconv.ParseInt(str, 10, 64)
			interfaceValue.SetInt(v)
		case typeFloat32:
			v, _ := strconv.ParseFloat(str, 32)
			interfaceValue.SetFloat(v)
		case typeFloat64:
			v, _ := strconv.ParseFloat(str, 64)
			interfaceValue.SetFloat(v)
		}
	} else if v, ok := (*inf).(int64); ok {
		switch structType {
		case typeUint, typeUint8, typeUint16, typeUint32, typeUint64:
			interfaceValue.SetUint(uint64(v))
		case typeFloat32, typeFloat64:
			interfaceValue.SetFloat(float64(v))
		case typeBool:
			interfaceValue.SetBool(v == 1)
		default:
			interfaceValue.SetInt(v)
		}
	} else if v, ok := (*inf).(float64); ok {
		switch structType {
		case typeBool:
			interfaceValue.SetBool(v == 1)
		default:
			interfaceValue.SetFloat((v))
		}
	} else if v, ok := (*inf).(uint64); ok {
		switch structType {
		case typeInt, typeInt8, typeInt16, typeInt32, typeInt64:
			interfaceValue.SetInt(int64(v))
		case typeFloat32, typeFloat64:
			interfaceValue.SetFloat(float64(v))
		case typeBool:
			interfaceValue.SetBool(v == 1)
		default:
			interfaceValue.SetUint((v))
		}
	} else if v, ok := (*inf).(bool); ok {
		interfaceValue.SetBool(v)
	}
}

//Count ..
func (o *Orm) Count() (int64, error) {
	var count int64
	query, args := o.getCountQuery()
	o.query.explain = query
	err := o.mysql.QueryRow(query, args...).Scan(&count)
	return count, err
}

//Exist ..
func (o *Orm) Exist() (bool, error) {
	n, err := o.Count()
	return n > 0, err
}
