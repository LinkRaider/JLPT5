import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';
import { LoginComponent } from './features/auth/components/login/login.component';
import { RegisterComponent } from './features/auth/components/register/register.component';
import { DashboardComponent } from './features/progress/components/dashboard/dashboard.component';
import { VocabularyListComponent } from './features/vocabulary/components/vocabulary-list/vocabulary-list.component';
import { VocabularyReviewComponent } from './features/vocabulary/components/vocabulary-review/vocabulary-review.component';
import { GrammarListComponent } from './features/grammar/components/grammar-list/grammar-list.component';
import { GrammarDetailComponent } from './features/grammar/components/grammar-detail/grammar-detail.component';
import { QuizListComponent } from './features/quiz/components/quiz-list/quiz-list.component';
import { QuizSessionComponent } from './features/quiz/components/quiz-session/quiz-session.component';
import { QuizResultComponent } from './features/quiz/components/quiz-result/quiz-result.component';
import { QuizHistoryComponent } from './features/quiz/components/quiz-history/quiz-history.component';

export const routes: Routes = [
  // Redirect root to dashboard
  {
    path: '',
    redirectTo: '/dashboard',
    pathMatch: 'full'
  },

  // Auth routes (no guard)
  {
    path: 'login',
    component: LoginComponent
  },
  {
    path: 'register',
    component: RegisterComponent
  },

  // Protected routes (with auth guard)
  {
    path: 'dashboard',
    component: DashboardComponent,
    canActivate: [authGuard]
  },

  // Vocabulary routes
  {
    path: 'vocabulary',
    canActivate: [authGuard],
    children: [
      {
        path: '',
        component: VocabularyListComponent
      },
      {
        path: 'review',
        component: VocabularyReviewComponent
      },
      {
        path: ':id',
        component: VocabularyListComponent
      }
    ]
  },

  // Grammar routes
  {
    path: 'grammar',
    canActivate: [authGuard],
    children: [
      {
        path: '',
        component: GrammarListComponent
      },
      {
        path: ':id',
        component: GrammarDetailComponent
      }
    ]
  },

  // Quiz routes
  {
    path: 'quiz',
    canActivate: [authGuard],
    children: [
      {
        path: '',
        component: QuizListComponent
      },
      {
        path: 'history',
        component: QuizHistoryComponent
      },
      {
        path: 'result/:id',
        component: QuizResultComponent
      },
      {
        path: ':id',
        component: QuizSessionComponent
      }
    ]
  },

  // Wildcard route - redirect to dashboard
  {
    path: '**',
    redirectTo: '/dashboard'
  }
];
