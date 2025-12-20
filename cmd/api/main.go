package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VJ-2303/email-worker/internal/worker"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config     config
	logger     *log.Logger
	workerPool *worker.Pool
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "devolopment", "Environment (devolopment|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	pool := worker.NewPool(4, 10, logger)

	pool.Run()

	app := &application{
		config:     cfg,
		logger:     logger,
		workerPool: pool,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheclHandler)
	mux.HandleFunc("/v1/send", app.sendEmailHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
