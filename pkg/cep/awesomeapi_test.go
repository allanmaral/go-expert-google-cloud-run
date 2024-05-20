package cep

import (
	"context"
	"errors"
	"testing"
)

func TestAwesomeAPILoader_Load(t *testing.T) {
	t.Run("AwesomeAPI should return error on invalid CEP", func(t *testing.T) {
		sut := NewAwesomeAPILoader()
		ctx := context.Background()

		_, err := sut.Load(ctx, "invalid-cep")

		if !errors.Is(err, ErrInvalidCEP) {
			t.Errorf("expected invalid CEP error, got '%v' instead", err)
		}
	})

	t.Run("AwesomeAPI should return cep not found error on non-existent cep", func(t *testing.T) {
		sut := NewAwesomeAPILoader()
		ctx := context.Background()

		_, err := sut.Load(ctx, "99999999")

		if !errors.Is(err, ErrCEPNotFound) {
			t.Errorf("expected CEP not found error, got '%v' instead", err)
		}
	})

	t.Run("AwesomeAPI should return address on valid cep", func(t *testing.T) {
		sut := NewAwesomeAPILoader()
		ctx := context.Background()

		got, err := sut.Load(ctx, "25808110")

		if err != nil {
			t.Errorf("expected error to be nil, got '%v' instead", err)
		}

		if got.Latitude == "" {
			t.Errorf("expect latitude to be defined, got empty string instead")
		}

		if got.Longitude == "" {
			t.Errorf("expect longitude to be defined, got empty string instead")
		}
	})
}
