package dto

import (
	"eniqilo-store-be/internal/entity"
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
		CreatedAt      string                  `json:"createdAt" default:"desc"`
	}
)

func (d *ReqCheckoutPost) ToEntity(transactionId string) entity.Checkout {
	return entity.Checkout{TransactionId: transactionId, CustomerId: d.CustomerId, ProductDetails: d.ProductDetails, Paid: d.Paid, Change: d.Change}
}
