#!/usr/bin/env bash

set -euo pipefail

rm -rf ./k8s
mkdir -p ./k8s

# Remove restart policy from all services, it defaults to always in k8s.
# Set controller type to statefulset for apis, which maintain a cache on disk
yq \
    '(del(.services[].restart)) |
     ((.services | (.billing-api, .main-api)).labels += {"kompose.controller.type":"statefulset"})' \
    docker-compose-split.yaml | \
    kompose -f - -o k8s convert -n oms --with-kompose-annotation=false

# Translate kompose labels to more standard kubernetes labels
for f in ./k8s/*.yaml; do
    yq -i \
    '((.. | select(has("io.kompose.service")).["io.kompose.service"] | key) = "app.kubernetes.io/component") |
     ((.. | select(has("app.kubernetes.io/component"))) += {"app.kubernetes.io/name":"oms"})
    ' $f
done

# For Kubernetes we need to reference our published Docker images, we can't build in-place like docker-compose does.
for f in ./k8s/*-worker-deployment.yaml; do
    yq -i '.spec.template.spec.containers[0].image |= "ghcr.io/temporalio/reference-app-orders-go-worker:latest" |
           .spec.template.spec.containers[0].imagePullPolicy = "Always"' $f
done
for f in ./k8s/*-api-statefulset.yaml; do
    yq -i '.spec.template.spec.containers[0].image |= "ghcr.io/temporalio/reference-app-orders-go-api:latest" |
           .spec.template.spec.containers[0].imagePullPolicy = "Always"' $f
done
yq -i '.spec.template.spec.containers[0].image |= "ghcr.io/temporalio/reference-app-orders-go-codec-server:latest" |
       .spec.template.spec.containers[0].imagePullPolicy = "Always"' k8s/codec-server-deployment.yaml

# We don't rely on service links, so disable them to avoid collisions with our configuration environment variables.
for f in ./k8s/*-{deployment,statefulset}.yaml; do
    yq -i '.spec.template.spec.enableServiceLinks = false' $f
done

# Remove redundant defaults to make the manifests easier to read
for f in ./k8s/*.yaml; do
    yq -i 'del(.spec.template.spec.restartPolicy)' $f
done

# Update Temporal Address to assume Temporal is deployed in this Kubernetes cluster
for f in ./k8s/*-{statefulset,deployment}.yaml; do
    yq -i '(.spec.template.spec.containers[0].env[] | select(.name == "TEMPORAL_ADDRESS").value) |= "temporal-frontend.temporal:7233"' $f
done
