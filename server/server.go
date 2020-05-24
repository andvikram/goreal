package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/andvikram/goreal/app/service"
	"github.com/andvikram/goreal/configuration"
	"github.com/andvikram/goreal/ds"
	"github.com/andvikram/goreal/logger"
	"github.com/andvikram/goreal/routes"
)

var (
	srv *http.Server
	log = logger.GoRealLog{}
	err error
)

// Start function incorporates the logic around starting the server
func Start(env string) {
	configuration.Initialize()
	logger.Start(env)

	router := routes.ServiceRoutes()
	addr := fmt.Sprintf("%s%s", routes.Host, routes.Port)

	srv = &http.Server{
		Addr:    addr,
		Handler: router,
		// ReadTimeout:    30 * time.Second,
		// WriteTimeout:   30 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	// Initiate service
	initDSVars()
	service.GR = service.NewGoReal()

	servConns := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// Received an interrupt signal, shut down.
		shutdownConnections(servConns)
		shutdownServer()
	}()

	fmt.Println(
		fmt.Sprintf(`
		*!*!* GoReal server running on %s *!*!*

		Ctrl+C to stop
		`, addr),
	)

	if routes.Scheme == "https" {
		err = srv.ListenAndServeTLS(
			configuration.Config.CertFilePath,
			configuration.Config.KeyFilePath,
		)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		log.WithFields(logger.Fields{
			"event": "server.Start()",
			"error": err,
		}).Error("Error starting server")
		shutdownConnections(servConns)
	}

	<-servConns
}

func shutdownConnections(servConns chan struct{}) {
	service.Discontinue = true
	ds.CloseDS()
	logger.Stop()
	close(servConns)
}

func shutdownServer() {
	fmt.Print("\n\tShutting down server ...\n\n")
	err := srv.Shutdown(context.Background())
	if err != nil {
		// Error from closing listeners, or context timeout:
		log.WithFields(logger.Fields{
			"event": "server.Start()",
			"error": err,
		}).Error("Error shutting down server")
	}
}

func initDSVars() {
	ds.Env = configuration.Config.Environment
	ds.DSName = configuration.Config.Datasink
	ds.DSUrl = configuration.Config.DatasinkURL
}
