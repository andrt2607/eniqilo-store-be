package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"eniqilo-store-be/internal/dto"
	"eniqilo-store-be/internal/ierr"
	"eniqilo-store-be/internal/service"
	global_constant "eniqilo-store-be/pkg/constant"
	response "eniqilo-store-be/pkg/resp"

	"github.com/go-chi/chi/v5"
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

func (h *productHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var req dto.ReqCreateProduct

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, global_constant.FAILED_PARSE_REQ_BODY)
		return
	}
	statusCode, message, res, err := h.productSvc.UpdateByID(r.Context(), id, req)
	if err != nil {
		response.RespondWithError(w, statusCode, res)
		return
	}
	response.RespondWithJSON(w, http.StatusOK, message, res)
}

func (h *productHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	_, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	var param dto.ReqParamProductGet

	param.ID = queryParams.Get("id")
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	param.Limit = limit
	offset, _ := strconv.Atoi(queryParams.Get("offset"))
	param.Offset = offset
	param.Name = queryParams.Get("name")
	param.IsAvailable = queryParams.Get("isAvailable")
	param.Category = dto.Category(queryParams.Get("category"))
	param.Sku = queryParams.Get("sku")
	param.Price = dto.Sort(queryParams.Get("price"))
	param.InStock = queryParams.Get("inStock")
	param.CreatedAt = dto.Sort(queryParams.Get("createdAt"))

	statusCode, message, res, err := h.productSvc.Get(r.Context(), param)
	if err != nil {
		response.RespondWithError(w, statusCode, res)
		return
	}
	response.RespondWithJSON(w, statusCode, message, res)
}

func (h *productHandler) GetProductSKU(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	var param dto.ReqParamProductSKUGet

	param.Name = queryParams.Get("name")
	param.SKU = queryParams.Get("sku")
	param.Stock = queryParams.Get("inStock")
	param.Price = dto.Sort(queryParams.Get("price"))
	param.Category = dto.Category(queryParams.Get("category"))
	param.CreatedAt = dto.Sort(queryParams.Get("createdAt"))
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

func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	_, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	statusCode, message, res, err := h.productSvc.DeleteProduct(r.Context(), id)
	if err != nil {
		response.RespondWithError(w, statusCode, res)
		return
	}
	response.RespondWithJSON(w, statusCode, message, res)
}
