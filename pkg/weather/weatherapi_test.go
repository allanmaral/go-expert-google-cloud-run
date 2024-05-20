package weather

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestWeatherAPILoader_Load(t *testing.T) {
	apikey := os.Getenv("WEATHER_APIKEY")

	t.Run("WeatherAPI should return unauthorized error on invalid API KEY", func(t *testing.T) {
		sut := NewWeatherAPILoader("invalid-key")
		ctx := context.Background()

		_, err := sut.Load(ctx, "0.000", "0.000")

		if !errors.Is(err, ErrUnauthorized) {
			t.Errorf("expected unauthorized error, got '%v' instead", err)
		}
	})

	t.Run("WeatherAPI should return error on invalid latitude and longitude", func(t *testing.T) {
		sut := NewWeatherAPILoader(apikey)
		ctx := context.Background()

		_, err := sut.Load(ctx, "", "")

		if !errors.Is(err, ErrInvalidLocation) {
			t.Errorf("expected invalid location error, got '%v' instead", err)
		}
	})

	t.Run("WeatherAPI should return temperature on valid location", func(t *testing.T) {
		sut := NewWeatherAPILoader(apikey)
		ctx := context.Background()

		got, err := sut.Load(ctx, "-22.09967", "-43.2116")

		if err != nil {
			t.Errorf("expected error to be nil, got '%v' instead", err)
		}

		c := got.TempC
		wantK := CelsiusToKelvin(c)
		if got.TempK != wantK {
			t.Errorf("expected TempK to be %f, got %f instead", got.TempK, wantK)
		}

		wantF := CelsiusToFahrenheit(c)
		if got.TempF != wantF {
			t.Errorf("expected TempF to be %f, got %f instead", got.TempF, wantF)
		}
	})
}
