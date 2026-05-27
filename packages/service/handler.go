package service

import (
	"encoding/json"
	"net/http"
)

type GQLRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

type GQLResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []GQLError  `json:"errors,omitempty"`
}

type GQLError struct {
	Message string   `json:"message"`
	Path    []string `json:"path,omitempty"`
}

func WriteGQL(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GQLResponse{Data: data})
}

func WriteGQLError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(GQLResponse{
		Errors: []GQLError{{Message: message}},
	})
}
