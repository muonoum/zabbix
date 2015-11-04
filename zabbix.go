package zabbix

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

const (
	OK = iota
	Problem
)

var Priorities = []string{
	"unclassified",
	"info",
	"warning",
	"average",
	"high",
	"disaster",
}

type Params map[string]interface{}

func (p Params) WithGroups(groups []int) Params {
	if len(groups) > 0 {
		p["groupids"] = groups
	}

	return p
}

var MaxPriority = len(Priorities)

type Command struct {
	JsonRPC string  `json:"jsonrpc"`
	ID      int     `json:"id"`
	Auth    *string `json:"auth"`
	Params  Params  `json:"params"`
	Method  string  `json:"method"`
}

type Client struct {
	user     string
	password string
	uri      string
	token    *string
	mu       sync.Mutex
}

type responseData struct {
	JsonRPC string
	ID      int
	Result  *json.RawMessage
	Error   *json.RawMessage
}

type Response struct {
	responseData
}

func (r *Response) Decode(object interface{}) error {
	if r.Result == nil {
		return errors.New("Result property is nil")
	} else if err := json.Unmarshal(*r.Result, object); err != nil {
		return fmt.Errorf("Could not decode result property: %s", err)
	} else {
		return nil
	}
}

type responseError struct {
	Code    int
	Message string
	Data    string
}

func (e responseError) Error() string {
	return e.Data
}

type authError struct {
	message string
}

func (e authError) Error() string {
	return e.message
}
