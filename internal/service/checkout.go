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

func (u *CheckoutService) Register(ctx context.Context, body dto.ReqCheckoutRegister) (int, string, interface{}, error) {
	result := dto.ResCheckoutRegister{}

	err := u.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	checkout := body.ToEntity(u.cfg.BCryptSalt)
	checkoutID, err := u.repo.Checkout.Insert(ctx, checkout)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusInternalServerError, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	result.CheckoutID = checkoutID
	result.PhoneNumber = body.PhoneNumber
	result.Name = body.Name

	return http.StatusOK, global_constant.SUCCESS_REGISTER_USER, result, nil
}

func (u *CheckoutService) GetCheckout(ctx context.Context, param dto.ReqParamCheckoutGet, sub string) ([]dto.ResCheckoutGet, error) {

	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Checkout.GetCheckout(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// func (u *CheckoutService) GetCheckoutByID(ctx context.Context, id, sub string) (dto.ResCheckoutGet, error) {
// 	res, err := u.repo.Checkout.GetCheckoutByID(ctx, id, sub)
// 	if err != nil {
// 		return res, err
// 	}

// 	return res, nil
// }

// func (u *CheckoutService) Login(ctx context.Context, body dto.ReqCheckoutLogin) (int, string, interface{}, error) {
// 	res := dto.ResCheckoutLogin{}

// 	err := u.validator.Struct(body)
// 	if err != nil {
// 		ierr.LogErrorWithLocation(err)
// 		return http.StatusBadRequest, global_constant.INTERNAL_SERVER_ERROR, err.Error(), err
// 	}

// 	// perlu parse phone
// 	// _, err = mail.ParseAddress(body.Email)
// 	// if err != nil {
// 	// 	return res, ierr.ErrBadRequest
// 	// }

// 	checkout, err := u.repo.Checkout.GetByPhoneNumber(ctx, body.PhoneNumber)
// 	if err != nil {
// 		ierr.LogErrorWithLocation(err)
// 		return http.StatusNotFound, global_constant.NOT_FOUND, err.Error(), err
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(checkout.Password), []byte(body.Password)); err != nil {
// 		if err == bcrypt.ErrMismatchedHashAndPassword {
// 			ierr.LogErrorWithLocation(err)
// 			return http.StatusBadRequest, global_constant.WRONG_PASSWORD, err.Error(), err
// 		}
// 		return http.StatusBadRequest, global_constant.WRONG_PASSWORD, err.Error(), err
// 	}

// 	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: checkout.ID})
// 	if err != nil {
// 		return http.StatusInternalServerError, global_constant.FAIL_GENERATE_TOKEN, err.Error(), err
// 	}

// 	res.PhoneNumber = checkout.PhoneNumber
// 	// res.Email = checkout.Email
// 	res.Name = checkout.Name
// 	res.AccessToken = token

// 	return http.StatusOK, global_constant.SUCCESS_LOGIN_USER, res, nil
// }
