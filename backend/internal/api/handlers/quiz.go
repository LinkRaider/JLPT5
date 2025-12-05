package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joaosantos/jlpt5/internal/api/dto"
	"github.com/joaosantos/jlpt5/internal/domain/models"
	"github.com/joaosantos/jlpt5/internal/domain/services"
	"github.com/joaosantos/jlpt5/internal/utils"
	pkgErrors "github.com/joaosantos/jlpt5/pkg/errors"
)

// QuizHandler handles quiz endpoints
type QuizHandler struct {
	quizService *services.QuizService
	logger      *utils.Logger
}

// NewQuizHandler creates a new quiz handler
func NewQuizHandler(quizService *services.QuizService, logger *utils.Logger) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
		logger:      logger,
	}
}

// ListQuizzes retrieves a list of quizzes
func (h *QuizHandler) ListQuizzes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, pkgErrors.BadRequest("Method not allowed"))
		return
	}

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

	quizzes, total, err := h.quizService.GetQuizzesList(r.Context(), jlptLevel, page, pageSize)
	if err != nil {
		sendError(w, err)
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	response := dto.QuizListResponse{
		Items:      toQuizResponseList(quizzes),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	sendSuccess(w, http.StatusOK, response)
}

// GetQuiz retrieves a specific quiz
func (h *QuizHandler) GetQuiz(w http.ResponseWriter, r *http.Request) {
	quizID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid quiz ID"))
		return
	}

	quiz, err := h.quizService.GetQuizByID(r.Context(), quizID)
	if err != nil {
		sendError(w, err)
		return
	}

	sendSuccess(w, http.StatusOK, toQuizResponse(*quiz))
}

// StartQuiz starts a new quiz session
func (h *QuizHandler) StartQuiz(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	quizID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid quiz ID"))
		return
	}

	quizWithQuestions, session, err := h.quizService.StartQuizSession(r.Context(), userID, quizID)
	if err != nil {
		sendError(w, err)
		return
	}

	response := dto.StartQuizResponse{
		SessionID: session.ID,
		Quiz:      toQuizResponse(quizWithQuestions.Quiz),
		Questions: toQuizQuestionResponseList(quizWithQuestions.Questions),
		StartedAt: session.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	sendSuccess(w, http.StatusOK, response)
}

// SubmitQuiz submits answers for a quiz session
func (h *QuizHandler) SubmitQuiz(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	sessionID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid session ID"))
		return
	}

	var req dto.SubmitQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid request body"))
		return
	}

	if len(req.Answers) == 0 {
		sendError(w, pkgErrors.BadRequest("No answers provided"))
		return
	}

	result, err := h.quizService.SubmitQuizAnswers(r.Context(), userID, sessionID, req.Answers)
	if err != nil {
		sendError(w, err)
		return
	}

	response := toQuizResultResponse(result)
	sendSuccess(w, http.StatusOK, response)
}

// GetQuizResult retrieves the result of a quiz session
func (h *QuizHandler) GetQuizResult(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	sessionID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, pkgErrors.BadRequest("Invalid session ID"))
		return
	}

	result, err := h.quizService.GetQuizSessionResult(r.Context(), userID, sessionID)
	if err != nil {
		sendError(w, err)
		return
	}

	response := toQuizResultResponse(result)
	sendSuccess(w, http.StatusOK, response)
}

// GetQuizHistory retrieves user's quiz history
func (h *QuizHandler) GetQuizHistory(w http.ResponseWriter, r *http.Request) {
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

	sessions, err := h.quizService.GetUserQuizHistory(r.Context(), userID, page, pageSize)
	if err != nil {
		sendError(w, err)
		return
	}

	response := dto.QuizHistoryResponse{
		Sessions: toQuizSessionResponseList(sessions),
		Page:     page,
		PageSize: pageSize,
	}

	sendSuccess(w, http.StatusOK, response)
}

// Helper functions

func toQuizResponse(quiz models.Quiz) dto.QuizResponse {
	return dto.QuizResponse{
		ID:           quiz.ID,
		Title:        quiz.Title,
		Description:  quiz.Description,
		JLPTLevel:    quiz.JLPTLevel,
		QuizType:     quiz.QuizType,
		PassingScore: quiz.PassingScore,
	}
}

func toQuizResponseList(quizzes []models.Quiz) []dto.QuizResponse {
	responses := make([]dto.QuizResponse, len(quizzes))
	for i, quiz := range quizzes {
		responses[i] = toQuizResponse(quiz)
	}
	return responses
}

func toQuizQuestionResponse(question models.QuizQuestion) dto.QuizQuestionResponse {
	// Convert individual option fields to slice
	var options []string
	if question.OptionA != nil {
		options = append(options, *question.OptionA)
	}
	if question.OptionB != nil {
		options = append(options, *question.OptionB)
	}
	if question.OptionC != nil {
		options = append(options, *question.OptionC)
	}
	if question.OptionD != nil {
		options = append(options, *question.OptionD)
	}

	return dto.QuizQuestionResponse{
		ID:           question.ID,
		QuestionText: question.QuestionText,
		QuestionType: string(question.QuestionType),
		Options:      options,
		Points:       question.Points,
	}
}

func toQuizQuestionResponseList(questions []models.QuizQuestion) []dto.QuizQuestionResponse {
	responses := make([]dto.QuizQuestionResponse, len(questions))
	for i, q := range questions {
		responses[i] = toQuizQuestionResponse(q)
	}
	return responses
}

func toQuizQuestionDetailResponse(question models.QuizQuestion) dto.QuizQuestionDetailResponse {
	// Convert individual option fields to slice
	var options []string
	if question.OptionA != nil {
		options = append(options, *question.OptionA)
	}
	if question.OptionB != nil {
		options = append(options, *question.OptionB)
	}
	if question.OptionC != nil {
		options = append(options, *question.OptionC)
	}
	if question.OptionD != nil {
		options = append(options, *question.OptionD)
	}

	return dto.QuizQuestionDetailResponse{
		ID:            question.ID,
		QuestionText:  question.QuestionText,
		QuestionType:  string(question.QuestionType),
		Options:       options,
		CorrectAnswer: question.CorrectAnswer,
		Explanation:   question.Explanation,
		Points:        question.Points,
	}
}

func toQuizQuestionDetailResponseList(questions []models.QuizQuestion) []dto.QuizQuestionDetailResponse {
	responses := make([]dto.QuizQuestionDetailResponse, len(questions))
	for i, q := range questions {
		responses[i] = toQuizQuestionDetailResponse(q)
	}
	return responses
}

func toQuizAnswerResponse(answer models.QuizAnswer) dto.QuizAnswerResponse {
	var userAnswer string
	var isCorrect bool

	if answer.UserAnswer != nil {
		userAnswer = *answer.UserAnswer
	}
	if answer.IsCorrect != nil {
		isCorrect = *answer.IsCorrect
	}

	return dto.QuizAnswerResponse{
		QuestionID: answer.QuizQuestionID,
		UserAnswer: userAnswer,
		IsCorrect:  isCorrect,
	}
}

func toQuizAnswerResponseList(answers []models.QuizAnswer) []dto.QuizAnswerResponse {
	responses := make([]dto.QuizAnswerResponse, len(answers))
	for i, a := range answers {
		responses[i] = toQuizAnswerResponse(a)
	}
	return responses
}

func toQuizResultResponse(result *services.QuizSessionResult) dto.QuizResultResponse {
	var score, totalPoints int
	var percentage float64
	var passed bool
	var completedAt string

	if result.Session.Score != nil {
		score = *result.Session.Score
	}
	if result.Session.TotalPoints != nil {
		totalPoints = *result.Session.TotalPoints
	}
	if result.Session.Percentage != nil {
		percentage = *result.Session.Percentage
	}
	if result.Session.Passed != nil {
		passed = *result.Session.Passed
	}
	if result.Session.CompletedAt != nil {
		completedAt = result.Session.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return dto.QuizResultResponse{
		SessionID:      result.Session.ID,
		Quiz:           toQuizResponse(result.Quiz),
		Score:          score,
		TotalPoints:    totalPoints,
		Percentage:     percentage,
		TotalQuestions: len(result.Questions),
		Passed:         passed,
		StartedAt:      result.Session.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
		CompletedAt:    completedAt,
		Questions:      toQuizQuestionDetailResponseList(result.Questions),
		Answers:        toQuizAnswerResponseList(result.Answers),
	}
}

func toQuizSessionResponse(session models.QuizSession) dto.QuizSessionResponse {
	response := dto.QuizSessionResponse{
		ID:         session.ID,
		QuizID:     session.QuizID,
		StartedAt:  session.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
		Score:      session.Score,
		Percentage: session.Percentage,
		Passed:     session.Passed,
	}

	if session.CompletedAt != nil {
		completedStr := session.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		response.CompletedAt = &completedStr
	}

	return response
}

func toQuizSessionResponseList(sessions []models.QuizSession) []dto.QuizSessionResponse {
	responses := make([]dto.QuizSessionResponse, len(sessions))
	for i, s := range sessions {
		responses[i] = toQuizSessionResponse(s)
	}
	return responses
}
