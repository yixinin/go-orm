package orm

import (
	"errors"
	"reflect"
)

func (o *Orm) parseTableName(v interface{}) error {
	if name, ok := v.(string); ok && name != "" {
		o.tableName = name
		return nil
	}
	table, ok := v.(TableNameble)
	if !ok {
		return errors.New("must impl TableNameble")
	}
	o.tableName = table.TableName()
	return nil
}

func (o *Orm) parseTableNameBySlice(t reflect.Type, vs ...reflect.Value) error {
	method, ok := t.MethodByName("TableName")
	if ok {
		var args = []reflect.Value{}
		if len(vs) > 0 {
			args = append(args, vs[0])
		}
		v := method.Func.Call(args)
		if len(v) > 0 {
			o.tableName = v[0].String()
			return nil
		}
	}
	return errors.New("must impl TableNameble")
}
