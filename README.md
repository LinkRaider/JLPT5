# JLPT5 Learning Application

A full-stack web application for learning Japanese at JLPT N5 level, featuring vocabulary flashcards with spaced repetition, grammar lessons, practice quizzes, and progress tracking.

## Technology Stack

- **Frontend**: Angular (TypeScript)
- **Backend**: Go (Golang)
- **Database**: PostgreSQL 16
- **Infrastructure**: Kubernetes (Minikube) + Helm
- **Containerization**: Docker

## Features

- **Vocabulary Flashcards**: Study JLPT N5 vocabulary with intelligent spaced repetition (SM-2 algorithm)
- **Grammar Lessons**: Comprehensive grammar points with examples and explanations
- **Practice Quizzes**: Multiple choice and fill-in-the-blank questions to test knowledge
- **Progress Tracking**: Track learning statistics, study streaks, and overall progress

## Architecture

### Monorepo Structure

```
JLPT5/
├── frontend/          # Angular application
├── backend/           # Go REST API
├── helm/jlpt5/       # Helm chart for Kubernetes deployment
├── scripts/          # Setup and deployment scripts
└── docker-compose.yml # Local development environment
```

### Backend Architecture

- **Layered Design**: Handler → Service → Repository
- **Clean Architecture**: Domain-driven design with dependency inversion
- **RESTful API**: Versioned endpoints (`/api/v1`)
- **JWT Authentication**: Secure token-based authentication

### Frontend Architecture

- **Feature Modules**: Lazy-loaded modules for each feature area
- **Reactive Design**: RxJS-based state management
- **Angular Material**: Modern, responsive UI components

### Database Design

- **PostgreSQL**: Relational database with proper normalization
- **Migration-Based**: Version-controlled schema management
- **Spaced Repetition**: Dedicated tables for learning algorithm data

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker Desktop
- Minikube
- kubectl
- Helm 3

### Quick Start (Local Development with Docker Compose)

```bash
# Start all services
docker-compose up

# Access the application
# Frontend: http://localhost:4200
# Backend: http://localhost:8080
```

### Kubernetes Deployment (Minikube)

```bash
# Setup Minikube environment
make setup

# Build Docker images
make build

# Deploy to dev environment
make deploy-dev

# Access the application
# http://jlpt5-dev.local
```

## Development

### Backend Development

```bash
cd backend

# Initialize dependencies
go mod download

# Run locally
go run cmd/api/main.go

# Run tests
go test ./...

# Run with hot reload
air
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm start

# Run tests
npm test

# Build for production
npm run build
```

### Database Migrations

Migrations are located in `backend/internal/infrastructure/postgres/migrations/`

```bash
# Migrations run automatically on backend startup
# Or manually via migration tool
```

## Deployment Environments

### Development
- **Namespace**: `jlpt5-dev`
- **URL**: http://jlpt5-dev.local
- **Replicas**: 1 (all services)
- **Resources**: Minimal

### Production
- **Namespace**: `jlpt5`
- **URL**: http://jlpt5.local
- **Replicas**: Frontend (2), Backend (3), PostgreSQL (1)
- **Resources**: Production-grade with autoscaling

## API Documentation

### Base URL
```
http://localhost:8080/api/v1  (local)
http://jlpt5-dev.local/api    (dev)
```

### Endpoints

#### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login and get JWT token
- `POST /auth/refresh` - Refresh JWT token

#### Vocabulary
- `GET /vocabulary` - List vocabulary items
- `GET /vocabulary/due` - Get items due for review
- `POST /vocabulary/:id/review` - Submit review result

#### Grammar
- `GET /grammar` - List grammar lessons
- `GET /grammar/:id` - Get lesson details

#### Quizzes
- `GET /quizzes` - List available quizzes
- `POST /quizzes/:id/start` - Start a quiz
- `POST /quizzes/:id/submit` - Submit answers

#### Progress
- `GET /progress/stats` - Get user statistics
- `GET /progress/vocabulary` - Vocabulary progress
- `GET /progress/streak` - Study streak info

## Makefile Commands

```bash
make setup        # Setup Minikube environment
make build        # Build Docker images
make deploy-dev   # Deploy to dev environment
make deploy-prod  # Deploy to prod environment
make clean        # Clean up resources
```

## Project Status

Current Phase: **Foundation & Backend Development**

- [x] Project structure setup
- [x] Git initialization
- [ ] Backend API implementation
- [ ] Frontend development
- [ ] Database migrations and seeding
- [ ] Docker containerization
- [ ] Kubernetes deployment
- [ ] Documentation

## Contributing

This is a personal learning project. Contributions are welcome but please open an issue first to discuss proposed changes.

## License

MIT License - See LICENSE file for details

## Resources

- [JLPT Official Website](https://www.jlpt.jp/e/)
- [SM-2 Spaced Repetition Algorithm](https://en.wikipedia.org/wiki/SuperMemo#SM-2_algorithm)
- [Angular Documentation](https://angular.io/docs)
- [Go Documentation](https://go.dev/doc/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
