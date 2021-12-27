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
	"github.com/KaiserWerk/Maestro/internal/shutdownManager"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	bindAddr   string
	configFile = flag.String("config", "", "The configuration file to use")
	logDir     = flag.String("logDir", ".", "The directory to save log files to")
)

func main() {
	flag.Parse()

	logging.Init(*logDir)

	logger := logging.New(logrus.InfoLevel, "main", logging.ModeBoth)
	defer shutdownManager.Initiate()
	defer panicHandler.HandlePanic(logger)

	if *configFile != "" {
		configuration.SetFile(*configFile)
	}

	conf, created, err := configuration.Setup()
	if err != nil {
		logger.WithField("error", err.Error()).Panic("error setting up configuration")
	}
	if created {
		logger.Info("configuration file was created; exiting")
		return
	}
	cache.Init(conf)

	u, err := url.ParseRequestURI(conf.App.BindAddress)
	if err != nil {
		logger.WithField("error", err.Error()).Panic("invalid bind address")
	}

	setupBindAddr(u, &bindAddr)
	s := &http.Server{
		Addr:           bindAddr,
		Handler:        getRouter(),
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
		if err := s.ListenAndServeTLS(conf.App.CertificateFile, conf.App.KeyFile); err != nil && err != http.ErrServerClosed {
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

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "This resource could not be found", http.StatusNotFound)
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "This is the Maestro start page.", http.StatusNoContent)
	})

	hd := &handler.HttpHandler{
		Logger: logging.New(logrus.InfoLevel, "main", logging.ModeBoth),
	}

	routerV1 := router.PathPrefix("/api/v1").Subrouter()

	routerV1.HandleFunc("/register", middleware.Auth(hd.RegistrationHandler)).Methods(http.MethodPost)
	routerV1.HandleFunc("/deregister", middleware.Auth(hd.DeregistrationHandler)).Methods(http.MethodDelete)
	routerV1.HandleFunc("/ping", middleware.Auth(hd.PingHandler)).Methods(http.MethodPut)
	routerV1.HandleFunc("/query", middleware.Auth(hd.QueryHandler)).Methods(http.MethodGet)

	return router
}
