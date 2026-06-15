package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/marianamachado/cloudrun-deploy/internal/cep"
	"github.com/marianamachado/cloudrun-deploy/internal/service"
)

type WeatherHandler struct {
	Service *service.WeatherService
}

func NewWeatherHandler(svc *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{Service: svc}
}

func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	zipcode := strings.TrimPrefix(r.URL.Path, "/")
	if zipcode == "" {
		http.Error(w, service.ErrInvalidZipcode.Error(), http.StatusUnprocessableEntity)
		return
	}

	response, err := h.Service.GetTemperatureByZipcode(r.Context(), zipcode)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidZipcode):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case errors.Is(err, cep.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}