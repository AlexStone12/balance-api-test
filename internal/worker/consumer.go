package worker

import (
	"bank-api-test/internal/database"
	"bank-api-test/internal/types"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"time"

	"github.com/adjust/rmq/v5"
)

const (
	consumeDuration = time.Millisecond
)

type Consumer struct {
	name   string
	count  int
	before time.Time
	db     *sqlx.DB
}

func NewConsumer(tag string, db *sqlx.DB) *Consumer {
	return &Consumer{
		name:   fmt.Sprintf("consumer-%s", tag),
		before: time.Now(),
		db:     db,
	}
}

func (consumer *Consumer) Consume(delivery rmq.Delivery) {
	payload := delivery.Payload()
	transaction := types.RedisTransaction{}
	err := json.Unmarshal([]byte(payload), &transaction)
	if err != nil {
		log.Printf("%e", err)
		return
	}
	fmt.Printf("Consumed: %+v", transaction)
	amount, err := strconv.Atoi(transaction.Amount)
	if err != nil {
		log.Printf("%e", err)
		return
	}

	historyTransaction := types.HistoryTransaction{
		Id:     transaction.Id,
		UserId: transaction.UserId,
		Amount: amount,
		Action: transaction.Action,
		Status: types.TransactionStatusCreated,
	}
	balance, err := database.GetOrCreateUserBalance(consumer.db, transaction.UserId)
	if err != nil {
		log.Printf("%e", err)
		return
	}

	err = database.CreateTransaction(consumer.db, historyTransaction)
	if err != nil {
		log.Printf("Create Transaction err: %e", err)
		return
	}

	switch transaction.Action {
	case types.ActionTypeIncrease:
		err = database.SetUserBalance(consumer.db, transaction.UserId, balance+amount)
		if err != nil {
			return
		}
		err = database.SetTransactionStatus(consumer.db, transaction.Id, types.TransactionStatusProcessed)
		if err != nil {
			return
		}
		err = delivery.Ack()
		if err != nil {
			return
		}
	case types.ActionTypeDecrease:
		if balance-amount < 0 {
			err = delivery.Reject()
			if err != nil {
				return
			}
			err = database.SetTransactionStatus(consumer.db, transaction.Id, types.TransactionStatusRejected)
			if err != nil {
				return
			}
		}

		err = database.SetUserBalance(consumer.db, transaction.UserId, balance-amount)
		if err != nil {
			return
		}
		err = database.SetTransactionStatus(consumer.db, transaction.Id, types.TransactionStatusProcessed)
		if err != nil {
			return
		}
		err = delivery.Ack()
		if err != nil {
			return
		}
	}
	time.Sleep(consumeDuration)
}
