package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PemasukanController interface {
	AddPemasukan(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	UpdatePemasukan(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetPemasukan(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	DeletePemasukan(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}