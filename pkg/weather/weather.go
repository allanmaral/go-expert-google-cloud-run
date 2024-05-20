package weather

import (
	"context"
	"errors"
)

type Weather struct {
	TempC   float64 `json:"temp_C"`
	TempF   float64 `json:"temp_F"`
	TempK   float64 `json:"temp_K"`
	Service string  `json:"-"`
}

var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidLocation    = errors.New("invalid location")
	ErrServiceUnavailable = errors.New("service unavailable")
)

type Loader interface {
	Load(ctx context.Context, lat, lng string) (Weather, error)
}

func CelsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func CelsiusToKelvin(c float64) float64 {
	return c + 273
}
