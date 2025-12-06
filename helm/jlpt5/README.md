# JLPT5 Helm Chart

A Helm chart for deploying the JLPT5 Japanese learning application on Kubernetes.

## Prerequisites

- Kubernetes 1.24+
- Helm 4.0+
- PV provisioner support in the underlying infrastructure (for PostgreSQL persistence)

## Components

This chart deploys the following components:

- **Frontend**: Angular 18 application serving the user interface
- **Backend**: Go API server handling business logic
- **PostgreSQL**: Database for storing application data

## Installing the Chart

To install the chart with the release name `jlpt5`:

```bash
helm install jlpt5 ./jlpt5
```

Or with custom values:

```bash
helm install jlpt5 ./jlpt5 -f custom-values.yaml
```

## Uninstalling the Chart

To uninstall/delete the `jlpt5` deployment:

```bash
helm uninstall jlpt5
```

## Configuration

The following table lists the configurable parameters and their default values.

### Global Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `global.nameOverride` | Override chart name | `""` |
| `global.fullnameOverride` | Override full name | `""` |

### PostgreSQL Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `postgresql.enabled` | Enable PostgreSQL | `true` |
| `postgresql.image.repository` | PostgreSQL image repository | `postgres` |
| `postgresql.image.tag` | PostgreSQL image tag | `16-alpine` |
| `postgresql.service.type` | Service type | `ClusterIP` |
| `postgresql.service.port` | Service port | `5432` |
| `postgresql.persistence.enabled` | Enable persistence | `true` |
| `postgresql.persistence.size` | PVC size | `10Gi` |
| `postgresql.auth.database` | Database name | `jlpt5_dev` |
| `postgresql.auth.username` | Database username | `jlpt5_user` |
| `postgresql.auth.password` | Database password | `jlpt5_password` |

### Backend Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `backend.enabled` | Enable backend | `true` |
| `backend.replicaCount` | Number of replicas | `2` |
| `backend.image.repository` | Backend image repository | `jlpt5-backend` |
| `backend.image.tag` | Backend image tag | `latest` |
| `backend.service.type` | Service type | `ClusterIP` |
| `backend.service.port` | Service port | `8080` |
| `backend.autoscaling.enabled` | Enable HPA | `false` |
| `backend.autoscaling.minReplicas` | Minimum replicas | `2` |
| `backend.autoscaling.maxReplicas` | Maximum replicas | `10` |

### Frontend Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `frontend.enabled` | Enable frontend | `true` |
| `frontend.replicaCount` | Number of replicas | `2` |
| `frontend.image.repository` | Frontend image repository | `jlpt5-frontend` |
| `frontend.image.tag` | Frontend image tag | `latest` |
| `frontend.service.type` | Service type | `LoadBalancer` |
| `frontend.service.port` | Service port | `80` |
| `frontend.autoscaling.enabled` | Enable HPA | `false` |

### Ingress Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ingress.enabled` | Enable Ingress | `false` |
| `ingress.className` | Ingress class name | `nginx` |
| `ingress.hosts` | Ingress hosts configuration | See values.yaml |

## Example Configurations

### Production Configuration

```yaml
postgresql:
  persistence:
    size: 50Gi
  resources:
    limits:
      cpu: 2000m
      memory: 2Gi
    requests:
      cpu: 1000m
      memory: 1Gi

backend:
  replicaCount: 3
  autoscaling:
    enabled: true
    minReplicas: 3
    maxReplicas: 20
  resources:
    limits:
      cpu: 2000m
      memory: 2Gi
    requests:
      cpu: 1000m
      memory: 1Gi

frontend:
  replicaCount: 3
  autoscaling:
    enabled: true
    minReplicas: 3
    maxReplicas: 20

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: jlpt5.yourdomain.com
      paths:
        - path: /api
          pathType: Prefix
          service: backend
        - path: /
          pathType: Prefix
          service: frontend
  tls:
    - secretName: jlpt5-tls
      hosts:
        - jlpt5.yourdomain.com
```

### Development Configuration

```yaml
postgresql:
  persistence:
    enabled: false

backend:
  replicaCount: 1
  image:
    pullPolicy: Always

frontend:
  replicaCount: 1
  image:
    pullPolicy: Always
  service:
    type: ClusterIP
```

## Building and Pushing Docker Images

Before deploying with Helm, you need to build and push the Docker images:

```bash
# Build backend image
cd backend
docker build -t your-registry/jlpt5-backend:latest .
docker push your-registry/jlpt5-backend:latest

# Build frontend image
cd ../frontend
docker build -t your-registry/jlpt5-frontend:latest .
docker push your-registry/jlpt5-frontend:latest
```

Then update your values.yaml or use --set:

```bash
helm install jlpt5 ./jlpt5 \
  --set backend.image.repository=your-registry/jlpt5-backend \
  --set frontend.image.repository=your-registry/jlpt5-frontend
```

## Accessing the Application

### Via LoadBalancer (default)

```bash
kubectl get svc jlpt5-frontend
```

### Via Port Forward

```bash
# Frontend
kubectl port-forward svc/jlpt5-frontend 4200:80

# Backend
kubectl port-forward svc/jlpt5-backend 8080:8080
```

### Via Ingress

If ingress is enabled, access via the configured hostname.

## Upgrading

To upgrade the deployment:

```bash
helm upgrade jlpt5 ./jlpt5 -f custom-values.yaml
```

## Troubleshooting

### Check pod status

```bash
kubectl get pods -l app.kubernetes.io/instance=jlpt5
```

### View logs

```bash
# Frontend logs
kubectl logs -l app.kubernetes.io/component=frontend

# Backend logs
kubectl logs -l app.kubernetes.io/component=backend

# Database logs
kubectl logs -l app.kubernetes.io/component=database
```

### Check events

```bash
kubectl get events --sort-by='.lastTimestamp'
```

## Security Considerations

**Important**: The default values include hardcoded passwords for development purposes. In production:

1. Use Kubernetes Secrets or a secrets management solution
2. Set `postgresql.auth.existingSecret` to use an existing secret
3. Update JWT secret in backend configuration
4. Enable TLS/SSL for all services
5. Configure network policies
6. Enable pod security policies

## License

This chart is provided as-is for the JLPT5 application.
