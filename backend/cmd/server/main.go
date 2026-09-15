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

	"github.com/go-chi/cors"
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

	devMode := os.Getenv("DEV_MODE") == "true"

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
	groupService := service.NewGroupService(collegeClient, groupRepo)

	// Handlers

	scheduleHandler := handler.NewScheduleHandler(
		scheduleService,
		groupRepo,
	)
	groupHandler := handler.NewGroupHandler(groupService)

	// Router

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Telegram-Init-Data"},
	}))

	r.Get("/health", handler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		if devMode {
			r.Use(auth.DevMiddleware(userRepo))
		} else {
			r.Use(auth.Middleware(userRepo, botToken))
		}

		r.Get("/me", handler.Me)

		r.Get("/groups", groupHandler.List)

		r.Put("/me/group", handler.UpdateGroup(userRepo, groupRepo))

		r.Route("/schedule", func(r chi.Router) {
			r.Get("/today", scheduleHandler.Today)
			r.Get("/tomorrow", scheduleHandler.Tomorrow)
			r.Get("/week", scheduleHandler.Week)
			r.Get("/week/next", scheduleHandler.NextWeek)
		})
	})

	log.Println("server started on :3000")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}