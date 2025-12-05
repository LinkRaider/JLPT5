import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from './api.service';
import {
  GrammarLesson,
  GrammarListResponse,
  MarkCompletedRequest
} from '../models/grammar.model';
import { HttpParams } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class GrammarService {
  constructor(private apiService: ApiService) {}

  getGrammarList(page: number = 1, pageSize: number = 20, jlptLevel?: number): Observable<GrammarListResponse> {
    let params = new HttpParams()
      .set('page', page.toString())
      .set('page_size', pageSize.toString());

    if (jlptLevel) {
      params = params.set('jlpt_level', jlptLevel.toString());
    }

    return this.apiService.get<GrammarListResponse>('/grammar', params);
  }

  getGrammarLesson(id: number): Observable<GrammarLesson> {
    return this.apiService.get<GrammarLesson>(`/grammar/${id}`);
  }

  markCompleted(lessonId: number, notes?: string): Observable<any> {
    const request: MarkCompletedRequest = { notes };
    return this.apiService.post(`/grammar/${lessonId}/complete`, request);
  }
}
