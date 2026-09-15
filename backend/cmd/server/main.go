package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	
	"github.com/go-chi/chi/v5"

	"github.com/mrDisa/Raspy/backend/internal/auth"
	"github.com/mrDisa/Raspy/backend/internal/collegeapi"
	"github.com/mrDisa/Raspy/backend/internal/handler"
	"github.com/mrDisa/Raspy/backend/internal/repository"
	"github.com/mrDisa/Raspy/backend/internal/service"
)

func openDB(databaseURL string) (*sql.DB, error) {

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	db, err := openDB(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Repositories

	userRepo := repository.NewUserRepository(db)
	groupRepo := repository.NewGroupRepository(db)

	// College API

	collegeClient := collegeapi.NewClient(
		"https://api.stavmk.ru/widget/schedule/",
	)

	// Services

	scheduleService := service.NewScheduleService(collegeClient)

	// Handlers

	scheduleHandler := handler.NewScheduleHandler(
		scheduleService,
		groupRepo,
	)

	// Router

	r := chi.NewRouter()

	r.Get("/health", handler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(auth.Middleware(userRepo, botToken))

		r.Get("/me", handler.Me)

		r.Route("/schedule", func(r chi.Router) {
			r.Get("/today", scheduleHandler.Today)
			r.Get("/tomorrow", scheduleHandler.Tomorrow)
		})
	})

	log.Println("server started on :3000")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}