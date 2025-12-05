export interface ProgressStats {
  user_id: number;
  vocabulary_learned: number;
  vocabulary_mastered: number;
  vocabulary_due: number;
  grammar_completed: number;
  grammar_total: number;
  quizzes_taken: number;
  quizzes_passed: number;
  average_quiz_score: number;
  study_streak_days: number;
  last_study_date?: string;
  total_study_time_minutes: number;
}
