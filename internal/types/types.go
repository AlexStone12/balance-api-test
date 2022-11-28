package types

type RedisTransaction struct {
	Id     string `json:"id"`
	Amount string `json:"amount"`
	UserId string `json:"userId"`
	Action string `json:"action"`
}

type HistoryTransaction struct {
	Id     string `json:"id"`
	Amount int    `json:"amount"`
	UserId string `json:"userId"`
	Action string `json:"action"`
	Status string `json:"status"`
}

const (
	ActionTypeDecrease string = "decrease"
	ActionTypeIncrease string = "increase"

	TransactionStatusRejected  string = "rejected"
	TransactionStatusProcessed string = "processed"
	TransactionStatusCreated   string = "created"
)
