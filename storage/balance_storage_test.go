package storage

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBalanceStorage_Save(t *testing.T) {
	db := BalanceConnection(t)
	bt := NewBalanceStorage(db)

	balance := Balance{
		UserID:         1,
		AccountBalance: 10,
	}
	id, err := bt.Save(balance)
	require.NoError(t, err)
	dbBalance, err := bt.GetBalanceByID(id)
	require.NoError(t, err)
	require.NotEmpty(t, dbBalance)
}
func TestBalanceStorage_GetBalanceByID_Error(t *testing.T) {
	db := BalanceConnection(t)
	bt := NewBalanceStorage(db)
	_, err := bt.GetBalanceByID(1234)
	require.Error(t, err)
}

func TestBalanceStorage_Update(t *testing.T) {
	db := BalanceConnection(t)
	bt := NewBalanceStorage(db)

	balance := Balance{
		UserID: 11, AccountBalance: 20,
	}

	id, err := bt.Save(balance)
	require.NoError(t, err)

	updated_at := time.Now().Round(time.Millisecond)

	balance1 := Replenishment{
		ID:        id,
		UserID:    11,
		UpdatedAt: updated_at,
		Amount:    10,
	}
	balance.AccountBalance = balance.AccountBalance + balance1.Amount

	id1, err := bt.Update(balance1)
	require.NoError(t, err)

	dbBalance, err := bt.GetBalanceByID(id1)
	require.NoError(t, err)
	require.Equal(t, balance1.UserID, dbBalance.UserID)
	require.Equal(t, balance.AccountBalance, dbBalance.AccountBalance)
	require.Equal(t, balance1.UpdatedAt.Unix(), dbBalance.UpdatedAt.Unix())
}

func TestBalanceStorage_Transfer(t *testing.T) {
	db := BalanceConnection(t)
	bt := NewBalanceStorage(db)
	balance1 := Balance{UserID: 5, AccountBalance: 20}
	balance2 := Balance{UserID: 6, AccountBalance: 0}

	id1, err := bt.Save(balance1)
	require.NoError(t, err)
	id2, err := bt.Save(balance2)
	require.NoError(t, err)

	updated_at := time.Now().Round(time.Millisecond)

	tr := Transfer{
		Sender:    id1,
		Recipient: id2,
		Amount:    20,
		UpdatedAt: updated_at,
	}
	ctx := context.Background()

	balance1.AccountBalance = balance1.AccountBalance - tr.Amount
	balance2.AccountBalance = balance2.AccountBalance + tr.Amount
	err = bt.Transfer(tr, ctx)
	require.NoError(t, err)

	dbBalance1, err := bt.GetBalanceByID(id1)
	require.NoError(t, err)
	dbBalance2, err := bt.GetBalanceByID(id2)
	require.NoError(t, err)
	require.Equal(t, balance1.AccountBalance, dbBalance1.AccountBalance)
	require.Equal(t, balance2.AccountBalance, dbBalance2.AccountBalance)
}

//перевела денег
//поулчила балансы обоих юзеров  по очереди и проверила что денег норм

func BalanceConnection(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", "postgres://postgres:dev@localhost:8002/postgres?sslmode=disable")
	require.NoError(t, err)

	t.Cleanup(func() {
		err = db.Close()
		require.NoError(t, err)
	})
	err = db.Ping()
	require.NoError(t, err)
	return db
}
