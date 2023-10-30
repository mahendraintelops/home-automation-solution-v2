package client

import (
	"io"
	"net/http"
)

func ExecuteNode14() ([]byte, error) {
	response, err := http.Get("device-service:4550/ping")

	if err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}
