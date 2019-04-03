package webserver

import (
	"errors"
	"net/http"
	"strings"
)

func parseCommand(input string) (string, string, error) {
	// Split out command and arguments
	var command, args string
	var err error

	// commands are all 3 or 4 characters
	if len(input) < 3 {
		return command, args, errors.New("Syntax error")
	}

	response := strings.SplitAfterN(input, " ", 2)

	switch {
	case len(response) == 2:
		command = strings.TrimSpace(response[0])
		args = strings.TrimSpace(response[1])
	case len(response) == 1:
		command = response[0]
	}
	return command, args, err
}

func handleCommand(input string) (string, error) {
	// Handles input after authentication
	input = strings.TrimSpace(input)
	return "", nil
}

// StartServer - start REST API server
func StartServer() {
	http.HandleFunc("/", getCommandFromURL)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// getCommandFromURL - parce URL and get command to execute
func getCommandFromURL(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	// write to web-page
	w.Write([]byte(message))
	// return command
}
