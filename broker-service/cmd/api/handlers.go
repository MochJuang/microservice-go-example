package main

import (
	"broker/logs/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type RequestPaylaod struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

const (
	ActionAuth = "auth"
	ActionLog  = "log"
	ActionMail = "mail"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var payload = jsonResponse{
		Error:   false,
		Message: "Hit the broker service",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPaylaod
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println("error read payload")
		log.Println(err.Error())
		app.errorJson(w, err)
		return
	}

	switch requestPayload.Action {
	case ActionAuth:
		app.authenticate(w, requestPayload.Auth)
	case ActionLog:
		app.logItemViaListener(w, requestPayload.Log)
	case ActionMail:
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJson(w, errors.New("unknow action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json for authenticate to auth service
	jsonData, _ := json.Marshal(a)
	path := fmt.Sprintf("http://%s:%s/authenticate", os.Getenv("BASE_URL"), os.Getenv("AUTH_SERVICE_PORT"))
	request, err := http.NewRequest("POST", path, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	}

	if response.StatusCode != http.StatusAccepted {
		log.Println(response)
		app.errorJson(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}

	var payload = jsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    jsonFromService.Data,
	}

	log.Println("success response")
	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {
	// create some json for authenticate to auth service
	request := HttpClientRequest{
		Url:     fmt.Sprintf("http://%s:%s/log", os.Getenv("BASE_URL"), os.Getenv("LOGGER_SERVICE_PORT")),
		Method:  "POST",
		Payload: l,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := app.HttpClient(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	var payload = jsonResponse{
		Error:   false,
		Message: "logged!",
	}

	log.Println("success response")
	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logItemViaListener(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	var payload = jsonResponse{
		Error:   false,
		Message: "logged!",
	}

	log.Println("success response")
	app.writeJson(w, http.StatusAccepted, payload)
}

func (app Config) logItemViaGPRC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPaylaod

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	log.Println(requestPayload)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	_, err = app.GRPCCLient.LogClient.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJson(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged from grpc"

	app.writeJson(w, http.StatusOK, payload)

}

func (app *Config) pushToQueue(name, msg string) error {
	logPayload := LogPayload{
		Name: name,
		Data: msg,
	}

	payload, _ := json.Marshal(logPayload)
	return app.Emitter.Push("log.INFO", string(payload))
}

func (app Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	request := HttpClientRequest{
		Url:     fmt.Sprintf("http://%s:%s/send", os.Getenv("BASE_URL"), os.Getenv("MAIL_SERVICE_PORT")),
		Method:  "POST",
		Payload: msg,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	response, err := app.HttpClient(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	successMessage := fmt.Sprintf("send email to %s is success", msg.To)
	log.Println(successMessage)
	app.writeJson(w, http.StatusAccepted, jsonResponse{
		Error:   false,
		Message: successMessage,
		Data:    response,
	})

}
