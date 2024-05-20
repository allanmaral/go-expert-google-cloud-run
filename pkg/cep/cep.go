package cep

import (
	"context"
	"errors"
)

type CEP struct {
	Cep          string
	Street       string
	Neighborhood string
	City         string
	State        string
	Latitude     string
	Longitude    string
	Service      string
}

var ErrCEPNotFound = errors.New("CEP not found")
var ErrInvalidCEP = errors.New("invalid CEP")
var ErrServiceUnavailable = errors.New("service unavailable")

type Loader interface {
	Load(ctx context.Context, cep string) (CEP, error)
}

func Valid(cep string) bool {
	if cep == "" {
		return false
	}

	if len(cep) != 8 {
		return false
	}

	for _, d := range cep {
		if d < '0' || d > '9' {
			return false
		}
	}

	return true
}
