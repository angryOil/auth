package main

import (
	"auth/cmd/app/handler"
	"auth/controller"
	"auth/jwt"
	"auth/repository"
	"auth/repository/infla"
	"auth/service"
	"net/http"
)

func main() {
	c := getController()
	h := handler.NewHandler(c)
	http.ListenAndServe(":8081", h)
}

func getController() controller.IController {
	return controller.NewController(service.NewUserService(repository.NewRepository(infla.NewDB()), jwt.NewProvider("this_is_secretKey")))
}
