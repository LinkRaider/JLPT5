package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/joaosantos/jlpt5/internal/api/dto"
	"github.com/joaosantos/jlpt5/internal/domain/models"
	"github.com/joaosantos/jlpt5/internal/domain/services"
	"github.com/joaosantos/jlpt5/internal/utils"
	pkgErrors "github.com/joaosantos/jlpt5/pkg/errors"
)

// VocabularyHandler handles vocabulary endpoints
type VocabularyHandler struct {
	vocabService *services.VocabularyService
	logger       *utils.Logger
}

// NewVocabularyHandler creates a new vocabulary handler
func NewVocabularyHandler(vocabService *services.VocabularyService, logger *utils.Logger) *VocabularyHandler {
	return &VocabularyHandler{
		vocabService: vocabService,
		logger:       logger,
	}
}

// ListVocabulary retrieves a list of vocabulary items
func (h *VocabularyHandler) ListVocabulary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, pkgErrors.BadRequest("Method not allowed"))
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r)

	// Parse query parameters
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var jlptLevel *int
	if levelStr := query.Get("jlpt_level"); levelStr != "" {
		level, err := strconv.Atoi(levelStr)
		if err == nil && level >= 1 && level <= 5 {
			jlptLevel = &level
		}
	}

	items, total, err := h.vocabService.GetVocabularyList(r.Context(), userID, jlptLevel, page, pageSize)
	if err != nil {
		sendError(w, err)
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	response := dto.VocabularyListResponse{
		Items:      toVocabularyResponseList(items),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	sendSuccess(w, http.StatusOK, response)
}

// GetVocabulary retrieves a specific vocabulary item
func (h *VocabularyHandler) GetVocabulary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, pkgErrors.BadRequest("Method not allowed"))
		return
	}

	userID := getUserIDFromContext(r)
	vocabID, err := extractIDFromPath(r.URL.Path, "/api/v1/vocabulary/")
	if err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid vocabulary ID"))
		return
	}

	item, err := h.vocabService.GetVocabularyByID(r.Context(), userID, vocabID)
	if err != nil {
		sendError(w, err)
		return
	}

	sendSuccess(w, http.StatusOK, toVocabularyResponse(*item))
}

// GetDueVocabulary retrieves vocabulary items due for review
func (h *VocabularyHandler) GetDueVocabulary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, pkgErrors.BadRequest("Method not allowed"))
		return
	}

	userID := getUserIDFromContext(r)

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	items, err := h.vocabService.GetDueVocabulary(r.Context(), userID, limit)
	if err != nil {
		sendError(w, err)
		return
	}

	sendSuccess(w, http.StatusOK, map[string]interface{}{
		"items": toVocabularyResponseList(items),
		"count": len(items),
	})
}

// SubmitReview handles vocabulary review submission
func (h *VocabularyHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, pkgErrors.BadRequest("Method not allowed"))
		return
	}

	userID := getUserIDFromContext(r)
	vocabID, err := extractIDFromPath(r.URL.Path, "/api/v1/vocabulary/")
	if err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid vocabulary ID"))
		return
	}

	var req dto.ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid request body"))
		return
	}

	progress, err := h.vocabService.SubmitReview(r.Context(), userID, vocabID, req.IsCorrect)
	if err != nil {
		sendError(w, err)
		return
	}

	message := "Great job! Keep it up!"
	if !req.IsCorrect {
		message = "Don't worry, you'll get it next time!"
	}

	response := dto.ReviewResponse{
		Success:        true,
		Progress:       toProgressResponse(progress),
		NextReviewDate: progress.NextReviewDate.Format(time.RFC3339),
		Message:        message,
	}

	sendSuccess(w, http.StatusOK, response)
}

// Helper functions

func toVocabularyResponse(item models.VocabularyWithProgress) dto.VocabularyResponse {
	response := dto.VocabularyResponse{
		ID:                 item.ID,
		Word:               item.Word,
		Reading:            item.Reading,
		Meaning:            item.Meaning,
		PartOfSpeech:       item.PartOfSpeech,
		JLPTLevel:          item.JLPTLevel,
		ExampleSentence:    item.ExampleSentence,
		ExampleTranslation: item.ExampleTranslation,
		AudioURL:           item.AudioURL,
	}

	if item.Progress != nil {
		response.Progress = toProgressResponse(item.Progress)
	}

	return response
}

func toProgressResponse(progress *models.UserVocabularyProgress) *dto.ProgressResponse {
	successRate := 0.0
	if progress.TotalReviews > 0 {
		successRate = float64(progress.CorrectReviews) / float64(progress.TotalReviews) * 100
	}

	isDue := time.Now().After(progress.NextReviewDate)

	response := &dto.ProgressResponse{
		ID:             progress.ID,
		EaseFactor:     progress.EaseFactor,
		Interval:       progress.Interval,
		Repetitions:    progress.Repetitions,
		NextReviewDate: progress.NextReviewDate.Format(time.RFC3339),
		TotalReviews:   progress.TotalReviews,
		CorrectReviews: progress.CorrectReviews,
		SuccessRate:    successRate,
		IsDue:          isDue,
	}

	if progress.LastReviewedAt != nil {
		formatted := progress.LastReviewedAt.Format(time.RFC3339)
		response.LastReviewedAt = &formatted
	}

	return response
}

func toVocabularyResponseList(items []models.VocabularyWithProgress) []dto.VocabularyResponse {
	responses := make([]dto.VocabularyResponse, len(items))
	for i, item := range items {
		responses[i] = toVocabularyResponse(item)
	}
	return responses
}

func getUserIDFromContext(r *http.Request) int {
	// This will be set by auth middleware
	// For now, return a default user ID (we'll implement auth middleware later)
	return 1 // TODO: Get from context after auth middleware is implemented
}

func extractIDFromPath(path, prefix string) (int, error) {
	// Remove prefix and any trailing slashes
	idStr := strings.TrimPrefix(path, prefix)
	idStr = strings.TrimSuffix(idStr, "/review")
	idStr = strings.Trim(idStr, "/")

	return strconv.Atoi(idStr)
}

// Helper functions for response handling
func sendSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if appErr, ok := err.(*pkgErrors.AppError); ok {
		w.WriteHeader(appErr.StatusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": appErr.Message,
			"code":  appErr.Code,
		})
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": "Internal server error",
		"code":  "INTERNAL_ERROR",
	})
}
