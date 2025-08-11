package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AdminController interface {
	SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	FindByNik(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}