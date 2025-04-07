FROM golang:1.23.8 AS oms-builder

WORKDIR /usr/src/oms

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY app ./app
COPY cmd ./cmd

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -v -o /usr/local/bin/oms ./cmd/oms

FROM busybox AS oms-worker

COPY --from=oms-builder /usr/local/bin/oms /usr/local/bin/oms

ENTRYPOINT ["/usr/local/bin/oms", "worker"]

FROM busybox AS oms-api

EXPOSE 8081
EXPOSE 8082
EXPOSE 8083
EXPOSE 8084

COPY --from=oms-builder /usr/local/bin/oms /usr/local/bin/oms

ENTRYPOINT ["/usr/local/bin/oms", "api"]

FROM busybox AS oms-codec-server

EXPOSE 8089

COPY --from=oms-builder /usr/local/bin/oms /usr/local/bin/oms

ENTRYPOINT ["/usr/local/bin/oms", "codec-server"]
