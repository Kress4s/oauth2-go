package controllers

import (
	"auth-server/models"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-session/session"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	//
	if err := srv.HandleAuthorizeRequest(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodPost {
		if r.Form == nil {
			r.ParseForm()
		}
		// 从登陆用户的form表单进来的
		uid, err := srv.PasswordAuthorizationHandler(r.Context(), r.Form.Get("client_id"), r.Form.Get("username"),
			r.Form.Get("password"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// 跳转到业务系统，让他请求 authorize 认证接口
		// clientInfo, err := manager.GetClient(r.Context(), r.Form.Get("client_id"))
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusUnauthorized)
		// 	return
		// }

		store.Set("LoggedInUserId", uid)
		store.Save()

		w.Header().Set("Location", `oauth2/authorize?redirect_uri=http://localhost:8001/auth/token&response_type=code&
		client_id=xys&scope=read`)
		w.WriteHeader(http.StatusFound)
	}
	// 若请求方法错误, 提供login.html页面
	outputHTML(w, r, "static/login.html")
}

func PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	// 这里可以判断 clientID 是否授权
	// username, password是否正确
	// 返回 user_id

	// 本次是测试单个client授权
	// if username == "xys" && password == "123456" && clientID == "xys-client" {
	if username == "xys" && password == "123456" {
		return "100", nil
	}
	return "", errors.New("client refuse or password is not right")
}

// 内部使用，用来查看是否登陆过,主要是取user_id 或者 跳转登陆页
func userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uid, ok := store.Get("LoggedInUserId")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		// 没有登陆，直接重定向到oauth-server的登陆页
		w.Header().Set("Location", "oauth2/login")
		w.WriteHeader(http.StatusFound)
		return "", nil
	}
	// 若有登陆过，直接返回user_id
	return uid.(string), nil
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	err := srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 提供 HTML 文件显示
func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 获取 access token
	access_token, ok := srv.BearerAuth(r)
	if !ok {
		log.Println("Failed to get access token from request")
		return
	}

	root_ctx := context.Background()
	ctx, cancle_func := context.WithTimeout(root_ctx, time.Second)
	defer cancle_func()

	// 从 access token 中获取 信息
	token_info, err := srv.Manager.LoadAccessToken(ctx, access_token)
	if err != nil {
		log.Println(err)
		return
	}

	// 获取 user id
	user_id := token_info.GetUserID()
	grant_scope := token_info.GetScope()

	user_info := models.User{}

	// 根据 grant scope 决定获取哪些用户信息
	if grant_scope != "read" {
		log.Println("invalid grant scope")
		w.Write([]byte("invalid grant scope"))
		return
	}

	user_info = UserMap[user_id]
	resp, err := json.Marshal(user_info)
	w.Write(resp)
	return
}
