import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-grammar-list',
  imports: [CommonModule],
  template: `
    <div class="grammar-container">
      <h1>Grammar Points</h1>
      <p>Study JLPT5 grammar patterns and structures.</p>

      <div class="placeholder">
        <p>Grammar points list coming soon...</p>
      </div>
    </div>
  `,
  styles: [`
    .grammar-container {
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
export class GrammarListComponent {
  constructor(private router: Router) {}
}
