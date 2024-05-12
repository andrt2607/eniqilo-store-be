package service

import (
	"context"
	"fmt"
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

	fmt.Println(body)
	// //register custom validation phone
	errAmount := c.validator.RegisterValidation("isvalidamount", dto.IsValidAmount)
	if errAmount != nil {
		ierr.LogErrorWithLocation(errAmount)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, errAmount.Error(), errAmount
	}
	errInteger := c.validator.RegisterValidation("isvalidinteger", dto.IsValidInteger)
	if errInteger != nil {
		ierr.LogErrorWithLocation(errInteger)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, errInteger.Error(), errInteger
	}

	fmt.Println("masuk sini")
	err := c.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}
	fmt.Println("masuk bos")

	responseCodeValidation, err := c.repo.Checkout.PostValidateCheckout(ctx, body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return responseCodeValidation, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	responseCode, responseData, errPost := c.repo.Checkout.PostCheckout(ctx, body)
	if errPost != nil {
		ierr.LogErrorWithLocation(errPost)
		return responseCode, global_constant.FAIL_VALIDATE_REQ_BODY, errPost.Error(), errPost
	}

	return responseCode, global_constant.SUCCESS, responseData, nil
}

func (u *CheckoutService) GetCheckout(ctx context.Context, param dto.ReqParamCheckoutGet) ([]dto.ResCheckoutGet, error) {

	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Checkout.GetCheckout(ctx, param)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return nil, err
	}

	return res, nil
}
