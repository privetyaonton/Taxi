package model

const (
	StatusCreated string = "created"
	StatusDeleted string = "deleted"
)

type User struct {
	ID          uint64  `json:"-"`
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Email       string  `json:"email"`
	Raiting     float64 `json:"raiting"`
	Status      string  `json:"-"`
}
