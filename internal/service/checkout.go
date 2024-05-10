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

type CheckoutService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newCheckoutService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *CheckoutService {
	return &CheckoutService{repo, validator, cfg}
}

func (c *CheckoutService) PostCheckout(ctx context.Context, body dto.ReqCheckoutPost) (int, string, interface{}, error) {

	err := c.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	responseCode, err := c.repo.Checkout.PostCheckout(ctx, body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return responseCode, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	return responseCode, global_constant.SUCCESS, "", nil
}

func (u *CheckoutService) GetCheckout(ctx context.Context, param dto.ReqParamCheckoutGet) ([]dto.ResCheckoutGet, error) {

	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Checkout.GetCheckout(ctx, param)
	if err != nil {
		return nil, err
	}

	return res, nil
}
