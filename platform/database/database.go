package db

import (
	"database/sql"
	"fmt"

	"accounts/api/pkg/config"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	DB *sql.DB
}

func ConnectDB() (*Mysql, error) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DBCfg().User,
		config.DBCfg().Pass,
		config.DBCfg().Host,
		config.DBCfg().Port,
		config.DBCfg().DB,
	))
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Mysql{DB: db}, nil
}
