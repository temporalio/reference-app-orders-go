#!/usr/bin/env bash

set -euo pipefail

rm -rf ./k8s
mkdir -p ./k8s

# Remove restart policy from all services
# Set controller type to statefulset for apis, which maintain a cache on disk
# Set service type to loadbalancer for web service, which should be exposed outside the cluster
yq \
    '(del(.services[].restart)) |
     ((.services | (.billing-api, .main-api)).labels += {"kompose.controller.type":"statefulset"}) |
     (.services.web.labels += {"kompose.service.type":"loadbalancer"})
     ' \
    docker-compose-split.yaml | \
    kompose -f - -o k8s convert -n oms --with-kompose-annotation=false

# Rename the web-tcp service to web
mv ./k8s/web-tcp-service.yaml ./k8s/web-service.yaml
yq '(.metadata.name, .metadata.labels.["io.kompose.service"]) |= "web"' -i ./k8s/web-service.yaml

# Use standard kubernetes labels
for f in ./k8s/*.yaml; do
    yq -i \
    '((.. | select(has("io.kompose.service")).["io.kompose.service"] | key) = "app.kubernetes.io/component") |
     ((.. | select(has("app.kubernetes.io/component"))) += {"app.kubernetes.io/name":"oms"})
    ' $f
done

# Correct images
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

# Disable service links
for f in ./k8s/*-{deployment,statefulset}.yaml; do
    yq -i '.spec.template.spec.enableServiceLinks = false' $f
done

# Remove redundant defaults
for f in ./k8s/*.yaml; do
    yq -i 'del(.spec.template.spec.restartPolicy)' $f
done

# Update Temporal Address
for f in ./k8s/*-{statefulset,deployment}.yaml; do
    yq -i '(.spec.template.spec.containers[0].env[] | select(.name == "TEMPORAL_ADDRESS").value) |= "temporal-frontend.temporal:7233"' $f
done
