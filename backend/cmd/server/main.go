package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrDisa/Raspy/backend/internal/handler"
)

func main() {
	r := chi.NewRouter()
	r.Get("/health", handler.Health)

	log.Println("server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
