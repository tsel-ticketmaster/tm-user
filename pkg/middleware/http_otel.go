package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func HTTPOpenTelemetryTracer(handler http.Handler) http.Handler {
	handler = otelhttp.NewHandler(handler, "http-server", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}
