package storage

import (
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserStorage_Save(t *testing.T) {
	db := UserConnection(t)
	ut := NewUserStorage(db)

	user := User{
		Name:     uuid.NewString(),
		Login:    uuid.NewString(),
		Password: uuid.NewString(),
	}
	id, err := ut.Save(user)
	require.NoError(t, err)

	dbUser, err := ut.GetUserByID(id)
	require.NoError(t, err)
	require.NotEmpty(t, dbUser)
}
func TestUserStorage_GetUserByID_Error(t *testing.T) {
	db := UserConnection(t)
	ut := NewUserStorage(db)
	_, err := ut.GetUserByID(123)
	require.Error(t, err)
}

func UserConnection(t *testing.T) *sql.DB {
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
