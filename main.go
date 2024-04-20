package main

import (
	"fmt"
	"net/http"
	"oauth2-go/server"
)

func main() {
	server.Init()

	// auth_server 授权入口
	http.HandleFunc("/oauth2/authorize", server.AuthorizeHandler)

	// auth_server 发现未登录状态, 跳转到的登录handler
	http.HandleFunc("/oauth2/login", server.LoginHandler)

	// auth_server拿到 client以后重定向到的地址, 也就是 auth_client 获取到了code, 准备用code换取accesstoken
	// http.HandleFunc("/oauth2/code_to_token", server.CodeToToken)

	// auth_server 处理由code 换取access token 的handler
	http.HandleFunc("/oauth2/token", server.TokenHandler)

	// 登录完成, 同意授权的页面
	http.HandleFunc("/oauth2/agree-auth", server.AgreeAuthHandler)

	// access token 换取用户信息的handler
	http.HandleFunc("/oauth2/getuserinfo", server.GetUserInfoHandler)

	http.Handle("/", http.FileServer(http.Dir("./static"))) //http://localhost:9000/juejin.html

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(":9000", nil)
	}()
	err := <-errChan
	if err != nil {
		fmt.Println("Hello server stop running.")
	}

}



// http://localhost:8000/?redirect_uri=http://localhost:8001/auth/token&response_type=code&client_id=xys&scope=read