FROM golang:1.22.2 AS dev-server-builder

WORKDIR /usr/src/dev-server

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY app ./app
COPY cmd ./cmd
RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/dev-server ./cmd/dev-server

FROM scratch AS dev-server
EXPOSE 8081
EXPOSE 8082
EXPOSE 8083
EXPOSE 8084

COPY --from=dev-server-builder /usr/local/bin/dev-server /usr/local/bin/dev-server

ENTRYPOINT ["/usr/local/bin/dev-server"]

FROM node:20-slim AS web-builder

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY web /app
WORKDIR /app

FROM web-builder AS web-builder-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM web-builder AS web-builder-build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM web-builder AS web
EXPOSE 3000
COPY --from=web-builder-build /app/build /app
COPY --from=web-builder-deps /app/node_modules /app/node_modules
CMD ["node", "/app"]
