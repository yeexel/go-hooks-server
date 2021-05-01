package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebhookTestHandler(t *testing.T) {
	tt := []struct {
		name        string
		method      string
		statusCode  int
		postData    *webhookTestInboundRequest
		want        string
		compareBody bool
		inclHeader  bool
		webhookId   string
	}{
		{
			name:        "invalid request: wrong http method",
			method:      http.MethodGet,
			statusCode:  http.StatusMethodNotAllowed,
			postData:    &webhookTestInboundRequest{},
			want:        "method not allowed",
			compareBody: true,
			inclHeader:  false,
			webhookId:   "",
		},
		{
			name:        "invalid request: validation failed",
			method:      http.MethodPost,
			statusCode:  http.StatusBadRequest,
			postData:    &webhookTestInboundRequest{},
			want:        "",
			compareBody: false,
			inclHeader:  false,
			webhookId:   "",
		},
		{
			name:        "invalid request: header missing",
			method:      http.MethodPost,
			statusCode:  http.StatusBadRequest,
			postData:    &webhookTestInboundRequest{},
			want:        "request header not found: X-WebhookId",
			compareBody: true,
			inclHeader:  false,
			webhookId:   "",
		},
		{
			name:       "valid request",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
			postData: &webhookTestInboundRequest{
				Payload: []byte(`{"some": 123}`),
			},
			want:        "",
			compareBody: false,
			inclHeader:  true,
			webhookId:   "myid",
		},
		{
			name:       "invalid request: wrong webhook ID",
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			postData: &webhookTestInboundRequest{
				Payload: []byte(`{"some": 123}`),
			},
			want:        "webhook not found",
			compareBody: true,
			inclHeader:  true,
			webhookId:   "myid-wrong",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			payloadBuf := new(bytes.Buffer)
			json.NewEncoder(payloadBuf).Encode(tc.postData)

			request := httptest.NewRequest(tc.method, "/api/webhooks/test", payloadBuf)

			if tc.inclHeader {
				request.Header.Set(customHeaderName, tc.webhookId)
			}

			responseRecorder := httptest.NewRecorder()

			WebhookTestHandler{
				webhooks: map[string]webhook{"myid": {
					Url: "https://google.com",
				}},
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
