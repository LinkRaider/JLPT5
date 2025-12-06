# JLPT5 Learning Application

A production-ready full-stack web application for learning Japanese at JLPT N5 level, featuring vocabulary flashcards with spaced repetition, grammar lessons, practice quizzes, and progress tracking.

[![Backend CI](https://github.com/yourusername/jlpt5/actions/workflows/backend-ci.yml/badge.svg)](https://github.com/yourusername/jlpt5/actions/workflows/backend-ci.yml)
[![Frontend CI](https://github.com/yourusername/jlpt5/actions/workflows/frontend-ci.yml/badge.svg)](https://github.com/yourusername/jlpt5/actions/workflows/frontend-ci.yml)

## Technology Stack

- **Frontend**: Angular 18+ (TypeScript) with standalone components
- **Backend**: Go 1.24+ with clean architecture
- **Database**: PostgreSQL 16-alpine
- **Container Orchestration**: Kubernetes with Helm 4
- **CI/CD**: GitHub Actions with semantic commit triggers
- **Containerization**: Docker with multi-stage builds

## Features

- **Vocabulary Flashcards**: Study JLPT N5 vocabulary with intelligent spaced repetition (SM-2 algorithm)
- **Grammar Lessons**: Comprehensive grammar points with examples and explanations
- **Practice Quizzes**: Multiple choice and fill-in-the-blank questions to test knowledge
- **Progress Tracking**: Track learning statistics, study streaks, and overall progress
- **User Authentication**: Secure JWT-based authentication with refresh tokens
- **RESTful API**: Clean, versioned API endpoints

## Quick Start

### Option 1: Docker Compose (Recommended for Development)

```bash
# Clone the repository
git clone https://github.com/yourusername/jlpt5.git
cd jlpt5

# Start all services
docker compose up -d

# Access the application
# Frontend: http://localhost:4200
# Backend API: http://localhost:8080
# Database: localhost:5432
```

### Option 2: Kubernetes with Helm (Production)

```bash
# Install with Helm
helm install jlpt5 ./helm/jlpt5

# Or with custom values
helm install jlpt5 ./helm/jlpt5 -f custom-values.yaml

# Access the application (depending on service type)
kubectl get svc
```

See [Helm Chart Documentation](./helm/jlpt5/README.md) for detailed deployment options.

## Architecture

### Monorepo Structure

```
JLPT5/
├── .github/
│   └── workflows/          # GitHub Actions CI/CD pipelines
├── backend/
│   ├── cmd/api/           # Application entry point
│   ├── internal/          # Private application code
│   │   ├── domain/        # Domain models
│   │   ├── handler/       # HTTP handlers
│   │   ├── service/       # Business logic
│   │   ├── repository/    # Data access
│   │   └── infrastructure/ # Database, migrations
│   ├── Dockerfile         # Backend container image
│   └── go.mod             # Go dependencies
├── frontend/
│   ├── src/app/
│   │   ├── core/          # Core services, guards
│   │   ├── features/      # Feature modules
│   │   │   ├── auth/      # Authentication
│   │   │   ├── vocabulary/ # Vocabulary flashcards
│   │   │   ├── grammar/   # Grammar lessons
│   │   │   ├── quiz/      # Practice quizzes
│   │   │   └── progress/  # Progress tracking
│   │   └── shared/        # Shared components
│   ├── Dockerfile         # Frontend container image
│   └── package.json       # Node dependencies
├── helm/jlpt5/            # Helm chart for K8s deployment
│   ├── Chart.yaml         # Chart metadata
│   ├── values.yaml        # Default configuration
│   ├── templates/         # Kubernetes manifests
│   └── README.md          # Helm documentation
└── docker-compose.yml     # Local development environment
```

### Backend Architecture

- **Clean Architecture**: Domain-driven design with dependency inversion
- **Layered Design**: Handler → Service → Repository
- **RESTful API**: Versioned endpoints (`/api/v1`)
- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **Database Migrations**: Automatic migrations using golang-migrate
- **Health Checks**: `/health` endpoint for liveness and readiness probes

### Frontend Architecture

- **Standalone Components**: Modern Angular 18+ architecture
- **Lazy Loading**: Feature modules loaded on demand
- **Reactive State**: RxJS-based state management
- **Inline Templates**: Component-scoped templates and styles
- **JWT Interceptor**: Automatic token attachment to requests

### Database Design

- **PostgreSQL 16**: Relational database with proper normalization
- **Migration-Based**: Version-controlled schema management
- **Tables**: users, vocabulary, grammar, quiz, progress, reviews
- **Indexes**: Optimized for spaced repetition queries

## Development

### Prerequisites

- **Docker Desktop** (for local development)
- **Go 1.24+** (for backend development)
- **Node.js 20+** (for frontend development)
- **kubectl** (for Kubernetes deployment)
- **Helm 4+** (for Kubernetes deployment)

### Local Development with Docker

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down

# Rebuild after code changes
docker compose up -d --build
```

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Run locally (requires PostgreSQL)
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=jlpt5_user
export DB_PASSWORD=jlpt5_password
export DB_NAME=jlpt5_dev
go run cmd/api/main.go

# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm start
# or
ng serve

# Run tests
npm test

# Run linting
npm run lint

# Build for production
npm run build
```

### Database Migrations

Migrations are automatically applied on backend startup. Manual migration:

```bash
# Install migrate tool
brew install golang-migrate  # macOS
# or download from https://github.com/golang-migrate/migrate

# Run migrations
migrate -path backend/internal/infrastructure/postgres/migrations \
        -database "postgres://jlpt5_user:jlpt5_password@localhost:5432/jlpt5_dev?sslmode=disable" \
        up

# Rollback
migrate -path backend/internal/infrastructure/postgres/migrations \
        -database "postgres://jlpt5_user:jlpt5_password@localhost:5432/jlpt5_dev?sslmode=disable" \
        down 1
```

## CI/CD with GitHub Actions

This project uses **semantic commit messages** to trigger automated builds and deployments.

### Commit Message Keywords

| Keyword | Description | Triggers |
|---------|-------------|----------|
| `feat:` | New feature | Build + Test |
| `fix:` | Bug fix | Build + Test |
| `build:` | Build system changes | Build + Test |
| `backend:` | Backend-specific | Build + Test Backend |
| `frontend:` | Frontend-specific | Build + Test Frontend |
| `deploy:` | Deploy changes | Build + Test + Deploy |
| `deploy:all` | Deploy everything | Build All + Deploy All |
| `release:` | Create release | Build All + Deploy All |
| `hotfix:` | Emergency fix | Build All + Deploy All |

### Example Workflows

```bash
# Develop a new feature (builds but doesn't deploy)
git commit -m "feat: add vocabulary quiz feature"
git push

# Deploy backend changes
git commit -m "deploy: update authentication service"
git push

# Full release deployment
git commit -m "release: v2.0.0 - quiz improvements"
git push

# Emergency hotfix
git commit -m "hotfix: critical database connection issue"
git push
```

See [GitHub Actions Documentation](./.github/workflows/README.md) for detailed CI/CD information.

## Deployment

### Docker Compose (Development)

```bash
# Start
docker compose up -d

# Stop
docker compose down

# Clean everything (including volumes)
docker compose down -v
```

### Helm/Kubernetes (Production)

```bash
# Install to production
helm install jlpt5 ./helm/jlpt5 \
  --namespace jlpt5-prod \
  --create-namespace \
  --set backend.replicaCount=3 \
  --set frontend.replicaCount=3

# Upgrade
helm upgrade jlpt5 ./helm/jlpt5 \
  --namespace jlpt5-prod

# Rollback
helm rollback jlpt5 --namespace jlpt5-prod

# Uninstall
helm uninstall jlpt5 --namespace jlpt5-prod
```

### Environment Configuration

| Environment | Namespace | Replicas (BE/FE) | Resources |
|-------------|-----------|------------------|-----------|
| Development | `jlpt5-dev` | 1/1 | Minimal |
| Staging | `jlpt5-staging` | 2/2 | Medium |
| Production | `jlpt5-prod` | 3/3 | High + Autoscaling |

## API Documentation

### Base URL
```
Local:      http://localhost:8080/api/v1
Docker:     http://localhost:8080/api/v1
Kubernetes: http://<service-ip>:8080/api/v1
```

### Authentication Endpoints

```bash
# Register
POST /api/v1/auth/register
Content-Type: application/json
{
  "username": "user@example.com",
  "password": "SecurePass123!"
}

# Login
POST /api/v1/auth/login
Content-Type: application/json
{
  "username": "user@example.com",
  "password": "SecurePass123!"
}

# Refresh Token
POST /api/v1/auth/refresh
Authorization: Bearer <refresh_token>
```

### Vocabulary Endpoints

```bash
# Get all vocabulary
GET /api/v1/vocabulary

# Get vocabulary due for review
GET /api/v1/vocabulary/due

# Submit review
POST /api/v1/vocabulary/:id/review
Content-Type: application/json
{
  "rating": 3,  # 0-5 (SM-2 algorithm)
  "timeSpent": 5000  # milliseconds
}
```

### Health Check

```bash
# Check API health
GET /health

# Response
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-12-06T16:00:00Z"
}
```

## Testing

### Backend Tests

```bash
cd backend

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...
```

### Frontend Tests

```bash
cd frontend

# Run unit tests
npm test

# Run with coverage
npm test -- --code-coverage

# Run e2e tests
npm run e2e
```

## Project Status

**Current Phase: ✅ Production Ready**

- [x] Project structure and architecture
- [x] Backend API (14+ endpoints, all tested)
- [x] Frontend development (Angular 18)
- [x] Database design and migrations
- [x] JWT authentication with refresh tokens
- [x] Docker containerization
- [x] Docker Compose setup
- [x] Kubernetes Helm charts
- [x] GitHub Actions CI/CD
- [x] Automated testing
- [x] Security scanning
- [x] Production deployment
- [ ] End-to-end tests
- [ ] Performance optimization
- [ ] Monitoring and observability

## Production Deployment TODO

Before deploying to production, complete these critical tasks. See [PRODUCTION-CHECKLIST.md](./PRODUCTION-CHECKLIST.md) for detailed instructions.

### 🔴 Critical (Must Complete)

- [ ] **Secrets Management** - Replace hardcoded passwords with Kubernetes Secrets or external secrets manager
- [ ] **Container Registry** - Push images to production registry (GHCR, ECR, Docker Hub) with proper versioning
- [ ] **SSL/TLS Certificates** - Install cert-manager and configure Let's Encrypt for HTTPS
- [ ] **Domain & DNS** - Configure production domain and DNS records
- [ ] **Database Backups** - Set up automated daily backups with tested restore procedure

### 🟡 Important (Should Complete)

- [ ] **Monitoring** - Deploy Prometheus + Grafana for metrics and alerting
- [ ] **Logging** - Set up centralized logging (Loki, ELK stack)
- [ ] **Rate Limiting** - Configure API rate limits to prevent abuse
- [ ] **Security Hardening** - Enable network policies, pod security standards, non-root containers
- [ ] **Performance Testing** - Run load tests to verify scalability
- [ ] **Resource Tuning** - Optimize CPU/memory limits and autoscaling thresholds

### 🟢 Nice to Have

- [ ] **CDN** - Configure CloudFlare or AWS CloudFront for static assets
- [ ] **Caching** - Add Redis for session management and API caching
- [ ] **Error Tracking** - Integrate Sentry or Rollbar for error monitoring
- [ ] **APM** - Set up application performance monitoring
- [ ] **Analytics** - Add user analytics (Google Analytics, Mixpanel)
- [ ] **Email Service** - Configure transactional email (SendGrid, AWS SES)

**Quick Start**: For minimum viable production, focus on the 5 critical items. Estimated time: 2-3 days.

## Security

- JWT tokens with secure signing
- Password hashing with bcrypt
- SQL injection prevention (prepared statements)
- XSS protection
- CORS configuration
- Security scanning with Trivy
- Container image vulnerability scanning

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Commit using semantic messages (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

All PRs must pass:
- Automated tests
- Code linting
- Security scanning
- Commit message validation

## License

MIT License - See [LICENSE](LICENSE) file for details

## Resources

- [JLPT Official Website](https://www.jlpt.jp/e/)
- [SM-2 Spaced Repetition Algorithm](https://en.wikipedia.org/wiki/SuperMemo#SM-2_algorithm)
- [Angular Documentation](https://angular.io/docs)
- [Go Documentation](https://go.dev/doc/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Helm Documentation](https://helm.sh/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

## Support

For issues, questions, or contributions:
- Open an issue on GitHub
- Check existing documentation in `/docs`
- Review Helm chart docs: [./helm/jlpt5/README.md](./helm/jlpt5/README.md)
- Review CI/CD docs: [./.github/workflows/README.md](./.github/workflows/README.md)

---

**Happy Learning! 🎌📚**
