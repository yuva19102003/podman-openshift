# OpenShift Configuration Files

This directory contains all the Kubernetes/OpenShift manifests needed to deploy the Go Monolith application.

## Files Overview

### Configuration
- `secret.yaml` - Sensitive data (DB credentials, API keys)
- `configmap.yaml` - Environment variables as key-value pairs
- `configmap-env.yaml` - Environment variables as .env file
- `configmap-json.yaml` - Application configuration as config.json

### Storage
- `pvc.yaml` - PersistentVolumeClaim for data storage

### Application
- `deployment.yaml` - Application deployment with volume mounts
- `service.yaml` - Service to expose the application
- `route.yaml` - OpenShift Route for external access

## Deployment Order

Apply the manifests in this order:

### 1. Create Secrets and ConfigMaps
```bash
oc apply -f secret.yaml
oc apply -f configmap.yaml
oc apply -f configmap-env.yaml
oc apply -f configmap-json.yaml
```

### 2. Create PersistentVolumeClaim
```bash
oc apply -f pvc.yaml
```

### 3. Deploy Application
```bash
oc apply -f deployment.yaml
oc apply -f service.yaml
oc apply -f route.yaml
```

## Quick Deploy (All at Once)
```bash
oc apply -f .
```

## Verify Deployment

Check all resources:
```bash
oc get all -l app=go-monolith
```

Check ConfigMaps:
```bash
oc get configmap -l app=go-monolith
```

Check Secrets:
```bash
oc get secret -l app=go-monolith
```

Check PVC:
```bash
oc get pvc app-storage
```

View pod logs:
```bash
oc logs -f deployment/go-monolith
```

Get route URL:
```bash
oc get route go-monolith-route -o jsonpath='{.spec.host}'
```

## Access the Application

Once deployed, get the route:
```bash
echo "https://$(oc get route go-monolith-route -o jsonpath='{.spec.host}')"
```

## Update Configuration

### Update ConfigMap
```bash
oc apply -f configmap.yaml
# Restart pods to pick up changes
oc rollout restart deployment/go-monolith
```

### Update Secret
```bash
oc apply -f secret.yaml
# Restart pods to pick up changes
oc rollout restart deployment/go-monolith
```

## Troubleshooting

### Check pod status
```bash
oc get pods -l app=go-monolith
```

### View pod events
```bash
oc describe pod -l app=go-monolith
```

### Check logs
```bash
oc logs -f deployment/go-monolith
```

### Check mounted volumes
```bash
oc exec deployment/go-monolith -- ls -la /app
oc exec deployment/go-monolith -- cat /app/config.json
oc exec deployment/go-monolith -- cat /app/.env
oc exec deployment/go-monolith -- ls -la /app/data/log
```

### Check environment variables
```bash
oc exec deployment/go-monolith -- env | grep APP
```

## Clean Up

Remove all resources:
```bash
oc delete -f .
```

Or individually:
```bash
oc delete route go-monolith-route
oc delete service go-monolith-service
oc delete deployment go-monolith
oc delete pvc app-storage
oc delete configmap go-monolith-config go-monolith-env-file go-monolith-json-config
oc delete secret go-monolith-secrets
```

## Notes

- The namespace is set to `romanyuvan-dev` - update if needed
- PVC uses the existing `app-storage` claim
- ConfigMaps and Secrets are mounted as files in the container
- Environment variables are injected from ConfigMap and Secret
- Health checks are configured on `/health` endpoint
- TLS is enabled with edge termination on the route
