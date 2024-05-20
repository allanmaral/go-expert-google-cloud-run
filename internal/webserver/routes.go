package webserver

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/allanmaral/go-expert-google-cloud-run/pkg/cep"
	"github.com/allanmaral/go-expert-google-cloud-run/pkg/weather"
)

func addRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	cepLoader cep.Loader,
	weatherLoader weather.Loader,
) {
	mux.Handle("GET /api/weather/{cep}", handleGetTemperature(logger, cepLoader, weatherLoader))
	mux.Handle("GET /api/health", handleHealth(cepLoader, weatherLoader))
	mux.Handle("GET /api/ready", handleReady())
}

type errorResponse struct {
	Message string `json:"Message"`
}

func handleGetTemperature(
	logger *log.Logger,
	cepLoader cep.Loader,
	weatherLoader weather.Loader,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := r.PathValue("cep")

		cepRes, err := cepLoader.Load(r.Context(), input)
		if err != nil {
			if errors.Is(err, cep.ErrInvalidCEP) {
				_ = encode(w, r, http.StatusUnprocessableEntity, errorResponse{Message: "invalid zipcode"})
				return
			} else if errors.Is(err, cep.ErrCEPNotFound) {
				_ = encode(w, r, http.StatusNotFound, errorResponse{Message: "can not find zipcode"})
				return
			} else if errors.Is(err, cep.ErrServiceUnavailable) {
				_ = encode(w, r, http.StatusBadGateway, errorResponse{Message: "cep service is unavailable, try again later"})
				logger.Printf("cep service is unavailable %s\n", err)
				return
			} else {
				_ = encode(w, r, http.StatusInternalServerError, errorResponse{Message: "internal server error"})
				logger.Printf("unhandled error while loading cep %s\n", err)
				return
			}
		}

		weatherRes, err := weatherLoader.Load(r.Context(), cepRes.Latitude, cepRes.Longitude)
		if err != nil {
			if errors.Is(err, weather.ErrServiceUnavailable) {
				_ = encode(w, r, http.StatusBadGateway, errorResponse{Message: "weather service is unavailable, try again later"})
				logger.Printf("weather service in unavailable %s\n", err)
				return
			} else {
				_ = encode(w, r, http.StatusInternalServerError, errorResponse{Message: "internal server error"})
				logger.Printf("unhandled error while loading weather %s\n", err)
				return
			}
		}

		_ = encode(w, r, http.StatusOK, weatherRes)
	})
}

func handleHealth(cepLoader cep.Loader, weatherLoader weather.Loader) http.Handler {
	type response struct {
		OK      bool   `json:"ok"`
		Weather string `json:"weather"`
		CEP     string `json:"cep"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res := response{OK: true}

			_, err := cepLoader.Load(r.Context(), "70150900")
			if err != nil {
				res.CEP = fmt.Sprintf("Error: %s", err)
				res.OK = false
			} else {
				res.CEP = "Working!"
			}

			_, err = weatherLoader.Load(r.Context(), "-15.80097", "-47.86072")
			if err != nil {
				res.Weather = fmt.Sprintf("Error: %s", err)
				res.OK = false
			} else {
				res.Weather = "Working!"
			}

			status := http.StatusOK
			if !res.OK {
				status = http.StatusInternalServerError
			}

			_ = encode(w, r, status, res)
		},
	)
}

func handleReady() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)
}
