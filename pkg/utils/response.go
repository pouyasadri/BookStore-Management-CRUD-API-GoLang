package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func RespondWithError(w http.ResponseWriter, code int, message string, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
	json.NewEncoder(w).Encode(response)
}

func RespondWithSuccess(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := SuccessResponse{
		Code: code,
		Data: data,
	}
	json.NewEncoder(w).Encode(response)
}

func RespondWithMessage(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := map[string]interface{}{
		"code":    code,
		"message": message,
	}
	json.NewEncoder(w).Encode(response)
}
