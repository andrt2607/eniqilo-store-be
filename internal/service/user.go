package service

import (
	"context"
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

type UserService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newUserService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *UserService {
	return &UserService{repo, validator, cfg}
}

func (u *UserService) Register(ctx context.Context, body dto.ReqRegister) (int, string, interface{}, error) {
	res := dto.ResRegister{}

	err := u.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	// perlu parse phone
	// _, err = mail.ParseAddress(body.Email)
	// if err != nil {
	// 	return res, ierr.ErrBadRequest
	// }

	user := body.ToEntity(u.cfg.BCryptSalt)
	userID, err := u.repo.User.Insert(ctx, user)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusInternalServerError, global_constant.FAIL_VALIDATE_REQ_BODY, err.Error(), err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: userID})
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusInternalServerError, global_constant.INTERNAL_SERVER_ERROR, err.Error(), err
	}

	// res.Email = body.Email
	res.PhoneNumber = body.PhoneNumber
	res.Name = body.Name
	res.AccessToken = token
	result := dto.ResRegister{}
	result.UserID = userID
	result.PhoneNumber = body.PhoneNumber
	result.Name = body.Name
	result.AccessToken = token
	return http.StatusOK, global_constant.SUCCESS_REGISTER_USER, result, nil
}

func (u *UserService) Login(ctx context.Context, body dto.ReqLogin) (int, string, interface{}, error) {
	res := dto.ResLogin{}

	err := u.validator.Struct(body)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusBadRequest, global_constant.INTERNAL_SERVER_ERROR, err.Error(), err
	}

	// perlu parse phone
	// _, err = mail.ParseAddress(body.Email)
	// if err != nil {
	// 	return res, ierr.ErrBadRequest
	// }

	user, err := u.repo.User.GetByPhoneNumber(ctx, body.PhoneNumber)
	if err != nil {
		ierr.LogErrorWithLocation(err)
		return http.StatusNotFound, global_constant.NOT_FOUND, err.Error(), err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			ierr.LogErrorWithLocation(err)
			return http.StatusBadRequest, global_constant.WRONG_PASSWORD, err.Error(), err
		}
		return http.StatusBadRequest, global_constant.WRONG_PASSWORD, err.Error(), err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: user.ID})
	if err != nil {
		return http.StatusInternalServerError, global_constant.FAIL_GENERATE_TOKEN, err.Error(), err
	}

	res.PhoneNumber = user.PhoneNumber
	// res.Email = user.Email
	res.Name = user.Name
	res.AccessToken = token

	return http.StatusOK, global_constant.SUCCESS_LOGIN_USER, res, nil
}

// func (u *UserService) UpdateAccount(ctx context.Context, body dto.ReqUpdateAccount, sub string) error {
// 	err := u.validator.Struct(body)
// 	if err != nil {
// 		return ierr.ErrBadRequest
// 	}

// 	if body.ImageURL == "http://incomplete" {
// 		return ierr.ErrBadRequest
// 	}

// 	err = u.repo.User.LookUp(ctx, sub)
// 	if err != nil {
// 		return err
// 	}

// 	err = u.repo.User.UpdateAccount(ctx, sub, body.Name, body.ImageURL)
// 	return err
// }
