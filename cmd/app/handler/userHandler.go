package handler

import (
	"auth/controller"
	"auth/controller/req"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
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

// createUser godoc
// @Summary Request createUser
// @Description 유저를 생성합니다
// @Tags users
// @Param createUser body req.CreateDto true "Create user"
// @Accept  json
// @Produce  json
// @Success 201 {object} bool
// @Failure 400 "bad request"
// @Failure 409 "duplicate email"
// @Router /users [post]
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
		if strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// login godoc
// @Summary Request login
// @Description 로그인을합니다
// @Tags users
// @Param loginUser body req.LoginDto true "login User"
// @Param q query string false "name search by q" Format(email)
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Failure 400 "잘못된요청시"
// @Failure 401 "아이디 혹은 비밀번호가 다를경우"
// @Failure 500 "예상치못한 error"
// @Router /users/login [post]
func (uh UserHandler) login(w http.ResponseWriter, r *http.Request) {
	loginDto := &req.LoginDto{}
	err := json.NewDecoder(r.Body).Decode(loginDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	token, err := uh.c.Login(r.Context(), *loginDto)
	fmt.Println("loginDto?", loginDto)

	if err != nil {
		// id , password 가 없는경우
		if strings.Contains(err.Error(), "login fail") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("login fail password or email is not matched"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error sorry"))
		log.Println(err)
		return
	}
	w.Write([]byte(token))
}

// verifyToken godoc
// @Summary Request verifyToken
// @Description 토큰을 검증합니다
// @Tags users
// @Param loginUser body string true "verifyToken"
// @Accept  json
// @Produce  json
// @Success 200 {object} bool
// @Failure 400 "잘못된 토큰이없을경우"
// @Failure 401 "잘못된 토큰일경우"
// @Router /users/verify [post]
func (uh UserHandler) verifyToken(w http.ResponseWriter, r *http.Request) {
	//token := r.URL.Query().Get("token")
	readBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid token"))
		return
	}
	token := string(readBody)
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no token"))
		return
	}
	result, err := uh.c.Verify(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(string(err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%t", result)
}
