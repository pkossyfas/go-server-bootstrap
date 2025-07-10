/*
Package server is responsible for the setting up the http server.
*/
package server

import (
	"context"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkossyfas/go-server-bootstrap/config"
	"github.com/pkossyfas/go-server-bootstrap/controller"
	"github.com/pkossyfas/go-server-bootstrap/logger"
	"github.com/pkossyfas/go-server-bootstrap/metrics"
)

var requestRouter *http.ServeMux

// LoadEndpoints loads the default endpoints and their handlers.
func LoadEndpoints() {
	if requestRouter == nil {
		requestRouter = http.NewServeMux()
	}

	requestRouter.Handle("/metrics", metrics.Handler())
	requestRouter.Handle("/version", controller.VersionEndpointMetrics)
	requestRouter.Handle("/health", controller.HealthEndpointMetrics)
	requestRouter.Handle("/ready", controller.ReadyEndpointMetrics)
}

// StartServer starts the main http server.
func StartServer() {
	// Create the server
	port := config.GetConfig.ServerPort
	readTimeout, err := time.ParseDuration(config.GetConfig.ServerReadTimeout)
	if err != nil {
		logger.Fatal(err)
	}
	writeTimeout, err := time.ParseDuration(config.GetConfig.ServerWriteTimeout)
	if err != nil {
		logger.Fatal(err)
	}

	srv := &http.Server{
		Handler:      requestRouter,
		Addr:         ":" + port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Create a channel to be closed after graceful shutdown.
	idleConnsClosed := make(chan struct{})

	// Start the thread responsible for graceful shutdown.
	go gracefulServerShutdown(srv, idleConnsClosed)

	// Start the server.
	logger.Info("starting the http server with port %v", port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.Fatal(err)
	}

	// Wait for graceful shutdown
	<-idleConnsClosed

	// It will reach this point when server is being stopped
	logger.Info("server stopped")

}

// gracefulServerShutdown waits for an os/sycall interruption signal.
// When an os interrupt signal (e.g. Ctrl-c) or a SIGTERM system call (kill -15)
// is received, the server will be gracefully shutting down within a given time.
func gracefulServerShutdown(srv *http.Server, idleConnsClosed chan struct{}) {
	// make the channel which waits for an os.Signal
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	// We received an interrupt signal, shut down within a given time.
	logger.Info("shutting down the http server gracefully..")
	shutdownTimeout, _ := time.ParseDuration(config.GetConfig.ShutdownTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		logger.Error(err, "graceful shutdown failed")
	}
	close(idleConnsClosed)
}
