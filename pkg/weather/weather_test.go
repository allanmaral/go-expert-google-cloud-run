package weather

import "testing"

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		celsius    float64
		fahrenheit float64
	}{
		{celsius: -40, fahrenheit: -40},
		{celsius: 0, fahrenheit: 32},
		{celsius: 100, fahrenheit: 212},
	}
	for i, test := range tests {
		got := CelsiusToFahrenheit(test.celsius)
		if got != test.fahrenheit {
			t.Errorf("(%d): expected %f, got %f instead", i, test.fahrenheit, got)
		}
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		celsius float64
		kelvin  float64
	}{
		{celsius: -273, kelvin: 0},
		{celsius: 0, kelvin: 273},
		{celsius: 100, kelvin: 373},
	}
	for i, test := range tests {
		got := CelsiusToKelvin(test.celsius)
		if got != test.kelvin {
			t.Errorf("(%d): expected %f, got %f instead", i, test.kelvin, got)
		}
	}
}
