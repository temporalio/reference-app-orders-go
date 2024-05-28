FROM golang:1.22.2 AS dev-server-builder

WORKDIR /usr/src/dev-server

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY app ./app
COPY cmd ./cmd
RUN go build -v -o /usr/local/bin/dev-server ./cmd/dev-server

FROM scratch AS dev-server-builder

COPY --from=builder /usr/local/bin/dev-server /usr/local/bin/dev-server

CMD ["/usr/local/bin/dev-server"]

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
COPY --from=web-builder-deps /app/node_modules /app/node_modules
COPY --from=web-builder-build /app/build /app/build
EXPOSE 5173
CMD ["pnpm", "start"]