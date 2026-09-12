package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/mrDisa/Raspy/backend/internal/handler"
	"github.com/mrDisa/Raspy/backend/internal/repository"
)

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://raspy:raspy_dev_password@localhost:5432/raspy?sslmode=disable"
	}
	db, err := repository.NewDB(connString)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	userRepo := repository.NewPostgresUserRepository(db)
	_ = userRepo
	
	r := chi.NewRouter()
	r.Get("/health", handler.Health)

	log.Println("server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
