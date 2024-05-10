package service

import (
	"context"
	"errors"
	"net/http"

	"eniqilo-store-be/internal/cfg"
	"eniqilo-store-be/internal/dto"
	"eniqilo-store-be/internal/ierr"
	"eniqilo-store-be/internal/repo"
	"eniqilo-store-be/pkg/auth"
	global_constant "eniqilo-store-be/pkg/constant"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type StaffService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newStaffService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *StaffService {
	return &StaffService{repo, validator, cfg}
}

func (u *StaffService) Register(ctx context.Context, body dto.ReqStaffRegister) (int, string, interface{}, error) {
	//check if phone number already exist
	checkedStaff, _ := u.repo.Staff.GetByPhoneNumber(ctx, body.PhoneNumber)
	if checkedStaff.PhoneNumber == body.PhoneNumber {
		existPhoneError := errors.New("phone number already exist")
		ierr.LogErrorWithLocation(existPhoneError)
		return http.StatusConflict, global_constant.EXISTING_PHONE_NUMBER, existPhoneError.Error(), existPhoneError
	}
	//register custom validation phone
	errPhone := u.validator.RegisterValidation("phone", dto.PhoneValidation)
	if errPhone != nil {
		ierr.LogErrorWithLocation(errPhone)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, errPhone.Error(), errPhone
	}
	//validate request body
	err := u.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	staff := body.ToEntity(u.cfg.BCryptSalt)
	staffID, err := u.repo.Staff.Insert(ctx, staff)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusInternalServerError, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}
	//generate token
	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: staffID})
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusInternalServerError, global_constant.INTERNAL_SERVER_ERROR, err.Error(), err
	}
	result := dto.ResStaffRegister{}
	result.StaffID = staffID
	result.PhoneNumber = body.PhoneNumber
	result.Name = body.Name
	result.AccessToken = token
	return http.StatusOK, global_constant.SUCCESS_REGISTER_USER, result, nil
}

func (u *StaffService) Login(ctx context.Context, body dto.ReqStaffLogin) (int, string, interface{}, error) {
	//validate request body
	err := u.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	//check if phone number not found
	staff, err := u.repo.Staff.GetByPhoneNumber(ctx, body.PhoneNumber)
	if err != nil {
		if err.Error() == "no rows in result set" {
			ierr.LogErrorWithLocation(err)
			return http.StatusNotFound, global_constant.NOT_FOUND, global_constant.PHONE_NUMBER_NOT_FOUND, err
		}
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			ierr.LogErrorWithLocation(err)
			return http.StatusBadRequest, global_constant.WRONG_PASSWORD, err.Error(), err
		}
		return http.StatusBadRequest, global_constant.WRONG_PASSWORD, err.Error(), err
	}
	//generate token
	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: staff.ID})
	if err != nil {
		return http.StatusInternalServerError, global_constant.FAIL_GENERATE_TOKEN, err.Error(), err
	}
	res := dto.ResStaffLogin{}
	res.PhoneNumber = staff.PhoneNumber
	res.Name = staff.Name
	res.AccessToken = token

	return http.StatusOK, global_constant.SUCCESS_LOGIN_USER, res, nil
}
