package main

import (
	"auth/cmd/app/handler"
	"auth/controller"
	_ "auth/docs"
	"auth/jwt"
	"auth/repository"
	"auth/repository/infla"
	"auth/service"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// @title Auth API
// @version 1.0
// @description This is a sample serice for managing todos
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @host localhost:8081
// @BasePath /

func main() {

	r := mux.NewRouter()
	// swagger
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	c := getController()
	h := handler.NewHandler(c)
	r.PathPrefix("/").Handler(h)
	http.ListenAndServe(":8081", r)
}

func getController() controller.IController {
	return controller.NewController(service.NewUserService(repository.NewRepository(infla.NewDB()), jwt.NewProvider("this_is_secretKey")))
}
