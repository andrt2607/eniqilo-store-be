package entity

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
	Notes       string `json:"notes"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"is_available"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Category string
type CategoryProduct []Category

const (
	Clothing    Category = "Clothing"
	Accessories Category = "Accessories"
	Footwear    Category = "Footwear"
	Beverages   Category = "Beverages"
)
