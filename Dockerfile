# --- STAGE 1: Base (Copy the source code and install packages) ---
FROM golang:1.24-alpine as base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# --- STAGE 2: Development Runtime (Includes and run 'air'for hot-reloading) ---
FROM base as dev

RUN go install github.com/air-verse/air@v1.62.0

CMD ["air", "-c", ".air.toml"]

# --- STAGE 3: Build stage (Compiles the go binary) ---
FROM base AS build

RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot/main.go

# --- STAGE 4: Production Runtime (Minimal, only the compiled binary) ---
FROM scratch as prod

WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/bot .

CMD ["/app/bot"]
