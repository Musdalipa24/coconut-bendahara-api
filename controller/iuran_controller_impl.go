package controller

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type IuranController interface {
	CreateMember(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetAllMembers(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetMemberById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateIuran(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type IuranControllerImpl struct {
	IuranService service.IuranService
}

func NewIuranController(s service.IuranService) IuranController {
	return &IuranControllerImpl{IuranService: s}
}

// CreateMember implements IuranController.
func (i *IuranControllerImpl) CreateMember(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	memberReq := dto.MemberRequest{}
	err := util.ReadFromRequestBody(r, &memberReq)
	if err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("failed to decode request body (%v)", err))
		return
	}

	responseDTO, code, err := i.IuranService.CreateMember(r.Context(), memberReq)
	if err != nil {
		helper.WriteJSONError(w, code, err.Error())
		return
	}

	helper.WriteJSONSuccess(w, responseDTO, code, "successfully created member")
}

// GetAllMembers implements MemberController.
func (i *IuranControllerImpl) GetAllMembers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	responseDTO, code, err := i.IuranService.GetAllMember(r.Context())
	if err != nil {
		// helper.WriteJSONError(w, code, err.Error())
		// return
		panic(err)
	}

	helper.WriteJSONSuccess(w, responseDTO, code, "successfully retrieved all members")
}

// GetMemberById implements MemberController.
func (i *IuranControllerImpl) GetMemberById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	responseDTO, code, err := i.IuranService.GetMemberById(r.Context(), id)
	if err != nil {
		helper.WriteJSONError(w, code, err.Error())
		return
	}

	helper.WriteJSONSuccess(w, responseDTO, code, "successfully retrieved member by ID")
}

// UpdateIuran implements IuranController.
func (i *IuranControllerImpl) UpdateIuran(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id_member")
	pembayaranReq := dto.IuranRequest{}
	err := util.ReadFromRequestBody(r, &pembayaranReq)
	if err != nil {
		helper.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("failed to decode request body (%v)", err))
		return
	}

	responseDTO, code, err := i.IuranService.UpdateIuran(r.Context(), pembayaranReq, id)
	if err != nil {
		helper.WriteJSONError(w, code, err.Error())
		return
	}

	helper.WriteJSONSuccess(w, responseDTO, code, "successfully updated iuran")
}

// DeleteMember implements IuranController.
func (i *IuranControllerImpl) DeleteMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id_member")

	code, err := i.IuranService.DeleteMember(r.Context(), id)
	if err != nil {
		helper.WriteJSONError(w, code, err.Error())
		return
	}

	helper.WriteJSONSuccess(w, nil, code, "successfully deleted member")
}