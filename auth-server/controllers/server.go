package controllers

import (
	mod "auth-server/models"
	"fmt"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

// manager
var manager manage.Manager

// auth server
var srv server.Server

// 测试创建一个用户
var user = mod.GenerateTestUser()

// 创建一个临时存储user信息的map
var UserMap = map[string]mod.User{
	fmt.Sprintf("%d", user.ID): *user,
}

func init() {
	clientStore := store.NewClientStore()
	clientStore.Set("xys", &models.Client{
		ID:     "xys",
		Secret: "xxx",
		Domain: "xys.com",
	})

	manager = *manage.NewDefaultManager()

	// TODO
	// 检测 redirect_url 和 client 的 Domain，待定
	manager.SetValidateURIHandler(func(baseURI, redirectURI string) error {
		return nil
	})

	// 默认设置超时时间 token两个小时有效
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// 选择token的存储方式，这里使用 内存存储, 也可以是用redis/mysql 等
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// 设置授权的client内容信息
	// manager包含client信息
	manager.MapClientStorage(clientStore)

	// server
	srv = *server.NewServer(server.NewConfig(), &manager)

	// 根据 client id 从 manager 中获取 client info, 在获取 access token 校验过程中会被用到
	srv.SetClientInfoHandler(server.ClientFormHandler)

	//  设置为 authorization code 模式, 可以穿多个，
	//  目前之前code授权码模式和refresh token模式
	srv.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.Refreshing)

	// authorization code 模式,  第一步获取code,然后再用code换取 access token, 而不是直接获取 access token
	srv.SetAllowedResponseType(oauth2.Code)

	// 允许使用GET方法，请求授权
	srv.SetAllowGetAccessRequest(true)

	// 校验请求用户的handler
	srv.SetUserAuthorizationHandler(userAuthorizationHandler)

	// 校验请求client_id 和 username，password合法性
	srv.SetPasswordAuthorizationHandler(PasswordAuthorizationHandler)
}
