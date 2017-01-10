package zabbix

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Params map[string]interface{}

type Command struct {
	JsonRPC string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Auth    *string     `json:"auth"`
	Params  interface{} `json:"params"`
	Method  string      `json:"method"`
}

func (client *Client) Base(method string, params interface{}) (Response, error) {
	command := Command{
		JsonRPC: "2.0",
		Id:      0,
		Params:  params,
		Method:  method,
	}

	if method != "user.login" {
		command.Auth = client.token
	}

	encoded, err := json.Marshal(command)
	if err != nil {
		return Response{}, err
	}

	poster := &http.Client{Timeout: client.timeout}
	reader := bytes.NewReader(encoded)
	response, err := poster.Post(client.uri, "application/json", reader)
	if err != nil {
		return Response{}, err
	}

	return ResponseFromReader(response.Body)
}

func (client *Client) Call(method string, params interface{}) (Response, error) {
	response, err := client.Base(method, params)
	if _, ok := err.(AuthenticationError); ok {
		client.token = nil
		if err = client.Login(); err != nil {
			return response, err
		}

		response, err = client.Base(method, params)
	}

	return response, err
}
