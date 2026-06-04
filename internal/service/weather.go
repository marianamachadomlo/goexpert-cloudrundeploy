package service

import (
	"context"

	"github.com/marianamachado/cloudrun-deploy/internal/cep"
	"github.com/marianamachado/cloudrun-deploy/internal/temperature"
	"github.com/marianamachado/cloudrun-deploy/internal/weather"
)

type LocationLookup interface {
	Lookup(ctx context.Context, zipcode string) (*cep.Location, error)
}

type WeatherLookup interface {
	CurrentCelsius(ctx context.Context, city, uf string) (float64, error)
}

type WeatherService struct {
	CEPClient     LocationLookup
	WeatherClient WeatherLookup
}

func (s *WeatherService) GetTemperatureByZipcode(ctx context.Context, zipcode string) (temperature.Response, error) {
	if !cep.IsValid(zipcode) {
		return temperature.Response{}, ErrInvalidZipcode
	}

	location, err := s.CEPClient.Lookup(ctx, zipcode)
	if err != nil {
		return temperature.Response{}, err
	}

	celsius, err := s.WeatherClient.CurrentCelsius(ctx, location.City, location.UF)
	if err != nil {
		return temperature.Response{}, err
	}

	return temperature.FromCelsius(celsius), nil
}

var _ LocationLookup = (*cep.ViaCEPClient)(nil)
var _ WeatherLookup = (*weather.Client)(nil)
