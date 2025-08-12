package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
)

// Interface harus dideklarasikan sebelum struct
type IuranController interface {
	AddIuran(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateIuran(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetAllIuran(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetIuranById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteIuran(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}
type IuranControllerImpl struct {
	IuranService service.IuranService
}

func NewIuranController(s service.IuranService) IuranController {
	return &IuranControllerImpl{IuranService: s}
}

func (c *IuranControllerImpl) AddIuran(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req dto.IuranRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := c.IuranService.AddIuran(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Iuran berhasil ditambahkan", "id": id})
}

func (c *IuranControllerImpl) UpdateIuran(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")

	var req dto.IuranRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.IuranService.UpdateIuran(idStr, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Ambil data iuran yang sudah diupdate
	updated, err := c.IuranService.GetIuranById(idStr)
	if err != nil {
		http.Error(w, "Gagal mengambil data setelah update", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Iuran berhasil diperbarui",
		"data":    updated,
	})
}

func (c *IuranControllerImpl) GetAllIuran(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := c.IuranService.GetAllIuran()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (c *IuranControllerImpl) GetIuranById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	data, err := c.IuranService.GetIuranById(idStr)
	if err != nil {
		http.Error(w, "Data iuran tidak ditemukan", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (c *IuranControllerImpl) DeleteIuran(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	err := c.IuranService.DeleteIuran(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Iuran berhasil dihapus"})
}
