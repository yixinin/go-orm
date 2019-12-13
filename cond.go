package orm

import (
	"fmt"
	"strings"
)

//Cond ..
type Cond struct {
	wheres    []string
	whereArgs []interface{}
}

//NewCond ..
func NewCond() *Cond {
	return &Cond{}
}

//Or ..
func (c *Cond) Or(cond *Cond) *Cond {
	c.wheres = []string{
		fmt.Sprintf("((%s) OR (%s))", strings.Join(c.wheres, " AND "), strings.Join(cond.wheres, " AND ")),
	}
	c.whereArgs = append(c.whereArgs, cond.whereArgs...)
	return c
}

//And ..
func (c *Cond) And(cond *Cond) *Cond {
	c.wheres = []string{
		fmt.Sprintf("(%s) AND (%s)", strings.Join(c.wheres, " AND "), strings.Join(cond.wheres, " AND ")),
	}
	c.whereArgs = append(c.whereArgs, cond.whereArgs...)
	return c
}

//Where Where("k=?"",v) or Where("k = v")
func (c *Cond) Where(where string, v ...interface{}) *Cond {
	c.wheres = append(c.wheres, fmt.Sprintf("(%s)", where))
	if len(v) > 0 {
		c.whereArgs = append(c.whereArgs, v...)
	}
	return c
}

//Eq =
func (c *Cond) Eq(k string, v interface{}) *Cond {
	c.wheres = append(c.wheres, fmt.Sprintf("(`%s` = ?)", k))
	c.whereArgs = append(c.whereArgs, v)
	return c
}

//Ne <>
func (c *Cond) Ne(k string, v interface{}) *Cond {
	c.wheres = append(c.wheres, fmt.Sprintf("(`%s` <> ?)", k))
	c.whereArgs = append(c.whereArgs, v)
	return c
}

//Gt >=
func (c *Cond) Gt(k string, v interface{}) *Cond {
	c.wheres = append(c.wheres, fmt.Sprintf("(`%s` > ?)", k))
	c.whereArgs = append(c.whereArgs, v)
	return c
}

//Gte >=
func (c *Cond) Gte(k string, v interface{}) *Cond {
	c.wheres = append(c.wheres, fmt.Sprintf("(`%s` >= ?)", k))
	c.whereArgs = append(c.whereArgs, v)
	return c
}

//Lt <
func (c *Cond) Lt(k string, v interface{}) *Cond {
	c.wheres = append(c.wheres, fmt.Sprintf("(`%s` < ?)", k))
	c.whereArgs = append(c.whereArgs, v)
	return c
}

//Lte <=
func (c *Cond) Lte(k string, v interface{}) *Cond {
	c.wheres = append(c.wheres, fmt.Sprintf("(`%s` <= ?)", k))
	c.whereArgs = append(c.whereArgs, v)
	return c
}

//In ...
func (c *Cond) In(k string, s interface{}) *Cond {
	var slice = joinSlice(s)
	c.wheres = append(c.wheres, fmt.Sprintf("(`%s` in(%s))", k, slice))
	return c
}

//NotIn ..
func (c *Cond) NotIn(k string, s interface{}) *Cond {
	var slice = joinSlice(s)
	c.wheres = append(c.wheres, fmt.Sprintf("(`%s` not in(%s))", k, slice))
	return c
}

//Between ...
func (c *Cond) Between(k string, v1, v2 interface{}) *Cond {
	c.wheres = append(c.wheres, fmt.Sprintf("(`%s` BETWEEN ? AND ?)", k))
	c.whereArgs = append(c.whereArgs, v1, v2)
	return c
}

//IsNull  ..
func (c *Cond) IsNull(k string, isNull bool) *Cond {
	if isNull {
		c.wheres = append(c.wheres, fmt.Sprintf("(`%s` is null)", k))
		return c
	}
	c.wheres = append(c.wheres, fmt.Sprintf("(`%s` is not null)", k))
	return c
}

//Like Like("name","xx","<>")
func (c *Cond) Like(k, v, side string) *Cond {
	switch side {
	case "<":
		c.wheres = append(c.wheres, fmt.Sprintf("(`%s` like '%%%s')", k, v))
	case ">":
		c.wheres = append(c.wheres, fmt.Sprintf("(`%s` like '%s%%')", k, v))
	case "<>":
		c.wheres = append(c.wheres, fmt.Sprintf("(`%s` like '%%%s%%')", k, v))
	}
	return c
}
