package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func AuthToken(w http.ResponseWriter, r *http.Request) {
	/*
		1. 接到code
		2. 拿code取换access_token
	*/
	// if r.Form == nil {
	// 	r.ParseForm()
	// }
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, errors.New("code is empty").Error(), http.StatusUnauthorized)
	}

	reqUrl := fmt.Sprintf("http://localhost:8000/oauth2/access_token?code=%s", code)

	var reader io.Reader
	req, _ := http.NewRequest(http.MethodGet, reqUrl, reader)

	ww, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer ww.Body.Close()
	content, err := io.ReadAll(ww.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(content)
}
