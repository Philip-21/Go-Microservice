package main

import (
	"net/http"
)

// the homepage handler for the frontend
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.Writejson(w, http.StatusOK, payload)
}
