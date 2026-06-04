package weather

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientCurrentCelsius(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") == "" {
			t.Fatal("missing q query parameter")
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"current": map[string]float64{
				"temp_c": 28.5,
			},
		})
	}))
	defer server.Close()

	client := &Client{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}

	temp, err := client.CurrentCelsius(context.Background(), "São Paulo", "SP")
	if err != nil {
		t.Fatalf("CurrentCelsius() error = %v", err)
	}

	if temp != 28.5 {
		t.Fatalf("CurrentCelsius() = %v, want 28.5", temp)
	}
}
