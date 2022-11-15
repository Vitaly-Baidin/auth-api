package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Response Base response
type Response struct {
	StatusCode int            `json:"-"`
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

func (response *Response) SendResponse(rw http.ResponseWriter) {
	rw.WriteHeader(response.StatusCode)
	rw.Header().Add("Content-Type", "application/json")

	marshal, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marhsal failed: %v", err)
		response.StatusCode = http.StatusInternalServerError
		newMessage := fmt.Sprintf("%s: %v", response.Message, err)
		response.Message = newMessage
	}

	_, err = rw.Write(marshal)
	if err != nil {
		log.Printf("Write to ResponseWriter failed: %v", err)
	}
}

func SendResponseData(rw http.ResponseWriter, data map[string]any) {
	response := &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       data,
	}
	response.SendResponse(rw)
}

func SendErrorResponse(rw http.ResponseWriter, status int, message string) {
	response := &Response{
		StatusCode: status,
		Success:    false,
		Message:    message,
	}
	response.SendResponse(rw)
}
