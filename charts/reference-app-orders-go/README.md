# Reference App Orders Go Helm Chart

A Helm chart for deploying the reference-app-orders-go application, which demonstrates a distributed order management system built with Temporal.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- A running Temporal cluster

## Installation

### Install from OCI Registry

```bash
helm install my-orders-app oci://ghcr.io/temporalio/charts/reference-app-orders-go
```

### Install from Local Chart

```bash
helm install my-orders-app ./charts/reference-app-orders-go
```

## Limitations

⚠️ **Important**: The `main.api.replicaCount` is currently limited to `1` until database support is implemented. Setting it higher will cause the chart installation to fail with an error message. This prevents potential data consistency issues when multiple API instances are running without proper database coordination.

## Configuration

The following table lists the configurable parameters and their default values:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `nameOverride` | Override the name of the chart | `""` |
| `fullnameOverride` | Override the full name of the chart | `""` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `main.worker.replicaCount` | Number of main worker replicas | `1` |
| `main.worker.image.repository` | Main worker image repository | `ghcr.io/temporalio/reference-app-orders-go` |
| `main.worker.image.tag` | Main worker image tag | `latest` |
| `main.api.replicaCount` | Number of main API replicas (limited to 1) | `1` |
| `main.api.image.repository` | Main API image repository | `ghcr.io/temporalio/reference-app-orders-go` |
| `main.api.image.tag` | Main API image tag | `latest` |
| `billing.worker.replicaCount` | Number of billing worker replicas | `1` |
| `billing.api.replicaCount` | Number of billing API replicas | `1` |
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

### Workers
- **Main Worker**: Handles order and shipment workflows
- **Billing Worker**: Handles billing workflows

### APIs
- **Main API**: Exposes order, shipment, and fraud APIs
- **Billing API**: Exposes billing API

## Example Values

```yaml
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

## Monitoring

When `serviceMonitor.enabled` is set to `true`, the chart creates ServiceMonitor resources for Prometheus to scrape metrics from the application endpoints. ServiceMonitors are deployed in the same namespace as the application services.

## Uninstallation

```bash
helm uninstall my-orders-app
``` 