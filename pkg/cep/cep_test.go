package cep

import "testing"

func TestValid(t *testing.T) {
	t.Run("Valid should return false on empty cep", func(t *testing.T) {
		cep := ""

		got := Valid(cep)

		if got != false {
			t.Errorf("expected Valid to return false, got %v instead", got)
		}
	})

	t.Run("Valid should return false on cep with length different then 8 digits", func(t *testing.T) {
		tests := []string{
			"1", "12", "123", "1234", "12345", "123456", "1234567", "123456789", "1234567890",
		}
		for _, cep := range tests {
			got := Valid(cep)

			if got != false {
				t.Errorf("(%s): expected Valid to return false, got %v instead", cep, got)
			}
		}
	})

	t.Run("Valid should return false on cep with non-numeric characters", func(t *testing.T) {
		tests := []string{
			"12345 12", "1234 012", "12345 123", "12345-12", "1234-012", "12345-123",
		}
		for _, cep := range tests {
			got := Valid(cep)

			if got != false {
				t.Errorf("(%s): expected Valid to return false, got %v instead", cep, got)
			}
		}
	})

	t.Run("Valid should return true on cep with 8 numeric characters", func(t *testing.T) {
		tests := []string{
			"12345678", "25808110", "25809600", "25800000",
		}
		for _, cep := range tests {
			got := Valid(cep)

			if got != true {
				t.Errorf("(%s): expected Valid to return true, got %v instead", cep, got)
			}
		}
	})
}
