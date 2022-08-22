package bank

import uuid "github.com/satori/go.uuid"

type Transaction struct {
	ID     string
	From   string `json:"from"`
	To     string `json:"to"`
	Amount float64
}

func NewTransaction() *Transaction {
	t := &Transaction{}
	t.ID = uuid.NewV4().String()
	return t
}
