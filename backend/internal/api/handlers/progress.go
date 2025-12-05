package handlers

import (
	"net/http"

	"github.com/joaosantos/jlpt5/internal/api/dto"
	"github.com/joaosantos/jlpt5/internal/domain/services"
	"github.com/joaosantos/jlpt5/internal/utils"
)

// ProgressHandler handles progress and statistics endpoints
type ProgressHandler struct {
	progressService *services.ProgressService
	logger          *utils.Logger
}

// NewProgressHandler creates a new progress handler
func NewProgressHandler(progressService *services.ProgressService, logger *utils.Logger) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
		logger:          logger,
	}
}

// GetStats retrieves overall user statistics
func (h *ProgressHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	stats, err := h.progressService.GetUserStatistics(r.Context(), userID)
	if err != nil {
		sendError(w, err)
		return
	}

	// Calculate grammar progress percentage
	var grammarProgress float64
	if stats.GrammarTotal > 0 {
		grammarProgress = (float64(stats.GrammarCompleted) / float64(stats.GrammarTotal)) * 100
	}

	response := dto.ProgressStatsResponse{
		UserID:                stats.UserID,
		StudyStreakDays:       stats.StudyStreakDays,
		LastStudyDate:         stats.LastStudyDate,
		TotalStudyTimeMinutes: stats.TotalStudyTimeMinutes,
		VocabularyLearned:     stats.VocabularyLearned,
		VocabularyDueCount:    stats.VocabularyDueCount,
		GrammarCompleted:      stats.GrammarCompleted,
		GrammarTotal:          stats.GrammarTotal,
		GrammarProgress:       grammarProgress,
		QuizzesTaken:          stats.QuizzesTaken,
		QuizzesPassed:         stats.QuizzesPassed,
		AverageQuizScore:      stats.AverageQuizScore,
	}

	sendSuccess(w, http.StatusOK, response)
}
