package entity

type User struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}
