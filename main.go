package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)

// webhook private struct defines basic properties of a single webhook.
type webhook struct {
	Id    string `json:"id"`
	Url   string `json:"url"`
	Token string `json:"token"`
}

// webhookLocalStorage represents in-memory storage for all incoming webhooks.
type webhookLocalStorage map[string]webhook

func main() {
	// init in-memory storage
	storage := make(webhookLocalStorage)

	mux := http.NewServeMux()

	mux.Handle("/api/webhooks", WebhookRegisterHandler{storage})
	mux.Handle("/api/webhooks/test", WebhookTestHandler{storage})

	log.Fatal(http.ListenAndServe(":9876", logRequest(mux)))
}

// validateRequest decodes request body into struct and performs validation afterwords.
func validateRequest(requestBody io.ReadCloser, target interface{}) error {
	validate := validator.New()

	json.NewDecoder(requestBody).Decode(&target)

	if err := validate.Struct(target); err != nil {
		return err
	}

	return nil
}

// logRequest adds basic logging functionality to the server.
// Log entry example: 2021/05/01 19:52:40 [::1]:58694 POST /api/webhooks
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
