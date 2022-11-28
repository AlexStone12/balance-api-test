package transaction_history

import (
	"bank-api-test/internal/database"
	"github.com/jmoiron/sqlx"
)

type responseBuilder struct {
	db *sqlx.DB
}

func NewResponseBuilder(db *sqlx.DB) *responseBuilder {
	return &responseBuilder{
		db: db,
	}
}

func (r *responseBuilder) GetHistory(userId string) (*Response, error) {
	transactions, err := database.GetTransactionsByUserId(r.db, userId)
	if err != nil {
		return nil, err
	}
	return &Response{Transactions: transactions}, nil
}
