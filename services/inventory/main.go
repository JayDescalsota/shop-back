package main

import (
	"net/http"

	"backend/packages/service"

	"github.com/go-chi/chi/v5"
)

func main() {
	service.Serve("inventory", func(r chi.Router, deps *service.Dependencies) {
		r.Post("/graphql", handleInventory(deps))
		r.Post("/query", handleInventory(deps))
	})
}

func handleInventory(deps *service.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"__typename":"Query"}}`))
	}
}
