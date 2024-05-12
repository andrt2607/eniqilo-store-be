package entity

type CheckoutDetail struct {
	ProductId string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity"  validate:"required,min=0"`
}

type Checkout struct {
	TransactionId  string           `json:"transactionId"`
	CustomerId     string           `json:"customerId"`
	ProductDetails []CheckoutDetail `json:"productDetails"`
	Paid           *int             `json:"paid"`
	Change         *int             `json:"change"`
}
