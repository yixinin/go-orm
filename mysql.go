package orm

import (
	"database/sql"
	"fmt"
	// _ "github.com/go-sql-driver/mysql"
)

//MysqlConfig ..
type MysqlConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DB       string

	MaxConn int `yaml:"max_conn"`
	MaxIdle int `yaml:"max_idle"`
}

func openMysql(cfg *MysqlConfig) (*sql.DB, error) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)
	db, err := sql.Open("mysql", source)
	if err == nil {
		db.SetMaxOpenConns(cfg.MaxConn)
		db.SetMaxIdleConns(cfg.MaxIdle)
		//err = db.Ping()
	}

	return db, err
}
