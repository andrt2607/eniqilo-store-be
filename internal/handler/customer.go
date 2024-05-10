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

type customerHandler struct {
	customerSvc *service.CustomerService
}

func newCustomerHandler(customerSvc *service.CustomerService) *customerHandler {
	return &customerHandler{customerSvc}
}

func (h *customerHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqCustomerRegister

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, global_constant.FAILED_PARSE_REQ_BODY)
		return
	}
	statusCode, message, res, err := h.customerSvc.Register(r.Context(), req)
	if err != nil {
		response.RespondWithError(w, statusCode, res)
		return
	}
	response.RespondWithJSON(w, http.StatusCreated, message, res)
}

func (h *customerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ReqParamCustomerGet

	fmt.Println(queryParams)

	param.PhoneNumber = queryParams.Get("phoneNumber")
	param.Name = queryParams.Get("name")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}
	fmt.Println(token.Subject())

	customers, err := h.customerSvc.GetCustomer(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = customers

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}

// func (h *customerHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("id")
// 	token, _, err := jwtauth.FromContext(r.Context())
// 	if err != nil {
// 		http.Error(w, "failed to get token from request", http.StatusBadRequest)
// 		return
// 	}

// 	customer, err := h.customerSvc.GetCustomerByID(r.Context(), id, token.Subject())
// 	if err != nil {
// 		code, msg := ierr.TranslateError(err)
// 		http.Error(w, msg, code)
// 		return
// 	}

// 	successRes := response.SuccessReponse{}
// 	successRes.Message = "success"
// 	successRes.Data = customer

// 	json.NewEncoder(w).Encode(successRes)
// 	w.WriteHeader(http.StatusOK)
// }

// func (h *customerHandler) Login(w http.ResponseWriter, r *http.Request) {
// 	var req dto.ReqCustomerLogin

// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		response.RespondWithError(w, http.StatusBadRequest, global_constant.FAILED_PARSE_REQ_BODY)
// 		return
// 	}

// 	statusCode, message, res, err := h.customerSvc.Login(r.Context(), req)
// 	if err != nil {
// 		response.RespondWithError(w, statusCode, res)
// 		return
// 	}
// 	response.RespondWithJSON(w, http.StatusCreated, message, res)
// }
