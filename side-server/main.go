package main

import (
	"fmt"
	"net/http"
	"oauth2-go/side-server/controllers"
)

func main() {
	http.HandleFunc("/auth/token", controllers.AuthToken)
	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(":8001", nil)
	}()
	err := <-errChan
	if err != nil {
		fmt.Println("Hello server stop running.")
	}
}
