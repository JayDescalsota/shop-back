package main

import (
	"net/http"

	"backend/packages/service"

	"github.com/go-chi/chi/v5"
)

func main() {
	service.Serve("payroll", func(r chi.Router, deps *service.Dependencies) {
		r.Post("/graphql", handlePayroll(deps))
		r.Post("/query", handlePayroll(deps))
	})
}

func handlePayroll(deps *service.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"__typename":"Query"}}`))
	}
}
