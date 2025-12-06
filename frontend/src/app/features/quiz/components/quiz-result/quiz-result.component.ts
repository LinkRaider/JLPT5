import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-quiz-result',
  imports: [CommonModule],
  template: `
    <div class="result-container">
      <h1>Quiz Results</h1>
      <p>View your quiz performance and score.</p>

      <div class="placeholder">
        <p>Quiz results feature coming soon...</p>
        <button class="btn-secondary" (click)="goBack()">Back to Quizzes</button>
      </div>
    </div>
  `,
  styles: [`
    .result-container {
      padding: 20px;
      max-width: 900px;
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

    .placeholder {
      background: white;
      border-radius: 10px;
      padding: 60px;
      text-align: center;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    }

    .placeholder p {
      color: #999;
      font-style: italic;
      margin-bottom: 20px;
    }

    .btn-secondary {
      padding: 12px 24px;
      background: #6c757d;
      color: white;
      border: none;
      border-radius: 5px;
      font-size: 16px;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.3s;
    }

    .btn-secondary:hover {
      background: #5a6268;
    }
  `]
})
export class QuizResultComponent {
  constructor(
    private router: Router,
    private route: ActivatedRoute
  ) {}

  goBack(): void {
    this.router.navigate(['/quiz']);
  }
}
