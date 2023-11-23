package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"user_balance/api"
	"user_balance/storage"
)

func main() {
	//$ docker run -d -p 8002:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DATABASE=postgres postgres
	db, err := sql.Open("postgres", "postgres://postgres:dev@localhost:8002/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db - OK")

	us := storage.NewUserStorage(db)
	uh := api.NewUserHandler(us)

	bs := storage.NewBalanceStorage(db)
	bh := api.NewBalanceHandler(bs)

	router := http.NewServeMux()
	router.HandleFunc("/", uh.HomePage)
	router.HandleFunc("/users/create", uh.CreateUser)
	router.HandleFunc("/users/getOn", uh.GetUserByID)

	router.HandleFunc("/balances/create", bh.CreateBalance)
	router.HandleFunc("/balances/GetOne", bh.GetBalanceByID)
	router.HandleFunc("/balances/Replenishment", bh.Replenishment)
	router.HandleFunc("/balances/Trasfer", bh.Transfer)

	err = http.ListenAndServe(":8094", router)
}
