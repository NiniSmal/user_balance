package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"user_balance/storage"
)

type UserHandler struct {
	storage *storage.UserStorage
}

func NewUserHandler(s *storage.UserStorage) *UserHandler {
	return &UserHandler{
		storage: s,
	}
}

func HandlerError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.Write([]byte("The technical error"))
	return
}

func (u *UserHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	log.Println("User home page")
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user storage.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(w, err)
		return
	}
	_, err = u.storage.Save(user)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idUs := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idUs)
	if err != nil {
		HandlerError(w, err)
	}
	user, err := u.storage.GetUserByID(int64(id))
	if err != nil {
		HandlerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
