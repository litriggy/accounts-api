package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type App struct {
	Host        string
	Port        int
	Debug       bool
	ReadTimeout time.Duration
	Memcached   string
	SecretKey   string
	Expiry      int
	Stage       string
}

type DBConf struct {
	Host string
	Port int
	User string
	Pass string
	DB   string
}

var app = &App{}
var dbConf = &DBConf{}
var rpcs map[string]string

func DBCfg() *DBConf {
	return dbConf
}
func AppCfg() *App {
	return app
}

func GetRPC(net string) (string, error) {
	switch net {
	case "pol":
		return os.Getenv("POL"), nil
	case "eth":
		return os.Getenv("ETH"), nil
	default:
		return "", errors.New("invalid net")
	}

}

func LoadApp() {
	fmt.Println("asdf", os.Getenv("PORT"))
	app.Host = os.Getenv("HOST")
	app.Port, _ = strconv.Atoi(os.Getenv("PORT"))
	app.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	timeOut, _ := strconv.Atoi(os.Getenv("READ_TIMEOUT"))
	app.ReadTimeout = time.Duration(timeOut) * time.Second
	app.SecretKey = os.Getenv("SECRET_KEY")
	app.Expiry, _ = strconv.Atoi(os.Getenv("EXPIRY"))
	app.Memcached = os.Getenv("MEMCACHED")

	app.Stage = os.Getenv("STAGE")

}
func LoadDb() {
	dbConf.Host = os.Getenv("DB_HOST")
	dbConf.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	dbConf.User = os.Getenv("DB_USER")
	dbConf.Pass = os.Getenv("DB_PASSWORD")
	dbConf.DB = os.Getenv("DB_DATABASE")
}
