package main

import (
	"log"
	"net/http"
	"os"

	"github.com/marianamachado/cloudrun-deploy/internal/cep"
	"github.com/marianamachado/cloudrun-deploy/internal/handler"
	"github.com/marianamachado/cloudrun-deploy/internal/service"
	"github.com/marianamachado/cloudrun-deploy/internal/weather"
)

func main() {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	svc := &service.WeatherService{
		CEPClient:     cep.NewViaCEPClient(),
		WeatherClient: weather.NewClient(apiKey),
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler.NewWeatherHandler(svc))

	addr := ":" + port
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
