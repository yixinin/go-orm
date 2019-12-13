package orm

import (
	"database/sql"
	"fmt"
	"reflect"
)

//Op op type
type Op uint8

const (
	//OpQuery ..
	OpQuery Op = 1
	//OpInsert ..
	OpInsert Op = 2
	//OpUpdate ..
	OpUpdate Op = 3
	//OpDelete ..
	OpDelete Op = 4
)

//Query ..
type Query struct {
	queryString string
	args        []interface{}

	colm map[string]Col
	cols []string

	// wheres    []string
	// whereArgs []interface{}

	sorts []string

	explain string
}

//Col ..
type Col struct {
	t reflect.Type
	i int
}

//Insert ..
type Insert struct {
	queryString string
	args        []interface{}
	explain     string
}

//Update ..
type Update struct {
	sets    []string
	setArgs []interface{}

	queryString string
	args        []interface{}

	explain string
}

//Delete ..
type Delete struct {
	queryString string
	args        []interface{}

	explain string
}

//Orm ..
type Orm struct {
	tableName string

	op Op

	mysql *sql.DB

	wheres    []string
	whereArgs []interface{}

	limit     string
	limitArgs []interface{}

	pk          string
	pkValue     interface{}
	structValue reflect.Value
	usePk       bool

	query  *Query
	insert *Insert
	update *Update
	delete *Delete
}

var mysqlDB *sql.DB

//TableNameble ..
type TableNameble interface {
	TableName() string
}

//Init ..
func Init(cfg *MysqlConfig) {
	var err error
	mysqlDB, err = openMysql(cfg)
	if err != nil {
		panic(fmt.Errorf("conn mysql error:%v", err))
	}
}

//NewOrm ..
func NewOrm() *Orm {
	return &Orm{
		mysql:  mysqlDB,
		query:  &Query{},
		insert: &Insert{},
		update: &Update{},
		delete: &Delete{},
	}
}

//Table ..
func (o *Orm) Table(name string) *Orm {
	o.tableName = name
	return o
}

//DB ...
func (o *Orm) DB() *sql.DB {
	return o.mysql
}
