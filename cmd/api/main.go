package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VJ-2303/email-worker/internal/mailer"
	"github.com/VJ-2303/email-worker/internal/worker"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
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

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.gmail.com", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-user", "vanaraj1018@gmail.com", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-pass", "your smtp password", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "test@email-worker.com", "SMTP sender email")

	flag.Parse()

	mailerInstance := mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	pool := worker.NewPool(4, 10, logger, mailerInstance)

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
