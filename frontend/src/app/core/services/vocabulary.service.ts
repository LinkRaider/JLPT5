import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from './api.service';
import {
  Vocabulary,
  VocabularyListResponse,
  ReviewRequest,
  ReviewResponse
} from '../models/vocabulary.model';
import { HttpParams } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class VocabularyService {
  constructor(private apiService: ApiService) {}

  getVocabularyList(page: number = 1, pageSize: number = 20, jlptLevel?: number): Observable<VocabularyListResponse> {
    let params = new HttpParams()
      .set('page', page.toString())
      .set('page_size', pageSize.toString());

    if (jlptLevel) {
      params = params.set('jlpt_level', jlptLevel.toString());
    }

    return this.apiService.get<VocabularyListResponse>('/vocabulary', params);
  }

  getVocabularyItem(id: number): Observable<Vocabulary> {
    return this.apiService.get<Vocabulary>(`/vocabulary/${id}`);
  }

  getDueVocabulary(): Observable<VocabularyListResponse> {
    return this.apiService.get<VocabularyListResponse>('/vocabulary/due');
  }

  submitReview(vocabularyId: number, quality: number): Observable<ReviewResponse> {
    const request: ReviewRequest = { quality };
    return this.apiService.post<ReviewResponse>(`/vocabulary/${vocabularyId}/review`, request);
  }
}
