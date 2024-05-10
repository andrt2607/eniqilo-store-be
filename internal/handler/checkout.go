package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"eniqilo-store-be/internal/dto"
	"eniqilo-store-be/internal/ierr"
	"eniqilo-store-be/internal/service"
	global_constant "eniqilo-store-be/pkg/constant"
	response "eniqilo-store-be/pkg/resp"

	"github.com/go-chi/jwtauth/v5"
)

type checkoutHandler struct {
	checkoutSvc *service.CheckoutService
}

func newCheckoutHandler(checkoutSvc *service.CheckoutService) *checkoutHandler {
	return &checkoutHandler{checkoutSvc}
}

func (h *checkoutHandler) PostCheckout(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqCheckoutPost

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, global_constant.FAILED_PARSE_REQ_BODY)
		return
	}
	statusCode, message, res, err := h.checkoutSvc.Register(r.Context(), req)
	if err != nil {
		response.RespondWithError(w, statusCode, res)
		return
	}
	response.RespondWithJSON(w, http.StatusCreated, message, res)
}

func (h *checkoutHandler) GetCheckout(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ReqParamCheckoutGet

	fmt.Println(queryParams)

	param.PhoneNumber = queryParams.Get("phoneNumber")
	param.Name = queryParams.Get("name")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}
	fmt.Println(token.Subject())

	checkouts, err := h.checkoutSvc.GetCheckout(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = checkouts

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}

// func (h *checkoutHandler) GetCheckoutByID(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("id")
// 	token, _, err := jwtauth.FromContext(r.Context())
// 	if err != nil {
// 		http.Error(w, "failed to get token from request", http.StatusBadRequest)
// 		return
// 	}

// 	checkout, err := h.checkoutSvc.GetCheckoutByID(r.Context(), id, token.Subject())
// 	if err != nil {
// 		code, msg := ierr.TranslateError(err)
// 		http.Error(w, msg, code)
// 		return
// 	}

// 	successRes := response.SuccessReponse{}
// 	successRes.Message = "success"
// 	successRes.Data = checkout

// 	json.NewEncoder(w).Encode(successRes)
// 	w.WriteHeader(http.StatusOK)
// }

// func (h *checkoutHandler) Login(w http.ResponseWriter, r *http.Request) {
// 	var req dto.ReqCheckoutLogin

// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		response.RespondWithError(w, http.StatusBadRequest, global_constant.FAILED_PARSE_REQ_BODY)
// 		return
// 	}

// 	statusCode, message, res, err := h.checkoutSvc.Login(r.Context(), req)
// 	if err != nil {
// 		response.RespondWithError(w, statusCode, res)
// 		return
// 	}
// 	response.RespondWithJSON(w, http.StatusCreated, message, res)
// }
