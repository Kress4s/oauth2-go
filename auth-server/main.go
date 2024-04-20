package main

import (
	"auth-server/controllers"
	"fmt"
	"net/http"
)

func main() {
	// auth-server 授权入口
	http.HandleFunc("/oauth2/authorize", controllers.AuthorizeHandler)

	// 跳转登陆页 或者 直接请求authorize 拿code返回给 redirect_url
	http.HandleFunc("/oauth2/login", controllers.LoginHandler)

	// code 换  access_token
	http.HandleFunc("/oauth2/access_token", controllers.LoginHandler)

	// 获取用户信息
	http.HandleFunc("/oauth2/userinfo", controllers.GetUserInfoHandler)

	// 开启文件代理，处理静态页面
	// http.Handle("/tmpfiles/", http.FileServer(http.Dir("./static")))
	http.Handle("/", http.FileServer(http.Dir("./static")))

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(":8000", nil)
	}()
	err := <-errChan
	if err != nil {
		fmt.Println("Hello server stop running.")
	}
}
