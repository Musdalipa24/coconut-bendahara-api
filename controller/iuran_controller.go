package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type IuranController interface {
	AddIuran(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	UpdateIuran(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetAllIuran(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetIuranById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteIuran(w http.ResponseWriter, r *http.Request, ps httprouter.Params)

}
