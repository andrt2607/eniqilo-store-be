package dto

import (
	"eniqilo-store-be/internal/entity"
)

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
		Limit      string `json:"limit"`
		Offset     string `json:"offset"`
		CreatedAt  string `json:"createdAt"`
	}
	ResCheckoutGet struct {
		TransactionId  string                  `json:"transactionId"`
		CustomerId     string                  `json:"customerId" validate:"required"`
		ProductDetails []entity.CheckoutDetail `json:"productDetails" validate:"required"`
		Paid           int                     `json:"paid" validate:"required,min=0"`
		Change         int                     `json:"change" validate:"required,min=0"`
	}
)

func (d *ReqCheckoutPost) ToEntity(transactionId string) entity.Checkout {
	return entity.Checkout{TransactionId: transactionId, CustomerId: d.CustomerId, ProductDetails: d.ProductDetails, Paid: d.Paid, Change: d.Change}
}
