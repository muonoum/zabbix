package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/tj/go-debug"
)

var debug = Debug("zabbix")

func NewClient(uri, user, password string) *Client {
	return &Client{
		user:     user,
		password: password,
		uri:      uri,
		token:    nil,
	}
}

func (c *Client) SetToken(token string) {
	c.token = &token
}

func (c *Client) Login() error {
	debug("Trying to log in as `%s'", c.user)

	defer c.mu.Unlock()
	c.mu.Lock()

	if c.token == nil {
		if res, err := c.Base("user.login", Params{"user": c.user, "password": c.password}); err != nil {
			return err
		} else if err := res.Decode(&c.token); err != nil {
			return fmt.Errorf("Could not decode authentication token: %s", err)
		}

		debug("Logged in as `%s', token `%s'", c.user, *c.token)
	} else {
		debug("Using existing token `%s'", *c.token)
	}

	return nil
}

func (c *Client) Base(method string, params Params) (*Response, error) {
	cmd := &Command{JsonRPC: "2.0", ID: 0, Params: params, Method: method}

	if method != "user.login" {
		cmd.Auth = c.token
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(c.uri, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var res *Response
	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		return nil, err
	}

	if res.Error != nil {
		var e responseError
		if err := json.Unmarshal(*res.Error, &e); err != nil {
			return res, fmt.Errorf("Could not decode error object: %s", err)
		}

		switch e.Data {
		case "Session terminated, re-login, please.":
			return res, authError{e.Data}
		case "Not authorised.":
			return res, authError{e.Data}
		default:
			return res, e
		}
	}

	return res, nil
}

func (c *Client) Call(method string, params Params) (res *Response, err error) {
	res, err = c.Base(method, params)
	if _, ok := err.(authError); ok {
		debug("Authentication error: %s", err)

		c.token = nil

		if err = c.Login(); err != nil {
			return
		} else if res, err = c.Base(method, params); err != nil {
			return
		}
	} else if err != nil {
		return
	}

	return
}

func (c *Client) Decode(method string, object interface{}, params Params) error {
	if res, err := c.Call(method, params); err != nil {
		return err
	} else {
		return res.Decode(object)
	}
}
