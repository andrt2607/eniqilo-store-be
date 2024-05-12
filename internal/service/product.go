package service

import (
	"context"
	"net/http"

	"eniqilo-store-be/internal/cfg"
	"eniqilo-store-be/internal/dto"
	"eniqilo-store-be/internal/ierr"
	"eniqilo-store-be/internal/repo"
	global_constant "eniqilo-store-be/pkg/constant"

	"github.com/go-playground/validator/v10"
)

type ProductService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newProductService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *ProductService {
	return &ProductService{repo, validator, cfg}
}

func (p *ProductService) Create(ctx context.Context, body dto.ReqCreateProduct) (int, string, interface{}, error) {
	err := p.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	product, err := p.repo.Product.Insert(ctx, body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusInternalServerError, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	return http.StatusOK, global_constant.SUCCESS_CREATE_PRODUCT, product, nil
}

func (p *ProductService) UpdateByID(ctx context.Context, productId string, body dto.ReqCreateProduct) (int, string, interface{}, error) {
	err := p.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	product, err := p.repo.Product.UpdateByID(ctx, productId, body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		if err == ierr.ErrNotFound {
			return http.StatusNotFound, global_constant.NOT_FOUND, err.Error(), err
		}
		return http.StatusInternalServerError, global_constant.NOT_FOUND, err.Error(), err
	}

	return http.StatusOK, global_constant.SUCCESS_UPDATE_PRODUCT, product, nil
}

func (p *ProductService) Get(ctx context.Context, param dto.ReqParamProductGet) (int, string, interface{}, error) {
	err := p.validator.Struct(param)
	if err != nil {
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	products, err := p.repo.Product.Get(ctx, param)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusInternalServerError, global_constant.FAIL_GET_PRODUCT, err.Error(), err
	}

	return http.StatusOK, global_constant.SUCCESS, products, nil
}

func (p *ProductService) GetProductSKU(ctx context.Context, param dto.ReqParamProductSKUGet) ([]dto.ResProductSKUGet, error) {
	err := p.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := p.repo.Product.GetProductSKU(ctx, param)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (p *ProductService) DeleteProduct(ctx context.Context, id string) (int, string, interface{}, error) {

	err := p.repo.Product.DeleteProduct(ctx, id)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusNotFound, global_constant.FAIL_DELETE_PRODUCT, global_constant.FAIL_DELETE_PRODUCT, err
	}

	return http.StatusOK, global_constant.SUCCESS, global_constant.SUCCESS_DELETE_PRODUCT, nil
}
