# Reference App Orders Go Helm Chart

A Helm chart for deploying the reference-app-orders-go application, which demonstrates a distributed order management system built with Temporal.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- A running Temporal Service

## Installation

### Install from OCI Registry

```bash
helm install -n oms --create-namespace oms oci://ghcr.io/temporalio/charts/reference-app-orders-go
```

### Install from Local Chart

```bash
helm install -n oms --create-namespace oms ./charts/reference-app-orders-go
```

## Configuration

The following table lists the configurable parameters and their default values:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `nameOverride` | Override the name of the chart | `""` |
| `fullnameOverride` | Override the full name of the chart | `""` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `mongodb.enabled` | Enable MongoDB deployment | `true` |
| `mongodb.image.repository` | MongoDB image repository | `mongo` |
| `mongodb.image.tag` | MongoDB image tag | `6` |
| `mongodb.resources.limits.cpu` | MongoDB CPU limit | `500m` |
| `mongodb.resources.limits.memory` | MongoDB memory limit | `512Mi` |
| `mongodb.resources.requests.cpu` | MongoDB CPU request | `250m` |
| `mongodb.resources.requests.memory` | MongoDB memory request | `256Mi` |
| `mongodb.persistence.enabled` | Enable MongoDB persistent storage | `true` |
| `mongodb.persistence.size` | MongoDB storage size | `100Mi` |
| `mongodb.persistence.storageClass` | MongoDB storage class | `""` |
| `main.worker.replicaCount` | Number of main worker replicas | `1` |
| `main.worker.image.repository` | Main worker image repository | `ghcr.io/temporalio/reference-app-orders-go` |
| `main.worker.image.tag` | Main worker image tag | `Chart.appVersion` |
| `main.api.replicaCount` | Number of main API replicas | `1` |
| `main.api.image.repository` | Main API image repository | `ghcr.io/temporalio/reference-app-orders-go` |
| `main.api.image.tag` | Main API image tag | `Chart.appVersion` |
| `billing.worker.replicaCount` | Number of billing worker replicas | `1` |
| `billing.worker.image.repository` | Billing worker image repository | `ghcr.io/temporalio/reference-app-orders-go` |
| `billing.worker.image.tag` | Billing worker image tag | `Chart.appVersion` |
| `billing.api.replicaCount` | Number of billing API replicas | `1` |
| `billing.api.image.repository` | Billing API image repository | `ghcr.io/temporalio/reference-app-orders-go` |
| `billing.api.image.tag` | Billing API image tag | `Chart.appVersion` |
| `web.replicaCount` | Number of web application replicas | `1` |
| `web.image.repository` | Web application image repository | `ghcr.io/temporalio/reference-app-orders-web` |
| `web.image.tag` | Web application image tag | `latest` |
| `temporal.address` | Temporal frontend address | `temporal-frontend:7233` |
| `temporal.namespace` | Temporal namespace | `default` |
| `services.billing.port` | Billing API port | `8081` |
| `services.order.port` | Order API port | `8082` |
| `services.shipment.port` | Shipment API port | `8083` |
| `services.fraud.port` | Fraud API port | `8084` |
| `metrics.enabled` | Enable metrics collection | `true` |
| `metrics.port` | Metrics port | `9090` |
| `serviceMonitor.enabled` | Enable ServiceMonitor for Prometheus | `false` |
| `serviceMonitor.interval` | Scrape interval | `30s` |
| `encryptionKeyID` | Optional encryption key ID for payload encryption | `""` |

## Services

This chart deploys the following services:

### Database
- **MongoDB**: Provides a shared cache for the main API service

### Workers
- **Main Worker**: Handles order and shipment workflows
- **Billing Worker**: Handles billing workflows

### APIs
- **Main API**: Exposes order, shipment, and fraud APIs
- **Billing API**: Exposes billing API

### Web Application
- **Web**: Frontend web application that provides a user interface for the order management system

## Example Values

```yaml
# Scale API services horizontally
main:
  api:
    replicaCount: 3

# Enable web application
web:
  replicaCount: 2
  image:
    tag: "v2.0.0"

# Custom MongoDB configuration
mongodb:
  resources:
    limits:
      cpu: "1000m"
      memory: "1Gi"
  persistence:
    size: 1Gi
    storageClass: "fast-ssd"

# Custom image
main:
  worker:
    image:
      repository: my-registry/orders-app
      tag: v1.0.0

# Enable monitoring
serviceMonitor:
  enabled: true
  interval: 15s

# Connect to external Temporal cluster
temporal:
  address: my-temporal-cluster:7233
  namespace: orders

# Enable payload encryption
encryptionKeyID: "my-encryption-key"
```

## Database

The chart includes a MongoDB deployment that serves as a cache for the main API.

## Monitoring

When `serviceMonitor.enabled` is set to `true`, the chart creates ServiceMonitor resources for Prometheus to scrape metrics from the application endpoints. ServiceMonitors are deployed in the same namespace as the application services.

## Uninstallation

```bash
helm uninstall oms
```

**Note**: This will also remove the MongoDB StatefulSet and its associated PersistentVolumeClaim. Make sure to backup any important data before uninstalling. 