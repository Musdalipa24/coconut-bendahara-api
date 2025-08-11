package controller

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type PengeluaranController interface {
	AddPengeluaran(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	UpdatePengeluaran(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetPengeluaran(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	DeletePengeluaran(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}