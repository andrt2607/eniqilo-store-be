package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"eniqilo-store-be/internal/dto"
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
	// category := ToCategory()
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

func ToCategory(category string) string {
	switch category {
	case string(dto.Clothing):
		return string(dto.Clothing)
	case string(dto.Accessories):
		return string(dto.Accessories)
	case string(dto.Footwear):
		return string(dto.Footwear)
	case string(dto.Beverages):
		return string(dto.Beverages)
	default:
		return ""
	}
}
