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
	"github.com/silvan-talos/cookie-syncer/syncing"
)

const (
	defaultPort = ":8080"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(cookie_syncer_db)/syncer")
	if err != nil {
		log.Println("failed to connect to the database:", err)
		return
	}
	defer db.Close()

	partners := mysql.NewPartnerRepository(db)
	syncs := mysql.NewSyncRepository(db)
	ps := partner.NewService(partners)
	ss := syncing.NewService(syncs)
	server := server.New(ps, ss)
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
