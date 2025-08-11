package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type TransactionController interface {
	GetAllTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetLastTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}
