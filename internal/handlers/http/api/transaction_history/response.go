package transaction_history

import "bank-api-test/internal/types"

type Response struct {
	Transactions []types.HistoryTransaction `json:"transactions"`
}
