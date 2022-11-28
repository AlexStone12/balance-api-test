package database

import (
	"fmt"
	"github.com/adjust/rmq/v5"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

func RedisConnect() rmq.Connection {
	rHost, exists := os.LookupEnv("REDIS_HOST")
	if !exists {
		rHost = "localhost"
	}
	rPort, exists := os.LookupEnv("REDIS_PORT")
	if !exists {
		rHost = "6379"
	}
	rAddress := fmt.Sprintf("%s:%s", rHost, rPort)
	redisConnection, err := rmq.OpenConnection("service", "tcp", rAddress, 1, nil)
	if err != nil {
		log.Fatal(err)
	}
	return redisConnection
}

func PostgresConnect() *sqlx.DB {
	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		dbName = "balance"
	}
	dbHost, exists := os.LookupEnv("DB_HOST")
	if !exists {
		dbHost = "postgres"
	}
	dbPort, exists := os.LookupEnv("DB_PORT")
	if !exists {
		dbPort = "5432"
	}
	dbUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		dbUser = "balance"
	}
	dbPass, exists := os.LookupEnv("DB_PASS")
	if !exists {
		dbPass = "pass"
	}
	address := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	dbConnection, err := sqlx.Connect("postgres", address)
	if err != nil {
		log.Fatal(err)
	}
	return dbConnection
}
