package services

import (
	"context"
	"time"

	"github.com/joaosantos/jlpt5/internal/domain/models"
	"github.com/joaosantos/jlpt5/internal/domain/repository"
	"github.com/joaosantos/jlpt5/internal/utils"
	pkgErrors "github.com/joaosantos/jlpt5/pkg/errors"
)

// QuizService handles quiz business logic
type QuizService struct {
	quizRepo repository.QuizRepository
	logger   *utils.Logger
}

// NewQuizService creates a new quiz service
func NewQuizService(quizRepo repository.QuizRepository, logger *utils.Logger) *QuizService {
	return &QuizService{
		quizRepo: quizRepo,
		logger:   logger,
	}
}

// QuizWithQuestions represents a quiz with its questions
type QuizWithQuestions struct {
	Quiz      models.Quiz
	Questions []models.QuizQuestion
}

// QuizSessionResult represents a completed quiz session with detailed results
type QuizSessionResult struct {
	Session   models.QuizSession
	Quiz      models.Quiz
	Questions []models.QuizQuestion
	Answers   []models.QuizAnswer
}

// GetQuizzesList retrieves quizzes with optional filtering
func (s *QuizService) GetQuizzesList(ctx context.Context, jlptLevel *int, page, pageSize int) ([]models.Quiz, int, error) {
	offset := (page - 1) * pageSize

	quizzes, err := s.quizRepo.GetAllQuizzes(ctx, jlptLevel, pageSize, offset)
	if err != nil {
		s.logger.Error("Failed to get quizzes", utils.WithContext("error", err.Error()))
		return nil, 0, pkgErrors.Internal("Failed to retrieve quizzes", err)
	}

	total, err := s.quizRepo.CountQuizzes(ctx, jlptLevel)
	if err != nil {
		s.logger.Error("Failed to count quizzes", utils.WithContext("error", err.Error()))
		return nil, 0, pkgErrors.Internal("Failed to count quizzes", err)
	}

	return quizzes, total, nil
}

// GetQuizByID retrieves a quiz by ID
func (s *QuizService) GetQuizByID(ctx context.Context, quizID int) (*models.Quiz, error) {
	return s.quizRepo.GetQuizByID(ctx, quizID)
}

// StartQuizSession starts a new quiz session
func (s *QuizService) StartQuizSession(ctx context.Context, userID, quizID int) (*QuizWithQuestions, *models.QuizSession, error) {
	// Verify quiz exists
	quiz, err := s.quizRepo.GetQuizByID(ctx, quizID)
	if err != nil {
		return nil, nil, err
	}

	// Get quiz questions
	questions, err := s.quizRepo.GetQuizQuestions(ctx, quizID)
	if err != nil {
		s.logger.Error("Failed to get quiz questions", utils.WithContext("error", err.Error()))
		return nil, nil, pkgErrors.Internal("Failed to retrieve quiz questions", err)
	}

	if len(questions) == 0 {
		return nil, nil, pkgErrors.BadRequest("Quiz has no questions")
	}

	// Create quiz session
	now := time.Now()
	session := &models.QuizSession{
		UserID:    userID,
		QuizID:    quizID,
		StartedAt: now,
	}

	if err := s.quizRepo.CreateQuizSession(ctx, session); err != nil {
		s.logger.Error("Failed to create quiz session", utils.WithContext("error", err.Error()))
		return nil, nil, pkgErrors.Internal("Failed to create quiz session", err)
	}

	quizWithQuestions := &QuizWithQuestions{
		Quiz:      *quiz,
		Questions: questions,
	}

	s.logger.Info("Quiz session started", utils.WithContext("user_id", userID, "quiz_id", quizID, "session_id", session.ID))
	return quizWithQuestions, session, nil
}

// SubmitQuizAnswers submits answers for a quiz session and calculates the score
func (s *QuizService) SubmitQuizAnswers(ctx context.Context, userID, sessionID int, answers map[int]string) (*QuizSessionResult, error) {
	// Get session
	session, err := s.quizRepo.GetQuizSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Verify session belongs to user
	if session.UserID != userID {
		return nil, pkgErrors.Forbidden("Not authorized to submit answers for this session")
	}

	// Verify session is not already completed
	if session.CompletedAt != nil {
		return nil, pkgErrors.BadRequest("Quiz session already completed")
	}

	// Get quiz and questions
	quiz, err := s.quizRepo.GetQuizByID(ctx, session.QuizID)
	if err != nil {
		return nil, err
	}

	questions, err := s.quizRepo.GetQuizQuestions(ctx, session.QuizID)
	if err != nil {
		s.logger.Error("Failed to get quiz questions", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to retrieve quiz questions", err)
	}

	// Create a map of question ID to question for easy lookup
	questionMap := make(map[int]models.QuizQuestion)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	// Score the answers
	now := time.Now()
	var totalPoints int
	var earnedPoints int
	var submittedAnswers []models.QuizAnswer

	for questionID, userAnswer := range answers {
		question, exists := questionMap[questionID]
		if !exists {
			s.logger.Warn("Answer submitted for non-existent question", utils.WithContext("question_id", questionID))
			continue
		}

		isCorrect := userAnswer == question.CorrectAnswer
		if isCorrect {
			earnedPoints += question.Points
		}
		totalPoints += question.Points

		isCorrectPtr := &isCorrect
		userAnswerPtr := &userAnswer

		answer := models.QuizAnswer{
			QuizSessionID:  sessionID,
			QuizQuestionID: questionID,
			UserAnswer:     userAnswerPtr,
			IsCorrect:      isCorrectPtr,
			AnsweredAt:     now,
		}

		if err := s.quizRepo.CreateQuizAnswer(ctx, &answer); err != nil {
			s.logger.Error("Failed to create quiz answer", utils.WithContext("error", err.Error()))
			return nil, pkgErrors.Internal("Failed to save quiz answer", err)
		}

		submittedAnswers = append(submittedAnswers, answer)
	}

	// Calculate final score and percentage
	var score int
	var percentage float64
	if totalPoints > 0 {
		score = earnedPoints
		percentage = (float64(earnedPoints) / float64(totalPoints)) * 100
	}

	// Determine if passed
	passed := percentage >= float64(quiz.PassingScore)

	// Calculate time spent
	timeSpent := int(now.Sub(session.StartedAt).Seconds())

	// Update session
	session.CompletedAt = &now
	session.Score = &score
	session.TotalPoints = &totalPoints
	percentagePtr := &percentage
	session.Percentage = percentagePtr
	session.Passed = &passed
	session.TimeSpentSeconds = &timeSpent

	if err := s.quizRepo.UpdateQuizSession(ctx, session); err != nil {
		s.logger.Error("Failed to update quiz session", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to update quiz session", err)
	}

	result := &QuizSessionResult{
		Session:   *session,
		Quiz:      *quiz,
		Questions: questions,
		Answers:   submittedAnswers,
	}

	s.logger.Info("Quiz session completed", utils.WithContext(
		"user_id", userID,
		"session_id", sessionID,
		"score", score,
		"passed", passed,
	))

	return result, nil
}

// GetQuizSessionResult retrieves the result of a completed quiz session
func (s *QuizService) GetQuizSessionResult(ctx context.Context, userID, sessionID int) (*QuizSessionResult, error) {
	// Get session
	session, err := s.quizRepo.GetQuizSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Verify session belongs to user
	if session.UserID != userID {
		return nil, pkgErrors.Forbidden("Not authorized to view this session")
	}

	// Get quiz
	quiz, err := s.quizRepo.GetQuizByID(ctx, session.QuizID)
	if err != nil {
		return nil, err
	}

	// Get questions
	questions, err := s.quizRepo.GetQuizQuestions(ctx, session.QuizID)
	if err != nil {
		s.logger.Error("Failed to get quiz questions", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to retrieve quiz questions", err)
	}

	// Get answers
	answers, err := s.quizRepo.GetSessionAnswers(ctx, sessionID)
	if err != nil {
		s.logger.Error("Failed to get session answers", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to retrieve session answers", err)
	}

	result := &QuizSessionResult{
		Session:   *session,
		Quiz:      *quiz,
		Questions: questions,
		Answers:   answers,
	}

	return result, nil
}

// GetUserQuizHistory retrieves a user's quiz history
func (s *QuizService) GetUserQuizHistory(ctx context.Context, userID, page, pageSize int) ([]models.QuizSession, error) {
	offset := (page - 1) * pageSize

	sessions, err := s.quizRepo.GetUserQuizSessions(ctx, userID, pageSize, offset)
	if err != nil {
		s.logger.Error("Failed to get user quiz sessions", utils.WithContext("error", err.Error()))
		return nil, pkgErrors.Internal("Failed to retrieve quiz history", err)
	}

	return sessions, nil
}
