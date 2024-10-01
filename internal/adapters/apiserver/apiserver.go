package apiserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"go.uber.org/zap"

	"kl-webhook-adapter/internal/adapters/apiserver/internal"
	"kl-webhook-adapter/internal/adapters/broker"
	"kl-webhook-adapter/internal/core/app"
)

func Start(_ context.Context, cfg *app.Config) {
	_, err := broker.JetStreamInit(&cfg.Broker)

	if err != nil {
		zap.S().Fatal(err)
		return
	}
	defer broker.NATSConnection.Close()
	zap.S().Infof("Successfully connected to the Broker")

	listenAddr := fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port)
	zap.S().Infof("Starting http server on %s", listenAddr)

	mux := http.NewServeMux()

	apiRoutes(mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	})
	hdl := c.Handler(mux)
	hdl = recoveryMiddleware()(hdl)
	hdl = zapLoggerMiddleware(zap.S())(hdl)

	srv := NewHttpServer(listenAddr, hdl)
	srv.ShutdownCallback = func() {
		zap.S().Info("Cleaning up resources")
		// add cleaning up resources
		zap.S().Infof("Resources has been cleaned up")
	}
	zap.S().Info("Starting HTTP server...")
	go func() {
		srv.start()
	}()

	srv.waitWithGracefulShutdown()
}

func apiRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/version", internal.VersionRouteHandler())

	mux.HandleFunc(internal.AdapterRoutePath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			internal.SendToBroker()(w, r)
		default:
			internal.RenderMethodNotAllowed(w, r)
		}
	})
}
