package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VJ-2303/email-worker/internal/mailer"
	"github.com/VJ-2303/email-worker/internal/worker"
	"github.com/joho/godotenv"
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
	godotenv.Load()
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "devolopment", "Environment (devolopment|staging|production)")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.gmail.com", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-user", "vanaraj1018@gmail.com", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-pass", os.Getenv("SMTPPASS"), "SMTP password")
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

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Printf("server starting on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)

	sig := <-quit

	logger.Printf("Caught Signal: %s, Shutting down server....", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Forced Shutdown")
	}

	logger.Println("Waiting for Email workers to finish")
	pool.Shutdown()
	logger.Println("Email workers stopped")

	logger.Println("Server stopped")
}
