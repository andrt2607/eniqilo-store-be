package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	// _, _, err := jwtauth.FromContext(r.Context())
	// if err != nil {
	// 	http.Error(w, "failed to get token from request", http.StatusBadRequest)
	// 	return
	// }

	var req dto.ReqCheckoutPost

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, global_constant.FAILED_PARSE_REQ_BODY)
		return
	}

	statusCode, message, res, err := h.checkoutSvc.PostCheckout(r.Context(), req)
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

	param.CustomerId = queryParams.Get("customerId")
	param.CreatedAt = queryParams.Get("createdAt")
	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))

	_, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	checkouts, err := h.checkoutSvc.GetCheckout(r.Context(), param)
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
