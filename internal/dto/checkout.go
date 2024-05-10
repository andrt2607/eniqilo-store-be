package dto

import (
	"eniqilo-store-be/internal/entity"
)

// type createdAt string

// const (
// 	Asc  createdAt = "asc"
// 	Desc createdAt = "desc"
// )

type (
	ReqCheckoutPost struct {
		CustomerId     string                  `json:"customerId" validate:"required"`
		ProductDetails []entity.CheckoutDetail `json:"productDetails" validate:"required"`
		Paid           int                     `json:"paid" validate:"required,min=0"`
		Change         int                     `json:"change" validate:"required,min=0"`
	}
	ResCheckoutPost struct {
		TransactionId string `json:"transactionId"`
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
	ResPostValidateCheckout struct {
		Stock  int `json:"stock"`
		Charge int
	}
)

func (d *ReqCheckoutPost) ToEntity(transactionId string) entity.Checkout {
	return entity.Checkout{TransactionId: transactionId, CustomerId: d.CustomerId, ProductDetails: d.ProductDetails, Paid: d.Paid, Change: d.Change}
}
