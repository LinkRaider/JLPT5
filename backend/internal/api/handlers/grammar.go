package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/joaosantos/jlpt5/internal/api/dto"
	"github.com/joaosantos/jlpt5/internal/domain/models"
	"github.com/joaosantos/jlpt5/internal/domain/services"
	"github.com/joaosantos/jlpt5/internal/utils"
	pkgErrors "github.com/joaosantos/jlpt5/pkg/errors"
)

// GrammarHandler handles grammar endpoints
type GrammarHandler struct {
	grammarService *services.GrammarService
	logger         *utils.Logger
}

// NewGrammarHandler creates a new grammar handler
func NewGrammarHandler(grammarService *services.GrammarService, logger *utils.Logger) *GrammarHandler {
	return &GrammarHandler{
		grammarService: grammarService,
		logger:         logger,
	}
}

// ListGrammar retrieves a list of grammar lessons
func (h *GrammarHandler) ListGrammar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, pkgErrors.BadRequest("Method not allowed"))
		return
	}

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

	lessons, total, err := h.grammarService.GetLessonsList(r.Context(), userID, jlptLevel, page, pageSize)
	if err != nil {
		sendError(w, err)
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	response := dto.GrammarListResponse{
		Items:      toGrammarLessonResponseList(lessons),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	sendSuccess(w, http.StatusOK, response)
}

// GetGrammarLesson retrieves a specific grammar lesson
func (h *GrammarHandler) GetGrammarLesson(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	lessonID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid lesson ID"))
		return
	}

	lesson, err := h.grammarService.GetLessonByID(r.Context(), userID, lessonID)
	if err != nil {
		sendError(w, err)
		return
	}

	sendSuccess(w, http.StatusOK, toGrammarLessonResponse(*lesson))
}

// MarkCompleted marks a grammar lesson as completed
func (h *GrammarHandler) MarkCompleted(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	lessonID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid lesson ID"))
		return
	}

	var req dto.MarkCompletedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid request body"))
		return
	}

	progress, err := h.grammarService.MarkAsCompleted(r.Context(), userID, lessonID, req.Notes)
	if err != nil {
		sendError(w, err)
		return
	}

	sendSuccess(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"message":  "Lesson marked as completed",
		"progress": toGrammarProgressResponse(progress),
	})
}

// Helper functions

func toGrammarLessonResponse(lesson models.GrammarLessonWithExamples) dto.GrammarLessonResponse {
	response := dto.GrammarLessonResponse{
		ID:           lesson.ID,
		Title:        lesson.Title,
		GrammarPoint: lesson.GrammarPoint,
		Explanation:  lesson.Explanation,
		UsageNotes:   lesson.UsageNotes,
		JLPTLevel:    lesson.JLPTLevel,
		LessonOrder:  lesson.LessonOrder,
		Examples:     toGrammarExampleResponseList(lesson.Examples),
	}

	if lesson.Progress != nil {
		response.Progress = toGrammarProgressResponse(lesson.Progress)
	}

	return response
}

func toGrammarExampleResponseList(examples []models.GrammarExample) []dto.GrammarExampleResponse {
	responses := make([]dto.GrammarExampleResponse, len(examples))
	for i, ex := range examples {
		responses[i] = dto.GrammarExampleResponse{
			ID:                 ex.ID,
			JapaneseSentence:   ex.JapaneseSentence,
			EnglishTranslation: ex.EnglishTranslation,
			Notes:              ex.Notes,
		}
	}
	return responses
}

func toGrammarProgressResponse(progress *models.UserGrammarProgress) *dto.GrammarProgressResponse {
	response := &dto.GrammarProgressResponse{
		ID:        progress.ID,
		Completed: progress.Completed,
		Notes:     progress.Notes,
	}

	if progress.CompletedAt != nil {
		formatted := progress.CompletedAt.Format(time.RFC3339)
		response.CompletedAt = &formatted
	}

	return response
}

func toGrammarLessonResponseList(lessons []models.GrammarLessonWithExamples) []dto.GrammarLessonResponse {
	responses := make([]dto.GrammarLessonResponse, len(lessons))
	for i, lesson := range lessons {
		responses[i] = toGrammarLessonResponse(lesson)
	}
	return responses
}
