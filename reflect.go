package orm

import (
	"errors"
	"reflect"
)

func (o *Orm) parseInterface(p interface{}) (interfacevalue reflect.Value, structType reflect.Type, err error) {
	defer func() {
		if o.pk == "" {
			pk, ok := pkMap[o.tableName]
			if ok {
				o.pk = pk
			}
		}
		if o.pk != "" && (o.pkValue == nil || reflect.ValueOf(o.pkValue).IsZero()) {
			if o.structValue.IsValid() && !o.structValue.IsZero() {
				pkValue := o.structValue.FieldByName(o.pk)
				if pkValue.IsValid() && !pkValue.IsZero() {
					o.pkValue = pkValue.Interface()
				}
			}
		}
	}()
	value := reflect.ValueOf(p)
	switch value.Kind() {
	case reflect.String:
		o.tableName = p.(string)
	case reflect.Struct:
		o.structValue = value
		if table, ok := p.(TableNameble); ok {
			o.tableName = table.TableName()
		}
	case reflect.Ptr:
		ptrValue := reflect.Indirect(value)
		switch ptrValue.Kind() {
		case reflect.Slice:
			interfacevalue = ptrValue
			sliceElementType := interfacevalue.Type().Elem()
			switch sliceElementType.Kind() {
			case reflect.Ptr:
				structType = sliceElementType.Elem()
				if structType.Kind() == reflect.Struct && structType.Name() != typeTime {
					pv := reflect.New(structType)
					if err = o.parseTableNameBySlice(sliceElementType, pv); err != nil {
						return
					}
				} else {
					return
				}
			case reflect.Struct:
				structType = sliceElementType
				pv := reflect.Indirect(reflect.New(structType))
				if err = o.parseTableNameBySlice(sliceElementType, pv); err != nil {
					return
				}
			default:
				structType = sliceElementType
				return
			}

		case reflect.Struct:
			structType = ptrValue.Type()
			interfacevalue = ptrValue
			if ptrValue.Type().Name() != typeTime {
				o.structValue = ptrValue
				if table, ok := p.(TableNameble); ok {
					o.tableName = table.TableName()
				}
			}
		default: //接受值的指针
			structType = ptrValue.Type()
			interfacevalue = ptrValue
			return
		}
	case reflect.Slice:
		interfacevalue = value
		sliceElementType := interfacevalue.Type().Elem()
		switch sliceElementType.Kind() {
		case reflect.Ptr:
			structType = sliceElementType.Elem()
			if structType.Kind() == reflect.Struct {
				pv := reflect.New(structType)
				if err = o.parseTableNameBySlice(sliceElementType, pv); err != nil {
					return
				}
			} else {
				return
			}
		case reflect.Struct:
			structType = sliceElementType
			pv := reflect.Indirect(reflect.New(structType))
			if err = o.parseTableNameBySlice(sliceElementType, pv); err != nil {
				return
			}
		}
	default:
		err = errors.New("unsurpported type")
		return
	}
	return
}
