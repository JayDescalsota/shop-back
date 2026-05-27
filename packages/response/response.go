package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type APIResponse struct {
	Data      interface{} `json:"data,omitempty"`
	Errors    []APIError  `json:"errors,omitempty"`
	Meta      Meta        `json:"meta"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type Meta struct {
	RequestID string    `json:"requestId"`
	Timestamp time.Time `json:"timestamp"`
	TenantID  string    `json:"tenantId,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	resp := APIResponse{
		Data: data,
		Meta: Meta{
			RequestID: uuid.New().String(),
			Timestamp: time.Now().UTC(),
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func Error(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	resp := APIResponse{
		Errors: []APIError{
			{Code: code, Message: message},
		},
		Meta: Meta{
			RequestID: uuid.New().String(),
			Timestamp: time.Now().UTC(),
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func ValidationError(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnprocessableEntity)

	apiErrors := make([]APIError, 0, len(errors))
	for field, msg := range errors {
		apiErrors = append(apiErrors, APIError{
			Code:    "VALIDATION_ERROR",
			Message: msg,
			Field:   field,
		})
	}

	resp := APIResponse{
		Errors: apiErrors,
		Meta: Meta{
			RequestID: uuid.New().String(),
			Timestamp: time.Now().UTC(),
		},
	}

	json.NewEncoder(w).Encode(resp)
}
