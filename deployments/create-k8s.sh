#!/usr/bin/env bash

set -euo pipefail

rm -rf ./k8s
mkdir -p ./k8s

# Remove restart policy from all services, it defaults to always in k8s.
# Set controller type to statefulset for apis, which maintain a cache on disk
yq \
    '(del(.services[].restart)) |
     (.services.mongo.labels += {"kompose.controller.type":"statefulset"})' \
    docker-compose-split.yaml | \
    kompose -f - -o k8s convert -n oms --with-kompose-annotation=false

# Add "oms-" prefix to deployment and service names to align with Helm chart naming
for f in ./k8s/*-deployment.yaml ./k8s/*-service.yaml ./k8s/*-statefulset.yaml ./k8s/*-persistentvolumeclaim.yaml; do
    yq -i '.metadata.name |= "oms-" + .' "$f"
done

# Update PVC references in statefulsets to use the new prefixed names
for f in ./k8s/*-statefulset.yaml; do
    yq -i '(.spec.template.spec.volumes[]?.persistentVolumeClaim.claimName) |= "oms-" + .' "$f"
done

# Update service references in deployments to use the new prefixed service names
for f in ./k8s/*-deployment.yaml; do
    yq -i '
    (.spec.template.spec.containers[0].env[]? | select(.name == "MONGO_URL").value) |= sub("mongo:27017"; "oms-mongo:27017") |
    (.spec.template.spec.containers[0].env[]? | select(.value | test("http://[^/]+")).value) |= sub("http://([^:/]+)"; "http://oms-\1")
    ' "$f"
done

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
for f in ./k8s/*-api-deployment.yaml; do
    yq -i '.spec.template.spec.containers[0].image |= "ghcr.io/temporalio/reference-app-orders-go-api:latest" |
           .spec.template.spec.containers[0].imagePullPolicy = "Always"' $f
done
yq -i '.spec.template.spec.containers[0].image |= "ghcr.io/temporalio/reference-app-orders-go-codec-server:latest" |
       .spec.template.spec.containers[0].imagePullPolicy = "Always"' k8s/codec-server-deployment.yaml

# We don't rely on service links, so disable them to avoid collisions with our configuration environment variables.
for f in ./k8s/*-deployment.yaml; do
    yq -i '.spec.template.spec.enableServiceLinks = false' $f
done

# Remove redundant defaults to make the manifests easier to read
for f in ./k8s/*.yaml; do
    yq -i 'del(.spec.template.spec.restartPolicy)' $f
done

# Update Temporal Address to assume Temporal is deployed in this Kubernetes cluster
for f in ./k8s/*-deployment.yaml; do
    yq -i '(.spec.template.spec.containers[0].env[] | select(.name == "TEMPORAL_ADDRESS").value) |= "temporal-frontend.temporal:7233"' $f
done
