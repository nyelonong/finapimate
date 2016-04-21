package jsonapi

import (
	"fmt"
	"log"

	"encoding/json"
	"net/http"
)

type Link struct {
	Self string `json:"self"`
	Prev string `json:"prev"`
	Next string `json:"next"`
}

type Data struct {
	Id        string      `json:"id"`
	Type      string      `json:"type,omitempty"`
	Attribute interface{} `json:"attributes,omitempty"`
}

type ResponseMultiData struct {
	Link Link   `json:"links,omitempty"`
	Data []Data `json:"data,omitempty"`
}

type ResponseData struct {
	Data Data `json:"data,omitempty"`
}

type Error struct {
	Message string `json:"message"`
}

func SuccessWriter(res http.ResponseWriter, dataResponse interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)

	dataBytes, err := json.Marshal(dataResponse)
	if err != nil {
		log.Fatalln(err)
	}
	res.Write(dataBytes)
}

func ErrorsWriter(res http.ResponseWriter, status int, message string) {

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Status-Code", fmt.Sprintf("%d", status))
	res.Header().Set("Access-Control-Allow-Headers", "Content-type, Authorization")
	res.WriteHeader(status)

	errorsResponse := Error{
		Message: message,
	}
	errorsBytes, err := json.Marshal(errorsResponse)
	if err != nil {
		log.Fatalln(err)
	}

	res.Write(errorsBytes)
}
