# Step 1: Build the Go API
FROM golang:1.23-alpine AS builder

# Install Reflex
RUN go install github.com/cespare/reflex@latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o absensi-api

# Step 2: Run the Go API in a lightweight container
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/absensi-api /app/absensi-api

EXPOSE 8080

CMD ["/app/absensi-api"]
