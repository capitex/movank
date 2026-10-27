# ── Stage 1: build
FROM golang:1.23-alpine AS builder
WORKDIR /src

# Cache de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Código fuente
COPY . .

# Binario estático (CGO off → funciona en cualquier linux)
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/backend .

# ── Stage 2: runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /out/backend /usr/local/bin/backend

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/backend"]