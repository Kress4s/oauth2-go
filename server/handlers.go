package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-session/session"
)

// 授权入口, juejin.html 和 agree-auth.html 按下 button 后
func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	err := srv.HandleAuthorizeRequest(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// AuthorizeHandler 内部使用, 用于查看是否有登陆状态
func userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (user_id string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO
	// 判断clientID和Secret是否合法

	uid, ok := store.Get("LoggedInUserId")
	// 如果没有查询到登陆状态, 则跳转到 登陆页面
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		w.Header().Set("Location", "/oauth2/login")
		w.WriteHeader(http.StatusFound)
		return "", nil
	}
	// 若有登录状态, 返回 user id
	user_id = uid.(string)
	return user_id, nil
}

// 登录页面的handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		user_id, err := srv.PasswordAuthorizationHandler(context.TODO(), r.Form.Get("client_id"), r.Form.Get("username"),
			r.Form.Get("password"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		store.Set("LoggedInUserId", user_id) // 保存登录状态
		store.Save()

		// 跳转到 同意授权页面
		w.Header().Set("Location", "/oauth2/agree-auth")
		w.WriteHeader(http.StatusFound)
		return
	}

	// 若请求方法错误, 提供login.html页面
	outputHTML(w, r, "static/login.html")
}

// 若发现登录状态则提供 agree-auth.html, 否则跳转到 登陆页面
func AgreeAuthHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 如果没有查询到登陆状态, 则跳转到 登陆页面
	if _, ok := store.Get("LoggedInUserId"); !ok {
		w.Header().Set("Location", "/oauth2/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	// 如果有登陆状态, 会跳转到 确认授权页面
	outputHTML(w, r, "static/agree-auth.html")
}

// code 换取 access token
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	err := srv.HandleTokenRequest(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// access token 换取用户信息
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

	user_info := UserInfo{}

	// 根据 grant scope 决定获取哪些用户信息
	if grant_scope != "read_user_info" {
		log.Println("invalid grant scope")
		w.Write([]byte("invalid grant scope"))
		return
	}

	user_info = user_info_map[user_id]
	resp, err := json.Marshal(user_info)
	w.Write(resp)
	return
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
