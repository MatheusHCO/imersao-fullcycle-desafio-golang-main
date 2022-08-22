package bank

import uuid "github.com/satori/go.uuid"

type Account struct {
	ID      string  `json:"id"`
	Number  string  `json:"account_number"`
	Balance float64 `json:"-"`
}

func NewAccount() *Account {
	a := &Account{}
	a.ID = uuid.NewV4().String()
	a.Balance = 1000
	return a
}
