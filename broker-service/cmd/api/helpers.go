package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) (err error) {
	var maxBytes = 104856

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	var dec = json.NewDecoder(r.Body)

	err = dec.Decode(data)
	if err != nil {
		return
	}

	//err = dec.Decode(&struct{}{})
	//if err != io.EOF {
	//	return errors.New("body must have only a single JSON value")
	//}

	return
}

func (app *Config) writeJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) (err error) {
	var out []byte
	out, err = json.Marshal(data)
	if err != nil {
		return
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)

	return
}

func (app *Config) errorJson(w http.ResponseWriter, errData error, status ...int) (err error) {
	var statusCode = http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = errData.Error()

	log.Println("Error :", errData.Error())
	err = app.writeJson(w, statusCode, payload)
	return
}

type HttpClientRequest struct {
	Url, Method string
	Payload     interface{}
	Headers     map[string]string
}

func (app *Config) HttpClient(params HttpClientRequest) (result map[string]interface{}, err error) {

	requestBody, err := json.Marshal(params.Payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest(params.Method, params.Url, bytes.NewBuffer(requestBody))

	if err != nil {
		return
	}
	if len(params.Headers) > 0 {
		for key, val := range params.Headers {
			req.Header.Set(key, val)
		}
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("error calling %s", params.Url)
		return
	}

	result = make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	return
}
