package main

/*
	This is just an quick setup of the main file
*/

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/IanVinkHub/gojapi/app/jio"
	"github.com/IanVinkHub/gojapi/controllers"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func fallback(w http.ResponseWriter, r *http.Request) {
	jio.Error(w, errors.New("Endpoint not found"), 404)
}

func channelHandler() {
	for {
		select {}
	}
}

func main() {
	http.HandleFunc("/", fallback)
	http.HandleFunc("/auth/login", controllers.Login)
	http.HandleFunc("/auth/logout", controllers.Logout)
	http.HandleFunc("/auth/register", controllers.Register)
	http.HandleFunc("/auth/changepassword", controllers.Changepassword)

	http.HandleFunc("/auth/token", controllers.GetTokenInfo)
	http.HandleFunc("/auth/refresh", controllers.RefreshTokenPost)

	go channelHandler()
	log.Fatal(http.ListenAndServe(*addr, nil))
}
