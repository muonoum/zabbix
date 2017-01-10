package zabbix

import (
	"fmt"
	"sync"
	"time"
)

type Client struct {
	sync.Mutex
	user     string
	password string
	uri      string
	token    *string
	timeout  time.Duration
}

func New(uri, user, password string, timeout time.Duration) *Client {
	return &Client{
		user:     user,
		password: password,
		uri:      uri,
		token:    nil,
		timeout:  timeout,
	}
}

func (client *Client) Login() error {
	defer client.Unlock()
	client.Lock()

	if client.token != nil {
		return nil
	}

	response, err := client.Call("user.login", Params{
		"user": client.user, "password": client.password,
	})
	if err != nil {
		return err
	}

	if err = response.Decode(&client.token); err != nil {
		return fmt.Errorf("Could not decode authentication token: %s", err)
	}

	return nil
}

func (client *Client) Logout() error {
	defer client.Unlock()
	client.Lock()

	if client.token == nil {
		return nil
	}

	_, err := client.Call("user.logout", Params{"token": *client.token})
	return err
}

func (client *Client) Decode(method string, object interface{}, params interface{}) error {
	response, err := client.Call(method, params)
	if err != nil {
		return err
	}

	return response.Decode(object)
}
