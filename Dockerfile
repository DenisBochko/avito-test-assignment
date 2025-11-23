# Stage 1: Builder
FROM golang:1.25.2 AS builder

WORKDIR /app

COPY . .

RUN GOOS=linux go build -o /app/bin/avito-test-assignment ./cmd/avito_test_assignment

# Stage 2: Run
FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/bin/avito-test-assignment /app/bin/avito-test-assignment
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/bin/avito-test-assignment"]
