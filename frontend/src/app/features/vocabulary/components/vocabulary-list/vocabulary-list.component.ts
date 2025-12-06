import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-vocabulary-list',
  imports: [CommonModule],
  template: `
    <div class="vocabulary-container">
      <h1>Vocabulary List</h1>
      <p>Browse and study JLPT5 vocabulary words.</p>

      <div class="actions">
        <button class="btn-primary" (click)="startReview()">Start Review</button>
      </div>

      <div class="placeholder">
        <p>Vocabulary list coming soon...</p>
      </div>
    </div>
  `,
  styles: [`
    .vocabulary-container {
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
export class VocabularyListComponent {
  constructor(private router: Router) {}

  startReview(): void {
    this.router.navigate(['/vocabulary/review']);
  }
}
