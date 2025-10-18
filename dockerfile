FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /evv-logger-backend ./cmd/api

FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser
COPY --from=builder /evv-logger-backend .

USER appuser
EXPOSE 8080
CMD ["./evv-logger-backend"]
