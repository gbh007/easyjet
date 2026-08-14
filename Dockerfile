FROM node:22-alpine AS vue-builder

WORKDIR /build

COPY frontend/package*.json ./
RUN --mount=type=cache,target="/build/node_modules" npm install

COPY frontend/ ./
RUN --mount=type=cache,target="/build/node_modules" npm run build

FROM golang:1.26.5 AS go-builder

WORKDIR /build

COPY go.mod go.sum ./
RUN --mount=type=cache,target="$(go env GOMODCACHE)" \
    --mount=type=cache,target="$(go env GOCACHE)" \
    go mod download 

COPY --exclude=frontend . .
RUN --mount=type=cache,target="$(go env GOMODCACHE)" \
    --mount=type=cache,target="$(go env GOCACHE)" \
    CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /build/main \
    cmd/server/main.go

FROM docker:29

COPY --from=go-builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=go-builder /build/main /app/main
COPY --from=vue-builder /build/dist /app/static

WORKDIR /app

ENTRYPOINT ["/app/main"]