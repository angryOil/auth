package handler

import (
	"auth/controller"
	"auth/controller/req"
	"auth/repository"
	"auth/repository/infla"
	"auth/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func NewHandler() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/users", createUser).Methods(http.MethodPost)
	m.HandleFunc("/users/verify", verifyToken).Methods(http.MethodPost)
	m.HandleFunc("/users/login", login).Methods(http.MethodPost)
	return m
}

// token 으로 응답을 할것이므로 [string] type
func getController() controller.UserController {
	return controller.NewController(service.NewUserService(repository.NewRepository(infla.NewDB())))
}

var c = getController()

func createUser(w http.ResponseWriter, r *http.Request) {
	reqDto := &req.CreateDto{}

	err := json.NewDecoder(r.Body).Decode(reqDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = c.CreateUser(r.Context(), *reqDto)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func login(w http.ResponseWriter, r *http.Request) {
	loginDto := &req.LoginDto{}
	err := json.NewDecoder(r.Body).Decode(&loginDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	token, err := c.Login(r.Context(), *loginDto)
	fmt.Println("loginDto?", loginDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(token))

}

func verifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no token"))
		return
	}
	result, err := c.Verify(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{result:%b}", result)
}
