package handler

import (
	"auth/controller"
	"auth/controller/req"
	"auth/jwt"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func NewHandler() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/users", createUser).Methods(http.MethodPost)
	m.HandleFunc("/users/verify", verifyToken).Methods(http.MethodPost)
	return m
}

func getController() controller.UserController {
	return controller.NewController()
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
	result, err := c.CreateUser(r.Context(), *reqDto)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(string(result)))
}

var p = jwt.NewProvider("hello_world_this_is_secretKey")

func verifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	fmt.Println("token::")
	p.GetPayLoad(token)
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no token"))
		return
	}
	result, err := p.ValidToken(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintln("", result)))
}
