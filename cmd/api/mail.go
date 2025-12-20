package main

import (
	"net/http"

	"github.com/VJ-2303/email-worker/internal/worker"
)

type EmailRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (app *application) sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var input EmailRequest

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if input.To == "" || input.From == "" {
		http.Error(w, "To and sbject fields are required", http.StatusBadRequest)
		return
	}
	task := worker.Task{
		From:    input.From,
		To:      input.To,
		Subject: input.Subject,
		Body:    input.Body,
	}
	app.workerPool.Submit(task)

	app.writeJSON(w, http.StatusAccepted, map[string]string{
		"status":  "queued",
		"message": "Email has been added to the queue",
	})
}
