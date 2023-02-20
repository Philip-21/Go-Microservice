package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"log"
	"net/http"
)

func (app *Config) AuthTemplate(body []byte, w http.ResponseWriter, r *http.Request) (AuthPayload, error) {
	var data AuthPayload
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// store data in the struct
		data := AuthPayload{
			Email:    email,
			Password: password,
		}
		fmt.Fprintf(w, "Email: %s, Password: %s\n", data.Email, data.Password)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Invalid request method")
	}
	return data, nil
}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(w http.ResponseWriter, r *http.Request) {
	// Make an HTTP request to retrieve the template from the front-end service
	resp, err := http.Get("http://front-end-service/auth")
	if err != nil {
		log.Printf("error in getting service", err)
	}
	defer resp.Body.Close()

	// Read the response body into a string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//create some json to send to the Authentication Microservice
	auth, err := app.AuthTemplate(body, w, r)
	if err != nil {
		log.Println(err)
	}
	jsonData, _ := json.MarshalIndent(auth, "", "\t")

	//call the Authentication service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData)) //must be the same as the docker compose yml
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, errors.New("Error in Calling Auth Service"))
		return
	}
	//create a variable We'll read response.body into
	var jsonFromService jsonResponse

	//decode the json from the Auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		log.Println(err)
		app.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	///writes the actual data from the service
	app.Writejson(w, http.StatusAccepted, payload)

}
func (app *Config) SignUpTemplate(body []byte, w http.ResponseWriter, r *http.Request) (SignUpPayload, error) {
	var sign SignUpPayload

	if r.Method == http.MethodPost {
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		email := r.FormValue("email")
		password := r.FormValue("password")

		sign := SignUpPayload{
			FirstName: firstname,
			LastName:  lastname,
			Email:     email,
			Password:  password,
		}
		fmt.Fprintf(w, "FirstName:%s", "LastName:%s", "Email: %s, Password: %s\n", sign.FirstName, sign.LastName, sign.Email, sign.Password)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return sign, nil

}

func (app *Config) SignUp(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://front-end-service/signup")
	if err != nil {
		log.Printf("Error in getting service", err)

	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

	}
	sign, err := app.SignUpTemplate(body, w, r)
	if err != nil {
		log.Println(err)
	}

	jsonData, _ := json.MarshalIndent(sign, "", "/t")

	request, err := http.NewRequest("POST", "http://authentication-service/signup", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("error in requesting signup service", err)
		app.ErrorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		app.ErrorJSON(w, errors.New("Invalid inputs"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, errors.New("Error in calling signup handlers from auth service"))
		return
	}
	//create a variable We'll read response.body into
	var jsonFromService jsonResponse
	//decode the json from the Auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		log.Println(err)
		app.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Registered!"
	payload.Data = jsonFromService.Data

	///writes the actual data from the service
	app.Writejson(w, http.StatusAccepted, payload)

}
