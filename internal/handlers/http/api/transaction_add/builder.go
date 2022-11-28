package transaction_add

import (
	"bank-api-test/internal/database"
	"bank-api-test/internal/types"
	"encoding/json"
	"fmt"
	"github.com/adjust/rmq/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type responseBuilder struct {
	db    *sqlx.DB
	redis rmq.Connection
}

func NewResponseBuilder(db *sqlx.DB, connection rmq.Connection) *responseBuilder {
	return &responseBuilder{
		db:    db,
		redis: connection,
	}
}

func (r *responseBuilder) AddTransaction(id string, amount string, action string) (*Response, error) {
	balance, err := database.GetOrCreateUserBalance(r.db, id)
	if err != nil {
		return nil, err
	}
	queue, err := r.redis.OpenQueue(id)
	if err != nil {
		return nil, err
	}

	transactionId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	if action != types.ActionTypeDecrease && action != types.ActionTypeIncrease {
		return nil, fmt.Errorf("this action (%s) is not allowed", action)
	}

	amountInt, _ := strconv.Atoi(amount)
	if action == types.ActionTypeDecrease && amountInt > balance {
		return nil, fmt.Errorf("balance (%d) is not enough (%d)", balance, amountInt)
	}

	t := types.RedisTransaction{
		Id:     transactionId.String(),
		UserId: id,
		Amount: amount,
		Action: action,
	}

	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	err = queue.PublishBytes(data)
	if err != nil {
		return nil, err
	}

	return &Response{TransactionId: transactionId.String()}, nil
}
