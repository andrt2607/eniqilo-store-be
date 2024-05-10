package dto

import (
	"eniqilo-store-be/internal/entity"
)

type (
	ReqCustomerRegister struct {
		PhoneNumber string `json:"phoneNumber" validate:"required"`
		Name        string `json:"name" validate:"required,min=5,max=50"`
	}
	ResCustomerRegister struct {
		CustomerID  string `json:"userId,omitempty"`
		PhoneNumber string `json:"phoneNumber"`
		Name        string `json:"name"`
	}
	ReqParamCustomerGet struct {
		PhoneNumber string `json:"phoneNumber"`
		Name        string `json:"name"`
	}
	ResCustomerGet struct {
		CustomerID  string `json:"userId,omitempty"`
		PhoneNumber string `json:"phoneNumber"`
		Name        string `json:"name"`
	}
)

func (d *ReqCustomerRegister) ToEntity(cryptCost int) entity.Customer {
	return entity.Customer{PhoneNumber: d.PhoneNumber, Name: d.Name}
}
