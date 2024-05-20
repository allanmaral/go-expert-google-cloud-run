package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/allanmaral/go-expert-google-cloud-run/internal/webserver"
	"github.com/allanmaral/go-expert-google-cloud-run/pkg/cep"
	"github.com/allanmaral/go-expert-google-cloud-run/pkg/weather"
)

func run(
	ctx context.Context,
	getEnv func(key string) string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := log.New(stdout, "SERVER: ", log.LstdFlags)
	cepLoader := cep.NewAwesomeAPILoader()
	weatherLoader := weather.NewWeatherAPILoader(getEnv("WEATHER_APIKEY"))

	srv := webserver.New(logger, cepLoader, weatherLoader)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(getEnv("HOST"), getEnv("PORT")),
		Handler: srv,
	}

	go func() {
		logger.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			_, _ = fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Printf("Shutting http server down...\n")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			_, _ = fmt.Fprintf(stderr, "error shutting http server down: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()

	if err := run(ctx, os.Getenv, os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
