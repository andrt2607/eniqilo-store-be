package dto

import (
	"eniqilo-store-be/internal/entity"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// type createdAt string

// const (
// 	Asc  createdAt = "asc"
// 	Desc createdAt = "desc"
// )

// ,dive,keys,eq=1|eq=2,endkeys,required,isvalidinteger"

type (
	ReqCheckoutPost struct {
		CustomerId     string                  `json:"customerId" validate:"required"`
		ProductDetails []entity.CheckoutDetail `json:"productDetails" validate:"required,dive"`
		Paid           *int                    `json:"paid" validate:"required"`
		Change         *int                    `json:"change" validate:"required"`
	}
	ResCheckoutPost struct {
		TransactionId string `json:"transactionId"`
	}
	ResValidateCheckoutPost struct {
		Stock  *int `json:"stock"`
		Charge *int `json:"change"`
	}
	ResCheckoutPos struct {
		ID        string `json:"id" validate:"omitempty"`
		CreatedAt string `json:"createdAt" validate:"omitempty,iso8601"`
	}

	ReqParamCheckoutGet struct {
		CustomerId string `json:"customerId"`
		Limit      int    `json:"limit"`
		Offset     int    `json:"offset"`
		CreatedAt  string `json:"createdAt" default:"desc"`
	}
	ResCheckoutGet struct {
		OrderId        string                  `json:"transactionId"`
		CustomerId     string                  `json:"customerId"`
		ProductDetails []entity.CheckoutDetail `json:"productDetails"`
		Paid           int                     `json:"paid" validate:"min=0"`
		Change         int                     `json:"change" validate:"min=0"`
	}
)

func (d *ReqCheckoutPost) ToEntity(transactionId string) entity.Checkout {
	return entity.Checkout{TransactionId: transactionId, CustomerId: d.CustomerId, ProductDetails: d.ProductDetails, Paid: d.Paid, Change: d.Change}
}

// func IsValidAmount(fl validator.FieldLevel) bool {
// 	id_input := fl.Field().Int()
// 	if id_input >= 0 {
// 		return true
// 	}
// 	return false
// }

func IsValidAmount(fl validator.FieldLevel) bool {
	fmt.Println(fl)
	// Get the field value
	fieldValue := fl.Field().Interface()
	fmt.Println(fieldValue)

	// Check if the field value is nil
	if fieldValue == nil {
		return false // Fail validation
	}

	// Check if the field value is a non-zero integer
	if intValue, ok := fieldValue.(int); ok && intValue >= 0 {
		return true // Pass validation
	}

	return false // Fail validation
}

func IsValidInteger(fl validator.FieldLevel) bool {

	id_input := fl.Field().String()

	_, err := strconv.Atoi(id_input)
	if err != nil {
		return false
	}
	return true
}
