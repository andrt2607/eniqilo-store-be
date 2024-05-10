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

type productHandler struct {
	productSvc *service.ProductService
}

func newProductHandler(productSvc *service.ProductService) *productHandler {
	return &productHandler{productSvc}
}

func (h *productHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqCreateProduct

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, global_constant.FAILED_PARSE_REQ_BODY)
		return
	}
	statusCode, message, res, err := h.productSvc.Create(r.Context(), req)
	if err != nil {
		response.RespondWithError(w, statusCode, res)
		return
	}
	response.RespondWithJSON(w, http.StatusCreated, message, res)
}

func (h *productHandler) GetProductSKU(w http.ResponseWriter, r *http.Request) {

	fmt.Println("masuk sini")

	_, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	var param dto.ReqParamProductSKUGet

	param.Name = queryParams.Get("name")
	param.SKU = queryParams.Get("sku")
	param.Stock = queryParams.Get("inStock")
	param.Price = queryParams.Get("price")
	param.CreatedAt = queryParams.Get("createdAt")
	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))

	products, err := h.productSvc.GetProductSKU(r.Context(), param)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = products

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}
