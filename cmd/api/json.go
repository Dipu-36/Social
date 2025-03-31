package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init(){
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1_048_578 //one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	//Here above we are limiting the size of the request body that will be read to prevent DDOS attacks 

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}
//When building an API responses  should be consisitent meaning that whatevere the client is consuming it should expect the error in the same 
func writeJSONError( w http.ResponseWriter, status int, message string) error {
	type envelope struct{
		Error string `json:"error"`
	}
	return writeJSON(w, status, &envelope{Error: message})
}