package orm

import (
	"fmt"
	"reflect"
	"strings"
)

func (o *Orm) parseCols(cols []string) []string {
	if o.tableName != "" {
		if m, ok := colMap[o.tableName]; ok {
			o.query.colm = m
			if len(cols) > 0 {
				return cols
			}
			var cols = make([]string, 0, len(m))
			for k := range m {
				cols = append(cols, k)
			}
			return cols
		}
	}

	return nil
}

func (o *Orm) parseDelete() bool {
	if o.tableName == "" {
		return false
	}
	var query = fmt.Sprintf("delete from `%s` ", o.tableName)
	o.parsePkWhere()
	if len(o.wheres) > 0 {
		query = fmt.Sprintf("%s where %s", query, strings.Join(o.wheres, " AND "))
		o.delete.args = append(o.delete.args, o.whereArgs...)
	}
	if len(o.limit) > 0 && len(o.limitArgs) > 0 {
		query = fmt.Sprintf("%s %s", query, o.limit)
		o.delete.args = append(o.delete.args, o.limitArgs[0])
	}
	o.delete.queryString = query
	return true
}

func (o *Orm) parseUpdateSets(v interface{}) {
	var value = reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}
	if value.Kind() != reflect.Struct {
		return
	}
	var cols = colMap[o.tableName]
	var snakePk = toSnake(o.pk)
	for k, col := range cols {
		if k == snakePk {
			continue
		}
		if pkValue := value.Field(col.i).Interface(); k == o.pk && pkValue != nil {
			o.pkValue = pkValue
			continue
		}
		o.update.sets = append(o.update.sets, fmt.Sprintf("`%s` = ?", k))
		o.update.setArgs = append(o.update.setArgs, value.Field(col.i).Interface())
	}
}

func (o *Orm) parseUpdateParams(m UpdateParam) {
	for k, v := range m {
		o.update.sets = append(o.update.sets, fmt.Sprintf("%s = ?", k))
		o.update.setArgs = append(o.update.setArgs, v)
	}
}

func (o *Orm) parseUpdate() bool {
	if o.tableName == "" {
		return false
	}
	var query = fmt.Sprintf("update `%s`", o.tableName)
	if len(o.update.sets) > 0 {
		query = fmt.Sprintf("%s set %s ", query, strings.Join(o.update.sets, ","))
		o.update.args = append(o.update.args, o.update.setArgs...)
	} else {
		return false
	}
	o.parsePkWhere()
	if len(o.wheres) > 0 {
		query = fmt.Sprintf("%s where %s", query, strings.Join(o.wheres, " AND "))
		o.update.args = append(o.update.args, o.whereArgs...)
	}
	if len(o.limit) > 0 && len(o.limitArgs) > 0 {
		query = fmt.Sprintf("%s %s", query, o.limit)
		o.update.args = append(o.update.args, o.limitArgs[0])
	}
	o.update.queryString = query
	return true
}

func (o *Orm) parsePk(args ...interface{}) {
	if len(args) > 0 {
		o.pkValue = args[0]
		return
	}
	if o.tableName == "" {
		return
	}
	if o.pk == "" {
		pk, ok := pkMap[o.tableName]
		if !ok {
			return
		}
		o.pk = pk
	}

	if o.structValue.IsValid() && !o.structValue.IsZero() {
		f := o.structValue.FieldByName(o.pk)
		if !f.IsValid() || f.IsZero() {
			return
		}
		o.pkValue = f.Interface()
	}
}

func (o *Orm) parseQuery() bool {
	if o.tableName == "" {
		return false
	}
	var query = fmt.Sprintf("select `%s` from `%s` ", strings.Join(o.query.cols, "`,`"), o.tableName)
	o.parsePkWhere()
	if len(o.wheres) > 0 {
		query = fmt.Sprintf("%s where %s", query, strings.Join(o.wheres, " AND "))
		o.query.args = append(o.query.args, o.whereArgs...)
	}

	if len(o.query.sorts) > 0 {
		query = fmt.Sprintf("%s order by %s", query, strings.Join(o.query.sorts, ","))
	}
	if len(o.limit) > 0 {
		query = fmt.Sprintf("%s %s", query, o.limit)
		o.query.args = append(o.query.args, o.limitArgs...)
	}
	o.query.queryString = query
	return true
}

func (o *Orm) getCountQuery() (string, []interface{}) {
	var query = fmt.Sprintf("select count(*) from `%s`", o.tableName)
	o.parsePkWhere()
	if len(o.wheres) > 0 {
		query = query + " where " + strings.Join(o.wheres, " AND ")
	}
	return query, o.whereArgs
}

func (o *Orm) parsePkWhere() {
	if (len(o.wheres) == 0 || o.usePk) && o.pkValue != nil {
		if value := reflect.ValueOf(o.pkValue); value.IsValid() && !value.IsZero() {
			o.wheres = append(o.wheres, fmt.Sprintf("%s = ?", toSnake(o.pk)))
			o.whereArgs = append(o.whereArgs, o.pkValue)
		}
	}
}
