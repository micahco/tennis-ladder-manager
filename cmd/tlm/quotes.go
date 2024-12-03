package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/abiosoft/ishell/v2"
	"github.com/fatih/color"
	"github.com/pebbe/zmq4"
)

type QuotesRequest struct {
	Filters map[string]interface{} `json:"filters"`
}

type QuotesResponse struct {
	Code    int                    `json:"code"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
}

func (app *application) quote(c *ishell.Context) error {
	context, err := zmq4.NewContext()
	if err != nil {
		return err
	}
	defer context.Term()

	socket, err := context.NewSocket(zmq4.REQ)
	if err != nil {
		return err
	}
	defer socket.Close()

	err = socket.Connect("tcp://localhost:9876")
	if err != nil {
		return err
	}

	req := QuotesRequest{
		Filters: map[string]interface{}{},
	}

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// Send the request to the server
	_, err = socket.SendBytes(b, 0)
	if err != nil {
		return err
	}

	// Wait for a response
	res, err := socket.Recv(0)
	if err != nil {
		return err
	}

	var response QuotesResponse
	err = json.Unmarshal([]byte(res), &response)
	if err != nil {
		return err
	}

	if response.Code != 200 {
		return errors.New(response.Message)
	}

	q := fmt.Sprintln(response.Data["quote"], "\n\t-", response.Data["author"])
	italic := color.New(color.Italic).SprintFunc()
	app.shell.Println()
	app.shell.Println(italic(q))

	return nil
}
