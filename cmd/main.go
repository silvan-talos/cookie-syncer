package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	"github.com/silvan-talos/cookie-syncer/mysql"
	"github.com/silvan-talos/cookie-syncer/partner"
	"github.com/silvan-talos/cookie-syncer/server"
)

const (
	defaultPort = ":8080"
)

func main() {
	db, err := sql.Open("mysql", "root:root@/syncer")
	if err != nil {
		log.Println("failed to connect to the database:", err)
		return
	}
	partners := mysql.NewPartnerRepository(db)
	ps := partner.NewService(partners)
	server := server.New(ps)
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
