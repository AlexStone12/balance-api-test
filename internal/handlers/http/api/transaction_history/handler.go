package transaction_history

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	Path   = "/transaction/history"
	Method = "GET"
)

type handler struct {
	builder *responseBuilder
}

func NewHandler(builder *responseBuilder) *handler {
	return &handler{
		builder: builder,
	}
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-UserId")

	if id == "" {
		w.WriteHeader(400)
		return
	}

	transactions, err := h.builder.GetHistory(id)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("%e", err)
		return
	}

	data, err := json.Marshal(transactions)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("%e", err)
		return
	}

	headers := w.Header()
	headers.Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Printf("%e", err)
		return
	}
}
