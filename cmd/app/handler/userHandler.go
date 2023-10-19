package handler

import (
	"auth/controller"
	"auth/controller/req"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandler struct {
	c controller.IController
}

func NewHandler(c controller.IController) http.Handler {
	h := UserHandler{c: c}
	m := mux.NewRouter()
	m.HandleFunc("/users", h.createUser).Methods(http.MethodPost)
	m.HandleFunc("/users/verify", h.verifyToken).Methods(http.MethodPost)
	m.HandleFunc("/users/login", h.login).Methods(http.MethodPost)
	return m
}

func (uh UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	reqDto := &req.CreateDto{}

	err := json.NewDecoder(r.Body).Decode(reqDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = uh.c.CreateUser(r.Context(), *reqDto)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (uh UserHandler) login(w http.ResponseWriter, r *http.Request) {
	loginDto := &req.LoginDto{}
	err := json.NewDecoder(r.Body).Decode(&loginDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	token, err := uh.c.Login(r.Context(), *loginDto)
	fmt.Println("loginDto?", loginDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(token))

}

func (uh UserHandler) verifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no token"))
		return
	}
	result, err := uh.c.Verify(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{result:%b}", result)
}
