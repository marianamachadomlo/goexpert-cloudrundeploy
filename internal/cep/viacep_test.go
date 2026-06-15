package cep

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViaCEPClientLookupSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/01310100/json/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(map[string]string{
			"localidade": "São Paulo",
			"uf":         "SP",
		})
	}))
	defer server.Close()

	client := &ViaCEPClient{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}

	location, err := client.Lookup(context.Background(), "01310100")
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}

	if location.City != "São Paulo" || location.UF != "SP" {
		t.Fatalf("unexpected location: %+v", location)
	}
}

func TestViaCEPClientLookupNotFound(t *testing.T) {
	tests := []struct {
		name string
		body map[string]any
	}{
		{
			name: "erro as bool",
			body: map[string]any{"erro": true},
		},
		{
			name: "erro as string",
			body: map[string]any{"erro": "true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(tt.body)
			}))
			defer server.Close()

			client := &ViaCEPClient{
				BaseURL:    server.URL,
				HTTPClient: server.Client(),
			}

			_, err := client.Lookup(context.Background(), "99999999")
			if err != ErrNotFound {
				t.Fatalf("Lookup() error = %v, want %v", err, ErrNotFound)
			}
		})
	}
}
