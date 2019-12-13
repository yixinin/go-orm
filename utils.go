package orm

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	layout = "2006-01-02 15:04:05"
)

func joinSlice(vs interface{}, elemName ...string) string {

	var value = reflect.ValueOf(vs)
	var isInterface = false

	var e string
	if len(elemName) < 1 || elemName[0] == "" {
		v := reflect.Indirect(value)
		t := v.Type().Elem()
		e = t.Name()
		if t.Kind() == reflect.Interface {
			var inf = v.Index(0).Interface()
			e = reflect.ValueOf(inf).Type().Name()
			isInterface = true
		}

	}
	var kind reflect.Kind
	if e == "string" {
		kind = reflect.String
	} else if strings.Contains(e, "uint") {
		kind = reflect.Uint
	} else if strings.Contains(e, "int") {
		kind = reflect.Int
	} else if strings.Contains(e, "float") {
		kind = reflect.Float64
	} else if e == "bool" {
		kind = reflect.Bool
	} else {
		return ""
	}

	// var kind = t.Kind()
	var s = make([]string, 0, value.Len())
	for i := 0; i < value.Len(); i++ {
		v := value.Index(i)
		switch kind {
		case reflect.String:
			if isInterface {
				s = append(s, fmt.Sprintf("'%s'", v.Interface()))
			} else {
				s = append(s, fmt.Sprintf("'%s'", v.String()))
			}
		case reflect.Int:
			if isInterface {
				s = append(s, fmt.Sprintf("%d", v.Interface()))
			} else {
				s = append(s, fmt.Sprintf("%d", v.Int()))
			}
		case reflect.Uint:
			if isInterface {
				s = append(s, fmt.Sprintf("%d", v.Interface()))
			} else {
				s = append(s, fmt.Sprintf("%d", v.Uint()))
			}

		case reflect.Float64:
			if isInterface {
				s = append(s, fmt.Sprintf("%f", v.Interface()))
			} else {
				s = append(s, fmt.Sprintf("%f", v.Float()))
			}
		case reflect.Bool:
			if isInterface {
				if v.Interface().(bool) {
					s = append(s, fmt.Sprintf("%d", 1))
				} else {
					s = append(s, fmt.Sprintf("%d", 0))
				}
			} else {
				if v.Bool() {
					s = append(s, fmt.Sprintf("%d", 1))
				} else {
					s = append(s, fmt.Sprintf("%d", 0))
				}
			}
		}

	}

	return strings.Join(s, ",")
}

func toSnake(str string) string {
	var s = make([]byte, 0, len(str)*2)
	for i, v := range []byte(str) {
		if v >= 'a' {
			s = append(s, v)
		} else if v <= 'z' && v >= 'A' {
			if i > 0 {
				s = append(s, '_')
			}
			s = append(s, v+32)
		}
	}
	return string(s)
}

func fromSnake(str string) string {
	var s = make([]byte, 0, len(str))
	var sIndex = 0
	for i, v := range []byte(str) {
		if v == '_' {
			sIndex = i
			continue
		}
		if i == 0 || (i != 1 && sIndex == i-1) {
			s = append(s, v-32)
		} else {
			s = append(s, v)
		}
	}
	return string(s)
}
