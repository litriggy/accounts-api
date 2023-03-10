package config

import (
	"database/sql"
	"fmt"
	"log"

	db "accounts/api/platform/database/models"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*db.Queries, error) {
	var (
		database *sql.DB
		err      error
	)

	database, err = MysqlConnection()

	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}
	return db.New(database), nil
}

func MysqlConnection() (*sql.DB, error) {
	database, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		DBCfg().User,
		DBCfg().Pass,
		DBCfg().Host,
		DBCfg().Port,
		DBCfg().DB,
	))
	if err != nil {
		log.Fatalf("ERR(database:13) connection failed to mysql database: %v", err)
		return nil, err
	}
	if err := database.Ping(); err != nil {
		defer database.Close()
		log.Fatalf("ERR(database:18) DB connection failed to mysql: %v", err)
		return nil, err
	}
	return database, nil
}

func DBConn() (*db.Queries, *sql.DB, error) {
	var (
		database *sql.DB
		err      error
	)
	database, err = sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		DBCfg().User,
		DBCfg().Pass,
		DBCfg().Host,
		DBCfg().Port,
		DBCfg().DB,
	))
	if err != nil {
		log.Fatalf("ERR(database:13) connection failed to mysql database: %v", err)
		return nil, nil, err
	}
	if err := database.Ping(); err != nil {
		defer database.Close()
		log.Fatalf("ERR(database:18) DB connection failed to mysql: %v", err)
		return nil, nil, err
	}
	//return db.New(database), nil
	return db.New(database), database, nil
}
