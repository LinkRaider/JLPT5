export interface GrammarLesson {
  id: number;
  title: string;
  grammar_point: string;
  explanation: string;
  usage_notes?: string;
  jlpt_level: number;
  lesson_order: number;
  examples: GrammarExample[];
  progress?: GrammarProgress;
}

export interface GrammarExample {
  id: number;
  japanese_sentence: string;
  english_translation: string;
  notes?: string;
}

export interface GrammarProgress {
  id: number;
  completed: boolean;
  completed_at?: string;
  notes?: string;
}

export interface GrammarListResponse {
  items: GrammarLesson[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface MarkCompletedRequest {
  notes?: string;
}
