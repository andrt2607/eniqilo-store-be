package entity

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Sku         string `json:"sku"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
	Notes       string `json:"notes"`
	Price       int    `json:"price"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"is_available"`
	UserID      string `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
