package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/nik/mthen-api/internal/db"
	"github.com/nik/mthen-api/internal/handler"
	mw "github.com/nik/mthen-api/internal/middleware"
	"github.com/nik/mthen-api/internal/service"
)

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Database connection
	ctx := context.Background()
	if err := db.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL")

	// Services
	authSvc := service.NewAuthService()
	timelineSvc := service.NewTimelineService()
	albumSvc := service.NewAlbumService()
	artistSvc := service.NewArtistService()
	searchSvc := service.NewSearchService()
	userSvc := service.NewUserService()
	genreSvc := service.NewGenreService()

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	timelineH := handler.NewTimelineHandler(timelineSvc)
	albumH := handler.NewAlbumHandler(albumSvc)
	artistH := handler.NewArtistHandler(artistSvc)
	searchH := handler.NewSearchHandler(searchSvc)
	userH := handler.NewUserHandler(userSvc)
	genreH := handler.NewGenreHandler(genreSvc)

	// Router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(mw.Logger)
	r.Use(mw.CORSHandler())
	r.Use(chimiddleware.Recoverer)
	r.Use(mw.JSONContentType)

	// Health check
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		mw.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Auth routes (public)
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authH.Register)
		r.Post("/login", authH.Login)
	})

	// Public routes
	r.Route("/api/timeline", func(r chi.Router) {
		r.Get("/", timelineH.GetYears)
		r.Get("/{year}", timelineH.GetYear)
		r.Get("/{year}/{month}", timelineH.GetMonth)
	})

	r.Route("/api/albums", func(r chi.Router) {
		r.Get("/", albumH.List)
		r.Get("/{id}", albumH.GetByID)
	})

	r.Route("/api/artists", func(r chi.Router) {
		r.Get("/", artistH.List)
		r.Get("/{id}", artistH.GetByID)
		r.Get("/{id}/years", artistH.GetReleaseYears)
	})

	r.Get("/api/search", searchH.Search)
	r.Get("/api/genres", genreH.List)

	// Protected user routes
	r.Route("/api/me", func(r chi.Router) {
		r.Use(mw.JWTAuth)

		r.Get("/profile", userH.GetProfile)
		r.Put("/profile", userH.UpdateProfile)

		r.Get("/top-albums", userH.GetTopAlbums)
		r.Put("/top-albums", userH.UpdateTopAlbums)

		r.Get("/top-songs", userH.GetTopSongs)
		r.Put("/top-songs", userH.UpdateTopSongs)

		r.Get("/top-artists", userH.GetTopArtists)
		r.Put("/top-artists", userH.UpdateTopArtists)

		r.Post("/listening", userH.CreateListening)
		r.Get("/listening", userH.ListListening)

		r.Get("/monthly-set/{year}/{month}", userH.GetMonthlySet)
		r.Put("/monthly-set/{year}/{month}", userH.UpsertMonthlySet)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("MThen API server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
