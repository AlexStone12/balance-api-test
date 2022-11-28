package transaction_add

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	Path   = "/transaction/add"
	Method = "POST"
)

type handler struct {
	builder *responseBuilder
}

func NewHandler(builder *responseBuilder) *handler {
	return &handler{builder: builder}
}

// Handle X-UserId пробрасывается от сервиса авторизации, его наличие подтверждает наличие такового в базе
func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-UserId")
	if id == "" {
		log.Printf("X-UserId not found: ")
		w.WriteHeader(400)
		return
	}
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("JSON Unmarshalling error: %e", err)
		w.WriteHeader(400)
		return
	}

	response, err := h.builder.AddTransaction(id, req.Amount, req.Action)
	if err != nil {
		log.Printf("AddTransaction error: %e", err)
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(500)
		log.Print(err)
		return
	}

	headers := w.Header()
	headers.Add("Content-Type", "application/json")

	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}
