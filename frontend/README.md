# JLPT5 Frontend

Angular 18+ frontend application for the JLPT5 Japanese learning platform.

## Overview

This is a modern Angular application built with standalone components, featuring vocabulary flashcards, grammar lessons, practice quizzes, and progress tracking for JLPT N5 level Japanese learners.

## Technology Stack

- **Framework**: Angular 18+
- **Language**: TypeScript 5+
- **Architecture**: Standalone Components
- **State Management**: RxJS
- **HTTP Client**: Angular HttpClient with JWT interceptor
- **Routing**: Angular Router with lazy loading
- **Testing**: Vitest + Angular Testing Library
- **Build Tool**: Angular CLI with esbuild

## Project Structure

```
src/
├── app/
│   ├── core/                    # Core services and guards
│   │   ├── guards/             # Route guards (auth)
│   │   ├── interceptors/       # HTTP interceptors (JWT)
│   │   └── services/           # Core services (auth, API)
│   ├── features/               # Feature modules
│   │   ├── auth/              # Authentication
│   │   │   ├── components/
│   │   │   │   ├── login/
│   │   │   │   └── register/
│   │   │   └── services/
│   │   ├── vocabulary/        # Vocabulary flashcards
│   │   │   ├── components/
│   │   │   │   ├── vocabulary-list/
│   │   │   │   └── vocabulary-review/
│   │   │   └── services/
│   │   ├── grammar/           # Grammar lessons
│   │   │   ├── components/
│   │   │   │   ├── grammar-list/
│   │   │   │   └── grammar-detail/
│   │   │   └── services/
│   │   ├── quiz/              # Practice quizzes
│   │   │   ├── components/
│   │   │   │   ├── quiz-list/
│   │   │   │   ├── quiz-session/
│   │   │   │   ├── quiz-result/
│   │   │   │   └── quiz-history/
│   │   │   └── services/
│   │   └── progress/          # Progress tracking
│   │       ├── components/
│   │       │   └── dashboard/
│   │       └── services/
│   ├── shared/                # Shared components and utilities
│   ├── app.component.ts       # Root component
│   └── app.routes.ts          # Application routes
├── assets/                    # Static assets
├── environments/              # Environment configurations
└── styles/                    # Global styles
```

## Development

### Prerequisites

- Node.js 20+
- npm 10+

### Setup

```bash
# Install dependencies
npm install

# Start development server
npm start
# or
ng serve

# Application runs at http://localhost:4200
```

### Development Commands

```bash
# Start dev server
npm start

# Build for production
npm run build

# Run unit tests
npm test

# Run tests with coverage
npm test -- --code-coverage

# Run linting
npm run lint

# Fix linting issues
npm run lint -- --fix

# Generate a new component
ng generate component features/quiz/components/my-component

# Generate a new service
ng generate service features/quiz/services/my-service
```

## Features

### Authentication

- Login/Register forms with validation
- JWT token management
- Automatic token refresh
- Route guards for protected pages
- HTTP interceptor for adding auth headers

### Vocabulary

- Flashcard-based learning interface
- Spaced repetition review system (SM-2 algorithm)
- Progress tracking per word
- Due date calculation
- Review statistics

### Grammar

- Lesson browser with categorization
- Detailed grammar explanations
- Example sentences with translations
- Practice exercises

### Quizzes

- Multiple choice questions
- Fill-in-the-blank exercises
- Timed quiz sessions
- Immediate feedback
- Score tracking and history

### Progress Dashboard

- Overall statistics
- Study streak tracking
- Vocabulary mastery levels
- Quiz performance graphs
- Daily/weekly/monthly views

## Routing

```typescript
// Application routes
const routes = [
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'dashboard', component: DashboardComponent, canActivate: [AuthGuard] },
  { path: 'vocabulary', component: VocabularyListComponent, canActivate: [AuthGuard] },
  { path: 'vocabulary/review', component: VocabularyReviewComponent, canActivate: [AuthGuard] },
  { path: 'grammar', component: GrammarListComponent, canActivate: [AuthGuard] },
  { path: 'grammar/:id', component: GrammarDetailComponent, canActivate: [AuthGuard] },
  { path: 'quiz', component: QuizListComponent, canActivate: [AuthGuard] },
  { path: 'quiz/:id', component: QuizSessionComponent, canActivate: [AuthGuard] },
  { path: 'quiz/:id/result', component: QuizResultComponent, canActivate: [AuthGuard] },
  { path: 'quiz/history', component: QuizHistoryComponent, canActivate: [AuthGuard] },
];
```

## API Integration

### Base Configuration

```typescript
// Environment configuration
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api/v1'
};
```

### HTTP Interceptor

All API requests automatically include:
- JWT token in Authorization header
- Content-Type: application/json
- Error handling and retry logic

### Example Service

```typescript
@Injectable({ providedIn: 'root' })
export class VocabularyService {
  constructor(private http: HttpClient) {}

  getVocabulary() {
    return this.http.get<Vocabulary[]>(`${environment.apiUrl}/vocabulary`);
  }

  getDueForReview() {
    return this.http.get<Vocabulary[]>(`${environment.apiUrl}/vocabulary/due`);
  }

  submitReview(id: string, rating: number) {
    return this.http.post(`${environment.apiUrl}/vocabulary/${id}/review`, { rating });
  }
}
```

## Component Architecture

### Standalone Components

All components use Angular's standalone component syntax:

```typescript
@Component({
  selector: 'app-vocabulary-list',
  imports: [CommonModule, RouterModule],
  template: `...`,
  styles: [`...`]
})
export class VocabularyListComponent {
  // Component logic
}
```

### Inline Templates

Components use inline templates and styles for better encapsulation and maintainability.

## Testing

### Unit Tests

```bash
# Run all tests
npm test

# Run with coverage
npm test -- --code-coverage

# Run specific test file
npm test -- vocabulary-list.component.spec.ts
```

### Test Structure

```typescript
describe('VocabularyListComponent', () => {
  let component: VocabularyListComponent;
  let fixture: ComponentFixture<VocabularyListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VocabularyListComponent]
    }).compileComponents();

    fixture = TestBed.createComponent(VocabularyListComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
```

## Building for Production

### Docker Build

```bash
# Build Docker image
docker build -t jlpt5-frontend:latest .

# Run container
docker run -p 4200:4200 jlpt5-frontend:latest
```

### Production Build

```bash
# Build with production configuration
npm run build

# Output in dist/frontend/browser/
```

### Build Optimization

The production build includes:
- Ahead-of-Time (AOT) compilation
- Tree shaking
- Minification
- Compression
- Source maps (optional)

## Environment Configuration

### Development

```typescript
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api/v1'
};
```

### Production

```typescript
export const environment = {
  production: true,
  apiUrl: 'https://api.jlpt5.com/api/v1'
};
```

## Code Style

### TypeScript

- Strict mode enabled
- ESLint for code quality
- Prettier for formatting

### Angular Best Practices

- Use standalone components
- Implement OnPush change detection where possible
- Unsubscribe from observables
- Use async pipe for subscriptions
- Lazy load feature modules

## Deployment

### Docker Compose (Development)

```bash
cd ..
docker compose up -d frontend
```

### Kubernetes with Helm

```bash
cd ../helm/jlpt5
helm install jlpt5 . --set frontend.image.tag=latest
```

### GitHub Actions CI/CD

Commits with specific keywords trigger automated builds:

```bash
# Build frontend
git commit -m "frontend: update quiz component"
git push

# Deploy frontend
git commit -m "deploy: frontend styling updates"
git push
```

See [GitHub Actions Documentation](../.github/workflows/README.md) for details.

## Troubleshooting

### Port Already in Use

```bash
# Kill process on port 4200
lsof -ti:4200 | xargs kill -9

# Or use a different port
ng serve --port 4201
```

### Module Not Found

```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Build Errors

```bash
# Clear Angular cache
rm -rf .angular

# Rebuild
npm run build
```

## Additional Resources

- [Angular Documentation](https://angular.dev)
- [Angular CLI Documentation](https://angular.dev/tools/cli)
- [RxJS Documentation](https://rxjs.dev)
- [TypeScript Documentation](https://www.typescriptlang.org/docs)
- [Main Project README](../README.md)

## Contributing

Follow semantic commit messages:
- `feat:` - New feature
- `fix:` - Bug fix
- `style:` - UI/styling changes
- `refactor:` - Code refactoring
- `test:` - Test updates
- `docs:` - Documentation updates

All PRs must pass linting and tests.
