package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joaosantos/jlpt5/internal/api/handlers"
	"github.com/joaosantos/jlpt5/internal/api/routes"
	"github.com/joaosantos/jlpt5/internal/config"
	"github.com/joaosantos/jlpt5/internal/domain/services"
	"github.com/joaosantos/jlpt5/internal/infrastructure/database"
	"github.com/joaosantos/jlpt5/internal/infrastructure/postgres"
	"github.com/joaosantos/jlpt5/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger := utils.NewLogger(cfg.Log.Level)
	logger.Info("Starting JLPT5 Backend API", utils.WithContext(
		"version", "1.0.0",
		"port", cfg.Server.Port,
	))

	// Connect to database
	db, err := database.NewPostgresConnection(&cfg.Database, logger)
	if err != nil {
		logger.Error("Failed to connect to database", utils.WithContext("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		logger.Error("Failed to run migrations", utils.WithContext("error", err.Error()))
		os.Exit(1)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	vocabRepo := postgres.NewVocabularyRepository(db)
	grammarRepo := postgres.NewGrammarRepository(db)
	quizRepo := postgres.NewQuizRepository(db)

	// Initialize utilities
	jwtManager := utils.NewJWTManager(&cfg.JWT)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtManager, logger)
	spacedRepetitionService := services.NewSpacedRepetitionService()
	vocabService := services.NewVocabularyService(vocabRepo, spacedRepetitionService, logger)
	grammarService := services.NewGrammarService(grammarRepo, logger)
	quizService := services.NewQuizService(quizRepo, logger)
	progressService := services.NewProgressService(db, logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, logger)
	vocabHandler := handlers.NewVocabularyHandler(vocabService, logger)
	grammarHandler := handlers.NewGrammarHandler(grammarService, logger)
	quizHandler := handlers.NewQuizHandler(quizService, logger)
	progressHandler := handlers.NewProgressHandler(progressService, logger)

	// Setup routes
	router := routes.NewRouter(db, logger, authHandler, vocabHandler, grammarHandler, quizHandler, progressHandler)
	handler := router.SetupRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server listening", utils.WithContext("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", utils.WithContext("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", utils.WithContext("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Server stopped gracefully")
}
