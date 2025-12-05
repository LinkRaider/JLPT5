import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { StorageService } from './storage.service';

export interface ApiError {
  code: string;
  error: string;
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private readonly baseUrl = 'http://localhost:8080/api/v1';

  constructor(
    private http: HttpClient,
    private storageService: StorageService
  ) {}

  private getHeaders(includeAuth: boolean = true): HttpHeaders {
    let headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });

    if (includeAuth) {
      const token = this.storageService.getAccessToken();
      if (token) {
        headers = headers.set('Authorization', `Bearer ${token}`);
      }
    }

    return headers;
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    let errorMessage = 'An error occurred';

    if (error.error instanceof ErrorEvent) {
      // Client-side error
      errorMessage = `Error: ${error.error.message}`;
    } else {
      // Server-side error
      if (error.error && typeof error.error === 'object' && 'error' in error.error) {
        const apiError = error.error as ApiError;
        errorMessage = apiError.error;
      } else {
        errorMessage = `Error Code: ${error.status}\nMessage: ${error.message}`;
      }
    }

    console.error('API Error:', errorMessage);
    return throwError(() => new Error(errorMessage));
  }

  get<T>(endpoint: string, params?: HttpParams, includeAuth: boolean = true): Observable<T> {
    return this.http.get<T>(`${this.baseUrl}${endpoint}`, {
      headers: this.getHeaders(includeAuth),
      params
    }).pipe(catchError(this.handleError));
  }

  post<T>(endpoint: string, body: any, includeAuth: boolean = true): Observable<T> {
    return this.http.post<T>(`${this.baseUrl}${endpoint}`, body, {
      headers: this.getHeaders(includeAuth)
    }).pipe(catchError(this.handleError));
  }

  put<T>(endpoint: string, body: any, includeAuth: boolean = true): Observable<T> {
    return this.http.put<T>(`${this.baseUrl}${endpoint}`, body, {
      headers: this.getHeaders(includeAuth)
    }).pipe(catchError(this.handleError));
  }

  delete<T>(endpoint: string, includeAuth: boolean = true): Observable<T> {
    return this.http.delete<T>(`${this.baseUrl}${endpoint}`, {
      headers: this.getHeaders(includeAuth)
    }).pipe(catchError(this.handleError));
  }
}
