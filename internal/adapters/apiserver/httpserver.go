package apiserver

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	// Graceful shutdown interval, this is a time to wait for HTTP handlers to finish processing
	// data when process interrupt signal received.
	gracefulShutdownTimeout = 3 * time.Second

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	httpReadTimeout = 30 * time.Second

	// WriteTimeout is the maximum duration before timing out
	// writes of the response.
	httpWriteTimeout = 30 * time.Second

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	httpIdleTimeout = 30 * time.Second
)

// Server provides an http.Server
type Server struct {
	*http.Server
	ShutdownCallback func()
	done             chan bool
}

func NewHttpServer(listenAddr string, handler http.Handler) *Server {
	// Create a new logger using zap
	errLogger, err := zap.NewStdLogAt(zap.S().Desugar(), zapcore.ErrorLevel)
	if err != nil {
		zap.S().Fatalf("failed to create standard log: %v", err)
	}
	srv := http.Server{
		Addr:         listenAddr,
		Handler:      handler,
		ErrorLog:     errLogger,
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
		IdleTimeout:  httpIdleTimeout,
	}

	return &Server{
		Server: &srv,
		done:   make(chan bool),
	}
}

// Start runs ListenAndServe on the http.Server with graceful shutdown
func (srv *Server) start() {
	done := make(chan struct{})
	go func() {
		l, err := net.Listen("tcp", srv.Addr)
		if err != nil {
			zap.S().Fatalf("error creating server: %v", err)
		}
		done <- struct{}{}
		if err := srv.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Fatalln("Error serving http content:", err)
		}
	}()
	<-done
	zap.S().Infof("HTTP server is ready to handle requests on address %s", srv.Addr)
}

func (srv *Server) stop() {
	srv.done <- true
}

// WaitWithGracefulShutdown performs waiting for server to receive interrupt signal in order to perform
// graceful shutdown. Normally if application receives a signal to be stopped,
// it will do just that - quit. We give some HTTP server some time to close open connections (and may
// be finish serving HTTP calls).
func (srv *Server) waitWithGracefulShutdown() {
	srv.waitForServerToStop()
	zap.S().Infof("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Warnf("Could not gracefully shutdown the server")
	}
	zap.S().Infof("Server stopped")

	(srv.ShutdownCallback)()
}

func (srv *Server) waitForServerToStop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	select {
	case <-srv.done:
		return
	case sig := <-quit:
		zap.S().Infof("Server received interrupt signal: %s", sig.String())
		return
	}
}
