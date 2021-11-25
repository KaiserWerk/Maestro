package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/KaiserWerk/Maestro/internal/global"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/KaiserWerk/Maestro/internal/configuration"
	"github.com/KaiserWerk/Maestro/internal/logging"
	"github.com/KaiserWerk/Maestro/internal/shutdownManager"

	"github.com/gorilla/mux"
)

var (
	bindAddr string
	configFile = flag.String("config", "", "The configuration file to use")
	logDir = flag.String("logDir", ".", "The directory to save log files to")
)

func main() {
	flag.Parse()

	//defer panicHandler.HandlePanic()
	defer shutdownManager.Initiate()
	logging.Init(*logDir)

	logger := logging.New(logrus.InfoLevel, "main", logging.ModeBoth)

	if *configFile != "" {
		configuration.SetFile(*configFile)
	}

	conf, created, err := configuration.Setup()
	if err != nil {
		logger.WithField("error", err.Error()).Panic("error setting up configuration")
	}
	if created {
		logger.Info("configuration file was created; exiting")
	}

	u, err := url.ParseRequestURI(conf.BindAddress)
	if err != nil {
		logger.WithField("error", err.Error()).Panic("invalid bind address")
	}

	fmt.Printf("bind addr: %s, %s, %s\n", u.Scheme, u.Host, u.Port())

	setupBindAddr(u, &bindAddr)






	//c := cache.New(5 * time.Minute)
	//
	//// Put something into the cache
	//c.Set("a", "b", 1 * time.Minute)
	//
	//// Read from the cache
	//obj, found := c.Get("a")
	//if found {
	//	// Convert the type
	//	fmt.Println(obj.(string))
	//}

	router := getRouter()
	s := &http.Server{
		Addr: bindAddr,
		Handler: router,
		ReadTimeout: time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout: 2 * time.Second,
		MaxHeaderBytes: 3 << 10,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		s.SetKeepAlivesEnabled(false)
		if err := s.Shutdown(ctx); err != nil {
			panic("Could not gracefully shutdown the server: %v" + err.Error())
		}
	}()

	fmt.Printf("Starting up server with binding address %s...\n", bindAddr)
	if u.Scheme == "http" {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panic("Could not start server")
		}
	} else if u.Scheme == "https" {
		if err := s.ListenAndServeTLS(conf.CertificateFile, conf.KeyFile); err != nil && err != http.ErrServerClosed {
			logger.Panic("Could not start server")
		}
	}


}

func setupBindAddr(u *url.URL, addr *string) {
	if u.Host != "" {
		*addr = u.Host
	} else if u.Port() != "" {
		*addr = ":" + u.Port()
	} else {
		*addr = ":" + global.GetDefaultPort()
	}
}

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/register", nil).Methods("POST")
	router.Handle("/deregister", nil).Methods("POST")
	router.Handle("/query", nil).Methods("GET")
	router.Handle("/queryall", nil).Methods("GET")

	return router
}