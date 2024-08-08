package handler

import (
	"fmt"
	"net/http"

	"github.com/teddys48/kmpro/helper"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := http.NewServeMux()
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")

	server.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		helper.ReturnResponse(w, "Welcomeee!")
	})

	server.ServeHTTP(w, r)
}
