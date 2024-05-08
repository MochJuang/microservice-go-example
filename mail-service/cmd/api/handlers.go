package main

import (
	"log"
	"net/http"
)

type MailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	var requestPaylaod MailMessage

	log.Println("start request send email")
	err := app.readJSON(w, r, &requestPaylaod)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	log.Println(requestPaylaod)

	msg := Message{
		From:    requestPaylaod.From,
		To:      requestPaylaod.To,
		Subject: requestPaylaod.Subject,
		Data:    requestPaylaod.Message,
	}

	err = app.Mailer.SenSMTPMessage(msg)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "send to " + requestPaylaod.To,
	}

	log.Println("success send email")
	app.writeJson(w, http.StatusAccepted, response)
}
