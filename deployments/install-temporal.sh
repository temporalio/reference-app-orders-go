#!/bin/sh

echo "Installing Temporal via Helm..."

helm install \
    --namespace temporal \
    --create-namespace \
    --set server.replicaCount=1 \
    --set cassandra.config.cluster_size=1 \
    --set elasticsearch.replicas=1 \
    --set prometheus.enabled=false \
    --set grafana.enabled=false \
    --repo https://go.temporal.io/helm-charts \
    temporal temporal --timeout 15m

echo "Waiting for Temporal admintools to be ready..."

kubectl rollout status -n temporal deployment/temporal-admintools                                                                                                                                                            git:(rh-aws-eks) ✗

echo "Creating default namespace..."

kubectl exec -ti -n temporal deployment/temporal-admintools -- \
	temporal operator namespace create -n default

echo "Done."