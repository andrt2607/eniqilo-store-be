package dto

import (
	"eniqilo-store-be/internal/entity"
	"eniqilo-store-be/pkg/auth"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type (
	ReqStaffRegister struct {
		PhoneNumber string `json:"phoneNumber" validate:"required"`
		Name        string `json:"name" validate:"required,min=5,max=50"`
		Password    string `json:"password" validate:"required,min=5,max=15"`
	}
	ResStaffRegister struct {
		StaffID     string `json:"userID,omitempty"`
		PhoneNumber string `json:"phoneNumber"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
	ReqStaffLogin struct {
		PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=16"`
		Password    string `json:"password" validate:"required,min=5,max=15"`
	}
	ResStaffLogin struct {
		StaffID     string `json:"userID,omitempty"`
		PhoneNumber string `json:"phoneNumber"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
)

func PhoneValidation(fl validator.FieldLevel) bool {
	// Pola regex untuk nomor telepon
	phonePattern := `^\+\d{1,3}(-\d+)?$`
	phoneNumber := fl.Field().String()
	matched, _ := regexp.MatchString(phonePattern, phoneNumber)
	return matched
}

func (d *ReqStaffRegister) ToEntity(cryptCost int) entity.Staff {
	return entity.Staff{Name: d.Name, Password: auth.HashPassword(d.Password, cryptCost), PhoneNumber: d.PhoneNumber}
}
