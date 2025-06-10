# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd/memory-service && go build -o /memory-service

# Final image
FROM alpine:latest
COPY --from=builder /memory-service /memory-service
ENTRYPOINT ["/memory-service"]