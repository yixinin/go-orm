package orm

import (
	"fmt"
	"strings"
)

//Where Where("k=?"",v) or Where("k = v")
func (o *Orm) Where(where string, v ...interface{}) *Orm {
	o.wheres = append(o.wheres, fmt.Sprintf("(%s)", where))
	if len(v) > 0 {
		o.whereArgs = append(o.whereArgs, v...)
	}
	return o
}

//Eq =
func (o *Orm) Eq(k string, v interface{}) *Orm {
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` = ?)", k))
	o.whereArgs = append(o.whereArgs, v)
	return o
}

//Ne <>
func (o *Orm) Ne(k string, v interface{}) *Orm {
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` <> ?)", k))
	o.whereArgs = append(o.whereArgs, v)
	return o
}

//Gt >=
func (o *Orm) Gt(k string, v interface{}) *Orm {
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` > ?)", k))
	o.whereArgs = append(o.whereArgs, v)
	return o
}

//Gte >=
func (o *Orm) Gte(k string, v interface{}) *Orm {
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` >= ?)", k))
	o.whereArgs = append(o.whereArgs, v)
	return o
}

//Lt <
func (o *Orm) Lt(k string, v interface{}) *Orm {
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` < ?)", k))
	o.whereArgs = append(o.whereArgs, v)
	return o
}

//Lte <=
func (o *Orm) Lte(k string, v interface{}) *Orm {
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` <= ?)", k))
	o.whereArgs = append(o.whereArgs, v)
	return o
}

//In ...
func (o *Orm) In(k string, s interface{}) *Orm {
	var slice = joinSlice(s)
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` in(%s))", k, slice))
	return o
}

//NotIn ..
func (o *Orm) NotIn(k string, s interface{}) *Orm {
	var slice = joinSlice(s)
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` not in(%s))", k, slice))
	return o
}

//Between ...
func (o *Orm) Between(k string, v1, v2 interface{}) *Orm {
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` BETWEEN ? AND ?)", k))
	o.whereArgs = append(o.whereArgs, v1, v2)
	return o
}

//IsNull  ..
func (o *Orm) IsNull(k string, isNull bool) *Orm {
	if isNull {
		o.wheres = append(o.wheres, fmt.Sprintf("(`%s` is null)", k))
		return o
	}
	o.wheres = append(o.wheres, fmt.Sprintf("(`%s` is not null)", k))
	return o
}

//Like Like("name","xx","<>")
func (o *Orm) Like(k, v, side string) *Orm {
	switch side {
	case "<":
		o.wheres = append(o.wheres, fmt.Sprintf("(`%s` like '%%%s')", k, v))
	case ">":
		o.wheres = append(o.wheres, fmt.Sprintf("(`%s` like '%s%%')", k, v))
	case "<>":
		o.wheres = append(o.wheres, fmt.Sprintf("(`%s` like '%%%s%%')", k, v))
	}
	return o
}

//Cond ..
func (o *Orm) Cond(c *Cond) *Orm {
	o.wheres = append(o.wheres, strings.Join(c.wheres, " AND "))
	o.whereArgs = append(o.whereArgs, c.whereArgs...)
	return o
}
