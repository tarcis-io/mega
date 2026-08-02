package server

import (
	"net/http"
)

func Run() error {
	addr := ":8080"

	serveMux := http.NewServeMux()

	if err := http.ListenAndServe(addr, serveMux); err != nil {
		return err
	}

	return nil
}
