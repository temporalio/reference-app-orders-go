#!/bin/sh

echo "Installing Temporal via Helm..."

helm install --wait \
    --namespace temporal \
    --create-namespace \
    --set server.replicaCount=1 \
    --set server.config.namespaces.create=true \
    --set cassandra.config.cluster_size=1 \
    --set elasticsearch.replicas=1 \
    --set prometheus.enabled=false \
    --set grafana.enabled=false \
    --repo https://go.temporal.io/helm-charts \
    temporal temporal --timeout 15m

echo "Done."
