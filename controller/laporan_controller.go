package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type LaporanKeuanganController interface {
	GetAllLaporan(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetLastBalance(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetTotalIncome(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetTotalExpenditure(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetLaporanByDateRange(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}
