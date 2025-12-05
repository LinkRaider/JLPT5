export interface Vocabulary {
  id: number;
  word: string;
  reading: string;
  meaning: string;
  part_of_speech: string;
  example_sentence?: string;
  jlpt_level: number;
  progress?: VocabularyProgress;
}

export interface VocabularyProgress {
  id: number;
  ease_factor: number;
  interval: number;
  repetitions: number;
  next_review_date: string;
  last_reviewed_at?: string;
  total_reviews: number;
  correct_reviews: number;
}

export interface VocabularyListResponse {
  items: Vocabulary[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ReviewRequest {
  quality: number; // 0-5 for SM-2 algorithm
}

export interface ReviewResponse {
  success: boolean;
  message: string;
  progress: VocabularyProgress;
}
