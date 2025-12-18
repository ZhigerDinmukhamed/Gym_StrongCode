FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Папки для данных и логов
RUN mkdir -p /app/data /app/logs

EXPOSE 8080

CMD ["./main"]