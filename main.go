package main

import (
	"accounts/api/cmd/server"
	"accounts/api/pkg/config"
)

//	@title			Accounts API
//	@version		1.0
//	@description	API 서버 입니다. session 키 활용 데이터 확인, 업데이트 시 헤더 authorization 키 값으로 재발급된 세션키를 확인할 수 있습니다.
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//
//	@host		localhost:3000
//	@BasePath	/api
//	@schemes	http
func main() {
	// go chain.Test()
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	// <-signalChan
	// result, _ := service.GetUser(1)

	// fmt.Println(result)
	config.LoadConfigs(".env")
	server.Serve()
}
