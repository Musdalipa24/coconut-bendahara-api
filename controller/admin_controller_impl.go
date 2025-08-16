package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type AdminController interface {
	SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	FindByUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateAdmin(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type adminControllerImpl struct {
	AdminService service.AdminService
}

func NewAdminController(adminService service.AdminService) AdminController {
	return adminControllerImpl{
		AdminService: adminService,
	}
}

// SignUp implements AdminController.
func (a adminControllerImpl) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminRequest := dto.AdminRequest{}
	util.ReadFromRequestBody(r, &adminRequest)

	responseDTO,code, err := a.AdminService.SignUp(r.Context(), adminRequest)
	if err != nil {
		helper.WriteJSONError(w, code, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, code, "registration successfully")
}

// SignIn implements AdminController.
func (a adminControllerImpl) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	loginRequest := dto.LoginRequest{}
	util.ReadFromRequestBody(r, &loginRequest)

	responseDTO, code, err := a.AdminService.SignIn(r.Context(), loginRequest)
	if err != nil {
		helper.WriteJSONError(w, code, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, code, "registration successfully")
}

// FindByNik implements AdminController.
func (a adminControllerImpl) FindByUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	responseDTO, code, err := a.AdminService.GetAdminByUsername(r.Context(), username)
	if err != nil {
		helper.WriteJSONError(w, code, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, code, "get admin successfully")
}

// UpdateAdmin implements AdminController.
func (a adminControllerImpl) UpdateAdmin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	adminRequest := dto.UpdateAdminRequest{}
	util.ReadFromRequestBody(r, &adminRequest)

	responseDTO, code, err := a.AdminService.UpdateAdmin(r.Context(), adminRequest, username)
	if err != nil {
		helper.WriteJSONError(w, code, err.Error())
		return
	}

	helper.WriteJSONSuccess(w, responseDTO, code, "update admin successfully")
}