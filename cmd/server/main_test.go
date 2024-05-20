package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestServerHealth(t *testing.T) {
	t.Run("server should return 200 on /api/health", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		client := http.Client{}
		getEnv := func(key string) string {
			if key == "WEATHER_APIKEY" {
				return os.Getenv(key)
			}
			return ""
		}
		s := spawnServer(ctx, t, getEnv)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/health", s.endpoint), nil)
		if err != nil {
			t.Errorf("failed to create request: %s", err)
		}
		resp, err := client.Do(req)

		if err != nil {
			t.Errorf("failed to make request: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Errorf("expected status code 200, got %d instead", resp.StatusCode)
		}
	})

	t.Run("server should return 500 calling /api/health on invalid API_KEY", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		client := http.Client{}
		getEnv := func(key string) string {
			if key == "WEATHER_APIKEY" {
				return "invalid-key"
			}
			return ""
		}
		s := spawnServer(ctx, t, getEnv)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/health", s.endpoint), nil)
		if err != nil {
			t.Errorf("failed to create request: %s", err)
		}
		resp, err := client.Do(req)

		if err != nil {
			t.Errorf("failed to make request: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 500 {
			t.Errorf("expected status code 500, got %d instead", resp.StatusCode)
		}
	})
}

func TestServerWeather(t *testing.T) {
	type response struct {
		TempC float64 `json:"temp_C"`
		TempF float64 `json:"temp_F"`
		TempK float64 `json:"temp_K"`
	}
	type errorResponse struct {
		Message string `json:"message"`
	}

	t.Run("server should return 200 calling /api/weather on valid cep", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		client := http.Client{}
		getEnv := func(key string) string {
			if key == "WEATHER_APIKEY" {
				return os.Getenv(key)
			}
			return ""
		}
		s := spawnServer(ctx, t, getEnv)
		validCEP := "70150900"

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/weather/%s", s.endpoint, validCEP), nil)
		if err != nil {
			t.Errorf("failed to create request: %s", err)
		}
		resp, err := client.Do(req)

		if err != nil {
			t.Errorf("failed to make request: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Errorf("expected status code 200, got %d instead", resp.StatusCode)
		}

		var r response
		if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
			t.Errorf("failed to parse response: %s", err)
		}
		if r.TempK == 0 {
			t.Errorf("invalid temperature data C: %f, F: %f, K: %f", r.TempC, r.TempF, r.TempK)
		}
	})

	t.Run("server should return 422 calling /api/weather on invalid cep", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		client := http.Client{}
		getEnv := func(key string) string {
			if key == "WEATHER_APIKEY" {
				return os.Getenv(key)
			}
			return ""
		}
		s := spawnServer(ctx, t, getEnv)
		validCEP := "invalid-cep"

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/weather/%s", s.endpoint, validCEP), nil)
		if err != nil {
			t.Errorf("failed to create request: %s", err)
		}
		resp, err := client.Do(req)

		if err != nil {
			t.Errorf("failed to make request: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 422 {
			t.Errorf("expected status code 422, got %d instead", resp.StatusCode)
		}

		var r errorResponse
		if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
			t.Errorf("failed to parse response: %s", err)
		}
		if r.Message != "invalid zipcode" {
			t.Errorf("invalid error message: %s", r.Message)
		}
	})

	t.Run("server should return 404 calling /api/weather on non-existent cep", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		client := http.Client{}
		getEnv := func(key string) string {
			if key == "WEATHER_APIKEY" {
				return os.Getenv(key)
			}
			return ""
		}
		s := spawnServer(ctx, t, getEnv)
		validCEP := "99999999"

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/weather/%s", s.endpoint, validCEP), nil)
		if err != nil {
			t.Errorf("failed to create request: %s", err)
		}
		resp, err := client.Do(req)

		if err != nil {
			t.Errorf("failed to make request: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 404 {
			t.Errorf("expected status code 404, got %d instead", resp.StatusCode)
		}

		var r errorResponse
		if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
			t.Errorf("failed to parse response: %s", err)
		}
		if r.Message != "can not find zipcode" {
			t.Errorf("invalid error message: %s", r.Message)
		}
	})
}

type server struct {
	endpoint string
	stdout   *bytes.Buffer
	stderr   *bytes.Buffer
}

// spawnServer starts a new server using the run function and waits
// for the server to start answering the requests.
func spawnServer(ctx context.Context, t testing.TB, getEnv func(string) string) *server {
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("failed to find a open tcp port: %s", err)
	}

	env := func(key string) string {
		if key == "HOST" {
			return "localhost"
		}
		if key == "PORT" {
			return fmt.Sprintf("%d", port)
		}
		return getEnv(key)
	}

	s := server{
		endpoint: fmt.Sprintf("http://localhost:%d", port),
		stdout:   new(bytes.Buffer),
		stderr:   new(bytes.Buffer),
	}

	go func() {
		if err := run(ctx, env, s.stdout, s.stderr); err != nil {
			t.Errorf("failed to run server: %s\n", err)
		}
	}()

	if err := waitForReady(ctx, 1*time.Second, fmt.Sprintf("%s/api/ready", s.endpoint)); err != nil {
		t.Fatalf("%s\n", err)
	}

	return &s
}

// getFreePort asks the kernel for an available tcp port.
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	if err = listener.Close(); err != nil {
		return 0, err
	}

	return port, nil
}

// waitForReady calls the specified endpoint until it gets a 200
// response or until the context is cancelled or the timeout is
// reached.
func waitForReady(
	ctx context.Context,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			if time.Since(start) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}

			time.Sleep(250 * time.Millisecond)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %s\n", err.Error())
			continue
		}

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Endpoint is ready!")
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()
	}
}
