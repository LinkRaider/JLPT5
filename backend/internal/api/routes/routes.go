package routes

import (
	"encoding/json"
	"net/http"

	"github.com/joaosantos/jlpt5/internal/api/handlers"
	"github.com/joaosantos/jlpt5/internal/infrastructure/database"
	"github.com/joaosantos/jlpt5/internal/utils"
)

// Router holds the dependencies for routing
type Router struct {
	db              *database.DB
	logger          *utils.Logger
	authHandler     *handlers.AuthHandler
	vocabHandler    *handlers.VocabularyHandler
	grammarHandler  *handlers.GrammarHandler
	quizHandler     *handlers.QuizHandler
	progressHandler *handlers.ProgressHandler
}

// NewRouter creates a new router with dependencies
func NewRouter(
	db *database.DB,
	logger *utils.Logger,
	authHandler *handlers.AuthHandler,
	vocabHandler *handlers.VocabularyHandler,
	grammarHandler *handlers.GrammarHandler,
	quizHandler *handlers.QuizHandler,
	progressHandler *handlers.ProgressHandler,
) *Router {
	return &Router{
		db:              db,
		logger:          logger,
		authHandler:     authHandler,
		vocabHandler:    vocabHandler,
		grammarHandler:  grammarHandler,
		quizHandler:     quizHandler,
		progressHandler: progressHandler,
	}
}

// SetupRoutes configures all application routes
func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoints
	mux.HandleFunc("/health", r.healthCheckHandler)
	mux.HandleFunc("/ready", r.readinessCheckHandler)

	// Authentication routes
	mux.HandleFunc("POST /api/v1/auth/register", r.authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", r.authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", r.authHandler.RefreshToken)

	// Vocabulary routes
	mux.HandleFunc("GET /api/v1/vocabulary", r.vocabHandler.ListVocabulary)
	mux.HandleFunc("GET /api/v1/vocabulary/due", r.vocabHandler.GetDueVocabulary)
	mux.HandleFunc("/api/v1/vocabulary/", r.vocabHandler.GetVocabulary) // Handles GET /api/v1/vocabulary/{id}
	mux.HandleFunc("POST /api/v1/vocabulary/", r.vocabHandler.SubmitReview) // Handles POST /api/v1/vocabulary/{id}/review

	// Grammar routes
	mux.HandleFunc("GET /api/v1/grammar", r.grammarHandler.ListGrammar)
	mux.HandleFunc("GET /api/v1/grammar/{id}", r.grammarHandler.GetGrammarLesson)
	mux.HandleFunc("POST /api/v1/grammar/{id}/complete", r.grammarHandler.MarkCompleted)

	// Quiz routes
	mux.HandleFunc("GET /api/v1/quizzes", r.quizHandler.ListQuizzes)
	mux.HandleFunc("GET /api/v1/quizzes/history", r.quizHandler.GetQuizHistory)
	mux.HandleFunc("GET /api/v1/quizzes/{id}", r.quizHandler.GetQuiz)
	mux.HandleFunc("POST /api/v1/quizzes/{id}/start", r.quizHandler.StartQuiz)
	mux.HandleFunc("GET /api/v1/quizzes/sessions/{id}", r.quizHandler.GetQuizResult)
	mux.HandleFunc("POST /api/v1/quizzes/sessions/{id}/submit", r.quizHandler.SubmitQuiz)

	// Progress routes
	mux.HandleFunc("GET /api/v1/progress/stats", r.progressHandler.GetStats)

	r.logger.Info("Routes registered successfully")

	// Apply middleware
	handler := r.loggingMiddleware(mux)
	handler = r.recoveryMiddleware(handler)
	handler = r.corsMiddleware(handler)

	return handler
}

// healthCheckHandler returns the health status of the application
func (r *Router) healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "jlpt5-backend",
	})
}

// readinessCheckHandler checks if the application is ready to serve requests
func (r *Router) readinessCheckHandler(w http.ResponseWriter, req *http.Request) {
	// Check database connection
	if err := r.db.HealthCheck(req.Context()); err != nil {
		r.logger.Error("Readiness check failed", utils.WithContext("error", err.Error()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "not ready",
			"reason": "database connection failed",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ready",
		"service": "jlpt5-backend",
	})
}
