package webserver

import (
	"log"
	"net/http"

	"github.com/allanmaral/go-expert-google-cloud-run/pkg/cep"
	"github.com/allanmaral/go-expert-google-cloud-run/pkg/weather"
)

func New(
	logger *log.Logger,
	cepLoader cep.Loader,
	weatherLoader weather.Loader,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, logger, cepLoader, weatherLoader)

	var handler http.Handler = mux
	handler = withLogging(logger, handler)

	return handler
}
