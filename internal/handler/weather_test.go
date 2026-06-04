package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marianamachado/cloudrun-deploy/internal/cep"
	"github.com/marianamachado/cloudrun-deploy/internal/service"
	"github.com/marianamachado/cloudrun-deploy/internal/temperature"
)

type mockCEPClient struct {
	location *cep.Location
	err      error
}

func (m *mockCEPClient) Lookup(_ context.Context, _ string) (*cep.Location, error) {
	return m.location, m.err
}

type mockWeatherClient struct {
	celsius float64
	err     error
}

func (m *mockWeatherClient) CurrentCelsius(_ context.Context, _, _ string) (float64, error) {
	return m.celsius, m.err
}

func TestWeatherHandlerSuccess(t *testing.T) {
	svc := &service.WeatherService{
		CEPClient: &mockCEPClient{
			location: &cep.Location{City: "São Paulo", UF: "SP"},
		},
		WeatherClient: &mockWeatherClient{celsius: 28.5},
	}

	handler := NewWeatherHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/01310100", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var response temperature.Response
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.TempC != 28.5 || response.TempF != 83.3 || response.TempK != 301.5 {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestWeatherHandlerInvalidZipcode(t *testing.T) {
	svc := &service.WeatherService{
		CEPClient:     &mockCEPClient{},
		WeatherClient: &mockWeatherClient{},
	}

	handler := NewWeatherHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/1234", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
	if rec.Body.String() != "invalid zipcode\n" {
		t.Fatalf("body = %q, want %q", rec.Body.String(), "invalid zipcode\n")
	}
}

func TestWeatherHandlerZipcodeNotFound(t *testing.T) {
	svc := &service.WeatherService{
		CEPClient: &mockCEPClient{
			err: cep.ErrNotFound,
		},
		WeatherClient: &mockWeatherClient{},
	}

	handler := NewWeatherHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/99999999", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if rec.Body.String() != "can not find zipcode\n" {
		t.Fatalf("body = %q, want %q", rec.Body.String(), "can not find zipcode\n")
	}
}

func TestWeatherHandlerInternalError(t *testing.T) {
	svc := &service.WeatherService{
		CEPClient: &mockCEPClient{
			location: &cep.Location{City: "São Paulo", UF: "SP"},
		},
		WeatherClient: &mockWeatherClient{
			err: errors.New("weather unavailable"),
		},
	}

	handler := NewWeatherHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/01310100", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
