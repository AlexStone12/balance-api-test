package main

import (
	"bank-api-test/internal/database"
	"bank-api-test/internal/handlers/http/api/transaction_add"
	"bank-api-test/internal/handlers/http/api/transaction_history"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	//sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env файл не найден")
	}

}

func main() {
	log.Printf("Service started")
	redis := database.RedisConnect()
	db := database.PostgresConnect()

	if err := dbInit(db); err != nil {
		log.Fatal("Unable to initialize database:", err)
	}

	addHandler := transaction_add.NewHandler(transaction_add.NewResponseBuilder(db, redis))
	historyHandler := transaction_history.NewHandler(transaction_history.NewResponseBuilder(db))

	router := mux.NewRouter()
	router.HandleFunc(transaction_add.Path, addHandler.Handle).Methods(transaction_add.Method)
	router.HandleFunc(transaction_history.Path, historyHandler.Handle).Methods(transaction_history.Method)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func dbInit(db *sqlx.DB) error {
	_, err := db.Exec(database.WalletTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(database.ActionType)
	if err != nil {
		return err
	}
	_, err = db.Exec(database.TransactionStatus)
	if err != nil {
		return err
	}
	_, err = db.Exec(database.TransactionHistoryTable)
	if err != nil {
		return err
	}

	return nil
}
