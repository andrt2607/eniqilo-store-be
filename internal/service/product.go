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

	return http.StatusOK, global_constant.SUCCESS_REGISTER_USER, product, nil
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
