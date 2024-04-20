package server

import (
	"context"
	"errors"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

var manager *manage.Manager
var srv *server.Server

// 用户信息结构体
type UserInfo struct {
	Username string `json:"username"`
	Gender   string `json:"gender"`
}

// 用一个 map 存储用户信息
var user_info_map = make(map[string]UserInfo)

func Init() {
	// 设置client信息
	client_store := store.NewClientStore()
	client_store.Set("juejin", &models.Client{
		ID:     "juejin",
		Secret: "xxxxxx",
		Domain: "http://xys.com",
	})

	// 设置manager信息， 参与校验 code/access token 请求
	manager = manage.NewDefaultManager()

	// 校验 redirect_uri 和 client 的 Domain, 简单起见, 不做校验
	manager.SetValidateURIHandler(func(baseURI, redirectURI string) error {
		return nil
	})

	// 默认设置超时时间
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token memory store, support redis、mysql... .etc
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// manger 包含 client 信息
	manager.MapClientStorage(client_store)

	// server 也包含 manager，client信息
	srv = server.NewServer(server.NewConfig(), manager)

	// 根据 client id 从 manager 中获取 client info, 在获取 access token 校验过程中会被用到
	srv.SetClientInfoHandler(server.ClientFormHandler)

	//  设置为 authorization code 模式
	srv.SetAllowedGrantType(oauth2.AuthorizationCode)

	// authorization code 模式,  第一步获取code,然后再用code换取 access token, 而不是直接获取 access token
	srv.SetAllowedResponseType(oauth2.Code)

	// 校验授权请求用户的handler, 会重定向到 登陆页面, 返回"", nil
	srv.SetUserAuthorizationHandler(userAuthorizationHandler)

	// 校验授权请求的用户的账号密码, 给 LoginHandler 使用, 简单起见, 只允许一个用户授权
	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string,
		err error) {
		if username == "Tom" && password == "123456" {
			return "0001", nil
		}
		return "", errors.New("username or password error")
	})

	// 允许使用 get 方法请求授权
	srv.SetAllowGetAccessRequest(true)

	// 储存用户信息的一个 map
	user_info_map["0001"] = UserInfo{
		"Tom", "Male",
	}
}
