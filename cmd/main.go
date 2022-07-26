package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/silvan-talos/cookie-syncer/server"
)

const (
	defaultPort = ":8080"
)

func main() {
	server := server.New()
	errs := make(chan error, 2)
	go func() {
		log.Println("starting server")
		errs <- http.ListenAndServe(defaultPort, server)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Println("server shutting down:", <-errs)
}
