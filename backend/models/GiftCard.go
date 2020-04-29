package models

type GiftCard struct {
	Model
	PersonName string `json:"Name"`
	Email      string
	Phone      string
	Summ       int
	Approved   bool //If we approved this request and received money - true
}
