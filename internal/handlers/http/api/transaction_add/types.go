package transaction_add

type Response struct {
	TransactionId string `json:"transactionId"`
}

type Request struct {
	Amount string `json:"amount"`
	Action string `json:"action"`
}
