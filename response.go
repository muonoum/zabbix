package zabbix

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const (
	errorSessionTerminated = "Session terminated, re-login, please."
	errorNotAuthorised     = "Not authorized."
)

type ResponseError struct {
	Code    int
	Message string
	Data    string
}

func (err ResponseError) Error() string {
	return err.Data
}

type AuthenticationError struct {
	message string
}

func (err AuthenticationError) Error() string {
	return err.message
}

type Response struct {
	JsonRPC string
	Id      int
	Result  *json.RawMessage
	Error   *ResponseError
}

func ResponseFromReader(reader io.ReadCloser) (rsp Response, _ error) {
	defer reader.Close()
	if err := json.NewDecoder(reader).Decode(&rsp); err != nil {
		return rsp, err
	} else if rsp.Error == nil {
		return rsp, nil
	}

	switch err := rsp.Error; err.Error() {
	case errorSessionTerminated, errorNotAuthorised:
		return rsp, AuthenticationError{err.Error()}
	default:
		return rsp, err
	}
}

func (r Response) Decode(object interface{}) error {
	if r.Result == nil {
		return errors.New("Result property is nil")
	} else if err := json.Unmarshal(*r.Result, object); err != nil {
		return fmt.Errorf("Could not decode result property: %s", err)
	}

	return nil
}
