package main

import (
	"encoding/json"

	"github.com/pebbe/zmq4"
)

func (app *application) validate_set(aGamesWon, bGamesWon int) (bool, error) {
	context, err := zmq4.NewContext()
	if err != nil {
		return false, err
	}
	defer context.Term()

	socket, err := context.NewSocket(zmq4.REQ)
	if err != nil {
		return false, err
	}
	defer socket.Close()

	err = socket.Connect("tcp://localhost:9878")
	if err != nil {
		return false, err
	}

	req := []int{aGamesWon, bGamesWon}

	b, err := json.Marshal(req)
	if err != nil {
		return false, err
	}

	// Send the request to the server
	_, err = socket.SendBytes(b, 0)
	if err != nil {
		return false, err
	}

	// Wait for a response
	res, err := socket.Recv(0)
	if err != nil {
		return false, err
	}

	var isValid bool
	err = json.Unmarshal([]byte(res), &isValid)
	if err != nil {
		return false, err
	}

	return isValid, nil
}
