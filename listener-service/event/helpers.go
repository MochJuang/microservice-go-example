package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type HttpClientRequest struct {
	Url, Method string
	Payload     interface{}
	Headers     map[string]string
}

func HttpClient(params HttpClientRequest) (result map[string]interface{}, err error) {

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
