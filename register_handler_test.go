package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebhookRegisterHandler(t *testing.T) {
	tt := []struct {
		name        string
		method      string
		statusCode  int
		postData    *webhookRegisterRequest
		want        string
		compareBody bool
	}{
		{
			name:        "invalid request: wrong http method",
			method:      http.MethodGet,
			statusCode:  http.StatusMethodNotAllowed,
			postData:    &webhookRegisterRequest{},
			want:        "method not allowed",
			compareBody: true,
		},
		{
			name:        "invalid request: validation failed",
			method:      http.MethodPost,
			statusCode:  http.StatusBadRequest,
			postData:    &webhookRegisterRequest{},
			want:        "",
			compareBody: false,
		},
		{
			name:       "valid request",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
			postData: &webhookRegisterRequest{
				Url:   "http://example.com",
				Token: "token",
			},
			want:        "",
			compareBody: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			payloadBuf := new(bytes.Buffer)
			json.NewEncoder(payloadBuf).Encode(tc.postData)

			request := httptest.NewRequest(tc.method, "/api/webhooks", payloadBuf)
			responseRecorder := httptest.NewRecorder()

			WebhookRegisterHandler{
				webhooks: make(webhookLocalStorage),
			}.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if tc.compareBody {
				if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
					t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
				}
			}
		})
	}
}

func TestGenerateNanoId(t *testing.T) {
	got := len(generateNanoId())
	want := defaultIdLen

	if got != want {
		t.Errorf("wrong ID length; got: %d; want: %d", got, want)
	}
}
