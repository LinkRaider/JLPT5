# GitHub Actions Workflows

This repository uses GitHub Actions for CI/CD with **commit message-based triggering**. Different workflows are triggered based on keywords in your commit messages.

## Commit Message Keywords

### Build Keywords

These keywords trigger builds and tests:

- `build:` - Trigger a build for changed components
- `feat:` - New feature (triggers build)
- `fix:` - Bug fix (triggers build)
- `chore:` - Maintenance tasks (triggers build)
- `backend:` - Backend-specific changes
- `frontend:` - Frontend-specific changes

### Deploy Keywords

These keywords trigger deployments:

- `deploy:` - Deploy changed components
- `deploy:all` - Deploy all services (backend, frontend, database)
- `release:` - Create a release and deploy everything
- `hotfix:` - Emergency fix deployment

## Workflows

### 1. Backend CI/CD (`backend-ci.yml`)

**Triggers:**
- Push to `main` or `develop` with backend file changes
- Commit messages starting with: `build:`, `backend:`, `feat:`, `fix:`, `chore:`
- Pull requests affecting backend files

**Jobs:**
1. **Test** - Runs Go tests with PostgreSQL
2. **Build & Push** - Builds Docker image and pushes to GitHub Container Registry
3. **Deploy** - Deploys to Kubernetes using Helm (on `deploy:` or `release:`)

**Example:**
```bash
git commit -m "feat: add user authentication endpoint"
git push  # Triggers test + build + push

git commit -m "deploy: update authentication service"
git push  # Triggers test + build + push + deploy
```

### 2. Frontend CI/CD (`frontend-ci.yml`)

**Triggers:**
- Push to `main` or `develop` with frontend file changes
- Commit messages starting with: `build:`, `frontend:`, `feat:`, `fix:`, `chore:`, `style:`
- Pull requests affecting frontend files

**Jobs:**
1. **Test** - Runs Angular tests and linting
2. **Build & Push** - Builds Docker image and pushes to registry
3. **Deploy** - Deploys to Kubernetes using Helm

**Example:**
```bash
git commit -m "feat: add quiz component"
git push  # Triggers test + build + push

git commit -m "deploy: update frontend styling"
git push  # Triggers test + build + push + deploy
```

### 3. Deploy All Services (`deploy-all.yml`)

**Triggers:**
- Commit messages starting with: `deploy:all`, `release:`, `hotfix:`
- Manual workflow dispatch from GitHub UI

**Jobs:**
1. **Build Backend** - Builds backend Docker image
2. **Build Frontend** - Builds frontend Docker image
3. **Deploy** - Deploys all services to Kubernetes
4. **Notify** - Sends deployment status notification

**Example:**
```bash
git commit -m "release: v1.2.0 with new quiz features"
git push  # Builds and deploys everything

# Or manually from GitHub:
# Actions → Deploy All Services → Run workflow
```

### 4. Pull Request Checks (`pr-check.yml`)

**Triggers:**
- All pull requests to `main` or `develop`

**Jobs:**
1. **Validate Commits** - Ensures commit messages follow semantic format
2. **Lint Backend** - Runs Go linting
3. **Lint Frontend** - Runs Angular linting
4. **Security Scan** - Runs Trivy vulnerability scanner
5. **Check Helm** - Validates Helm charts

## Semantic Commit Message Format

All commits should follow this format:

```
<type>: <description>

[optional body]

[optional footer]
```

### Valid Types

| Type | Description | Triggers Build | Triggers Deploy |
|------|-------------|----------------|-----------------|
| `feat` | New feature | ✅ | ❌ |
| `fix` | Bug fix | ✅ | ❌ |
| `docs` | Documentation | ❌ | ❌ |
| `style` | Code style changes | ✅ (FE) | ❌ |
| `refactor` | Code refactoring | ✅ | ❌ |
| `test` | Add/update tests | ❌ | ❌ |
| `chore` | Maintenance | ✅ | ❌ |
| `build` | Build system changes | ✅ | ❌ |
| `ci` | CI configuration | ❌ | ❌ |
| `perf` | Performance improvement | ✅ | ❌ |
| `backend` | Backend-specific | ✅ | ❌ |
| `frontend` | Frontend-specific | ✅ | ❌ |
| `deploy` | Deploy changes | ✅ | ✅ |
| `deploy:all` | Deploy everything | ✅ | ✅ |
| `release` | Create release | ✅ | ✅ |
| `hotfix` | Emergency fix | ✅ | ✅ |

### Examples

```bash
# Feature development (builds but doesn't deploy)
git commit -m "feat: add vocabulary quiz feature"

# Bug fix (builds but doesn't deploy)
git commit -m "fix: resolve authentication token expiry issue"

# Deploy backend only
git commit -m "deploy: update backend API endpoints"

# Deploy everything
git commit -m "deploy:all: roll out v2.0 with new features"

# Release
git commit -m "release: v2.1.0 - quiz improvements and bug fixes"

# Hotfix
git commit -m "hotfix: critical database connection issue"
```

## Environment Configuration

### Required Secrets

Configure these in GitHub Settings → Secrets and variables → Actions:

| Secret | Description |
|--------|-------------|
| `KUBE_CONFIG` | Base64-encoded Kubernetes config file |
| `GITHUB_TOKEN` | Automatically provided by GitHub |

### Optional Secrets

| Secret | Description |
|--------|-------------|
| `CODECOV_TOKEN` | Code coverage reporting token |
| `SLACK_WEBHOOK` | Slack notification webhook |

## Manual Deployment

You can manually trigger deployment from the GitHub Actions UI:

1. Go to **Actions** tab
2. Select **Deploy All Services**
3. Click **Run workflow**
4. Choose environment (production/staging/development)
5. Click **Run workflow**

## Deployment Environments

| Environment | Namespace | Backend Replicas | Frontend Replicas |
|-------------|-----------|------------------|-------------------|
| Production | `jlpt5-prod` | 3 | 3 |
| Staging | `jlpt5-staging` | 2 | 2 |
| Development | `jlpt5-dev` | 1 | 1 |

## Workflow Status

Check workflow status:
- Repository badges show build status
- **Actions** tab shows detailed logs
- Failed workflows send notifications

## Best Practices

1. **Always use semantic commit messages** - PRs will fail without them
2. **Test locally first** - Run `docker compose up` before pushing
3. **Use feature branches** - Create PRs for review
4. **Deploy incrementally** - Use `deploy:` for single service, `deploy:all` for full deployment
5. **Monitor deployments** - Check Kubernetes pods after deployment

## Troubleshooting

### Build fails but tests pass locally

- Check GitHub Actions logs for specific errors
- Ensure Docker image builds successfully
- Verify all dependencies are in package files

### Deployment fails

- Check `KUBE_CONFIG` secret is correctly set
- Verify Kubernetes cluster is accessible
- Check Helm chart values are correct
- Review pod logs: `kubectl logs -n jlpt5-prod <pod-name>`

### Commit message validation fails

- Ensure commit message starts with valid type
- Use lowercase for type
- Add colon and space after type
- Example: `feat: add new feature` ✅
- Not: `Feat add new feature` ❌

## Additional Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
