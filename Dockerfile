# Frontend build
FROM oven/bun:1 AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile
COPY frontend/ .
RUN bun run build

# Backend build
FROM golang:1.22-alpine AS backend
WORKDIR /app
COPY backend/go.* ./
RUN go mod download
COPY backend/ .
COPY --from=frontend /app/frontend/dist ./dist
RUN go build -o server ./cmd/server

# Final
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=backend /app/server /server
EXPOSE 8080
CMD ["/server"]
