import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule],
  template: `
    <div class="dashboard-container">
      <h1>Dashboard</h1>
      <div class="dashboard-grid">
        <div class="dashboard-card" (click)="navigateTo('/vocabulary')">
          <h2>Vocabulary</h2>
          <p>Practice JLPT5 vocabulary</p>
        </div>
        <div class="dashboard-card" (click)="navigateTo('/grammar')">
          <h2>Grammar</h2>
          <p>Study JLPT5 grammar points</p>
        </div>
        <div class="dashboard-card" (click)="navigateTo('/quiz')">
          <h2>Quizzes</h2>
          <p>Test your knowledge</p>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .dashboard-container {
      padding: 20px;
      max-width: 1200px;
      margin: 0 auto;
    }

    h1 {
      color: #333;
      margin-bottom: 30px;
    }

    .dashboard-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
      gap: 20px;
    }

    .dashboard-card {
      background: white;
      border-radius: 10px;
      padding: 30px;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
      cursor: pointer;
      transition: transform 0.3s, box-shadow 0.3s;
    }

    .dashboard-card:hover {
      transform: translateY(-5px);
      box-shadow: 0 5px 20px rgba(0, 0, 0, 0.15);
    }

    .dashboard-card h2 {
      color: #667eea;
      margin: 0 0 10px 0;
    }

    .dashboard-card p {
      color: #666;
      margin: 0;
    }
  `]
})
export class DashboardComponent {
  constructor(private router: Router) {}

  navigateTo(path: string): void {
    this.router.navigate([path]);
  }
}
