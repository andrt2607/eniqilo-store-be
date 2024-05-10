package entity

type CheckoutDetail struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type Checkout struct {
	TransactionId  string           `json:"transactionId"`
	CustomerId     string           `json:"customerId"`
	ProductDetails []CheckoutDetail `json:"productDetails"`
	Paid           int              `json:"paid"`
	Change         int              `json:"change"`
}
