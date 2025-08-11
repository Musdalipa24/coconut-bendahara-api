package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
)

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
	id, _ := strconv.ParseInt(idStr, 10, 64)
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
	err = c.IuranService.UpdateIuran(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Iuran berhasil diperbarui"})
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
	id, _ := strconv.ParseInt(idStr, 10, 64)
	data, err := c.IuranService.GetIuranById(id)
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
