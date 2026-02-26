# OpenShift Deployment Guide

## Overview
This application is designed to run on OpenShift with external configuration and persistent storage.

## Required OpenShift Resources

### 1. ConfigMap (for .env variables)
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-monolith-config
data:
  APP_NAME: "OpenShift Go Monolith"
  APP_ENV: "production"
```

### 2. Secret (for sensitive data)
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: go-monolith-secrets
type: Opaque
stringData:
  DB_USER: "app_user"
```

### 3. ConfigMap (for config.json)
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-monolith-json-config
data:
  config.json: |
    {
      "application": {
        "name": "OpenShift Go Monolith",
        "environment": "production"
      },
      "database": {
        "user": "app_user"
      }
    }
```

### 4. PersistentVolumeClaim (for data directory)
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: go-monolith-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

## Deployment Configuration

### Deployment with Volume Mounts
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-monolith
spec:
  replicas: 1
  selector:
    matchLabels:
      app: go-monolith
  template:
    metadata:
      labels:
        app: go-monolith
    spec:
      containers:
      - name: go-monolith
        image: your-registry/openshift-go-monolith:latest
        ports:
        - containerPort: 8080
        
        # Environment variables from ConfigMap
        envFrom:
        - configMapRef:
            name: go-monolith-config
        
        # Sensitive environment variables from Secret
        env:
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: go-monolith-secrets
              key: DB_USER
        
        # Volume mounts
        volumeMounts:
        # Mount config.json from ConfigMap
        - name: json-config
          mountPath: /app/config.json
          subPath: config.json
          readOnly: true
        
        # Mount .env from ConfigMap (optional, since we use envFrom)
        - name: env-config
          mountPath: /app/.env
          subPath: .env
          readOnly: true
        
        # Mount PersistentVolume for data directory
        - name: data-volume
          mountPath: /app/data
        
        # Health checks
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        
        # Resource limits
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
      
      volumes:
      # ConfigMap volume for config.json
      - name: json-config
        configMap:
          name: go-monolith-json-config
      
      # ConfigMap volume for .env (optional)
      - name: env-config
        configMap:
          name: go-monolith-config
          items:
          - key: .env
            path: .env
      
      # PersistentVolume for data
      - name: data-volume
        persistentVolumeClaim:
          claimName: go-monolith-data
```

### Service
```yaml
apiVersion: v1
kind: Service
metadata:
  name: go-monolith-service
spec:
  selector:
    app: go-monolith
  ports:
  - protocol: TCP
    port: 8080
    targetPort: 8080
  type: ClusterIP
```

### Route (OpenShift)
```yaml
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: go-monolith-route
spec:
  to:
    kind: Service
    name: go-monolith-service
  port:
    targetPort: 8080
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
```

## Build and Push Image

### Using Podman
```bash
# Build the image
podman build -t openshift-go-monolith:latest -f Containerfile .

# Tag for your registry
podman tag openshift-go-monolith:latest your-registry/openshift-go-monolith:latest

# Push to registry
podman push your-registry/openshift-go-monolith:latest
```

### Using OpenShift BuildConfig
```yaml
apiVersion: build.openshift.io/v1
kind: BuildConfig
metadata:
  name: go-monolith-build
spec:
  source:
    type: Git
    git:
      uri: https://github.com/your-org/your-repo.git
      ref: main
    contextDir: openshift-go-monolith
  strategy:
    type: Docker
    dockerStrategy:
      dockerfilePath: Containerfile
  output:
    to:
      kind: ImageStreamTag
      name: go-monolith:latest
```

## Deployment Steps

1. Create the ConfigMaps and Secrets:
```bash
oc create -f configmap.yaml
oc create -f secret.yaml
oc create -f config-json-configmap.yaml
```

2. Create the PersistentVolumeClaim:
```bash
oc create -f pvc.yaml
```

3. Deploy the application:
```bash
oc create -f deployment.yaml
oc create -f service.yaml
oc create -f route.yaml
```

4. Verify deployment:
```bash
oc get pods
oc get svc
oc get route
```

## Accessing the Application

Get the route URL:
```bash
oc get route go-monolith-route -o jsonpath='{.spec.host}'
```

Access the dashboard:
```
https://<route-url>
```

## Monitoring

View logs:
```bash
oc logs -f deployment/go-monolith
```

Check health:
```bash
curl https://<route-url>/health
```

View stats:
```bash
curl https://<route-url>/api/stats
```

## Notes

- The `.env` and `config.json` files in the repository are for local development only
- In OpenShift, these are provided via ConfigMaps and Secrets
- The `data/` directory is mounted as a PersistentVolume for log storage
- The application runs as non-root user (UID 1001) for security
- All volumes support arbitrary UIDs for OpenShift compatibility
