package main

import (
	"encoding/json"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const defaultIdLen = 7

type webhookRegisterRequest struct {
	Url   string `json:"url" validate:"required,url"`
	Token string `json:"token" validate:"required"`
}

type WebhookRegisterHandler struct {
	webhooks webhookLocalStorage
}

func (h WebhookRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var whData webhookRegisterRequest
		if err := validateRequest(r.Body, &whData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := generateNanoId()
		wh := webhook{
			Id:    id,
			Url:   whData.Url,
			Token: whData.Token,
		}

		// Persist webhook inside the map.
		h.webhooks[id] = wh

		// Return webhook data as JSON with webhook ID in place.
		json.NewEncoder(w).Encode(wh)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// generateNanoId generates alphanumeric string with defined length.
func generateNanoId() string {
	id, _ := gonanoid.New(defaultIdLen)

	return id
}
