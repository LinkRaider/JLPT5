export interface Quiz {
  id: number;
  title: string;
  description?: string;
  quiz_type: string;
  jlpt_level: number;
  difficulty_level: string;
  time_limit_minutes?: number;
  passing_score: number;
  question_count: number;
}

export interface QuizQuestion {
  id: number;
  question_text: string;
  question_type: string;
  options: string[];
  points: number;
}

export interface QuizSession {
  session_id: number;
  quiz_id: number;
  quiz_title: string;
  questions: QuizQuestion[];
  started_at: string;
  time_limit_minutes?: number;
}

export interface QuizSubmission {
  answers: { [questionId: number]: string };
}

export interface QuizResult {
  session_id: number;
  quiz_id: number;
  quiz_title: string;
  total_points: number;
  earned_points: number;
  percentage: number;
  passed: boolean;
  started_at: string;
  completed_at?: string;
  answers: QuizAnswer[];
}

export interface QuizAnswer {
  question_id: number;
  question_text: string;
  user_answer: string;
  correct_answer: string;
  is_correct: boolean;
  points: number;
}

export interface QuizListResponse {
  items: Quiz[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface QuizHistory {
  sessions: QuizHistorySession[];
  total_quizzes_taken: number;
  average_score: number;
  pass_rate: number;
}

export interface QuizHistorySession {
  session_id: number;
  quiz_id: number;
  quiz_title: string;
  percentage: number;
  passed: boolean;
  completed_at: string;
}
