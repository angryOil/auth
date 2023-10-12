package main

import (
	"auth/cmd/app/handler"
	"net/http"
)

func main() {
	h := handler.NewHandler()
	http.ListenAndServe(":8081", h)
}
