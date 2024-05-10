package dto

type Category string

const (
	Clothing    Category = "Clothing"
	Accessories Category = "Accessories"
	Footwear    Category = "Footwear"
	Beverages   Category = "Beverages"
)

type Sort string

const (
	ASC  Sort = "ASC"
	DESC Sort = "DESC"
)

type (
	ReqCreateProduct struct {
		Name        string `json:"name" validate:"required,min=1,max=30"`
		SKU         string `json:"sku" validate:"required,min=1,max=30"`
		Category    string `json:"category" validate:"required,oneof=Clothing Accessories Footwear Beverages"`
		ImageURL    string `json:"imageUrl" validate:"required"`
		Notes       string `json:"notes" validate:"required,min=1,max=200"`
		Price       int    `json:"price" validate:"required,min=1"`
		Stock       int    `json:"stock" validate:"required,min=0,max=100000"`
		Location    string `json:"location" validate:"required,min=1,max=200"`
		IsAvailable bool   `json:"isAvailable" validate:"required"`
	}

	ResCreateProduct struct {
		ID        string `json:"id" validate:"omitempty"`
		CreatedAt string `json:"createdAt" validate:"omitempty,iso8601"`
	}

	ResGetProduct struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		SKU         string `json:"sku"`
		Category    string `json:"category"`
		ImageURL    string `json:"imageUrl"`
		Notes       string `json:"notes"`
		Price       int    `json:"price"`
		Stock       int    `json:"stock"`
		Location    string `json:"location"`
		IsAvailable bool   `json:"isAvailable"`
		CreatedAt   string `json:"createdAt"`
	}

	ReqParamProductGet struct {
		ID          string   `json:"id"`
		Limit       int      `json:"limit"`
		Offset      int      `json:"offset"`
		Name        string   `json:"name"`
		IsAvailable string   `json:"isAvailable"`
		Category    Category `json:"category"`
		Sku         string   `json:"sku"`
		Price       Sort     `json:"price"`
		InStock     string   `json:"inStock"`
		CreatedAt   Sort     `json:"createdAt"`
  }
	ReqParamProductSKUGet struct {
		Name      string `json:"phoneNumber"`
		SKU       string `json:"sku"`
		Price     string `json:"price"`
		Category  string `json:"category"`
		Stock     string `json:"inStock"`
		CreatedAt string `json:"createdAt"`
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
	}
	ResProductSKUGet struct {
		Id        string `json:"id,omitempty"`
		Name      string `json:"name"`
		SKU       string `json:"sku"`
		Category  string `json:"category"`
		ImageURL  string `json:"imageUrl"`
		Stock     int    `json:"stock"`
		Price     int    `json:"price"`
		Location  string `json:"location"`
		CreatedAt string `json:"createdAt"`

	}
)
