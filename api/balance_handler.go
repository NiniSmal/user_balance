package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"user_balance/storage"
)

type BalanceHandler struct {
	storage *storage.BalanceStorage
}

func NewBalanceHandler(s *storage.BalanceStorage) *BalanceHandler {
	return &BalanceHandler{
		storage: s,
	}
}

func (b *BalanceHandler) CreateBalance(w http.ResponseWriter, r *http.Request) {
	var balance storage.Balance

	err := json.NewDecoder(r.Body).Decode(&balance)
	if err != nil {
		HandlerError(w, err)
		return
	}
	err = balance.Validate()
	if err != nil {
		HandlerError(w, err)
		return
	}

	balance.CreatedAt = time.Now()
	balance.UpdatedAt = balance.CreatedAt

	_, err = b.storage.Save(balance)
	if err != nil {
		HandlerError(w, err)
	}
}

func (b *BalanceHandler) GetBalanceByID(w http.ResponseWriter, r *http.Request) {
	idUs := r.URL.Query().Get("user_id")

	idUser, err := strconv.Atoi(idUs)
	if err != nil {
		HandlerError(w, err)
		return
	}

	balance, err := b.storage.GetBalanceByID(int64(idUser))
	if err != nil {
		HandlerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(balance)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (b *BalanceHandler) Replenishment(w http.ResponseWriter, r *http.Request) {
	var req storage.Replenishment

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandlerError(w, err)
		return
	}
	req.UpdatedAt = time.Now()

	_, err = b.storage.Update(req)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (b *BalanceHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var tr storage.Transfer

	err := json.NewDecoder(r.Body).Decode(&tr)
	if err != nil {
		HandlerError(w, err)
		return
	}
	err = tr.Validate()
	if err != nil {
		HandlerError(w, err)
		return
	}

	ctx := r.Context()
	tr.UpdatedAt = time.Now()

	err = b.storage.Transfer(tr, ctx)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
