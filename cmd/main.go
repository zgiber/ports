package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zgiber/ports/service"
	"github.com/zgiber/ports/store"
	"github.com/zgiber/ports/transport/rest"
)

var (
	addr   = flag.String("addr", ":8000", "the address where the REST service listens for incoming calls")
	dbPath = flag.String("db_path", "", "the path to the directory for the database")

	shutdownTimeout = 30 * time.Second
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())

	store, err := store.New(*dbPath)
	if err != nil {
		log.Fatal(err)
	}

	rest := rest.NewServer(ctx, service.New(store), *addr)

	go listenForSignals(func() error {
		cancel()
		shutdownCtx, cf := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cf()
		return rest.Shutdown(shutdownCtx)
	})

	if err := rest.ListenAndServe(); err != nil {
		cancel()
		store.Close()
		log.Fatal(err)
	}
}

func listenForSignals(shutdown func() error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	signal.Stop(c)
	if err := shutdown(); err != nil {
		log.Println("error during shutdown", err.Error())
	}
}
