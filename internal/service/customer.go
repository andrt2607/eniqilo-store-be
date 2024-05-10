package service

import (
	"context"
	"errors"
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

	//register custom validation phone
	errPhone := u.validator.RegisterValidation("phone", dto.PhoneValidation)
	if errPhone != nil {
		ierr.LogErrorWithLocation(errPhone)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, errPhone.Error(), errPhone
	}

	err := u.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	customer := body.ToEntity(u.cfg.BCryptSalt)
	customerID, err := u.repo.Customer.Insert(ctx, customer)
	if err != nil {
		existPhoneError := errors.New("phone number already exist")
		ierr.LogErrorWithLocation(existPhoneError)
		return http.StatusConflict, global_constant.EXISTING_PHONE_NUMBER, existPhoneError.Error(), existPhoneError
	}

	result.CustomerID = customerID
	result.PhoneNumber = body.PhoneNumber
	result.Name = body.Name

	return http.StatusCreated, global_constant.SUCCESS_REGISTER_USER, result, nil
}

func (u *CustomerService) GetCustomer(ctx context.Context, param dto.ReqParamCustomerGet) ([]dto.ResCustomerGet, error) {

	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Customer.GetCustomer(ctx, param)
	if err != nil {
		return nil, err
	}

	return res, nil
}
