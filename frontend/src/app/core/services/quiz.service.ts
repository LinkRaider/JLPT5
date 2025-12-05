import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from './api.service';
import {
  Quiz,
  QuizListResponse,
  QuizSession,
  QuizSubmission,
  QuizResult,
  QuizHistory
} from '../models/quiz.model';
import { HttpParams } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class QuizService {
  constructor(private apiService: ApiService) {}

  getQuizList(page: number = 1, pageSize: number = 20, jlptLevel?: number): Observable<QuizListResponse> {
    let params = new HttpParams()
      .set('page', page.toString())
      .set('page_size', pageSize.toString());

    if (jlptLevel) {
      params = params.set('jlpt_level', jlptLevel.toString());
    }

    return this.apiService.get<QuizListResponse>('/quizzes', params);
  }

  getQuiz(id: number): Observable<Quiz> {
    return this.apiService.get<Quiz>(`/quizzes/${id}`);
  }

  startQuiz(quizId: number): Observable<QuizSession> {
    return this.apiService.post<QuizSession>(`/quizzes/${quizId}/start`, {});
  }

  submitQuiz(sessionId: number, answers: { [questionId: number]: string }): Observable<QuizResult> {
    const submission: QuizSubmission = { answers };
    return this.apiService.post<QuizResult>(`/quizzes/sessions/${sessionId}/submit`, submission);
  }

  getQuizResult(sessionId: number): Observable<QuizResult> {
    return this.apiService.get<QuizResult>(`/quizzes/sessions/${sessionId}`);
  }

  getQuizHistory(): Observable<QuizHistory> {
    return this.apiService.get<QuizHistory>('/quizzes/history');
  }
}
