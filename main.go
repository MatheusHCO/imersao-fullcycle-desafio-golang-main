package main

import (
	"database/sql"
	"desafio/bank"
	"errors"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	e := echo.New()
	e.POST("/bank-accounts", bankAccounts)
	e.POST("/bank-accounts/transfer", transfer)
	e.Logger.Fatal(e.Start(":8000"))
}

func bankAccounts(c echo.Context) error {
	account := bank.NewAccount()
	if err := c.Bind(account); err != nil {
		println(err.Error())
		return err
	}
	if account.Number == "" {
		return c.JSON(500, map[string]string{"error": "Missing parameters"})
	}

	_, err := loadAccount(account.Number)
	if err == nil {
		return c.JSON(500, map[string]string{"error": "Account already exists"})
	}

	err = saveAccount(*account)
	if err != nil {
		println(err.Error())
		return c.JSON(500, map[string]string{"error": "Fail"})
	}
	return c.JSON(201, account)
}

func saveAccount(account bank.Account) error {
	db, err := sql.Open("sqlite3", "bank.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO accounts (id, number, balance) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(account.ID, account.Number, account.Balance)
	if err != nil {
		return err
	}

	return nil
}

func loadAccount(accountNumber string) (bank.Account, error) {
	db, err := sql.Open("sqlite3", "bank.db")
	if err != nil {
		return bank.Account{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM accounts WHERE number = $1")
	if err != nil {
		return bank.Account{}, err
	}

	rows, err := stmt.Query(accountNumber)
	if err != nil {
		return bank.Account{}, err
	}

	defer rows.Close()

	var (
		id      string
		number  string
		balance float64
	)
	rows.Next()
	err = rows.Scan(&id, &number, &balance)
	if err != nil {
		return bank.Account{}, err
	}

	return bank.Account{
		ID:      id,
		Number:  number,
		Balance: balance,
	}, nil
}

func transfer(c echo.Context) error {

	transaction := bank.NewTransaction()
	if err := c.Bind(transaction); err != nil {
		println(err.Error())
		return err
	}

	if transaction.From == "" || transaction.To == "" {
		return c.JSON(500, map[string]string{"error": "Missing parameters"})
	}

	if transaction.From == transaction.To {
		return c.JSON(501, map[string]string{"error": "Cannot transfer to the same account"})
	}

	if transaction.Amount <= 0 {
		return c.JSON(502, map[string]string{"error": "Transaction amount must be more than 0"})
	}

	from, err := loadAccount(transaction.From)
	if err != nil {
		return c.JSON(503, map[string]string{"error": "Origin account not found"})
	}

	to, err := loadAccount(transaction.To)
	if err != nil {
		return c.JSON(504, map[string]string{"error": "Destination account not found"})
	}

	err = commitTransfer(&from, &to, transaction)
	if err != nil {
		return c.JSON(505, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, map[string]float64{
		"from_balance": from.Balance,
		"to_balance":   to.Balance,
	})
}

func commitTransfer(from *bank.Account, to *bank.Account, transaction *bank.Transaction) error {
	if from.Balance < transaction.Amount {
		return errors.New("Transfer amounts exceeds account balance")
	}

	db, err := sql.Open("sqlite3", "bank.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO transactions (id, `from`, `to`, amount) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(transaction.ID, transaction.From, transaction.To, transaction.Amount)
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("UPDATE accounts SET balance = $1 WHERE number = $2")
	if err != nil {
		return err
	}

	from.Balance = from.Balance - transaction.Amount
	_, err = stmt.Exec(from.Balance, from.Number)
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("UPDATE accounts SET balance = $1 WHERE number = $2")
	if err != nil {
		return err
	}

	to.Balance = to.Balance + transaction.Amount
	_, err = stmt.Exec(to.Balance, to.Number)
	if err != nil {
		return err
	}

	return nil
}
