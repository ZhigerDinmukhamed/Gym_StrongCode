FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=1
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata sqlite
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY .env .

RUN mkdir -p /app/data /app/logs

EXPOSE 8080

CMD ["./main"]
