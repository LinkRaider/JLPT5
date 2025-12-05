import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from './api.service';
import { ProgressStats } from '../models/progress.model';

@Injectable({
  providedIn: 'root'
})
export class ProgressService {
  constructor(private apiService: ApiService) {}

  getStats(): Observable<ProgressStats> {
    return this.apiService.get<ProgressStats>('/progress/stats');
  }
}
