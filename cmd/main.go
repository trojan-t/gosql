package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/trojan-t/gosql/cmd/app"
	"github.com/trojan-t/gosql/pkg/customers"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@localhost:5432/db"
	if err := execute(host, port, dsn); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string, dsn string) (err error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()
	mux := http.NewServeMux()
	customerSvs := customers.NewService(db)
	serverHandler := app.NewServer(mux, customerSvs)
	serverHandler.Init()

	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: serverHandler,
	}
	return srv.ListenAndServe()
}
