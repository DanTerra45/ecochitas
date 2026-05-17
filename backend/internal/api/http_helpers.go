package api

import (
	"encoding/json"
	"net/http"
)

type Error_response struct {
	Message string `json:"message"`
}

func write_json_response(response_writer http.ResponseWriter, http_status_code int, response_payload any) {
	response_writer.Header().Set("Content-Type", "application/json")
	response_writer.WriteHeader(http_status_code)
	_ = json.NewEncoder(response_writer).Encode(response_payload)
}

func write_json_error(response_writer http.ResponseWriter, http_status_code int, error_message string) {
	write_json_response(response_writer, http_status_code, Error_response{
		Message: error_message,
	})
}
