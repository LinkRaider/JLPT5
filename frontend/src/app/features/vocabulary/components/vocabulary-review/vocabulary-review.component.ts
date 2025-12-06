import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-vocabulary-review',
  imports: [CommonModule],
  template: `
    <div class="review-container">
      <h1>Vocabulary Review</h1>
      <p>Practice your JLPT5 vocabulary with flashcards.</p>

      <div class="placeholder">
        <p>Vocabulary review feature coming soon...</p>
        <button class="btn-secondary" (click)="goBack()">Back to List</button>
      </div>
    </div>
  `,
  styles: [`
    .review-container {
      padding: 20px;
      max-width: 800px;
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
export class VocabularyReviewComponent {
  constructor(private router: Router) {}

  goBack(): void {
    this.router.navigate(['/vocabulary']);
  }
}
