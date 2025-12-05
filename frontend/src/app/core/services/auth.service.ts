import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { ApiService } from './api.service';
import { StorageService } from './storage.service';
import { User, AuthResponse, LoginRequest, RegisterRequest } from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  private isAuthenticatedSubject = new BehaviorSubject<boolean>(false);
  public isAuthenticated$ = this.isAuthenticatedSubject.asObservable();

  constructor(
    private apiService: ApiService,
    private storageService: StorageService
  ) {
    this.loadUserFromStorage();
  }

  private loadUserFromStorage(): void {
    const token = this.storageService.getAccessToken();
    const user = this.storageService.getUser();

    if (token && user) {
      this.currentUserSubject.next(user);
      this.isAuthenticatedSubject.next(true);
    }
  }

  register(request: RegisterRequest): Observable<AuthResponse> {
    return this.apiService.post<AuthResponse>('/auth/register', request, false).pipe(
      tap(response => this.handleAuthSuccess(response))
    );
  }

  login(request: LoginRequest): Observable<AuthResponse> {
    return this.apiService.post<AuthResponse>('/auth/login', request, false).pipe(
      tap(response => this.handleAuthSuccess(response))
    );
  }

  refreshToken(): Observable<AuthResponse> {
    const refreshToken = this.storageService.getRefreshToken();
    return this.apiService.post<AuthResponse>('/auth/refresh', { refresh_token: refreshToken }, false).pipe(
      tap(response => this.handleAuthSuccess(response))
    );
  }

  logout(): void {
    this.storageService.clearAll();
    this.currentUserSubject.next(null);
    this.isAuthenticatedSubject.next(false);
  }

  private handleAuthSuccess(response: AuthResponse): void {
    this.storageService.setAccessToken(response.access_token);
    this.storageService.setRefreshToken(response.refresh_token);
    this.storageService.setUser(response.user);
    this.currentUserSubject.next(response.user);
    this.isAuthenticatedSubject.next(true);
  }

  getCurrentUser(): User | null {
    return this.currentUserSubject.value;
  }

  isAuthenticated(): boolean {
    return this.isAuthenticatedSubject.value;
  }
}
