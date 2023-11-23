package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Balance struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"id_user"`
	AccountBalance int64     `json:"account_balance"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Replenishment struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"id_user"`
	Amount    int64     `json:"amount"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Transfer struct {
	Sender    int64     `json:"sender"`
	Recipient int64     `json:"recipient"`
	Amount    int64     `json:"amount"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BalanceStorage struct {
	db *sql.DB
}

func NewBalanceStorage(b *sql.DB) *BalanceStorage {
	return &BalanceStorage{
		db: b,
	}
}

func (b *Balance) Validate() error {

	if b.AccountBalance < 0 {
		return errors.New("The balance can't be less  0")
	}

	return nil
}
func (t *Transfer) Validate() error {
	if t.Sender <= 0 || t.Recipient <= 0 {
		return errors.New("This Id is incorrect")
	}
	return nil
}

func (b *BalanceStorage) Save(balance Balance) (int64, error) {
	query := "INSERT INTO balances(id_user, account_balance, created_at, updated_at)VALUES($1,$2,$3,$4) RETURNING id_user"

	var id int64

	err := b.db.QueryRow(query, balance.UserID, balance.AccountBalance, balance.CreatedAt, balance.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return balance.UserID, nil
}

func (b *BalanceStorage) GetBalanceByID(userID int64) (Balance, error) {
	query := "SELECT id, id_user, account_balance, created_at, updated_at FROM balances WHERE id_user = $1 "

	var balance Balance

	err := b.db.QueryRow(query, userID).Scan(&balance.ID, &balance.UserID, &balance.AccountBalance, &balance.CreatedAt, &balance.UpdatedAt)
	if err != nil {
		return Balance{}, err
	}
	return balance, nil
}

func (b *BalanceStorage) Update(repl Replenishment) (int64, error) {
	query := "UPDATE balances SET account_balance = account_balance + $1, updated_at = $2 WHERE id_user = $3 RETURNING id_user"
	var id int64

	err := b.db.QueryRow(query, repl.Amount, repl.UpdatedAt, repl.UserID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (b *BalanceStorage) Transfer(tr Transfer, ctx context.Context) error {
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE balances SET account_balance = account_balance -  $1, updated_at = $2 WHERE id_user = $3", tr.Amount, tr.UpdatedAt, tr.Sender)
	if err != nil {
		return fmt.Errorf("update balance for user %d: %w", tr.Sender, err)
	}

	_, err = b.db.Exec("UPDATE balances SET account_balance = account_balance + $1, updated_at = $2 WHERE id_user = $3", tr.Amount, tr.UpdatedAt, tr.Recipient)
	if err != nil {
		return fmt.Errorf("update balance for user %d: %w", tr.Recipient, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
