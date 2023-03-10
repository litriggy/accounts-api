package memcached

import (
	"accounts/api/pkg/config"
	"fmt"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/bradfitz/gomemcache/memcache"
)

var mc *memcache.Client

func Init() {

	mc = memcache.New(config.AppCfg().Memcached)
	err := mc.Ping()
	if err != nil {
		panic(err)
	}
}

func CreateSession(userID string) (string, error) {
	id, err := nanoid.New()
	if err != nil {
		return "", err
	}
	if err := mc.Set(&memcache.Item{Key: id, Value: []byte(userID), Expiration: int32(60 * 60 * 12)}); err != nil {
		fmt.Println(err)
		return "", err
	}
	return id, nil
}

func GetSession(sessionKey string) (string, string, error) {
	userID, err := mc.Get(sessionKey)
	if err != nil {
		return "", "", err
	}

	if err := mc.Delete(sessionKey); err != nil {
		return "", "", err
	}
	newSessionKey, _ := CreateSession(string(userID.Value))

	return string(userID.Value), newSessionKey, nil
}
