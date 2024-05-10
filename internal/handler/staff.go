package handler

import (
	"encoding/json"
	"net/http"

	"eniqilo-store-be/internal/dto"
	"eniqilo-store-be/internal/service"
	global_constant "eniqilo-store-be/pkg/constant"
	response "eniqilo-store-be/pkg/resp"
)

type staffHandler struct {
	staffSvc *service.StaffService
}

func newStaffHandler(staffSvc *service.StaffService) *staffHandler {
	return &staffHandler{staffSvc}
}

func (h *staffHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqStaffRegister

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, global_constant.FAILED_PARSE_REQ_BODY)
		return
	}
	statusCode, message, res, err := h.staffSvc.Register(r.Context(), req)
	if err != nil {
		response.RespondWithError(w, statusCode, res)
		return
	}
	response.RespondWithJSON(w, http.StatusCreated, message, res)
}

func (h *staffHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqStaffLogin

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, global_constant.FAILED_PARSE_REQ_BODY)
		return
	}

	statusCode, message, res, err := h.staffSvc.Login(r.Context(), req)
	if err != nil {
		response.RespondWithError(w, statusCode, res)
		return
	}
	response.RespondWithJSON(w, http.StatusCreated, message, res)
}
