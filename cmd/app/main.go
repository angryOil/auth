package main

import (
	"auth/cmd/app/handler"
	"auth/controller"
	handler2 "auth/deco/handler"
	"auth/jwt"
	"auth/repository"
	"auth/repository/infla"
	"auth/service"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	c := getController()
	h := handler.NewHandler(c)
	wrappedLogger := handler2.NewDecoHandler(h, handler2.Logger)

	r.PathPrefix("/users").Handler(wrappedLogger)

	http.ListenAndServe(":8081", r)
}

func getController() controller.IController {
	return controller.NewController(service.NewUserService(repository.NewRepository(infla.NewDB()), jwt.NewProvider("this_is_secretKey")))
}
