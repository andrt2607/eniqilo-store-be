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

type CustomerService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newCustomerService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *CustomerService {
	return &CustomerService{repo, validator, cfg}
}

func (u *CustomerService) Register(ctx context.Context, body dto.ReqCustomerRegister) (int, string, interface{}, error) {
	result := dto.ResCustomerRegister{}

	err := u.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	customer := body.ToEntity(u.cfg.BCryptSalt)
	customerID, err := u.repo.Customer.Insert(ctx, customer)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusInternalServerError, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	result.CustomerID = customerID
	result.PhoneNumber = body.PhoneNumber
	result.Name = body.Name

	return http.StatusOK, global_constant.SUCCESS_REGISTER_USER, result, nil
}

func (u *CustomerService) GetCustomer(ctx context.Context, param dto.ReqParamCustomerGet, sub string) ([]dto.ResCustomerGet, error) {

	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Customer.GetCustomer(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// func (u *CustomerService) GetCustomerByID(ctx context.Context, id, sub string) (dto.ResCustomerGet, error) {
// 	res, err := u.repo.Customer.GetCustomerByID(ctx, id, sub)
// 	if err != nil {
// 		return res, err
// 	}

// 	return res, nil
// }

// func (u *CustomerService) Login(ctx context.Context, body dto.ReqCustomerLogin) (int, string, interface{}, error) {
// 	res := dto.ResCustomerLogin{}

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

// 	customer, err := u.repo.Customer.GetByPhoneNumber(ctx, body.PhoneNumber)
// 	if err != nil {
// 		ierr.LogErrorWithLocation(err)
// 		return http.StatusNotFound, global_constant.NOT_FOUND, err.Error(), err
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(body.Password)); err != nil {
// 		if err == bcrypt.ErrMismatchedHashAndPassword {
// 			ierr.LogErrorWithLocation(err)
// 			return http.StatusBadRequest, global_constant.WRONG_PASSWORD, err.Error(), err
// 		}
// 		return http.StatusBadRequest, global_constant.WRONG_PASSWORD, err.Error(), err
// 	}

// 	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: customer.ID})
// 	if err != nil {
// 		return http.StatusInternalServerError, global_constant.FAIL_GENERATE_TOKEN, err.Error(), err
// 	}

// 	res.PhoneNumber = customer.PhoneNumber
// 	// res.Email = customer.Email
// 	res.Name = customer.Name
// 	res.AccessToken = token

// 	return http.StatusOK, global_constant.SUCCESS_LOGIN_USER, res, nil
// }
