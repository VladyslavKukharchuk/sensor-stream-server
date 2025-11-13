# Build stage
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o sensor-stream-server main.go

# Run stage
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/sensor-stream-server ./sensor-stream-server
COPY internal/views ./views

EXPOSE 8080

CMD ["./sensor-stream-server"]
