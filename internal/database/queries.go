package database

import (
	"bank-api-test/internal/types"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	WalletTableName  = "wallet"
	HistoryTableName = "t_history"

	defaultBalance = 0
)

const WalletTable = `CREATE TABLE  IF NOT EXISTS wallet (
user_id VARCHAR(50) PRIMARY KEY,
amount INT
);
`

const TransactionHistoryTable = `CREATE TABLE IF NOT EXISTS t_history (
id UUID PRIMARY KEY,
user_id VARCHAR(50) REFERENCES wallet (user_id),
action_type action,
amount INT,
status status
);
`

const ActionType = `DO $$
begin
IF NOT EXISTS (
SELECT 1 FROM pg_type WHERE typname = 'action'
)
THEN  CREATE TYPE action AS ENUM ('increase', 'decrease');
END IF;
end$$;
`

const TransactionStatus = `DO $$
begin
IF NOT EXISTS (
SELECT 1 FROM pg_type WHERE typname = 'status'
)
THEN  CREATE TYPE status AS ENUM ('rejected', 'processed', 'created');
END IF;
end$$;
`

func GetTransactionsByUserId(db *sqlx.DB, userId string) ([]types.HistoryTransaction, error) {
	rows, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"user_id",
			"action_type",
			"amount",
			"status",
		).
		From(HistoryTableName).
		Where(sq.Eq{
			"user_id": userId,
		}).
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}

	var transactions []types.HistoryTransaction
	for rows.Next() {
		var (
			id         string
			userId     string
			actionType string
			amount     int
			status     string
		)

		err := rows.Scan(&id, &userId, &actionType, &amount, &status)
		if err != nil {
			return nil, err
		}
		t := types.HistoryTransaction{
			Id:     id,
			UserId: userId,
			Amount: amount,
			Action: actionType,
			Status: status,
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// GetOrCreateUserBalance Т.К X-UserId передается от сервиса авторизации, то это существующий пользователь,
// если у него нет баланса, то создаем кошелек с нулевым балансом
func GetOrCreateUserBalance(db *sqlx.DB, userId string) (int, error) {
	rows, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"amount",
		).
		From(WalletTableName).
		Where(sq.Eq{
			"user_id": userId,
		}).
		RunWith(db).Query()
	if err != nil {
		return 0, err
	}

	amount := -1
	for rows.Next() {
		err = rows.Scan(&amount)
		if err != nil {
			return 0, err
		}
	}
	if amount < 0 {
		amount = defaultBalance
		_, err = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
			Insert(WalletTableName).
			Columns("user_id", "amount").
			Values(
				userId,
				amount,
			).
			RunWith(db).
			Exec()
		if err != nil {
			return 0, err
		}
	}
	return amount, nil
}

func SetTransactionStatus(db *sqlx.DB, transactionId string, status string) error {
	_, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(HistoryTableName).
		Set("status", status).
		Where(sq.Eq{
			"id": transactionId,
		}).
		RunWith(db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func SetUserBalance(db *sqlx.DB, userId string, balance int) error {
	_, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(WalletTableName).
		Set("amount", balance).
		Where(sq.Eq{
			"user_id": userId,
		}).
		RunWith(db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func CreateTransaction(db *sqlx.DB, transaction types.HistoryTransaction) error {
	_, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(HistoryTableName).
		Columns("id", "amount", "user_id", "action_type", "status").
		Values(
			transaction.Id,
			transaction.Amount,
			transaction.UserId,
			transaction.Action,
			types.TransactionStatusCreated,
		).
		RunWith(db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
