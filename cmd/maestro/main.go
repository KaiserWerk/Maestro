package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/KaiserWerk/Maestro/internal/cache"
	"github.com/KaiserWerk/Maestro/internal/configuration"
	"github.com/KaiserWerk/Maestro/internal/global"
	"github.com/KaiserWerk/Maestro/internal/handler"
	"github.com/KaiserWerk/Maestro/internal/logging"
	"github.com/KaiserWerk/Maestro/internal/middleware"
	"github.com/KaiserWerk/Maestro/internal/panicHandler"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	bindAddr   string
	configFile = flag.String("config", "app.yaml", "The configuration file to use")
	logDir     = flag.String("logdir", ".", "The directory to save log files to")
)

func main() {
	flag.Parse()

	if *configFile == "" {
		fmt.Println("The configuration file parameter is empty, please supply a valid value.")
		return
	}

	logger, cleanup, err := logging.New(*logDir, "main", logrus.InfoLevel, logging.ModeBoth)
	if err != nil {
		fmt.Println("could not instantiate logger:", err.Error())
		return
	}
	defer cleanup()
	defer panicHandler.HandlePanic(logger)

	conf, created, err := configuration.Setup(*configFile)
	if err != nil {
		fmt.Println("error setting up configuration:", err.Error())
		return
	}
	if created {
		logger.Info("configuration file was created; exiting")
	}
	u, err := url.ParseRequestURI(conf.BindAddress)
	if err != nil {
		logger.WithField("error", err.Error()).Panic("invalid bind address")
	}

	setupBindAddr(u, &bindAddr)
	s := &http.Server{
		Addr:           bindAddr,
		Handler:        getRouter(conf, logger),
		ReadTimeout:    time.Second,
		WriteTimeout:   2 * time.Second,
		MaxHeaderBytes: 3 << 10,
	}
	s.SetKeepAlivesEnabled(false)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			panic("Could not gracefully shutdown the server: %v" + err.Error())
		}
	}()

	fmt.Printf("Starting up server with bind address %s...\n", bindAddr)
	if u.Scheme == "http" {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panic("Could not start server: " + err.Error())
		}
	} else if u.Scheme == "https" {
		s.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
		if err := s.ListenAndServeTLS(conf.CertificateFile, conf.KeyFile); err != nil && err != http.ErrServerClosed {
			logger.Panic("Could not start server with TLS: " + err.Error())
		}
	}
}

func setupBindAddr(u *url.URL, addr *string) {
	if u.Host != "" {
		*addr = u.Host
	} else if u.Port() != "" {
		*addr = ":" + u.Port()
	} else {
		*addr = ":" + global.DefaultPort
	}
}

func getRouter(appConfig *configuration.AppConfig, logger *logrus.Entry) *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "This resource could not be found", http.StatusNotFound)
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "This is the Maestro start page.", http.StatusNoContent)
	})

	bh := &handler.BaseHandler{
		Config:       appConfig,
		Logger:       logger,
		MaestroCache: cache.New(appConfig),
	}

	mwh := middleware.MWHandler{
		Config: appConfig,
		Logger: logger,
	}

	routerV1 := router.PathPrefix("/api/v1").Subrouter()
	routerV1.HandleFunc("/register", mwh.Auth(bh.RegistrationHandler)).Methods(http.MethodPost)
	routerV1.HandleFunc("/deregister", mwh.Auth(bh.DeregistrationHandler)).Methods(http.MethodDelete)
	routerV1.HandleFunc("/ping", mwh.Auth(bh.PingHandler)).Methods(http.MethodPut)
	routerV1.HandleFunc("/query", mwh.Auth(bh.QueryHandler)).Methods(http.MethodGet)

	return router
}
