package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"` //info needed to authenticate
	Log    LogPayload  `json:"log,omitempty"`  //info needed by to show a user is logged in
	//LogGRPc LogPayload  `json:"loggrpc,omitempty"`
	Mail MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
type SignUpPayload struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"email"`
	Data string `json:"data"`
}

// the output diplayed in the frontend
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` //used any data type  instead of an interface , cause im parsin  insmall values
}

// Read json
func (app *Config) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	//make limitation of the size of uploaded json file
	maxBytes := 1048576 //one megabyte

	//request body is executed on the serverside
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	//unpack data from json to struct
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF { //EOF is the error returned by Read when no more input is available
		return errors.New("body must have single json value")
	}
	return nil

}

// Write Json
func (app *Config) Writejson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

// error json
func (app *Config) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	//check to see if status has been specified
	if len(status) > 0 { //specified
		statusCode = status[0]
	}

	var payload *jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.Writejson(w, statusCode, payload)

}
