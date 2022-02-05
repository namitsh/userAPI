package router

import (
	controller "UserMicroservice/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/signup", controller.SignUpHandler).Methods(http.MethodPost)
	r.HandleFunc("/login", controller.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/user/{id}", controller.GetUserHandler).Methods(http.MethodGet)
	r.HandleFunc("/signup", controller.UpdateUserHandler).Methods(http.MethodPatch)
	return r
}
