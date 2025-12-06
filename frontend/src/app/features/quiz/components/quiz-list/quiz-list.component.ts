import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-quiz-list',
  imports: [CommonModule],
  template: `
    <div class="quiz-container">
      <h1>Quizzes</h1>
      <p>Test your JLPT5 knowledge with practice quizzes.</p>

      <div class="actions">
        <button class="btn-primary" (click)="viewHistory()">View Quiz History</button>
      </div>

      <div class="placeholder">
        <p>Quiz list coming soon...</p>
      </div>
    </div>
  `,
  styles: [`
    .quiz-container {
      padding: 20px;
      max-width: 1200px;
      margin: 0 auto;
    }

    h1 {
      color: #333;
      margin-bottom: 10px;
    }

    p {
      color: #666;
      margin-bottom: 20px;
    }

    .actions {
      margin-bottom: 20px;
    }

    .btn-primary {
      padding: 12px 24px;
      background: #667eea;
      color: white;
      border: none;
      border-radius: 5px;
      font-size: 16px;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.3s;
    }

    .btn-primary:hover {
      background: #5568d3;
    }

    .placeholder {
      background: white;
      border-radius: 10px;
      padding: 40px;
      text-align: center;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    }

    .placeholder p {
      color: #999;
      font-style: italic;
    }
  `]
})
export class QuizListComponent {
  constructor(private router: Router) {}

  viewHistory(): void {
    this.router.navigate(['/quiz/history']);
  }
}
