package service

import (
	"context"
	"errors"
	"testing"

	"github.com/marianamachado/cloudrun-deploy/internal/cep"
)

type stubCEPClient struct {
	location *cep.Location
	err      error
}

func (s *stubCEPClient) Lookup(_ context.Context, _ string) (*cep.Location, error) {
	return s.location, s.err
}

type stubWeatherClient struct {
	celsius float64
	err     error
}

func (s *stubWeatherClient) CurrentCelsius(_ context.Context, _, _ string) (float64, error) {
	return s.celsius, s.err
}

func TestGetTemperatureByZipcodeSuccess(t *testing.T) {
	svc := &WeatherService{
		CEPClient: &stubCEPClient{
			location: &cep.Location{City: "Rio de Janeiro", UF: "RJ"},
		},
		WeatherClient: &stubWeatherClient{celsius: 30},
	}

	response, err := svc.GetTemperatureByZipcode(context.Background(), "20040020")
	if err != nil {
		t.Fatalf("GetTemperatureByZipcode() error = %v", err)
	}

	if response.TempC != 30 || response.TempF != 86 || response.TempK != 303 {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestGetTemperatureByZipcodeInvalidZipcode(t *testing.T) {
	svc := &WeatherService{
		CEPClient:     &stubCEPClient{},
		WeatherClient: &stubWeatherClient{},
	}

	_, err := svc.GetTemperatureByZipcode(context.Background(), "abc")
	if !errors.Is(err, ErrInvalidZipcode) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidZipcode)
	}
}

func TestGetTemperatureByZipcodeNotFound(t *testing.T) {
	svc := &WeatherService{
		CEPClient: &stubCEPClient{
			err: cep.ErrNotFound,
		},
		WeatherClient: &stubWeatherClient{},
	}

	_, err := svc.GetTemperatureByZipcode(context.Background(), "99999999")
	if !errors.Is(err, cep.ErrNotFound) {
		t.Fatalf("error = %v, want %v", err, cep.ErrNotFound)
	}
}
