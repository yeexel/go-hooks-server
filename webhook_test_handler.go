package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// customHeaderName sets a name for custom request header used for authentication.
const customHeaderName = "X-WebhookId"

type webhookTestInboundRequest struct {
	Payload json.RawMessage `json:"payload" validate:"required"`
}

type webhookTestOutboundRequest struct {
	Token   string          `json:"token"`
	Payload json.RawMessage `json:"payload"`
}

type WebhookTestHandler struct {
	webhooks webhookLocalStorage
}

func (h WebhookTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		whId := r.Header.Get(customHeaderName)
		if whId == "" {
			http.Error(w, fmt.Sprintf("request header not found: %s", customHeaderName), http.StatusBadRequest)
			return
		}

		var whPostData webhookTestInboundRequest
		if err := validateRequest(r.Body, &whPostData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Find existing webhook with given ID.
		wh, exists := h.webhooks[whId]
		if !exists {
			http.Error(w, "webhook not found", http.StatusBadRequest)
			return
		}

		payloadBuf := new(bytes.Buffer)
		outboundData := &webhookTestOutboundRequest{
			Token:   wh.Token,
			Payload: whPostData.Payload,
		}
		json.NewEncoder(payloadBuf).Encode(outboundData)
		if _, err := http.Post(wh.Url, "application/json", payloadBuf); err != nil {
			// redirect output to stdout
			fmt.Println(err)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
